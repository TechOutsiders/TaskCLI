package cmd

import (
	"fmt"
	"os"

	"github.com/TechOutsiders/TaskCLI/internal/repository"
	"github.com/TechOutsiders/TaskCLI/internal/service"
	"github.com/TechOutsiders/TaskCLI/internal/storage"
)

// Run initializes application dependencies and processes the command-line arguments.
func Run() {
	const storagePath = "./data/tasks.json"

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
	service := service.NewService(repository)

	switch os.Args[1] {
	case "add":
		// TODO: implement add comand
	case "list":
		// TODO: implement list comand
	case "show":
		// TODO: implement show command
	case "delete":
		// TODO: implement delete command
	case "edit":
		// TODO: implement edit command
	case "status":
		// TODO: implement status command
	case "priority":
		// TODO: implement priority command
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s", os.Args[1])
		os.Exit(1)
	}

	// TODO: Use when the command handlers will be implemented.
	_ = service
}
