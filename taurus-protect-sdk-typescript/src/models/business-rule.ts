/**
 * Business rule models for Taurus-PROTECT SDK.
 *
 * Business rules define automated policies that apply to wallets, addresses,
 * or currencies. They can enforce constraints like spending limits, approval
 * requirements, or allowed transaction types.
 */

/**
 * Currency information associated with a business rule.
 */
export interface BusinessRuleCurrency {
  /** Currency ID */
  readonly id?: string;
  /** Currency symbol (e.g., ETH, BTC) */
  readonly symbol?: string;
  /** Currency name */
  readonly name?: string;
}

/**
 * Business rule model representing automated policies in Taurus-PROTECT.
 *
 * Business rules can enforce constraints like spending limits, approval requirements,
 * or allowed transaction types at various levels (global, currency, wallet, address).
 */
export interface BusinessRule {
  /** Unique identifier for the rule */
  readonly id?: string;
  /** Tenant ID the rule belongs to */
  readonly tenantId?: number;
  /** Currency ID (deprecated, use entityType/entityId) */
  readonly currency?: string;
  /** Wallet ID (deprecated, use entityType/entityId) */
  readonly walletId?: string;
  /** Address ID (deprecated, use entityType/entityId) */
  readonly addressId?: string;
  /** Rule key identifier */
  readonly ruleKey?: string;
  /** Rule value */
  readonly ruleValue?: string;
  /** Rule group for categorization */
  readonly ruleGroup?: string;
  /** Human-readable description of the rule */
  readonly ruleDescription?: string;
  /** Validation expression for the rule */
  readonly ruleValidation?: string;
  /** Currency information */
  readonly currencyInfo?: BusinessRuleCurrency;
  /** Entity type: global, currency, wallet, address, exchange, exchange_account, tn_participant */
  readonly entityType?: string;
  /** Entity ID (wallet ID, address ID, currency ID, exchange label, etc.) */
  readonly entityId?: string;
}

/**
 * Options for listing business rules.
 */
export interface ListBusinessRulesOptions {
  /** Filter by rule IDs */
  ids?: string[];
  /** Filter by rule keys */
  ruleKeys?: string[];
  /** Filter by rule groups */
  ruleGroups?: string[];
  /** Filter by wallet IDs (deprecated, use entityType/entityIds) */
  walletIds?: string[];
  /** Filter by currency IDs */
  currencyIds?: string[];
  /** Filter by address IDs (deprecated, use entityType/entityIds) */
  addressIds?: string[];
  /** Filter by level: global, currency, wallet, address (deprecated, use entityType) */
  level?: string;
  /** Filter by entity type: global, currency, wallet, address, exchange, exchange_account, tn_participant */
  entityType?: string;
  /** Filter by entity IDs */
  entityIds?: string[];
  /** Page size for pagination */
  pageSize?: number;
  /** Current page cursor for pagination */
  currentPage?: string;
  /** Page request direction: FIRST, PREVIOUS, NEXT, LAST */
  pageRequest?: 'FIRST' | 'PREVIOUS' | 'NEXT' | 'LAST';
}

/**
 * Result of listing business rules with pagination cursor.
 */
export interface ListBusinessRulesResult {
  /** List of business rules */
  readonly rules: BusinessRule[];
  /** Cursor for next page (base64-encoded) */
  readonly nextCursor?: string;
  /** Whether there are more pages */
  readonly hasMore: boolean;
}
