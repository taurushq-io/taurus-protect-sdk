/**
 * Balance mapper functions for converting OpenAPI DTOs to domain models.
 */

import type { TgvalidatordAssetBalance } from '../internal/openapi/models/TgvalidatordAssetBalance';
import type { TgvalidatordNFTCollectionBalance } from '../internal/openapi/models/TgvalidatordNFTCollectionBalance';
import type { AssetBalance, NFTCollectionBalance } from '../models/balance';
import { safeInt, safeMap } from './base';

/**
 * Maps a balance row to an AssetBalance domain model.
 *
 * The row is `{asset, balance}`: the currency lives in `asset.currencyInfo` and the amounts
 * in the nested `balance` object. The flat `currency` / `balance` keys this mapper used to
 * read do not exist on the wire, so every row came back without a currency and with a
 * `"[object Object]"` balance.
 */
export function assetBalanceFromDto(
  dto: TgvalidatordAssetBalance | null | undefined
): AssetBalance | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const info = dto.asset?.currencyInfo;
  return {
    currencyId: info?.id,
    currency: info?.symbol ?? dto.asset?.currency,
    blockchain: info?.blockchain,
    network: info?.network,
    contractAddress: info?.contractAddress,
    tokenId: dto.asset?.nft?.tokenid ?? info?.tokenID,
    balance: dto.balance?.totalConfirmed,
    fiatValue: undefined,
    fiatCurrency: undefined,
  };
}

/**
 * Maps an array of balance rows to AssetBalance domain models.
 */
export function assetBalancesFromDto(
  dtos: TgvalidatordAssetBalance[] | null | undefined
): AssetBalance[] {
  return safeMap(dtos, assetBalanceFromDto);
}

/**
 * Maps an NFT collection balance row to an NFTCollectionBalance domain model.
 *
 * The row is `{currencyInfo, balance}`: the collection is described by its currency and
 * the count is the confirmed balance.
 */
export function nftCollectionBalanceFromDto(
  dto: TgvalidatordNFTCollectionBalance | null | undefined
): NFTCollectionBalance | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const info = dto.currencyInfo;
  return {
    name: info?.name,
    symbol: info?.symbol,
    blockchain: info?.blockchain,
    network: info?.network,
    contractAddress: info?.contractAddress,
    count: safeInt(dto.balance?.totalConfirmed),
    logoUrl: info?.logo,
  };
}

/**
 * Maps an array of NFT collection balance rows to NFTCollectionBalance domain models.
 */
export function nftCollectionBalancesFromDto(
  dtos: TgvalidatordNFTCollectionBalance[] | null | undefined
): NFTCollectionBalance[] {
  return safeMap(dtos, nftCollectionBalanceFromDto);
}
