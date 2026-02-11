/**
 * Multi-factor signature mappers.
 */

import type { TgvalidatordGetMultiFactorSignatureEntitiesInfoReply } from '../internal/openapi/models/TgvalidatordGetMultiFactorSignatureEntitiesInfoReply';
import type { TgvalidatordMultiFactorSignaturesEntityType } from '../internal/openapi/models/TgvalidatordMultiFactorSignaturesEntityType';
import {
  MultiFactorSignatureEntityType,
  type MultiFactorSignatureInfo,
} from '../models/multi-factor-signature';

/**
 * Map entity type from DTO.
 */
export function multiFactorSignatureEntityTypeFromDto(
  dto: TgvalidatordMultiFactorSignaturesEntityType | undefined | null
): MultiFactorSignatureEntityType {
  switch (dto) {
    case 'REQUEST':
      return MultiFactorSignatureEntityType.REQUEST;
    case 'WHITELISTED_ADDRESS':
      return MultiFactorSignatureEntityType.WHITELISTED_ADDRESS;
    case 'WHITELISTED_CONTRACT':
      return MultiFactorSignatureEntityType.WHITELISTED_CONTRACT;
    default:
      return MultiFactorSignatureEntityType.REQUEST;
  }
}

/**
 * Map entity type to DTO.
 */
export function multiFactorSignatureEntityTypeToDto(
  entityType: MultiFactorSignatureEntityType
): TgvalidatordMultiFactorSignaturesEntityType {
  return entityType as TgvalidatordMultiFactorSignaturesEntityType;
}

/**
 * Map multi-factor signature info from DTO.
 */
export function multiFactorSignatureInfoFromDto(
  dto: TgvalidatordGetMultiFactorSignatureEntitiesInfoReply | undefined | null
): MultiFactorSignatureInfo | null {
  if (!dto) {
    return null;
  }
  return {
    id: dto.id ?? '',
    payloadToSign: dto.payloadToSign ?? [],
    entityType: multiFactorSignatureEntityTypeFromDto(dto.entityType),
  };
}
