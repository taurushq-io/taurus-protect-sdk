/**
 * Whitelisted asset (contract) models for Taurus-PROTECT SDK.
 *
 * This module provides domain models for whitelisted assets and
 * their cryptographic verification envelopes.
 */

import { guardSignedPayload } from "../helpers/signed-payload-guard";
import type { Verified } from "../helpers/verified";
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
  /**
   * The payload {@link verifiedHash} covers, and the bytes {@link verifiedAsset} was
   * parsed from. When a legacy variant matched, this is that variant rather than
   * `metadata.payloadAsString`.
   */
  readonly verifiedPayload: string;
  /**
   * The envelope, marked as having passed verification.
   *
   * Anything that reads raw envelope fields takes this type rather than the bare
   * envelope, so a read path that skipped `verify()` will not compile.
   */
  readonly verifiedEnvelope: Verified<SignedWhitelistedAssetEnvelope>;
}

/**
 * The asset peer of `WhitelistedAddressApproval`: the rows an approver reviewed, each
 * carrying the metadata hash it had at review time.
 *
 * Same reasoning, same guarantee — see that type. Ids are numbers here because
 * `WhitelistedAsset.id` is.
 */
export class WhitelistedAssetApproval {
  /** Row id -> the metadata hash that row carried when it was reviewed. */
  readonly #pinned: ReadonlyMap<number, string>;

  /**
   * Minted by a verified read. Prefer `result.select(...)` / `result.selectAll()` on a
   * whitelisted-asset listing over calling this directly.
   *
   * @param pinned - row id -> reviewed metadata hash
   */
  constructor(pinned: ReadonlyMap<number, string>) {
    this.#pinned = new Map(pinned);
  }

  /** The pinned row ids, in no particular order; the approval path sorts them. */
  ids(): number[] {
    return [...this.#pinned.keys()];
  }

  /**
   * The reviewed metadata hash for `id`, or `undefined` when that id was not pinned.
   *
   * @param id - the row id
   */
  pinnedHash(id: number): string | undefined {
    return this.#pinned.get(id);
  }

  /** True when this selection pins nothing. */
  isEmpty(): boolean {
    return this.#pinned.size === 0;
  }
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
 * This is step 6 of the asset flow: steps 1-5 prove the envelope is authentic, and this
 * is what stops an unsigned value reaching the caller.
 *
 * Bounded and duplicate-key-checked before parsing, exactly as
 * `parseWhitelistedAddressFromJson` is: the payload is hash-checked in step 1 but not
 * AUTHENTICATED until step 5's signatures verify, and `JSON.parse` silently keeps the
 * last of two duplicate keys.
 *
 * @param jsonPayload - JSON string of the signed payload
 * @returns Parsed WhitelistedAsset
 * @throws IntegrityError if the payload is oversized or carries a duplicate key
 * @throws Error if the payload cannot be parsed
 */
export function parseWhitelistedAssetFromJson(
  jsonPayload: string
): WhitelistedAsset {
  if (!jsonPayload) {
    throw new Error("JSON payload cannot be empty");
  }
  guardSignedPayload(jsonPayload, "whitelisted asset");

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
