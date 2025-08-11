package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestComposeConfigUnmarshal(t *testing.T) {
	yamlData := `
templates:
  frontend:
    template-url: "https://github.com/example/react-template"
    output-folder: "./frontend"
    vars:
      project_name: "my-app"
      author: "john-doe"
    non-interactive: true
  backend:
    template-url: "https://github.com/example/go-api-template"
    output-folder: "./backend"
    var-file: "backend-vars.yaml"
    missing-key-action: "error"
include:
  - path: "other-compose.yaml"
extends:
  file: "base-compose.yaml"
  template: "base-template"
`

	var config ComposeConfig
	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	if len(config.Templates) != 2 {
		t.Errorf("Expected 2 templates, got %d", len(config.Templates))
	}

	frontend, exists := config.Templates["frontend"]
	if !exists {
		t.Fatal("Frontend template not found")
	}

	if frontend.TemplateURL != "https://github.com/example/react-template" {
		t.Errorf("Expected template URL 'https://github.com/example/react-template', got '%s'", frontend.TemplateURL)
	}

	if frontend.OutputFolder != "./frontend" {
		t.Errorf("Expected output folder './frontend', got '%s'", frontend.OutputFolder)
	}

	if !frontend.NonInteractive {
		t.Error("Expected non-interactive to be true")
	}

	if len(frontend.Vars) != 2 {
		t.Errorf("Expected 2 vars, got %d", len(frontend.Vars))
	}

	if frontend.Vars["project_name"] != "my-app" {
		t.Errorf("Expected var project_name to be 'my-app', got '%s'", frontend.Vars["project_name"])
	}

	backend, exists := config.Templates["backend"]
	if !exists {
		t.Fatal("Backend template not found")
	}

	if backend.MissingKeyAction != "error" {
		t.Errorf("Expected missing-key-action 'error', got '%s'", backend.MissingKeyAction)
	}

	if len(config.Include) != 1 {
		t.Errorf("Expected 1 include, got %d", len(config.Include))
	}

	if config.Include[0].Path != "other-compose.yaml" {
		t.Errorf("Expected include path 'other-compose.yaml', got '%s'", config.Include[0].Path)
	}

	if config.Extends == nil {
		t.Fatal("Expected extends config to be present")
	}

	if config.Extends.File != "base-compose.yaml" {
		t.Errorf("Expected extends file 'base-compose.yaml', got '%s'", config.Extends.File)
	}

	if config.Extends.Template != "base-template" {
		t.Errorf("Expected extends template 'base-template', got '%s'", config.Extends.Template)
	}
}

func TestTemplateVarFileTypes(t *testing.T) {
	t.Run("string var file", func(t *testing.T) {
		yamlData := `
templates:
  test:
    template-url: "https://example.com"
    output-folder: "./test"
    var-file: "vars.yaml"
`
		var config ComposeConfig
		err := yaml.Unmarshal([]byte(yamlData), &config)
		if err != nil {
			t.Fatalf("Failed to unmarshal YAML: %v", err)
		}

		template := config.Templates["test"]
		if template.VarFile != "vars.yaml" {
			t.Errorf("Expected var-file to be 'vars.yaml', got %v", template.VarFile)
		}
	})

	t.Run("array var file", func(t *testing.T) {
		yamlData := `
templates:
  test:
    template-url: "https://example.com"
    output-folder: "./test"
    var-file:
      - "vars1.yaml"
      - "vars2.yaml"
`
		var config ComposeConfig
		err := yaml.Unmarshal([]byte(yamlData), &config)
		if err != nil {
			t.Fatalf("Failed to unmarshal YAML: %v", err)
		}

		template := config.Templates["test"]
		varFiles, ok := template.VarFile.([]interface{})
		if !ok {
			t.Fatalf("Expected var-file to be an array, got %T", template.VarFile)
		}

		if len(varFiles) != 2 {
			t.Errorf("Expected 2 var files, got %d", len(varFiles))
		}
	})
}