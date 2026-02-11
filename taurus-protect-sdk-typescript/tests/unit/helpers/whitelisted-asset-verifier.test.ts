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

// =============================================================================
// Helpers
// =============================================================================

function generateP256KeyPair(): {
  privateKey: crypto.KeyObject;
  publicKey: crypto.KeyObject;
} {
  const { privateKey, publicKey } = crypto.generateKeyPairSync("ec", {
    namedCurve: "P-256",
  });
  return { privateKey, publicKey };
}

function keyToPem(publicKey: crypto.KeyObject): string {
  return encodePublicKeyPem(publicKey);
}

function buildAssetPayload(overrides?: Partial<{
  blockchain: string;
  network: string;
  contractAddress: string;
  name: string;
  symbol: string;
  decimals: number;
  isNFT: boolean;
  kindType: string;
}>): Record<string, unknown> {
  return {
    blockchain: "ETH",
    network: "mainnet",
    contractAddress: "0xUSDC",
    name: "USDC",
    symbol: "USDC",
    decimals: 6,
    ...overrides,
  };
}

function payloadToString(payload: Record<string, unknown>): string {
  return JSON.stringify(payload);
}

interface AssetTestFixture {
  saPriv: crypto.KeyObject;
  saPub: crypto.KeyObject;
  userPriv: crypto.KeyObject;
  userPub: crypto.KeyObject;
  saPem: string;
  userPem: string;
  envelope: SignedWhitelistedAssetEnvelope;
  rulesContainerDecoder: (b64: string) => DecodedRulesContainer;
  userSignaturesDecoder: (b64: string) => RuleUserSignature[];
}

function buildFullAssetFixture(overrides?: {
  blockchain?: string;
  network?: string;
  groupId?: string;
  userId?: string;
}): AssetTestFixture {
  const { privateKey: saPriv, publicKey: saPub } = generateP256KeyPair();
  const { privateKey: userPriv, publicKey: userPub } = generateP256KeyPair();
  const saPem = keyToPem(saPub);
  const userPem = keyToPem(userPub);

  const blockchain = overrides?.blockchain ?? "ETH";
  const network = overrides?.network ?? "mainnet";
  const groupId = overrides?.groupId ?? "approvers";
  const userId = overrides?.userId ?? "user1@bank.com";

  const payload = buildAssetPayload({ blockchain, network });
  const payloadStr = payloadToString(payload);
  const metadataHash = calculateHexHash(payloadStr);

  // Build rules container
  const rulesJson = JSON.stringify({
    users: [{ id: userId, publicKey: userPem, roles: ["USER"] }],
    groups: [{ id: groupId, userIds: [userId] }],
    contractAddressWhitelistingRules: [
      {
        blockchain,
        network,
        parallelThresholds: [{ groupId, minimumSignatures: 1 }],
      },
    ],
  });
  const rulesB64 = Buffer.from(rulesJson).toString("base64");
  const rulesData = Buffer.from(rulesB64, "base64");

  // Sign rules container with SuperAdmin key
  const saSig = signData(saPriv, rulesData);

  // Sign hashes array with user key
  const hashes = [metadataHash];
  const hashesJson = JSON.stringify(hashes);
  const userSig = signData(userPriv, Buffer.from(hashesJson, "utf-8"));

  const envelope: SignedWhitelistedAssetEnvelope = {
    id: 1,
    metadata: {
      hash: metadataHash,
      payloadAsString: payloadStr,
    },
    rulesContainerBase64: rulesB64,
    rulesSignaturesBase64: Buffer.from("dummy").toString("base64"),
    signedContractAddress: {
      payload: undefined,
      signatures: [
        {
          userSignature: {
            userId,
            signature: userSig,
            comment: undefined,
          },
          hashes,
        },
      ],
    },
    blockchain,
    network,
  };

  const rulesContainerDecoder = (_b64: string): DecodedRulesContainer => ({
    users: [
      {
        id: userId,
        name: "User 1",
        publicKeyPem: userPem,
        roles: ["USER"],
      },
    ],
    groups: [
      { id: groupId, name: "Approvers", userIds: [userId] },
    ],
    addressWhitelistingRules: [],
    contractAddressWhitelistingRules: [
      {
        blockchain,
        network,
        parallelThresholds: [
          {
            thresholds: [
              { groupId, minimumSignatures: 1, threshold: 0 },
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

  const userSignaturesDecoder = (_b64: string): RuleUserSignature[] => [
    { userId: "sa@bank.com", signature: saSig },
  ];

  return {
    saPriv,
    saPub,
    userPriv,
    userPub,
    saPem,
    userPem,
    envelope,
    rulesContainerDecoder,
    userSignaturesDecoder,
  };
}

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
