/**
 * Unit tests for ScoreService.
 */

import { ScoreService } from '../../../src/services/score-service';
import { ValidationError } from '../../../src/errors';
import type { ScoresApi } from '../../../src/internal/openapi/apis/ScoresApi';

function createMockApi(): jest.Mocked<ScoresApi> {
  return {
    scoreServiceRefreshAddressScore: jest.fn(),
    scoreServiceRefreshWLAScore: jest.fn(),
  } as unknown as jest.Mocked<ScoresApi>;
}

describe('ScoreService', () => {
  let mockApi: jest.Mocked<ScoresApi>;
  let service: ScoreService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new ScoreService(mockApi);
  });

  describe('refreshAddressScore', () => {
    it('should throw ValidationError when addressId is 0', async () => {
      await expect(service.refreshAddressScore(0, 'chainalysis')).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when addressId is negative', async () => {
      await expect(service.refreshAddressScore(-1, 'chainalysis')).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when scoreProvider is empty', async () => {
      await expect(service.refreshAddressScore(1, '')).rejects.toThrow(ValidationError);
    });

    it('should return scores', async () => {
      mockApi.scoreServiceRefreshAddressScore.mockResolvedValue({
        scores: [{ provider: 'chainalysis', score: '85', risk: 'LOW' }],
      } as never);

      const scores = await service.refreshAddressScore(1, 'chainalysis');
      expect(scores).toBeDefined();
    });
  });

  describe('refreshWhitelistedAddressScore', () => {
    it('should throw ValidationError when addressId is 0', async () => {
      await expect(
        service.refreshWhitelistedAddressScore(0, 'chainalysis')
      ).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when scoreProvider is empty', async () => {
      await expect(
        service.refreshWhitelistedAddressScore(1, '')
      ).rejects.toThrow(ValidationError);
    });
  });
});
