package helper

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewWhitelistedAddressVerifier(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	keys := []*ecdsa.PublicKey{&key.PublicKey}

	t.Run("with keys", func(t *testing.T) {
		v := NewWhitelistedAddressVerifier(keys, 1)
		if v == nil {
			t.Error("NewWhitelistedAddressVerifier() returned nil")
		}
		if len(v.superAdminKeys) != 1 {
			t.Errorf("superAdminKeys length = %d, want 1", len(v.superAdminKeys))
		}
		if v.minValidSignatures != 1 {
			t.Errorf("minValidSignatures = %d, want 1", v.minValidSignatures)
		}
	})

	t.Run("with nil keys", func(t *testing.T) {
		v := NewWhitelistedAddressVerifier(nil, 0)
		if v == nil {
			t.Error("NewWhitelistedAddressVerifier() returned nil")
		}
		if len(v.superAdminKeys) != 0 {
			t.Errorf("superAdminKeys length = %d, want 0", len(v.superAdminKeys))
		}
	})
}

func TestVerifyWhitelistedAddress_NilInput(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAddressVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

	t.Run("nil address", func(t *testing.T) {
		_, err := v.VerifyWhitelistedAddress(nil, nil, nil)
		if err == nil {
			t.Error("VerifyWhitelistedAddress() expected error for nil address")
		}
	})

	t.Run("nil metadata", func(t *testing.T) {
		addr := &model.WhitelistedAddress{}
		_, err := v.VerifyWhitelistedAddress(addr, nil, nil)
		if err == nil {
			t.Error("VerifyWhitelistedAddress() expected error for nil metadata")
		}
	})
}

func TestVerifyMetadataHash(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAddressVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

	// Compute expected hash
	payload := `{"currency":"ETH","network":"mainnet","address":"0x123"}`
	expectedHash := crypto.CalculateHexHash(payload)

	tests := []struct {
		name        string
		addr        *model.WhitelistedAddress
		expectError bool
	}{
		{
			name: "valid hash",
			addr: &model.WhitelistedAddress{
				Metadata: &model.WhitelistedAssetMetadata{
					PayloadAsString: payload,
					Hash:            expectedHash,
				},
			},
			expectError: false,
		},
		{
			name: "empty payload",
			addr: &model.WhitelistedAddress{
				Metadata: &model.WhitelistedAssetMetadata{
					PayloadAsString: "",
					Hash:            "somehash",
				},
			},
			expectError: true,
		},
		{
			name: "empty hash",
			addr: &model.WhitelistedAddress{
				Metadata: &model.WhitelistedAssetMetadata{
					PayloadAsString: payload,
					Hash:            "",
				},
			},
			expectError: true,
		},
		{
			name: "mismatched hash",
			addr: &model.WhitelistedAddress{
				Metadata: &model.WhitelistedAssetMetadata{
					PayloadAsString: payload,
					Hash:            "wronghash",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.verifyMetadataHash(tt.addr)
			if tt.expectError && err == nil {
				t.Error("verifyMetadataHash() expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("verifyMetadataHash() unexpected error: %v", err)
			}
		})
	}
}

func TestVerifyRulesContainerSignatures_NoKeys(t *testing.T) {
	v := NewWhitelistedAddressVerifier(nil, 0)

	addr := &model.WhitelistedAddress{
		RulesContainer:  "somecontainer",
		RulesSignatures: "somesigs",
	}

	err := v.verifyRulesContainerSignatures(addr, nil)
	if err == nil {
		t.Error("verifyRulesContainerSignatures() expected error when no keys configured")
	}
}

func TestVerifyRulesContainerSignatures_EmptyFields(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAddressVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

	t.Run("empty rules container", func(t *testing.T) {
		addr := &model.WhitelistedAddress{
			RulesContainer:  "",
			RulesSignatures: "somesigs",
		}
		err := v.verifyRulesContainerSignatures(addr, nil)
		if err == nil {
			t.Error("expected error for empty rules container")
		}
	})

	t.Run("empty rules signatures", func(t *testing.T) {
		addr := &model.WhitelistedAddress{
			RulesContainer:  "somecontainer",
			RulesSignatures: "",
		}
		err := v.verifyRulesContainerSignatures(addr, nil)
		if err == nil {
			t.Error("expected error for empty rules signatures")
		}
	})
}

func TestGetApplicableThresholds(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAddressVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

	defaultThresholds := []*model.SequentialThresholds{
		{Thresholds: []*model.GroupThreshold{{GroupID: "default", MinimumSignatures: 1}}},
	}

	lineThresholds := []*model.SequentialThresholds{
		{Thresholds: []*model.GroupThreshold{{GroupID: "line-specific", MinimumSignatures: 2}}},
	}

	rules := &model.AddressWhitelistingRules{
		ParallelThresholds: defaultThresholds,
		Lines: []*model.AddressWhitelistingLine{
			{
				Cells: []*model.RuleSource{
					{
						Type:           model.RuleSourceTypeInternalWallet,
						InternalWallet: &model.RuleSourceInternalWallet{Path: "m/44/60/0"},
					},
				},
				ParallelThresholds: lineThresholds,
			},
		},
	}

	t.Run("use default when has linked addresses", func(t *testing.T) {
		addr := &model.WhitelistedAddress{
			LinkedInternalAddresses: []model.InternalAddress{{ID: 1}},
			LinkedWallets:           []model.InternalWallet{{ID: 1, Path: "m/44/60/0"}},
		}
		result := v.getApplicableThresholds(rules, addr)
		if len(result) != 1 || result[0].Thresholds[0].GroupID != "default" {
			t.Error("expected default thresholds when linked addresses present")
		}
	})

	t.Run("use default when multiple wallets", func(t *testing.T) {
		addr := &model.WhitelistedAddress{
			LinkedWallets: []model.InternalWallet{
				{ID: 1, Path: "m/44/60/0"},
				{ID: 2, Path: "m/44/60/1"},
			},
		}
		result := v.getApplicableThresholds(rules, addr)
		if len(result) != 1 || result[0].Thresholds[0].GroupID != "default" {
			t.Error("expected default thresholds when multiple wallets present")
		}
	})

	t.Run("use line threshold when matching wallet path", func(t *testing.T) {
		addr := &model.WhitelistedAddress{
			LinkedWallets: []model.InternalWallet{{ID: 1, Path: "m/44/60/0"}},
		}
		result := v.getApplicableThresholds(rules, addr)
		if len(result) != 1 || result[0].Thresholds[0].GroupID != "line-specific" {
			t.Error("expected line-specific thresholds when wallet path matches")
		}
	})

	t.Run("use default when no matching line", func(t *testing.T) {
		addr := &model.WhitelistedAddress{
			LinkedWallets: []model.InternalWallet{{ID: 1, Path: "m/44/60/999"}},
		}
		result := v.getApplicableThresholds(rules, addr)
		if len(result) != 1 || result[0].Thresholds[0].GroupID != "default" {
			t.Error("expected default thresholds when no line matches")
		}
	})
}

func TestMatchesWalletPath(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAddressVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

	tests := []struct {
		name       string
		line       *model.AddressWhitelistingLine
		walletPath string
		expected   bool
	}{
		{
			name:       "empty cells",
			line:       &model.AddressWhitelistingLine{Cells: nil},
			walletPath: "m/44/60/0",
			expected:   false,
		},
		{
			name: "wrong source type",
			line: &model.AddressWhitelistingLine{
				Cells: []*model.RuleSource{{Type: model.RuleSourceTypeUnknown}},
			},
			walletPath: "m/44/60/0",
			expected:   false,
		},
		{
			name: "nil internal wallet",
			line: &model.AddressWhitelistingLine{
				Cells: []*model.RuleSource{{Type: model.RuleSourceTypeInternalWallet}},
			},
			walletPath: "m/44/60/0",
			expected:   false,
		},
		{
			name: "matching path",
			line: &model.AddressWhitelistingLine{
				Cells: []*model.RuleSource{
					{
						Type:           model.RuleSourceTypeInternalWallet,
						InternalWallet: &model.RuleSourceInternalWallet{Path: "m/44/60/0"},
					},
				},
			},
			walletPath: "m/44/60/0",
			expected:   true,
		},
		{
			name: "non-matching path",
			line: &model.AddressWhitelistingLine{
				Cells: []*model.RuleSource{
					{
						Type:           model.RuleSourceTypeInternalWallet,
						InternalWallet: &model.RuleSourceInternalWallet{Path: "m/44/60/1"},
					},
				},
			},
			walletPath: "m/44/60/0",
			expected:   false,
		},
		{
			name: "empty wallet path",
			line: &model.AddressWhitelistingLine{
				Cells: []*model.RuleSource{
					{
						Type:           model.RuleSourceTypeInternalWallet,
						InternalWallet: &model.RuleSourceInternalWallet{Path: "m/44/60/0"},
					},
				},
			},
			walletPath: "",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.matchesWalletPath(tt.line, tt.walletPath)
			if result != tt.expected {
				t.Errorf("matchesWalletPath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVerifySequentialThresholds_NilInput(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAddressVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

	t.Run("nil threshold", func(t *testing.T) {
		err := v.verifySequentialThresholds(nil, nil, nil, "hash", nil)
		if err == nil {
			t.Error("expected error for nil threshold")
		}
	})

	t.Run("empty thresholds", func(t *testing.T) {
		err := v.verifySequentialThresholds(&model.SequentialThresholds{}, nil, nil, "hash", nil)
		if err == nil {
			t.Error("expected error for empty thresholds")
		}
	})
}

func TestContainsHash(t *testing.T) {
	tests := []struct {
		name     string
		hashes   []string
		hash     string
		expected bool
	}{
		{"found", []string{"a", "b", "c"}, "b", true},
		{"not found", []string{"a", "b", "c"}, "d", false},
		{"empty list", []string{}, "a", false},
		{"nil list", nil, "a", false},
		{"empty hash", []string{"a", "b"}, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsHash(tt.hashes, tt.hash)
			if result != tt.expected {
				t.Errorf("containsHash() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVerifyHashInSignedHashes_NoSignedAddress(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAddressVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

	t.Run("nil signed address", func(t *testing.T) {
		addr := &model.WhitelistedAddress{
			SignedAddress: nil,
			Metadata:      &model.WhitelistedAssetMetadata{Hash: "hash"},
		}
		_, err := v.verifyHashInSignedHashes(addr)
		if err == nil {
			t.Error("expected error for nil signed address")
		}
	})

	t.Run("empty signatures", func(t *testing.T) {
		addr := &model.WhitelistedAddress{
			SignedAddress: &model.SignedWhitelistedAddress{Signatures: nil},
			Metadata:      &model.WhitelistedAssetMetadata{Hash: "hash"},
		}
		_, err := v.verifyHashInSignedHashes(addr)
		if err == nil {
			t.Error("expected error for empty signatures")
		}
	})
}

func TestVerifyHashInSignedHashes_HashCoverage(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAddressVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

	payload := `{"currency":"ETH"}`
	correctHash := crypto.CalculateHexHash(payload)

	t.Run("hash found", func(t *testing.T) {
		addr := &model.WhitelistedAddress{
			SignedAddress: &model.SignedWhitelistedAddress{
				Signatures: []model.WhitelistSignature{
					{Hashes: []string{correctHash}},
				},
			},
			Metadata: &model.WhitelistedAssetMetadata{
				Hash:            correctHash,
				PayloadAsString: payload,
			},
		}
		foundHash, err := v.verifyHashInSignedHashes(addr)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if foundHash != correctHash {
			t.Errorf("foundHash = %q, want %q", foundHash, correctHash)
		}
	})

	t.Run("hash not found", func(t *testing.T) {
		addr := &model.WhitelistedAddress{
			SignedAddress: &model.SignedWhitelistedAddress{
				Signatures: []model.WhitelistSignature{
					{Hashes: []string{"otherhash"}},
				},
			},
			Metadata: &model.WhitelistedAssetMetadata{
				Hash:            correctHash,
				PayloadAsString: payload,
			},
		}
		_, err := v.verifyHashInSignedHashes(addr)
		if err == nil {
			t.Error("expected error when hash not found")
		}
	})
}
