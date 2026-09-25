/**
 * Asset mapper functions for converting the v2 asset registry DTOs to domain models.
 */

import type { TgvalidatordAddressTargetV2 } from '../internal/openapi/models/TgvalidatordAddressTargetV2';
import type { TgvalidatordAssetAddressV2 } from '../internal/openapi/models/TgvalidatordAssetAddressV2';
import type { TgvalidatordAssetOperationV2 } from '../internal/openapi/models/TgvalidatordAssetOperationV2';
import type { TgvalidatordAssetResourceV2 } from '../internal/openapi/models/TgvalidatordAssetResourceV2';
import type {
  AssetAddressTargetV2,
  AssetAddressV2,
  AssetOperationV2,
  AssetV2,
} from '../models/asset';
import { safeMap } from './base';

function targetFromDto(
  dto: TgvalidatordAddressTargetV2 | undefined
): AssetAddressTargetV2 | undefined {
  if (!dto) {
    return undefined;
  }
  return { addressId: dto.addressID, whitelistedAddressId: dto.whitelistedAddressID };
}

/**
 * Maps a v2 asset registry row to an AssetV2 domain model.
 */
export function assetV2FromDto(
  dto: TgvalidatordAssetResourceV2 | null | undefined
): AssetV2 | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }
  const canton = dto.blockchainAsset?.cantonNativeTokenAsset;
  const configuration = canton?._configuration;
  return {
    id: dto.id,
    tenantId: dto.tenantID,
    version: dto.version,
    createdAt: dto.createdAt,
    updatedAt: dto.updatedAt,
    label: dto.label,
    assetType: dto.assetType,
    status: dto.status,
    blockchain: dto.blockchain,
    network: dto.network,
    currencyId: dto.currencyID,
    name: dto.name,
    symbol: dto.symbol,
    decimals: dto.decimals,
    contractAddress: dto.contractAddress,
    attributes: (dto.attributes ?? []).map((attr) => ({ key: attr.key, value: attr.value })),
    cantonInstrumentId: canton?.instrumentID,
    cantonConfiguration: configuration
      ? {
          cid: configuration.cid,
          requireCredentials: configuration.requireCredentials,
          paused: configuration.paused,
          operator: configuration.operator,
        }
      : undefined,
  };
}

/**
 * Maps v2 asset registry rows to AssetV2 domain models.
 */
export function assetsV2FromDto(
  dtos: TgvalidatordAssetResourceV2[] | null | undefined
): AssetV2[] {
  return safeMap(dtos, assetV2FromDto);
}

/**
 * Maps a v2 asset holder row to an AssetAddressV2 domain model, unverified: only
 * `AssetService.queryAssetAddresses` marks a row verified, after re-reading it.
 */
export function assetAddressV2FromDto(
  dto: TgvalidatordAssetAddressV2 | null | undefined
): AssetAddressV2 | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }
  return {
    address: dto.address,
    kycStatus: dto.kycStatus,
    balance: dto.balance,
    addressType: dto.addressType,
    addressId: dto.addressID,
    whitelistedAddressId: dto.whitelistedAddressID,
    verified: false,
  };
}

/**
 * Maps v2 asset holder rows to AssetAddressV2 domain models.
 */
export function assetAddressesV2FromDto(
  dtos: TgvalidatordAssetAddressV2[] | null | undefined
): AssetAddressV2[] {
  return safeMap(dtos, assetAddressV2FromDto);
}

/**
 * Maps a v2 asset operation row to an AssetOperationV2 domain model.
 */
export function assetOperationV2FromDto(
  dto: TgvalidatordAssetOperationV2 | null | undefined
): AssetOperationV2 | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }
  const create = dto.create;
  const update = dto.update;
  const imported = dto._import;
  return {
    id: dto.id,
    assetId: dto.assetID,
    type: dto.type,
    status: dto.status,
    createdAt: dto.createdAt,
    updatedAt: dto.updatedAt,
    initiatedByAddressId: dto.initiatedByAddressID,
    failureReason: dto.failureReason,
    blockingReason: dto.blockingReason,
    create: create
      ? {
          label: create.label,
          price: create.price,
          decimals: create.decimals,
          blockchain: create.blockchain,
          network: create.network,
          assetType: create.assetType,
          cantonInstrumentId: create.params?.cantonNativeTokenParams?.instrumentID,
          cantonName: create.params?.cantonNativeTokenParams?.name,
          cantonSymbol: create.params?.cantonNativeTokenParams?.symbol,
        }
      : undefined,
    update: update
      ? {
          label: update.label,
          price: update.price,
          cantonRequireCredentials:
            update.params?.cantonUpdateNativeTokenParams?.requireCredentials,
        }
      : undefined,
    import: imported
      ? {
          blockchain: imported.blockchain,
          network: imported.network,
          label: imported.label,
          price: imported.price,
          decimals: imported.decimals,
          address: imported.address,
        }
      : undefined,
    mint: dto.mint
      ? {
          destination: targetFromDto(dto.mint.destination),
          amount: dto.mint.amount,
          nftMetadata: dto.mint.nftMetadata,
        }
      : undefined,
    burn: dto.burn
      ? {
          destination: targetFromDto(dto.burn.destination),
          amount: dto.burn.amount,
          nftTokenIds: dto.burn.nftTokenIDs,
        }
      : undefined,
    pauseAccount: dto.pauseAccount
      ? { target: targetFromDto(dto.pauseAccount.target) }
      : undefined,
    unpauseAccount: dto.unpauseAccount
      ? { target: targetFromDto(dto.unpauseAccount.target) }
      : undefined,
    setKyc: dto.setKyc
      ? { target: targetFromDto(dto.setKyc.target), status: dto.setKyc.status }
      : undefined,
  };
}

/**
 * Maps v2 asset operation rows to AssetOperationV2 domain models.
 */
export function assetOperationsV2FromDto(
  dtos: TgvalidatordAssetOperationV2[] | null | undefined
): AssetOperationV2[] {
  return safeMap(dtos, assetOperationV2FromDto);
}
