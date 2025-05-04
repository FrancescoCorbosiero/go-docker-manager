package main

import (
	"encoding/json"
	"fmt"
//	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ContainerStatus represents the status info for a container
type ContainerStatus struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Image   string `json:"image"`
	Created string `json:"created"`
	Ports   string `json:"ports"`
}

// ModuleInfo represents info about a module
type ModuleInfo struct {
	Name      string                 `json:"name"`
	Template  string                 `json:"template"`
	EnvConfig map[string]string      `json:"env_config"`
	Services  map[string]interface{} `json:"services"`
	Status    string                 `json:"status"`
}

// ServerConfig holds the web server configuration
type ServerConfig struct {
	Port     string `json:"port"`
	BasePath string `json:"base_path"`
}

// APIServer implements a simple web server for container management
func startAPIServer(config Configuration) {
	// Load server configuration
	serverConfig := ServerConfig{
		Port:     "8080",
		BasePath: "/api",
	}

	// Create API endpoints
	http.HandleFunc(serverConfig.BasePath+"/modules", func(w http.ResponseWriter, r *http.Request) {
		modules, err := listModules(config)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list modules: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(modules)
	})

	http.HandleFunc(serverConfig.BasePath+"/templates", func(w http.ResponseWriter, r *http.Request) {
		templates, err := listTemplates(config)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list templates: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(templates)
	})

	http.HandleFunc(serverConfig.BasePath+"/dock", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var moduleConfig ModuleConfig
		err := json.NewDecoder(r.Body).Decode(&moduleConfig)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusBadRequest)
			return
		}

		err = dockContainer(config, moduleConfig.Name, moduleConfig.Template)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to dock container: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Container started successfully"})
	})

	// Start the web server
	log.Printf("Starting API server on port %s", serverConfig.Port)
	log.Fatal(http.ListenAndServe(":"+serverConfig.Port, nil))
}

// listModules returns information about all modules in the compose directory
func listModules(config Configuration) ([]ModuleInfo, error) {
	var modules []ModuleInfo

	// Read the compose directory
	composeEntries, err := os.ReadDir(config.ComposeDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read compose directory: %v", err)
	}

	// Process each entry
	for _, entry := range composeEntries {
		if !entry.IsDir() {
			continue
		}

		moduleName := entry.Name()
		moduleDir := filepath.Join(config.ComposeDir, moduleName)

		// Check for docker-compose.yml
		composeFile := filepath.Join(moduleDir, "docker-compose.yml")
		if _, err := os.Stat(composeFile); os.IsNotExist(err) {
			continue
		}

		// Get template name by checking the content (this is an approximation)
		templateName := "unknown"
		for _, templateEntry := range mustListTemplates(config) {
			// Compare docker-compose files to determine the template
			templateComposeFile := filepath.Join(config.TemplatesDir, templateEntry, "docker-compose.yml")
			if compareFiles(composeFile, templateComposeFile) {
				templateName = templateEntry
				break
			}
		}

		// Read .env file
		envConfig := make(map[string]string)
		envFile := filepath.Join(moduleDir, ".env")
		if envContent, err := os.ReadFile(envFile); err == nil {
			lines := strings.Split(string(envContent), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}

				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					envConfig[parts[0]] = parts[1]
				}
			}
		}

		// Add module info
		modules = append(modules, ModuleInfo{
			Name:      moduleName,
			Template:  templateName,
			EnvConfig: envConfig,
			Services:  nil, // This would need to parse the docker-compose.yml
			Status:    getModuleStatus(moduleName),
		})
	}

	return modules, nil
}

// listTemplates returns the list of available templates
func listTemplates(config Configuration) ([]string, error) {
	templates := []string{}

	// Read the templates directory
	entries, err := os.ReadDir(config.TemplatesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read templates directory: %v", err)
	}

	// Filter for directories containing docker-compose.yml
	for _, entry := range entries {
		if entry.IsDir() {
			templateDir := filepath.Join(config.TemplatesDir, entry.Name())
			composePath := filepath.Join(templateDir, "docker-compose.yml")
			
			if _, err := os.Stat(composePath); err == nil {
				templates = append(templates, entry.Name())
			}
		}
	}

	return templates, nil
}

// mustListTemplates is like listTemplates but panics on error
func mustListTemplates(config Configuration) []string {
	templates, err := listTemplates(config)
	if err != nil {
		log.Fatalf("Failed to list templates: %v", err)
	}
	return templates
}

// compareFiles returns true if the files have similar content
func compareFiles(file1, file2 string) bool {
	content1, err := os.ReadFile(file1)
	if err != nil {
		return false
	}

	content2, err := os.ReadFile(file2)
	if err != nil {
		return false
	}

	// This is a very simple comparison - in a real app you might want to 
	// implement a more sophisticated comparison
	return strings.Contains(string(content1), string(content2)) || 
	       strings.Contains(string(content2), string(content1))
}

// getModuleStatus checks if containers for this module are running
func getModuleStatus(moduleName string) string {
	// In a real implementation, you would use docker API or exec to check container status
	// This is a placeholder
	return "unknown"
}

// These functions would be added to the main application to extend its functionality
// for future web UI integration. They are not implemented in detail here but serve
// as a starting point for extending the application.

// backupModule creates a backup of a module
func backupModule(config Configuration, moduleName string) error {
	// Implementation would create a backup of the module configuration
	return nil
}

// restoreModule restores a module from backup
func restoreModule(config Configuration, moduleName, backupName string) error {
	// Implementation would restore a module from backup
	return nil
}

// updateEnvVar updates an environment variable in a module's .env file
func updateEnvVar(config Configuration, moduleName, key, value string) error {
	// Implementation would update a specific environment variable
	return nil
}

// getContainerLogs returns logs for a specific container
func getContainerLogs(moduleName, containerName string, tail int) (string, error) {
	// Implementation would fetch logs using docker API or exec
	return "", nil
}