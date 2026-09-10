/**
 * Signature verification utilities for governance rules.
 *
 * This module provides functions to verify SuperAdmin signatures on
 * governance rules and other cryptographically protected data.
 */

import { createHash, type KeyObject } from "crypto";

import { verifySignature, decodePublicKeysPem } from "../crypto";
import { ConfigurationError, IntegrityError } from "../errors";
import { isCryptoVerificationError } from "../crypto";
import { strictBase64Decode } from "./strict-base64";

/**
 * Returns the first key that verifies the signature, or undefined.
 */
function matchingKey(
  data: Buffer | string,
  signatureBase64: string,
  publicKeys: KeyObject[]
): KeyObject | undefined {
  const dataBuffer =
    typeof data === "string" ? Buffer.from(data, "utf-8") : data;

  for (const publicKey of publicKeys) {
    try {
      if (verifySignature(publicKey, dataBuffer, signatureBase64)) {
        return publicKey;
      }
    } catch (error: unknown) {
      // This key did not verify the signature; try the next configured key.
      if (isCryptoVerificationError(error)) {
        continue;
      }
      throw error;
    }
  }

  return undefined;
}

/**
 * Identifies a key by its encoded bytes, so the same key configured twice
 * counts as one signer.
 */
export function keyFingerprint(publicKey: KeyObject): string {
  const der = publicKey.export({ type: "spki", format: "der" });
  return createHash("sha256").update(der).digest("hex");
}

/**
 * Verifies a single signature against a set of public keys.
 *
 * @param data - The signed data (as buffer or string)
 * @param signatureBase64 - Base64-encoded signature
 * @param publicKeys - Array of public keys to verify against
 * @returns true if the signature is valid for any of the keys
 */
export function isValidSignature(
  data: Buffer | string,
  signatureBase64: string,
  publicKeys: KeyObject[]
): boolean {
  return matchingKey(data, signatureBase64, publicKeys) !== undefined;
}

/**
 * Verifies a base64 signature against ONE public key.
 *
 * The single-key form the other three SDKs expose (Go VerifySignatureWithKey, Java
 * SignatureVerifier.verifySignature, Python verify_raw_signature). Without it a caller
 * checking a known signer had to pass a one-element array to
 * {@link isValidSignature}, which reads as a quorum check.
 *
 * Returns false rather than throwing for a missing key or malformed signature: the data,
 * signature and key are all public here, so there is nothing to leak and nothing to
 * distinguish between failure modes.
 *
 * @param data - The signed payload
 * @param signatureBase64 - The signature, base64-encoded raw r||s
 * @param publicKey - The public key to check against
 * @returns true when the signature verifies under this key
 */
export function verifySignatureWithKey(
  data: Buffer | string,
  signatureBase64: string,
  publicKey: KeyObject | undefined | null
): boolean {
  if (!publicKey) {
    return false;
  }
  const dataBuffer = typeof data === "string" ? Buffer.from(data, "utf-8") : data;
  return verifySignature(publicKey, dataBuffer, signatureBase64);
}

/**
 * Verifies that the rules container is signed by the required number of DISTINCT
 * SuperAdmin keys. Throws on failure, returns normally on success.
 *
 * This is the single place the threshold is evaluated — every caller routes here.
 * minValidSignatures counts distinct signing keys, not signature entries: ECDSA is
 * randomized, so counting entries would let a single key produce as many valid
 * signatures as any threshold demands, making a 2-of-N quorum no stronger than 1-of-N.
 *
 * @throws {@link ConfigurationError} If minValidSignatures is not positive
 * @throws {@link IntegrityError} If too few distinct SuperAdmin keys signed
 */
export function verifyGovernanceRulesSignatures(
  rulesContainerData: Buffer,
  signatures: ReadonlyArray<{ readonly signature?: string }>,
  superAdminKeys: KeyObject[],
  minValidSignatures: number
): void {
  if (minValidSignatures <= 0) {
    throw new ConfigurationError("minValidSignatures must be positive");
  }

  if (rulesContainerData.length === 0) {
    throw new IntegrityError("rules container data cannot be empty");
  }

  if (superAdminKeys.length === 0) {
    throw new IntegrityError("no SuperAdmin keys configured for verification");
  }

  if (signatures.length === 0) {
    throw new IntegrityError("no signatures provided");
  }

  const signers = new Set<string>();

  for (const sig of signatures) {
    if (!sig.signature) continue;

    const key = matchingKey(rulesContainerData, sig.signature, superAdminKeys);
    if (key) {
      signers.add(keyFingerprint(key));
    }
  }

  if (signers.size < minValidSignatures) {
    throw new IntegrityError(
      `insufficient distinct valid SuperAdmin signers: got ${signers.size}, need ${minValidSignatures}`
    );
  }
}

/**
 * Verifies governance rules signatures using SuperAdmin keys in PEM form.
 *
 * @param rulesContainerBase64 - Base64-encoded rules container
 * @param signaturesBase64 - Base64-encoded signatures array
 * @param minValidSignatures - Minimum number of DISTINCT valid signers required
 * @param superAdminKeysPem - Array of SuperAdmin public keys in PEM format
 * @returns true if enough distinct signers are present
 * @throws {@link ConfigurationError} If minValidSignatures is not positive
 */
export function verifyGovernanceRules(
  rulesContainerBase64: string,
  signaturesBase64: string,
  minValidSignatures: number,
  superAdminKeysPem: string[]
): boolean {
  if (minValidSignatures <= 0) {
    throw new ConfigurationError("minValidSignatures must be positive");
  }

  if (!rulesContainerBase64 || !signaturesBase64) {
    return false;
  }

  if (superAdminKeysPem.length === 0) {
    return false;
  }

  // Decode public keys
  let publicKeys: KeyObject[];
  try {
    publicKeys = decodePublicKeysPem(superAdminKeysPem);
  } catch (error: unknown) {
    if (isCryptoVerificationError(error)) {
      return false;
    }
    throw error;
  }

  // Decode signatures from base64
  let signatures: Array<{ userId?: string; signature?: string }>;
  try {
    const decoded = Buffer.from(signaturesBase64, "base64").toString("utf-8");
    const data = JSON.parse(decoded) as unknown;
    if (Array.isArray(data)) {
      signatures = data as Array<{ userId?: string; signature?: string }>;
    } else if (data && typeof data === "object" && "signatures" in data) {
      const dataObj = data as { signatures?: unknown };
      if (!Array.isArray(dataObj.signatures)) {
        // Invalid format: signatures property exists but is not an array
        return false;
      }
      signatures = dataObj.signatures as Array<{
        userId?: string;
        signature?: string;
      }>;
    } else {
      // Invalid format: not an array and no signatures property
      return false;
    }
  } catch (error: unknown) {
    if (error instanceof SyntaxError ||
        (error instanceof Error &&
         (error.message.includes('decode') ||
          error.message.includes('JSON') ||
          error.message.includes('base64')))) {
      return false;
    }
    throw error;
  }

  // The signed data is the raw rules container bytes
  const rulesData = strictBase64Decode(rulesContainerBase64);

  try {
    verifyGovernanceRulesSignatures(
      rulesData,
      signatures,
      publicKeys,
      minValidSignatures
    );
    return true;
  } catch (error: unknown) {
    if (error instanceof IntegrityError) {
      return false;
    }
    throw error;
  }
}
