package repository

import (
	"fmt"

	"github.com/TechOutsiders/TaskCLI/internal/model"
	"github.com/TechOutsiders/TaskCLI/internal/service"
	"github.com/google/uuid"
)

// Storage defines the interface for task storage.
type Storage interface {
	Save(tasks []model.Task) error
	Load() ([]model.Task, error)
}

// TasksRepository implements the [service.Repository] interface. It is
// based on simple file storage.
type TasksRepository struct {
	storage Storage
}

// New creates a new *TasksRepository instance.
func New(storage Storage) *TasksRepository {
	return &TasksRepository{
		storage: storage,
	}
}

var _ service.Repository = (*TasksRepository)(nil)

// GetTasks retrieves all tasks from the storage.
func (ts *TasksRepository) GetTasks() ([]model.Task, error) {
	tasks, err := ts.storage.Load()
	if err != nil {
		return nil, fmt.Errorf("loading tasks: %w", err)
	}

	return tasks, nil
}

// GetTask retrieves a single task by its ID. It returns an error if
// the task is not found.
func (ts *TasksRepository) GetTask(id uuid.UUID) (*model.Task, error) {
	tasks, err := ts.storage.Load()
	if err != nil {
		return nil, fmt.Errorf("loading tasks: %w", err)
	}

	for _, task := range tasks {
		if task.ID == id {
			return &task, nil
		}
	}

	return nil, fmt.Errorf("task with ID %s not found", id)
}

// CreateTask adds a new task to the storage.
func (ts *TasksRepository) CreateTask(task *model.Task) error {
	tasks, err := ts.storage.Load()
	if err != nil {
		return fmt.Errorf("loading tasks: %w", err)
	}

	for _, t := range tasks {
		if t.ID == task.ID {
			return fmt.Errorf("task with ID %s already exists", task.ID)
		}
	}

	tasks = append(tasks, *task)

	err = ts.storage.Save(tasks)
	if err != nil {
		return fmt.Errorf("saving tasks: %w", err)
	}

	return nil
}

// DeleteTask removes a task by its ID. It returns an error if
// the task is not found.
func (ts *TasksRepository) DeleteTask(id uuid.UUID) error {
	tasks, err := ts.storage.Load()
	if err != nil {
		return fmt.Errorf("loading tasks: %w", err)
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)

			err = ts.storage.Save(tasks)
			if err != nil {
				return fmt.Errorf("saving tasks: %w", err)
			}

			return nil
		}
	}

	return fmt.Errorf("task with ID %s not found", id)
}

// UpdateTask updates an existing task with the provided data.
func (ts *TasksRepository) UpdateTask(task *model.Task) error {
	tasks, err := ts.storage.Load()
	if err != nil {
		return fmt.Errorf("loading tasks: %w", err)
	}

	for i, t := range tasks {
		if t.ID != task.ID {
			continue
		}

		tasks[i].Title = task.Title
		tasks[i].Description = task.Description
		tasks[i].Status = task.Status
		tasks[i].Priority = task.Priority

		err = ts.storage.Save(tasks)
		if err != nil {
			return fmt.Errorf("saving tasks: %w", err)
		}

		return nil
	}

	return fmt.Errorf("task with ID %s not found", task.ID)
}
