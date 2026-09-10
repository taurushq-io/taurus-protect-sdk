/**
 * Unit tests for WhitelistedAssetVerifier.
 *
 * Tests the complete 5-step verification flow for whitelisted assets:
 * 1. Verify metadata hash
 * 2. Verify rules container signatures (SuperAdmin keys)
 * 3. Decode rules container
 * 4. Verify hash coverage
 * 5. Verify whitelist signatures meet governance thresholds
 */

import * as crypto from "crypto";

import { calculateHexHash, signData, verifySignature, encodePublicKeyPem } from "../../../src/crypto";
import { IntegrityError, WhitelistError } from "../../../src/errors";
import { WhitelistedAssetVerifier } from "../../../src/helpers/whitelisted-asset-verifier";
import type {
  DecodedRulesContainer,
  RuleUserSignature,
} from "../../../src/models/governance-rules";
import type {
  SignedWhitelistedAssetEnvelope,
} from "../../../src/models/whitelisted-asset";

import {
  generateP256KeyPair,
  keyToPem,
  buildAssetPayload,
  payloadToString,
  buildFullAssetFixture,
  type AssetTestFixture,
} from "../fixtures/whitelisted-asset-fixtures";
import { computeAssetLegacyHashes } from "../../../src/helpers/whitelist-hash-helper";

// =============================================================================
// Helpers
// =============================================================================

// =============================================================================
// Constructor Tests
// =============================================================================

describe("WhitelistedAssetVerifier constructor", () => {
  it("should create verifier with valid config", () => {
    const { saPem } = buildFullAssetFixture();
    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [saPem],
      minValidSignatures: 1,
    });
    expect(verifier).toBeDefined();
  });

  it("should throw if no SuperAdmin keys", () => {
    expect(() => {
      new WhitelistedAssetVerifier({
        superAdminKeysPem: [],
        minValidSignatures: 1,
      });
    }).toThrow("At least one SuperAdmin key is required");
  });

  it("should throw if minValidSignatures < 1", () => {
    const { saPem } = buildFullAssetFixture();
    expect(() => {
      new WhitelistedAssetVerifier({
        superAdminKeysPem: [saPem],
        minValidSignatures: 0,
      });
    }).toThrow("minValidSignatures must be at least 1");
  });

  it("should throw for invalid PEM key", () => {
    expect(() => {
      new WhitelistedAssetVerifier({
        superAdminKeysPem: ["not a valid PEM key"],
        minValidSignatures: 1,
      });
    }).toThrow();
  });
});

// =============================================================================
// Step 1 Tests: Metadata Hash
// =============================================================================

describe("WhitelistedAssetVerifier - Step 1: Metadata Hash", () => {
  it("should pass with valid hash", () => {
    const f = buildFullAssetFixture();
    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const result = verifier.verify(
      f.envelope,
      f.rulesContainerDecoder,
      f.userSignaturesDecoder
    );
    expect(result).toBeDefined();
    expect(result.verifiedHash).toBe(f.envelope.metadata.hash);
  });

  it("should throw IntegrityError for mismatched hash", () => {
    const f = buildFullAssetFixture();
    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const tampered: SignedWhitelistedAssetEnvelope = {
      ...f.envelope,
      metadata: {
        hash: "0".repeat(64),
        payloadAsString: f.envelope.metadata.payloadAsString,
      },
    };

    expect(() => {
      verifier.verify(tampered, f.rulesContainerDecoder, f.userSignaturesDecoder);
    }).toThrow(IntegrityError);
  });

  it("should throw IntegrityError for empty payload", () => {
    const f = buildFullAssetFixture();
    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const tampered: SignedWhitelistedAssetEnvelope = {
      ...f.envelope,
      metadata: { hash: "abc", payloadAsString: "" },
    };

    expect(() => {
      verifier.verify(tampered, f.rulesContainerDecoder, f.userSignaturesDecoder);
    }).toThrow("payloadAsString is empty");
  });

  it("should throw IntegrityError for empty hash", () => {
    const f = buildFullAssetFixture();
    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const tampered: SignedWhitelistedAssetEnvelope = {
      ...f.envelope,
      metadata: { hash: "", payloadAsString: '{"foo":"bar"}' },
    };

    expect(() => {
      verifier.verify(tampered, f.rulesContainerDecoder, f.userSignaturesDecoder);
    }).toThrow("metadata hash is empty");
  });

  it("should throw for null envelope", () => {
    const f = buildFullAssetFixture();
    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    expect(() => {
      verifier.verify(null as any, f.rulesContainerDecoder, f.userSignaturesDecoder);
    }).toThrow(IntegrityError);
  });

  it("should throw for null metadata", () => {
    const f = buildFullAssetFixture();
    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const tampered = { ...f.envelope, metadata: null } as any;

    expect(() => {
      verifier.verify(tampered, f.rulesContainerDecoder, f.userSignaturesDecoder);
    }).toThrow(IntegrityError);
  });
});

// =============================================================================
// Step 2 Tests: Rules Container Signatures
// =============================================================================

describe("WhitelistedAssetVerifier - Step 2: Rules Container Signatures", () => {
  it("should throw for empty rulesContainer", () => {
    const f = buildFullAssetFixture();
    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const tampered: SignedWhitelistedAssetEnvelope = {
      ...f.envelope,
      rulesContainerBase64: "",
    };

    expect(() => {
      verifier.verify(tampered, f.rulesContainerDecoder, f.userSignaturesDecoder);
    }).toThrow("rulesContainer is empty");
  });

  it("should throw for empty rulesSignatures", () => {
    const f = buildFullAssetFixture();
    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const tampered: SignedWhitelistedAssetEnvelope = {
      ...f.envelope,
      rulesSignaturesBase64: "",
    };

    expect(() => {
      verifier.verify(tampered, f.rulesContainerDecoder, f.userSignaturesDecoder);
    }).toThrow("rulesSignatures is empty");
  });

  it("should throw when insufficient valid signatures", () => {
    const f = buildFullAssetFixture();
    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 2, // Require 2 but only 1 is valid
    });

    expect(() => {
      verifier.verify(f.envelope, f.rulesContainerDecoder, f.userSignaturesDecoder);
    }).toThrow("rules container signature verification failed");
  });

  it("should throw when signature decoder fails", () => {
    const f = buildFullAssetFixture();
    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const badDecoder = (_b64: string): RuleUserSignature[] => {
      throw new Error("decode failed");
    };

    expect(() => {
      verifier.verify(f.envelope, f.rulesContainerDecoder, badDecoder);
    }).toThrow("failed to decode rules signatures");
  });
});

// =============================================================================
// Step 3 Tests: Decode Rules Container
// =============================================================================

describe("WhitelistedAssetVerifier - Step 3: Decode Rules Container", () => {
  it("should throw when decoder fails", () => {
    const f = buildFullAssetFixture();
    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const badDecoder = (_b64: string): DecodedRulesContainer => {
      throw new Error("invalid protobuf");
    };

    expect(() => {
      verifier.verify(f.envelope, badDecoder, f.userSignaturesDecoder);
    }).toThrow("failed to decode rules container");
  });
});

// =============================================================================
// Step 4 Tests: Hash Coverage
// =============================================================================

describe("WhitelistedAssetVerifier - Step 4: Hash Coverage", () => {
  it("should throw for null signedContractAddress", () => {
    const f = buildFullAssetFixture();
    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const tampered = {
      ...f.envelope,
      signedContractAddress: null as any,
    };

    expect(() => {
      verifier.verify(tampered, f.rulesContainerDecoder, f.userSignaturesDecoder);
    }).toThrow("signedContractAddress");
  });

  it("should throw for empty signatures in signedContractAddress", () => {
    const f = buildFullAssetFixture();
    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const tampered: SignedWhitelistedAssetEnvelope = {
      ...f.envelope,
      signedContractAddress: {
        payload: undefined,
        signatures: [],
      },
    };

    expect(() => {
      verifier.verify(tampered, f.rulesContainerDecoder, f.userSignaturesDecoder);
    }).toThrow("no signatures in signedContractAddress");
  });

  it("should throw when hash not covered by any signature", () => {
    const f = buildFullAssetFixture();
    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const tampered: SignedWhitelistedAssetEnvelope = {
      ...f.envelope,
      signedContractAddress: {
        payload: undefined,
        signatures: [
          {
            userSignature: { userId: "u1", signature: "sig", comment: undefined },
            hashes: ["wrong_hash"],
          },
        ],
      },
    };

    expect(() => {
      verifier.verify(tampered, f.rulesContainerDecoder, f.userSignaturesDecoder);
    }).toThrow("metadata hash is not covered by any signature");
  });

  it("should pass when signatures contain current hash among other hashes", () => {
    const f = buildFullAssetFixture();
    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const metadataHash = f.envelope.metadata.hash;

    // Hashes array has the current hash plus other hashes
    const hashes = ["other_hash_1", metadataHash, "other_hash_2"];
    const hashesJson = JSON.stringify(hashes);
    const userSig = signData(f.userPriv, Buffer.from(hashesJson, "utf-8"));

    const envelope: SignedWhitelistedAssetEnvelope = {
      ...f.envelope,
      signedContractAddress: {
        payload: undefined,
        signatures: [
          {
            userSignature: {
              userId: "user1@bank.com",
              signature: userSig,
              comment: undefined,
            },
            hashes,
          },
        ],
      },
    };

    const result = verifier.verify(
      envelope,
      f.rulesContainerDecoder,
      f.userSignaturesDecoder
    );
    expect(result).toBeDefined();
    expect(result.verifiedHash).toBe(metadataHash);
  });
});

// =============================================================================
// Step 5 Tests: Whitelist Signatures
// =============================================================================

describe("WhitelistedAssetVerifier - Step 5: Whitelist Signatures", () => {
  it("should throw when no rules for blockchain", () => {
    const f = buildFullAssetFixture();
    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    // Decoder returns rules for BTC, not ETH
    const btcDecoder = (_b64: string): DecodedRulesContainer => ({
      users: [],
      groups: [],
      addressWhitelistingRules: [],
      contractAddressWhitelistingRules: [
        {
          blockchain: "BTC",
          network: "mainnet",
          parallelThresholds: [
            { thresholds: [{ groupId: "approvers", minimumSignatures: 1, threshold: 0 }] },
          ],
        },
      ],
      transactionRules: [],
      minimumDistinctUserSignatures: 0,
      minimumDistinctGroupSignatures: 0,
      enforcedRulesHash: "",
      timestamp: 0,
      hsmSlotId: 0,
      minimumCommitmentSignatures: 0,
      engineIdentities: [],
    });

    expect(() => {
      verifier.verify(f.envelope, btcDecoder, f.userSignaturesDecoder);
    }).toThrow(WhitelistError);
  });

  it("should throw when no thresholds defined", () => {
    const f = buildFullAssetFixture();
    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const emptyThresholdsDecoder = (_b64: string): DecodedRulesContainer => ({
      users: [],
      groups: [],
      addressWhitelistingRules: [],
      contractAddressWhitelistingRules: [
        {
          blockchain: "ETH",
          network: "mainnet",
          parallelThresholds: [],
        },
      ],
      transactionRules: [],
      minimumDistinctUserSignatures: 0,
      minimumDistinctGroupSignatures: 0,
      enforcedRulesHash: "",
      timestamp: 0,
      hsmSlotId: 0,
      minimumCommitmentSignatures: 0,
      engineIdentities: [],
    });

    expect(() => {
      verifier.verify(f.envelope, emptyThresholdsDecoder, f.userSignaturesDecoder);
    }).toThrow("no threshold rules defined");
  });

  it("should throw when threshold requires more sigs than available", () => {
    const f = buildFullAssetFixture();
    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const twoSigsDecoder = (_b64: string): DecodedRulesContainer => ({
      users: [
        { id: "user1@bank.com", name: "User 1", publicKeyPem: f.userPem, roles: ["USER"] },
      ],
      groups: [
        { id: "approvers", name: "Approvers", userIds: ["user1@bank.com"] },
      ],
      addressWhitelistingRules: [],
      contractAddressWhitelistingRules: [
        {
          blockchain: "ETH",
          network: "mainnet",
          parallelThresholds: [
            { thresholds: [{ groupId: "approvers", minimumSignatures: 2, threshold: 0 }] },
          ],
        },
      ],
      transactionRules: [],
      minimumDistinctUserSignatures: 0,
      minimumDistinctGroupSignatures: 0,
      enforcedRulesHash: "",
      timestamp: 0,
      hsmSlotId: 0,
      minimumCommitmentSignatures: 0,
      engineIdentities: [],
    });

    expect(() => {
      verifier.verify(f.envelope, twoSigsDecoder, f.userSignaturesDecoder);
    }).toThrow(WhitelistError);
  });
});

// =============================================================================
// End-to-End Happy Path
// =============================================================================

describe("WhitelistedAssetVerifier - End-to-End", () => {
  it("should succeed with valid data", () => {
    const f = buildFullAssetFixture();
    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const result = verifier.verify(
      f.envelope,
      f.rulesContainerDecoder,
      f.userSignaturesDecoder
    );

    expect(result).toBeDefined();
    expect(result.verifiedAsset).toBeDefined();
    expect(result.verifiedHash).toBe(f.envelope.metadata.hash);
  });

  it("should verify with two users and sequential thresholds (AND logic)", () => {
    const { privateKey: saPriv, publicKey: saPub } = generateP256KeyPair();
    const { privateKey: u1Priv, publicKey: u1Pub } = generateP256KeyPair();
    const { privateKey: u2Priv, publicKey: u2Pub } = generateP256KeyPair();
    const saPem = keyToPem(saPub);
    const u1Pem = keyToPem(u1Pub);
    const u2Pem = keyToPem(u2Pub);

    const payload = buildAssetPayload();
    const payloadStr = payloadToString(payload);
    const metadataHash = calculateHexHash(payloadStr);

    const hashes = [metadataHash];
    const hashesJson = JSON.stringify(hashes);
    const u1Sig = signData(u1Priv, Buffer.from(hashesJson));
    const u2Sig = signData(u2Priv, Buffer.from(hashesJson));

    const rulesB64 = Buffer.from("{}").toString("base64");
    const saSig = signData(saPriv, Buffer.from(rulesB64, "base64"));

    const envelope: SignedWhitelistedAssetEnvelope = {
      id: 2,
      metadata: { hash: metadataHash, payloadAsString: payloadStr },
      rulesContainerBase64: rulesB64,
      rulesSignaturesBase64: Buffer.from("dummy").toString("base64"),
      signedContractAddress: {
        payload: undefined,
        signatures: [
          { userSignature: { userId: "user1@bank.com", signature: u1Sig, comment: undefined }, hashes },
          { userSignature: { userId: "user2@bank.com", signature: u2Sig, comment: undefined }, hashes },
        ],
      },
      blockchain: "ETH",
      network: "mainnet",
    };

    const rcDecoder = (_b64: string): DecodedRulesContainer => ({
      users: [
        { id: "user1@bank.com", name: "User 1", publicKeyPem: u1Pem, roles: ["USER"] },
        { id: "user2@bank.com", name: "User 2", publicKeyPem: u2Pem, roles: ["USER"] },
      ],
      groups: [
        { id: "group_a", name: "Group A", userIds: ["user1@bank.com"] },
        { id: "group_b", name: "Group B", userIds: ["user2@bank.com"] },
      ],
      addressWhitelistingRules: [],
      contractAddressWhitelistingRules: [
        {
          blockchain: "ETH",
          network: "mainnet",
          parallelThresholds: [
            {
              thresholds: [
                { groupId: "group_a", minimumSignatures: 1, threshold: 0 },
                { groupId: "group_b", minimumSignatures: 1, threshold: 0 },
              ],
            },
          ],
        },
      ],
      transactionRules: [],
      minimumDistinctUserSignatures: 0,
      minimumDistinctGroupSignatures: 0,
      enforcedRulesHash: "",
      timestamp: 0,
      hsmSlotId: 0,
      minimumCommitmentSignatures: 0,
      engineIdentities: [],
    });

    const usDecoder = (_b64: string): RuleUserSignature[] => [
      { userId: "sa", signature: saSig },
    ];

    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [saPem],
      minValidSignatures: 1,
    });

    const result = verifier.verify(envelope, rcDecoder, usDecoder);
    expect(result).toBeDefined();
    expect(result.verifiedHash).toBe(metadataHash);
  });

  it("should succeed with parallel paths (OR logic)", () => {
    const f = buildFullAssetFixture();
    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    // First path fails (wrong group), second path succeeds
    const twoPathDecoder = (_b64: string): DecodedRulesContainer => ({
      users: [
        { id: "user1@bank.com", name: "User 1", publicKeyPem: f.userPem, roles: ["USER"] },
      ],
      groups: [
        { id: "approvers", name: "Approvers", userIds: ["user1@bank.com"] },
        { id: "other_team", name: "Other Team", userIds: ["nobody@bank.com"] },
      ],
      addressWhitelistingRules: [],
      contractAddressWhitelistingRules: [
        {
          blockchain: "ETH",
          network: "mainnet",
          parallelThresholds: [
            // Path 1: will fail
            { thresholds: [{ groupId: "other_team", minimumSignatures: 1, threshold: 0 }] },
            // Path 2: will pass
            { thresholds: [{ groupId: "approvers", minimumSignatures: 1, threshold: 0 }] },
          ],
        },
      ],
      transactionRules: [],
      minimumDistinctUserSignatures: 0,
      minimumDistinctGroupSignatures: 0,
      enforcedRulesHash: "",
      timestamp: 0,
      hsmSlotId: 0,
      minimumCommitmentSignatures: 0,
      engineIdentities: [],
    });

    const result = verifier.verify(f.envelope, twoPathDecoder, f.userSignaturesDecoder);
    expect(result).toBeDefined();
  });

  it("should fail when ALL parallel paths fail", () => {
    const f = buildFullAssetFixture();
    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const allFailDecoder = (_b64: string): DecodedRulesContainer => ({
      users: [
        { id: "user1@bank.com", name: "User 1", publicKeyPem: f.userPem, roles: ["USER"] },
      ],
      groups: [
        { id: "team_a", name: "Team A", userIds: ["nobody_a@bank.com"] },
        { id: "team_b", name: "Team B", userIds: ["nobody_b@bank.com"] },
      ],
      addressWhitelistingRules: [],
      contractAddressWhitelistingRules: [
        {
          blockchain: "ETH",
          network: "mainnet",
          parallelThresholds: [
            { thresholds: [{ groupId: "team_a", minimumSignatures: 1, threshold: 0 }] },
            { thresholds: [{ groupId: "team_b", minimumSignatures: 1, threshold: 0 }] },
          ],
        },
      ],
      transactionRules: [],
      minimumDistinctUserSignatures: 0,
      minimumDistinctGroupSignatures: 0,
      enforcedRulesHash: "",
      timestamp: 0,
      hsmSlotId: 0,
      minimumCommitmentSignatures: 0,
      engineIdentities: [],
    });

    expect(() => {
      verifier.verify(f.envelope, allFailDecoder, f.userSignaturesDecoder);
    }).toThrow(WhitelistError);
  });
});

// =============================================================================
// Batch Verification
// =============================================================================

describe("WhitelistedAssetVerifier - Batch Verification", () => {
  it("should fail on first invalid envelope in strict mode", () => {
    const f = buildFullAssetFixture();
    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const validEnvelope = f.envelope;
    const invalidEnvelope: SignedWhitelistedAssetEnvelope = {
      ...f.envelope,
      metadata: {
        hash: "tampered_hash",
        payloadAsString: f.envelope.metadata.payloadAsString,
      },
    };

    const results: string[] = [];
    for (const env of [validEnvelope, invalidEnvelope]) {
      try {
        verifier.verify(env, f.rulesContainerDecoder, f.userSignaturesDecoder);
        results.push("pass");
      } catch (error) {
        if (error instanceof IntegrityError) {
          results.push("integrity_fail");
        } else {
          results.push("other_fail");
        }
      }
    }

    expect(results).toEqual(["pass", "integrity_fail"]);
  });
});

// =============================================================================
// Legacy hash carried into step 5
// =============================================================================

describe("WhitelistedAssetVerifier - legacy hash threading", () => {
  /**
   * An asset signed before `isNFT` entered the schema is covered by the LEGACY
   * hash, not by SHA-256 of today's payload. Step 4 matches the legacy hash, so
   * step 5 must look for that same hash. While step 4 returned void and step 5
   * re-read `metadata.hash`, such an asset passed step 4 and then failed step 5.
   */
  it("verifies an asset whose signature covers only the legacy hash", () => {
    const { privateKey: saPriv, publicKey: saPub } = generateP256KeyPair();
    const { privateKey: userPriv, publicKey: userPub } = generateP256KeyPair();

    // Today's payload carries isNFT; the signature predates it.
    const currentPayload =
      '{"blockchain":"ETH","network":"mainnet","contractAddress":"0xUSDC","name":"USDC","symbol":"USDC","decimals":6,"isNFT":false}';
    const legacyPayload =
      '{"blockchain":"ETH","network":"mainnet","contractAddress":"0xUSDC","name":"USDC","symbol":"USDC","decimals":6}';

    const currentHash = calculateHexHash(currentPayload);
    const legacyHash = calculateHexHash(legacyPayload);
    expect(legacyHash).not.toBe(currentHash);
    expect(computeAssetLegacyHashes(currentPayload)).toContain(legacyHash);

    const userId = "user1@bank.com";
    const groupId = "approvers";

    // The user signed the LEGACY hash only.
    const hashes = [legacyHash];
    const userSig = signData(userPriv, Buffer.from(JSON.stringify(hashes), "utf-8"));
    const rulesB64 = Buffer.from("{}").toString("base64");
    const saSig = signData(saPriv, Buffer.from(rulesB64, "base64"));

    const envelope: SignedWhitelistedAssetEnvelope = {
      id: 7,
      metadata: { hash: currentHash, payloadAsString: currentPayload },
      rulesContainerBase64: rulesB64,
      rulesSignaturesBase64: Buffer.from("dummy").toString("base64"),
      signedContractAddress: {
        payload: undefined,
        signatures: [
          { userSignature: { userId, signature: userSig, comment: undefined }, hashes },
        ],
      },
      blockchain: "ETH",
      network: "mainnet",
    };

    const verifier = new WhitelistedAssetVerifier({
      superAdminKeysPem: [keyToPem(saPub)],
      minValidSignatures: 1,
    });

    const result = verifier.verify(
      envelope,
      () => ({
        users: [{ id: userId, name: "U", publicKeyPem: keyToPem(userPub), roles: ["USER"] }],
        groups: [{ id: groupId, name: "G", userIds: [userId] }],
        addressWhitelistingRules: [],
        contractAddressWhitelistingRules: [
          {
            blockchain: "ETH",
            network: "mainnet",
            parallelThresholds: [
              { thresholds: [{ groupId, minimumSignatures: 1, threshold: 0 }] },
            ],
          },
        ],
        transactionRules: [],
        minimumDistinctUserSignatures: 0,
        minimumDistinctGroupSignatures: 0,
        enforcedRulesHash: "",
        timestamp: 0,
        hsmSlotId: 0,
        minimumCommitmentSignatures: 0,
        engineIdentities: [],
      }),
      () => [{ userId: "sa@bank.com", signature: saSig }]
    );

    // The reported hash is the one actually covered, not the current one.
    expect(result.verifiedHash).toBe(legacyHash);
    expect(result.verifiedAsset.symbol).toBe("USDC");
  });
});
