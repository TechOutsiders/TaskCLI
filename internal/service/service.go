package service

import (
	"fmt"
	"time"

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

// Service provides the business logic for the entire application.
type Service struct {
	repository Repository
}

// NewService creates a new *Service instance.
func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

// GetTasks gets all tasks.
func (s *Service) GetTasks() ([]model.Task, error) {
	tasks, err := s.repository.GetTasks()
	if err != nil {
		return nil, fmt.Errorf("getting tasks: %w", err)
	}

	return tasks, nil
}

// GetTask gets a task by id.
func (s *Service) GetTask(id uuid.UUID) (task *model.Task, err error) {
	task, err = s.repository.GetTask(id)
	if err != nil {
		return nil, fmt.Errorf("getting task: %s: %w", id, err)
	}

	return task, nil
}

// CreateTaskData stores all data required for task creation.
type CreateTaskData struct {
	Title       string
	Description string
	Priority    model.Priority
}

// CreateTask creates a new task.
// It generates a task's id and creation time.
func (s *Service) CreateTask(data *CreateTaskData) (*model.Task, error) {
	task := &model.Task{
		ID:          uuid.New(),
		Title:       data.Title,
		Description: data.Description,
		Status:      model.StatusInProgress,
		Priority:    data.Priority,
		CreatedAt:   time.Now(),
	}

	err := s.repository.CreateTask(task)
	if err != nil {
		return nil, fmt.Errorf("creating task: %w", err)
	}

	return task, nil
}

// DeleteTask deletes task by id.
func (s *Service) DeleteTask(id uuid.UUID) (err error) {
	err = s.repository.DeleteTask(id)
	if err != nil {
		return fmt.Errorf("deleting task: %w", err)
	}

	return nil
}

// UpdateTaskData stores all data required for task updating.
type UpdateTaskData struct {
	ID          uuid.UUID
	Title       string
	Description string
}

// UpdateTask updates task's title and description.
// If the value of data.Title or data.Description is an empty string,
// the value of the corresponding field is not changed.
func (s *Service) UpdateTask(data *UpdateTaskData) error {
	task, err := s.repository.GetTask(data.ID)
	if err != nil {
		return fmt.Errorf("getting tasks: %w", err)
	}

	if data.Title != "" {
		task.Title = data.Title
	}

	if data.Description != "" {
		task.Description = data.Description
	}

	err = s.repository.UpdateTask(task)
	if err != nil {
		return fmt.Errorf("updating task: %w", err)
	}

	return nil
}
