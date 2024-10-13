// Package service provides the business logic for the todo endpoint.
package service

import (
	"github.com/zuu-development/fullstack-examination-2024/internal/model"
	"github.com/zuu-development/fullstack-examination-2024/internal/repository"
)

// Todo is the service for the todo endpoint.
type Todo interface {
	Create(task string, priority string) (*model.Todo, error)                              // Added priority parameter
	Update(id int, task string, status model.Status, priority string) (*model.Todo, error) // Updated to include priority
	Delete(id int) error
	Find(id int) (*model.Todo, error)
	FindAll(task string, status string) (incompleteTasks []*model.Todo, completedTasks []*model.Todo, err error)
}

type todo struct {
	todoRepository repository.Todo
}

// NewTodo creates a new Todo service.
func NewTodo(r repository.Todo) Todo {
	return &todo{todoRepository: r}
}

// ResponseError represents the structure of the error response.
type ResponseError struct {
	Errors []Error `json:"errors"`
}

// Error represents a single error item.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ResponseData represents the structure of the successful response containing todos.
type ResponseData struct {
	Data []*model.Todo `json:"data"`
}

// Create inserts a new todo into the database.
func (t *todo) Create(task string, priority string) (*model.Todo, error) { // Updated method signature
	todo := model.NewTodo(task, priority) // Pass priority to the model
	if err := t.todoRepository.Create(todo); err != nil {
		return nil, err
	}
	return todo, nil
}

// Update modifies an existing todo in the database.
func (t *todo) Update(id int, task string, status model.Status, priority string) (*model.Todo, error) { // Updated method signature
	currentTodo, err := t.Find(id)
	if err != nil {
		return nil, err
	}

	// Create a new Todo with updated values
	todo := model.NewUpdateTodo(id, task, status) // Assuming this method is modified accordingly
	if todo.Task == "" {
		todo.Task = currentTodo.Task
	}
	if todo.Status == "" {
		todo.Status = currentTodo.Status
	}
	if priority != "" { // Update priority only if it's provided
		todo.Priority = priority
	}

	if err := t.todoRepository.Update(todo); err != nil {
		return nil, err
	}
	return todo, nil
}

// Delete removes a todo by its ID from the database.
func (t *todo) Delete(id int) error {
	if err := t.todoRepository.Delete(id); err != nil {
		return err
	}
	return nil
}

// Find retrieves a todo by its ID from the database.
func (t *todo) Find(id int) (*model.Todo, error) {
	todo, err := t.todoRepository.Find(id)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

// FindAll retrieves all todos with optional filtering by task and status.
// @Summary	Find all todos
// @Tags		todos
// @Param		task	query		string	false	"Task keyword to search"
// @Param		status	query		string	false	"Status to filter by"
// @Success	200	{object}	ResponseData{Data=[]model.Todo}
// @Failure	500	{object}	ResponseError
// @Router		/todos [get]
// FindAll retrieves all todos with optional filtering by task and status.
// Returns separate lists for incomplete and completed tasks.
func (t *todo) FindAll(task string, status string) (incompleteTasks []*model.Todo, completedTasks []*model.Todo, err error) {
	allTodos, err := t.todoRepository.FindAll(task, status) // Call the repository method
	if err != nil {
		return nil, nil, err
	}

	// Separate tasks into completed and incomplete
	for _, todo := range allTodos {
		if todo.Status == model.Done {
			completedTasks = append(completedTasks, todo)
		} else {
			incompleteTasks = append(incompleteTasks, todo)
		}
	}

	return incompleteTasks, completedTasks, nil
}
