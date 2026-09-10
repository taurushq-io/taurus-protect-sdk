/**
 * Contract whitelisting mappers.
 *
 * Attribute mappers only. The envelope and model mappers that lived here built a
 * WhitelistedContract from the DTO — and parsed the signed payload as fact via
 * extractContractDetailsFromPayload — with no verification of any kind, which made
 * WhitelistedAssetService's six-step chain avoidable for the same entity. Do not
 * reintroduce a mapper that takes a bare contract envelope.
 */

import { safeBool, safeMap, safeString } from './base';
import type { WhitelistedContractAttribute } from '../models/contract-whitelist';

/**
 * Maps a whitelisted contract attribute from DTO.
 */
export function whitelistedContractAttributeFromDto(
  dto: unknown
): WhitelistedContractAttribute | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    id: safeString(d.id),
    key: safeString(d.key),
    value: safeString(d.value),
    contentType: safeString(d.contentType ?? d.content_type),
    type: safeString(d.type),
    subType: safeString(d.subtype ?? d.subType ?? d.sub_type),
    isFile: safeBool(d.isfile ?? d.isFile ?? d.is_file),
  };
}

/**
 * Maps an array of whitelisted contract attributes from DTOs.
 */
export function whitelistedContractAttributesFromDto(
  dtos: unknown[] | null | undefined
): WhitelistedContractAttribute[] {
  return safeMap(dtos, whitelistedContractAttributeFromDto);
}
