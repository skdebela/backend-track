package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID		  	string		`json:"id"`
	Title       string		`json:"title"`
	Description string		`json:"description"`
	DueDate	 	time.Time	`json:"due_date"`
	Status    	string		`json:"status"`
}

var tasks = []Task{
	{ID: "1", Title: "Task 1", Description: "Description 1", DueDate: time.Now().Add(48 * time.Hour), Status: "Pending"},
	{ID: "2", Title: "Task 2", Description: "Description 2", DueDate: time.Now().Add(72 * time.Hour), Status: "In Progress"},
	{ID: "3", Title: "Task 3", Description: "Description 3", DueDate: time.Now().Add(24 * time.Hour), Status: "Completed"},
}

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/tasks", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})
	})

	router.GET("/tasks/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		for _, task := range tasks {
			if task.ID == id {
				ctx.JSON(http.StatusOK, gin.H{"task": task})
				return
			}
		}
		ctx.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
	})

	router.PUT("/tasks/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		var updatedTask Task

		if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"errror": err.Error()})
			return
		}

		for i, task := range tasks {
			if task.ID == id {
				if updatedTask.Title != "" {
					tasks[i].Title = updatedTask.Title
				}
				if updatedTask.Description != "" {
					tasks[i].Description = updatedTask.Description
				}
				ctx.JSON(http.StatusOK, gin.H{"message": "Task updated"})
				return
			}
		}

		ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
	})

	router.DELETE("/tasks/:id", func(ctx *gin.Context) {
		id := ctx.Param(("id"))

		for i, val := range tasks {
			if val.ID == id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				ctx.JSON(http.StatusOK, gin.H{"message": "Task removed"})
				return
			}
		}

		ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
	})

	router.POST("/tasks/:id", func(ctx *gin.Context) {
		var newTask Task

		if err := ctx.ShouldBindJSON(&newTask); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tasks = append(tasks, newTask)
		ctx.JSON(http.StatusCreated, gin.H{"message": "Task creaed"})
	})

	err := router.Run(":8080")
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}