package helper

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/testdata"
)

// =============================================================================
// Hash Verification Tests
// =============================================================================
// These tests verify that SHA-256 hash computation produces the expected results
// and that hash mismatches are properly detected.

// TestMetadataHashMatchesComputedHash verifies that SHA256(payloadAsString) equals metadata.hash
// using real API data from the fixtures.
func TestMetadataHashMatchesComputedHash(t *testing.T) {
	computedHash := crypto.CalculateHexHash(testdata.RealPayloadAsString)

	if computedHash != testdata.RealMetadataHash {
		t.Errorf("Hash mismatch:\n  computed: %s\n  expected: %s", computedHash, testdata.RealMetadataHash)
	}
}

// TestMetadataHashMismatchDetected verifies that a tampered payload is detected
// by hash verification.
func TestMetadataHashMismatchDetected(t *testing.T) {
	// Tamper with the payload by changing the address
	tamperedPayload := `{"currency":"ALGO","addressType":"individual","address":"TAMPERED_ADDRESS_12345","memo":"","label":"TN_Bank ACC Cockroach_WTRTest","customerId":"","exchangeAccountId":"","linkedInternalAddresses":[],"contractType":"","tnParticipantID":"84dc35e3-0af8-4b6b-be75-785f4b149d16"}`

	computedHash := crypto.CalculateHexHash(tamperedPayload)

	if computedHash == testdata.RealMetadataHash {
		t.Error("Tampered payload should NOT produce the same hash as the original")
	}
}

// TestEmptyPayloadRaisesError verifies that an empty payloadAsString fails verification.
func TestEmptyPayloadRaisesError(t *testing.T) {
	emptyPayload := ""

	// ComputeLegacyHashes returns nil for empty payload
	hashes := ComputeLegacyHashes(emptyPayload)
	if hashes != nil {
		t.Error("ComputeLegacyHashes should return nil for empty payload")
	}

	// ParseWhitelistedAddressFromJSON returns an error for empty payload
	_, err := ParseWhitelistedAddressFromJSON(emptyPayload)
	if err == nil {
		t.Error("ParseWhitelistedAddressFromJSON should return error for empty payload")
	}
}

// TestEmptyHashRaisesError verifies that verification fails when the metadata hash is empty.
func TestEmptyHashRaisesError(t *testing.T) {
	// When verifying, an empty hash should not match any computed hash
	emptyHash := ""
	computedHash := crypto.CalculateHexHash(testdata.RealPayloadAsString)

	if ConstantTimeCompare(computedHash, emptyHash) {
		t.Error("Empty metadata hash should not match any computed hash")
	}
}

// =============================================================================
// Legacy Hash Case 1 Tests: contractType removed
// =============================================================================
// These tests verify the legacy hash computation for addresses signed before
// the contractType field was added to the schema.

// TestCase1_RemovesContractType verifies that removing contractType from the current
// payload produces the legacy hash.
func TestCase1_RemovesContractType(t *testing.T) {
	// Compute the hash of the legacy payload directly
	legacyHash := crypto.CalculateHexHash(testdata.Case1LegacyPayload)

	// Verify it matches the expected legacy hash
	if legacyHash != testdata.Case1LegacyHash {
		t.Errorf("Legacy hash mismatch:\n  computed: %s\n  expected: %s", legacyHash, testdata.Case1LegacyHash)
	}

	// Verify that ComputeLegacyHashes produces this hash from the current payload
	hashes := ComputeLegacyHashes(testdata.Case1CurrentPayload)
	if len(hashes) == 0 {
		t.Fatal("ComputeLegacyHashes should produce at least one hash")
	}

	found := false
	for _, h := range hashes {
		if h == testdata.Case1LegacyHash {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("ComputeLegacyHashes did not produce the expected legacy hash.\n  Expected: %s\n  Got: %v", testdata.Case1LegacyHash, hashes)
	}
}

// TestCase1_CurrentPayloadProducesCurrentHash is a sanity check that the current
// payload produces the expected current hash.
func TestCase1_CurrentPayloadProducesCurrentHash(t *testing.T) {
	currentHash := crypto.CalculateHexHash(testdata.Case1CurrentPayload)

	if currentHash != testdata.Case1CurrentHash {
		t.Errorf("Current hash mismatch:\n  computed: %s\n  expected: %s", currentHash, testdata.Case1CurrentHash)
	}
}

// TestCase1_LegacyPayloadProducesLegacyHash is a sanity check that the legacy
// payload produces the expected legacy hash.
func TestCase1_LegacyPayloadProducesLegacyHash(t *testing.T) {
	legacyHash := crypto.CalculateHexHash(testdata.Case1LegacyPayload)

	if legacyHash != testdata.Case1LegacyHash {
		t.Errorf("Legacy hash mismatch:\n  computed: %s\n  expected: %s", legacyHash, testdata.Case1LegacyHash)
	}
}

// TestCase1_HashDifference verifies that the current and legacy hashes are different.
func TestCase1_HashDifference(t *testing.T) {
	if testdata.Case1CurrentHash == testdata.Case1LegacyHash {
		t.Error("Current and legacy hashes should be different")
	}

	currentHash := crypto.CalculateHexHash(testdata.Case1CurrentPayload)
	legacyHash := crypto.CalculateHexHash(testdata.Case1LegacyPayload)

	if currentHash == legacyHash {
		t.Error("Current and legacy hashes should be different")
	}
}

// =============================================================================
// Legacy Hash Case 2 Tests: contractType AND labels removed
// =============================================================================
// These tests verify the legacy hash computation for addresses signed before
// both contractType and labels in linkedInternalAddresses were added.

// TestCase2_RemovesContractTypeAndLabelsInObjects verifies that removing both
// contractType and labels from linkedInternalAddresses produces the legacy hash.
func TestCase2_RemovesContractTypeAndLabelsInObjects(t *testing.T) {
	// Compute the hash of the legacy payload directly
	legacyHash := crypto.CalculateHexHash(testdata.Case2LegacyPayload)

	// Verify it matches the expected legacy hash
	if legacyHash != testdata.Case2LegacyHash {
		t.Errorf("Legacy hash mismatch:\n  computed: %s\n  expected: %s", legacyHash, testdata.Case2LegacyHash)
	}

	// Verify that ComputeLegacyHashes produces this hash from the current payload
	hashes := ComputeLegacyHashes(testdata.Case2CurrentPayload)
	if len(hashes) == 0 {
		t.Fatal("ComputeLegacyHashes should produce at least one hash")
	}

	found := false
	for _, h := range hashes {
		if h == testdata.Case2LegacyHash {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("ComputeLegacyHashes did not produce the expected legacy hash.\n  Expected: %s\n  Got: %v", testdata.Case2LegacyHash, hashes)
	}
}

// TestCase2_LabelPatternDoesNotAffectMainLabel verifies that the main address label
// (not inside linkedInternalAddresses) is preserved when computing legacy hashes.
func TestCase2_LabelPatternDoesNotAffectMainLabel(t *testing.T) {
	// The main label "20200324 test address 2" should be preserved
	mainLabel := `"label":"20200324 test address 2"`

	// The legacy payload should still contain the main label
	if !containsString(testdata.Case2LegacyPayload, mainLabel) {
		t.Error("Legacy payload should still contain the main label")
	}

	// Apply the label removal pattern and verify main label is preserved
	withoutLabels := labelInObjectPattern.ReplaceAllString(testdata.Case2CurrentPayload, "}")
	if !containsString(withoutLabels, mainLabel) {
		t.Errorf("Main label should be preserved after removing object labels.\n  Result: %s", withoutLabels)
	}
}

// TestCase2_OnlyRemovingContractTypeIsNotEnough verifies that only removing
// contractType does not produce the correct legacy hash for Case 2.
func TestCase2_OnlyRemovingContractTypeIsNotEnough(t *testing.T) {
	// Remove only contractType
	withoutContractType := contractTypePattern.ReplaceAllString(testdata.Case2CurrentPayload, "")
	hashWithoutContractType := crypto.CalculateHexHash(withoutContractType)

	// This should NOT match the legacy hash (which requires both removals)
	if hashWithoutContractType == testdata.Case2LegacyHash {
		t.Error("Only removing contractType should NOT produce the Case 2 legacy hash")
	}
}

// TestCase2_HashDifference verifies that the current and legacy hashes are different.
func TestCase2_HashDifference(t *testing.T) {
	if testdata.Case2CurrentHash == testdata.Case2LegacyHash {
		t.Error("Current and legacy hashes should be different")
	}

	currentHash := crypto.CalculateHexHash(testdata.Case2CurrentPayload)
	legacyHash := crypto.CalculateHexHash(testdata.Case2LegacyPayload)

	if currentHash == legacyHash {
		t.Error("Current and legacy hashes should be different")
	}
}

// =============================================================================
// Strategy Tests
// =============================================================================
// These tests verify that all legacy hash strategies produce different hashes
// and that Strategy 3 (both removed) matches the original legacy hash.

// TestAllStrategiesProduceDifferentHashes verifies that Current, Strategy 1,
// Strategy 2, and Strategy 3 all produce different hashes.
func TestAllStrategiesProduceDifferentHashes(t *testing.T) {
	// Compute all strategy hashes for Case 2 (has both contractType and labels)
	currentHash := crypto.CalculateHexHash(testdata.Case2CurrentPayload)

	// Strategy 1: Remove contractType only
	withoutContractType := contractTypePattern.ReplaceAllString(testdata.Case2CurrentPayload, "")
	strategy1Hash := crypto.CalculateHexHash(withoutContractType)

	// Strategy 2: Remove labels only
	withoutLabels := labelInObjectPattern.ReplaceAllString(testdata.Case2CurrentPayload, "}")
	strategy2Hash := crypto.CalculateHexHash(withoutLabels)

	// Strategy 3: Remove both contractType AND labels
	withoutBoth := labelInObjectPattern.ReplaceAllString(testdata.Case2CurrentPayload, "}")
	withoutBoth = contractTypePattern.ReplaceAllString(withoutBoth, "")
	strategy3Hash := crypto.CalculateHexHash(withoutBoth)

	// All hashes should be different
	hashes := map[string]string{
		"current":   currentHash,
		"strategy1": strategy1Hash,
		"strategy2": strategy2Hash,
		"strategy3": strategy3Hash,
	}

	seen := make(map[string]string)
	for name, hash := range hashes {
		if prevName, exists := seen[hash]; exists {
			t.Errorf("Duplicate hash detected: %s and %s both produce %s", prevName, name, hash)
		}
		seen[hash] = name
	}

	// Log the hashes for debugging
	t.Logf("Current hash:   %s", currentHash)
	t.Logf("Strategy 1 hash (no contractType): %s", strategy1Hash)
	t.Logf("Strategy 2 hash (no labels):       %s", strategy2Hash)
	t.Logf("Strategy 3 hash (no both):         %s", strategy3Hash)
}

// TestStrategy3MatchesOriginalLegacyHash verifies that Strategy 3 (both removed)
// produces the expected CASE2_LEGACY_HASH.
func TestStrategy3MatchesOriginalLegacyHash(t *testing.T) {
	// Strategy 3: Remove both contractType AND labels
	withoutBoth := labelInObjectPattern.ReplaceAllString(testdata.Case2CurrentPayload, "}")
	withoutBoth = contractTypePattern.ReplaceAllString(withoutBoth, "")
	strategy3Hash := crypto.CalculateHexHash(withoutBoth)

	if strategy3Hash != testdata.Case2LegacyHash {
		t.Errorf("Strategy 3 hash does not match CASE2_LEGACY_HASH.\n  computed: %s\n  expected: %s", strategy3Hash, testdata.Case2LegacyHash)
	}
}

// TestStrategy2MatchesIntermediateLegacyPayload verifies that Strategy 2 (labels removed only)
// produces the hash of the intermediate legacy payload format.
func TestStrategy2MatchesIntermediateLegacyPayload(t *testing.T) {
	// Strategy 2: Remove labels only (keep contractType)
	withoutLabels := labelInObjectPattern.ReplaceAllString(testdata.Case2CurrentPayload, "}")
	strategy2Hash := crypto.CalculateHexHash(withoutLabels)

	// Compute the expected hash from the intermediate legacy payload
	expectedHash := crypto.CalculateHexHash(testdata.Strategy2LegacyPayload)

	if strategy2Hash != expectedHash {
		t.Errorf("Strategy 2 hash does not match Strategy2LegacyPayload hash.\n  computed: %s\n  expected: %s", strategy2Hash, expectedHash)
	}
}

// TestComputeLegacyHashesReturnsAllStrategies verifies that ComputeLegacyHashes
// returns hashes for all applicable strategies.
func TestComputeLegacyHashesReturnsAllStrategies(t *testing.T) {
	// For Case 2 which has both contractType and labels, we expect 3 hashes:
	// 1. contractType removed only
	// 2. labels removed only
	// 3. both removed
	hashes := ComputeLegacyHashes(testdata.Case2CurrentPayload)

	if len(hashes) != 3 {
		t.Errorf("ComputeLegacyHashes should return 3 hashes for Case 2, got %d: %v", len(hashes), hashes)
	}

	// Verify the Case 2 legacy hash is included
	found := false
	for _, h := range hashes {
		if h == testdata.Case2LegacyHash {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Case 2 legacy hash not found in ComputeLegacyHashes result.\n  Expected: %s\n  Got: %v", testdata.Case2LegacyHash, hashes)
	}
}

// TestComputeLegacyHashesForCase1 verifies ComputeLegacyHashes works correctly
// for Case 1 (contractType only, no labels in linkedInternalAddresses).
func TestComputeLegacyHashesForCase1(t *testing.T) {
	// Case 1 only has contractType, no labels in linkedInternalAddresses
	hashes := ComputeLegacyHashes(testdata.Case1CurrentPayload)

	// Should return 1 hash (contractType removed only)
	if len(hashes) != 1 {
		t.Errorf("ComputeLegacyHashes should return 1 hash for Case 1, got %d: %v", len(hashes), hashes)
	}

	// Verify the Case 1 legacy hash is included
	if len(hashes) > 0 && hashes[0] != testdata.Case1LegacyHash {
		t.Errorf("Case 1 legacy hash mismatch.\n  Expected: %s\n  Got: %s", testdata.Case1LegacyHash, hashes[0])
	}
}

// =============================================================================
// Helper functions
// =============================================================================

// containsString checks if a string contains a substring.
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

// containsSubstring is a helper for containsString.
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
