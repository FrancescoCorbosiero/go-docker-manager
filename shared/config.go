package shared

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