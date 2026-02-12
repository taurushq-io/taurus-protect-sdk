/**
 * Token metadata service for Taurus-PROTECT SDK.
 *
 * Provides methods for retrieving token metadata for various token standards
 * including ERC-20, ERC-721, ERC-1155 (Ethereum), FA tokens (Tezos), and CryptoPunks.
 */

import { ValidationError } from '../errors';
import type { TokenMetadataApi } from '../internal/openapi/apis/TokenMetadataApi';
import type {
  TokenMetadata,
  CryptoPunkMetadata,
  GetERCTokenMetadataOptions,
  GetEVMERCTokenMetadataOptions,
  GetFATokenMetadataOptions,
  GetCryptoPunkMetadataOptions,
} from '../models/token-metadata';
import { BaseService } from './base';

// Re-export types for backwards compatibility
export type {
  TokenMetadata,
  CryptoPunkMetadata,
  GetERCTokenMetadataOptions,
  GetEVMERCTokenMetadataOptions,
  GetFATokenMetadataOptions,
  GetCryptoPunkMetadataOptions,
} from '../models/token-metadata';

/**
 * Service for retrieving token metadata.
 *
 * Provides access to token metadata for various token standards
 * including ERC-20, ERC-721, ERC-1155 (Ethereum), and FA tokens (Tezos).
 *
 * @example
 * ```typescript
 * // Get ERC-20 token metadata
 * const usdc = await tokenMetadataService.getERCTokenMetadata({
 *   network: 'mainnet',
 *   contract: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48',
 *   blockchain: 'ETH',
 * });
 * console.log(`Token: ${usdc.name}, Decimals: ${usdc.decimals}`);
 *
 * // Get ERC-721 NFT metadata with image data
 * const nft = await tokenMetadataService.getEVMERCTokenMetadata({
 *   network: 'mainnet',
 *   contract: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D',
 *   tokenId: '1234',
 *   withData: true,
 *   blockchain: 'ETH',
 * });
 * console.log(`NFT: ${nft.name}, URI: ${nft.uri}`);
 *
 * // Get FA token metadata (Tezos)
 * const tzBtc = await tokenMetadataService.getFATokenMetadata({
 *   network: 'mainnet',
 *   contract: 'KT1PWx2mnDueood7fEmfbBDKx1D9BAnnXitn',
 *   tokenId: '0',
 * });
 * console.log(`Token: ${tzBtc.name}, Symbol: ${tzBtc.symbol}`);
 *
 * // Get CryptoPunk metadata
 * const punk = await tokenMetadataService.getCryptoPunkMetadata({
 *   network: 'mainnet',
 *   contract: '0xb47e3cd837dDF8e4c57F05d70Ab865de6e193BBB',
 *   punkId: '7804',
 * });
 * console.log(`Punk #${punk.punkId}: ${punk.punkAttributes}`);
 * ```
 */
export class TokenMetadataService extends BaseService {
  private readonly tokenMetadataApi: TokenMetadataApi;

  /**
   * Creates a new TokenMetadataService instance.
   *
   * @param tokenMetadataApi - The TokenMetadataApi instance from the OpenAPI client
   */
  constructor(tokenMetadataApi: TokenMetadataApi) {
    super();
    this.tokenMetadataApi = tokenMetadataApi;
  }

  /**
   * Retrieves ERC token metadata (ERC-20, ERC-721, ERC-1155).
   *
   * @deprecated Use getEVMERCTokenMetadata instead
   * @param options - The options for retrieving token metadata
   * @returns The token metadata
   * @throws {@link ValidationError} If required arguments are missing
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const metadata = await tokenMetadataService.getERCTokenMetadata({
   *   network: 'mainnet',
   *   contract: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48',
   *   blockchain: 'ETH',
   * });
   * console.log(`Token: ${metadata.name}, Decimals: ${metadata.decimals}`);
   * ```
   */
  async getERCTokenMetadata(options: GetERCTokenMetadataOptions): Promise<TokenMetadata> {
    if (!options.network || options.network.trim() === '') {
      throw new ValidationError('network is required');
    }
    if (!options.contract || options.contract.trim() === '') {
      throw new ValidationError('contract is required');
    }

    return this.execute(async () => {
      const response = await this.tokenMetadataApi.tokenMetadataServiceGetERCTokenMetadata({
        network: options.network,
        contract: options.contract,
        token: options.tokenId ?? '',
        withData: options.withData,
        blockchain: options.blockchain,
      });

      return this.mapERCTokenMetadata(response.result);
    });
  }

  /**
   * Retrieves ERC token metadata for EVM-compatible chains.
   *
   * This is the preferred method for getting ERC token metadata
   * on any EVM-compatible blockchain.
   *
   * @param options - The options for retrieving token metadata
   * @returns The token metadata
   * @throws {@link ValidationError} If required arguments are missing
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // Get ERC-20 token metadata on Ethereum
   * const eth = await tokenMetadataService.getEVMERCTokenMetadata({
   *   network: 'mainnet',
   *   contract: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48',
   *   blockchain: 'ETH',
   * });
   *
   * // Get ERC-721 NFT metadata on Polygon
   * const nft = await tokenMetadataService.getEVMERCTokenMetadata({
   *   network: 'mainnet',
   *   contract: '0x...',
   *   tokenId: '1234',
   *   withData: true,
   *   blockchain: 'MATIC',
   * });
   * ```
   */
  async getEVMERCTokenMetadata(options: GetEVMERCTokenMetadataOptions): Promise<TokenMetadata> {
    if (!options.network || options.network.trim() === '') {
      throw new ValidationError('network is required');
    }
    if (!options.contract || options.contract.trim() === '') {
      throw new ValidationError('contract is required');
    }
    if (!options.blockchain || options.blockchain.trim() === '') {
      throw new ValidationError('blockchain is required');
    }

    return this.execute(async () => {
      const response = await this.tokenMetadataApi.tokenMetadataServiceGetEVMERCTokenMetadata({
        network: options.network,
        contract: options.contract,
        token: options.tokenId ?? '',
        withData: options.withData,
        blockchain: options.blockchain,
      });

      return this.mapERCTokenMetadata(response.result);
    });
  }

  /**
   * Retrieves FA token metadata (Tezos FA1.2/FA2 standards).
   *
   * @param options - The options for retrieving token metadata
   * @returns The token metadata
   * @throws {@link ValidationError} If required arguments are missing
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const metadata = await tokenMetadataService.getFATokenMetadata({
   *   network: 'mainnet',
   *   contract: 'KT1PWx2mnDueood7fEmfbBDKx1D9BAnnXitn',
   *   tokenId: '0',
   * });
   * console.log(`Token: ${metadata.name}, Symbol: ${metadata.symbol}`);
   * ```
   */
  async getFATokenMetadata(options: GetFATokenMetadataOptions): Promise<TokenMetadata> {
    if (!options.network || options.network.trim() === '') {
      throw new ValidationError('network is required');
    }
    if (!options.contract || options.contract.trim() === '') {
      throw new ValidationError('contract is required');
    }

    return this.execute(async () => {
      const response = await this.tokenMetadataApi.tokenMetadataServiceGetFATokenMetadata({
        network: options.network,
        contract: options.contract,
        token: options.tokenId ?? '0',
        withData: options.withData,
      });

      return this.mapFATokenMetadata(response.result);
    });
  }

  /**
   * Retrieves CryptoPunk metadata.
   *
   * Returns the metadata for a CryptoPunk NFT including its
   * attributes and image data.
   *
   * @param options - The options for retrieving punk metadata
   * @returns The CryptoPunk metadata
   * @throws {@link ValidationError} If required arguments are missing
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const punk = await tokenMetadataService.getCryptoPunkMetadata({
   *   network: 'mainnet',
   *   contract: '0xb47e3cd837dDF8e4c57F05d70Ab865de6e193BBB',
   *   punkId: '7804',
   * });
   * console.log(`Punk #${punk.punkId}: ${punk.punkAttributes}`);
   * ```
   */
  async getCryptoPunkMetadata(options: GetCryptoPunkMetadataOptions): Promise<CryptoPunkMetadata> {
    if (!options.network || options.network.trim() === '') {
      throw new ValidationError('network is required');
    }
    if (!options.contract || options.contract.trim() === '') {
      throw new ValidationError('contract is required');
    }
    if (!options.punkId || options.punkId.trim() === '') {
      throw new ValidationError('punkId is required');
    }

    return this.execute(async () => {
      const response =
        await this.tokenMetadataApi.tokenMetadataServiceGetCryptoPunksTokenMetadata({
          network: options.network,
          contract: options.contract,
          token: options.punkId,
          blockchain: options.blockchain,
        });

      return this.mapCryptoPunkMetadata(response.result);
    });
  }

  /**
   * Maps ERC token metadata from the API response to the domain model.
   */
  private mapERCTokenMetadata(dto: unknown): TokenMetadata {
    if (!dto) {
      return {};
    }
    const d = dto as Record<string, unknown>;
    return {
      name: d.name as string | undefined,
      description: d.description as string | undefined,
      decimals: d.decimals as string | undefined,
      dataType: d.dataType as string | undefined,
      base64Data: d.base64Data as string | undefined,
      uri: d.uri as string | undefined,
    };
  }

  /**
   * Maps FA token metadata from the API response to the domain model.
   */
  private mapFATokenMetadata(dto: unknown): TokenMetadata {
    if (!dto) {
      return {};
    }
    const d = dto as Record<string, unknown>;
    return {
      name: d.name as string | undefined,
      symbol: d.symbol as string | undefined,
      decimals: d.decimals as string | undefined,
      dataType: d.dataType as string | undefined,
      base64Data: d.base64Data as string | undefined,
      uri: d.uri as string | undefined,
    };
  }

  /**
   * Maps CryptoPunk metadata from the API response to the domain model.
   */
  private mapCryptoPunkMetadata(dto: unknown): CryptoPunkMetadata {
    if (!dto) {
      return {};
    }
    const d = dto as Record<string, unknown>;
    return {
      punkId: d.punkId as string | undefined,
      punkAttributes: d.punkAttributes as string | undefined,
      image: d.image as string | undefined,
    };
  }
}
