package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// TestIntegration_CaptureRawWhitelistedAddressResponse captures the raw API response
// for a whitelisted address to use as test fixtures in unit tests across all SDKs.
//
// This test outputs the raw JSON to stdout and optionally writes to a file.
// Run with: go test -v ./test/integration/... -run CaptureRawWhitelistedAddressResponse
func TestIntegration_CaptureRawWhitelistedAddressResponse(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// Get whitelisted addresses
	addresses, _, err := client.WhitelistedAddresses().ListWhitelistedAddresses(ctx, &model.ListWhitelistedAddressesOptions{
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("ListWhitelistedAddresses error: %v", err)
	}

	if len(addresses) == 0 {
		t.Skip("No whitelisted addresses available for capture")
	}

	// Find an address with full metadata (has payloadAsString and signatures)
	var selectedAddr *model.WhitelistedAddress
	for _, addr := range addresses {
		if addr.Metadata != nil && addr.Metadata.PayloadAsString != "" {
			if addr.SignedAddress != nil && len(addr.SignedAddress.Signatures) > 0 {
				selectedAddr = addr
				break
			}
		}
	}

	if selectedAddr == nil {
		t.Skip("No whitelisted address with complete metadata found")
	}

	// Marshal the model to JSON
	rawJSON, err := json.MarshalIndent(selectedAddr, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal address: %v", err)
	}

	t.Logf("Whitelisted Address model:\n%s", rawJSON)

	// Print key fields for quick reference
	t.Logf("\n=== Key Fields ===")
	t.Logf("ID: %s", selectedAddr.ID)
	t.Logf("Blockchain: %s", selectedAddr.Blockchain)
	t.Logf("Network: %s", selectedAddr.Network)
	t.Logf("Status: %s", selectedAddr.Status)
	t.Logf("Address: %s", selectedAddr.Address)
	t.Logf("Label: %s", selectedAddr.Label)

	if selectedAddr.Metadata != nil {
		t.Logf("Metadata.Hash: %s", selectedAddr.Metadata.Hash)
		payloadLen := len(selectedAddr.Metadata.PayloadAsString)
		if payloadLen > 100 {
			t.Logf("Metadata.PayloadAsString (truncated): %s...", selectedAddr.Metadata.PayloadAsString[:100])
		} else {
			t.Logf("Metadata.PayloadAsString: %s", selectedAddr.Metadata.PayloadAsString)
		}
	}

	if selectedAddr.SignedAddress != nil {
		t.Logf("SignedAddress.Signatures count: %d", len(selectedAddr.SignedAddress.Signatures))
		for i, sig := range selectedAddr.SignedAddress.Signatures {
			t.Logf("  Signature[%d]: hashes=%v", i, sig.Hashes)
		}
	}

	if selectedAddr.RulesContainer != "" {
		containerLen := len(selectedAddr.RulesContainer)
		t.Logf("RulesContainer length: %d", containerLen)
	}

	if selectedAddr.RulesSignatures != "" {
		sigLen := len(selectedAddr.RulesSignatures)
		t.Logf("RulesSignatures length: %d", sigLen)
	}

	// Log trails
	if len(selectedAddr.Trails) > 0 {
		t.Logf("Trails count: %d", len(selectedAddr.Trails))
		for i, trail := range selectedAddr.Trails {
			t.Logf("  Trail[%d]: action=%s, date=%v", i, trail.Action, trail.Date)
		}
	}

	// Log attributes
	if len(selectedAddr.Attributes) > 0 {
		t.Logf("Attributes count: %d", len(selectedAddr.Attributes))
		for i, attr := range selectedAddr.Attributes {
			t.Logf("  Attribute[%d]: key=%s, value=%s", i, attr.Key, attr.Value)
		}
	}

	// Write to file for external use
	outputPath := "/tmp/whitelisted_address_model.json"
	if err := os.WriteFile(outputPath, rawJSON, 0644); err != nil {
		t.Logf("Warning: could not write to %s: %v", outputPath, err)
	} else {
		t.Logf("\nWritten to: %s", outputPath)
	}

	// Also capture multiple addresses to have variety
	if len(addresses) > 1 {
		allJSON, _ := json.MarshalIndent(addresses, "", "  ")
		allOutputPath := "/tmp/whitelisted_addresses_model.json"
		if err := os.WriteFile(allOutputPath, allJSON, 0644); err != nil {
			t.Logf("Warning: could not write to %s: %v", allOutputPath, err)
		} else {
			t.Logf("All addresses written to: %s", allOutputPath)
		}
	}
}

// TestIntegration_CaptureRawWhitelistedAssetResponse captures the raw API response
// for a whitelisted asset (contract address) to use as test fixtures.
func TestIntegration_CaptureRawWhitelistedAssetResponse(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// Get whitelisted assets
	assets, _, err := client.WhitelistedAssets().ListWhitelistedAssets(ctx, &model.ListWhitelistedAssetsOptions{
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("ListWhitelistedAssets error: %v", err)
	}

	if len(assets) == 0 {
		t.Skip("No whitelisted assets available for capture")
	}

	// Find an asset with full metadata
	var selectedAsset *model.WhitelistedAsset
	for _, asset := range assets {
		if asset.Metadata != nil && asset.Metadata.PayloadAsString != "" {
			if asset.SignedContractAddress != nil && len(asset.SignedContractAddress.Signatures) > 0 {
				selectedAsset = asset
				break
			}
		}
	}

	if selectedAsset == nil {
		t.Skip("No whitelisted asset with complete metadata found")
	}

	// Marshal the model to JSON
	rawJSON, err := json.MarshalIndent(selectedAsset, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal asset: %v", err)
	}

	t.Logf("Whitelisted Asset model:\n%s", rawJSON)

	// Print key fields
	t.Logf("\n=== Key Fields ===")
	t.Logf("ID: %s", selectedAsset.ID)
	t.Logf("Blockchain: %s", selectedAsset.Blockchain)
	t.Logf("Network: %s", selectedAsset.Network)
	t.Logf("Status: %s", selectedAsset.Status)

	if selectedAsset.Metadata != nil {
		t.Logf("Metadata.Hash: %s", selectedAsset.Metadata.Hash)
		// SECURITY: Extract name and symbol from PayloadAsString (verified source)
		if selectedAsset.Metadata.PayloadAsString != "" {
			var payload map[string]interface{}
			if err := json.Unmarshal([]byte(selectedAsset.Metadata.PayloadAsString), &payload); err == nil {
				if name, ok := payload["name"].(string); ok {
					t.Logf("Name (from PayloadAsString): %s", name)
				}
				if symbol, ok := payload["symbol"].(string); ok {
					t.Logf("Symbol (from PayloadAsString): %s", symbol)
				}
			}
		}
	}

	// Log trails
	if len(selectedAsset.Trails) > 0 {
		t.Logf("Trails count: %d", len(selectedAsset.Trails))
		for i, trail := range selectedAsset.Trails {
			t.Logf("  Trail[%d]: action=%s, date=%v", i, trail.Action, trail.Date)
		}
	}

	// Log attributes
	if len(selectedAsset.Attributes) > 0 {
		t.Logf("Attributes count: %d", len(selectedAsset.Attributes))
		for i, attr := range selectedAsset.Attributes {
			t.Logf("  Attribute[%d]: key=%s, value=%s", i, attr.Key, attr.Value)
		}
	}

	// Write to file
	outputPath := "/tmp/whitelisted_asset_model.json"
	if err := os.WriteFile(outputPath, rawJSON, 0644); err != nil {
		t.Logf("Warning: could not write to %s: %v", outputPath, err)
	} else {
		t.Logf("\nWritten to: %s", outputPath)
	}
}

// TestIntegration_PrintFixtureGoCode generates Go code for the fixture file.
// This test reads the captured JSON and outputs Go code that can be copied into fixtures.
func TestIntegration_PrintFixtureGoCode(t *testing.T) {
	skipIfNotIntegration(t)

	jsonPath := "/tmp/whitelisted_address_model.json"
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Skipf("No captured data at %s: %v (run CaptureRawWhitelistedAddressResponse first)", jsonPath, err)
	}

	// Parse the JSON to extract key fields
	var envelope map[string]interface{}
	if err := json.Unmarshal(data, &envelope); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Extract and print key values for fixture
	fmt.Println("// === FIXTURE DATA (copy to testdata/whitelisted_address_fixtures.go) ===")
	fmt.Println()

	metadata, _ := envelope["metadata"].(map[string]interface{})
	if metadata != nil {
		if hash, ok := metadata["hash"].(string); ok {
			fmt.Printf("const RealMetadataHash = %q\n", hash)
		}
		if payload, ok := metadata["payloadAsString"].(string); ok {
			fmt.Printf("const RealPayloadAsString = `%s`\n", payload)
		}
	}

	if rc, ok := envelope["rulesContainer"].(string); ok && rc != "" {
		// Truncate for display
		if len(rc) > 100 {
			fmt.Printf("// RulesContainer (base64, length=%d): %s...\n", len(rc), rc[:100])
		}
	}

	if rs, ok := envelope["rulesSignatures"].(string); ok && rs != "" {
		if len(rs) > 100 {
			fmt.Printf("// RulesSignatures (base64, length=%d): %s...\n", len(rs), rs[:100])
		}
	}

	fmt.Println()
	fmt.Println("// Full JSON for reference:")
	fmt.Printf("const RealAPIResponseJSON = `%s`\n", string(data))
}
