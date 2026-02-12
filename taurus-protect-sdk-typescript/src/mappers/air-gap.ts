/**
 * Air gap mapper functions for converting domain models to OpenAPI DTOs.
 *
 * Note: Air gap operations primarily map from domain models to DTOs (for requests),
 * rather than from DTOs to domain models (the response is a binary Blob or empty object).
 */

import type { TgvalidatordGetOutgoingAirGapRequest } from '../internal/openapi/models/TgvalidatordGetOutgoingAirGapRequest';
import type { TgvalidatordSubmitIncomingAirGapRequest } from '../internal/openapi/models/TgvalidatordSubmitIncomingAirGapRequest';
import type {
  GetOutgoingAirGapOptions,
  GetOutgoingAirGapAddressOptions,
  SubmitIncomingAirGapOptions,
} from '../models/air-gap';

/**
 * Maps GetOutgoingAirGapOptions to OpenAPI request DTO.
 *
 * @param options - The outgoing air gap options with request IDs
 * @returns The OpenAPI request DTO
 */
export function toGetOutgoingAirGapRequest(
  options: GetOutgoingAirGapOptions
): TgvalidatordGetOutgoingAirGapRequest {
  return {
    requests: {
      ids: options.requestIds,
      signature: options.signature,
    },
  };
}

/**
 * Maps GetOutgoingAirGapAddressOptions to OpenAPI request DTO.
 *
 * @param options - The outgoing air gap options with address IDs
 * @returns The OpenAPI request DTO
 */
export function toGetOutgoingAirGapAddressRequest(
  options: GetOutgoingAirGapAddressOptions
): TgvalidatordGetOutgoingAirGapRequest {
  return {
    addresses: {
      ids: options.addressIds,
    },
  };
}

/**
 * Maps SubmitIncomingAirGapOptions to OpenAPI request DTO.
 *
 * @param options - The incoming air gap options with signed payload
 * @returns The OpenAPI request DTO
 */
export function toSubmitIncomingAirGapRequest(
  options: SubmitIncomingAirGapOptions
): TgvalidatordSubmitIncomingAirGapRequest {
  return {
    payload: options.payload,
    signature: options.signature,
  };
}
