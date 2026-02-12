package service

import (
	"context"
	"fmt"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// TokenMetadataService provides token metadata retrieval operations.
type TokenMetadataService struct {
	api       *openapi.TokenMetadataAPIService
	errMapper *ErrorMapper
}

// NewTokenMetadataService creates a new TokenMetadataService.
func NewTokenMetadataService(client *openapi.APIClient) *TokenMetadataService {
	return &TokenMetadataService{
		api:       client.TokenMetadataAPI,
		errMapper: NewErrorMapper(),
	}
}

// GetCryptoPunkMetadata retrieves metadata for a CryptoPunk NFT.
// The network is typically "mainnet", contract should be the CryptoPunks contract address
// (0xb47e3cd837dDF8e4c57F05d70Ab865de6e193BBB for ETH mainnet), and tokenID is the punk ID (0-9999).
func (s *TokenMetadataService) GetCryptoPunkMetadata(ctx context.Context, network, contract, tokenID string) (*model.CryptoPunkMetadata, error) {
	if network == "" {
		return nil, fmt.Errorf("network cannot be empty")
	}
	if contract == "" {
		return nil, fmt.Errorf("contract cannot be empty")
	}
	if tokenID == "" {
		return nil, fmt.Errorf("tokenID cannot be empty")
	}

	resp, httpResp, err := s.api.TokenMetadataServiceGetCryptoPunksTokenMetadata(ctx, network, contract, tokenID).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, fmt.Errorf("CryptoPunk metadata not found")
	}

	return mapper.CryptoPunkMetadataFromDTO(resp.Result), nil
}

// GetERCTokenMetadata retrieves metadata for an ERC-721 or ERC-1155 token.
// This uses the newer EVM endpoint that supports multiple blockchains.
// The network specifies the network (e.g., "mainnet"), contract is the token contract address,
// and tokenID is the token ID.
func (s *TokenMetadataService) GetERCTokenMetadata(ctx context.Context, network, contract, tokenID string, opts *model.GetERCTokenMetadataOptions) (*model.ERCTokenMetadata, error) {
	if network == "" {
		return nil, fmt.Errorf("network cannot be empty")
	}
	if contract == "" {
		return nil, fmt.Errorf("contract cannot be empty")
	}
	if tokenID == "" {
		return nil, fmt.Errorf("tokenID cannot be empty")
	}

	req := s.api.TokenMetadataServiceGetEVMERCTokenMetadata(ctx, network, contract, tokenID)

	if opts != nil {
		if opts.WithData {
			req = req.WithData(true)
		}
		if opts.Blockchain != "" {
			req = req.Blockchain(opts.Blockchain)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, fmt.Errorf("ERC token metadata not found")
	}

	return mapper.ERCTokenMetadataFromDTO(resp.Result), nil
}
