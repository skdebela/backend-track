package data

import (
	"sync"
	"time"

	"github.com/skdebela/task_manager_api/models"
)

var (
	tasks []models.Task
	mu    sync.RWMutex
	once  sync.Once
)

func initData() {
	tasks = []models.Task{
		{ID: "1", Title: "Task 1", Description: "Description 1", DueDate: time.Now().Add(48*time.Hour), Status: "Pending"},
		{ID: "2", Title: "Task 2", Description: "Description 2", DueDate: time.Now().Add(72*time.Hour), Status: "In Progress"},
		{ID: "3", Title: "Task 3", Description: "Description 3", DueDate: time.Now().Add(24*time.Hour), Status: "Completed"},
	}
}

func GetTasks() []models.Task {
	mu.RLock()
	defer mu.RUnlock()
	once.Do(initData)
	return append([]models.Task{}, tasks...) // return copy
}

func GetTaskByID(id string) (models.Task, bool) {
	mu.RLock()
	defer mu.RUnlock()
	once.Do(initData)
	for _, t := range tasks {
		if t.ID == id {
			return t, true
		}
	}
	return models.Task{}, false
}

func AddTask(task models.Task) {
	mu.Lock()
	defer mu.Unlock()
	once.Do(initData)
	tasks = append(tasks, task)
}

func UpdateTask(id string, updater func(*models.Task)) bool {
	mu.Lock()
	defer mu.Unlock()
	once.Do(initData)
	for i := range tasks {
		if tasks[i].ID == id {
			updater(&tasks[i])
			return true
		}
	}
	return false
}

func DeleteTask(id string) bool {
	mu.Lock()
	defer mu.Unlock()
	once.Do(initData)
	for i, t := range tasks {
		if t.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return true
		}
	}
	return false
}