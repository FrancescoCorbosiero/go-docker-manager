package shared

// Represents the application configuration
type Configuration struct {
	TemplatesDir string
	ComposeDir   string
}

// Holds the data needed to create a new module
type ModuleConfig struct {
	Name     string
	Template string
	EnvVars  map[string]string
}