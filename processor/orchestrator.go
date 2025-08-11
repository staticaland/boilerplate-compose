package processor

import (
	"fmt"
	"log"
	"strings"
)

type Orchestrator struct {
	processor *TemplateProcessor
	dryRun    bool
}

func NewOrchestrator(processor *TemplateProcessor, dryRun bool) *Orchestrator {
	return &Orchestrator{
		processor: processor,
		dryRun:    dryRun,
	}
}

func (o *Orchestrator) Process() error {
	jobs, err := o.processor.BuildProcessingJobs()
	if err != nil {
		return fmt.Errorf("failed to build processing jobs: %w", err)
	}

	log.Printf("Processing %d templates", len(jobs))

	for _, job := range jobs {
		if err := o.processJob(job); err != nil {
			return fmt.Errorf("failed to process template '%s': %w", job.Name, err)
		}
	}

	return nil
}

func (o *Orchestrator) processJob(job ProcessingJob) error {
	log.Printf("Processing template: %s", job.Name)

	if o.dryRun {
		return o.dryRunJob(job)
	}

	// TODO: Actual execution will be implemented in the next step
	log.Printf("Would execute: boilerplate %s", strings.Join(job.Args, " "))

	return nil
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