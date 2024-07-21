package main

import (
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/joshua468/hng_boilerplate_golang_web/internal/routes"
    "github.com/joshua468/hng_boilerplate_golang_web/utility"
)

func main() {
    // Load environment variables
    if err := utility.LoadEnv(); err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Initialize Gin router
    router := gin.Default()

    // Set up routes
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        log.Fatal("DATABASE_URL environment variable not set")
    }
    db, err := utility.ConnectDatabase(dbURL)
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }

    // Set up routes and middleware
    routes.SetupRoutes(router, db)

    // Start the server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Default to port 8080 if not set
    }
    log.Printf("Starting server on port %s...", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
}
