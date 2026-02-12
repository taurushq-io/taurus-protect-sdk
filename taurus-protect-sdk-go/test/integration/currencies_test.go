package integration

import (
	"context"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestIntegration_ListCurrencies(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	currencies, err := client.Currencies().GetCurrencies(ctx, nil)
	if err != nil {
		t.Fatalf("GetCurrencies() error = %v", err)
	}

	t.Logf("Found %d currencies", len(currencies))
	for _, c := range currencies {
		t.Logf("Currency: ID=%s, Symbol=%s, Name=%s", c.ID, c.Symbol, c.Name)
	}
}

func TestIntegration_GetCurrency(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// First list currencies to get a valid ID
	currencies, err := client.Currencies().GetCurrencies(ctx, nil)
	if err != nil {
		t.Fatalf("GetCurrencies() error = %v", err)
	}
	if len(currencies) == 0 {
		t.Skip("No currencies available for testing")
	}

	// Find a currency with both blockchain and network set (required by API)
	var source *model.Currency
	for _, c := range currencies {
		if c.Blockchain != "" && c.Network != "" {
			source = c
			break
		}
	}
	if source == nil {
		t.Skip("No currency with blockchain and network available for testing")
	}

	t.Logf("Selected currency for lookup: ID=%s, Symbol=%s, Blockchain=%s, Network=%s",
		source.ID, source.Symbol, source.Blockchain, source.Network)

	// Note: GetCurrency API looks up by blockchain+network, not by currencyID
	// The currencyID filter is optional and typically used for tokens
	currency, err := client.Currencies().GetCurrency(ctx, &model.GetCurrencyOptions{
		Blockchain: source.Blockchain,
		Network:    source.Network,
	})
	if err != nil {
		t.Fatalf("GetCurrency(blockchain=%s, network=%s) error = %v", source.Blockchain, source.Network, err)
	}

	t.Logf("Currency details:")
	t.Logf("  ID: %s", currency.ID)
	t.Logf("  Symbol: %s", currency.Symbol)
	t.Logf("  Name: %s", currency.Name)
	t.Logf("  Blockchain: %s", currency.Blockchain)
}
