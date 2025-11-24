package models

import "time"


type Task struct {
	ID		  	string		`json:"id" bson:"id"`
	Title       string		`json:"title" bson:"title"`
	Description string		`json:"description" bson:"description"`
	DueDate	 	time.Time	`json:"due_date" bson:"due_date"`
	Status    	string		`json:"status" bson:"status"`
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