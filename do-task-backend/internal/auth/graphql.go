// Package auth provides authentication utilities for GraphQL resolvers
package auth

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
)

// UserInfo represents authenticated user information
type UserInfo struct {
	ID            string
	Email         string
	Name          string
	Authenticated bool
}

// GetUserFromGraphQLContext extracts user information from GraphQL context
// This works with the Gin context that's passed through to GraphQL resolvers
func GetUserFromGraphQLContext(ctx context.Context) (*UserInfo, error) {
	// Get Gin context from GraphQL context
	ginCtx, exists := ctx.Value("GinContextKey").(*gin.Context)
	if !exists {
		// Try to get it directly if it's stored differently
		if directCtx, ok := ctx.(*gin.Context); ok {
			ginCtx = directCtx
		} else {
			return &UserInfo{Authenticated: false}, nil
		}
	}

	userID, email, name, authenticated := GetUserFromContext(ginCtx)

	return &UserInfo{
		ID:            userID,
		Email:         email,
		Name:          name,
		Authenticated: authenticated,
	}, nil
}

// Require Authentication checks if user is authenticated and returns user info
func RequireAuthentication(ctx context.Context) (*UserInfo, error) {
	userInfo, err := GetUserFromGraphQLContext(ctx)
	if err != nil {
		return nil, err
	}

	if !userInfo.Authenticated {
		return nil, errors.New("authentication required")
	}

	return userInfo, nil
}
