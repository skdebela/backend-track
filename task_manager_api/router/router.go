package router

import (
	"github.com/gin-gonic/gin"

	"github.com/skdebela/task_manager_api/controllers"
	"github.com/skdebela/task_manager_api/data"
)

func SetupRouter(taskService *data.TaskService) *gin.Engine {
	r := gin.Default()

	tc := controllers.NewTaskController(taskService)

	r.GET("/tasks", tc.GetTasks)
	r.GET("/tasks/:id", tc.GetTask)
	r.POST("/tasks", tc.CreateTask)
	r.PUT("/tasks/:id", tc.UpdateTask)
	r.DELETE("/tasks/:id", tc.DeleteTask)

	return r
}