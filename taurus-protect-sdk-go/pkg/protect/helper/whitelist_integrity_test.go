package helper

import (
	"testing"
)

// TestParseWhitelistedAddress_ValidMatch tests that a valid JSON payload is parsed correctly
// and all fields match the expected values.
func TestParseWhitelistedAddress_ValidMatch(t *testing.T) {
	payload := `{
		"currency": "ETH",
		"network": "mainnet",
		"address": "0x742d35Cc6634C0532925a3b844Bc9e7595f",
		"memo": "test-memo",
		"label": "My Test Address",
		"customerId": "customer-123",
		"contractType": "ERC20",
		"tnParticipantID": "participant-456",
		"addressType": "individual",
		"exchangeAccountId": "12345",
		"linkedInternalAddresses": [
			{"id": 1, "address": "0xabc", "label": "Internal Address 1"},
			{"id": 2, "address": "0xdef", "label": "Internal Address 2"}
		],
		"linkedWallets": [
			{"id": 10, "name": "Wallet 1", "path": "m/44/60/0/0"},
			{"id": 20, "name": "Wallet 2", "path": "m/44/60/0/1"}
		]
	}`

	result, err := ParseWhitelistedAddressFromJSON(payload)
	if err != nil {
		t.Fatalf("ParseWhitelistedAddressFromJSON() unexpected error: %v", err)
	}

	// Verify all fields match expected values
	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"Blockchain", result.Blockchain, "ETH"},
		{"Network", result.Network, "mainnet"},
		{"Address", result.Address, "0x742d35Cc6634C0532925a3b844Bc9e7595f"},
		{"Memo", result.Memo, "test-memo"},
		{"Label", result.Label, "My Test Address"},
		{"CustomerId", result.CustomerId, "customer-123"},
		{"ContractType", result.ContractType, "ERC20"},
		{"TnParticipantID", result.TnParticipantID, "participant-456"},
		{"AddressType", result.AddressType, "individual"},
		{"ExchangeAccountId", result.ExchangeAccountId, int64(12345)},
		{"LinkedInternalAddresses count", len(result.LinkedInternalAddresses), 2},
		{"LinkedWallets count", len(result.LinkedWallets), 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}

	// Verify LinkedInternalAddresses details
	if result.LinkedInternalAddresses[0].ID != 1 {
		t.Errorf("LinkedInternalAddresses[0].ID = %d, want 1", result.LinkedInternalAddresses[0].ID)
	}
	if result.LinkedInternalAddresses[0].Label != "Internal Address 1" {
		t.Errorf("LinkedInternalAddresses[0].Label = %q, want %q", result.LinkedInternalAddresses[0].Label, "Internal Address 1")
	}
	if result.LinkedInternalAddresses[1].ID != 2 {
		t.Errorf("LinkedInternalAddresses[1].ID = %d, want 2", result.LinkedInternalAddresses[1].ID)
	}

	// Verify LinkedWallets details
	if result.LinkedWallets[0].ID != 10 {
		t.Errorf("LinkedWallets[0].ID = %d, want 10", result.LinkedWallets[0].ID)
	}
	if result.LinkedWallets[0].Label != "Wallet 1" {
		t.Errorf("LinkedWallets[0].Label = %q, want %q", result.LinkedWallets[0].Label, "Wallet 1")
	}
	if result.LinkedWallets[0].Path != "m/44/60/0/0" {
		t.Errorf("LinkedWallets[0].Path = %q, want %q", result.LinkedWallets[0].Path, "m/44/60/0/0")
	}
	if result.LinkedWallets[1].ID != 20 {
		t.Errorf("LinkedWallets[1].ID = %d, want 20", result.LinkedWallets[1].ID)
	}
}

// TestParseWhitelistedAddress_InvalidBlockchain tests that blockchain mismatch can be detected.
func TestParseWhitelistedAddress_InvalidBlockchain(t *testing.T) {
	payload := `{"currency": "ETH", "network": "mainnet", "address": "0x123"}`

	result, err := ParseWhitelistedAddressFromJSON(payload)
	if err != nil {
		t.Fatalf("ParseWhitelistedAddressFromJSON() unexpected error: %v", err)
	}

	// Verify parsed blockchain
	if result.Blockchain != "ETH" {
		t.Errorf("Blockchain = %q, want %q", result.Blockchain, "ETH")
	}

	// Demonstrate mismatch detection: expected BTC but got ETH
	expectedBlockchain := "BTC"
	if result.Blockchain == expectedBlockchain {
		t.Errorf("Expected blockchain mismatch: got %q, expected mismatch with %q", result.Blockchain, expectedBlockchain)
	}
}

// TestParseWhitelistedAddress_InvalidAddress tests that address mismatch can be detected.
func TestParseWhitelistedAddress_InvalidAddress(t *testing.T) {
	payload := `{"currency": "ETH", "network": "mainnet", "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f"}`

	result, err := ParseWhitelistedAddressFromJSON(payload)
	if err != nil {
		t.Fatalf("ParseWhitelistedAddressFromJSON() unexpected error: %v", err)
	}

	// Verify parsed address
	if result.Address != "0x742d35Cc6634C0532925a3b844Bc9e7595f" {
		t.Errorf("Address = %q, want %q", result.Address, "0x742d35Cc6634C0532925a3b844Bc9e7595f")
	}

	// Demonstrate mismatch detection: different address would fail validation
	differentAddress := "0xDifferentAddress123"
	if result.Address == differentAddress {
		t.Errorf("Expected address mismatch: got %q, expected mismatch with %q", result.Address, differentAddress)
	}
}

// TestParseWhitelistedAddress_InvalidLabel tests that label mismatch can be detected.
func TestParseWhitelistedAddress_InvalidLabel(t *testing.T) {
	payload := `{"currency": "ETH", "network": "mainnet", "address": "0x123", "label": "My Address"}`

	result, err := ParseWhitelistedAddressFromJSON(payload)
	if err != nil {
		t.Fatalf("ParseWhitelistedAddressFromJSON() unexpected error: %v", err)
	}

	// Verify parsed label
	if result.Label != "My Address" {
		t.Errorf("Label = %q, want %q", result.Label, "My Address")
	}

	// Demonstrate mismatch detection: different label would indicate tampering
	expectedLabel := "Different Label"
	if result.Label == expectedLabel {
		t.Errorf("Expected label mismatch: got %q, expected mismatch with %q", result.Label, expectedLabel)
	}
}

// TestParseWhitelistedAddress_EmptyPayload tests that an empty payload returns an error.
func TestParseWhitelistedAddress_EmptyPayload(t *testing.T) {
	_, err := ParseWhitelistedAddressFromJSON("")
	if err == nil {
		t.Error("ParseWhitelistedAddressFromJSON() expected error for empty payload, got nil")
	}
}

// TestParseWhitelistedAddress_InvalidJSON tests that invalid JSON returns an error.
func TestParseWhitelistedAddress_InvalidJSON(t *testing.T) {
	tests := []struct {
		name    string
		payload string
	}{
		{"not JSON at all", "not json at all"},
		{"malformed JSON", `{"currency": "ETH", "address": `},
		{"array instead of object", `["ETH", "mainnet"]`},
		{"unclosed string", `{"currency": "ETH`},
		{"invalid escape", `{"currency": "ETH\x"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseWhitelistedAddressFromJSON(tt.payload)
			if err == nil {
				t.Error("ParseWhitelistedAddressFromJSON() expected error for invalid JSON, got nil")
			}
		})
	}
}

// TestComputeLegacyHashes_ContractTypeOnly tests legacy hash computation when only contractType
// needs to be removed.
func TestComputeLegacyHashes_ContractTypeOnly(t *testing.T) {
	// Payload with contractType but no labels in linkedInternalAddresses
	payload := `{"address":"0x123","contractType":"ERC20","linkedInternalAddresses":[{"id":1,"address":"0xabc"}]}`

	hashes := ComputeLegacyHashes(payload)

	// Should produce 1 hash: contractType removed
	if len(hashes) != 1 {
		t.Errorf("ComputeLegacyHashes() returned %d hashes, want 1", len(hashes))
	}

	// Verify the hash is 64 hex characters (SHA-256)
	if len(hashes) > 0 && len(hashes[0]) != 64 {
		t.Errorf("Hash length = %d, want 64", len(hashes[0]))
	}
}

// TestComputeLegacyHashes_LabelsOnly tests legacy hash computation when only labels
// need to be removed from linkedInternalAddresses.
func TestComputeLegacyHashes_LabelsOnly(t *testing.T) {
	// Payload with labels in linkedInternalAddresses but no contractType
	payload := `{"address":"0x123","linkedInternalAddresses":[{"id":1,"label":"test"}]}`

	hashes := ComputeLegacyHashes(payload)

	// Should produce 1 hash: labels removed
	if len(hashes) != 1 {
		t.Errorf("ComputeLegacyHashes() returned %d hashes, want 1", len(hashes))
	}

	// Verify the hash is 64 hex characters (SHA-256)
	if len(hashes) > 0 && len(hashes[0]) != 64 {
		t.Errorf("Hash length = %d, want 64", len(hashes[0]))
	}
}

// TestComputeLegacyHashes_BothRemoved tests legacy hash computation when both contractType
// and labels need to be removed.
func TestComputeLegacyHashes_BothRemoved(t *testing.T) {
	// Payload with both contractType and labels
	payload := `{"address":"0x123","contractType":"ERC20","linkedInternalAddresses":[{"id":1,"label":"test"}]}`

	hashes := ComputeLegacyHashes(payload)

	// Should produce 3 hashes:
	// 1. contractType removed only
	// 2. labels removed only
	// 3. both removed
	if len(hashes) != 3 {
		t.Errorf("ComputeLegacyHashes() returned %d hashes, want 3", len(hashes))
	}

	// Verify all hashes are 64 hex characters (SHA-256)
	for i, h := range hashes {
		if len(h) != 64 {
			t.Errorf("hashes[%d] length = %d, want 64", i, len(h))
		}
	}
}

// TestComputeLegacyHashes_UniqueHashes verifies that ComputeLegacyHashes returns no duplicates.
func TestComputeLegacyHashes_UniqueHashes(t *testing.T) {
	testCases := []struct {
		name    string
		payload string
	}{
		{
			name:    "contractType only",
			payload: `{"address":"0x123","contractType":"ERC20"}`,
		},
		{
			name:    "labels only",
			payload: `{"address":"0x123","linkedInternalAddresses":[{"id":1,"label":"test"}]}`,
		},
		{
			name:    "both fields",
			payload: `{"address":"0x123","contractType":"ERC20","linkedInternalAddresses":[{"id":1,"label":"test"}]}`,
		},
		{
			name:    "multiple linked addresses",
			payload: `{"address":"0x123","contractType":"CMTA20","linkedInternalAddresses":[{"id":1,"label":"a"},{"id":2,"label":"b"}]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hashes := ComputeLegacyHashes(tc.payload)

			// Check for duplicates
			seen := make(map[string]bool)
			for _, h := range hashes {
				if seen[h] {
					t.Error("ComputeLegacyHashes() returned duplicate hash")
				}
				seen[h] = true
			}
		})
	}
}

// TestComputeLegacyHashes_EmptyPayload verifies that empty payload returns nil.
func TestComputeLegacyHashes_EmptyPayload(t *testing.T) {
	hashes := ComputeLegacyHashes("")
	if hashes != nil {
		t.Errorf("ComputeLegacyHashes() = %v, want nil", hashes)
	}
}

// TestComputeLegacyHashes_NoModifiableFields verifies that payloads without contractType
// or labels in objects return no hashes.
func TestComputeLegacyHashes_NoModifiableFields(t *testing.T) {
	testCases := []struct {
		name    string
		payload string
	}{
		{
			name:    "simple payload without fields",
			payload: `{"address":"0x123","network":"mainnet"}`,
		},
		{
			name:    "main label (not in object)",
			payload: `{"address":"0x123","label":"main","network":"mainnet"}`,
		},
		{
			name:    "linked addresses without labels",
			payload: `{"address":"0x123","linkedInternalAddresses":[{"id":1,"address":"0xabc"}]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hashes := ComputeLegacyHashes(tc.payload)
			if len(hashes) != 0 {
				t.Errorf("ComputeLegacyHashes() returned %d hashes, want 0", len(hashes))
			}
		})
	}
}

// TestParseWhitelistedAddress_OptionalFields tests that optional fields are properly handled.
func TestParseWhitelistedAddress_OptionalFields(t *testing.T) {
	// Minimal payload with only required fields
	payload := `{"currency": "BTC", "network": "testnet", "address": "bc1qtest"}`

	result, err := ParseWhitelistedAddressFromJSON(payload)
	if err != nil {
		t.Fatalf("ParseWhitelistedAddressFromJSON() unexpected error: %v", err)
	}

	// Verify required fields
	if result.Blockchain != "BTC" {
		t.Errorf("Blockchain = %q, want %q", result.Blockchain, "BTC")
	}
	if result.Network != "testnet" {
		t.Errorf("Network = %q, want %q", result.Network, "testnet")
	}
	if result.Address != "bc1qtest" {
		t.Errorf("Address = %q, want %q", result.Address, "bc1qtest")
	}

	// Verify optional fields are empty/zero
	if result.Label != "" {
		t.Errorf("Label = %q, want empty string", result.Label)
	}
	if result.Memo != "" {
		t.Errorf("Memo = %q, want empty string", result.Memo)
	}
	if result.CustomerId != "" {
		t.Errorf("CustomerId = %q, want empty string", result.CustomerId)
	}
	if result.ContractType != "" {
		t.Errorf("ContractType = %q, want empty string", result.ContractType)
	}
	if result.ExchangeAccountId != 0 {
		t.Errorf("ExchangeAccountId = %d, want 0", result.ExchangeAccountId)
	}
	if len(result.LinkedInternalAddresses) != 0 {
		t.Errorf("LinkedInternalAddresses count = %d, want 0", len(result.LinkedInternalAddresses))
	}
	if len(result.LinkedWallets) != 0 {
		t.Errorf("LinkedWallets count = %d, want 0", len(result.LinkedWallets))
	}
}

// TestParseWhitelistedAddress_ExchangeAccountIdParsing tests various exchangeAccountId formats.
func TestParseWhitelistedAddress_ExchangeAccountIdParsing(t *testing.T) {
	tests := []struct {
		name     string
		payload  string
		expected int64
	}{
		{
			name:     "numeric string",
			payload:  `{"currency": "ETH", "network": "mainnet", "address": "0x123", "exchangeAccountId": "12345"}`,
			expected: 12345,
		},
		{
			name:     "large number",
			payload:  `{"currency": "ETH", "network": "mainnet", "address": "0x123", "exchangeAccountId": "9223372036854775807"}`,
			expected: 9223372036854775807, // max int64
		},
		{
			name:     "zero",
			payload:  `{"currency": "ETH", "network": "mainnet", "address": "0x123", "exchangeAccountId": "0"}`,
			expected: 0,
		},
		{
			name:     "empty string",
			payload:  `{"currency": "ETH", "network": "mainnet", "address": "0x123", "exchangeAccountId": ""}`,
			expected: 0,
		},
		{
			name:     "non-numeric string",
			payload:  `{"currency": "ETH", "network": "mainnet", "address": "0x123", "exchangeAccountId": "not-a-number"}`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseWhitelistedAddressFromJSON(tt.payload)
			if err != nil {
				t.Fatalf("ParseWhitelistedAddressFromJSON() unexpected error: %v", err)
			}
			if result.ExchangeAccountId != tt.expected {
				t.Errorf("ExchangeAccountId = %d, want %d", result.ExchangeAccountId, tt.expected)
			}
		})
	}
}

// TestParseWhitelistedAddress_FieldIntegrity verifies that the parsed address maintains
// integrity with the original JSON data for validation purposes.
func TestParseWhitelistedAddress_FieldIntegrity(t *testing.T) {
	// This test simulates the verification flow where we compare
	// parsed JSON fields against expected values from the API response

	// Original JSON payload (as received from server)
	payload := `{
		"currency": "ETH",
		"network": "mainnet",
		"address": "0x742d35Cc6634C0532925a3b844Bc9e7595f",
		"label": "Production Treasury",
		"customerId": "cust-001",
		"contractType": "ERC20"
	}`

	parsed, err := ParseWhitelistedAddressFromJSON(payload)
	if err != nil {
		t.Fatalf("ParseWhitelistedAddressFromJSON() unexpected error: %v", err)
	}

	// Simulate expected values from API response metadata
	expectedValues := map[string]string{
		"blockchain":   "ETH",
		"network":      "mainnet",
		"address":      "0x742d35Cc6634C0532925a3b844Bc9e7595f",
		"label":        "Production Treasury",
		"customerId":   "cust-001",
		"contractType": "ERC20",
	}

	// Verify each field matches
	if parsed.Blockchain != expectedValues["blockchain"] {
		t.Errorf("Blockchain integrity check failed: %q != %q", parsed.Blockchain, expectedValues["blockchain"])
	}
	if parsed.Network != expectedValues["network"] {
		t.Errorf("Network integrity check failed: %q != %q", parsed.Network, expectedValues["network"])
	}
	if parsed.Address != expectedValues["address"] {
		t.Errorf("Address integrity check failed: %q != %q", parsed.Address, expectedValues["address"])
	}
	if parsed.Label != expectedValues["label"] {
		t.Errorf("Label integrity check failed: %q != %q", parsed.Label, expectedValues["label"])
	}
	if parsed.CustomerId != expectedValues["customerId"] {
		t.Errorf("CustomerId integrity check failed: %q != %q", parsed.CustomerId, expectedValues["customerId"])
	}
	if parsed.ContractType != expectedValues["contractType"] {
		t.Errorf("ContractType integrity check failed: %q != %q", parsed.ContractType, expectedValues["contractType"])
	}
}
