package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	zLogger "github.com/iabdulzahid/go-logger/logger"
	"github.com/iabdulzahid/golang_task_manager/internal/models"
	"github.com/iabdulzahid/golang_task_manager/pkg/globals"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var db *sql.DB

// InitDB initializes the database connection and ensures the "tasks" table exists.
func InitDB() (*sql.DB, error) {
	// Get the database URL from environment variables
	// databaseURL := os.Getenv("DATABASE_URL")
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Now you can access environment variables
	databaseURL := os.Getenv("DATABASE_URL")
	fmt.Println("DATABASE_URL:", databaseURL)
	fmt.Println("InitDB.........: " + databaseURL)
	if databaseURL == "" {
		// Default URL for local development
		databaseURL = "postgres://postgres:abc123@172.19.35.0:5432/taskdb?sslmode=disable"
	}

	// Open the database connection
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
		return nil, err
	}

	// Test the database connection
	err = db.Ping()
	if err != nil {
		log.Printf("Failed to ping the database: %v\n", err)
		return nil, err
	}
	log.Println("Database connection established successfully.")

	// Check if the "tasks" table exists
	var tableName string
	query := `
	SELECT table_name 
	FROM information_schema.tables 
	WHERE table_schema = 'public' AND table_name = 'tasks';
	`
	err = db.QueryRow(query).Scan(&tableName)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Table 'tasks' does not exist. Creating it now...")
		} else {
			log.Printf("Failed to check table existence: %v\n", err)
			return nil, err
		}
	}

	// Create the table if it does not exist
	if tableName != "tasks" {
		createTableQuery := `
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,       -- Changed to TEXT and still keeps it as PRIMARY KEY
			title TEXT NOT NULL,
			description TEXT,
			priority TEXT,
			due_date TEXT,
			labels TEXT,
			created_at TEXT,
			updated_at TEXT,
			is_overdue BOOLEAN DEFAULT FALSE   -- Added is_overdue column with a default value of FALSE
		);`
		_, err = db.Exec(createTableQuery)
		if err != nil {
			log.Printf("Failed to create table: %v\n", err)
			return nil, err
		}
		log.Println("Table 'tasks' created successfully.")
	} else {
		log.Println("Table 'tasks' already exists. Skipping creation.")
	}

	return db, nil
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	return db
}

// CreateTask inserts a new task into the database
func CreateTask(task *models.Task) error {
	db := globals.DB
	logger := globals.Logger
	// Generate a unique ID (e.g., UUID)
	task.ID = uuid.New().String() // Assign a new UUID string to the task ID

	// Set timestamps
	task.CreatedAt = time.Now().Format(time.RFC3339) // Format time as string
	task.UpdatedAt = time.Now().Format(time.RFC3339)
	task.DueDate = time.Now().Add(10 * time.Minute).Format(time.RFC3339)

	// Convert Labels to a comma-separated string
	labelsStr := strings.Join(task.Labels, ",")

	globals.SetPriorityBasedOnDueDate(logger, task)

	// Prepare the SQL query to insert the task
	query := `
		INSERT INTO tasks (id, title, description, priority, due_date, labels, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := db.Exec(query, task.ID, task.Title, task.Description, task.Priority, task.DueDate, labelsStr, task.CreatedAt, task.UpdatedAt)
	if err != nil {
		log.Printf("Failed to create task: %v\n", err)
		return err
	}

	log.Println("Task created successfully")
	return nil
}

// GetAllTasks retrieves all tasks from the database
// func GetTasks() ([]models.Task, error) {
// 	rows, err := db.Query("SELECT * FROM tasks")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var tasks []models.Task
// 	for rows.Next() {
// 		var task models.Task
// 		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.DueDate, &task.Labels)
// 		if err != nil {
// 			return nil, err
// 		}
// 		tasks = append(tasks, task)
// 	}
// 	return tasks, nil
// }

// GetTasks retrieves tasks from the database, sorted by priority (High > Medium > Low).
// GetTasks retrieves tasks from the database, sorted by priority (High > Medium > Low).
func GetTasks(logger zLogger.Logger) ([]models.Task, error) {
	// Ensure that the db object is initialize
	db := globals.DB
	if db == nil {
		log.Println("Database connection is nil")
		return nil, fmt.Errorf("database connection is nil")
	}

	// Query to retrieve tasks sorted by priority
	query := `
        SELECT id, title, description, priority, due_date, labels, created_at, updated_at, is_overdue
        FROM tasks
        ORDER BY CASE
            WHEN priority = 'High' THEN 1
            WHEN priority = 'Medium' THEN 2
            WHEN priority = 'Low' THEN 3
            ELSE 4
        END`

	// Execute the query to retrieve tasks
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error fetching tasks: %v", err)
		return nil, fmt.Errorf("failed to fetch tasks from database: %v", err)
	}
	defer rows.Close()

	// Store the tasks
	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var labelsStr string // Temporarily hold the labels as a string

		// Scan the results into the task struct and labelsStr for labels column
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.DueDate, &labelsStr, &task.CreatedAt, &task.UpdatedAt, &task.IsOverdue)
		if err != nil {
			log.Printf("Error scanning task: %v", err)
			continue // Skip this task and continue with the next one
		}

		// Split the labels string into a slice of strings
		task.Labels = strings.Split(labelsStr, ",") // Convert comma-separated string to a slice
		// Set priority if it's not already set (based on due_date)
		globals.SetPriorityBasedOnDueDate(logger, &task)
		tasks = append(tasks, task)
	}

	// Check for errors after iterating through the rows
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	// Return the retrieved tasks
	return tasks, nil
}

// GetTaskByID retrieves a task by ID
func GetTaskByID(taskId string) (*models.Task, error) {
	db := globals.DB
	row := db.QueryRow("SELECT id, title, description, priority, due_date, labels, created_at, updated_at FROM tasks WHERE id = $1", taskId)

	var task models.Task
	var labelsStr string // Use a temporary variable for labels

	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.DueDate, &labelsStr, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	// Convert labels string to a slice of strings
	task.Labels = strings.Split(labelsStr, ",")

	return &task, nil
}

// UpdateTask updates an existing task by ID
func UpdateTask(taskId string, task *models.Task) (*models.Task, error) {
	db := globals.DB

	// Convert labels slice to a comma-separated string
	labelsStr := strings.Join(task.Labels, ",")

	// Execute the update query
	_, err := db.Exec(`UPDATE tasks SET title = $1, description = $2, priority = $3, due_date = $4, labels = $5 WHERE id = $6`,
		task.Title, task.Description, task.Priority, task.DueDate, labelsStr, taskId)
	if err != nil {
		return nil, err
	}

	// Fetch and return the updated task
	return GetTaskByID(taskId)
}

func UpdateTaskPriority(taskID string, newPriority string) error {
	if !globals.IsValidPriority(newPriority) {
		return fmt.Errorf("invalid priority: %s. Valid values are: %v", newPriority, globals.GetValidPriorityValues())
	}
	db := globals.DB
	query := `UPDATE tasks SET priority = $1, updated_at = $2 WHERE id = $3`
	_, err := db.Exec(query, newPriority, time.Now().Format(time.RFC3339), taskID)
	if err != nil {
		log.Printf("Error updating task priority: %v", err)
		return err
	}

	log.Println("Task priority updated successfully")
	return nil
}

// DeleteTask deletes a task by ID
func DeleteTask(taskId string) error {
	db := globals.DB
	_, err := db.Exec("DELETE FROM tasks WHERE id = $1", taskId)
	return err
}
