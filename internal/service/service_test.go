package service

import (
	"errors"
	"reflect"
	"testing"

	"github.com/TechOutsiders/TaskCLI/internal/model"
	"github.com/google/uuid"
)

// mockRepository is a test implementation of the repository interface.
type mockRepository struct {
	tasks        []model.Task
	err          error
	saveTask     *model.Task
	deleteTaskID uuid.UUID
}

// GetTask returns a task.
func (m *mockRepository) GetTask(id uuid.UUID) (*model.Task, error) {
	if m.err != nil {
		return nil, m.err
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

	repositoryErr := errors.New("repository error")

	tests := []struct {
		name              string
		repo              *mockRepository
		id                uuid.UUID
		wantErr           bool
		wantRepositoryErr error
	}{
		{
			name: "success",
			repo: &mockRepository{
				tasks: []model.Task{
					expectedTask,
				},
			},
			id: taskID,
		},
		{
			name:    "not found",
			repo:    &mockRepository{},
			id:      taskID,
			wantErr: true,
		},
		{
			name: "repository error",
			repo: &mockRepository{
				err: repositoryErr,
			},
			id:                taskID,
			wantErr:           true,
			wantRepositoryErr: repositoryErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.repo)

			task, err := s.GetTask(tt.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if tt.wantRepositoryErr != nil {
				if !errors.Is(err, tt.wantRepositoryErr) {
					t.Errorf("expected repository error %v, got %v", tt.wantRepositoryErr, err)
				}
			}

			if tt.wantErr {
				return
			}

			if task == nil {
				t.Fatal("got nil")
			}

			if !reflect.DeepEqual(*task, expectedTask) {
				t.Errorf("expected task %+v, got %+v", expectedTask, *task)
			}
		})
	}
}

// CreateTask saves the task received from the service.
func (m *mockRepository) CreateTask(task *model.Task) error {
	if m.err != nil {
		return m.err
	}

	m.saveTask = task
	return nil
}

func TestCreatTask(t *testing.T) {
	repositoryErr := errors.New("repository error")

	tests := []struct {
		name              string
		repo              *mockRepository
		data              *CreateTaskData
		wantErr           bool
		wantRepositoryErr error
	}{
		{
			name: "success",
			repo: &mockRepository{},
			data: &CreateTaskData{
				Title:       "Learning Spanish",
				Description: "First lesson: August 20, 2026",
				Priority:    model.PriorityMedium,
			},
		},
		{
			name: "repository error",
			repo: &mockRepository{
				err: repositoryErr,
			},
			data: &CreateTaskData{
				Title:       "Learning Spanish",
				Description: "First lesson",
				Priority:    model.PriorityMedium,
			},
			wantErr:           true,
			wantRepositoryErr: repositoryErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.repo)

			task, err := s.CreateTask(tt.data)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if tt.wantRepositoryErr != nil {
				if !errors.Is(err, tt.wantRepositoryErr) {
					t.Errorf("expected repository error %v, got %v", tt.wantRepositoryErr, err)
				}
			}

			if tt.wantErr {
				return
			}

			if task == nil {
				t.Fatal("got nil")
			}

			if task.ID == uuid.Nil {
				t.Error("unexpected ID")
			}

			if task.Title != tt.data.Title {
				t.Errorf("expected title %q, got %q", tt.data.Title, task.Title)
			}

			if task.Description != tt.data.Description {
				t.Errorf("expected title %q, got %q", tt.data.Description, task.Description)
			}

			if task.Priority != tt.data.Priority {
				t.Errorf("expected priority %q, got %q", tt.data.Priority, task.Priority)
			}

			if task.Status != model.StatusInProgress {
				t.Errorf("unexpected status: %v", task.Status)
			}

			if task.CreatedAt.IsZero() {
				t.Error("unexpected CreatedAt")
			}
		})
	}
}

// DeleteTask stores the ID of the deleted task.
func (m *mockRepository) DeleteTask(id uuid.UUID) error {
	m.deleteTaskID = id

	return m.err
}

func TestDeleteTask(t *testing.T) {
	repositoryErr := errors.New("repository error")
	taskID := uuid.New()

	tests := []struct {
		name              string
		repo              *mockRepository
		wantErr           bool
		wantRepositoryErr error
	}{
		{
			name: "success",
			repo: &mockRepository{},
		},
		{
			name: "repository error",
			repo: &mockRepository{
				err: repositoryErr,
			},
			wantErr:           true,
			wantRepositoryErr: repositoryErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.repo)

			err := s.DeleteTask(taskID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if tt.wantRepositoryErr != nil {
				if !errors.Is(err, tt.wantRepositoryErr) {
					t.Errorf("expected repository error %v, got %v", tt.wantRepositoryErr, err)
				}
			}

			if tt.wantErr {
				return
			}

			if tt.repo.deleteTaskID != taskID {
				t.Errorf("expected task ID %v, got %v", taskID, tt.repo.deleteTaskID)
			}
		})
	}
}

func (m *mockRepository) UpdateTask(task *model.Task) error {
	return nil
}

func (m *mockRepository) GetTasks() ([]model.Task, error) {
	return m.tasks, nil
}
