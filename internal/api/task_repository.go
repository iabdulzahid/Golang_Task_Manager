package api

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/iabdulzahid/golang_task_manager/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB() error {
	var err error
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "./tasks.db"
	}

	db, err = sql.Open("sqlite3", databaseURL)
	if err != nil {
		return err
	}

	// Create table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		description TEXT,
		priority TEXT,
		due_date TEXT,
		labels TEXT
	)`)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}
	return nil
}

// CreateTask inserts a new task into the database
func createTask(task *models.Task) error {
	_, err := db.Exec(`INSERT INTO tasks (title, description, priority, due_date, labels) 
	VALUES (?, ?, ?, ?, ?)`, task.Title, task.Description, task.Priority, task.DueDate, task.Labels)
	return err
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
func getTaskByID(id int) (*models.Task, error) {
	row := db.QueryRow("SELECT * FROM tasks WHERE id = ?", id)
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
func updateTask(id int, task *models.Task) (*models.Task, error) {
	_, err := db.Exec(`UPDATE tasks SET title = ?, description = ?, priority = ?, due_date = ?, labels = ? WHERE id = ?`,
		task.Title, task.Description, task.Priority, task.DueDate, task.Labels, id)
	if err != nil {
		return nil, err
	}
	return getTaskByID(id)
}

// DeleteTask deletes a task by ID
func deleteTask(id int) error {
	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}
