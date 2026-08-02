package model

// Task holds all the information about a single to-do item.
type Task struct {
	ID          int
	Title       string
	Description string
	Status      Status
	Priority    Priority
	Created     string
}
