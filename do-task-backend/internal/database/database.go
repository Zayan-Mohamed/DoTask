// Package database provides database operations for the application
package database

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Zayan-Mohamed/do-task-backend/internal/models"
)

// DB is a mock in-memory database
type DB struct {
	tasks      map[string]models.Task
	categories map[string]models.Category
	mutex      sync.RWMutex
	nextID     int
}

// Connect initializes and returns a database connection
func Connect() (*DB, error) {
	db := &DB{
		tasks:      make(map[string]models.Task),
		categories: make(map[string]models.Category),
		nextID:     1,
	}
	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	// No-op for the mock implementation
	return nil
}

// InitSchema initializes the database schema with some sample data
func InitSchema(db *DB) error {
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
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Validate that the category exists
	if _, exists := db.categories[input.CategoryID]; !exists {
		return models.Task{}, errors.New("category not found")
	}

	now := time.Now()
	dueDate, err := time.Parse(time.RFC3339, input.DueDate)
	if err != nil {
		return models.Task{}, err
	}

	id := db.generateID()
	task := models.Task{
		ID:          id,
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Priority:    input.Priority,
		DueDate:     dueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
		CategoryID:  input.CategoryID,
		Tags:        input.Tags,
	}

	db.tasks[id] = task
	return task, nil
}

// GetTask retrieves a task by ID
func (db *DB) GetTask(id string) (models.Task, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	task, exists := db.tasks[id]
	if !exists {
		return models.Task{}, errors.New("task not found")
	}
	return task, nil
}

// GetAllTasks retrieves all tasks
func (db *DB) GetAllTasks() ([]models.Task, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	tasks := make([]models.Task, 0, len(db.tasks))
	for _, task := range db.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// UpdateTask updates an existing task
func (db *DB) UpdateTask(id string, input models.UpdateTaskInput) (models.Task, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	task, exists := db.tasks[id]
	if !exists {
		return models.Task{}, errors.New("task not found")
	}

	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.Status != nil {
		task.Status = *input.Status
	}
	if input.Priority != nil {
		task.Priority = *input.Priority
	}
	if input.DueDate != nil {
		dueDate, err := time.Parse(time.RFC3339, *input.DueDate)
		if err != nil {
			return models.Task{}, err
		}
		task.DueDate = dueDate
	}
	if input.CategoryID != nil {
		if _, exists := db.categories[*input.CategoryID]; !exists {
			return models.Task{}, errors.New("category not found")
		}
		task.CategoryID = *input.CategoryID
	}
	if input.Tags != nil {
		task.Tags = input.Tags
	}

	task.UpdatedAt = time.Now()
	db.tasks[id] = task
	return task, nil
}

// UpdateTaskStatus updates the status of a task
func (db *DB) UpdateTaskStatus(id string, status models.TaskStatus) (models.Task, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	task, exists := db.tasks[id]
	if !exists {
		return models.Task{}, errors.New("task not found")
	}

	task.Status = status
	task.UpdatedAt = time.Now()
	db.tasks[id] = task
	return task, nil
}

// DeleteTask deletes a task by ID
func (db *DB) DeleteTask(id string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, exists := db.tasks[id]; !exists {
		return errors.New("task not found")
	}

	delete(db.tasks, id)
	return nil
}

// CreateCategory creates a new category
func (db *DB) CreateCategory(name string) (models.Category, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Check for duplicate name
	for _, category := range db.categories {
		if category.Name == name {
			return models.Category{}, errors.New("category with this name already exists")
		}
	}

	id := db.generateID()
	category := models.Category{
		ID:   id,
		Name: name,
	}

	db.categories[id] = category
	return category, nil
}

// GetCategory retrieves a category by ID
func (db *DB) GetCategory(id string) (models.Category, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	category, exists := db.categories[id]
	if !exists {
		return models.Category{}, errors.New("category not found")
	}
	return category, nil
}

// GetAllCategories retrieves all categories
func (db *DB) GetAllCategories() ([]models.Category, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	categories := make([]models.Category, 0, len(db.categories))
	for _, category := range db.categories {
		categories = append(categories, category)
	}
	return categories, nil
}

// UpdateCategory updates an existing category
func (db *DB) UpdateCategory(id string, name string) (models.Category, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	category, exists := db.categories[id]
	if !exists {
		return models.Category{}, errors.New("category not found")
	}

	// Check for duplicate name
	for cid, cat := range db.categories {
		if cid != id && cat.Name == name {
			return models.Category{}, errors.New("category with this name already exists")
		}
	}

	category.Name = name
	db.categories[id] = category
	return category, nil
}

// DeleteCategory deletes a category by ID
func (db *DB) DeleteCategory(id string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, exists := db.categories[id]; !exists {
		return errors.New("category not found")
	}

	// Check if any tasks are using this category
	for _, task := range db.tasks {
		if task.CategoryID == id {
			return errors.New("cannot delete category with associated tasks")
		}
	}

	delete(db.categories, id)
	return nil
}

// GetTasksInCategory retrieves all tasks in a category
func (db *DB) GetTasksInCategory(categoryID string) ([]models.Task, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	if _, exists := db.categories[categoryID]; !exists {
		return nil, errors.New("category not found")
	}

	var tasks []models.Task
	for _, task := range db.tasks {
		if task.CategoryID == categoryID {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

// Helper to generate unique IDs
func (db *DB) generateID() string {
	id := db.nextID
	db.nextID++
	return fmt.Sprintf("%d", id)
}
