/**
 * Business rule service for Taurus-PROTECT SDK.
 *
 * Provides methods for retrieving business rules that define automated policies
 * applied to wallets, addresses, or currencies.
 */

import { NotFoundError, ValidationError } from '../errors';
import type { BusinessRulesApi } from '../internal/openapi/apis/BusinessRulesApi';
import { businessRulesFromDto } from '../mappers/business-rule';
import type {
  BusinessRule,
  ListBusinessRulesOptions,
  ListBusinessRulesResult,
} from '../models/business-rule';
import { buildCursorPage, cursorRequest } from '../models/pagination';
import { BaseService } from './base';
import { cursorQuery } from './paging';

/**
 * Service for managing business rules in Taurus-PROTECT.
 *
 * Business rules define automated policies that apply to wallets, addresses, or currencies.
 * They can enforce constraints like spending limits, approval requirements, or allowed
 * transaction types.
 *
 * @example
 * ```typescript
 * // List all business rules
 * const result = await businessRuleService.list();
 * for (const rule of result.rules) {
 *   console.log(`${rule.ruleKey}: ${rule.ruleValue}`);
 * }
 *
 * // Get business rules for a specific wallet
 * const walletRules = await businessRuleService.list({
 *   walletIds: ['123'],
 * });
 *
 * // Get business rules for a specific currency
 * const currencyRules = await businessRuleService.list({
 *   currencyIds: ['ETH'],
 * });
 *
 * // Get a specific rule by ID
 * const rule = await businessRuleService.get('rule-123');
 * console.log(`Rule: ${rule.ruleDescription}`);
 * ```
 */
export class BusinessRuleService extends BaseService {
  private readonly businessRulesApi: BusinessRulesApi;

  /**
   * Creates a new BusinessRuleService instance.
   *
   * @param businessRulesApi - The BusinessRulesApi instance from the OpenAPI client
   */
  constructor(businessRulesApi: BusinessRulesApi) {
    super();
    this.businessRulesApi = businessRulesApi;
  }

  /**
   * Lists business rules with optional filtering and pagination.
   *
   * @param options - Filters, `pageSize` (1-100, default 20) and `cursor`
   * @returns The page of business rules and its cursor pagination
   * @throws {@link ValidationError} If the page size or cursor options are invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // List all business rules
   * const result = await businessRuleService.list();
   *
   * // Filter by wallet
   * const walletRules = await businessRuleService.list({
   *   walletIds: ['123'],
   * });
   *
   * // Filter by currency
   * const currencyRules = await businessRuleService.list({
   *   currencyIds: ['ETH'],
   * });
   *
   * // Filter by entity type
   * const globalRules = await businessRuleService.list({
   *   entityType: 'global',
   * });
   *
   * // Paginate through results
   * let result = await businessRuleService.list({ pageSize: 50 });
   * while (result.pagination.hasMore) {
   *   result = await businessRuleService.list({
   *     pageSize: 50,
   *     cursor: result.pagination.nextCursor,
   *   });
   * }
   * ```
   */
  async list(options?: ListBusinessRulesOptions): Promise<ListBusinessRulesResult> {
    const page = cursorRequest(options);

    return this.execute(async () => {
      const response = await this.businessRulesApi.ruleServiceGetBusinessRulesV2({
        ids: options?.ids,
        ruleKeys: options?.ruleKeys,
        ruleGroups: options?.ruleGroups,
        walletIds: options?.walletIds,
        currencyIds: options?.currencyIds,
        addressIds: options?.addressIds,
        level: options?.level,
        entityType: options?.entityType,
        entityIDs: options?.entityIds,
        ...cursorQuery(page),
      });

      return {
        rules: businessRulesFromDto(response.result),
        pagination: buildCursorPage(page.pageSize, response.cursor),
      };
    });
  }

  /**
   * Gets a business rule by ID.
   *
   * @param ruleId - The unique rule identifier
   * @returns The business rule
   * @throws {@link ValidationError} If ruleId is empty
   * @throws {@link NotFoundError} If rule is not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const rule = await businessRuleService.get('rule-123');
   * console.log(`Rule: ${rule.ruleKey} = ${rule.ruleValue}`);
   * console.log(`Description: ${rule.ruleDescription}`);
   * ```
   */
  async get(ruleId: string): Promise<BusinessRule> {
    if (!ruleId || ruleId.trim() === '') {
      throw new ValidationError('ruleId is required');
    }

    return this.execute(async () => {
      const response = await this.businessRulesApi.ruleServiceGetBusinessRulesV2({
        ids: [ruleId],
        cursorPageSize: '1',
      });

      const rules = businessRulesFromDto(response.result);

      if (rules.length === 0) {
        throw new NotFoundError(`Business rule with id '${ruleId}' not found`);
      }

      return rules[0];
    });
  }

  /**
   * Enables or disables transaction processing for the tenant (the
   * transactions-enabled business rule).
   *
   * This is the kill switch tg-protect-mcpd drives, so it must exist in every SDK. Peer
   * of Go BusinessRuleService.UpdateTransactionsEnabled, Java updateTransactionsEnabled
   * and Python update_transactions_enabled.
   *
   * @param enabled - true to allow transactions, false to halt them
   * @throws {@link APIError} If the API request fails
   *
   * @example
   * ```typescript
   * await businessRuleService.updateTransactionsEnabled(false);
   * ```
   */
  async updateTransactionsEnabled(enabled: boolean): Promise<void> {
    return this.execute(async () => {
      await this.businessRulesApi.ruleServiceUpdateTransactionsEnabledBusinessRule({
        body: { enabled },
      });
    });
  }
}
