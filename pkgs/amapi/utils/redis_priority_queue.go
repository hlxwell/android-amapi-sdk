// Package utils provides utility functions for the amapi package.
package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisPriorityQueue provides a priority queue implementation using Redis sorted set.
//
// 使用 Redis sorted set 实现优先级队列。
// Score = 优先级（数字，越大优先级越高）
// Member = Task JSON 字符串
//
// # 使用方法
//
//	client := redis.NewClient(&redis.Options{
//	    Addr: "localhost:6379",
//	})
//
//	queue := NewRedisPriorityQueue(client, "amapi:")
//	defer queue.Close()
//
//	// 入队
//	task := &Task{...}
//	err := queue.Enqueue(ctx, task, 500)
//
//	// 阻塞出队
//	popped, err := queue.DequeueBlocking(ctx, 5*time.Second)
//
//	// 非阻塞出队
//	popped, err := queue.Dequeue(ctx)
type RedisPriorityQueue struct {
	client    *redis.Client
	keyPrefix string
	queueKey  string
}

// NewRedisPriorityQueue creates a new Redis priority queue.
func NewRedisPriorityQueue(client *redis.Client, keyPrefix string) *RedisPriorityQueue {
	if keyPrefix == "" {
		keyPrefix = "amapi:"
	}

	return &RedisPriorityQueue{
		client:    client,
		keyPrefix: keyPrefix,
		queueKey:  keyPrefix + "queue:priority",
	}
}

// Enqueue adds a task to the priority queue.
//
// priority 是任务的优先级（0-1000，越大优先级越高）。
// 如果 priority 小于 0，会被设置为 0；如果大于 1000，会被设置为 1000。
func (q *RedisPriorityQueue) Enqueue(ctx context.Context, task *Task, priority int) error {
	if task == nil {
		return fmt.Errorf("task cannot be nil")
	}

	// Validate and clamp priority
	if priority < 0 {
		priority = 0
	} else if priority > 1000 {
		priority = 1000
	}

	// Ensure task priority matches
	task.Priority = priority

	// Serialize task
	taskJSON, err := task.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize task: %w", err)
	}

	// Add to sorted set with score = priority
	// Use negative priority to reverse order (higher priority = higher score)
	// Actually, we want higher priority to be popped first, so we use higher score
	// Redis ZPOPMAX returns highest score first, which is what we want
	err = q.client.ZAdd(ctx, q.queueKey, redis.Z{
		Score:  float64(priority),
		Member: taskJSON,
	}).Err()

	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	return nil
}

// DequeueBlocking removes and returns the highest priority task, blocking if the queue is empty.
//
// timeout 是阻塞超时时间。如果队列为空，会阻塞直到有任务或超时。
// 如果超时，返回 redis.Nil 错误。
func (q *RedisPriorityQueue) DequeueBlocking(ctx context.Context, timeout time.Duration) (*Task, error) {
	// Use BZPOPMAX to block until an item is available
	// BZPOPMAX returns the item with the highest score (highest priority)
	result, err := q.client.BZPopMax(ctx, timeout, q.queueKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("queue is empty or timeout: %w", err)
		}
		return nil, fmt.Errorf("failed to dequeue task: %w", err)
	}

	// Parse task from JSON
	task, err := DeserializeTask(result.Member.(string))
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize task: %w", err)
	}

	return task, nil
}

// Dequeue removes and returns the highest priority task without blocking.
//
// 如果队列为空，返回 redis.Nil 错误。
func (q *RedisPriorityQueue) Dequeue(ctx context.Context) (*Task, error) {
	// Use ZPOPMAX to get and remove the item with the highest score
	result, err := q.client.ZPopMax(ctx, q.queueKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("queue is empty: %w", err)
		}
		return nil, fmt.Errorf("failed to dequeue task: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("queue is empty")
	}

	// Parse task from JSON
	taskJSON := result[0].Member.(string)
	task, err := DeserializeTask(taskJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize task: %w", err)
	}

	return task, nil
}

// Peek returns the highest priority task without removing it.
//
// 如果队列为空，返回 redis.Nil 错误。
func (q *RedisPriorityQueue) Peek(ctx context.Context) (*Task, error) {
	// Use ZRANGE to get the highest score item without removing it
	result, err := q.client.ZRangeWithScores(ctx, q.queueKey, -1, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to peek task: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("queue is empty")
	}

	// Parse task from JSON
	taskJSON := result[0].Member.(string)
	task, err := DeserializeTask(taskJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize task: %w", err)
	}

	return task, nil
}

// Size returns the number of tasks in the queue.
func (q *RedisPriorityQueue) Size(ctx context.Context) (int64, error) {
	count, err := q.client.ZCard(ctx, q.queueKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get queue size: %w", err)
	}
	return count, nil
}

// Clear removes all tasks from the queue.
func (q *RedisPriorityQueue) Clear(ctx context.Context) error {
	err := q.client.Del(ctx, q.queueKey).Err()
	if err != nil {
		return fmt.Errorf("failed to clear queue: %w", err)
	}
	return nil
}

// Close closes the priority queue (no-op for Redis implementation).
func (q *RedisPriorityQueue) Close() error {
	// Don't close the Redis client as it may be shared
	return nil
}


