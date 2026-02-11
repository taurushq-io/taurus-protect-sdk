/**
 * Business rule mapper functions for converting OpenAPI DTOs to domain models.
 */

import type { BusinessRule, BusinessRuleCurrency } from '../models/business-rule';
import { safeInt, safeString } from './base';

/**
 * Maps a currency DTO to a BusinessRuleCurrency model.
 */
export function businessRuleCurrencyFromDto(dto: unknown): BusinessRuleCurrency | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }
  const d = dto as Record<string, unknown>;
  return {
    id: safeString(d.id),
    symbol: safeString(d.symbol),
    name: safeString(d.name),
  };
}

/**
 * Maps a business rule DTO to a BusinessRule model.
 *
 * Handles the tenantId conversion from string to number to match the Java SDK.
 * Also normalizes entityID/entityId field name differences.
 */
export function businessRuleFromDto(dto: unknown): BusinessRule | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;

  return {
    id: safeString(d.id),
    tenantId: safeInt(d.tenantId),
    currency: safeString(d.currency),
    walletId: safeString(d.walletId),
    addressId: safeString(d.addressId),
    ruleKey: safeString(d.ruleKey),
    ruleValue: safeString(d.ruleValue),
    ruleGroup: safeString(d.ruleGroup),
    ruleDescription: safeString(d.ruleDescription),
    ruleValidation: safeString(d.ruleValidation),
    currencyInfo: businessRuleCurrencyFromDto(d.currencyInfo),
    entityType: safeString(d.entityType),
    entityId: safeString(d.entityID ?? d.entityId),
  };
}

/**
 * Maps an array of business rule DTOs to BusinessRule models.
 */
export function businessRulesFromDto(dtos: unknown[] | null | undefined): BusinessRule[] {
  if (!dtos) {
    return [];
  }
  const result: BusinessRule[] = [];
  for (const dto of dtos) {
    const mapped = businessRuleFromDto(dto);
    if (mapped !== undefined) {
      result.push(mapped);
    }
  }
  return result;
}
