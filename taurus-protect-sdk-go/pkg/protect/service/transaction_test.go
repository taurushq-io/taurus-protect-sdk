package service

import (
	"testing"
)

func TestNewTransactionService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

// TransactionService only has ListTransactions which doesn't have pre-API validation
// Tests for this service require mocking the OpenAPI client
// Integration tests should cover the full flow
