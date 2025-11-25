package router

import (
	"github.com/gin-gonic/gin"

	"github.com/skdebela/task_manager_api/controllers"
	"github.com/skdebela/task_manager_api/data"
	"github.com/skdebela/task_manager_api/middleware"
)

func SetupRouter(taskService *data.TaskService, userService *data.UserService) *gin.Engine {
	r := gin.Default()

	tc := controllers.NewTaskController(taskService)
	uc := controllers.NewUserController(userService)

	r.POST("/register", uc.Register)
	r.POST("/login", uc.Login)

	authorized := r.Group("/")
	authorized.Use(middleware.AuthRequired(userService))
	{
		authorized.GET("/tasks", tc.GetTasks)
		authorized.GET("/tasks/:id", tc.GetTask)
		
		admin := authorized.Group("/")
		admin.Use(middleware.AdminRequired())
		{
			admin.POST("/tasks", tc.CreateTask)
			admin.PUT("/tasks/:id", tc.UpdateTask)
			admin.DELETE("/tasks/:id", tc.DeleteTask)
			admin.POST("/promote/:id", uc.PromoteUser)


		}
	}

	return r
}