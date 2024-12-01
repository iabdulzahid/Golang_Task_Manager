package export

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	dbFunc "github.com/iabdulzahid/golang_task_manager/internal/database"
	"github.com/iabdulzahid/golang_task_manager/internal/models"
	"github.com/iabdulzahid/golang_task_manager/pkg/globals"
)

// ExportTasks godoc
// @Summary Export tasks to JSON or CSV
// @Description Export all tasks to JSON or CSV format based on the requested file format
// @Tags tasks
// @Produce json
// @Param format query string true "Export format" Enums(json, csv)
// @Success 200 {string} string "File exported successfully"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasks/export [get]
func ExportTasks(c *gin.Context) {
	format := c.DefaultQuery("format", "json")
	logger := globals.Logger
	// Fetch tasks from the database
	tasks, err := dbFunc.GetTasks(logger)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch tasks"})
		return
	}

	// Export based on the requested format (json or csv)
	switch format {
	case "json":
		exportTasksToJSON(c, tasks)
	case "csv":
		exportTasksToCSV(c, tasks)
	default:
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid format. Use 'json' or 'csv'"})
	}
}

func exportTasksToJSON(c *gin.Context, tasks []models.Task) {
	// Set content type and file name for JSON export
	c.Header("Content-Disposition", "attachment; filename=tasks.json")
	c.Header("Content-Type", "application/json")

	// Write tasks to response as JSON
	if err := json.NewEncoder(c.Writer).Encode(tasks); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to export tasks to JSON"})
		return
	}
}

func exportTasksToCSV(c *gin.Context, tasks []models.Task) {
	// Set content type and file name for CSV export
	c.Header("Content-Disposition", "attachment; filename=tasks.csv")
	c.Header("Content-Type", "text/csv")

	// Create CSV writer
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write header
	err := writer.Write([]string{"ID", "Title", "Description", "Priority", "DueDate", "Labels", "CreatedAt", "UpdatedAt"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to write CSV header"})
		return
	}

	// Write task data
	for _, task := range tasks {
		prior := globals.GetAddress(task.Priority)
		err := writer.Write([]string{
			task.ID,
			task.Title,
			task.Description,
			string(**prior),
			task.DueDate,
			strings.Join(task.Labels, ","),
			task.CreatedAt,
			task.UpdatedAt,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to write task data to CSV"})
			return
		}
	}
}
