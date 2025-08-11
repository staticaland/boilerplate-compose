package processor

import (
	"fmt"
	"log"
	"strings"
	"time"

	"boilerplate-compose/executor"
)

type Orchestrator struct {
	processor *TemplateProcessor
	executor  *executor.CliExecutor
	dryRun    bool
}

func NewOrchestrator(processor *TemplateProcessor, exec *executor.CliExecutor, dryRun bool) *Orchestrator {
	return &Orchestrator{
		processor: processor,
		executor:  exec,
		dryRun:    dryRun,
	}
}

func (o *Orchestrator) Process() error {
	jobs, err := o.processor.BuildProcessingJobs()
	if err != nil {
		return fmt.Errorf("failed to build processing jobs: %w", err)
	}

	if !o.dryRun {
		// Check if boilerplate CLI is available
		if err := o.executor.CheckBoilerplateAvailable(); err != nil {
			return fmt.Errorf("boilerplate CLI check failed: %w", err)
		}
	}

	log.Printf("Processing %d templates", len(jobs))

	summary := executor.NewExecutionSummary()
	startTime := time.Now()

	for _, job := range jobs {
		result := o.processJob(job)
		summary.AddResult(result)

		// Stop on first failure unless in dry-run mode
		if !result.Success && !o.dryRun {
			summary.Print()
			return fmt.Errorf("template processing failed, stopping execution")
		}
	}

	summary.TotalDuration = time.Since(startTime)
	summary.Print()

	if summary.FailureCount > 0 && !o.dryRun {
		return fmt.Errorf("%d template(s) failed", summary.FailureCount)
	}

	return nil
}

func (o *Orchestrator) processJob(job ProcessingJob) executor.ExecutionResult {
	startTime := time.Now()
	result := executor.ExecutionResult{
		TemplateName: job.Name,
		StartTime:    startTime,
	}

	log.Printf("Processing template: %s", job.Name)

	if o.dryRun {
		err := o.dryRunJob(job)
		result.Success = err == nil
		result.Error = err
	} else {
		err := o.executor.Execute(job.Args, job.Name)
		result.Success = err == nil
		result.Error = err
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	return result
}

func (o *Orchestrator) dryRunJob(job ProcessingJob) error {
	fmt.Printf("\n=== Template: %s ===\n", job.Name)
	fmt.Printf("Command that would be executed:\n")
	fmt.Printf("boilerplate %s\n", strings.Join(job.Args, " "))
	fmt.Printf("\nTemplate details:\n")
	fmt.Printf("  URL: %s\n", job.Template.TemplateURL)
	fmt.Printf("  Output: %s\n", job.Template.OutputFolder)

	if len(job.Template.Vars) > 0 {
		fmt.Printf("  Variables:\n")
		for k, v := range job.Template.Vars {
			fmt.Printf("    %s = %s\n", k, v)
		}
	}

	return nil
}