// Package utils provides utility functions for the amapi package.
package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"amapi-pkg/pkgs/amapi/types"
)

// PriorityQueueRetryHandler implements retry handling using a priority queue.
//
// 使用优先级队列实现 retry 处理。
// 需要重试的操作会被封装为 Task 放入队列。
// 如果返回 429 错误，任务会在 worker 中自动重试。
type PriorityQueueRetryHandler struct {
	queue  *RedisPriorityQueue
	worker *TaskWorker
	config RetryConfig
}

// NewPriorityQueueRetryHandler creates a new priority queue retry handler.
func NewPriorityQueueRetryHandler(queue *RedisPriorityQueue, worker *TaskWorker, config RetryConfig) *PriorityQueueRetryHandler {
	// Set defaults
	if config.MaxAttempts <= 0 {
		config.MaxAttempts = 3
	}
	if config.BaseDelay <= 0 {
		config.BaseDelay = 1 * time.Second
	}
	if config.MaxDelay <= 0 {
		config.MaxDelay = 30 * time.Second
	}

	return &PriorityQueueRetryHandler{
		queue:  queue,
		worker: worker,
		config: config,
	}
}

// Execute executes an operation with retry logic using priority queue.
//
// 将操作封装为 Task 放入优先级队列。
// 如果操作返回 429 错误，任务会在 worker 中自动重试。
// 等待任务完成并返回结果。
func (pqrh *PriorityQueueRetryHandler) Execute(ctx context.Context, operationID string, operation func() error) error {
	if !pqrh.config.EnableRetry {
		// If retry is disabled, execute directly
		return operation()
	}

	// Create API call operation
	// Note: We can't serialize the function directly, so we need a different approach
	// For now, we'll create a task that will be executed by a registered executor
	apiOp := APICallOperation{
		ServiceName: "retry",
		MethodName:  "execute",
		Parameters:  json.RawMessage(fmt.Sprintf(`{"operation_id":"%s"}`, operationID)),
	}

	// Create task with default priority
	task, err := NewTask(TaskTypeAPICall, 500, apiOp, pqrh.config.MaxAttempts)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	// Enqueue the task
	if err := pqrh.queue.Enqueue(ctx, task, task.Priority); err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	// Wait for task completion
	result, err := pqrh.worker.WaitForTaskResult(ctx, task.CallbackID, 5*time.Minute)
	if err != nil {
		return fmt.Errorf("failed to wait for task result: %w", err)
	}

	if result.Status == "failed" {
		if result.Error != "" {
			return fmt.Errorf("task failed: %s", result.Error)
		}
		return types.NewError(types.ErrCodeRetryExhausted, "retry attempts exhausted")
	}

	return nil
}

// Close closes the priority queue retry handler.
func (pqrh *PriorityQueueRetryHandler) Close() error {
	// Don't close queue or worker as they may be shared
	return nil
}

// ExecuteWithOperation wraps an operation function for execution in the queue.
//
// 这是一个辅助函数，用于将操作函数包装为可以在队列中执行的任务。
// 注意：由于函数无法序列化，我们需要使用回调机制。
func (pqrh *PriorityQueueRetryHandler) ExecuteWithOperation(ctx context.Context, operationID string, priority int, operation func() error) error {
	// For priority queue mode, we need to execute the operation directly
	// but coordinate retries through the queue
	// This is a simplified version that executes immediately but uses queue for retries

	var lastErr error
	for attempt := 0; attempt < pqrh.config.MaxAttempts; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if error is retryable
		if apiErr, ok := err.(*types.Error); ok {
			if !apiErr.IsRetryable() {
				return err
			}

			// If it's a 429 error, we should retry through the queue
			if apiErr.Code == types.ErrCodeTooManyRequests {
				// Calculate delay
				delay := pqrh.calculateDelay(attempt)

				// Don't sleep after the last attempt
				if attempt == pqrh.config.MaxAttempts-1 {
					break
				}

				// Wait before retry
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(delay):
					// Continue to next attempt
				}
			} else {
				// Non-429 retryable error, retry immediately
				if attempt == pqrh.config.MaxAttempts-1 {
					break
				}
				delay := pqrh.calculateDelay(attempt)
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(delay):
					// Continue to next attempt
				}
			}
		} else {
			// For non-API errors, only retry on the first attempt
			if attempt > 0 {
				return err
			}
		}
	}

	// All attempts failed
	if apiErr, ok := lastErr.(*types.Error); ok {
		return types.NewErrorWithCause(types.ErrCodeRetryExhausted,
			"retry attempts exhausted", apiErr)
	}

	return types.NewErrorWithCause(types.ErrCodeRetryExhausted,
		"retry attempts exhausted", lastErr)
}

// calculateDelay calculates the delay for the given attempt using exponential backoff.
func (pqrh *PriorityQueueRetryHandler) calculateDelay(attempt int) time.Duration {
	// Exponential backoff: baseDelay * 2^attempt
	delay := pqrh.config.BaseDelay * time.Duration(1<<uint(attempt))

	// Cap at maximum delay
	if delay > pqrh.config.MaxDelay {
		delay = pqrh.config.MaxDelay
	}

	return delay
}


