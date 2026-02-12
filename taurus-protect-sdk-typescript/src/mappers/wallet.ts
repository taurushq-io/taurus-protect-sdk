/**
 * Wallet mapper functions for converting OpenAPI DTOs to domain models.
 */

import type { TgvalidatordBalance } from '../internal/openapi/models/TgvalidatordBalance';
import type { TgvalidatordBalanceHistoryPoint } from '../internal/openapi/models/TgvalidatordBalanceHistoryPoint';
import type { TgvalidatordWallet } from '../internal/openapi/models/TgvalidatordWallet';
import type { TgvalidatordWalletAttribute } from '../internal/openapi/models/TgvalidatordWalletAttribute';
import type { TgvalidatordWalletInfo } from '../internal/openapi/models/TgvalidatordWalletInfo';
import type {
  BalanceHistoryPoint,
  Wallet,
  WalletAttribute,
  WalletBalance,
} from '../models/wallet';
import { WalletStatus, WalletType } from '../models/wallet';
import { currencyFromDto } from './currency';
import {
  safeBool,
  safeBoolDefault,
  safeDate,
  safeInt,
  safeMap,
  safeString,
  safeStringDefault,
} from './base';

/**
 * Maps a TgvalidatordBalance DTO to a WalletBalance domain model.
 *
 * @param dto - The balance DTO from the OpenAPI response
 * @returns The WalletBalance domain model, or undefined if dto is null/undefined
 */
export function walletBalanceFromDto(
  dto: TgvalidatordBalance | null | undefined
): WalletBalance | undefined {
  if (!dto) {
    return undefined;
  }
  return {
    totalConfirmed: safeString(dto.totalConfirmed),
    totalUnconfirmed: safeString(dto.totalUnconfirmed),
    availableConfirmed: safeString(dto.availableConfirmed),
    availableUnconfirmed: safeString(dto.availableUnconfirmed),
    reservedConfirmed: safeString(dto.reservedConfirmed),
    reservedUnconfirmed: safeString(dto.reservedUnconfirmed),
  };
}

/**
 * Maps a TgvalidatordWalletAttribute DTO to a WalletAttribute domain model.
 *
 * @param dto - The attribute DTO from the OpenAPI response
 * @returns The WalletAttribute domain model, or undefined if dto is null/undefined
 */
export function walletAttributeFromDto(
  dto: TgvalidatordWalletAttribute | null | undefined
): WalletAttribute | undefined {
  if (!dto) {
    return undefined;
  }
  return {
    id: safeString(dto.id),
    key: safeString(dto.key),
    value: safeString(dto.value),
    contentType: safeString(dto.contentType),
    owner: safeString(dto.owner),
    type: safeString(dto.type),
    subtype: safeString(dto.subtype),
    isFile: safeBool(dto.isfile),
  };
}

/**
 * Determines wallet status from DTO fields.
 *
 * @param disabled - Whether the wallet is disabled
 * @returns The appropriate WalletStatus enum value
 */
function determineWalletStatus(disabled: boolean | undefined): WalletStatus {
  if (disabled === true) {
    return WalletStatus.DISABLED;
  }
  // Default to ACTIVE if not explicitly disabled
  return WalletStatus.ACTIVE;
}

/**
 * Maps a TgvalidatordWalletInfo DTO to a Wallet domain model.
 * This is the response type from GetWalletV2 and GetWalletsV2.
 *
 * @param dto - The wallet info DTO from the OpenAPI response
 * @returns The Wallet domain model, or undefined if dto is null/undefined
 */
export function walletFromDto(
  dto: TgvalidatordWalletInfo | null | undefined
): Wallet | undefined {
  if (!dto) {
    return undefined;
  }

  return {
    id: safeStringDefault(dto.id, ''),
    name: safeStringDefault(dto.name, ''),
    status: determineWalletStatus(dto.disabled),
    type: WalletType.STANDARD, // Default type, not provided in API response
    blockchain: safeString(dto.blockchain),
    network: safeString(dto.network),
    currency: safeString(dto.currency),
    disabled: safeBoolDefault(dto.disabled, false),
    isOmnibus: safeBoolDefault(dto.isOmnibus, false),
    balance: walletBalanceFromDto(dto.balance),
    createdAt: safeDate(dto.creationDate),
    updatedAt: safeDate(dto.updateDate),
    comment: safeString(dto.comment),
    customerId: safeString(dto.customerId),
    addressesCount: safeInt(dto.addressesCount),
    attributes: safeMap(dto.attributes, walletAttributeFromDto),
    visibilityGroupId: safeString(dto.visibilityGroupID),
    externalWalletId: safeString(dto.externalWalletId),
    accountPath: safeString(dto.accountPath),
    currencyInfo: dto.currencyInfo ? currencyFromDto(dto.currencyInfo) : undefined,
    tags: [], // Tags not included in WalletInfo response
  };
}

/**
 * Maps a TgvalidatordWallet DTO to a Wallet domain model.
 * This is the response type from CreateWallet.
 *
 * @param dto - The wallet DTO from the OpenAPI response
 * @returns The Wallet domain model, or undefined if dto is null/undefined
 */
export function walletFromCreateDto(
  dto: TgvalidatordWallet | null | undefined
): Wallet | undefined {
  if (!dto) {
    return undefined;
  }

  return {
    id: safeStringDefault(dto.id, ''),
    name: safeStringDefault(dto.name, ''),
    status: determineWalletStatus(dto.disabled),
    type: WalletType.STANDARD, // Default type, not provided in API response
    blockchain: safeString(dto.blockchain),
    network: undefined, // Not provided in TgvalidatordWallet
    currency: safeString(dto.currency),
    disabled: safeBoolDefault(dto.disabled, false),
    isOmnibus: safeBoolDefault(dto.isOmnibus, false),
    balance: walletBalanceFromDto(dto.balance),
    createdAt: safeDate(dto.creationDate),
    updatedAt: safeDate(dto.updateDate),
    comment: safeString(dto.comment),
    customerId: safeString(dto.customerId),
    addressesCount: safeInt(dto.addressesCount),
    attributes: safeMap(dto.attributes, walletAttributeFromDto),
    visibilityGroupId: undefined, // Not provided in TgvalidatordWallet
    externalWalletId: safeString(dto.externalWalletId),
    accountPath: safeString(dto.accountPath),
    currencyInfo: dto.currencyInfo ? currencyFromDto(dto.currencyInfo) : undefined,
    tags: [], // Tags not included in Wallet response
  };
}

/**
 * Maps an array of TgvalidatordWalletInfo DTOs to Wallet domain models.
 *
 * @param dtos - The array of wallet info DTOs
 * @returns Array of Wallet domain models (undefined entries filtered out)
 */
export function walletsFromDto(
  dtos: TgvalidatordWalletInfo[] | null | undefined
): Wallet[] {
  return safeMap(dtos, walletFromDto);
}

/**
 * Maps a TgvalidatordBalanceHistoryPoint DTO to a BalanceHistoryPoint domain model.
 *
 * @param dto - The balance history point DTO from the OpenAPI response
 * @returns The BalanceHistoryPoint domain model, or undefined if dto is null/undefined
 */
export function balanceHistoryPointFromDto(
  dto: TgvalidatordBalanceHistoryPoint | null | undefined
): BalanceHistoryPoint | undefined {
  if (!dto) {
    return undefined;
  }
  return {
    pointDate: safeDate(dto.pointDate),
    balance: walletBalanceFromDto(dto.balance),
  };
}

/**
 * Maps an array of TgvalidatordBalanceHistoryPoint DTOs to BalanceHistoryPoint domain models.
 *
 * @param dtos - The array of balance history point DTOs
 * @returns Array of BalanceHistoryPoint domain models (undefined entries filtered out)
 */
export function balanceHistoryFromDto(
  dtos: TgvalidatordBalanceHistoryPoint[] | null | undefined
): BalanceHistoryPoint[] {
  return safeMap(dtos, balanceHistoryPointFromDto);
}
