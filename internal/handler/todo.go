package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zuu-development/fullstack-examination-2024/internal/errors"
	"github.com/zuu-development/fullstack-examination-2024/internal/model"
	"github.com/zuu-development/fullstack-examination-2024/internal/service"
)

// TodoHandler is the request handler for the todo endpoint.
type TodoHandler interface {
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	Find(c echo.Context) error
	FindAll(c echo.Context) error
}

type todoHandler struct {
	Handler
	service service.Todo
}

// NewTodo returns a new instance of the todo handler.
func NewTodo(s service.Todo) TodoHandler {
	return &todoHandler{service: s}
}

// CreateRequest is the request parameter for creating a new todo
type CreateRequest struct {
	Task     string `json:"task" validate:"required"`
	Priority string `json:"priority" validate:"required"` // New Filed for priority
}

// Create handles the creation of a new todo.
// @Summary	Create a new todo
// @Tags		todos
// @Accept		json
// @Produce	json
// @Param		todo	body		model.Todo	true	"Todo object to create"
// @Success	201	{object}	model.Todo
// @Failure	400	{object}	ResponseError
// @Failure	500	{object}	ResponseError
// @Router		/todos [post]
// Create handles the creation of a new todo.
func (h *todoHandler) Create(c echo.Context) error {
	var req CreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{
			Errors: []Error{{Code: errors.CodeBadRequest, Message: "Invalid input"}},
		})
	}

	createdTodo, err := h.service.Create(req.Task, req.Priority) // Pass priority here
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{
			Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}},
		})
	}
	return c.JSON(http.StatusCreated, createdTodo)
}

// UpdateRequest is the request parameter for updating a todo
type UpdateRequest struct {
	UpdateRequestBody
	UpdateRequestPath
}

// UpdateRequestBody is the request body for updating a todo
type UpdateRequestBody struct {
	Task     string       `json:"task,omitempty"`
	Status   model.Status `json:"status,omitempty"`
	Priority string       `json:"priority,omitempty"` // Make it a pointer to allow omitting
}

// UpdateRequestPath is the request parameter for updating a todo
type UpdateRequestPath struct {
	ID int `param:"id" validate:"required"`
}

// Update handles the update of an existing todo.
// @Summary	Update a todo
// @Tags		todos
// @Accept		json
// @Produce	json
// @Param		body	body		UpdateRequest	true	"body"
// @Param		id	path	int	true	"Todo ID"
// @Success	200	{object}	ResponseData{Data=model.Todo}
// @Failure	400	{object}	ResponseError
// @Failure	404	{object}	ResponseError
// @Failure	500	{object}	ResponseError
// @Router		/todos/:id [put]
func (h *todoHandler) Update(c echo.Context) error {
	var req UpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{
			Errors: []Error{{Code: errors.CodeBadRequest, Message: "Invalid input"}},
		})
	}

	// Get the Todo ID from the URL parameters
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{
			Errors: []Error{{Code: errors.CodeBadRequest, Message: "Invalid Todo ID"}},
		})
	}

	// Call the service to update the todo
	todo, err := h.service.Update(id, req.Task, req.Status, req.Priority) // Pass priority here
	if err != nil {
		// Handle any other errors
		return c.JSON(http.StatusInternalServerError, ResponseError{
			Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}},
		})
	}

	return c.JSON(http.StatusOK, ResponseData{Data: todo})
}

// DeleteRequest is the request parameter for deleting a todo
type DeleteRequest struct {
	ID int `param:"id" validate:"required"`
}

// Delete handles the deletion of a todo.
// @Summary	Delete a todo
// @Tags		todos
// @Param		id	path	int	true	"Todo ID"
// @Success	204
// @Failure	400	{object}	ResponseError
// @Failure	404	{object}	ResponseError
// @Failure	500	{object}	ResponseError
// @Router		/todos/:id [delete]
func (h *todoHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{
			Errors: []Error{{Code: errors.CodeBadRequest, Message: "Invalid ID"}},
		})
	}

	if err := h.service.Delete(id); err != nil {
		if err == model.ErrNotFound {
			return c.JSON(http.StatusNotFound, ResponseError{
				Errors: []Error{{Code: errors.CodeNotFound, Message: "Todo not found"}},
			})
		}
		return c.JSON(http.StatusInternalServerError, ResponseError{
			Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}},
		})
	}
	return c.NoContent(http.StatusNoContent)
}

// Find handles finding a todo by ID.
// @Summary	Find a todo
// @Tags		todos
// @Param		id	path	int	true	"Todo ID"
// @Success	200	{object}	ResponseData{Data=model.Todo}
// @Failure	400	{object}	ResponseError
// @Failure	404	{object}	ResponseError
// @Failure	500	{object}	ResponseError
// @Router		/todos/:id [get]
func (h *todoHandler) Find(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{
			Errors: []Error{{Code: errors.CodeBadRequest, Message: "Invalid ID"}},
		})
	}

	res, err := h.service.Find(id)
	if err != nil {
		if err == model.ErrNotFound {
			return c.JSON(http.StatusNotFound, ResponseError{
				Errors: []Error{{Code: errors.CodeNotFound, Message: "Todo not found"}},
			})
		}
		return c.JSON(http.StatusInternalServerError, ResponseError{
			Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}},
		})
	}
	return c.JSON(http.StatusOK, ResponseData{Data: res})
}

// FindAll handles finding all todos.
// @Summary	Find all todos
// @Tags		todos
// @Param		task	query	string	false	"Filter by task"
// @Param		status	query	string	false	"Filter by status"
// @Success	200	{object}	ResponseData{Data=[]model.Todo}
// @Failure	500	{object}	ResponseError
// @Router		/todos [get]
func (h *todoHandler) FindAll(c echo.Context) error {
	task := c.QueryParam("task")
	status := c.QueryParam("status")

	incompleteTasks, completedTasks, err := h.service.FindAll(task, status) // Get both lists
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{
			Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}},
		})
	}

	// Structure the response
	response := struct {
		IncompleteTasks []*model.Todo `json:"incomplete_tasks"`
		CompletedTasks  []*model.Todo `json:"completed_tasks"`
	}{
		IncompleteTasks: incompleteTasks,
		CompletedTasks:  completedTasks,
	}

	return c.JSON(http.StatusOK, response)
}
