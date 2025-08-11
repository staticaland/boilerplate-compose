package processor

import (
	"reflect"
	"testing"
	"path/filepath"

	"boilerplate-compose/config"
)

func TestNewTemplateProcessor(t *testing.T) {
	cfg := &config.ComposeConfig{}
	configPath := "/path/to/config.yaml"
	tp := NewTemplateProcessor(cfg, configPath)

	if tp.config != cfg {
		t.Error("Expected config to be set")
	}
	if tp.configPath != configPath {
		t.Error("Expected configPath to be set")
	}
}

func TestBuildBoilerplateArgs(t *testing.T) {
	tests := []struct {
		name     string
		template config.Template
		expected []string
	}{
		{
			name: "basic template with vars",
			template: config.Template{
				TemplateURL:    "https://github.com/example/template",
				OutputFolder:   "./output",
				Vars:           map[string]string{"key1": "value1", "key2": "value2"},
				NonInteractive: true,
			},
			expected: []string{
				"--template-url", "https://github.com/example/template",
				"--output-folder", "/test/output",
				"--var", "key1=value1",
				"--var", "key2=value2",
				"--non-interactive",
			},
		},
		{
			name: "template with single var-file",
			template: config.Template{
				TemplateURL:  "https://github.com/example/template",
				OutputFolder: "./output",
				VarFile:      "vars.yaml",
			},
			expected: []string{
				"--template-url", "https://github.com/example/template",
				"--output-folder", "/test/output",
				"--var-file", "vars.yaml",
			},
		},
		{
			name: "template with multiple var-files",
			template: config.Template{
				TemplateURL:  "https://github.com/example/template",
				OutputFolder: "./output",
				VarFile:      []interface{}{"vars1.yaml", "vars2.yaml"},
			},
			expected: []string{
				"--template-url", "https://github.com/example/template",
				"--output-folder", "/test/output",
				"--var-file", "vars1.yaml",
				"--var-file", "vars2.yaml",
			},
		},
		{
			name: "template with all boolean flags",
			template: config.Template{
				TemplateURL:               "https://github.com/example/template",
				OutputFolder:              "./output",
				NonInteractive:            true,
				NoHooks:                   true,
				NoShell:                   true,
				DisableDependencyPrompt:   true,
			},
			expected: []string{
				"--template-url", "https://github.com/example/template",
				"--output-folder", "/test/output",
				"--non-interactive",
				"--no-hooks",
				"--no-shell",
				"--disable-dependency-prompt",
			},
		},
		{
			name: "template with action flags",
			template: config.Template{
				TemplateURL:           "https://github.com/example/template",
				OutputFolder:          "./output",
				MissingKeyAction:      "zero",
				MissingConfigAction:   "ignore",
			},
			expected: []string{
				"--template-url", "https://github.com/example/template",
				"--output-folder", "/test/output",
				"--missing-key-action", "zero",
				"--missing-config-action", "ignore",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := &TemplateProcessor{configPath: "/test/config.yaml"}
			args, err := tp.buildBoilerplateArgs(tt.template)
			if err != nil {
				t.Fatalf("buildBoilerplateArgs() error = %v", err)
			}

			// Check if all expected args are present (order may vary for vars)
			if !containsAllArgs(args, tt.expected) {
				t.Errorf("buildBoilerplateArgs() = %v, want %v", args, tt.expected)
			}
		})
	}
}

func TestBuildProcessingJobs(t *testing.T) {
	cfg := &config.ComposeConfig{
		Templates: map[string]config.Template{
			"frontend": {
				TemplateURL:    "https://github.com/example/react-template",
				OutputFolder:   "./frontend",
				Vars:           map[string]string{"project_name": "my-react-app"},
				NonInteractive: true,
			},
			"backend": {
				TemplateURL:  "https://github.com/example/go-template",
				OutputFolder: "./backend",
				NoHooks:      true,
			},
		},
	}

	tp := NewTemplateProcessor(cfg, "/test/config.yaml")
	jobs, err := tp.BuildProcessingJobs()

	if err != nil {
		t.Fatalf("BuildProcessingJobs() error = %v", err)
	}

	if len(jobs) != 2 {
		t.Fatalf("Expected 2 jobs, got %d", len(jobs))
	}

	// Check that both templates are processed
	jobNames := make(map[string]bool)
	for _, job := range jobs {
		jobNames[job.Name] = true
		if job.Template.TemplateURL == "" {
			t.Errorf("Job %s missing template URL", job.Name)
		}
		if len(job.Args) == 0 {
			t.Errorf("Job %s missing args", job.Name)
		}
	}

	if !jobNames["frontend"] || !jobNames["backend"] {
		t.Error("Expected both frontend and backend jobs")
	}
}

// Helper function to check if args contains all expected arguments
func containsAllArgs(args, expected []string) bool {
	argMap := make(map[string]int)
	for _, arg := range args {
		argMap[arg]++
	}

	expectedMap := make(map[string]int)
	for _, exp := range expected {
		expectedMap[exp]++
	}

	return reflect.DeepEqual(argMap, expectedMap)
}

func TestResolveOutputPath(t *testing.T) {
	tests := []struct {
		name         string
		configPath   string
		outputFolder string
		expected     string
	}{
		{
			name:         "relative path in same directory",
			configPath:   "/home/user/project/config.yaml",
			outputFolder: "./output",
			expected:     "/home/user/project/output",
		},
		{
			name:         "relative path in subdirectory",
			configPath:   "/home/user/project/config.yaml", 
			outputFolder: "./frontend/build",
			expected:     "/home/user/project/frontend/build",
		},
		{
			name:         "absolute path unchanged",
			configPath:   "/home/user/project/config.yaml",
			outputFolder: "/tmp/output",
			expected:     "/tmp/output",
		},
		{
			name:         "config in subdirectory with relative output",
			configPath:   "/home/user/project/configs/app.yaml",
			outputFolder: "../output", 
			expected:     "/home/user/project/output",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := &TemplateProcessor{configPath: tt.configPath}
			result := tp.resolveOutputPath(tt.outputFolder)
			
			// Clean paths to handle different OS path separators
			result = filepath.Clean(result)
			expected := filepath.Clean(tt.expected)
			
			if result != expected {
				t.Errorf("resolveOutputPath() = %v, want %v", result, expected)
			}
		})
	}
}