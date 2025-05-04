package utils

import (
    "fmt"
    "os"
    "os/exec"
    "strings"
	"bufio"
    "io"
	"context"
	"errors"
	"time"
)

// RunShell runs a shell command and returns its output or error
func RunShell(cmdName string, args ...string) (string, error) {
	// **Security Check:** Whitelist allowed commands
	allowedCommands := map[string]bool{
		"ls":   true,
		"grep": true,
		"cat":  true,
		// Add other explicitly allowed commands
	}
	if !allowedCommands[cmdName] {
		return "", errors.New("command not allowed: " + cmdName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Example timeout
	defer cancel()

	cmd := exec.CommandContext(ctx, cmdName, args...)
	outputBytes, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return "", errors.New("command execution timed out")
	}

	if err != nil {
		return strings.TrimSpace(string(outputBytes)), errors.New(err.Error() + ": " + strings.TrimSpace(string(outputBytes)))
	}

	return strings.TrimSpace(string(outputBytes)), nil
}

// EnsureDockerNetworks ensures the required Docker networks exist
func EnsureDockerNetworks(networks []string) {
    for _, net := range networks {
        _, err := RunShell("docker", "network", "inspect", net)
        if err != nil {
            fmt.Printf("🔧 Creating missing network: %s\n", net)
            if _, err := RunShell("docker", "network", "create", net); err != nil {
                fmt.Printf("❌ Failed to create network %s: %v\n", net, err)
                os.Exit(1)
            }
        }
    }
}

// IsContainerRunning checks if a Docker container is running
func IsContainerRunning(name string) bool {
    output, err := RunShell("docker", "ps", "--format", "{{.Names}}")
    if err != nil {
        return false
    }
    lines := strings.Split(output, "\n")
    for _, line := range lines {
        if line == name {
            return true
        }
    }
    return false
}

// CheckTraefikHealth checks if Traefik container is running and healthy
func CheckTraefikHealth(containerName string) {
    if !IsContainerRunning(containerName) {
        fmt.Printf("❌ Traefik container \"%s\" is not running.\n", containerName)
        fmt.Println("Showing last 20 log lines (if available):")
        logs, err := RunShell("docker", "logs", "--tail", "20", containerName)
        if err != nil {
            fmt.Println("⚠️ No logs found. Container may not exist.")
        } else {
            fmt.Println(logs)
        }
        os.Exit(1)
    }

    fmt.Println("⏳ Checking health of Traefik container...")

    status, err := RunShell("docker", "inspect", "--format={{.State.Health.Status}}", containerName)
    if err != nil {
        fmt.Println("⚠️ Could not inspect container health. Proceeding anyway.")
        return
    }

    status = strings.TrimSpace(status)

    switch status {
    case "healthy":
        fmt.Println("✅ Traefik is healthy.")
    case "none":
        fmt.Println("⚠️ Traefik does not have a health check defined. Proceeding anyway.")
    default:
        fmt.Printf("❌ Traefik container is not healthy (status: %s).\n", status)
        fmt.Println("🪵 Showing last 20 log lines:")
        logs, _ := RunShell("docker", "logs", "--tail", "20", containerName)
        fmt.Println(logs)
        os.Exit(1)
    }
}

func ProcessEnvTemplate(templateEnvContent string) map[string]string {
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
func CopyFile(src, dst string) error {
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