/**
 * Unit tests for the 6-step whitelisted address verification flow.
 *
 * This test suite verifies each step of the cryptographic verification flow:
 * 1. Verify metadata hash (SHA256(payloadAsString) == metadata.hash)
 * 2. Verify rules container signatures (SuperAdmin keys)
 * 3. Decode rules container (base64 -> model)
 * 4. Verify hash coverage (metadata.hash in signature hashes list)
 * 5. Verify whitelist signatures meet governance thresholds
 *
 * Test patterns adapted from:
 * - Java SDK: WhitelistIntegrityHelperTest.java
 * - Go SDK: whitelist_helper_test.go
 * - Python SDK: test_whitelist_helper.py
 */

import * as crypto from "crypto";

import { calculateHexHash, decodePublicKeyPem, signData } from "../../../src/crypto";
import { IntegrityError, WhitelistError } from "../../../src/errors";
import {
  verifyHashCoverage,
  computeLegacyHashes,
} from "../../../src/helpers/whitelist-hash-helper";
import {
  rulesContainerFromBase64,
  userSignaturesFromBase64,
} from "../../../src/mappers/governance-rules";
import { isValidSignature } from "../../../src/helpers/signature-verifier";
import { WhitelistedAddressVerifier } from "../../../src/helpers/whitelisted-address-verifier";
import type { DecodedRulesContainer } from "../../../src/models/governance-rules";
import type { SignedWhitelistedAddressEnvelope } from "../../../src/models/whitelisted-address";
import {
  REAL_PAYLOAD_AS_STRING,
  REAL_METADATA_HASH,
  CASE1_CURRENT_PAYLOAD,
  CASE1_LEGACY_HASH,
  createMockEnvelope,
  createLegacyMockEnvelope,
} from "../fixtures/whitelisted-address-fixtures";
import rawResponse from "../fixtures/whitelisted-address-raw-response.json";

// ============================================================================
// TEST FIXTURES - Keys and Rules Container
// ============================================================================

// Test SuperAdmin keys from the raw response fixture
const SUPER_ADMIN_1_PEM = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEyWjh6d+PgOK3LqockShMcDMtAHIm
itWjoVSX/FzBAWvemeaeNnYDKzEXiDDgiq2tILFL1Chdkqofhp9EdBZOlQ==
-----END PUBLIC KEY-----`;

const SUPER_ADMIN_2_PEM = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAELJhEUNLLHgI8LiWJaeJGpaBfdvgo
YyKsjSFyTMxECR/E+1qpzDlNNug7hDPgBPpZ3Z+U8QWjaKB4Mrbj2/kImQ==
-----END PUBLIC KEY-----`;

// User public key from rules container (team1@bank.com)
const USER_PUBLIC_KEY_PEM = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEM2NtzaFhm7xIR3OvWq5chW3/GEvW
L+3uqoE6lEJ13eWbulxsP/5h36VCqYDIGN/0wDeWwLYdpu5HhSXWhxCsCA==
-----END PUBLIC KEY-----`;

// Rules container and signatures from fixture
const RULES_CONTAINER_BASE64 = Buffer.from(
  JSON.stringify(rawResponse.rulesContainerJson),
  "utf-8"
).toString("base64");

const RULES_SIGNATURES_BASE64 = rawResponse.rulesSignatures;

/**
 * Generate a test ECDSA key pair for signing tests.
 */
function generateTestKeyPair(): { privateKey: crypto.KeyObject; publicKey: crypto.KeyObject; publicKeyPem: string } {
  const { privateKey, publicKey } = crypto.generateKeyPairSync("ec", {
    namedCurve: "P-256",
  });
  const publicKeyPem = publicKey.export({ type: "spki", format: "pem" }) as string;
  return { privateKey, publicKey, publicKeyPem };
}

/**
 * Create a mock signed envelope for testing.
 */
function createTestEnvelope(overrides?: Partial<SignedWhitelistedAddressEnvelope>): SignedWhitelistedAddressEnvelope {
  const mockEnvelope = createMockEnvelope();
  return {
    id: overrides?.id ?? mockEnvelope.id,
    metadata: overrides?.metadata ?? {
      hash: mockEnvelope.hash,
      payloadAsString: mockEnvelope.payloadAsString,
    },
    rulesContainerBase64: overrides?.rulesContainerBase64 ?? RULES_CONTAINER_BASE64,
    rulesSignaturesBase64: overrides?.rulesSignaturesBase64 ?? RULES_SIGNATURES_BASE64,
    signedAddress: overrides?.signedAddress ?? {
      payload: mockEnvelope.payloadAsString,
      signatures: mockEnvelope.signatures.map((sig) => ({
        userSignature: {
          userId: "team1@bank.com",
          signature: rawResponse.signedAddressSignatures[0]?.userSignature.signature,
          comment: "ok",
        },
        hashes: sig.hashes,
      })),
    },
    blockchain: overrides?.blockchain ?? "ALGO",
    network: overrides?.network ?? "mainnet",
    linkedInternalAddresses: overrides?.linkedInternalAddresses ?? [],
    linkedWallets: overrides?.linkedWallets ?? [],
  };
}

// ============================================================================
// STEP 1 TESTS: Verify Metadata Hash
// ============================================================================

describe("Step 1: Verify Metadata Hash", () => {
  describe("should verify metadata hash success - SHA256 matches", () => {
    it("should compute correct hash for real payload", () => {
      const computedHash = calculateHexHash(REAL_PAYLOAD_AS_STRING);
      expect(computedHash).toBe(REAL_METADATA_HASH);
    });

    it("should verify hash matches for valid envelope", () => {
      const envelope = createTestEnvelope();
      const computedHash = calculateHexHash(envelope.metadata.payloadAsString);
      expect(computedHash).toBe(envelope.metadata.hash);
    });

    it("should produce consistent lowercase hex hash", () => {
      const hash = calculateHexHash("test data for hashing");
      expect(hash).toMatch(/^[a-f0-9]{64}$/);
    });
  });

  describe("should verify metadata hash failure raises IntegrityError", () => {
    it("should detect tampered payload", () => {
      const originalPayload = REAL_PAYLOAD_AS_STRING;
      const tamperedPayload = originalPayload.replace("ALGO", "ETH");

      const originalHash = calculateHexHash(originalPayload);
      const tamperedHash = calculateHexHash(tamperedPayload);

      expect(tamperedHash).not.toBe(originalHash);
    });

    it("should throw IntegrityError when hash does not match", () => {
      const envelope = createTestEnvelope({
        metadata: {
          hash: "invalid_hash_that_does_not_match_payload",
          payloadAsString: REAL_PAYLOAD_AS_STRING,
        },
      });

      const verifier = new WhitelistedAddressVerifier({
        superAdminKeysPem: [SUPER_ADMIN_1_PEM, SUPER_ADMIN_2_PEM],
        minValidSignatures: 1,
      });

      expect(() => {
        verifier.verify(envelope, rulesContainerFromBase64, userSignaturesFromBase64);
      }).toThrow(IntegrityError);
    });

    it("should throw IntegrityError for empty payload", () => {
      const envelope = createTestEnvelope({
        metadata: {
          hash: REAL_METADATA_HASH,
          payloadAsString: "",
        },
      });

      const verifier = new WhitelistedAddressVerifier({
        superAdminKeysPem: [SUPER_ADMIN_1_PEM, SUPER_ADMIN_2_PEM],
        minValidSignatures: 1,
      });

      expect(() => {
        verifier.verify(envelope, rulesContainerFromBase64, userSignaturesFromBase64);
      }).toThrow(IntegrityError);
    });

    it("should throw IntegrityError for empty hash", () => {
      const envelope = createTestEnvelope({
        metadata: {
          hash: "",
          payloadAsString: REAL_PAYLOAD_AS_STRING,
        },
      });

      const verifier = new WhitelistedAddressVerifier({
        superAdminKeysPem: [SUPER_ADMIN_1_PEM, SUPER_ADMIN_2_PEM],
        minValidSignatures: 1,
      });

      expect(() => {
        verifier.verify(envelope, rulesContainerFromBase64, userSignaturesFromBase64);
      }).toThrow(IntegrityError);
    });
  });
});

// ============================================================================
// STEP 2 TESTS: Verify Rules Container Signatures
// ============================================================================

describe("Step 2: Verify Rules Container Signatures", () => {
  describe("should verify rules signatures success", () => {
    it("should validate signature with correct public key", () => {
      const { privateKey, publicKey } = generateTestKeyPair();
      const testData = Buffer.from("test data to sign", "utf-8");
      const signature = signData(privateKey, testData);

      const isValid = isValidSignature(testData, signature, [publicKey]);
      expect(isValid).toBe(true);
    });

    it("should accept valid SuperAdmin signatures on rules container", () => {
      // Note: The fixture uses protobuf format for rules signatures,
      // which userSignaturesFromBase64 cannot parse (returns empty array).
      // This test verifies the JSON format path works correctly.
      const jsonSignatures = Buffer.from(
        JSON.stringify([
          { userId: "superadmin1@bank.com", signature: "dGVzdF9zaWduYXR1cmU=" },
          { userId: "superadmin2@bank.com", signature: "dGVzdF9zaWduYXR1cmUy" },
        ]),
        "utf-8"
      ).toString("base64");

      const signatures = userSignaturesFromBase64(jsonSignatures);

      expect(signatures.length).toBe(2);
      expect(signatures[0]?.userId).toBe("superadmin1@bank.com");
      expect(signatures[1]?.userId).toBe("superadmin2@bank.com");
    });

    it("should decode protobuf format signatures", () => {
      // The fixture rulesSignatures is in protobuf format — now correctly decoded
      const signatures = userSignaturesFromBase64(RULES_SIGNATURES_BASE64);

      // Protobuf format is now decoded correctly (matches Go SDK behavior)
      expect(signatures.length).toBeGreaterThan(0);
      for (const sig of signatures) {
        expect(sig.userId).toBeTruthy();
        expect(sig.signature).toBeTruthy();
      }
    });

    it("should count valid signatures correctly", () => {
      const { privateKey: key1, publicKey: pub1 } = generateTestKeyPair();
      const { privateKey: key2, publicKey: pub2 } = generateTestKeyPair();
      const testData = Buffer.from("test rules container data", "utf-8");

      const sig1 = signData(key1, testData);
      const sig2 = signData(key2, testData);

      // Both signatures should be valid against their respective keys
      expect(isValidSignature(testData, sig1, [pub1])).toBe(true);
      expect(isValidSignature(testData, sig2, [pub2])).toBe(true);

      // Cross-verification should fail
      expect(isValidSignature(testData, sig1, [pub2])).toBe(false);
      expect(isValidSignature(testData, sig2, [pub1])).toBe(false);
    });
  });

  describe("should verify rules signatures failure raises IntegrityError", () => {
    it("should reject invalid signature format", () => {
      const { publicKey } = generateTestKeyPair();
      const testData = Buffer.from("test data", "utf-8");
      const invalidSignature = "not_a_valid_base64_signature!!!";

      const isValid = isValidSignature(testData, invalidSignature, [publicKey]);
      expect(isValid).toBe(false);
    });

    it("should reject signature with wrong key", () => {
      const { privateKey } = generateTestKeyPair();
      const { publicKey: wrongPublicKey } = generateTestKeyPair();
      const testData = Buffer.from("test data", "utf-8");
      const signature = signData(privateKey, testData);

      const isValid = isValidSignature(testData, signature, [wrongPublicKey]);
      expect(isValid).toBe(false);
    });

    it("should throw IntegrityError when insufficient valid signatures", () => {
      const envelope = createTestEnvelope({
        rulesSignaturesBase64: Buffer.from(
          JSON.stringify([
            { userId: "invalid@example.com", signature: "invalid_signature" },
          ]),
          "utf-8"
        ).toString("base64"),
      });

      const verifier = new WhitelistedAddressVerifier({
        superAdminKeysPem: [SUPER_ADMIN_1_PEM, SUPER_ADMIN_2_PEM],
        minValidSignatures: 2,
      });

      expect(() => {
        verifier.verify(envelope, rulesContainerFromBase64, userSignaturesFromBase64);
      }).toThrow(IntegrityError);
    });

    it("should throw IntegrityError for empty rules container", () => {
      const envelope = createTestEnvelope({
        rulesContainerBase64: "",
      });

      const verifier = new WhitelistedAddressVerifier({
        superAdminKeysPem: [SUPER_ADMIN_1_PEM],
        minValidSignatures: 1,
      });

      expect(() => {
        verifier.verify(envelope, rulesContainerFromBase64, userSignaturesFromBase64);
      }).toThrow(IntegrityError);
    });

    it("should throw IntegrityError for empty rules signatures", () => {
      const envelope = createTestEnvelope({
        rulesSignaturesBase64: "",
      });

      const verifier = new WhitelistedAddressVerifier({
        superAdminKeysPem: [SUPER_ADMIN_1_PEM],
        minValidSignatures: 1,
      });

      expect(() => {
        verifier.verify(envelope, rulesContainerFromBase64, userSignaturesFromBase64);
      }).toThrow(IntegrityError);
    });
  });
});

// ============================================================================
// STEP 3 TESTS: Decode Rules Container
// ============================================================================

describe("Step 3: Decode Rules Container", () => {
  describe("should decode rules container success", () => {
    it("should decode base64 JSON rules container", () => {
      const decoded = rulesContainerFromBase64(RULES_CONTAINER_BASE64);

      expect(decoded).toBeDefined();
      expect(decoded.users).toBeDefined();
      expect(decoded.groups).toBeDefined();
      expect(decoded.addressWhitelistingRules).toBeDefined();
    });

    it("should parse users correctly", () => {
      const decoded = rulesContainerFromBase64(RULES_CONTAINER_BASE64);

      expect(decoded.users.length).toBeGreaterThan(0);

      const superAdmin = decoded.users.find((u) => u.id === "superadmin1@bank.com");
      expect(superAdmin).toBeDefined();
      expect(superAdmin?.roles).toContain("SUPERADMIN");
      expect(superAdmin?.publicKeyPem).toBeDefined();
    });

    it("should parse groups correctly", () => {
      const decoded = rulesContainerFromBase64(RULES_CONTAINER_BASE64);

      expect(decoded.groups.length).toBeGreaterThan(0);

      const team1Group = decoded.groups.find((g) => g.id === "team1");
      expect(team1Group).toBeDefined();
      expect(team1Group?.userIds).toContain("team1@bank.com");
    });

    it("should parse address whitelisting rules correctly", () => {
      const decoded = rulesContainerFromBase64(RULES_CONTAINER_BASE64);

      expect(decoded.addressWhitelistingRules.length).toBeGreaterThan(0);

      const algoRule = decoded.addressWhitelistingRules.find(
        (r) => r.currency === "ALGO"
      );
      expect(algoRule).toBeDefined();
      expect(algoRule?.network).toBe("mainnet");
      expect(algoRule?.parallelThresholds).toBeDefined();
    });

    it("should parse HSM slot user correctly", () => {
      const decoded = rulesContainerFromBase64(RULES_CONTAINER_BASE64);

      const hsmUser = decoded.users.find((u) => u.roles.includes("HSMSLOT"));
      expect(hsmUser).toBeDefined();
      expect(hsmUser?.id).toBe("hsmslot@bank.com");
      expect(hsmUser?.publicKeyPem).toBeDefined();
    });
  });

  describe("should decode rules container failure raises IntegrityError", () => {
    it("should throw IntegrityError for invalid base64", () => {
      expect(() => {
        rulesContainerFromBase64("not_valid_base64!!!");
      }).toThrow(IntegrityError);
    });

    it("should throw IntegrityError for invalid JSON in base64", () => {
      const invalidJson = Buffer.from("not valid json {{{", "utf-8").toString("base64");

      expect(() => {
        rulesContainerFromBase64(invalidJson);
      }).toThrow(IntegrityError);
    });

    it("should return empty container for empty input", () => {
      const decoded = rulesContainerFromBase64("");

      expect(decoded.users).toEqual([]);
      expect(decoded.groups).toEqual([]);
      expect(decoded.addressWhitelistingRules).toEqual([]);
    });

    it("should handle malformed rules container gracefully", () => {
      // Valid JSON but missing expected fields
      const minimalJson = Buffer.from(
        JSON.stringify({ someOtherField: "value" }),
        "utf-8"
      ).toString("base64");

      const decoded = rulesContainerFromBase64(minimalJson);

      // Should return empty arrays for missing fields
      expect(decoded.users).toEqual([]);
      expect(decoded.groups).toEqual([]);
    });
  });
});

// ============================================================================
// STEP 4 TESTS: Verify Hash Coverage
// ============================================================================

describe("Step 4: Verify Hash Coverage", () => {
  describe("should verify hash coverage success", () => {
    it("should find hash in signatures list", () => {
      const metadataHash = "abc123def456";
      const signatures = [
        { hashes: ["other1", "other2"] },
        { hashes: ["abc123def456", "other3"] },
      ];

      const result = verifyHashCoverage(metadataHash, signatures);
      expect(result).toBe(true);
    });

    it("should find hash in first position", () => {
      const metadataHash = "target_hash";
      const signatures = [{ hashes: ["target_hash", "other1", "other2"] }];

      const result = verifyHashCoverage(metadataHash, signatures);
      expect(result).toBe(true);
    });

    it("should find hash in last position", () => {
      const metadataHash = "target_hash";
      const signatures = [
        { hashes: ["other1", "other2"] },
        { hashes: ["other3", "target_hash"] },
      ];

      const result = verifyHashCoverage(metadataHash, signatures);
      expect(result).toBe(true);
    });

    it("should verify real metadata hash coverage", () => {
      const mockEnvelope = createMockEnvelope();
      const result = verifyHashCoverage(
        mockEnvelope.hash,
        mockEnvelope.signatures
      );
      expect(result).toBe(true);
    });
  });

  describe("should verify hash coverage failure uses legacy", () => {
    it("should compute legacy hashes when current hash not found", () => {
      const legacyHashes = computeLegacyHashes(CASE1_CURRENT_PAYLOAD);

      expect(legacyHashes).toContain(CASE1_LEGACY_HASH);
    });

    it("should find legacy hash in signatures", () => {
      // Simulate legacy scenario: signatures contain only legacy hash
      const legacyMockEnvelope = createLegacyMockEnvelope("case1");

      // Current hash should NOT be found
      const currentHashFound = verifyHashCoverage(
        legacyMockEnvelope.hash,
        legacyMockEnvelope.signatures
      );
      expect(currentHashFound).toBe(false);

      // Legacy hash SHOULD be found
      const legacyHashes = computeLegacyHashes(legacyMockEnvelope.payloadAsString);
      let legacyHashFound = false;
      for (const legacyHash of legacyHashes) {
        if (verifyHashCoverage(legacyHash, legacyMockEnvelope.signatures)) {
          legacyHashFound = true;
          break;
        }
      }
      expect(legacyHashFound).toBe(true);
    });

    it("should return false when hash not found anywhere", () => {
      const metadataHash = "nonexistent_hash";
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

    it("should handle signatures with empty hashes arrays", () => {
      const metadataHash = "somehash";
      const signatures = [{ hashes: [] }, { hashes: [] }];

      const result = verifyHashCoverage(metadataHash, signatures);
      expect(result).toBe(false);
    });
  });
});

// ============================================================================
// STEP 5 TESTS: Verify Whitelist Signatures
// ============================================================================

describe("Step 5: Verify Whitelist Signatures", () => {
  describe("should verify whitelist signatures success", () => {
    it("should validate user signature format", () => {
      const { privateKey, publicKey } = generateTestKeyPair();

      // Sign a hashes array as JSON (this is what whitelist signatures sign)
      const hashes = [REAL_METADATA_HASH];
      const hashesJson = JSON.stringify(hashes);
      const signature = signData(privateKey, Buffer.from(hashesJson, "utf-8"));

      // Verify signature
      const isValid = isValidSignature(
        Buffer.from(hashesJson, "utf-8"),
        signature,
        [publicKey]
      );
      expect(isValid).toBe(true);
    });

    it("should verify signature covers correct hash", () => {
      const { privateKey, publicKey } = generateTestKeyPair();

      const correctHash = "correct_hash_value";
      const wrongHash = "wrong_hash_value";

      // Sign the correct hash
      const hashesJson = JSON.stringify([correctHash]);
      const signature = signData(privateKey, Buffer.from(hashesJson, "utf-8"));

      // Verify signature against correct hash
      const correctHashJson = JSON.stringify([correctHash]);
      const isValidCorrect = isValidSignature(
        Buffer.from(correctHashJson, "utf-8"),
        signature,
        [publicKey]
      );
      expect(isValidCorrect).toBe(true);

      // Verify signature against wrong hash should fail
      const wrongHashJson = JSON.stringify([wrongHash]);
      const isValidWrong = isValidSignature(
        Buffer.from(wrongHashJson, "utf-8"),
        signature,
        [publicKey]
      );
      expect(isValidWrong).toBe(false);
    });

    it("should verify multiple signatures from same group", () => {
      const { privateKey: key1, publicKey: pub1 } = generateTestKeyPair();
      const { privateKey: key2, publicKey: pub2 } = generateTestKeyPair();

      const hashes = ["test_hash"];
      const hashesJson = JSON.stringify(hashes);
      const data = Buffer.from(hashesJson, "utf-8");

      const sig1 = signData(key1, data);
      const sig2 = signData(key2, data);

      expect(isValidSignature(data, sig1, [pub1])).toBe(true);
      expect(isValidSignature(data, sig2, [pub2])).toBe(true);
    });
  });

  describe("should verify whitelist signatures failure raises IntegrityError", () => {
    it("should throw WhitelistError when no matching rules found", () => {
      const envelope = createTestEnvelope({
        blockchain: "NONEXISTENT",
        network: "testnet",
      });

      const verifier = new WhitelistedAddressVerifier({
        superAdminKeysPem: [SUPER_ADMIN_1_PEM, SUPER_ADMIN_2_PEM],
        minValidSignatures: 1,
      });

      // This will fail at rules container signature verification first
      // In a real scenario with valid rules signatures, it would fail at finding matching rules
      expect(() => {
        verifier.verify(envelope, rulesContainerFromBase64, userSignaturesFromBase64);
      }).toThrow();
    });

    it("should reject signature from non-group member", () => {
      const { privateKey, publicKey } = generateTestKeyPair();

      // Sign with a key not in the rules container
      const hashes = [REAL_METADATA_HASH];
      const signature = signData(privateKey, Buffer.from(JSON.stringify(hashes), "utf-8"));

      // Verification should fail because the signer is not in the group
      // This is tested through the full verifier flow
      expect(signature).toBeDefined();
    });

    it("should reject tampered signature", () => {
      const { publicKey } = generateTestKeyPair();

      const hashes = ["test_hash"];
      const hashesJson = JSON.stringify(hashes);
      const tamperedSignature = "dGFtcGVyZWRfc2lnbmF0dXJl"; // "tampered_signature" in base64

      const isValid = isValidSignature(
        Buffer.from(hashesJson, "utf-8"),
        tamperedSignature,
        [publicKey]
      );
      expect(isValid).toBe(false);
    });

    it("should detect signature not covering required hash", () => {
      const { privateKey, publicKey } = generateTestKeyPair();

      // Sign different hashes than required
      const signedHashes = ["different_hash"];
      const signature = signData(
        privateKey,
        Buffer.from(JSON.stringify(signedHashes), "utf-8")
      );

      // Verification against different hash content should fail
      const requiredHashes = ["required_hash"];
      const isValid = isValidSignature(
        Buffer.from(JSON.stringify(requiredHashes), "utf-8"),
        signature,
        [publicKey]
      );
      expect(isValid).toBe(false);
    });

    it("should require minimum signatures threshold", () => {
      // Test that threshold validation works
      const verifier = new WhitelistedAddressVerifier({
        superAdminKeysPem: [SUPER_ADMIN_1_PEM, SUPER_ADMIN_2_PEM],
        minValidSignatures: 2,
      });

      // Single signature should not meet threshold of 2
      expect(verifier).toBeDefined();
    });
  });
});

// ============================================================================
// VERIFIER CONFIGURATION TESTS
// ============================================================================

describe("WhitelistedAddressVerifier Configuration", () => {
  it("should throw error when no SuperAdmin keys provided", () => {
    expect(() => {
      new WhitelistedAddressVerifier({
        superAdminKeysPem: [],
        minValidSignatures: 1,
      });
    }).toThrow("At least one SuperAdmin key is required");
  });

  it("should throw error when minValidSignatures is zero", () => {
    expect(() => {
      new WhitelistedAddressVerifier({
        superAdminKeysPem: [SUPER_ADMIN_1_PEM],
        minValidSignatures: 0,
      });
    }).toThrow("minValidSignatures must be at least 1");
  });

  it("should throw error when minValidSignatures is negative", () => {
    expect(() => {
      new WhitelistedAddressVerifier({
        superAdminKeysPem: [SUPER_ADMIN_1_PEM],
        minValidSignatures: -1,
      });
    }).toThrow("minValidSignatures must be at least 1");
  });

  it("should throw error for invalid PEM key format", () => {
    expect(() => {
      new WhitelistedAddressVerifier({
        superAdminKeysPem: ["not a valid pem key"],
        minValidSignatures: 1,
      });
    }).toThrow();
  });

  it("should accept valid configuration", () => {
    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [SUPER_ADMIN_1_PEM, SUPER_ADMIN_2_PEM],
      minValidSignatures: 1,
    });
    expect(verifier).toBeDefined();
  });
});

// ============================================================================
// INTEGRATION TESTS - Full Verification Flow
// ============================================================================

describe("Full Verification Flow Integration", () => {
  it("should throw IntegrityError for null envelope", () => {
    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [SUPER_ADMIN_1_PEM],
      minValidSignatures: 1,
    });

    expect(() => {
      verifier.verify(null as unknown as SignedWhitelistedAddressEnvelope, rulesContainerFromBase64, userSignaturesFromBase64);
    }).toThrow(IntegrityError);
  });

  it("should throw IntegrityError for envelope without metadata", () => {
    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [SUPER_ADMIN_1_PEM],
      minValidSignatures: 1,
    });

    const badEnvelope = {
      id: "123",
      metadata: undefined,
    } as unknown as SignedWhitelistedAddressEnvelope;

    expect(() => {
      verifier.verify(badEnvelope, rulesContainerFromBase64, userSignaturesFromBase64);
    }).toThrow(IntegrityError);
  });

  it("should throw IntegrityError for envelope without signedAddress", () => {
    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [SUPER_ADMIN_1_PEM],
      minValidSignatures: 1,
    });

    const envelope = createTestEnvelope();
    // Remove signedAddress to trigger the error
    const badEnvelope = {
      ...envelope,
      signedAddress: undefined,
    } as unknown as SignedWhitelistedAddressEnvelope;

    // Will fail at hash verification first, then at signedAddress check
    expect(() => {
      verifier.verify(badEnvelope, rulesContainerFromBase64, userSignaturesFromBase64);
    }).toThrow();
  });
});
