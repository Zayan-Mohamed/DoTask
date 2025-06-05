// Package models defines the data structures used in the application
package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"` // This will be hashed
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// LoginInput represents the input for user login
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterInput represents the input for user registration
type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateProfileInput represents the input for updating user profile
type UpdateProfileInput struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

// ChangePasswordInput represents the input for changing user password
type ChangePasswordInput struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

// AuthResponse represents the response after successful authentication
type AuthResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

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
	UserID      string       `json:"userId"` // Associate task with user
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
	UserID      string       `json:"userId"` // Add UserID for authentication
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
