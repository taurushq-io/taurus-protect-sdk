/**
 * Unit tests for signature-verifier.ts.
 *
 * Tests the SuperAdmin signature verification functions:
 * - isValidSignature(): single signature verification against multiple keys
 * - verifyGovernanceRules(): multi-signature threshold verification
 */

import * as crypto from "crypto";

import { signData, encodePublicKeyPem } from "../../../src/crypto";
import { isValidSignature, verifyGovernanceRules } from "../../../src/helpers/signature-verifier";

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
  it("should return true when minValidSignatures is 0", () => {
    expect(verifyGovernanceRules("abc", "abc", 0, [])).toBe(true);
  });

  it("should return true when minValidSignatures is negative", () => {
    expect(verifyGovernanceRules("abc", "abc", -1, [])).toBe(true);
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

  it("should handle signatures without userId (no dedup)", () => {
    const { privateKey, publicKey } = generateP256KeyPair();
    const pem = keyToPem(publicKey);

    const rulesB64 = Buffer.from("no-userid-rules").toString("base64");
    const rulesData = Buffer.from(rulesB64, "base64");
    const sig = signData(privateKey, rulesData);

    // Two signatures without userId: both should be counted
    const signaturesB64 = buildSignaturesBase64([
      { signature: sig },
      { signature: sig },
    ]);

    expect(verifyGovernanceRules(rulesB64, signaturesB64, 2, [pem])).toBe(true);
  });
});
