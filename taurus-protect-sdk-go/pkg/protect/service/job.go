package service

import (
	"context"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// JobService provides job management operations.
type JobService struct {
	api       *openapi.JobsAPIService
	errMapper *ErrorMapper
}

// NewJobService creates a new JobService.
func NewJobService(client *openapi.APIClient) *JobService {
	return &JobService{
		api:       client.JobsAPI,
		errMapper: NewErrorMapper(),
	}
}

// GetJob retrieves a single job by name.
func (s *JobService) GetJob(ctx context.Context, name string) (*model.Job, error) {
	req := s.api.JobServiceGetJob(ctx, name)

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.JobFromDTO(resp.Job), nil
}

// ListJobs retrieves all registered jobs.
func (s *JobService) ListJobs(ctx context.Context) ([]*model.Job, error) {
	req := s.api.JobServiceGetJobs(ctx)

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.JobsFromDTO(resp.Jobs), nil
}

// GetJobStatus retrieves the status of a specific job execution.
func (s *JobService) GetJobStatus(ctx context.Context, name string, id string) (*model.JobStatus, error) {
	req := s.api.JobServiceGetJobStatus(ctx, name, id)

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.JobStatusFromDTO(resp.Status), nil
}
