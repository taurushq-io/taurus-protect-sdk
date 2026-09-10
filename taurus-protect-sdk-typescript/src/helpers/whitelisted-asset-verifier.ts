/**
 * Whitelisted asset (contract) verification utilities.
 *
 * This module implements the complete 5-step verification flow for whitelisted assets:
 * 1. Verify metadata hash (SHA-256 of payloadAsString == metadata.hash)
 * 2. Verify rules container signatures (SuperAdmin ECDSA signatures)
 * 3. Decode rules container (base64 -> model)
 * 4. Verify hash in signed hashes list
 * 5. Verify whitelist signatures (user signatures meet governance thresholds)
 *
 * This is a CRITICAL security feature that ensures whitelisted assets have
 * been properly approved according to governance rules before use.
 *
 * Note: Unlike WhitelistedAddressVerifier, this does NOT support legacy hashes
 * since the asset schema is stable.
 */

import type { KeyObject } from "crypto";

import { calculateHexHash, verifySignature, decodePublicKeyPem } from "../crypto";
import { IntegrityError, WhitelistError } from "../errors";
import type {
  DecodedRulesContainer,
  RuleUserSignature,
  GroupThreshold,
  SequentialThresholds,
} from "../models/governance-rules";
import {
  findContractAddressWhitelistingRules,
  findUserById,
  findGroupById,
} from "../models/governance-rules";
import type {
  SignedWhitelistedAssetEnvelope,
  WhitelistedAsset,
  WhitelistedAssetVerificationResult,
} from "../models/whitelisted-asset";
import { parseWhitelistedAssetFromJson } from "../models/whitelisted-asset";
import type { WhitelistSignatureEntry } from "../models/whitelisted-address";
import { constantTimeCompare } from "./constant-time";
import { keyFingerprint, verifyGovernanceRulesSignatures } from "./signature-verifier";
import {
  computeAssetLegacyHashes,
  verifyHashCoverage,
  containsHash,
  resolveRuleKey,
} from "./whitelist-hash-helper";
import { attestVerified } from "./verified";
import { strictBase64Decode } from "./strict-base64";

/**
 * Configuration for WhitelistedAssetVerifier.
 */
export interface WhitelistedAssetVerifierConfig {
  /** SuperAdmin public keys in PEM format for rules container verification. */
  superAdminKeysPem: string[];
  /** Minimum number of valid SuperAdmin signatures required. */
  minValidSignatures: number;
}

/**
 * Type for rules container decoder function.
 */
export type RulesContainerDecoder = (base64Data: string) => DecodedRulesContainer;

/**
 * Type for user signatures decoder function.
 */
export type UserSignaturesDecoder = (base64Data: string) => RuleUserSignature[];

/**
 * Verifier for whitelisted assets (contracts).
 *
 * Performs the complete 5-step cryptographic verification to ensure
 * asset integrity and proper approval according to governance rules.
 *
 * @example
 * ```typescript
 * import { WhitelistedAssetVerifier } from '@taurushq/protect-sdk';
 * import { rulesContainerFromBase64, userSignaturesFromBase64 } from '@taurushq/protect-sdk';
 *
 * const verifier = new WhitelistedAssetVerifier({
 *   superAdminKeysPem: [superAdminKey1, superAdminKey2],
 *   minValidSignatures: 1,
 * });
 *
 * try {
 *   const result = verifier.verify(
 *     envelope,
 *     rulesContainerFromBase64,
 *     userSignaturesFromBase64
 *   );
 *   console.log('Verified asset:', result.verifiedAsset);
 * } catch (error) {
 *   if (error instanceof IntegrityError) {
 *     console.error('Integrity verification failed:', error.message);
 *   }
 * }
 * ```
 */
export class WhitelistedAssetVerifier {
  private readonly superAdminKeys: KeyObject[];
  private readonly minValidSignatures: number;

  /**
   * Creates a new WhitelistedAssetVerifier.
   *
   * @param config - Configuration with SuperAdmin keys and minimum signatures
   * @throws Error if SuperAdmin keys cannot be decoded
   */
  constructor(config: WhitelistedAssetVerifierConfig) {
    if (!config.superAdminKeysPem || config.superAdminKeysPem.length === 0) {
      throw new Error("At least one SuperAdmin key is required");
    }

    if (config.minValidSignatures < 1) {
      throw new Error("minValidSignatures must be at least 1");
    }

    this.superAdminKeys = config.superAdminKeysPem.map((pem) =>
      decodePublicKeyPem(pem)
    );
    this.minValidSignatures = config.minValidSignatures;
  }

  /**
   * Performs the complete 5-step verification of a whitelisted asset.
   *
   * Steps:
   * 1. Verify metadata hash (SHA-256 of payloadAsString == metadata.hash)
   * 2. Verify rules container signatures (SuperAdmin ECDSA signatures)
   * 3. Decode rules container (base64 -> model)
   * 4. Verify hash in signed hashes list
   * 5. Verify whitelist signatures (user signatures meet governance thresholds)
   *
   * @param envelope - The signed whitelisted asset envelope
   * @param rulesContainerDecoder - Function to decode base64 rules container
   * @param userSignaturesDecoder - Function to decode base64 user signatures
   * @returns Verification result with verified asset and hash
   * @throws IntegrityError if any verification step fails
   * @throws WhitelistError if governance thresholds are not met
   */
  // 5-step verification for a whitelisted asset, plus the parse that makes it usable.
  //
  //   envelope ──▶ 1 hash ──▶ 2 SuperAdmin sigs ──▶ 3 decode rules
  //                                                      │
  //                6 parse VERIFIED payload ◀── 5 thresholds ◀── 4 coverage
  //                           │                       │                │
  //                           ▼                       │      returns the hash it
  //                    returned to caller             │      matched, which is
  //                                                   │      what step 5 checks
  //                                                   ▼
  //                                per group: DISTINCT signers, keyed on the
  //                                container's public key, never the entry's userId
  //
  // Steps 1-5 prove the envelope is authentic; step 6 is what stops an unsigned
  // value reaching the caller. Assets use ContractAddressWhitelistingRules and the
  // ASSET legacy hashes (isNFT / kindType) — not the address ones.
  //
  // All six run here: verify() returns the parsed asset, not just a verdict.
  verify(
    envelope: SignedWhitelistedAssetEnvelope,
    rulesContainerDecoder: RulesContainerDecoder,
    userSignaturesDecoder: UserSignaturesDecoder
  ): WhitelistedAssetVerificationResult {
    if (!envelope) {
      throw new IntegrityError("envelope cannot be null or undefined");
    }
    if (!envelope.metadata) {
      throw new IntegrityError("metadata cannot be null or undefined");
    }

    // Step 1: Verify metadata hash
    this.verifyMetadataHash(envelope);

    // Step 2: Verify rules container signatures
    this.verifyRulesContainerSignatures(envelope, userSignaturesDecoder);

    // Step 3: Decode rules container
    const rulesContainer = this.decodeRulesContainer(
      envelope,
      rulesContainerDecoder
    );

    // Step 4: Verify hash in signed hashes list. Returns the hash actually
    // covered, which may be a legacy one.
    const verifiedHash = this.verifyHashInSignedHashes(envelope);

    // Step 5: Verify whitelist signatures against the hash step 4 matched
    this.verifyWhitelistSignatures(envelope, rulesContainer, verifiedHash);

    // Parse and return verified asset
    const verifiedAsset = parseWhitelistedAssetFromJson(
      envelope.metadata.payloadAsString
    );

    // Set the ID from the envelope (not in the signed payload)
    const assetWithId: WhitelistedAsset = {
      ...verifiedAsset,
      id: envelope.id,
    };

    return {
      verifiedAsset: assetWithId,
      verifiedHash,
      // The single point where the marker is applied: every step above has passed.
      verifiedEnvelope: attestVerified(envelope),
    };
  }

  /**
   * Step 1: Verify that the computed hash matches the provided hash.
   *
   * Uses constant-time comparison to prevent timing attacks.
   *
   * @param envelope - The signed whitelisted asset envelope
   * @throws IntegrityError if hash verification fails
   */
  private verifyMetadataHash(envelope: SignedWhitelistedAssetEnvelope): void {
    if (!envelope.metadata.payloadAsString) {
      throw new IntegrityError("payloadAsString is empty");
    }
    if (!envelope.metadata.hash) {
      throw new IntegrityError("metadata hash is empty");
    }

    const computedHash = calculateHexHash(envelope.metadata.payloadAsString);
    if (!constantTimeCompare(computedHash, envelope.metadata.hash)) {
      throw new IntegrityError("metadata hash verification failed");
    }
  }

  /**
   * Step 2: Verify SuperAdmin signatures on the rules container.
   *
   * @param envelope - The signed whitelisted asset envelope
   * @param userSignaturesDecoder - Function to decode user signatures
   * @throws IntegrityError if signature verification fails
   */
  private verifyRulesContainerSignatures(
    envelope: SignedWhitelistedAssetEnvelope,
    userSignaturesDecoder: UserSignaturesDecoder
  ): void {
    if (this.superAdminKeys.length === 0) {
      throw new IntegrityError(
        "no SuperAdmin keys configured for verification"
      );
    }

    if (!envelope.rulesContainerBase64) {
      throw new IntegrityError("rulesContainer is empty");
    }
    if (!envelope.rulesSignaturesBase64) {
      throw new IntegrityError("rulesSignatures is empty");
    }

    // Decode rules signatures
    let signatures: RuleUserSignature[];
    try {
      signatures = userSignaturesDecoder(envelope.rulesSignaturesBase64);
    } catch (error) {
      throw new IntegrityError(
        `failed to decode rules signatures: ${error instanceof Error ? error.message : "unknown error"}`
      );
    }

    // Decode rules container data (raw bytes)
    let rulesData: Buffer;
    try {
      rulesData = strictBase64Decode(envelope.rulesContainerBase64);
    } catch (error) {
      throw new IntegrityError(
        `failed to decode rules container: ${error instanceof Error ? error.message : "unknown error"}`
      );
    }

    try {
      verifyGovernanceRulesSignatures(
        rulesData,
        signatures,
        this.superAdminKeys,
        this.minValidSignatures
      );
    } catch (error) {
      throw new IntegrityError(
        `rules container signature verification failed: ${error instanceof Error ? error.message : "unknown error"}`
      );
    }
  }

  /**
   * Step 3: Decode the rules container from base64.
   *
   * @param envelope - The signed whitelisted asset envelope
   * @param rulesContainerDecoder - Function to decode rules container
   * @returns Decoded rules container
   * @throws IntegrityError if decoding fails
   */
  private decodeRulesContainer(
    envelope: SignedWhitelistedAssetEnvelope,
    rulesContainerDecoder: RulesContainerDecoder
  ): DecodedRulesContainer {
    try {
      return rulesContainerDecoder(envelope.rulesContainerBase64);
    } catch (error) {
      throw new IntegrityError(
        `failed to decode rules container: ${error instanceof Error ? error.message : "unknown error"}`
      );
    }
  }

  /**
   * Step 4: Verify that the metadata hash is covered by at least one signature.
   *
   * Supports legacy hashes for backward compatibility with assets signed
   * before schema changes (e.g., before isNFT or kindType was added).
   *
   * Returns the hash that was actually covered — the current one, or the legacy
   * variant that matched. Returning `void` and letting the caller re-read
   * `metadata.hash` meant a legacy-signed asset passed step 4 and then failed
   * step 5, because step 5 looked for a hash no signature covers. Go and Python
   * both carry the matched hash forward.
   *
   * @param envelope - The signed whitelisted asset envelope
   * @returns the covered hash
   * @throws IntegrityError if hash is not covered by any signature
   */
  private verifyHashInSignedHashes(
    envelope: SignedWhitelistedAssetEnvelope
  ): string {
    if (!envelope.signedContractAddress) {
      throw new IntegrityError("signedContractAddress is null or undefined");
    }

    const signatures = envelope.signedContractAddress.signatures;
    if (!signatures || signatures.length === 0) {
      throw new IntegrityError("no signatures in signedContractAddress");
    }

    const metadataHash = envelope.metadata.hash;

    // Try the provided hash first
    if (verifyHashCoverage(metadataHash, signatures)) {
      return metadataHash;
    }

    // Try legacy hashes for backward compatibility
    // This handles assets signed before schema changes (e.g., before isNFT or kindType was added)
    const legacyHashes = computeAssetLegacyHashes(
      envelope.metadata.payloadAsString
    );
    for (const legacyHash of legacyHashes) {
      if (verifyHashCoverage(legacyHash, signatures)) {
        return legacyHash;
      }
    }

    throw new IntegrityError("metadata hash is not covered by any signature");
  }

  /**
   * Step 5: Verify that user signatures meet governance threshold requirements.
   *
   * @param envelope - The signed whitelisted asset envelope
   * @param rulesContainer - The decoded rules container
   * @throws WhitelistError if thresholds are not met
   */
  private verifyWhitelistSignatures(
    envelope: SignedWhitelistedAssetEnvelope,
    rulesContainer: DecodedRulesContainer,
    metadataHash: string
  ): void {

    // Keyed off the SIGNED payload, not the response. See the address verifier.
    const { blockchain, network } = resolveRuleKey(
      envelope.metadata?.payloadAsString,
      envelope.blockchain,
      envelope.network
    );

    // Find matching contract address whitelisting rules
    const whitelistRules = findContractAddressWhitelistingRules(
      rulesContainer,
      blockchain,
      network
    );

    if (!whitelistRules) {
      throw new WhitelistError(
        `no contract address whitelisting rules found for blockchain=${blockchain} network=${network}`
      );
    }

    // Contract whitelisting uses parallelThresholds directly (no rule lines matching)
    const parallelThresholds = whitelistRules.parallelThresholds;
    if (!parallelThresholds || parallelThresholds.length === 0) {
      throw new WhitelistError("no threshold rules defined");
    }

    // Try to verify all paths (OR logic - only one needs to succeed)
    const pathFailures = this.tryVerifyAllPaths(
      parallelThresholds,
      rulesContainer,
      envelope.signedContractAddress.signatures,
      metadataHash
    );

    if (pathFailures.length > 0) {
      throw new WhitelistError(
        `signature verification failed for whitelisted asset (ID: ${envelope.id}): ` +
          `no approval path satisfied the threshold requirements. ${pathFailures.join("; ")}`
      );
    }
  }

  /**
   * Tries to verify all parallel threshold paths.
   *
   * @param parallelThresholds - List of sequential thresholds (OR logic between them)
   * @param rulesContainer - The decoded rules container
   * @param signatures - User signatures on the asset
   * @param metadataHash - The verified metadata hash
   * @returns Empty array if verification passed, or list of failure messages
   */
  private tryVerifyAllPaths(
    parallelThresholds: SequentialThresholds[],
    rulesContainer: DecodedRulesContainer,
    signatures: WhitelistSignatureEntry[],
    metadataHash: string
  ): string[] {
    // Pre-compute JSON serialization of each signature's hashes array once,
    // so it can be reused across all group threshold checks.
    const precomputedHashesJson = new Map<number, Buffer>();
    for (let idx = 0; idx < signatures.length; idx++) {
      const sig = signatures[idx]!;
      if (sig.hashes) {
        precomputedHashesJson.set(
          idx,
          Buffer.from(JSON.stringify(sig.hashes), "utf-8")
        );
      }
    }

    // Cache for decoded public keys (keyed by PEM string) to avoid
    // repeated PEM parsing across group threshold checks.
    const keyCache = new Map<string, KeyObject>();

    const pathFailures: string[] = [];

    for (let i = 0; i < parallelThresholds.length; i++) {
      const seqThreshold = parallelThresholds[i];
      const error = this.verifySequentialThresholds(
        seqThreshold!,
        rulesContainer,
        signatures,
        metadataHash,
        precomputedHashesJson,
        keyCache
      );

      if (error === null) {
        return []; // Verification passed
      }
      pathFailures.push(`Path ${i + 1}: ${error}`);
    }

    return pathFailures;
  }

  /**
   * Verifies all group thresholds in a sequential threshold path.
   *
   * All group thresholds must be satisfied (AND logic).
   *
   * @param seqThreshold - Sequential thresholds to verify
   * @param rulesContainer - The decoded rules container
   * @param signatures - User signatures
   * @param metadataHash - The verified metadata hash
   * @returns null if successful, or error message on failure
   */
  private verifySequentialThresholds(
    seqThreshold: SequentialThresholds,
    rulesContainer: DecodedRulesContainer,
    signatures: WhitelistSignatureEntry[],
    metadataHash: string,
    precomputedHashesJson: Map<number, Buffer>,
    keyCache: Map<string, KeyObject>
  ): string | null {
    if (!seqThreshold || !seqThreshold.thresholds || seqThreshold.thresholds.length === 0) {
      return "no group thresholds defined";
    }

    // ALL group thresholds must be satisfied (AND logic)
    for (const groupThreshold of seqThreshold.thresholds) {
      const error = this.verifyGroupThreshold(
        groupThreshold,
        rulesContainer,
        signatures,
        metadataHash,
        precomputedHashesJson,
        keyCache
      );
      if (error !== null) {
        return error;
      }
    }

    return null;
  }

  /**
   * Verifies that a group threshold is met.
   *
   * @param groupThreshold - The group threshold requirement
   * @param rulesContainer - The decoded rules container
   * @param signatures - User signatures
   * @param metadataHash - The verified metadata hash
   * @returns null if successful, or error message on failure
   */
  private verifyGroupThreshold(
    groupThreshold: GroupThreshold,
    rulesContainer: DecodedRulesContainer,
    signatures: WhitelistSignatureEntry[],
    metadataHash: string,
    precomputedHashesJson: Map<number, Buffer>,
    keyCache: Map<string, KeyObject>
  ): string | null {
    const groupId = groupThreshold.groupId;
    // Prefer minimumSignatures, fall back to threshold
    const minSigs =
      groupThreshold.minimumSignatures > 0
        ? groupThreshold.minimumSignatures
        : groupThreshold.threshold;

    if (!groupId) {
      return "group threshold has no groupId";
    }

    const group = findGroupById(rulesContainer, groupId);
    if (!group) {
      return `group '${groupId}' not found in rules container`;
    }

    if (!group.userIds || group.userIds.length === 0) {
      if (minSigs > 0) {
        return `group '${groupId}' has no users but requires ${minSigs} signature(s)`;
      }
      return null; // minSignatures == 0, so empty group is OK
    }

    // A populated group with a zero threshold is a malformed container, not a
    // group anyone may satisfy. There is no post-loop threshold check — the only
    // success exit is inside the loop after an increment — so a zero here
    // silently means "one signature suffices", turning a 2-of-N group into
    // 1-of-N. Fail closed.
    if (minSigs <= 0) {
      return `group '${groupId}' has ${group.userIds.length} user(s) but requires 0 signature(s): minimumSignatures must be positive`;
    }

    // Build set for faster lookup
    const groupUserIdSet = new Set(group.userIds);

    // Count valid signatures from users in this group
    // Count DISTINCT signers, not signature entries — see the address verifier.
    const signers = new Set<string>();
    const skippedReasons: string[] = [];

    for (let sigIdx = 0; sigIdx < signatures.length; sigIdx++) {
      const sig = signatures[sigIdx]!;
      if (!sig.userSignature) {
        skippedReasons.push("signature has nil userSig");
        continue;
      }

      const sigUserId = sig.userSignature.userId;
      if (!sigUserId || !groupUserIdSet.has(sigUserId)) {
        continue; // Signer not in this group - not an error, just not relevant
      }

      // Check that metadata hash is covered by this signature
      if (!containsHash(sig.hashes, metadataHash)) {
        skippedReasons.push(
          `user '${sigUserId}' signature does not cover metadata hash '${metadataHash}' (signed hashes=${JSON.stringify(sig.hashes)})`
        );
        continue;
      }

      const user = findUserById(rulesContainer, sigUserId);
      if (!user) {
        skippedReasons.push(
          `user '${sigUserId}' not found in rules container`
        );
        continue;
      }
      if (!user.publicKeyPem) {
        skippedReasons.push(`user '${sigUserId}' has no public key`);
        continue;
      }

      // Decode user's public key (cached across group threshold checks)
      let publicKey: KeyObject | undefined = keyCache.get(user.publicKeyPem);
      if (!publicKey) {
        try {
          publicKey = decodePublicKeyPem(user.publicKeyPem);
          keyCache.set(user.publicKeyPem, publicKey);
        } catch (error) {
          skippedReasons.push(
            `failed to decode public key for user '${sigUserId}': ${error instanceof Error ? error.message : "unknown error"}`
          );
          continue;
        }
      }

      // Use pre-computed JSON-encoded hashes buffer
      const hashesData = precomputedHashesJson.get(sigIdx);
      if (!hashesData) {
        skippedReasons.push(`user '${sigUserId}' has no hashes`);
        continue;
      }

      if (!sig.userSignature.signature) {
        skippedReasons.push(`user '${sigUserId}' signature is empty`);
        continue;
      }

      try {
        if (verifySignature(publicKey, hashesData, sig.userSignature.signature)) {
          signers.add(keyFingerprint(publicKey));
          if (signers.size >= minSigs) {
            return null; // Threshold met
          }
        } else {
          skippedReasons.push(
            `user '${sigUserId}' signature verification failed`
          );
        }
      } catch (error) {
        skippedReasons.push(
          `user '${sigUserId}' signature verification error: ${error instanceof Error ? error.message : "unknown error"}`
        );
      }
    }

    // Threshold not met
    let message = `group '${groupId}' requires ${minSigs} distinct signer(s) but only ${signers.size} valid`;
    if (skippedReasons.length > 0) {
      message += ` [${skippedReasons.join("; ")}]`;
    }
    return message;
  }

}
