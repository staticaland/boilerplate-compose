package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("valid config file", func(t *testing.T) {
		configContent := `
templates:
  frontend:
    template-url: "https://github.com/example/react-template"
    output-folder: "./frontend"
    vars:
      project_name: "my-app"
  backend:
    template-url: "https://github.com/example/go-api-template"
    output-folder: "./backend"
`
		tempFile := createTempConfigFile(t, configContent)
		defer os.Remove(tempFile)

		config, err := LoadConfig(tempFile)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if len(config.Templates) != 2 {
			t.Errorf("Expected 2 templates, got %d", len(config.Templates))
		}

		frontend := config.Templates["frontend"]
		if frontend.TemplateURL != "https://github.com/example/react-template" {
			t.Errorf("Expected template URL 'https://github.com/example/react-template', got '%s'", frontend.TemplateURL)
		}
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := LoadConfig("nonexistent-file.yaml")
		if err == nil {
			t.Fatal("Expected error for nonexistent file")
		}
		if !strings.Contains(err.Error(), "config file not found") {
			t.Errorf("Expected 'config file not found' error, got: %v", err)
		}
	})

	t.Run("invalid YAML", func(t *testing.T) {
		invalidYAML := `
templates:
  frontend:
    template-url: "https://example.com"
    output-folder: ./frontend
    invalid: [unclosed array
`
		tempFile := createTempConfigFile(t, invalidYAML)
		defer os.Remove(tempFile)

		_, err := LoadConfig(tempFile)
		if err == nil {
			t.Fatal("Expected error for invalid YAML")
		}
		if !strings.Contains(err.Error(), "failed to parse YAML") {
			t.Errorf("Expected 'failed to parse YAML' error, got: %v", err)
		}
	})

	t.Run("no templates defined", func(t *testing.T) {
		emptyConfig := `
include:
  - path: "other.yaml"
`
		tempFile := createTempConfigFile(t, emptyConfig)
		defer os.Remove(tempFile)

		_, err := LoadConfig(tempFile)
		if err == nil {
			t.Fatal("Expected validation error for no templates")
		}
		if !strings.Contains(err.Error(), "no templates defined") {
			t.Errorf("Expected 'no templates defined' error, got: %v", err)
		}
	})

	t.Run("missing template-url", func(t *testing.T) {
		invalidConfig := `
templates:
  frontend:
    output-folder: "./frontend"
`
		tempFile := createTempConfigFile(t, invalidConfig)
		defer os.Remove(tempFile)

		_, err := LoadConfig(tempFile)
		if err == nil {
			t.Fatal("Expected validation error for missing template-url")
		}
		if !strings.Contains(err.Error(), "template-url is required") {
			t.Errorf("Expected 'template-url is required' error, got: %v", err)
		}
	})

	t.Run("missing output-folder", func(t *testing.T) {
		invalidConfig := `
templates:
  frontend:
    template-url: "https://example.com"
`
		tempFile := createTempConfigFile(t, invalidConfig)
		defer os.Remove(tempFile)

		_, err := LoadConfig(tempFile)
		if err == nil {
			t.Fatal("Expected validation error for missing output-folder")
		}
		if !strings.Contains(err.Error(), "output-folder is required") {
			t.Errorf("Expected 'output-folder is required' error, got: %v", err)
		}
	})
}

func TestValidateConfig(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		config := &ComposeConfig{
			Templates: map[string]Template{
				"test": {
					TemplateURL:  "https://example.com",
					OutputFolder: "./test",
				},
			},
		}

		err := validateConfig(config)
		if err != nil {
			t.Errorf("Expected no error for valid config, got: %v", err)
		}
	})

	t.Run("empty templates", func(t *testing.T) {
		config := &ComposeConfig{
			Templates: map[string]Template{},
		}

		err := validateConfig(config)
		if err == nil {
			t.Fatal("Expected error for empty templates")
		}
		if !strings.Contains(err.Error(), "no templates defined") {
			t.Errorf("Expected 'no templates defined' error, got: %v", err)
		}
	})

	t.Run("missing required fields", func(t *testing.T) {
		config := &ComposeConfig{
			Templates: map[string]Template{
				"test": {
					TemplateURL: "https://example.com",
					// OutputFolder missing
				},
			},
		}

		err := validateConfig(config)
		if err == nil {
			t.Fatal("Expected validation error")
		}
		if !strings.Contains(err.Error(), "output-folder is required") {
			t.Errorf("Expected 'output-folder is required' error, got: %v", err)
		}
	})
}

func createTempConfigFile(t *testing.T, content string) string {
	t.Helper()
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test-config.yaml")
	
	if err := os.WriteFile(tempFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	
	return tempFile
}