package processor

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"boilerplate-compose/config"
	"boilerplate-compose/executor"
)

func TestNewOrchestrator(t *testing.T) {
	cfg := &config.ComposeConfig{}
	tp := NewTemplateProcessor(cfg, "/test/config.yaml")
	exec := executor.NewCliExecutor("", false)
	orch := NewOrchestrator(tp, exec, true)

	if orch.processor != tp {
		t.Error("Expected processor to be set")
	}

	if orch.executor != exec {
		t.Error("Expected executor to be set")
	}

	if !orch.dryRun {
		t.Error("Expected dryRun to be true")
	}
}

func TestOrchestrator_Process(t *testing.T) {
	cfg := &config.ComposeConfig{
		Templates: map[string]config.Template{
			"test-template": {
				TemplateURL:  "https://github.com/example/template",
				OutputFolder: "./output",
				Vars:         map[string]string{"key": "value"},
			},
		},
	}

	tp := NewTemplateProcessor(cfg, "/test/config.yaml")
	exec := executor.NewCliExecutor("", false)

	t.Run("dry run mode", func(t *testing.T) {
		// Capture stdout to verify dry run output
		var buf bytes.Buffer
		originalStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		orch := NewOrchestrator(tp, exec, true)
		err := orch.Process()

		w.Close()
		os.Stdout = originalStdout

		// Read the captured output
		go func() {
			bytes := make([]byte, 1024)
			n, _ := r.Read(bytes)
			buf.Write(bytes[:n])
		}()

		if err != nil {
			t.Fatalf("Process() error = %v", err)
		}
	})

	t.Run("normal mode with CLI not available", func(t *testing.T) {
		// Capture log output
		var logBuf bytes.Buffer
		log.SetOutput(&logBuf)
		defer log.SetOutput(os.Stderr)

		orch := NewOrchestrator(tp, exec, false)
		err := orch.Process()

		// Expect this to fail since boilerplate CLI is not installed
		if err == nil {
			t.Fatal("Expected error when boilerplate CLI is not available")
		}

		if !strings.Contains(err.Error(), "boilerplate CLI check failed") {
			t.Errorf("Expected CLI check failure, got: %v", err)
		}
	})
}

func TestOrchestrator_ProcessJob(t *testing.T) {
	cfg := &config.ComposeConfig{}
	tp := NewTemplateProcessor(cfg, "/test/config.yaml")
	exec := executor.NewCliExecutor("", false)

	job := ProcessingJob{
		Name: "test-job",
		Template: config.Template{
			TemplateURL:  "https://github.com/example/template",
			OutputFolder: "./output",
		},
		Args: []string{"--template-url", "https://github.com/example/template", "--output-folder", "./output"},
	}

	t.Run("dry run job", func(t *testing.T) {
		// Capture stdout to verify dry run output
		var buf bytes.Buffer
		originalStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		orch := NewOrchestrator(tp, exec, true)
		result := orch.processJob(job)
		
		if !result.Success {
			t.Fatalf("processJob() result error = %v", result.Error)
		}

		w.Close()
		os.Stdout = originalStdout

		// Read the captured output
		output := make([]byte, 1024)
		n, _ := r.Read(output)
		buf.Write(output[:n])

		outputStr := string(output[:n])
		if !strings.Contains(outputStr, "=== Template: test-job ===") {
			t.Error("Expected template header in dry run output")
		}
		if !strings.Contains(outputStr, "boilerplate --template-url") {
			t.Error("Expected boilerplate command in dry run output")
		}
	})

	t.Run("normal mode job with CLI not available", func(t *testing.T) {
		// Capture log output
		var logBuf bytes.Buffer
		log.SetOutput(&logBuf)
		defer log.SetOutput(os.Stderr)

		orch := NewOrchestrator(tp, exec, false)
		result := orch.processJob(job)

		if result.Success {
			t.Fatal("Expected job to fail when boilerplate CLI is not available")
		}

		if result.Error == nil {
			t.Error("Expected error when boilerplate CLI is not available")
		}

		logOutput := logBuf.String()
		if !strings.Contains(logOutput, "Processing template: test-job") {
			t.Error("Expected log message about processing template")
		}
	})
}