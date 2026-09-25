/**
 * Multi-factor signature mappers.
 */

import type { TgvalidatordGetMultiFactorSignatureEntitiesInfoReply } from '../internal/openapi/models/TgvalidatordGetMultiFactorSignatureEntitiesInfoReply';
import type { TgvalidatordMultiFactorSignaturesEntityType } from '../internal/openapi/models/TgvalidatordMultiFactorSignaturesEntityType';
import type {
  MultiFactorSignatureEntityType,
  MultiFactorSignatureInfo,
} from '../models/multi-factor-signature';

/**
 * Map entity type from DTO.
 *
 * The wire value is passed through verbatim, like in every other SDK: a kind this client does
 * not know stays that raw string instead of being read as a known kind.
 */
export function multiFactorSignatureEntityTypeFromDto(
  dto: TgvalidatordMultiFactorSignaturesEntityType | undefined | null
): MultiFactorSignatureEntityType | undefined {
  return dto == null ? undefined : (dto as MultiFactorSignatureEntityType);
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
