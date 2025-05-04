package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"bufio"
	"github.com/FrancescoCorbosiero/go-docker-manager/shared"
	dockerops "github.com/FrancescoCorbosiero/go-docker-manager/internal"
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
		err := dockerops.dockContainer(config, *container, *template)
		if err != nil {
			log.Fatalf("Failed to dock container: %v", err)
		}
	case "list":
		err := dockerops.listContainers()
		if err != nil {
			log.Fatalf("Failed to list containers: %v", err)
		}
	case "logs":
		if *container == "" {
			log.Fatal("Container name is required for logs command")
		}
		err := dockerops.showLogs(*container)
		if err != nil {
			log.Fatalf("Failed to show logs: %v", err)
		}
	case "down":
		if *container == "" {
			log.Fatal("Container name is required for down command")
		}
		err := dockerops.stopContainer(*container)
		if err != nil {
			log.Fatalf("Failed to stop container: %v", err)
		}
	case "restart":
		if *container == "" {
			log.Fatal("Container name is required for restart command")
		}
		err := dockerops.restartContainer(*container)
		if err != nil {
			log.Fatalf("Failed to restart container: %v", err)
		}
	default:
		printHelp()
	}
}

func processEnvTemplate(templateEnvContent string) map[string]string {
	moduleEnvVars := make(map[string]string)
	placeholders := make(map[string]bool)
	placeholderValues := make(map[string]string)

	scanner := bufio.NewScanner(strings.NewReader(templateEnvContent))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
		continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
		continue
		}

		key := parts[0]
		defaultValue := parts[1]

		if strings.HasPrefix(defaultValue, "<") && strings.HasSuffix(defaultValue, ">") {
		placeholder := strings.TrimPrefix(strings.TrimSuffix(defaultValue, ">"), "<")
		placeholders[placeholder] = true
		} else {
		moduleEnvVars[key] = defaultValue
		}
	}

	reader := bufio.NewReader(os.Stdin)
	for placeholder := range placeholders {
		fmt.Printf("Enter value for %s: ", placeholder)
		value, _ := reader.ReadString('\n')
		value = strings.TrimSpace(value)
		if value == "" {
		placeholderValues[placeholder] = "<" + placeholder + ">" // Keep default if no input
		} else {
		placeholderValues[placeholder] = value
		}
	}

	scanner = bufio.NewScanner(strings.NewReader(templateEnvContent))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
		continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
		continue
		}

		key := parts[0]
		defaultValue := parts[1]

		if strings.HasPrefix(defaultValue, "<") && strings.HasSuffix(defaultValue, ">") {
		placeholder := strings.TrimPrefix(strings.TrimSuffix(defaultValue, ">"), "<")
		moduleEnvVars[key] = placeholderValues[placeholder]
		}
	}

	return moduleEnvVars
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
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