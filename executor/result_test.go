package executor

import (
	"errors"
	"testing"
	"time"
)

func TestExecutionResult(t *testing.T) {
	result := ExecutionResult{
		TemplateName: "test-template",
		Success:      true,
		Error:        nil,
		Duration:     100 * time.Millisecond,
		StartTime:    time.Now(),
		EndTime:      time.Now().Add(100 * time.Millisecond),
	}

	if result.TemplateName != "test-template" {
		t.Errorf("expected TemplateName 'test-template', got %q", result.TemplateName)
	}

	if !result.Success {
		t.Error("expected Success to be true")
	}

	if result.Error != nil {
		t.Errorf("expected Error to be nil, got %v", result.Error)
	}
}

func TestNewExecutionSummary(t *testing.T) {
	summary := NewExecutionSummary()

	if summary == nil {
		t.Fatal("expected non-nil ExecutionSummary")
	}

	if len(summary.Results) != 0 {
		t.Errorf("expected empty Results slice, got length %d", len(summary.Results))
	}

	if summary.SuccessCount != 0 {
		t.Errorf("expected SuccessCount 0, got %d", summary.SuccessCount)
	}

	if summary.FailureCount != 0 {
		t.Errorf("expected FailureCount 0, got %d", summary.FailureCount)
	}

	if summary.TotalDuration != 0 {
		t.Errorf("expected TotalDuration 0, got %v", summary.TotalDuration)
	}
}

func TestExecutionSummary_AddResult(t *testing.T) {
	summary := NewExecutionSummary()
	
	// Add successful result
	successResult := ExecutionResult{
		TemplateName: "success-template",
		Success:      true,
		Duration:     50 * time.Millisecond,
	}
	summary.AddResult(successResult)

	if len(summary.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(summary.Results))
	}

	if summary.SuccessCount != 1 {
		t.Errorf("expected SuccessCount 1, got %d", summary.SuccessCount)
	}

	if summary.FailureCount != 0 {
		t.Errorf("expected FailureCount 0, got %d", summary.FailureCount)
	}

	// Add failed result
	failureResult := ExecutionResult{
		TemplateName: "failure-template",
		Success:      false,
		Error:        errors.New("test error"),
		Duration:     30 * time.Millisecond,
	}
	summary.AddResult(failureResult)

	if len(summary.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(summary.Results))
	}

	if summary.SuccessCount != 1 {
		t.Errorf("expected SuccessCount 1, got %d", summary.SuccessCount)
	}

	if summary.FailureCount != 1 {
		t.Errorf("expected FailureCount 1, got %d", summary.FailureCount)
	}

	expectedDuration := 80 * time.Millisecond
	if summary.TotalDuration != expectedDuration {
		t.Errorf("expected TotalDuration %v, got %v", expectedDuration, summary.TotalDuration)
	}
}