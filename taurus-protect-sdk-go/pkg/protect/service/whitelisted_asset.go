package service

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"strconv"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/helper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// WhitelistedAssetService provides whitelisted asset (contract) management operations.
type WhitelistedAssetService struct {
	api       *openapi.ContractWhitelistingAPIService
	errMapper *ErrorMapper
	verifier  *helper.WhitelistedAssetVerifier
}

// WhitelistedAssetServiceConfig holds configuration for the WhitelistedAssetService.
type WhitelistedAssetServiceConfig struct {
	// SuperAdminKeys are the public keys used to verify governance rules signatures.
	SuperAdminKeys []*ecdsa.PublicKey
	// MinValidSignatures is the minimum number of valid SuperAdmin signatures required.
	MinValidSignatures int
}

// NewWhitelistedAssetServiceWithVerification creates a new WhitelistedAssetService with
// integrity verification enabled. When SuperAdmin keys are provided, all retrieved assets
// will be cryptographically verified before being returned.
func NewWhitelistedAssetServiceWithVerification(
	client *openapi.APIClient,
	config *WhitelistedAssetServiceConfig,
) *WhitelistedAssetService {
	svc := &WhitelistedAssetService{
		api:       client.ContractWhitelistingAPI,
		errMapper: NewErrorMapper(),
	}

	// Create verifier if SuperAdmin keys are configured
	if config != nil && len(config.SuperAdminKeys) > 0 {
		svc.verifier = helper.NewWhitelistedAssetVerifier(
			config.SuperAdminKeys,
			config.MinValidSignatures,
		)
	}

	return svc
}

// GetWhitelistedAsset retrieves a whitelisted asset by ID.
// If verification is enabled (SuperAdmin keys configured), the asset integrity
// will be cryptographically verified before being returned.
func (s *WhitelistedAssetService) GetWhitelistedAsset(ctx context.Context, id string) (*model.WhitelistedAsset, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	resp, httpResp, err := s.api.WhitelistServiceGetWhitelistedContract(ctx, id).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, nil
	}

	asset := mapper.WhitelistedAssetFromDTO(resp.Result)

	// Verify integrity — always enforced
	if asset != nil {
		if err := s.verifyAsset(asset); err != nil {
			return nil, err
		}
	}

	return asset, nil
}

// ListWhitelistedAssets retrieves a list of whitelisted assets.
func (s *WhitelistedAssetService) ListWhitelistedAssets(ctx context.Context, opts *model.ListWhitelistedAssetsOptions) ([]*model.WhitelistedAsset, *model.Pagination, error) {
	req := s.api.WhitelistServiceGetWhitelistedContracts(ctx)

	if opts != nil {
		if opts.Limit > 0 {
			req = req.Limit(fmt.Sprintf("%d", opts.Limit))
		}
		if opts.Offset > 0 {
			req = req.Offset(fmt.Sprintf("%d", opts.Offset))
		}
		if opts.Query != "" {
			req = req.Query(opts.Query)
		}
		if opts.Blockchain != "" {
			req = req.Blockchain(opts.Blockchain)
		}
		if opts.Network != "" {
			req = req.Network(opts.Network)
		}
		if opts.IncludeForApproval {
			req = req.IncludeForApproval(true)
		}
		if len(opts.KindTypes) > 0 {
			req = req.KindTypes(opts.KindTypes)
		}
		if len(opts.IDs) > 0 {
			req = req.WhitelistedContractAddressIds(opts.IDs)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, nil, s.errMapper.MapError(err, httpResp)
	}

	assets := mapper.WhitelistedAssetsFromDTO(resp.Result)

	// Verify integrity of each asset — always enforced
	for _, asset := range assets {
		if asset != nil {
			if err := s.verifyAsset(asset); err != nil {
				return nil, nil, fmt.Errorf("verification failed for asset %s: %w", asset.ID, err)
			}
		}
	}

	var pagination *model.Pagination
	if resp.TotalItems != nil {
		pagination = &model.Pagination{}
		if total, parseErr := strconv.ParseInt(*resp.TotalItems, 10, 64); parseErr == nil {
			pagination.TotalItems = total
		}
		if opts != nil {
			pagination.Limit = opts.Limit
			pagination.Offset = opts.Offset
			pagination.HasMore = pagination.Offset+pagination.Limit < pagination.TotalItems
		}
	}

	return assets, pagination, nil
}

// GetWhitelistedAssetEnvelope retrieves a whitelisted asset envelope by ID and performs
// the complete 5-step verification flow. The envelope contains both the verified whitelisted
// asset and the decoded rules container.
//
// This method requires verification to be enabled (SuperAdmin keys must be configured).
// If verification is not configured, an error is returned.
func (s *WhitelistedAssetService) GetWhitelistedAssetEnvelope(
	ctx context.Context,
	id string,
) (*model.WhitelistedAssetEnvelope, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	if s.verifier == nil {
		return nil, &model.IntegrityError{
			Message: "verification is required for GetWhitelistedAssetEnvelope but no SuperAdmin keys are configured",
		}
	}

	resp, httpResp, err := s.api.WhitelistServiceGetWhitelistedContract(ctx, id).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, nil
	}

	// Create envelope from DTO
	envelope := mapper.WhitelistedAssetEnvelopeFromDTO(resp.Result)

	// Initialize and verify the envelope
	if err := s.initializeAssetEnvelope(envelope); err != nil {
		return nil, err
	}

	return envelope, nil
}

// initializeAssetEnvelope performs the 5-step verification and populates the verified fields.
func (s *WhitelistedAssetService) initializeAssetEnvelope(envelope *model.WhitelistedAssetEnvelope) error {
	if envelope == nil {
		return fmt.Errorf("envelope cannot be nil")
	}

	// Check required data for verification
	if envelope.Metadata == nil {
		return &model.IntegrityError{Message: "metadata is required for verification"}
	}
	if envelope.RulesContainer == "" {
		return &model.IntegrityError{Message: "rules container is required for verification"}
	}
	if envelope.SignedContractAddress == nil {
		return &model.IntegrityError{Message: "signed contract address is required for verification"}
	}

	// Create a temporary WhitelistedAsset for verification
	tempAsset := &model.WhitelistedAsset{
		ID:                   envelope.ID,
		Blockchain:           envelope.Blockchain,
		Network:              envelope.Network,
		Metadata:             envelope.Metadata,
		SignedContractAddress: envelope.SignedContractAddress,
		RulesContainer:       envelope.RulesContainer,
		RulesSignatures:      envelope.RulesSignatures,
	}

	// Perform the 5-step verification
	result, err := s.verifier.VerifyWhitelistedAsset(
		tempAsset,
		mapper.RulesContainerFromBase64,
		mapper.UserSignaturesFromBase64,
	)
	if err != nil {
		return err
	}

	// Set verified data on envelope
	// For assets, the verified "asset" is the tempAsset with verified fields
	envelope.SetVerified(tempAsset, result.RulesContainer)

	return nil
}

// verifyAsset performs the 5-step integrity verification on a whitelisted asset.
// Returns nil if verification passes, or an error describing the failure.
func (s *WhitelistedAssetService) verifyAsset(asset *model.WhitelistedAsset) error {
	if s.verifier == nil {
		return &model.IntegrityError{Message: "verification is required but no verifier is configured"}
	}

	// Verification is enabled but required data is missing — this is an error.
	// An attacker could strip verification data to bypass checks.
	if asset.Metadata == nil || asset.RulesContainer == "" || asset.SignedContractAddress == nil {
		return &model.IntegrityError{Message: "verification enabled but required data missing"}
	}

	_, err := s.verifier.VerifyWhitelistedAsset(
		asset,
		mapper.RulesContainerFromBase64,
		mapper.UserSignaturesFromBase64,
	)
	return err
}
