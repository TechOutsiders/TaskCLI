package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/TechOutsiders/TaskCLI/internal/model"
	"github.com/google/uuid"
)

type saveTestCase struct {
	name  string
	tasks []model.Task
}

type loadTestCase struct {
	name     string
	data     string
	expected []model.Task
}

var (
	firstCreatedAt  = time.Date(2026, 8, 17, 12, 0, 0, 0, time.UTC)
	secondCreatedAt = time.Date(2026, 8, 17, 13, 0, 0, 0, time.UTC)
	firstTaskID     = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	secondTaskID    = uuid.MustParse("22222222-2222-2222-2222-222222222222")
)

func TestFileStorage_Save(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.json")

	storage, err := NewFileStorage(path)
	if err != nil {
		t.Fatalf("NewfileStorage() error: %v", err)
	}

	testCases := []saveTestCase{
		{
			name: "single task",
			tasks: []model.Task{
				{
					ID:          firstTaskID,
					Title:       "Testing",
					Description: "Learning how to write tests",
					Status:      model.StatusToDo,
					Priority:    model.PriorityHigh,
					CreatedAt:   firstCreatedAt,
				},
			},
		},
		{
			name: "two tasks",
			tasks: []model.Task{
				{
					ID:          firstTaskID,
					Title:       "Testing",
					Description: "Learning how to write tests",
					Status:      model.StatusToDo,
					Priority:    model.PriorityHigh,
					CreatedAt:   firstCreatedAt,
				},
				{
					ID:          secondTaskID,
					Title:       "Testing number 2",
					Description: "Learning how to write tests",
					Status:      model.StatusInProgress,
					Priority:    model.PriorityMedium,
					CreatedAt:   secondCreatedAt,
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
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.json")

	storage, err := NewFileStorage(path)
	if err != nil {
		t.Fatalf("NewfileStorage() error: %v", err)
	}

	testCases := []loadTestCase{
		{
			name: "single task",
			data: `[
			  { 
			    "ID": "` + firstTaskID.String() + `", 
				"Title": "Learn testing",
				"Description": "Learn how to write tests",
				"Status": "ToDo",
				"Priority": "HIGH",
				"CreatedAt": "2026-08-17T12:00:00Z"
			  } 
		    ]`,
			expected: []model.Task{
				{
					ID:          firstTaskID,
					Title:       "Learn testing",
					Description: "Learn how to write tests",
					Status:      model.StatusToDo,
					Priority:    model.PriorityHigh,
					CreatedAt:   firstCreatedAt,
				},
			},
		},
		{
			name: "multiple tasks",
			data: `[
			  { 
			    "ID": "` + firstTaskID.String() + `", 
				"Title": "Learn Go",
				"Description": "Learn Go testing",
				"Status": "ToDo",
				"Priority": "HIGH",
				"CreatedAt": "2026-08-17T12:00:00Z"
			  },
			  {
			    "ID": "` + secondTaskID.String() + `", 
				"Title": "Write tests",
				"Description": "Write unit tests",
				"Status": "In Progress",
				"Priority": "MEDIUM",
				"CreatedAt": "2026-08-17T13:00:00Z"
			  }
		    ]`,
			expected: []model.Task{
				{
					ID:          firstTaskID,
					Title:       "Learn Go",
					Description: "Learn Go testing",
					Status:      model.StatusToDo,
					Priority:    model.PriorityHigh,
					CreatedAt:   firstCreatedAt,
				},
				{
					ID:          secondTaskID,
					Title:       "Write tests",
					Description: "Write unit tests",
					Status:      model.StatusInProgress,
					Priority:    model.PriorityMedium,
					CreatedAt:   secondCreatedAt,
				},
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
