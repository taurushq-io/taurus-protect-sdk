/**
 * Currency models for Taurus-PROTECT SDK.
 */

/**
 * Currency information.
 */
export interface Currency {
  /** Unique currency identifier */
  readonly id?: string;
  /** Currency symbol (e.g., "ETH", "BTC") */
  readonly symbol?: string;
  /** Currency name */
  readonly name?: string;
  /** Display name for UI presentation */
  readonly displayName?: string;
  /** Currency type */
  readonly type?: string;
  /** Blockchain type */
  readonly blockchain?: string;
  /** Network (e.g., "mainnet", "testnet") */
  readonly network?: string;
  /** Number of decimal places */
  readonly decimals?: number;
  /** BIP-44 coin type index */
  readonly coinTypeIndex?: string;
  /** Token contract address (for tokens) */
  readonly contractAddress?: string;
  /** Token ID (for multi-asset contracts) */
  readonly tokenId?: string;
  /** WLCA identifier */
  readonly wlcaId?: number;
  /** Logo URL or path */
  readonly logo?: string;
  /** Whether this is a token (not the native currency) */
  readonly isToken: boolean;
  /** Whether this is an ERC-20 token */
  readonly isERC20: boolean;
  /** Whether this is an FA1.2 token (Tezos) */
  readonly isFA12: boolean;
  /** Whether this is an FA2.0 token (Tezos) */
  readonly isFA20: boolean;
  /** Whether this is an NFT */
  readonly isNFT: boolean;
  /** Whether this currency uses UTXO-based transactions */
  readonly isUTXOBased: boolean;
  /** Whether this currency uses account-based transactions */
  readonly isAccountBased: boolean;
  /** Whether this is a fiat currency */
  readonly isFiat: boolean;
  /** Whether this currency supports staking */
  readonly hasStaking: boolean;
  /** Whether this currency is enabled */
  readonly enabled: boolean;
  /** Whether this is the native currency of the blockchain (backward compat) */
  readonly isNative?: boolean;
  /** Whether this currency is disabled (backward compat) */
  readonly isDisabled?: boolean;
  /** Logo URL (backward compat, prefer logo) */
  readonly logoUrl?: string;
  /** Price in base currency (backward compat) */
  readonly price?: string;
  /** Base currency used for pricing (backward compat) */
  readonly priceCurrency?: string;
}

/**
 * Options for listing currencies.
 */
export interface ListCurrenciesOptions {
  /** Include disabled currencies */
  showDisabled?: boolean;
  /** Include logo URLs in response */
  includeLogo?: boolean;
}

/**
 * Options for getting a currency by blockchain.
 */
export interface GetCurrencyByBlockchainOptions {
  /** Blockchain type (required) */
  blockchain: string;
  /** Network (required) */
  network: string;
  /** Token contract address (optional, returns native currency if not specified) */
  contractAddress?: string;
  /** Token ID (optional, for multi-asset contracts) */
  tokenId?: string;
}
