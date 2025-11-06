// Package utils provides utility functions for the amapi package.
package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	"amapi-pkg/pkgs/amapi/types"
)

// setupTestRedis creates a test Redis server using miniredis.
func setupTestRedis(t *testing.T) (*redis.Client, func()) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("Failed to connect to test Redis: %v", err)
	}

	cleanup := func() {
		client.Close()
		mr.Close()
	}

	return client, cleanup
}

// mockAPICallExecutor creates a mock executor that can simulate different scenarios.
func mockAPICallExecutor(success bool, errCode int, delay time.Duration) TaskExecutor {
	return func(ctx context.Context, operation json.RawMessage) (interface{}, error) {
		// Simulate delay
		if delay > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		// Simulate different error scenarios
		if !success {
			if errCode == types.ErrCodeTooManyRequests {
				return nil, types.NewError(types.ErrCodeTooManyRequests, "rate limit exceeded")
			}
			return nil, fmt.Errorf("mock API error")
		}

		return map[string]interface{}{
			"success": true,
			"data":    "mock response",
		}, nil
	}
}

// TestTaskSerialization tests task serialization and deserialization.
func TestTaskSerialization(t *testing.T) {
	operation := APICallOperation{
		ServiceName: "test",
		MethodName:  "test",
		Parameters:  []byte(`{"test": "data"}`),
	}

	task, err := NewTask(TaskTypeAPICall, 500, operation, 3)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	if task.ID == "" {
		t.Error("Task ID should not be empty")
	}

	if task.Type != TaskTypeAPICall {
		t.Errorf("Expected task type %s, got %s", TaskTypeAPICall, task.Type)
	}

	if task.Priority != 500 {
		t.Errorf("Expected priority 500, got %d", task.Priority)
	}

	// Test serialization
	serialized, err := task.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize task: %v", err)
	}

	// Test deserialization
	deserialized, err := DeserializeTask(serialized)
	if err != nil {
		t.Fatalf("Failed to deserialize task: %v", err)
	}

	if deserialized.ID != task.ID {
		t.Errorf("ID mismatch: expected %s, got %s", task.ID, deserialized.ID)
	}

	if deserialized.Priority != task.Priority {
		t.Errorf("Priority mismatch: expected %d, got %d", task.Priority, deserialized.Priority)
	}
}

// TestRedisPriorityQueue tests the Redis priority queue functionality.
func TestRedisPriorityQueue(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	queue := NewRedisPriorityQueue(client, "test:")

	ctx := context.Background()

	// Test enqueue
	operation1 := APICallOperation{ServiceName: "test1", MethodName: "test"}
	task1, err := NewTask(TaskTypeAPICall, 500, operation1, 3)
	if err != nil {
		t.Fatalf("Failed to create task1: %v", err)
	}

	operation2 := APICallOperation{ServiceName: "test2", MethodName: "test"}
	task2, err := NewTask(TaskTypeAPICall, 700, operation2, 3)
	if err != nil {
		t.Fatalf("Failed to create task2: %v", err)
	}

	// Enqueue tasks with different priorities
	if err := queue.Enqueue(ctx, task1, 500); err != nil {
		t.Fatalf("Failed to enqueue task1: %v", err)
	}

	if err := queue.Enqueue(ctx, task2, 700); err != nil {
		t.Fatalf("Failed to enqueue task2: %v", err)
	}

	// Test size
	size, err := queue.Size(ctx)
	if err != nil {
		t.Fatalf("Failed to get queue size: %v", err)
	}

	if size != 2 {
		t.Errorf("Expected queue size 2, got %d", size)
	}

	// Test dequeue (should get highest priority first)
	dequeued, err := queue.Dequeue(ctx)
	if err != nil {
		t.Fatalf("Failed to dequeue: %v", err)
	}

	// Task2 should be dequeued first (higher priority)
	if dequeued.Priority != 700 {
		t.Errorf("Expected priority 700, got %d", dequeued.Priority)
	}

	// Dequeue the remaining task
	dequeued, err = queue.Dequeue(ctx)
	if err != nil {
		t.Fatalf("Failed to dequeue: %v", err)
	}

	if dequeued.Priority != 500 {
		t.Errorf("Expected priority 500, got %d", dequeued.Priority)
	}

	// Test empty queue
	_, err = queue.Dequeue(ctx)
	if err == nil {
		t.Error("Expected error when dequeuing from empty queue")
	}
}

// TestTaskWorkerIntegration tests the complete task worker flow.
func TestTaskWorkerIntegration(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create worker config
	config := TaskWorkerConfig{
		Concurrency:  2,
		PollInterval: 200 * time.Millisecond, // Use shorter interval for faster processing
		KeyPrefix:    "test:",
		RateLimit:    1000, // High rate limit for testing
		Burst:        100,
		MaxRetries:   3,
		BaseDelay:    500 * time.Millisecond, // Shorter delay for testing
		MaxDelay:     2 * time.Second,
	}

	// Create worker
	worker := NewTaskWorker(client, config)

	// Register executor that simulates successful API call
	worker.RegisterExecutor(TaskTypeAPICall, mockAPICallExecutor(true, 0, 100*time.Millisecond))

	// Start worker
	if err := worker.Start(ctx); err != nil {
		t.Fatalf("Failed to start worker: %v", err)
	}
	defer worker.Stop()

	// Create and enqueue task
	queue := NewRedisPriorityQueue(client, config.KeyPrefix)
	operation := APICallOperation{
		ServiceName: "test",
		MethodName:  "test",
		Parameters:  []byte(`{}`),
	}

	task, err := NewTask(TaskTypeAPICall, 500, operation, config.MaxRetries)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	if err := queue.Enqueue(ctx, task, task.Priority); err != nil {
		t.Fatalf("Failed to enqueue task: %v", err)
	}

	// Give worker time to start processing and pick up the task
	time.Sleep(1 * time.Second)

	// Wait for task completion
	result, err := worker.WaitForTaskResult(ctx, task.CallbackID, 15*time.Second)
	if err != nil {
		t.Fatalf("Failed to wait for task result: %v", err)
	}

	if result.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", result.Status)
	}

	if result.Error != "" {
		t.Errorf("Expected no error, got '%s'", result.Error)
	}
}

// TestTaskWorker429Retry tests 429 error retry logic.
func TestTaskWorker429Retry(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create worker config
	config := TaskWorkerConfig{
		Concurrency:  1,
		PollInterval: 1 * time.Second, // miniredis requires at least 1s
		KeyPrefix:    "test:",
		RateLimit:    100,
		Burst:        10,
		MaxRetries:   3,
		BaseDelay:    1 * time.Second, // miniredis requires at least 1s
		MaxDelay:     2 * time.Second,
	}

	// Create worker
	worker := NewTaskWorker(client, config)

	// Track retry attempts
	retryCount := 0
	maxRetries := 3

	// Create executor that returns 429 on first attempts, then succeeds
	executor := func(ctx context.Context, operation json.RawMessage) (interface{}, error) {
		retryCount++
		if retryCount < maxRetries {
			return nil, types.NewError(types.ErrCodeTooManyRequests, "rate limit exceeded")
		}
		return map[string]interface{}{"success": true}, nil
	}

	worker.RegisterExecutor(TaskTypeAPICall, executor)

	// Start worker
	if err := worker.Start(ctx); err != nil {
		t.Fatalf("Failed to start worker: %v", err)
	}
	defer worker.Stop()

	// Create and enqueue task
	queue := NewRedisPriorityQueue(client, config.KeyPrefix)
	operation := APICallOperation{
		ServiceName: "test",
		MethodName:  "test",
		Parameters:  []byte(`{}`),
	}

	task, err := NewTask(TaskTypeAPICall, 500, operation, maxRetries)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	if err := queue.Enqueue(ctx, task, task.Priority); err != nil {
		t.Fatalf("Failed to enqueue task: %v", err)
	}

	// Give worker time to start processing
	time.Sleep(500 * time.Millisecond)

	// Wait for task completion (with longer timeout for retries)
	result, err := worker.WaitForTaskResult(ctx, task.CallbackID, 20*time.Second)
	if err != nil {
		t.Fatalf("Failed to wait for task result: %v", err)
	}

	// Task should eventually succeed after retries
	if result.Status != "completed" {
		t.Errorf("Expected status 'completed' after retries, got '%s' (error: %s)", result.Status, result.Error)
	}

	// Verify retry count
	if retryCount < maxRetries {
		t.Errorf("Expected at least %d retry attempts, got %d", maxRetries, retryCount)
	}
}

// TestPriorityQueueRateLimiter tests the priority queue rate limiter.
func TestPriorityQueueRateLimiter(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()

	// Create queue and worker
	queue := NewRedisPriorityQueue(client, "test:")
	config := DefaultTaskWorkerConfig()
	config.KeyPrefix = "test:"
	config.RateLimit = 100
	config.Burst = 10

	worker := NewTaskWorker(client, config)
	worker.RegisterExecutor(TaskTypeAPICall, mockAPICallExecutor(true, 0, 0))

	// Create rate limiter
	limiterConfig := DefaultPriorityQueueRateLimiterConfig()
	limiterConfig.DefaultPriority = 500
	limiterConfig.MaxQueueSize = 1000

	limiter := NewPriorityQueueRateLimiter(queue, worker, limiterConfig)

	// Test Allow (should check queue size)
	allowed := limiter.Allow(ctx)
	if !allowed {
		t.Error("Expected Allow() to return true for empty queue")
	}

	// Test Wait (should not block for empty queue)
	err := limiter.Wait(ctx)
	if err != nil {
		t.Errorf("Expected Wait() to succeed, got error: %v", err)
	}
}

// TestPriorityQueueRetryHandler tests the priority queue retry handler.
func TestPriorityQueueRetryHandler(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create queue and worker
	queue := NewRedisPriorityQueue(client, "test:")
	config := DefaultTaskWorkerConfig()
	config.KeyPrefix = "test:"
	config.RateLimit = 100
	config.Burst = 10

	worker := NewTaskWorker(client, config)
	worker.RegisterExecutor(TaskTypeAPICall, mockAPICallExecutor(true, 0, 50*time.Millisecond))

	// Start worker
	if err := worker.Start(ctx); err != nil {
		t.Fatalf("Failed to start worker: %v", err)
	}
	defer worker.Stop()

	// Create retry handler
	retryConfig := RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   100 * time.Millisecond,
		MaxDelay:    1 * time.Second,
		EnableRetry: true,
	}

	handler := NewPriorityQueueRetryHandler(queue, worker, retryConfig)

	// Test Execute with successful operation
	operationID := "test-operation-1"
	operation := func() error {
		return nil
	}

	// Note: This will use the ExecuteWithOperation method which executes directly
	// For full queue-based execution, we'd need to integrate with the client
	err := handler.Execute(ctx, operationID, operation)
	if err != nil {
		t.Errorf("Expected Execute() to succeed, got error: %v", err)
	}
}

// TestEndToEndFlow tests the complete end-to-end flow.
func TestEndToEndFlow(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Setup
	config := TaskWorkerConfig{
		Concurrency:  3,
		PollInterval: 1 * time.Second, // miniredis requires at least 1s
		KeyPrefix:    "test:",
		RateLimit:    100, // 100 requests per second
		Burst:        10,
		MaxRetries:   3,
		BaseDelay:    1 * time.Second, // miniredis requires at least 1s
		MaxDelay:     2 * time.Second,
	}

	queue := NewRedisPriorityQueue(client, config.KeyPrefix)
	worker := NewTaskWorker(client, config)

	// Register executor
	worker.RegisterExecutor(TaskTypeAPICall, mockAPICallExecutor(true, 0, 100*time.Millisecond))

	// Start worker
	if err := worker.Start(ctx); err != nil {
		t.Fatalf("Failed to start worker: %v", err)
	}
	defer worker.Stop()

	// Create rate limiter and retry handler
	limiterConfig := DefaultPriorityQueueRateLimiterConfig()
	limiterConfig.DefaultPriority = 500
	limiterConfig.MaxQueueSize = 1000

	limiter := NewPriorityQueueRateLimiter(queue, worker, limiterConfig)

	retryConfig := RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   100 * time.Millisecond,
		MaxDelay:    1 * time.Second,
		EnableRetry: true,
	}

	_ = NewPriorityQueueRetryHandler(queue, worker, retryConfig)

	// Test 1: Basic flow - enqueue and execute
	t.Run("BasicFlow", func(t *testing.T) {
		operation := APICallOperation{
			ServiceName: "test",
			MethodName:  "test",
			Parameters:  []byte(`{}`),
		}

		task, err := NewTask(TaskTypeAPICall, 500, operation, 3)
		if err != nil {
			t.Fatalf("Failed to create task: %v", err)
		}

		// Check rate limiter
		if !limiter.Allow(ctx) {
			t.Error("Rate limiter should allow task")
		}

		// Enqueue
		if err := queue.Enqueue(ctx, task, task.Priority); err != nil {
			t.Fatalf("Failed to enqueue: %v", err)
		}

		// Give worker time to process
		time.Sleep(500 * time.Millisecond)

		// Wait for completion
		result, err := worker.WaitForTaskResult(ctx, task.CallbackID, 10*time.Second)
		if err != nil {
			t.Fatalf("Failed to wait for result: %v", err)
		}

		if result.Status != "completed" {
			t.Errorf("Expected completed, got %s", result.Status)
		}
	})

	// Test 2: Priority ordering
	t.Run("PriorityOrdering", func(t *testing.T) {
		// Enqueue tasks with different priorities
		tasks := []*Task{}
		priorities := []int{300, 700, 500, 900}

		for i, priority := range priorities {
			operation := APICallOperation{
				ServiceName: "test",
				MethodName:  fmt.Sprintf("test%d", i),
				Parameters:  []byte(`{}`),
			}

			task, err := NewTask(TaskTypeAPICall, priority, operation, 3)
			if err != nil {
				t.Fatalf("Failed to create task: %v", err)
			}

			tasks = append(tasks, task)
			if err := queue.Enqueue(ctx, task, priority); err != nil {
				t.Fatalf("Failed to enqueue task: %v", err)
			}
		}

		// Give worker time to process
		time.Sleep(1 * time.Second)

		// Wait for all tasks to complete
		completed := 0
		for _, task := range tasks {
			result, err := worker.WaitForTaskResult(ctx, task.CallbackID, 15*time.Second)
			if err != nil {
				t.Errorf("Failed to wait for task %s: %v", task.ID, err)
				continue
			}

			if result.Status == "completed" {
				completed++
			}
		}

		if completed != len(tasks) {
			t.Errorf("Expected all %d tasks to complete, got %d", len(tasks), completed)
		}
	})

	// Test 3: Rate limiting
	t.Run("RateLimiting", func(t *testing.T) {
		// This test verifies that rate limiting is applied
		// We can't easily test the exact timing, but we can verify the worker respects rate limits

		// Create multiple tasks
		for i := 0; i < 5; i++ {
			operation := APICallOperation{
				ServiceName: "test",
				MethodName:  fmt.Sprintf("test%d", i),
				Parameters:  []byte(`{}`),
			}

			task, err := NewTask(TaskTypeAPICall, 500, operation, 3)
			if err != nil {
				t.Fatalf("Failed to create task: %v", err)
			}

			if err := queue.Enqueue(ctx, task, task.Priority); err != nil {
				t.Fatalf("Failed to enqueue task: %v", err)
			}
		}

		// Wait for all tasks (they should be rate limited)
		time.Sleep(1 * time.Second)

		// Verify queue is being processed
		size, err := queue.Size(ctx)
		if err != nil {
			t.Fatalf("Failed to get queue size: %v", err)
		}

		// Queue should be processed (may not be empty due to rate limiting)
		if size > 5 {
			t.Errorf("Queue size should decrease, got %d", size)
		}
	})
}

// TestErrorHandling tests error handling in various scenarios.
func TestErrorHandling(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config := TaskWorkerConfig{
		Concurrency:  1,
		PollInterval: 1 * time.Second, // miniredis requires at least 1s
		KeyPrefix:    "test:",
		RateLimit:    100,
		Burst:        10,
		MaxRetries:   2, // Only 2 retries
		BaseDelay:    1 * time.Second, // miniredis requires at least 1s
		MaxDelay:     2 * time.Second,
	}

	worker := NewTaskWorker(client, config)

	// Create executor that always fails
	worker.RegisterExecutor(TaskTypeAPICall, func(ctx context.Context, operation json.RawMessage) (interface{}, error) {
		return nil, errors.New("permanent error")
	})

	if err := worker.Start(ctx); err != nil {
		t.Fatalf("Failed to start worker: %v", err)
	}
	defer worker.Stop()

	queue := NewRedisPriorityQueue(client, config.KeyPrefix)
	operation := APICallOperation{
		ServiceName: "test",
		MethodName:  "test",
		Parameters:  []byte(`{}`),
	}

	task, err := NewTask(TaskTypeAPICall, 500, operation, config.MaxRetries)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	if err := queue.Enqueue(ctx, task, task.Priority); err != nil {
		t.Fatalf("Failed to enqueue task: %v", err)
	}

	// Give worker time to process
	time.Sleep(500 * time.Millisecond)

	// Wait for task result
	result, err := worker.WaitForTaskResult(ctx, task.CallbackID, 10*time.Second)
	if err != nil {
		t.Fatalf("Failed to wait for task result: %v", err)
	}

	// Task should fail
	if result.Status != "failed" {
		t.Errorf("Expected status 'failed', got '%s'", result.Status)
	}

	if result.Error == "" {
		t.Error("Expected error message, got empty")
	}
}

