package service

import (
	"context"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// ConfigService provides configuration management operations.
type ConfigService struct {
	api       *openapi.ConfigAPIService
	errMapper *ErrorMapper
}

// NewConfigService creates a new ConfigService.
func NewConfigService(client *openapi.APIClient) *ConfigService {
	return &ConfigService{
		api:       client.ConfigAPI,
		errMapper: NewErrorMapper(),
	}
}

// GetTenantConfig retrieves the configuration of the tenant for the connected user.
func (s *ConfigService) GetTenantConfig(ctx context.Context) (*model.TenantConfig, error) {
	req := s.api.StatusServiceGetConfigTenant(ctx)

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.TenantConfigFromDTO(resp), nil
}
