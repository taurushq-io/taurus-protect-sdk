/**
 * Balance models for Taurus-PROTECT SDK.
 */

import type { CursorPage, CursorPageOptions } from './pagination';

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
 * Options for listing balances. A cursor list: `pageSize` 1-100 (default 20) and the
 * `cursor` of a previous page.
 */
export interface ListBalancesOptions extends CursorPageOptions {
  /** Filter by currency ID or symbol */
  readonly currency?: string;
  /** Filter by token ID */
  readonly tokenId?: string;
}

/**
 * A page of tenant asset balances.
 */
export interface ListBalancesResult {
  /** The asset balances of this page */
  readonly items: AssetBalance[];
  /** Cursor pagination, including the server's total */
  readonly pagination: CursorPage;
}

/**
 * Options for listing NFT collection balances. A cursor list: `pageSize` 1-100
 * (default 20) and the `cursor` of a previous page.
 */
export interface ListNFTCollectionBalancesOptions extends CursorPageOptions {
  /** Filter by blockchain */
  readonly blockchain?: string;
  /** Filter by network */
  readonly network?: string;
  /** Search the collection name or symbol */
  readonly query?: string;
  /** Whether to exclude collections with zero balance */
  readonly onlyPositiveBalance?: boolean;
}

/**
 * A page of NFT collection balances.
 */
export interface ListNFTCollectionBalancesResult {
  /** The NFT collection balances of this page */
  readonly items: NFTCollectionBalance[];
  /** Cursor pagination */
  readonly pagination: CursorPage;
}

/**
 * Options for listing the tokens a wallet holds. A token list: `pageSize` 1-100
 * (default 20) and the `cursor` of a previous page.
 */
export type ListWalletTokensOptions = CursorPageOptions;

/**
 * A page of the tokens a wallet holds.
 */
export interface ListWalletTokensResult {
  /** The asset balances of this page */
  readonly items: AssetBalance[];
  /** Cursor pagination, including the server's total */
  readonly pagination: CursorPage;
}
