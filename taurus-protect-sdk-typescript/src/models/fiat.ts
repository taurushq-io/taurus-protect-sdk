/**
 * Fiat provider models for Taurus-PROTECT SDK.
 */

import type { Currency } from './currency';

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
 * Response cursor for paginated fiat provider results.
 */
export interface FiatResponseCursor {
  /** Current page token */
  readonly currentPage?: string;
  /** Next page token */
  readonly nextPage?: string;
  /** Whether there are more pages */
  readonly hasMore: boolean;
}

/**
 * Paginated result for fiat provider accounts.
 */
export interface FiatProviderAccountResult {
  /** The fiat provider accounts */
  readonly accounts: FiatProviderAccount[];
  /** Pagination cursor */
  readonly cursor?: FiatResponseCursor;
}

/**
 * Paginated result for fiat provider counterparty accounts.
 */
export interface FiatProviderCounterpartyAccountResult {
  /** The fiat provider counterparty accounts */
  readonly accounts: FiatProviderCounterpartyAccount[];
  /** Pagination cursor */
  readonly cursor?: FiatResponseCursor;
}

/**
 * Paginated result for fiat provider operations.
 */
export interface FiatProviderOperationResult {
  /** The fiat provider operations */
  readonly operations: FiatProviderOperation[];
  /** Pagination cursor */
  readonly cursor?: FiatResponseCursor;
}

/**
 * Options for listing fiat provider accounts.
 */
export interface ListFiatProviderAccountsOptions {
  /** Filter by provider */
  provider: string;
  /** Filter by label */
  label: string;
  /** Filter by account type (e.g., 'wallet', 'bank') */
  accountType?: string;
  /** Sort order for results ('ASC' or 'DESC') */
  sortOrder?: string;
  /** Pagination cursor */
  cursor?: FiatRequestCursor;
}

/**
 * Options for listing fiat provider counterparty accounts.
 */
export interface ListFiatProviderCounterpartyAccountsOptions {
  /** Filter by provider */
  provider: string;
  /** Filter by label */
  label: string;
  /** Filter by counterparty ID */
  counterpartyId?: string;
  /** Sort order for results ('ASC' or 'DESC') */
  sortOrder?: string;
  /** Pagination cursor */
  cursor?: FiatRequestCursor;
}

/**
 * Options for listing fiat provider operations.
 */
export interface ListFiatProviderOperationsOptions {
  /** Filter by provider */
  provider?: string;
  /** Filter by label */
  label?: string;
  /** Sort order for results ('ASC' or 'DESC') */
  sortOrder?: string;
  /** Pagination cursor */
  cursor?: FiatRequestCursor;
}

/**
 * Request cursor for paginating fiat provider results.
 */
export interface FiatRequestCursor {
  /** Current page token */
  currentPage?: string;
  /** Page request type ('FIRST', 'PREVIOUS', 'NEXT', 'LAST') */
  pageRequest?: 'FIRST' | 'PREVIOUS' | 'NEXT' | 'LAST';
  /** Number of items per page */
  pageSize?: number;
}
