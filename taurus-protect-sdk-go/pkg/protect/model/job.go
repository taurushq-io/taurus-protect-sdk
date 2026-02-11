package model

import "time"

// Job represents a scheduled job in the system.
type Job struct {
	// Name is the unique identifier for the job.
	Name string `json:"name"`
	// Statistics contains statistics about the job's execution history.
	Statistics *JobStatistics `json:"statistics,omitempty"`
}

// JobStatistics contains execution statistics for a job.
type JobStatistics struct {
	// Pending is the number of pending job executions.
	Pending int64 `json:"pending"`
	// Successes is the total number of successful job executions.
	Successes int64 `json:"successes"`
	// Failures is the total number of failed job executions.
	Failures int64 `json:"failures"`
	// LastSuccess contains details of the last successful execution.
	LastSuccess *JobStatus `json:"last_success,omitempty"`
	// LastFailure contains details of the last failed execution.
	LastFailure *JobStatus `json:"last_failure,omitempty"`
	// AvgDuration is the average duration of job executions.
	AvgDuration string `json:"avg_duration,omitempty"`
	// MaxDuration is the maximum duration of job executions.
	MaxDuration string `json:"max_duration,omitempty"`
	// MinDuration is the minimum duration of job executions.
	MinDuration string `json:"min_duration,omitempty"`
}

// JobStatus represents the status of a specific job execution.
type JobStatus struct {
	// ID is the unique identifier for this job execution.
	ID string `json:"id"`
	// StartedAt is when the job execution started.
	StartedAt time.Time `json:"started_at"`
	// UpdatedAt is when the job status was last updated.
	UpdatedAt time.Time `json:"updated_at"`
	// TimeoutAt is when the job will timeout if not completed.
	TimeoutAt time.Time `json:"timeout_at"`
	// Message contains details about the job execution.
	Message string `json:"message,omitempty"`
	// Status is the current status of the job execution (e.g., "running", "completed", "failed").
	Status string `json:"status"`
}
