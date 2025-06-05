// Package auth provides authentication middleware and utilities
package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Zayan-Mohamed/do-task-backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret []byte

func InitJWT() {
	// Load JWT secret from environment variable
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	jwtSecret = []byte(secret)
}

// CustomClaims represents the JWT token claims
type CustomClaims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for a user
func GenerateToken(user *models.User) (string, error) {
	claims := CustomClaims{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword compares a password with its hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// AuthMiddleware provides JWT authentication middleware for GraphQL
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from cookies first (for HTTP-only cookies)
		tokenString, err := c.Cookie("accessToken")
		if err != nil {
			// If no cookie, try Authorization header
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString == "" {
			c.Next()
			return
		}

		// Validate token
		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		// Add user info to Gin - context
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("userName", claims.Name)
		c.Set("authenticated", true)

		c.Next()
	}
}

// RequireAuth middleware that requires authentication
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authenticated, exists := c.Get("authenticated")
		if !exists || !authenticated.(bool) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetUserFromContext extracts user information from Gin context
func GetUserFromContext(c *gin.Context) (userID, email, name string, authenticated bool) {
	if auth, exists := c.Get("authenticated"); exists && auth.(bool) {
		if uid, exists := c.Get("userID"); exists {
			userID = uid.(string)
		}
		if e, exists := c.Get("userEmail"); exists {
			email = e.(string)
		}
		if n, exists := c.Get("userName"); exists {
			name = n.(string)
		}
		authenticated = true
	}
	return
}

// SetTokenCookie sets the JWT token as an HTTP-only cookie
func SetTokenCookie(c *gin.Context, token string) {
	// Set cookie to expire in 7 days (same as token expiration)
	c.SetCookie(
		"accessToken",
		token,
		int(7*24*time.Hour.Seconds()), // 7 days in seconds
		"/",                           // Path
		"",                            // Domain (empty for current domain)
		false,                         // Secure (should be true in production with HTTPS)
		true,                          // HTTP only
	)
}
