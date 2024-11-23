package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/iabdulzahid/golang_task_manager/docs" // Import Swagger docs
	"github.com/iabdulzahid/golang_task_manager/internal/api"
	"github.com/iabdulzahid/golang_task_manager/internal/middleware"
	"github.com/iabdulzahid/golang_task_manager/pkg/export"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Task Manager API
// @version 1.0
// @description This is a task management API built with Go (Golang).
// @termsOfService https://example.com/terms
// @contact.name API Support
// @contact.url https://example.com/support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api
func main() {
	// Initialize Database
	err := api.InitDB()
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}

	// Create a new Gin router
	r := gin.Default()

	// Apply rate limiting middleware
	r.Use(middleware.RateLimiter())

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Define routes
	r.POST("/tasks", api.CreateTask)
	r.GET("/tasks", api.GetAllTasks)
	r.GET("/tasks/:id", api.GetTaskByID)
	r.PUT("/tasks/:id", api.UpdateTask)
	r.DELETE("/tasks/:id", api.DeleteTask)
	r.GET("/tasks/export", export.ExportTasks)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
