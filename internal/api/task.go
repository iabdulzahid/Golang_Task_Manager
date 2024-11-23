package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iabdulzahid/golang_task_manager/internal/models"
)

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task with title, description, priority, and due date
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body model.Task true "Task data"
// @Success 201 {object} model.Task
// @Failure 400 {object} gin.H{"error": "Bad request"}
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /tasks [post]
func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate task
	if task.Title == "" || task.Priority == "" || task.DueDate.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	// Set default value for empty labels
	if task.Labels == nil {
		task.Labels = []string{}
	}

	// Create task
	err := createTask(&task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the created task
	c.JSON(http.StatusCreated, task)
}

// GetAllTasks godoc
// @Summary Get all tasks
// @Description Get a list of all tasks in the system
// @Tags tasks
// @Produce json
// @Success 200 {array} model.Task
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /tasks [get]
func GetAllTasks(c *gin.Context) {
	// Call the function to get tasks
	tasks, err := GetTasks()
	if err != nil {
		// If there's an error, return 500 with the error message
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks: " + err.Error()})
		return
	}

	// Return the list of tasks as a JSON response with status 200
	c.JSON(http.StatusOK, tasks)
}

// GetTaskByID godoc
// @Summary Get task by ID
// @Description Get task details by task ID
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} model.Task
// @Failure 404 {object} gin.H{"error": "Task not found"}
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /tasks/{id} [get]
func GetTaskByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	task, err := getTaskByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask godoc
// @Summary Update an existing task
// @Description Update task details by task ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body model.Task true "Task data"
// @Success 200 {object} model.Task
// @Failure 400 {object} gin.H{"error": "Bad request"}
// @Failure 404 {object} gin.H{"error": "Task not found"}
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /tasks/{id} [put]
func UpdateTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update task
	updatedTask, err := updateTask(id, &task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete a task by its ID
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 200 {object} gin.H{"message": "Task deleted"}
// @Failure 404 {object} gin.H{"error": "Task not found"}
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /tasks/{id} [delete]
func DeleteTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := deleteTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
