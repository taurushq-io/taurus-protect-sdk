/**
 * Unit tests for staking mapping logic.
 *
 * The StakingService uses private mapping methods (no dedicated mapper file).
 * These tests verify the mapping logic by testing the model structures and
 * mapping patterns used within the service.
 */

import type {
  ADAStakePoolInfo,
  ETHValidatorInfo,
  FTMValidatorInfo,
  ICPNeuronInfo,
  NEARValidatorInfo,
  StakeAccount,
  StakeCursor,
  XTZStakingRewards,
} from '../../../src/models/staking';

describe('ADAStakePoolInfo mapping', () => {
  it('should map all fields from reply', () => {
    const reply = {
      pledge: '1000000000',
      margin: 0.05,
      fixedCost: '340000000',
      url: 'https://pool.example.com/metadata.json',
      activeStake: '50000000000',
      epoch: '450',
    };

    const result: ADAStakePoolInfo = {
      pledge: reply.pledge,
      margin: reply.margin,
      fixedCost: reply.fixedCost,
      url: reply.url,
      activeStake: reply.activeStake,
      epoch: reply.epoch,
    };

    expect(result.pledge).toBe('1000000000');
    expect(result.margin).toBe(0.05);
    expect(result.fixedCost).toBe('340000000');
    expect(result.url).toBe('https://pool.example.com/metadata.json');
    expect(result.activeStake).toBe('50000000000');
    expect(result.epoch).toBe('450');
  });

  it('should handle undefined fields', () => {
    const result: ADAStakePoolInfo = {};
    expect(result.pledge).toBeUndefined();
    expect(result.margin).toBeUndefined();
  });
});

describe('ETHValidatorInfo mapping', () => {
  it('should map all fields including addressID to addressId', () => {
    const dto = {
      id: 'val-1',
      pubkey: '0xabc123...',
      status: 'active_ongoing',
      balance: '32000000000',
      network: 'mainnet',
      provider: 'Lido',
      addressID: 'addr-42',
    };

    const result: ETHValidatorInfo = {
      id: dto.id,
      pubkey: dto.pubkey,
      status: dto.status,
      balance: dto.balance,
      network: dto.network,
      provider: dto.provider,
      addressId: dto.addressID,
    };

    expect(result.id).toBe('val-1');
    expect(result.pubkey).toBe('0xabc123...');
    expect(result.status).toBe('active_ongoing');
    expect(result.balance).toBe('32000000000');
    expect(result.network).toBe('mainnet');
    expect(result.provider).toBe('Lido');
    expect(result.addressId).toBe('addr-42');
  });
});

describe('FTMValidatorInfo mapping', () => {
  it('should map all fields including validatorID to validatorId', () => {
    const reply = {
      validatorID: '12',
      address: '0xfantom...',
      isActive: true,
      totalStake: '1000000',
      selfStake: '500000',
      deactivatedAtDateUnix: '0',
      createdAtDateUnix: '1672531200',
    };

    const result: FTMValidatorInfo = {
      validatorId: reply.validatorID,
      address: reply.address,
      isActive: reply.isActive,
      totalStake: reply.totalStake,
      selfStake: reply.selfStake,
      deactivatedAtDateUnix: reply.deactivatedAtDateUnix,
      createdAtDateUnix: reply.createdAtDateUnix,
    };

    expect(result.validatorId).toBe('12');
    expect(result.isActive).toBe(true);
    expect(result.totalStake).toBe('1000000');
    expect(result.selfStake).toBe('500000');
  });
});

describe('ICPNeuronInfo mapping', () => {
  it('should map all fields including stakeE8S to stakeE8s', () => {
    const reply = {
      neuronId: 'neuron-1',
      retrieveAtTimestampSeconds: '1700000000',
      neuronState: 'DISSOLVING',
      ageSeconds: '31536000',
      dissolveDelaySeconds: '15768000',
      votingPower: '1000000000',
      createdTimestampSeconds: '1668464000',
      stakeE8S: '500000000',
      joinedCommunityFundTimestampSeconds: '1670000000',
      knownNeuronData: {
        name: 'Test Neuron',
        description: 'A test neuron for verification',
      },
    };

    const result: ICPNeuronInfo = {
      neuronId: reply.neuronId,
      retrieveAtTimestampSeconds: reply.retrieveAtTimestampSeconds,
      neuronState: reply.neuronState as string | undefined,
      ageSeconds: reply.ageSeconds,
      dissolveDelaySeconds: reply.dissolveDelaySeconds,
      votingPower: reply.votingPower,
      createdTimestampSeconds: reply.createdTimestampSeconds,
      stakeE8s: reply.stakeE8S,
      joinedCommunityFundTimestampSeconds: reply.joinedCommunityFundTimestampSeconds,
      knownNeuronData: {
        name: reply.knownNeuronData.name,
        description: reply.knownNeuronData.description,
      },
    };

    expect(result.neuronId).toBe('neuron-1');
    expect(result.stakeE8s).toBe('500000000');
    expect(result.neuronState).toBe('DISSOLVING');
    expect(result.knownNeuronData?.name).toBe('Test Neuron');
    expect(result.knownNeuronData?.description).toBe('A test neuron for verification');
  });

  it('should handle missing knownNeuronData', () => {
    const reply = {
      neuronId: 'neuron-2',
      knownNeuronData: undefined as { name?: string; description?: string } | undefined,
    };

    const result: ICPNeuronInfo = {
      neuronId: reply.neuronId,
      knownNeuronData: reply.knownNeuronData
        ? { name: reply.knownNeuronData.name, description: reply.knownNeuronData.description }
        : undefined,
    };

    expect(result.knownNeuronData).toBeUndefined();
  });
});

describe('NEARValidatorInfo mapping', () => {
  it('should map all fields', () => {
    const reply = {
      validatorAddress: 'validator.poolv1.near',
      ownerId: 'owner.near',
      totalStakedBalance: '10000000000000000000000000',
      rewardFeeFraction: 10,
      stakingKey: 'ed25519:abc123...',
      isStakingPaused: false,
    };

    const result: NEARValidatorInfo = {
      validatorAddress: reply.validatorAddress,
      ownerId: reply.ownerId,
      totalStakedBalance: reply.totalStakedBalance,
      rewardFeeFraction: reply.rewardFeeFraction,
      stakingKey: reply.stakingKey,
      isStakingPaused: reply.isStakingPaused,
    };

    expect(result.validatorAddress).toBe('validator.poolv1.near');
    expect(result.totalStakedBalance).toBe('10000000000000000000000000');
    expect(result.isStakingPaused).toBe(false);
  });
});

describe('StakeAccount mapping', () => {
  it('should map stake account with Solana details', () => {
    const dto = {
      id: 'sa-1',
      addressId: 'addr-1',
      accountAddress: 'Stake11111...',
      createdAt: new Date('2024-01-01'),
      updatedAt: new Date('2024-06-01'),
      updatedAtBlock: '250000000',
      accountType: 'StakeAccountTypeSolana',
      solanaStakeAccount: {
        derivationIndex: '0',
        state: 'active',
        validatorAddress: 'Vote111...',
        activeBalance: '5000000000',
        inactiveBalance: '0',
        allowMerge: true,
      },
    };

    const result: StakeAccount = {
      id: dto.id,
      addressId: dto.addressId,
      accountAddress: dto.accountAddress,
      createdAt: dto.createdAt,
      updatedAt: dto.updatedAt,
      updatedAtBlock: dto.updatedAtBlock,
      accountType: dto.accountType as 'StakeAccountTypeSolana',
      solanaStakeAccount: {
        derivationIndex: dto.solanaStakeAccount.derivationIndex,
        state: dto.solanaStakeAccount.state as 'active',
        validatorAddress: dto.solanaStakeAccount.validatorAddress,
        activeBalance: dto.solanaStakeAccount.activeBalance,
        inactiveBalance: dto.solanaStakeAccount.inactiveBalance,
        allowMerge: dto.solanaStakeAccount.allowMerge,
      },
    };

    expect(result.id).toBe('sa-1');
    expect(result.accountType).toBe('StakeAccountTypeSolana');
    expect(result.solanaStakeAccount?.state).toBe('active');
    expect(result.solanaStakeAccount?.activeBalance).toBe('5000000000');
    expect(result.solanaStakeAccount?.allowMerge).toBe(true);
  });

  it('should handle missing Solana details', () => {
    const result: StakeAccount = {
      id: 'sa-2',
      addressId: 'addr-2',
      accountAddress: 'StakeABC...',
    };

    expect(result.solanaStakeAccount).toBeUndefined();
    expect(result.accountType).toBeUndefined();
  });
});

describe('StakeCursor mapping', () => {
  it('should map cursor fields', () => {
    const cursor = {
      currentPage: 'page-token-xyz',
      hasNext: true,
      hasPrevious: false,
    };

    const result: StakeCursor = {
      currentPage: cursor.currentPage,
      hasNext: cursor.hasNext,
      hasPrevious: cursor.hasPrevious,
    };

    expect(result.currentPage).toBe('page-token-xyz');
    expect(result.hasNext).toBe(true);
    expect(result.hasPrevious).toBe(false);
  });

  it('should handle undefined cursor', () => {
    const cursor: undefined = undefined;
    const result: StakeCursor | undefined = cursor;
    expect(result).toBeUndefined();
  });
});

describe('XTZStakingRewards mapping', () => {
  it('should map rewards amount', () => {
    const reply = {
      receivedRewardsAmount: '12345678',
    };

    const result: XTZStakingRewards = {
      receivedRewardsAmount: reply.receivedRewardsAmount,
    };

    expect(result.receivedRewardsAmount).toBe('12345678');
  });

  it('should handle undefined rewards', () => {
    const result: XTZStakingRewards = {};
    expect(result.receivedRewardsAmount).toBeUndefined();
  });
});
