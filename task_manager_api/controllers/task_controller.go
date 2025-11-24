package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/skdebela/task_manager_api/data"
	"github.com/skdebela/task_manager_api/models"
)

type TaskController struct {
	Service *data.TaskService
}

func NewTaskController(s *data.TaskService) *TaskController {
	return &TaskController{Service: s}
}

func (tc *TaskController) GetTasks(c *gin.Context) {
	tasks, err := tc.Service.GetTasks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (tc *TaskController) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, found, err := tc.Service.GetTaskByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (tc *TaskController) CreateTask(c *gin.Context) {
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

	if err := tc.Service.AddTask(c.Request.Context(), task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Task created successfully",
		"task":    task,
	})
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	set := bson.D{}
    if input.Title != "" {
        set = append(set, bson.E{Key: "title", Value: input.Title})
    }
    if input.Description != "" {
        set = append(set, bson.E{Key: "description", Value: input.Description})
    }
    if !input.DueDate.IsZero() {
        set = append(set, bson.E{Key: "due_date", Value: input.DueDate})
    }
    if input.Status != "" {
        set = append(set, bson.E{Key: "status", Value: input.Status})
    }

    if len(set) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
        return
    }

    updated, err := tc.Service.UpdateTask(c.Request.Context(), id, bson.D{{Key: "$set", Value: set}})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	if !updated {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if _, err := tc.Service.DeleteTask(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
