package integration

import (
	"context"
	"testing"
)

func TestIntegration_ListBalances(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	result, err := client.Balances().GetBalances(ctx, nil)
	if err != nil {
		t.Fatalf("GetBalances() error = %v", err)
	}

	t.Logf("Found %d balances", len(result.Balances))
	for _, b := range result.Balances {
		currency := ""
		totalConfirmed := ""
		if b.Asset != nil {
			currency = b.Asset.Currency
		}
		if b.Balance != nil {
			totalConfirmed = b.Balance.TotalConfirmed
		}
		t.Logf("Balance: Currency=%s, TotalConfirmed=%s", currency, totalConfirmed)
	}
}
