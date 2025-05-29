// Package main is the entry point for the GraphQL API server
package main

import (
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
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
	// Set up the GraphQL handler
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &resolvers.Resolver{DB: db},
	}))

	// Routes
	r.POST("/query", func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/", func(c *gin.Context) {
		playground := playground.Handler("GraphQL playground", "/query")
		playground.ServeHTTP(c.Writer, c.Request)
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
