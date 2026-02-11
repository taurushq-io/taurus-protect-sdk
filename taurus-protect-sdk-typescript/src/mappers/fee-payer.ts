/**
 * Fee payer mapper functions for converting OpenAPI DTOs to domain models.
 */

import type {
  FeePayer,
  FeePayerEth,
  FeePayerEthLocal,
  FeePayerEthRemote,
  FeePayerInfo,
} from '../models/fee-payer';
import { safeBool, safeDate, safeMap, safeString } from './base';

/**
 * Maps an ETHLocal DTO to a FeePayerEthLocal domain model.
 */
export function feePayerEthLocalFromDto(dto: unknown): FeePayerEthLocal | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    addressId: safeString(d.addressId ?? d.address_id),
    forwarderAddressId: safeString(d.forwarderAddressId ?? d.forwarder_address_id),
    autoApprove: safeBool(d.autoApprove ?? d.auto_approve),
    creatorAddressId: safeString(d.creatorAddressId ?? d.creator_address_id),
    forwarderKind: safeString(d.forwarderKind ?? d.forwarder_kind),
    domainSeparator: safeString(d.domainSeparator ?? d.domain_separator),
  };
}

/**
 * Maps an ETHRemote DTO to a FeePayerEthRemote domain model.
 */
export function feePayerEthRemoteFromDto(dto: unknown): FeePayerEthRemote | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    url: safeString(d.url),
    username: safeString(d.username),
    fromAddressId: safeString(d.fromAddressId ?? d.from_address_id),
    forwarderAddress: safeString(d.forwarderAddress ?? d.forwarder_address),
    forwarderAddressId: safeString(d.forwarderAddressId ?? d.forwarder_address_id),
    creatorAddress: safeString(d.creatorAddress ?? d.creator_address),
    creatorAddressId: safeString(d.creatorAddressId ?? d.creator_address_id),
    forwarderKind: safeString(d.forwarderKind ?? d.forwarder_kind),
    domainSeparator: safeString(d.domainSeparator ?? d.domain_separator),
  };
}

/**
 * Maps a FeePayerETH DTO to a FeePayerEth domain model.
 */
export function feePayerEthFromDto(dto: unknown): FeePayerEth | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    kind: safeString(d.kind),
    local: feePayerEthLocalFromDto(d.local),
    remote: feePayerEthRemoteFromDto(d.remote),
    remoteEncrypted: safeString(d.remoteEncrypted ?? d.remote_encrypted),
  };
}

/**
 * Maps a TgvalidatordFeePayer DTO to a FeePayerInfo domain model.
 */
export function feePayerInfoFromDto(dto: unknown): FeePayerInfo | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    blockchain: safeString(d.blockchain),
    eth: feePayerEthFromDto(d.eth),
  };
}

/**
 * Maps a TgvalidatordFeePayerEnvelope DTO to a FeePayer domain model.
 */
export function feePayerFromDto(dto: unknown): FeePayer | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    id: safeString(d.id),
    tenantId: safeString(d.tenantId ?? d.tenant_id),
    blockchain: safeString(d.blockchain),
    network: safeString(d.network),
    name: safeString(d.name),
    creationDate: safeDate(d.creationDate ?? d.creation_date),
    feePayerInfo: feePayerInfoFromDto(d.feePayer ?? d.fee_payer),
  };
}

/**
 * Maps an array of fee payer DTOs to FeePayer domain models.
 */
export function feePayersFromDto(dtos: unknown[] | null | undefined): FeePayer[] {
  return safeMap(dtos, feePayerFromDto);
}
