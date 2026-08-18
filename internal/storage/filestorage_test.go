package storage

import (
	"os"
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

func TestFileStorage_Load(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.json")

	storage := &FileStorage{
		path: path,
	}

	firstID := uuid.New()
	secondID := uuid.New()

	testCases := []struct {
		name     string
		data     string
		expected []model.Task
	}{
		{
			name: "one task",
			data: `[
			  { 
			    "ID": "` + firstID.String() + `", 
				"Title": "Learn testing",
				"Description": "Learn how to write tests",
				"Status": "ToDo",
				"Priority": "HIGH",
				"CreatedAt": "2026-08-17T12:00:00Z"
			  } 
		    ]`,
			expected: []model.Task{
				{
					ID:          firstID,
					Title:       "Learn testing",
					Description: "Learn how to write tests",
					Status:      model.StatusToDo,
					Priority:    model.PriorityHigh,
					CreatedAt:   time.Date(2026, 8, 17, 12, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "multiple tasks",
			data: `[
			  { 
			    "ID": "` + firstID.String() + `", 
				"Title": "Learn Go",
				"Description": "Learn Go testing",
				"Status": "ToDo",
				"Priority": "HIGH",
				"CreatedAt": "2026-08-17T12:00:00Z"
			  },
			  {
			    "ID": "` + secondID.String() + `", 
				"Title": "Write tests",
				"Description": "Write unit tests",
				"Status": "In Progress",
				"Priority": "MEDIUM",
				"CreatedAt": "2026-08-17T13:00:00Z"
			  }
		    ]`,
			expected: []model.Task{
				{
					ID:          firstID,
					Title:       "Learn Go",
					Description: "Learn Go testing",
					Status:      model.StatusToDo,
					Priority:    model.PriorityHigh,
					CreatedAt:   time.Date(2026, 8, 17, 12, 0, 0, 0, time.UTC),
				},
				{
					ID:          secondID,
					Title:       "Write tests",
					Description: "Write unit tests",
					Status:      model.StatusInProgress,
					Priority:    model.PriorityMedium,
					CreatedAt:   time.Date(2026, 8, 17, 13, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name:     "empty task list",
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
