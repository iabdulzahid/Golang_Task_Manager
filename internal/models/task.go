package models

// Task struct for task model
type Task struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    string   `json:"priority"`
	DueDate     string   `json:"due_date"`
	Labels      []string `json:"labels"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

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
