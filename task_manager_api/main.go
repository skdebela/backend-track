// main.go
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
	mongoURI := os.Getenv("MONGO_URI")
    if mongoURI == "" {
        mongoURI = "mongodb://localhost:27017"
    }
    dbName := os.Getenv("MONGO_DB")
    if dbName == "" {
        dbName = "task_manager_db"
    }
    collName := os.Getenv("MONGO_COLLECTION")
    if collName == "" {
        collName = "tasks"
    }

    ctx := context.Background()
    svc, err := data.NewTaskService(ctx, mongoURI, dbName)
    if err != nil {
        log.Fatalf("failed to connect to mongo: %v", err)
    }
    defer func() {
        _ = svc.Disconnect(ctx)
    }()

    r := router.SetupRouter(svc)

    srv := &http.Server{
        Addr:    "8080",
        Handler: r,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    go func() {
        log.Printf("server listening on :8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()

    // graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("shutting down server...")

    ctxShut, cancel := context.WithTimeout(ctx, 15*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctxShut); err != nil {
        log.Fatalf("server forced to shutdown: %v", err)
    }

    if err := svc.Disconnect(ctxShut); err != nil {
        log.Printf("error disconnecting mongo: %v", err)
    }
    log.Println("server exited properly")
}