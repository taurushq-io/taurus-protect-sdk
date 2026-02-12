package helper

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// =============================================================================
// Helpers
// =============================================================================

func generateAssetTestKeyPair(t *testing.T) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	t.Helper()
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key pair: %v", err)
	}
	return priv, &priv.PublicKey
}

func publicKeyToPEM(t *testing.T, pub *ecdsa.PublicKey) string {
	t.Helper()
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}))
}

func encodeRulesContainerJSON(t *testing.T, rules interface{}) string {
	t.Helper()
	jsonBytes, err := json.Marshal(rules)
	if err != nil {
		t.Fatalf("failed to marshal rules: %v", err)
	}
	return base64.StdEncoding.EncodeToString(jsonBytes)
}

type assetTestFixture struct {
	saPriv    *ecdsa.PrivateKey
	saPub     *ecdsa.PublicKey
	userPriv  *ecdsa.PrivateKey
	userPub   *ecdsa.PublicKey
	verifier  *WhitelistedAssetVerifier
	asset     *model.WhitelistedAsset
	rcDecoder func(string) (*model.DecodedRulesContainer, error)
	usDecoder func(string) ([]*model.RuleUserSignature, error)
}

func buildAssetTestFixture(t *testing.T) *assetTestFixture {
	t.Helper()
	saPriv, saPub := generateAssetTestKeyPair(t)
	userPriv, userPub := generateAssetTestKeyPair(t)

	// Build payload
	payload := `{"blockchain":"ETH","network":"mainnet","contractAddress":"0xUSDC","name":"USDC","symbol":"USDC","decimals":6}`
	metadataHash := crypto.CalculateHexHash(payload)

	// Build rules container
	userPubPEM := publicKeyToPEM(t, userPub)
	rulesB64 := encodeRulesContainerJSON(t, map[string]interface{}{
		"users": []map[string]interface{}{
			{"id": "user1@bank.com", "publicKey": userPubPEM, "roles": []string{"USER"}},
		},
		"groups": []map[string]interface{}{
			{"id": "approvers", "userIds": []string{"user1@bank.com"}},
		},
		"contractAddressWhitelistingRules": []map[string]interface{}{
			{
				"blockchain": "ETH",
				"network":    "mainnet",
				"parallelThresholds": []map[string]interface{}{
					{"groupId": "approvers", "minimumSignatures": 1},
				},
			},
		},
	})

	rulesData, _ := base64.StdEncoding.DecodeString(rulesB64)

	// Sign rules container with SuperAdmin key
	saSig, err := crypto.SignData(saPriv, rulesData)
	if err != nil {
		t.Fatalf("failed to sign rules: %v", err)
	}

	// Sign hashes array with user key
	hashes := []string{metadataHash}
	hashesJSON, _ := json.Marshal(hashes)
	userSig, err := crypto.SignData(userPriv, hashesJSON)
	if err != nil {
		t.Fatalf("failed to sign hashes: %v", err)
	}

	asset := &model.WhitelistedAsset{
		ID:         "asset-1",
		Blockchain: "ETH",
		Network:    "mainnet",
		Metadata: &model.WhitelistedAssetMetadata{
			Hash:            metadataHash,
			PayloadAsString: payload,
		},
		RulesContainer:  rulesB64,
		RulesSignatures: base64.StdEncoding.EncodeToString([]byte("dummy")),
		SignedContractAddress: &model.SignedContractAddress{
			Signatures: []model.WhitelistSignature{
				{
					UserSignature: &model.WhitelistUserSignature{
						UserID:    "user1@bank.com",
						Signature: userSig,
					},
					Hashes: hashes,
				},
			},
		},
	}

	rcDecoder := func(b64 string) (*model.DecodedRulesContainer, error) {
		return &model.DecodedRulesContainer{
			Users: []*model.RuleUser{
				{
					ID:        "user1@bank.com",
					PublicKey: userPub,
					Roles:     []string{"USER"},
				},
			},
			Groups: []*model.RuleGroup{
				{ID: "approvers", UserIDs: []string{"user1@bank.com"}},
			},
			ContractAddressWhitelistingRules: []*model.ContractAddressWhitelistingRules{
				{
					Blockchain: "ETH",
					Network:    "mainnet",
					ParallelThresholds: []*model.SequentialThresholds{
						{
							Thresholds: []*model.GroupThreshold{
								{GroupID: "approvers", MinimumSignatures: 1},
							},
						},
					},
				},
			},
		}, nil
	}

	usDecoder := func(b64 string) ([]*model.RuleUserSignature, error) {
		return []*model.RuleUserSignature{
			{UserID: "sa@bank.com", Signature: saSig},
		}, nil
	}

	verifier := NewWhitelistedAssetVerifier([]*ecdsa.PublicKey{saPub}, 1)

	return &assetTestFixture{
		saPriv:    saPriv,
		saPub:     saPub,
		userPriv:  userPriv,
		userPub:   userPub,
		verifier:  verifier,
		asset:     asset,
		rcDecoder: rcDecoder,
		usDecoder: usDecoder,
	}
}

// =============================================================================
// Constructor Tests
// =============================================================================

func TestNewWhitelistedAssetVerifier(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	keys := []*ecdsa.PublicKey{&key.PublicKey}

	t.Run("with keys", func(t *testing.T) {
		v := NewWhitelistedAssetVerifier(keys, 1)
		if v == nil {
			t.Error("NewWhitelistedAssetVerifier() returned nil")
		}
		if len(v.superAdminKeys) != 1 {
			t.Errorf("superAdminKeys length = %d, want 1", len(v.superAdminKeys))
		}
		if v.minValidSignatures != 1 {
			t.Errorf("minValidSignatures = %d, want 1", v.minValidSignatures)
		}
	})

	t.Run("with nil keys", func(t *testing.T) {
		v := NewWhitelistedAssetVerifier(nil, 0)
		if v == nil {
			t.Error("NewWhitelistedAssetVerifier() returned nil")
		}
		if len(v.superAdminKeys) != 0 {
			t.Errorf("superAdminKeys length = %d, want 0", len(v.superAdminKeys))
		}
	})
}

// =============================================================================
// Nil/Missing Input Tests
// =============================================================================

func TestWhitelistedAssetVerifier_NilInput(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAssetVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

	t.Run("nil asset", func(t *testing.T) {
		_, err := v.VerifyWhitelistedAsset(nil, nil, nil)
		if err == nil {
			t.Error("expected error for nil asset")
		}
	})

	t.Run("nil metadata", func(t *testing.T) {
		asset := &model.WhitelistedAsset{}
		_, err := v.VerifyWhitelistedAsset(asset, nil, nil)
		if err == nil {
			t.Error("expected error for nil metadata")
		}
	})
}

// =============================================================================
// Step 1 Tests: Metadata Hash Verification
// =============================================================================

func TestAssetVerifier_Step1_MetadataHash(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAssetVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

	payload := `{"blockchain":"ETH","contractAddress":"0x123"}`
	expectedHash := crypto.CalculateHexHash(payload)

	tests := []struct {
		name        string
		asset       *model.WhitelistedAsset
		expectError bool
		errContains string
	}{
		{
			name: "valid hash",
			asset: &model.WhitelistedAsset{
				Metadata: &model.WhitelistedAssetMetadata{
					PayloadAsString: payload,
					Hash:            expectedHash,
				},
			},
			expectError: false,
		},
		{
			name: "empty payload",
			asset: &model.WhitelistedAsset{
				Metadata: &model.WhitelistedAssetMetadata{
					PayloadAsString: "",
					Hash:            "somehash",
				},
			},
			expectError: true,
			errContains: "payloadAsString is empty",
		},
		{
			name: "empty hash",
			asset: &model.WhitelistedAsset{
				Metadata: &model.WhitelistedAssetMetadata{
					PayloadAsString: payload,
					Hash:            "",
				},
			},
			expectError: true,
			errContains: "metadata hash is empty",
		},
		{
			name: "mismatched hash",
			asset: &model.WhitelistedAsset{
				Metadata: &model.WhitelistedAssetMetadata{
					PayloadAsString: payload,
					Hash:            "deadbeef",
				},
			},
			expectError: true,
			errContains: "metadata hash verification failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.verifyMetadataHash(tt.asset)
			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.errContains != "" {
					intErr, ok := err.(*model.IntegrityError)
					if !ok {
						t.Errorf("expected IntegrityError, got %T", err)
					} else if intErr.Message == "" {
						t.Errorf("expected error containing %q, got empty message", tt.errContains)
					}
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// =============================================================================
// Step 2 Tests: Rules Container Signatures
// =============================================================================

func TestAssetVerifier_Step2_RulesContainerSignatures(t *testing.T) {
	t.Run("no SuperAdmin keys", func(t *testing.T) {
		v := NewWhitelistedAssetVerifier(nil, 0)
		asset := &model.WhitelistedAsset{
			RulesContainer:  "c29tZQ==",
			RulesSignatures: "c29tZQ==",
		}
		err := v.verifyRulesContainerSignatures(asset, func(b64 string) ([]*model.RuleUserSignature, error) {
			return nil, nil
		})
		if err == nil {
			t.Error("expected error when no keys configured")
		}
	})

	t.Run("empty rules container", func(t *testing.T) {
		key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		v := NewWhitelistedAssetVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)
		asset := &model.WhitelistedAsset{
			RulesContainer:  "",
			RulesSignatures: "c29tZQ==",
		}
		err := v.verifyRulesContainerSignatures(asset, nil)
		if err == nil {
			t.Error("expected error for empty rules container")
		}
	})

	t.Run("empty rules signatures", func(t *testing.T) {
		key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		v := NewWhitelistedAssetVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)
		asset := &model.WhitelistedAsset{
			RulesContainer:  "c29tZQ==",
			RulesSignatures: "",
		}
		err := v.verifyRulesContainerSignatures(asset, nil)
		if err == nil {
			t.Error("expected error for empty rules signatures")
		}
	})

	t.Run("decode failure", func(t *testing.T) {
		key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		v := NewWhitelistedAssetVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)
		asset := &model.WhitelistedAsset{
			RulesContainer:  "c29tZQ==",
			RulesSignatures: "c29tZQ==",
		}
		err := v.verifyRulesContainerSignatures(asset, func(b64 string) ([]*model.RuleUserSignature, error) {
			return nil, fmt.Errorf("decode failed")
		})
		if err == nil {
			t.Error("expected error when decode fails")
		}
	})

	t.Run("insufficient signatures", func(t *testing.T) {
		f := buildAssetTestFixture(t)
		v := NewWhitelistedAssetVerifier([]*ecdsa.PublicKey{f.saPub}, 2) // Require 2 but only 1
		err := v.verifyRulesContainerSignatures(f.asset, f.usDecoder)
		if err == nil {
			t.Error("expected error when insufficient signatures")
		}
	})
}

// =============================================================================
// Step 3 Tests: Decode Rules Container
// =============================================================================

func TestAssetVerifier_Step3_DecodeRulesContainer(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAssetVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

	t.Run("nil decoder", func(t *testing.T) {
		asset := &model.WhitelistedAsset{
			RulesContainer: "c29tZQ==",
		}
		_, err := v.decodeRulesContainer(asset, nil)
		if err == nil {
			t.Error("expected error for nil decoder")
		}
	})

	t.Run("decoder failure", func(t *testing.T) {
		asset := &model.WhitelistedAsset{
			RulesContainer: "c29tZQ==",
		}
		_, err := v.decodeRulesContainer(asset, func(b64 string) (*model.DecodedRulesContainer, error) {
			return nil, fmt.Errorf("decode failed")
		})
		if err == nil {
			t.Error("expected error when decoder fails")
		}
	})

	t.Run("successful decode", func(t *testing.T) {
		asset := &model.WhitelistedAsset{
			RulesContainer: "c29tZQ==",
		}
		expected := &model.DecodedRulesContainer{}
		result, err := v.decodeRulesContainer(asset, func(b64 string) (*model.DecodedRulesContainer, error) {
			return expected, nil
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != expected {
			t.Error("returned container does not match expected")
		}
	})
}

// =============================================================================
// Step 4 Tests: Hash Coverage
// =============================================================================

func TestAssetVerifier_Step4_HashCoverage(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAssetVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

	payload := `{"blockchain":"ETH","contractAddress":"0x123"}`
	correctHash := crypto.CalculateHexHash(payload)

	t.Run("nil signedContractAddress", func(t *testing.T) {
		asset := &model.WhitelistedAsset{
			SignedContractAddress: nil,
			Metadata: &model.WhitelistedAssetMetadata{
				Hash:            correctHash,
				PayloadAsString: payload,
			},
		}
		_, err := v.verifyHashInSignedHashes(asset)
		if err == nil {
			t.Error("expected error for nil signedContractAddress")
		}
	})

	t.Run("empty signatures", func(t *testing.T) {
		asset := &model.WhitelistedAsset{
			SignedContractAddress: &model.SignedContractAddress{Signatures: nil},
			Metadata: &model.WhitelistedAssetMetadata{
				Hash:            correctHash,
				PayloadAsString: payload,
			},
		}
		_, err := v.verifyHashInSignedHashes(asset)
		if err == nil {
			t.Error("expected error for empty signatures")
		}
	})

	t.Run("hash found in signatures", func(t *testing.T) {
		asset := &model.WhitelistedAsset{
			SignedContractAddress: &model.SignedContractAddress{
				Signatures: []model.WhitelistSignature{
					{Hashes: []string{correctHash}},
				},
			},
			Metadata: &model.WhitelistedAssetMetadata{
				Hash:            correctHash,
				PayloadAsString: payload,
			},
		}
		foundHash, err := v.verifyHashInSignedHashes(asset)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if foundHash != correctHash {
			t.Errorf("foundHash = %q, want %q", foundHash, correctHash)
		}
	})

	t.Run("hash not found", func(t *testing.T) {
		asset := &model.WhitelistedAsset{
			SignedContractAddress: &model.SignedContractAddress{
				Signatures: []model.WhitelistSignature{
					{Hashes: []string{"wronghash"}},
				},
			},
			Metadata: &model.WhitelistedAssetMetadata{
				Hash:            correctHash,
				PayloadAsString: payload,
			},
		}
		_, err := v.verifyHashInSignedHashes(asset)
		if err == nil {
			t.Error("expected error when hash not covered")
		}
	})
}

// =============================================================================
// Step 4 Legacy Hash Tests
// =============================================================================

func TestAssetVerifier_Step4_LegacyHashFallback(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAssetVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

	// Build payload with isNFT field (legacy hashes are computed without it)
	payload := `{"blockchain":"ETH","network":"mainnet","contractAddress":"0xUSDC","name":"USDC","symbol":"USDC","decimals":6,"isNFT":false}`
	currentHash := crypto.CalculateHexHash(payload)

	// Compute legacy hashes
	legacyHashes := ComputeAssetLegacyHashes(payload)
	if len(legacyHashes) == 0 {
		t.Skip("No legacy hashes computed for this payload - skipping")
	}

	t.Run("legacy hash match succeeds", func(t *testing.T) {
		asset := &model.WhitelistedAsset{
			SignedContractAddress: &model.SignedContractAddress{
				Signatures: []model.WhitelistSignature{
					{Hashes: []string{legacyHashes[0]}},
				},
			},
			Metadata: &model.WhitelistedAssetMetadata{
				Hash:            currentHash,
				PayloadAsString: payload,
			},
		}
		foundHash, err := v.verifyHashInSignedHashes(asset)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if foundHash != legacyHashes[0] {
			t.Errorf("foundHash = %q, want legacy hash %q", foundHash, legacyHashes[0])
		}
		if foundHash == currentHash {
			t.Error("foundHash should not match current hash when using legacy fallback")
		}
	})
}

// =============================================================================
// Step 5 Tests: Whitelist Signatures
// =============================================================================

func TestAssetVerifier_Step5_WhitelistSignatures(t *testing.T) {
	t.Run("no rules for blockchain", func(t *testing.T) {
		f := buildAssetTestFixture(t)
		// Use a container with no matching rules
		container := &model.DecodedRulesContainer{
			ContractAddressWhitelistingRules: []*model.ContractAddressWhitelistingRules{
				{
					Blockchain:         "BTC",
					Network:            "mainnet",
					ParallelThresholds: []*model.SequentialThresholds{},
				},
			},
		}
		err := f.verifier.verifyWhitelistSignatures(f.asset, container, f.asset.Metadata.Hash)
		if err == nil {
			t.Error("expected error when no matching rules")
		}
		if _, ok := err.(*model.WhitelistError); !ok {
			t.Errorf("expected WhitelistError, got %T", err)
		}
	})

	t.Run("empty parallel thresholds", func(t *testing.T) {
		f := buildAssetTestFixture(t)
		container := &model.DecodedRulesContainer{
			ContractAddressWhitelistingRules: []*model.ContractAddressWhitelistingRules{
				{
					Blockchain:         "ETH",
					Network:            "mainnet",
					ParallelThresholds: []*model.SequentialThresholds{},
				},
			},
		}
		err := f.verifier.verifyWhitelistSignatures(f.asset, container, f.asset.Metadata.Hash)
		if err == nil {
			t.Error("expected error when no thresholds defined")
		}
	})

	t.Run("group not found in container", func(t *testing.T) {
		f := buildAssetTestFixture(t)
		container := &model.DecodedRulesContainer{
			Groups: []*model.RuleGroup{},
			ContractAddressWhitelistingRules: []*model.ContractAddressWhitelistingRules{
				{
					Blockchain: "ETH",
					Network:    "mainnet",
					ParallelThresholds: []*model.SequentialThresholds{
						{
							Thresholds: []*model.GroupThreshold{
								{GroupID: "ghost_group", MinimumSignatures: 1},
							},
						},
					},
				},
			},
		}
		err := f.verifier.verifyWhitelistSignatures(f.asset, container, f.asset.Metadata.Hash)
		if err == nil {
			t.Error("expected error when group not found")
		}
	})

	t.Run("threshold requires 2 sigs but only 1 valid", func(t *testing.T) {
		f := buildAssetTestFixture(t)
		container := &model.DecodedRulesContainer{
			Users: []*model.RuleUser{
				{ID: "user1@bank.com", PublicKey: f.userPub, Roles: []string{"USER"}},
			},
			Groups: []*model.RuleGroup{
				{ID: "approvers", UserIDs: []string{"user1@bank.com"}},
			},
			ContractAddressWhitelistingRules: []*model.ContractAddressWhitelistingRules{
				{
					Blockchain: "ETH",
					Network:    "mainnet",
					ParallelThresholds: []*model.SequentialThresholds{
						{
							Thresholds: []*model.GroupThreshold{
								{GroupID: "approvers", MinimumSignatures: 2},
							},
						},
					},
				},
			},
		}
		err := f.verifier.verifyWhitelistSignatures(f.asset, container, f.asset.Metadata.Hash)
		if err == nil {
			t.Error("expected error when threshold not met")
		}
	})
}

// =============================================================================
// End-to-End Happy Path
// =============================================================================

func TestAssetVerifier_EndToEnd_HappyPath(t *testing.T) {
	f := buildAssetTestFixture(t)

	result, err := f.verifier.VerifyWhitelistedAsset(f.asset, f.rcDecoder, f.usDecoder)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("result should not be nil")
	}
	if result.RulesContainer == nil {
		t.Error("rules container should not be nil")
	}
	if result.VerifiedHash != f.asset.Metadata.Hash {
		t.Errorf("verifiedHash = %q, want %q", result.VerifiedHash, f.asset.Metadata.Hash)
	}
}

func TestAssetVerifier_EndToEnd_TwoUsersSequentialThresholds(t *testing.T) {
	saPriv, saPub := generateAssetTestKeyPair(t)
	u1Priv, u1Pub := generateAssetTestKeyPair(t)
	u2Priv, u2Pub := generateAssetTestKeyPair(t)

	payload := `{"blockchain":"ETH","network":"mainnet","contractAddress":"0xUSDC","name":"USDC"}`
	metadataHash := crypto.CalculateHexHash(payload)

	hashes := []string{metadataHash}
	hashesJSON, _ := json.Marshal(hashes)
	u1Sig, _ := crypto.SignData(u1Priv, hashesJSON)
	u2Sig, _ := crypto.SignData(u2Priv, hashesJSON)

	rulesB64 := encodeRulesContainerJSON(t, map[string]interface{}{
		"users":  []map[string]interface{}{},
		"groups": []map[string]interface{}{},
	})
	rulesData, _ := base64.StdEncoding.DecodeString(rulesB64)
	saSig, _ := crypto.SignData(saPriv, rulesData)

	asset := &model.WhitelistedAsset{
		ID:         "asset-2",
		Blockchain: "ETH",
		Network:    "mainnet",
		Metadata: &model.WhitelistedAssetMetadata{
			Hash:            metadataHash,
			PayloadAsString: payload,
		},
		RulesContainer:  rulesB64,
		RulesSignatures: base64.StdEncoding.EncodeToString([]byte("dummy")),
		SignedContractAddress: &model.SignedContractAddress{
			Signatures: []model.WhitelistSignature{
				{
					UserSignature: &model.WhitelistUserSignature{UserID: "user1@bank.com", Signature: u1Sig},
					Hashes:        hashes,
				},
				{
					UserSignature: &model.WhitelistUserSignature{UserID: "user2@bank.com", Signature: u2Sig},
					Hashes:        hashes,
				},
			},
		},
	}

	rcDecoder := func(b64 string) (*model.DecodedRulesContainer, error) {
		return &model.DecodedRulesContainer{
			Users: []*model.RuleUser{
				{ID: "user1@bank.com", PublicKey: u1Pub, Roles: []string{"USER"}},
				{ID: "user2@bank.com", PublicKey: u2Pub, Roles: []string{"USER"}},
			},
			Groups: []*model.RuleGroup{
				{ID: "group_a", UserIDs: []string{"user1@bank.com"}},
				{ID: "group_b", UserIDs: []string{"user2@bank.com"}},
			},
			ContractAddressWhitelistingRules: []*model.ContractAddressWhitelistingRules{
				{
					Blockchain: "ETH",
					Network:    "mainnet",
					ParallelThresholds: []*model.SequentialThresholds{
						{
							Thresholds: []*model.GroupThreshold{
								{GroupID: "group_a", MinimumSignatures: 1},
								{GroupID: "group_b", MinimumSignatures: 1},
							},
						},
					},
				},
			},
		}, nil
	}

	usDecoder := func(b64 string) ([]*model.RuleUserSignature, error) {
		return []*model.RuleUserSignature{
			{UserID: "sa", Signature: saSig},
		}, nil
	}

	verifier := NewWhitelistedAssetVerifier([]*ecdsa.PublicKey{saPub}, 1)
	result, err := verifier.VerifyWhitelistedAsset(asset, rcDecoder, usDecoder)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("result should not be nil")
	}
	if result.VerifiedHash != metadataHash {
		t.Errorf("verifiedHash = %q, want %q", result.VerifiedHash, metadataHash)
	}
}

func TestAssetVerifier_EndToEnd_ParallelPathsORLogic(t *testing.T) {
	f := buildAssetTestFixture(t)

	// Override rcDecoder with two parallel paths - first fails, second succeeds
	rcDecoder := func(b64 string) (*model.DecodedRulesContainer, error) {
		return &model.DecodedRulesContainer{
			Users: []*model.RuleUser{
				{ID: "user1@bank.com", PublicKey: f.userPub, Roles: []string{"USER"}},
			},
			Groups: []*model.RuleGroup{
				{ID: "approvers", UserIDs: []string{"user1@bank.com"}},
				{ID: "other_team", UserIDs: []string{"nobody@bank.com"}},
			},
			ContractAddressWhitelistingRules: []*model.ContractAddressWhitelistingRules{
				{
					Blockchain: "ETH",
					Network:    "mainnet",
					ParallelThresholds: []*model.SequentialThresholds{
						// Path 1: will fail (user not in group)
						{Thresholds: []*model.GroupThreshold{{GroupID: "other_team", MinimumSignatures: 1}}},
						// Path 2: will pass
						{Thresholds: []*model.GroupThreshold{{GroupID: "approvers", MinimumSignatures: 1}}},
					},
				},
			},
		}, nil
	}

	result, err := f.verifier.VerifyWhitelistedAsset(f.asset, rcDecoder, f.usDecoder)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("result should not be nil")
	}
}

// =============================================================================
// Sequential Threshold Edge Cases
// =============================================================================

func TestAssetVerifier_SequentialThresholds_NilInput(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAssetVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

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

// =============================================================================
// Group Threshold Edge Cases
// =============================================================================

func TestAssetVerifier_GroupThreshold_EmptyGroupZeroMinSigs(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAssetVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

	container := &model.DecodedRulesContainer{
		Groups: []*model.RuleGroup{
			{ID: "empty_group", UserIDs: []string{}},
		},
	}

	t.Run("empty group with 0 min sigs is OK", func(t *testing.T) {
		gt := &model.GroupThreshold{GroupID: "empty_group", MinimumSignatures: 0}
		err := v.verifyGroupThreshold(gt, container, nil, "hash", nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("empty group with min sigs > 0 fails", func(t *testing.T) {
		gt := &model.GroupThreshold{GroupID: "empty_group", MinimumSignatures: 1}
		err := v.verifyGroupThreshold(gt, container, nil, "hash", nil)
		if err == nil {
			t.Error("expected error for empty group with min sigs > 0")
		}
	})
}

func TestAssetVerifier_GroupThreshold_NilUserSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAssetVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

	container := &model.DecodedRulesContainer{
		Groups: []*model.RuleGroup{
			{ID: "group1", UserIDs: []string{"user1@bank.com"}},
		},
	}

	signatures := []model.WhitelistSignature{
		{UserSignature: nil, Hashes: []string{"hash1"}},
	}

	gt := &model.GroupThreshold{GroupID: "group1", MinimumSignatures: 1}
	err := v.verifyGroupThreshold(gt, container, signatures, "hash1", precomputeHashesJSON(signatures))
	if err == nil {
		t.Error("expected error when all signatures have nil userSig")
	}
}

// =============================================================================
// Named Tests Matching Requested Convention
// =============================================================================

// TestWhitelistedAssetVerifier_ValidEnvelope verifies that a fully valid asset
// envelope (correct hash, valid SuperAdmin sig, valid user sig meeting threshold)
// passes verification end-to-end.
func TestWhitelistedAssetVerifier_ValidEnvelope(t *testing.T) {
	f := buildAssetTestFixture(t)

	result, err := f.verifier.VerifyWhitelistedAsset(f.asset, f.rcDecoder, f.usDecoder)
	if err != nil {
		t.Fatalf("expected no error for valid envelope, got: %v", err)
	}
	if result == nil {
		t.Fatal("result should not be nil for valid envelope")
	}
	if result.RulesContainer == nil {
		t.Error("rules container should not be nil in result")
	}
	if result.VerifiedHash != f.asset.Metadata.Hash {
		t.Errorf("verifiedHash = %q, want %q", result.VerifiedHash, f.asset.Metadata.Hash)
	}
}

// TestWhitelistedAssetVerifier_InvalidHash verifies that a hash mismatch between
// the computed hash (from PayloadAsString) and the provided hash produces an
// IntegrityError.
func TestWhitelistedAssetVerifier_InvalidHash(t *testing.T) {
	f := buildAssetTestFixture(t)

	// Corrupt the metadata hash so it no longer matches the payload
	f.asset.Metadata.Hash = "0000000000000000000000000000000000000000000000000000000000000000"

	_, err := f.verifier.VerifyWhitelistedAsset(f.asset, f.rcDecoder, f.usDecoder)
	if err == nil {
		t.Fatal("expected error for invalid hash, got nil")
	}

	intErr, ok := err.(*model.IntegrityError)
	if !ok {
		t.Fatalf("expected IntegrityError, got %T: %v", err, err)
	}
	if intErr.Message == "" {
		t.Error("IntegrityError message should not be empty")
	}
}

// TestWhitelistedAssetVerifier_InsufficientSignatures verifies that when the
// number of valid user signatures is below the governance threshold, verification
// fails.
func TestWhitelistedAssetVerifier_InsufficientSignatures(t *testing.T) {
	f := buildAssetTestFixture(t)

	// Override rcDecoder to require 2 signatures but only 1 user exists
	rcDecoder := func(b64 string) (*model.DecodedRulesContainer, error) {
		return &model.DecodedRulesContainer{
			Users: []*model.RuleUser{
				{ID: "user1@bank.com", PublicKey: f.userPub, Roles: []string{"USER"}},
			},
			Groups: []*model.RuleGroup{
				{ID: "approvers", UserIDs: []string{"user1@bank.com"}},
			},
			ContractAddressWhitelistingRules: []*model.ContractAddressWhitelistingRules{
				{
					Blockchain: "ETH",
					Network:    "mainnet",
					ParallelThresholds: []*model.SequentialThresholds{
						{
							Thresholds: []*model.GroupThreshold{
								{GroupID: "approvers", MinimumSignatures: 2},
							},
						},
					},
				},
			},
		}, nil
	}

	_, err := f.verifier.VerifyWhitelistedAsset(f.asset, rcDecoder, f.usDecoder)
	if err == nil {
		t.Fatal("expected error when insufficient signatures, got nil")
	}
}

// TestWhitelistedAssetVerifier_NilMetadata verifies that a nil Metadata field
// on the asset produces an error immediately.
func TestWhitelistedAssetVerifier_NilMetadata(t *testing.T) {
	f := buildAssetTestFixture(t)
	f.asset.Metadata = nil

	_, err := f.verifier.VerifyWhitelistedAsset(f.asset, f.rcDecoder, f.usDecoder)
	if err == nil {
		t.Fatal("expected error for nil metadata, got nil")
	}
}

// TestWhitelistedAssetVerifier_EmptySignatures verifies that when
// SignedContractAddress.Signatures is empty, verification fails at Step 4.
func TestWhitelistedAssetVerifier_EmptySignatures(t *testing.T) {
	f := buildAssetTestFixture(t)
	f.asset.SignedContractAddress = &model.SignedContractAddress{
		Signatures: []model.WhitelistSignature{},
	}

	_, err := f.verifier.VerifyWhitelistedAsset(f.asset, f.rcDecoder, f.usDecoder)
	if err == nil {
		t.Fatal("expected error for empty signatures, got nil")
	}

	intErr, ok := err.(*model.IntegrityError)
	if !ok {
		t.Fatalf("expected IntegrityError, got %T: %v", err, err)
	}
	if intErr.Message == "" {
		t.Error("IntegrityError message should not be empty")
	}
}

// TestWhitelistedAssetVerifier_LegacyHashFallback verifies that when the
// current hash is not found in signatures but a legacy hash (computed without
// isNFT) IS found, verification succeeds with the legacy hash as VerifiedHash.
func TestWhitelistedAssetVerifier_LegacyHashFallback(t *testing.T) {
	saPriv, saPub := generateAssetTestKeyPair(t)
	userPriv, userPub := generateAssetTestKeyPair(t)

	// Payload with isNFT - legacy hash will remove it
	payload := `{"blockchain":"ETH","network":"mainnet","contractAddress":"0xUSDC","name":"USDC","symbol":"USDC","decimals":6,"isNFT":false}`
	currentHash := crypto.CalculateHexHash(payload)

	legacyHashes := ComputeAssetLegacyHashes(payload)
	if len(legacyHashes) == 0 {
		t.Skip("No legacy hashes computed for this payload - skipping")
	}
	legacyHash := legacyHashes[0]

	// Sign with legacy hash (simulating an asset signed before isNFT was added)
	hashes := []string{legacyHash}
	hashesJSON, _ := json.Marshal(hashes)
	userSig, err := crypto.SignData(userPriv, hashesJSON)
	if err != nil {
		t.Fatalf("failed to sign hashes: %v", err)
	}

	rulesB64 := encodeRulesContainerJSON(t, map[string]interface{}{
		"users":  []map[string]interface{}{},
		"groups": []map[string]interface{}{},
	})
	rulesData, _ := base64.StdEncoding.DecodeString(rulesB64)
	saSig, _ := crypto.SignData(saPriv, rulesData)

	asset := &model.WhitelistedAsset{
		ID:         "legacy-asset",
		Blockchain: "ETH",
		Network:    "mainnet",
		Metadata: &model.WhitelistedAssetMetadata{
			Hash:            currentHash,
			PayloadAsString: payload,
		},
		RulesContainer:  rulesB64,
		RulesSignatures: base64.StdEncoding.EncodeToString([]byte("dummy")),
		SignedContractAddress: &model.SignedContractAddress{
			Signatures: []model.WhitelistSignature{
				{
					UserSignature: &model.WhitelistUserSignature{
						UserID:    "user1@bank.com",
						Signature: userSig,
					},
					Hashes: hashes,
				},
			},
		},
	}

	rcDecoder := func(b64 string) (*model.DecodedRulesContainer, error) {
		return &model.DecodedRulesContainer{
			Users: []*model.RuleUser{
				{ID: "user1@bank.com", PublicKey: userPub, Roles: []string{"USER"}},
			},
			Groups: []*model.RuleGroup{
				{ID: "approvers", UserIDs: []string{"user1@bank.com"}},
			},
			ContractAddressWhitelistingRules: []*model.ContractAddressWhitelistingRules{
				{
					Blockchain: "ETH",
					Network:    "mainnet",
					ParallelThresholds: []*model.SequentialThresholds{
						{
							Thresholds: []*model.GroupThreshold{
								{GroupID: "approvers", MinimumSignatures: 1},
							},
						},
					},
				},
			},
		}, nil
	}

	usDecoder := func(b64 string) ([]*model.RuleUserSignature, error) {
		return []*model.RuleUserSignature{
			{UserID: "sa@bank.com", Signature: saSig},
		}, nil
	}

	verifier := NewWhitelistedAssetVerifier([]*ecdsa.PublicKey{saPub}, 1)
	result, err := verifier.VerifyWhitelistedAsset(asset, rcDecoder, usDecoder)
	if err != nil {
		t.Fatalf("expected no error for legacy hash fallback, got: %v", err)
	}
	if result == nil {
		t.Fatal("result should not be nil")
	}
	if result.VerifiedHash == currentHash {
		t.Error("verifiedHash should not match current hash when using legacy fallback")
	}
	if result.VerifiedHash != legacyHash {
		t.Errorf("verifiedHash = %q, want legacy hash %q", result.VerifiedHash, legacyHash)
	}
}

// TestWhitelistedAssetVerifier_NilRulesContainer verifies that when the rules
// container decoder returns nil, verification fails with an appropriate error.
func TestWhitelistedAssetVerifier_NilRulesContainer(t *testing.T) {
	f := buildAssetTestFixture(t)

	// Override rcDecoder to return nil container
	nilRCDecoder := func(b64 string) (*model.DecodedRulesContainer, error) {
		return nil, fmt.Errorf("rules container not available")
	}

	_, err := f.verifier.VerifyWhitelistedAsset(f.asset, nilRCDecoder, f.usDecoder)
	if err == nil {
		t.Fatal("expected error when rules container decoder returns nil, got nil")
	}
}

// TestWhitelistedAssetVerifier_ValidEnvelopeList verifies batch verification
// of multiple valid assets. Each asset is independently verified by calling
// VerifyWhitelistedAsset in sequence.
func TestWhitelistedAssetVerifier_ValidEnvelopeList(t *testing.T) {
	saPriv, saPub := generateAssetTestKeyPair(t)
	userPriv, userPub := generateAssetTestKeyPair(t)

	payloads := []struct {
		id         string
		blockchain string
		network    string
		payload    string
	}{
		{"asset-A", "ETH", "mainnet", `{"blockchain":"ETH","network":"mainnet","contractAddress":"0xUSDC","name":"USDC"}`},
		{"asset-B", "ETH", "mainnet", `{"blockchain":"ETH","network":"mainnet","contractAddress":"0xDAI","name":"DAI"}`},
		{"asset-C", "ETH", "mainnet", `{"blockchain":"ETH","network":"mainnet","contractAddress":"0xWBTC","name":"WBTC"}`},
	}

	rulesB64 := encodeRulesContainerJSON(t, map[string]interface{}{
		"users":  []map[string]interface{}{},
		"groups": []map[string]interface{}{},
	})
	rulesData, _ := base64.StdEncoding.DecodeString(rulesB64)
	saSig, _ := crypto.SignData(saPriv, rulesData)

	var assets []*model.WhitelistedAsset
	for _, p := range payloads {
		metadataHash := crypto.CalculateHexHash(p.payload)
		hashes := []string{metadataHash}
		hashesJSON, _ := json.Marshal(hashes)
		userSig, err := crypto.SignData(userPriv, hashesJSON)
		if err != nil {
			t.Fatalf("failed to sign hashes for %s: %v", p.id, err)
		}

		assets = append(assets, &model.WhitelistedAsset{
			ID:         p.id,
			Blockchain: p.blockchain,
			Network:    p.network,
			Metadata: &model.WhitelistedAssetMetadata{
				Hash:            metadataHash,
				PayloadAsString: p.payload,
			},
			RulesContainer:  rulesB64,
			RulesSignatures: base64.StdEncoding.EncodeToString([]byte("dummy")),
			SignedContractAddress: &model.SignedContractAddress{
				Signatures: []model.WhitelistSignature{
					{
						UserSignature: &model.WhitelistUserSignature{
							UserID:    "user1@bank.com",
							Signature: userSig,
						},
						Hashes: hashes,
					},
				},
			},
		})
	}

	rcDecoder := func(b64 string) (*model.DecodedRulesContainer, error) {
		return &model.DecodedRulesContainer{
			Users: []*model.RuleUser{
				{ID: "user1@bank.com", PublicKey: userPub, Roles: []string{"USER"}},
			},
			Groups: []*model.RuleGroup{
				{ID: "approvers", UserIDs: []string{"user1@bank.com"}},
			},
			ContractAddressWhitelistingRules: []*model.ContractAddressWhitelistingRules{
				{
					Blockchain: "ETH",
					Network:    "mainnet",
					ParallelThresholds: []*model.SequentialThresholds{
						{
							Thresholds: []*model.GroupThreshold{
								{GroupID: "approvers", MinimumSignatures: 1},
							},
						},
					},
				},
			},
		}, nil
	}

	usDecoder := func(b64 string) ([]*model.RuleUserSignature, error) {
		return []*model.RuleUserSignature{
			{UserID: "sa@bank.com", Signature: saSig},
		}, nil
	}

	verifier := NewWhitelistedAssetVerifier([]*ecdsa.PublicKey{saPub}, 1)

	for i, asset := range assets {
		t.Run(fmt.Sprintf("asset_%d_%s", i, asset.ID), func(t *testing.T) {
			result, err := verifier.VerifyWhitelistedAsset(asset, rcDecoder, usDecoder)
			if err != nil {
				t.Fatalf("unexpected error for asset %s: %v", asset.ID, err)
			}
			if result == nil {
				t.Fatalf("result should not be nil for asset %s", asset.ID)
			}
			if result.RulesContainer == nil {
				t.Errorf("rules container should not be nil for asset %s", asset.ID)
			}
			if result.VerifiedHash != asset.Metadata.Hash {
				t.Errorf("asset %s: verifiedHash = %q, want %q",
					asset.ID, result.VerifiedHash, asset.Metadata.Hash)
			}
		})
	}
}
