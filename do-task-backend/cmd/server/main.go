// Package main is the entry point for the GraphQL API server
package main

import (
	"context"
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Zayan-Mohamed/do-task-backend/internal/auth"
	"github.com/Zayan-Mohamed/do-task-backend/internal/database"
	"github.com/Zayan-Mohamed/do-task-backend/internal/graph/generated"
	"github.com/Zayan-Mohamed/do-task-backend/internal/resolvers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const defaultPort = "8080"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Failed to load .env file: %v", err)
	}

	auth.InitJWT()

	// Set up the database connection
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if err := db.InitSchema(); err != nil {
		log.Printf("Warning: Failed to initialize sample data: %v", err)
	}

	// Create a new Gin router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Add authentication middleware
	r.Use(auth.AuthMiddleware())

	// Set up the GraphQL handler
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &resolvers.Resolver{DB: db},
	}))

	// Routes
	r.POST("/query", func(c *gin.Context) {
		// Add Gin context to the request context so GraphQL resolvers can access it
		ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		srv.ServeHTTP(c.Writer, c.Request)
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "dotask-backend"})
	})
	r.HEAD("/health", func(c *gin.Context) {
		c.Status(200)
	})

	// Root endpoint - GraphQL playground (also handle HEAD for health checks)
	r.GET("/", func(c *gin.Context) {
		playground := playground.Handler("GraphQL playground", "/query")
		playground.ServeHTTP(c.Writer, c.Request)
	})
	r.HEAD("/", func(c *gin.Context) {
		c.Status(200)
	})

	// Determine port to run the server on
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Start the server
	log.Printf("Server is running on http://localhost:%s/", port)
	log.Printf("GraphQL playground is available at http://localhost:%s/", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
