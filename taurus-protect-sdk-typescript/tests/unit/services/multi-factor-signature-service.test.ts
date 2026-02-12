/**
 * Unit tests for MultiFactorSignatureService.
 */

import { MultiFactorSignatureService } from '../../../src/services/multi-factor-signature-service';
import { NotFoundError, ValidationError } from '../../../src/errors';
import { MultiFactorSignatureEntityType } from '../../../src/models/multi-factor-signature';
import type { MultiFactorSignatureApi } from '../../../src/internal/openapi/apis/MultiFactorSignatureApi';

function createMockApi(): jest.Mocked<MultiFactorSignatureApi> {
  return {
    multiFactorSignatureServiceGetMultiFactorSignatureEntitiesInfo: jest.fn(),
    multiFactorSignatureServiceCreateMultiFactorSignatureBatch: jest.fn(),
    multiFactorSignatureServiceApproveMultiFactorSignature: jest.fn(),
    multiFactorSignatureServiceRejectMultiFactorSignature: jest.fn(),
  } as unknown as jest.Mocked<MultiFactorSignatureApi>;
}

describe('MultiFactorSignatureService', () => {
  let mockApi: jest.Mocked<MultiFactorSignatureApi>;
  let service: MultiFactorSignatureService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new MultiFactorSignatureService(mockApi);
  });

  describe('get', () => {
    it('should throw ValidationError when id is empty', async () => {
      await expect(service.get('')).rejects.toThrow(ValidationError);
      await expect(service.get('')).rejects.toThrow('id is required');
    });

    it('should throw ValidationError when id is whitespace', async () => {
      await expect(service.get('  ')).rejects.toThrow(ValidationError);
    });

    it('should return multi-factor signature info', async () => {
      mockApi.multiFactorSignatureServiceGetMultiFactorSignatureEntitiesInfo.mockResolvedValue({
        id: 'mfs-123',
        entities: [
          { entityId: 'req-1', payloadToSign: 'payload1' },
        ],
      } as never);

      const result = await service.get('mfs-123');
      expect(result).toBeDefined();
      expect(mockApi.multiFactorSignatureServiceGetMultiFactorSignatureEntitiesInfo).toHaveBeenCalledWith({
        id: 'mfs-123',
      });
    });
  });

  describe('create', () => {
    it('should create a multi-factor signature batch and return ID', async () => {
      mockApi.multiFactorSignatureServiceCreateMultiFactorSignatureBatch.mockResolvedValue({
        id: 'mfs-456',
      });

      const id = await service.create({
        entityType: MultiFactorSignatureEntityType.REQUEST,
        entityIds: ['req-1', 'req-2'],
      });

      expect(id).toBe('mfs-456');
      expect(mockApi.multiFactorSignatureServiceCreateMultiFactorSignatureBatch).toHaveBeenCalledWith({
        body: expect.objectContaining({
          entityIDs: ['req-1', 'req-2'],
        }),
      });
    });

    it('should throw ValidationError when entityIds is empty', async () => {
      await expect(
        service.create({ entityType: MultiFactorSignatureEntityType.REQUEST, entityIds: [] })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.create({ entityType: MultiFactorSignatureEntityType.REQUEST, entityIds: [] })
      ).rejects.toThrow('entityIds cannot be empty');
    });

    it('should throw ValidationError when entityType is missing', async () => {
      await expect(
        service.create({ entityType: '' as never, entityIds: ['1'] })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.create({ entityType: '' as never, entityIds: ['1'] })
      ).rejects.toThrow('entityType is required');
    });

    it('should handle REQUEST entity type', async () => {
      mockApi.multiFactorSignatureServiceCreateMultiFactorSignatureBatch.mockResolvedValue({
        id: 'mfs-1',
      });

      await service.create({ entityType: MultiFactorSignatureEntityType.REQUEST, entityIds: ['1'] });

      expect(mockApi.multiFactorSignatureServiceCreateMultiFactorSignatureBatch).toHaveBeenCalledWith({
        body: expect.objectContaining({
          entityType: expect.anything(),
        }),
      });
    });

    it('should handle WHITELISTED_ADDRESS entity type', async () => {
      mockApi.multiFactorSignatureServiceCreateMultiFactorSignatureBatch.mockResolvedValue({
        id: 'mfs-2',
      });

      await service.create({
        entityType: MultiFactorSignatureEntityType.WHITELISTED_ADDRESS,
        entityIds: ['addr-1'],
      });
      expect(mockApi.multiFactorSignatureServiceCreateMultiFactorSignatureBatch).toHaveBeenCalled();
    });

    it('should handle WHITELISTED_CONTRACT entity type', async () => {
      mockApi.multiFactorSignatureServiceCreateMultiFactorSignatureBatch.mockResolvedValue({
        id: 'mfs-3',
      });

      await service.create({
        entityType: MultiFactorSignatureEntityType.WHITELISTED_CONTRACT,
        entityIds: ['c-1'],
      });
      expect(mockApi.multiFactorSignatureServiceCreateMultiFactorSignatureBatch).toHaveBeenCalled();
    });
  });

  describe('approve', () => {
    it('should approve a multi-factor signature', async () => {
      mockApi.multiFactorSignatureServiceApproveMultiFactorSignature.mockResolvedValue({} as never);

      await service.approve({
        id: 'mfs-123',
        signature: 'base64-signature',
        comment: 'Approved via SDK',
      });

      expect(mockApi.multiFactorSignatureServiceApproveMultiFactorSignature).toHaveBeenCalledWith({
        id: 'mfs-123',
        body: {
          signature: 'base64-signature',
          comment: 'Approved via SDK',
        },
      });
    });

    it('should throw ValidationError when id is empty', async () => {
      await expect(
        service.approve({ id: '', signature: 'sig', comment: 'ok' })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.approve({ id: '', signature: 'sig', comment: 'ok' })
      ).rejects.toThrow('id is required');
    });

    it('should throw ValidationError when signature is empty', async () => {
      await expect(
        service.approve({ id: 'mfs-1', signature: '', comment: 'ok' })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.approve({ id: 'mfs-1', signature: '', comment: 'ok' })
      ).rejects.toThrow('signature is required');
    });

    it('should default comment to empty string when not provided', async () => {
      mockApi.multiFactorSignatureServiceApproveMultiFactorSignature.mockResolvedValue({} as never);

      await service.approve({ id: 'mfs-1', signature: 'sig' });

      expect(mockApi.multiFactorSignatureServiceApproveMultiFactorSignature).toHaveBeenCalledWith({
        id: 'mfs-1',
        body: { signature: 'sig', comment: '' },
      });
    });
  });

  describe('reject', () => {
    it('should reject a multi-factor signature', async () => {
      mockApi.multiFactorSignatureServiceRejectMultiFactorSignature.mockResolvedValue({} as never);

      await service.reject({ id: 'mfs-123', comment: 'Not needed' });

      expect(mockApi.multiFactorSignatureServiceRejectMultiFactorSignature).toHaveBeenCalledWith({
        id: 'mfs-123',
        body: { comment: 'Not needed' },
      });
    });

    it('should throw ValidationError when id is empty', async () => {
      await expect(
        service.reject({ id: '', comment: 'reason' })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.reject({ id: '', comment: 'reason' })
      ).rejects.toThrow('id is required');
    });

    it('should default comment to empty string when not provided', async () => {
      mockApi.multiFactorSignatureServiceRejectMultiFactorSignature.mockResolvedValue({} as never);

      await service.reject({ id: 'mfs-1' });

      expect(mockApi.multiFactorSignatureServiceRejectMultiFactorSignature).toHaveBeenCalledWith({
        id: 'mfs-1',
        body: { comment: '' },
      });
    });
  });
});
