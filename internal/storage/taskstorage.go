package storage

import (
	"encoding/json"
	"fmt"

	"github.com/TechOutsiders/TaskCLI/internal/model"
	"github.com/google/uuid"
)

// TaskStorage implements the service.Repository interface using FileStorage.
// It provides CRUD operations for tasks.
type TaskStorage struct {
	fileStorage *FileStorage
}

// NewTaskStorage creates a new TaskStorage instance.
// It ensures the storage file exists and is properly initialized.
func NewTaskStorage(path string) (*TaskStorage, error) {
	fs, err := NewFileStorage(path)
	if err != nil {
		return nil, fmt.Errorf("Creating file storage: %w", err)
	}

	ts := &TaskStorage{
		fileStorage: fs,
	}

	return ts, nil
}

// loadTasks reads all tasks from the storage file.
// Returns an empty slice if the file is empty.
func (ts *TaskStorage) loadTasks() ([]model.Task, error) {
	data, err := ts.fileStorage.Load()
	if err != nil {
		return nil, fmt.Errorf("Loading tasks: %w", err)
	}

	if len(data) == 0 {
		return []model.Task{}, nil
	}

	var tasks []model.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("Unmarshaling tasks: %w", err)
	}

	return tasks, nil
}

// saveTasks writes the provided tasks to the storage file.
// It marshals the tasks to JSON with indentation.
func (ts *TaskStorage) saveTasks(tasks []model.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("Marshaling tasks: %w", err)
	}

	if err := ts.fileStorage.Save(data); err != nil {
		return fmt.Errorf("Saving tasks: %w", err)
	}

	return nil
}

// GetTasks retrieves all tasks from the storage.
// Returns an empty slice if no tasks exist.
func (ts *TaskStorage) GetTasks() ([]model.Task, error) {
	return ts.loadTasks()
}

// GetTask retrieves a single task by its UUID.
// Returns an error if the task is not found.
func (ts *TaskStorage) GetTask(id uuid.UUID) (*model.Task, error) {
	tasks, err := ts.loadTasks()
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		if task.ID == id {
			return &task, nil
		}
	}

	return nil, fmt.Errorf("Task with ID %s not found", id)
}

// CreateTask adds a new task to the storage.
// Generates a new UUID if the task ID is not set.
func (ts *TaskStorage) CreateTask(task *model.Task) error {

	tasks, err := ts.loadTasks()
	if err != nil {
		return err
	}

	for _, t := range tasks {
		if t.ID == task.ID {
			return fmt.Errorf("Task with ID %s already exists", task.ID)
		}
	}

	tasks = append(tasks, *task)
	return ts.saveTasks(tasks)
}

// DeleteTask removes a task by its UUID.
// Returns an error if the task cannot be loaded or if the task is not found.
func (ts *TaskStorage) DeleteTask(id uuid.UUID) error {
	tasks, err := ts.loadTasks()
	if err != nil {
		return err
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return ts.saveTasks(tasks)
		}
	}
	return fmt.Errorf("Task with ID %s not found", id)
}

// UpdateTask updates an existing task with the provided data.
// All fields of the task are persisted.
func (ts *TaskStorage) UpdateTask(task *model.Task) error {

	tasks, err := ts.loadTasks()
	if err != nil {
		return err
	}

	for i, t := range tasks {

		if t.ID != task.ID {
			continue
		}

		tasks[i].Title = task.Title

		tasks[i].Description = task.Description

		tasks[i].Status = task.Status

		tasks[i].Priority = task.Priority

		return ts.saveTasks(tasks)
	}
	return fmt.Errorf("task with ID %s not found", task.ID)
}
