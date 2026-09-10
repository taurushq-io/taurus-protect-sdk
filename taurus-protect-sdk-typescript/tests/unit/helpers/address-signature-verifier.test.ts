/**
 * Unit tests for address-signature-verifier.ts.
 *
 * Tests HSM signature verification for addresses.
 *
 * Both helpers THROW on failure rather than returning booleans, matching the Go,
 * Java and Python peers. The boolean contract they replaced could not distinguish
 * "nothing to verify" from "verification failed", and the batch form returned an
 * array a caller could ignore entirely — silently accepting unverified addresses.
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
  it("passes for a valid HSM signature", () => {
    const { privateKey, publicKey } = generateP256KeyPair();
    const rulesContainer = buildRulesContainer(encodePublicKeyPem(publicKey));

    const address = "0xabc123def456";
    const sig = signData(privateKey, Buffer.from(address, "utf-8"));

    expect(() => verifyAddressSignature(address, sig, rulesContainer)).not.toThrow();
  });

  it("throws when the signature is empty", () => {
    const { publicKey } = generateP256KeyPair();
    const rulesContainer = buildRulesContainer(encodePublicKeyPem(publicKey));

    expect(() => verifyAddressSignature("0xabc", "", rulesContainer)).toThrow(IntegrityError);
    expect(() => verifyAddressSignature("0xabc", "", rulesContainer)).toThrow(/has no signature/);
  });

  // Absent before: the helper never checked the address at all, so an empty
  // address was verified against whatever signature happened to be supplied.
  it("throws when the blockchain address is empty", () => {
    const { privateKey, publicKey } = generateP256KeyPair();
    const rulesContainer = buildRulesContainer(encodePublicKeyPem(publicKey));
    const sig = signData(privateKey, Buffer.from("", "utf-8"));

    expect(() => verifyAddressSignature("", sig, rulesContainer)).toThrow(
      /has no blockchain address/
    );
  });

  it("throws for invalid signature data", () => {
    const { publicKey } = generateP256KeyPair();
    const rulesContainer = buildRulesContainer(encodePublicKeyPem(publicKey));

    expect(() =>
      verifyAddressSignature("0xabc", "bm90LWEtc2lnbmF0dXJl", rulesContainer)
    ).toThrow(IntegrityError);
  });

  it("throws when the signature is for a different address", () => {
    const { privateKey, publicKey } = generateP256KeyPair();
    const rulesContainer = buildRulesContainer(encodePublicKeyPem(publicKey));
    const sig = signData(privateKey, Buffer.from("0xaaa", "utf-8"));

    expect(() => verifyAddressSignature("0xbbb", sig, rulesContainer)).toThrow(
      /verification failed/
    );
  });

  it("throws when no HSM key is present in the rules container", () => {
    const { privateKey } = generateP256KeyPair();
    const sig = signData(privateKey, Buffer.from("0xabc", "utf-8"));

    expect(() =>
      verifyAddressSignature("0xabc", sig, buildEmptyRulesContainer())
    ).toThrow(IntegrityError);
  });

  it("verifies against the UTF-8 encoding of the address", () => {
    const { privateKey, publicKey } = generateP256KeyPair();
    const rulesContainer = buildRulesContainer(encodePublicKeyPem(publicKey));

    const address = "tz1VSUr8wwNhLAzempoch5d6hLRiTh8Cjcjb";
    const sig = signData(privateKey, Buffer.from(address, "utf-8"));

    expect(() => verifyAddressSignature(address, sig, rulesContainer)).not.toThrow();
  });

  it("names the address id in the error when one is supplied", () => {
    const { publicKey } = generateP256KeyPair();
    const rulesContainer = buildRulesContainer(encodePublicKeyPem(publicKey));

    expect(() => verifyAddressSignature("0xabc", "", rulesContainer, 42)).toThrow(
      /Address 42 has no signature/
    );
  });
});

describe("verifyAddressSignatures", () => {
  it("passes when every address verifies", () => {
    const { privateKey, publicKey } = generateP256KeyPair();
    const rulesContainer = buildRulesContainer(encodePublicKeyPem(publicKey));

    const addresses = ["0xaaa", "0xbbb"].map((address) => ({
      address,
      signature: signData(privateKey, Buffer.from(address, "utf-8")),
    }));

    expect(() => verifyAddressSignatures(addresses, rulesContainer)).not.toThrow();
  });

  // The batch form used to return booleans, so a caller who ignored the array
  // accepted every unverified address without noticing.
  it("throws on the first address that does not verify", () => {
    const { privateKey, publicKey } = generateP256KeyPair();
    const rulesContainer = buildRulesContainer(encodePublicKeyPem(publicKey));

    const good = "0xaaa";
    const addresses = [
      { address: good, signature: signData(privateKey, Buffer.from(good, "utf-8")), id: 1 },
      { address: "0xbbb", signature: signData(privateKey, Buffer.from("0xccc", "utf-8")), id: 2 },
    ];

    expect(() => verifyAddressSignatures(addresses, rulesContainer)).toThrow(
      /verification failed for address 2/
    );
  });

  it("throws for an address whose signature is undefined", () => {
    const { publicKey } = generateP256KeyPair();
    const rulesContainer = buildRulesContainer(encodePublicKeyPem(publicKey));

    expect(() =>
      verifyAddressSignatures([{ address: "0xaaa", signature: undefined }], rulesContainer)
    ).toThrow(/has no signature/);
  });

  it("accepts an empty list", () => {
    expect(() => verifyAddressSignatures([], buildEmptyRulesContainer())).not.toThrow();
  });
});
