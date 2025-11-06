// Package utils provides utility functions for the amapi package.
package utils

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// TaskType represents the type of task operation.
type TaskType string

const (
	// TaskTypeAPICall represents an API call operation.
	TaskTypeAPICall TaskType = "api_call"
)

// Task represents a task to be executed in the priority queue.
//
// Task 包含了在优先级队列中执行的任务信息。
// 任务会被序列化为 JSON 存储在 Redis sorted set 中。
type Task struct {
	// ID is the unique identifier for the task.
	ID string `json:"id"`

	// Type is the type of operation.
	Type TaskType `json:"type"`

	// Priority is the task priority (0-1000, higher is better).
	// Default: 500
	Priority int `json:"priority"`

	// Operation contains the operation metadata (serialized as JSON).
	// For API calls, this contains the API endpoint, method, parameters, etc.
	Operation json.RawMessage `json:"operation"`

	// CreatedAt is when the task was created.
	CreatedAt time.Time `json:"created_at"`

	// MaxRetries is the maximum number of retry attempts.
	MaxRetries int `json:"max_retries"`

	// RetryCount is the current retry attempt number.
	RetryCount int `json:"retry_count"`

	// CallbackID is an identifier for retrieving the task result.
	// Results are stored in Redis with key: {prefix}task:result:{callbackID}
	CallbackID string `json:"callback_id"`
}

// NewTask creates a new task with the given parameters.
func NewTask(taskType TaskType, priority int, operation interface{}, maxRetries int) (*Task, error) {
	// Serialize operation to JSON
	operationJSON, err := json.Marshal(operation)
	if err != nil {
		return nil, err
	}

	// Validate priority range
	if priority < 0 {
		priority = 0
	} else if priority > 1000 {
		priority = 1000
	}

	// Generate unique IDs
	taskID := uuid.New().String()
	callbackID := uuid.New().String()

	return &Task{
		ID:         taskID,
		Type:       taskType,
		Priority:   priority,
		Operation:  operationJSON,
		CreatedAt: time.Now(),
		MaxRetries: maxRetries,
		RetryCount: 0,
		CallbackID: callbackID,
	}, nil
}

// Serialize converts the task to JSON string.
func (t *Task) Serialize() (string, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// DeserializeTask parses a JSON string into a Task.
func DeserializeTask(data string) (*Task, error) {
	var task Task
	if err := json.Unmarshal([]byte(data), &task); err != nil {
		return nil, err
	}
	return &task, nil
}

// TaskResult represents the result of a task execution.
type TaskResult struct {
	// TaskID is the ID of the task.
	TaskID string `json:"task_id"`

	// Status is the task status: "pending", "processing", "completed", "failed"
	Status string `json:"status"`

	// Result contains the result data (serialized as JSON).
	Result json.RawMessage `json:"result,omitempty"`

	// Error contains the error message if the task failed.
	Error string `json:"error,omitempty"`

	// CreatedAt is when the task was created.
	CreatedAt time.Time `json:"created_at"`

	// CompletedAt is when the task completed (or failed).
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// NewTaskResult creates a new task result.
func NewTaskResult(taskID string, status string) *TaskResult {
	return &TaskResult{
		TaskID:    taskID,
		Status:    status,
		CreatedAt: time.Now(),
	}
}

// Serialize converts the task result to JSON string.
func (tr *TaskResult) Serialize() (string, error) {
	data, err := json.Marshal(tr)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// DeserializeTaskResult parses a JSON string into a TaskResult.
func DeserializeTaskResult(data string) (*TaskResult, error) {
	var result TaskResult
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// APICallOperation represents an API call operation.
type APICallOperation struct {
	// ServiceName is the service name (e.g., "enterprises", "devices")
	ServiceName string `json:"service_name"`

	// MethodName is the method name (e.g., "List", "Get", "Create")
	MethodName string `json:"method_name"`

	// Parameters are the method parameters (serialized as JSON).
	Parameters json.RawMessage `json:"parameters"`

	// ResourceName is the resource name (e.g., enterprise name, device name)
	ResourceName string `json:"resource_name,omitempty"`
}


