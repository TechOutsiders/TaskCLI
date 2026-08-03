package service

import (
	"github.com/TechOutsiders/TaskCLI/internal/model"
	"github.com/google/uuid"
)

// Repository defines the contract for task persistence.
// It's responsible only for storing and retrieving tasks.
type Repository interface {
	// GetTasks retrieves all tasks.
	GetTasks() (tasks []model.Task, err error)

	// GetTask retrieves a single task by ID. If no task is found,
	// a nil task and a non-nil error are returned.
	GetTask(id uuid.UUID) (task *model.Task, err error)

	// CreateTask creates a new task.
	CreateTask(task *model.Task) (err error)

	// DeleteTask deletes a task by ID. If no task is found,
	// a non-nil error is returned.
	DeleteTask(id uuid.UUID) (err error)

	// UpdateTask updates an existing task with the provided data.
	// All fields of the task are persisted.
	UpdateTask(task *model.Task) (err error)
}
