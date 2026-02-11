/**
 * Unit tests for JobService.
 */

import { JobService } from '../../../src/services/job-service';
import { ValidationError, NotFoundError } from '../../../src/errors';
import type { JobsApi } from '../../../src/internal/openapi/apis/JobsApi';

function createMockApi(): jest.Mocked<JobsApi> {
  return {
    jobServiceGetJobs: jest.fn(),
    jobServiceGetJob: jest.fn(),
    jobServiceGetJobStatus: jest.fn(),
  } as unknown as jest.Mocked<JobsApi>;
}

describe('JobService', () => {
  let mockApi: jest.Mocked<JobsApi>;
  let service: JobService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new JobService(mockApi);
  });

  describe('list', () => {
    it('should return jobs', async () => {
      mockApi.jobServiceGetJobs.mockResolvedValue({
        result: [
          { name: 'balance-sync', status: 'RUNNING' },
          { name: 'tx-monitor', status: 'IDLE' },
        ],
      } as never);

      const jobs = await service.list();
      expect(jobs).toHaveLength(2);
    });

    it('should handle empty results', async () => {
      mockApi.jobServiceGetJobs.mockResolvedValue({
        result: [],
      } as never);

      const jobs = await service.list();
      expect(jobs).toHaveLength(0);
    });
  });

  describe('get', () => {
    it('should throw ValidationError when name is empty', async () => {
      await expect(service.get('')).rejects.toThrow(ValidationError);
      await expect(service.get('')).rejects.toThrow('name is required');
    });

    it('should return job when found', async () => {
      mockApi.jobServiceGetJob.mockResolvedValue({
        result: { name: 'balance-sync', status: 'RUNNING' },
      } as never);

      const job = await service.get('balance-sync');
      expect(job).toBeDefined();
    });
  });

  describe('getStatus', () => {
    it('should throw ValidationError when name is empty', async () => {
      await expect(service.getStatus('', 'id-1')).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when id is empty', async () => {
      await expect(service.getStatus('balance-sync', '')).rejects.toThrow(ValidationError);
    });

    it('should return job status', async () => {
      mockApi.jobServiceGetJobStatus.mockResolvedValue({
        result: { status: 'COMPLETED', progress: '100' },
      } as never);

      const status = await service.getStatus('balance-sync', 'run-1');
      expect(status).toBeDefined();
    });
  });
});
