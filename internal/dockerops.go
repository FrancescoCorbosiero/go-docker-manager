package internal

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"github.com/FrancescoCorbosiero/go-docker-manager/shared"
	utils "github.com/FrancescoCorbosiero/go-docker-manager/pkg/utils"
)

// dockContainer creates a new module from a template and runs it
func DockContainer(config shared.Configuration, containerName, templateName string) error {
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
	err := utils.CopyFile(
		filepath.Join(templateDir, "docker-compose.yml"),
		filepath.Join(moduleDir, "docker-compose.yml"),
	)
	if err != nil {
		return fmt.Errorf("failed to copy docker-compose.yml: %v", err)
	}

	// Read template .env file
	templateEnvPath := filepath.Join(templateDir, ".env.template")
	templateEnvContent, err := os.ReadFile(templateEnvPath)
	if err != nil {
		return fmt.Errorf("failed to read template .env file: %v", err)
	}

	// Process the .env template with user input for variable values
	moduleEnvVars := utils.ProcessEnvTemplate(string(templateEnvContent))

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
func ListContainers() error {
	cmd := exec.Command("docker", "ps")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// showLogs shows logs for a specific container
func ShowLogs(containerName string) error {
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
func StopContainer(containerName string) error {
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
func RestartContainer(containerName string) error {
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
