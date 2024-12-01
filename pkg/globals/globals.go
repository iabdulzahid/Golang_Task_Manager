package globals

import (
	"database/sql"
	"time"

	zLogger "github.com/iabdulzahid/go-logger/logger"
	"github.com/iabdulzahid/golang_task_manager/internal/models"
)

var Logger zLogger.Logger
var DB *sql.DB

func IsValidPriority(priority string) bool {
	switch priority {
	case string(models.Low), string(models.Medium), string(models.High):
		return true
	}
	return false
}

// GetValidPriorityValues returns the valid priority values.
func GetValidPriorityValues() []string {
	return []string{string(models.Low), string(models.Medium), string(models.High)}
}

// func SetPriorityBasedOnDueDate(task models.Task) {
// 	if task.Priority == "" { // Check if priority is not set by the user
// 		now := time.Now()
// 		if task.DueDate != "" {
// 			due, err := time.Parse(time.RFC3339, task.DueDate)
// 			if err == nil {
// 				// Set priority based on due date
// 				daysRemaining := due.Sub(now).Hours() / 24
// 				if daysRemaining <= 1 {
// 					task.Priority = models.High
// 				} else if daysRemaining <= 7 {
// 					task.Priority = models.Medium
// 				} else {
// 					task.Priority = models.Low
// 				}
// 			}
// 		} else {
// 			task.Priority = models.Low // Default priority if no due_date is set
// 		}
// 	}
// }

// Function to set priority based on due_date if not already set
func SetPriorityBasedOnDueDate(logger zLogger.Logger, task *models.Task) {
	// If priority is not set, calculate based on due_date
	// if task.Priority == nil || *task.Priority == "" {
	if task.DueDate != "" {
		// Parse the due_date
		dueDate, err := time.Parse(time.RFC3339, task.DueDate)
		if err != nil {
			// log.Printf("Error parsing due date for task %s: %v", task.ID, err)
			return
		}

		// Calculate the number of days until the due date
		daysRemaining := int(time.Until(dueDate).Hours() / 24)
		logger.Info("SetPriorityBasedOnDueDate......if", "daysRemaining", daysRemaining)
		// Set priority based on the number of days remaining
		switch {
		case daysRemaining <= 2:
			task.Priority = GetAddress(models.High)
		case daysRemaining <= 5:
			task.Priority = GetAddress(models.Medium)
		default:
			task.Priority = GetAddress(models.Low)
		}
	} else {
		logger.Info("SetPriorityBasedOnDueDate......else", "task.Priority", task.Priority)
		// If due_date is not set, set priority to Low
		task.Priority = GetAddress(models.Low)
	}
	// }
	logger.Info("SetPriorityBasedOnDueDate", "task.Priority", task.Priority)
}

func GetAddress[T any](param T) *T {
	return &param
}
