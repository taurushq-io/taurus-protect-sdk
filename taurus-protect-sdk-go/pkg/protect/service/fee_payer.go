package service

import (
	"context"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// FeePayerService provides fee payer management operations.
type FeePayerService struct {
	api       *openapi.FeePayersAPIService
	errMapper *ErrorMapper
}

// NewFeePayerService creates a new FeePayerService.
func NewFeePayerService(client *openapi.APIClient) *FeePayerService {
	return &FeePayerService{
		api:       client.FeePayersAPI,
		errMapper: NewErrorMapper(),
	}
}

// ListFeePayers retrieves one page of fee payers. Result.Pagination is never nil; continue with
// its NextOffset until HasMore is false.
func (s *FeePayerService) ListFeePayers(ctx context.Context, opts *model.ListFeePayersOptions) (*model.ListFeePayersResult, error) {
	if opts == nil {
		opts = &model.ListFeePayersOptions{}
	}
	window, err := resolveOffsetWindow(opts.Limit, opts.Offset)
	if err != nil {
		return nil, err
	}

	req := applyOffsetWindow(s.api.FeePayerServiceGetFeePayers(ctx), window)
	if len(opts.IDs) > 0 {
		req = req.Ids(opts.IDs)
	}
	if opts.Blockchain != "" {
		req = req.Blockchain(opts.Blockchain)
	}
	if opts.Network != "" {
		req = req.Network(opts.Network)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	pagination, err := offsetPagination(rulePlusRows, window, len(resp.Result), 0,
		offsetReply{TotalItems: resp.TotalItems})
	if err != nil {
		return nil, err
	}
	return &model.ListFeePayersResult{
		FeePayers:  mapper.FeePayersFromDTO(resp.Result),
		Pagination: pagination,
	}, nil
}

// GetFeePayer retrieves a single fee payer by ID.
func (s *FeePayerService) GetFeePayer(ctx context.Context, id string) (*model.FeePayer, error) {
	resp, httpResp, err := s.api.FeePayerServiceGetFeePayer(ctx, id).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.FeePayerFromDTO(resp.Feepayer), nil
}

// GetChecksum computes a checksum for the provided data.
func (s *FeePayerService) GetChecksum(ctx context.Context, req *model.ChecksumRequest) (*model.ChecksumResult, error) {
	body := mapper.ChecksumRequestToDTO(req)

	resp, httpResp, err := s.api.FeePayerServiceChecksum(ctx).Body(body).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.ChecksumResult{}
	if resp.Checksum != nil {
		result.Checksum = *resp.Checksum
	}

	return result, nil
}
