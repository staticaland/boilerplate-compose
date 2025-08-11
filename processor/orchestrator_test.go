package processor

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"boilerplate-compose/config"
)

func TestNewOrchestrator(t *testing.T) {
	cfg := &config.ComposeConfig{}
	tp := NewTemplateProcessor(cfg)
	orch := NewOrchestrator(tp, true)

	if orch.processor != tp {
		t.Error("Expected processor to be set")
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

	tp := NewTemplateProcessor(cfg)

	t.Run("dry run mode", func(t *testing.T) {
		// Capture stdout to verify dry run output
		var buf bytes.Buffer
		originalStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		orch := NewOrchestrator(tp, true)
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

	t.Run("normal mode", func(t *testing.T) {
		// Capture log output
		var logBuf bytes.Buffer
		log.SetOutput(&logBuf)
		defer log.SetOutput(os.Stderr)

		orch := NewOrchestrator(tp, false)
		err := orch.Process()

		if err != nil {
			t.Fatalf("Process() error = %v", err)
		}

		logOutput := logBuf.String()
		if !strings.Contains(logOutput, "Processing 1 templates") {
			t.Error("Expected log message about processing templates")
		}
		if !strings.Contains(logOutput, "Processing template: test-template") {
			t.Error("Expected log message about processing specific template")
		}
	})
}

func TestOrchestrator_ProcessJob(t *testing.T) {
	cfg := &config.ComposeConfig{}
	tp := NewTemplateProcessor(cfg)

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

		orch := NewOrchestrator(tp, true)
		err := orch.processJob(job)

		w.Close()
		os.Stdout = originalStdout

		// Read the captured output
		output := make([]byte, 1024)
		n, _ := r.Read(output)
		buf.Write(output[:n])

		if err != nil {
			t.Fatalf("processJob() error = %v", err)
		}

		outputStr := string(output[:n])
		if !strings.Contains(outputStr, "=== Template: test-job ===") {
			t.Error("Expected template header in dry run output")
		}
		if !strings.Contains(outputStr, "boilerplate --template-url") {
			t.Error("Expected boilerplate command in dry run output")
		}
	})

	t.Run("normal mode job", func(t *testing.T) {
		// Capture log output
		var logBuf bytes.Buffer
		log.SetOutput(&logBuf)
		defer log.SetOutput(os.Stderr)

		orch := NewOrchestrator(tp, false)
		err := orch.processJob(job)

		if err != nil {
			t.Fatalf("processJob() error = %v", err)
		}

		logOutput := logBuf.String()
		if !strings.Contains(logOutput, "Processing template: test-job") {
			t.Error("Expected log message about processing template")
		}
		if !strings.Contains(logOutput, "Would execute: boilerplate") {
			t.Error("Expected log message about command execution")
		}
	})
}