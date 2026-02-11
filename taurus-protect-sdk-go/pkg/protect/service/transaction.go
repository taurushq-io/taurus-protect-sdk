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

// ListTransactions retrieves a list of transactions.
func (s *TransactionService) ListTransactions(ctx context.Context, opts *model.ListTransactionsOptions) ([]*model.Transaction, *model.Pagination, error) {
	req := s.api.TransactionServiceGetTransactions(ctx)

	if opts != nil {
		if opts.Limit > 0 {
			req = req.Limit(fmt.Sprintf("%d", opts.Limit))
		}
		if opts.Offset > 0 {
			req = req.Offset(fmt.Sprintf("%d", opts.Offset))
		}
		if opts.Currency != "" {
			req = req.Currency(opts.Currency)
		}
		if opts.Direction != "" {
			req = req.Direction(opts.Direction)
		}
		if opts.Blockchain != "" {
			req = req.Blockchain(opts.Blockchain)
		}
		if opts.Query != "" {
			req = req.Query(opts.Query)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, nil, s.errMapper.MapError(err, httpResp)
	}

	transactions := mapper.TransactionsFromDTO(resp.Result)

	var pagination *model.Pagination
	if resp.TotalItems != nil {
		pagination = &model.Pagination{}
		if total, err := strconv.ParseInt(*resp.TotalItems, 10, 64); err == nil {
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

	return transactions, pagination, nil
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

	if resp.Result == nil || len(resp.Result) == 0 {
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

	if resp.Result == nil || len(resp.Result) == 0 {
		return nil, fmt.Errorf("transaction not found")
	}

	return mapper.TransactionFromDTO(&resp.Result[0]), nil
}

// ListTransactionsByAddress retrieves a list of transactions for a specific address.
func (s *TransactionService) ListTransactionsByAddress(ctx context.Context, address string, opts *model.ListTransactionsByAddressOptions) ([]*model.Transaction, *model.Pagination, error) {
	if address == "" {
		return nil, nil, fmt.Errorf("address cannot be empty")
	}

	req := s.api.TransactionServiceGetTransactions(ctx).
		Address(address)

	if opts != nil {
		if opts.Limit > 0 {
			req = req.Limit(fmt.Sprintf("%d", opts.Limit))
		}
		if opts.Offset > 0 {
			req = req.Offset(fmt.Sprintf("%d", opts.Offset))
		}
		if opts.Currency != "" {
			req = req.Currency(opts.Currency)
		}
		if opts.Direction != "" {
			req = req.Direction(opts.Direction)
		}
		if opts.Blockchain != "" {
			req = req.Blockchain(opts.Blockchain)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, nil, s.errMapper.MapError(err, httpResp)
	}

	transactions := mapper.TransactionsFromDTO(resp.Result)

	var pagination *model.Pagination
	if resp.TotalItems != nil {
		pagination = &model.Pagination{}
		if total, err := strconv.ParseInt(*resp.TotalItems, 10, 64); err == nil {
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

	return transactions, pagination, nil
}

// ExportTransactions exports transactions in the specified format.
// The format can be "csv", "json", or "csv_simple". If not specified, defaults to "csv".
// Returns the exported data as a string.
func (s *TransactionService) ExportTransactions(ctx context.Context, opts *model.ExportTransactionsOptions) (string, error) {
	req := s.api.TransactionServiceExportTransactions(ctx)

	if opts != nil {
		if opts.Limit > 0 {
			req = req.Limit(fmt.Sprintf("%d", opts.Limit))
		}
		if opts.Offset > 0 {
			req = req.Offset(fmt.Sprintf("%d", opts.Offset))
		}
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
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return "", s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return "", nil
	}

	return *resp.Result, nil
}
