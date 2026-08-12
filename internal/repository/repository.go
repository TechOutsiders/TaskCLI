package repository

import (
	"fmt"

	"github.com/TechOutsiders/TaskCLI/internal/model"
	"github.com/TechOutsiders/TaskCLI/internal/service"
	"github.com/TechOutsiders/TaskCLI/internal/storage"
	"github.com/google/uuid"
)

// TaskStorage implements the [service.Repository] interface. It is
// based on simple file storage.
type TaskStorage struct {
	fileStorage *storage.FileStorage
}

// NewTaskStorage creates a new *TaskStorage instance.
func NewTaskStorage(path string) (*TaskStorage, error) {
	fs, err := storage.NewFileStorage(path)
	if err != nil {
		return nil, fmt.Errorf("creating file storage: %w", err)
	}

	ts := &TaskStorage{
		fileStorage: fs,
	}

	return ts, nil
}

var _ service.Repository = (*TaskStorage)(nil)

// GetTasks retrieves all tasks from the storage.
func (ts *TaskStorage) GetTasks() ([]model.Task, error) {
	tasks, err := ts.fileStorage.Load()
	if err != nil {
		return nil, fmt.Errorf("loading tasks: %w", err)
	}

	return tasks, nil
}

// GetTask retrieves a single task by its ID. It returns an error if
// the task is not found.
func (ts *TaskStorage) GetTask(id uuid.UUID) (*model.Task, error) {
	tasks, err := ts.fileStorage.Load()
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
func (ts *TaskStorage) CreateTask(task *model.Task) error {
	tasks, err := ts.fileStorage.Load()
	if err != nil {
		return fmt.Errorf("loading tasks: %w", err)
	}

	for _, t := range tasks {
		if t.ID == task.ID {
			return fmt.Errorf("task with ID %s already exists", task.ID)
		}
	}

	tasks = append(tasks, *task)

	err = ts.fileStorage.Save(tasks)
	if err != nil {
		return fmt.Errorf("saving tasks: %w", err)
	}

	return nil
}

// DeleteTask removes a task by its ID. It returns an error if
// the task is not found.
func (ts *TaskStorage) DeleteTask(id uuid.UUID) error {
	tasks, err := ts.fileStorage.Load()
	if err != nil {
		return fmt.Errorf("loading tasks: %w", err)
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)

			err = ts.fileStorage.Save(tasks)
			if err != nil {
				return fmt.Errorf("saving tasks: %w", err)
			}

			return nil
		}
	}

	return fmt.Errorf("task with ID %s not found", id)
}

// UpdateTask updates an existing task with the provided data.
func (ts *TaskStorage) UpdateTask(task *model.Task) error {
	tasks, err := ts.fileStorage.Load()
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

		err = ts.fileStorage.Save(tasks)
		if err != nil {
			return fmt.Errorf("saving tasks: %w", err)
		}

		return nil
	}

	return fmt.Errorf("task with ID %s not found", task.ID)
}
