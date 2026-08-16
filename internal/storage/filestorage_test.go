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
	t.Run("saves one task", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "tasks.json")

		storage, err := NewFileStorage(path)
		if err != nil {
			t.Fatalf("NewfileStorage() error: %v", err)
		}

		createdAt := time.Now().UTC().Round(0)

		expected := []model.Task{
			{
				ID:          uuid.New(),
				Title:       "Testing",
				Description: "Learning how to write tests",
				Status:      model.StatusToDo,
				Priority:    model.PriorityHigh,
				CreatedAt:   createdAt,
			},
		}

		err = storage.Save(expected)
		if err != nil {
			t.Fatalf("Save() error: %v", err)
		}

		got, err := storage.Load()
		if err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("tasks = %+v, want = %+v", got, expected)
		}
	})
}
