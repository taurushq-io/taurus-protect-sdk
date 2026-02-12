package service

import (
	"context"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// HealthService provides health check operations.
type HealthService struct {
	api       *openapi.HealthAPIService
	errMapper *ErrorMapper
}

// NewHealthService creates a new HealthService.
func NewHealthService(client *openapi.APIClient) *HealthService {
	return &HealthService{
		api:       client.HealthAPI,
		errMapper: NewErrorMapper(),
	}
}

// GetAllHealthChecks retrieves all health checks with optional filtering.
func (s *HealthService) GetAllHealthChecks(ctx context.Context, opts *model.GetAllHealthChecksOptions) (*model.GetAllHealthChecksResult, error) {
	req := s.api.HealthServiceGetAllHealthChecks(ctx)

	if opts != nil {
		if opts.TenantID != "" {
			req = req.TenantId(opts.TenantID)
		}
		if opts.FailIfUnhealthy {
			req = req.FailIfUnhealthy(opts.FailIfUnhealthy)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.GetAllHealthChecksResult{
		Components: mapper.HealthComponentsFromDTO(resp.Components),
	}

	return result, nil
}
