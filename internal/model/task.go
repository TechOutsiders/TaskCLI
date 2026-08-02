package model

type Status string

const (
	ToDo       Status = "ToDo"
	InProgress Status = "In Progress"
	Done       Status = "Done"
)

type Priority string

const (
	Low    Priority = "LOW"
	Medium Priority = "MEDIUM"
	High   Priority = "HIGH"
)

type Task struct {
	ID          int
	Title       string
	Description string
	Status      Status
	Priority    Priority
	Created     string
}
