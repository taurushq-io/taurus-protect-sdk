/**
 * User device mapper functions for converting OpenAPI DTOs to domain models.
 */

import type {
  UserDevicePairing,
  UserDevicePairingInfo,
  UserDevicePairingStatus,
} from '../models/user-device';
import { safeString } from './base';

/**
 * Maps a CreateUserDevicePairingReply DTO to a UserDevicePairing domain model.
 *
 * @param dto - The DTO from the OpenAPI client
 * @returns The domain model
 */
export function userDevicePairingFromDto(dto: unknown): UserDevicePairing | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  const pairingId = safeString(d.pairingID ?? d.pairingId ?? d.pairing_id);

  if (!pairingId) {
    return undefined;
  }

  return {
    pairingId,
  };
}

/**
 * Maps a UserDevicePairingInfo DTO to a UserDevicePairingInfo domain model.
 *
 * @param dto - The DTO from the OpenAPI client
 * @returns The domain model
 */
export function userDevicePairingInfoFromDto(dto: unknown): UserDevicePairingInfo | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  const pairingId = safeString(d.pairingID ?? d.pairingId ?? d.pairing_id);
  const statusValue = d.status;

  if (!pairingId || statusValue === undefined) {
    return undefined;
  }

  // Extract status value - handle both string and object formats
  let status: UserDevicePairingStatus;
  if (typeof statusValue === 'string') {
    status = statusValue as UserDevicePairingStatus;
  } else if (typeof statusValue === 'object' && statusValue !== null) {
    // Handle case where status might be an enum object
    status = String(statusValue) as UserDevicePairingStatus;
  } else {
    return undefined;
  }

  return {
    pairingId,
    status,
    apiKey: safeString(d.apiKey ?? d.api_key),
  };
}
