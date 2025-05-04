package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
//	"text/template"
)

// Configuration represents the application configuration
type Configuration struct {
	TemplatesDir string
	ComposeDir   string
}

// ModuleConfig holds the data needed to create a new module
type ModuleConfig struct {
	Name     string
	Template string
	EnvVars  map[string]string
}

func main() {
	// Set up logging
	logFile, err := os.OpenFile("docker-manager.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))

	// Configuration
	config := Configuration{
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
		err := dockContainer(config, *container, *template)
		if err != nil {
			log.Fatalf("Failed to dock container: %v", err)
		}
	case "list":
		err := listContainers()
		if err != nil {
			log.Fatalf("Failed to list containers: %v", err)
		}
	case "logs":
		if *container == "" {
			log.Fatal("Container name is required for logs command")
		}
		err := showLogs(*container)
		if err != nil {
			log.Fatalf("Failed to show logs: %v", err)
		}
	case "down":
		if *container == "" {
			log.Fatal("Container name is required for down command")
		}
		err := stopContainer(*container)
		if err != nil {
			log.Fatalf("Failed to stop container: %v", err)
		}
	case "restart":
		if *container == "" {
			log.Fatal("Container name is required for restart command")
		}
		err := restartContainer(*container)
		if err != nil {
			log.Fatalf("Failed to restart container: %v", err)
		}
	default:
		printHelp()
	}
}

// dockContainer creates a new module from a template and runs it
func dockContainer(config Configuration, containerName, templateName string) error {
	log.Printf("Docking container %s using template %s", containerName, templateName)

	// Check if template exists
	templateDir := filepath.Join(config.TemplatesDir, templateName)
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return fmt.Errorf("template %s does not exist", templateName)
	}

	// Create module directory if it doesn't exist
	moduleDir := filepath.Join(config.ComposeDir, containerName)
	if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
		err = os.MkdirAll(moduleDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create module directory: %v", err)
		}
	}

	// Copy docker-compose.yml from template to module
	err := copyFile(
		filepath.Join(templateDir, "docker-compose.yml"),
		filepath.Join(moduleDir, "docker-compose.yml"),
	)
	if err != nil {
		return fmt.Errorf("failed to copy docker-compose.yml: %v", err)
	}

	// Read template .env file
	templateEnvPath := filepath.Join(templateDir, ".env")
	templateEnvContent, err := os.ReadFile(templateEnvPath)
	if err != nil {
		return fmt.Errorf("failed to read template .env file: %v", err)
	}

	// Process the .env template with user input for variable values
	moduleEnvVars := make(map[string]string)
	lines := strings.Split(string(templateEnvContent), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		defaultValue := parts[1]

		// If the key is surrounded by <>, it's a placeholder that needs user input
		if strings.HasPrefix(defaultValue, "<") && strings.HasSuffix(defaultValue, ">") {
			fmt.Printf("Enter value for %s [%s]: ", key, defaultValue)
			var value string
			fmt.Scanln(&value)
			if value == "" {
				value = defaultValue
			}
			moduleEnvVars[key] = value
		} else {
			moduleEnvVars[key] = defaultValue
		}
	}

	// Create .env file with user-provided values
	moduleEnvPath := filepath.Join(moduleDir, ".env")
	moduleEnvFile, err := os.Create(moduleEnvPath)
	if err != nil {
		return fmt.Errorf("failed to create module .env file: %v", err)
	}
	defer moduleEnvFile.Close()

	for key, value := range moduleEnvVars {
		_, err := moduleEnvFile.WriteString(fmt.Sprintf("%s=%s\n", key, value))
		if err != nil {
			return fmt.Errorf("failed to write to module .env file: %v", err)
		}
	}

	// Run the container with docker-compose
	cmd := exec.Command("docker", "compose", "-p", containerName, "up", "-d")
	cmd.Dir = moduleDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start container: %v, output: %s", err, output)
	}

	log.Printf("Container %s started successfully", containerName)
	fmt.Printf("Container %s started successfully\n", containerName)
	return nil
}

// listContainers lists all running docker containers
func listContainers() error {
	cmd := exec.Command("docker", "ps")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// showLogs shows logs for a specific container
func showLogs(containerName string) error {
	// Find the directory for the container
	moduleDir := filepath.Join("compose", containerName)
	if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
		return fmt.Errorf("module directory for %s does not exist", containerName)
	}

	cmd := exec.Command("docker", "compose", "-p", containerName, "logs", "-f")
	cmd.Dir = moduleDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// stopContainer stops and removes a container
func stopContainer(containerName string) error {
	moduleDir := filepath.Join("compose", containerName)
	if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
		return fmt.Errorf("module directory for %s does not exist", containerName)
	}

	cmd := exec.Command("docker", "compose", "-p", containerName, "down")
	cmd.Dir = moduleDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// restartContainer restarts a container
func restartContainer(containerName string) error {
	moduleDir := filepath.Join("compose", containerName)
	if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
		return fmt.Errorf("module directory for %s does not exist", containerName)
	}

	// Stop first
	stopCmd := exec.Command("docker", "compose", "-p", containerName, "down")
	stopCmd.Dir = moduleDir
	_, err := stopCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop container: %v", err)
	}

	// Then start
	startCmd := exec.Command("docker", "compose", "-p", containerName, "up", "-d")
	startCmd.Dir = moduleDir
	_, err = startCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}

	log.Printf("Container %s restarted successfully", containerName)
	fmt.Printf("Container %s restarted successfully\n", containerName)
	return nil
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