/**
 * Address mapper for converting OpenAPI DTOs to domain models.
 */

import type { TgvalidatordAddress, TgvalidatordAddressAttribute } from "../internal/openapi";
import type { Address, AddressAttribute, Balance } from "../models/address";
import { safeString, safeStringDefault, safeBoolDefault, safeDate, safeMap } from "./base";

/**
 * Convert OpenAPI TgvalidatordAddressAttribute to domain AddressAttribute.
 *
 * @param dto - OpenAPI address attribute DTO
 * @returns Domain AddressAttribute model
 */
export function addressAttributeFromDto(dto: TgvalidatordAddressAttribute): AddressAttribute {
  return {
    id: safeStringDefault(dto.id, ""),
    key: safeStringDefault(dto.key, ""),
    value: safeStringDefault(dto.value, ""),
  };
}

/**
 * Convert OpenAPI TgvalidatordAddress to domain Address.
 *
 * @param dto - OpenAPI address DTO (TgvalidatordAddress)
 * @returns Domain Address model or undefined if dto is null/undefined
 */
export function addressFromDto(dto: TgvalidatordAddress | null | undefined): Address | undefined {
  if (dto == null) {
    return undefined;
  }

  // Extract attributes if present
  const attributes = safeMap(dto.attributes, addressAttributeFromDto);

  // Extract linked whitelisted address IDs
  const linkedWhitelistedAddressIds = dto.linkedWhitelistedAddressIds ?? [];

  return {
    id: safeStringDefault(dto.id, ""),
    walletId: safeStringDefault(dto.walletId, ""),
    address: safeStringDefault(dto.address, ""),
    alternateAddress: safeString(dto.alternateAddress),
    label: safeString(dto.label),
    comment: safeString(dto.comment),
    currency: safeStringDefault(dto.currency, ""),
    customerId: safeString(dto.customerId),
    externalAddressId: safeString(dto.externalAddressId),
    addressPath: safeString(dto.addressPath),
    addressIndex: safeString(dto.addressIndex),
    nonce: safeString(dto.nonce),
    status: safeString(dto.status),
    signature: safeString(dto.signature),
    disabled: safeBoolDefault(dto.disabled, false),
    canUseAllFunds: safeBoolDefault(dto.canUseAllFunds, false),
    createdAt: safeDate(dto.creationDate),
    updatedAt: safeDate(dto.updateDate),
    attributes,
    linkedWhitelistedAddressIds,
    balance: dto.balance ? {
      totalConfirmed: safeStringDefault(dto.balance.totalConfirmed, "0"),
      totalUnconfirmed: safeStringDefault(dto.balance.totalUnconfirmed, "0"),
      availableConfirmed: safeStringDefault(dto.balance.availableConfirmed, "0"),
      availableUnconfirmed: safeStringDefault(dto.balance.availableUnconfirmed, "0"),
      reservedConfirmed: safeStringDefault(dto.balance.reservedConfirmed, "0"),
      reservedUnconfirmed: safeStringDefault(dto.balance.reservedUnconfirmed, "0"),
    } : undefined,
  };
}

/**
 * Convert list of OpenAPI address DTOs to domain Addresses.
 *
 * @param dtos - List of OpenAPI address DTOs
 * @returns List of domain Address models (filters out undefined)
 */
export function addressesFromDto(
  dtos: TgvalidatordAddress[] | null | undefined
): Address[] {
  if (dtos == null) {
    return [];
  }
  const result: Address[] = [];
  for (const dto of dtos) {
    const address = addressFromDto(dto);
    if (address !== undefined) {
      result.push(address);
    }
  }
  return result;
}
