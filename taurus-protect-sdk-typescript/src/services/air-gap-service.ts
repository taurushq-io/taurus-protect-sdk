/**
 * Air gap service for Taurus-PROTECT SDK.
 *
 * Provides operations for transferring data to and from a cold HSM
 * (Hardware Security Module) in an air-gapped environment.
 */

import { ValidationError } from '../errors';
import type { AirGapApi } from '../internal/openapi/apis/AirGapApi';
import {
  toGetOutgoingAirGapRequest,
  toGetOutgoingAirGapAddressRequest,
  toSubmitIncomingAirGapRequest,
} from '../mappers/air-gap';
import type {
  GetOutgoingAirGapOptions,
  GetOutgoingAirGapAddressOptions,
  SubmitIncomingAirGapOptions,
} from '../models/air-gap';
import { BaseService } from './base';

/**
 * Service for air gap operations in the Taurus-PROTECT system.
 *
 * This service provides operations for transferring data to and from a cold HSM
 * (Hardware Security Module) in an air-gapped environment.
 *
 * @example
 * ```typescript
 * // Export requests for cold HSM signing
 * const payload = await client.airGap.getOutgoingAirGap({
 *   requestIds: ['request-1', 'request-2'],
 * });
 *
 * // Transfer payload to cold HSM, get signed response, then submit
 * await client.airGap.submitIncomingAirGap({
 *   payload: signedPayloadBase64,
 * });
 * ```
 */
export class AirGapService extends BaseService {
  private readonly airGapApi: AirGapApi;

  /**
   * Creates a new AirGapService instance.
   *
   * @param airGapApi - The AirGapApi instance from the OpenAPI client
   */
  constructor(airGapApi: AirGapApi) {
    super();
    this.airGapApi = airGapApi;
  }

  /**
   * Exports HSM-ready requests for cold HSM signing.
   *
   * This endpoint returns the payload to be transmitted to the cold HSM
   * for offline signing.
   *
   * @param options - Options containing the request IDs to export
   * @returns A Blob containing the payload for the cold HSM
   * @throws {@link ValidationError} If requestIds is empty
   * @throws {@link APIError} If the API request fails
   *
   * @example
   * ```typescript
   * const payload = await airGapService.getOutgoingAirGap({
   *   requestIds: ['123', '456'],
   * });
   *
   * // Save the blob to a file for air-gap transfer
   * const arrayBuffer = await payload.arrayBuffer();
   * fs.writeFileSync('hsm-payload.bin', Buffer.from(arrayBuffer));
   * ```
   */
  async getOutgoingAirGap(options: GetOutgoingAirGapOptions): Promise<Blob> {
    if (!options.requestIds || options.requestIds.length === 0) {
      throw new ValidationError('requestIds cannot be empty');
    }

    return this.execute(async () => {
      const request = toGetOutgoingAirGapRequest(options);
      return this.airGapApi.airGapServiceGetOutgoingAirGap({ body: request });
    });
  }

  /**
   * Exports addresses for cold HSM signing.
   *
   * This endpoint returns the payload containing addresses to be signed
   * by the cold HSM for offline signing.
   *
   * @param options - Options containing the address IDs to export
   * @returns A Blob containing the payload for the cold HSM
   * @throws {@link ValidationError} If addressIds is empty
   * @throws {@link APIError} If the API request fails
   *
   * @example
   * ```typescript
   * const payload = await airGapService.getOutgoingAirGapAddresses({
   *   addressIds: ['addr-1', 'addr-2'],
   * });
   * ```
   */
  async getOutgoingAirGapAddresses(
    options: GetOutgoingAirGapAddressOptions
  ): Promise<Blob> {
    if (!options.addressIds || options.addressIds.length === 0) {
      throw new ValidationError('addressIds cannot be empty');
    }

    return this.execute(async () => {
      const request = toGetOutgoingAirGapAddressRequest(options);
      return this.airGapApi.airGapServiceGetOutgoingAirGap({ body: request });
    });
  }

  /**
   * Imports signed requests from the cold HSM.
   *
   * This endpoint accepts an envelope of signed requests from the cold HSM
   * after offline signing.
   *
   * @param options - Options containing the signed payload from the cold HSM
   * @throws {@link ValidationError} If payload is empty
   * @throws {@link APIError} If the API request fails
   *
   * @example
   * ```typescript
   * // Read the signed payload from the cold HSM
   * const signedPayload = fs.readFileSync('signed-payload.bin');
   * const payloadBase64 = signedPayload.toString('base64');
   *
   * await airGapService.submitIncomingAirGap({
   *   payload: payloadBase64,
   * });
   * ```
   */
  async submitIncomingAirGap(options: SubmitIncomingAirGapOptions): Promise<void> {
    if (!options.payload || options.payload.trim() === '') {
      throw new ValidationError('payload cannot be empty');
    }

    return this.execute(async () => {
      const request = toSubmitIncomingAirGapRequest(options);
      await this.airGapApi.airGapServiceSubmitIncomingAirGap({ body: request });
    });
  }
}
