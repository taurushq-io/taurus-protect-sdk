package mapper

import (
	"strconv"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// JobFromDTO converts an OpenAPI TgvalidatordJob to a domain Job.
func JobFromDTO(dto *openapi.TgvalidatordJob) *model.Job {
	if dto == nil {
		return nil
	}

	job := &model.Job{
		Name: safeString(dto.Name),
	}

	if dto.Statistics != nil {
		job.Statistics = JobStatisticsFromDTO(dto.Statistics)
	}

	return job
}

// JobsFromDTO converts a slice of OpenAPI TgvalidatordJob to domain Jobs.
func JobsFromDTO(dtos []openapi.TgvalidatordJob) []*model.Job {
	if dtos == nil {
		return nil
	}
	jobs := make([]*model.Job, len(dtos))
	for i := range dtos {
		jobs[i] = JobFromDTO(&dtos[i])
	}
	return jobs
}

// JobStatisticsFromDTO converts an OpenAPI TgvalidatordJobStatistics to a domain JobStatistics.
func JobStatisticsFromDTO(dto *openapi.TgvalidatordJobStatistics) *model.JobStatistics {
	if dto == nil {
		return nil
	}

	stats := &model.JobStatistics{
		AvgDuration: safeString(dto.AvgDuration),
		MaxDuration: safeString(dto.MaxDuration),
		MinDuration: safeString(dto.MinDuration),
	}

	// Parse numeric string fields
	if dto.Pending != nil {
		if pending, err := strconv.ParseInt(*dto.Pending, 10, 64); err == nil {
			stats.Pending = pending
		}
	}
	if dto.Successes != nil {
		if successes, err := strconv.ParseInt(*dto.Successes, 10, 64); err == nil {
			stats.Successes = successes
		}
	}
	if dto.Failures != nil {
		if failures, err := strconv.ParseInt(*dto.Failures, 10, 64); err == nil {
			stats.Failures = failures
		}
	}

	// Convert job status fields
	if dto.LastSuccess != nil {
		stats.LastSuccess = JobStatusFromDTO(dto.LastSuccess)
	}
	if dto.LastFailure != nil {
		stats.LastFailure = JobStatusFromDTO(dto.LastFailure)
	}

	return stats
}

// JobStatusFromDTO converts an OpenAPI TgvalidatordJobStatus to a domain JobStatus.
func JobStatusFromDTO(dto *openapi.TgvalidatordJobStatus) *model.JobStatus {
	if dto == nil {
		return nil
	}

	status := &model.JobStatus{
		ID:      safeString(dto.Id),
		Message: safeString(dto.Message),
		Status:  safeString(dto.Status),
	}

	// Convert time fields
	if dto.StartedAt != nil {
		status.StartedAt = *dto.StartedAt
	}
	if dto.UpdatedAt != nil {
		status.UpdatedAt = *dto.UpdatedAt
	}
	if dto.TimeoutAt != nil {
		status.TimeoutAt = *dto.TimeoutAt
	}

	return status
}
