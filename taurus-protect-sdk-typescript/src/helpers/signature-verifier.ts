/**
 * Signature verification utilities for governance rules.
 *
 * This module provides functions to verify SuperAdmin signatures on
 * governance rules and other cryptographically protected data.
 */

import type { KeyObject } from "crypto";

import { verifySignature, decodePublicKeysPem } from "../crypto";

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
  const dataBuffer =
    typeof data === "string" ? Buffer.from(data, "utf-8") : data;

  for (const publicKey of publicKeys) {
    try {
      if (verifySignature(publicKey, dataBuffer, signatureBase64)) {
        return true;
      }
    } catch (error: unknown) {
      // Only swallow expected crypto/signature errors
      if (error instanceof Error &&
          (error.message.includes('signature') ||
           error.message.includes('key') ||
           error.message.includes('Invalid') ||
           error.message.includes('decode') ||
           error.message.includes('ERR_OSSL'))) {
        continue;
      }
      throw error;
    }
  }

  return false;
}

/**
 * Verifies governance rules signatures using SuperAdmin keys.
 *
 * @param rulesContainerBase64 - Base64-encoded rules container
 * @param signaturesBase64 - Base64-encoded signatures array
 * @param minValidSignatures - Minimum number of valid signatures required
 * @param superAdminKeysPem - Array of SuperAdmin public keys in PEM format
 * @returns true if enough valid signatures are present
 */
export function verifyGovernanceRules(
  rulesContainerBase64: string,
  signaturesBase64: string,
  minValidSignatures: number,
  superAdminKeysPem: string[]
): boolean {
  if (minValidSignatures <= 0) {
    return true;
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
    if (error instanceof Error &&
        (error.message.includes('key') ||
         error.message.includes('Invalid') ||
         error.message.includes('decode') ||
         error.message.includes('ERR_OSSL') ||
         error.message.includes('PEM'))) {
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
  const rulesData = Buffer.from(rulesContainerBase64, "base64");

  // Count valid signatures
  let validCount = 0;
  const seenUserIds = new Set<string>();

  for (const sig of signatures) {
    if (!sig.signature) continue;

    // Ensure distinct user signatures
    if (sig.userId && seenUserIds.has(sig.userId)) {
      continue;
    }

    if (isValidSignature(rulesData, sig.signature, publicKeys)) {
      validCount++;
      if (sig.userId) {
        seenUserIds.add(sig.userId);
      }

      // Early return once threshold is met
      if (validCount >= minValidSignatures) {
        return true;
      }
    }
  }

  return validCount >= minValidSignatures;
}
