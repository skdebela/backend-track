package main

import (
	"library_management/controllers"
	"library_management/services"
	"library_management/concurrency"
)

func main() {
	lib := services.NewLibrary()
	concurrency.StartWorkers(lib, 5)
	cli := controllers.NewCLIController(lib)
	cli.Start()
}
