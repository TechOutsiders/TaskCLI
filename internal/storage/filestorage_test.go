package storage_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/TechOutsiders/TaskCLI/internal/model"
	"github.com/TechOutsiders/TaskCLI/internal/storage"
	"github.com/google/uuid"
)

// Shared test data
var (
	firstCreatedAt  = time.Date(2026, 8, 17, 12, 0, 0, 0, time.UTC)
	secondCreatedAt = time.Date(2026, 8, 17, 13, 0, 0, 0, time.UTC)

	firstTask = model.Task{
		ID:          uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Title:       "Testing",
		Description: "Learning how to write tests",
		Status:      model.StatusToDo,
		Priority:    model.PriorityHigh,
		CreatedAt:   firstCreatedAt,
	}
	secondTask = model.Task{
		ID:          uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		Title:       "Testing number 2",
		Description: "Learning how to write tests",
		Status:      model.StatusInProgress,
		Priority:    model.PriorityMedium,
		CreatedAt:   secondCreatedAt,
	}
)

func TestFileStorage_Save(t *testing.T) {
	storage, path := newTestStorage(t)

	testCases := []struct {
		name  string
		tasks []model.Task
	}{
		{
			name: "single task",
			tasks: []model.Task{
				firstTask,
			},
		},
		{
			name: "multiple tasks",
			tasks: []model.Task{
				firstTask,
				secondTask,
			},
		},
		{
			name:  "no tasks",
			tasks: []model.Task{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := storage.Save(tc.tasks)
			if err != nil {
				t.Fatalf("Save() error: %v", err)
			}

			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("ReadFile() error: %v", err)
			}

			var got []model.Task
			err = json.Unmarshal(data, &got)
			if err != nil {
				t.Fatalf("Unmarshal() error: %v", err)
			}

			if !reflect.DeepEqual(got, tc.tasks) {
				t.Errorf("tasks = %+v, want = %+v", got, tc.tasks)
			}
		})
	}
}

func TestFileStorage_Load(t *testing.T) {
	storage, path := newTestStorage(t)

	testCases := []struct {
		name     string
		data     string
		expected []model.Task
	}{
		{
			name: "single task",
			data: `[
				{
					"ID": "` + firstTask.ID.String() + `",
					"Title": "Testing",
					"Description": "Learning how to write tests",
					"Status": "ToDo",
					"Priority": "HIGH",
					"CreatedAt": "2026-08-17T12:00:00Z"
				}
			]`,
			expected: []model.Task{
				firstTask,
			},
		},
		{
			name: "multiple tasks",
			data: `[
				{
					"ID": "` + firstTask.ID.String() + `",
					"Title": "Testing",
					"Description": "Learning how to write tests",
					"Status": "ToDo",
					"Priority": "HIGH",
					"CreatedAt": "2026-08-17T12:00:00Z"
				},
				{
					"ID": "` + secondTask.ID.String() + `",
					"Title": "Testing number 2",
					"Description": "Learning how to write tests",
					"Status": "In Progress",
					"Priority": "MEDIUM",
					"CreatedAt": "2026-08-17T13:00:00Z"
				}
			]`,
			expected: []model.Task{
				firstTask,
				secondTask,
			},
		},
		{
			name:     "no tasks",
			data:     `[]`,
			expected: []model.Task{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := os.WriteFile(path, []byte(tc.data), 0o644)
			if err != nil {
				t.Fatalf("WriteFile() error: %v", err)
			}

			got, err := storage.Load()
			if err != nil {
				t.Fatalf("Load() error: %v ", err)
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Load() = %+v, want %+v", got, tc.expected)
			}
		})
	}
}

// newTestStorage creates a FileStorage instance backed by a temporary file.
func newTestStorage(t *testing.T) (*storage.FileStorage, string) {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.json")

	storage, err := storage.NewFileStorage(path)
	if err != nil {
		t.Fatalf("NewFileStorage() error: %v", err)
	}

	return storage, path
}
