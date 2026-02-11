/**
 * Pledge models for Taurus Network.
 *
 * This module provides domain models for pledges, which represent
 * funds that are reserved and pledged from one participant (pledgor)
 * to another (pledgee) on a shared address.
 */

/**
 * Pledge status enum.
 */
export enum PledgeStatus {
  PENDING = "PENDING",
  CONFIRMED = "CONFIRMED",
  REJECTED = "REJECTED",
  CANCELED = "CANCELED",
  UNPLEDGED = "UNPLEDGED",
}

/**
 * Pledge type for withdrawal rights.
 */
export enum PledgeType {
  NO_WITHDRAWALS_RIGHTS = "NO_WITHDRAWALS_RIGHTS",
  PLEDGEE_WITHDRAWALS_RIGHTS = "PLEDGEE_WITHDRAWALS_RIGHTS",
  PLEDGEE_AUTO_APPROVED_WITHDRAWALS_RIGHTS = "PLEDGEE_AUTO_APPROVED_WITHDRAWALS_RIGHTS",
}

/**
 * Pledge action status enum.
 */
export enum PledgeActionStatus {
  PENDING = "PENDING",
  APPROVED = "APPROVED",
  REJECTED = "REJECTED",
  EXECUTED = "EXECUTED",
  CANCELED = "CANCELED",
}

/**
 * Pledge action type enum.
 */
export enum PledgeActionType {
  CREATE_PLEDGE = "CREATE_PLEDGE",
  ADD_COLLATERAL = "ADD_COLLATERAL",
  WITHDRAW = "WITHDRAW",
  INITIATE_WITHDRAW = "INITIATE_WITHDRAW",
  UNPLEDGE = "UNPLEDGE",
  REJECT = "REJECT",
}

/**
 * Pledge withdrawal status enum.
 */
export enum PledgeWithdrawalStatus {
  PENDING = "PENDING",
  PENDING_APPROVAL = "PENDING_APPROVAL",
  APPROVED = "APPROVED",
  REJECTED = "REJECTED",
  EXECUTED = "EXECUTED",
  CANCELED = "CANCELED",
}

/**
 * Custom attribute on a pledge.
 */
export interface PledgeAttribute {
  /** Attribute key. */
  readonly key: string;
  /** Attribute value. */
  readonly value: string;
}

/**
 * Duration configuration for a pledge.
 */
export interface PledgeDurationSetup {
  /** Start date of the pledge. */
  readonly startDate: Date | undefined;
  /** End date of the pledge. */
  readonly endDate: Date | undefined;
}

/**
 * Audit trail entry for a pledge.
 */
export interface PledgeTrail {
  /** Trail entry ID. */
  readonly id: string;
  /** Action performed. */
  readonly action: string;
  /** Who performed the action. */
  readonly actor: string;
  /** When it occurred. */
  readonly timestamp: Date | undefined;
  /** Optional comment. */
  readonly comment: string | undefined;
}

/**
 * Taurus Network pledge.
 *
 * A pledge represents funds that are reserved and pledged from one
 * participant (pledgor) to another (pledgee) on a shared address.
 */
export interface Pledge {
  /** Unique pledge identifier. */
  readonly id: string;
  /** ID of the shared address holding pledged funds. */
  readonly sharedAddressId: string;
  /** ID of the pledgor (owner) participant. */
  readonly ownerParticipantId: string;
  /** ID of the pledgee (target) participant. */
  readonly targetParticipantId: string;
  /** Currency identifier. */
  readonly currencyId: string;
  /** Blockchain name. */
  readonly blockchain: string;
  /** Network name. */
  readonly network: string;
  /** Additional argument 1. */
  readonly arg1: string | undefined;
  /** Additional argument 2. */
  readonly arg2: string | undefined;
  /** Pledged amount in smallest currency unit. */
  readonly amount: string;
  /** Current pledge status. */
  readonly status: string;
  /** Pledge type for withdrawal rights. */
  readonly pledgeType: string;
  /** Direction of the pledge (incoming/outgoing). */
  readonly direction: string;
  /** External reference for reconciliation. */
  readonly externalReferenceId: string | undefined;
  /** Internal reconciliation note. */
  readonly reconciliationNote: string | undefined;
  /** Whitelisted address ID if pledgee. */
  readonly wlAddressId: string | undefined;
  /** Original creation timestamp. */
  readonly originCreationDate: Date | undefined;
  /** When pledge was unpledged (if applicable). */
  readonly unpledgeDate: Date | undefined;
  /** Duration configuration. */
  readonly durationSetup: PledgeDurationSetup | undefined;
  /** Custom key-value attributes. */
  readonly attributes: PledgeAttribute[];
  /** Audit trail entries. */
  readonly trails: PledgeTrail[];
  /** When the pledge was created. */
  readonly createdAt: Date | undefined;
  /** When the pledge was last updated. */
  readonly updatedAt: Date | undefined;
}

/**
 * Metadata for a pledge action.
 */
export interface PledgeActionMetadata {
  /** Hash of the action metadata. */
  readonly hash: string;
  /** Payload data as string. */
  readonly payload: string | undefined;
}

/**
 * Audit trail entry for a pledge action.
 */
export interface PledgeActionTrail {
  /** Trail entry ID. */
  readonly id: string;
  /** Action performed. */
  readonly action: string;
  /** Who performed the action. */
  readonly actor: string;
  /** When it occurred. */
  readonly timestamp: Date | undefined;
  /** Optional comment. */
  readonly comment: string | undefined;
}

/**
 * Pledge action requiring approval.
 *
 * Pledge actions represent pending operations on pledges that require
 * approval before execution.
 */
export interface PledgeAction {
  /** Unique action identifier. */
  readonly id: string;
  /** Associated pledge ID. */
  readonly pledgeId: string;
  /** Type of action. */
  readonly actionType: string;
  /** Current action status. */
  readonly status: string;
  /** Action metadata including hash for signing. */
  readonly metadata: PledgeActionMetadata | undefined;
  /** Rule applied to this action. */
  readonly rule: string | undefined;
  /** List of approvers needed. */
  readonly needsApprovalFrom: string[];
  /** Associated withdrawal ID if applicable. */
  readonly pledgeWithdrawalId: string | undefined;
  /** Envelope data. */
  readonly envelope: string | undefined;
  /** Audit trail. */
  readonly trails: PledgeActionTrail[];
  /** Creation timestamp. */
  readonly createdAt: Date | undefined;
  /** Last approval timestamp. */
  readonly lastApprovalDate: Date | undefined;
}

/**
 * Audit trail entry for a pledge withdrawal.
 */
export interface PledgeWithdrawalTrail {
  /** Trail entry ID. */
  readonly id: string;
  /** Action performed. */
  readonly action: string;
  /** Who performed the action. */
  readonly actor: string;
  /** When it occurred. */
  readonly timestamp: Date | undefined;
  /** Optional comment. */
  readonly comment: string | undefined;
}

/**
 * Pledge withdrawal record.
 *
 * Represents a withdrawal request or execution from a pledge.
 */
export interface PledgeWithdrawal {
  /** Unique withdrawal identifier. */
  readonly id: string;
  /** Associated pledge ID. */
  readonly pledgeId: string;
  /** Destination shared address ID. */
  readonly destinationSharedAddressId: string;
  /** Withdrawal amount. */
  readonly amount: string;
  /** Withdrawal status. */
  readonly status: string;
  /** Blockchain transaction hash. */
  readonly txHash: string | undefined;
  /** Internal transaction ID. */
  readonly txId: string | undefined;
  /** Associated request ID. */
  readonly requestId: string | undefined;
  /** Block number of transaction. */
  readonly txBlockNumber: string | undefined;
  /** Who initiated the withdrawal. */
  readonly initiatorParticipantId: string | undefined;
  /** External reference. */
  readonly externalReferenceId: string | undefined;
  /** Audit trail. */
  readonly trails: PledgeWithdrawalTrail[];
  /** Creation timestamp. */
  readonly createdAt: Date | undefined;
}

// Request models

/**
 * Request to create a new pledge.
 */
export interface CreatePledgeRequest {
  /** Shared address ID for pledge. */
  readonly sharedAddressId: string;
  /** Currency identifier. */
  readonly currencyId: string;
  /** Amount in smallest currency unit. */
  readonly amount: string;
  /** Withdrawal rights type. */
  readonly pledgeType?: string;
  /** Duration configuration. */
  readonly pledgeDurationSetup?: PledgeDurationSetup;
  /** Custom attributes. */
  readonly keyValueAttributes?: PledgeAttribute[];
  /** External reference. */
  readonly externalReferenceId?: string;
  /** Reconciliation note. */
  readonly reconciliationNote?: string;
}

/**
 * Request to update a pledge's default destination.
 */
export interface UpdatePledgeRequest {
  /** Default destination shared address ID. */
  readonly defaultDestinationSharedAddressId?: string;
  /** Default destination internal address ID. */
  readonly defaultDestinationInternalAddressId?: string;
}

/**
 * Request to add collateral to an existing pledge.
 */
export interface AddPledgeCollateralRequest {
  /** Amount to add in smallest currency unit. */
  readonly amount: string;
}

/**
 * Request for pledgee withdrawal from pledge.
 *
 * Use either destinationSharedAddressId or destinationInternalAddressId.
 */
export interface WithdrawPledgeRequest {
  /** Amount to withdraw. */
  readonly amount: string;
  /** Destination shared address ID. */
  readonly destinationSharedAddressId?: string;
  /** Destination internal address ID. */
  readonly destinationInternalAddressId?: string;
  /** External reference ID. */
  readonly externalReferenceId?: string;
}

/**
 * Request for pledgor-initiated withdrawal from pledge.
 */
export interface InitiateWithdrawPledgeRequest {
  /** Amount to withdraw. */
  readonly amount: string;
  /** Destination shared address ID. */
  readonly destinationSharedAddressId?: string;
}

/**
 * Request to reject a pledge.
 */
export interface RejectPledgeRequest {
  /** Rejection comment. */
  readonly comment: string;
}

/**
 * Request to approve pledge actions.
 */
export interface ApprovePledgeActionsRequest {
  /** IDs of pledge actions to approve. */
  readonly ids: string[];
  /** Approval comment. */
  readonly comment?: string;
  /** ECDSA signature over action hashes. */
  readonly signature: string;
}

/**
 * Request to reject pledge actions.
 */
export interface RejectPledgeActionsRequest {
  /** IDs of pledge actions to reject. */
  readonly ids: string[];
  /** Rejection comment. */
  readonly comment: string;
}

// Filter options

/**
 * Options for listing pledges.
 */
export interface ListPledgesOptions {
  /** Maximum items to return (default: 50, max: 1000). */
  readonly limit?: number;
  /** Number of items to skip. */
  readonly offset?: number;
  /** Filter by statuses. */
  readonly statuses?: string[];
  /** Filter by currency. */
  readonly currencyId?: string;
  /** Filter by direction (incoming/outgoing). */
  readonly direction?: string;
  /** Filter by participant ID. */
  readonly participantId?: string;
}

/**
 * Options for listing pledge actions.
 */
export interface ListPledgeActionsOptions {
  /** Maximum items to return (default: 50, max: 1000). */
  readonly limit?: number;
  /** Number of items to skip. */
  readonly offset?: number;
  /** Filter by statuses. */
  readonly statuses?: string[];
  /** Filter by action types. */
  readonly actionTypes?: string[];
  /** Filter by pledge ID. */
  readonly pledgeId?: string;
}

/**
 * Options for listing pledge withdrawals.
 */
export interface ListPledgeWithdrawalsOptions {
  /** Maximum items to return (default: 50, max: 1000). */
  readonly limit?: number;
  /** Number of items to skip. */
  readonly offset?: number;
  /** Filter by statuses. */
  readonly statuses?: string[];
  /** Filter by pledge ID. */
  readonly pledgeId?: string;
}
