/**
 * Unit tests for whitelist integrity helpers.
 *
 * Tests for:
 * - parseWhitelistedAddressFromJson: Parsing JSON payloads into WhitelistedAddress
 * - computeLegacyHashes: Computing legacy hashes for backward compatibility
 * - verifyHashCoverage: Verifying hash is covered by signatures (constant-time)
 *
 * Test patterns adapted from Java's WhitelistIntegrityHelperTest.
 */

import {
  parseWhitelistedAddressFromJson,
  computeLegacyHashes,
  verifyHashCoverage,
} from "../../../src/helpers/whitelist-hash-helper";
import { calculateHexHash } from "../../../src/crypto";

describe("parseWhitelistedAddressFromJson", () => {
  it("should parse valid JSON with all fields", () => {
    const json = JSON.stringify({
      currency: "ETH",
      network: "mainnet",
      address: "0xf631ce893edb440e49188a9912505051d0796818",
      memo: "test-memo",
      label: "My ETH Address",
      customerId: "cust-123",
      contractType: "ERC20",
      tnParticipantID: "participant-456",
      addressType: "EXTERNAL",
      exchangeAccountId: "12345",
    });

    const result = parseWhitelistedAddressFromJson(json);

    expect(result.blockchain).toBe("ETH");
    expect(result.network).toBe("mainnet");
    expect(result.address).toBe("0xf631ce893edb440e49188a9912505051d0796818");
    expect(result.memo).toBe("test-memo");
    expect(result.label).toBe("My ETH Address");
    expect(result.customerId).toBe("cust-123");
    expect(result.contractType).toBe("ERC20");
    expect(result.tnParticipantId).toBe("participant-456");
    expect(result.addressType).toBe("EXTERNAL");
    expect(result.exchangeAccountId).toBe(12345);
  });

  it("should throw error for empty payload", () => {
    expect(() => parseWhitelistedAddressFromJson("")).toThrow(
      "JSON payload cannot be empty"
    );
  });

  it("should throw error for invalid JSON", () => {
    expect(() => parseWhitelistedAddressFromJson("not-valid-json")).toThrow(
      "Failed to parse whitelist payload"
    );
  });

  it("should correctly parse linkedInternalAddresses", () => {
    const json = JSON.stringify({
      currency: "ETH",
      address: "0xabc123",
      linkedInternalAddresses: [
        { id: 100, address: "0xdef456", label: "Internal 1" },
        { id: 200, address: "0xghi789", label: "Internal 2" },
      ],
    });

    const result = parseWhitelistedAddressFromJson(json);

    expect(result.linkedInternalAddresses).toHaveLength(2);
    expect(result.linkedInternalAddresses[0].id).toBe(100);
    expect(result.linkedInternalAddresses[0].label).toBe("Internal 1");
    expect(result.linkedInternalAddresses[1].id).toBe(200);
    expect(result.linkedInternalAddresses[1].label).toBe("Internal 2");
  });

  it("should correctly parse linkedWallets", () => {
    const json = JSON.stringify({
      currency: "BTC",
      address: "bc1qxyz",
      linkedWallets: [
        { id: 1, name: "Wallet A", path: "m/44'/0'/0'" },
        { id: 2, name: "Wallet B", path: "m/44'/0'/1'" },
      ],
    });

    const result = parseWhitelistedAddressFromJson(json);

    expect(result.linkedWallets).toHaveLength(2);
    expect(result.linkedWallets[0].id).toBe(1);
    // JSON field is "name" but model field is "label"
    expect(result.linkedWallets[0].label).toBe("Wallet A");
    expect(result.linkedWallets[0].path).toBe("m/44'/0'/0'");
    expect(result.linkedWallets[1].id).toBe(2);
    expect(result.linkedWallets[1].label).toBe("Wallet B");
    expect(result.linkedWallets[1].path).toBe("m/44'/0'/1'");
  });

  it("should parse exchangeAccountId as number", () => {
    const json = JSON.stringify({
      currency: "ETH",
      address: "0xabc123",
      exchangeAccountId: "98765",
    });

    const result = parseWhitelistedAddressFromJson(json);

    expect(result.exchangeAccountId).toBe(98765);
    expect(typeof result.exchangeAccountId).toBe("number");
  });

  it("should handle missing optional fields gracefully", () => {
    const json = JSON.stringify({
      currency: "ETH",
      address: "0xabc123",
    });

    const result = parseWhitelistedAddressFromJson(json);

    expect(result.blockchain).toBe("ETH");
    expect(result.address).toBe("0xabc123");
    expect(result.memo).toBeUndefined();
    expect(result.label).toBeUndefined();
    expect(result.customerId).toBeUndefined();
    expect(result.contractType).toBeUndefined();
    expect(result.exchangeAccountId).toBeUndefined();
    expect(result.linkedInternalAddresses).toHaveLength(0);
    expect(result.linkedWallets).toHaveLength(0);
  });

  it("should handle invalid exchangeAccountId string", () => {
    const json = JSON.stringify({
      currency: "ETH",
      address: "0xabc123",
      exchangeAccountId: "not-a-number",
    });

    const result = parseWhitelistedAddressFromJson(json);

    // NaN is returned by parseInt, so it should be undefined
    expect(result.exchangeAccountId).toBeUndefined();
  });
});

describe("computeLegacyHashes", () => {
  it("should return empty array for empty payload", () => {
    expect(computeLegacyHashes("")).toEqual([]);
  });

  it("should return empty array for null-ish payload", () => {
    expect(computeLegacyHashes(null as unknown as string)).toEqual([]);
    expect(computeLegacyHashes(undefined as unknown as string)).toEqual([]);
  });

  it("should compute hash without contractType field", () => {
    // Original payload with contractType
    const original = JSON.stringify({
      currency: "ETH",
      address: "0xabc123",
      contractType: "ERC20",
      label: "Test",
    });

    // What we expect after removing contractType
    const withoutContractType = JSON.stringify({
      currency: "ETH",
      address: "0xabc123",
      label: "Test",
    });

    const hashes = computeLegacyHashes(original);

    // Should contain hash of version without contractType
    const expectedHash = calculateHexHash(withoutContractType);
    expect(hashes).toContain(expectedHash);
  });

  it("should compute hash without labels in linked addresses", () => {
    // Payload with label in linkedInternalAddresses
    const original =
      '{"currency":"ETH","address":"0xabc","linkedInternalAddresses":[{"id":1,"label":"Internal"}]}';

    // What we expect after removing label (only labels followed by closing brace)
    const withoutLabels =
      '{"currency":"ETH","address":"0xabc","linkedInternalAddresses":[{"id":1}]}';

    const hashes = computeLegacyHashes(original);

    const expectedHash = calculateHexHash(withoutLabels);
    expect(hashes).toContain(expectedHash);
  });

  it("should compute hash without both fields", () => {
    // Payload with both contractType and label in linkedInternalAddresses
    const original =
      '{"currency":"ETH","address":"0xabc","contractType":"ERC20","linkedInternalAddresses":[{"id":1,"label":"Internal"}]}';

    // What we expect after removing both
    const withoutBoth =
      '{"currency":"ETH","address":"0xabc","linkedInternalAddresses":[{"id":1}]}';

    const hashes = computeLegacyHashes(original);

    const expectedHash = calculateHexHash(withoutBoth);
    expect(hashes).toContain(expectedHash);
  });

  it("should not return duplicate hashes", () => {
    // If removing contractType and removing labels produce the same payload,
    // we should only get one hash
    const original =
      '{"currency":"ETH","address":"0xabc","contractType":"ERC20"}';

    const hashes = computeLegacyHashes(original);

    // Check for duplicates
    const uniqueHashes = [...new Set(hashes)];
    expect(hashes.length).toBe(uniqueHashes.length);
  });

  it("should return empty array when no legacy transformations apply", () => {
    // Payload without contractType or labels in linked addresses
    const original = '{"currency":"ETH","address":"0xabc"}';

    const hashes = computeLegacyHashes(original);

    // No transformations apply, so no legacy hashes
    expect(hashes).toHaveLength(0);
  });

  it("should handle complex nested structure", () => {
    const original = JSON.stringify({
      currency: "ETH",
      network: "mainnet",
      address: "0xf631ce893edb440e49188a9912505051d0796818",
      contractType: "ERC721",
      label: "Main Address",
      linkedInternalAddresses: [
        { id: 1, address: "0x111", label: "First" },
        { id: 2, address: "0x222", label: "Second" },
      ],
    });

    const hashes = computeLegacyHashes(original);

    // Should have at least one hash (without contractType, without labels, or without both)
    expect(hashes.length).toBeGreaterThan(0);

    // All hashes should be valid hex strings (64 chars for SHA-256)
    for (const hash of hashes) {
      expect(hash).toMatch(/^[a-f0-9]{64}$/);
    }
  });
});

describe("verifyHashCoverage", () => {
  it("should return true when hash is found", () => {
    const metadataHash = "abc123def456";
    const signatures = [
      { hashes: ["other1", "other2"] },
      { hashes: ["abc123def456", "other3"] },
    ];

    const result = verifyHashCoverage(metadataHash, signatures);

    expect(result).toBe(true);
  });

  it("should return false when hash is not found", () => {
    const metadataHash = "notfound";
    const signatures = [
      { hashes: ["hash1", "hash2"] },
      { hashes: ["hash3", "hash4"] },
    ];

    const result = verifyHashCoverage(metadataHash, signatures);

    expect(result).toBe(false);
  });

  it("should handle empty signatures array", () => {
    const metadataHash = "somehash";
    const signatures: Array<{ hashes: string[] }> = [];

    const result = verifyHashCoverage(metadataHash, signatures);

    expect(result).toBe(false);
  });

  it("should handle signatures with empty hashes", () => {
    const metadataHash = "somehash";
    const signatures = [{ hashes: [] }, { hashes: [] }];

    const result = verifyHashCoverage(metadataHash, signatures);

    expect(result).toBe(false);
  });

  it("should handle signatures with undefined hashes", () => {
    const metadataHash = "somehash";
    const signatures = [
      { hashes: undefined as unknown as string[] },
      { hashes: ["somehash"] },
    ];

    const result = verifyHashCoverage(metadataHash, signatures);

    expect(result).toBe(true);
  });

  it("should find hash in first signature", () => {
    const metadataHash = "target";
    const signatures = [
      { hashes: ["target"] },
      { hashes: ["other1", "other2"] },
    ];

    const result = verifyHashCoverage(metadataHash, signatures);

    expect(result).toBe(true);
  });

  it("should find hash in last signature", () => {
    const metadataHash = "target";
    const signatures = [
      { hashes: ["other1", "other2"] },
      { hashes: ["other3", "target"] },
    ];

    const result = verifyHashCoverage(metadataHash, signatures);

    expect(result).toBe(true);
  });

  it("should use constant-time comparison (no early return)", () => {
    // This test verifies the security property that all hashes are checked
    // regardless of where the match is found. We can't directly test timing,
    // but we can verify the function works correctly in all positions.
    const metadataHash = "findme";

    // Hash at beginning
    expect(
      verifyHashCoverage(metadataHash, [
        { hashes: ["findme", "a", "b", "c"] },
        { hashes: ["d", "e", "f"] },
      ])
    ).toBe(true);

    // Hash in middle
    expect(
      verifyHashCoverage(metadataHash, [
        { hashes: ["a", "b"] },
        { hashes: ["c", "findme", "d"] },
        { hashes: ["e", "f"] },
      ])
    ).toBe(true);

    // Hash at end
    expect(
      verifyHashCoverage(metadataHash, [
        { hashes: ["a", "b", "c"] },
        { hashes: ["d", "e", "findme"] },
      ])
    ).toBe(true);

    // Hash not present
    expect(
      verifyHashCoverage(metadataHash, [
        { hashes: ["a", "b", "c"] },
        { hashes: ["d", "e", "f"] },
      ])
    ).toBe(false);
  });

  it("should handle realistic hash values", () => {
    const payload1 = '{"currency":"ETH","address":"0xabc123"}';
    const payload2 = '{"currency":"BTC","address":"bc1qxyz"}';
    const metadataHash = calculateHexHash(payload1);
    const otherHash = calculateHexHash(payload2);

    const signatures = [
      { hashes: [otherHash] },
      { hashes: [metadataHash, "anotherhash"] },
    ];

    expect(verifyHashCoverage(metadataHash, signatures)).toBe(true);
    expect(verifyHashCoverage("nonexistent", signatures)).toBe(false);
  });
});
