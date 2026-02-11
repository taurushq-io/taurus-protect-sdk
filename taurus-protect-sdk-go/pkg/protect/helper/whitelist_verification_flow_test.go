// Package helper provides signature verification and validation utilities.
//
// This file contains tests for the 6-step whitelisted address verification flow:
// Step 1: Verify metadata hash (SHA256(payloadAsString) == metadata.hash)
// Step 2: Verify rules container signatures (SuperAdmin keys)
// Step 3: Decode rules container (base64 → protobuf → model)
// Step 4: Verify hash coverage (metadata.hash in signature hashes list)
// Step 5: Verify whitelist signatures meet governance thresholds
// Step 6: Parse WhitelistedAddress from verified payload
package helper

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/testdata"
)

// =============================================================================
// Step 1 Tests: Verify Metadata Hash
// =============================================================================

// TestStep1_VerifyMetadataHashSuccess verifies that step 1 passes when
// SHA256(payloadAsString) equals metadata.hash.
func TestStep1_VerifyMetadataHashSuccess(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAddressVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

	// Use real fixture data
	payload := testdata.RealPayloadAsString
	expectedHash := testdata.RealMetadataHash

	// Verify computed hash matches expected
	computedHash := crypto.CalculateHexHash(payload)
	if computedHash != expectedHash {
		t.Fatalf("Test setup error: computed hash %q != expected %q", computedHash, expectedHash)
	}

	addr := &model.WhitelistedAddress{
		Metadata: &model.WhitelistedAssetMetadata{
			PayloadAsString: payload,
			Hash:            expectedHash,
		},
	}

	err := v.verifyMetadataHash(addr)
	if err != nil {
		t.Errorf("Step 1 should pass: verifyMetadataHash() unexpected error: %v", err)
	}
}

// TestStep1_VerifyMetadataHashFailure verifies that step 1 fails with IntegrityError
// when the computed hash does not match the provided hash.
func TestStep1_VerifyMetadataHashFailure(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAddressVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

	tests := []struct {
		name    string
		payload string
		hash    string
	}{
		{
			name:    "hash mismatch - wrong hash",
			payload: testdata.RealPayloadAsString,
			hash:    "0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			name:    "hash mismatch - empty hash",
			payload: testdata.RealPayloadAsString,
			hash:    "",
		},
		{
			name:    "hash mismatch - empty payload",
			payload: "",
			hash:    testdata.RealMetadataHash,
		},
		{
			name:    "hash mismatch - modified payload",
			payload: testdata.RealPayloadAsString + " ",
			hash:    testdata.RealMetadataHash,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr := &model.WhitelistedAddress{
				Metadata: &model.WhitelistedAssetMetadata{
					PayloadAsString: tt.payload,
					Hash:            tt.hash,
				},
			}

			err := v.verifyMetadataHash(addr)
			if err == nil {
				t.Error("Step 1 should fail: verifyMetadataHash() expected IntegrityError, got nil")
				return
			}

			// Verify it's an IntegrityError
			if _, ok := err.(*model.IntegrityError); !ok {
				t.Errorf("Step 1 should return IntegrityError, got %T", err)
			}
		})
	}
}

// =============================================================================
// Step 2 Tests: Verify Rules Container Signatures
// =============================================================================

// TestStep2_VerifyRulesSignaturesSuccess verifies that step 2 passes when
// the rules container is properly signed by SuperAdmin keys.
func TestStep2_VerifyRulesSignaturesSuccess(t *testing.T) {
	// Generate test key pairs for SuperAdmins
	key1, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	key2, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	superAdminKeys := []*ecdsa.PublicKey{&key1.PublicKey, &key2.PublicKey}

	v := NewWhitelistedAddressVerifier(superAdminKeys, 2)

	// Create test rules container data
	rulesContainerData := []byte("test rules container data")
	rulesContainer := EncodeBase64(rulesContainerData)

	// Sign the rules container with both keys
	sig1, _ := crypto.SignData(key1, rulesContainerData)
	sig2, _ := crypto.SignData(key2, rulesContainerData)

	// Create base64-encoded rules signatures (any value since we mock the decoder)
	rulesSignatures := EncodeBase64([]byte("mock signatures"))

	addr := &model.WhitelistedAddress{
		RulesContainer:  rulesContainer,
		RulesSignatures: rulesSignatures,
	}

	// Create a mock decoder that returns the valid signatures
	mockDecoder := func(base64Data string) ([]*model.RuleUserSignature, error) {
		return []*model.RuleUserSignature{
			{UserID: "superadmin1@bank.com", Signature: sig1},
			{UserID: "superadmin2@bank.com", Signature: sig2},
		}, nil
	}

	err := v.verifyRulesContainerSignatures(addr, mockDecoder)
	if err != nil {
		t.Errorf("Step 2 should pass: verifyRulesContainerSignatures() unexpected error: %v", err)
	}
}

// TestStep2_VerifyRulesSignaturesFailure verifies that step 2 fails with IntegrityError
// when the rules container signatures are invalid.
func TestStep2_VerifyRulesSignaturesFailure(t *testing.T) {
	// Generate random keys that won't match the fixture signatures
	wrongKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAddressVerifier([]*ecdsa.PublicKey{&wrongKey.PublicKey}, 1)

	// Load rules container and signatures from fixture
	rulesContainer, rulesSignatures := loadRulesContainerFromFixture(t)

	addr := &model.WhitelistedAddress{
		RulesContainer:  rulesContainer,
		RulesSignatures: rulesSignatures,
	}

	mockDecoder := func(base64Data string) ([]*model.RuleUserSignature, error) {
		return decodeUserSignaturesFromBase64(base64Data)
	}

	err := v.verifyRulesContainerSignatures(addr, mockDecoder)
	if err == nil {
		t.Error("Step 2 should fail: verifyRulesContainerSignatures() expected IntegrityError, got nil")
		return
	}

	// Verify it's an IntegrityError
	if _, ok := err.(*model.IntegrityError); !ok {
		t.Errorf("Step 2 should return IntegrityError, got %T", err)
	}
}

// TestStep2_VerifyRulesSignaturesFailure_NoKeys verifies that step 2 fails
// when no SuperAdmin keys are configured.
func TestStep2_VerifyRulesSignaturesFailure_NoKeys(t *testing.T) {
	v := NewWhitelistedAddressVerifier(nil, 0)

	addr := &model.WhitelistedAddress{
		RulesContainer:  "somebase64data",
		RulesSignatures: "somebase64sigs",
	}

	err := v.verifyRulesContainerSignatures(addr, nil)
	if err == nil {
		t.Error("Step 2 should fail: expected error when no keys configured")
		return
	}

	if _, ok := err.(*model.IntegrityError); !ok {
		t.Errorf("Step 2 should return IntegrityError, got %T", err)
	}
}

// =============================================================================
// Step 3 Tests: Decode Rules Container
// =============================================================================

// TestStep3_DecodeRulesContainerSuccess verifies that step 3 passes when
// the rules container can be decoded from base64 protobuf to model.
func TestStep3_DecodeRulesContainerSuccess(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAddressVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

	// Load rules container from fixture
	rulesContainer, _ := loadRulesContainerFromFixture(t)

	addr := &model.WhitelistedAddress{
		RulesContainer: rulesContainer,
	}

	// Create a mock decoder that returns a valid container
	mockDecoder := func(base64Data string) (*model.DecodedRulesContainer, error) {
		return &model.DecodedRulesContainer{
			Users: []*model.RuleUser{
				{ID: "team1@bank.com", Roles: []string{"USER", "OPERATOR"}},
			},
			Groups: []*model.RuleGroup{
				{ID: "team1", UserIDs: []string{"team1@bank.com"}},
			},
			AddressWhitelistingRules: []*model.AddressWhitelistingRules{
				{
					Currency: "ALGO",
					Network:  "mainnet",
					ParallelThresholds: []*model.SequentialThresholds{
						{Thresholds: []*model.GroupThreshold{
							{GroupID: "team1", MinimumSignatures: 1},
						}},
					},
				},
			},
		}, nil
	}

	container, err := v.decodeRulesContainer(addr, mockDecoder)
	if err != nil {
		t.Errorf("Step 3 should pass: decodeRulesContainer() unexpected error: %v", err)
		return
	}

	if container == nil {
		t.Error("Step 3 should return non-nil container")
		return
	}

	// Verify the container has expected content
	if len(container.Users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(container.Users))
	}
	if len(container.Groups) != 1 {
		t.Errorf("Expected 1 group, got %d", len(container.Groups))
	}
}

// TestStep3_DecodeRulesContainerFailure verifies that step 3 fails with IntegrityError
// when the rules container cannot be decoded.
func TestStep3_DecodeRulesContainerFailure(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAddressVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

	addr := &model.WhitelistedAddress{
		RulesContainer: "invalid-base64-data",
	}

	// Create a mock decoder that fails
	mockDecoder := func(base64Data string) (*model.DecodedRulesContainer, error) {
		return nil, fmt.Errorf("failed to decode protobuf: invalid data")
	}

	_, err := v.decodeRulesContainer(addr, mockDecoder)
	if err == nil {
		t.Error("Step 3 should fail: decodeRulesContainer() expected IntegrityError, got nil")
		return
	}

	// Verify it's an IntegrityError
	if _, ok := err.(*model.IntegrityError); !ok {
		t.Errorf("Step 3 should return IntegrityError, got %T", err)
	}
}

// TestStep3_DecodeRulesContainerFailure_NilDecoder verifies that step 3 fails
// when no decoder is provided.
func TestStep3_DecodeRulesContainerFailure_NilDecoder(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAddressVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

	addr := &model.WhitelistedAddress{
		RulesContainer: "somedata",
	}

	_, err := v.decodeRulesContainer(addr, nil)
	if err == nil {
		t.Error("Step 3 should fail: expected error when decoder is nil")
	}
}

// =============================================================================
// Step 4 Tests: Verify Hash Coverage
// =============================================================================

// TestStep4_VerifyHashCoverageSuccess verifies that step 4 passes when
// the metadata hash is found in at least one signature's hashes list.
func TestStep4_VerifyHashCoverageSuccess(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAddressVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

	metadataHash := testdata.RealMetadataHash

	addr := &model.WhitelistedAddress{
		SignedAddress: &model.SignedWhitelistedAddress{
			Signatures: []model.WhitelistSignature{
				{
					Hashes: []string{metadataHash},
					UserSignature: &model.WhitelistUserSignature{
						UserID:    "team1@bank.com",
						Signature: "somesignature",
					},
				},
			},
		},
		Metadata: &model.WhitelistedAssetMetadata{
			Hash:            metadataHash,
			PayloadAsString: testdata.RealPayloadAsString,
		},
	}

	foundHash, err := v.verifyHashInSignedHashes(addr)
	if err != nil {
		t.Errorf("Step 4 should pass: verifyHashInSignedHashes() unexpected error: %v", err)
		return
	}

	if foundHash != metadataHash {
		t.Errorf("Step 4 should return matching hash: got %q, want %q", foundHash, metadataHash)
	}
}

// TestStep4_VerifyHashCoverageFailure_UsesLegacy verifies that step 4 falls back
// to legacy hash computation when the current hash is not found.
func TestStep4_VerifyHashCoverageFailure_UsesLegacy(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAddressVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

	// Use legacy test case: Case1 where contractType was added after signing
	legacyHash := testdata.Case1LegacyHash
	currentHash := testdata.Case1CurrentHash
	currentPayload := testdata.Case1CurrentPayload

	// Verify test setup: current and legacy hashes should differ
	if currentHash == legacyHash {
		t.Fatal("Test setup error: current and legacy hashes should differ")
	}

	addr := &model.WhitelistedAddress{
		SignedAddress: &model.SignedWhitelistedAddress{
			Signatures: []model.WhitelistSignature{
				{
					// Signature covers the LEGACY hash (what was originally signed)
					Hashes: []string{legacyHash},
					UserSignature: &model.WhitelistUserSignature{
						UserID:    "team1@bank.com",
						Signature: "somesignature",
					},
				},
			},
		},
		Metadata: &model.WhitelistedAssetMetadata{
			// Current hash computed from current payload (won't match directly)
			Hash:            currentHash,
			PayloadAsString: currentPayload,
		},
	}

	foundHash, err := v.verifyHashInSignedHashes(addr)
	if err != nil {
		t.Errorf("Step 4 should pass with legacy hash: verifyHashInSignedHashes() unexpected error: %v", err)
		return
	}

	// Should return the legacy hash since that's what was signed
	if foundHash != legacyHash {
		t.Errorf("Step 4 should return legacy hash: got %q, want %q", foundHash, legacyHash)
	}
}

// TestStep4_VerifyHashCoverageFailure verifies that step 4 fails with IntegrityError
// when the metadata hash is not found in any signature's hashes list.
func TestStep4_VerifyHashCoverageFailure(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	v := NewWhitelistedAddressVerifier([]*ecdsa.PublicKey{&key.PublicKey}, 1)

	// Use a hash that won't match any signatures
	unrelatedHash := "1111111111111111111111111111111111111111111111111111111111111111"

	addr := &model.WhitelistedAddress{
		SignedAddress: &model.SignedWhitelistedAddress{
			Signatures: []model.WhitelistSignature{
				{
					Hashes: []string{"different_hash_1", "different_hash_2"},
					UserSignature: &model.WhitelistUserSignature{
						UserID:    "team1@bank.com",
						Signature: "somesignature",
					},
				},
			},
		},
		Metadata: &model.WhitelistedAssetMetadata{
			Hash: unrelatedHash,
			// Payload that also won't produce legacy hashes matching the signatures
			PayloadAsString: `{"currency":"TEST","address":"test"}`,
		},
	}

	_, err := v.verifyHashInSignedHashes(addr)
	if err == nil {
		t.Error("Step 4 should fail: verifyHashInSignedHashes() expected IntegrityError, got nil")
		return
	}

	// Verify it's an IntegrityError
	if _, ok := err.(*model.IntegrityError); !ok {
		t.Errorf("Step 4 should return IntegrityError, got %T", err)
	}
}

// =============================================================================
// Step 5 Tests: Verify Whitelist Signatures
// =============================================================================

// TestStep5_VerifyWhitelistSignaturesSuccess verifies that step 5 passes when
// signatures meet the governance threshold requirements.
func TestStep5_VerifyWhitelistSignaturesSuccess(t *testing.T) {
	// Generate a key pair for signing
	userKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	v := NewWhitelistedAddressVerifier(nil, 0) // SuperAdmin keys not needed for step 5

	metadataHash := testdata.RealMetadataHash
	hashes := []string{metadataHash}

	// Sign the hashes array with the user's key
	hashesJSON, _ := json.Marshal(hashes)
	signature, err := crypto.SignData(userKey, hashesJSON)
	if err != nil {
		t.Fatalf("Failed to sign hashes: %v", err)
	}

	addr := &model.WhitelistedAddress{
		ID:         "36663",
		Blockchain: "ALGO",
		Network:    "mainnet",
		SignedAddress: &model.SignedWhitelistedAddress{
			Signatures: []model.WhitelistSignature{
				{
					Hashes: hashes,
					UserSignature: &model.WhitelistUserSignature{
						UserID:    "team1@bank.com",
						Signature: signature,
					},
				},
			},
		},
	}

	rulesContainer := &model.DecodedRulesContainer{
		Users: []*model.RuleUser{
			{
				ID:        "team1@bank.com",
				PublicKey: &userKey.PublicKey,
				Roles:     []string{"USER", "OPERATOR"},
			},
		},
		Groups: []*model.RuleGroup{
			{ID: "team1", UserIDs: []string{"team1@bank.com"}},
		},
		AddressWhitelistingRules: []*model.AddressWhitelistingRules{
			{
				Currency: "ALGO",
				Network:  "mainnet",
				ParallelThresholds: []*model.SequentialThresholds{
					{Thresholds: []*model.GroupThreshold{
						{GroupID: "team1", MinimumSignatures: 1},
					}},
				},
			},
		},
	}

	err = v.verifyWhitelistSignatures(addr, rulesContainer, metadataHash)
	if err != nil {
		t.Errorf("Step 5 should pass: verifyWhitelistSignatures() unexpected error: %v", err)
	}
}

// TestStep5_VerifyWhitelistSignaturesFailure verifies that step 5 fails with
// WhitelistError when signatures do not meet the governance threshold.
func TestStep5_VerifyWhitelistSignaturesFailure(t *testing.T) {
	v := NewWhitelistedAddressVerifier(nil, 0) // SuperAdmin keys not needed for step 5

	metadataHash := testdata.RealMetadataHash

	addr := &model.WhitelistedAddress{
		ID:         "36663",
		Blockchain: "ALGO",
		Network:    "mainnet",
		SignedAddress: &model.SignedWhitelistedAddress{
			Signatures: []model.WhitelistSignature{
				{
					Hashes: []string{metadataHash},
					UserSignature: &model.WhitelistUserSignature{
						UserID:    "team1@bank.com",
						Signature: "invalid_signature", // Invalid signature
					},
				},
			},
		},
	}

	// Generate a different key than what was used to sign
	differentKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	rulesContainer := &model.DecodedRulesContainer{
		Users: []*model.RuleUser{
			{
				ID:        "team1@bank.com",
				PublicKey: &differentKey.PublicKey, // Different key = verification will fail
				Roles:     []string{"USER", "OPERATOR"},
			},
		},
		Groups: []*model.RuleGroup{
			{ID: "team1", UserIDs: []string{"team1@bank.com"}},
		},
		AddressWhitelistingRules: []*model.AddressWhitelistingRules{
			{
				Currency: "ALGO",
				Network:  "mainnet",
				ParallelThresholds: []*model.SequentialThresholds{
					{Thresholds: []*model.GroupThreshold{
						{GroupID: "team1", MinimumSignatures: 1},
					}},
				},
			},
		},
	}

	err := v.verifyWhitelistSignatures(addr, rulesContainer, metadataHash)
	if err == nil {
		t.Error("Step 5 should fail: verifyWhitelistSignatures() expected WhitelistError, got nil")
		return
	}

	// Verify it's a WhitelistError or IntegrityError (signature verification failure)
	if _, isWhitelist := err.(*model.WhitelistError); !isWhitelist {
		if _, isIntegrity := err.(*model.IntegrityError); !isIntegrity {
			t.Errorf("Step 5 should return WhitelistError or IntegrityError, got %T", err)
		}
	}
}

// TestStep5_VerifyWhitelistSignaturesFailure_InsufficientSignatures verifies that
// step 5 fails when there are not enough signatures to meet the threshold.
func TestStep5_VerifyWhitelistSignaturesFailure_InsufficientSignatures(t *testing.T) {
	// Generate a key pair for signing
	userKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	v := NewWhitelistedAddressVerifier(nil, 0)

	metadataHash := testdata.RealMetadataHash
	hashes := []string{metadataHash}

	// Sign with only one user
	hashesJSON, _ := json.Marshal(hashes)
	signature, _ := crypto.SignData(userKey, hashesJSON)

	addr := &model.WhitelistedAddress{
		ID:         "36663",
		Blockchain: "ALGO",
		Network:    "mainnet",
		SignedAddress: &model.SignedWhitelistedAddress{
			Signatures: []model.WhitelistSignature{
				{
					Hashes: hashes,
					UserSignature: &model.WhitelistUserSignature{
						UserID:    "user1@bank.com",
						Signature: signature,
					},
				},
			},
		},
	}

	rulesContainer := &model.DecodedRulesContainer{
		Users: []*model.RuleUser{
			{ID: "user1@bank.com", PublicKey: &userKey.PublicKey, Roles: []string{"USER"}},
			{ID: "user2@bank.com", PublicKey: &userKey.PublicKey, Roles: []string{"USER"}},
		},
		Groups: []*model.RuleGroup{
			{ID: "team1", UserIDs: []string{"user1@bank.com", "user2@bank.com"}},
		},
		AddressWhitelistingRules: []*model.AddressWhitelistingRules{
			{
				Currency: "ALGO",
				Network:  "mainnet",
				ParallelThresholds: []*model.SequentialThresholds{
					{Thresholds: []*model.GroupThreshold{
						// Requires 2 signatures but we only have 1
						{GroupID: "team1", MinimumSignatures: 2},
					}},
				},
			},
		},
	}

	err := v.verifyWhitelistSignatures(addr, rulesContainer, metadataHash)
	if err == nil {
		t.Error("Step 5 should fail: expected error for insufficient signatures")
		return
	}

	// Verify error message mentions the threshold requirement
	if _, isIntegrity := err.(*model.IntegrityError); !isIntegrity {
		if _, isWhitelist := err.(*model.WhitelistError); !isWhitelist {
			t.Errorf("Step 5 should return WhitelistError or IntegrityError, got %T", err)
		}
	}
}

// TestStep5_VerifyWhitelistSignaturesFailure_NoMatchingRules verifies that
// step 5 fails when there are no matching address whitelisting rules.
func TestStep5_VerifyWhitelistSignaturesFailure_NoMatchingRules(t *testing.T) {
	v := NewWhitelistedAddressVerifier(nil, 0)

	addr := &model.WhitelistedAddress{
		ID:         "36663",
		Blockchain: "ALGO",
		Network:    "mainnet",
		SignedAddress: &model.SignedWhitelistedAddress{
			Signatures: []model.WhitelistSignature{
				{Hashes: []string{"hash"}},
			},
		},
	}

	rulesContainer := &model.DecodedRulesContainer{
		// No AddressWhitelistingRules for ALGO/mainnet
		AddressWhitelistingRules: []*model.AddressWhitelistingRules{
			{
				Currency: "ETH",
				Network:  "mainnet",
				ParallelThresholds: []*model.SequentialThresholds{
					{Thresholds: []*model.GroupThreshold{
						{GroupID: "team1", MinimumSignatures: 1},
					}},
				},
			},
		},
	}

	err := v.verifyWhitelistSignatures(addr, rulesContainer, "hash")
	if err == nil {
		t.Error("Step 5 should fail: expected error when no matching rules")
		return
	}

	// Verify it's a WhitelistError
	if _, ok := err.(*model.WhitelistError); !ok {
		t.Errorf("Step 5 should return WhitelistError, got %T", err)
	}
}

// =============================================================================
// Helper Functions
// =============================================================================

// fixtureJSON represents the structure of the JSON fixture file.
type fixtureJSON struct {
	RulesSignatures         string          `json:"rulesSignatures"`
	SignedAddressSignatures json.RawMessage `json:"signedAddressSignatures"`
	RulesContainerJSON      json.RawMessage `json:"rulesContainerJson"`
}

// parseSuperAdminKeysFromFixture parses SuperAdmin public keys from the fixture.
func parseSuperAdminKeysFromFixture() ([]*ecdsa.PublicKey, error) {
	// These PEM keys are from the fixture's rulesContainerJson.users
	superAdminKeysPEM := []string{
		`-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEyWjh6d+PgOK3LqockShMcDMtAHIm
itWjoVSX/FzBAWvemeaeNnYDKzEXiDDgiq2tILFL1Chdkqofhp9EdBZOlQ==
-----END PUBLIC KEY-----`,
		`-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAELJhEUNLLHgI8LiWJaeJGpaBfdvgo
YyKsjSFyTMxECR/E+1qpzDlNNug7hDPgBPpZ3Z+U8QWjaKB4Mrbj2/kImQ==
-----END PUBLIC KEY-----`,
	}

	var keys []*ecdsa.PublicKey
	for _, pem := range superAdminKeysPEM {
		key, err := crypto.DecodePublicKeyPEM(pem)
		if err != nil {
			return nil, fmt.Errorf("failed to decode PEM: %w", err)
		}
		keys = append(keys, key)
	}
	return keys, nil
}

// loadRulesContainerFromFixture loads the rules container and signatures from the fixture.
func loadRulesContainerFromFixture(t *testing.T) (rulesContainer, rulesSignatures string) {
	t.Helper()
	// These are the base64-encoded values from whitelisted_address_raw_response.json
	// The rulesSignatures is the base64-encoded protobuf of UserSignatures
	return "CIQCEoABClgKFHN1cGVyYWRtaW4xQGJhbmsuY29tEkAST4lY73f1ffMV+dvzv3NaiWpGWRu1NM1BcbHU+NFBSIFamyiwWbbty+yjYqO31K8vaots3jRmbWmaL8iH167NClgKFHN1cGVyYWRtaW4yQGJhbmsuY29tEkBhYKingez/vk0OSOCOw2GpfPFe7TDE9AVvBcDJ4tkih/eW1guV1QG4GzJlpUUUNoLkqEp8tJksWV7+NKm1jxLd",
		"ClgKFHN1cGVyYWRtaW4xQGJhbmsuY29tEkAST4lY73f1ffMV+dvzv3NaiWpGWRu1NM1BcbHU+NFBSIFamyiwWbbty+yjYqO31K8vaots3jRmbWmaL8iH167NClgKFHN1cGVyYWRtaW4yQGJhbmsuY29tEkBhYKingez/vk0OSOCOw2GpfPFe7TDE9AVvBcDJ4tkih/eW1guV1QG4GzJlpUUUNoLkqEp8tJksWV7+NKm1jxLd"
}

// decodeUserSignaturesFromBase64 decodes user signatures from base64.
// This is a simplified implementation that returns mock signatures for testing.
func decodeUserSignaturesFromBase64(base64Data string) ([]*model.RuleUserSignature, error) {
	// Decode the base64 data to verify it's valid
	_, err := DecodeBase64(base64Data)
	if err != nil {
		return nil, err
	}

	// Return mock signatures that match the fixture
	// In a real implementation, this would parse the protobuf
	return []*model.RuleUserSignature{
		{
			UserID:    "superadmin1@bank.com",
			Signature: "Ek+JWO939X3zFfnb879zWolqRlkbtTTNQXGx1PjRQUiBWpsosVm27cvsY2Kjt9SvL2qLbN40Zm1pmi/Ih9euzQ==",
		},
		{
			UserID:    "superadmin2@bank.com",
			Signature: "YWCop4Hs/75NDkjgjsNhqXzxXuwwxPQFbwXAyeLZIof3ltYLldUBuBsyZaVFFDaC5KhKfLSZLFle/jSptY8S3Q==",
		},
	}, nil
}
