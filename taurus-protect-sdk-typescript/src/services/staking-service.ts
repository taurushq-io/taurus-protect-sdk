/**
 * Staking service for Taurus-PROTECT SDK.
 *
 * Provides methods for retrieving staking information across multiple blockchain networks.
 */

import { ValidationError } from '../errors';
import type { StakingApi } from '../internal/openapi/apis/StakingApi';
import { StakingServiceGetStakeAccountsAccountTypeEnum } from '../internal/openapi/apis/StakingApi';
import type {
  TgvalidatordGetADAStakePoolInfoReply,
  TgvalidatordETHValidatorInfo,
  TgvalidatordGetFTMValidatorInfoReply,
  TgvalidatordGetICPNeuronInfoReply,
  TgvalidatordGetNEARValidatorInfoReply,
  TgvalidatordStakeAccount,
  TgvalidatordGetXTZAddressStakingRewardsReply,
  TgvalidatordResponseCursor,
} from '../internal/openapi/models';
import type {
  ADAStakePoolInfo,
  ETHValidatorInfo,
  FTMValidatorInfo,
  ICPNeuronInfo,
  NEARValidatorInfo,
  StakeAccount,
  StakeAccountResult,
  StakeAccountType,
  StakeCursor,
  SolanaStakeAccount,
  SolanaStakeAccountState,
  ListStakeAccountsOptions,
  XTZStakingRewards,
  GetXTZStakingRewardsOptions,
} from '../models/staking';
import { BaseService } from './base';

/**
 * Service for retrieving staking information across multiple blockchain networks.
 *
 * This service provides access to validator information, stake accounts, and staking
 * rewards for various proof-of-stake blockchains including Cardano (ADA), Ethereum (ETH),
 * Fantom (FTM), Internet Computer (ICP), NEAR Protocol, Solana, and Tezos (XTZ).
 *
 * @example
 * ```typescript
 * // Get Cardano stake pool info
 * const poolInfo = await stakingService.getADAStakePoolInfo('mainnet', 'pool1abc123...');
 *
 * // Get Ethereum validator info
 * const validators = await stakingService.getETHValidatorsInfo('mainnet', ['validator1', 'validator2']);
 *
 * // Get stake accounts with pagination
 * const result = await stakingService.getStakeAccounts({ addressId: 'address-123' });
 * ```
 */
export class StakingService extends BaseService {
  private readonly stakingApi: StakingApi;

  /**
   * Creates a new StakingService instance.
   *
   * @param stakingApi - The StakingApi instance from the OpenAPI client
   */
  constructor(stakingApi: StakingApi) {
    super();
    this.stakingApi = stakingApi;
  }

  /**
   * Retrieves information about a Cardano stake pool.
   *
   * Returns details including the pool's pledge, margin, fixed costs, and active stake.
   *
   * @param network - The network (e.g., "mainnet", "preprod")
   * @param stakePoolId - The stake pool ID (Bech32 format, starting with "pool1")
   * @returns The stake pool information
   * @throws {@link ValidationError} If network or stakePoolId is empty
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const poolInfo = await stakingService.getADAStakePoolInfo('mainnet', 'pool1abc123...');
   * console.log(`Active stake: ${poolInfo.activeStake}`);
   * console.log(`Margin: ${poolInfo.margin}`);
   * ```
   */
  async getADAStakePoolInfo(network: string, stakePoolId: string): Promise<ADAStakePoolInfo> {
    if (!network || network.trim() === '') {
      throw new ValidationError('network is required');
    }
    if (!stakePoolId || stakePoolId.trim() === '') {
      throw new ValidationError('stakePoolId is required');
    }

    return this.execute(async () => {
      const reply = await this.stakingApi.stakingServiceGetADAStakePoolInfo({
        network,
        stakePoolId,
      });

      return this.mapADAStakePoolInfo(reply);
    });
  }

  /**
   * Retrieves information about Ethereum validators.
   *
   * Returns details including validator public keys, balances, and status.
   * Maximum 500 IDs can be provided per request.
   *
   * @param network - The network (e.g., "mainnet", "goerli")
   * @param ids - The list of validator IDs to query
   * @returns The list of validator information
   * @throws {@link ValidationError} If network is empty or ids is empty
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const validators = await stakingService.getETHValidatorsInfo('mainnet', ['id1', 'id2']);
   * for (const v of validators) {
   *   console.log(`Validator ${v.id}: ${v.status}, balance: ${v.balance}`);
   * }
   * ```
   */
  async getETHValidatorsInfo(network: string, ids: string[]): Promise<ETHValidatorInfo[]> {
    if (!network || network.trim() === '') {
      throw new ValidationError('network is required');
    }
    if (!ids || ids.length === 0) {
      throw new ValidationError('ids cannot be empty');
    }

    return this.execute(async () => {
      const reply = await this.stakingApi.stakingServiceGetETHValidatorsInfo({
        network,
        ids,
      });

      const validators = reply.validators ?? [];
      return validators.map((v) => this.mapETHValidatorInfo(v));
    });
  }

  /**
   * Retrieves information about a Fantom validator.
   *
   * Returns details including the validator's stake amounts and status.
   *
   * @param network - The network (e.g., "mainnet")
   * @param validatorAddress - The validator's address
   * @returns The validator information
   * @throws {@link ValidationError} If network or validatorAddress is empty
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const info = await stakingService.getFTMValidatorInfo('mainnet', '0x...');
   * console.log(`Total stake: ${info.totalStake}`);
   * console.log(`Is active: ${info.isActive}`);
   * ```
   */
  async getFTMValidatorInfo(network: string, validatorAddress: string): Promise<FTMValidatorInfo> {
    if (!network || network.trim() === '') {
      throw new ValidationError('network is required');
    }
    if (!validatorAddress || validatorAddress.trim() === '') {
      throw new ValidationError('validatorAddress is required');
    }

    return this.execute(async () => {
      const reply = await this.stakingApi.stakingServiceGetFTMValidatorInfo({
        network,
        validatorAddress,
      });

      return this.mapFTMValidatorInfo(reply);
    });
  }

  /**
   * Retrieves information about an Internet Computer Protocol neuron.
   *
   * Returns details including the neuron's stake, voting power, and dissolve delay.
   *
   * @param network - The network (e.g., "mainnet")
   * @param neuronId - The neuron ID
   * @returns The neuron information
   * @throws {@link ValidationError} If network or neuronId is empty
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const neuron = await stakingService.getICPNeuronInfo('mainnet', '123456789');
   * console.log(`Stake: ${neuron.stakeE8s}`);
   * console.log(`Voting power: ${neuron.votingPower}`);
   * ```
   */
  async getICPNeuronInfo(network: string, neuronId: string): Promise<ICPNeuronInfo> {
    if (!network || network.trim() === '') {
      throw new ValidationError('network is required');
    }
    if (!neuronId || neuronId.trim() === '') {
      throw new ValidationError('neuronId is required');
    }

    return this.execute(async () => {
      const reply = await this.stakingApi.stakingServiceGetICPNeuronInfo({
        network,
        neuronID: neuronId,
      });

      return this.mapICPNeuronInfo(reply);
    });
  }

  /**
   * Retrieves information about a NEAR Protocol validator.
   *
   * Returns details including the validator's total staked balance and fee structure.
   *
   * @param network - The network (e.g., "mainnet", "testnet")
   * @param validatorAddress - The validator's contract address
   * @returns The validator information
   * @throws {@link ValidationError} If network or validatorAddress is empty
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const validator = await stakingService.getNEARValidatorInfo('mainnet', 'validator.poolv1.near');
   * console.log(`Total staked: ${validator.totalStakedBalance}`);
   * console.log(`Fee: ${validator.rewardFeeFraction}`);
   * ```
   */
  async getNEARValidatorInfo(
    network: string,
    validatorAddress: string
  ): Promise<NEARValidatorInfo> {
    if (!network || network.trim() === '') {
      throw new ValidationError('network is required');
    }
    if (!validatorAddress || validatorAddress.trim() === '') {
      throw new ValidationError('validatorAddress is required');
    }

    return this.execute(async () => {
      const reply = await this.stakingApi.stakingServiceGetNEARValidatorInfo({
        network,
        validatorAddress,
      });

      return this.mapNEARValidatorInfo(reply);
    });
  }

  /**
   * Retrieves stake accounts with optional filtering.
   *
   * Returns a paginated list of stake accounts that can be filtered by address,
   * account type, or account address.
   *
   * @param options - Optional filtering and pagination options
   * @returns A paginated result containing stake accounts
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // Get all stake accounts for an address
   * const result = await stakingService.getStakeAccounts({
   *   addressId: 'address-123',
   * });
   *
   * for (const account of result.stakeAccounts) {
   *   console.log(`Account ${account.accountAddress}: ${account.accountType}`);
   * }
   *
   * // Paginate through results
   * if (result.cursor?.hasNext) {
   *   const nextPage = await stakingService.getStakeAccounts({
   *     cursorCurrentPage: result.cursor.currentPage,
   *     cursorPageRequest: 'NEXT',
   *   });
   * }
   * ```
   */
  async getStakeAccounts(options?: ListStakeAccountsOptions): Promise<StakeAccountResult> {
    return this.execute(async () => {
      const reply = await this.stakingApi.stakingServiceGetStakeAccounts({
        addressId: options?.addressId,
        accountType: options?.accountType as StakingServiceGetStakeAccountsAccountTypeEnum,
        accountAddress: options?.accountAddress,
        cursorCurrentPage: options?.cursorCurrentPage,
        cursorPageRequest: options?.cursorPageRequest,
        cursorPageSize: options?.cursorPageSize?.toString(),
      });

      const stakeAccounts = (reply.stakeAccounts ?? []).map((sa) => this.mapStakeAccount(sa));
      const cursor = this.mapCursor(reply.cursor);

      return {
        stakeAccounts,
        cursor,
      };
    });
  }

  /**
   * Retrieves Tezos staking rewards for an address over a time period.
   *
   * Returns the total rewards received by the specified address within the given
   * date range.
   *
   * @param options - The query options including network, addressId, and optional date range
   * @returns The staking rewards information
   * @throws {@link ValidationError} If network or addressId is empty
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // Get all-time rewards
   * const rewards = await stakingService.getXTZStakingRewards({
   *   network: 'mainnet',
   *   addressId: 'address-123',
   * });
   * console.log(`Total rewards: ${rewards.receivedRewardsAmount}`);
   *
   * // Get rewards for a specific period
   * const periodRewards = await stakingService.getXTZStakingRewards({
   *   network: 'mainnet',
   *   addressId: 'address-123',
   *   from: new Date('2024-01-01'),
   *   to: new Date('2024-12-31'),
   * });
   * ```
   */
  async getXTZStakingRewards(options: GetXTZStakingRewardsOptions): Promise<XTZStakingRewards> {
    if (!options.network || options.network.trim() === '') {
      throw new ValidationError('network is required');
    }
    if (!options.addressId || options.addressId.trim() === '') {
      throw new ValidationError('addressId is required');
    }

    return this.execute(async () => {
      const reply = await this.stakingApi.stakingServiceGetXTZAddressStakingRewards({
        network: options.network,
        addressID: options.addressId,
        from: options.from,
        to: options.to,
      });

      return this.mapXTZStakingRewards(reply);
    });
  }

  // Private mapping methods

  private mapADAStakePoolInfo(reply: TgvalidatordGetADAStakePoolInfoReply): ADAStakePoolInfo {
    return {
      pledge: reply.pledge,
      margin: reply.margin,
      fixedCost: reply.fixedCost,
      url: reply.url,
      activeStake: reply.activeStake,
      epoch: reply.epoch,
    };
  }

  private mapETHValidatorInfo(dto: TgvalidatordETHValidatorInfo): ETHValidatorInfo {
    return {
      id: dto.id,
      pubkey: dto.pubkey,
      status: dto.status,
      balance: dto.balance,
      network: dto.network,
      provider: dto.provider,
      addressId: dto.addressID,
    };
  }

  private mapFTMValidatorInfo(reply: TgvalidatordGetFTMValidatorInfoReply): FTMValidatorInfo {
    return {
      validatorId: reply.validatorID,
      address: reply.address,
      isActive: reply.isActive,
      totalStake: reply.totalStake,
      selfStake: reply.selfStake,
      deactivatedAtDateUnix: reply.deactivatedAtDateUnix,
      createdAtDateUnix: reply.createdAtDateUnix,
    };
  }

  private mapICPNeuronInfo(reply: TgvalidatordGetICPNeuronInfoReply): ICPNeuronInfo {
    return {
      neuronId: reply.neuronId,
      retrieveAtTimestampSeconds: reply.retrieveAtTimestampSeconds,
      neuronState: reply.neuronState as string | undefined,
      ageSeconds: reply.ageSeconds,
      dissolveDelaySeconds: reply.dissolveDelaySeconds,
      votingPower: reply.votingPower,
      createdTimestampSeconds: reply.createdTimestampSeconds,
      stakeE8s: reply.stakeE8S,
      joinedCommunityFundTimestampSeconds: reply.joinedCommunityFundTimestampSeconds,
      knownNeuronData: reply.knownNeuronData
        ? {
            name: reply.knownNeuronData.name,
            description: reply.knownNeuronData.description,
          }
        : undefined,
    };
  }

  private mapNEARValidatorInfo(reply: TgvalidatordGetNEARValidatorInfoReply): NEARValidatorInfo {
    return {
      validatorAddress: reply.validatorAddress,
      ownerId: reply.ownerId,
      totalStakedBalance: reply.totalStakedBalance,
      rewardFeeFraction: reply.rewardFeeFraction,
      stakingKey: reply.stakingKey,
      isStakingPaused: reply.isStakingPaused,
    };
  }

  private mapStakeAccount(dto: TgvalidatordStakeAccount): StakeAccount {
    return {
      id: dto.id,
      addressId: dto.addressId,
      accountAddress: dto.accountAddress,
      createdAt: dto.createdAt,
      updatedAt: dto.updatedAt,
      updatedAtBlock: dto.updatedAtBlock,
      accountType: dto.accountType as StakeAccountType | undefined,
      solanaStakeAccount: dto.solanaStakeAccount
        ? {
            derivationIndex: dto.solanaStakeAccount.derivationIndex,
            state: dto.solanaStakeAccount.state as SolanaStakeAccountState | undefined,
            validatorAddress: dto.solanaStakeAccount.validatorAddress,
            activeBalance: dto.solanaStakeAccount.activeBalance,
            inactiveBalance: dto.solanaStakeAccount.inactiveBalance,
            allowMerge: dto.solanaStakeAccount.allowMerge,
          }
        : undefined,
    };
  }

  private mapCursor(cursor?: TgvalidatordResponseCursor): StakeCursor | undefined {
    if (!cursor) {
      return undefined;
    }

    return {
      currentPage: cursor.currentPage,
      hasNext: cursor.hasNext,
      hasPrevious: cursor.hasPrevious,
    };
  }

  private mapXTZStakingRewards(reply: TgvalidatordGetXTZAddressStakingRewardsReply): XTZStakingRewards {
    return {
      receivedRewardsAmount: reply.receivedRewardsAmount,
    };
  }
}
