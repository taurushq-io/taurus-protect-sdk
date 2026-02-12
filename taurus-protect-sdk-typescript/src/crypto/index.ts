/**
 * Cryptographic utilities for Taurus-PROTECT SDK.
 *
 * This module provides:
 * - SHA-256 hashing and HMAC-SHA256 functions
 * - PEM key parsing for ECDSA P-256 keys
 * - ECDSA signing and verification (raw r||s format)
 * - TPV1-HMAC-SHA256 authentication for API requests
 */

// Hashing functions
export {
  calculateHexHash,
  calculateSha256Bytes,
  calculateBase64Hmac,
  verifyBase64Hmac,
  constantTimeCompare,
  constantTimeCompareBytes,
} from "./hashing";

// Key handling
export {
  decodePrivateKeyPem,
  decodePublicKeyPem,
  decodePublicKeysPem,
  encodePublicKeyPem,
  getPublicKeyFromPrivate,
} from "./keys";

// ECDSA signing
export { signData, verifySignature } from "./signing";

// TPV1 authentication
export { TPV1Auth, calculateSignedHeader } from "./tpv1";
