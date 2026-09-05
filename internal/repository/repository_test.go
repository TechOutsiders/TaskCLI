package repository_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/TechOutsiders/TaskCLI/internal/model"
	"github.com/TechOutsiders/TaskCLI/internal/repository"
	"github.com/google/uuid"
)

// mockStorage implements the [repository.Storage] interface for testing.
type mockStorage struct {
	err        error
	savedTasks []model.Task
	tasks      []model.Task
}

// Load implements the [repository.Storage] interface.
func (m *mockStorage) Load() ([]model.Task, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.tasks, nil
}

// Save implements the [repository.Storage] interface.
func (m *mockStorage) Save(tasks []model.Task) error {
	if m.err != nil {
		return m.err
	}

	m.savedTasks = tasks

	return nil
}

// Shared test data
//
// TODO: Move shared test data to a separate `cli_test` package.
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

	notFoundTask = model.Task{
		ID:    notFoundID,
		Title: "New task",
	}
)

// storageErr is used to simulate a storage error in repository tests.
var storageErr = errors.New("storage error")

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
			name:    "storage error",
			err:     storageErr,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage := &mockStorage{
				tasks: tc.tasks,
				err:   tc.err,
			}

			repository := repository.New(storage)

			got, err := repository.GetTasks()
			if (err != nil) != tc.wantErr {
				t.Errorf("GetTasks() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf(
					"got = %+v, want %+v",
					got,
					tc.expected,
				)
			}
		})
	}
}

func TestRepository_GetTask(t *testing.T) {
	testCases := []struct {
		err      error
		expected *model.Task
		name     string
		tasks    []model.Task
		id       uuid.UUID
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
			name:    "storage error",
			id:      firstTaskID,
			err:     storageErr,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage := &mockStorage{
				tasks: tc.tasks,
				err:   tc.err,
			}

			repository := repository.New(storage)

			got, err := repository.GetTask(tc.id)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetTask() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf(
					"got = %+v, want %+v",
					got,
					tc.expected,
				)
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
			name:    "storage error",
			task:    firstTask,
			err:     storageErr,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage := &mockStorage{
				tasks: tc.tasks,
				err:   tc.err,
			}

			repository := repository.New(storage)

			err := repository.CreateTask(&tc.task)
			if (err != nil) != tc.wantErr {
				t.Errorf("CreateTask() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !reflect.DeepEqual(storage.savedTasks, tc.expectedTasks) {
				t.Errorf(
					"got = %+v, want %+v",
					storage.savedTasks,
					tc.expectedTasks,
				)
			}
		})
	}
}

func TestRepository_DeleteTask(t *testing.T) {
	testCases := []struct {
		err           error
		name          string
		tasks         []model.Task
		expectedTasks []model.Task
		id            uuid.UUID
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
			name:    "storage error",
			id:      firstTaskID,
			err:     storageErr,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage := &mockStorage{
				tasks: tc.tasks,
				err:   tc.err,
			}

			repository := repository.New(storage)

			err := repository.DeleteTask(tc.id)
			if (err != nil) != tc.wantErr {
				t.Errorf("DeleteTask() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !reflect.DeepEqual(storage.savedTasks, tc.expectedTasks) {
				t.Errorf(
					"got = %+v, want %+v",
					storage.savedTasks,
					tc.expectedTasks,
				)
			}
		})
	}
}

func TestRepository_UpdateTask(t *testing.T) {
	updatedTask := model.Task{
		ID:          firstTask.ID,
		Title:       "Updated title",
		Description: "Updated description",
		Status:      model.StatusDone,
		Priority:    model.PriorityLow,
		CreatedAt:   firstCreatedAt,
	}

	testCases := []struct {
		name          string
		task          model.Task
		storageErr    error
		expectedTasks []model.Task
		wantErr       bool
	}{
		{
			name: "update task",
			task: updatedTask,
			expectedTasks: []model.Task{
				updatedTask,
				secondTask,
			},
		},
		{
			name:    "task not found",
			task:    notFoundTask,
			wantErr: true,
		},
		{
			name:       "storage error",
			task:       updatedTask,
			storageErr: storageErr,
			wantErr:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage := &mockStorage{
				tasks: []model.Task{firstTask, secondTask},
				err:   tc.storageErr,
			}

			repository := repository.New(storage)

			err := repository.UpdateTask(&tc.task)
			if (err != nil) != tc.wantErr {
				t.Errorf("UpdateTask() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !reflect.DeepEqual(storage.savedTasks, tc.expectedTasks) {
				t.Errorf(
					"got = %+v, want %+v",
					storage.savedTasks,
					tc.expectedTasks,
				)
			}
		})
	}
}
