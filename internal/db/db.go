package db

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/iabdulzahid/golang_task_manager/internal/models"
	"github.com/iabdulzahid/golang_task_manager/pkg/globals"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var db *sql.DB

// InitDB initializes the PostgreSQL database connection
// InitDB initializes the database connection and ensures the "tasks" table exists.
func InitDB() (*sql.DB, error) {
	// Get the database URL from environment variables
	databaseURL := os.Getenv("DATABASE_URL")
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
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			priority TEXT,
			due_date TEXT,
			labels TEXT,
			created_at TEXT,
			updated_at TEXT
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
	// Set timestamps
	task.CreatedAt = time.Now().Format(time.RFC3339) // Format time as string
	task.UpdatedAt = time.Now().Format(time.RFC3339)
	task.DueDate = time.Now().Add(10 * time.Minute).Format(time.RFC3339)
	logger := globals.Logger
	// Convert Labels to a comma-separated string
	labelsStr := strings.Join(task.Labels, ",")
	logger.Info("CreateTask", "task", task)
	// Prepare the SQL query
	query := `
		INSERT INTO tasks (title, description, priority, due_date, labels, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := db.Exec(query, task.Title, task.Description, task.Priority, task.DueDate, labelsStr, task.CreatedAt, task.UpdatedAt)
	if err != nil {
		logger.Info("CreateTask", "err", err.Error())
		log.Printf("Failed to create task: %v\n", err)
		return err
	}

	log.Println("Task created successfully")
	return nil
}

// GetAllTasks retrieves all tasks from the database
func GetTasks() ([]models.Task, error) {
	rows, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.DueDate, &task.Labels)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// GetTaskByID retrieves a task by ID
func GetTaskByID(id int) (*models.Task, error) {
	row := db.QueryRow("SELECT * FROM tasks WHERE id = $1", id)
	var task models.Task
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.DueDate, &task.Labels)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("task not found")
		}
		return nil, err
	}
	return &task, nil
}

// UpdateTask updates an existing task by ID
func UpdateTask(id int, task *models.Task) (*models.Task, error) {
	_, err := db.Exec(`UPDATE tasks SET title = $1, description = $2, priority = $3, due_date = $4, labels = $5 WHERE id = $6`,
		task.Title, task.Description, task.Priority, task.DueDate, task.Labels, id)
	if err != nil {
		return nil, err
	}
	return GetTaskByID(id)
}

// DeleteTask deletes a task by ID
func DeleteTask(id int) error {
	_, err := db.Exec("DELETE FROM tasks WHERE id = $1", id)
	return err
}
