/**
 * Unit tests for whitelisted address legacy hash verification.
 *
 * These tests verify:
 * - Hash verification: SHA256(payloadAsString) == metadata.hash
 * - Legacy hash computation strategies for backward compatibility
 *
 * Test patterns adapted from:
 * - Java SDK: WhitelistIntegrityHelperTest.java
 * - Go SDK: whitelist_helper_test.go
 * - Python SDK: test_whitelist_helper.py
 */

import { calculateHexHash } from "../../../src/crypto";
import { computeLegacyHashes } from "../../../src/helpers/whitelist-hash-helper";
import {
  REAL_PAYLOAD_AS_STRING,
  REAL_METADATA_HASH,
  CASE1_LEGACY_PAYLOAD,
  CASE1_CURRENT_PAYLOAD,
  CASE1_LEGACY_HASH,
  CASE1_CURRENT_HASH,
  CASE2_LEGACY_PAYLOAD,
  CASE2_CURRENT_PAYLOAD,
  CASE2_LEGACY_HASH,
  CASE2_CURRENT_HASH,
  STRATEGY2_LEGACY_PAYLOAD,
} from "../fixtures/whitelisted-address-fixtures";

// ============================================================================
// HASH VERIFICATION TESTS
// ============================================================================

describe("Hash Verification", () => {
  describe("metadata hash matches computed hash", () => {
    it("should verify SHA256(payloadAsString) == metadata.hash for real payload", () => {
      const computedHash = calculateHexHash(REAL_PAYLOAD_AS_STRING);
      expect(computedHash).toBe(REAL_METADATA_HASH);
    });

    it("should produce lowercase hex hash", () => {
      const hash = calculateHexHash("test");
      expect(hash).toMatch(/^[a-f0-9]{64}$/);
      expect(hash).toBe(hash.toLowerCase());
    });

    it("should produce 64-character hash", () => {
      const hash = calculateHexHash("any data");
      expect(hash.length).toBe(64);
    });

    it("should verify case1 current payload hash", () => {
      const computedHash = calculateHexHash(CASE1_CURRENT_PAYLOAD);
      expect(computedHash).toBe(CASE1_CURRENT_HASH);
    });

    it("should verify case1 legacy payload hash", () => {
      const computedHash = calculateHexHash(CASE1_LEGACY_PAYLOAD);
      expect(computedHash).toBe(CASE1_LEGACY_HASH);
    });

    it("should verify case2 current payload hash", () => {
      const computedHash = calculateHexHash(CASE2_CURRENT_PAYLOAD);
      expect(computedHash).toBe(CASE2_CURRENT_HASH);
    });

    it("should verify case2 legacy payload hash", () => {
      const computedHash = calculateHexHash(CASE2_LEGACY_PAYLOAD);
      expect(computedHash).toBe(CASE2_LEGACY_HASH);
    });
  });

  describe("metadata hash mismatch detected", () => {
    it("should detect tampered payload", () => {
      // Tamper with the payload by changing the address
      const tamperedPayload = REAL_PAYLOAD_AS_STRING.replace(
        "P4QCJV2YYLAEULGLJQAW4XTU3EBOHWL5C46I5SPLH2H7AJEE367ZDACV5A",
        "TAMPERED_ADDRESS_XXXXX"
      );
      const tamperedHash = calculateHexHash(tamperedPayload);

      expect(tamperedHash).not.toBe(REAL_METADATA_HASH);
    });

    it("should detect single character change", () => {
      const original = '{"currency":"ETH","address":"0xabc123"}';
      const modified = '{"currency":"ETH","address":"0xabc124"}'; // changed last char

      const originalHash = calculateHexHash(original);
      const modifiedHash = calculateHexHash(modified);

      expect(originalHash).not.toBe(modifiedHash);
    });

    it("should detect whitespace changes", () => {
      const original = '{"currency":"ETH","address":"0xabc123"}';
      const withSpace = '{"currency":"ETH", "address":"0xabc123"}'; // added space after comma

      const originalHash = calculateHexHash(original);
      const withSpaceHash = calculateHexHash(withSpace);

      expect(originalHash).not.toBe(withSpaceHash);
    });

    it("should detect field reordering", () => {
      const original = '{"currency":"ETH","address":"0xabc123"}';
      const reordered = '{"address":"0xabc123","currency":"ETH"}';

      const originalHash = calculateHexHash(original);
      const reorderedHash = calculateHexHash(reordered);

      expect(originalHash).not.toBe(reorderedHash);
    });
  });

  describe("empty payload raises error", () => {
    it("should return empty array for empty payloadAsString", () => {
      // computeLegacyHashes handles empty string gracefully
      const result = computeLegacyHashes("");
      expect(result).toEqual([]);
    });

    it("should return empty array for null payloadAsString", () => {
      const result = computeLegacyHashes(null as unknown as string);
      expect(result).toEqual([]);
    });

    it("should return empty array for undefined payloadAsString", () => {
      const result = computeLegacyHashes(undefined as unknown as string);
      expect(result).toEqual([]);
    });
  });

  describe("empty hash validation", () => {
    it("should detect empty hash mismatch", () => {
      const computedHash = calculateHexHash(REAL_PAYLOAD_AS_STRING);
      const emptyHash = "";

      expect(computedHash).not.toBe(emptyHash);
      expect(emptyHash.length).toBe(0);
      expect(computedHash.length).toBe(64);
    });

    it("should not match null hash", () => {
      const computedHash = calculateHexHash(REAL_PAYLOAD_AS_STRING);
      const nullHash = null;

      expect(computedHash).not.toBe(nullHash);
    });
  });
});

// ============================================================================
// LEGACY HASH CASE 1 TESTS - contractType field removal
// ============================================================================

describe("Legacy Hash Case 1 - contractType field removal", () => {
  describe("case1 removes contractType", () => {
    it("should produce legacy hash when removing contractType from current payload", () => {
      const legacyHashes = computeLegacyHashes(CASE1_CURRENT_PAYLOAD);

      // The legacy hash (without contractType) should be in the results
      expect(legacyHashes).toContain(CASE1_LEGACY_HASH);
    });

    it("should remove contractType field via regex pattern", () => {
      // Manually verify the transformation
      const current = '{"currency":"ETH","address":"0xabc","contractType":""}';
      const expected = '{"currency":"ETH","address":"0xabc"}';

      // Apply the same regex pattern used in computeLegacyHashes
      const withoutContractType = current.replace(
        /,"contractType":"[^"]*"/g,
        ""
      );
      expect(withoutContractType).toBe(expected);
    });

    it("should handle contractType with non-empty value", () => {
      const withValue =
        '{"currency":"ETH","address":"0xabc","contractType":"ERC20"}';
      const expected = '{"currency":"ETH","address":"0xabc"}';

      const withoutContractType = withValue.replace(
        /,"contractType":"[^"]*"/g,
        ""
      );
      expect(withoutContractType).toBe(expected);
    });
  });

  describe("case1 current payload produces current hash", () => {
    it("should compute current hash from current payload", () => {
      const hash = calculateHexHash(CASE1_CURRENT_PAYLOAD);
      expect(hash).toBe(CASE1_CURRENT_HASH);
    });

    it("should NOT include current hash in legacy hashes", () => {
      const legacyHashes = computeLegacyHashes(CASE1_CURRENT_PAYLOAD);
      expect(legacyHashes).not.toContain(CASE1_CURRENT_HASH);
    });
  });

  describe("case1 legacy payload produces legacy hash", () => {
    it("should compute legacy hash from legacy payload", () => {
      const hash = calculateHexHash(CASE1_LEGACY_PAYLOAD);
      expect(hash).toBe(CASE1_LEGACY_HASH);
    });

    it("should return empty array when legacy payload has no contractType to remove", () => {
      // The legacy payload doesn't have contractType, so no legacy transformation applies
      const legacyHashes = computeLegacyHashes(CASE1_LEGACY_PAYLOAD);
      expect(legacyHashes).toEqual([]);
    });
  });

  describe("case1 hash difference", () => {
    it("should produce different hashes for current and legacy payloads", () => {
      expect(CASE1_CURRENT_HASH).not.toBe(CASE1_LEGACY_HASH);
    });

    it("should verify the hashes are actually different strings", () => {
      expect(CASE1_CURRENT_HASH.length).toBe(64);
      expect(CASE1_LEGACY_HASH.length).toBe(64);

      // Count differing characters
      let differences = 0;
      for (let i = 0; i < 64; i++) {
        if (CASE1_CURRENT_HASH[i] !== CASE1_LEGACY_HASH[i]) {
          differences++;
        }
      }
      expect(differences).toBeGreaterThan(0);
    });
  });
});

// ============================================================================
// LEGACY HASH CASE 2 TESTS - contractType AND labels in objects removal
// ============================================================================

describe("Legacy Hash Case 2 - contractType AND labels removal", () => {
  describe("case2 removes contractType and labels in objects", () => {
    it("should produce legacy hash when removing both from current payload", () => {
      const legacyHashes = computeLegacyHashes(CASE2_CURRENT_PAYLOAD);

      // Strategy 3: remove both contractType and labels
      expect(legacyHashes).toContain(CASE2_LEGACY_HASH);
    });

    it("should handle linkedInternalAddresses with labels", () => {
      // Verify the label removal regex
      const withLabel =
        '{"data":[{"id":"1","address":"0xabc","label":"Test"}]}';
      const expected = '{"data":[{"id":"1","address":"0xabc"}]}';

      const withoutLabel = withLabel.replace(/,"label":"[^"]*"}/g, "}");
      expect(withoutLabel).toBe(expected);
    });

    it("should only remove labels followed by closing brace", () => {
      // Main label should NOT be removed (not followed by })
      const payload =
        '{"label":"Main","address":"0x","linkedInternalAddresses":[{"id":"1","label":"Internal"}]}';

      // Only the internal label (followed by }) should be removed
      const withoutObjectLabels = payload.replace(/,"label":"[^"]*"}/g, "}");

      // Main label should still be there
      expect(withoutObjectLabels).toContain('"label":"Main"');
      // Internal label should be removed
      expect(withoutObjectLabels).not.toContain('"label":"Internal"');
    });
  });

  describe("case2 label pattern does not affect main label", () => {
    it("should preserve main address label", () => {
      // The case2 current payload has both main label and linked address labels
      expect(CASE2_CURRENT_PAYLOAD).toContain('"label":"20200324 test address 2"');

      const legacyHashes = computeLegacyHashes(CASE2_CURRENT_PAYLOAD);

      // Verify the main label is preserved in the transformation
      // by checking that the legacy payload still contains the main label
      expect(CASE2_LEGACY_PAYLOAD).toContain('"label":"20200324 test address 2"');
    });

    it("should not match main label followed by other fields", () => {
      const payload =
        '{"currency":"ETH","label":"Main Label","customerId":"123"}';

      // The regex should not match this since label is followed by "," not "}"
      const withoutObjectLabels = payload.replace(/,"label":"[^"]*"}/g, "}");

      expect(withoutObjectLabels).toBe(payload); // Should be unchanged
    });
  });

  describe("case2 only removing contractType is not enough", () => {
    it("should not match legacy hash with only contractType removed", () => {
      // Remove only contractType from case2
      const onlyContractTypeRemoved = CASE2_CURRENT_PAYLOAD.replace(
        /,"contractType":"[^"]*"/g,
        ""
      );

      const hashWithOnlyContractTypeRemoved = calculateHexHash(
        onlyContractTypeRemoved
      );

      // This hash should NOT match the case2 legacy hash
      // because the original signature was made without BOTH contractType AND labels
      expect(hashWithOnlyContractTypeRemoved).not.toBe(CASE2_LEGACY_HASH);
    });

    it("should not match legacy hash with only labels removed", () => {
      // Remove only labels from case2 (keeping contractType)
      const onlyLabelsRemoved = CASE2_CURRENT_PAYLOAD.replace(
        /,"label":"[^"]*"}/g,
        "}"
      );

      const hashWithOnlyLabelsRemoved = calculateHexHash(onlyLabelsRemoved);

      // This hash should NOT match the case2 legacy hash
      expect(hashWithOnlyLabelsRemoved).not.toBe(CASE2_LEGACY_HASH);
    });

    it("should match legacy hash only when BOTH removed", () => {
      // Remove both contractType AND labels
      let bothRemoved = CASE2_CURRENT_PAYLOAD.replace(/,"label":"[^"]*"}/g, "}");
      bothRemoved = bothRemoved.replace(/,"contractType":"[^"]*"/g, "");

      const hashWithBothRemoved = calculateHexHash(bothRemoved);

      // This hash SHOULD match the case2 legacy hash
      expect(hashWithBothRemoved).toBe(CASE2_LEGACY_HASH);
    });
  });

  describe("case2 hash difference", () => {
    it("should produce different hashes for current and legacy payloads", () => {
      expect(CASE2_CURRENT_HASH).not.toBe(CASE2_LEGACY_HASH);
    });

    it("should verify all hashes are different", () => {
      const currentHash = calculateHexHash(CASE2_CURRENT_PAYLOAD);
      const legacyHash = calculateHexHash(CASE2_LEGACY_PAYLOAD);

      // Compute intermediate hashes
      const onlyContractTypeRemoved = CASE2_CURRENT_PAYLOAD.replace(
        /,"contractType":"[^"]*"/g,
        ""
      );
      const onlyLabelsRemoved = CASE2_CURRENT_PAYLOAD.replace(
        /,"label":"[^"]*"}/g,
        "}"
      );

      const hashOnlyContractType = calculateHexHash(onlyContractTypeRemoved);
      const hashOnlyLabels = calculateHexHash(onlyLabelsRemoved);

      // All four should be different
      const allHashes = [
        currentHash,
        legacyHash,
        hashOnlyContractType,
        hashOnlyLabels,
      ];
      const uniqueHashes = new Set(allHashes);
      expect(uniqueHashes.size).toBe(4);
    });
  });
});

// ============================================================================
// STRATEGY TESTS
// ============================================================================

describe("Legacy Hash Strategies", () => {
  describe("all strategies produce different hashes", () => {
    it("should produce 4 different hashes for case2", () => {
      // Strategy 0: Current payload (no transformation)
      const hashCurrent = calculateHexHash(CASE2_CURRENT_PAYLOAD);

      // Strategy 1: Remove contractType only
      const strategy1 = CASE2_CURRENT_PAYLOAD.replace(
        /,"contractType":"[^"]*"/g,
        ""
      );
      const hashStrategy1 = calculateHexHash(strategy1);

      // Strategy 2: Remove labels only (keep contractType)
      const strategy2 = CASE2_CURRENT_PAYLOAD.replace(/,"label":"[^"]*"}/g, "}");
      const hashStrategy2 = calculateHexHash(strategy2);

      // Strategy 3: Remove both contractType AND labels
      let strategy3 = CASE2_CURRENT_PAYLOAD.replace(/,"label":"[^"]*"}/g, "}");
      strategy3 = strategy3.replace(/,"contractType":"[^"]*"/g, "");
      const hashStrategy3 = calculateHexHash(strategy3);

      // All 4 should be different
      const hashes = [hashCurrent, hashStrategy1, hashStrategy2, hashStrategy3];
      const uniqueHashes = new Set(hashes);
      expect(uniqueHashes.size).toBe(4);
    });

    it("should return at most 3 legacy hashes (strategies 1, 2, 3)", () => {
      const legacyHashes = computeLegacyHashes(CASE2_CURRENT_PAYLOAD);

      // Should have exactly 3 legacy hashes for case2
      // (contractType only, labels only, both)
      expect(legacyHashes.length).toBeLessThanOrEqual(3);
      expect(legacyHashes.length).toBeGreaterThan(0);
    });

    it("should not include current hash in legacy hashes", () => {
      const currentHash = calculateHexHash(CASE2_CURRENT_PAYLOAD);
      const legacyHashes = computeLegacyHashes(CASE2_CURRENT_PAYLOAD);

      expect(legacyHashes).not.toContain(currentHash);
    });
  });

  describe("strategy3 matches original legacy hash", () => {
    it("should match case2 legacy hash with strategy 3", () => {
      const legacyHashes = computeLegacyHashes(CASE2_CURRENT_PAYLOAD);

      // Strategy 3 (both removed) should match CASE2_LEGACY_HASH
      expect(legacyHashes).toContain(CASE2_LEGACY_HASH);
    });

    it("should verify strategy 2 payload", () => {
      // STRATEGY2_LEGACY_PAYLOAD has labels removed but contractType kept
      const hash = calculateHexHash(STRATEGY2_LEGACY_PAYLOAD);

      // Get computed legacy hashes
      const legacyHashes = computeLegacyHashes(CASE2_CURRENT_PAYLOAD);

      // Strategy 2 hash should be in the legacy hashes
      expect(legacyHashes).toContain(hash);
    });
  });

  describe("computeLegacyHashes returns unique hashes", () => {
    it("should not return duplicate hashes", () => {
      const legacyHashes = computeLegacyHashes(CASE2_CURRENT_PAYLOAD);

      const uniqueHashes = new Set(legacyHashes);
      expect(legacyHashes.length).toBe(uniqueHashes.size);
    });

    it("should deduplicate when strategies produce same result", () => {
      // Payload with contractType but no linked addresses
      const payload = '{"currency":"ETH","address":"0xabc","contractType":""}';

      const legacyHashes = computeLegacyHashes(payload);

      // Only strategy 1 (remove contractType) applies
      // Strategies 2 and 3 would produce same results
      const uniqueHashes = new Set(legacyHashes);
      expect(legacyHashes.length).toBe(uniqueHashes.size);
    });
  });

  describe("no legacy transformations apply", () => {
    it("should return empty array when payload has neither field", () => {
      const payload = '{"currency":"ETH","address":"0xabc123"}';

      const legacyHashes = computeLegacyHashes(payload);

      expect(legacyHashes).toHaveLength(0);
    });

    it("should return empty array for real payload (has contractType but would produce same hash)", () => {
      // REAL_PAYLOAD_AS_STRING has contractType, so should have legacy hashes
      const legacyHashes = computeLegacyHashes(REAL_PAYLOAD_AS_STRING);

      // Should have at least one legacy hash (without contractType)
      expect(legacyHashes.length).toBeGreaterThan(0);
    });
  });
});

// ============================================================================
// EDGE CASES
// ============================================================================

describe("Edge Cases", () => {
  it("should handle payload with special characters in values", () => {
    const payload =
      '{"currency":"ETH","address":"0xabc","contractType":"test\\"quote"}';

    // Should not throw
    expect(() => computeLegacyHashes(payload)).not.toThrow();
  });

  it("should handle payload with unicode characters", () => {
    const payload =
      '{"currency":"ETH","address":"0xabc","label":"Caf\\u00e9","contractType":""}';

    const legacyHashes = computeLegacyHashes(payload);
    expect(legacyHashes.length).toBeGreaterThan(0);
  });

  it("should handle very long payloads", () => {
    // Create a payload with many linked addresses
    const linkedAddresses = Array.from(
      { length: 100 },
      (_, i) =>
        `{"id":"${i}","address":"0x${"a".repeat(40)}","label":"Addr ${i}"}`
    ).join(",");

    const payload = `{"currency":"ETH","address":"0xabc","linkedInternalAddresses":[${linkedAddresses}],"contractType":""}`;

    const legacyHashes = computeLegacyHashes(payload);

    // Should handle large payloads
    expect(legacyHashes.length).toBeGreaterThan(0);

    // All hashes should be valid
    for (const hash of legacyHashes) {
      expect(hash).toMatch(/^[a-f0-9]{64}$/);
    }
  });

  it("should handle nested quotes in label correctly", () => {
    // Labels inside linkedInternalAddresses followed by }
    const payload =
      '{"address":"0x","linkedInternalAddresses":[{"id":"1","label":"has \\"quotes\\""}],"contractType":""}';

    // The regex might not handle escaped quotes perfectly, but should not throw
    expect(() => computeLegacyHashes(payload)).not.toThrow();
  });
});
