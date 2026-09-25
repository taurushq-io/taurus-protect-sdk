package service

import (
	"context"
	"errors"
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

	return s.verifiedAddress(ctx, resp.Result)
}

// verifiedAddress is the construction seam for a single *model.Address: GetAddress and
// CreateAddress both go through it, because "remember to verify" was a rule rather than the
// only available construction path, and CreateAddress is what that cost — GetAssetAddresses
// had already been fixed for exactly this, and the create path was missed. Same reasoning as
// RequestService.verifiedRequest.
//
// The two PAGE paths — ListAddresses and AssetService.GetAssetAddresses — verify through
// helper.VerifyAddressSignatures instead, so this is not literally the only route to an
// Address. That is a recorded divergence, not an oversight: the batch verifier is STRICTER
// than this seam (it errors on an empty address string, where this seam returns it), so the
// invariant "never return a non-empty address that has not been verified" holds on all four
// paths — but a page containing an address still being created fails here and succeeds in
// Java, whose list path shares the seam. See TODOS.md; do not "align" it by loosening the
// batch verifier without deciding which behaviour is wanted.
//
// Asynchronous creation is the one case that is NOT simply verify-and-throw. The DTO carries a
// `status` of `created`/`creating`/`signed`/`observed`/`confirmed`, so a reply can legitimately
// arrive before the HSM has signed the address. The rule that holds either way: never hand back
// a non-empty Address.Address that has not been verified. So a signature present must verify,
// and a signature absent means the address string is withheld rather than returned unchecked —
// the caller re-reads through this same seam once the status advances.
func (s *AddressService) verifiedAddress(
	ctx context.Context,
	dto *openapi.TgvalidatordAddress,
) (*model.Address, error) {
	address := mapper.AddressFromDTO(dto)
	if address == nil {
		return nil, fmt.Errorf("address not found")
	}

	if address.Address == "" {
		// Nothing to verify and nothing to misuse: an address that has not been generated
		// yet carries no destination. Status tells the caller to come back.
		return address, nil
	}

	if address.Signature == "" {
		// A non-empty address with no signature is NOT a silent skip. Returning it would
		// hand the caller an attacker-controllable destination in the same type as a
		// verified one, which is the whole defect.
		return nil, &model.IntegrityError{
			Message: fmt.Sprintf("address %s (status %q) carries an address string but no HSM "+
				"signature; refusing to return an unverified destination. Re-read once the "+
				"status reaches \"signed\"", address.ID, address.Status),
		}
	}

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

// maxAddressIDsPerRequest is validatord's cap on addressIds in one GetAddresses request.
const maxAddressIDsPerRequest = 50

// verifiedAddressesByID re-reads managed addresses by id, at most maxAddressIDsPerRequest ids
// per request, and runs every returned row through verifiedAddress. It returns the rows that
// verified keyed by id, plus why each other returned row did not.
//
// A row whose signature is missing or does not verify is reported, not fatal — the caller
// decides what to do without it. A rules container that cannot verify any address (none, or no
// HSMSLOT key) is not a property of one row, so it aborts the call, as does any request error.
func (s *AddressService) verifiedAddressesByID(ctx context.Context, ids []string) (map[string]*model.Address, map[string]string, error) {
	rulesContainer, err := s.rulesCache.Get(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get rules container for verification: %w", err)
	}
	if rulesContainer == nil || rulesContainer.GetHsmPublicKey() == nil {
		return nil, nil, &model.IntegrityError{
			Message: "the rules container has no HSMSLOT key, so no address can be verified",
		}
	}

	verified := make(map[string]*model.Address, len(ids))
	failed := make(map[string]string)
	for _, batch := range chunkIDs(ids, maxAddressIDsPerRequest) {
		resp, httpResp, err := s.api.WalletServiceGetAddresses(ctx).
			AddressIds(batch).
			Limit(strconv.Itoa(len(batch))).
			Execute()
		if err != nil {
			return nil, nil, s.errMapper.MapError(err, httpResp)
		}
		for i := range resp.Result {
			address, err := s.verifiedAddress(ctx, &resp.Result[i])
			if err != nil {
				var integrity *model.IntegrityError
				if !errors.As(err, &integrity) {
					return nil, nil, err
				}
				failed[safeString(resp.Result[i].Id)] = err.Error()
				continue
			}
			verified[address.ID] = address
		}
	}
	return verified, failed, nil
}

// ListAddresses retrieves one page of addresses with mandatory signature verification.
// Returns an IntegrityError if any address fails signature verification. The returned
// pagination is never nil; continue with its NextOffset until HasMore is false.
func (s *AddressService) ListAddresses(ctx context.Context, opts *model.ListAddressesOptions) ([]*model.Address, *model.Pagination, error) {
	if opts == nil {
		opts = &model.ListAddressesOptions{}
	}
	window, err := resolveOffsetWindow(opts.Limit, opts.Offset)
	if err != nil {
		return nil, nil, err
	}

	req := applyOffsetWindow(s.api.WalletServiceGetAddresses(ctx), window)
	if opts.WalletID != "" {
		req = req.WalletId(opts.WalletID)
	}
	// Query matches address, alternate address, comment, label and customer id.
	if opts.Query != "" {
		req = req.Query(opts.Query)
	}
	if opts.ExcludeDisabled {
		req = req.IncludeDisabledAddresses("exclude")
	}
	if len(opts.AddressIDs) > 0 {
		req = req.AddressIds(opts.AddressIDs)
	}
	if len(opts.Addresses) > 0 {
		req = req.Addresses(opts.Addresses)
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
	if opts.BalanceAbove != "" {
		req = req.BalanceAbove(opts.BalanceAbove)
	}
	if opts.BalanceBelow != "" {
		req = req.BalanceBelow(opts.BalanceBelow)
	}
	if opts.SortBy != "" {
		req = req.SortBy(opts.SortBy)
	}
	if opts.SortOrder != "" {
		req = req.SortOrder(opts.SortOrder)
	}
	req = applyAddressScoreFilter(req, opts.Score)

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

	// Through the same seam as every read. The create reply carries the same
	// TgvalidatordAddress the read paths return, signature included, so there was never a
	// reason for this path to be the unverified one — and callers of create are precisely
	// the ones about to publish or fund a fresh deposit address.
	return s.verifiedAddress(ctx, resp.Result)
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

// applyAddressScoreFilter wires the scoreFilter group. The flat scoreProvider /
// scoreInBelow / coinfirmScoreGreater parameters the client also exposes are
// deprecated in favour of these.
//
// A nil filter sets nothing: a zero-valued struct would send empty provider
// parameters on every address call.
func applyAddressScoreFilter(
	req openapi.ApiWalletServiceGetAddressesRequest,
	f *model.AddressScoreFilter,
) openapi.ApiWalletServiceGetAddressesRequest {
	if f == nil {
		return req
	}
	if f.Provider != "" {
		req = req.ScoreFilterScoreProvider(f.Provider)
	}
	if f.ScorechainInBelow != "" {
		req = req.ScoreFilterScorechainFiltersScoreInBelow(f.ScorechainInBelow)
	}
	if f.ScorechainOutBelow != "" {
		req = req.ScoreFilterScorechainFiltersScoreOutBelow(f.ScorechainOutBelow)
	}
	if f.ScorechainExclusive {
		req = req.ScoreFilterScorechainFiltersScoreExclusive(true)
	}
	if f.CoinfirmScoreGreater != "" {
		req = req.ScoreFilterCoinfirmFiltersScoreGreater(f.CoinfirmScoreGreater)
	}
	if f.ChainalysisScoreGreater != "" {
		req = req.ScoreFilterChainalysisFiltersScoreGreater(f.ChainalysisScoreGreater)
	}
	if f.EllipticScoreGreater != "" {
		req = req.ScoreFilterEllipticFiltersScoreGreater(f.EllipticScoreGreater)
	}
	if f.TRMLabsScoreGreater != "" {
		req = req.ScoreFilterTrmlabsFiltersScoreGreater(f.TRMLabsScoreGreater)
	}
	return req
}
