package service

import (
	"errors"
	"testing"

	"github.com/TechOutsiders/TaskCLI/internal/model"
	"github.com/google/uuid"
)

type mockRepository struct {
	tasks         []model.Task
	getTaskErr    error
	createTaskErr error
	deleteTaskErr error
	saveTask      *model.Task
	deleteTaskID  uuid.UUID
}

func (m *mockRepository) GetTask(id uuid.UUID) (*model.Task, error) {
	if m.getTaskErr != nil {
		return nil, m.getTaskErr
	}

	for i := range m.tasks {
		if m.tasks[i].ID == id {
			return &m.tasks[i], nil
		}
	}
	return nil, errors.New("task not found")
}

func TestGetTask(t *testing.T) {
	taskID := uuid.New()

	expectedTask := model.Task{
		ID:          taskID,
		Title:       "Learn English",
		Description: "Repeat Present Simple",
		Status:      model.StatusToDo,
		Priority:    model.PriorityMedium,
	}

	repo := &mockRepository{
		tasks: []model.Task{
			expectedTask,
		},
	}

	s := NewService(repo)

	task, err := s.GetTask(taskID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if task == nil {
		t.Fatal("got nil")
	}

	if task.ID != expectedTask.ID {
		t.Errorf(
			"expected ID %s, got %s",
			expectedTask.ID,
			task.ID,
		)
	}

	if task.Title != expectedTask.Title {
		t.Errorf(
			"expected title %q, got %q",
			expectedTask.Title,
			task.Title,
		)
	}

	if task.Description != expectedTask.Description {
		t.Errorf(
			"expected description %q, got %q",
			expectedTask.Description,
			task.Description,
		)
	}

	if task.Status != expectedTask.Status {
		t.Errorf(
			"expected status %q, got %q",
			expectedTask.Status,
			task.Status,
		)
	}

	if task.Priority != expectedTask.Priority {
		t.Errorf(
			"expected priority %q, got %q",
			expectedTask.Priority,
			task.Priority,
		)
	}
}

func TestGetTaskNotFound(t *testing.T) {
	repo := &mockRepository{}

	s := NewService(repo)

	taskID := uuid.New()

	task, err := s.GetTask(taskID)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if task != nil {
		t.Errorf("expected nil task, got %v", task)
	}
}

func TestGetTaskRepositoryError(t *testing.T) {
	repositoryErr := errors.New("repository error")

	repo := &mockRepository{
		getTaskErr: repositoryErr,
	}

	s := NewService(repo)

	task, err := s.GetTask(uuid.New())

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if task != nil {
		t.Errorf("expected nil task, got %v", task)
	}

	if !errors.Is(err, repositoryErr) {
		t.Errorf(
			"expected repository error, got %v",
			err,
		)
	}
}

func (m *mockRepository) GetTasks() ([]model.Task, error) {
	return m.tasks, nil
}

func (m *mockRepository) CreateTask(task *model.Task) error {
	if m.createTaskErr != nil {
		return m.createTaskErr
	}

	m.saveTask = task
	return nil
}

func TestCreatTask(t *testing.T) {
	repo := &mockRepository{}

	s := NewService(repo)

	data := &CreateTaskData{
		Title:       "Learning Spanish",
		Description: "First lesson: August 20, 2026",
		Priority:    model.PriorityMedium,
	}
	task, err := s.CreateTask(data)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if task == nil {
		t.Fatal("got nil")
	}

	if repo.saveTask == nil {
		t.Fatal("repository did not receive task")
	}

	if task.ID == uuid.Nil {
		t.Error("unexpected ID")
	}

	if task.Title != data.Title {
		t.Error("unexpected Title")
	}

	if task.Description != data.Description {
		t.Error("unexpected Description")
	}

	if task.Priority != data.Priority {
		t.Error("unexpected priority")
	}

	if task.Status != model.StatusInProgress {
		t.Error("unexpected status")
	}

	if task.CreatedAt.IsZero() {
		t.Error("unexpected CreatedAt")
	}
}

func TestCreateTaskRepositoryError(t *testing.T) {
	repositoryErr := errors.New("repository error")

	repo := &mockRepository{
		createTaskErr: repositoryErr,
	}

	s := NewService(repo)

	data := &CreateTaskData{
		Title:       "Learning Spanish",
		Description: "First lesson",
		Priority:    model.PriorityMedium,
	}

	task, err := s.CreateTask(data)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if task != nil {
		t.Errorf("expected nil task, got %v", task)
	}

	if !errors.Is(err, repositoryErr) {
		t.Errorf("expected repository error, got %v", err)
	}
}

func (m *mockRepository) DeleteTask(id uuid.UUID) error {
	m.deleteTaskID = id

	return m.deleteTaskErr
}

func TestDeleteTask(t *testing.T) {
	taskID := uuid.New()

	repo := &mockRepository{}

	s := NewService(repo)

	err := s.DeleteTask(taskID)

	if err != nil {
		t.Fatal("an error occurred during deletion ")
	}

	if taskID != repo.deleteTaskID {
		t.Fatal("error transmitting the Service ID to the Repository ")
	}
}

func TestDeleteTaskError(t *testing.T) {
	taskID := uuid.New()

	repositoryErr := errors.New("repository error")

	repo := &mockRepository{
		deleteTaskErr: repositoryErr,
	}

	s := NewService(repo)

	err := s.DeleteTask(taskID)

	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, repositoryErr) {
		t.Errorf("expected repository error, got %v", err)
	}
}

func (m *mockRepository) UpdateTask(task *model.Task) error {
	return nil
}
