/**
 * Whitelisted asset (contract) service for Taurus-PROTECT SDK.
 *
 * Provides operations for retrieving and verifying whitelisted assets.
 * All assets retrieved through this service can be automatically verified
 * for cryptographic integrity using governance rules when verification is enabled.
 */

import type { ContractWhitelistingApi } from "../internal/openapi";
import { IntegrityError, NotFoundError, ValidationError } from "../errors";
import {
  WhitelistedAssetVerifier,
  type WhitelistedAssetVerifierConfig,
  type RulesContainerDecoder,
  type UserSignaturesDecoder,
} from "../helpers/whitelisted-asset-verifier";
import type {
  WhitelistedAsset,
  SignedWhitelistedAssetEnvelope,
  WhitelistedAssetVerificationResult,
} from "../models/whitelisted-asset";
import { parseWhitelistedAssetFromJson } from "../models/whitelisted-asset";
import type { Pagination } from "../models/pagination";
import { BaseService } from "./base";

/**
 * Options for listing whitelisted assets.
 */
export interface ListWhitelistedAssetsOptions {
  /** Maximum number of items to return (max 100). */
  limit?: number;
  /** Offset for pagination. */
  offset?: number;
  /** Search query. */
  query?: string;
  /** Filter by blockchain. */
  blockchain?: string;
  /** Filter by network. */
  network?: string;
  /** Include assets pending approval. */
  includeForApproval?: boolean;
  /** Filter by contract kind types (e.g., 'nft', 'token'). */
  kindTypes?: string[];
}

/**
 * Result of listing whitelisted assets.
 */
export interface ListWhitelistedAssetsResult {
  /** List of whitelisted assets. */
  items: WhitelistedAsset[];
  /** Pagination information. */
  pagination: Pagination | undefined;
}

/**
 * Configuration for WhitelistedAssetService with verification.
 */
export interface WhitelistedAssetServiceConfig extends WhitelistedAssetVerifierConfig {
  /** Function to decode rules container from base64. */
  rulesContainerDecoder: RulesContainerDecoder;
  /** Function to decode user signatures from base64. */
  userSignaturesDecoder: UserSignaturesDecoder;
}

/**
 * Service for managing whitelisted assets (contracts).
 *
 * Provides operations for retrieving and verifying whitelisted assets.
 * When verification is enabled, all get operations will verify the
 * cryptographic signatures according to governance rules.
 *
 * @example
 * ```typescript
 * const service = WhitelistedAssetService.withVerification(api, {
 *   superAdminKeys: [key1, key2],
 *   minValidSignatures: 1,
 *   rulesContainerDecoder,
 *   userSignaturesDecoder,
 * });
 * const result = await service.getWithVerification(123);
 * ```
 */
export class WhitelistedAssetService extends BaseService {
  private readonly verifier: WhitelistedAssetVerifier;
  private readonly rulesContainerDecoder: RulesContainerDecoder;
  private readonly userSignaturesDecoder: UserSignaturesDecoder;

  /**
   * Creates a new WhitelistedAssetService with mandatory verification configuration.
   *
   * Verification is always enforced. All get/list operations use the provided
   * SuperAdmin keys to verify cryptographic integrity.
   *
   * @param api - The ContractWhitelistingApi instance from OpenAPI client
   * @param config - Verification configuration (required)
   */
  constructor(
    private readonly api: ContractWhitelistingApi,
    config: WhitelistedAssetServiceConfig
  ) {
    super();
    this.verifier = new WhitelistedAssetVerifier(config);
    this.rulesContainerDecoder = config.rulesContainerDecoder;
    this.userSignaturesDecoder = config.userSignaturesDecoder;
  }

  /**
   * Creates a WhitelistedAssetService with verification enabled.
   *
   * @param api - The ContractWhitelistingApi instance
   * @param config - Verification configuration
   * @returns WhitelistedAssetService with verification enabled
   */
  static withVerification(
    api: ContractWhitelistingApi,
    config: WhitelistedAssetServiceConfig
  ): WhitelistedAssetService {
    return new WhitelistedAssetService(api, config);
  }

  /**
   * Gets a whitelisted asset by ID.
   *
   * @param assetId - The asset ID to retrieve
   * @returns The whitelisted asset
   * @throws ValidationError if assetId is invalid
   * @throws NotFoundError if asset not found
   * @throws APIError if API request fails
   */
  async get(assetId: number): Promise<WhitelistedAsset> {
    if (assetId <= 0) {
      throw new ValidationError("assetId must be positive");
    }

    return this.execute(async () => {
      const response = await this.api.whitelistServiceGetWhitelistedContract({
        id: String(assetId),
      });

      const envelope = response.result;
      if (envelope == null) {
        throw new NotFoundError(`Whitelisted asset ${assetId} not found`);
      }

      return this.mapEnvelopeToAsset(envelope);
    });
  }

  /**
   * Gets a whitelisted asset by ID with full verification.
   *
   * Performs the complete 5-step verification of the asset.
   *
   * @param assetId - The asset ID to retrieve
   * @returns Verification result with verified asset
   * @throws ValidationError if assetId is invalid or verification not configured
   * @throws NotFoundError if asset not found
   * @throws IntegrityError if verification fails
   * @throws WhitelistError if governance thresholds not met
   * @throws APIError if API request fails
   */
  async getWithVerification(
    assetId: number
  ): Promise<WhitelistedAssetVerificationResult> {
    if (assetId <= 0) {
      throw new ValidationError("assetId must be positive");
    }

    return this.execute(async () => {
      const response = await this.api.whitelistServiceGetWhitelistedContract({
        id: String(assetId),
      });

      const dto = response.result;
      if (dto == null) {
        throw new NotFoundError(`Whitelisted asset ${assetId} not found`);
      }

      // Map DTO to envelope
      const envelope = this.mapDtoToEnvelope(dto, assetId);

      // Perform verification
      return this.verifier.verify(
        envelope,
        this.rulesContainerDecoder,
        this.userSignaturesDecoder
      );
    });
  }

  /**
   * Gets the signed envelope for a whitelisted asset.
   *
   * @param assetId - The asset ID to retrieve
   * @returns The signed envelope
   * @throws ValidationError if assetId is invalid
   * @throws NotFoundError if asset not found
   * @throws APIError if API request fails
   */
  async getEnvelope(assetId: number): Promise<SignedWhitelistedAssetEnvelope> {
    if (assetId <= 0) {
      throw new ValidationError("assetId must be positive");
    }

    return this.execute(async () => {
      const response = await this.api.whitelistServiceGetWhitelistedContract({
        id: String(assetId),
      });

      const dto = response.result;
      if (dto == null) {
        throw new NotFoundError(`Whitelisted asset ${assetId} not found`);
      }

      return this.mapDtoToEnvelope(dto, assetId);
    });
  }

  /**
   * Lists whitelisted assets.
   *
   * @param options - Optional filtering and pagination options
   * @returns Object with assets array and pagination info
   * @throws APIError if API request fails
   */
  async list(
    options?: ListWhitelistedAssetsOptions
  ): Promise<ListWhitelistedAssetsResult> {
    const limit = options?.limit ?? 50;
    const offset = options?.offset ?? 0;

    if (limit <= 0) {
      throw new ValidationError("limit must be positive");
    }
    if (offset < 0) {
      throw new ValidationError("offset cannot be negative");
    }

    return this.execute(async () => {
      const response = await this.api.whitelistServiceGetWhitelistedContracts({
        limit: String(limit),
        offset: String(offset),
        query: options?.query,
        blockchain: options?.blockchain,
        network: options?.network,
        includeForApproval: options?.includeForApproval,
        kindTypes: options?.kindTypes,
      });

      const items: WhitelistedAsset[] = [];
      if (response.result) {
        for (const envelope of response.result) {
          items.push(this.mapEnvelopeToAsset(envelope));
        }
      }

      // Extract pagination
      const totalItems = response.totalItems
        ? parseInt(response.totalItems, 10)
        : undefined;
      const pagination: Pagination | undefined = totalItems !== undefined
        ? { totalItems, offset, limit }
        : undefined;

      return { items, pagination };
    });
  }

  /**
   * Maps an envelope DTO to a WhitelistedAsset.
   *
   * SECURITY: Security-critical fields (blockchain, network, contractAddress, name, symbol)
   * are sourced ONLY from the verified payload, never from unverified DTO fields.
   * This prevents an attacker from manipulating DTO fields to bypass verification.
   *
   * @throws IntegrityError if payload is missing (security requirement)
   */
  private mapEnvelopeToAsset(
    envelope: {
      id?: string;
      metadata?: { payloadAsString?: string };
      blockchain?: string;
      network?: string;
    }
  ): WhitelistedAsset {
    // Parse from verified payload - this is the cryptographically signed data
    if (envelope.metadata?.payloadAsString) {
      const asset = parseWhitelistedAssetFromJson(envelope.metadata.payloadAsString);
      return {
        ...asset,
        id: envelope.id ? parseInt(envelope.id, 10) : 0,
        // SECURITY: blockchain and network come ONLY from verified payload
        // Do NOT use envelope.blockchain/network as they come from unverified API response
      };
    }

    // No payload means we cannot extract verified asset data - throw error
    // This is a security requirement: we must not return unverified data
    throw new IntegrityError(
      `Whitelisted asset ${envelope.id ?? 'unknown'}: metadata.payloadAsString is missing, cannot extract verified asset`
    );
  }

  /**
   * Maps a DTO to a SignedWhitelistedAssetEnvelope.
   */
  private mapDtoToEnvelope(
    dto: {
      id?: string;
      metadata?: { hash?: string; payloadAsString?: string };
      rulesContainer?: string;
      rulesSignatures?: string;
      signedContractAddress?: {
        payload?: string;
        signatures?: Array<{
          signature?: { userId?: string; signature?: string; comment?: string };
          hashes?: string[];
        }>;
      };
      blockchain?: string;
      network?: string;
    },
    assetId: number
  ): SignedWhitelistedAssetEnvelope {
    return {
      id: assetId,
      metadata: {
        hash: dto.metadata?.hash ?? "",
        payloadAsString: dto.metadata?.payloadAsString ?? "",
      },
      rulesContainerBase64: dto.rulesContainer ?? "",
      rulesSignaturesBase64: dto.rulesSignatures ?? "",
      signedContractAddress: {
        payload: dto.signedContractAddress?.payload,
        signatures: (dto.signedContractAddress?.signatures ?? []).map((sig) => ({
          userSignature: sig.signature
            ? {
                userId: sig.signature.userId,
                signature: sig.signature.signature,
                comment: sig.signature.comment,
              }
            : undefined,
          hashes: sig.hashes ?? [],
        })),
      },
      blockchain: dto.blockchain ?? "",
      network: dto.network ?? "",
    };
  }
}
