/**
 * Unit tests for ChangeService.
 */

import { ChangeService } from '../../../src/services/change-service';
import { ValidationError, NotFoundError } from '../../../src/errors';
import type { ChangesApi } from '../../../src/internal/openapi/apis/ChangesApi';

function createMockChangesApi(): jest.Mocked<ChangesApi> {
  return {
    changeServiceGetChange: jest.fn(),
    changeServiceGetChanges: jest.fn(),
    changeServiceGetChangesForApproval: jest.fn(),
    changeServiceApproveChange: jest.fn(),
    changeServiceApproveChanges: jest.fn(),
    changeServiceRejectChange: jest.fn(),
    changeServiceRejectChanges: jest.fn(),
    changeServiceCreateChange: jest.fn(),
  } as unknown as jest.Mocked<ChangesApi>;
}

describe('ChangeService', () => {
  let mockApi: jest.Mocked<ChangesApi>;
  let service: ChangeService;

  beforeEach(() => {
    mockApi = createMockChangesApi();
    service = new ChangeService(mockApi);
  });

  describe('get', () => {
    it('should return a change for a valid ID', async () => {
      mockApi.changeServiceGetChange.mockResolvedValue({
        result: {
          id: '123',
          action: 'CREATE',
          entity: 'user',
          status: 'Created',
        },
      } as any);

      const change = await service.get('123');

      expect(change).toBeDefined();
      expect(mockApi.changeServiceGetChange).toHaveBeenCalledWith({ id: '123' });
    });

    it('should throw ValidationError when id is empty', async () => {
      await expect(service.get('')).rejects.toThrow(ValidationError);
      await expect(service.get('')).rejects.toThrow('id is required');
    });

    it('should throw ValidationError when id is whitespace', async () => {
      await expect(service.get('   ')).rejects.toThrow(ValidationError);
    });

    it('should throw NotFoundError when change is not found', async () => {
      mockApi.changeServiceGetChange.mockResolvedValue({
        result: undefined,
      } as any);

      await expect(service.get('nonexistent')).rejects.toThrow(NotFoundError);
    });
  });

  describe('list', () => {
    it('should return a list of changes', async () => {
      mockApi.changeServiceGetChanges.mockResolvedValue({
        result: [
          { id: '1', action: 'CREATE', entity: 'user', status: 'Created' },
          { id: '2', action: 'UPDATE', entity: 'group', status: 'Created' },
        ],
      } as any);

      const result = await service.list();

      expect(result).toBeDefined();
    });

    it('should pass filter options to API', async () => {
      mockApi.changeServiceGetChanges.mockResolvedValue({
        result: [],
      } as any);

      await service.list({ entity: 'user', status: 'Created' });

      expect(mockApi.changeServiceGetChanges).toHaveBeenCalledWith(
        expect.objectContaining({
          entity: 'user',
          status: 'Created',
        })
      );
    });
  });

  describe('approve', () => {
    it('should approve a change', async () => {
      mockApi.changeServiceApproveChange.mockResolvedValue({} as any);

      await service.approve('123');

      expect(mockApi.changeServiceApproveChange).toHaveBeenCalledWith({
        id: '123',
        body: {},
      });
    });

    it('should throw ValidationError when id is empty', async () => {
      await expect(service.approve('')).rejects.toThrow(ValidationError);
      await expect(service.approve('')).rejects.toThrow('id is required');
    });

    it('should throw ValidationError when id is whitespace', async () => {
      await expect(service.approve('   ')).rejects.toThrow(ValidationError);
    });
  });

  describe('approveMany', () => {
    it('should approve multiple changes', async () => {
      mockApi.changeServiceApproveChanges.mockResolvedValue({} as any);

      await service.approveMany(['123', '456']);

      expect(mockApi.changeServiceApproveChanges).toHaveBeenCalledWith({
        body: { ids: ['123', '456'] },
      });
    });

    it('should throw ValidationError when ids is empty', async () => {
      await expect(service.approveMany([])).rejects.toThrow(ValidationError);
      await expect(service.approveMany([])).rejects.toThrow('ids cannot be empty');
    });
  });

  describe('reject', () => {
    it('should reject a change', async () => {
      mockApi.changeServiceRejectChange.mockResolvedValue({} as any);

      await service.reject('123');

      expect(mockApi.changeServiceRejectChange).toHaveBeenCalledWith({
        id: '123',
        body: {},
      });
    });

    it('should throw ValidationError when id is empty', async () => {
      await expect(service.reject('')).rejects.toThrow(ValidationError);
      await expect(service.reject('')).rejects.toThrow('id is required');
    });
  });

  describe('rejectMany', () => {
    it('should reject multiple changes', async () => {
      mockApi.changeServiceRejectChanges.mockResolvedValue({} as any);

      await service.rejectMany(['123', '456']);

      expect(mockApi.changeServiceRejectChanges).toHaveBeenCalledWith({
        body: { ids: ['123', '456'] },
      });
    });

    it('should throw ValidationError when ids is empty', async () => {
      await expect(service.rejectMany([])).rejects.toThrow(ValidationError);
      await expect(service.rejectMany([])).rejects.toThrow('ids cannot be empty');
    });
  });

  describe('listForApproval', () => {
    it('should return changes pending approval', async () => {
      mockApi.changeServiceGetChangesForApproval.mockResolvedValue({
        result: [
          { id: '1', action: 'CREATE', entity: 'user', status: 'Created' },
        ],
      } as any);

      const result = await service.listForApproval();

      expect(result).toBeDefined();
    });

    it('should pass filter options', async () => {
      mockApi.changeServiceGetChangesForApproval.mockResolvedValue({
        result: [],
      } as any);

      await service.listForApproval({ entities: ['user', 'group'] });

      expect(mockApi.changeServiceGetChangesForApproval).toHaveBeenCalledWith(
        expect.objectContaining({
          entities: ['user', 'group'],
        })
      );
    });
  });

  describe('create', () => {
    it('should create a change and return the ID', async () => {
      mockApi.changeServiceCreateChange.mockResolvedValue({
        result: { id: '42' },
      } as any);

      const changeId = await service.create({
        action: 'update',
        entity: 'businessrule',
        entityId: '10',
        changes: { rulevalue: '100' },
        comment: 'test change',
      });

      expect(changeId).toBe('42');
      expect(mockApi.changeServiceCreateChange).toHaveBeenCalledWith({
        body: expect.objectContaining({
          action: 'update',
          entity: 'businessrule',
          entityId: '10',
          changes: { rulevalue: '100' },
          changeComment: 'test change',
        }),
      });
    });

    it('should map comment to changeComment in DTO', async () => {
      mockApi.changeServiceCreateChange.mockResolvedValue({
        result: { id: '1' },
      } as any);

      await service.create({
        action: 'create',
        entity: 'user',
        comment: 'my comment',
      });

      expect(mockApi.changeServiceCreateChange).toHaveBeenCalledWith({
        body: expect.objectContaining({
          changeComment: 'my comment',
        }),
      });
    });

    it('should return empty string when result has no id', async () => {
      mockApi.changeServiceCreateChange.mockResolvedValue({
        result: {},
      } as any);

      const changeId = await service.create({
        action: 'update',
        entity: 'businessrule',
      });

      expect(changeId).toBe('');
    });

    it('should throw ValidationError when request is null', async () => {
      await expect(service.create(null as any)).rejects.toThrow(ValidationError);
      await expect(service.create(null as any)).rejects.toThrow('request is required');
    });

    it('should throw ValidationError when action is empty', async () => {
      await expect(service.create({ action: '', entity: 'user' })).rejects.toThrow(ValidationError);
      await expect(service.create({ action: '', entity: 'user' })).rejects.toThrow('action is required');
    });

    it('should throw ValidationError when action is whitespace', async () => {
      await expect(service.create({ action: '  ', entity: 'user' })).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when entity is empty', async () => {
      await expect(service.create({ action: 'update', entity: '' })).rejects.toThrow(ValidationError);
      await expect(service.create({ action: 'update', entity: '' })).rejects.toThrow('entity is required');
    });

    it('should throw ValidationError when entity is whitespace', async () => {
      await expect(service.create({ action: 'update', entity: '   ' })).rejects.toThrow(ValidationError);
    });
  });
});
