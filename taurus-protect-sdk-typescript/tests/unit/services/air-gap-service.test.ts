/**
 * Unit tests for AirGapService.
 */

import { AirGapService } from '../../../src/services/air-gap-service';
import { ValidationError } from '../../../src/errors';
import type { AirGapApi } from '../../../src/internal/openapi/apis/AirGapApi';

function createMockApi(): jest.Mocked<AirGapApi> {
  return {
    airGapServiceGetOutgoingAirGap: jest.fn(),
    airGapServiceSubmitIncomingAirGap: jest.fn(),
  } as unknown as jest.Mocked<AirGapApi>;
}

describe('AirGapService', () => {
  let mockApi: jest.Mocked<AirGapApi>;
  let service: AirGapService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new AirGapService(mockApi);
  });

  describe('getOutgoingAirGap', () => {
    it('should throw ValidationError when requestIds is empty', async () => {
      await expect(
        service.getOutgoingAirGap({ requestIds: [] })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.getOutgoingAirGap({ requestIds: [] })
      ).rejects.toThrow('requestIds cannot be empty');
    });

    it('should call API with request IDs', async () => {
      const mockBlob = new Blob(['test']);
      mockApi.airGapServiceGetOutgoingAirGap.mockResolvedValue(mockBlob as never);

      const result = await service.getOutgoingAirGap({ requestIds: ['req-1', 'req-2'] });
      expect(result).toBeDefined();
      expect(mockApi.airGapServiceGetOutgoingAirGap).toHaveBeenCalledWith({
        body: {
          requests: {
            ids: ['req-1', 'req-2'],
            signature: undefined,
          },
        },
      });
    });
  });

  describe('getOutgoingAirGapAddresses', () => {
    it('should throw ValidationError when addressIds is empty', async () => {
      await expect(
        service.getOutgoingAirGapAddresses({ addressIds: [] })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.getOutgoingAirGapAddresses({ addressIds: [] })
      ).rejects.toThrow('addressIds cannot be empty');
    });

    it('should call API with address IDs', async () => {
      const mockBlob = new Blob(['test']);
      mockApi.airGapServiceGetOutgoingAirGap.mockResolvedValue(mockBlob as never);

      const result = await service.getOutgoingAirGapAddresses({ addressIds: ['addr-1'] });
      expect(result).toBeDefined();
      expect(mockApi.airGapServiceGetOutgoingAirGap).toHaveBeenCalled();
    });
  });

  describe('submitIncomingAirGap', () => {
    it('should throw ValidationError when payload is empty', async () => {
      await expect(
        service.submitIncomingAirGap({ payload: '' })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.submitIncomingAirGap({ payload: '' })
      ).rejects.toThrow('payload cannot be empty');
    });

    it('should throw ValidationError when payload is whitespace', async () => {
      await expect(
        service.submitIncomingAirGap({ payload: '  ' })
      ).rejects.toThrow(ValidationError);
    });

    it('should submit incoming payload', async () => {
      mockApi.airGapServiceSubmitIncomingAirGap.mockResolvedValue({} as never);

      await service.submitIncomingAirGap({ payload: 'base64-signed-payload' });
      expect(mockApi.airGapServiceSubmitIncomingAirGap).toHaveBeenCalledWith({
        body: expect.objectContaining({
          payload: 'base64-signed-payload',
        }),
      });
    });
  });
});
