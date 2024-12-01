package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	zLogger "github.com/iabdulzahid/go-logger/logger"
	_ "github.com/iabdulzahid/golang_task_manager/docs" // Import Swagger docs
	"github.com/iabdulzahid/golang_task_manager/internal/api"
	taskDB "github.com/iabdulzahid/golang_task_manager/internal/database"
	"github.com/iabdulzahid/golang_task_manager/internal/export"
	"github.com/iabdulzahid/golang_task_manager/internal/middleware"
	"github.com/iabdulzahid/golang_task_manager/internal/monitor"
	"github.com/iabdulzahid/golang_task_manager/pkg/globals"
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
	db, err := taskDB.InitDB()
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}

	goLogger, err := zLogger.NewLogger(
		zLogger.Config{
			AppName:            "golang-task-manager",
			LogLevel:           "info",
			LogFormat:          "json",
			LogFilePath:        fmt.Sprintf("./logs/%s.log", "golang-task-manager"),
			LogFilePermissions: "0644",
			TimeFormat:         "2006-01-02T15:04:05Z07:005",
			LogOutput:          []string{"stdout", "file"},
		},
	)

	if err != nil {
		log.Fatal("Failed to create logger: ", err)
	}

	globals.DB = db

	logger := goLogger.WithContext(
		zap.String("requestId", "12345"),
		zap.String("userId", "admin"),
	)
	logger.Debug("welcome to Golang Task Manager", "time", time.Now())
	globals.Logger = *logger
	// Create a new Gin router
	r := gin.Default()

	// Apply rate limiting middleware
	r.Use(middleware.RateLimiter())

	go monitor.TaskMonitor(*logger)

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
