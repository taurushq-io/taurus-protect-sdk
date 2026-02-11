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
  BusinessRuleCurrency,
  ListBusinessRulesOptions,
  ListBusinessRulesResult,
} from '../models/business-rule';
import { BaseService } from './base';

// Re-export types for convenience
export type { BusinessRule, BusinessRuleCurrency, ListBusinessRulesOptions, ListBusinessRulesResult };

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
   * @param options - Optional filtering and pagination options
   * @returns A result containing business rules and pagination cursor
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
   * while (result.hasMore) {
   *   result = await businessRuleService.list({
   *     pageSize: 50,
   *     currentPage: result.nextCursor,
   *     pageRequest: 'NEXT',
   *   });
   * }
   * ```
   */
  async list(options?: ListBusinessRulesOptions): Promise<ListBusinessRulesResult> {
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
        cursorPageSize: options?.pageSize != null ? String(options.pageSize) : undefined,
        cursorCurrentPage: options?.currentPage,
        cursorPageRequest: options?.pageRequest,
      });

      const resp = response as Record<string, unknown>;
      const rules = businessRulesFromDto(resp.result as unknown[]);

      // Extract cursor information
      const cursor = resp.cursor as Record<string, unknown> | undefined;
      const nextCursor = cursor?.currentPage != null ? String(cursor.currentPage) : undefined;
      const hasMore = cursor?.hasNextPage === true || cursor?.hasMore === true;

      return {
        rules,
        nextCursor,
        hasMore,
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
      });

      const resp = response as Record<string, unknown>;
      const rules = businessRulesFromDto(resp.result as unknown[]);

      if (rules.length === 0) {
        throw new NotFoundError(`Business rule with id '${ruleId}' not found`);
      }

      return rules[0];
    });
  }
}
