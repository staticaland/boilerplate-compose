package processor

import (
	"reflect"
	"testing"

	"boilerplate-compose/config"
)

func TestNewTemplateProcessor(t *testing.T) {
	cfg := &config.ComposeConfig{}
	tp := NewTemplateProcessor(cfg)

	if tp.config != cfg {
		t.Error("Expected config to be set")
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
				"--output-folder", "./output",
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
				"--output-folder", "./output",
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
				"--output-folder", "./output",
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
				"--output-folder", "./output",
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
				"--output-folder", "./output",
				"--missing-key-action", "zero",
				"--missing-config-action", "ignore",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := &TemplateProcessor{}
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

	tp := NewTemplateProcessor(cfg)
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