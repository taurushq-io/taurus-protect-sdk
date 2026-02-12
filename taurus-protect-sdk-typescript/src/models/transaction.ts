/**
 * Transaction domain models for Taurus-PROTECT SDK.
 */

/**
 * Transaction status enum.
 */
export enum TransactionStatus {
  NEW = "NEW",
  SUCCESS = "SUCCESS",
  PENDING = "PENDING",
  TEMPORARY_FAILURE = "TEMPORARY_FAILURE",
  INVALID = "INVALID",
  TIMEOUT = "TIMEOUT",
  MINED = "MINED",
  LOST = "LOST",
  EXPIRED = "EXPIRED",
}

/**
 * Transaction direction enum.
 */
export enum TransactionDirection {
  INCOMING = "incoming",
  OUTGOING = "outgoing",
}

/**
 * Transaction type enum.
 */
export enum TransactionType {
  TRANSFER = "transfer",
  BURN = "burn",
  APPROVE = "approve",
  STAKE = "stake",
  UNSTAKE = "unstake",
}

/**
 * Represents address information within a transaction.
 */
export interface TransactionAddressInfo {
  /** The blockchain address */
  readonly address: string | undefined;
  /** Label of the address if known to Taurus-PROTECT */
  readonly label: string | undefined;
  /** Container identifier */
  readonly container: string | undefined;
  /** Customer identifier */
  readonly customerId: string | undefined;
  /** Amount in smallest currency unit */
  readonly amount: string | undefined;
  /** Amount in main currency unit */
  readonly amountMainUnit: string | undefined;
  /** Type of the transaction element (source or destination) */
  readonly type: string | undefined;
  /** Index within transaction */
  readonly idx: string | undefined;
  /** Internal address ID if this is an internal address */
  readonly internalAddressId: string | undefined;
  /** Whitelisted address ID if this is a whitelisted address */
  readonly whitelistedAddressId: string | undefined;
}

/**
 * Currency information for a transaction.
 */
export interface TransactionCurrencyInfo {
  /** Currency ID */
  readonly id: string | undefined;
  /** Currency symbol (e.g., "ETH", "BTC") */
  readonly symbol: string | undefined;
  /** Currency name */
  readonly name: string | undefined;
  /** Number of decimal places */
  readonly decimals: number | undefined;
  /** Blockchain identifier */
  readonly blockchain: string | undefined;
  /** Network identifier */
  readonly network: string | undefined;
}

/**
 * Transaction attribute (key-value metadata).
 */
export interface TransactionAttribute {
  /** Attribute key */
  readonly key: string | undefined;
  /** Attribute value */
  readonly value: string | undefined;
}

/**
 * Represents a blockchain transaction in Taurus-PROTECT.
 *
 * Transactions represent the movement of cryptocurrency on the blockchain
 * and can be either incoming (received) or outgoing (sent).
 */
export interface Transaction {
  /** Unique transaction identifier */
  readonly id: string;
  /** Transaction direction (incoming or outgoing) */
  readonly direction: TransactionDirection | string | undefined;
  /** Currency symbol or ID */
  readonly currency: string | undefined;
  /** Currency information */
  readonly currencyInfo: TransactionCurrencyInfo | undefined;
  /** Blockchain identifier (e.g., "ETH", "BTC") */
  readonly blockchain: string | undefined;
  /** Network identifier (e.g., "mainnet", "testnet") */
  readonly network: string | undefined;
  /** Blockchain transaction hash */
  readonly hash: string | undefined;
  /** Block number */
  readonly block: string | undefined;
  /** Confirmation block number */
  readonly confirmationBlock: string | undefined;
  /** Amount in smallest currency unit */
  readonly amount: string | undefined;
  /** Amount in main currency unit */
  readonly amountMainUnit: string | undefined;
  /** Fee in smallest currency unit */
  readonly fee: string | undefined;
  /** Fee in main currency unit */
  readonly feeMainUnit: string | undefined;
  /** Transaction type (e.g., "burn", "approve") */
  readonly type: string | undefined;
  /** Current transaction status */
  readonly status: TransactionStatus | string | undefined;
  /** Whether the transaction is fully confirmed */
  readonly isConfirmed: boolean;
  /** Source addresses */
  readonly sources: TransactionAddressInfo[];
  /** Destination addresses */
  readonly destinations: TransactionAddressInfo[];
  /** Transaction ID (internal) */
  readonly transactionId: string | undefined;
  /** Unique ID */
  readonly uniqueId: string | undefined;
  /** Associated request ID */
  readonly requestId: string | undefined;
  /** Whether the associated request is visible */
  readonly requestVisible: boolean | undefined;
  /** When the transaction was received by the blockchain */
  readonly receptionDate: Date | undefined;
  /** When the transaction was confirmed */
  readonly confirmationDate: Date | undefined;
  /** Transaction attributes */
  readonly attributes: TransactionAttribute[];
  /** Fork number (typically 0, 1 for new platforms like DOT AssetHub) */
  readonly forkNumber: string | undefined;
}

/**
 * Options for listing transactions.
 */
export interface ListTransactionsOptions {
  /** Maximum number of transactions to return (default: 50, max: 100) */
  readonly limit?: number;
  /** Offset for pagination (default: 0) */
  readonly offset?: number;
  /** Filter by currency ID or symbol */
  readonly currency?: string;
  /** Filter by direction ("incoming" or "outgoing") */
  readonly direction?: TransactionDirection | string;
  /** Filter by blockchain */
  readonly blockchain?: string;
  /** Filter by network */
  readonly network?: string;
  /** Filter transactions after this date */
  readonly fromDate?: Date;
  /** Filter transactions before this date */
  readonly toDate?: Date;
  /** Search query string */
  readonly query?: string;
  /** Filter by transaction type */
  readonly type?: string;
  /** Filter by source address */
  readonly source?: string;
  /** Filter by destination address */
  readonly destination?: string;
  /** Filter by address (either source or destination) */
  readonly address?: string;
  /** Filter by transaction hashes */
  readonly hashes?: string[];
  /** Filter by transaction IDs */
  readonly ids?: string[];
  /** Filter by from block number */
  readonly fromBlockNumber?: string;
  /** Filter by to block number */
  readonly toBlockNumber?: string;
  /** Filter by amount above this value */
  readonly amountAbove?: string;
  /** Exclude transactions with unknown source/destination */
  readonly excludeUnknownSourceDestination?: boolean;
  /** Filter by customer ID */
  readonly customerId?: string;
}
