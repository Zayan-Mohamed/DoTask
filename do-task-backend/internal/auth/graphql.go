// Package auth provides authentication utilities for GraphQL resolvers
package auth

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	ID            string
	Email         string
	Name          string
	Authenticated bool
}

func GetUserFromGraphQLContext(ctx context.Context) (*UserInfo, error) {
	// Get Gin context from GraphQL context
	ginCtx, exists := ctx.Value("GinContextKey").(*gin.Context)
	if !exists {
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
