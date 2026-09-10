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

// ListWallets retrieves a list of wallets.
func (s *WalletService) ListWallets(ctx context.Context, opts *model.ListWalletsOptions) ([]*model.Wallet, *model.Pagination, error) {
	// v2 takes the same request and reply as the deprecated GetWalletsInfo; only the
	// path differs.
	req := s.api.WalletServiceGetWalletsV2(ctx)

	if opts != nil {
		if opts.Limit > 0 {
			req = req.Limit(fmt.Sprintf("%d", opts.Limit))
		}
		if opts.Offset > 0 {
			req = req.Offset(fmt.Sprintf("%d", opts.Offset))
		}
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
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, nil, s.errMapper.MapError(err, httpResp)
	}

	wallets := mapper.WalletsFromDTO(resp.Result)

	var pagination *model.Pagination
	if resp.TotalItems != nil || resp.Offset != nil {
		pagination = &model.Pagination{}
		if resp.TotalItems != nil {
			if total, err := strconv.ParseInt(*resp.TotalItems, 10, 64); err == nil {
				pagination.TotalItems = total
			}
		}
		if resp.Offset != nil {
			if offset, err := strconv.ParseInt(*resp.Offset, 10, 64); err == nil {
				pagination.Offset = offset
			}
		}
		if opts != nil {
			pagination.Limit = opts.Limit
		}
		pagination.HasMore = pagination.Offset+pagination.Limit < pagination.TotalItems
	}

	return wallets, pagination, nil
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

// GetWalletTokens retrieves token balances for a wallet.
// limit specifies the maximum number of tokens to return.
func (s *WalletService) GetWalletTokens(ctx context.Context, walletID string, opts *model.GetWalletTokensOptions) ([]*model.AssetBalance, error) {
	if walletID == "" {
		return nil, fmt.Errorf("walletID cannot be empty")
	}

	req := s.api.WalletServiceGetWalletTokens(ctx, walletID)
	if opts != nil {
		if opts.Limit > 0 {
			req = req.Limit(fmt.Sprintf("%d", opts.Limit))
		}
		if opts.Cursor != "" {
			req = req.Cursor(opts.Cursor)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.AssetBalancesFromDTO(resp.Balances), nil
}
