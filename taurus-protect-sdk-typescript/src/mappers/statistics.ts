/**
 * Statistics mapper functions for converting OpenAPI DTOs to domain models.
 */

import type { PortfolioStatistics } from '../models/statistics';
import { safeString } from './base';

/**
 * Maps an aggregated stats data DTO to a PortfolioStatistics domain model.
 *
 * @param dto - The OpenAPI DTO (TgvalidatordAggregatedStatsData or compatible)
 * @returns The domain model, or undefined if dto is null/undefined
 */
export function portfolioStatisticsFromDto(dto: unknown): PortfolioStatistics | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    avgBalancePerAddress: safeString(d.avgBalancePerAddress ?? d.avg_balance_per_address),
    addressesCount: safeString(d.addressesCount ?? d.addresses_count),
    walletsCount: safeString(d.walletsCount ?? d.wallets_count),
    totalBalance: safeString(d.totalBalance ?? d.total_balance),
    totalBalanceBaseCurrency: safeString(d.totalBalanceBaseCurrency ?? d.total_balance_base_currency),
  };
}
