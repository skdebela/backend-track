package models

import "time"


type Task struct {
	ID		  	string		`json:"id"`
	Title       string		`json:"title"`
	Description string		`json:"description"`
	DueDate	 	time.Time	`json:"due_date"`
	Status    	string		`json:"status"`
}

func NewTask(id, title, description string, dueDate time.Time, status string) Task {
	return Task{
		ID:          id,
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		Status:      status,
	}
}