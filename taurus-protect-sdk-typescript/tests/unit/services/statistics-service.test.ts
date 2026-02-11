/**
 * Unit tests for StatisticsService.
 */

import { StatisticsService } from '../../../src/services/statistics-service';
import type { StatisticsApi } from '../../../src/internal/openapi/apis/StatisticsApi';

function createMockApi(): jest.Mocked<StatisticsApi> {
  return {
    statisticsServiceGetPortfolioStatistics: jest.fn(),
  } as unknown as jest.Mocked<StatisticsApi>;
}

describe('StatisticsService', () => {
  let mockApi: jest.Mocked<StatisticsApi>;
  let service: StatisticsService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new StatisticsService(mockApi);
  });

  describe('getPortfolioStatistics', () => {
    it('should return portfolio statistics', async () => {
      mockApi.statisticsServiceGetPortfolioStatistics.mockResolvedValue({
        result: {
          totalBalance: '1000000',
          walletsCount: '10',
          addressesCount: '50',
        },
      } as never);

      const stats = await service.getPortfolioStatistics();
      expect(stats).toBeDefined();
    });

    it('should handle response with missing fields', async () => {
      mockApi.statisticsServiceGetPortfolioStatistics.mockResolvedValue({
        result: {},
      } as never);

      const stats = await service.getPortfolioStatistics();
      expect(stats).toBeDefined();
    });
  });
});
