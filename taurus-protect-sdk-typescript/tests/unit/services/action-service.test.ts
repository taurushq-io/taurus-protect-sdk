/**
 * Unit tests for ActionService.
 */

import { ActionService } from '../../../src/services/action-service';
import { NotFoundError, ValidationError } from '../../../src/errors';
import type { ActionsApi } from '../../../src/internal/openapi/apis/ActionsApi';

function createMockApi(): jest.Mocked<ActionsApi> {
  return {
    actionServiceGetActions: jest.fn(),
    actionServiceGetAction: jest.fn(),
  } as unknown as jest.Mocked<ActionsApi>;
}

describe('ActionService', () => {
  let mockApi: jest.Mocked<ActionsApi>;
  let service: ActionService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new ActionService(mockApi);
  });

  describe('list', () => {
    it('should return actions', async () => {
      mockApi.actionServiceGetActions.mockResolvedValue({
        result: [
          { id: 'action-1', label: 'Auto-transfer', status: 'ACTIVE' },
          { id: 'action-2', label: 'Balance alert', status: 'INACTIVE' },
        ],
      } as never);

      const actions = await service.list();
      expect(actions).toBeDefined();
      expect(actions.length).toBeGreaterThanOrEqual(0);
    });

    it('should pass options to API', async () => {
      mockApi.actionServiceGetActions.mockResolvedValue({
        result: [],
      } as never);

      await service.list({ limit: '10', offset: '5', ids: ['action-1'] });

      expect(mockApi.actionServiceGetActions).toHaveBeenCalledWith({
        limit: '10',
        offset: '5',
        ids: ['action-1'],
      });
    });

    it('should handle empty results', async () => {
      mockApi.actionServiceGetActions.mockResolvedValue({
        result: [],
      } as never);

      const actions = await service.list();
      expect(actions).toHaveLength(0);
    });

    it('should work without options', async () => {
      mockApi.actionServiceGetActions.mockResolvedValue({
        result: [],
      } as never);

      const actions = await service.list();
      expect(actions).toBeDefined();
      expect(mockApi.actionServiceGetActions).toHaveBeenCalledWith({
        limit: undefined,
        offset: undefined,
        ids: undefined,
      });
    });
  });

  describe('get', () => {
    it('should throw ValidationError when actionId is empty', async () => {
      await expect(service.get('')).rejects.toThrow(ValidationError);
      await expect(service.get('')).rejects.toThrow('actionId is required');
    });

    it('should throw ValidationError when actionId is whitespace', async () => {
      await expect(service.get('  ')).rejects.toThrow(ValidationError);
    });

    it('should throw NotFoundError when action is not found', async () => {
      mockApi.actionServiceGetAction.mockResolvedValue({
        action: undefined,
        result: undefined,
      } as never);

      await expect(service.get('nonexistent')).rejects.toThrow(NotFoundError);
    });

    it('should return action for valid id', async () => {
      mockApi.actionServiceGetAction.mockResolvedValue({
        action: {
          id: 'action-123',
          label: 'Auto-transfer',
          status: 'ACTIVE',
          autoApprove: true,
        },
      } as never);

      const action = await service.get('action-123');
      expect(action).toBeDefined();
      expect(mockApi.actionServiceGetAction).toHaveBeenCalledWith({ id: 'action-123' });
    });
  });
});
