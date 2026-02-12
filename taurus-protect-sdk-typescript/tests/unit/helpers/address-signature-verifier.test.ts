/**
 * Unit tests for address-signature-verifier.ts.
 *
 * Tests HSM signature verification for addresses:
 * - verifyAddressSignature(): single address verification
 * - verifyAddressSignatures(): batch address verification
 */

import * as crypto from "crypto";

import { signData, encodePublicKeyPem } from "../../../src/crypto";
import { IntegrityError } from "../../../src/errors";
import {
  verifyAddressSignature,
  verifyAddressSignatures,
} from "../../../src/helpers/address-signature-verifier";
import type { DecodedRulesContainer } from "../../../src/models/governance-rules";

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

function buildRulesContainer(hsmPublicKeyPem: string): DecodedRulesContainer {
  return {
    users: [
      {
        id: "hsm-slot-1",
        name: "HSM Slot",
        publicKeyPem: hsmPublicKeyPem,
        roles: ["HSMSLOT"],
      },
    ],
    groups: [],
    minimumDistinctUserSignatures: 0,
    minimumDistinctGroupSignatures: 0,
    transactionRules: [],
    addressWhitelistingRules: [],
    contractAddressWhitelistingRules: [],
    enforcedRulesHash: "",
    timestamp: 0,
    hsmSlotId: 0,
    minimumCommitmentSignatures: 0,
    engineIdentities: [],
  };
}

function buildEmptyRulesContainer(): DecodedRulesContainer {
  return {
    users: [],
    groups: [],
    minimumDistinctUserSignatures: 0,
    minimumDistinctGroupSignatures: 0,
    transactionRules: [],
    addressWhitelistingRules: [],
    contractAddressWhitelistingRules: [],
    enforcedRulesHash: "",
    timestamp: 0,
    hsmSlotId: 0,
    minimumCommitmentSignatures: 0,
    engineIdentities: [],
  };
}

function buildRulesContainerWithNonHsmUser(publicKeyPem: string): DecodedRulesContainer {
  return {
    users: [
      {
        id: "user-1",
        name: "Regular User",
        publicKeyPem,
        roles: ["USER"],
      },
    ],
    groups: [],
    minimumDistinctUserSignatures: 0,
    minimumDistinctGroupSignatures: 0,
    transactionRules: [],
    addressWhitelistingRules: [],
    contractAddressWhitelistingRules: [],
    enforcedRulesHash: "",
    timestamp: 0,
    hsmSlotId: 0,
    minimumCommitmentSignatures: 0,
    engineIdentities: [],
  };
}

// =============================================================================
// verifyAddressSignature
// =============================================================================

describe("verifyAddressSignature", () => {
  it("should return true for a valid HSM signature", () => {
    const { privateKey, publicKey } = generateP256KeyPair();
    const pem = encodePublicKeyPem(publicKey);
    const rulesContainer = buildRulesContainer(pem);

    const address = "0xabc123def456";
    const sig = signData(privateKey, Buffer.from(address, "utf-8"));

    expect(verifyAddressSignature(address, sig, rulesContainer)).toBe(true);
  });

  it("should return false for empty signature", () => {
    const { publicKey } = generateP256KeyPair();
    const pem = encodePublicKeyPem(publicKey);
    const rulesContainer = buildRulesContainer(pem);

    expect(verifyAddressSignature("0xabc", "", rulesContainer)).toBe(false);
  });

  it("should throw IntegrityError when HSM public key is not found", () => {
    const rulesContainer = buildEmptyRulesContainer();

    expect(() => {
      verifyAddressSignature("0xabc", "some-sig", rulesContainer);
    }).toThrow(IntegrityError);
    expect(() => {
      verifyAddressSignature("0xabc", "some-sig", rulesContainer);
    }).toThrow("HSM public key not found");
  });

  it("should throw IntegrityError when no user has HSMSLOT role", () => {
    const { publicKey } = generateP256KeyPair();
    const pem = encodePublicKeyPem(publicKey);
    const rulesContainer = buildRulesContainerWithNonHsmUser(pem);

    expect(() => {
      verifyAddressSignature("0xabc", "some-sig", rulesContainer);
    }).toThrow(IntegrityError);
    expect(() => {
      verifyAddressSignature("0xabc", "some-sig", rulesContainer);
    }).toThrow("HSM public key not found");
  });

  it("should return false for invalid signature data", () => {
    const { publicKey } = generateP256KeyPair();
    const pem = encodePublicKeyPem(publicKey);
    const rulesContainer = buildRulesContainer(pem);

    // A 64-byte random signature that won't verify
    const fakeSig = Buffer.alloc(64, 0x42).toString("base64");

    expect(verifyAddressSignature("0xabc", fakeSig, rulesContainer)).toBe(false);
  });

  it("should return false when signature is for a different address", () => {
    const { privateKey, publicKey } = generateP256KeyPair();
    const pem = encodePublicKeyPem(publicKey);
    const rulesContainer = buildRulesContainer(pem);

    const sig = signData(privateKey, Buffer.from("original-address", "utf-8"));

    expect(verifyAddressSignature("different-address", sig, rulesContainer)).toBe(false);
  });

  it("should verify signature against correct address encoding (UTF-8)", () => {
    const { privateKey, publicKey } = generateP256KeyPair();
    const pem = encodePublicKeyPem(publicKey);
    const rulesContainer = buildRulesContainer(pem);

    // Unicode address for edge case
    const address = "addr_with_unicode_\u00e9";
    const sig = signData(privateKey, Buffer.from(address, "utf-8"));

    expect(verifyAddressSignature(address, sig, rulesContainer)).toBe(true);
  });
});

// =============================================================================
// verifyAddressSignatures (batch)
// =============================================================================

describe("verifyAddressSignatures", () => {
  it("should verify multiple addresses", () => {
    const { privateKey, publicKey } = generateP256KeyPair();
    const pem = encodePublicKeyPem(publicKey);
    const rulesContainer = buildRulesContainer(pem);

    const addr1 = "0xaddr1";
    const addr2 = "0xaddr2";
    const sig1 = signData(privateKey, Buffer.from(addr1, "utf-8"));
    const sig2 = signData(privateKey, Buffer.from(addr2, "utf-8"));

    const results = verifyAddressSignatures(
      [
        { address: addr1, signature: sig1 },
        { address: addr2, signature: sig2 },
      ],
      rulesContainer
    );

    expect(results).toEqual([true, true]);
  });

  it("should return false for addresses with undefined signature", () => {
    const { publicKey } = generateP256KeyPair();
    const pem = encodePublicKeyPem(publicKey);
    const rulesContainer = buildRulesContainer(pem);

    const results = verifyAddressSignatures(
      [{ address: "0xabc", signature: undefined }],
      rulesContainer
    );

    expect(results).toEqual([false]);
  });

  it("should return mixed results for mixed valid/invalid signatures", () => {
    const { privateKey, publicKey } = generateP256KeyPair();
    const pem = encodePublicKeyPem(publicKey);
    const rulesContainer = buildRulesContainer(pem);

    const goodAddr = "0xgood";
    const goodSig = signData(privateKey, Buffer.from(goodAddr, "utf-8"));
    const badSig = Buffer.alloc(64, 0x42).toString("base64");

    const results = verifyAddressSignatures(
      [
        { address: goodAddr, signature: goodSig },
        { address: "0xbad", signature: badSig },
      ],
      rulesContainer
    );

    expect(results).toEqual([true, false]);
  });

  it("should return empty array for empty input", () => {
    const { publicKey } = generateP256KeyPair();
    const pem = encodePublicKeyPem(publicKey);
    const rulesContainer = buildRulesContainer(pem);

    const results = verifyAddressSignatures([], rulesContainer);

    expect(results).toEqual([]);
  });

  it("should return all false when HSM key is missing (catches IntegrityError)", () => {
    const rulesContainer = buildEmptyRulesContainer();

    const results = verifyAddressSignatures(
      [
        { address: "0xabc", signature: "some-sig" },
        { address: "0xdef", signature: "other-sig" },
      ],
      rulesContainer
    );

    expect(results).toEqual([false, false]);
  });
});
