/**
 * Wallet domain models for Taurus-PROTECT SDK.
 */

import type { Currency } from './currency';

/**
 * Wallet status enum.
 */
export enum WalletStatus {
  UNKNOWN = 'UNKNOWN',
  ACTIVE = 'ACTIVE',
  PENDING = 'PENDING',
  DISABLED = 'DISABLED',
}

/**
 * Wallet type enum.
 */
export enum WalletType {
  UNKNOWN = 'UNKNOWN',
  STANDARD = 'STANDARD',
  AIRGAP = 'AIRGAP',
  STAKING = 'STAKING',
}

/**
 * Represents wallet balance information.
 */
export interface WalletBalance {
  /** Total confirmed balance in the smallest currency unit (e.g., WEI for ETH) */
  readonly totalConfirmed: string | undefined;
  /** Total balance including unconfirmed transactions */
  readonly totalUnconfirmed: string | undefined;
  /** Available confirmed balance ready to be spent */
  readonly availableConfirmed: string | undefined;
  /** Available balance including unconfirmed transactions */
  readonly availableUnconfirmed: string | undefined;
  /** Confirmed reserved balance held for pending transactions */
  readonly reservedConfirmed: string | undefined;
  /** Reserved unconfirmed balance not yet validated */
  readonly reservedUnconfirmed: string | undefined;
}

/**
 * Represents a wallet attribute (key-value metadata).
 */
export interface WalletAttribute {
  /** Unique attribute identifier */
  readonly id: string | undefined;
  /** Attribute key */
  readonly key: string | undefined;
  /** Attribute value */
  readonly value: string | undefined;
  /** Content type of the attribute */
  readonly contentType: string | undefined;
  /** Owner of the attribute */
  readonly owner: string | undefined;
  /** Attribute type */
  readonly type: string | undefined;
  /** Attribute subtype */
  readonly subtype: string | undefined;
  /** Whether this attribute is a file */
  readonly isFile: boolean | undefined;
}

/**
 * Represents a cryptocurrency wallet in Taurus-PROTECT.
 */
export interface Wallet {
  /** Unique identifier */
  readonly id: string;
  /** Human-readable name */
  readonly name: string;
  /** Current status */
  readonly status: WalletStatus;
  /** Wallet type */
  readonly type: WalletType;
  /** Blockchain identifier (e.g., "ETH", "BTC") */
  readonly blockchain: string | undefined;
  /** Network identifier (e.g., "mainnet", "testnet") */
  readonly network: string | undefined;
  /** Currency identifier or symbol */
  readonly currency: string | undefined;
  /** Whether this wallet is disabled */
  readonly disabled: boolean;
  /** Whether this is an omnibus wallet */
  readonly isOmnibus: boolean;
  /** Wallet balance information */
  readonly balance: WalletBalance | undefined;
  /** Creation timestamp */
  readonly createdAt: Date | undefined;
  /** Last update timestamp */
  readonly updatedAt: Date | undefined;
  /** Optional description/comment */
  readonly comment: string | undefined;
  /** Customer identifier for omnibus wallets */
  readonly customerId: string | undefined;
  /** Number of addresses in the wallet */
  readonly addressesCount: number | undefined;
  /** Wallet attributes (key-value metadata) */
  readonly attributes: WalletAttribute[];
  /** Visibility group ID */
  readonly visibilityGroupId: string | undefined;
  /** External wallet identifier */
  readonly externalWalletId: string | undefined;
  /** HD wallet account derivation path */
  readonly accountPath: string | undefined;
  /** Detailed currency information */
  readonly currencyInfo: Currency | undefined;
  /** Associated tags (for backward compatibility) */
  readonly tags: string[];
}

/**
 * Request parameters for creating a new wallet.
 */
export interface CreateWalletRequest {
  /**
   * Wallet name (must be unique per blockchain/network).
   * Required.
   */
  readonly name: string;
  /**
   * Blockchain identifier (e.g., "ETH", "BTC").
   * Required if currency is not provided.
   */
  readonly blockchain?: string;
  /**
   * Network identifier (e.g., "mainnet", "testnet").
   * Required if blockchain is enabled on multiple networks.
   */
  readonly network?: string;
  /**
   * Currency ID or symbol. Required if blockchain is not provided.
   * Must be the native currency of a blockchain.
   */
  readonly currency?: string;
  /**
   * Whether this is an omnibus wallet (single owner).
   * When true, addresses can be used as transaction sources.
   */
  readonly isOmnibus?: boolean;
  /**
   * Optional comment or description (max 254 characters).
   */
  readonly comment?: string;
  /**
   * Customer identifier (max 254 characters).
   * Applied to all addresses when isOmnibus is true.
   */
  readonly customerId?: string;
  /**
   * Visibility group UUID.
   * Restricts wallet visibility and management to users in this group.
   */
  readonly visibilityGroupId?: string;
  /**
   * External wallet identifier.
   */
  readonly externalWalletId?: string;
}

/**
 * Options for listing wallets.
 */
export interface ListWalletsOptions {
  /** Maximum number of wallets to return (default: 50) */
  readonly limit?: number;
  /** Number of wallets to skip for pagination (default: 0) */
  readonly offset?: number;
  /** Filter by currency ID or symbol */
  readonly currency?: string;
  /** Filter by currencies (array of IDs or symbols) */
  readonly currencies?: string[];
  /** Search query string */
  readonly query?: string;
  /** Filter by wallet name (case-insensitive partial match) */
  readonly name?: string;
  /** Sort order: "ASC" or "DESC" */
  readonly sortOrder?: string;
  /** Exclude wallets disabled by business rules */
  readonly excludeDisabled?: boolean;
  /** Filter by tag IDs (OR combination) */
  readonly tagIds?: string[];
  /** Only include wallets with positive balance */
  readonly onlyPositiveBalance?: boolean;
  /** Filter by blockchain */
  readonly blockchain?: string;
  /** Filter by network */
  readonly network?: string;
  /** Filter by specific wallet IDs (max 100) */
  readonly ids?: string[];
}

/**
 * Options for getting wallet balance history.
 */
export interface GetBalanceHistoryOptions {
  /** Interval in hours for balance snapshots */
  readonly intervalHours: number;
}

/**
 * A point in the wallet balance history.
 */
export interface BalanceHistoryPoint {
  /** Timestamp of this balance snapshot */
  readonly pointDate: Date | undefined;
  /** Balance at this point in time */
  readonly balance: WalletBalance | undefined;
}
