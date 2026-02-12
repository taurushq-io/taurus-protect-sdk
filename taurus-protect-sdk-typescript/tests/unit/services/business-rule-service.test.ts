/**
 * Unit tests for BusinessRuleService.
 */

import { BusinessRuleService } from '../../../src/services/business-rule-service';
import { ValidationError, NotFoundError } from '../../../src/errors';
import type { BusinessRulesApi } from '../../../src/internal/openapi/apis/BusinessRulesApi';

function createMockApi(): jest.Mocked<BusinessRulesApi> {
  return {
    ruleServiceGetBusinessRulesV2: jest.fn(),
  } as unknown as jest.Mocked<BusinessRulesApi>;
}

describe('BusinessRuleService', () => {
  let mockApi: jest.Mocked<BusinessRulesApi>;
  let service: BusinessRuleService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new BusinessRuleService(mockApi);
  });

  describe('list', () => {
    it('should return business rules', async () => {
      mockApi.ruleServiceGetBusinessRulesV2.mockResolvedValue({
        result: [
          { id: 'rule-1', ruleKey: 'MAX_AMOUNT', ruleValue: '1000' },
          { id: 'rule-2', ruleKey: 'MIN_APPROVAL', ruleValue: '2' },
        ],
      } as never);

      const result = await service.list();
      expect(result.rules).toHaveLength(2);
    });

    it('should handle empty results', async () => {
      mockApi.ruleServiceGetBusinessRulesV2.mockResolvedValue({
        result: [],
      } as never);

      const result = await service.list();
      expect(result.rules).toHaveLength(0);
    });

    it('should pass filter options to API', async () => {
      mockApi.ruleServiceGetBusinessRulesV2.mockResolvedValue({
        result: [],
      } as never);

      await service.list({ walletIds: ['123'], pageSize: 25 });

      expect(mockApi.ruleServiceGetBusinessRulesV2).toHaveBeenCalledWith(
        expect.objectContaining({
          walletIds: ['123'],
          cursorPageSize: '25',
        })
      );
    });
  });

  describe('get', () => {
    it('should throw ValidationError when ruleId is empty', async () => {
      await expect(service.get('')).rejects.toThrow(ValidationError);
      await expect(service.get('')).rejects.toThrow('ruleId is required');
    });

    it('should throw ValidationError when ruleId is whitespace', async () => {
      await expect(service.get('   ')).rejects.toThrow(ValidationError);
    });

    it('should throw NotFoundError when rule not found', async () => {
      mockApi.ruleServiceGetBusinessRulesV2.mockResolvedValue({
        result: [],
      } as never);

      await expect(service.get('nonexistent')).rejects.toThrow(NotFoundError);
    });

    it('should return rule when found', async () => {
      mockApi.ruleServiceGetBusinessRulesV2.mockResolvedValue({
        result: [{ id: 'rule-123', ruleKey: 'MAX_AMOUNT', ruleValue: '1000' }],
      } as never);

      const rule = await service.get('rule-123');
      expect(rule).toBeDefined();
      expect(rule.id).toBe('rule-123');
    });
  });
});
