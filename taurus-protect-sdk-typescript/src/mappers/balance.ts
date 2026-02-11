/**
 * Balance mapper functions for converting OpenAPI DTOs to domain models.
 */

import type { AssetBalance, NFTCollectionBalance } from '../models/balance';
import { safeInt, safeMap, safeString } from './base';

/**
 * Maps a balance DTO to an AssetBalance domain model.
 */
export function assetBalanceFromDto(dto: unknown): AssetBalance | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    currencyId: safeString(d.currencyId ?? d.currency_id),
    currency: safeString(d.currency ?? d.symbol),
    blockchain: safeString(d.blockchain),
    network: safeString(d.network),
    contractAddress: safeString(d.contractAddress ?? d.contract_address),
    tokenId: safeString(d.tokenId ?? d.token_id),
    balance: safeString(d.balance ?? d.totalConfirmed ?? d.total_confirmed),
    fiatValue: safeString(d.fiatValue ?? d.fiat_value),
    fiatCurrency: safeString(d.fiatCurrency ?? d.fiat_currency),
  };
}

/**
 * Maps an array of balance DTOs to AssetBalance domain models.
 */
export function assetBalancesFromDto(dtos: unknown[] | null | undefined): AssetBalance[] {
  return safeMap(dtos, assetBalanceFromDto);
}

/**
 * Maps an NFT collection balance DTO to an NFTCollectionBalance domain model.
 */
export function nftCollectionBalanceFromDto(dto: unknown): NFTCollectionBalance | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    name: safeString(d.name),
    symbol: safeString(d.symbol),
    blockchain: safeString(d.blockchain),
    network: safeString(d.network),
    contractAddress: safeString(d.contractAddress ?? d.contract_address),
    count: safeInt(d.count ?? d.balance),
    logoUrl: safeString(d.logoUrl ?? d.logo_url ?? d.logo),
  };
}

/**
 * Maps an array of NFT collection balance DTOs to NFTCollectionBalance domain models.
 */
export function nftCollectionBalancesFromDto(
  dtos: unknown[] | null | undefined
): NFTCollectionBalance[] {
  return safeMap(dtos, nftCollectionBalanceFromDto);
}
