package main

import (
	"fmt"
	"flag"
	"io"
	"log"
	"os"
	"github.com/FrancescoCorbosiero/go-docker-manager/shared"
	"github.com/FrancescoCorbosiero/go-docker-manager/internal"
	//"github.com/FrancescoCorbosiero/go-docker-manager/pkg/utils"
)

func main() {
	// Set up logging
	logFile, err := os.OpenFile("go-docker-manager.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))

	// Configuration
	config := shared.Configuration{
		TemplatesDir: "templates",
		ComposeDir:   "compose",
	}

	// Parse command-line arguments
	command := flag.String("command", "", "Command to execute (dock, list, logs, down, restart)")
	container := flag.String("container", "", "Container/module name")
	template := flag.String("template", "", "Template name to use")
	flag.Parse()

	// Execute the requested command
	switch *command {
	case "dock":
		if *container == "" || *template == "" {
			log.Fatal("Container name and template are required for dock command")
		}
		err := internal.DockContainer(config, *container, *template)
		if err != nil {
			log.Fatalf("Failed to dock container: %v", err)
		}
	case "list":
		err := internal.ListContainers()
		if err != nil {
			log.Fatalf("Failed to list containers: %v", err)
		}
	case "logs":
		if *container == "" {
			log.Fatal("Container name is required for logs command")
		}
		err := internal.ShowLogs(*container)
		if err != nil {
			log.Fatalf("Failed to show logs: %v", err)
		}
	case "down":
		if *container == "" {
			log.Fatal("Container name is required for down command")
		}
		err := internal.StopContainer(*container)
		if err != nil {
			log.Fatalf("Failed to stop container: %v", err)
		}
	case "restart":
		if *container == "" {
			log.Fatal("Container name is required for restart command")
		}
		err := internal.RestartContainer(*container)
		if err != nil {
			log.Fatalf("Failed to restart container: %v", err)
		}
	default:
		printHelp()
	}
}

// printHelp prints the help message
func printHelp() {
	fmt.Println("Docker Manager - Container orchestration tool")
	fmt.Println("Commands:")
	fmt.Println("  -command=dock -container=NAME -template=TEMPLATE  Create and start a new container")
	fmt.Println("  -command=list                                    List running containers")
	fmt.Println("  -command=logs -container=NAME                    Show logs for a container")
	fmt.Println("  -command=down -container=NAME                    Stop and remove a container")
	fmt.Println("  -command=restart -container=NAME                 Restart a container")
}