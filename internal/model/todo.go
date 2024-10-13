package model

import "time"

// Todo represents a task in the todo application.
type Todo struct {
	ID        int       `json:"id"`
	Task      string    `json:"task"`
	Status    Status    `json:"status"`
	Priority  string    `json:"priority" gorm:"type:text;check:priority IN ('Low', 'Medium', 'High')"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NewTodo creates a new Todo item.
func NewTodo(task string, priority string) *Todo {
	return &Todo{
		Task:     task,
		Status:   Created,  // Default status, assuming you have a Pending status
		Priority: priority, // Set priority
	}
}

// NewUpdateTodo returns a new instance of the todo model for updating.
func NewUpdateTodo(id int, task string, status Status) *Todo {
	return &Todo{
		ID:     id,
		Task:   task,
		Status: status,
	}
}

// Status is the status of the task.
type Status string

const (
	// Created is the status for a created task.
	Created = Status("created")
	// Processing is the status for a processing task.
	Processing = Status("processing")
	// Done is the status for a done task.
	Done = Status("done")
)

// StatusMap is a map of task status.
var StatusMap = map[Status]bool{
	Created:    true,
	Processing: true,
	Done:       true,
}
