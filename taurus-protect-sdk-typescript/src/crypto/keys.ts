/**
 * Key handling utilities for Taurus-PROTECT SDK.
 *
 * Provides PEM key parsing for ECDSA P-256 keys.
 */

import * as crypto from "crypto";

/**
 * Decode a PEM-encoded ECDSA private key.
 *
 * Supports both EC PRIVATE KEY and PKCS#8 formats.
 *
 * @param pem - PEM-encoded private key string
 * @returns Node.js KeyObject for the private key
 * @throws Error if the key cannot be decoded or is not an EC key
 *
 * @example
 * ```typescript
 * const privateKey = decodePrivateKeyPem(`-----BEGIN EC PRIVATE KEY-----
 * MHQCAQEEIDfN...
 * -----END EC PRIVATE KEY-----`);
 * ```
 */
export function decodePrivateKeyPem(pem: string): crypto.KeyObject {
  try {
    const key = crypto.createPrivateKey({
      key: pem,
      format: "pem",
    });

    // Verify it's an EC key
    const keyDetails = key.asymmetricKeyDetails;
    if (!keyDetails || keyDetails.namedCurve !== "prime256v1") {
      // Also accept P-256 alias
      if (!keyDetails || keyDetails.namedCurve !== "P-256") {
        throw new Error(
          `Expected P-256/prime256v1 EC key, got: ${keyDetails?.namedCurve ?? "unknown"}`
        );
      }
    }

    return key;
  } catch (error) {
    if (error instanceof Error) {
      throw new Error(`Failed to decode private key: ${error.message}`);
    }
    throw new Error("Failed to decode private key: unknown error");
  }
}

/**
 * Decode a PEM-encoded ECDSA public key.
 *
 * @param pem - PEM-encoded public key string
 * @returns Node.js KeyObject for the public key
 * @throws Error if the key cannot be decoded or is not an EC key
 *
 * @example
 * ```typescript
 * const publicKey = decodePublicKeyPem(`-----BEGIN PUBLIC KEY-----
 * MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...
 * -----END PUBLIC KEY-----`);
 * ```
 */
export function decodePublicKeyPem(pem: string): crypto.KeyObject {
  try {
    const key = crypto.createPublicKey({
      key: pem,
      format: "pem",
    });

    // Verify it's an EC key
    const keyDetails = key.asymmetricKeyDetails;
    if (!keyDetails || keyDetails.namedCurve !== "prime256v1") {
      // Also accept P-256 alias
      if (!keyDetails || keyDetails.namedCurve !== "P-256") {
        throw new Error(
          `Expected P-256/prime256v1 EC key, got: ${keyDetails?.namedCurve ?? "unknown"}`
        );
      }
    }

    return key;
  } catch (error) {
    if (error instanceof Error) {
      throw new Error(`Failed to decode public key: ${error.message}`);
    }
    throw new Error("Failed to decode public key: unknown error");
  }
}

/**
 * Decode multiple PEM-encoded public keys.
 *
 * @param pems - Array of PEM-encoded public key strings
 * @returns Array of Node.js KeyObject for the public keys
 * @throws Error if any key cannot be decoded
 *
 * @example
 * ```typescript
 * const publicKeys = decodePublicKeysPem([pem1, pem2, pem3]);
 * ```
 */
export function decodePublicKeysPem(pems: string[]): crypto.KeyObject[] {
  return pems.map((pem) => decodePublicKeyPem(pem));
}

/**
 * Encode a public key to PEM format.
 *
 * @param key - Node.js KeyObject for the public key
 * @returns PEM-encoded public key string
 */
export function encodePublicKeyPem(key: crypto.KeyObject): string {
  return key.export({ type: "spki", format: "pem" }) as string;
}

/**
 * Extract the public key from a private key.
 *
 * @param privateKey - Node.js KeyObject for the private key
 * @returns Node.js KeyObject for the corresponding public key
 */
export function getPublicKeyFromPrivate(
  privateKey: crypto.KeyObject
): crypto.KeyObject {
  return crypto.createPublicKey(privateKey);
}
