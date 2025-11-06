// Package utils provides utility functions for the amapi package.
package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"amapi-pkg/pkgs/amapi/types"
)

// TaskExecutor is a function that executes a task and returns the result or error.
type TaskExecutor func(ctx context.Context, operation json.RawMessage) (interface{}, error)

// TaskWorkerConfig contains configuration for the task worker.
type TaskWorkerConfig struct {
	// Concurrency is the number of concurrent workers.
	Concurrency int

	// PollInterval is the interval between queue polls when queue is empty.
	PollInterval time.Duration

	// KeyPrefix is the Redis key prefix.
	KeyPrefix string

	// RateLimit is the rate limit (requests per second).
	RateLimit int

	// Burst is the burst capacity.
	Burst int

	// MaxRetries is the maximum number of retries for failed tasks.
	MaxRetries int

	// BaseDelay is the base delay for retry backoff.
	BaseDelay time.Duration

	// MaxDelay is the maximum delay for retry backoff.
	MaxDelay time.Duration
}

// DefaultTaskWorkerConfig returns default configuration.
func DefaultTaskWorkerConfig() TaskWorkerConfig {
	return TaskWorkerConfig{
		Concurrency:  10,
		PollInterval: 100 * time.Millisecond,
		KeyPrefix:    "amapi:",
		RateLimit:    1000, // 1000 requests per second
		Burst:        100,
		MaxRetries:   3,
		BaseDelay:    1 * time.Second,
		MaxDelay:     30 * time.Second,
	}
}

// TaskWorker consumes tasks from a priority queue and executes them.
//
// TaskWorker 从优先级队列中消费任务并执行。
// 执行前会检查 rate limit，执行后会处理 429 错误并重试。
type TaskWorker struct {
	client      *redis.Client
	queue       *RedisPriorityQueue
	rateLimiter *RedisRateLimiter
	config      TaskWorkerConfig
	executors   map[TaskType]TaskExecutor
	mu          sync.RWMutex

	// Control
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewTaskWorker creates a new task worker.
func NewTaskWorker(client *redis.Client, config TaskWorkerConfig) *TaskWorker {
	if config.Concurrency <= 0 {
		config.Concurrency = 10
	}
	if config.PollInterval <= 0 {
		config.PollInterval = 100 * time.Millisecond
	}

	// Create rate limiter with 1 second window for per-second rate limiting
	// Convert rate limit from per-second to per-minute for RedisRateLimiter
	rateLimiter := NewRedisRateLimiterWithWindow(client, config.KeyPrefix, config.RateLimit*60, config.Burst, 1*time.Second)

	// Create priority queue
	queue := NewRedisPriorityQueue(client, config.KeyPrefix)

	return &TaskWorker{
		client:      client,
		queue:       queue,
		rateLimiter: rateLimiter,
		config:      config,
		executors:   make(map[TaskType]TaskExecutor),
	}
}

// RegisterExecutor registers a task executor for a given task type.
func (tw *TaskWorker) RegisterExecutor(taskType TaskType, executor TaskExecutor) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	tw.executors[taskType] = executor
}

// Start starts the task worker.
//
// Worker 会启动多个 goroutine 并发消费队列。
// 每个 worker 会：
// 1. 从队列中取出任务（按优先级）
// 2. 检查并等待 rate limit
// 3. 执行任务
// 4. 处理 429 错误并重试
// 5. 存储结果
func (tw *TaskWorker) Start(ctx context.Context) error {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	if tw.ctx != nil {
		return fmt.Errorf("worker is already running")
	}

	tw.ctx, tw.cancel = context.WithCancel(ctx)

	// Start worker goroutines
	for i := 0; i < tw.config.Concurrency; i++ {
		tw.wg.Add(1)
		go tw.worker(i)
	}

	return nil
}

// Stop stops the task worker gracefully.
func (tw *TaskWorker) Stop() {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	if tw.cancel != nil {
		tw.cancel()
		tw.cancel = nil
	}

	tw.wg.Wait()
	tw.ctx = nil
}

// worker is the main worker loop.
func (tw *TaskWorker) worker(id int) {
	defer tw.wg.Done()

	ticker := time.NewTicker(tw.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-tw.ctx.Done():
			return
		case <-ticker.C:
			// Try to dequeue a task (non-blocking)
			task, err := tw.queue.Dequeue(tw.ctx)
			if err != nil {
				// Error dequeuing, continue to next poll
				continue
			}

			if task == nil {
				// Queue is empty, continue polling
				continue
			}

			// Process the task
			tw.processTask(tw.ctx, task)
		}
	}
}

// processTask processes a single task.
func (tw *TaskWorker) processTask(ctx context.Context, task *Task) {
	// Update task status to processing immediately
	// Use background context to avoid cancellation during status update
	statusCtx := context.Background()
	tw.updateTaskStatus(statusCtx, task.CallbackID, "processing", nil, nil)

	// Wait for rate limit before executing
	if err := tw.rateLimiter.Wait(ctx); err != nil {
		statusCtx := context.Background()
		tw.updateTaskStatus(statusCtx, task.CallbackID, "failed", nil, fmt.Errorf("rate limit error: %w", err))
		return
	}

	// Execute the task
	result, err := tw.executeTask(ctx, task)

	// Handle retry for 429 errors
	if err != nil && tw.is429Error(err) {
		if task.RetryCount < task.MaxRetries {
			// Calculate retry delay
			delay := tw.calculateRetryDelay(task.RetryCount)

			// Reduce priority for retry
			newPriority := task.Priority - 50
			if newPriority < 0 {
				newPriority = 0
			}

			// Increment retry count
			task.RetryCount++

			// Wait for delay
			select {
			case <-ctx.Done():
				statusCtx := context.Background()
				tw.updateTaskStatus(statusCtx, task.CallbackID, "failed", nil, ctx.Err())
				return
			case <-time.After(delay):
				// Re-enqueue task with lower priority
				if err := tw.queue.Enqueue(ctx, task, newPriority); err != nil {
					statusCtx := context.Background()
					tw.updateTaskStatus(statusCtx, task.CallbackID, "failed", nil, fmt.Errorf("failed to re-enqueue: %w", err))
					return
				}
				// Update status to pending for retry
				statusCtx := context.Background()
				tw.updateTaskStatus(statusCtx, task.CallbackID, "pending", nil, nil)
				return
			}
		}
	}

	// Task completed (success or non-retryable error)
	finalStatusCtx := context.Background()
	if err != nil {
		tw.updateTaskStatus(finalStatusCtx, task.CallbackID, "failed", nil, err)
	} else {
		tw.updateTaskStatus(finalStatusCtx, task.CallbackID, "completed", result, nil)
	}
}

// executeTask executes a task using the registered executor.
func (tw *TaskWorker) executeTask(ctx context.Context, task *Task) (interface{}, error) {
	tw.mu.RLock()
	executor, exists := tw.executors[task.Type]
	tw.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no executor registered for task type: %s", task.Type)
	}

	return executor(ctx, task.Operation)
}

// is429Error checks if an error is a 429 Too Many Requests error.
func (tw *TaskWorker) is429Error(err error) bool {
	if err == nil {
		return false
	}

	// Check if it's a types.Error with 429 code
	if apiErr, ok := err.(*types.Error); ok {
		return apiErr.Code == types.ErrCodeTooManyRequests
	}

	// Check error message for 429
	errStr := err.Error()
	return contains429(errStr)
}

// contains429 checks if error string contains 429.
func contains429(s string) bool {
	// Simple check for 429 in error message
	if len(s) < 3 {
		return false
	}
	// Check for "429" or "Too Many Requests" in the string
	for i := 0; i <= len(s)-3; i++ {
		if s[i:i+3] == "429" {
			return true
		}
	}
	// Check for "Too Many Requests"
	for i := 0; i <= len(s)-17; i++ {
		if s[i:i+17] == "Too Many Requests" {
			return true
		}
	}
	return false
}

// calculateRetryDelay calculates the retry delay using exponential backoff.
func (tw *TaskWorker) calculateRetryDelay(attempt int) time.Duration {
	// Exponential backoff: baseDelay * 2^attempt
	delay := tw.config.BaseDelay * time.Duration(1<<uint(attempt))

	// Cap at maximum delay
	if delay > tw.config.MaxDelay {
		delay = tw.config.MaxDelay
	}

	return delay
}

// updateTaskStatus updates the task status in Redis.
func (tw *TaskWorker) updateTaskStatus(ctx context.Context, callbackID string, status string, result interface{}, err error) {
	resultKey := tw.config.KeyPrefix + "task:result:" + callbackID

	now := time.Now()
	resultData := map[string]interface{}{
		"status":     status,
		"created_at": now.Format(time.RFC3339),
		"updated_at": now.Format(time.RFC3339),
	}

	if result != nil {
		resultJSON, marshalErr := json.Marshal(result)
		if marshalErr == nil {
			resultData["result"] = string(resultJSON)
		}
	}

	if err != nil {
		resultData["error"] = err.Error()
	}

	if status == "completed" || status == "failed" {
		resultData["completed_at"] = now.Format(time.RFC3339)
	}

	// Store as hash in Redis
	pipe := tw.client.Pipeline()
	for k, v := range resultData {
		pipe.HSet(ctx, resultKey, k, v)
	}
	pipe.Expire(ctx, resultKey, 1*time.Hour) // Expire after 1 hour
	if _, err := pipe.Exec(ctx); err != nil {
		// Log error but don't fail the task
		// In production, you might want to use a logger here
		_ = err
	}
}

// GetTaskResult retrieves the result of a task.
func (tw *TaskWorker) GetTaskResult(ctx context.Context, callbackID string) (*TaskResult, error) {
	resultKey := tw.config.KeyPrefix + "task:result:" + callbackID

	// Get all fields from hash
	result, err := tw.client.HGetAll(ctx, resultKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get task result: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("task result not found")
	}

	// Parse result
	taskResult := &TaskResult{
		TaskID: callbackID,
		Status: result["status"],
	}

	if result["result"] != "" {
		taskResult.Result = json.RawMessage(result["result"])
	}

	if result["error"] != "" {
		taskResult.Error = result["error"]
	}

	// Parse timestamps
	if created, ok := result["created_at"]; ok && created != "" {
		if t, err := time.Parse(time.RFC3339, created); err == nil {
			taskResult.CreatedAt = t
		}
	}

	if completed, ok := result["completed_at"]; ok && completed != "" {
		if t, err := time.Parse(time.RFC3339, completed); err == nil {
			taskResult.CompletedAt = &t
		}
	}

	return taskResult, nil
}

// WaitForTaskResult waits for a task to complete and returns the result.
func (tw *TaskWorker) WaitForTaskResult(ctx context.Context, callbackID string, timeout time.Duration) (*TaskResult, error) {
	deadline := time.Now().Add(timeout)
	pollInterval := 100 * time.Millisecond

	for time.Now().Before(deadline) {
		result, err := tw.GetTaskResult(ctx, callbackID)
		if err == nil {
			if result.Status == "completed" || result.Status == "failed" {
				return result, nil
			}
		}

		// Wait before next poll
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(pollInterval):
			continue
		}
	}

	return nil, fmt.Errorf("timeout waiting for task result")
}

// Close closes the task worker and releases resources.
func (tw *TaskWorker) Close() error {
	tw.Stop()
	if tw.rateLimiter != nil {
		return tw.rateLimiter.Close()
	}
	return nil
}

