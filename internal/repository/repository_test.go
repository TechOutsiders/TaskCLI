package repository

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/TechOutsiders/TaskCLI/internal/model"
	"github.com/google/uuid"
)

// mockStorage implements the [Storage] interface for testing.
type mockStorage struct {
	tasks      []model.Task
	err        error
	savedTasks []model.Task
}

// Load implements the [Storage] interface.
func (m *mockStorage) Load() ([]model.Task, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.tasks, nil
}

// Save implements the [Storage] interface.
func (m *mockStorage) Save(tasks []model.Task) error {
	if m.err != nil {
		return nil
	}

	m.savedTasks = tasks

	return nil
}

// Shared test data
var (
	firstCreatedAt  = time.Date(2026, 8, 17, 12, 0, 0, 0, time.UTC)
	secondCreatedAt = time.Date(2026, 8, 17, 13, 0, 0, 0, time.UTC)

	firstTaskID  = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	secondTaskID = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	notFoundID   = uuid.MustParse("33333333-3333-3333-3333-333333333333")

	firstTask = model.Task{
		ID:          firstTaskID,
		Title:       "Testing",
		Description: "Learning how to write tests",
		Status:      model.StatusToDo,
		Priority:    model.PriorityHigh,
		CreatedAt:   firstCreatedAt,
	}

	secondTask = model.Task{
		ID:          secondTaskID,
		Title:       "Testing number 2",
		Description: "Learning how to write tests",
		Status:      model.StatusInProgress,
		Priority:    model.PriorityMedium,
		CreatedAt:   secondCreatedAt,
	}
)

func TestRepository_GetTasks(t *testing.T) {
	testCases := []struct {
		name     string
		tasks    []model.Task
		err      error
		expected []model.Task
		wantErr  bool
	}{
		{
			name: "multiple tasks",
			tasks: []model.Task{
				firstTask,
				secondTask,
			},
			expected: []model.Task{
				firstTask,
				secondTask,
			},
		},
		{
			name:     "empty task list",
			tasks:    []model.Task{},
			expected: []model.Task{},
		},
		{
			name:    "storage load error",
			err:     errors.New("storage load error"),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage := &mockStorage{
				tasks: tc.tasks,
				err:   tc.err,
			}

			repository := New(storage)

			got, err := repository.GetTasks()
			if tc.wantErr {
				if err == nil {
					t.Fatalf("GetTasks() error = nil, want error")
				}

				return
			}

			if err != nil {
				t.Fatalf("GetTasks() error: %v", err)
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("GetTasks() = %+v, want %+v", got, tc.expected)
			}

		})
	}
}

func TestRepository_GetTask(t *testing.T) {
	testCases := []struct {
		name     string
		id       uuid.UUID
		tasks    []model.Task
		err      error
		expected *model.Task
		wantErr  bool
	}{
		{
			name: "task found",
			id:   firstTaskID,
			tasks: []model.Task{
				firstTask,
				secondTask,
			},
			expected: &firstTask,
		},
		{
			name: "task not found",
			id:   notFoundID,
			tasks: []model.Task{
				firstTask,
				secondTask,
			},
			wantErr: true,
		},
		{
			name:    "storage load error",
			id:      firstTaskID,
			err:     errors.New("storage error"),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage := &mockStorage{
				tasks: tc.tasks,
				err:   tc.err,
			}

			repository := New(storage)

			got, err := repository.GetTask(tc.id)
			if tc.wantErr {
				if err == nil {
					t.Fatal("GetTask() error = nil, want error")
				}

				return
			}

			if err != nil {
				t.Fatalf("GetTask() error: %v", err)
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("GetTask() = %+v, want %+v", got, tc.expected)
			}
		})
	}
}

func TestRepository_CreateTask(t *testing.T) {
	testCases := []struct {
		name          string
		tasks         []model.Task
		task          model.Task
		err           error
		expectedTasks []model.Task
		wantErr       bool
	}{
		{
			name: "create task",
			tasks: []model.Task{
				firstTask,
			},
			task: secondTask,
			expectedTasks: []model.Task{
				firstTask,
				secondTask,
			},
		},
		{
			name: "duplicate task ID",
			tasks: []model.Task{
				firstTask,
			},
			task:    firstTask,
			wantErr: true,
		},
		{
			name:    "storage load error",
			task:    firstTask,
			err:     errors.New("storage error"),
			wantErr: true,
		},
		{
			name: "storage save error",
			tasks: []model.Task{
				firstTask,
			},
			task:    secondTask,
			err:     errors.New("storage error"),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage := &mockStorage{
				tasks: tc.tasks,
				err:   tc.err,
			}

			repository := New(storage)

			err := repository.CreateTask(&tc.task)
			if tc.wantErr {
				if err == nil {
					t.Fatal("CreateTask() error = nil, want error")
				}

				return
			}

			if err != nil {
				t.Fatalf("CreateTask() error: %v", err)
			}

			if !reflect.DeepEqual(storage.savedTasks, tc.expectedTasks) {
				t.Errorf(
					"savedTasks = %+v, want %+v",
					storage.savedTasks,
					tc.expectedTasks)
			}
		})
	}
}

func TestRepository_DeleteTask(t *testing.T) {
	testCases := []struct {
		name          string
		id            uuid.UUID
		tasks         []model.Task
		err           error
		expectedTasks []model.Task
		wantErr       bool
	}{
		{
			name: "delete task",
			id:   firstTaskID,
			tasks: []model.Task{
				firstTask,
				secondTask,
			},
			expectedTasks: []model.Task{
				secondTask,
			},
		},
		{
			name: "task not found",
			id:   notFoundID,
			tasks: []model.Task{
				firstTask,
				secondTask,
			},
			wantErr: true,
		},
		{
			name:    "storage load error",
			id:      firstTaskID,
			err:     errors.New("storage error"),
			wantErr: true,
		},
		{
			name: "storage save error",
			id:   firstTaskID,
			tasks: []model.Task{
				firstTask,
				secondTask,
			},
			err:     errors.New("storage error"),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage := &mockStorage{
				tasks: tc.tasks,
				err:   tc.err,
			}

			repository := New(storage)

			err := repository.DeleteTask(tc.id)

			if tc.wantErr {
				if err == nil {
					t.Fatal("DeleteTask() error = nil, want error")
				}

				return
			}

			if err != nil {
				t.Fatalf("DeleteTask() error: %v", err)
			}

			if !reflect.DeepEqual(storage.savedTasks, tc.expectedTasks) {
				t.Errorf(
					"savedTasks = %+v, want %+v",
					storage.savedTasks,
					tc.expectedTasks,
				)
			}
		})
	}
}

func TestRepository_UpdateTask(t *testing.T) {
	updatedTask := firstTask
	updatedTask.Title = "Updated title"
	updatedTask.Description = "Updated description"
	updatedTask.Status = model.StatusDone
	updatedTask.Priority = model.PriorityLow

	testCases := []struct {
		name          string
		tasks         []model.Task
		task          model.Task
		err           error
		expectedTasks []model.Task
		wantErr       bool
	}{
		{
			name: "update task",
			tasks: []model.Task{
				firstTask,
				secondTask,
			},
			task: updatedTask,
			expectedTasks: []model.Task{
				updatedTask,
				secondTask,
			},
		},
		{
			name: "task not found",
			tasks: []model.Task{
				firstTask,
				secondTask,
			},
			task: model.Task{
				ID:    notFoundID,
				Title: "New task",
			},
			wantErr: true,
		},
		{
			name:    "storage load error",
			task:    firstTask,
			err:     errors.New("storage error"),
			wantErr: true,
		},
		{
			name: "storage save error",
			tasks: []model.Task{
				firstTask,
				secondTask,
			},
			task:    updatedTask,
			err:     errors.New("storage error"),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage := &mockStorage{
				tasks: tc.tasks,
				err:   tc.err,
			}

			repository := New(storage)

			err := repository.UpdateTask(&tc.task)

			if tc.wantErr {
				if err == nil {
					t.Fatal("UpdateTask() error = nil, want error")
				}

				return
			}

			if err != nil {
				t.Fatalf("UpdateTask() error: %v", err)
			}

			if !reflect.DeepEqual(storage.savedTasks, tc.expectedTasks) {
				t.Errorf(
					"savedTasks = %+v, want %+v",
					storage.savedTasks,
					tc.expectedTasks,
				)
			}
		})
	}
}
