package service

import (
	"context"
	"fmt"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// BusinessRuleService provides business rule management operations.
type BusinessRuleService struct {
	api       *openapi.BusinessRulesAPIService
	errMapper *ErrorMapper
}

// NewBusinessRuleService creates a new BusinessRuleService.
func NewBusinessRuleService(client *openapi.APIClient) *BusinessRuleService {
	return &BusinessRuleService{
		api:       client.BusinessRulesAPI,
		errMapper: NewErrorMapper(),
	}
}

// ListBusinessRules retrieves a list of business rules with optional filtering and pagination.
// Uses the v2 API with cursor-based pagination.
func (s *BusinessRuleService) ListBusinessRules(ctx context.Context, opts *model.ListBusinessRulesOptions) (*model.ListBusinessRulesResult, error) {
	req := s.api.RuleServiceGetBusinessRulesV2(ctx)

	if opts != nil {
		if len(opts.IDs) > 0 {
			req = req.Ids(opts.IDs)
		}
		if len(opts.RuleKeys) > 0 {
			req = req.RuleKeys(opts.RuleKeys)
		}
		if len(opts.RuleGroups) > 0 {
			req = req.RuleGroups(opts.RuleGroups)
		}
		if len(opts.WalletIDs) > 0 {
			req = req.WalletIds(opts.WalletIDs)
		}
		if len(opts.CurrencyIDs) > 0 {
			req = req.CurrencyIds(opts.CurrencyIDs)
		}
		if len(opts.AddressIDs) > 0 {
			req = req.AddressIds(opts.AddressIDs)
		}
		if opts.Level != "" {
			req = req.Level(opts.Level)
		}
		if opts.EntityType != "" {
			req = req.EntityType(opts.EntityType)
		}
		if len(opts.EntityIDs) > 0 {
			req = req.EntityIDs(opts.EntityIDs)
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

	result := &model.ListBusinessRulesResult{
		BusinessRules: mapper.BusinessRulesFromDTO(resp.Result),
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

// UpdateTransactionsEnabled toggles the transactions enabled business rule.
func (s *BusinessRuleService) UpdateTransactionsEnabled(ctx context.Context, req *model.UpdateTransactionsEnabledRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	body := openapi.TgvalidatordUpdateTransactionsEnabledBusinessRuleRequest{
		Enabled: &req.Enabled,
	}

	_, httpResp, err := s.api.RuleServiceUpdateTransactionsEnabledBusinessRule(ctx).
		Body(body).
		Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}
