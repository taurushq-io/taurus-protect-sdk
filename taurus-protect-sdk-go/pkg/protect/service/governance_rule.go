package service

import (
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/helper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// GovernanceRuleService provides governance rules management operations.
type GovernanceRuleService struct {
	api                *openapi.GovernanceRulesAPIService
	errMapper          *ErrorMapper
	superAdminKeys     []*ecdsa.PublicKey
	minValidSignatures int
}

// GovernanceRuleServiceConfig holds configuration for signature verification.
type GovernanceRuleServiceConfig struct {
	// SuperAdminKeys are the public keys used to verify governance rules signatures.
	SuperAdminKeys []*ecdsa.PublicKey
	// MinValidSignatures is the minimum number of valid SuperAdmin signatures required.
	MinValidSignatures int
}

// NewGovernanceRuleService creates a new GovernanceRuleService without verification.
func NewGovernanceRuleService(client *openapi.APIClient) *GovernanceRuleService {
	return &GovernanceRuleService{
		api:       client.GovernanceRulesAPI,
		errMapper: NewErrorMapper(),
	}
}

// NewGovernanceRuleServiceWithVerification creates a new GovernanceRuleService with
// signature verification enabled. When SuperAdmin keys are provided, GetDecodedRulesContainer
// will verify signatures before returning the decoded container.
func NewGovernanceRuleServiceWithVerification(
	client *openapi.APIClient,
	config *GovernanceRuleServiceConfig,
) *GovernanceRuleService {
	svc := &GovernanceRuleService{
		api:       client.GovernanceRulesAPI,
		errMapper: NewErrorMapper(),
	}

	if config != nil {
		svc.superAdminKeys = config.SuperAdminKeys
		svc.minValidSignatures = config.MinValidSignatures
	}

	return svc
}

// SuperAdminKeys returns the configured SuperAdmin public keys.
func (s *GovernanceRuleService) SuperAdminKeys() []*ecdsa.PublicKey {
	return s.superAdminKeys
}

// MinValidSignatures returns the minimum number of valid signatures required.
func (s *GovernanceRuleService) MinValidSignatures() int {
	return s.minValidSignatures
}

// GetRules retrieves the currently enforced governance rules.
func (s *GovernanceRuleService) GetRules(ctx context.Context) (*model.GovernanceRuleset, error) {
	resp, httpResp, err := s.api.RuleServiceGetRules(ctx).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, nil
	}

	return mapper.GovernanceRulesetFromDTO(resp.Result), nil
}

// GetRulesByID retrieves a governance ruleset by its ID.
func (s *GovernanceRuleService) GetRulesByID(ctx context.Context, id string) (*model.GovernanceRuleset, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	resp, httpResp, err := s.api.RuleServiceGetRulesByID(ctx, id).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, nil
	}

	return mapper.GovernanceRulesetFromDTO(resp.Result), nil
}

// GetRulesHistory retrieves the history of governance rules with pagination.
func (s *GovernanceRuleService) GetRulesHistory(ctx context.Context, opts *model.ListRulesHistoryOptions) (*model.GovernanceRulesHistoryResult, error) {
	req := s.api.RuleServiceGetRulesHistory(ctx)

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

	result := &model.GovernanceRulesHistoryResult{
		Rules: mapper.GovernanceRulesetsFromDTO(resp.Result),
	}

	if resp.TotalItems != nil {
		if total, parseErr := strconv.ParseInt(*resp.TotalItems, 10, 64); parseErr == nil {
			result.TotalItems = total
		}
	}

	if resp.Cursor != nil {
		result.Cursor = *resp.Cursor
	}

	return result, nil
}

// GetRulesProposal retrieves the proposed governance rules.
// Requires SuperAdmin or SuperAdminReadOnly role.
func (s *GovernanceRuleService) GetRulesProposal(ctx context.Context) (*model.GovernanceRuleset, error) {
	resp, httpResp, err := s.api.RuleServiceGetRulesProposal(ctx).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, nil
	}

	return mapper.GovernanceRulesetFromDTO(resp.Result), nil
}

// GetPublicKeys retrieves the list of SuperAdmin public keys.
func (s *GovernanceRuleService) GetPublicKeys(ctx context.Context) ([]*model.SuperAdminPublicKey, error) {
	resp, httpResp, err := s.api.RuleServiceGetPublicKeys(ctx).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.SuperAdminPublicKeysFromDTO(resp.PublicKeys), nil
}

// GetDecodedRulesContainer decodes and verifies a GovernanceRuleset's rules container.
// If SuperAdmin keys are configured, it verifies the signatures before returning.
// Returns the decoded rules container or an error if verification fails.
func (s *GovernanceRuleService) GetDecodedRulesContainer(
	rules *model.GovernanceRuleset,
) (*model.DecodedRulesContainer, error) {
	if rules == nil {
		return nil, fmt.Errorf("governance rules cannot be nil")
	}

	if rules.RulesContainer == "" {
		return nil, fmt.Errorf("rules container is empty")
	}

	// Verify signatures if SuperAdmin keys are configured
	if len(s.superAdminKeys) > 0 {
		if err := s.VerifyGovernanceRules(rules); err != nil {
			return nil, err
		}
	}

	// Decode the rules container
	return mapper.RulesContainerFromBase64(rules.RulesContainer)
}

// VerifyGovernanceRules verifies the SuperAdmin signatures on the governance rules.
func (s *GovernanceRuleService) VerifyGovernanceRules(rules *model.GovernanceRuleset) error {
	if len(rules.Signatures) == 0 {
		return &model.IntegrityError{Message: "no signatures provided for governance rules"}
	}

	// Decode the rules container data
	rulesData, err := base64.StdEncoding.DecodeString(rules.RulesContainer)
	if err != nil {
		return &model.IntegrityError{
			Message: fmt.Sprintf("failed to decode rules container: %v", err),
		}
	}

	// Convert signatures to the format expected by the verifier
	signatures := make([]*model.RuleUserSignature, len(rules.Signatures))
	for i := range rules.Signatures {
		signatures[i] = &model.RuleUserSignature{
			UserID:    rules.Signatures[i].UserID,
			Signature: rules.Signatures[i].Signature,
		}
	}

	// Verify signatures
	if err := helper.VerifyGovernanceRulesSignatures(
		rulesData,
		signatures,
		s.superAdminKeys,
		s.minValidSignatures,
	); err != nil {
		return &model.IntegrityError{
			Message: fmt.Sprintf("governance rules signature verification failed: %v", err),
		}
	}

	return nil
}
