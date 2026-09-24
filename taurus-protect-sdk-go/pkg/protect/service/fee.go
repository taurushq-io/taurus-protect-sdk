package service

import (
	"context"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// FeeService provides fee estimation operations.
type FeeService struct {
	api       *openapi.FeeAPIService
	errMapper *ErrorMapper
}

// NewFeeService creates a new FeeService.
func NewFeeService(client *openapi.APIClient) *FeeService {
	return &FeeService{
		api:       client.FeeAPI,
		errMapper: NewErrorMapper(),
	}
}

// GetFeesV2 retrieves the native currency fee estimates (FeeService_GetFeesV2). The deprecated
// v1 GetFees is not wrapped.
func (s *FeeService) GetFeesV2(ctx context.Context) (*model.GetFeesV2Result, error) {
	req := s.api.FeeServiceGetFeesV2(ctx)

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return &model.GetFeesV2Result{
		Fees: mapper.FeesV2FromDTO(resp.Result),
	}, nil
}
