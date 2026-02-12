/**
 * Address models for Taurus-PROTECT SDK.
 */

/**
 * Address status enum.
 */
export enum AddressStatus {
  UNKNOWN = "UNKNOWN",
  CREATED = "CREATED",
  CREATING = "CREATING",
  SIGNED = "SIGNED",
  OBSERVED = "OBSERVED",
  CONFIRMED = "CONFIRMED",
  ACTIVE = "ACTIVE",
  PENDING = "PENDING",
  DISABLED = "DISABLED",
}

/**
 * Custom attribute on an address.
 */
export interface AddressAttribute {
  /** Attribute identifier */
  readonly id: string;
  /** Attribute name */
  readonly key: string;
  /** Attribute value */
  readonly value: string;
}

/** Balance information for an address. */
export interface Balance {
  readonly totalConfirmed: string;
  readonly totalUnconfirmed: string;
  readonly availableConfirmed: string;
  readonly availableUnconfirmed: string;
  readonly reservedConfirmed: string;
  readonly reservedUnconfirmed: string;
}

/**
 * Represents a blockchain address in Taurus-PROTECT.
 *
 * Addresses are created within wallets and are used to send and receive
 * cryptocurrency. Each address has an HSM-generated signature that can
 * be verified to ensure integrity.
 */
export interface Address {
  /** Unique address identifier */
  readonly id: string;
  /** Parent wallet ID */
  readonly walletId: string;
  /** The blockchain address string */
  readonly address: string;
  /** Alternate address format if available */
  readonly alternateAddress: string | undefined;
  /** Human-readable label */
  readonly label: string | undefined;
  /** Optional description */
  readonly comment: string | undefined;
  /** Currency symbol (e.g., 'ETH', 'BTC') */
  readonly currency: string;
  /** Optional customer identifier */
  readonly customerId: string | undefined;
  /** Optional external identifier */
  readonly externalAddressId: string | undefined;
  /** HD derivation path */
  readonly addressPath: string | undefined;
  /** Index used for address generation */
  readonly addressIndex: string | undefined;
  /** Current address nonce */
  readonly nonce: string | undefined;
  /** Creation status (created, creating, signed, observed, confirmed) */
  readonly status: string | undefined;
  /** HSM signature for integrity verification */
  readonly signature: string | undefined;
  /** Whether the address is disabled */
  readonly disabled: boolean;
  /** Whether all funds can be used */
  readonly canUseAllFunds: boolean;
  /** Creation timestamp */
  readonly createdAt: Date | undefined;
  /** Last update timestamp */
  readonly updatedAt: Date | undefined;
  /** Custom key-value attributes */
  readonly attributes: AddressAttribute[];
  /** IDs of linked whitelisted addresses */
  readonly linkedWhitelistedAddressIds: string[];
  /** Balance information */
  readonly balance?: Balance;
}

/**
 * Request to create a new address.
 */
export interface CreateAddressRequest {
  /** Wallet ID (required) */
  readonly walletId: string;
  /** Address label (required) */
  readonly label: string;
  /** Optional description */
  readonly comment?: string;
  /** Optional customer identifier */
  readonly customerId?: string;
  /** Optional external identifier */
  readonly externalAddressId?: string;
  /** Address type (e.g., p2pkh, p2wpkh) */
  readonly addressType?: string;
  /** Use non-hardened derivation */
  readonly nonHardenedDerivation?: boolean;
}

/**
 * Options for listing addresses.
 */
export interface ListAddressesOptions {
  /** Filter by wallet ID */
  readonly walletId?: string;
  /** Maximum items to return (default: 50, max: 1000) */
  readonly limit?: number;
  /** Number of items to skip */
  readonly offset?: number;
  /** Search query */
  readonly query?: string;
  /** Filter by blockchain */
  readonly blockchain?: string;
  /** Filter by network */
  readonly network?: string;
}
