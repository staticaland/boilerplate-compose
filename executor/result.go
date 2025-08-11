package executor

import (
	"fmt"
	"time"
)

type ExecutionResult struct {
	TemplateName string
	Success      bool
	Error        error
	Duration     time.Duration
	StartTime    time.Time
	EndTime      time.Time
}

type ExecutionSummary struct {
	Results       []ExecutionResult
	TotalDuration time.Duration
	SuccessCount  int
	FailureCount  int
}

func NewExecutionSummary() *ExecutionSummary {
	return &ExecutionSummary{
		Results: make([]ExecutionResult, 0),
	}
}

func (s *ExecutionSummary) AddResult(result ExecutionResult) {
	s.Results = append(s.Results, result)
	if result.Success {
		s.SuccessCount++
	} else {
		s.FailureCount++
	}
	s.TotalDuration += result.Duration
}

func (s *ExecutionSummary) Print() {
	fmt.Printf("\n=== Execution Summary ===\n")
	fmt.Printf("Total templates: %d\n", len(s.Results))
	fmt.Printf("Successful: %d\n", s.SuccessCount)
	fmt.Printf("Failed: %d\n", s.FailureCount)
	fmt.Printf("Total duration: %v\n", s.TotalDuration)

	if s.FailureCount > 0 {
		fmt.Printf("\nFailed templates:\n")
		for _, result := range s.Results {
			if !result.Success {
				fmt.Printf("  - %s: %v\n", result.TemplateName, result.Error)
			}
		}
	}

	fmt.Printf("\nTemplate execution times:\n")
	for _, result := range s.Results {
		status := "✓"
		if !result.Success {
			status = "✗"
		}
		fmt.Printf("  %s %s: %v\n", status, result.TemplateName, result.Duration)
	}
}