/**
 * ECDSA signing utilities for Taurus-PROTECT SDK.
 *
 * Provides ECDSA P-256 signing and verification with raw r||s format,
 * compatible with Java's SHA256withPLAIN-ECDSA.
 */

import * as crypto from "crypto";

import { IntegrityError } from "../errors";

/**
 * Size of P-256 curve parameters (r and s) in bytes.
 */
const P256_COMPONENT_SIZE = 32;

/**
 * Total size of raw r||s signature for P-256.
 */
const P256_RAW_SIGNATURE_SIZE = P256_COMPONENT_SIZE * 2;

/**
 * Sign data using ECDSA with SHA-256 (plain format).
 *
 * Returns signature in raw r||s format (base64-encoded), matching
 * the Java SDK's SHA256withPLAIN-ECDSA output.
 *
 * @param privateKey - ECDSA private key (P-256)
 * @param data - Data to sign
 * @returns Base64-encoded raw r||s signature (64 bytes before encoding)
 *
 * @example
 * ```typescript
 * const signature = signData(privateKey, Buffer.from("request hash data"));
 * // signature is base64-encoded 64-byte raw r||s format
 * ```
 */
export function signData(privateKey: crypto.KeyObject, data: Uint8Array): string {
  // Validate P-256 curve for defense-in-depth
  const keyDetails = privateKey.asymmetricKeyDetails;
  if (keyDetails && keyDetails.namedCurve !== 'prime256v1' && keyDetails.namedCurve !== 'P-256') {
    throw new Error(`Only P-256 (prime256v1) keys are supported for signing, got: ${keyDetails.namedCurve ?? 'unknown'}`);
  }

  // Sign with ECDSA - Node.js returns DER-encoded signature by default
  const derSignature = crypto.sign("sha256", Buffer.from(data), {
    key: privateKey,
    dsaEncoding: "der",
  });

  // Convert DER to raw r||s format (matching Java's PLAIN-ECDSA)
  const rawSignature = derToRaw(derSignature);

  return Buffer.from(rawSignature).toString("base64");
}

/**
 * Verify an ECDSA signature.
 *
 * Expects signature in raw r||s format (base64-encoded), matching
 * the Java SDK's SHA256withPLAIN-ECDSA format.
 *
 * @param publicKey - ECDSA public key (P-256)
 * @param data - The signed data
 * @param signatureB64 - Base64-encoded raw r||s signature
 * @returns True if signature is valid, false otherwise
 *
 * @example
 * ```typescript
 * const valid = verifySignature(publicKey, Buffer.from("data"), signature);
 * if (!valid) {
 *   throw new IntegrityError("Signature verification failed");
 * }
 * ```
 */
export function verifySignature(
  publicKey: crypto.KeyObject,
  data: Uint8Array,
  signatureB64: string
): boolean {
  try {
    // Validate P-256 curve for defense-in-depth
    const keyDetails = publicKey.asymmetricKeyDetails;
    if (keyDetails && keyDetails.namedCurve !== 'prime256v1' && keyDetails.namedCurve !== 'P-256') {
      throw new IntegrityError('Invalid key curve: only P-256 (secp256r1) is supported');
    }

    // Decode base64 signature
    const rawSignature = Buffer.from(signatureB64, "base64");

    // Signature should be exactly 64 bytes (r||s for P-256)
    if (rawSignature.length !== P256_RAW_SIGNATURE_SIZE) {
      return false;
    }

    // Convert raw r||s to DER format for verification
    const derSignature = rawToDer(rawSignature);

    // Verify signature
    return crypto.verify(
      "sha256",
      Buffer.from(data),
      {
        key: publicKey,
        dsaEncoding: "der",
      },
      derSignature
    );
  } catch (error: unknown) {
    // Always re-throw IntegrityError (e.g. invalid curve)
    if (error instanceof IntegrityError) {
      throw error;
    }
    // Only swallow expected crypto/signature errors
    if (error instanceof Error &&
        (error.message.includes('signature') ||
         error.message.includes('key') ||
         error.message.includes('Invalid') ||
         error.message.includes('decode') ||
         error.message.includes('ERR_OSSL') ||
         error.message.includes('DER'))) {
      return false;
    }
    throw error;
  }
}

/**
 * Convert DER-encoded ECDSA signature to raw r||s format.
 *
 * DER format: 0x30 [total-length] 0x02 [r-length] [r] 0x02 [s-length] [s]
 * Raw format: [r (32 bytes)] [s (32 bytes)]
 *
 * SECURITY: This function validates all buffer accesses to prevent
 * buffer over-reads from malformed DER input.
 *
 * @param der - DER-encoded signature
 * @returns Raw r||s signature (64 bytes for P-256)
 * @throws Error if DER format is invalid or buffer bounds are exceeded
 */
function derToRaw(der: Buffer): Uint8Array {
  // Minimum DER signature size: 0x30 + len + 0x02 + len + r(1) + 0x02 + len + s(1) = 8 bytes
  if (der.length < 8) {
    throw new Error("Invalid DER signature: too short");
  }

  // Parse DER structure
  if (der[0] !== 0x30) {
    throw new Error("Invalid DER signature: expected SEQUENCE tag");
  }

  let offset = 2; // Skip SEQUENCE tag and length
  let sequenceLength = der[1]!;

  // Handle extended length encoding (lengths > 127 bytes)
  if (sequenceLength > 0x80) {
    const numLengthBytes = sequenceLength - 0x80;

    // Validate: 0x80 alone (indefinite length) is invalid for DER
    if (numLengthBytes === 0) {
      throw new Error(
        "Invalid DER signature: indefinite length encoding not allowed"
      );
    }

    // Validate: we have enough bytes for the extended length
    if (2 + numLengthBytes > der.length) {
      throw new Error(
        "Invalid DER signature: extended length bytes exceed buffer"
      );
    }

    // Read the actual length from the extended length bytes
    sequenceLength = 0;
    for (let i = 0; i < numLengthBytes; i++) {
      const byteIndex = 2 + i;
      if (byteIndex >= der.length) {
        throw new Error('DER signature truncated: extended length exceeds buffer');
      }
      sequenceLength = (sequenceLength << 8) | der[byteIndex]!;
    }

    offset = 2 + numLengthBytes;

    // Validate the declared sequence length matches remaining buffer
    if (sequenceLength !== der.length - offset) {
      throw new Error(
        "Invalid DER signature: declared length does not match actual content"
      );
    }
  } else if (sequenceLength === 0x80) {
    // 0x80 alone means indefinite length - not valid in DER
    throw new Error(
      "Invalid DER signature: indefinite length encoding not allowed"
    );
  } else {
    // Short form: validate sequence length
    if (sequenceLength !== der.length - 2) {
      throw new Error(
        "Invalid DER signature: declared length does not match actual content"
      );
    }
  }

  // Parse r - check bounds before accessing
  if (offset >= der.length) {
    throw new Error("Invalid DER signature: truncated before r INTEGER tag");
  }
  if (der[offset] !== 0x02) {
    throw new Error("Invalid DER signature: expected INTEGER tag for r");
  }
  offset++;

  if (offset >= der.length) {
    throw new Error("Invalid DER signature: truncated before r length");
  }
  const rLength = der[offset]!;
  offset++;

  // Validate r length doesn't exceed remaining buffer
  if (rLength > der.length - offset) {
    throw new Error("Invalid DER signature: r length exceeds buffer");
  }
  const rStart = offset;
  offset += rLength;

  // Parse s - check bounds before accessing
  if (offset >= der.length) {
    throw new Error("Invalid DER signature: truncated before s INTEGER tag");
  }
  if (der[offset] !== 0x02) {
    throw new Error("Invalid DER signature: expected INTEGER tag for s");
  }
  offset++;

  if (offset >= der.length) {
    throw new Error("Invalid DER signature: truncated before s length");
  }
  const sLength = der[offset]!;
  offset++;

  // Validate s length doesn't exceed remaining buffer
  if (sLength > der.length - offset) {
    throw new Error("Invalid DER signature: s length exceeds buffer");
  }
  const sStart = offset;

  // Extract r and s, handling leading zeros
  const r = extractInteger(der.subarray(rStart, rStart + rLength));
  const s = extractInteger(der.subarray(sStart, sStart + sLength));

  // Combine into raw format
  const raw = new Uint8Array(P256_RAW_SIGNATURE_SIZE);
  raw.set(r, 0);
  raw.set(s, P256_COMPONENT_SIZE);

  return raw;
}

/**
 * Convert raw r||s signature to DER format.
 *
 * @param raw - Raw r||s signature (64 bytes for P-256)
 * @returns DER-encoded signature
 */
function rawToDer(raw: Buffer): Buffer {
  // Split into r and s
  const r = raw.subarray(0, P256_COMPONENT_SIZE);
  const s = raw.subarray(P256_COMPONENT_SIZE);

  // Encode as integers (add leading zero if high bit set)
  const rEncoded = encodeInteger(r);
  const sEncoded = encodeInteger(s);

  // Build DER structure
  const totalLength = rEncoded.length + sEncoded.length;
  const der = Buffer.alloc(2 + totalLength);

  let offset = 0;
  der[offset++] = 0x30; // SEQUENCE tag
  der[offset++] = totalLength;
  rEncoded.copy(der, offset);
  offset += rEncoded.length;
  sEncoded.copy(der, offset);

  return der;
}

/**
 * Extract an integer from DER INTEGER encoding to fixed-size bytes.
 *
 * @param data - DER INTEGER value (may have leading zeros)
 * @returns Fixed-size integer (32 bytes for P-256)
 */
function extractInteger(data: Uint8Array): Uint8Array {
  const result = new Uint8Array(P256_COMPONENT_SIZE);

  if (data.length > P256_COMPONENT_SIZE) {
    // Strip leading zeros (DER adds 0x00 prefix if high bit is set)
    const startOffset = data.length - P256_COMPONENT_SIZE;
    result.set(data.subarray(startOffset));
  } else if (data.length < P256_COMPONENT_SIZE) {
    // Pad with leading zeros
    result.set(data, P256_COMPONENT_SIZE - data.length);
  } else {
    result.set(data);
  }

  return result;
}

/**
 * Encode a fixed-size integer to DER INTEGER format.
 *
 * @param data - Fixed-size integer (32 bytes for P-256)
 * @returns DER INTEGER encoding (0x02 [length] [value])
 */
function encodeInteger(data: Uint8Array): Buffer {
  // Find first non-zero byte
  let start = 0;
  while (start < data.length - 1 && data[start] === 0) {
    start++;
  }

  // Check if we need to add a leading zero (if high bit is set)
  const needsLeadingZero = (data[start]! & 0x80) !== 0;
  const valueLength = data.length - start;
  const totalLength = needsLeadingZero ? valueLength + 1 : valueLength;

  const encoded = Buffer.alloc(2 + totalLength);
  let offset = 0;
  encoded[offset++] = 0x02; // INTEGER tag
  encoded[offset++] = totalLength;

  if (needsLeadingZero) {
    encoded[offset++] = 0x00;
  }

  data.subarray(start).forEach((byte) => {
    encoded[offset++] = byte;
  });

  return encoded;
}
