package storage

import (
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/TechOutsiders/TaskCLI/internal/model"
	"github.com/google/uuid"
)

func TestFileStorageSave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.json")

	storage, err := NewFileStorage(path)
	if err != nil {
		t.Fatalf("NewfileStorage() error: %v", err)
	}

	createdAt := time.Now().UTC().Round(0)

	testCases := []struct {
		name  string
		tasks []model.Task
	}{
		{
			name: "one task",
			tasks: []model.Task{
				{
					ID:          uuid.New(),
					Title:       "Testing",
					Description: "Learning how to write tests",
					Status:      model.StatusToDo,
					Priority:    model.PriorityHigh,
					CreatedAt:   createdAt,
				},
			},
		},
		{
			name: "multiple task",
			tasks: []model.Task{
				{
					ID:          uuid.New(),
					Title:       "Testing",
					Description: "Learning how to write tests",
					Status:      model.StatusToDo,
					Priority:    model.PriorityHigh,
					CreatedAt:   createdAt,
				},
				{
					ID:          uuid.New(),
					Title:       "Testing number 2",
					Description: "Learning how to write tests",
					Status:      model.StatusInProgress,
					Priority:    model.PriorityMedium,
					CreatedAt:   createdAt.Add(time.Hour),
				},
				{
					ID:          uuid.New(),
					Title:       "Testing number 2",
					Description: "Learning how to write tests",
					Status:      model.StatusInProgress,
					Priority:    model.PriorityMedium,
					CreatedAt:   createdAt.Add(2 * time.Hour),
				},
			},
		},
		{
			name:  "empty task list",
			tasks: []model.Task{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err = storage.Save(tc.tasks)
			if err != nil {
				t.Fatalf("Save() error: %v", err)
			}

			got, err := storage.Load()
			if err != nil {
				t.Fatalf("Load() error: %v", err)
			}

			if !reflect.DeepEqual(got, tc.tasks) {
				t.Errorf("tasks = %+v, want = %+v", got, tc.tasks)
			}
		})
	}
}
