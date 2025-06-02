// Package database provides database operations for the application
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Zayan-Mohamed/do-task-backend/internal/models"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
)

// DB wraps the database connection
type DB struct {
	*sql.DB
}

// Connect initializes and returns a database connection
func Connect() (*DB, error) {
	// Get database URL from environment or use default
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/dotask?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &DB{db}, nil
}

// RunMigrations runs database migrations
func (db *DB) RunMigrations() error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// InitSchema initializes the database with sample data
func (db *DB) InitSchema() error {
	// Note: This function is for initial setup and doesn't use user-specific data
	// In production, categories and tasks should be created through authenticated requests

	// Skip initialization if this is a production environment
	// Categories and tasks will be created by authenticated users
	return nil
}

// CreateTask creates a new task
func (db *DB) CreateTask(input models.CreateTaskInput) (models.Task, error) {
	query := `
		INSERT INTO tasks (title, description, status, priority, due_date, category_id, user_id, tags)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, title, description, status, priority, due_date, created_at, updated_at, category_id, user_id, tags`

	dueDate, err := time.Parse(time.RFC3339, input.DueDate)
	if err != nil {
		return models.Task{}, err
	}

	var task models.Task
	err = db.QueryRow(query,
		input.Title,
		input.Description,
		input.Status,
		input.Priority,
		dueDate,
		input.CategoryID,
		input.UserID, // Add user_id parameter
		pq.Array(input.Tags),
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.CategoryID,
		&task.UserID,
		pq.Array(&task.Tags),
	)

	if err != nil {
		return models.Task{}, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

// GetTask retrieves a task by ID for a specific user
func (db *DB) GetTask(id string, userID string) (models.Task, error) {
	query := `
		SELECT id, title, description, status, priority, due_date, created_at, updated_at, category_id, user_id, tags
		FROM tasks WHERE id = $1 AND user_id = $2`

	var task models.Task
	err := db.QueryRow(query, id, userID).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.CategoryID,
		&task.UserID,
		pq.Array(&task.Tags),
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Task{}, errors.New("task not found")
		}
		return models.Task{}, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}

// GetAllTasks retrieves all tasks for a specific user
func (db *DB) GetAllTasks() ([]models.Task, error) {
	query := `
		SELECT id, title, description, status, priority, due_date, created_at, updated_at, category_id, user_id, tags
		FROM tasks ORDER BY created_at DESC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.CategoryID,
			&task.UserID,
			pq.Array(&task.Tags),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetAllTasksByUser retrieves all tasks for a specific user
func (db *DB) GetAllTasksByUser(userID string) ([]models.Task, error) {
	query := `
		SELECT id, title, description, status, priority, due_date, created_at, updated_at, category_id, user_id, tags
		FROM tasks WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks for user: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.CategoryID,
			&task.UserID,
			pq.Array(&task.Tags),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// UpdateTask updates an existing task for a specific user
func (db *DB) UpdateTask(id string, input models.UpdateTaskInput, userID string) (models.Task, error) {
	// Build dynamic query based on provided fields
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}
	argIndex := 1

	if input.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, *input.Title)
		argIndex++
	}
	if input.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *input.Description)
		argIndex++
	}
	if input.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *input.Status)
		argIndex++
	}
	if input.Priority != nil {
		setParts = append(setParts, fmt.Sprintf("priority = $%d", argIndex))
		args = append(args, *input.Priority)
		argIndex++
	}
	if input.DueDate != nil {
		dueDate, err := time.Parse(time.RFC3339, *input.DueDate)
		if err != nil {
			return models.Task{}, err
		}
		setParts = append(setParts, fmt.Sprintf("due_date = $%d", argIndex))
		args = append(args, dueDate)
		argIndex++
	}
	if input.CategoryID != nil {
		setParts = append(setParts, fmt.Sprintf("category_id = $%d", argIndex))
		args = append(args, *input.CategoryID)
		argIndex++
	}
	if input.Tags != nil {
		setParts = append(setParts, fmt.Sprintf("tags = $%d", argIndex))
		args = append(args, pq.Array(input.Tags))
		argIndex++
	}

	// Add ID and userID as the last arguments
	args = append(args, id, userID)

	// Join all setParts
	setClause := ""
	for i, part := range setParts {
		if i > 0 {
			setClause += ", "
		}
		setClause += part
	}

	query := fmt.Sprintf(`
		UPDATE tasks SET %s
		WHERE id = $%d AND user_id = $%d
		RETURNING id, title, description, status, priority, due_date, created_at, updated_at, category_id, user_id, tags`,
		setClause, argIndex, argIndex+1)

	var task models.Task
	err := db.QueryRow(query, args...).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.CategoryID,
		&task.UserID,
		pq.Array(&task.Tags),
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Task{}, errors.New("task not found")
		}
		return models.Task{}, fmt.Errorf("failed to update task: %w", err)
	}

	return task, nil
}

// UpdateTaskStatus updates the status of a task for a specific user
func (db *DB) UpdateTaskStatus(id string, status models.TaskStatus, userID string) (models.Task, error) {
	query := `
		UPDATE tasks SET status = $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3
		RETURNING id, title, description, status, priority, due_date, created_at, updated_at, category_id, user_id, tags`

	var task models.Task
	err := db.QueryRow(query, status, id, userID).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.CategoryID,
		&task.UserID,
		pq.Array(&task.Tags),
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Task{}, errors.New("task not found")
		}
		return models.Task{}, fmt.Errorf("failed to update task status: %w", err)
	}

	return task, nil
}

// DeleteTask deletes a task by ID for a specific user
func (db *DB) DeleteTask(id string, userID string) error {
	query := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`
	result, err := db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

// CreateCategory creates a new category for a specific user
func (db *DB) CreateCategory(name string, userID string) (models.Category, error) {
	query := `
		INSERT INTO categories (name, user_id)
		VALUES ($1, $2)
		RETURNING id, name, created_at, updated_at`

	var category models.Category
	err := db.QueryRow(query, name, userID).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" { // unique violation
			return models.Category{}, errors.New("category with this name already exists")
		}
		return models.Category{}, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

// GetCategory retrieves a category by ID for a specific user
func (db *DB) GetCategory(id string, userID string) (models.Category, error) {
	query := `SELECT id, name, created_at, updated_at FROM categories WHERE id = $1 AND user_id = $2`

	var category models.Category
	err := db.QueryRow(query, id, userID).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Category{}, errors.New("category not found")
		}
		return models.Category{}, fmt.Errorf("failed to get category: %w", err)
	}

	return category, nil
}

// GetAllCategories retrieves all categories for a specific user
func (db *DB) GetAllCategories(userID string) ([]models.Category, error) {
	query := `SELECT id, name, created_at, updated_at FROM categories WHERE user_id = $1 ORDER BY name`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// UpdateCategory updates an existing category for a specific user
func (db *DB) UpdateCategory(id string, name string, userID string) (models.Category, error) {
	query := `
		UPDATE categories SET name = $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3
		RETURNING id, name, created_at, updated_at`

	var category models.Category
	err := db.QueryRow(query, name, id, userID).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Category{}, errors.New("category not found")
		}
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return models.Category{}, errors.New("category with this name already exists")
		}
		return models.Category{}, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

// DeleteCategory deletes a category by ID for a specific user
func (db *DB) DeleteCategory(id string, userID string) error {
	// Check if any tasks are using this category
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM tasks WHERE category_id = $1 AND user_id = $2", id, userID).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check category usage: %w", err)
	}

	if count > 0 {
		return errors.New("cannot delete category with associated tasks")
	}

	query := `DELETE FROM categories WHERE id = $1 AND user_id = $2`
	result, err := db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("category not found")
	}

	return nil
}

// GetTasksInCategory retrieves all tasks in a category for a specific user
func (db *DB) GetTasksInCategory(categoryID string, userID string) ([]models.Task, error) {
	// First check if category exists for this user
	_, err := db.GetCategory(categoryID, userID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, title, description, status, priority, due_date, created_at, updated_at, category_id, user_id, tags
		FROM tasks WHERE category_id = $1 AND user_id = $2 ORDER BY created_at DESC`

	rows, err := db.Query(query, categoryID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks in category: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.CategoryID,
			&task.UserID,
			pq.Array(&task.Tags),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// User management functions

// CreateUser creates a new user
func (db *DB) CreateUser(name, email, hashedPassword string) (models.User, error) {
	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, created_at, updated_at`

	var user models.User
	err := db.QueryRow(query, name, email, hashedPassword).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" { // unique violation
			return models.User{}, errors.New("user with this email already exists")
		}
		return models.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (db *DB) GetUserByEmail(email string) (models.User, error) {
	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1`

	var user models.User
	err := db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (db *DB) GetUserByID(id string) (models.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1`

	var user models.User
	err := db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
