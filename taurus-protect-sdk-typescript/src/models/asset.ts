/**
 * Asset models for Taurus-PROTECT SDK.
 *
 * Covers the asset balance views (which addresses / wallets hold an asset) and the v2
 * asset registry (assets, their holders and their lifecycle operations).
 */

import type { Address } from './address';
import type { CursorPage, CursorPageOptions } from './pagination';
import type { Wallet } from './wallet';

/**
 * Options for listing the addresses that hold an asset. A cursor list: `pageSize` 1-100
 * (default 20) and the `cursor` of a previous page.
 */
export interface GetAssetAddressesOptions extends CursorPageOptions {
  /** The currency code (e.g., "ETH", "BTC", "USDC") - required */
  readonly currency: string;
  /** Optional wallet ID to filter addresses */
  readonly walletId?: string;
  /** Optional address ID to filter */
  readonly addressId?: string;
  /** Optional blockchain addresses to filter */
  readonly addresses?: string[];
}

/**
 * Options for listing the wallets that hold an asset. A cursor list: `pageSize` 1-100
 * (default 20) and the `cursor` of a previous page.
 */
export interface GetAssetWalletsOptions extends CursorPageOptions {
  /** The currency code (e.g., "ETH", "BTC", "USDC") - required */
  readonly currency: string;
  /** Optional wallet ID to filter */
  readonly walletId?: string;
  /** Optional wallet name to filter */
  readonly walletName?: string;
}

/**
 * A page of the addresses holding an asset.
 */
export interface ListAssetAddressesResult {
  /** The verified addresses of this page */
  readonly items: Address[];
  /** Cursor pagination, including the server's total */
  readonly pagination: CursorPage;
}

/**
 * A page of the wallets holding an asset.
 */
export interface ListAssetWalletsResult {
  /** The wallets of this page */
  readonly items: Wallet[];
  /** Cursor pagination, including the server's total */
  readonly pagination: CursorPage;
}

/** Address type filter of the v2 asset-holder query. */
export type AssetAddressTypeV2 =
  | "ADDRESS_TYPE_V2_INTERNAL"
  | "ADDRESS_TYPE_V2_WHITELISTED"
  | "ADDRESS_TYPE_V2_EXTERNAL";

/** KYC status filter of the v2 asset-holder query. */
export type AssetKycStatusV2 = "KYC_STATUS_V2_APPROVED" | "KYC_STATUS_V2_REVOKED";

/** Operation type filter of the v2 asset operations list. */
export type AssetOperationTypeV2 =
  | "ASSET_OPERATION_TYPE_V2_CREATE"
  | "ASSET_OPERATION_TYPE_V2_UPDATE"
  | "ASSET_OPERATION_TYPE_V2_IMPORT"
  | "ASSET_OPERATION_TYPE_V2_MINT"
  | "ASSET_OPERATION_TYPE_V2_BURN"
  | "ASSET_OPERATION_TYPE_V2_PAUSE"
  | "ASSET_OPERATION_TYPE_V2_UNPAUSE"
  | "ASSET_OPERATION_TYPE_V2_PAUSE_ACCOUNT"
  | "ASSET_OPERATION_TYPE_V2_UNPAUSE_ACCOUNT"
  | "ASSET_OPERATION_TYPE_V2_SET_KYC";

/** Operation status filter of the v2 asset operations list. */
export type AssetOperationStatusV2 =
  | "ASSET_OPERATION_STATUS_V2_PENDING"
  | "ASSET_OPERATION_STATUS_V2_PAUSED"
  | "ASSET_OPERATION_STATUS_V2_COMPLETED"
  | "ASSET_OPERATION_STATUS_V2_FAILED"
  | "ASSET_OPERATION_STATUS_V2_CANCELED";

/**
 * Options for querying the v2 asset registry. A cursor list: `pageSize` 1-100
 * (default 20) and the `cursor` of a previous page.
 */
export interface QueryAssetsOptions extends CursorPageOptions {
  /** Filter by blockchain */
  readonly blockchain?: string;
  /** Filter by network */
  readonly network?: string;
  /** Filter by symbol */
  readonly symbol?: string;
  /** Filter by contract address */
  readonly contractAddress?: string;
  /** Filter by label */
  readonly label?: string;
  /** Filter by currency name */
  readonly currencyName?: string;
}

/**
 * Options for querying the holders of a v2 asset. A cursor list: `pageSize` 1-100
 * (default 20) and the `cursor` of a previous page.
 */
export interface QueryAssetAddressesOptions extends CursorPageOptions {
  /** Filter by address type */
  readonly addressType?: AssetAddressTypeV2;
  /** Filter by KYC status */
  readonly kycStatus?: AssetKycStatusV2;
}

/**
 * Options for listing the operations of a v2 asset. A cursor list: `pageSize` 1-100
 * (default 20) and the `cursor` of a previous page.
 */
export interface ListAssetOperationsOptions extends CursorPageOptions {
  /** Filter by operation type */
  readonly type?: AssetOperationTypeV2;
  /** Filter by operation status */
  readonly status?: AssetOperationStatusV2;
}

/** A key/value attribute of a v2 asset. */
export interface AssetAttributeV2 {
  readonly key?: string;
  readonly value?: string;
}

/** On-ledger configuration of a Canton instrument. */
export interface CantonInstrumentConfigurationV2 {
  /** Contract ID of the instrument configuration */
  readonly cid?: string;
  /** Whether holders need credentials */
  readonly requireCredentials?: boolean;
  /** Whether the instrument is paused */
  readonly paused?: boolean;
  /** Operator party */
  readonly operator?: string;
}

/**
 * An asset of the v2 asset registry.
 */
export interface AssetV2 {
  /** Asset ID */
  readonly id?: string;
  /** Tenant ID */
  readonly tenantId?: string;
  /** Record version */
  readonly version?: string;
  /** Creation date */
  readonly createdAt?: Date;
  /** Last update date */
  readonly updatedAt?: Date;
  /** Label */
  readonly label?: string;
  /** Asset type */
  readonly assetType?: string;
  /** Registry status (ASSET_VIEW_STATUS_V2_*) */
  readonly status?: string;
  /** Blockchain */
  readonly blockchain?: string;
  /** Network */
  readonly network?: string;
  /** Currency ID of the asset */
  readonly currencyId?: string;
  /** Name */
  readonly name?: string;
  /** Symbol */
  readonly symbol?: string;
  /** Decimals */
  readonly decimals?: string;
  /** Contract address */
  readonly contractAddress?: string;
  /** Attributes */
  readonly attributes: AssetAttributeV2[];
  /** Canton instrument ID, for a Canton native token */
  readonly cantonInstrumentId?: string;
  /** Canton instrument configuration, for a Canton native token */
  readonly cantonConfiguration?: CantonInstrumentConfigurationV2;
}

/**
 * A holder of a v2 asset.
 *
 * The holder rows carry no signature of their own. An internal or whitelisted holder is
 * returned only once its address was re-read through the verified reader for its type
 * (`verified: true`, and `address` is the verified reader's). Any other holder, such as
 * an external on-chain one, is returned with `verified: false`: on-chain data, never a
 * Taurus-PROTECT address.
 */
export interface AssetAddressV2 {
  /** Blockchain address: the verified reader's when `verified` is true */
  readonly address?: string;
  /** KYC status (KYC_STATUS_V2_*) */
  readonly kycStatus?: string;
  /** Balance of the asset at this address */
  readonly balance?: string;
  /** Address type (ADDRESS_TYPE_V2_*) */
  readonly addressType?: string;
  /** Internal address ID, for an internal address */
  readonly addressId?: string;
  /** Whitelisted address ID, for a whitelisted address */
  readonly whitelistedAddressId?: string;
  /**
   * True when `address` was verified: an internal holder's against its HSM-signed
   * managed address, a whitelisted holder's against its verified whitelist entry.
   */
  readonly verified: boolean;
}

/**
 * A holder row left out of a page because it could not be verified. The whitelist
 * `{id, reason}` shape: `id` is the row's addressID / whitelistedAddressID, else its
 * address.
 */
export interface ExcludedAssetAddressV2 {
  /** The row's address ID, whitelisted address ID, or address */
  readonly id: string;
  /** Why the row could not be verified */
  readonly reason: string;
}

/** The address an asset operation acts on. */
export interface AssetAddressTargetV2 {
  /** Internal address ID */
  readonly addressId?: string;
  /** Whitelisted address ID */
  readonly whitelistedAddressId?: string;
}

/**
 * A lifecycle operation of a v2 asset. Exactly one of the detail fields is set, matching
 * `type`.
 */
export interface AssetOperationV2 {
  /** Operation ID */
  readonly id?: string;
  /** The asset the operation acts on */
  readonly assetId?: string;
  /** Operation type (ASSET_OPERATION_TYPE_V2_*) */
  readonly type?: string;
  /** Operation status (ASSET_OPERATION_STATUS_V2_*) */
  readonly status?: string;
  /** Creation date */
  readonly createdAt?: Date;
  /** Last update date */
  readonly updatedAt?: Date;
  /** Address that initiated the operation */
  readonly initiatedByAddressId?: string;
  /** Why the operation failed, when it did */
  readonly failureReason?: string;
  /** What the operation waits for, when it is blocked */
  readonly blockingReason?: string;
  /** Details of a create operation */
  readonly create?: {
    readonly label?: string;
    readonly price?: string;
    readonly decimals?: string;
    readonly blockchain?: string;
    readonly network?: string;
    readonly assetType?: string;
    readonly cantonInstrumentId?: string;
    readonly cantonName?: string;
    readonly cantonSymbol?: string;
  };
  /** Details of an update operation */
  readonly update?: {
    readonly label?: string;
    readonly price?: string;
    readonly cantonRequireCredentials?: boolean;
  };
  /** Details of an import operation */
  readonly import?: {
    readonly blockchain?: string;
    readonly network?: string;
    readonly label?: string;
    readonly price?: string;
    readonly decimals?: string;
    readonly address?: string;
  };
  /** Details of a mint operation */
  readonly mint?: {
    readonly destination?: AssetAddressTargetV2;
    readonly amount?: string;
    readonly nftMetadata?: string[];
  };
  /** Details of a burn operation */
  readonly burn?: {
    readonly destination?: AssetAddressTargetV2;
    readonly amount?: string;
    readonly nftTokenIds?: string[];
  };
  /** Details of a pause-account operation */
  readonly pauseAccount?: { readonly target?: AssetAddressTargetV2 };
  /** Details of an unpause-account operation */
  readonly unpauseAccount?: { readonly target?: AssetAddressTargetV2 };
  /** Details of a set-KYC operation */
  readonly setKyc?: {
    readonly target?: AssetAddressTargetV2;
    /** KYC status (KYC_STATUS_V2_*) */
    readonly status?: string;
  };
}

/** A page of v2 assets. */
export interface ListAssetsV2Result {
  /** The assets of this page */
  readonly items: AssetV2[];
  /** Cursor pagination: continue with `cursor: pagination.nextCursor` while `hasMore` */
  readonly pagination: CursorPage;
}

/** A page of the holders of a v2 asset. */
export interface ListAssetAddressesV2Result {
  /** The holders of this page that survived verification (see {@link AssetAddressV2}) */
  readonly items: AssetAddressV2[];
  /**
   * Cursor pagination: continue with `cursor: pagination.nextCursor` while `hasMore`.
   * Exclusions never move it.
   */
  readonly pagination: CursorPage;
  /** The rows the server returned that failed verification and are absent from `items` */
  readonly excludedUnverified: ExcludedAssetAddressV2[];
}

/** A page of the operations of a v2 asset. */
export interface ListAssetOperationsResult {
  /** The operations of this page */
  readonly items: AssetOperationV2[];
  /** Cursor pagination: continue with `cursor: pagination.nextCursor` while `hasMore` */
  readonly pagination: CursorPage;
}
