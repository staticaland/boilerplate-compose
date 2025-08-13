package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewEnvironmentManager(t *testing.T) {
	em := NewEnvironmentManager()
	if em == nil {
		t.Fatal("NewEnvironmentManager returned nil")
	}
	if em.envVars == nil {
		t.Fatal("envVars map not initialized")
	}
}

func TestLoadEnvironmentFromFile(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "env_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test .env file
	envFile := filepath.Join(tempDir, ".env")
	envContent := `# Test environment file
TAG=v1.5
PROJECT_NAME=test-project
AUTHOR=john-doe
# Empty line above

QUOTED_VALUE="quoted value"
SINGLE_QUOTED='single quoted'
NO_QUOTES=no quotes value
`
	err = os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write env file: %v", err)
	}

	// Test loading the file
	em := NewEnvironmentManager()
	err = em.LoadEnvironmentFromFile(envFile)
	if err != nil {
		t.Fatalf("LoadEnvironmentFromFile failed: %v", err)
	}

	// Verify variables were loaded correctly
	tests := []struct {
		key      string
		expected string
	}{
		{"TAG", "v1.5"},
		{"PROJECT_NAME", "test-project"},
		{"AUTHOR", "john-doe"},
		{"QUOTED_VALUE", "quoted value"},
		{"SINGLE_QUOTED", "single quoted"},
		{"NO_QUOTES", "no quotes value"},
	}

	for _, test := range tests {
		value, exists := em.GetVariable(test.key)
		if !exists {
			t.Errorf("Variable %s not found", test.key)
			continue
		}
		if value != test.expected {
			t.Errorf("Variable %s: expected %q, got %q", test.key, test.expected, value)
		}
	}
}

func TestLoadEnvironmentFromFileNotFound(t *testing.T) {
	em := NewEnvironmentManager()
	err := em.LoadEnvironmentFromFile("/nonexistent/file.env")
	if err == nil {
		t.Fatal("Expected error for nonexistent file")
	}
}

func TestLoadEnvironmentFromFileEmpty(t *testing.T) {
	em := NewEnvironmentManager()
	err := em.LoadEnvironmentFromFile("")
	if err != nil {
		t.Fatalf("Empty path should not cause error: %v", err)
	}
}

func TestLoadSystemEnvironment(t *testing.T) {
	// Set a test environment variable
	testKey := "TEST_ENV_VAR_FOR_TESTING"
	testValue := "test_value_123"
	err := os.Setenv(testKey, testValue)
	if err != nil {
		t.Fatalf("Failed to set test env var: %v", err)
	}
	defer os.Unsetenv(testKey)

	em := NewEnvironmentManager()
	em.LoadSystemEnvironment()

	value, exists := em.GetVariable(testKey)
	if !exists {
		t.Errorf("System environment variable %s not loaded", testKey)
	}
	if value != testValue {
		t.Errorf("System environment variable %s: expected %q, got %q", testKey, testValue, value)
	}
}

func TestParseEnvLine(t *testing.T) {
	em := NewEnvironmentManager()

	tests := []struct {
		line        string
		expectError bool
		key         string
		value       string
	}{
		{"KEY=value", false, "KEY", "value"},
		{"KEY=", false, "KEY", ""},
		{"KEY=value with spaces", false, "KEY", "value with spaces"},
		{"KEY=\"quoted value\"", false, "KEY", "quoted value"},
		{"KEY='single quoted'", false, "KEY", "single quoted"},
		{"  KEY  =  value  ", false, "KEY", "value"},
		{"=value", true, "", ""},
		{"KEY", true, "", ""},
		{"", true, "", ""},
	}

	for i, test := range tests {
		err := em.parseEnvLine(test.line, i+1)
		if test.expectError {
			if err == nil {
				t.Errorf("Line %d: expected error for %q", i+1, test.line)
			}
			continue
		}
		if err != nil {
			t.Errorf("Line %d: unexpected error for %q: %v", i+1, test.line, err)
			continue
		}

		value, exists := em.GetVariable(test.key)
		if !exists {
			t.Errorf("Line %d: variable %s not found after parsing %q", i+1, test.key, test.line)
			continue
		}
		if value != test.value {
			t.Errorf("Line %d: variable %s: expected %q, got %q", i+1, test.key, test.value, value)
		}
	}
}

func TestInterpolateString(t *testing.T) {
	em := NewEnvironmentManager()
	em.SetVariable("TAG", "v1.5")
	em.SetVariable("PROJECT", "my-project")
	em.SetVariable("EMPTY", "")

	tests := []struct {
		input    string
		expected string
	}{
		{"simple string", "simple string"},
		{"${TAG}", "v1.5"},
		{"version:${TAG}", "version:v1.5"},
		{"${PROJECT}-${TAG}", "my-project-v1.5"},
		{"${NONEXISTENT}", "${NONEXISTENT}"},
		{"${TAG}/${PROJECT}", "v1.5/my-project"},
		{"${EMPTY}", ""},
		{"pre-${TAG}-post", "pre-v1.5-post"},
		{"multiple: ${TAG} and ${PROJECT}", "multiple: v1.5 and my-project"},
	}

	for _, test := range tests {
		result := em.InterpolateString(test.input)
		if result != test.expected {
			t.Errorf("InterpolateString(%q): expected %q, got %q", test.input, test.expected, result)
		}
	}
}

func TestInterpolateMapValues(t *testing.T) {
	em := NewEnvironmentManager()
	em.SetVariable("TAG", "v1.5")
	em.SetVariable("PROJECT", "my-project")

	input := map[string]string{
		"version":     "${TAG}",
		"name":        "${PROJECT}",
		"description": "Project ${PROJECT} version ${TAG}",
		"static":      "no variables here",
	}

	expected := map[string]string{
		"version":     "v1.5",
		"name":        "my-project",
		"description": "Project my-project version v1.5",
		"static":      "no variables here",
	}

	result := em.InterpolateMapValues(input)

	for key, expectedValue := range expected {
		if result[key] != expectedValue {
			t.Errorf("InterpolateMapValues: key %s: expected %q, got %q", key, expectedValue, result[key])
		}
	}
}

func TestGetSetVariable(t *testing.T) {
	em := NewEnvironmentManager()

	// Test setting and getting
	em.SetVariable("TEST_KEY", "test_value")
	value, exists := em.GetVariable("TEST_KEY")
	if !exists {
		t.Error("Variable should exist after setting")
	}
	if value != "test_value" {
		t.Errorf("Expected %q, got %q", "test_value", value)
	}

	// Test nonexistent variable
	_, exists = em.GetVariable("NONEXISTENT")
	if exists {
		t.Error("Nonexistent variable should not exist")
	}
}

func TestGetAllVariables(t *testing.T) {
	em := NewEnvironmentManager()
	em.SetVariable("KEY1", "value1")
	em.SetVariable("KEY2", "value2")

	all := em.GetAllVariables()
	if len(all) != 2 {
		t.Errorf("Expected 2 variables, got %d", len(all))
	}
	if all["KEY1"] != "value1" {
		t.Errorf("KEY1: expected %q, got %q", "value1", all["KEY1"])
	}
	if all["KEY2"] != "value2" {
		t.Errorf("KEY2: expected %q, got %q", "value2", all["KEY2"])
	}

	// Verify it's a copy (modifying the returned map shouldn't affect the original)
	all["KEY1"] = "modified"
	originalValue, _ := em.GetVariable("KEY1")
	if originalValue != "value1" {
		t.Error("GetAllVariables should return a copy, not the original map")
	}
}