/**
 * Unit tests for WhitelistedAddressVerifier.
 *
 * Tests the complete 6-step verification flow for whitelisted addresses:
 * 1. Verify metadata hash
 * 2. Verify rules container signatures (SuperAdmin keys)
 * 3. Decode rules container
 * 4. Verify hash coverage (with legacy hash support)
 * 5. Verify whitelist signatures meet governance thresholds
 * 6. Parse WhitelistedAddress from verified payload
 */

import * as crypto from "crypto";

import { calculateHexHash, signData, encodePublicKeyPem } from "../../../src/crypto";
import { IntegrityError, WhitelistError } from "../../../src/errors";
import { WhitelistedAddressVerifier } from "../../../src/helpers/whitelisted-address-verifier";
import type {
  DecodedRulesContainer,
  RuleUserSignature,
} from "../../../src/models/governance-rules";
import type {
  SignedWhitelistedAddressEnvelope,
} from "../../../src/models/whitelisted-address";

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

function buildAddressPayload(overrides?: Partial<{
  currency: string;
  address: string;
  label: string;
  contractType: string;
  linkedInternalAddresses: unknown[];
  linkedWallets: unknown[];
}>): Record<string, unknown> {
  return {
    currency: "ETH",
    addressType: "individual",
    address: "0xABCD1234",
    memo: "",
    label: "Test Addr",
    customerId: "",
    exchangeAccountId: "",
    linkedInternalAddresses: [],
    contractType: "",
    ...overrides,
  };
}

function payloadToString(payload: Record<string, unknown>): string {
  return JSON.stringify(payload);
}

interface AddressTestFixture {
  saPriv: crypto.KeyObject;
  saPub: crypto.KeyObject;
  userPriv: crypto.KeyObject;
  userPub: crypto.KeyObject;
  saPem: string;
  userPem: string;
  envelope: SignedWhitelistedAddressEnvelope;
  rulesContainerDecoder: (b64: string) => DecodedRulesContainer;
  userSignaturesDecoder: (b64: string) => RuleUserSignature[];
}

function buildFullAddressFixture(overrides?: {
  blockchain?: string;
  network?: string;
  groupId?: string;
  userId?: string;
  payload?: Record<string, unknown>;
  linkedWallets?: Array<{ id: number; path: string; label: string | undefined }>;
}): AddressTestFixture {
  const { privateKey: saPriv, publicKey: saPub } = generateP256KeyPair();
  const { privateKey: userPriv, publicKey: userPub } = generateP256KeyPair();
  const saPem = keyToPem(saPub);
  const userPem = keyToPem(userPub);

  const blockchain = overrides?.blockchain ?? "ETH";
  const network = overrides?.network ?? "mainnet";
  const groupId = overrides?.groupId ?? "approvers";
  const userId = overrides?.userId ?? "user1@bank.com";

  const payload = overrides?.payload ?? buildAddressPayload();
  const payloadStr = payloadToString(payload);
  const metadataHash = calculateHexHash(payloadStr);

  // Build rules container
  const rulesJson = JSON.stringify({
    users: [{ id: userId, publicKey: userPem, roles: ["USER"] }],
    groups: [{ id: groupId, userIds: [userId] }],
    addressWhitelistingRules: [
      {
        currency: blockchain,
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

  const envelope: SignedWhitelistedAddressEnvelope = {
    id: "addr-1",
    metadata: {
      hash: metadataHash,
      payloadAsString: payloadStr,
    },
    rulesContainerBase64: rulesB64,
    rulesSignaturesBase64: Buffer.from("dummy").toString("base64"),
    signedAddress: {
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
    linkedInternalAddresses: [],
    linkedWallets: overrides?.linkedWallets ?? [],
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
    addressWhitelistingRules: [
      {
        currency: blockchain,
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
    contractAddressWhitelistingRules: [],
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

describe("WhitelistedAddressVerifier constructor", () => {
  it("should create verifier with valid config", () => {
    const { saPem } = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [saPem],
      minValidSignatures: 1,
    });
    expect(verifier).toBeDefined();
  });

  it("should throw if no SuperAdmin keys", () => {
    expect(() => {
      new WhitelistedAddressVerifier({
        superAdminKeysPem: [],
        minValidSignatures: 1,
      });
    }).toThrow("At least one SuperAdmin key is required");
  });

  it("should throw if minValidSignatures < 1", () => {
    const { saPem } = buildFullAddressFixture();
    expect(() => {
      new WhitelistedAddressVerifier({
        superAdminKeysPem: [saPem],
        minValidSignatures: 0,
      });
    }).toThrow("minValidSignatures must be at least 1");
  });

  it("should throw for invalid PEM key", () => {
    expect(() => {
      new WhitelistedAddressVerifier({
        superAdminKeysPem: ["not a valid PEM key"],
        minValidSignatures: 1,
      });
    }).toThrow();
  });
});

// =============================================================================
// Step 1 Tests: Metadata Hash
// =============================================================================

describe("WhitelistedAddressVerifier - Step 1: Metadata Hash", () => {
  it("should pass with valid hash", () => {
    const f = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
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
    const f = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const tampered: SignedWhitelistedAddressEnvelope = {
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

  it("should throw for empty payload", () => {
    const f = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const tampered: SignedWhitelistedAddressEnvelope = {
      ...f.envelope,
      metadata: { hash: "abc", payloadAsString: "" },
    };

    expect(() => {
      verifier.verify(tampered, f.rulesContainerDecoder, f.userSignaturesDecoder);
    }).toThrow("payloadAsString is empty");
  });

  it("should throw for empty hash", () => {
    const f = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const tampered: SignedWhitelistedAddressEnvelope = {
      ...f.envelope,
      metadata: { hash: "", payloadAsString: '{"foo":"bar"}' },
    };

    expect(() => {
      verifier.verify(tampered, f.rulesContainerDecoder, f.userSignaturesDecoder);
    }).toThrow("metadata hash is empty");
  });

  it("should throw for null envelope", () => {
    const f = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    expect(() => {
      verifier.verify(null as any, f.rulesContainerDecoder, f.userSignaturesDecoder);
    }).toThrow(IntegrityError);
  });

  it("should throw for null metadata", () => {
    const f = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
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

describe("WhitelistedAddressVerifier - Step 2: Rules Container Signatures", () => {
  it("should throw for empty rulesContainer", () => {
    const f = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const tampered: SignedWhitelistedAddressEnvelope = {
      ...f.envelope,
      rulesContainerBase64: "",
    };

    expect(() => {
      verifier.verify(tampered, f.rulesContainerDecoder, f.userSignaturesDecoder);
    }).toThrow("rulesContainer is empty");
  });

  it("should throw for empty rulesSignatures", () => {
    const f = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const tampered: SignedWhitelistedAddressEnvelope = {
      ...f.envelope,
      rulesSignaturesBase64: "",
    };

    expect(() => {
      verifier.verify(tampered, f.rulesContainerDecoder, f.userSignaturesDecoder);
    }).toThrow("rulesSignatures is empty");
  });

  it("should throw when insufficient valid signatures", () => {
    const f = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 2,
    });

    expect(() => {
      verifier.verify(f.envelope, f.rulesContainerDecoder, f.userSignaturesDecoder);
    }).toThrow("rules container signature verification failed");
  });

  it("should throw when signature decoder fails", () => {
    const f = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
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

describe("WhitelistedAddressVerifier - Step 3: Decode Rules Container", () => {
  it("should throw when decoder fails", () => {
    const f = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
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

describe("WhitelistedAddressVerifier - Step 4: Hash Coverage", () => {
  it("should throw for null signedAddress", () => {
    const f = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const tampered = {
      ...f.envelope,
      signedAddress: null as any,
    };

    expect(() => {
      verifier.verify(tampered, f.rulesContainerDecoder, f.userSignaturesDecoder);
    }).toThrow("signedAddress");
  });

  it("should throw for empty signatures in signedAddress", () => {
    const f = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const tampered: SignedWhitelistedAddressEnvelope = {
      ...f.envelope,
      signedAddress: {
        payload: undefined,
        signatures: [],
      },
    };

    expect(() => {
      verifier.verify(tampered, f.rulesContainerDecoder, f.userSignaturesDecoder);
    }).toThrow("no signatures in signedAddress");
  });

  it("should throw when hash not covered by any signature", () => {
    const f = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const tampered: SignedWhitelistedAddressEnvelope = {
      ...f.envelope,
      signedAddress: {
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
    }).toThrow("is not covered by any signature");
  });

  it("should pass with legacy hash fallback (contractType removed)", () => {
    const f = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    // Build a payload WITH contractType (non-empty so removal is detectable)
    const payloadWithContract = buildAddressPayload({ contractType: "ERC20" });
    const payloadStr = payloadToString(payloadWithContract);
    const currentHash = calculateHexHash(payloadStr);

    // Compute the legacy hash (without contractType field)
    const payloadObj = JSON.parse(payloadStr);
    delete payloadObj.contractType;
    const legacyPayloadStr = JSON.stringify(payloadObj);
    const legacyHash = calculateHexHash(legacyPayloadStr);

    // Sign hashes array with user key using the LEGACY hash
    const hashes = [legacyHash];
    const hashesJson = JSON.stringify(hashes);
    const userSig = signData(f.userPriv, Buffer.from(hashesJson, "utf-8"));

    const legacyEnvelope: SignedWhitelistedAddressEnvelope = {
      ...f.envelope,
      metadata: {
        hash: currentHash,
        payloadAsString: payloadStr,
      },
      signedAddress: {
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
      legacyEnvelope,
      f.rulesContainerDecoder,
      f.userSignaturesDecoder
    );
    expect(result).toBeDefined();
    // The verified hash should be the legacy hash (the one found in signatures)
    expect(result.verifiedHash).toBe(legacyHash);
    expect(result.verifiedHash).not.toBe(currentHash);
  });
});

// =============================================================================
// Step 5 Tests: Whitelist Signatures
// =============================================================================

describe("WhitelistedAddressVerifier - Step 5: Whitelist Signatures", () => {
  it("should throw when no rules for blockchain", () => {
    const f = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const btcDecoder = (_b64: string): DecodedRulesContainer => ({
      users: [],
      groups: [],
      addressWhitelistingRules: [
        {
          currency: "BTC",
          network: "mainnet",
          parallelThresholds: [
            { thresholds: [{ groupId: "approvers", minimumSignatures: 1, threshold: 0 }] },
          ],
        },
      ],
      contractAddressWhitelistingRules: [],
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
    const f = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const emptyThresholdsDecoder = (_b64: string): DecodedRulesContainer => ({
      users: [],
      groups: [],
      addressWhitelistingRules: [
        {
          currency: "ETH",
          network: "mainnet",
          parallelThresholds: [],
        },
      ],
      contractAddressWhitelistingRules: [],
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
    const f = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
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
      addressWhitelistingRules: [
        {
          currency: "ETH",
          network: "mainnet",
          parallelThresholds: [
            { thresholds: [{ groupId: "approvers", minimumSignatures: 2, threshold: 0 }] },
          ],
        },
      ],
      contractAddressWhitelistingRules: [],
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

describe("WhitelistedAddressVerifier - End-to-End", () => {
  it("should succeed with valid data", () => {
    const f = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const result = verifier.verify(
      f.envelope,
      f.rulesContainerDecoder,
      f.userSignaturesDecoder
    );

    expect(result).toBeDefined();
    expect(result.verifiedWhitelistedAddress).toBeDefined();
    expect(result.verifiedHash).toBe(f.envelope.metadata.hash);
    expect(result.verifiedRulesContainer).toBeDefined();
  });

  it("should return correct address from verified payload", () => {
    const f = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const result = verifier.verify(
      f.envelope,
      f.rulesContainerDecoder,
      f.userSignaturesDecoder
    );

    expect(result.verifiedWhitelistedAddress.address).toBe("0xABCD1234");
    expect(result.verifiedWhitelistedAddress.id).toBe("addr-1");
  });

  it("should succeed with two users and sequential thresholds (AND logic)", () => {
    const { privateKey: saPriv, publicKey: saPub } = generateP256KeyPair();
    const { privateKey: u1Priv, publicKey: u1Pub } = generateP256KeyPair();
    const { privateKey: u2Priv, publicKey: u2Pub } = generateP256KeyPair();
    const saPem = keyToPem(saPub);
    const u1Pem = keyToPem(u1Pub);
    const u2Pem = keyToPem(u2Pub);

    const payload = buildAddressPayload();
    const payloadStr = payloadToString(payload);
    const metadataHash = calculateHexHash(payloadStr);

    const hashes = [metadataHash];
    const hashesJson = JSON.stringify(hashes);
    const u1Sig = signData(u1Priv, Buffer.from(hashesJson));
    const u2Sig = signData(u2Priv, Buffer.from(hashesJson));

    const rulesB64 = Buffer.from("{}").toString("base64");
    const saSig = signData(saPriv, Buffer.from(rulesB64, "base64"));

    const envelope: SignedWhitelistedAddressEnvelope = {
      id: "addr-2",
      metadata: { hash: metadataHash, payloadAsString: payloadStr },
      rulesContainerBase64: rulesB64,
      rulesSignaturesBase64: Buffer.from("dummy").toString("base64"),
      signedAddress: {
        payload: undefined,
        signatures: [
          { userSignature: { userId: "user1@bank.com", signature: u1Sig, comment: undefined }, hashes },
          { userSignature: { userId: "user2@bank.com", signature: u2Sig, comment: undefined }, hashes },
        ],
      },
      blockchain: "ETH",
      network: "mainnet",
      linkedInternalAddresses: [],
      linkedWallets: [],
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
      addressWhitelistingRules: [
        {
          currency: "ETH",
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
      contractAddressWhitelistingRules: [],
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

    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [saPem],
      minValidSignatures: 1,
    });

    const result = verifier.verify(envelope, rcDecoder, usDecoder);
    expect(result).toBeDefined();
    expect(result.verifiedHash).toBe(metadataHash);
  });

  it("should succeed with parallel paths (OR logic)", () => {
    const f = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const twoPathDecoder = (_b64: string): DecodedRulesContainer => ({
      users: [
        { id: "user1@bank.com", name: "User 1", publicKeyPem: f.userPem, roles: ["USER"] },
      ],
      groups: [
        { id: "approvers", name: "Approvers", userIds: ["user1@bank.com"] },
        { id: "other_team", name: "Other Team", userIds: ["nobody@bank.com"] },
      ],
      addressWhitelistingRules: [
        {
          currency: "ETH",
          network: "mainnet",
          parallelThresholds: [
            // Path 1: will fail
            { thresholds: [{ groupId: "other_team", minimumSignatures: 1, threshold: 0 }] },
            // Path 2: will pass
            { thresholds: [{ groupId: "approvers", minimumSignatures: 1, threshold: 0 }] },
          ],
        },
      ],
      contractAddressWhitelistingRules: [],
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
    const f = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
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
      transactionRules: [],
      addressWhitelistingRules: [
        {
          currency: "ETH",
          network: "mainnet",
          parallelThresholds: [
            { thresholds: [{ groupId: "team_a", minimumSignatures: 1, threshold: 0 }] },
            { thresholds: [{ groupId: "team_b", minimumSignatures: 1, threshold: 0 }] },
          ],
        },
      ],
      contractAddressWhitelistingRules: [],
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
// Rule Lines Tests
// =============================================================================

describe("WhitelistedAddressVerifier - Rule Lines", () => {
  it("should use rule line thresholds when no linked addresses and exactly 1 linked wallet", () => {
    const f = buildFullAddressFixture();
    const { privateKey: u2Priv, publicKey: u2Pub } = generateP256KeyPair();
    const u2Pem = keyToPem(u2Pub);

    const metadataHash = f.envelope.metadata.hash;
    const hashes = [metadataHash];
    const hashesJson = JSON.stringify(hashes);
    const u2Sig = signData(u2Priv, Buffer.from(hashesJson, "utf-8"));

    // Envelope has 1 wallet, no linked addresses
    const envelopeWithWallet: SignedWhitelistedAddressEnvelope = {
      ...f.envelope,
      linkedInternalAddresses: [],
      linkedWallets: [{ id: 1, path: "ETH/wallet1", label: "My Wallet" }],
      signedAddress: {
        payload: undefined,
        signatures: [
          {
            userSignature: {
              userId: "user2@bank.com",
              signature: u2Sig,
              comment: undefined,
            },
            hashes,
          },
        ],
      },
    };

    // Rules with a line that matches the wallet path
    const ruleLineDecoder = (_b64: string): DecodedRulesContainer => ({
      users: [
        { id: "user2@bank.com", name: "User 2", publicKeyPem: u2Pem, roles: ["USER"] },
      ],
      groups: [
        { id: "wallet_approvers", name: "Wallet Approvers", userIds: ["user2@bank.com"] },
      ],
      transactionRules: [],
      addressWhitelistingRules: [
        {
          currency: "ETH",
          network: "mainnet",
          parallelThresholds: [
            // Default thresholds require a group user2 is NOT in
            { thresholds: [{ groupId: "default_group", minimumSignatures: 1, threshold: 0 }] },
          ],
          // Lines with wallet-specific thresholds
          lines: [
            {
              cells: [{ type: "INTERNAL_WALLET", internalWallet: { path: "ETH/wallet1" } }],
              parallelThresholds: [
                { thresholds: [{ groupId: "wallet_approvers", minimumSignatures: 1, threshold: 0 }] },
              ],
            },
          ],
        } as any,
      ],
      contractAddressWhitelistingRules: [],
      minimumDistinctUserSignatures: 0,
      minimumDistinctGroupSignatures: 0,
      enforcedRulesHash: "",
      timestamp: 0,
      hsmSlotId: 0,
      minimumCommitmentSignatures: 0,
      engineIdentities: [],
    });

    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const result = verifier.verify(
      envelopeWithWallet,
      ruleLineDecoder,
      f.userSignaturesDecoder
    );
    expect(result).toBeDefined();
    expect(result.verifiedWhitelistedAddress.blockchain).toBe("ETH");
  });

  it("should use default thresholds when multiple wallets are linked", () => {
    const f = buildFullAddressFixture();

    // Envelope with 2 wallets - should fall back to default thresholds
    const envelopeWithTwoWallets: SignedWhitelistedAddressEnvelope = {
      ...f.envelope,
      linkedInternalAddresses: [],
      linkedWallets: [
        { id: 1, path: "ETH/wallet1", label: "Wallet 1" },
        { id: 2, path: "ETH/wallet2", label: "Wallet 2" },
      ],
    };

    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const result = verifier.verify(
      envelopeWithTwoWallets,
      f.rulesContainerDecoder,
      f.userSignaturesDecoder
    );
    expect(result).toBeDefined();
  });
});

// =============================================================================
// Batch Verification
// =============================================================================

describe("WhitelistedAddressVerifier - Batch Verification", () => {
  it("should verify multiple envelopes independently", () => {
    const f1 = buildFullAddressFixture({ blockchain: "ETH", payload: buildAddressPayload({ currency: "ETH" }) });
    const f2 = buildFullAddressFixture({ blockchain: "BTC", payload: buildAddressPayload({ currency: "BTC" }) });

    const verifier1 = new WhitelistedAddressVerifier({
      superAdminKeysPem: [f1.saPem],
      minValidSignatures: 1,
    });
    const verifier2 = new WhitelistedAddressVerifier({
      superAdminKeysPem: [f2.saPem],
      minValidSignatures: 1,
    });

    const result1 = verifier1.verify(f1.envelope, f1.rulesContainerDecoder, f1.userSignaturesDecoder);
    const result2 = verifier2.verify(f2.envelope, f2.rulesContainerDecoder, f2.userSignaturesDecoder);

    expect(result1.verifiedWhitelistedAddress.blockchain).toBe("ETH");
    expect(result2.verifiedWhitelistedAddress.blockchain).toBe("BTC");
  });

  it("should fail on first invalid envelope in strict mode", () => {
    const f = buildFullAddressFixture();
    const verifier = new WhitelistedAddressVerifier({
      superAdminKeysPem: [f.saPem],
      minValidSignatures: 1,
    });

    const validEnvelope = f.envelope;
    const invalidEnvelope: SignedWhitelistedAddressEnvelope = {
      ...f.envelope,
      metadata: {
        hash: "tampered_hash",
        payloadAsString: f.envelope.metadata.payloadAsString,
      },
    };

    // Simulate strict batch: verify each, expect failure on invalid
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
