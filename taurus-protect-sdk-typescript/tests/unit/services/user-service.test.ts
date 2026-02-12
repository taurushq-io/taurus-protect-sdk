/**
 * Unit tests for UserService.
 */

import { UserService } from '../../../src/services/user-service';
import { ValidationError, NotFoundError } from '../../../src/errors';
import type { UsersApi } from '../../../src/internal/openapi/apis/UsersApi';

function createMockUsersApi(): jest.Mocked<UsersApi> {
  return {
    userServiceGetUser: jest.fn(),
    userServiceGetMe: jest.fn(),
    userServiceGetUsers: jest.fn(),
  } as unknown as jest.Mocked<UsersApi>;
}

describe('UserService', () => {
  let mockApi: jest.Mocked<UsersApi>;
  let service: UserService;

  beforeEach(() => {
    mockApi = createMockUsersApi();
    service = new UserService(mockApi);
  });

  describe('getCurrentUser', () => {
    it('should return the current user', async () => {
      mockApi.userServiceGetMe.mockResolvedValue({
        result: {
          id: 'user-1',
          firstName: 'John',
          lastName: 'Doe',
          email: 'john@example.com',
        },
      } as any);

      const user = await service.getCurrentUser();

      expect(user).toBeDefined();
      expect(user.firstName).toBe('John');
      expect(user.lastName).toBe('Doe');
      expect(user.email).toBe('john@example.com');
    });

    it('should throw when no result is returned', async () => {
      mockApi.userServiceGetMe.mockResolvedValue({
        result: undefined,
      } as any);

      await expect(service.getCurrentUser()).rejects.toThrow();
    });
  });

  describe('get', () => {
    it('should return a user for a valid ID', async () => {
      mockApi.userServiceGetUser.mockResolvedValue({
        result: {
          id: 'user-123',
          firstName: 'Jane',
          lastName: 'Smith',
          email: 'jane@example.com',
        },
      } as any);

      const user = await service.get('user-123');

      expect(user).toBeDefined();
      expect(user.firstName).toBe('Jane');
      expect(mockApi.userServiceGetUser).toHaveBeenCalledWith({ id: 'user-123' });
    });

    it('should throw ValidationError when userId is empty', async () => {
      await expect(service.get('')).rejects.toThrow(ValidationError);
      await expect(service.get('')).rejects.toThrow('userId is required');
    });

    it('should throw ValidationError when userId is whitespace', async () => {
      await expect(service.get('   ')).rejects.toThrow(ValidationError);
    });

    it('should throw NotFoundError when user is not found', async () => {
      mockApi.userServiceGetUser.mockResolvedValue({
        result: undefined,
      } as any);

      await expect(service.get('nonexistent')).rejects.toThrow(NotFoundError);
    });
  });

  describe('list', () => {
    it('should return paginated users', async () => {
      mockApi.userServiceGetUsers.mockResolvedValue({
        result: [
          { id: '1', firstName: 'User', lastName: 'One', email: 'user1@example.com' },
          { id: '2', firstName: 'User', lastName: 'Two', email: 'user2@example.com' },
        ],
        totalItems: '10',
      } as any);

      const result = await service.list({ limit: 50 });

      expect(result.items).toHaveLength(2);
      expect(result.pagination.totalItems).toBe(10);
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
      mockApi.userServiceGetUsers.mockResolvedValue({
        result: [],
        totalItems: '0',
      } as any);

      await service.list();

      expect(mockApi.userServiceGetUsers).toHaveBeenCalledWith(
        expect.objectContaining({
          limit: '50',
          offset: '0',
        })
      );
    });
  });
});
