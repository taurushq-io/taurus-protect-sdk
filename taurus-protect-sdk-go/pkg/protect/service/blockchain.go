package service

import (
	"context"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// BlockchainService provides blockchain information operations.
type BlockchainService struct {
	api       *openapi.BlockchainAPIService
	errMapper *ErrorMapper
}

// NewBlockchainService creates a new BlockchainService.
func NewBlockchainService(client *openapi.APIClient) *BlockchainService {
	return &BlockchainService{
		api:       client.BlockchainAPI,
		errMapper: NewErrorMapper(),
	}
}

// ListBlockchains retrieves a list of all enabled blockchains.
func (s *BlockchainService) ListBlockchains(ctx context.Context, opts *model.ListBlockchainsOptions) ([]*model.Blockchain, error) {
	req := s.api.BlockchainServiceGetBlockchains(ctx)

	if opts != nil {
		if opts.Blockchain != "" {
			req = req.Blockchain(opts.Blockchain)
		}
		if opts.Network != "" {
			req = req.Network(opts.Network)
		}
		if opts.IncludeBlockHeight {
			req = req.IncludeBlockHeight(true)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.BlockchainsFromDTO(resp.Blockchains), nil
}
