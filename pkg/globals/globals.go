package globals

import (
	"database/sql"

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
