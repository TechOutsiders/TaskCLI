package main

import (
	"log"
	"os"

	"github.com/TechOutsiders/TaskCLI/internal/cmd"
)

func main() {
	err := cmd.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
