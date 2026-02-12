package helper

import (
	"strings"
	"testing"
)

func TestComputeLegacyHashes(t *testing.T) {
	tests := []struct {
		name           string
		payload        string
		expectHashes   int
		containsCheck  string // If set, at least one hash should be computed from a payload containing this string removed
	}{
		{
			name:         "empty payload",
			payload:      "",
			expectHashes: 0,
		},
		{
			name:         "no contract type or labels",
			payload:      `{"currency":"ETH","network":"mainnet","address":"0x123"}`,
			expectHashes: 0,
		},
		{
			name:           "with contract type",
			payload:        `{"currency":"ETH","address":"0x123","contractType":"ERC20"}`,
			expectHashes:   1,
			containsCheck:  "contractType",
		},
		{
			name:           "with labels in linked addresses",
			payload:        `{"currency":"ETH","linkedInternalAddresses":[{"id":1,"label":"test"}]}`,
			expectHashes:   1,
			containsCheck:  `"label"`,
		},
		{
			name:          "with both contract type and labels",
			payload:       `{"currency":"ETH","contractType":"ERC20","linkedInternalAddresses":[{"id":1,"label":"test"}]}`,
			expectHashes:  3, // contractType only, labels only, both removed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashes := ComputeLegacyHashes(tt.payload)
			if len(hashes) != tt.expectHashes {
				t.Errorf("ComputeLegacyHashes() returned %d hashes, want %d", len(hashes), tt.expectHashes)
			}

			// Check for uniqueness
			seen := make(map[string]bool)
			for _, h := range hashes {
				if seen[h] {
					t.Error("ComputeLegacyHashes() returned duplicate hash")
				}
				seen[h] = true

				// Verify hashes are hex strings of expected length (SHA-256 = 64 hex chars)
				if len(h) != 64 {
					t.Errorf("Hash length = %d, want 64", len(h))
				}
			}
		})
	}
}

func TestParseWhitelistedAddressFromJSON(t *testing.T) {
	tests := []struct {
		name        string
		payload     string
		expectError bool
		checkAddr   func(t *testing.T, result interface{})
	}{
		{
			name:        "empty payload",
			payload:     "",
			expectError: true,
		},
		{
			name:        "invalid JSON",
			payload:     "not json",
			expectError: true,
		},
		{
			name:        "minimal valid payload",
			payload:     `{"currency":"ETH","network":"mainnet","address":"0x123"}`,
			expectError: false,
			checkAddr: func(t *testing.T, result interface{}) {
				addr := result.(*testWhitelistedAddress)
				if addr.Blockchain != "ETH" {
					t.Errorf("Blockchain = %q, want %q", addr.Blockchain, "ETH")
				}
				if addr.Network != "mainnet" {
					t.Errorf("Network = %q, want %q", addr.Network, "mainnet")
				}
				if addr.Address != "0x123" {
					t.Errorf("Address = %q, want %q", addr.Address, "0x123")
				}
			},
		},
		{
			name: "full payload",
			payload: `{
				"currency": "ETH",
				"network": "mainnet",
				"address": "0x742d35Cc6634C0532925a3b844Bc9e7595f",
				"memo": "test memo",
				"label": "My Address",
				"customerId": "customer-123",
				"contractType": "ERC20",
				"tnParticipantID": "participant-456",
				"addressType": "individual",
				"exchangeAccountId": "12345",
				"linkedInternalAddresses": [
					{"id": 1, "address": "0xabc", "label": "Internal 1"}
				],
				"linkedWallets": [
					{"id": 10, "name": "Wallet 1", "path": "m/44/60/0/0"}
				]
			}`,
			expectError: false,
			checkAddr: func(t *testing.T, result interface{}) {
				addr := result.(*testWhitelistedAddress)
				if addr.Label != "My Address" {
					t.Errorf("Label = %q, want %q", addr.Label, "My Address")
				}
				if addr.ContractType != "ERC20" {
					t.Errorf("ContractType = %q, want %q", addr.ContractType, "ERC20")
				}
				if len(addr.LinkedInternalAddresses) != 1 {
					t.Errorf("LinkedInternalAddresses count = %d, want 1", len(addr.LinkedInternalAddresses))
				}
				if len(addr.LinkedWallets) != 1 {
					t.Errorf("LinkedWallets count = %d, want 1", len(addr.LinkedWallets))
				}
				if addr.LinkedWallets[0].Path != "m/44/60/0/0" {
					t.Errorf("LinkedWallets[0].Path = %q, want %q", addr.LinkedWallets[0].Path, "m/44/60/0/0")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseWhitelistedAddressFromJSON(tt.payload)
			if tt.expectError {
				if err == nil {
					t.Error("ParseWhitelistedAddressFromJSON() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("ParseWhitelistedAddressFromJSON() unexpected error: %v", err)
				return
			}

			if tt.checkAddr != nil {
				// Convert to our test wrapper type
				var linkedWallets []testInternalWallet
				for _, w := range result.LinkedWallets {
					linkedWallets = append(linkedWallets, testInternalWallet{
						ID:    w.ID,
						Path:  w.Path,
						Label: w.Label,
					})
				}

				var linkedAddrs []interface{}
				for _, a := range result.LinkedInternalAddresses {
					linkedAddrs = append(linkedAddrs, a)
				}

				wrapper := &testWhitelistedAddress{
					Blockchain:              result.Blockchain,
					Network:                 result.Network,
					Address:                 result.Address,
					Label:                   result.Label,
					ContractType:            result.ContractType,
					LinkedInternalAddresses: linkedAddrs,
					LinkedWallets:           linkedWallets,
				}
				tt.checkAddr(t, wrapper)
			}
		})
	}
}

// testWhitelistedAddress is a wrapper for testing
type testWhitelistedAddress struct {
	Blockchain              string
	Network                 string
	Address                 string
	Label                   string
	ContractType            string
	LinkedInternalAddresses []interface{}
	LinkedWallets           []testInternalWallet
}

type testInternalWallet struct {
	ID    int64
	Path  string
	Label string
}

func TestContractTypePattern(t *testing.T) {
	// The contractTypePattern matches ,\"contractType\":\"...\" (note: comma BEFORE the field)
	// This only removes contractType when it appears AFTER other fields
	tests := []struct {
		input    string
		expected string
	}{
		{
			// contractType after another field - should be removed
			input:    `{"address":"0x123","contractType":"ERC20"}`,
			expected: `{"address":"0x123"}`,
		},
		{
			// contractType first - comma pattern won't match, field stays
			input:    `{"contractType":"CMTA20","address":"0x123"}`,
			expected: `{"contractType":"CMTA20","address":"0x123"}`,
		},
		{
			// No contractType at all
			input:    `{"address":"0x123"}`,
			expected: `{"address":"0x123"}`,
		},
		{
			// contractType in middle - should be removed
			input:    `{"address":"0x123","contractType":"ERC20","network":"mainnet"}`,
			expected: `{"address":"0x123","network":"mainnet"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := contractTypePattern.ReplaceAllString(tt.input, "")
			if result != tt.expected {
				t.Errorf("contractTypePattern result = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestLabelInObjectPattern(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			// Label inside an object (followed by })
			input:    `{"id":1,"label":"test"}`,
			expected: `{"id":1}`,
		},
		{
			// Multiple objects with labels
			input:    `[{"id":1,"label":"a"},{"id":2,"label":"b"}]`,
			expected: `[{"id":1},{"id":2}]`,
		},
		{
			// Main label not in object (followed by other fields) - should not match
			input:    `{"label":"main","address":"0x123"}`,
			expected: `{"label":"main","address":"0x123"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := labelInObjectPattern.ReplaceAllString(tt.input, "}")
			if result != tt.expected {
				t.Errorf("labelInObjectPattern result = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestComputeAssetLegacyHashes(t *testing.T) {
	tests := []struct {
		name           string
		payload        string
		expectHashes   int
		containsCheck  string // If set, at least one hash should be computed from a payload containing this string removed
	}{
		{
			name:         "empty payload",
			payload:      "",
			expectHashes: 0,
		},
		{
			name:         "no isNFT or kindType",
			payload:      `{"currency":"ETH","network":"mainnet","address":"0x123"}`,
			expectHashes: 0,
		},
		{
			name:           "with isNFT only",
			payload:        `{"currency":"ETH","address":"0x123","isNFT":true}`,
			expectHashes:   1,
			containsCheck:  "isNFT",
		},
		{
			name:           "with isNFT false",
			payload:        `{"currency":"ETH","address":"0x123","isNFT":false}`,
			expectHashes:   1,
			containsCheck:  "isNFT",
		},
		{
			name:           "with kindType only",
			payload:        `{"currency":"ETH","kindType":"fungible","address":"0x123"}`,
			expectHashes:   1,
			containsCheck:  "kindType",
		},
		{
			name:         "with both isNFT and kindType",
			payload:      `{"currency":"ETH","isNFT":true,"kindType":"nft","address":"0x123"}`,
			expectHashes: 3, // isNFT only, kindType only, both removed
		},
		{
			name:         "with isNFT at start (trailing comma pattern)",
			payload:      `{"isNFT":true,"currency":"ETH","address":"0x123"}`,
			expectHashes: 1,
		},
		{
			name:         "with kindType at start (trailing comma pattern)",
			payload:      `{"kindType":"fungible","currency":"ETH","address":"0x123"}`,
			expectHashes: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashes := ComputeAssetLegacyHashes(tt.payload)
			if len(hashes) != tt.expectHashes {
				t.Errorf("ComputeAssetLegacyHashes() returned %d hashes, want %d", len(hashes), tt.expectHashes)
			}

			// Check for uniqueness
			seen := make(map[string]bool)
			for _, h := range hashes {
				if seen[h] {
					t.Error("ComputeAssetLegacyHashes() returned duplicate hash")
				}
				seen[h] = true

				// Verify hashes are hex strings of expected length (SHA-256 = 64 hex chars)
				if len(h) != 64 {
					t.Errorf("Hash length = %d, want 64", len(h))
				}
			}
		})
	}
}

func TestIsNFTPattern(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			// isNFT after another field - should be removed with leading comma
			input:    `{"address":"0x123","isNFT":true}`,
			expected: `{"address":"0x123"}`,
		},
		{
			// isNFT first - trailing comma pattern
			input:    `{"isNFT":true,"address":"0x123"}`,
			expected: `{"address":"0x123"}`,
		},
		{
			// isNFT false
			input:    `{"address":"0x123","isNFT":false}`,
			expected: `{"address":"0x123"}`,
		},
		{
			// isNFT in middle
			input:    `{"address":"0x123","isNFT":true,"network":"mainnet"}`,
			expected: `{"address":"0x123","network":"mainnet"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isNFTPatternWithLeadingComma.ReplaceAllString(tt.input, "")
			result = isNFTPatternWithTrailingComma.ReplaceAllString(result, "")
			if result != tt.expected {
				t.Errorf("isNFT patterns result = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestKindTypePattern(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			// kindType after another field - should be removed with leading comma
			input:    `{"address":"0x123","kindType":"fungible"}`,
			expected: `{"address":"0x123"}`,
		},
		{
			// kindType first - trailing comma pattern
			input:    `{"kindType":"nft","address":"0x123"}`,
			expected: `{"address":"0x123"}`,
		},
		{
			// kindType in middle
			input:    `{"address":"0x123","kindType":"fungible","network":"mainnet"}`,
			expected: `{"address":"0x123","network":"mainnet"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := kindTypePatternWithLeadingComma.ReplaceAllString(tt.input, "")
			result = kindTypePatternWithTrailingComma.ReplaceAllString(result, "")
			if result != tt.expected {
				t.Errorf("kindType patterns result = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestCheckHashesSignature(t *testing.T) {
	t.Run("empty hashes", func(t *testing.T) {
		err := CheckHashesSignature(nil, "sig", nil)
		if err == nil || !strings.Contains(err.Error(), "empty") {
			t.Error("CheckHashesSignature() expected error for empty hashes")
		}
	})

	t.Run("empty signature", func(t *testing.T) {
		err := CheckHashesSignature([]string{"hash"}, "", nil)
		if err == nil || !strings.Contains(err.Error(), "empty") {
			t.Error("CheckHashesSignature() expected error for empty signature")
		}
	})

	t.Run("unsupported key type", func(t *testing.T) {
		err := CheckHashesSignature([]string{"hash"}, "sig", "not a key")
		if err == nil || !strings.Contains(err.Error(), "unsupported") {
			t.Error("CheckHashesSignature() expected error for unsupported key type")
		}
	})
}

