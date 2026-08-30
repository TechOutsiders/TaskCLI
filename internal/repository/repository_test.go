package repository

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/TechOutsiders/TaskCLI/internal/model"
	"github.com/google/uuid"
)

type mockStorage struct {
	tasks      []model.Task
	loadErr    error
	saveErr    error
	savedTasks []model.Task
}

func (m *mockStorage) Load() ([]model.Task, error) {
	if m.loadErr != nil {
		return nil, m.loadErr
	}

	return m.tasks, nil
}

func (m *mockStorage) Save(tasks []model.Task) error {
	if m.saveErr != nil {
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
		loadErr  error
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
			loadErr: errors.New("storage load error"),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage := &mockStorage{
				tasks:   tc.tasks,
				loadErr: tc.loadErr,
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
