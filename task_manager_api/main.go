// main.go
package main

import (
	"fmt"

	"github.com/skdebela/task_manager_api/router"
)

func main() {
	r := router.SetupRouter()

	fmt.Println("Task Manager API running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}