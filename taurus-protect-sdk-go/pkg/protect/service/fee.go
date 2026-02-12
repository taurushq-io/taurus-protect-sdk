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

// GetFees retrieves a list of fee estimates.
//
// Deprecated: Use GetFeesV2 instead. This endpoint is deprecated and may be removed in a future version.
func (s *FeeService) GetFees(ctx context.Context) (*model.GetFeesResult, error) {
	req := s.api.FeeServiceGetFees(ctx)

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return &model.GetFeesResult{
		Fees: mapper.FeesFromDTO(resp.Result),
	}, nil
}

// GetFeesV2 retrieves a list of native currency fee estimates.
// This is the recommended method for getting fee information.
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
