package integration

import (
	"context"
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestIntegration_ListTransactions(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	transactions, pagination, err := client.Transactions().ListTransactions(ctx, &model.ListTransactionsOptions{
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("ListTransactions() error = %v", err)
	}

	t.Logf("Found %d transactions", len(transactions))
	if pagination != nil {
		t.Logf("Total items: %d, HasMore: %v", pagination.TotalItems, pagination.HasMore)
	}

	for _, tx := range transactions {
		t.Logf("Transaction: ID=%s, Hash=%s, Amount=%s", tx.ID, tx.Hash, tx.Amount)
	}
}

func TestIntegration_ListTransactionsByCurrency(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	transactions, _, err := client.Transactions().ListTransactions(ctx, &model.ListTransactionsOptions{
		Limit:    10,
		Currency: "ETH",
	})
	if err != nil {
		t.Fatalf("ListTransactions(Currency=ETH) error = %v", err)
	}

	t.Logf("Found %d ETH transactions", len(transactions))
	for _, tx := range transactions {
		t.Logf("Transaction: ID=%s, Direction=%s, Amount=%s", tx.ID, tx.Direction, tx.Amount)
	}
}

func TestIntegration_GetTransactionById(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// First list transactions to find a valid ID
	transactions, _, err := client.Transactions().ListTransactions(ctx, &model.ListTransactionsOptions{
		Limit: 1,
	})
	if err != nil {
		t.Fatalf("ListTransactions() error = %v", err)
	}

	if len(transactions) == 0 {
		t.Skip("No transactions found, skipping GetTransaction test")
	}

	txID := transactions[0].ID
	t.Logf("Getting transaction by ID: %s", txID)

	// Get the transaction by ID
	tx, err := client.Transactions().GetTransaction(ctx, txID)
	if err != nil {
		t.Fatalf("GetTransaction(%s) error = %v", txID, err)
	}

	t.Logf("Transaction retrieved: ID=%s, Hash=%s, Amount=%s, Currency=%s, Direction=%s",
		tx.ID, tx.Hash, tx.Amount, tx.Currency, tx.Direction)

	if tx.ID != txID {
		t.Errorf("Transaction ID mismatch: expected %s, got %s", txID, tx.ID)
	}
}

func TestIntegration_GetTransactionByHash(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// First list transactions to find a valid hash
	transactions, _, err := client.Transactions().ListTransactions(ctx, &model.ListTransactionsOptions{
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("ListTransactions() error = %v", err)
	}

	// Find a transaction with a non-empty hash
	var txHash string
	for _, tx := range transactions {
		if tx.Hash != "" {
			txHash = tx.Hash
			break
		}
	}

	if txHash == "" {
		t.Skip("No transactions with hash found, skipping GetTransactionByHash test")
	}

	t.Logf("Getting transaction by hash: %s", txHash)

	// Get the transaction by hash
	tx, err := client.Transactions().GetTransactionByHash(ctx, txHash)
	if err != nil {
		t.Fatalf("GetTransactionByHash(%s) error = %v", txHash, err)
	}

	t.Logf("Transaction retrieved: ID=%s, Hash=%s, Amount=%s, Currency=%s, Direction=%s",
		tx.ID, tx.Hash, tx.Amount, tx.Currency, tx.Direction)

	if tx.Hash != txHash {
		t.Errorf("Transaction hash mismatch: expected %s, got %s", txHash, tx.Hash)
	}
}

func TestIntegration_ListTransactionsByAddress(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// First list transactions to find a source address
	transactions, _, err := client.Transactions().ListTransactions(ctx, &model.ListTransactionsOptions{
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("ListTransactions() error = %v", err)
	}

	// Find a transaction with a source address
	var sourceAddress string
	for _, tx := range transactions {
		if len(tx.Sources) > 0 && tx.Sources[0].Address != "" {
			sourceAddress = tx.Sources[0].Address
			break
		}
	}

	if sourceAddress == "" {
		t.Skip("No transactions with source address found, skipping ListTransactionsByAddress test")
	}

	t.Logf("Listing transactions by address: %s", sourceAddress)

	// Get transactions by address
	txByAddr, pagination, err := client.Transactions().ListTransactionsByAddress(ctx, sourceAddress, &model.ListTransactionsByAddressOptions{
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("ListTransactionsByAddress(%s) error = %v", sourceAddress, err)
	}

	t.Logf("Found %d transactions for address %s", len(txByAddr), sourceAddress)
	if pagination != nil {
		t.Logf("Total items: %d, HasMore: %v", pagination.TotalItems, pagination.HasMore)
	}

	for _, tx := range txByAddr {
		t.Logf("Transaction: ID=%s, Hash=%s, Direction=%s, Amount=%s", tx.ID, tx.Hash, tx.Direction, tx.Amount)
	}
}

func TestIntegration_ExportTransactions(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// Set time range: 30 days back to now
	now := time.Now()
	thirtyDaysAgo := now.AddDate(0, 0, -30)

	t.Logf("Exporting transactions from %s to %s", thirtyDaysAgo.Format(time.RFC3339), now.Format(time.RFC3339))

	// Export transactions to CSV
	csvData, err := client.Transactions().ExportTransactions(ctx, &model.ExportTransactionsOptions{
		From:   &thirtyDaysAgo,
		To:     &now,
		Format: "csv",
		Limit:  100,
	})
	if err != nil {
		t.Fatalf("ExportTransactions() error = %v", err)
	}

	if csvData == "" {
		t.Log("No transactions exported (empty CSV)")
		return
	}

	// Log a snippet of the CSV (first 500 characters or full content if shorter)
	snippet := csvData
	if len(snippet) > 500 {
		snippet = snippet[:500] + "..."
	}
	t.Logf("CSV export snippet:\n%s", snippet)
	t.Logf("Total CSV length: %d characters", len(csvData))
}
