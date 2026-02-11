/**
 * Unit tests for UserDeviceService.
 */

import { UserDeviceService } from '../../../src/services/user-device-service';
import { ValidationError } from '../../../src/errors';
import type { UserDeviceApi } from '../../../src/internal/openapi/apis/UserDeviceApi';

function createMockApi(): jest.Mocked<UserDeviceApi> {
  return {
    userDeviceServiceCreateUserDevicePairing: jest.fn(),
    userDeviceServiceGetUserDevicePairingStatus: jest.fn(),
    userDeviceServiceStartUserDevicePairing: jest.fn(),
    userDeviceServiceApproveUserDevicePairing: jest.fn(),
  } as unknown as jest.Mocked<UserDeviceApi>;
}

describe('UserDeviceService', () => {
  let mockApi: jest.Mocked<UserDeviceApi>;
  let service: UserDeviceService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new UserDeviceService(mockApi);
  });

  describe('createPairing', () => {
    it('should create a device pairing and return the result', async () => {
      mockApi.userDeviceServiceCreateUserDevicePairing.mockResolvedValue({
        pairingId: 'pairing-123',
      } as never);

      const pairing = await service.createPairing();
      expect(pairing).toBeDefined();
      expect(mockApi.userDeviceServiceCreateUserDevicePairing).toHaveBeenCalled();
    });

    it('should throw ValidationError when API returns invalid response', async () => {
      mockApi.userDeviceServiceCreateUserDevicePairing.mockResolvedValue({} as never);

      await expect(service.createPairing()).rejects.toThrow(ValidationError);
    });
  });

  describe('getPairingStatus', () => {
    it('should throw ValidationError when pairingId is empty', async () => {
      await expect(service.getPairingStatus('', '123456')).rejects.toThrow(ValidationError);
      await expect(service.getPairingStatus('', '123456')).rejects.toThrow('pairingId is required');
    });

    it('should throw ValidationError when nonce is empty', async () => {
      await expect(service.getPairingStatus('pairing-1', '')).rejects.toThrow(ValidationError);
      await expect(service.getPairingStatus('pairing-1', '')).rejects.toThrow('nonce is required');
    });

    it('should throw ValidationError when pairingId is whitespace', async () => {
      await expect(service.getPairingStatus('  ', '123456')).rejects.toThrow(ValidationError);
    });

    it('should return pairing status', async () => {
      mockApi.userDeviceServiceGetUserDevicePairingStatus.mockResolvedValue({
        status: 'PENDING',
        pairingId: 'pairing-123',
      } as never);

      const info = await service.getPairingStatus('pairing-123', '654321');
      expect(info).toBeDefined();
      expect(mockApi.userDeviceServiceGetUserDevicePairingStatus).toHaveBeenCalledWith({
        pairingID: 'pairing-123',
        nonce: '654321',
      });
    });
  });

  describe('startPairing', () => {
    it('should throw ValidationError when pairingId is empty', async () => {
      await expect(
        service.startPairing('', { nonce: '123456', publicKey: 'key' })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.startPairing('', { nonce: '123456', publicKey: 'key' })
      ).rejects.toThrow('pairingId is required');
    });

    it('should throw ValidationError when pairingId is whitespace', async () => {
      await expect(
        service.startPairing('  ', { nonce: '123456', publicKey: 'key' })
      ).rejects.toThrow(ValidationError);
    });

    it('should start pairing with correct parameters', async () => {
      mockApi.userDeviceServiceStartUserDevicePairing.mockResolvedValue({} as never);

      await service.startPairing('pairing-123', {
        nonce: '654321',
        publicKey: 'base64-public-key',
      });

      expect(mockApi.userDeviceServiceStartUserDevicePairing).toHaveBeenCalledWith({
        pairingID: 'pairing-123',
        body: {
          nonce: '654321',
          publicKey: 'base64-public-key',
        },
      });
    });
  });

  describe('approvePairing', () => {
    it('should throw ValidationError when pairingId is empty', async () => {
      await expect(
        service.approvePairing('', { nonce: '123456' })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.approvePairing('', { nonce: '123456' })
      ).rejects.toThrow('pairingId is required');
    });

    it('should throw ValidationError when nonce is empty', async () => {
      await expect(
        service.approvePairing('pairing-1', { nonce: '' })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.approvePairing('pairing-1', { nonce: '' })
      ).rejects.toThrow('nonce is required');
    });

    it('should approve pairing with correct parameters', async () => {
      mockApi.userDeviceServiceApproveUserDevicePairing.mockResolvedValue({} as never);

      await service.approvePairing('pairing-123', { nonce: '654321' });

      expect(mockApi.userDeviceServiceApproveUserDevicePairing).toHaveBeenCalledWith({
        pairingID: 'pairing-123',
        body: {
          nonce: '654321',
        },
      });
    });
  });
});
