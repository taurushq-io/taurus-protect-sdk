/**
 * Unit tests for GroupService.
 */

import { GroupService } from '../../../src/services/group-service';
import { ValidationError, NotFoundError } from '../../../src/errors';
import type { GroupsApi } from '../../../src/internal/openapi/apis/GroupsApi';

function createMockGroupsApi(): jest.Mocked<GroupsApi> {
  return {
    userServiceGetGroups: jest.fn(),
  } as unknown as jest.Mocked<GroupsApi>;
}

describe('GroupService', () => {
  let mockApi: jest.Mocked<GroupsApi>;
  let service: GroupService;

  beforeEach(() => {
    mockApi = createMockGroupsApi();
    service = new GroupService(mockApi);
  });

  describe('get', () => {
    it('should return a group for a valid ID', async () => {
      mockApi.userServiceGetGroups.mockResolvedValue({
        result: [
          {
            id: 'group-123',
            name: 'Admins',
            userIds: ['user-1', 'user-2'],
          },
        ],
        totalItems: '1',
      } as any);

      const group = await service.get('group-123');

      expect(group).toBeDefined();
      expect(group.name).toBe('Admins');
      expect(mockApi.userServiceGetGroups).toHaveBeenCalledWith({
        limit: '1',
        offset: '0',
        ids: ['group-123'],
      });
    });

    it('should throw ValidationError when groupId is empty', async () => {
      await expect(service.get('')).rejects.toThrow(ValidationError);
      await expect(service.get('')).rejects.toThrow('groupId is required');
    });

    it('should throw ValidationError when groupId is whitespace', async () => {
      await expect(service.get('   ')).rejects.toThrow(ValidationError);
    });

    it('should throw NotFoundError when group is not found', async () => {
      mockApi.userServiceGetGroups.mockResolvedValue({
        result: [],
        totalItems: '0',
      } as any);

      await expect(service.get('nonexistent')).rejects.toThrow(NotFoundError);
    });
  });

  describe('list', () => {
    it('should return paginated groups', async () => {
      mockApi.userServiceGetGroups.mockResolvedValue({
        result: [
          { id: '1', name: 'Admins' },
          { id: '2', name: 'Traders' },
        ],
        totalItems: '5',
      } as any);

      const result = await service.list({ limit: 50 });

      expect(result.items).toHaveLength(2);
      expect(result.pagination.totalItems).toBe(5);
      expect(result.pagination.limit).toBe(50);
    });

    it('should throw ValidationError when limit is invalid', async () => {
      await expect(service.list({ limit: 0 })).rejects.toThrow(ValidationError);
      await expect(service.list({ limit: 0 })).rejects.toThrow('limit must be positive');
    });

    it('should throw ValidationError when offset is negative', async () => {
      await expect(service.list({ offset: -1 })).rejects.toThrow(ValidationError);
      await expect(service.list({ offset: -1 })).rejects.toThrow('offset cannot be negative');
    });

    it('should use default limit and offset when not provided', async () => {
      mockApi.userServiceGetGroups.mockResolvedValue({
        result: [],
        totalItems: '0',
      } as any);

      await service.list();

      expect(mockApi.userServiceGetGroups).toHaveBeenCalledWith(
        expect.objectContaining({
          limit: '50',
          offset: '0',
        })
      );
    });

    it('should pass query filter to API', async () => {
      mockApi.userServiceGetGroups.mockResolvedValue({
        result: [],
        totalItems: '0',
      } as any);

      await service.list({ query: 'admin' });

      expect(mockApi.userServiceGetGroups).toHaveBeenCalledWith(
        expect.objectContaining({
          query: 'admin',
        })
      );
    });
  });
});
