package model

// Status defines the possible states of a task.
type Status string

const (
	StatusToDo       Status = "ToDo"
	StatusInProgress Status = "In Progress"
	StatusDone       Status = "Done"
)
