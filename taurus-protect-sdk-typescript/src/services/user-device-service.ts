/**
 * User device service for Taurus-PROTECT SDK.
 *
 * Provides methods for managing user device pairing in the Taurus-PROTECT system.
 * Device pairing is used to register and authenticate user devices for secure
 * access to the platform.
 */

import { ValidationError } from '../errors';
import type { UserDeviceApi } from '../internal/openapi/apis/UserDeviceApi';
import {
  userDevicePairingFromDto,
  userDevicePairingInfoFromDto,
} from '../mappers/user-device';
import type {
  ApprovePairingOptions,
  StartPairingOptions,
  UserDevicePairing,
  UserDevicePairingInfo,
} from '../models/user-device';
import { BaseService } from './base';

/**
 * Service for managing user device pairing.
 *
 * Device pairing follows a 3-step process:
 * 1. Create a pairing request (returns a pairing ID)
 * 2. Start the pairing process (with nonce and public key)
 * 3. Approve the pairing (finalizes the device registration)
 *
 * @example
 * ```typescript
 * // Step 1: Create a new device pairing
 * const pairing = await userDeviceService.createPairing();
 * console.log(`Pairing ID: ${pairing.pairingId}`);
 *
 * // Step 2: Start the pairing process
 * await userDeviceService.startPairing(pairing.pairingId, {
 *   nonce: '123456',
 *   publicKey: 'base64-encoded-public-key',
 * });
 *
 * // Check pairing status
 * const info = await userDeviceService.getPairingStatus(pairing.pairingId, '123456');
 * console.log(`Status: ${info.status}`);
 *
 * // Step 3: Approve the pairing
 * await userDeviceService.approvePairing(pairing.pairingId, { nonce: '123456' });
 *
 * // Get final status with API key
 * const finalInfo = await userDeviceService.getPairingStatus(pairing.pairingId, '123456');
 * if (finalInfo.apiKey) {
 *   console.log(`Device paired! API Key: ${finalInfo.apiKey}`);
 * }
 * ```
 */
export class UserDeviceService extends BaseService {
  private readonly userDeviceApi: UserDeviceApi;

  /**
   * Creates a new UserDeviceService instance.
   *
   * @param userDeviceApi - The UserDeviceApi instance from the OpenAPI client
   */
  constructor(userDeviceApi: UserDeviceApi) {
    super();
    this.userDeviceApi = userDeviceApi;
  }

  /**
   * Creates a new device pairing request.
   *
   * This is Step 1 of the pairing process. It creates a pairing request
   * for the current user and returns a pairing ID that should be used
   * in subsequent steps.
   *
   * @returns The created pairing with its ID
   * @throws {@link APIError} If the API request fails
   *
   * @example
   * ```typescript
   * const pairing = await userDeviceService.createPairing();
   * console.log(`Use pairing ID: ${pairing.pairingId}`);
   * ```
   */
  async createPairing(): Promise<UserDevicePairing> {
    return this.execute(async () => {
      const response = await this.userDeviceApi.userDeviceServiceCreateUserDevicePairing();
      const pairing = userDevicePairingFromDto(response);

      if (!pairing) {
        throw new ValidationError('Invalid response from create pairing API');
      }

      return pairing;
    });
  }

  /**
   * Gets the status of a device pairing request.
   *
   * Use this method to check the current state of a pairing process.
   * After approval, this will return the API key for the paired device.
   *
   * @param pairingId - The pairing ID from createPairing()
   * @param nonce - The nonce used to start the pairing (6 digits)
   * @returns The pairing status and info
   * @throws {@link ValidationError} If pairingId or nonce is empty
   * @throws {@link APIError} If the API request fails
   *
   * @example
   * ```typescript
   * const info = await userDeviceService.getPairingStatus(pairingId, '123456');
   * console.log(`Status: ${info.status}`);
   * if (info.apiKey) {
   *   console.log(`API Key: ${info.apiKey}`);
   * }
   * ```
   */
  async getPairingStatus(pairingId: string, nonce: string): Promise<UserDevicePairingInfo> {
    if (!pairingId || pairingId.trim() === '') {
      throw new ValidationError('pairingId is required');
    }
    if (!nonce || nonce.trim() === '') {
      throw new ValidationError('nonce is required');
    }

    return this.execute(async () => {
      const response = await this.userDeviceApi.userDeviceServiceGetUserDevicePairingStatus({
        pairingID: pairingId,
        nonce,
      });
      const info = userDevicePairingInfoFromDto(response);

      if (!info) {
        throw new ValidationError('Invalid response from get pairing status API');
      }

      return info;
    });
  }

  /**
   * Starts the device pairing process.
   *
   * This is Step 2 of the pairing process. It initiates the actual pairing
   * using the provided nonce and public key. The nonce and public key will
   * be used to validate future requests.
   *
   * @param pairingId - The pairing ID from createPairing()
   * @param options - The start pairing options (nonce and publicKey)
   * @throws {@link ValidationError} If pairingId is empty
   * @throws {@link APIError} If the API request fails
   *
   * @example
   * ```typescript
   * await userDeviceService.startPairing(pairingId, {
   *   nonce: '123456',
   *   publicKey: 'base64-encoded-public-key',
   * });
   * ```
   */
  async startPairing(pairingId: string, options: StartPairingOptions): Promise<void> {
    if (!pairingId || pairingId.trim() === '') {
      throw new ValidationError('pairingId is required');
    }

    return this.execute(async () => {
      await this.userDeviceApi.userDeviceServiceStartUserDevicePairing({
        pairingID: pairingId,
        body: {
          nonce: options.nonce,
          publicKey: options.publicKey,
        },
      });
    });
  }

  /**
   * Approves a device pairing request.
   *
   * This is Step 3 (final step) of the pairing process. It checks the
   * validity of the nonce and finalizes the device pairing. After approval,
   * the pairing status will include the API key for the device.
   *
   * @param pairingId - The pairing ID from createPairing()
   * @param options - The approve pairing options (nonce)
   * @throws {@link ValidationError} If pairingId or nonce is empty
   * @throws {@link APIError} If the API request fails
   *
   * @example
   * ```typescript
   * await userDeviceService.approvePairing(pairingId, { nonce: '123456' });
   *
   * // After approval, get the API key
   * const info = await userDeviceService.getPairingStatus(pairingId, '123456');
   * console.log(`API Key: ${info.apiKey}`);
   * ```
   */
  async approvePairing(pairingId: string, options: ApprovePairingOptions): Promise<void> {
    if (!pairingId || pairingId.trim() === '') {
      throw new ValidationError('pairingId is required');
    }
    if (!options.nonce || options.nonce.trim() === '') {
      throw new ValidationError('nonce is required');
    }

    return this.execute(async () => {
      await this.userDeviceApi.userDeviceServiceApproveUserDevicePairing({
        pairingID: pairingId,
        body: {
          nonce: options.nonce,
        },
      });
    });
  }
}
