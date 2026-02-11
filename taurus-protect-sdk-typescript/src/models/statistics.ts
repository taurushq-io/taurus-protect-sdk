/**
 * Statistics models for Taurus-PROTECT SDK.
 */

/**
 * Portfolio statistics containing aggregated information about the portfolio.
 *
 * Includes total balance, address counts, wallet counts, and other aggregated metrics.
 */
export interface PortfolioStatistics {
  /**
   * Average balance per address in the smallest currency unit.
   * Example: 1500000000000000000 WEI corresponds to 1.5 ETH (ETH has 18 decimals).
   */
  readonly avgBalancePerAddress?: string;

  /** Total number of addresses in the portfolio */
  readonly addressesCount?: string;

  /** Total number of wallets in the portfolio */
  readonly walletsCount?: string;

  /**
   * Total balance in the smallest currency unit.
   * Example: 1500000000000000000 WEI corresponds to 1.5 ETH (ETH has 18 decimals).
   */
  readonly totalBalance?: string;

  /**
   * Total balance converted to the base currency (fiat currency like CHF, EUR, USD).
   * Value is in the main unit of the base currency.
   */
  readonly totalBalanceBaseCurrency?: string;
}
