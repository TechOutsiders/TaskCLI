package cmd

import (
	"fmt"

	"github.com/TechOutsiders/TaskCLI/internal/repository"
	"github.com/TechOutsiders/TaskCLI/internal/service"
)

const storagePath = "./data/tasks.json"

// Run initializes application dependencies and processes the command-line arguments.
func Run(args []string) error {
	repository, err := repository.NewTaskStorage(storagePath)
	if err != nil {
		return fmt.Errorf("creating repository: %w", err)
	}

	service := service.NewService(repository)

	if len(args) < 2 {
		return fmt.Errorf("command is required")
	}

	switch args[1] {
	case "add":
		//TODO: implement add comand
	case "list":
		//TODO: implement list comand
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
		return fmt.Errorf("unknown command: %s", args[1])
	}

	_ = service

	return nil
}
