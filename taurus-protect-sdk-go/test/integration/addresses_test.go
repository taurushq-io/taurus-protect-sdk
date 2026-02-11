package integration

import (
	"context"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestIntegration_ListAddresses(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClientWithVerification(t)
	defer client.Close()

	ctx := context.Background()
	addresses, pagination, err := client.Addresses().ListAddresses(ctx, &model.ListAddressesOptions{
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("ListAddresses() error = %v", err)
	}

	t.Logf("Found %d addresses", len(addresses))
	if pagination != nil {
		t.Logf("Total items: %d, HasMore: %v", pagination.TotalItems, pagination.HasMore)
	}

	for _, a := range addresses {
		t.Logf("Address: ID=%s, Label=%s, Currency=%s", a.ID, a.Label, a.Currency)
	}
}

func TestIntegration_GetAddress(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClientWithVerification(t)
	defer client.Close()

	// First, get a list to find a valid address ID
	ctx := context.Background()
	addresses, _, err := client.Addresses().ListAddresses(ctx, &model.ListAddressesOptions{
		Limit: 1,
	})
	if err != nil {
		t.Fatalf("ListAddresses() error = %v", err)
	}
	if len(addresses) == 0 {
		t.Skip("No addresses available for testing")
	}

	addressID := addresses[0].ID
	address, err := client.Addresses().GetAddress(ctx, addressID)
	if err != nil {
		t.Fatalf("GetAddress(%s) error = %v", addressID, err)
	}

	t.Logf("Address details:")
	t.Logf("  ID: %s", address.ID)
	t.Logf("  Address: %s", address.Address)
	t.Logf("  Label: %s", address.Label)
	t.Logf("  Currency: %s", address.Currency)
	t.Logf("  WalletID: %s", address.WalletID)
	if address.Balance != nil {
		t.Logf("  Balance (confirmed): %s", address.Balance.TotalConfirmed)
	}
}
