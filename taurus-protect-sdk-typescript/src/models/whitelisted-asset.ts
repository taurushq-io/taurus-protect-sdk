/**
 * Whitelisted asset (contract) models for Taurus-PROTECT SDK.
 *
 * This module provides domain models for whitelisted assets and
 * their cryptographic verification envelopes.
 */

import type {
  WhitelistMetadata,
  WhitelistSignatureEntry,
} from "./whitelisted-address";

// Re-export shared types for convenience
export type { WhitelistMetadata, WhitelistSignatureEntry };

/**
 * A whitelisted contract asset (token).
 *
 * Whitelisted assets are pre-approved tokens/contracts that can be
 * used in transactions. They must be verified with cryptographic
 * signatures before use.
 */
export interface WhitelistedAsset {
  /** Unique identifier for the whitelisted asset. */
  readonly id: number;
  /** Contract address on the blockchain. */
  readonly contractAddress: string;
  /** Blockchain identifier (e.g., "ETH", "MATIC"). */
  readonly blockchain: string;
  /** Network type (e.g., "mainnet", "testnet"). */
  readonly network: string;
  /** Human-readable name of the token/contract. */
  readonly name: string | undefined;
  /** Token symbol (e.g., "USDC", "WETH"). */
  readonly symbol: string | undefined;
  /** Number of decimal places for the token. */
  readonly decimals: number | undefined;
  /** Creation timestamp. */
  readonly createdAt: Date | undefined;
}

/**
 * Signed contract address data containing signatures.
 */
export interface SignedContractAddress {
  /** Base64-encoded signed payload. */
  readonly payload: string | undefined;
  /** List of signatures on this contract address. */
  readonly signatures: WhitelistSignatureEntry[];
}

/**
 * Envelope containing a whitelisted asset with all data needed for verification.
 *
 * This envelope contains:
 * - Metadata with hash and payload for step 1 (hash verification)
 * - Rules container and signatures for steps 2-3 (SuperAdmin verification)
 * - Signed contract address with user signatures for steps 4-5 (hash coverage and threshold verification)
 */
export interface SignedWhitelistedAssetEnvelope {
  /** Unique identifier of the whitelisted asset. */
  readonly id: number;
  /** Metadata with hash and payload. */
  readonly metadata: WhitelistMetadata;
  /** Base64-encoded rules container. */
  readonly rulesContainerBase64: string;
  /** Base64-encoded rules signatures. */
  readonly rulesSignaturesBase64: string;
  /** Signed contract address data with user signatures. */
  readonly signedContractAddress: SignedContractAddress;
  /** Blockchain identifier. */
  readonly blockchain: string;
  /** Network identifier. */
  readonly network: string;
}

/**
 * Result of whitelisted asset verification.
 */
export interface WhitelistedAssetVerificationResult {
  /** The verified whitelisted asset parsed from the payload. */
  readonly verifiedAsset: WhitelistedAsset;
  /** The hash that was verified. */
  readonly verifiedHash: string;
}

/**
 * Creates an empty WhitelistedAsset with required fields.
 *
 * @returns An empty WhitelistedAsset object
 */
export function createEmptyWhitelistedAsset(): WhitelistedAsset {
  return {
    id: 0,
    contractAddress: "",
    blockchain: "",
    network: "",
    name: undefined,
    symbol: undefined,
    decimals: undefined,
    createdAt: undefined,
  };
}

/**
 * JSON payload structure for whitelisted asset.
 */
export interface WhitelistedAssetPayload {
  readonly blockchain?: string;
  readonly network?: string;
  readonly contractAddress?: string;
  readonly name?: string;
  readonly symbol?: string;
  readonly decimals?: number;
  readonly tokenId?: string;
}

/**
 * Parses a WhitelistedAsset from a verified JSON payload.
 *
 * This extracts the signed fields from the cryptographically verified payload.
 *
 * @param jsonPayload - JSON string of the signed payload
 * @returns Parsed WhitelistedAsset
 * @throws Error if the payload cannot be parsed
 */
export function parseWhitelistedAssetFromJson(
  jsonPayload: string
): WhitelistedAsset {
  if (!jsonPayload) {
    throw new Error("JSON payload cannot be empty");
  }

  let payload: WhitelistedAssetPayload;
  try {
    payload = JSON.parse(jsonPayload) as WhitelistedAssetPayload;
  } catch (error) {
    throw new Error(
      `Failed to parse whitelist payload: ${error instanceof Error ? error.message : "unknown error"}`
    );
  }

  return {
    id: 0, // ID is not in the signed payload
    blockchain: payload.blockchain ?? "",
    network: payload.network ?? "",
    contractAddress: payload.contractAddress ?? "",
    name: payload.name,
    symbol: payload.symbol,
    decimals: payload.decimals,
    createdAt: undefined,
  };
}
