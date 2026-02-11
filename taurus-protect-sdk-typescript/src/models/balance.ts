/**
 * Balance models for Taurus-PROTECT SDK.
 */

/**
 * Asset balance representing the total balance for a specific asset.
 */
export interface AssetBalance {
  /** Currency ID */
  readonly currencyId?: string;
  /** Currency symbol (e.g., "ETH", "BTC") */
  readonly currency?: string;
  /** Blockchain type */
  readonly blockchain?: string;
  /** Network (e.g., "mainnet", "testnet") */
  readonly network?: string;
  /** Token contract address (for tokens) */
  readonly contractAddress?: string;
  /** Token ID (for multi-asset contracts) */
  readonly tokenId?: string;
  /** Total confirmed balance */
  readonly balance?: string;
  /** Fiat value of the balance */
  readonly fiatValue?: string;
  /** Fiat currency used for valuation */
  readonly fiatCurrency?: string;
}

/**
 * NFT collection balance.
 */
export interface NFTCollectionBalance {
  /** Collection name */
  readonly name?: string;
  /** Collection symbol */
  readonly symbol?: string;
  /** Blockchain type */
  readonly blockchain?: string;
  /** Network (e.g., "mainnet", "testnet") */
  readonly network?: string;
  /** Contract address */
  readonly contractAddress?: string;
  /** Number of NFTs owned */
  readonly count?: number;
  /** Logo URL */
  readonly logoUrl?: string;
}

/**
 * Options for listing balances.
 */
export interface ListBalancesOptions {
  /** Filter by currency ID or symbol */
  currency?: string;
  /** Maximum number of items to return */
  limit?: number;
}

/**
 * Options for listing NFT collection balances.
 */
export interface ListNFTCollectionBalancesOptions {
  /** Blockchain to filter by (required) */
  blockchain: string;
  /** Network to filter by (required) */
  network: string;
  /** Maximum number of items to return */
  limit?: number;
  /** Whether to exclude collections with zero balance */
  onlyPositiveBalance?: boolean;
}
