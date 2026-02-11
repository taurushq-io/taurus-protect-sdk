package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// WhitelistedContractService provides whitelisted contract management operations.
// It wraps the ContractWhitelistingAPI for managing whitelisted smart contracts
// (tokens and NFTs) in the Taurus-PROTECT platform.
type WhitelistedContractService struct {
	api       *openapi.ContractWhitelistingAPIService
	errMapper *ErrorMapper
}

// NewWhitelistedContractService creates a new WhitelistedContractService.
func NewWhitelistedContractService(client *openapi.APIClient) *WhitelistedContractService {
	return &WhitelistedContractService{
		api:       client.ContractWhitelistingAPI,
		errMapper: NewErrorMapper(),
	}
}

// GetWhitelistedContract retrieves a whitelisted contract by ID.
func (s *WhitelistedContractService) GetWhitelistedContract(ctx context.Context, id string) (*model.WhitelistedContract, error) {
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

	return mapper.WhitelistedContractFromDTO(resp.Result), nil
}

// ListWhitelistedContracts retrieves a list of whitelisted contracts.
func (s *WhitelistedContractService) ListWhitelistedContracts(ctx context.Context, opts *model.ListWhitelistedContractsOptions) (*model.ListWhitelistedContractsResult, error) {
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
		return nil, s.errMapper.MapError(err, httpResp)
	}

	contracts := mapper.WhitelistedContractsFromDTO(resp.Result)

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

	return &model.ListWhitelistedContractsResult{
		Contracts:  contracts,
		Pagination: pagination,
	}, nil
}

// ListWhitelistedContractsForApproval retrieves contracts pending approval.
func (s *WhitelistedContractService) ListWhitelistedContractsForApproval(ctx context.Context, opts *model.ListWhitelistedContractsForApprovalOptions) (*model.ListWhitelistedContractsResult, error) {
	req := s.api.WhitelistServiceGetWhitelistedContractsForApproval(ctx)

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
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	contracts := mapper.WhitelistedContractsFromDTO(resp.Result)

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

	return &model.ListWhitelistedContractsResult{
		Contracts:  contracts,
		Pagination: pagination,
	}, nil
}

// CreateWhitelistedContract creates a new whitelisted contract.
func (s *WhitelistedContractService) CreateWhitelistedContract(ctx context.Context, req *model.CreateWhitelistedContractRequest) (string, error) {
	if req == nil {
		return "", fmt.Errorf("request cannot be nil")
	}
	if req.Blockchain == "" {
		return "", fmt.Errorf("blockchain is required")
	}
	if req.Symbol == "" {
		return "", fmt.Errorf("symbol is required")
	}

	createReq := openapi.TgvalidatordCreateWhitelistedContractAddressRequest{
		Blockchain: req.Blockchain,
		Symbol:     req.Symbol,
	}

	if req.ContractAddress != "" {
		createReq.ContractAddress = &req.ContractAddress
	}
	if req.Name != "" {
		createReq.Name = &req.Name
	}
	if req.Decimals != "" {
		createReq.Decimals = &req.Decimals
	}
	if req.TokenID != "" {
		createReq.TokenId = &req.TokenID
	}
	if req.Kind != "" {
		createReq.Kind = &req.Kind
	}
	if req.Network != "" {
		createReq.Network = &req.Network
	}

	resp, httpResp, err := s.api.WhitelistServiceCreateWhitelistedContract(ctx).
		Body(createReq).
		Execute()
	if err != nil {
		return "", s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil || resp.Result.Id == nil {
		return "", fmt.Errorf("failed to create whitelisted contract: no ID returned")
	}

	return *resp.Result.Id, nil
}

// UpdateWhitelistedContract updates an existing whitelisted contract.
// Note: Only symbol, name, and decimals can be updated.
// ALGO whitelisted contracts are not editable as they are populated with on-chain data at creation.
// Returns the temporary ID of the updated contract (approval may be required).
func (s *WhitelistedContractService) UpdateWhitelistedContract(ctx context.Context, id string, req *model.UpdateWhitelistedContractRequest) (string, error) {
	if id == "" {
		return "", fmt.Errorf("id cannot be empty")
	}
	if req == nil {
		return "", fmt.Errorf("request cannot be nil")
	}

	updateReq := openapi.WhitelistServiceUpdateWhitelistedContractBody{
		Symbol:   req.Symbol,
		Name:     req.Name,
		Decimals: req.Decimals,
	}

	resp, httpResp, err := s.api.WhitelistServiceUpdateWhitelistedContract(ctx, id).
		Body(updateReq).
		Execute()
	if err != nil {
		return "", s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil || resp.Result.UpdatedContractTemporaryId == nil {
		return "", nil
	}

	return *resp.Result.UpdatedContractTemporaryId, nil
}

// DeleteWhitelistedContract deletes a whitelisted contract.
// Note: This operation is deprecated due to complex dependencies.
func (s *WhitelistedContractService) DeleteWhitelistedContract(ctx context.Context, id string, comment string) (string, error) {
	if id == "" {
		return "", fmt.Errorf("id cannot be empty")
	}

	deleteReq := openapi.TgvalidatordDeleteWhitelistedContractAddressRequest{
		Id: id,
	}
	if comment != "" {
		deleteReq.Comment = &comment
	}

	resp, httpResp, err := s.api.WhitelistServiceDeleteWhitelistedContract(ctx).
		Body(deleteReq).
		Execute()
	if err != nil {
		return "", s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil || resp.Result.Id == nil {
		return "", nil
	}

	return *resp.Result.Id, nil
}

// ApproveWhitelistedContract approves whitelisted contracts.
// The signature is base64(ecdsa_sign(sha256([hex(sha256(req1_metadata)),hex(sha256(req2_metadata)),...]))).
func (s *WhitelistedContractService) ApproveWhitelistedContract(ctx context.Context, ids []string, signature string, comment string) error {
	if len(ids) == 0 {
		return fmt.Errorf("ids cannot be empty")
	}
	if signature == "" {
		return fmt.Errorf("signature is required")
	}
	if comment == "" {
		return fmt.Errorf("comment is required")
	}

	approveReq := openapi.TgvalidatordApproveWhitelistedContractAddressRequest{
		Ids:       ids,
		Signature: signature,
		Comment:   comment,
	}

	_, httpResp, err := s.api.WhitelistServiceApproveWhitelistedContract(ctx).
		Body(approveReq).
		Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// RejectWhitelistedContract rejects whitelisted contracts.
func (s *WhitelistedContractService) RejectWhitelistedContract(ctx context.Context, ids []string, comment string) error {
	if len(ids) == 0 {
		return fmt.Errorf("ids cannot be empty")
	}
	if comment == "" {
		return fmt.Errorf("comment is required")
	}

	rejectReq := openapi.TgvalidatordRejectWhitelistedContractAddressRequest{
		Ids:     ids,
		Comment: comment,
	}

	_, httpResp, err := s.api.WhitelistServiceRejectWhitelistedContract(ctx).
		Body(rejectReq).
		Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// CreateWhitelistedContractAttribute creates an attribute on a whitelisted contract.
func (s *WhitelistedContractService) CreateWhitelistedContractAttribute(ctx context.Context, contractID string, req *model.CreateWhitelistedContractAttributeRequest) ([]model.WhitelistedContractAttribute, error) {
	if contractID == "" {
		return nil, fmt.Errorf("contractID cannot be empty")
	}
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.Key == "" {
		return nil, fmt.Errorf("key is required")
	}

	attrReq := openapi.TgvalidatordCreateWhitelistedContractAddressAttributeRequest{}
	attrReq.Key = &req.Key
	attrReq.Value = &req.Value
	if req.ContentType != "" {
		attrReq.ContentType = &req.ContentType
	}
	if req.Type != "" {
		attrReq.Type = &req.Type
	}
	if req.Subtype != "" {
		attrReq.Subtype = &req.Subtype
	}
	if req.IsFile {
		attrReq.Isfile = &req.IsFile
	}

	createBody := openapi.WhitelistServiceCreateWhitelistedContractAttributesBody{
		Attributes: []openapi.TgvalidatordCreateWhitelistedContractAddressAttributeRequest{attrReq},
	}

	resp, httpResp, err := s.api.WhitelistServiceCreateWhitelistedContractAttributes(ctx, contractID).
		Body(createBody).
		Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, nil
	}

	return mapper.WhitelistedContractAttributesFromDTO(resp.Result), nil
}

// CreateWhitelistedContractAttributes creates multiple attributes on a whitelisted contract.
func (s *WhitelistedContractService) CreateWhitelistedContractAttributes(ctx context.Context, contractID string, reqs []model.CreateWhitelistedContractAttributeRequest) ([]model.WhitelistedContractAttribute, error) {
	if contractID == "" {
		return nil, fmt.Errorf("contractID cannot be empty")
	}
	if len(reqs) == 0 {
		return nil, fmt.Errorf("at least one attribute request is required")
	}

	attrReqs := make([]openapi.TgvalidatordCreateWhitelistedContractAddressAttributeRequest, len(reqs))
	for i, req := range reqs {
		if req.Key == "" {
			return nil, fmt.Errorf("key is required for attribute at index %d", i)
		}
		attrReq := openapi.TgvalidatordCreateWhitelistedContractAddressAttributeRequest{}
		attrReq.Key = &req.Key
		attrReq.Value = &req.Value
		if req.ContentType != "" {
			attrReq.ContentType = &req.ContentType
		}
		if req.Type != "" {
			attrReq.Type = &req.Type
		}
		if req.Subtype != "" {
			attrReq.Subtype = &req.Subtype
		}
		if req.IsFile {
			attrReq.Isfile = &req.IsFile
		}
		attrReqs[i] = attrReq
	}

	createBody := openapi.WhitelistServiceCreateWhitelistedContractAttributesBody{
		Attributes: attrReqs,
	}

	resp, httpResp, err := s.api.WhitelistServiceCreateWhitelistedContractAttributes(ctx, contractID).
		Body(createBody).
		Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, nil
	}

	return mapper.WhitelistedContractAttributesFromDTO(resp.Result), nil
}

// GetWhitelistedContractAttribute retrieves an attribute from a whitelisted contract.
func (s *WhitelistedContractService) GetWhitelistedContractAttribute(ctx context.Context, contractID string, attributeID string) (*model.WhitelistedContractAttribute, error) {
	if contractID == "" {
		return nil, fmt.Errorf("contractID cannot be empty")
	}
	if attributeID == "" {
		return nil, fmt.Errorf("attributeID cannot be empty")
	}

	resp, httpResp, err := s.api.WhitelistServiceGetWhitelistedContractAttribute(ctx, contractID, attributeID).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, nil
	}

	attr := mapper.WhitelistedContractAttributeFromDTO(resp.Result)
	return &attr, nil
}

// DeleteWhitelistedContractAttribute deletes an attribute from a whitelisted contract.
func (s *WhitelistedContractService) DeleteWhitelistedContractAttribute(ctx context.Context, contractID string, attributeID string) error {
	if contractID == "" {
		return fmt.Errorf("contractID cannot be empty")
	}
	if attributeID == "" {
		return fmt.Errorf("attributeID cannot be empty")
	}

	_, httpResp, err := s.api.WhitelistServiceDeleteWhitelistedContractAttribute(ctx, contractID, attributeID).Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}
