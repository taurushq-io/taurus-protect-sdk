package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/cache"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/helper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// AssetService provides asset balance retrieval operations.
type AssetService struct {
	api        *openapi.AssetsAPIService
	errMapper  *ErrorMapper
	rulesCache *cache.RulesContainerCache
}

// NewAssetService creates a new AssetService with mandatory address signature
// verification.
//
// GetAssetAddresses returns the same Address entity AddressService does, signature and
// all, and used to hand it over unverified — so the mandatory verification on
// AddressService could be walked around by asking for the same rows here.
//
// Panics if rulesCache is nil, as NewAddressService does.
func NewAssetService(client *openapi.APIClient, rulesCache *cache.RulesContainerCache) *AssetService {
	if rulesCache == nil {
		panic("rulesCache cannot be nil - address signature verification is mandatory")
	}
	return &AssetService{
		api:        client.AssetsAPI,
		errMapper:  NewErrorMapper(),
		rulesCache: rulesCache,
	}
}

// GetAssetAddresses retrieves address-level balances for a specific asset.
func (s *AssetService) GetAssetAddresses(ctx context.Context, req *model.GetAssetAddressesRequest) (*model.GetAssetAddressesResult, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.Asset.Currency == "" {
		return nil, fmt.Errorf("asset currency is required")
	}

	// Build the OpenAPI request
	apiReq := openapi.TgvalidatordGetAssetAddressesRequest{
		Asset: mapper.AssetFilterToDTO(&req.Asset),
	}

	// Set optional fields
	if req.Limit > 0 {
		limit := fmt.Sprintf("%d", req.Limit)
		apiReq.Limit = &limit
	}
	if req.Cursor != "" {
		apiReq.Cursor = &req.Cursor
	}
	if req.SortOrder != "" {
		sorting := openapi.TgvalidatordGetAssetAddressesRequestSorting{}
		sortOrder := openapi.TgvalidatordGetAssetAddressesRequestSortingSortOrder(req.SortOrder)
		sorting.SortOrder = &sortOrder
		apiReq.Sorting = &sorting
	}
	if req.WalletID != "" {
		apiReq.WalletId = &req.WalletID
	}
	if req.AddressID != "" {
		apiReq.AddressId = &req.AddressID
	}
	if len(req.Addresses) > 0 {
		apiReq.Addresses = req.Addresses
	}

	// Handle cursor-based pagination
	if req.CurrentPage != "" || req.PageRequest != "" || req.PageSize > 0 {
		cursor := openapi.TgvalidatordRequestCursor{}
		if req.CurrentPage != "" {
			cursor.CurrentPage = &req.CurrentPage
		}
		if req.PageRequest != "" {
			cursor.PageRequest = &req.PageRequest
		}
		if req.PageSize > 0 {
			pageSize := fmt.Sprintf("%d", req.PageSize)
			cursor.PageSize = &pageSize
		}
		apiReq.RequestCursor = &cursor
	}

	resp, httpResp, err := s.api.WalletServiceGetAssetAddresses(ctx).Body(apiReq).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	addresses := mapper.AddressesFromDTO(resp.Addresses)

	// Fail-fast, as AddressService does: one unverifiable address is not a row to skip
	// past when the caller is choosing where funds go.
	rulesContainer, err := s.rulesCache.Get(ctx)
	if err != nil {
		return nil, err
	}
	if rulesContainer == nil {
		return nil, &model.IntegrityError{
			Message: "rules container required for address signature verification",
		}
	}
	if err := helper.VerifyAddressSignatures(addresses, rulesContainer); err != nil {
		return nil, err
	}

	result := &model.GetAssetAddressesResult{
		Addresses: addresses,
	}

	// Parse total items
	if resp.TotalItems != nil {
		if total, parseErr := strconv.ParseInt(*resp.TotalItems, 10, 64); parseErr == nil {
			result.TotalItems = total
		}
	}

	// Set next cursor
	if resp.Next != nil {
		result.NextCursor = *resp.Next
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

// GetAssetWallets retrieves wallet-level balances for a specific asset.
func (s *AssetService) GetAssetWallets(ctx context.Context, req *model.GetAssetWalletsRequest) (*model.GetAssetWalletsResult, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.Asset.Currency == "" {
		return nil, fmt.Errorf("asset currency is required")
	}

	// Build the OpenAPI request
	apiReq := openapi.TgvalidatordGetAssetWalletsRequest{
		Asset: mapper.AssetFilterToDTO(&req.Asset),
	}

	// Set optional fields
	if req.Limit > 0 {
		limit := fmt.Sprintf("%d", req.Limit)
		apiReq.Limit = &limit
	}
	if req.Cursor != "" {
		apiReq.Cursor = &req.Cursor
	}
	if req.WalletID != "" {
		apiReq.WalletId = &req.WalletID
	}
	if req.WalletName != "" {
		apiReq.WalletName = &req.WalletName
	}

	// Handle cursor-based pagination
	if req.CurrentPage != "" || req.PageRequest != "" || req.PageSize > 0 {
		cursor := openapi.TgvalidatordRequestCursor{}
		if req.CurrentPage != "" {
			cursor.CurrentPage = &req.CurrentPage
		}
		if req.PageRequest != "" {
			cursor.PageRequest = &req.PageRequest
		}
		if req.PageSize > 0 {
			pageSize := fmt.Sprintf("%d", req.PageSize)
			cursor.PageSize = &pageSize
		}
		apiReq.RequestCursor = &cursor
	}

	resp, httpResp, err := s.api.WalletServiceGetAssetWallets(ctx).Body(apiReq).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.GetAssetWalletsResult{
		Wallets: mapper.WalletsFromDTO(resp.Wallets),
	}

	// Parse total items
	if resp.TotalItems != nil {
		if total, parseErr := strconv.ParseInt(*resp.TotalItems, 10, 64); parseErr == nil {
			result.TotalItems = total
		}
	}

	// Set next cursor
	if resp.Next != nil {
		result.NextCursor = *resp.Next
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
