/**
 * Unit tests for AuditService.
 */

import { AuditService } from '../../../src/services/audit-service';
import { ValidationError } from '../../../src/errors';
import type { AuditApi } from '../../../src/internal/openapi/apis/AuditApi';

function createMockAuditApi(): jest.Mocked<AuditApi> {
  return {
    auditServiceGetAuditTrails: jest.fn(),
    auditServiceExportAuditTrails: jest.fn(),
  } as unknown as jest.Mocked<AuditApi>;
}

describe('AuditService', () => {
  let mockApi: jest.Mocked<AuditApi>;
  let service: AuditService;

  beforeEach(() => {
    mockApi = createMockAuditApi();
    service = new AuditService(mockApi);
  });

  describe('list', () => {
    it('should read the rows from `result` and the actor from the nested user', async () => {
      mockApi.auditServiceGetAuditTrails.mockResolvedValue({
        result: [
          {
            id: '1',
            action: 'CREATE',
            entity: 'WALLET',
            user: { id: 'u1', email: 'user@example.com', externalUserId: 'ext-1' },
            creationDate: new Date('2024-01-01'),
          },
          {
            id: '2',
            action: 'UPDATE',
            entity: 'ADDRESS',
            user: { email: 'admin@example.com' },
            creationDate: new Date('2024-01-02'),
          },
        ],
        cursor: { currentPage: 'next-page', hasNext: true },
      });

      const result = await service.list();

      expect(result.items).toHaveLength(2);
      expect(result.items[0]).toMatchObject({
        id: '1',
        userId: 'u1',
        userEmail: 'user@example.com',
        externalUserId: 'ext-1',
      });
      expect(result.pagination).toEqual({
        pageSize: 20,
        nextCursor: 'next-page',
        hasMore: true,
      });
    });

    it('should pass the page size and continue from a cursor', async () => {
      mockApi.auditServiceGetAuditTrails.mockResolvedValue({ result: [] });

      await service.list({ pageSize: 100, cursor: 'next-page' });

      expect(mockApi.auditServiceGetAuditTrails).toHaveBeenCalledWith(
        expect.objectContaining({
          cursorPageSize: '100',
          cursorCurrentPage: 'next-page',
          cursorPageRequest: 'NEXT',
        })
      );
    });

    it('should throw ValidationError when the page size is out of bounds', async () => {
      await expect(service.list({ pageSize: 101 })).rejects.toThrow(ValidationError);
      await expect(service.list({ pageSize: 101 })).rejects.toThrow('pageSize must be at most 100, got 101');
      await expect(service.list({ pageSize: -1 })).rejects.toThrow(ValidationError);
      expect(mockApi.auditServiceGetAuditTrails).not.toHaveBeenCalled();
    });

    it('should send the default page size when none is given', async () => {
      mockApi.auditServiceGetAuditTrails.mockResolvedValue({});

      const result = await service.list();

      expect(mockApi.auditServiceGetAuditTrails).toHaveBeenCalledWith(
        expect.objectContaining({
          cursorPageSize: '20',
        })
      );
      expect(result).toEqual({
        items: [],
        pagination: { pageSize: 20, nextCursor: '', hasMore: false },
      });
    });

    it('should pass filter options to API', async () => {
      mockApi.auditServiceGetAuditTrails.mockResolvedValue({ result: [] });

      await service.list({
        entities: ['WALLET'],
        actions: ['CREATE'],
        externalUserId: 'user-123',
      });

      expect(mockApi.auditServiceGetAuditTrails).toHaveBeenCalledWith(
        expect.objectContaining({
          entities: ['WALLET'],
          actions: ['CREATE'],
          externalUserId: 'user-123',
        })
      );
    });
  });

  describe('exportAuditTrails', () => {
    it('should export audit trails', async () => {
      mockApi.auditServiceExportAuditTrails.mockResolvedValue({
        result: 'id,action,entity\n1,CREATE,WALLET',
      });

      const result = await service.exportAuditTrails();

      expect(result).toContain('CREATE');
    });

    it('should pass filter options to export', async () => {
      mockApi.auditServiceExportAuditTrails.mockResolvedValue({
        result: '',
      });

      await service.exportAuditTrails({
        entities: ['WALLET'],
        format: 'json',
      });

      expect(mockApi.auditServiceExportAuditTrails).toHaveBeenCalledWith(
        expect.objectContaining({
          entities: ['WALLET'],
          format: 'json',
        })
      );
    });
  });
});
