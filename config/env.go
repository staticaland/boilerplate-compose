package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

// EnvironmentManager handles environment variables and .env file parsing
type EnvironmentManager struct {
	envVars map[string]string
}

// NewEnvironmentManager creates a new environment manager
func NewEnvironmentManager() *EnvironmentManager {
	return &EnvironmentManager{
		envVars: make(map[string]string),
	}
}

// LoadEnvironmentFromFile loads variables from a .env file
func (em *EnvironmentManager) LoadEnvironmentFromFile(envFilePath string) error {
	if envFilePath == "" {
		return nil
	}

	fileEnv, err := godotenv.Read(envFilePath)
	if err != nil {
		return fmt.Errorf("failed to load env file %s: %w", envFilePath, err)
	}

	for key, value := range fileEnv {
		em.envVars[key] = value
	}

	return nil
}

// LoadSystemEnvironment loads variables from the system environment
func (em *EnvironmentManager) LoadSystemEnvironment() {
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			em.envVars[parts[0]] = parts[1]
		}
	}
}


// InterpolateString performs variable interpolation on a string using ${VAR} syntax
func (em *EnvironmentManager) InterpolateString(input string) string {
	// Regex to match ${VAR} patterns
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	
	return re.ReplaceAllStringFunc(input, func(match string) string {
		// Extract variable name from ${VAR}
		varName := match[2 : len(match)-1] // Remove ${ and }
		
		// Look up the variable value
		if value, exists := em.envVars[varName]; exists {
			return value
		}
		
		// If variable not found, return the original match
		return match
	})
}

// InterpolateMapValues performs variable interpolation on all string values in a map
func (em *EnvironmentManager) InterpolateMapValues(m map[string]string) map[string]string {
	result := make(map[string]string)
	for key, value := range m {
		result[key] = em.InterpolateString(value)
	}
	return result
}

// GetVariable returns the value of an environment variable
func (em *EnvironmentManager) GetVariable(name string) (string, bool) {
	value, exists := em.envVars[name]
	return value, exists
}

// SetVariable sets an environment variable
func (em *EnvironmentManager) SetVariable(name, value string) {
	em.envVars[name] = value
}

// GetAllVariables returns a copy of all environment variables
func (em *EnvironmentManager) GetAllVariables() map[string]string {
	result := make(map[string]string)
	for k, v := range em.envVars {
		result[k] = v
	}
	return result
}