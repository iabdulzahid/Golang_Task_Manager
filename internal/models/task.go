package models

import "time"

// Task struct for task model
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	DueDate     time.Time `json:"due_date"`
	Labels      []string  `json:"labels"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
