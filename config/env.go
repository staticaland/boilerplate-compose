package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

	if !filepath.IsAbs(envFilePath) {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
		envFilePath = filepath.Join(wd, envFilePath)
	}

	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		return fmt.Errorf("env file not found: %s", envFilePath)
	}

	file, err := os.Open(envFilePath)
	if err != nil {
		return fmt.Errorf("failed to open env file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE format
		if err := em.parseEnvLine(line, lineNum); err != nil {
			return fmt.Errorf("error parsing line %d in %s: %w", lineNum, envFilePath, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading env file: %w", err)
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

// parseEnvLine parses a single line from the .env file
func (em *EnvironmentManager) parseEnvLine(line string, lineNum int) error {
	// Basic KEY=VALUE parsing
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid format, expected KEY=VALUE")
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	// Validate key format (basic validation)
	if key == "" {
		return fmt.Errorf("empty key")
	}

	// Remove quotes from value if present
	if len(value) >= 2 {
		if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
			(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}
	}

	em.envVars[key] = value
	return nil
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