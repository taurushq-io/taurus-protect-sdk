/**
 * Constant-time comparison utilities to prevent timing attacks.
 *
 * These functions use Node.js crypto.timingSafeEqual() to ensure that
 * the comparison time does not leak information about the content,
 * which is critical for security-sensitive comparisons like hash
 * and signature verification.
 */

import { timingSafeEqual } from "crypto";

/**
 * Compare two strings in constant time.
 *
 * Uses crypto.timingSafeEqual to prevent timing attacks by ensuring
 * the comparison time does not leak information about the content.
 *
 * This function should be used for security-sensitive comparisons such as:
 * - Hash verification
 * - Signature verification
 * - Token comparison
 *
 * @param a - First string to compare
 * @param b - Second string to compare
 * @returns true if strings are equal, false otherwise
 *
 * @example
 * ```typescript
 * const computedHash = calculateHash(data);
 * const providedHash = response.metadata.hash;
 *
 * if (!constantTimeCompare(computedHash, providedHash)) {
 *   throw new IntegrityError("Hash mismatch");
 * }
 * ```
 */
export function constantTimeCompare(a: string, b: string): boolean {
  // Convert strings to buffers using UTF-8 encoding
  const bufferA = Buffer.from(a, "utf-8");
  const bufferB = Buffer.from(b, "utf-8");

  // timingSafeEqual requires buffers of equal length
  // If lengths differ, we still need constant-time behavior
  if (bufferA.length !== bufferB.length) {
    // Compare against a buffer of the same length to maintain constant time
    // We use bufferA as a dummy to ensure we still do the comparison work
    timingSafeEqual(bufferA, bufferA);
    return false;
  }

  return timingSafeEqual(bufferA, bufferB);
}

/**
 * Compare two byte arrays in constant time.
 *
 * Uses crypto.timingSafeEqual to prevent timing attacks by ensuring
 * the comparison time does not leak information about the content.
 *
 * @param a - First byte array to compare
 * @param b - Second byte array to compare
 * @returns true if bytes are equal, false otherwise
 *
 * @example
 * ```typescript
 * const computedSignature = sign(data, privateKey);
 * const providedSignature = response.signature;
 *
 * if (!constantTimeCompareBytes(computedSignature, providedSignature)) {
 *   throw new IntegrityError("Signature mismatch");
 * }
 * ```
 */
export function constantTimeCompareBytes(
  a: Uint8Array,
  b: Uint8Array
): boolean {
  // Convert Uint8Array to Buffer for timingSafeEqual
  const bufferA = Buffer.from(a);
  const bufferB = Buffer.from(b);

  // timingSafeEqual requires buffers of equal length
  // If lengths differ, we still need constant-time behavior
  if (bufferA.length !== bufferB.length) {
    // Compare against a buffer of the same length to maintain constant time
    timingSafeEqual(bufferA, bufferA);
    return false;
  }

  return timingSafeEqual(bufferA, bufferB);
}
