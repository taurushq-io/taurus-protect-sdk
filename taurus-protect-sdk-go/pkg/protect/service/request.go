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
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// RequestService provides transaction request management operations.
type RequestService struct {
	api       *openapi.RequestsAPIService
	errMapper *ErrorMapper
	logger    Logger
}

// NewRequestService creates a new RequestService.
func NewRequestService(client *openapi.APIClient, opts ...ServiceOption) *RequestService {
	return &RequestService{
		api:       client.RequestsAPI,
		errMapper: NewErrorMapper(),
		logger:    applyServiceOptions(opts).logger,
	}
}

// verifiedRequest maps one request DTO and verifies its metadata, and is the ONLY way
// a *model.Request is built in this package.
//
//	GetRequest ─────────┐
//	List/ForApproval ───┤
//	CreateOutgoing ─────┤
//	CreateCancel ───────┼─▶ verifiedRequest ─▶ RequestFromDTO ─▶ VerifyAndMaterialise
//	CreateIncoming ─────┘                                            │
//	(+ the 4 convenience wrappers, which delegate)                   │
//	                                          ok ◀──────────────────┴─────▶ IntegrityError
//	                                    (entries populated)              (wrapped, names the id)
//
// Verification lived inline in GetRequest only. The list paths skipped it in one pass
// and the six create paths in the next — both times because "remember to verify" was a
// rule rather than the only available construction path. So: do NOT call
// mapper.RequestFromDTO or mapper.RequestsFromDTO anywhere else in this file.
//
// The id is wrapped into the error because a failing create has already succeeded
// server-side; the caller needs the id to reconcile rather than retrying and
// double-creating. %w keeps the underlying *model.IntegrityError matchable.
func (s *RequestService) verifiedRequest(dto *openapi.TgvalidatordRequest) (*model.Request, error) {
	r := mapper.RequestFromDTO(dto)
	if r == nil {
		return nil, fmt.Errorf("request not found")
	}
	if err := r.Metadata.VerifyAndMaterialise(); err != nil {
		return nil, fmt.Errorf("request %s: %w", r.ID, err)
	}
	return r, nil
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

	return s.verifiedRequest(resp.Result)
}

// verifiedRequests keeps the rows whose metadata verifies and materialises their
// payload, dropping the rest.
//
//	rows ──▶ verifiedRequest ──┬─ ok ──────▶ kept (entries populated)
//	                           └─ error ───▶ excluded + logged + named in the result
//
// Excluding rather than failing the whole call keeps one corrupt row from denying
// access to every good one; naming the exclusions keeps a shortened list from
// reading as a complete one. Both list paths share this — duplicating it is how the
// two drifted apart before. It takes DTOs, not models, so the mapper stays behind the
// seam.
func (s *RequestService) verifiedRequests(ctx context.Context, dtos []openapi.TgvalidatordRequest) ([]*model.Request, []string) {
	kept := make([]*model.Request, 0, len(dtos))
	var excluded []string

	for i := range dtos {
		r, err := s.verifiedRequest(&dtos[i])
		if err != nil {
			id := ""
			if dtos[i].Id != nil {
				id = *dtos[i].Id
			}
			excluded = append(excluded, id)
			s.logger.Warn(ctx, "request excluded: metadata integrity verification failed",
				Field{Key: "resource", Value: "request"},
				Field{Key: "request_id", Value: id},
				Field{Key: "reason", Value: err.Error()})
			continue
		}
		kept = append(kept, r)
	}
	return kept, excluded
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
		if len(opts.Statuses) > 0 {
			req = req.Statuses(opts.Statuses)
		}
		if len(opts.IDs) > 0 {
			req = req.Ids(opts.IDs)
		}
		if len(opts.Types) > 0 {
			req = req.Types(opts.Types)
		}
		if len(opts.ExternalRequestIDs) > 0 {
			req = req.ExternalRequestIDs(opts.ExternalRequestIDs)
		}
		if opts.Currency != "" {
			req = req.CurrencyID(opts.Currency)
		}
		if opts.FromDate != nil {
			req = req.From(*opts.FromDate)
		}
		if opts.ToDate != nil {
			req = req.To(*opts.ToDate)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	kept, excluded := s.verifiedRequests(ctx, resp.Result)
	result := &model.RequestResult{
		Requests:           kept,
		ExcludedUnverified: excluded,
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

	kept, excluded := s.verifiedRequests(ctx, resp.Result)
	result := &model.RequestResult{
		Requests:           kept,
		ExcludedUnverified: excluded,
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

	return s.verifiedRequest(resp.Result)
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

	return s.verifiedRequest(resp.Result)
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

	return s.verifiedRequest(resp.Result)
}

// ApproveRequests approves multiple requests using a private key for signing.
// The requests are sorted by ID before signing. Returns the number of requests signed.
//
//	requests ─▶ presence check ─▶ RE-VERIFY hash vs payload ──fail──▶ refuse, sign nothing
//	                                          │ok
//	                                          ▼
//	                              numeric id ─▶ sort ─▶ sign(JSON(hashes)) ─▶ POST once
//
// HashVerified is deliberately NOT consulted — see the loop below.
// The optional comment is what gets recorded against the approval. It is variadic
// rather than a new positional parameter so existing callers keep compiling; Python
// and TypeScript take it as an optional argument, and it used to be hardcoded here.
func (s *RequestService) ApproveRequests(ctx context.Context, requests []*model.Request, privateKey *ecdsa.PrivateKey, comment ...string) (int, error) {
	if len(requests) == 0 {
		return 0, fmt.Errorf("requests list cannot be empty")
	}
	if privateKey == nil {
		return 0, fmt.Errorf("privateKey cannot be nil")
	}

	// Validate all requests have metadata with hash and valid numeric IDs, and
	// RE-VERIFY every hash this signature will attest to.
	//
	//	presence check ─▶ re-verify hash vs payload ─▶ numeric id ─▶ sort ─▶ sign once
	//
	// HashVerified is deliberately NOT consulted. It is a serialized field
	// (`json:"hash_verified,omitempty"`), so a *model.Request decoded from a queue,
	// webhook or cached blob can arrive claiming true and get an attacker-chosen hash
	// signed by the approver's real key. Re-verification is a SHA-256 over a string we
	// already hold — no network — and it is the same guarantee
	// WhitelistedAssetService.ApproveWhitelistedAssets gets by re-reading, so both
	// signing paths now rest on the same rule rather than two different ones.
	for _, r := range requests {
		if r.Metadata == nil || r.Metadata.Hash == "" {
			return 0, fmt.Errorf("request %s has no metadata hash", r.ID)
		}
		// Ordered after the presence check so an early-status request still reports the
		// clearer "no metadata hash".
		if err := r.Metadata.VerifyAndMaterialise(); err != nil {
			return 0, fmt.Errorf("refusing to sign request %s: %w", r.ID, err)
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
		Comment:   approvalComment(comment),
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
func (s *RequestService) ApproveRequest(ctx context.Context, request *model.Request, privateKey *ecdsa.PrivateKey, comment ...string) (int, error) {
	if request == nil {
		return 0, fmt.Errorf("request cannot be nil")
	}
	return s.ApproveRequests(ctx, []*model.Request{request}, privateKey, comment...)
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

// approvalComment picks the caller's comment, falling back to a default so an
// approval always carries some rationale.
func approvalComment(comment []string) string {
	if len(comment) > 0 && comment[0] != "" {
		return comment[0]
	}
	return "approved via taurus-protect-sdk-go"
}
