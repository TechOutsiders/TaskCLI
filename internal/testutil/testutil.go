package testutil

import (
	"errors"
	"time"

	"github.com/TechOutsiders/TaskCLI/internal/model"
	"github.com/google/uuid"
)

// Shared test data
var (
	FirstCreatedAt  = time.Date(2026, 8, 17, 12, 0, 0, 0, time.UTC)
	SecondCreatedAt = time.Date(2026, 8, 17, 13, 0, 0, 0, time.UTC)

	FirstTaskID  = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	SecondTaskID = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	NotFoundID   = uuid.MustParse("33333333-3333-3333-3333-333333333333")

	FirstTask = model.Task{
		ID:          FirstTaskID,
		Title:       "Testing",
		Description: "Learning how to write tests",
		Status:      model.StatusToDo,
		Priority:    model.PriorityHigh,
		CreatedAt:   FirstCreatedAt,
	}

	SecondTask = model.Task{
		ID:          SecondTaskID,
		Title:       "Testing number 2",
		Description: "Learning how to write tests",
		Status:      model.StatusInProgress,
		Priority:    model.PriorityMedium,
		CreatedAt:   SecondCreatedAt,
	}

	NotFoundTask = model.Task{
		ID:    NotFoundID,
		Title: "New task",
	}
)

// StorageErr is used to simulate a storage error in repository tests.
var StorageErr = errors.New("storage error")
