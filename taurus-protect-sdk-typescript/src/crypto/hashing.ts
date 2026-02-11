/**
 * Hashing utilities for Taurus-PROTECT SDK.
 *
 * Provides SHA-256 and HMAC-SHA256 functions compatible with the Java SDK.
 */

import * as crypto from "crypto";

/**
 * Calculate SHA-256 hash and return hex-encoded string.
 *
 * This is used for request hash verification.
 *
 * @param data - The data to hash
 * @returns Hex-encoded SHA-256 hash (lowercase)
 *
 * @example
 * ```typescript
 * const hash = calculateHexHash("hello");
 * // Returns: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
 * ```
 */
export function calculateHexHash(data: string): string {
  return crypto.createHash("sha256").update(data, "utf8").digest("hex");
}

/**
 * Calculate SHA-256 hash of bytes.
 *
 * @param data - The bytes to hash
 * @returns Raw SHA-256 hash bytes (32 bytes)
 */
export function calculateSha256Bytes(data: Uint8Array): Uint8Array {
  return new Uint8Array(crypto.createHash("sha256").update(data).digest());
}

/**
 * Calculate HMAC-SHA256 and return base64-encoded result.
 *
 * @param secret - The secret key as bytes
 * @param data - The data to sign
 * @returns Base64-encoded HMAC-SHA256 signature
 *
 * @example
 * ```typescript
 * const secret = Buffer.from("my-secret-key");
 * const hmac = calculateBase64Hmac(secret, "data-to-sign");
 * ```
 */
export function calculateBase64Hmac(secret: Uint8Array, data: string): string {
  return crypto
    .createHmac("sha256", Buffer.from(secret))
    .update(data, "utf8")
    .digest("base64");
}

/**
 * Verify HMAC-SHA256 signature using timing-safe comparison.
 *
 * This MUST be used when comparing cryptographic HMACs to prevent timing attacks.
 *
 * @param secret - The secret key as bytes
 * @param data - The original data
 * @param expectedHmac - The expected base64-encoded HMAC
 * @returns True if the HMAC matches
 *
 * @example
 * ```typescript
 * const secret = Buffer.from("my-secret-key");
 * const isValid = verifyBase64Hmac(secret, "data", "expected-hmac-base64");
 * if (!isValid) {
 *   throw new Error("HMAC verification failed");
 * }
 * ```
 */
export function verifyBase64Hmac(
  secret: Uint8Array,
  data: string,
  expectedHmac: string
): boolean {
  const computed = calculateBase64Hmac(secret, data);
  return constantTimeCompare(computed, expectedHmac);
}

/**
 * Compare two strings in constant time to prevent timing attacks.
 *
 * This MUST be used when comparing cryptographic hashes or signatures.
 *
 * @param a - First string
 * @param b - Second string
 * @returns True if strings are equal
 *
 * @example
 * ```typescript
 * const isEqual = constantTimeCompare(computedHash, expectedHash);
 * ```
 */
export function constantTimeCompare(a: string, b: string): boolean {
  const bufA = Buffer.from(a, "utf8");
  const bufB = Buffer.from(b, "utf8");

  // timingSafeEqual requires equal length buffers
  // If lengths differ, we still need constant-time behavior to prevent timing attacks
  if (bufA.length !== bufB.length) {
    // Perform a dummy comparison to maintain constant time
    crypto.timingSafeEqual(bufA, bufA);
    return false;
  }

  return crypto.timingSafeEqual(bufA, bufB);
}

/**
 * Compare two byte arrays in constant time to prevent timing attacks.
 *
 * @param a - First byte array
 * @param b - Second byte array
 * @returns True if byte arrays are equal
 */
export function constantTimeCompareBytes(
  a: Uint8Array,
  b: Uint8Array
): boolean {
  const bufA = Buffer.from(a);
  const bufB = Buffer.from(b);

  // timingSafeEqual requires equal length buffers
  // If lengths differ, we still need constant-time behavior to prevent timing attacks
  if (bufA.length !== bufB.length) {
    // Perform a dummy comparison to maintain constant time
    crypto.timingSafeEqual(bufA, bufA);
    return false;
  }

  return crypto.timingSafeEqual(bufA, bufB);
}
