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
  TgvalidatordAddress,
  TgvalidatordCreateAddressRequest,
  TgvalidatordGetAddressProofOfReserveReply,
  WalletServiceCreateAddressAttributesBody,
} from "../internal/openapi";
import {
  ConfigurationError,
  IntegrityError,
  NotFoundError,
  ServerError,
  ValidationError,
} from "../errors";
import { verifyAddressSignature } from "../helpers";
import type { RulesContainerCache } from "../cache";
import type { DecodedRulesContainer } from "../models/governance-rules";
import type { Address, CreateAddressRequest, ListAddressesOptions } from "../models/address";
import type { Pagination } from "../models/pagination";
import { addressFromDto } from "../mappers/address";
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

      // Through the ONE verification seam. See verifiedAddress.
      const address = await this.verifiedAddress(result);
      if (address == null) {
        throw new NotFoundError(`Address ${addressId} not found`);
      }

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

      // CRITICAL: every row goes through the ONE verification seam.
      const addresses = await this.verifiedAddresses(response.result);

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

      // CRITICAL: every row goes through the ONE verification seam.
      const addresses = await this.verifiedAddresses(response.result);

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
   * Creates a new address, with mandatory signature verification of the reply.
   *
   * The created address is returned only once its HSM signature has been verified.
   * Creation can be asynchronous, so a reply may legitimately carry no address string
   * yet — that comes back with its `status` and nothing to verify. What never happens
   * is a non-empty `address` being returned unverified. See `verifiedAddress`.
   *
   * @param request - Address creation parameters
   * @returns The created address
   * @throws ValidationError if required fields are missing
   * @throws IntegrityError if the reply carries an address string that cannot be
   *   verified — either it has no HSM signature, or the signature does not verify
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

      // Through the same seam as every read. The create reply carries the same
      // TgvalidatordAddress the read paths return, signature included, so there was
      // never a reason for this path to be the unverified one — and callers of create
      // are precisely the ones about to publish or fund a fresh deposit address.
      const address = await this.verifiedAddress(result);
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
   * @throws IntegrityError if the reply carries an unverifiable address string
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
   * The ONE construction seam for an {@link Address}.
   *
   * Every path that returns an Address goes through here — `get`, `list`,
   * `listWithOptions` and `create` — because "remember to verify" was a rule rather
   * than the only available construction path, and `create` is what that cost:
   * `AssetService.getAssetAddresses` had already been fixed for exactly this and the
   * create path was missed anyway. Same reasoning as the request service's
   * verification seam.
   *
   * The rule it enforces: **never return a non-empty `Address.address` string that
   * has not been verified.** It branches on the address STRING, not on `status`,
   * because `status` is server-controlled and so cannot be the thing that decides
   * whether a check runs:
   *
   * - address empty → return as-is. Asynchronous creation is real (`status` is one of
   *   created/creating/signed/observed/confirmed), and an address that has not been
   *   generated yet carries no destination: nothing to verify, nothing to misuse.
   * - address present, signature absent → refuse. Returning it would hand the caller
   *   an attacker-controllable destination in the same type as a verified one, which
   *   is the whole defect. The caller re-reads through this same seam once the status
   *   advances.
   * - address present, signature present → verify against the HSMSLOT key.
   *
   * The refusal is raised BEFORE the rules container is fetched: there is nothing the
   * container could say that would make an unsigned address string usable.
   *
   * @param dto - The address DTO straight off the wire
   * @param rulesContainer - Optional pre-fetched container. Pass it when verifying
   *   several addresses to avoid an N+1 of cache lookups.
   * @returns The verified address, or undefined when the DTO maps to nothing (each
   *   caller raises its own not-found/invalid-response error, which differ)
   * @throws IntegrityError if the address string cannot be shown to be genuine
   */
  private async verifiedAddress(
    dto: TgvalidatordAddress,
    rulesContainer?: DecodedRulesContainer
  ): Promise<Address | undefined> {
    const address = addressFromDto(dto);
    if (address == null) {
      return undefined;
    }

    if (!address.address) {
      return address;
    }

    if (!address.signature) {
      throw new IntegrityError(
        `address ${address.id} (status ${JSON.stringify(address.status ?? "")}) carries an ` +
          `address string but no HSM signature; refusing to return an unverified ` +
          `destination. Re-read once the status reaches "signed"`
      );
    }

    const rules = rulesContainer ?? (await this.rulesCache.get());

    // The helper throws on every failure — absent HSM key, or a signature that does
    // not verify — so those checks live in one place shared with the peer services.
    verifyAddressSignature(address.address, address.signature, rules, address.id);

    return address;
  }

  /**
   * Maps and verifies a page of address DTOs through {@link verifiedAddress}.
   *
   * The container is fetched once for the whole page rather than per row.
   *
   * @param dtos - The address DTOs from a list reply
   * @returns The verified addresses
   * @throws IntegrityError if any row cannot be shown to be genuine
   */
  private async verifiedAddresses(
    dtos: TgvalidatordAddress[] | null | undefined
  ): Promise<Address[]> {
    if (dtos == null || dtos.length === 0) {
      return [];
    }

    // Pre-fetch the rules container once to avoid N+1 cache lookups — but only when a
    // row actually carries a signature to check. Fetching unconditionally would make a
    // page of unsigned rows fail with a governance-fetch error instead of the integrity
    // refusal that names the offending row, which is the wrong diagnosis and breaks the
    // ordering `verifiedAddress` documents (refuse before fetching).
    const anySignatureToCheck = dtos.some(
      (dto) => Boolean(dto.address) && Boolean(dto.signature)
    );
    const rulesContainer = anySignatureToCheck
      ? await this.rulesCache.get()
      : undefined;

    const addresses: Address[] = [];
    for (const dto of dtos) {
      const address = await this.verifiedAddress(dto, rulesContainer);
      if (address !== undefined) {
        addresses.push(address);
      }
    }
    return addresses;
  }
}
