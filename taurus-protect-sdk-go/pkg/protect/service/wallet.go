// Package service provides high-level service wrappers for the Taurus-PROTECT API.
package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// WalletService provides wallet management operations.
type WalletService struct {
	api       *openapi.WalletsAPIService
	errMapper *ErrorMapper
}

// NewWalletService creates a new WalletService.
func NewWalletService(client *openapi.APIClient) *WalletService {
	return &WalletService{
		api:       client.WalletsAPI,
		errMapper: NewErrorMapper(),
	}
}

// GetWallet retrieves a wallet by ID.
func (s *WalletService) GetWallet(ctx context.Context, walletID string) (*model.Wallet, error) {
	if walletID == "" {
		return nil, fmt.Errorf("walletID cannot be empty")
	}

	resp, httpResp, err := s.api.WalletServiceGetWalletV2(ctx, walletID).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, fmt.Errorf("wallet not found")
	}

	return mapper.WalletFromDTO(resp.Result), nil
}

// ListWallets retrieves one page of wallets. The returned pagination is never nil; continue
// with its NextOffset until HasMore is false.
func (s *WalletService) ListWallets(ctx context.Context, opts *model.ListWalletsOptions) ([]*model.Wallet, *model.Pagination, error) {
	if opts == nil {
		opts = &model.ListWalletsOptions{}
	}
	window, err := resolveOffsetWindow(opts.Limit, opts.Offset)
	if err != nil {
		return nil, nil, err
	}

	// v2 takes the same request and reply as the deprecated GetWalletsInfo; only the
	// path differs.
	req := applyOffsetWindow(s.api.WalletServiceGetWalletsV2(ctx), window)
	if opts.Currency != "" {
		req = req.Currencies([]string{opts.Currency})
	}
	// Query matches seven columns (currency, customerid, blockchain, name,
	// container, accountpath, comment); Name matches the name alone.
	if opts.Query != "" {
		req = req.Query(opts.Query)
	}
	if opts.Name != "" {
		req = req.Name(opts.Name)
	}
	if opts.ExcludeDisabled {
		req = req.ExcludeDisabled(true)
	}
	if len(opts.IDs) > 0 {
		req = req.Ids(opts.IDs)
	}
	if opts.Blockchain != "" {
		req = req.Blockchain(opts.Blockchain)
	}
	if opts.Network != "" {
		req = req.Network(opts.Network)
	}
	if len(opts.TagIDs) > 0 {
		req = req.TagIDs(opts.TagIDs)
	}
	if opts.OnlyPositiveBalance {
		req = req.OnlyPositiveBalance(true)
	}
	if opts.SortOrder != "" {
		req = req.SortOrder(opts.SortOrder)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, nil, s.errMapper.MapError(err, httpResp)
	}

	// The reply offset is the NEXT page's offset, not the current one.
	pagination, err := offsetPagination(ruleReplyOffset, window, len(resp.Result), 0,
		offsetReply{TotalItems: resp.TotalItems, Offset: resp.Offset})
	if err != nil {
		return nil, nil, err
	}
	return mapper.WalletsFromDTO(resp.Result), pagination, nil
}

// CreateWallet creates a new wallet.
func (s *WalletService) CreateWallet(ctx context.Context, req *model.CreateWalletRequest) (*model.Wallet, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Currency == "" {
		return nil, fmt.Errorf("currency is required")
	}

	createReq := openapi.TgvalidatordCreateWalletRequest{
		Name:     req.Name,
		Currency: &req.Currency,
	}

	if req.Comment != "" {
		createReq.Comment = &req.Comment
	}
	if req.CustomerID != "" {
		createReq.CustomerId = &req.CustomerID
	}
	if req.ExternalWalletID != "" {
		createReq.ExternalWalletId = &req.ExternalWalletID
	}
	if req.VisibilityGroupID != "" {
		createReq.VisibilityGroupID = &req.VisibilityGroupID
	}

	resp, httpResp, err := s.api.WalletServiceCreateWallet(ctx).
		Body(createReq).
		Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, fmt.Errorf("failed to create wallet")
	}

	return mapper.WalletFromCreateDTO(resp.Result), nil
}

// CreateWalletAttribute creates a custom attribute on a wallet.
func (s *WalletService) CreateWalletAttribute(ctx context.Context, walletID, key, value string) error {
	if walletID == "" {
		return fmt.Errorf("walletID cannot be empty")
	}
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	attrReq := openapi.TgvalidatordCreateWalletAttributeRequest{}
	attrReq.SetKey(key)
	attrReq.SetValue(value)

	body := openapi.WalletServiceCreateWalletAttributesBody{
		Attributes: []openapi.TgvalidatordCreateWalletAttributeRequest{attrReq},
	}

	_, httpResp, err := s.api.WalletServiceCreateWalletAttributes(ctx, walletID).
		Body(body).
		Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// DeleteWalletAttribute deletes a custom attribute from a wallet.
func (s *WalletService) DeleteWalletAttribute(ctx context.Context, walletID, attributeID string) error {
	if walletID == "" {
		return fmt.Errorf("walletID cannot be empty")
	}
	if attributeID == "" {
		return fmt.Errorf("attributeID cannot be empty")
	}

	_, httpResp, err := s.api.WalletServiceDeleteWalletAttribute(ctx, walletID, attributeID).
		Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// GetWalletBalanceHistory retrieves balance history for a wallet.
// intervalHours specifies the interval between balance snapshots in hours.
func (s *WalletService) GetWalletBalanceHistory(ctx context.Context, walletID string, intervalHours int) ([]*model.BalanceHistoryPoint, error) {
	if walletID == "" {
		return nil, fmt.Errorf("walletID cannot be empty")
	}

	req := s.api.WalletServiceGetWalletBalanceHistory(ctx, walletID)
	if intervalHours > 0 {
		req = req.IntervalHours(fmt.Sprintf("%d", intervalHours))
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.BalanceHistoryPointsFromDTO(resp.Result), nil
}

// GetWalletTokens retrieves one page of a wallet's token balances. Continue with
// Page.NextCursor until Page.HasMore is false.
func (s *WalletService) GetWalletTokens(ctx context.Context, walletID string, opts *model.GetWalletTokensOptions) (*model.WalletTokensResult, error) {
	if walletID == "" {
		return nil, fmt.Errorf("walletID cannot be empty")
	}
	if opts == nil {
		opts = &model.GetWalletTokensOptions{}
	}
	pageSize, err := resolvePageSize("PageSize", opts.PageSize)
	if err != nil {
		return nil, err
	}

	req := s.api.WalletServiceGetWalletTokens(ctx, walletID).Limit(strconv.FormatInt(pageSize, 10))
	if opts.Cursor != "" {
		req = req.Cursor(opts.Cursor)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	page, err := cursorPage(pageSize, cursorReply{TokenOnly: true, Token: resp.Next, HasTotal: true, Total: resp.Total})
	if err != nil {
		return nil, err
	}
	return &model.WalletTokensResult{
		Tokens: mapper.AssetBalancesFromDTO(resp.Balances),
		Page:   page,
	}, nil
}
