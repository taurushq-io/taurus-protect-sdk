/**
 * Unit tests for VisibilityGroupService.
 */

import { VisibilityGroupService } from '../../../src/services/visibility-group-service';
import { ValidationError } from '../../../src/errors';
import type { RestrictedVisibilityGroupsApi } from '../../../src/internal/openapi/apis/RestrictedVisibilityGroupsApi';

function createMockApi(): jest.Mocked<RestrictedVisibilityGroupsApi> {
  return {
    userServiceGetVisibilityGroups: jest.fn(),
    userServiceGetUsersByVisibilityGroupID: jest.fn(),
  } as unknown as jest.Mocked<RestrictedVisibilityGroupsApi>;
}

describe('VisibilityGroupService', () => {
  let mockApi: jest.Mocked<RestrictedVisibilityGroupsApi>;
  let service: VisibilityGroupService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new VisibilityGroupService(mockApi);
  });

  describe('list', () => {
    it('should return visibility groups', async () => {
      mockApi.userServiceGetVisibilityGroups.mockResolvedValue({
        result: [
          { id: 'vg-1', name: 'Group A' },
          { id: 'vg-2', name: 'Group B' },
        ],
      } as never);

      const groups = await service.list();
      expect(groups).toHaveLength(2);
    });

    it('should handle empty results', async () => {
      mockApi.userServiceGetVisibilityGroups.mockResolvedValue({
        result: [],
      } as never);

      const groups = await service.list();
      expect(groups).toHaveLength(0);
    });
  });

  describe('getUsersByVisibilityGroup', () => {
    it('should throw ValidationError when visibilityGroupId is empty', async () => {
      await expect(service.getUsersByVisibilityGroup('')).rejects.toThrow(ValidationError);
    });

    it('should return users for a visibility group', async () => {
      mockApi.userServiceGetUsersByVisibilityGroupID.mockResolvedValue({
        result: [
          { id: 'user-1', firstName: 'Alice' },
        ],
      } as never);

      const users = await service.getUsersByVisibilityGroup('vg-123');
      expect(users).toBeDefined();
    });
  });
});
