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
	_ "github.com/lib/pq"
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
	// Create initial categories
	development, err := db.CreateCategory("Development")
	if err != nil {
		return err
	}

	design, err := db.CreateCategory("Design")
	if err != nil {
		return err
	}

	_, err = db.CreateCategory("Marketing")
	if err != nil {
		return err
	}

	_, err = db.CreateCategory("Personal")
	if err != nil {
		return err
	}

	// Create some initial tasks
	_, err = db.CreateTask(models.CreateTaskInput{
		Title:       "Complete project setup",
		Description: "Set up the SvelteKit and Golang GraphQL project structure",
		Status:      models.TaskStatusInProgress,
		Priority:    models.TaskPriorityHigh,
		DueDate:     time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		CategoryID:  development.ID,
		Tags:        []string{"setup", "project"},
	})

	if err != nil {
		return err
	}

	_, err = db.CreateTask(models.CreateTaskInput{
		Title:       "Design task interface",
		Description: "Create a clean and intuitive interface for managing tasks",
		Status:      models.TaskStatusTodo,
		Priority:    models.TaskPriorityMedium,
		DueDate:     time.Now().Add(48 * time.Hour).Format(time.RFC3339),
		CategoryID:  design.ID,
		Tags:        []string{"ui", "ux"},
	})

	if err != nil {
		return err
	}

	return nil
}

// CreateTask creates a new task
func (db *DB) CreateTask(input models.CreateTaskInput) (models.Task, error) {
	query := `
		INSERT INTO tasks (title, description, status, priority, due_date, category_id, tags)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, title, description, status, priority, due_date, created_at, updated_at, category_id, tags`

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
		pq.Array(&task.Tags),
	)

	if err != nil {
		return models.Task{}, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

// GetTask retrieves a task by ID
func (db *DB) GetTask(id string) (models.Task, error) {
	query := `
		SELECT id, title, description, status, priority, due_date, created_at, updated_at, category_id, tags
		FROM tasks WHERE id = $1`

	var task models.Task
	err := db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.CategoryID,
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

// GetAllTasks retrieves all tasks
func (db *DB) GetAllTasks() ([]models.Task, error) {
	query := `
		SELECT id, title, description, status, priority, due_date, created_at, updated_at, category_id, tags
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
			pq.Array(&task.Tags),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// UpdateTask updates an existing task
func (db *DB) UpdateTask(id string, input models.UpdateTaskInput) (models.Task, error) {
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

	// Add ID as the last argument
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE tasks SET %s
		WHERE id = $%d
		RETURNING id, title, description, status, priority, due_date, created_at, updated_at, category_id, tags`,
		fmt.Sprintf("%s", setParts[0]), // Handle the case where there might be no updates
		argIndex)

	// Join all setParts
	setClause := ""
	for i, part := range setParts {
		if i > 0 {
			setClause += ", "
		}
		setClause += part
	}

	query = fmt.Sprintf(`
		UPDATE tasks SET %s
		WHERE id = $%d
		RETURNING id, title, description, status, priority, due_date, created_at, updated_at, category_id, tags`,
		setClause, argIndex)

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

// UpdateTaskStatus updates the status of a task
func (db *DB) UpdateTaskStatus(id string, status models.TaskStatus) (models.Task, error) {
	query := `
		UPDATE tasks SET status = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, title, description, status, priority, due_date, created_at, updated_at, category_id, tags`

	var task models.Task
	err := db.QueryRow(query, status, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.CategoryID,
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

// DeleteTask deletes a task by ID
func (db *DB) DeleteTask(id string) error {
	query := `DELETE FROM tasks WHERE id = $1`
	result, err := db.Exec(query, id)
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

// CreateCategory creates a new category
func (db *DB) CreateCategory(name string) (models.Category, error) {
	query := `
		INSERT INTO categories (name)
		VALUES ($1)
		RETURNING id, name, created_at, updated_at`

	var category models.Category
	err := db.QueryRow(query, name).Scan(
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

// GetCategory retrieves a category by ID
func (db *DB) GetCategory(id string) (models.Category, error) {
	query := `SELECT id, name, created_at, updated_at FROM categories WHERE id = $1`

	var category models.Category
	err := db.QueryRow(query, id).Scan(
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

// GetAllCategories retrieves all categories
func (db *DB) GetAllCategories() ([]models.Category, error) {
	query := `SELECT id, name, created_at, updated_at FROM categories ORDER BY name`

	rows, err := db.Query(query)
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

// UpdateCategory updates an existing category
func (db *DB) UpdateCategory(id string, name string) (models.Category, error) {
	query := `
		UPDATE categories SET name = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, name, created_at, updated_at`

	var category models.Category
	err := db.QueryRow(query, name, id).Scan(
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

// DeleteCategory deletes a category by ID
func (db *DB) DeleteCategory(id string) error {
	// Check if any tasks are using this category
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM tasks WHERE category_id = $1", id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check category usage: %w", err)
	}

	if count > 0 {
		return errors.New("cannot delete category with associated tasks")
	}

	query := `DELETE FROM categories WHERE id = $1`
	result, err := db.Exec(query, id)
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

// GetTasksInCategory retrieves all tasks in a category
func (db *DB) GetTasksInCategory(categoryID string) ([]models.Task, error) {
	// First check if category exists
	_, err := db.GetCategory(categoryID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, title, description, status, priority, due_date, created_at, updated_at, category_id, tags
		FROM tasks WHERE category_id = $1 ORDER BY created_at DESC`

	rows, err := db.Query(query, categoryID)
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
			pq.Array(&task.Tags),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
