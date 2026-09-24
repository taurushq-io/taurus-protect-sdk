package service

import (
	"context"
	"fmt"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/cache"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/helper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// AssetService provides asset balance retrieval operations and the v2 asset registry reads.
type AssetService struct {
	api        *openapi.AssetsAPIService
	v2         *openapi.AssetV2APIService
	errMapper  *ErrorMapper
	rulesCache *cache.RulesContainerCache
	// The verified readers QueryAssetAddresses confirms its unsigned rows against.
	addresses   *AddressService
	whitelisted *WhitelistedAddressService
}

// NewAssetService creates a new AssetService with mandatory address verification.
//
// GetAssetAddresses returns the same Address entity AddressService does, signature and
// all, and used to hand it over unverified — so the mandatory verification on
// AddressService could be walked around by asking for the same rows here.
//
// QueryAssetAddresses rows carry no signature at all, so they are confirmed through the two
// verified readers passed in — the same instances the client hands out — rather than through a
// second copy of the HSM or whitelist verification.
//
// Panics if rulesCache, addresses or whitelisted is nil: address verification is mandatory.
func NewAssetService(
	client *openapi.APIClient,
	rulesCache *cache.RulesContainerCache,
	addresses *AddressService,
	whitelisted *WhitelistedAddressService,
) *AssetService {
	if rulesCache == nil {
		panic("rulesCache cannot be nil - address signature verification is mandatory")
	}
	if addresses == nil || whitelisted == nil {
		panic("the address and whitelisted-address services cannot be nil - asset holders are verified through them")
	}
	return &AssetService{
		api:         client.AssetsAPI,
		v2:          client.AssetV2API,
		errMapper:   NewErrorMapper(),
		rulesCache:  rulesCache,
		addresses:   addresses,
		whitelisted: whitelisted,
	}
}

// GetAssetAddresses retrieves one page of address-level balances for a specific asset, each
// address signature-verified. Continue with Page.NextCursor until Page.HasMore is false.
func (s *AssetService) GetAssetAddresses(ctx context.Context, req *model.GetAssetAddressesRequest) (*model.GetAssetAddressesResult, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.Asset.Currency == "" {
		return nil, fmt.Errorf("asset currency is required")
	}
	window, err := resolveCursorWindow(req.PageSize, req.Cursor, req.CurrentPage, req.PageRequest)
	if err != nil {
		return nil, err
	}

	// requestCursor only: the legacy limit + bytes cursor pair is not sent.
	apiReq := openapi.TgvalidatordGetAssetAddressesRequest{
		Asset:         mapper.AssetFilterToDTO(&req.Asset),
		RequestCursor: window.body(),
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

	resp, httpResp, err := s.api.WalletServiceGetAssetAddresses(ctx).Body(apiReq).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	page, err := cursorPage(window.pageSize, cursorReply{Cursor: resp.Cursor, HasTotal: true, Total: resp.TotalItems})
	if err != nil {
		return nil, err
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

	return &model.GetAssetAddressesResult{
		Addresses: addresses,
		Page:      page,
	}, nil
}

// GetAssetWallets retrieves one page of wallet-level balances for a specific asset. Continue
// with Page.NextCursor until Page.HasMore is false.
func (s *AssetService) GetAssetWallets(ctx context.Context, req *model.GetAssetWalletsRequest) (*model.GetAssetWalletsResult, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.Asset.Currency == "" {
		return nil, fmt.Errorf("asset currency is required")
	}
	window, err := resolveCursorWindow(req.PageSize, req.Cursor, req.CurrentPage, req.PageRequest)
	if err != nil {
		return nil, err
	}

	// requestCursor only: the legacy limit + bytes cursor pair is not sent.
	apiReq := openapi.TgvalidatordGetAssetWalletsRequest{
		Asset:         mapper.AssetFilterToDTO(&req.Asset),
		RequestCursor: window.body(),
	}
	if req.WalletID != "" {
		apiReq.WalletId = &req.WalletID
	}
	if req.WalletName != "" {
		apiReq.WalletName = &req.WalletName
	}

	resp, httpResp, err := s.api.WalletServiceGetAssetWallets(ctx).Body(apiReq).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	page, err := cursorPage(window.pageSize, cursorReply{Cursor: resp.Cursor, HasTotal: true, Total: resp.TotalItems})
	if err != nil {
		return nil, err
	}
	return &model.GetAssetWalletsResult{
		Wallets: mapper.WalletsFromDTO(resp.Wallets),
		Page:    page,
	}, nil
}

// QueryAssets retrieves one page of the v2 asset registry (AssetServiceV2_QueryAssetsV2, which
// supersedes GetAssetsV2). Continue with Page.NextCursor until Page.HasMore is false.
func (s *AssetService) QueryAssets(ctx context.Context, opts *model.QueryAssetsOptions) (*model.QueryAssetsResult, error) {
	if opts == nil {
		opts = &model.QueryAssetsOptions{}
	}
	window, err := resolveCursorWindow(opts.PageSize, opts.Cursor, "", "")
	if err != nil {
		return nil, err
	}

	body := openapi.TgvalidatordGetAssetsRequestV2{Cursor: window.body()}
	if opts.Blockchain != "" {
		body.Blockchain = &opts.Blockchain
	}
	if opts.Network != "" {
		body.Network = &opts.Network
	}
	if opts.Symbol != "" {
		body.Symbol = &opts.Symbol
	}
	if opts.ContractAddress != "" {
		body.ContractAddress = &opts.ContractAddress
	}
	if opts.Label != "" {
		body.Label = &opts.Label
	}
	if opts.CurrencyName != "" {
		body.CurrencyName = &opts.CurrencyName
	}

	resp, httpResp, err := s.v2.AssetServiceV2QueryAssetsV2(ctx).Body(body).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}
	page, err := cursorPage(window.pageSize, cursorReply{Cursor: resp.Cursor})
	if err != nil {
		return nil, err
	}
	return &model.QueryAssetsResult{
		Assets: mapper.AssetResourcesFromDTO(resp.Result),
		Page:   page,
	}, nil
}

// QueryAssetAddresses retrieves one page of the addresses holding a v2 asset, with their KYC
// status and balance. Continue with Page.NextCursor until Page.HasMore is false.
//
// The rows carry no signature, so each INTERNAL and WHITELISTED row is confirmed against its
// verified counterpart (see verifiedAssetAddresses) and marked Verified; a row that cannot be
// confirmed is withheld and named in ExcludedUnverified. Every other row is returned with
// Verified false.
func (s *AssetService) QueryAssetAddresses(ctx context.Context, assetID string, opts *model.QueryAssetAddressesOptions) (*model.QueryAssetAddressesResult, error) {
	if assetID == "" {
		return nil, fmt.Errorf("assetID cannot be empty")
	}
	if opts == nil {
		opts = &model.QueryAssetAddressesOptions{}
	}
	window, err := resolveCursorWindow(opts.PageSize, opts.Cursor, "", "")
	if err != nil {
		return nil, err
	}

	body := openapi.AssetServiceV2QueryAssetAddressesV2Body{Cursor: window.body()}
	if opts.AddressType != "" {
		addressType := openapi.TgvalidatordAddressTypeV2(opts.AddressType)
		body.AddressType = &addressType
	}
	if opts.KYCStatus != "" {
		kycStatus := openapi.TgvalidatordKYCStatusV2(opts.KYCStatus)
		body.KycStatus = &kycStatus
	}

	resp, httpResp, err := s.v2.AssetServiceV2QueryAssetAddressesV2(ctx, assetID).Body(body).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}
	// Exclusions below never move the cursor.
	page, err := cursorPage(window.pageSize, cursorReply{Cursor: resp.Cursor})
	if err != nil {
		return nil, err
	}
	addresses, excluded, err := s.verifiedAssetAddresses(ctx, mapper.AssetAddressesFromDTO(resp.Result))
	if err != nil {
		return nil, err
	}
	return &model.QueryAssetAddressesResult{
		Addresses:          addresses,
		ExcludedUnverified: excluded,
		Page:               page,
	}, nil
}

// verifiedAssetAddresses confirms a page of v2 asset holders through the verified readers.
//
//	INTERNAL    ─▶ addressID            ─▶ verified managed-address read (HSM, ≤ 50 ids) ──┐
//	WHITELISTED ─▶ whitelistedAddressID ─▶ verified whitelist read (6 steps, ≤ 100 ids) ───┤
//	                                          same address string? ── yes ─▶ Verified, address from the reader
//	                                                                └ no ──▶ excluded
//	any other type ─────────────────────────────────────────────────────▶ Verified = false
//
// A reader is called only when the page has a row of its type. Its call-level failures abort;
// its row-level ones become exclusions here. Rows came back but none survived is an error, so a
// filtered page never reads as an empty one.
func (s *AssetService) verifiedAssetAddresses(
	ctx context.Context,
	rows []*model.AssetAddress,
) ([]*model.AssetAddress, []model.ExcludedWhitelistedAddress, error) {
	internalType := string(openapi.TGVALIDATORDADDRESSTYPEV2_INTERNAL)
	whitelistedType := string(openapi.TGVALIDATORDADDRESSTYPEV2_WHITELISTED)

	var internalIDs, whitelistedIDs []string
	seen := make(map[string]bool)
	for _, row := range rows {
		switch {
		case row.AddressType == internalType && hasPlatformID(row.AddressID) && !seen["a"+row.AddressID]:
			seen["a"+row.AddressID] = true
			internalIDs = append(internalIDs, row.AddressID)
		case row.AddressType == whitelistedType && hasPlatformID(row.WhitelistedAddressID) && !seen["w"+row.WhitelistedAddressID]:
			seen["w"+row.WhitelistedAddressID] = true
			whitelistedIDs = append(whitelistedIDs, row.WhitelistedAddressID)
		}
	}

	var internal, whitelisted map[string]string
	var internalFailed, whitelistedFailed map[string]string
	if len(internalIDs) > 0 {
		verified, failed, err := s.addresses.verifiedAddressesByID(ctx, internalIDs)
		if err != nil {
			return nil, nil, err
		}
		internal, internalFailed = make(map[string]string, len(verified)), failed
		for id, address := range verified {
			internal[id] = address.Address
		}
	}
	if len(whitelistedIDs) > 0 {
		verified, failed, err := s.whitelisted.verifiedAddressesByID(ctx, whitelistedIDs, false)
		if err != nil {
			return nil, nil, err
		}
		whitelisted, whitelistedFailed = make(map[string]string, len(verified)), failed
		for id, address := range verified {
			whitelisted[id] = address.Address
		}
	}

	kept := make([]*model.AssetAddress, 0, len(rows))
	excluded := make([]model.ExcludedWhitelistedAddress, 0)
	for _, row := range rows {
		var id, verifiedAddress, reason string
		switch row.AddressType {
		case internalType:
			id, verifiedAddress, reason = confirmHolder(row.Address, row.AddressID, "addressID",
				internal, internalFailed, "managed address")
		case whitelistedType:
			id, verifiedAddress, reason = confirmHolder(row.Address, row.WhitelistedAddressID, "whitelistedAddressID",
				whitelisted, whitelistedFailed, "whitelisted address")
		default:
			// An EXTERNAL holder, or an untyped row: on-chain data nothing signs.
			row.Verified = false
			kept = append(kept, row)
			continue
		}
		if reason != "" {
			excluded = append(excluded, model.ExcludedWhitelistedAddress{ID: id, Reason: reason})
			continue
		}
		row.Address = verifiedAddress
		row.Verified = true
		kept = append(kept, row)
	}

	if len(rows) > 0 && len(kept) == 0 {
		return nil, nil, &model.IntegrityError{
			Message: fmt.Sprintf("all %d asset address(es) failed verification; first failure: %s",
				len(rows), excluded[0].Reason),
		}
	}
	return kept, excluded, nil
}

// confirmHolder checks one INTERNAL or WHITELISTED holder row against the addresses its verified
// reader returned, keyed by id. It returns the id an exclusion is filed under, the verified
// address, and why the row cannot be confirmed ("" when it can).
func confirmHolder(
	rowAddress, platformID, idField string,
	verified, failed map[string]string,
	kind string,
) (id, verifiedAddress, reason string) {
	if !hasPlatformID(platformID) {
		return rowAddress, "", fmt.Sprintf("the row carries no %s to verify it by", idField)
	}
	address, ok := verified[platformID]
	if !ok {
		if why, rejected := failed[platformID]; rejected {
			return platformID, "", fmt.Sprintf("the %s did not verify: %s", kind, why)
		}
		return platformID, "", fmt.Sprintf("the verified %s read did not return %s %s", kind, idField, platformID)
	}
	if address != rowAddress {
		return platformID, "", fmt.Sprintf("the row's address differs from the verified %s %s", kind, platformID)
	}
	return platformID, address, ""
}

// hasPlatformID reports whether a holder row carries a platform id; validatord sends zero (or
// omits the field) when the address is not managed or not whitelisted.
func hasPlatformID(id string) bool {
	return id != "" && id != "0"
}

// ListAssetOperations retrieves one page of the operations on a v2 asset. Continue with
// Page.NextCursor until Page.HasMore is false.
func (s *AssetService) ListAssetOperations(ctx context.Context, assetID string, opts *model.ListAssetOperationsOptions) (*model.ListAssetOperationsResult, error) {
	if assetID == "" {
		return nil, fmt.Errorf("assetID cannot be empty")
	}
	if opts == nil {
		opts = &model.ListAssetOperationsOptions{}
	}
	window, err := resolveCursorWindow(opts.PageSize, opts.Cursor, "", "")
	if err != nil {
		return nil, err
	}

	req := applyCursorQuery(s.v2.AssetServiceV2ListAssetOperationsV2(ctx, assetID), window)
	if opts.Type != "" {
		req = req.Type_(opts.Type)
	}
	if opts.Status != "" {
		req = req.Status(opts.Status)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}
	page, err := cursorPage(window.pageSize, cursorReply{Cursor: resp.Cursor})
	if err != nil {
		return nil, err
	}
	return &model.ListAssetOperationsResult{
		Operations: mapper.AssetOperationsFromDTO(resp.Result),
		Page:       page,
	}, nil
}
