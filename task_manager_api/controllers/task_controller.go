package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/skdebela/task_manager_api/data"
	"github.com/skdebela/task_manager_api/models"
)

func GetTasks(c *gin.Context) {
	tasks := data.GetTasks()
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func GetTask(c *gin.Context) {
	id := c.Param("id")
	task, found := data.GetTaskByID(id)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func CreateTask(c *gin.Context) {
	var input struct {
		Title       string    `json:"title" binding:"required"`
		Description string    `json:"description"`
		DueDate     string    `json:"due_date" binding:"required"`
		Status      string    `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dueDate, err := time.Parse(time.RFC3339, input.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due_date format. Use RFC3339"})
		return
	}

	task := models.NewTask(
		uuid.New().String(),
		input.Title,
		input.Description,
		dueDate,
		input.Status,
	)

	data.AddTask(task)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Task created successfully",
		"task":    task,
	})
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated := data.UpdateTask(id, func(t *models.Task) {
		if input.Title != "" {
			t.Title = input.Title
		}
		if input.Description != "" {
			t.Description = input.Description
		}
		if !input.DueDate.IsZero() {
			t.DueDate = input.DueDate
		}
		if input.Status != "" {
			t.Status = input.Status
		}
	})

	if !updated {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if !data.DeleteTask(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
