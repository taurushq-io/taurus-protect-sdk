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

// AddressService provides address management operations.
// Address signature verification is mandatory - all GetAddress and ListAddresses
// calls verify signatures using the configured rules container cache.
type AddressService struct {
	api        *openapi.AddressesAPIService
	errMapper  *ErrorMapper
	rulesCache *cache.RulesContainerCache
}

// NewAddressService creates a new AddressService with mandatory signature verification.
// The rulesCache parameter is required and must not be nil.
// Address signature verification is performed on all GetAddress and ListAddresses calls.
//
// Panics if rulesCache is nil - address signature verification is mandatory for security.
func NewAddressService(client *openapi.APIClient, rulesCache *cache.RulesContainerCache) *AddressService {
	if rulesCache == nil {
		panic("rulesCache cannot be nil - address signature verification is mandatory")
	}
	return &AddressService{
		api:        client.AddressesAPI,
		errMapper:  NewErrorMapper(),
		rulesCache: rulesCache,
	}
}

// GetAddress retrieves an address by ID with mandatory signature verification.
// Returns an IntegrityError if signature verification fails.
func (s *AddressService) GetAddress(ctx context.Context, addressID string) (*model.Address, error) {
	if addressID == "" {
		return nil, fmt.Errorf("addressID cannot be empty")
	}

	resp, httpResp, err := s.api.WalletServiceGetAddress(ctx, addressID).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, fmt.Errorf("address not found")
	}

	address := mapper.AddressFromDTO(resp.Result)

	// Mandatory address signature verification
	rulesContainer, err := s.rulesCache.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get rules container for verification: %w", err)
	}
	if rulesContainer == nil {
		return nil, fmt.Errorf("rules container required for address signature verification")
	}
	if err := helper.VerifyAddressSignature(address, rulesContainer); err != nil {
		return nil, err
	}

	return address, nil
}

// ListAddresses retrieves a list of addresses with mandatory signature verification.
// Returns an IntegrityError if any address fails signature verification.
func (s *AddressService) ListAddresses(ctx context.Context, opts *model.ListAddressesOptions) ([]*model.Address, *model.Pagination, error) {
	req := s.api.WalletServiceGetAddresses(ctx)

	if opts != nil {
		if opts.WalletID != "" {
			req = req.WalletId(opts.WalletID)
		}
		if opts.Limit > 0 {
			req = req.Limit(fmt.Sprintf("%d", opts.Limit))
		}
		if opts.Offset > 0 {
			req = req.Offset(fmt.Sprintf("%d", opts.Offset))
		}
		if opts.Query != "" {
			req = req.Query(opts.Query)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, nil, s.errMapper.MapError(err, httpResp)
	}

	addresses := mapper.AddressesFromDTO(resp.Result)

	// Mandatory address signature verification
	rulesContainer, err := s.rulesCache.Get(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get rules container for verification: %w", err)
	}
	if rulesContainer == nil {
		return nil, nil, fmt.Errorf("rules container required for address signature verification")
	}
	if err := helper.VerifyAddressSignatures(addresses, rulesContainer); err != nil {
		return nil, nil, err
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
		}
		// Use overflow-safe comparison: check if there are more items beyond offset+limit
		// Instead of: offset+limit < totalItems (which can overflow)
		// We use: totalItems > offset && totalItems-offset > limit
		pagination.HasMore = pagination.TotalItems > pagination.Offset &&
			pagination.TotalItems-pagination.Offset > pagination.Limit
	}

	return addresses, pagination, nil
}

// CreateAddress creates a new address in a wallet.
func (s *AddressService) CreateAddress(ctx context.Context, req *model.CreateAddressRequest) (*model.Address, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.WalletID == "" {
		return nil, fmt.Errorf("walletID is required")
	}
	if req.Label == "" {
		return nil, fmt.Errorf("label is required")
	}

	createReq := openapi.TgvalidatordCreateAddressRequest{
		WalletId: req.WalletID,
		Label:    req.Label,
	}

	if req.Comment != "" {
		createReq.Comment = &req.Comment
	}
	if req.CustomerID != "" {
		createReq.CustomerId = &req.CustomerID
	}
	if req.ExternalAddressID != "" {
		createReq.ExternalAddressId = &req.ExternalAddressID
	}
	if req.Type != "" {
		createReq.Type = &req.Type
	}
	if req.NonHardenedDerivation {
		createReq.NonHardenedDerivation = &req.NonHardenedDerivation
	}

	resp, httpResp, err := s.api.WalletServiceCreateAddress(ctx).
		Body(createReq).
		Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, fmt.Errorf("failed to create address")
	}

	return mapper.AddressFromDTO(resp.Result), nil
}

// CreateAddressAttribute creates an attribute on an address.
func (s *AddressService) CreateAddressAttribute(ctx context.Context, addressID string, key string, value string) ([]model.AddressAttribute, error) {
	if addressID == "" {
		return nil, fmt.Errorf("addressID cannot be empty")
	}
	if key == "" {
		return nil, fmt.Errorf("key cannot be empty")
	}

	attrReq := openapi.TgvalidatordCreateAddressAttributeRequest{
		Key:   &key,
		Value: &value,
	}

	body := openapi.WalletServiceCreateAddressAttributesBody{
		Attributes: []openapi.TgvalidatordCreateAddressAttributeRequest{attrReq},
	}

	resp, httpResp, err := s.api.WalletServiceCreateAddressAttributes(ctx, addressID).
		Body(body).
		Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, nil
	}

	return mapper.AddressAttributesFromDTO(resp.Result), nil
}

// DeleteAddressAttribute deletes an attribute from an address.
func (s *AddressService) DeleteAddressAttribute(ctx context.Context, addressID string, attributeID string) error {
	if addressID == "" {
		return fmt.Errorf("addressID cannot be empty")
	}
	if attributeID == "" {
		return fmt.Errorf("attributeID cannot be empty")
	}

	_, httpResp, err := s.api.WalletServiceDeleteAddressAttribute(ctx, addressID, attributeID).Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// GetAddressProofOfReserve retrieves the proof of reserve for an address.
func (s *AddressService) GetAddressProofOfReserve(ctx context.Context, addressID string, challenge string) (*model.ProofOfReserve, error) {
	if addressID == "" {
		return nil, fmt.Errorf("addressID cannot be empty")
	}

	req := s.api.WalletServiceGetAddressProofOfReserve(ctx, addressID)
	if challenge != "" {
		req = req.Challenge(challenge)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, nil
	}

	return mapper.ProofOfReserveFromDTO(resp.Result), nil
}
