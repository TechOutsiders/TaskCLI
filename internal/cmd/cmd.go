package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/TechOutsiders/TaskCLI/internal/model"
	"github.com/TechOutsiders/TaskCLI/internal/repository"
	"github.com/TechOutsiders/TaskCLI/internal/service"
	"github.com/TechOutsiders/TaskCLI/internal/storage"
	"github.com/google/uuid"
)

// Run initializes application dependencies and processes the command-line arguments.
func Run() {
	const storagePath = "../data/tasks.json"

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "command is required")
		os.Exit(1)
	}

	fileStorage, err := storage.NewFileStorage(storagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "creating file storage: %v\n", err)
		os.Exit(1)
	}

	repository := repository.New(fileStorage)
	svc := service.NewService(repository)

	handlers := map[string]func(*service.Service) error{
		"add":      handleAdd,
		"list":     handleList,
		"show":     handleShow,
		"delete":   handleDelete,
		"edit":     handleEdit,
		"status":   handleStatus,
		"priority": handlePriority,
	}

	handler, ok := handlers[os.Args[1]]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}

	err = handler(svc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// handleAdd handles the add command.
func handleAdd(s *service.Service) error {
	if len(os.Args) != 5 {
		return errors.New("usage: add <title> <description> <priority>")
	}

	title := os.Args[2]
	description := os.Args[3]
	priority := os.Args[4]

	var taskPriority model.Priority

	switch priority {
	case "low":
		taskPriority = model.PriorityLow
	case "medium":
		taskPriority = model.PriorityMedium
	case "high":
		taskPriority = model.PriorityHigh
	default:
		return errors.New("unknown priority")
	}

	createTaskData := &service.CreateTaskData{
		Title:       title,
		Description: description,
		Priority:    taskPriority,
	}

	_, err := s.CreateTask(createTaskData)

	return err
}

// handleList handles the list command.
func handleList(s *service.Service) error {
	if len(os.Args) != 2 {
		return errors.New("usage: taskcli list")
	}

	tasks, err := s.GetTasks()
	if err != nil {
		return err
	}

	fmt.Printf("%-38s %-15s %-10s %s\n", "ID", "Status", "Priority", "Title")

	for _, task := range tasks {
		fmt.Printf("%-38v %-15v %-10v %v\n",
			task.ID, task.Status, task.Priority, task.Title)
	}

	return nil
}

// handleShow handles the show command.
func handleShow(s *service.Service) error {
	if len(os.Args) != 3 {
		return errors.New("usage: taskcli show <id>")
	}

	id, err := uuid.Parse(os.Args[2])
	if err != nil {
		return errors.New("invalid task id")
	}

	task, err := s.GetTask(id)
	if err != nil {
		return err
	}

	fmt.Printf("ID: %v \n Status: %v \n Priority: %v \n Title: %v \n Description: %v \n Created: %v \n", task.ID, task.Status, task.Priority, task.Title, task.Description, task.CreatedAt)

	return nil
}

// handleDelete handles the delete command.
func handleDelete(s *service.Service) error {
	if len(os.Args) != 3 {
		return errors.New("usage: taskcli delete <id>")
	}

	id, err := uuid.Parse(os.Args[2])
	if err != nil {
		return errors.New("invalid task id")
	}

	err = s.DeleteTask(id)
	if err != nil {
		return err
	}

	return nil
}

// handleEdit handles the edit command.
func handleEdit(s *service.Service) error {
	if len(os.Args) < 3 {
		return errors.New("usage: taskcli edit <id> [--title <title>] [--description <description>]")
	}

	id, err := uuid.Parse(os.Args[2])
	if err != nil {
		return errors.New("invalid task id")
	}

	fs := flag.NewFlagSet("edit", flag.ContinueOnError)

	title := fs.String("title", "", "task title")
	description := fs.String("description", "", "task description")

	if err := fs.Parse(os.Args[3:]); err != nil {
		return err
	}

	hasUpdates := false
	fs.Visit(func(_ *flag.Flag) {
		hasUpdates = true
	})

	if !hasUpdates {
		return errors.New("nothing to update")
	}

	data := &service.UpdateTaskData{
		ID:          id,
		Title:       *title,
		Description: *description,
	}

	return s.UpdateTask(data)
}

// handleStatus handles the status command.
func handleStatus(s *service.Service) error {
	if len(os.Args) != 4 {
		return errors.New("usage: taskcli status <id> <status>")
	}

	id, err := uuid.Parse(os.Args[2])
	if err != nil {
		return errors.New("invalid task id")
	}

	var status model.Status

	switch os.Args[3] {
	case "todo":
		status = model.StatusToDo
	case "in_progress":
		status = model.StatusInProgress
	case "done":
		status = model.StatusDone
	default:
		return errors.New("unknown status")
	}

	err = s.UpdateStatus(id, status)
	if err != nil {
		return err
	}

	return nil
}

// handlePriority handles the priority command.
func handlePriority(s *service.Service) error {
	if len(os.Args) != 4 {
		return errors.New("usage: taskcli priority <id> <priority>")
	}

	id, err := uuid.Parse(os.Args[2])
	if err != nil {
		return errors.New("invalid task id")
	}

	var priority model.Priority

	switch os.Args[3] {
	case "low":
		priority = model.PriorityLow
	case "medium":
		priority = model.PriorityMedium
	case "high":
		priority = model.PriorityHigh
	default:
		return errors.New("unknown priority")
	}

	err = s.UpdatePriority(id, priority)
	if err != nil {
		return err
	}

	return nil
}
