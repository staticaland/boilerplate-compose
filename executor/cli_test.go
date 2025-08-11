package executor

import (
	"testing"
)

func TestNewCliExecutor(t *testing.T) {
	tests := []struct {
		name            string
		boilerplatePath string
		verbose         bool
	}{
		{"default path", "", false},
		{"custom path", "/usr/local/bin/boilerplate", true},
		{"verbose mode", "boilerplate", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewCliExecutor(tt.boilerplatePath, tt.verbose)
			
			if executor.boilerplatePath != tt.boilerplatePath {
				t.Errorf("expected boilerplatePath %q, got %q", tt.boilerplatePath, executor.boilerplatePath)
			}
			
			if executor.verbose != tt.verbose {
				t.Errorf("expected verbose %v, got %v", tt.verbose, executor.verbose)
			}
		})
	}
}

func TestCliExecutor_Execute_InvalidCommand(t *testing.T) {
	executor := NewCliExecutor("/nonexistent/boilerplate", false)
	
	err := executor.Execute([]string{"--help"}, "test-template")
	
	if err == nil {
		t.Error("expected error when executing nonexistent command, got nil")
	}
}

func TestCliExecutor_CheckBoilerplateAvailable_NotFound(t *testing.T) {
	executor := NewCliExecutor("/nonexistent/boilerplate", false)
	
	err := executor.CheckBoilerplateAvailable()
	
	if err == nil {
		t.Error("expected error when checking nonexistent boilerplate CLI, got nil")
	}
}

func TestCliExecutor_Execute_EmptyPath(t *testing.T) {
	executor := NewCliExecutor("", false)
	
	// This should set boilerplatePath to "boilerplate" as fallback
	err := executor.Execute([]string{"--help"}, "test-template")
	
	// We expect this to fail since boilerplate CLI is not installed in test env
	if err == nil {
		t.Error("expected error when boilerplate CLI not available, got nil")
	}
}