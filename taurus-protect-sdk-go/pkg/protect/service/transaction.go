package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// TransactionService provides transaction query operations.
type TransactionService struct {
	api       *openapi.TransactionsAPIService
	errMapper *ErrorMapper
}

// NewTransactionService creates a new TransactionService.
func NewTransactionService(client *openapi.APIClient) *TransactionService {
	return &TransactionService{
		api:       client.TransactionsAPI,
		errMapper: NewErrorMapper(),
	}
}

// ListTransactions retrieves one page of transactions. The returned pagination is never nil;
// continue with its NextOffset until HasMore is false.
func (s *TransactionService) ListTransactions(ctx context.Context, opts *model.ListTransactionsOptions) ([]*model.Transaction, *model.Pagination, error) {
	if opts == nil {
		opts = &model.ListTransactionsOptions{}
	}
	window, err := resolveOffsetWindow(opts.Limit, opts.Offset)
	if err != nil {
		return nil, nil, err
	}

	req := applyOffsetWindow(s.api.TransactionServiceGetTransactions(ctx), window)
	if opts.Currency != "" {
		req = req.Currency(opts.Currency)
	}
	if opts.Direction != "" {
		req = req.Direction(opts.Direction)
	}
	if opts.Blockchain != "" {
		req = req.Blockchain(opts.Blockchain)
	}
	if opts.Network != "" {
		req = req.Network(opts.Network)
	}
	if opts.Query != "" {
		req = req.Query(opts.Query)
	}
	if opts.FromDate != nil {
		req = req.From(*opts.FromDate)
	}
	if opts.ToDate != nil {
		req = req.To(*opts.ToDate)
	}
	// Exact selectors. Query also matches hash and transaction id, but as a
	// six-column ILIKE scan — prefer these when the value is known exactly.
	if len(opts.IDs) > 0 {
		req = req.Ids(opts.IDs)
	}
	if len(opts.Hashes) > 0 {
		req = req.Hashes(opts.Hashes)
	}
	if len(opts.TransactionIDs) > 0 {
		req = req.TransactionIds(opts.TransactionIDs)
	}
	if opts.Address != "" {
		req = req.Address(opts.Address)
	}
	if opts.Source != "" {
		req = req.Source(opts.Source)
	}
	if opts.Destination != "" {
		req = req.Destination(opts.Destination)
	}
	if opts.AmountAbove != "" {
		req = req.AmountAbove(opts.AmountAbove)
	}

	return s.transactionsPage(req, window)
}

// transactionsPage executes a GetTransactions request. The next offset is offset + rows: on
// the post-filter path the total is an upper bound, and a page that made no progress ends the
// walk.
func (s *TransactionService) transactionsPage(req openapi.ApiTransactionServiceGetTransactionsRequest, window offsetWindow) ([]*model.Transaction, *model.Pagination, error) {
	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, nil, s.errMapper.MapError(err, httpResp)
	}

	pagination, err := offsetPagination(rulePlusRows, window, len(resp.Result), 0,
		offsetReply{TotalItems: resp.TotalItems})
	if err != nil {
		return nil, nil, err
	}
	return mapper.TransactionsFromDTO(resp.Result), pagination, nil
}

// GetTransaction retrieves a transaction by ID.
func (s *TransactionService) GetTransaction(ctx context.Context, txID string) (*model.Transaction, error) {
	if txID == "" {
		return nil, fmt.Errorf("txID cannot be empty")
	}

	// Use the list endpoint with IDs filter to get a single transaction
	req := s.api.TransactionServiceGetTransactions(ctx).
		Ids([]string{txID}).
		Limit("1")

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if len(resp.Result) == 0 {
		return nil, fmt.Errorf("transaction not found")
	}

	return mapper.TransactionFromDTO(&resp.Result[0]), nil
}

// GetTransactionByHash retrieves a transaction by its blockchain hash.
func (s *TransactionService) GetTransactionByHash(ctx context.Context, hash string) (*model.Transaction, error) {
	if hash == "" {
		return nil, fmt.Errorf("hash cannot be empty")
	}

	// Use the list endpoint with Hashes filter to get a single transaction
	req := s.api.TransactionServiceGetTransactions(ctx).
		Hashes([]string{hash}).
		Limit("1")

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if len(resp.Result) == 0 {
		return nil, fmt.Errorf("transaction not found")
	}

	return mapper.TransactionFromDTO(&resp.Result[0]), nil
}

// ListTransactionsByAddress retrieves one page of transactions for a specific address. The
// returned pagination is never nil; continue with its NextOffset until HasMore is false.
func (s *TransactionService) ListTransactionsByAddress(ctx context.Context, address string, opts *model.ListTransactionsByAddressOptions) ([]*model.Transaction, *model.Pagination, error) {
	if address == "" {
		return nil, nil, fmt.Errorf("address cannot be empty")
	}
	if opts == nil {
		opts = &model.ListTransactionsByAddressOptions{}
	}
	window, err := resolveOffsetWindow(opts.Limit, opts.Offset)
	if err != nil {
		return nil, nil, err
	}

	req := applyOffsetWindow(s.api.TransactionServiceGetTransactions(ctx), window).Address(address)
	if opts.Currency != "" {
		req = req.Currency(opts.Currency)
	}
	if opts.Direction != "" {
		req = req.Direction(opts.Direction)
	}
	if opts.Blockchain != "" {
		req = req.Blockchain(opts.Blockchain)
	}

	return s.transactionsPage(req, window)
}

// ExportTransactions exports transactions as text: json (the server's default when Format is
// empty), csv or csv_simple. The export cannot page — validatord always exports from the first
// matching transaction — so Limit (default 20, no SDK maximum) is the only size control, and
// TotalItems above Limit means the export was truncated.
func (s *TransactionService) ExportTransactions(ctx context.Context, opts *model.ExportTransactionsOptions) (*model.ExportTransactionsResult, error) {
	if opts == nil {
		opts = &model.ExportTransactionsOptions{}
	}
	limit, err := resolveSize("Limit", opts.Limit, 0)
	if err != nil {
		return nil, err
	}

	req := s.api.TransactionServiceExportTransactions(ctx).Limit(strconv.FormatInt(limit, 10))
	if opts.Currency != "" {
		req = req.Currency(opts.Currency)
	}
	if opts.Direction != "" {
		req = req.Direction(opts.Direction)
	}
	if opts.Blockchain != "" {
		req = req.Blockchain(opts.Blockchain)
	}
	if opts.Format != "" {
		req = req.Format(opts.Format)
	}
	if opts.From != nil {
		req = req.From(*opts.From)
	}
	if opts.To != nil {
		req = req.To(*opts.To)
	}
	if opts.Address != "" {
		req = req.Address(opts.Address)
	}
	if opts.Query != "" {
		req = req.Query(opts.Query)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	total, err := parseCount("totalItems", resp.TotalItems)
	if err != nil {
		return nil, err
	}
	result := &model.ExportTransactionsResult{TotalItems: total}
	if resp.Result != nil {
		result.Data = *resp.Result
	}
	return result, nil
}
