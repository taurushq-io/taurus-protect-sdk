/**
 * Unit tests for business rule mapper functions.
 */

import {
  businessRuleFromDto,
  businessRuleCurrencyFromDto,
  businessRulesFromDto,
} from '../../../src/mappers/business-rule';

describe('businessRuleCurrencyFromDto', () => {
  it('should map all fields', () => {
    const dto = { id: 'c1', symbol: 'ETH', name: 'Ethereum' };

    const result = businessRuleCurrencyFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.id).toBe('c1');
    expect(result!.symbol).toBe('ETH');
    expect(result!.name).toBe('Ethereum');
  });

  it('should return undefined for null input', () => {
    expect(businessRuleCurrencyFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(businessRuleCurrencyFromDto(undefined)).toBeUndefined();
  });
});

describe('businessRuleFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      id: 'br-1',
      tenantId: 42,
      currency: 'ETH',
      walletId: 'w-1',
      addressId: 'a-1',
      ruleKey: 'maxAmount',
      ruleValue: '1000',
      ruleGroup: 'limits',
      ruleDescription: 'Max transfer amount',
      ruleValidation: 'numeric',
      currencyInfo: { id: 'c1', symbol: 'ETH', name: 'Ethereum' },
      entityType: 'wallet',
      entityId: 'w-1',
    };

    const result = businessRuleFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.id).toBe('br-1');
    expect(result!.tenantId).toBe(42);
    expect(result!.currency).toBe('ETH');
    expect(result!.walletId).toBe('w-1');
    expect(result!.addressId).toBe('a-1');
    expect(result!.ruleKey).toBe('maxAmount');
    expect(result!.ruleValue).toBe('1000');
    expect(result!.ruleGroup).toBe('limits');
    expect(result!.ruleDescription).toBe('Max transfer amount');
    expect(result!.ruleValidation).toBe('numeric');
    expect(result!.currencyInfo).toBeDefined();
    expect(result!.currencyInfo!.symbol).toBe('ETH');
    expect(result!.entityType).toBe('wallet');
    expect(result!.entityId).toBe('w-1');
  });

  it('should handle entityID alias (uppercase ID)', () => {
    const dto = { id: 'br-2', entityID: 'e-123' };

    const result = businessRuleFromDto(dto);

    expect(result!.entityId).toBe('e-123');
  });

  it('should prefer entityID over entityId', () => {
    const dto = { id: 'br-3', entityID: 'from-ID', entityId: 'from-Id' };

    const result = businessRuleFromDto(dto);

    expect(result!.entityId).toBe('from-ID');
  });

  it('should return undefined for null input', () => {
    expect(businessRuleFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(businessRuleFromDto(undefined)).toBeUndefined();
  });

  it('should handle empty object', () => {
    const result = businessRuleFromDto({});
    expect(result).toBeDefined();
    expect(result!.id).toBeUndefined();
  });
});

describe('businessRulesFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { id: 'br-1', ruleKey: 'maxAmount' },
      { id: 'br-2', ruleKey: 'minAmount' },
    ];

    const result = businessRulesFromDto(dtos);

    expect(result).toHaveLength(2);
    expect(result[0].ruleKey).toBe('maxAmount');
    expect(result[1].ruleKey).toBe('minAmount');
  });

  it('should return empty array for null input', () => {
    expect(businessRulesFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(businessRulesFromDto(undefined)).toEqual([]);
  });

  it('should return empty array for empty array', () => {
    expect(businessRulesFromDto([])).toEqual([]);
  });
});
