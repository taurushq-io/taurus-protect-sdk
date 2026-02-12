/**
 * Whitelisted address verification utilities.
 *
 * This module implements the complete 6-step verification flow for whitelisted addresses:
 * 1. Verify metadata hash (SHA-256 of payloadAsString == metadata.hash)
 * 2. Verify rules container signatures (SuperAdmin ECDSA signatures)
 * 3. Decode rules container (base64 -> model)
 * 4. Verify hash in signed hashes list (with legacy hash support)
 * 5. Verify whitelist signatures (user signatures meet governance thresholds)
 * 6. Parse and return verified WhitelistedAddress from payload
 *
 * This is a CRITICAL security feature that ensures whitelisted addresses have
 * been properly approved according to governance rules before use.
 */

import type { KeyObject } from "crypto";

import { calculateHexHash, verifySignature, decodePublicKeyPem } from "../crypto";
import { IntegrityError, WhitelistError } from "../errors";
import type {
  DecodedRulesContainer,
  RuleUserSignature,
  AddressWhitelistingRules,
  GroupThreshold,
  SequentialThresholds,
} from "../models/governance-rules";
import {
  findAddressWhitelistingRules,
  findUserById,
  findGroupById,
} from "../models/governance-rules";
import type {
  SignedWhitelistedAddressEnvelope,
  WhitelistedAddress,
  WhitelistedAddressVerificationResult,
  WhitelistSignatureEntry,
  InternalWallet,
} from "../models/whitelisted-address";
import { constantTimeCompare } from "./constant-time";
import { isValidSignature } from "./signature-verifier";
import {
  computeLegacyHashes,
  parseWhitelistedAddressFromJson,
  verifyHashCoverage,
} from "./whitelist-hash-helper";

/**
 * Configuration for WhitelistedAddressVerifier.
 */
export interface WhitelistedAddressVerifierConfig {
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
 * Address whitelisting line for rule matching.
 */
interface AddressWhitelistingLine {
  cells: Array<{
    type: string;
    internalWallet?: {
      path?: string;
    };
  }>;
  parallelThresholds: SequentialThresholds[];
}

/**
 * Extended address whitelisting rules with lines.
 */
interface ExtendedAddressWhitelistingRules extends AddressWhitelistingRules {
  lines?: AddressWhitelistingLine[];
}

/**
 * Verifier for whitelisted addresses.
 *
 * Performs the complete 6-step cryptographic verification to ensure
 * address integrity and proper approval according to governance rules.
 *
 * @example
 * ```typescript
 * import { WhitelistedAddressVerifier } from '@taurushq/protect-sdk';
 * import { rulesContainerFromBase64, userSignaturesFromBase64 } from '@taurushq/protect-sdk';
 *
 * const verifier = new WhitelistedAddressVerifier({
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
 *   console.log('Verified address:', result.verifiedWhitelistedAddress);
 * } catch (error) {
 *   if (error instanceof IntegrityError) {
 *     console.error('Integrity verification failed:', error.message);
 *   }
 * }
 * ```
 */
export class WhitelistedAddressVerifier {
  private readonly superAdminKeys: KeyObject[];
  private readonly minValidSignatures: number;

  /**
   * Creates a new WhitelistedAddressVerifier.
   *
   * @param config - Configuration with SuperAdmin keys and minimum signatures
   * @throws Error if SuperAdmin keys cannot be decoded
   */
  constructor(config: WhitelistedAddressVerifierConfig) {
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
   * Performs the complete 6-step verification of a whitelisted address.
   *
   * Steps:
   * 1. Verify metadata hash (SHA-256 of payloadAsString == metadata.hash)
   * 2. Verify rules container signatures (SuperAdmin ECDSA signatures)
   * 3. Decode rules container (base64 -> model)
   * 4. Verify hash in signed hashes list (with legacy hash support)
   * 5. Verify whitelist signatures (user signatures meet governance thresholds)
   * 6. Parse and return whitelisted address from payload
   *
   * @param envelope - The signed whitelisted address envelope
   * @param rulesContainerDecoder - Function to decode base64 rules container
   * @param userSignaturesDecoder - Function to decode base64 user signatures
   * @returns Verification result with verified address and hash
   * @throws IntegrityError if any verification step fails
   * @throws WhitelistError if governance thresholds are not met
   */
  /**
   * Performs the complete 6-step verification of a whitelisted address.
   *
   * @param envelope - The signed whitelisted address envelope
   * @param rulesContainerDecoder - Function to decode base64 rules container
   * @param userSignaturesDecoder - Function to decode base64 user signatures
   * @param cachedRulesContainer - Pre-verified and decoded rules container.
   *   When provided, steps 2-3 are skipped (already done during cache building).
   * @returns Verification result with verified address and hash
   * @throws IntegrityError if any verification step fails
   * @throws WhitelistError if governance thresholds are not met
   */
  verify(
    envelope: SignedWhitelistedAddressEnvelope,
    rulesContainerDecoder: RulesContainerDecoder,
    userSignaturesDecoder: UserSignaturesDecoder,
    cachedRulesContainer?: DecodedRulesContainer
  ): WhitelistedAddressVerificationResult {
    if (!envelope) {
      throw new IntegrityError("envelope cannot be null or undefined");
    }
    if (!envelope.metadata) {
      throw new IntegrityError("metadata cannot be null or undefined");
    }

    // Step 1: Verify metadata hash
    this.verifyMetadataHash(envelope);

    let rulesContainer: DecodedRulesContainer;
    if (cachedRulesContainer) {
      // Steps 2-3 already done during cache building
      rulesContainer = cachedRulesContainer;
    } else {
      // Step 2: Verify rules container signatures
      this.verifyRulesContainerSignatures(envelope, userSignaturesDecoder);

      // Step 3: Decode rules container
      rulesContainer = this.decodeRulesContainer(
        envelope,
        rulesContainerDecoder
      );
    }

    // Step 4: Verify hash in signed hashes list (with legacy support)
    const verifiedHash = this.verifyHashInSignedHashes(envelope);

    // Step 5: Verify whitelist signatures
    this.verifyWhitelistSignatures(envelope, rulesContainer, verifiedHash);

    // Step 6: Parse and return verified address
    const verifiedAddress = parseWhitelistedAddressFromJson(
      envelope.metadata.payloadAsString
    );

    // Set the ID from the envelope (not in the signed payload)
    const addressWithId: WhitelistedAddress = {
      ...verifiedAddress,
      id: envelope.id,
    };

    return {
      verifiedWhitelistedAddress: addressWithId,
      verifiedRulesContainer: rulesContainer,
      verifiedHash,
    };
  }

  /**
   * Verifies SuperAdmin signatures on a rules container and decodes it.
   * This performs steps 2-3 of the verification flow for a single rules container,
   * used by the service to build the normalized cache.
   *
   * @param rulesContainerBase64 - Base64-encoded rules container
   * @param rulesSignaturesBase64 - Base64-encoded rules signatures
   * @param rulesContainerDecoder - Function to decode rules container
   * @param userSignaturesDecoder - Function to decode user signatures
   * @returns Decoded rules container
   * @throws IntegrityError if verification or decoding fails
   */
  verifyAndDecodeRulesContainer(
    rulesContainerBase64: string,
    rulesSignaturesBase64: string,
    rulesContainerDecoder: RulesContainerDecoder,
    userSignaturesDecoder: UserSignaturesDecoder
  ): DecodedRulesContainer {
    if (!rulesContainerBase64) {
      throw new IntegrityError("rulesContainer is empty");
    }
    if (!rulesSignaturesBase64) {
      throw new IntegrityError("rulesSignatures is empty");
    }

    // Decode rules signatures
    let signatures: RuleUserSignature[];
    try {
      signatures = userSignaturesDecoder(rulesSignaturesBase64);
    } catch (error) {
      throw new IntegrityError(
        `failed to decode rules signatures: ${error instanceof Error ? error.message : "unknown error"}`
      );
    }

    // Decode rules container data (raw bytes)
    let rulesData: Buffer;
    try {
      rulesData = Buffer.from(rulesContainerBase64, "base64");
    } catch (error) {
      throw new IntegrityError(
        `failed to decode rules container: ${error instanceof Error ? error.message : "unknown error"}`
      );
    }

    // Count valid signatures
    let validCount = 0;
    for (const sig of signatures) {
      if (
        sig.signature &&
        isValidSignature(rulesData, sig.signature, this.superAdminKeys)
      ) {
        validCount++;
      }
    }

    if (validCount < this.minValidSignatures) {
      throw new IntegrityError(
        `rules container signature verification failed: only ${validCount} valid signatures, ` +
          `minimum ${this.minValidSignatures} required`
      );
    }

    // Decode rules container
    try {
      return rulesContainerDecoder(rulesContainerBase64);
    } catch (error) {
      throw new IntegrityError(
        `failed to decode rules container: ${error instanceof Error ? error.message : "unknown error"}`
      );
    }
  }

  /**
   * Step 1: Verify that the computed hash matches the provided hash.
   *
   * Uses constant-time comparison to prevent timing attacks.
   *
   * @param envelope - The signed whitelisted address envelope
   * @throws IntegrityError if hash verification fails
   */
  private verifyMetadataHash(envelope: SignedWhitelistedAddressEnvelope): void {
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
   * @param envelope - The signed whitelisted address envelope
   * @param userSignaturesDecoder - Function to decode user signatures
   * @throws IntegrityError if signature verification fails
   */
  private verifyRulesContainerSignatures(
    envelope: SignedWhitelistedAddressEnvelope,
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
      rulesData = Buffer.from(envelope.rulesContainerBase64, "base64");
    } catch (error) {
      throw new IntegrityError(
        `failed to decode rules container: ${error instanceof Error ? error.message : "unknown error"}`
      );
    }

    // Count valid signatures
    let validCount = 0;
    for (const sig of signatures) {
      if (
        sig.signature &&
        isValidSignature(rulesData, sig.signature, this.superAdminKeys)
      ) {
        validCount++;
      }
    }

    if (validCount < this.minValidSignatures) {
      throw new IntegrityError(
        `rules container signature verification failed: only ${validCount} valid signatures, ` +
          `minimum ${this.minValidSignatures} required`
      );
    }
  }

  /**
   * Step 3: Decode the rules container from base64.
   *
   * @param envelope - The signed whitelisted address envelope
   * @param rulesContainerDecoder - Function to decode rules container
   * @returns Decoded rules container
   * @throws IntegrityError if decoding fails
   */
  private decodeRulesContainer(
    envelope: SignedWhitelistedAddressEnvelope,
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
   * This step also tries legacy hashes for backward compatibility with
   * addresses signed before schema changes.
   *
   * @param envelope - The signed whitelisted address envelope
   * @returns The hash that was found (may be a legacy hash)
   * @throws IntegrityError if hash is not covered by any signature
   */
  private verifyHashInSignedHashes(
    envelope: SignedWhitelistedAddressEnvelope
  ): string {
    if (!envelope.signedAddress) {
      throw new IntegrityError("signedAddress is null or undefined");
    }

    const signatures = envelope.signedAddress.signatures;
    if (!signatures || signatures.length === 0) {
      throw new IntegrityError("no signatures in signedAddress");
    }

    // Try the provided hash first
    const providedHash = envelope.metadata.hash;
    if (verifyHashCoverage(providedHash, signatures)) {
      return providedHash;
    }

    // Try legacy hashes for backward compatibility
    const legacyHashes = computeLegacyHashes(envelope.metadata.payloadAsString);
    for (const legacyHash of legacyHashes) {
      if (verifyHashCoverage(legacyHash, signatures)) {
        return legacyHash;
      }
    }

    throw new IntegrityError(
      `metadata hash '${providedHash}' is not covered by any signature`
    );
  }

  /**
   * Step 5: Verify that user signatures meet governance threshold requirements.
   *
   * @param envelope - The signed whitelisted address envelope
   * @param rulesContainer - The decoded rules container
   * @param metadataHash - The verified metadata hash
   * @throws WhitelistError if thresholds are not met
   */
  private verifyWhitelistSignatures(
    envelope: SignedWhitelistedAddressEnvelope,
    rulesContainer: DecodedRulesContainer,
    metadataHash: string
  ): void {
    // Find matching address whitelisting rules
    const whitelistRules = findAddressWhitelistingRules(
      rulesContainer,
      envelope.blockchain,
      envelope.network
    ) as ExtendedAddressWhitelistingRules | undefined;

    if (!whitelistRules) {
      throw new WhitelistError(
        `no address whitelisting rules found for blockchain=${envelope.blockchain} network=${envelope.network}`
      );
    }

    // Determine which thresholds to use based on rule lines
    const parallelThresholds = this.getApplicableThresholds(
      whitelistRules,
      envelope
    );
    if (!parallelThresholds || parallelThresholds.length === 0) {
      throw new WhitelistError("no threshold rules defined");
    }

    // Try to verify all paths (OR logic - only one needs to succeed)
    const pathFailures = this.tryVerifyAllPaths(
      parallelThresholds,
      rulesContainer,
      envelope.signedAddress.signatures,
      metadataHash
    );

    if (pathFailures.length > 0) {
      throw new WhitelistError(
        `signature verification failed for whitelisted address (ID: ${envelope.id}): ` +
          `no approval path satisfied the threshold requirements. ${pathFailures.join("; ")}`
      );
    }
  }

  /**
   * Determines which thresholds to use based on rule lines.
   *
   * Checks rule lines only when: NO linked addresses AND exactly 1 linked wallet.
   *
   * @param rules - Address whitelisting rules
   * @param envelope - The envelope with linked addresses/wallets
   * @returns Applicable thresholds
   */
  private getApplicableThresholds(
    rules: ExtendedAddressWhitelistingRules,
    envelope: SignedWhitelistedAddressEnvelope
  ): SequentialThresholds[] {
    const hasLinkedAddresses = envelope.linkedInternalAddresses.length > 0;
    const walletCount = envelope.linkedWallets.length;

    // Check rule lines only if: no linked addresses AND exactly 1 linked wallet
    const shouldCheckRuleLines = !hasLinkedAddresses && walletCount === 1;

    if (shouldCheckRuleLines && rules.lines && rules.lines.length > 0) {
      const walletPath = envelope.linkedWallets[0]?.path;

      // Find matching line by wallet path
      for (const line of rules.lines) {
        if (this.matchesWalletPath(line, walletPath)) {
          return line.parallelThresholds;
        }
      }
    }

    // Fallback to default thresholds
    return rules.parallelThresholds;
  }

  /**
   * Checks if a rule line matches the given wallet path.
   *
   * @param line - The rule line to check
   * @param walletPath - The wallet path to match
   * @returns true if the line matches
   */
  private matchesWalletPath(
    line: AddressWhitelistingLine,
    walletPath: string | undefined
  ): boolean {
    if (!line.cells || line.cells.length === 0) {
      return false;
    }

    const source = line.cells[0];
    if (!source || source.type !== "INTERNAL_WALLET") {
      return false;
    }

    if (!source.internalWallet) {
      return false;
    }

    return !!walletPath && walletPath === source.internalWallet.path;
  }

  /**
   * Tries to verify all parallel threshold paths.
   *
   * @param parallelThresholds - List of sequential thresholds (OR logic between them)
   * @param rulesContainer - The decoded rules container
   * @param signatures - User signatures on the address
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

    // Build set for faster lookup
    const groupUserIdSet = new Set(group.userIds);

    // Count valid signatures from users in this group
    let validCount = 0;
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
      if (!this.containsHash(sig.hashes, metadataHash)) {
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
          validCount++;
          if (validCount >= minSigs) {
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
    let message = `group '${groupId}' requires ${minSigs} signature(s) but only ${validCount} valid`;
    if (skippedReasons.length > 0) {
      message += ` [${skippedReasons.join("; ")}]`;
    }
    return message;
  }

  /**
   * Checks if a hash is in the list using constant-time comparison.
   *
   * @param hashes - List of hashes to search
   * @param hash - Hash to find
   * @returns true if the hash is found
   */
  private containsHash(hashes: string[], hash: string): boolean {
    let found = false;
    for (const h of hashes) {
      if (constantTimeCompare(h, hash)) {
        found = true;
      }
    }
    return found;
  }
}
