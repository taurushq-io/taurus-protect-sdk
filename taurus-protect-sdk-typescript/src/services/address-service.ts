/**
 * Address service for Taurus-PROTECT SDK.
 *
 * Provides operations for creating, retrieving, and managing blockchain addresses
 * within wallets. All addresses retrieved through this service are automatically
 * verified for cryptographic integrity using the rules container public keys.
 * A RulesContainerCache is required at construction time.
 */

import type {
  AddressesApi,
  TgvalidatordCreateAddressRequest,
  TgvalidatordGetAddressProofOfReserveReply,
  WalletServiceCreateAddressAttributesBody,
} from "../internal/openapi";
import { ConfigurationError, IntegrityError, NotFoundError, ServerError, ValidationError } from "../errors";
import { verifyAddressSignature } from "../helpers";
import type { RulesContainerCache } from "../cache";
import type { DecodedRulesContainer } from "../models/governance-rules";
import type { Address, CreateAddressRequest, ListAddressesOptions } from "../models/address";
import type { Pagination } from "../models/pagination";
import { addressFromDto, addressesFromDto } from "../mappers/address";
import { BaseService } from "./base";

/**
 * Service for managing blockchain addresses.
 *
 * Provides operations for creating, retrieving, and managing addresses
 * within wallets. All addresses retrieved through this service are
 * automatically verified for cryptographic integrity using the rules
 * container public keys. A RulesContainerCache must be provided at
 * construction time — verification is mandatory.
 *
 * @example
 * ```typescript
 * // Create a new address
 * const address = await addressService.create({
 *   walletId: "123",
 *   label: "Customer Deposit",
 *   comment: "Primary deposit address",
 * });
 *
 * // Get an address
 * const address = await addressService.get(456);
 * console.log(`Address: ${address.address}`);
 *
 * // List addresses for a wallet
 * const { items, pagination } = await addressService.list(123);
 * ```
 */
export class AddressService extends BaseService {
  private readonly rulesCache: RulesContainerCache;

  /**
   * Creates a new AddressService.
   *
   * Address signature verification is MANDATORY — a RulesContainerCache must
   * always be provided. This mirrors Java/Go/Python where the cache is required
   * at construction time.
   *
   * @param api - The AddressesApi instance from OpenAPI client
   * @param rulesCache - Rules container cache for signature verification (required)
   * @throws ConfigurationError if rulesCache is not provided
   */
  constructor(
    private readonly api: AddressesApi,
    rulesCache: RulesContainerCache
  ) {
    super();
    if (!rulesCache) {
      throw new ConfigurationError(
        "RulesContainerCache is required for AddressService — address signature verification is mandatory"
      );
    }
    this.rulesCache = rulesCache;
  }

  /**
   * Gets an address by ID with mandatory signature verification.
   *
   * @param addressId - The address ID to retrieve
   * @returns The verified address
   * @throws ValidationError if addressId is invalid
   * @throws NotFoundError if address not found
   * @throws IntegrityError if signature verification fails
   * @throws APIError if API request fails
   */
  async get(addressId: number): Promise<Address> {
    if (addressId <= 0) {
      throw new ValidationError("addressId must be positive");
    }

    return this.execute(async () => {
      const response = await this.api.walletServiceGetAddress({
        id: String(addressId),
      });

      const result = response.result;
      if (result == null) {
        throw new NotFoundError(`Address ${addressId} not found`);
      }

      const address = addressFromDto(result);
      if (address == null) {
        throw new NotFoundError(`Address ${addressId} not found`);
      }

      // CRITICAL: Verify signature using HSM public key from rules container
      await this.verifyAddressSignature(address);

      return address;
    });
  }

  /**
   * Lists addresses for a wallet with mandatory signature verification.
   *
   * @param walletId - The wallet ID to list addresses for
   * @param options - Optional pagination and filtering options
   * @returns Object with addresses array and pagination info
   * @throws ValidationError if walletId is invalid or limit/offset are invalid
   * @throws IntegrityError if signature verification fails for any address
   * @throws APIError if API request fails
   */
  async list(
    walletId: number,
    options?: Omit<ListAddressesOptions, "walletId">
  ): Promise<{ items: Address[]; pagination: Pagination | undefined }> {
    if (walletId <= 0) {
      throw new ValidationError("walletId must be positive");
    }

    const limit = options?.limit ?? 50;
    const offset = options?.offset ?? 0;

    if (limit <= 0) {
      throw new ValidationError("limit must be positive");
    }
    if (offset < 0) {
      throw new ValidationError("offset cannot be negative");
    }

    return this.execute(async () => {
      const response = await this.api.walletServiceGetAddresses({
        walletId: String(walletId),
        limit: String(limit),
        offset: String(offset),
        query: options?.query,
        blockchain: options?.blockchain,
        network: options?.network,
      });

      const addresses = addressesFromDto(response.result);

      // CRITICAL: Verify signatures for all addresses
      // Pre-fetch rules container once to avoid N+1 cache lookups
      let rulesContainer: DecodedRulesContainer | undefined;
      if (this.rulesCache && addresses.length > 0) {
        rulesContainer = await this.rulesCache.get();
      }

      for (const address of addresses) {
        await this.verifyAddressSignature(address, rulesContainer);
      }

      // Extract pagination
      const totalItems = response.totalItems
        ? parseInt(response.totalItems, 10)
        : undefined;
      const pagination: Pagination | undefined = totalItems !== undefined
        ? { totalItems, offset, limit }
        : undefined;

      return { items: addresses, pagination };
    });
  }

  /**
   * Lists addresses with full filtering options.
   *
   * @param options - Optional filtering and pagination options
   * @returns Object with addresses array and pagination info
   * @throws IntegrityError if signature verification fails for any address
   * @throws APIError if API request fails
   */
  async listWithOptions(
    options?: ListAddressesOptions
  ): Promise<{ items: Address[]; pagination: Pagination | undefined }> {
    const limit = options?.limit ?? 50;
    const offset = options?.offset ?? 0;

    return this.execute(async () => {
      const response = await this.api.walletServiceGetAddresses({
        walletId: options?.walletId,
        limit: limit > 0 ? String(limit) : undefined,
        offset: offset > 0 ? String(offset) : undefined,
        query: options?.query,
        blockchain: options?.blockchain,
        network: options?.network,
      });

      const addresses = addressesFromDto(response.result);

      // CRITICAL: Verify signatures for all addresses
      // Pre-fetch rules container once to avoid N+1 cache lookups
      let rulesContainer: DecodedRulesContainer | undefined;
      if (this.rulesCache && addresses.length > 0) {
        rulesContainer = await this.rulesCache.get();
      }

      for (const address of addresses) {
        await this.verifyAddressSignature(address, rulesContainer);
      }

      // Extract pagination
      const totalItems = response.totalItems
        ? parseInt(response.totalItems, 10)
        : undefined;
      const pagination: Pagination | undefined = totalItems !== undefined
        ? { totalItems, offset, limit }
        : undefined;

      return { items: addresses, pagination };
    });
  }

  /**
   * Creates a new address.
   *
   * @param request - Address creation parameters
   * @returns The created address
   * @throws ValidationError if required fields are missing
   * @throws APIError if API request fails
   */
  async create(request: CreateAddressRequest): Promise<Address> {
    if (!request) {
      throw new ValidationError("request cannot be null or undefined");
    }
    if (!request.walletId) {
      throw new ValidationError("walletId is required");
    }
    if (!request.label) {
      throw new ValidationError("label is required");
    }

    return this.execute(async () => {
      const body: TgvalidatordCreateAddressRequest = {
        walletId: request.walletId,
        label: request.label,
        comment: request.comment ?? "",
        customerId: request.customerId,
        externalAddressId: request.externalAddressId,
        type: request.addressType,
        nonHardenedDerivation: request.nonHardenedDerivation,
      };

      const response = await this.api.walletServiceCreateAddress({ body });

      const result = response.result;
      if (result == null) {
        throw new ServerError("Failed to create address: no result returned");
      }

      const address = addressFromDto(result);
      if (address == null) {
        throw new ServerError("Failed to create address: invalid response");
      }

      return address;
    });
  }

  /**
   * Creates an address with explicit parameters.
   *
   * @param walletId - The wallet ID to create the address in
   * @param label - Human-readable label for the address
   * @param comment - Optional description
   * @param customerId - Optional customer identifier
   * @returns The created address
   * @throws ValidationError if required fields are missing or invalid
   * @throws APIError if API request fails
   */
  async createAddress(
    walletId: number,
    label: string,
    comment: string = "",
    customerId?: string
  ): Promise<Address> {
    if (walletId <= 0) {
      throw new ValidationError("walletId must be positive");
    }
    if (!label) {
      throw new ValidationError("label is required");
    }

    return this.create({
      walletId: String(walletId),
      label,
      comment,
      customerId,
    });
  }

  /**
   * Creates an attribute for an address.
   *
   * @param addressId - The address ID
   * @param key - The attribute key
   * @param value - The attribute value
   * @throws ValidationError if any argument is invalid
   * @throws APIError if API request fails
   */
  async createAttribute(
    addressId: number,
    key: string,
    value: string
  ): Promise<void> {
    if (addressId <= 0) {
      throw new ValidationError("addressId must be positive");
    }
    if (!key) {
      throw new ValidationError("key is required");
    }
    if (!value) {
      throw new ValidationError("value is required");
    }

    return this.execute(async () => {
      const body: WalletServiceCreateAddressAttributesBody = {
        attributes: [{ key, value }],
      };

      await this.api.walletServiceCreateAddressAttributes({
        addressId: String(addressId),
        body,
      });
    });
  }

  /**
   * Deletes an attribute from an address.
   *
   * @param addressId - The address ID
   * @param attributeId - The attribute ID to delete
   * @throws ValidationError if any argument is invalid
   * @throws APIError if API request fails
   */
  async deleteAttribute(addressId: number, attributeId: number): Promise<void> {
    if (addressId <= 0) {
      throw new ValidationError("addressId must be positive");
    }
    if (attributeId <= 0) {
      throw new ValidationError("attributeId must be positive");
    }

    return this.execute(async () => {
      await this.api.walletServiceDeleteAddressAttribute({
        addressId: String(addressId),
        id: String(attributeId),
      });
    });
  }

  /**
   * Gets the proof of reserve for an address.
   *
   * @param addressId - The address ID
   * @param challenge - Optional challenge string
   * @returns The proof of reserve response
   * @throws ValidationError if addressId is invalid
   * @throws APIError if API request fails
   */
  async getProofOfReserve(
    addressId: number,
    challenge?: string
  ): Promise<TgvalidatordGetAddressProofOfReserveReply["result"]> {
    if (addressId <= 0) {
      throw new ValidationError("addressId must be positive");
    }

    return this.execute(async () => {
      const response = await this.api.walletServiceGetAddressProofOfReserve({
        id: String(addressId),
        challenge,
      });
      return response.result;
    });
  }

  /**
   * Verifies the signature of an address using the rules container.
   *
   * @param address - The address to verify
   * @param rulesContainer - Optional pre-fetched rules container. If undefined,
   *   will be fetched from cache. Pass this when verifying multiple addresses
   *   to avoid N+1 cache lookups.
   * @throws IntegrityError if signature verification fails
   */
  private async verifyAddressSignature(
    address: Address,
    rulesContainer?: DecodedRulesContainer
  ): Promise<void> {
    // Get rules container if not provided
    const rules = rulesContainer ?? (await this.rulesCache.get());

    // Address must have a signature
    if (!address.signature) {
      throw new IntegrityError(
        `Address ${address.id} is missing signature - cannot verify integrity`
      );
    }

    // Verify the signature
    const isValid = verifyAddressSignature(
      address.address,
      address.signature,
      rules
    );

    if (!isValid) {
      throw new IntegrityError(
        `Invalid signature for address ${address.id} - data may have been tampered with`
      );
    }
  }
}
