package service

import (
	"context"
	"fmt"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model/taurusnetwork"
)

// TaurusNetworkSharingService provides Taurus-NETWORK shared address and asset operations.
type TaurusNetworkSharingService struct {
	api       *openapi.TaurusNetworkSharedAddressAssetAPIService
	errMapper *ErrorMapper
}

// NewTaurusNetworkSharingService creates a new TaurusNetworkSharingService.
func NewTaurusNetworkSharingService(client *openapi.APIClient) *TaurusNetworkSharingService {
	return &TaurusNetworkSharingService{
		api:       client.TaurusNetworkSharedAddressAssetAPI,
		errMapper: NewErrorMapper(),
	}
}

// ListSharedAddresses retrieves a list of shared addresses with optional filtering and pagination.
func (s *TaurusNetworkSharingService) ListSharedAddresses(ctx context.Context, opts *taurusnetwork.ListSharedAddressesOptions) (*taurusnetwork.ListSharedAddressesResult, error) {
	req := s.api.TaurusNetworkServiceGetSharedAddresses(ctx)

	if opts != nil {
		if opts.ParticipantID != "" {
			req = req.ParticipantID(opts.ParticipantID)
		}
		if opts.OwnerParticipantID != "" {
			req = req.OwnerParticipantID(opts.OwnerParticipantID)
		}
		if opts.TargetParticipantID != "" {
			req = req.TargetParticipantID(opts.TargetParticipantID)
		}
		if opts.Blockchain != "" {
			req = req.Blockchain(opts.Blockchain)
		}
		if opts.Network != "" {
			req = req.Network(opts.Network)
		}
		if len(opts.IDs) > 0 {
			req = req.Ids(opts.IDs)
		}
		if len(opts.Statuses) > 0 {
			req = req.Statuses(opts.Statuses)
		}
		if opts.SortOrder != "" {
			req = req.SortOrder(opts.SortOrder)
		}
		if opts.CurrentPage != "" {
			req = req.CursorCurrentPage(opts.CurrentPage)
		}
		if opts.PageRequest != "" {
			req = req.CursorPageRequest(opts.PageRequest)
		}
		if opts.PageSize > 0 {
			req = req.CursorPageSize(fmt.Sprintf("%d", opts.PageSize))
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &taurusnetwork.ListSharedAddressesResult{
		SharedAddresses: mapper.SharedAddressesFromDTO(resp.SharedAddresses),
	}

	// Parse cursor pagination info
	if resp.Cursor != nil {
		if resp.Cursor.CurrentPage != nil {
			result.CurrentPage = *resp.Cursor.CurrentPage
		}
		if resp.Cursor.HasPrevious != nil {
			result.HasPrevious = *resp.Cursor.HasPrevious
		}
		if resp.Cursor.HasNext != nil {
			result.HasNext = *resp.Cursor.HasNext
		}
	}

	return result, nil
}

// ListSharedAssets retrieves a list of shared assets with optional filtering and pagination.
func (s *TaurusNetworkSharingService) ListSharedAssets(ctx context.Context, opts *taurusnetwork.ListSharedAssetsOptions) (*taurusnetwork.ListSharedAssetsResult, error) {
	req := s.api.TaurusNetworkServiceGetSharedAssets(ctx)

	if opts != nil {
		if opts.ParticipantID != "" {
			req = req.ParticipantID(opts.ParticipantID)
		}
		if opts.OwnerParticipantID != "" {
			req = req.OwnerParticipantID(opts.OwnerParticipantID)
		}
		if opts.TargetParticipantID != "" {
			req = req.TargetParticipantID(opts.TargetParticipantID)
		}
		if opts.Blockchain != "" {
			req = req.Blockchain(opts.Blockchain)
		}
		if opts.Network != "" {
			req = req.Network(opts.Network)
		}
		if len(opts.IDs) > 0 {
			req = req.Ids(opts.IDs)
		}
		if len(opts.Statuses) > 0 {
			req = req.Statuses(opts.Statuses)
		}
		if opts.SortOrder != "" {
			req = req.SortOrder(opts.SortOrder)
		}
		if opts.CurrentPage != "" {
			req = req.CursorCurrentPage(opts.CurrentPage)
		}
		if opts.PageRequest != "" {
			req = req.CursorPageRequest(opts.PageRequest)
		}
		if opts.PageSize > 0 {
			req = req.CursorPageSize(fmt.Sprintf("%d", opts.PageSize))
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &taurusnetwork.ListSharedAssetsResult{
		SharedAssets: mapper.SharedAssetsFromDTO(resp.SharedAssets),
	}

	// Parse cursor pagination info
	if resp.Cursor != nil {
		if resp.Cursor.CurrentPage != nil {
			result.CurrentPage = *resp.Cursor.CurrentPage
		}
		if resp.Cursor.HasPrevious != nil {
			result.HasPrevious = *resp.Cursor.HasPrevious
		}
		if resp.Cursor.HasNext != nil {
			result.HasNext = *resp.Cursor.HasNext
		}
	}

	return result, nil
}

// ShareAddress shares an internal address with a Taurus-NETWORK participant.
// It automatically creates a whitelisted address to be approved/rejected on the target participant side.
func (s *TaurusNetworkSharingService) ShareAddress(ctx context.Context, request *taurusnetwork.ShareAddressRequest) error {
	if request == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if request.ToParticipantID == "" {
		return fmt.Errorf("to_participant_id is required")
	}
	if request.AddressID == "" {
		return fmt.Errorf("address_id is required")
	}

	body := openapi.TgvalidatordShareAddressRequest{
		ToParticipantID: request.ToParticipantID,
		AddressID:       request.AddressID,
	}

	// Convert key-value attributes if present
	if len(request.KeyValueAttributes) > 0 {
		body.KeyValueAttributes = mapper.KeyValueAttributesToDTO(request.KeyValueAttributes)
	}

	req := s.api.TaurusNetworkServiceShareAddress(ctx).Body(body)

	_, httpResp, err := req.Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// ShareWhitelistedAsset shares an asset with a Taurus-NETWORK participant.
// It automatically creates a whitelisted asset to be approved/rejected on the target participant side.
func (s *TaurusNetworkSharingService) ShareWhitelistedAsset(ctx context.Context, request *taurusnetwork.ShareWhitelistedAssetRequest) error {
	if request == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if request.ToParticipantID == "" {
		return fmt.Errorf("to_participant_id is required")
	}
	if request.WhitelistedContractID == "" {
		return fmt.Errorf("whitelisted_contract_id is required")
	}

	body := openapi.TgvalidatordShareWhitelistedAssetRequest{
		ToParticipantID:       request.ToParticipantID,
		WhitelistedContractID: request.WhitelistedContractID,
	}

	req := s.api.TaurusNetworkServiceShareWhitelistedAsset(ctx).Body(body)

	_, httpResp, err := req.Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// UnshareAddress unshares an address with a Taurus-NETWORK participant.
// The address must be shared with the participant to be unshared.
// Unsharing an address will not delete the shared address in the participant registry,
// it will update the status of the shared address to 'unshared'.
func (s *TaurusNetworkSharingService) UnshareAddress(ctx context.Context, sharedAddressID string) error {
	if sharedAddressID == "" {
		return fmt.Errorf("shared_address_id is required")
	}

	// The API requires an empty body for the unshare operation
	body := make(map[string]interface{})
	req := s.api.TaurusNetworkServiceUnshareAddress(ctx, sharedAddressID).Body(body)

	_, httpResp, err := req.Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// UnshareWhitelistedAsset unshares an asset with a Taurus-NETWORK participant.
// The asset must be shared with the participant to be unshared.
// Unsharing an asset will not delete the shared asset in the participant registry,
// it will update the status of the shared asset to 'unshared'.
func (s *TaurusNetworkSharingService) UnshareWhitelistedAsset(ctx context.Context, sharedAssetID string) error {
	if sharedAssetID == "" {
		return fmt.Errorf("shared_asset_id is required")
	}

	// The API requires an empty body for the unshare operation
	body := make(map[string]interface{})
	req := s.api.TaurusNetworkServiceUnshareWhitelistedAsset(ctx, sharedAssetID).Body(body)

	_, httpResp, err := req.Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}
