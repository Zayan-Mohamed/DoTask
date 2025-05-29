// Package models defines the data structures used in the application
package models

import (
	"time"
)

// TaskStatus represents the status of a task
type TaskStatus string

// Task priority levels
type TaskPriority string

// Task statuses
const (
	TaskStatusTodo       TaskStatus = "TODO"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusCompleted  TaskStatus = "COMPLETED"
)

// Task priorities
const (
	TaskPriorityLow    TaskPriority = "LOW"
	TaskPriorityMedium TaskPriority = "MEDIUM"
	TaskPriorityHigh   TaskPriority = "HIGH"
)

// Task represents a task in the system
type Task struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	DueDate     time.Time    `json:"dueDate"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	CategoryID  string       `json:"categoryId"`
	Tags        []string     `json:"tags"`
}

// Category represents a task category
type Category struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateTaskInput represents the input for creating a task
type CreateTaskInput struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	DueDate     string       `json:"dueDate"`
	CategoryID  string       `json:"categoryId"`
	Tags        []string     `json:"tags"`
}

// UpdateTaskInput represents the input for updating a task
type UpdateTaskInput struct {
	Title       *string       `json:"title"`
	Description *string       `json:"description"`
	Status      *TaskStatus   `json:"status"`
	Priority    *TaskPriority `json:"priority"`
	DueDate     *string       `json:"dueDate"`
	CategoryID  *string       `json:"categoryId"`
	Tags        []string      `json:"tags"`
}
