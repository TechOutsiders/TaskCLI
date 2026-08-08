package storage

import (
	"encoding/json"
	"fmt"

	"github.com/TechOutsiders/TaskCLI/internal/model"
	"github.com/google/uuid"
)

type TaskStorage struct {
	fileStorage *FileStorage
}

func NewTaskStorage(path string) (*TaskStorage, error) {
	fs, err := NewFileStorage(path)
	if err != nil {
		return nil, fmt.Errorf("creating file storage: %w", err)
	}

	ts := &TaskStorage{
		fileStorage: fs,
	}

	if err := ts.initializeFile(); err != nil {
		return nil, fmt.Errorf("initializing storage file: %w", err)
	}

	return ts, nil
}

// initializeFile проверяет и инициализирует файл пустым массивом при необходимости
func (ts *TaskStorage) initializeFile() error {
	data, err := ts.fileStorage.Load()
	if err != nil {
		return ts.saveTasks([]model.Task{})
	}

	if len(data) == 0 {
		return ts.saveTasks([]model.Task{})
	}

	var tasks []model.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return ts.saveTasks([]model.Task{})
	}

	return nil
}

// loadTasks загружает все задачи из файла
func (ts *TaskStorage) loadTasks() ([]model.Task, error) {
	data, err := ts.fileStorage.Load()
	if err != nil {
		return nil, fmt.Errorf("loading tasks: %w", err)
	}

	if len(data) == 0 {
		return []model.Task{}, nil
	}

	var tasks []model.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("unmarshaling tasks: %w", err)
	}

	return tasks, nil
}

func (ts *TaskStorage) saveTasks(tasks []model.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling tasks: %w", err)
	}

	if err := ts.fileStorage.Save(data); err != nil {
		return fmt.Errorf("saving tasks: %w", err)
	}

	return nil
}

// GetTasks retrieves all tasks.
func (ts *TaskStorage) GetTasks() ([]model.Task, error) {
	return ts.loadTasks()
}

// GetTask retrieves a single task by ID.
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

	return nil, fmt.Errorf("task with ID %s not found", id)
}

// CreateTask creates a new task.
func (ts *TaskStorage) CreateTask(task *model.Task) error {
	if task == nil {
		return fmt.Errorf("task cannot be nil")
	}

	if task.ID == uuid.Nil {
		task.ID = uuid.New()
	}

	tasks, err := ts.loadTasks()
	if err != nil {
		return err
	}

	for _, t := range tasks {
		if t.ID == task.ID {
			return fmt.Errorf("task with ID %s already exists", task.ID)
		}
	}

	tasks = append(tasks, *task)
	return ts.saveTasks(tasks)
}
