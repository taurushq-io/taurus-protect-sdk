package service

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
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
		if err := s.verifyAsset(asset, nil); err != nil {
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

	// Verify integrity of each asset — always enforced. Rows on a page normally share one
	// rules container, so steps 2-3 run once per distinct container rather than per row.
	containers := newInlineContainerCache(s.verifier)
	for _, asset := range assets {
		if asset != nil {
			if err := s.verifyAsset(asset, containers); err != nil {
				return nil, nil, fmt.Errorf("verification failed for asset %s: %w", asset.ID, err)
			}
		}
	}

	var limit, offset int64
	if opts != nil {
		limit, offset = opts.Limit, opts.Offset
	}
	return assets, assetPagination(resp.TotalItems, limit, offset), nil
}

// assetPagination builds the page window from the server's total.
func assetPagination(totalItems *string, limit, offset int64) *model.Pagination {
	if totalItems == nil {
		return nil
	}
	pagination := &model.Pagination{Limit: limit, Offset: offset}
	if total, err := strconv.ParseInt(*totalItems, 10, 64); err == nil {
		pagination.TotalItems = total
		// Overflow-safe: offset+limit can wrap on caller-supplied values.
		pagination.HasMore = total > offset && total-offset > limit
	}
	return pagination
}

// ListWhitelistedAssetsForApproval retrieves assets awaiting approval, verified the same
// way as ListWhitelistedAssets.
//
// Without this, the only reader of the for-approval endpoint was the unverified contract
// service — so the rows an approver inspects before whitelisting a contract address were
// never checked against governance.
func (s *WhitelistedAssetService) ListWhitelistedAssetsForApproval(
	ctx context.Context,
	opts *model.ListWhitelistedAssetsForApprovalOptions,
) ([]*model.WhitelistedAsset, *model.Pagination, error) {
	req := s.api.WhitelistServiceGetWhitelistedContractsForApproval(ctx)

	var limit, offset int64
	if opts != nil {
		if opts.Limit > 0 {
			req = req.Limit(fmt.Sprintf("%d", opts.Limit))
		}
		if opts.Offset > 0 {
			req = req.Offset(fmt.Sprintf("%d", opts.Offset))
		}
		if len(opts.IDs) > 0 {
			req = req.Ids(opts.IDs)
		}
		limit, offset = opts.Limit, opts.Offset
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, nil, s.errMapper.MapError(err, httpResp)
	}

	assets := mapper.WhitelistedAssetsFromDTO(resp.Result)

	containers := newInlineContainerCache(s.verifier)
	for _, asset := range assets {
		if asset != nil {
			if err := s.verifyAsset(asset, containers); err != nil {
				return nil, nil, fmt.Errorf("verification failed for asset %s: %w", asset.ID, err)
			}
		}
	}

	return assets, assetPagination(resp.TotalItems, limit, offset), nil
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

	// Carries only what the five steps need. Its Blockchain/Network come from the
	// DTO and are used solely to select the governance rules; they are NOT what
	// gets returned — step 6 below re-derives every identifying field from the
	// verified payload.
	subject := &model.WhitelistedAsset{
		ID:                    envelope.ID,
		Blockchain:            envelope.Blockchain,
		Network:               envelope.Network,
		Metadata:              envelope.Metadata,
		SignedContractAddress: envelope.SignedContractAddress,
		RulesContainer:        envelope.RulesContainer,
		RulesSignatures:       envelope.RulesSignatures,
	}

	// Steps 1-5
	result, err := s.verifier.VerifyWhitelistedAsset(
		subject,
		mapper.RulesContainerFromBase64,
		mapper.UserSignaturesFromBase64,
	)
	if err != nil {
		return err
	}

	// Step 6: parse the asset from the VERIFIED payload. Steps 1-5 prove the
	// envelope is authentic; this is what stops an unsigned value reaching the
	// caller. Marking the DTO-built subject as verified — which is what this used
	// to do — meant blockchain, network and the asset identity were never checked.
	verified, err := helper.ParseWhitelistedAssetFromJSON(envelope.Metadata.PayloadAsString)
	if err != nil {
		return &model.IntegrityError{
			Message: fmt.Sprintf("asset payload is verified but unparseable: %v", err),
		}
	}

	// Non-security fields the signed payload does not carry.
	verified.ID = envelope.ID
	verified.TenantID = envelope.TenantID
	verified.Status = envelope.Status
	verified.Metadata = envelope.Metadata
	verified.SignedContractAddress = envelope.SignedContractAddress
	verified.RulesContainer = envelope.RulesContainer
	verified.RulesSignatures = envelope.RulesSignatures

	envelope.SetVerified(verified, result.RulesContainer)

	return nil
}

// verifyAsset performs the 5-step integrity verification on a whitelisted asset.
// Returns nil if verification passes, or an error describing the failure.
//
// containers may be nil; when supplied, steps 2-3 come from the per-page memo.
func (s *WhitelistedAssetService) verifyAsset(asset *model.WhitelistedAsset, containers *inlineContainerCache) error {
	if s.verifier == nil {
		return &model.IntegrityError{Message: "verification is required but no verifier is configured"}
	}

	// Verification is enabled but required data is missing — this is an error.
	// An attacker could strip verification data to bypass checks.
	if asset.Metadata == nil || asset.RulesContainer == "" || asset.SignedContractAddress == nil {
		return &model.IntegrityError{Message: "verification enabled but required data missing"}
	}

	if containers != nil {
		cached, err := containers.get(asset.RulesContainer, asset.RulesSignatures)
		if err != nil {
			return err
		}
		if _, err := s.verifier.VerifyWhitelistedAsset(
			asset,
			mapper.RulesContainerFromBase64,
			mapper.UserSignaturesFromBase64,
			cached,
		); err != nil {
			return err
		}
		return populateVerifiedIdentity(asset)
	}

	if _, err := s.verifier.VerifyWhitelistedAsset(
		asset,
		mapper.RulesContainerFromBase64,
		mapper.UserSignaturesFromBase64,
	); err != nil {
		return err
	}
	return populateVerifiedIdentity(asset)
}

// populateVerifiedIdentity fills the identity fields from the signed payload — step 6 for
// assets, which the verifier does not perform.
//
// The DTO mapper cannot supply these: they are payload-only by contract, so without this
// a verified asset came back with an empty ContractAddress and a caller asking "is this
// contract whitelisted?" had to reach for the unverified reader to get an answer.
func populateVerifiedIdentity(asset *model.WhitelistedAsset) error {
	if asset.Metadata == nil {
		return &model.IntegrityError{Message: "metadata missing after verification"}
	}
	verified, err := helper.ParseWhitelistedAssetFromJSON(asset.Metadata.PayloadAsString)
	if err != nil {
		return fmt.Errorf("failed to parse verified asset: %w", err)
	}

	asset.ContractAddress = verified.ContractAddress
	asset.Name = verified.Name
	asset.Symbol = verified.Symbol
	asset.Decimals = verified.Decimals
	asset.TokenID = verified.TokenID
	return nil
}

// verifiedAssetsByID re-reads a batch of assets through the verifying list path,
// filtered by id, and returns them keyed by id.
//
//	ids ─▶ ONE filtered page ─▶ verify every row ─▶ map by id
//
// One round trip and one rules-container fetch for the whole batch. The per-id GET this
// replaced cost both per id, so a 50-id approval was 50 sequential round trips each
// running the full 6-step chain.
func (s *WhitelistedAssetService) verifiedAssetsByID(
	ctx context.Context,
	ids []string,
) (map[string]*model.WhitelistedAsset, error) {
	assets, _, err := s.ListWhitelistedAssets(ctx, &model.ListWhitelistedAssetsOptions{
		IDs:                ids,
		Limit:              int64(len(ids)),
		IncludeForApproval: true,
	})
	if err != nil {
		return nil, fmt.Errorf("refusing to sign: the verified read failed: %w", err)
	}

	byID := make(map[string]*model.WhitelistedAsset, len(assets))
	for _, asset := range assets {
		if asset != nil {
			byID[asset.ID] = asset
		}
	}
	return byID, nil
}

// ApproveWhitelistedAssets signs and submits an approval for the given whitelisted
// assets, all-or-nothing.
//
// It re-reads the batch through the verified path and signs the hashes THOSE rows carry,
// so the approver's signature covers metadata this SDK checked rather than whatever a
// caller was handed. The raw-signature form on WhitelistedContractService accepted an
// opaque blob over hashes nothing had verified.
//
// Any asset that is missing or fails verification aborts the whole call and nothing is
// signed: one signature covers every hash in the batch, so a partial approval would mean
// the caller believes they approved more than they did.
//
//	ids ─▶ sort numerically ─▶ ONE id-filtered verified page
//	                                          │
//	                                          ▼
//	                            completeness check ──miss──▶ abort, nothing signed
//	                                          │ok
//	                                          ▼
//	                          sign(JSON(hashes)) ─▶ POST once for the whole batch
func (s *WhitelistedAssetService) ApproveWhitelistedAssets(
	ctx context.Context,
	ids []string,
	privateKey *ecdsa.PrivateKey,
	comment string,
) error {
	if len(ids) == 0 {
		return fmt.Errorf("ids cannot be empty")
	}
	if privateKey == nil {
		return fmt.Errorf("privateKey cannot be nil")
	}
	if comment == "" {
		return fmt.Errorf("comment is required")
	}

	// Sorted numerically, as the request-approval path does, so the signed order is
	// independent of the order the caller passed.
	sortedIDs := make([]string, len(ids))
	copy(sortedIDs, ids)
	for _, id := range sortedIDs {
		if _, err := strconv.ParseInt(id, 10, 64); err != nil {
			return fmt.Errorf("whitelisted asset ID %q is not a valid numeric ID: %w", id, err)
		}
	}
	sort.Slice(sortedIDs, func(i, j int) bool {
		a, _ := strconv.ParseInt(sortedIDs[i], 10, 64)
		b, _ := strconv.ParseInt(sortedIDs[j], 10, 64)
		return a < b
	})

	// ONE id-filtered page through the verifying list path, not one GET per id. The list
	// path verifies every row and fetches the rules container once per call, so a 50-id
	// approval costs one round trip and one container fetch instead of fifty of each.
	// IncludeForApproval is required: the rows being approved are pending, so the default
	// list does not return them.
	verified, err := s.verifiedAssetsByID(ctx, sortedIDs)
	if err != nil {
		return err
	}

	hashes := make([]string, 0, len(sortedIDs))
	for _, id := range sortedIDs {
		asset, ok := verified[id]
		if !ok {
			// A page that silently omits a row must not become an approval of fewer
			// rows than the caller asked for.
			return &model.IntegrityError{
				Message: fmt.Sprintf("refusing to sign: asset %s was not returned by the verified read", id),
			}
		}
		if asset.Metadata == nil || asset.Metadata.Hash == "" {
			return &model.IntegrityError{
				Message: fmt.Sprintf("refusing to sign: asset %s has no metadata hash", id),
			}
		}
		hashes = append(hashes, asset.Metadata.Hash)
	}

	hashesJSON, err := json.Marshal(hashes)
	if err != nil {
		return fmt.Errorf("failed to serialize hashes: %w", err)
	}

	signature, err := crypto.SignData(privateKey, hashesJSON)
	if err != nil {
		return fmt.Errorf("failed to sign asset hashes: %w", err)
	}

	approveReq := openapi.TgvalidatordApproveWhitelistedContractAddressRequest{
		Ids:       sortedIDs,
		Signature: signature,
		Comment:   comment,
	}

	if _, httpResp, err := s.api.WhitelistServiceApproveWhitelistedContract(ctx).
		Body(approveReq).
		Execute(); err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}
