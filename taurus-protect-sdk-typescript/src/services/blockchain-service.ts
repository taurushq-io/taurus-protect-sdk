/**
 * Blockchain service for Taurus-PROTECT SDK.
 *
 * Provides methods for retrieving blockchain information.
 */

import { NotFoundError, ValidationError } from '../errors';
import type { BlockchainApi } from '../internal/openapi/apis/BlockchainApi';
import { blockchainFromDto, blockchainsFromDto } from '../mappers/blockchain';
import type { Blockchain, ListBlockchainsOptions } from '../models/blockchain';
import { BaseService } from './base';

/**
 * Service for retrieving blockchain information.
 *
 * Provides access to supported blockchain networks and their
 * configuration, including chain IDs, block heights, and other metadata.
 *
 * @example
 * ```typescript
 * // List all supported blockchains
 * const blockchains = await blockchainService.list();
 * for (const bc of blockchains) {
 *   console.log(`${bc.symbol} (${bc.network}): ${bc.name}`);
 * }
 *
 * // Filter by blockchain and network
 * const ethMainnet = await blockchainService.list({
 *   blockchain: 'ETH',
 *   network: 'mainnet',
 * });
 *
 * // Get blockchains with block height info
 * const withHeight = await blockchainService.list({
 *   blockchain: 'ETH',
 *   network: 'mainnet',
 *   includeBlockHeight: true,
 * });
 * console.log(`Current block height: ${withHeight[0]?.blockHeight}`);
 *
 * // Get a specific blockchain by symbol and network
 * const btc = await blockchainService.get('BTC', 'mainnet');
 * console.log(`${btc.name}: ${btc.confirmations} confirmations required`);
 * ```
 */
export class BlockchainService extends BaseService {
  private readonly blockchainApi: BlockchainApi;

  /**
   * Creates a new BlockchainService instance.
   *
   * @param blockchainApi - The BlockchainApi instance from the OpenAPI client
   */
  constructor(blockchainApi: BlockchainApi) {
    super();
    this.blockchainApi = blockchainApi;
  }

  /**
   * Lists all supported blockchains.
   *
   * @param options - Optional filtering options
   * @returns Array of blockchains
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // List all blockchains
   * const blockchains = await blockchainService.list();
   *
   * // Filter by blockchain symbol
   * const ethBlockchains = await blockchainService.list({ blockchain: 'ETH' });
   *
   * // Filter by network
   * const testnets = await blockchainService.list({ network: 'testnet' });
   *
   * // Get with block height (requires blockchain/network)
   * const btcMainnet = await blockchainService.list({
   *   blockchain: 'BTC',
   *   network: 'mainnet',
   *   includeBlockHeight: true,
   * });
   * ```
   */
  async list(options?: ListBlockchainsOptions): Promise<Blockchain[]> {
    return this.execute(async () => {
      const response = await this.blockchainApi.blockchainServiceGetBlockchains({
        blockchain: options?.blockchain,
        network: options?.network,
        includeBlockHeight: options?.includeBlockHeight,
      });

      const result =
        (response as Record<string, unknown>).blockchains ??
        (response as Record<string, unknown>).result;
      return blockchainsFromDto(result as unknown[]);
    });
  }

  /**
   * Gets a blockchain by symbol and network.
   *
   * @param blockchain - The blockchain symbol (e.g., "BTC", "ETH")
   * @param network - The network type (e.g., "mainnet", "testnet")
   * @param includeBlockHeight - Whether to include current block height
   * @returns The blockchain
   * @throws {@link ValidationError} If blockchain or network is empty
   * @throws {@link NotFoundError} If blockchain is not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const btc = await blockchainService.get('BTC', 'mainnet');
   * console.log(`${btc.name}: ${btc.confirmations} confirmations`);
   *
   * // With block height
   * const eth = await blockchainService.get('ETH', 'mainnet', true);
   * console.log(`Current block: ${eth.blockHeight}`);
   * ```
   */
  async get(
    blockchain: string,
    network: string,
    includeBlockHeight?: boolean
  ): Promise<Blockchain> {
    if (!blockchain || blockchain.trim() === '') {
      throw new ValidationError('blockchain is required');
    }
    if (!network || network.trim() === '') {
      throw new ValidationError('network is required');
    }

    return this.execute(async () => {
      const response = await this.blockchainApi.blockchainServiceGetBlockchains({
        blockchain,
        network,
        includeBlockHeight,
      });

      const result =
        (response as Record<string, unknown>).blockchains ??
        (response as Record<string, unknown>).result;
      const blockchains = blockchainsFromDto(result as unknown[]);

      if (blockchains.length === 0) {
        throw new NotFoundError(
          `Blockchain '${blockchain}' with network '${network}' not found`
        );
      }

      // Return the first match
      const bc = blockchains[0];
      if (!bc) {
        throw new NotFoundError(
          `Blockchain '${blockchain}' with network '${network}' not found`
        );
      }

      return bc;
    });
  }
}
