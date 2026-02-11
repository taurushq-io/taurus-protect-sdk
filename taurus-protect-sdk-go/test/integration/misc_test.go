package integration

import (
	"context"
	"testing"
)

func TestIntegration_ListTags(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	tags, err := client.Tags().ListTags(ctx, nil)
	if err != nil {
		t.Fatalf("ListTags() error = %v", err)
	}

	t.Logf("Found %d tags", len(tags))
	for _, tag := range tags {
		t.Logf("Tag: ID=%s, Value=%s", tag.ID, tag.Value)
	}
}

func TestIntegration_GetPortfolioStatistics(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	stats, err := client.Statistics().GetPortfolioStatistics(ctx)
	if err != nil {
		t.Fatalf("GetPortfolioStatistics() error = %v", err)
	}

	t.Logf("Portfolio statistics:")
	t.Logf("  TotalBalance: %s", stats.TotalBalance)
	t.Logf("  TotalBalanceBaseCurrency: %s", stats.TotalBalanceBaseCurrency)
	t.Logf("  WalletsCount: %s", stats.WalletsCount)
	t.Logf("  AddressesCount: %s", stats.AddressesCount)
}

func TestIntegration_ClientLifecycle(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)

	// Test that services are lazily initialized
	_ = client.Wallets()
	_ = client.Addresses()
	_ = client.Requests()
	_ = client.Transactions()

	// Test that Close works
	err := client.Close()
	if err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	// Test that Close is idempotent
	err = client.Close()
	if err != nil {
		t.Fatalf("Second Close() error = %v", err)
	}
}
