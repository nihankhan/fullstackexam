// Package repository provides the database operations for the todo endpoint.
package repository

import (
	log "github.com/sirupsen/logrus"
	"github.com/zuu-development/fullstack-examination-2024/internal/model"
	"gorm.io/gorm"
)

// Todo is the repository for the todo endpoint.
type Todo interface {
	Create(t *model.Todo) error
	Delete(id int) error
	Update(t *model.Todo) error
	Find(id int) (*model.Todo, error)
	FindAll(task string, status string) ([]*model.Todo, error) // Modified to include parameters
}

// todo is the implementation of the Todo repository.
type todo struct {
	db *gorm.DB
}

// NewTodo returns a new instance of the todo repository.
func NewTodo(db *gorm.DB) Todo {
	return &todo{
		db: db,
	}
}

// Create inserts a new todo into the database.
func (td *todo) Create(t *model.Todo) error {
	if err := td.db.Create(t).Error; err != nil {
		return err
	}
	return nil
}

// Update modifies an existing todo in the database.
func (td *todo) Update(t *model.Todo) error {
	if err := td.db.Save(t).Error; err != nil {
		return err
	}
	return nil
}

// Delete removes a todo by its ID from the database.
func (td *todo) Delete(id int) error {
	result := td.db.Where("id = ?", id).Delete(&model.Todo{})
	if result.RowsAffected == 0 {
		return model.ErrNotFound
	}
	if result.Error != nil {
		return result.Error
	}
	log.Info("Deleted todo with id: ", id)
	return nil
}

// Find retrieves a todo by its ID from the database.
func (td *todo) Find(id int) (*model.Todo, error) {
	var todo model.Todo // Changed to a value instead of pointer
	err := td.db.Where("id = ?", id).Take(&todo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &todo, nil // Return pointer to todo
}

// FindAll retrieves all todos with optional filtering by task and status.
func (td *todo) FindAll(task string, status string) ([]*model.Todo, error) {
	var todos []*model.Todo
	query := td.db.Model(&model.Todo{})

	// Apply filtering based on task and status
	if task != "" {
		query = query.Where("task LIKE ?", "%"+task+"%") // Use LIKE for partial matching
	}
	if status != "" {
		query = query.Where("status = ?", status) // Direct match for status
	}

	// Execute the query and fetch results
	if err := query.Find(&todos).Error; err != nil {
		log.Errorf("Error retrieving todos: %v", err) // Log unexpected errors
		return nil, err
	}
	return todos, nil
}
