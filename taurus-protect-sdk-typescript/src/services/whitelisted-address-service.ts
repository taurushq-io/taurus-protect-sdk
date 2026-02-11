/**
 * Whitelisted address service for Taurus-PROTECT SDK.
 *
 * Provides operations for retrieving and verifying whitelisted addresses.
 * All addresses retrieved through this service can be automatically verified
 * for cryptographic integrity using governance rules when verification is enabled.
 */

import type { AddressWhitelistingApi } from "../internal/openapi";
import { IntegrityError, NotFoundError, ValidationError } from "../errors";
import {
  WhitelistedAddressVerifier,
  type WhitelistedAddressVerifierConfig,
  type RulesContainerDecoder,
  type UserSignaturesDecoder,
} from "../helpers/whitelisted-address-verifier";
import type { DecodedRulesContainer } from "../models/governance-rules";
import type {
  WhitelistedAddress,
  SignedWhitelistedAddressEnvelope,
  WhitelistedAddressVerificationResult,
  InternalAddress,
  InternalWallet,
} from "../models/whitelisted-address";
import type { Pagination } from "../models/pagination";
import { BaseService } from "./base";

/**
 * Options for listing whitelisted addresses.
 */
export interface ListWhitelistedAddressesOptions {
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
  /** Filter by address type. */
  addressType?: string;
  /** Include addresses pending approval. */
  includeForApproval?: boolean;
  /** Filter by specific addresses. */
  addresses?: string[];
  /** Filter by tag IDs. */
  tagIds?: string[];
  /** Filter by contract types. */
  contractTypes?: string[];
}

/**
 * Result of listing whitelisted addresses.
 */
export interface ListWhitelistedAddressesResult {
  /** List of whitelisted addresses. */
  items: WhitelistedAddress[];
  /** Pagination information. */
  pagination: Pagination | undefined;
}

/**
 * Configuration for WhitelistedAddressService with verification.
 */
export interface WhitelistedAddressServiceConfig extends WhitelistedAddressVerifierConfig {
  /** Function to decode rules container from base64. */
  rulesContainerDecoder: RulesContainerDecoder;
  /** Function to decode user signatures from base64. */
  userSignaturesDecoder: UserSignaturesDecoder;
}

/**
 * Service for managing whitelisted addresses.
 *
 * Provides operations for retrieving and verifying whitelisted addresses.
 * When verification is enabled, all get operations will verify the
 * cryptographic signatures according to governance rules.
 *
 * @example
 * ```typescript
 * const service = WhitelistedAddressService.withVerification(api, {
 *   superAdminKeys: [key1, key2],
 *   minValidSignatures: 1,
 *   rulesContainerDecoder,
 *   userSignaturesDecoder,
 * });
 * const result = await service.getWithVerification('123');
 * ```
 */
export class WhitelistedAddressService extends BaseService {
  private readonly verifier: WhitelistedAddressVerifier;
  private readonly rulesContainerDecoder: RulesContainerDecoder;
  private readonly userSignaturesDecoder: UserSignaturesDecoder;

  /**
   * Creates a new WhitelistedAddressService with mandatory verification configuration.
   *
   * Verification is always enforced. All get/list operations use the provided
   * SuperAdmin keys to verify cryptographic integrity.
   *
   * @param api - The AddressWhitelistingApi instance from OpenAPI client
   * @param config - Verification configuration (required)
   */
  constructor(
    private readonly api: AddressWhitelistingApi,
    config: WhitelistedAddressServiceConfig
  ) {
    super();
    this.verifier = new WhitelistedAddressVerifier(config);
    this.rulesContainerDecoder = config.rulesContainerDecoder;
    this.userSignaturesDecoder = config.userSignaturesDecoder;
  }

  /**
   * Creates a WhitelistedAddressService with verification enabled.
   *
   * @param api - The AddressWhitelistingApi instance
   * @param config - Verification configuration
   * @returns WhitelistedAddressService with verification enabled
   */
  static withVerification(
    api: AddressWhitelistingApi,
    config: WhitelistedAddressServiceConfig
  ): WhitelistedAddressService {
    return new WhitelistedAddressService(api, config);
  }

  /**
   * Gets a whitelisted address by ID with mandatory verification.
   *
   * Performs the full 6-step verification of the address, matching
   * Java/Python/Go behavior where verification is always enforced.
   *
   * @param addressId - The address ID to retrieve
   * @returns The verified whitelisted address
   * @throws ValidationError if addressId is invalid
   * @throws NotFoundError if address not found
   * @throws IntegrityError if verification fails
   * @throws APIError if API request fails
   */
  async get(addressId: string): Promise<WhitelistedAddress> {
    if (!addressId) {
      throw new ValidationError("addressId is required");
    }

    return this.execute(async () => {
      const response = await this.api.whitelistServiceGetWhitelistedAddress({
        id: addressId,
      });

      const dto = response.result;
      if (dto == null) {
        throw new NotFoundError(`Whitelisted address ${addressId} not found`);
      }

      // Map to envelope and verify (6-step verification)
      const envelope = this.mapDtoToEnvelope(dto, addressId);
      const result = this.verifier.verify(
        envelope,
        this.rulesContainerDecoder,
        this.userSignaturesDecoder
      );

      return result.verifiedWhitelistedAddress;
    });
  }

  /**
   * Gets a whitelisted address by ID with full verification.
   *
   * Performs the complete 6-step verification of the address.
   *
   * @param addressId - The address ID to retrieve
   * @returns Verification result with verified address
   * @throws ValidationError if addressId is invalid or verification not configured
   * @throws NotFoundError if address not found
   * @throws IntegrityError if verification fails
   * @throws WhitelistError if governance thresholds not met
   * @throws APIError if API request fails
   */
  async getWithVerification(
    addressId: string
  ): Promise<WhitelistedAddressVerificationResult> {
    if (!addressId) {
      throw new ValidationError("addressId is required");
    }

    return this.execute(async () => {
      const response = await this.api.whitelistServiceGetWhitelistedAddress({
        id: addressId,
      });

      const dto = response.result;
      if (dto == null) {
        throw new NotFoundError(`Whitelisted address ${addressId} not found`);
      }

      // Map DTO to envelope
      const envelope = this.mapDtoToEnvelope(dto, addressId);

      // Perform verification
      return this.verifier.verify(
        envelope,
        this.rulesContainerDecoder,
        this.userSignaturesDecoder
      );
    });
  }

  /**
   * Gets the signed envelope for a whitelisted address.
   *
   * @param addressId - The address ID to retrieve
   * @returns The signed envelope
   * @throws ValidationError if addressId is invalid
   * @throws NotFoundError if address not found
   * @throws APIError if API request fails
   */
  async getEnvelope(addressId: string): Promise<SignedWhitelistedAddressEnvelope> {
    if (!addressId) {
      throw new ValidationError("addressId is required");
    }

    return this.execute(async () => {
      const response = await this.api.whitelistServiceGetWhitelistedAddress({
        id: addressId,
      });

      const dto = response.result;
      if (dto == null) {
        throw new NotFoundError(`Whitelisted address ${addressId} not found`);
      }

      return this.mapDtoToEnvelope(dto, addressId);
    });
  }

  /**
   * Lists whitelisted addresses with mandatory verification.
   *
   * Each address undergoes the full 6-step verification, matching
   * Java/Python/Go behavior where verification is always enforced.
   * Fails on the first verification error (strict mode).
   *
   * @param options - Optional filtering and pagination options
   * @returns Object with verified addresses array and pagination info
   * @throws IntegrityError if verification fails for any address
   * @throws WhitelistError if governance thresholds not met
   * @throws APIError if API request fails
   */
  async list(
    options?: ListWhitelistedAddressesOptions
  ): Promise<ListWhitelistedAddressesResult> {
    const limit = options?.limit ?? 50;
    const offset = options?.offset ?? 0;

    if (limit <= 0) {
      throw new ValidationError("limit must be positive");
    }
    if (offset < 0) {
      throw new ValidationError("offset cannot be negative");
    }

    return this.execute(async () => {
      const response = await this.api.whitelistServiceGetWhitelistedAddresses({
        limit: String(limit),
        offset: String(offset),
        query: options?.query,
        blockchain: options?.blockchain,
        network: options?.network,
        addressType: options?.addressType,
        includeForApproval: options?.includeForApproval,
        addresses: options?.addresses,
        tagIDs: options?.tagIds,
        contractTypes: options?.contractTypes,
        rulesContainerNormalized: true,
      });

      // Build rules container cache from normalized response
      const rulesContainerCache = this.buildRulesContainerCache(response);

      const items: WhitelistedAddress[] = [];
      if (response.result) {
        for (const dto of response.result) {
          const dtoId = dto.id ?? "";
          const envelope = this.mapDtoToEnvelope(dto, dtoId);

          // Look up cached rules container by hash
          const cached = envelope.rulesContainerHash
            ? rulesContainerCache.get(envelope.rulesContainerHash)
            : undefined;

          const result = this.verifier.verify(
            envelope,
            this.rulesContainerDecoder,
            this.userSignaturesDecoder,
            cached
          );
          items.push(result.verifiedWhitelistedAddress);
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
   * Builds a cache of verified rules containers from the normalized response.
   * When rulesContainerNormalized=true, the API returns deduplicated rules containers.
   * Each is verified once and cached by hash.
   */
  private buildRulesContainerCache(
    response: { rulesContainers?: Array<{ hash?: string; rulesContainer?: string; rulesSignatures?: string }> }
  ): Map<string, DecodedRulesContainer> {
    const cache = new Map<string, DecodedRulesContainer>();
    if (!response.rulesContainers || response.rulesContainers.length === 0) {
      return cache;
    }

    // Deduplicate by base64 container string to avoid re-verifying identical containers
    const verifiedContainers = new Map<string, DecodedRulesContainer>();

    for (const hashContainer of response.rulesContainers) {
      const containerHash = hashContainer.hash;
      const containerBase64 = hashContainer.rulesContainer;
      const signaturesBase64 = hashContainer.rulesSignatures;

      if (!containerHash || !containerBase64) {
        continue;
      }

      // Check if already verified (dedup by content)
      let decoded = verifiedContainers.get(containerBase64);
      if (!decoded) {
        decoded = this.verifier.verifyAndDecodeRulesContainer(
          containerBase64,
          signaturesBase64 ?? "",
          this.rulesContainerDecoder,
          this.userSignaturesDecoder
        );
        verifiedContainers.set(containerBase64, decoded);
      }

      cache.set(containerHash, decoded);
    }

    return cache;
  }

  /**
   * Maps a DTO to a SignedWhitelistedAddressEnvelope.
   */
  private mapDtoToEnvelope(
    dto: {
      id?: string;
      metadata?: { hash?: string; payloadAsString?: string; payload?: object };
      rulesContainer?: string;
      rulesContainerHash?: string;
      rulesSignatures?: string;
      signedAddress?: {
        payload?: string;
        signatures?: Array<{
          signature?: { userId?: string; signature?: string; comment?: string };
          hashes?: string[];
        }>;
      };
      blockchain?: string;
      network?: string;
    },
    addressId: string
  ): SignedWhitelistedAddressEnvelope {
    // SECURITY: Extract linked addresses/wallets from verified payloadAsString
    //
    // The metadata contains two representations:
    //   - payload: Raw object from API response (UNVERIFIED - could be tampered)
    //   - payloadAsString: JSON string that is cryptographically hashed (VERIFIED)
    //
    // The hash verification chain is:
    //   1. Server: metadata.hash = SHA256(payloadAsString)
    //   2. Hash is signed by governance rules (SuperAdmin keys)
    //   3. Client verifies: computed_hash(payloadAsString) === metadata.hash
    //
    // ATTACK VECTOR (if using raw payload):
    // An attacker intercepting the API response could:
    //   1. Modify payload.linkedInternalAddresses to include malicious addresses
    //   2. Leave payloadAsString unchanged (hash still verifies)
    //   3. Client trusts the modified linked addresses → SECURITY BYPASS
    //
    // By parsing from payloadAsString, we ensure the integrity chain is preserved.
    let linkedInternalAddresses: InternalAddress[] = [];
    let linkedWallets: InternalWallet[] = [];

    if (dto.metadata?.payloadAsString) {
      try {
        const verifiedPayload = JSON.parse(dto.metadata.payloadAsString) as Record<string, unknown>;
        if (Array.isArray(verifiedPayload.linkedInternalAddresses)) {
          linkedInternalAddresses = verifiedPayload.linkedInternalAddresses.map(
            (addr: { id?: number; label?: string }) => ({
              id: addr.id ?? 0,
              label: addr.label,
            })
          );
        }
        if (Array.isArray(verifiedPayload.linkedWallets)) {
          linkedWallets = verifiedPayload.linkedWallets.map(
            (wallet: { id?: number; path?: string; name?: string }) => ({
              id: wallet.id ?? 0,
              path: wallet.path,
              label: wallet.name,
            })
          );
        }
      } catch (error: unknown) {
        // JSON parsing failure is expected for malformed payloads — leave arrays empty.
        // The main address data comes from parseWhitelistedAddressFromJson in the verifier.
        if (!(error instanceof SyntaxError)) {
          throw error;
        }
      }
    }

    return {
      id: addressId,
      metadata: {
        hash: dto.metadata?.hash ?? "",
        payloadAsString: dto.metadata?.payloadAsString ?? "",
      },
      rulesContainerBase64: dto.rulesContainer ?? "",
      rulesSignaturesBase64: dto.rulesSignatures ?? "",
      rulesContainerHash: dto.rulesContainerHash,
      signedAddress: {
        payload: dto.signedAddress?.payload,
        signatures: (dto.signedAddress?.signatures ?? []).map((sig) => ({
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
      linkedInternalAddresses,
      linkedWallets,
    };
  }
}
