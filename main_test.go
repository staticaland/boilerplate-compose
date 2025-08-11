package main

import (
	"testing"
)

func TestFindConfigFile(t *testing.T) {
	tests := []struct {
		name      string
		specified string
		expected  string
	}{
		{"specified file", "custom.yaml", "custom.yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findConfigFile(tt.specified)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestFindConfigFile_DefaultFiles(t *testing.T) {
	// Test with existing project file
	t.Run("finds existing boilerplate-compose.yaml", func(t *testing.T) {
		result := findConfigFile("")
		expected := "boilerplate-compose.yaml"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})
}