package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/skdebela/task_manager_api/data"
	"github.com/skdebela/task_manager_api/router"
)

func main() {
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	dbName := getEnv("MONGO_DB", "task_manager_db")
	port := getEnv("PORT", "8080")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	taskSvc, err := data.NewTaskService(ctx, mongoURI, dbName)
	if err != nil {
		log.Fatalf("Failed to initialize TaskService: %v", err)
	}
	defer func() {
		if err := taskSvc.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting TaskService: %v", err)
		}
	}()

	userSvc, err := data.NewUserService(ctx, mongoURI, dbName)
	if err != nil {
		log.Fatalf("Failed to initialize UserService: %v", err)
	}
	defer func() {
		if err := userSvc.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting UserService: %v", err)
		}
	}()

	r := router.SetupRouter(taskSvc, userSvc)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on http://localhost:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}