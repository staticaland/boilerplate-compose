package executor

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

type CliExecutor struct {
	boilerplatePath string
	verbose         bool
}

func NewCliExecutor(boilerplatePath string, verbose bool) *CliExecutor {
	return &CliExecutor{
		boilerplatePath: boilerplatePath,
		verbose:         verbose,
	}
}

func (e *CliExecutor) Execute(args []string, templateName string) error {
	if e.boilerplatePath == "" {
		e.boilerplatePath = "boilerplate" // Default to PATH lookup
	}

	cmd := exec.Command(e.boilerplatePath, args...)

	// Set up pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	log.Printf("Executing template '%s': %s %s", templateName, e.boilerplatePath, strings.Join(args, " "))

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start boilerplate command: %w", err)
	}

	// Stream output
	done := make(chan error, 2)

	go e.streamOutput(stdout, "STDOUT", templateName, done)
	go e.streamOutput(stderr, "STDERR", templateName, done)

	// Wait for streaming to complete
	for i := 0; i < 2; i++ {
		if err := <-done; err != nil {
			log.Printf("Warning: error streaming output for template '%s': %v", templateName, err)
		}
	}

	// Wait for command to complete
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("boilerplate command failed for template '%s': %w", templateName, err)
	}

	log.Printf("Template '%s' completed successfully", templateName)
	return nil
}

func (e *CliExecutor) streamOutput(reader io.Reader, prefix string, templateName string, done chan error) {
	defer func() { done <- nil }()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if e.verbose {
			log.Printf("[%s][%s] %s", templateName, prefix, line)
		} else if prefix == "STDERR" {
			// Always show errors
			log.Printf("[%s] ERROR: %s", templateName, line)
		}
	}

	if err := scanner.Err(); err != nil {
		done <- fmt.Errorf("error reading %s: %w", prefix, err)
		return
	}
}

func (e *CliExecutor) CheckBoilerplateAvailable() error {
	cmd := exec.Command(e.boilerplatePath, "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("boilerplate CLI not found at '%s': %w", e.boilerplatePath, err)
	}
	return nil
}