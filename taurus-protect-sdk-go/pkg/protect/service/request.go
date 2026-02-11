package service

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/helper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// RequestService provides transaction request management operations.
type RequestService struct {
	api       *openapi.RequestsAPIService
	errMapper *ErrorMapper
}

// NewRequestService creates a new RequestService.
func NewRequestService(client *openapi.APIClient) *RequestService {
	return &RequestService{
		api:       client.RequestsAPI,
		errMapper: NewErrorMapper(),
	}
}

// GetRequest retrieves a request by ID with hash verification.
// Returns an IntegrityError if the computed hash doesn't match the provided hash.
func (s *RequestService) GetRequest(ctx context.Context, requestID string) (*model.Request, error) {
	if requestID == "" {
		return nil, fmt.Errorf("requestID cannot be empty")
	}

	resp, httpResp, err := s.api.RequestServiceGetRequest(ctx, requestID).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, fmt.Errorf("request not found")
	}

	r := mapper.RequestFromDTO(resp.Result)

	if err := verifyRequestHash(r); err != nil {
		return nil, err
	}

	return r, nil
}

// verifyRequestHash verifies the integrity of a request's metadata hash.
// Returns an IntegrityError if the hash is empty or doesn't match.
func verifyRequestHash(r *model.Request) error {
	if r.Metadata == nil || (r.Metadata.Hash == "" && r.Metadata.PayloadAsString == "") {
		return nil
	}

	computedHash := crypto.CalculateHexHash(r.Metadata.PayloadAsString)
	providedHash := r.Metadata.Hash
	if computedHash == "" || providedHash == "" {
		return &model.IntegrityError{
			Message: "request hash verification failed: hash values must be non-empty",
		}
	}
	if !helper.ConstantTimeCompare(computedHash, providedHash) {
		return &model.IntegrityError{
			Message: fmt.Sprintf("request hash verification failed: computed=%s, provided=%s", computedHash, providedHash),
		}
	}
	return nil
}

// ListRequests retrieves a list of requests using cursor-based pagination.
func (s *RequestService) ListRequests(ctx context.Context, opts *model.ListRequestsOptions) (*model.RequestResult, error) {
	req := s.api.RequestServiceGetRequestsV2(ctx)

	if opts != nil {
		if opts.PageSize > 0 {
			req = req.CursorPageSize(fmt.Sprintf("%d", opts.PageSize))
		}
		if opts.Cursor != "" {
			req = req.CursorCurrentPage(opts.Cursor)
			req = req.CursorPageRequest("NEXT")
		}
		if opts.Status != "" {
			req = req.Statuses([]string{opts.Status})
		}
		if opts.Currency != "" {
			req = req.CurrencyID(opts.Currency)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.RequestResult{
		Requests: mapper.RequestsFromDTO(resp.Result),
	}

	if resp.Cursor != nil {
		if resp.Cursor.CurrentPage != nil {
			result.NextCursor = *resp.Cursor.CurrentPage
		}
		if resp.Cursor.HasNext != nil {
			result.HasNext = *resp.Cursor.HasNext
		}
	}

	return result, nil
}

// ListRequestsForApproval retrieves requests pending approval for the current user.
func (s *RequestService) ListRequestsForApproval(ctx context.Context, opts *model.ListRequestsOptions) (*model.RequestResult, error) {
	req := s.api.RequestServiceGetRequestsForApprovalV2(ctx)

	if opts != nil {
		if opts.PageSize > 0 {
			req = req.CursorPageSize(fmt.Sprintf("%d", opts.PageSize))
		}
		if opts.Cursor != "" {
			req = req.CursorCurrentPage(opts.Cursor)
			req = req.CursorPageRequest("NEXT")
		}
		if opts.Currency != "" {
			req = req.CurrencyID(opts.Currency)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.RequestResult{
		Requests: mapper.RequestsFromDTO(resp.Result),
	}

	if resp.Cursor != nil {
		if resp.Cursor.CurrentPage != nil {
			result.NextCursor = *resp.Cursor.CurrentPage
		}
		if resp.Cursor.HasNext != nil {
			result.HasNext = *resp.Cursor.HasNext
		}
	}

	return result, nil
}

// CreateOutgoingRequest creates a new outgoing (withdrawal) request.
func (s *RequestService) CreateOutgoingRequest(ctx context.Context, req *model.CreateOutgoingRequest) (*model.Request, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.Amount == "" {
		return nil, fmt.Errorf("amount is required")
	}
	if req.FromAddressID == "" && req.FromWalletID == "" {
		return nil, fmt.Errorf("either fromAddressID or fromWalletID is required")
	}
	if req.ToAddressID == "" && req.ToWhitelistedAddressID == "" {
		return nil, fmt.Errorf("either toAddressID or toWhitelistedAddressID is required")
	}

	createReq := openapi.TgvalidatordCreateOutgoingRequestRequest{
		Amount: req.Amount,
	}

	if req.FromAddressID != "" {
		createReq.FromAddressId = &req.FromAddressID
	}
	if req.FromWalletID != "" {
		createReq.FromWalletId = &req.FromWalletID
	}
	if req.ToAddressID != "" {
		createReq.ToAddressId = &req.ToAddressID
	}
	if req.ToWhitelistedAddressID != "" {
		createReq.ToWhitelistedAddressId = &req.ToWhitelistedAddressID
	}
	if req.FeeLimit != "" {
		createReq.FeeLimit = &req.FeeLimit
	}
	if req.GasLimit != "" {
		createReq.GasLimit = &req.GasLimit
	}
	if req.Comment != "" {
		createReq.Comment = &req.Comment
	}
	if req.TransactionComment != "" {
		createReq.TransactionComment = &req.TransactionComment
	}
	if req.ExternalRequestID != "" {
		createReq.ExternalRequestId = &req.ExternalRequestID
	}
	if req.UseUnconfirmedFunds {
		createReq.UseUnconfirmedFunds = &req.UseUnconfirmedFunds
	}
	if req.FeePaidByReceiver {
		createReq.FeePaidByReceiver = &req.FeePaidByReceiver
	}
	if req.UseAllFunds {
		createReq.UseAllFunds = &req.UseAllFunds
	}

	resp, httpResp, err := s.api.RequestServiceCreateOutgoingRequest(ctx).
		Body(createReq).
		Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, fmt.Errorf("failed to create request")
	}

	return mapper.RequestFromDTO(resp.Result), nil
}

// CreateInternalTransferRequest creates an internal transfer request from one address to another.
// This is a convenience method that calls CreateOutgoingRequest with the appropriate parameters.
func (s *RequestService) CreateInternalTransferRequest(ctx context.Context, fromAddressID, toAddressID, amount string) (*model.Request, error) {
	if fromAddressID == "" {
		return nil, fmt.Errorf("fromAddressID cannot be empty")
	}
	if toAddressID == "" {
		return nil, fmt.Errorf("toAddressID cannot be empty")
	}
	if amount == "" {
		return nil, fmt.Errorf("amount cannot be empty")
	}

	return s.CreateOutgoingRequest(ctx, &model.CreateOutgoingRequest{
		FromAddressID: fromAddressID,
		ToAddressID:   toAddressID,
		Amount:        amount,
	})
}

// CreateInternalTransferFromWalletRequest creates an internal transfer from an omnibus wallet.
func (s *RequestService) CreateInternalTransferFromWalletRequest(ctx context.Context, fromWalletID, toAddressID, amount string) (*model.Request, error) {
	if fromWalletID == "" {
		return nil, fmt.Errorf("fromWalletID cannot be empty")
	}
	if toAddressID == "" {
		return nil, fmt.Errorf("toAddressID cannot be empty")
	}
	if amount == "" {
		return nil, fmt.Errorf("amount cannot be empty")
	}

	return s.CreateOutgoingRequest(ctx, &model.CreateOutgoingRequest{
		FromWalletID: fromWalletID,
		ToAddressID:  toAddressID,
		Amount:       amount,
	})
}

// CreateExternalTransferRequest creates an external transfer to a whitelisted address.
func (s *RequestService) CreateExternalTransferRequest(ctx context.Context, fromAddressID, toWhitelistedAddressID, amount string) (*model.Request, error) {
	if fromAddressID == "" {
		return nil, fmt.Errorf("fromAddressID cannot be empty")
	}
	if toWhitelistedAddressID == "" {
		return nil, fmt.Errorf("toWhitelistedAddressID cannot be empty")
	}
	if amount == "" {
		return nil, fmt.Errorf("amount cannot be empty")
	}

	return s.CreateOutgoingRequest(ctx, &model.CreateOutgoingRequest{
		FromAddressID:          fromAddressID,
		ToWhitelistedAddressID: toWhitelistedAddressID,
		Amount:                 amount,
	})
}

// CreateExternalTransferFromWalletRequest creates an external transfer from an omnibus wallet.
func (s *RequestService) CreateExternalTransferFromWalletRequest(ctx context.Context, fromWalletID, toWhitelistedAddressID, amount string) (*model.Request, error) {
	if fromWalletID == "" {
		return nil, fmt.Errorf("fromWalletID cannot be empty")
	}
	if toWhitelistedAddressID == "" {
		return nil, fmt.Errorf("toWhitelistedAddressID cannot be empty")
	}
	if amount == "" {
		return nil, fmt.Errorf("amount cannot be empty")
	}

	return s.CreateOutgoingRequest(ctx, &model.CreateOutgoingRequest{
		FromWalletID:           fromWalletID,
		ToWhitelistedAddressID: toWhitelistedAddressID,
		Amount:                 amount,
	})
}

// CreateCancelRequest creates a cancel request for a pending transaction.
func (s *RequestService) CreateCancelRequest(ctx context.Context, addressID string, nonce string) (*model.Request, error) {
	if addressID == "" {
		return nil, fmt.Errorf("addressID cannot be empty")
	}

	createReq := openapi.TgvalidatordCreateOutgoingCancelRequestRequest{
		AddressId: addressID,
	}
	if nonce != "" {
		createReq.Nonce = &nonce
	}

	resp, httpResp, err := s.api.RequestServiceCreateOutgoingCancelRequest(ctx).
		Body(createReq).
		Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, fmt.Errorf("failed to create cancel request")
	}

	return mapper.RequestFromDTO(resp.Result), nil
}

// CreateIncomingRequest creates an incoming request to log an incoming transaction from an exchange.
func (s *RequestService) CreateIncomingRequest(ctx context.Context, req *model.CreateIncomingRequest) (*model.Request, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.FromExchangeID == "" {
		return nil, fmt.Errorf("fromExchangeID is required")
	}
	if req.ToAddressID == "" {
		return nil, fmt.Errorf("toAddressID is required")
	}
	if req.Amount == "" {
		return nil, fmt.Errorf("amount is required")
	}

	createReq := openapi.TgvalidatordCreateIncomingRequestRequest{
		Amount:         req.Amount,
		FromExchangeId: req.FromExchangeID,
		ToAddressId:    req.ToAddressID,
	}

	if req.Comment != "" {
		createReq.Comment = &req.Comment
	}
	if req.ExternalRequestID != "" {
		createReq.ExternalRequestId = &req.ExternalRequestID
	}

	resp, httpResp, err := s.api.RequestServiceCreateIncomingRequest(ctx).
		Body(createReq).
		Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, fmt.Errorf("failed to create incoming request")
	}

	return mapper.RequestFromDTO(resp.Result), nil
}

// ApproveRequests approves multiple requests using a private key for signing.
// The requests are sorted by ID before signing. Returns the number of requests signed.
func (s *RequestService) ApproveRequests(ctx context.Context, requests []*model.Request, privateKey *ecdsa.PrivateKey) (int, error) {
	if len(requests) == 0 {
		return 0, fmt.Errorf("requests list cannot be empty")
	}
	if privateKey == nil {
		return 0, fmt.Errorf("privateKey cannot be nil")
	}

	// Validate all requests have metadata with hash and valid numeric IDs
	for _, r := range requests {
		if r.Metadata == nil || r.Metadata.Hash == "" {
			return 0, fmt.Errorf("request %s has no metadata hash", r.ID)
		}
		// Validate ID is a valid numeric value for sorting
		if _, err := strconv.ParseInt(r.ID, 10, 64); err != nil {
			return 0, fmt.Errorf("request ID %q is not a valid numeric ID: %w", r.ID, err)
		}
	}

	// Sort requests by ID (numeric sort) - IDs already validated above
	sortedRequests := make([]*model.Request, len(requests))
	copy(sortedRequests, requests)
	sort.Slice(sortedRequests, func(i, j int) bool {
		idI, _ := strconv.ParseInt(sortedRequests[i].ID, 10, 64)
		idJ, _ := strconv.ParseInt(sortedRequests[j].ID, 10, 64)
		return idI < idJ
	})

	// Build JSON array of hashes
	hashes := make([]string, len(sortedRequests))
	for i, r := range sortedRequests {
		hashes[i] = r.Metadata.Hash
	}

	hashesJSON, err := json.Marshal(hashes)
	if err != nil {
		return 0, fmt.Errorf("failed to serialize hashes: %w", err)
	}

	// Sign with ECDSA
	signature, err := crypto.SignData(privateKey, hashesJSON)
	if err != nil {
		return 0, fmt.Errorf("failed to sign request hashes: %w", err)
	}

	// Build request IDs
	ids := make([]string, len(sortedRequests))
	for i, r := range sortedRequests {
		ids[i] = r.ID
	}

	// Submit approval
	approveReq := openapi.TgvalidatordApproveRequestsRequest{
		Signature: signature,
		Comment:   "approved via taurus-protect-sdk-go",
		Ids:       ids,
	}

	resp, httpResp, err := s.api.RequestServiceApproveRequests(ctx).
		Body(approveReq).
		Execute()
	if err != nil {
		return 0, s.errMapper.MapError(err, httpResp)
	}

	if resp.SignedRequests != nil {
		signed, _ := strconv.Atoi(*resp.SignedRequests)
		return signed, nil
	}

	return 0, nil
}

// ApproveRequest approves a single request using a private key for signing.
func (s *RequestService) ApproveRequest(ctx context.Context, request *model.Request, privateKey *ecdsa.PrivateKey) (int, error) {
	if request == nil {
		return 0, fmt.Errorf("request cannot be nil")
	}
	return s.ApproveRequests(ctx, []*model.Request{request}, privateKey)
}

// RejectRequests rejects multiple requests with a comment.
func (s *RequestService) RejectRequests(ctx context.Context, requestIDs []string, comment string) error {
	if len(requestIDs) == 0 {
		return fmt.Errorf("requestIDs list cannot be empty")
	}
	if comment == "" {
		return fmt.Errorf("comment cannot be empty")
	}

	rejectReq := openapi.TgvalidatordRejectRequestsRequest{
		Comment: comment,
		Ids:     requestIDs,
	}

	_, httpResp, err := s.api.RequestServiceRejectRequests(ctx).
		Body(rejectReq).
		Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// RejectRequest rejects a single request with a comment.
func (s *RequestService) RejectRequest(ctx context.Context, requestID string, comment string) error {
	if requestID == "" {
		return fmt.Errorf("requestID cannot be empty")
	}
	return s.RejectRequests(ctx, []string{requestID}, comment)
}
