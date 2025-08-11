package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig(filepath string) (*ComposeConfig, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", filepath)
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config ComposeConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func validateConfig(config *ComposeConfig) error {
	if len(config.Templates) == 0 {
		return fmt.Errorf("no templates defined")
	}

	for name, template := range config.Templates {
		if template.TemplateURL == "" {
			return fmt.Errorf("template '%s': template-url is required", name)
		}
		if template.OutputFolder == "" {
			return fmt.Errorf("template '%s': output-folder is required", name)
		}
	}

	return nil
}