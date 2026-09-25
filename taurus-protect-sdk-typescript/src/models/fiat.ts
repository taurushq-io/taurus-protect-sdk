/**
 * Fiat provider models for Taurus-PROTECT SDK.
 */

import type { Currency } from './currency';
import type { CursorNavigationOptions, CursorPage, CursorPageOptions } from './pagination';

/**
 * Represents a fiat currency provider in the Taurus-PROTECT system.
 */
export interface FiatProvider {
  /** The fiat provider identifier (e.g., 'circle') */
  readonly provider?: string;
  /** The label of the fiat provider set in the config */
  readonly label?: string;
  /** Valuation in the base currency main unit (CHF, EUR, USD, etc.) */
  readonly baseCurrencyValuation?: string;
}

/**
 * Represents a fiat provider account in the Taurus-PROTECT system.
 */
export interface FiatProviderAccount {
  /** Unique account identifier */
  readonly id?: string;
  /** The fiat provider identifier */
  readonly provider?: string;
  /** The label of the fiat provider set in the config */
  readonly label?: string;
  /** The type of account (e.g., 'wallet', 'bank') */
  readonly accountType?: string;
  /** The account identifier within the provider */
  readonly accountIdentifier?: string;
  /** The display name of the account */
  readonly accountName?: string;
  /** Balance in the smallest currency unit */
  readonly totalBalance?: string;
  /** Currency identifier */
  readonly currencyId?: string;
  /** Currency information */
  readonly currencyInfo?: Currency;
  /** Valuation in the base currency main unit */
  readonly baseCurrencyValuation?: string;
  /** Account creation date */
  readonly creationDate?: Date;
  /** Account last update date */
  readonly updateDate?: Date;
}

/**
 * Represents a fiat provider counterparty account in the Taurus-PROTECT system.
 */
export interface FiatProviderCounterpartyAccount {
  /** Unique counterparty account identifier */
  readonly id?: string;
  /** The fiat provider identifier */
  readonly provider?: string;
  /** The label of the fiat provider set in the config */
  readonly label?: string;
  /** The type of account */
  readonly accountType?: string;
  /** The account identifier within the provider */
  readonly accountIdentifier?: string;
  /** The display name of the account */
  readonly accountName?: string;
  /** The counterparty identifier */
  readonly counterpartyId?: string;
  /** The counterparty display name */
  readonly counterpartyName?: string;
  /** Currency identifier */
  readonly currencyId?: string;
  /** Currency information */
  readonly currencyInfo?: Currency;
  /** Account creation date */
  readonly creationDate?: Date;
  /** Account last update date */
  readonly updateDate?: Date;
}

/**
 * Represents a fiat provider operation in the Taurus-PROTECT system.
 */
export interface FiatProviderOperation {
  /** Unique operation identifier */
  readonly id?: string;
  /** The fiat provider identifier */
  readonly provider?: string;
  /** The label of the fiat provider set in the config */
  readonly label?: string;
  /** The type of operation */
  readonly operationType?: string;
  /** The operation identifier within the provider */
  readonly operationIdentifier?: string;
  /** The direction of the operation (e.g., 'incoming', 'outgoing') */
  readonly operationDirection?: string;
  /** The status of the operation (e.g., 'pending', 'completed', 'failed') */
  readonly status?: string;
  /** The operation amount */
  readonly amount?: string;
  /** Currency identifier */
  readonly currencyId?: string;
  /** Currency information */
  readonly currencyInfo?: Currency;
  /** The source account identifier */
  readonly fromAccountId?: string;
  /** The destination account identifier */
  readonly toAccountId?: string;
  /** Details about the source of the operation */
  readonly fromDetails?: string;
  /** Details about the destination of the operation */
  readonly toDetails?: string;
  /** A comment or description of the operation */
  readonly comment?: string;
  /** Additional operation details */
  readonly operationDetails?: string;
  /** Operation creation date */
  readonly creationDate?: Date;
  /** Operation last update date */
  readonly updateDate?: Date;
}

/**
 * A fiat provider entity (a legal entity a provider holds accounts for).
 */
export interface FiatProviderEntity {
  /** Unique identifier */
  readonly id?: string;
  /** Fiat provider name */
  readonly provider?: string;
  /** Provider label */
  readonly label?: string;
  /** Account identifier at the provider */
  readonly accountIdentifier?: string;
  /** Entity name */
  readonly name?: string;
  /** Entity details, as the provider reports them */
  readonly details?: string;
  /** Creation date */
  readonly creationDate?: Date;
  /** Last update date */
  readonly updateDate?: Date;
}

/**
 * A page of fiat provider accounts.
 */
export interface FiatProviderAccountResult {
  /** The fiat provider accounts */
  readonly accounts: FiatProviderAccount[];
  /** Cursor pagination: continue with `cursor: pagination.nextCursor` while `hasMore` */
  readonly pagination: CursorPage;
}

/**
 * A page of fiat provider counterparty accounts.
 */
export interface FiatProviderCounterpartyAccountResult {
  /** The fiat provider counterparty accounts */
  readonly accounts: FiatProviderCounterpartyAccount[];
  /** Cursor pagination: continue with `cursor: pagination.nextCursor` while `hasMore` */
  readonly pagination: CursorPage;
}

/**
 * A page of fiat provider operations.
 */
export interface FiatProviderOperationResult {
  /** The fiat provider operations */
  readonly operations: FiatProviderOperation[];
  /** Cursor pagination: continue with `cursor: pagination.nextCursor` while `hasMore` */
  readonly pagination: CursorPage;
}

/**
 * A page of fiat provider entities.
 */
export interface ListFiatProviderEntitiesResult {
  /** The fiat provider entities */
  readonly items: FiatProviderEntity[];
  /** Cursor pagination: continue with `cursor: pagination.nextCursor` while `hasMore` */
  readonly pagination: CursorPage;
}

/**
 * Options for listing fiat provider accounts. A cursor list: `pageSize` 1-100
 * (default 20) and the `cursor` of a previous page.
 */
export interface ListFiatProviderAccountsOptions extends CursorNavigationOptions {
  /** Fiat provider (required) */
  readonly provider: string;
  /** Provider label (required) */
  readonly label: string;
  /** Filter by account type (e.g., 'wallet', 'bank') */
  readonly accountType?: string;
  /** Sort order for results ('ASC' or 'DESC') */
  readonly sortOrder?: string;
}

/**
 * Options for listing fiat provider counterparty accounts. A cursor list: `pageSize`
 * 1-100 (default 20) and the `cursor` of a previous page.
 */
export interface ListFiatProviderCounterpartyAccountsOptions extends CursorNavigationOptions {
  /** Fiat provider (required) */
  readonly provider: string;
  /** Provider label (required) */
  readonly label: string;
  /** Filter by counterparty ID */
  readonly counterpartyId?: string;
  /** Sort order for results ('ASC' or 'DESC') */
  readonly sortOrder?: string;
}

/**
 * Options for listing fiat provider operations. A cursor list: `pageSize` 1-100
 * (default 20) and the `cursor` of a previous page.
 */
export interface ListFiatProviderOperationsOptions extends CursorNavigationOptions {
  /** Filter by provider */
  readonly provider?: string;
  /** Filter by label */
  readonly label?: string;
  /** Sort order for results ('ASC' or 'DESC') */
  readonly sortOrder?: string;
}

/**
 * Options for listing fiat provider entities. A cursor list: `pageSize` 1-100
 * (default 20) and the `cursor` of a previous page.
 */
export interface ListFiatProviderEntitiesOptions extends CursorPageOptions {
  /** Filter by provider */
  readonly provider?: string;
  /** Filter by label */
  readonly label?: string;
  /** Sort order for results ('ASC' or 'DESC') */
  readonly sortOrder?: string;
}
