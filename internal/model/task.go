package model

import (
	"time"

	"github.com/google/uuid"
)

// Task holds all the information about a single to-do item.
type Task struct {
	ID          uuid.UUID
	Title       string
	Description string
	Status      Status
	Priority    Priority
	CreatedAt   time.Time
}
