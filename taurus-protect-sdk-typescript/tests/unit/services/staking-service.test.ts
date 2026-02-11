/**
 * Unit tests for StakingService.
 */

import { StakingService } from '../../../src/services/staking-service';
import { ValidationError } from '../../../src/errors';
import type { StakingApi } from '../../../src/internal/openapi/apis/StakingApi';

function createMockApi(): jest.Mocked<StakingApi> {
  return {
    stakingServiceGetADAStakePoolInfo: jest.fn(),
    stakingServiceGetETHValidatorsInfo: jest.fn(),
    stakingServiceGetFTMValidatorInfo: jest.fn(),
    stakingServiceGetICPNeuronInfo: jest.fn(),
    stakingServiceGetNEARValidatorInfo: jest.fn(),
    stakingServiceGetStakeAccounts: jest.fn(),
    stakingServiceGetXTZStakingRewards: jest.fn(),
  } as unknown as jest.Mocked<StakingApi>;
}

describe('StakingService', () => {
  let mockApi: jest.Mocked<StakingApi>;
  let service: StakingService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new StakingService(mockApi);
  });

  describe('getADAStakePoolInfo', () => {
    it('should throw ValidationError when network is empty', async () => {
      await expect(service.getADAStakePoolInfo('', 'pool-123')).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when stakePoolId is empty', async () => {
      await expect(service.getADAStakePoolInfo('mainnet', '')).rejects.toThrow(ValidationError);
    });

    it('should return stake pool info', async () => {
      mockApi.stakingServiceGetADAStakePoolInfo.mockResolvedValue({
        result: { poolId: 'pool-123', margin: '0.02' },
      } as never);

      const info = await service.getADAStakePoolInfo('mainnet', 'pool-123');
      expect(info).toBeDefined();
    });
  });

  describe('getETHValidatorsInfo', () => {
    it('should throw ValidationError when network is empty', async () => {
      await expect(service.getETHValidatorsInfo('', ['val-1'])).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when ids is empty', async () => {
      await expect(service.getETHValidatorsInfo('mainnet', [])).rejects.toThrow(ValidationError);
    });

    it('should return validators info', async () => {
      mockApi.stakingServiceGetETHValidatorsInfo.mockResolvedValue({
        validators: [{ validatorId: 'val-1', status: 'ACTIVE' }],
      } as never);

      const info = await service.getETHValidatorsInfo('mainnet', ['val-1']);
      expect(info).toHaveLength(1);
    });
  });

  describe('getICPNeuronInfo', () => {
    it('should throw ValidationError when network is empty', async () => {
      await expect(service.getICPNeuronInfo('', 'neuron-1')).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when neuronId is empty', async () => {
      await expect(service.getICPNeuronInfo('mainnet', '')).rejects.toThrow(ValidationError);
    });
  });
});
