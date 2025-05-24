// Package resolvers implements the GraphQL resolvers
package resolvers

import (
	"context"

	"github.com/Zayan-Mohamed/do-task-backend/internal/database"
	"github.com/Zayan-Mohamed/do-task-backend/internal/graph/generated"
	"github.com/Zayan-Mohamed/do-task-backend/internal/models"
)

// Resolver is the root resolver for the GraphQL schema
type Resolver struct {
	DB *database.DB
}

// Query returns the query resolver
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

// Mutation returns the mutation resolver
func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

// Task returns the task resolver
func (r *Resolver) Task() generated.TaskResolver {
	return &taskResolver{r}
}

// Category returns the category resolver
func (r *Resolver) Category() generated.CategoryResolver {
	return &categoryResolver{r}
}

// Tasks returns all tasks
func (r *queryResolver) Tasks(ctx context.Context) ([]*models.Task, error) {
	tasks, err := r.DB.GetAllTasks()
	if err != nil {
		return nil, err
	}

	// Convert to pointer slice
	result := make([]*models.Task, len(tasks))
	for i := range tasks {
		task := tasks[i]
		result[i] = &task
	}
	return result, nil
}

// Task returns a task by ID
func (r *queryResolver) Task(ctx context.Context, id string) (*models.Task, error) {
	task, err := r.DB.GetTask(id)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// Categories returns all categories
func (r *queryResolver) Categories(ctx context.Context) ([]*models.Category, error) {
	categories, err := r.DB.GetAllCategories()
	if err != nil {
		return nil, err
	}

	// Convert to pointer slice
	result := make([]*models.Category, len(categories))
	for i := range categories {
		category := categories[i]
		result[i] = &category
	}
	return result, nil
}

// Category returns a category by ID
func (r *queryResolver) Category(ctx context.Context, id string) (*models.Category, error) {
	category, err := r.DB.GetCategory(id)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// CreateTask creates a new task
func (r *mutationResolver) CreateTask(ctx context.Context, input models.CreateTaskInput) (*models.Task, error) {
	task, err := r.DB.CreateTask(input)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// UpdateTask updates an existing task
func (r *mutationResolver) UpdateTask(ctx context.Context, id string, input models.UpdateTaskInput) (*models.Task, error) {
	task, err := r.DB.UpdateTask(id, input)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// DeleteTask deletes a task
func (r *mutationResolver) DeleteTask(ctx context.Context, id string) (bool, error) {
	if err := r.DB.DeleteTask(id); err != nil {
		return false, err
	}
	return true, nil
}

// UpdateTaskStatus updates a task's status
func (r *mutationResolver) UpdateTaskStatus(ctx context.Context, id string, status models.TaskStatus) (*models.Task, error) {
	task, err := r.DB.UpdateTaskStatus(id, status)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// CreateCategory creates a new category
func (r *mutationResolver) CreateCategory(ctx context.Context, name string) (*models.Category, error) {
	category, err := r.DB.CreateCategory(name)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// UpdateCategory updates an existing category
func (r *mutationResolver) UpdateCategory(ctx context.Context, id string, name string) (*models.Category, error) {
	category, err := r.DB.UpdateCategory(id, name)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// DeleteCategory deletes a category
func (r *mutationResolver) DeleteCategory(ctx context.Context, id string) (bool, error) {
	if err := r.DB.DeleteCategory(id); err != nil {
		return false, err
	}
	return true, nil
}

// Category returns the category associated with a task
func (r *taskResolver) Category(ctx context.Context, obj *models.Task) (*models.Category, error) {
	category, err := r.DB.GetCategory(obj.CategoryID)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Tasks returns all tasks in a category
func (r *categoryResolver) Tasks(ctx context.Context, obj *models.Category) ([]*models.Task, error) {
	tasks, err := r.DB.GetTasksInCategory(obj.ID)
	if err != nil {
		return nil, err
	}

	// Convert to pointer slice
	result := make([]*models.Task, len(tasks))
	for i := range tasks {
		task := tasks[i]
		result[i] = &task
	}
	return result, nil
}
