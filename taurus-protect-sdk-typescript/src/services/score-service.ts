/**
 * Score service for Taurus-PROTECT SDK.
 *
 * Provides methods for managing blockchain analytics scores.
 */

import { ValidationError } from '../errors';
import type { ScoresApi } from '../internal/openapi/apis/ScoresApi';
import { scoresFromDto } from '../mappers/score';
import type { Score } from '../models/score';
import { BaseService } from './base';

/**
 * Service for managing blockchain analytics scores in the Taurus-PROTECT system.
 *
 * Scores represent risk assessments provided by third-party blockchain analytics
 * providers (e.g., Chainalysis, Elliptic). This service allows refreshing scores
 * for internal addresses and whitelisted external addresses.
 *
 * @example
 * ```typescript
 * // Refresh score for an internal address
 * const scores = await scoreService.refreshAddressScore(addressId, 'chainalysis');
 *
 * // Refresh score for a whitelisted external address
 * const wlaScores = await scoreService.refreshWhitelistedAddressScore(
 *   whitelistedAddressId,
 *   'elliptic'
 * );
 * ```
 */
export class ScoreService extends BaseService {
  private readonly scoresApi: ScoresApi;

  /**
   * Creates a new ScoreService instance.
   *
   * @param scoresApi - The ScoresApi instance from the OpenAPI client
   */
  constructor(scoresApi: ScoresApi) {
    super();
    this.scoresApi = scoresApi;
  }

  /**
   * Refreshes the compliance scores for an internal address.
   *
   * This endpoint triggers a new score calculation from the specified
   * blockchain analytics provider for the given address.
   *
   * @param addressId - The internal address ID (must be positive)
   * @param scoreProvider - The score provider name (e.g., "chainalysis", "elliptic")
   * @returns Array of refreshed scores from the provider
   * @throws {@link ValidationError} If addressId is not positive or scoreProvider is empty
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const scores = await scoreService.refreshAddressScore(123, 'chainalysis');
   * for (const score of scores) {
   *   console.log(`Score: ${score.score}, Type: ${score.type}`);
   * }
   * ```
   */
  async refreshAddressScore(addressId: number, scoreProvider: string): Promise<Score[]> {
    if (!addressId || addressId <= 0) {
      throw new ValidationError('addressId must be positive');
    }
    if (!scoreProvider || scoreProvider.trim() === '') {
      throw new ValidationError('scoreProvider is required');
    }

    return this.execute(async () => {
      const response = await this.scoresApi.scoreServiceRefreshAddressScore({
        addressId: String(addressId),
        body: {
          scoreProvider,
        },
      });

      const result =
        (response as Record<string, unknown>).scores ??
        (response as Record<string, unknown>).result;
      return scoresFromDto(result as unknown[]);
    });
  }

  /**
   * Refreshes the compliance scores for a whitelisted external address.
   *
   * This endpoint triggers a new score calculation from the specified
   * blockchain analytics provider for the given whitelisted address.
   *
   * @param whitelistedAddressId - The whitelisted address ID (must be positive)
   * @param scoreProvider - The score provider name (e.g., "chainalysis", "elliptic")
   * @returns Array of refreshed scores from the provider
   * @throws {@link ValidationError} If whitelistedAddressId is not positive or scoreProvider is empty
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const scores = await scoreService.refreshWhitelistedAddressScore(456, 'elliptic');
   * for (const score of scores) {
   *   console.log(`Provider: ${score.provider}, Score: ${score.score}`);
   * }
   * ```
   */
  async refreshWhitelistedAddressScore(
    whitelistedAddressId: number,
    scoreProvider: string
  ): Promise<Score[]> {
    if (!whitelistedAddressId || whitelistedAddressId <= 0) {
      throw new ValidationError('whitelistedAddressId must be positive');
    }
    if (!scoreProvider || scoreProvider.trim() === '') {
      throw new ValidationError('scoreProvider is required');
    }

    return this.execute(async () => {
      const response = await this.scoresApi.scoreServiceRefreshWLAScore({
        addressId: String(whitelistedAddressId),
        body: {
          scoreProvider,
        },
      });

      const result =
        (response as Record<string, unknown>).scores ??
        (response as Record<string, unknown>).result;
      return scoresFromDto(result as unknown[]);
    });
  }
}
