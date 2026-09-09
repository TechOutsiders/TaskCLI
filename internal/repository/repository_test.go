package repository_test

import (
	"reflect"
	"testing"

	"github.com/TechOutsiders/TaskCLI/internal/model"
	"github.com/TechOutsiders/TaskCLI/internal/repository"
	"github.com/TechOutsiders/TaskCLI/internal/testutil"
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
				testutil.FirstTask,
				testutil.SecondTask,
			},
			expected: []model.Task{
				testutil.FirstTask,
				testutil.SecondTask,
			},
		},
		{
			name:     "empty task list",
			tasks:    []model.Task{},
			expected: []model.Task{},
		},
		{
			name:    "storage error",
			err:     testutil.StorageErr,
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
			name: "success found",
			id:   testutil.FirstTaskID,
			tasks: []model.Task{
				testutil.FirstTask,
				testutil.SecondTask,
			},
			expected: &testutil.FirstTask,
		},
		{
			name: "task not found",
			id:   testutil.NotFoundID,
			tasks: []model.Task{
				testutil.FirstTask,
				testutil.SecondTask,
			},
			wantErr: true,
		},
		{
			name:    "storage error",
			id:      testutil.FirstTaskID,
			err:     testutil.StorageErr,
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
			name: "success creating",
			tasks: []model.Task{
				testutil.FirstTask,
			},
			task: testutil.SecondTask,
			expectedTasks: []model.Task{
				testutil.FirstTask,
				testutil.SecondTask,
			},
		},
		{
			name: "duplicate task ID",
			tasks: []model.Task{
				testutil.FirstTask,
			},
			task:    testutil.FirstTask,
			wantErr: true,
		},
		{
			name:    "storage error",
			task:    testutil.FirstTask,
			err:     testutil.StorageErr,
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
			name: "success deleting",
			id:   testutil.FirstTaskID,
			tasks: []model.Task{
				testutil.FirstTask,
				testutil.SecondTask,
			},
			expectedTasks: []model.Task{
				testutil.SecondTask,
			},
		},
		{
			name: "task not found",
			id:   testutil.NotFoundID,
			tasks: []model.Task{
				testutil.FirstTask,
				testutil.SecondTask,
			},
			wantErr: true,
		},
		{
			name:    "storage error",
			id:      testutil.FirstTaskID,
			err:     testutil.StorageErr,
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
		ID:          testutil.FirstTaskID,
		Title:       "Updated title",
		Description: "Updated description",
		Status:      model.StatusDone,
		Priority:    model.PriorityLow,
		CreatedAt:   testutil.FirstCreatedAt,
	}

	testCases := []struct {
		name          string
		task          model.Task
		storageErr    error
		expectedTasks []model.Task
		wantErr       bool
	}{
		{
			name: "success updating",
			task: updatedTask,
			expectedTasks: []model.Task{
				updatedTask,
				testutil.SecondTask,
			},
		},
		{
			name:    "task not found",
			task:    testutil.NotFoundTask,
			wantErr: true,
		},
		{
			name:       "storage error",
			task:       updatedTask,
			storageErr: testutil.StorageErr,
			wantErr:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage := &mockStorage{
				tasks: []model.Task{testutil.FirstTask, testutil.SecondTask},
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
