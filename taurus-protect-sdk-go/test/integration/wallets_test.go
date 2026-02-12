package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestIntegration_ListWallets(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	wallets, pagination, err := client.Wallets().ListWallets(ctx, &model.ListWalletsOptions{
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("ListWallets() error = %v", err)
	}

	t.Logf("Found %d wallets", len(wallets))
	if pagination != nil {
		t.Logf("Total items: %d, HasMore: %v", pagination.TotalItems, pagination.HasMore)
	}

	for _, w := range wallets {
		t.Logf("Wallet: ID=%s, Name=%s, Currency=%s", w.ID, w.Name, w.Currency)
	}
}

func TestIntegration_GetWallet(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	// First, get a list to find a valid wallet ID
	ctx := context.Background()
	wallets, _, err := client.Wallets().ListWallets(ctx, &model.ListWalletsOptions{
		Limit: 1,
	})
	if err != nil {
		t.Fatalf("ListWallets() error = %v", err)
	}
	if len(wallets) == 0 {
		t.Skip("No wallets available for testing")
	}

	walletID := wallets[0].ID
	wallet, err := client.Wallets().GetWallet(ctx, walletID)
	if err != nil {
		t.Fatalf("GetWallet(%s) error = %v", walletID, err)
	}

	t.Logf("Wallet details:")
	t.Logf("  ID: %s", wallet.ID)
	t.Logf("  Name: %s", wallet.Name)
	t.Logf("  Currency: %s", wallet.Currency)
	t.Logf("  Blockchain: %s", wallet.Blockchain)
	t.Logf("  AddressesCount: %d", wallet.AddressesCount)
}

func TestIntegration_Pagination(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	pageSize := int64(2)
	allWallets := make([]*model.Wallet, 0)
	offset := int64(0)

	// Fetch all wallets using pagination
	for {
		wallets, pagination, err := client.Wallets().ListWallets(ctx, &model.ListWalletsOptions{
			Limit:  pageSize,
			Offset: offset,
		})
		if err != nil {
			t.Fatalf("ListWallets(offset=%d) error = %v", offset, err)
		}

		allWallets = append(allWallets, wallets...)
		t.Logf("Fetched %d wallets (offset=%d)", len(wallets), offset)
		for _, w := range wallets {
			fmt.Println(w)
		}

		if pagination == nil || !pagination.HasMore {
			break
		}
		offset += pageSize

		// Safety limit for tests
		if offset > 100 {
			t.Log("Stopping pagination test at 100 items")
			break
		}
	}

	t.Logf("Total wallets fetched via pagination: %d", len(allWallets))
}
