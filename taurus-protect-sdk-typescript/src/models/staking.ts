/**
 * Staking models for Taurus-PROTECT SDK.
 *
 * These models represent staking information across multiple blockchain networks.
 */

/**
 * Cardano (ADA) stake pool information.
 */
export interface ADAStakePoolInfo {
  /** The pool's pledge amount in lovelace */
  pledge?: string;
  /** The pool's margin (fee percentage) as a decimal */
  margin?: number;
  /** The pool's fixed cost per epoch in lovelace */
  fixedCost?: string;
  /** The pool's metadata URL */
  url?: string;
  /** The pool's active stake in lovelace */
  activeStake?: string;
  /** The current epoch number */
  epoch?: string;
}

/**
 * Ethereum validator information.
 */
export interface ETHValidatorInfo {
  /** The validator ID */
  id?: string;
  /** The validator's public key */
  pubkey?: string;
  /** The validator's status (e.g., "active_ongoing", "pending_queued") */
  status?: string;
  /** The validator's current balance in Gwei */
  balance?: string;
  /** The network (e.g., "mainnet", "goerli") */
  network?: string;
  /** The staking provider */
  provider?: string;
  /** The associated address ID */
  addressId?: string;
}

/**
 * Fantom (FTM) validator information.
 */
export interface FTMValidatorInfo {
  /** The validator ID */
  validatorId?: string;
  /** The validator's address */
  address?: string;
  /** Whether the validator is active */
  isActive?: boolean;
  /** The total stake amount */
  totalStake?: string;
  /** The validator's self stake */
  selfStake?: string;
  /** The deactivation timestamp (Unix) */
  deactivatedAtDateUnix?: string;
  /** The creation timestamp (Unix) */
  createdAtDateUnix?: string;
}

/**
 * Internet Computer Protocol (ICP) neuron information.
 */
export interface ICPNeuronInfo {
  /** The neuron ID */
  neuronId?: string;
  /** The timestamp when info was retrieved (seconds) */
  retrieveAtTimestampSeconds?: string;
  /** The neuron state */
  neuronState?: string;
  /** The neuron's age in seconds */
  ageSeconds?: string;
  /** The dissolve delay in seconds */
  dissolveDelaySeconds?: string;
  /** The voting power */
  votingPower?: string;
  /** The creation timestamp (seconds) */
  createdTimestampSeconds?: string;
  /** The stake amount in e8s */
  stakeE8s?: string;
  /** The timestamp when joined community fund (seconds) */
  joinedCommunityFundTimestampSeconds?: string;
  /** Known neuron data if applicable */
  knownNeuronData?: {
    name?: string;
    description?: string;
  };
}

/**
 * NEAR Protocol validator information.
 */
export interface NEARValidatorInfo {
  /** The validator's contract address */
  validatorAddress?: string;
  /** The owner ID */
  ownerId?: string;
  /** The total staked balance */
  totalStakedBalance?: string;
  /** The reward fee fraction */
  rewardFeeFraction?: number;
  /** The staking public key */
  stakingKey?: string;
  /** Whether staking is paused */
  isStakingPaused?: boolean;
}

/**
 * Stake account type.
 */
export type StakeAccountType = 'StakeAccountTypeSolana';

/**
 * Solana stake account state.
 */
export type SolanaStakeAccountState = 'inactive' | 'activating' | 'active' | 'deactivating';

/**
 * Solana stake account details.
 */
export interface SolanaStakeAccount {
  /** The derivation index used to generate this stake account */
  derivationIndex?: string;
  /** The stake account state */
  state?: SolanaStakeAccountState;
  /** The delegated validator's address */
  validatorAddress?: string;
  /** The active stake balance in lamports */
  activeBalance?: string;
  /** The inactive stake balance in lamports */
  inactiveBalance?: string;
  /** Whether this stake account can be merged with others */
  allowMerge?: boolean;
}

/**
 * Stake account information.
 */
export interface StakeAccount {
  /** The stake account ID */
  id?: string;
  /** The associated address ID */
  addressId?: string;
  /** The on-chain account address */
  accountAddress?: string;
  /** The creation timestamp */
  createdAt?: Date;
  /** The last update timestamp */
  updatedAt?: Date;
  /** The block at which the account was last updated */
  updatedAtBlock?: string;
  /** The account type */
  accountType?: StakeAccountType;
  /** Solana-specific stake account details */
  solanaStakeAccount?: SolanaStakeAccount;
}

/**
 * Cursor for pagination.
 */
export interface StakeCursor {
  /** The current page token */
  currentPage?: string;
  /** Whether there is a next page */
  hasNext?: boolean;
  /** Whether there is a previous page */
  hasPrevious?: boolean;
}

/**
 * Result of listing stake accounts.
 */
export interface StakeAccountResult {
  /** The list of stake accounts */
  stakeAccounts: StakeAccount[];
  /** Pagination cursor */
  cursor?: StakeCursor;
}

/**
 * Options for listing stake accounts.
 */
export interface ListStakeAccountsOptions {
  /** Filter by associated address ID */
  addressId?: string;
  /** Filter by account type */
  accountType?: StakeAccountType;
  /** Filter by on-chain account address */
  accountAddress?: string;
  /** Pagination cursor (current page token) */
  cursorCurrentPage?: string;
  /** Page request direction (FIRST, PREVIOUS, NEXT, LAST) */
  cursorPageRequest?: 'FIRST' | 'PREVIOUS' | 'NEXT' | 'LAST';
  /** Page size */
  cursorPageSize?: number;
}

/**
 * Tezos (XTZ) staking rewards information.
 */
export interface XTZStakingRewards {
  /** The total received rewards amount */
  receivedRewardsAmount?: string;
}

/**
 * Options for getting XTZ staking rewards.
 */
export interface GetXTZStakingRewardsOptions {
  /** The network (e.g., "mainnet") */
  network: string;
  /** The address ID in the Taurus-PROTECT system */
  addressId: string;
  /** The start date (optional) */
  from?: Date;
  /** The end date (optional) */
  to?: Date;
}
