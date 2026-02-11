/**
 * Contract whitelisting service for Taurus-PROTECT SDK.
 *
 * Provides operations for managing whitelisted contract addresses such as
 * ERC20 tokens, NFT collections (ERC721/ERC1155), FA2 tokens on Tezos,
 * and other smart contract-based assets.
 */

import type { ContractWhitelistingApi } from "../internal/openapi/apis/ContractWhitelistingApi";
import { NotFoundError, ValidationError } from "../errors";
import type { Pagination, PaginatedResult } from "../models/pagination";
import type {
  WhitelistedContract,
  WhitelistedContractAttribute,
  ListWhitelistedContractsOptions,
  ListForApprovalOptions,
  CreateWhitelistedContractRequest,
  UpdateWhitelistedContractRequest,
} from "../models/contract-whitelist";
import {
  whitelistedContractFromDto,
  whitelistedContractsFromDto,
  whitelistedContractAttributeFromDto,
} from "../mappers/contract-whitelist";
import { BaseService } from "./base";

// Re-export types for convenience
export type {
  WhitelistedContract,
  WhitelistedContractAttribute,
  ListWhitelistedContractsOptions,
  ListForApprovalOptions,
  CreateWhitelistedContractRequest,
  UpdateWhitelistedContractRequest,
} from "../models/contract-whitelist";

/**
 * Service for managing whitelisted contract addresses.
 *
 * Provides operations for creating, approving, updating, and deleting
 * whitelisted contract addresses such as ERC20 tokens, NFT collections
 * (ERC721/ERC1155), FA2 tokens on Tezos, and other smart contract-based assets.
 *
 * @example
 * ```typescript
 * // List whitelisted contracts
 * const result = await contractWhitelistingService.list({
 *   blockchain: 'ETH',
 *   network: 'mainnet',
 *   limit: 50,
 * });
 *
 * // Get a whitelisted contract by ID
 * const contract = await contractWhitelistingService.get('123');
 *
 * // Create a new whitelisted contract
 * const id = await contractWhitelistingService.create({
 *   blockchain: 'ETH',
 *   network: 'mainnet',
 *   contractAddress: '0x1234...',
 *   symbol: 'USDC',
 *   name: 'USD Coin',
 *   decimals: 6,
 *   kind: 'erc20',
 * });
 * ```
 */
export class ContractWhitelistingService extends BaseService {
  /**
   * Creates a new ContractWhitelistingService instance.
   *
   * @param api - The ContractWhitelistingApi instance from the OpenAPI client
   */
  constructor(private readonly api: ContractWhitelistingApi) {
    super();
  }

  /**
   * Gets a whitelisted contract by ID.
   *
   * @param id - The contract ID to retrieve
   * @returns The whitelisted contract
   * @throws {@link ValidationError} If id is invalid
   * @throws {@link NotFoundError} If contract not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const contract = await contractWhitelistingService.get('123');
   * console.log(`Contract: ${contract.name} (${contract.symbol})`);
   * ```
   */
  async get(id: string): Promise<WhitelistedContract> {
    if (!id || id.trim() === "") {
      throw new ValidationError("id is required");
    }

    return this.execute(async () => {
      const response = await this.api.whitelistServiceGetWhitelistedContract({
        id,
      });

      const envelope = response.result;
      if (envelope == null) {
        throw new NotFoundError(`Whitelisted contract ${id} not found`);
      }

      const contract = whitelistedContractFromDto(envelope);
      if (!contract) {
        throw new NotFoundError(`Whitelisted contract ${id} not found`);
      }

      return contract;
    });
  }

  /**
   * Lists whitelisted contracts with optional filtering and pagination.
   *
   * @param options - Optional filtering and pagination options
   * @returns Paginated result containing whitelisted contracts
   * @throws {@link ValidationError} If options are invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // List all contracts for ETH mainnet
   * const result = await contractWhitelistingService.list({
   *   blockchain: 'ETH',
   *   network: 'mainnet',
   *   limit: 50,
   * });
   *
   * // Filter by kind types
   * const nfts = await contractWhitelistingService.list({
   *   kindTypes: ['nft'],
   * });
   * ```
   */
  async list(
    options?: ListWhitelistedContractsOptions
  ): Promise<PaginatedResult<WhitelistedContract>> {
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
        isNFT: options?.isNFT,
        kindTypes: options?.kindTypes,
        whitelistedContractAddressIds: options?.contractIds,
        includeForApproval: false,
      });

      const items = whitelistedContractsFromDto(response.result);
      const totalItems = response.totalItems
        ? parseInt(response.totalItems, 10)
        : 0;

      const pagination: Pagination = {
        totalItems,
        offset,
        limit,
      };

      return { items, pagination };
    });
  }

  /**
   * Lists whitelisted contracts pending approval.
   *
   * @param options - Optional filtering and pagination options
   * @returns Paginated result containing whitelisted contracts for approval
   * @throws {@link ValidationError} If options are invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const pending = await contractWhitelistingService.listForApproval({
   *   limit: 50,
   * });
   * console.log(`${pending.pagination.totalItems} contracts pending approval`);
   * ```
   */
  async listForApproval(
    options?: ListForApprovalOptions
  ): Promise<PaginatedResult<WhitelistedContract>> {
    const limit = options?.limit ?? 50;
    const offset = options?.offset ?? 0;

    if (limit <= 0) {
      throw new ValidationError("limit must be positive");
    }
    if (offset < 0) {
      throw new ValidationError("offset cannot be negative");
    }

    return this.execute(async () => {
      const response =
        await this.api.whitelistServiceGetWhitelistedContractsForApproval({
          limit: String(limit),
          offset: String(offset),
          ids: options?.ids,
        });

      const items = whitelistedContractsFromDto(response.result);
      const totalItems = response.totalItems
        ? parseInt(response.totalItems, 10)
        : 0;

      const pagination: Pagination = {
        totalItems,
        offset,
        limit,
      };

      return { items, pagination };
    });
  }

  /**
   * Creates a new whitelisted contract address.
   *
   * The contract will be created in a pending state and require approval
   * according to the configured governance rules.
   *
   * @param request - The contract creation request
   * @returns The ID of the created whitelist entry
   * @throws {@link ValidationError} If required parameters are missing
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const id = await contractWhitelistingService.create({
   *   blockchain: 'ETH',
   *   network: 'mainnet',
   *   contractAddress: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48',
   *   symbol: 'USDC',
   *   name: 'USD Coin',
   *   decimals: 6,
   *   kind: 'erc20',
   * });
   * ```
   */
  async create(request: CreateWhitelistedContractRequest): Promise<string> {
    if (!request.blockchain || request.blockchain.trim() === "") {
      throw new ValidationError("blockchain is required");
    }
    if (!request.network || request.network.trim() === "") {
      throw new ValidationError("network is required");
    }
    if (!request.symbol || request.symbol.trim() === "") {
      throw new ValidationError("symbol is required");
    }
    if (!request.name || request.name.trim() === "") {
      throw new ValidationError("name is required");
    }
    if (!request.kind || request.kind.trim() === "") {
      throw new ValidationError("kind is required");
    }

    return this.execute(async () => {
      const response = await this.api.whitelistServiceCreateWhitelistedContract(
        {
          body: {
            blockchain: request.blockchain,
            network: request.network,
            contractAddress: request.contractAddress,
            symbol: request.symbol,
            name: request.name,
            decimals: String(request.decimals),
            kind: request.kind,
            tokenId: request.tokenId,
          },
        }
      );

      const id = response.result?.id;
      if (!id) {
        throw new ValidationError(
          "Failed to create whitelisted contract: no ID returned"
        );
      }

      return id;
    });
  }

  /**
   * Approves one or more whitelisted contract addresses.
   *
   * Requires a cryptographic signature computed over the metadata hashes
   * of the contracts being approved.
   *
   * @param ids - The list of contract IDs to approve
   * @param signature - The approval signature (base64-encoded)
   * @param comment - The approval comment
   * @throws {@link ValidationError} If parameters are invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * await contractWhitelistingService.approve(
   *   ['123', '456'],
   *   'base64-signature...',
   *   'Approved for production use'
   * );
   * ```
   */
  async approve(
    ids: string[],
    signature: string,
    comment: string
  ): Promise<void> {
    if (!ids || ids.length === 0) {
      throw new ValidationError("ids cannot be empty");
    }
    if (!signature || signature.trim() === "") {
      throw new ValidationError("signature is required");
    }
    if (!comment || comment.trim() === "") {
      throw new ValidationError("comment is required");
    }

    return this.execute(async () => {
      await this.api.whitelistServiceApproveWhitelistedContract({
        body: {
          ids,
          signature,
          comment,
        },
      });
    });
  }

  /**
   * Rejects a whitelisted contract.
   *
   * @param ids - The list of contract IDs to reject
   * @param comment - The rejection reason
   * @throws {@link ValidationError} If parameters are invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * await contractWhitelistingService.reject(['123'], 'Not needed');
   * ```
   */
  async reject(ids: string[], comment: string): Promise<void> {
    if (!ids || ids.length === 0) {
      throw new ValidationError("ids cannot be empty");
    }
    if (!comment || comment.trim() === "") {
      throw new ValidationError("comment is required");
    }

    return this.execute(async () => {
      await this.api.whitelistServiceRejectWhitelistedContract({
        body: {
          ids,
          comment,
        },
      });
    });
  }

  /**
   * Updates an existing whitelisted contract.
   *
   * Only the symbol, name, and decimals can be updated.
   * Note: ALGO whitelisted contracts are not editable.
   *
   * @param id - The contract ID to update
   * @param request - The update request
   * @throws {@link ValidationError} If parameters are invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * await contractWhitelistingService.update('123', {
   *   symbol: 'USDC',
   *   name: 'USD Coin (Updated)',
   *   decimals: 6,
   * });
   * ```
   */
  async update(
    id: string,
    request: UpdateWhitelistedContractRequest
  ): Promise<void> {
    if (!id || id.trim() === "") {
      throw new ValidationError("id is required");
    }
    if (!request.symbol || request.symbol.trim() === "") {
      throw new ValidationError("symbol is required");
    }
    if (!request.name || request.name.trim() === "") {
      throw new ValidationError("name is required");
    }

    return this.execute(async () => {
      await this.api.whitelistServiceUpdateWhitelistedContract({
        id,
        body: {
          symbol: request.symbol,
          name: request.name,
          decimals: String(request.decimals),
        },
      });
    });
  }

  /**
   * Creates an attribute on a whitelisted contract.
   *
   * @param contractId - The contract ID
   * @param key - The attribute key
   * @param value - The attribute value
   * @param options - Optional attribute properties
   * @throws {@link ValidationError} If required parameters are invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * await contractWhitelistingService.createAttribute(
   *   '123',
   *   'category',
   *   'stablecoin'
   * );
   * ```
   */
  async createAttribute(
    contractId: string,
    key: string,
    value: string,
    options?: {
      contentType?: string;
      type?: string;
      subType?: string;
    }
  ): Promise<void> {
    if (!contractId || contractId.trim() === "") {
      throw new ValidationError("contractId is required");
    }
    if (!key || key.trim() === "") {
      throw new ValidationError("key is required");
    }

    return this.execute(async () => {
      await this.api.whitelistServiceCreateWhitelistedContractAttributes({
        whitelistedContractAddressId: contractId,
        body: {
          attributes: [
            {
              key,
              value,
              contentType: options?.contentType,
              type: options?.type,
              subtype: options?.subType,
            },
          ],
        },
      });
    });
  }

  /**
   * Gets an attribute from a whitelisted contract.
   *
   * @param contractId - The contract ID
   * @param attributeId - The attribute ID
   * @returns The attribute
   * @throws {@link ValidationError} If parameters are invalid
   * @throws {@link NotFoundError} If attribute not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const attr = await contractWhitelistingService.getAttribute('123', '456');
   * console.log(`${attr.key}: ${attr.value}`);
   * ```
   */
  async getAttribute(
    contractId: string,
    attributeId: string
  ): Promise<WhitelistedContractAttribute> {
    if (!contractId || contractId.trim() === "") {
      throw new ValidationError("contractId is required");
    }
    if (!attributeId || attributeId.trim() === "") {
      throw new ValidationError("attributeId is required");
    }

    return this.execute(async () => {
      const response =
        await this.api.whitelistServiceGetWhitelistedContractAttribute({
          whitelistedContractAddressId: contractId,
          id: attributeId,
        });

      const result = response.result;
      if (result == null) {
        throw new NotFoundError(
          `Attribute ${attributeId} not found on contract ${contractId}`
        );
      }

      const attribute = whitelistedContractAttributeFromDto(result);
      if (!attribute) {
        throw new NotFoundError(
          `Attribute ${attributeId} not found on contract ${contractId}`
        );
      }

      return attribute;
    });
  }

  /**
   * Deletes an attribute from a whitelisted contract.
   *
   * @param contractId - The contract ID
   * @param attributeId - The attribute ID to delete
   * @throws {@link ValidationError} If parameters are invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * await contractWhitelistingService.deleteAttribute('123', '456');
   * ```
   */
  async deleteAttribute(
    contractId: string,
    attributeId: string
  ): Promise<void> {
    if (!contractId || contractId.trim() === "") {
      throw new ValidationError("contractId is required");
    }
    if (!attributeId || attributeId.trim() === "") {
      throw new ValidationError("attributeId is required");
    }

    return this.execute(async () => {
      await this.api.whitelistServiceDeleteWhitelistedContractAttribute({
        whitelistedContractAddressId: contractId,
        id: attributeId,
      });
    });
  }
}
