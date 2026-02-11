package integration

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestPaginateAllWhitelistedAddresses(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	pageSize := int64(50)
	allAddresses := make([]*model.WhitelistedAddress, 0)
	offset := int64(0)

	// Fetch all whitelisted addresses using pagination
	for {
		addresses, pagination, err := client.WhitelistedAddresses().ListWhitelistedAddresses(ctx, &model.ListWhitelistedAddressesOptions{
			Limit:  pageSize,
			Offset: offset,
		})
		if err != nil {
			t.Fatalf("ListWhitelistedAddresses(offset=%d) error = %v", offset, err)
		}

		allAddresses = append(allAddresses, addresses...)
		t.Logf("Fetched %d whitelisted addresses (offset=%d)", len(addresses), offset)

		for _, a := range addresses {
			t.Logf("  WhitelistedAddress: ID=%s, Address=%s, Label=%s, Blockchain=%s", a.ID, a.Address, a.Label, a.Blockchain)
		}

		if pagination == nil || !pagination.HasMore {
			break
		}
		offset += pageSize

		// Safety limit for tests
		if offset > 2000 {
			t.Log("Stopping pagination test at 2000 items")
			break
		}
	}

	t.Logf("Total whitelisted addresses fetched via pagination: %d", len(allAddresses))
}

func TestListWhitelistedAddresses(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()

	addresses, _, err := client.WhitelistedAddresses().ListWhitelistedAddresses(ctx, &model.ListWhitelistedAddressesOptions{
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("ListWhitelistedAddresses error = %v", err)
	}

	t.Logf("Found %d whitelisted addresses", len(addresses))
	for _, a := range addresses {
		t.Logf("  %s/%s: %s", a.Blockchain, a.Network, a.Address)
	}
}

func TestPaginateAllWhitelistedAssets(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	pageSize := int64(10)
	allAssets := make([]*model.WhitelistedAsset, 0)
	offset := int64(0)

	// Fetch all whitelisted assets using pagination
	for {
		assets, pagination, err := client.WhitelistedAssets().ListWhitelistedAssets(ctx, &model.ListWhitelistedAssetsOptions{
			Limit:  pageSize,
			Offset: offset,
		})
		if err != nil {
			t.Fatalf("ListWhitelistedAssets(offset=%d) error = %v", offset, err)
		}

		allAssets = append(allAssets, assets...)
		t.Logf("Fetched %d whitelisted assets (offset=%d)", len(assets), offset)

		for _, a := range allAssets[len(allAssets)-len(assets):] {
			t.Logf("  WhitelistedAsset: ID=%s, Blockchain=%s, Network=%s, Status=%s", a.ID, a.Blockchain, a.Network, a.Status)
		}

		if pagination == nil || !pagination.HasMore {
			break
		}
		offset += pageSize

		// Safety limit for tests
		if offset > 5000 {
			t.Log("Stopping pagination test at 5000 items")
			break
		}
	}

	t.Logf("Total whitelisted assets fetched via pagination: %d", len(allAssets))
}

func TestListWhitelistedAssets(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()

	assets, _, err := client.WhitelistedAssets().ListWhitelistedAssets(ctx, &model.ListWhitelistedAssetsOptions{
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("ListWhitelistedAssets error = %v", err)
	}

	t.Logf("Found %d whitelisted assets", len(assets))
	for _, a := range assets {
		t.Logf("  ID=%s, Blockchain=%s, Network=%s, Status=%s", a.ID, a.Blockchain, a.Network, a.Status)
	}
}

func TestGetWhitelistedAddress(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// First list addresses to find a valid ID
	addresses, _, err := client.WhitelistedAddresses().ListWhitelistedAddresses(ctx, &model.ListWhitelistedAddressesOptions{
		Limit: 1,
	})
	if err != nil {
		t.Fatalf("ListWhitelistedAddresses error = %v", err)
	}

	if len(addresses) == 0 {
		t.Skip("No whitelisted addresses available for testing")
	}

	// Get by ID
	addressID := addresses[0].ID
	address, err := client.WhitelistedAddresses().GetWhitelistedAddress(ctx, addressID)
	if err != nil {
		t.Fatalf("GetWhitelistedAddress(id=%s) error = %v", addressID, err)
	}

	t.Logf("GetWhitelistedAddress: ID=%s, Blockchain=%s, Network=%s, Address=%s, Label=%s, Status=%s",
		address.ID, address.Blockchain, address.Network, address.Address, address.Label, address.Status)
}

func TestListWhitelistedAddressesByBlockchain(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// List addresses filtered by blockchain ETH
	addresses, pagination, err := client.WhitelistedAddresses().ListWhitelistedAddresses(ctx, &model.ListWhitelistedAddressesOptions{
		Blockchain: "ETH",
		Limit:      10,
	})
	if err != nil {
		t.Fatalf("ListWhitelistedAddresses(blockchain=ETH) error = %v", err)
	}

	t.Logf("Listed %d whitelisted addresses for blockchain ETH", len(addresses))
	if pagination != nil {
		t.Logf("Pagination: TotalItems=%d, HasMore=%v", pagination.TotalItems, pagination.HasMore)
	}

	for _, a := range addresses {
		t.Logf("  WhitelistedAddress: ID=%s, Blockchain=%s, Network=%s, Address=%s",
			a.ID, a.Blockchain, a.Network, a.Address)
		// Verify filter is working
		if a.Blockchain != "ETH" {
			t.Errorf("Expected Blockchain=ETH, got %s", a.Blockchain)
		}
	}
}

func TestListWhitelistedAddressesByBlockchainAndNetwork(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// List addresses filtered by blockchain ETH and network mainnet
	addresses, pagination, err := client.WhitelistedAddresses().ListWhitelistedAddresses(ctx, &model.ListWhitelistedAddressesOptions{
		Blockchain: "ETH",
		Network:    "mainnet",
		Limit:      10,
	})
	if err != nil {
		t.Fatalf("ListWhitelistedAddresses(blockchain=ETH, network=mainnet) error = %v", err)
	}

	t.Logf("Listed %d whitelisted addresses for blockchain ETH, network mainnet", len(addresses))
	if pagination != nil {
		t.Logf("Pagination: TotalItems=%d, HasMore=%v", pagination.TotalItems, pagination.HasMore)
	}

	for _, a := range addresses {
		t.Logf("  WhitelistedAddress: ID=%s, Blockchain=%s, Network=%s, Address=%s",
			a.ID, a.Blockchain, a.Network, a.Address)
		// Verify filters are working
		if a.Blockchain != "ETH" {
			t.Errorf("Expected Blockchain=ETH, got %s", a.Blockchain)
		}
		if a.Network != "mainnet" {
			t.Errorf("Expected Network=mainnet, got %s", a.Network)
		}
	}
}

func TestGetWhitelistedAsset(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// First list assets to find a valid ID
	assets, _, err := client.WhitelistedAssets().ListWhitelistedAssets(ctx, &model.ListWhitelistedAssetsOptions{
		Limit: 1,
	})
	if err != nil {
		t.Fatalf("ListWhitelistedAssets error = %v", err)
	}

	if len(assets) == 0 {
		t.Skip("No whitelisted assets available for testing")
	}

	// Get by ID
	assetID := assets[0].ID
	asset, err := client.WhitelistedAssets().GetWhitelistedAsset(ctx, assetID)
	if err != nil {
		t.Fatalf("GetWhitelistedAsset(id=%s) error = %v", assetID, err)
	}

	t.Logf("GetWhitelistedAsset: ID=%s, Blockchain=%s, Network=%s, Status=%s",
		asset.ID, asset.Blockchain, asset.Network, asset.Status)

	// SECURITY: Log metadata payload fields by parsing from PayloadAsString (verified source)
	if asset.Metadata != nil && asset.Metadata.PayloadAsString != "" {
		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(asset.Metadata.PayloadAsString), &payload); err == nil {
			if symbol, ok := payload["symbol"].(string); ok {
				t.Logf("  Symbol=%s", symbol)
			}
			if contractAddress, ok := payload["contractAddress"].(string); ok {
				t.Logf("  ContractAddress=%s", contractAddress)
			}
			if name, ok := payload["name"].(string); ok {
				t.Logf("  Name=%s", name)
			}
			if decimals, ok := payload["decimals"]; ok {
				t.Logf("  Decimals=%v", decimals)
			}
		}
	}
}

func TestListWhitelistedAssetsByBlockchain(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// List assets filtered by blockchain ETH
	assets, pagination, err := client.WhitelistedAssets().ListWhitelistedAssets(ctx, &model.ListWhitelistedAssetsOptions{
		Blockchain: "ETH",
		Limit:      10,
	})
	if err != nil {
		t.Fatalf("ListWhitelistedAssets(blockchain=ETH) error = %v", err)
	}

	t.Logf("Listed %d whitelisted assets for blockchain ETH", len(assets))
	if pagination != nil {
		t.Logf("Pagination: TotalItems=%d, HasMore=%v", pagination.TotalItems, pagination.HasMore)
	}

	for _, a := range assets {
		t.Logf("  WhitelistedAsset: ID=%s, Blockchain=%s, Network=%s, Status=%s",
			a.ID, a.Blockchain, a.Network, a.Status)
		// Verify filter is working
		if a.Blockchain != "ETH" {
			t.Errorf("Expected Blockchain=ETH, got %s", a.Blockchain)
		}
	}
}

func TestListWhitelistedAssetsByBlockchainAndNetwork(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// List assets filtered by blockchain ETH and network mainnet
	assets, pagination, err := client.WhitelistedAssets().ListWhitelistedAssets(ctx, &model.ListWhitelistedAssetsOptions{
		Blockchain: "ETH",
		Network:    "mainnet",
		Limit:      10,
	})
	if err != nil {
		t.Fatalf("ListWhitelistedAssets(blockchain=ETH, network=mainnet) error = %v", err)
	}

	t.Logf("Listed %d whitelisted assets for blockchain ETH, network mainnet", len(assets))
	if pagination != nil {
		t.Logf("Pagination: TotalItems=%d, HasMore=%v", pagination.TotalItems, pagination.HasMore)
	}

	for _, a := range assets {
		t.Logf("  WhitelistedAsset: ID=%s, Blockchain=%s, Network=%s, Status=%s",
			a.ID, a.Blockchain, a.Network, a.Status)
		// Verify filters are working
		if a.Blockchain != "ETH" {
			t.Errorf("Expected Blockchain=ETH, got %s", a.Blockchain)
		}
		if a.Network != "mainnet" {
			t.Errorf("Expected Network=mainnet, got %s", a.Network)
		}
	}
}
