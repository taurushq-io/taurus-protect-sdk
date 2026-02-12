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
    it('should return audit trails', async () => {
      mockApi.auditServiceGetAuditTrails.mockResolvedValue({
        auditTrails: [
          {
            id: '1',
            action: 'CREATE',
            entity: 'WALLET',
            user: { email: 'user@example.com' },
            date: new Date('2024-01-01'),
          },
          {
            id: '2',
            action: 'UPDATE',
            entity: 'ADDRESS',
            user: { email: 'admin@example.com' },
            date: new Date('2024-01-02'),
          },
        ],
      } as any);

      const result = await service.list();

      expect(result).toHaveLength(2);
    });

    it('should pass pagination options', async () => {
      mockApi.auditServiceGetAuditTrails.mockResolvedValue({
        auditTrails: [],
      } as any);

      await service.list({ limit: 100 });

      expect(mockApi.auditServiceGetAuditTrails).toHaveBeenCalledWith(
        expect.objectContaining({
          cursorPageSize: '100',
        })
      );
    });

    it('should throw ValidationError when limit is invalid', async () => {
      await expect(service.list({ limit: 0 })).rejects.toThrow(ValidationError);
      await expect(service.list({ limit: 0 })).rejects.toThrow('limit must be positive');
      await expect(service.list({ limit: -1 })).rejects.toThrow(ValidationError);
    });

    it('should use default limit when not provided', async () => {
      mockApi.auditServiceGetAuditTrails.mockResolvedValue({
        auditTrails: [],
      } as any);

      await service.list();

      expect(mockApi.auditServiceGetAuditTrails).toHaveBeenCalledWith(
        expect.objectContaining({
          cursorPageSize: '50',
        })
      );
    });

    it('should pass filter options to API', async () => {
      mockApi.auditServiceGetAuditTrails.mockResolvedValue({
        auditTrails: [],
      } as any);

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
