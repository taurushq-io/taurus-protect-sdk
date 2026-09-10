/**
 * Unit tests for signature-verifier.ts.
 *
 * Tests the SuperAdmin signature verification functions:
 * - isValidSignature(): single signature verification against multiple keys
 * - verifyGovernanceRules(): multi-signature threshold verification
 */

import * as crypto from "crypto";

import { signData, encodePublicKeyPem } from "../../../src/crypto";
import { ConfigurationError, IntegrityError } from "../../../src/errors";
import {
  isValidSignature,
  verifyGovernanceRules,
  verifyGovernanceRulesSignatures,
  verifySignatureWithKey,
} from "../../../src/helpers/signature-verifier";

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

function buildSignaturesBase64(
  signatures: Array<{ userId?: string; signature?: string }>
): string {
  return Buffer.from(JSON.stringify(signatures)).toString("base64");
}

function buildWrappedSignaturesBase64(
  signatures: Array<{ userId?: string; signature?: string }>
): string {
  return Buffer.from(JSON.stringify({ signatures })).toString("base64");
}

// =============================================================================
// distinct-signer counting
// =============================================================================

// The threshold counts distinct signing keys. Counting entries — or deduping by the
// caller-supplied userId — would let one compromised key satisfy any threshold, since
// ECDSA is randomized and a single key can emit unlimited valid signatures.
describe("verifyGovernanceRules distinct-signer counting", () => {
  const rulesData = Buffer.from("governance rules payload");
  const rulesBase64 = rulesData.toString("base64");

  const key1 = generateP256KeyPair();
  const key2 = generateP256KeyPair();
  const unconfigured = generateP256KeyPair();

  const bothPems = [keyToPem(key1.publicKey), keyToPem(key2.publicKey)];
  const key1TwicePems = [keyToPem(key1.publicKey), keyToPem(key1.publicKey)];

  it("accepts two distinct signers at threshold two", () => {
    const sigs = buildSignaturesBase64([
      { userId: "admin1", signature: signData(key1.privateKey, rulesData) },
      { userId: "admin2", signature: signData(key2.privateKey, rulesData) },
    ]);
    expect(verifyGovernanceRules(rulesBase64, sigs, 2, bothPems)).toBe(true);
  });

  it("rejects two signatures from one key at threshold two, accepts at one", () => {
    const sigs = buildSignaturesBase64([
      { userId: "admin1", signature: signData(key1.privateKey, rulesData) },
      { userId: "admin1-again", signature: signData(key1.privateKey, rulesData) },
    ]);
    expect(verifyGovernanceRules(rulesBase64, sigs, 2, bothPems)).toBe(false);
    expect(verifyGovernanceRules(rulesBase64, sigs, 1, bothPems)).toBe(true);
  });

  it("counts a replayed identical signature once", () => {
    const replayed = signData(key1.privateKey, rulesData);
    const sigs = buildSignaturesBase64([
      { userId: "admin1", signature: replayed },
      { userId: "admin1-replay", signature: replayed },
    ]);
    expect(verifyGovernanceRules(rulesBase64, sigs, 2, bothPems)).toBe(false);
  });

  it("counts the same key configured twice once", () => {
    const sigs = buildSignaturesBase64([
      { userId: "admin1", signature: signData(key1.privateKey, rulesData) },
      { userId: "admin1-again", signature: signData(key1.privateKey, rulesData) },
    ]);
    expect(verifyGovernanceRules(rulesBase64, sigs, 2, key1TwicePems)).toBe(false);
  });

  it("ignores a signature from an unconfigured key", () => {
    const sigs = buildSignaturesBase64([
      { userId: "admin1", signature: signData(key1.privateKey, rulesData) },
      { userId: "stranger", signature: signData(unconfigured.privateKey, rulesData) },
    ]);
    expect(verifyGovernanceRules(rulesBase64, sigs, 2, bothPems)).toBe(false);
  });
});

// =============================================================================
// isValidSignature
// =============================================================================

describe("isValidSignature", () => {
  it("should return true for a valid signature", () => {
    const { privateKey, publicKey } = generateP256KeyPair();
    const data = Buffer.from("test data");
    const sig = signData(privateKey, data);

    expect(isValidSignature(data, sig, [publicKey])).toBe(true);
  });

  it("should return true when string data is provided", () => {
    const { privateKey, publicKey } = generateP256KeyPair();
    const data = "test string data";
    const sig = signData(privateKey, Buffer.from(data, "utf-8"));

    expect(isValidSignature(data, sig, [publicKey])).toBe(true);
  });

  it("should return false for an invalid signature", () => {
    const { publicKey } = generateP256KeyPair();
    const data = Buffer.from("test data");
    // Use a random base64 string that is 64 bytes when decoded (raw r||s size)
    const fakeSig = Buffer.alloc(64, 0x42).toString("base64");

    expect(isValidSignature(data, fakeSig, [publicKey])).toBe(false);
  });

  it("should return false for empty public keys array", () => {
    const data = Buffer.from("test data");
    const fakeSig = Buffer.alloc(64, 0x42).toString("base64");

    expect(isValidSignature(data, fakeSig, [])).toBe(false);
  });

  it("should return true if any key matches (second key)", () => {
    const pair1 = generateP256KeyPair();
    const pair2 = generateP256KeyPair();
    const data = Buffer.from("multi-key data");
    const sig = signData(pair2.privateKey, data);

    expect(isValidSignature(data, sig, [pair1.publicKey, pair2.publicKey])).toBe(true);
  });

  it("should return false when signature is for different data", () => {
    const { privateKey, publicKey } = generateP256KeyPair();
    const sig = signData(privateKey, Buffer.from("original data"));

    expect(isValidSignature(Buffer.from("tampered data"), sig, [publicKey])).toBe(false);
  });

  it("should return false for malformed base64 signature", () => {
    const { publicKey } = generateP256KeyPair();
    const data = Buffer.from("test data");

    // A short base64 string decoding to fewer than 64 bytes -> returns false
    expect(isValidSignature(data, "not-valid-base64!", [publicKey])).toBe(false);
  });
});

// =============================================================================
// verifyGovernanceRules - threshold logic
// =============================================================================

describe("verifyGovernanceRules", () => {
  // Fail closed: a non-positive threshold is a misconfiguration, not "nothing to check".
  it("should throw when minValidSignatures is 0", () => {
    expect(() => verifyGovernanceRules("abc", "abc", 0, [])).toThrow(
      ConfigurationError
    );
  });

  it("should throw when minValidSignatures is negative", () => {
    expect(() => verifyGovernanceRules("abc", "abc", -1, [])).toThrow(
      ConfigurationError
    );
  });

  it("should return false for empty rulesContainerBase64", () => {
    const { publicKey } = generateP256KeyPair();
    const pem = keyToPem(publicKey);

    expect(verifyGovernanceRules("", "sig", 1, [pem])).toBe(false);
  });

  it("should return false for empty signaturesBase64", () => {
    const { publicKey } = generateP256KeyPair();
    const pem = keyToPem(publicKey);

    expect(verifyGovernanceRules("rules", "", 1, [pem])).toBe(false);
  });

  it("should return false for empty superAdminKeysPem", () => {
    expect(verifyGovernanceRules("rules", "sigs", 1, [])).toBe(false);
  });

  it("should return false for invalid PEM keys", () => {
    expect(verifyGovernanceRules("rules", "sigs", 1, ["not-a-pem"])).toBe(false);
  });

  it("should return true with one valid signature meeting threshold of 1", () => {
    const { privateKey, publicKey } = generateP256KeyPair();
    const pem = keyToPem(publicKey);

    const rulesB64 = Buffer.from("some-rules-data").toString("base64");
    const rulesData = Buffer.from(rulesB64, "base64");
    const sig = signData(privateKey, rulesData);

    const signaturesB64 = buildSignaturesBase64([
      { userId: "user1", signature: sig },
    ]);

    expect(verifyGovernanceRules(rulesB64, signaturesB64, 1, [pem])).toBe(true);
  });

  it("should return true with wrapped signatures format", () => {
    const { privateKey, publicKey } = generateP256KeyPair();
    const pem = keyToPem(publicKey);

    const rulesB64 = Buffer.from("wrapped-rules-data").toString("base64");
    const rulesData = Buffer.from(rulesB64, "base64");
    const sig = signData(privateKey, rulesData);

    const signaturesB64 = buildWrappedSignaturesBase64([
      { userId: "user1", signature: sig },
    ]);

    expect(verifyGovernanceRules(rulesB64, signaturesB64, 1, [pem])).toBe(true);
  });

  it("should return false when threshold is not met", () => {
    const { privateKey, publicKey } = generateP256KeyPair();
    const pem = keyToPem(publicKey);

    const rulesB64 = Buffer.from("rules-threshold-test").toString("base64");
    const rulesData = Buffer.from(rulesB64, "base64");
    const sig = signData(privateKey, rulesData);

    const signaturesB64 = buildSignaturesBase64([
      { userId: "user1", signature: sig },
    ]);

    // Require 2 but only 1 valid signature
    expect(verifyGovernanceRules(rulesB64, signaturesB64, 2, [pem])).toBe(false);
  });

  it("should count multiple distinct user signatures", () => {
    const pair1 = generateP256KeyPair();
    const pair2 = generateP256KeyPair();

    const rulesB64 = Buffer.from("multi-user-rules").toString("base64");
    const rulesData = Buffer.from(rulesB64, "base64");

    const sig1 = signData(pair1.privateKey, rulesData);
    const sig2 = signData(pair2.privateKey, rulesData);

    const signaturesB64 = buildSignaturesBase64([
      { userId: "user1", signature: sig1 },
      { userId: "user2", signature: sig2 },
    ]);

    expect(
      verifyGovernanceRules(rulesB64, signaturesB64, 2, [
        keyToPem(pair1.publicKey),
        keyToPem(pair2.publicKey),
      ])
    ).toBe(true);
  });

  it("should deduplicate signatures from the same userId", () => {
    const { privateKey, publicKey } = generateP256KeyPair();
    const pem = keyToPem(publicKey);

    const rulesB64 = Buffer.from("dedup-rules").toString("base64");
    const rulesData = Buffer.from(rulesB64, "base64");
    const sig = signData(privateKey, rulesData);

    // Same userId appears twice
    const signaturesB64 = buildSignaturesBase64([
      { userId: "user1", signature: sig },
      { userId: "user1", signature: sig },
    ]);

    // Should only count 1 unique user
    expect(verifyGovernanceRules(rulesB64, signaturesB64, 2, [pem])).toBe(false);
    expect(verifyGovernanceRules(rulesB64, signaturesB64, 1, [pem])).toBe(true);
  });

  it("should skip signatures with empty signature field", () => {
    const { privateKey, publicKey } = generateP256KeyPair();
    const pem = keyToPem(publicKey);

    const rulesB64 = Buffer.from("empty-sig-rules").toString("base64");
    const rulesData = Buffer.from(rulesB64, "base64");
    const validSig = signData(privateKey, rulesData);

    const signaturesB64 = buildSignaturesBase64([
      { userId: "user1", signature: "" },
      { userId: "user2", signature: validSig },
    ]);

    // Only user2 has a valid signature
    expect(verifyGovernanceRules(rulesB64, signaturesB64, 1, [pem])).toBe(true);
    expect(verifyGovernanceRules(rulesB64, signaturesB64, 2, [pem])).toBe(false);
  });

  it("should return false for invalid JSON in signaturesBase64", () => {
    const { publicKey } = generateP256KeyPair();
    const pem = keyToPem(publicKey);

    const rulesB64 = Buffer.from("json-test").toString("base64");
    const badJsonB64 = Buffer.from("not-json{{{").toString("base64");

    expect(verifyGovernanceRules(rulesB64, badJsonB64, 1, [pem])).toBe(false);
  });

  it("should return false for non-array/non-object JSON", () => {
    const { publicKey } = generateP256KeyPair();
    const pem = keyToPem(publicKey);

    const rulesB64 = Buffer.from("scalar-json-test").toString("base64");
    const scalarB64 = Buffer.from('"just a string"').toString("base64");

    expect(verifyGovernanceRules(rulesB64, scalarB64, 1, [pem])).toBe(false);
  });

  it("should return false for object with non-array signatures property", () => {
    const { publicKey } = generateP256KeyPair();
    const pem = keyToPem(publicKey);

    const rulesB64 = Buffer.from("bad-structure").toString("base64");
    const badStructB64 = Buffer.from(JSON.stringify({ signatures: "not-an-array" })).toString("base64");

    expect(verifyGovernanceRules(rulesB64, badStructB64, 1, [pem])).toBe(false);
  });

  it("should not let one key satisfy a threshold of two, with or without userId", () => {
    const { privateKey, publicKey } = generateP256KeyPair();
    const pem = keyToPem(publicKey);

    const rulesB64 = Buffer.from("no-userid-rules").toString("base64");
    const rulesData = Buffer.from(rulesB64, "base64");
    const sig = signData(privateKey, rulesData);

    // Omitting userId used to bypass dedup, so one key reached any threshold. Signers
    // are now counted by key, which no label can influence.
    const signaturesB64 = buildSignaturesBase64([
      { signature: sig },
      { signature: sig },
    ]);

    expect(verifyGovernanceRules(rulesB64, signaturesB64, 2, [pem])).toBe(false);
    expect(verifyGovernanceRules(rulesB64, signaturesB64, 1, [pem])).toBe(true);
  });
});
// =============================================================================
// verifyGovernanceRulesSignatures - the one place the threshold is evaluated
// =============================================================================

// The five distinct-signer cases are covered against the PEM wrapper above, which now
// delegates here. These cover the core's own contract: it throws instead of returning.
describe("verifyGovernanceRulesSignatures", () => {
  const rulesData = Buffer.from("rules-container-bytes");

  it("reports the distinct-signer shortfall, not the entry count", () => {
    const { privateKey, publicKey } = generateP256KeyPair();

    expect(() =>
      verifyGovernanceRulesSignatures(
        rulesData,
        [
          { signature: signData(privateKey, rulesData) },
          { signature: signData(privateKey, rulesData) },
          { signature: signData(privateKey, rulesData) },
        ],
        [publicKey],
        2
      )
    ).toThrow("insufficient distinct valid SuperAdmin signers: got 1, need 2");
  });

  it("skips null and empty signature entries without aborting", () => {
    const { privateKey, publicKey } = generateP256KeyPair();

    expect(() =>
      verifyGovernanceRulesSignatures(
        rulesData,
        [
          { signature: undefined },
          { signature: "" },
          { signature: signData(privateKey, rulesData) },
        ],
        [publicKey],
        1
      )
    ).not.toThrow();
  });

  it("throws IntegrityError when no key verifies", () => {
    const configured = generateP256KeyPair();
    const stranger = generateP256KeyPair();

    expect(() =>
      verifyGovernanceRulesSignatures(
        rulesData,
        [{ signature: signData(stranger.privateKey, rulesData) }],
        [configured.publicKey],
        1
      )
    ).toThrow(IntegrityError);
  });

  it("throws ConfigurationError on a non-positive threshold", () => {
    const { publicKey } = generateP256KeyPair();

    expect(() =>
      verifyGovernanceRulesSignatures(rulesData, [], [publicKey], 0)
    ).toThrow(ConfigurationError);
  });
});

// The single-key form the other three SDKs expose (Go VerifySignatureWithKey, Java
// verifySignature, Python verify_raw_signature). TypeScript had only the quorum-shaped
// isValidSignature, so checking one known signer meant passing a one-element array.
describe("verifySignatureWithKey", () => {
  it("accepts a signature made by that key and rejects one made by another", () => {
    const a = crypto.generateKeyPairSync("ec", { namedCurve: "P-256" });
    const b = crypto.generateKeyPairSync("ec", { namedCurve: "P-256" });
    const data = Buffer.from("payload");
    const sig = signData(a.privateKey, data);

    expect(verifySignatureWithKey(data, sig, a.publicKey)).toBe(true);
    expect(verifySignatureWithKey(data, sig, b.publicKey)).toBe(false);
  });

  it("returns false for a missing key, a tampered payload or a malformed signature", () => {
    const kp = crypto.generateKeyPairSync("ec", { namedCurve: "P-256" });
    const data = Buffer.from("payload");
    const sig = signData(kp.privateKey, data);

    expect(verifySignatureWithKey(data, sig, undefined)).toBe(false);
    expect(verifySignatureWithKey(Buffer.from("tampered"), sig, kp.publicKey)).toBe(false);
    expect(verifySignatureWithKey(data, "not-base64!!", kp.publicKey)).toBe(false);
  });
});
