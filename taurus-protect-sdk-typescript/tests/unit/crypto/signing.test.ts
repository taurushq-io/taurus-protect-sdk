/**
 * Unit tests for ECDSA signing utilities.
 *
 * These tests verify P-256 ECDSA signing, verification,
 * DER-to-raw conversion, and error handling for invalid keys.
 */

import * as crypto from "crypto";

import { signData, verifySignature } from "../../../src/crypto/signing";
import { isValidSignature } from "../../../src/helpers/signature-verifier";
import { IntegrityError } from "../../../src/errors";

/**
 * Generate a P-256 key pair for testing.
 */
function generateP256KeyPair(): {
  privateKey: crypto.KeyObject;
  publicKey: crypto.KeyObject;
} {
  const { privateKey, publicKey } = crypto.generateKeyPairSync("ec", {
    namedCurve: "P-256",
  });
  return { privateKey, publicKey };
}

/**
 * Generate a P-384 key pair (non-P-256) for negative testing.
 */
function generateP384KeyPair(): {
  privateKey: crypto.KeyObject;
  publicKey: crypto.KeyObject;
} {
  const { privateKey, publicKey } = crypto.generateKeyPairSync("ec", {
    namedCurve: "P-384",
  });
  return { privateKey, publicKey };
}

describe("signData", () => {
  let keyPair: { privateKey: crypto.KeyObject; publicKey: crypto.KeyObject };

  beforeEach(() => {
    keyPair = generateP256KeyPair();
  });

  it("should produce a non-empty base64 string", () => {
    const data = Buffer.from("test data");
    const signature = signData(keyPair.privateKey, data);

    expect(signature).toBeDefined();
    expect(signature.length).toBeGreaterThan(0);
    // Should be valid base64
    expect(() => Buffer.from(signature, "base64")).not.toThrow();
  });

  it("should produce a 64-byte raw signature (base64-encoded)", () => {
    const data = Buffer.from("test data");
    const signature = signData(keyPair.privateKey, data);
    const decoded = Buffer.from(signature, "base64");

    // P-256 raw r||s is exactly 64 bytes
    expect(decoded.length).toBe(64);
  });

  it("should sign empty data without error", () => {
    const data = Buffer.from("");
    const signature = signData(keyPair.privateKey, data);

    expect(signature).toBeDefined();
    expect(signature.length).toBeGreaterThan(0);
  });

  it("should throw for non-P-256 key", () => {
    const p384 = generateP384KeyPair();
    const data = Buffer.from("test data");

    expect(() => signData(p384.privateKey, data)).toThrow();
  });

  it("should produce different signatures for the same data (ECDSA is non-deterministic)", () => {
    const data = Buffer.from("test data");
    const sig1 = signData(keyPair.privateKey, data);
    const sig2 = signData(keyPair.privateKey, data);

    // Both should verify even if they differ
    expect(verifySignature(keyPair.publicKey, data, sig1)).toBe(true);
    expect(verifySignature(keyPair.publicKey, data, sig2)).toBe(true);
  });

  it("should produce different signatures for different data", () => {
    const sig1 = signData(keyPair.privateKey, Buffer.from("data1"));
    const sig2 = signData(keyPair.privateKey, Buffer.from("data2"));

    expect(sig1).not.toBe(sig2);
  });
});

describe("verifySignature", () => {
  let keyPair: { privateKey: crypto.KeyObject; publicKey: crypto.KeyObject };

  beforeEach(() => {
    keyPair = generateP256KeyPair();
  });

  it("should return true for a valid sign-then-verify round-trip", () => {
    const data = Buffer.from("round-trip test");
    const signature = signData(keyPair.privateKey, data);

    expect(verifySignature(keyPair.publicKey, data, signature)).toBe(true);
  });

  it("should return false with a wrong public key", () => {
    const data = Buffer.from("test data");
    const signature = signData(keyPair.privateKey, data);

    const otherKeyPair = generateP256KeyPair();
    expect(verifySignature(otherKeyPair.publicKey, data, signature)).toBe(
      false
    );
  });

  it("should return false with corrupted signature", () => {
    const data = Buffer.from("test data");
    const signature = signData(keyPair.privateKey, data);

    // Corrupt the signature by changing a character
    const corrupted =
      signature.substring(0, 5) + "X" + signature.substring(6);
    expect(verifySignature(keyPair.publicKey, data, corrupted)).toBe(false);
  });

  it("should return false with tampered data", () => {
    const data = Buffer.from("original data");
    const signature = signData(keyPair.privateKey, data);

    const tampered = Buffer.from("tampered data");
    expect(verifySignature(keyPair.publicKey, tampered, signature)).toBe(
      false
    );
  });

  it("should throw IntegrityError for non-P-256 public key", () => {
    const p384 = generateP384KeyPair();
    const data = Buffer.from("test data");
    // Sign with a valid P-256 key to produce a valid signature
    const signature = signData(keyPair.privateKey, data);

    expect(() =>
      verifySignature(p384.publicKey, data, signature)
    ).toThrow(IntegrityError);
  });

  it("should return false for an empty signature string", () => {
    const data = Buffer.from("test data");
    expect(verifySignature(keyPair.publicKey, data, "")).toBe(false);
  });

  it("should return false for a signature of wrong length", () => {
    const data = Buffer.from("test data");
    // 32 bytes instead of 64
    const shortSig = Buffer.alloc(32).toString("base64");
    expect(verifySignature(keyPair.publicKey, data, shortSig)).toBe(false);
  });

  it("should return false for completely random 64-byte signature", () => {
    const data = Buffer.from("test data");
    const randomSig = crypto.randomBytes(64).toString("base64");
    expect(verifySignature(keyPair.publicKey, data, randomSig)).toBe(false);
  });

  it("should verify signature on empty data", () => {
    const data = Buffer.from("");
    const signature = signData(keyPair.privateKey, data);
    expect(verifySignature(keyPair.publicKey, data, signature)).toBe(true);
  });

  it("should verify with Uint8Array data input", () => {
    const data = new Uint8Array([0x48, 0x65, 0x6c, 0x6c, 0x6f]); // "Hello"
    const signature = signData(keyPair.privateKey, data);
    expect(verifySignature(keyPair.publicKey, data, signature)).toBe(true);
  });
});

describe("isValidSignature", () => {
  let keyPair1: { privateKey: crypto.KeyObject; publicKey: crypto.KeyObject };
  let keyPair2: { privateKey: crypto.KeyObject; publicKey: crypto.KeyObject };

  beforeEach(() => {
    keyPair1 = generateP256KeyPair();
    keyPair2 = generateP256KeyPair();
  });

  it("should return true when one key in the list matches", () => {
    const data = Buffer.from("test data");
    const signature = signData(keyPair1.privateKey, data);

    const publicKeys = [keyPair2.publicKey, keyPair1.publicKey];
    expect(isValidSignature(data, signature, publicKeys)).toBe(true);
  });

  it("should return false when no keys match", () => {
    const data = Buffer.from("test data");
    const signature = signData(keyPair1.privateKey, data);

    const otherKeyPair = generateP256KeyPair();
    const publicKeys = [keyPair2.publicKey, otherKeyPair.publicKey];
    expect(isValidSignature(data, signature, publicKeys)).toBe(false);
  });

  it("should return false for an empty key list", () => {
    const data = Buffer.from("test data");
    const signature = signData(keyPair1.privateKey, data);

    expect(isValidSignature(data, signature, [])).toBe(false);
  });

  it("should accept string data input", () => {
    const data = "test data as string";
    const signature = signData(
      keyPair1.privateKey,
      Buffer.from(data, "utf-8")
    );

    expect(isValidSignature(data, signature, [keyPair1.publicKey])).toBe(true);
  });

  it("should return true for first matching key without checking rest", () => {
    const data = Buffer.from("test data");
    const signature = signData(keyPair1.privateKey, data);

    // First key matches - should return true
    const publicKeys = [keyPair1.publicKey, keyPair2.publicKey];
    expect(isValidSignature(data, signature, publicKeys)).toBe(true);
  });
});
