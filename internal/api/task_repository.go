package api

// import (
// 	"database/sql"
// 	"errors"
// 	"log"
// 	"os"

// 	_ "github.com/lib/pq" // PostgreSQL driver

// 	"github.com/iabdulzahid/golang_task_manager/internal/models"
// )

// var db *sql.DB

// // InitDB initializes the PostgreSQL database connection
// func InitDB() error {
// 	var err error

// 	// Get the database URL from environment variables
// 	databaseURL := os.Getenv("DATABASE_URL")
// 	if databaseURL == "" {
// 		// Example for local development
// 		databaseURL = "postgres://username:password@localhost:5432/tasksdb?sslmode=disable"
// 	}

// 	// Open the database connection
// 	db, err = sql.Open("postgres", databaseURL)
// 	if err != nil {
// 		return err
// 	}

// 	// Create table if not exists
// 	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
// 		id SERIAL PRIMARY KEY,
// 		title TEXT NOT NULL,
// 		description TEXT,
// 		priority TEXT,
// 		due_date DATE,
// 		labels TEXT
// 	)`)
// 	if err != nil {
// 		log.Fatal("Error creating table:", err)
// 	}
// 	return nil
// }

// // CreateTask inserts a new task into the database
// func createTask(task *models.Task) error {
// 	_, err := db.Exec(`INSERT INTO tasks (title, description, priority, due_date, labels)
// 	VALUES ($1, $2, $3, $4, $5)`, task.Title, task.Description, task.Priority, task.DueDate, task.Labels)
// 	return err
// }

// // GetAllTasks retrieves all tasks from the database
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

// // GetTaskByID retrieves a task by ID
// func getTaskByID(id int) (*models.Task, error) {
// 	row := db.QueryRow("SELECT * FROM tasks WHERE id = $1", id)
// 	var task models.Task
// 	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.DueDate, &task.Labels)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, errors.New("task not found")
// 		}
// 		return nil, err
// 	}
// 	return &task, nil
// }

// // UpdateTask updates an existing task by ID
// func updateTask(id int, task *models.Task) (*models.Task, error) {
// 	_, err := db.Exec(`UPDATE tasks SET title = $1, description = $2, priority = $3, due_date = $4, labels = $5 WHERE id = $6`,
// 		task.Title, task.Description, task.Priority, task.DueDate, task.Labels, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return getTaskByID(id)
// }

// // DeleteTask deletes a task by ID
// func deleteTask(id int) error {
// 	_, err := db.Exec("DELETE FROM tasks WHERE id = $1", id)
// 	return err
// }
