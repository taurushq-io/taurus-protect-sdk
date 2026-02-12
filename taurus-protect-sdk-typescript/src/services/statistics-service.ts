/**
 * Statistics service for Taurus-PROTECT SDK.
 *
 * Provides methods for retrieving portfolio statistics.
 */

import type { StatisticsApi } from '../internal/openapi/apis/StatisticsApi';
import { portfolioStatisticsFromDto } from '../mappers/statistics';
import type { PortfolioStatistics } from '../models/statistics';
import { BaseService } from './base';

/**
 * Service for retrieving portfolio statistics in the Taurus-PROTECT system.
 *
 * Provides access to aggregated portfolio statistics including total balances,
 * address counts, and wallet counts.
 *
 * @example
 * ```typescript
 * // Get portfolio statistics
 * const stats = await client.statistics.getPortfolioStatistics();
 * console.log(`Total wallets: ${stats.walletsCount}`);
 * console.log(`Total addresses: ${stats.addressesCount}`);
 * console.log(`Total balance (base currency): ${stats.totalBalanceBaseCurrency}`);
 * ```
 */
export class StatisticsService extends BaseService {
  private readonly statisticsApi: StatisticsApi;

  /**
   * Creates a new StatisticsService instance.
   *
   * @param statisticsApi - The StatisticsApi instance from the OpenAPI client
   */
  constructor(statisticsApi: StatisticsApi) {
    super();
    this.statisticsApi = statisticsApi;
  }

  /**
   * Retrieves aggregated portfolio statistics.
   *
   * Returns summary statistics including total balance, number of wallets,
   * and number of addresses across the entire portfolio.
   *
   * @returns The portfolio statistics
   * @throws {@link APIError} If the API call fails
   *
   * @example
   * ```typescript
   * const stats = await statisticsService.getPortfolioStatistics();
   * console.log(`Wallets: ${stats.walletsCount}`);
   * console.log(`Addresses: ${stats.addressesCount}`);
   * console.log(`Total Balance: ${stats.totalBalance}`);
   * console.log(`Fiat Value: ${stats.totalBalanceBaseCurrency}`);
   * ```
   */
  async getPortfolioStatistics(): Promise<PortfolioStatistics> {
    return this.execute(async () => {
      const response = await this.statisticsApi.statisticsServiceGetPortfolioStatistics();
      const result = response.result;
      const stats = portfolioStatisticsFromDto(result);

      // Return empty stats object if no result (should not happen with valid API response)
      return stats ?? {};
    });
  }
}
