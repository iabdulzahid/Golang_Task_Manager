package models

// Task struct for task model
type Task struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    Priority `json:"priority" enum:"Low,Medium,High"` // Swagger annotation for enum
	DueDate     string   `json:"due_date"`
	Labels      []string `json:"labels"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

// Define the custom type for Priority
type Priority string

// Define constants for the priority values
const (
	Low    Priority = "Low"
	Medium Priority = "Medium"
	High   Priority = "High"
)

// SuccessResponse is a struct for success responses
type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type SuccessMessage struct {
	Message string `json:"message"`
}

// ErrorResponse struct for standardized error message format
type ErrorResponse struct {
	Error string `json:"error"`
}
