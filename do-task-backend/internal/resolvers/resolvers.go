// Package resolvers implements the GraphQL resolvers
package resolvers

import (
	"context"
	"fmt"
	"time"

	"github.com/Zayan-Mohamed/do-task-backend/internal/auth"
	"github.com/Zayan-Mohamed/do-task-backend/internal/database"
	"github.com/Zayan-Mohamed/do-task-backend/internal/graph/generated"
	"github.com/Zayan-Mohamed/do-task-backend/internal/models"
	"github.com/gin-gonic/gin"
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

// User returns the user resolver
func (r *Resolver) User() generated.UserResolver {
	return &userResolver{r}
}

// Tasks returns all tasks for the authenticated user
func (r *queryResolver) Tasks(ctx context.Context) ([]*models.Task, error) {
	userInfo, err := auth.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	tasks, err := r.DB.GetAllTasksByUser(userInfo.ID)
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

// Task returns a task by ID for the authenticated user
func (r *queryResolver) Task(ctx context.Context, id string) (*models.Task, error) {
	userInfo, err := auth.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	task, err := r.DB.GetTask(id, userInfo.ID)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// Categories returns all categories for the authenticated user
func (r *queryResolver) Categories(ctx context.Context) ([]*models.Category, error) {
	userInfo, err := auth.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	categories, err := r.DB.GetAllCategories(userInfo.ID)
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

// Category returns a category by ID for the authenticated user
func (r *queryResolver) Category(ctx context.Context, id string) (*models.Category, error) {
	userInfo, err := auth.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	category, err := r.DB.GetCategory(id, userInfo.ID)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// CreateTask creates a new task for the authenticated user
func (r *mutationResolver) CreateTask(ctx context.Context, input models.CreateTaskInput) (*models.Task, error) {
	userInfo, err := auth.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	// Set the user ID from authentication context
	input.UserID = userInfo.ID

	// Validate that categoryId is not empty
	if input.CategoryID == "" {
		// Attempt to find a category for this user
		categories, err := r.DB.GetAllCategories(userInfo.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to find categories: %w", err)
		}

		if len(categories) == 0 {
			return nil, fmt.Errorf("cannot create task without a category. Please create a category first")
		}

		// Use the first available category
		input.CategoryID = categories[0].ID
	}

	task, err := r.DB.CreateTask(input)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// UpdateTask updates an existing task for the authenticated user
func (r *mutationResolver) UpdateTask(ctx context.Context, id string, input models.UpdateTaskInput) (*models.Task, error) {
	userInfo, err := auth.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	// Check if categoryId is provided and empty
	if input.CategoryID != nil && *input.CategoryID == "" {
		// If an empty string is provided, find a valid category
		categories, err := r.DB.GetAllCategories(userInfo.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to find categories: %w", err)
		}

		if len(categories) > 0 {
			// Use the first available category
			input.CategoryID = &categories[0].ID
		} else {
			// If no categories exist, set to nil to avoid empty string error
			input.CategoryID = nil
		}
	}

	task, err := r.DB.UpdateTask(id, input, userInfo.ID)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// DeleteTask deletes a task for the authenticated user
func (r *mutationResolver) DeleteTask(ctx context.Context, id string) (bool, error) {
	userInfo, err := auth.RequireAuthentication(ctx)
	if err != nil {
		return false, err
	}

	if err := r.DB.DeleteTask(id, userInfo.ID); err != nil {
		return false, err
	}
	return true, nil
}

// UpdateTaskStatus updates a task's status for the authenticated user
func (r *mutationResolver) UpdateTaskStatus(ctx context.Context, id string, status models.TaskStatus) (*models.Task, error) {
	userInfo, err := auth.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	task, err := r.DB.UpdateTaskStatus(id, status, userInfo.ID)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// CreateCategory creates a new category for the authenticated user
func (r *mutationResolver) CreateCategory(ctx context.Context, name string) (*models.Category, error) {
	userInfo, err := auth.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	category, err := r.DB.CreateCategory(name, userInfo.ID)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// UpdateCategory updates an existing category for the authenticated user
func (r *mutationResolver) UpdateCategory(ctx context.Context, id string, name string) (*models.Category, error) {
	userInfo, err := auth.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	category, err := r.DB.UpdateCategory(id, name, userInfo.ID)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// DeleteCategory deletes a category for the authenticated user
func (r *mutationResolver) DeleteCategory(ctx context.Context, id string) (bool, error) {
	userInfo, err := auth.RequireAuthentication(ctx)
	if err != nil {
		return false, err
	}

	if err := r.DB.DeleteCategory(id, userInfo.ID); err != nil {
		return false, err
	}
	return true, nil
}

// Category returns the category associated with a task
func (r *taskResolver) Category(ctx context.Context, obj *models.Task) (*models.Category, error) {
	userInfo, err := auth.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	category, err := r.DB.GetCategory(obj.CategoryID, userInfo.ID)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Tasks returns all tasks in a category for the authenticated user
func (r *categoryResolver) Tasks(ctx context.Context, obj *models.Category) ([]*models.Task, error) {
	userInfo, err := auth.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	tasks, err := r.DB.GetTasksInCategory(obj.ID, userInfo.ID)
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

// DueDate returns a string representation of the task's due date
func (r *taskResolver) DueDate(ctx context.Context, obj *models.Task) (string, error) {
	return obj.DueDate.Format(time.RFC3339), nil
}

// CreatedAt returns a string representation of the task's creation time
func (r *taskResolver) CreatedAt(ctx context.Context, obj *models.Task) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}

// UpdatedAt returns a string representation of the task's update time
func (r *taskResolver) UpdatedAt(ctx context.Context, obj *models.Task) (string, error) {
	return obj.UpdatedAt.Format(time.RFC3339), nil
}

// Authentication resolvers

// Register creates a new user account
func (r *mutationResolver) Register(ctx context.Context, input models.RegisterInput) (*models.AuthResponse, error) {
	// Hash the password
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user in database
	user, err := r.DB.CreateUser(input.Name, input.Email, hashedPassword)
	if err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := auth.GenerateToken(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Set the token as an HTTP-only cookie if Gin context is available
	if ginContext, exists := ctx.Value("GinContextKey").(*gin.Context); exists {
		auth.SetTokenCookie(ginContext, token)
	}

	return &models.AuthResponse{
		User:  &user,
		Token: token,
	}, nil
}

// Login authenticates a user
func (r *mutationResolver) Login(ctx context.Context, input models.LoginInput) (*models.AuthResponse, error) {
	// Get user by email
	user, err := r.DB.GetUserByEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check password
	if !auth.CheckPassword(input.Password, user.Password) {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate JWT token
	token, err := auth.GenerateToken(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Set the token as an HTTP-only cookie if Gin context is available
	if ginContext, exists := ctx.Value("GinContextKey").(*gin.Context); exists {
		auth.SetTokenCookie(ginContext, token)
	}

	// Don't return password in response
	user.Password = ""

	return &models.AuthResponse{
		User:  &user,
		Token: token,
	}, nil
}

// Me returns the current authenticated user
func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
	userInfo, err := auth.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	user, err := r.DB.GetUserByID(userInfo.ID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

type categoryResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type taskResolver struct{ *Resolver }
type userResolver struct{ *Resolver }

// CreatedAt resolves the createdAt field for User
func (r *userResolver) CreatedAt(ctx context.Context, obj *models.User) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}

// UpdatedAt resolves the updatedAt field for User
func (r *userResolver) UpdatedAt(ctx context.Context, obj *models.User) (string, error) {
	return obj.UpdatedAt.Format(time.RFC3339), nil
}
