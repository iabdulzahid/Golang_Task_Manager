package monitor

import (
	"fmt"
	"time"

	zLogger "github.com/iabdulzahid/go-logger/logger"
	"github.com/iabdulzahid/golang_task_manager/internal/database"
	"github.com/iabdulzahid/golang_task_manager/pkg/globals"
)

// monitorOverdueTasks checks tasks and updates their overdue status in the database
func TaskMonitor(logger zLogger.Logger) {
	// logger := globals.Logger
	db := globals.DB

	ticker := time.NewTicker(5 * time.Second)
	logger.Info("TaskMonitor........")
	// Create a ticker that ticks every 12 hours
	// ticker := time.NewTicker(12 * time.Hour)
	// Run a goroutine to periodically check app status
	for {
		logger.Info("TaskMonitor........inside for loop")
		select {
		case <-ticker.C:
			logger.Info("TaskMonitor........inside ticker")
			// Fetch tasks from the database
			tasks, err := database.GetTasks(logger)
			if err != nil {
				logger.Info("TaskMonitor", "err", fmt.Sprintf("Error fetching tasks: %s", err))
				logger.Error("TaskMonitor: Error fetching tasks", err)
				return
			}
			logger.Info("TaskMonitor", "tasks", tasks)
			// Loop through tasks and update the overdue status
			for _, task := range tasks {
				logger.Info("TaskMonitor........", "task", task)
				if task.DueDate != "" {
					dueDate, err := time.Parse(time.RFC3339, task.DueDate)
					if err != nil {
						// logger.Info("TaskMonitor", "msg", fmt.Sprintf("Error parsing due date for task %s: %v", task.ID, err))
						logger.Error(fmt.Sprintf("TaskMonitor: Error parsing due date for task %s", task.ID), err)
						continue
					}
					logger.Info("TaskMonitor", "dueDate", dueDate)
					globals.SetPriorityBasedOnDueDate(logger, &task)
					logger.Info("TaskMonitor", "task.Priority", task.Priority)
					// If the task is overdue, set the IsOverdue flag to true
					if dueDate.Before(time.Now()) {
						// Update the task in the database to reflect the overdue status
						// You would need to run a query to update the task in the database
						// Assuming you have a `db.Exec` function for running SQL queries
						_, err := db.Exec(
							"UPDATE tasks SET is_overdue = $1, priority = $2 WHERE id = $3",
							true, task.Priority, task.ID,
						)
						if err != nil {
							// log.Printf("TaskMonitor: Error updating overdue status for task %s: %v", task.ID, err)
							logger.Error(fmt.Sprintf("TaskMonitor: Error updating overdue status for task %s", task.ID), err)
							continue
						}
						logger.Info("TaskMonitor", "dueDate.Before(time.Now())", dueDate.Before(time.Now()))
					}
				}

			}
			// default:
			// 	logger.Info("TaskMonitor........inside default")
			// 	continue
		}
	}

	logger.Info("TaskMonitor.......exiting ticker")
	defer ticker.Stop()

}
