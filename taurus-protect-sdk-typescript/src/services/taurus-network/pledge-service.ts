/**
 * Pledge service for Taurus Network in Taurus-PROTECT SDK.
 *
 * Provides methods for managing Taurus Network pledges between participants.
 */

import { NotFoundError, ValidationError } from '../../errors';
import type { TaurusNetworkPledgeApi } from '../../internal/openapi/apis/TaurusNetworkPledgeApi';
import type {
  TgvalidatordTnPledge,
  TgvalidatordTnPledgeAction,
  TgvalidatordTnPledgeWithdrawal,
} from '../../internal/openapi/models/index';
import { BaseService } from '../base';

/**
 * Cursor-based pagination information.
 */
export interface CursorPagination {
  currentPage?: string;
  hasNext: boolean;
  hasPrevious: boolean;
}

/**
 * Pledge trail entry.
 */
export interface PledgeTrail {
  id?: string;
  pledgeId?: string;
  participantId?: string;
  pledgeAmount?: string;
  action?: string;
  comment?: string;
  createdAt?: Date;
}

/**
 * Pledge duration setup.
 */
export interface PledgeDurationSetup {
  minimumDuration?: string;
  endOfMinimumDurationDate?: Date;
  noticePeriodDuration?: string;
  endOfNoticePeriodDate?: Date;
}

/**
 * Pledge attribute.
 */
export interface PledgeAttribute {
  key: string;
  value: string;
}

/**
 * A Taurus Network pledge.
 */
export interface Pledge {
  id: string;
  sharedAddressId?: string;
  ownerParticipantId?: string;
  targetParticipantId?: string;
  currencyId?: string;
  blockchain?: string;
  network?: string;
  amount?: string;
  status?: string;
  pledgeType?: string;
  direction?: string;
  externalReferenceId?: string;
  reconciliationNote?: string;
  durationSetup?: PledgeDurationSetup;
  attributes?: PledgeAttribute[];
  trails?: PledgeTrail[];
  originCreationDate?: Date;
  unpledgeDate?: Date;
  createdAt?: Date;
  updatedAt?: Date;
}

/**
 * Pledge action metadata.
 */
export interface PledgeActionMetadata {
  hash?: string;
  payloadAsString?: string;
}

/**
 * Pledge action trail entry.
 */
export interface PledgeActionTrail {
  id?: string;
  pledgeActionId?: string;
  userId?: string;
  externalUserId?: string;
  action?: string;
  comment?: string;
  createdAt?: Date;
}

/**
 * A Taurus Network pledge action.
 */
export interface PledgeAction {
  id: string;
  pledgeId?: string;
  actionType?: string;
  status?: string;
  metadata?: PledgeActionMetadata;
  rule?: string;
  needsApprovalFrom?: string[];
  pledgeWithdrawalId?: string;
  envelope?: string;
  trails?: PledgeActionTrail[];
  createdAt?: Date;
  lastApprovalDate?: Date;
}

/**
 * Pledge withdrawal trail entry.
 */
export interface PledgeWithdrawalTrail {
  id?: string;
  pledgeWithdrawalId?: string;
  addressCommandId?: string;
  participantId?: string;
  action?: string;
  comment?: string;
  createdAt?: Date;
}

/**
 * A Taurus Network pledge withdrawal.
 */
export interface PledgeWithdrawal {
  id: string;
  pledgeId?: string;
  destinationSharedAddressId?: string;
  amount?: string;
  status?: string;
  txHash?: string;
  txId?: string;
  requestId?: string;
  txBlockNumber?: string;
  initiatorParticipantId?: string;
  externalReferenceId?: string;
  trails?: PledgeWithdrawalTrail[];
  createdAt?: Date;
}

/**
 * Options for listing pledges.
 */
export interface ListPledgesOptions {
  ownerParticipantId?: string;
  targetParticipantId?: string;
  sharedAddressIds?: string[];
  currencyId?: string;
  statuses?: string[];
  sortOrder?: string;
  pageSize?: number;
  currentPage?: string;
  pageRequest?: string;
}

/**
 * Options for listing pledge actions.
 */
export interface ListPledgeActionsOptions {
  ids?: string[];
  pledgeId?: string;
  types?: string[];
  sortOrder?: string;
  pageSize?: number;
  currentPage?: string;
  pageRequest?: string;
}

/**
 * Options for listing pledge withdrawals.
 */
export interface ListPledgeWithdrawalsOptions {
  pledgeId?: string;
  withdrawalStatus?: string;
  sortOrder?: string;
  pageSize?: number;
  currentPage?: string;
  pageRequest?: string;
}

/**
 * Request to create a pledge.
 */
export interface CreatePledgeRequest {
  sharedAddressId: string;
  currencyId: string;
  amount: string;
  pledgeType?: string;
  externalReferenceId?: string;
  reconciliationNote?: string;
  pledgeDurationSetup?: {
    minimumDuration?: string;
    endOfMinimumDurationDate?: Date;
    noticePeriodDuration?: string;
  };
  keyValueAttributes?: PledgeAttribute[];
}

/**
 * Request to update a pledge.
 */
export interface UpdatePledgeRequest {
  defaultDestinationSharedAddressId?: string;
  defaultDestinationInternalAddressId?: string;
}

/**
 * Request to add collateral to a pledge.
 */
export interface AddPledgeCollateralRequest {
  amount: string;
}

/**
 * Request to withdraw from a pledge.
 */
export interface WithdrawPledgeRequest {
  amount: string;
  destinationSharedAddressId?: string;
  destinationInternalAddressId?: string;
  externalReferenceId?: string;
}

/**
 * Request to initiate a withdrawal from a pledge.
 */
export interface InitiateWithdrawPledgeRequest {
  amount: string;
  destinationSharedAddressId?: string;
}

/**
 * Request to reject a pledge.
 */
export interface RejectPledgeRequest {
  comment: string;
}

/**
 * Request to reject pledge actions.
 */
export interface RejectPledgeActionsRequest {
  ids: string[];
  comment: string;
}

/**
 * Result of creating a pledge.
 */
export interface CreatePledgeResult {
  pledge: Pledge;
  pledgeActionId?: string;
}

/**
 * Result of adding collateral to a pledge.
 */
export interface AddCollateralResult {
  pledgeActionId?: string;
}

/**
 * Result of withdrawing from a pledge.
 */
export interface WithdrawPledgeResult {
  pledgeWithdrawalId?: string;
  pledgeActionId?: string;
}

/**
 * Result of unpledging.
 */
export interface UnpledgeResult {
  pledgeActionId?: string;
}

/**
 * Maps a DTO to a Pledge.
 */
function pledgeFromDto(dto?: TgvalidatordTnPledge): Pledge | undefined {
  if (!dto) {
    return undefined;
  }

  return {
    id: dto.id ?? '',
    sharedAddressId: dto.sharedAddressID,
    ownerParticipantId: dto.ownerParticipantID,
    targetParticipantId: dto.targetParticipantID,
    currencyId: dto.currencyID,
    blockchain: dto.blockchain,
    network: dto.network,
    amount: dto.amount,
    status: dto.status,
    pledgeType: dto.pledgeType,
    direction: dto.direction,
    externalReferenceId: dto.externalReferenceId,
    reconciliationNote: dto.reconciliationNote,
    durationSetup: dto.durationSetup
      ? {
          minimumDuration: dto.durationSetup.minimumDuration,
          endOfMinimumDurationDate: dto.durationSetup.endOfMinimumDurationDate,
          noticePeriodDuration: dto.durationSetup.noticePeriodDuration,
          endOfNoticePeriodDate: dto.durationSetup.endOfNoticePeriodDate,
        }
      : undefined,
    attributes: dto.attributes?.map((attr) => ({
      key: attr.key ?? '',
      value: attr.value ?? '',
    })),
    trails: dto.trails?.map((trail) => ({
      id: trail.id,
      pledgeId: trail.pledgeID,
      participantId: trail.participantID,
      pledgeAmount: trail.pledgeAmount,
      action: trail.action,
      comment: trail.comment,
      createdAt: trail.createdAt,
    })),
    originCreationDate: dto.originCreationDate,
    unpledgeDate: dto.unpledgeDate,
    createdAt: dto.createdAt,
    updatedAt: dto.updatedAt,
  };
}

/**
 * Maps a DTO to a PledgeAction.
 */
function pledgeActionFromDto(dto?: TgvalidatordTnPledgeAction): PledgeAction | undefined {
  if (!dto) {
    return undefined;
  }

  return {
    id: dto.id ?? '',
    pledgeId: dto.pledgeID,
    actionType: dto.actionType,
    status: dto.status,
    metadata: dto.metadata
      ? {
          hash: dto.metadata.hash,
          payloadAsString: dto.metadata.payloadAsString,
        }
      : undefined,
    rule: dto.rule,
    needsApprovalFrom: dto.needsApprovalFrom,
    pledgeWithdrawalId: dto.pledgeWithdrawalID,
    envelope: dto.envelope,
    trails: dto.trails?.map((trail) => ({
      id: trail.id,
      pledgeActionId: trail.pledgeActionID,
      userId: trail.userID,
      externalUserId: trail.externalUserID,
      action: trail.action,
      comment: trail.comment,
      createdAt: trail.createdAt,
    })),
    createdAt: dto.createdAt,
    lastApprovalDate: dto.lastApprovalDate,
  };
}

/**
 * Maps a DTO to a PledgeWithdrawal.
 */
function pledgeWithdrawalFromDto(dto?: TgvalidatordTnPledgeWithdrawal): PledgeWithdrawal | undefined {
  if (!dto) {
    return undefined;
  }

  return {
    id: dto.id ?? '',
    pledgeId: dto.pledgeID,
    destinationSharedAddressId: dto.destinationSharedAddressID,
    amount: dto.amount,
    status: dto.status,
    txHash: dto.txHash,
    txId: dto.txID,
    requestId: dto.requestID,
    txBlockNumber: dto.txBlockNumber,
    initiatorParticipantId: dto.initiatorParticipantID,
    externalReferenceId: dto.externalReferenceID,
    trails: dto.trails?.map((trail) => ({
      id: trail.id,
      pledgeWithdrawalId: trail.pledgeWithdrawalID,
      addressCommandId: trail.addressCommandID,
      participantId: trail.participantID,
      action: trail.action,
      comment: trail.comment,
      createdAt: trail.createdAt,
    })),
    createdAt: dto.createdAt,
  };
}

/**
 * Extracts cursor pagination from response.
 */
function extractCursorPagination(cursor?: {
  currentPage?: string;
  hasNext?: boolean;
  hasPrevious?: boolean;
}): CursorPagination | undefined {
  if (!cursor) {
    return undefined;
  }

  return {
    currentPage: cursor.currentPage,
    hasNext: cursor.hasNext ?? false,
    hasPrevious: cursor.hasPrevious ?? false,
  };
}

/**
 * Service for Taurus Network pledge operations.
 *
 * Provides operations for creating, updating, and managing pledges between
 * Taurus Network participants. Pledges represent reserved funds from one
 * participant (pledgor) to another (pledgee).
 *
 * @example
 * ```typescript
 * // Create a pledge
 * const result = await pledgeService.createPledge({
 *   sharedAddressId: '123',
 *   currencyId: 'ETH',
 *   amount: '1000000000000000000',
 *   pledgeType: 'PLEDGEE_WITHDRAWALS_RIGHTS',
 * });
 * console.log(`Created pledge: ${result.pledge.id}`);
 *
 * // List pledges
 * const { pledges, pagination } = await pledgeService.list();
 * for (const pledge of pledges) {
 *   console.log(`${pledge.id}: ${pledge.amount}`);
 * }
 * ```
 */
export class PledgeService extends BaseService {
  private readonly pledgeApi: TaurusNetworkPledgeApi;

  /**
   * Creates a new PledgeService instance.
   *
   * @param pledgeApi - The TaurusNetworkPledgeApi instance from the OpenAPI client
   */
  constructor(pledgeApi: TaurusNetworkPledgeApi) {
    super();
    this.pledgeApi = pledgeApi;
  }

  /**
   * Gets a pledge by ID.
   *
   * @param pledgeId - The pledge ID to retrieve
   * @returns The pledge
   * @throws {@link ValidationError} If pledgeId is invalid
   * @throws {@link NotFoundError} If pledge not found
   * @throws {@link APIError} If API request fails
   */
  async get(pledgeId: string): Promise<Pledge> {
    if (!pledgeId || pledgeId.trim() === '') {
      throw new ValidationError('pledgeId is required');
    }

    return this.execute(async () => {
      const response = await this.pledgeApi.taurusNetworkServiceGetPledge({
        pledgeID: pledgeId,
      });

      const pledge = pledgeFromDto(response.result);
      if (!pledge) {
        throw new NotFoundError(`Pledge ${pledgeId} not found`);
      }

      return pledge;
    });
  }

  /**
   * Lists pledges with optional filtering.
   *
   * @param options - Optional filtering and pagination options
   * @returns Pledges list and pagination info
   * @throws {@link APIError} If API request fails
   */
  async list(
    options?: ListPledgesOptions
  ): Promise<{ pledges: Pledge[]; pagination?: CursorPagination }> {
    return this.execute(async () => {
      const response = await this.pledgeApi.taurusNetworkServiceGetPledges({
        ownerParticipantID: options?.ownerParticipantId,
        targetParticipantID: options?.targetParticipantId,
        sharedAddressIDs: options?.sharedAddressIds,
        currencyID: options?.currencyId,
        statuses: options?.statuses,
        sortOrder: options?.sortOrder,
        cursorCurrentPage: options?.currentPage,
        cursorPageRequest: options?.pageRequest,
        cursorPageSize: options?.pageSize ? String(options.pageSize) : undefined,
      });

      const pledges: Pledge[] = [];
      if (response.pledges) {
        for (const dto of response.pledges) {
          const pledge = pledgeFromDto(dto);
          if (pledge) {
            pledges.push(pledge);
          }
        }
      }

      return {
        pledges,
        pagination: extractCursorPagination(response.cursor),
      };
    });
  }

  /**
   * Creates a new pledge.
   *
   * Creates a pledge of funds from an internal address to a Taurus Network
   * participant. The funds will be reserved until unpledged or withdrawn.
   * Returns both the pledge and the action that needs approval.
   *
   * @param request - Pledge creation parameters
   * @returns The created pledge and action requiring approval
   * @throws {@link ValidationError} If required fields are missing
   * @throws {@link APIError} If API request fails
   */
  async createPledge(request: CreatePledgeRequest): Promise<CreatePledgeResult> {
    if (!request.sharedAddressId || request.sharedAddressId.trim() === '') {
      throw new ValidationError('sharedAddressId is required');
    }
    if (!request.currencyId || request.currencyId.trim() === '') {
      throw new ValidationError('currencyId is required');
    }
    if (!request.amount || request.amount.trim() === '') {
      throw new ValidationError('amount is required');
    }

    return this.execute(async () => {
      const response = await this.pledgeApi.taurusNetworkServiceCreatePledge({
        body: {
          sharedAddressID: request.sharedAddressId,
          currencyID: request.currencyId,
          amount: request.amount,
          pledgeType: request.pledgeType,
          externalReferenceId: request.externalReferenceId,
          reconciliationNote: request.reconciliationNote,
          pledgeDurationSetup: request.pledgeDurationSetup
            ? {
                minimumDuration: request.pledgeDurationSetup.minimumDuration,
                endOfMinimumDurationDate: request.pledgeDurationSetup.endOfMinimumDurationDate,
                noticePeriodDuration: request.pledgeDurationSetup.noticePeriodDuration,
              }
            : undefined,
          keyValueAttributes: request.keyValueAttributes?.map((attr) => ({
            key: attr.key,
            value: attr.value,
          })),
        },
      });

      const pledge = pledgeFromDto(response.result);
      if (!pledge) {
        throw new ValidationError('Failed to create pledge: no result returned');
      }

      return { pledge, pledgeActionId: response.pledgeActionID };
    });
  }

  /**
   * Updates a pledge's default destination.
   *
   * @param pledgeId - The pledge ID to update
   * @param request - Update parameters
   * @throws {@link ValidationError} If pledgeId is invalid
   * @throws {@link APIError} If API request fails
   */
  async updatePledge(pledgeId: string, request: UpdatePledgeRequest): Promise<void> {
    if (!pledgeId || pledgeId.trim() === '') {
      throw new ValidationError('pledgeId is required');
    }

    return this.execute(async () => {
      await this.pledgeApi.taurusNetworkServiceUpdatePledge({
        pledgeID: pledgeId,
        body: {
          defaultDestinationSharedAddressID: request.defaultDestinationSharedAddressId,
          defaultDestinationInternalAddressID: request.defaultDestinationInternalAddressId,
        },
      });
    });
  }

  /**
   * Adds collateral to an existing pledge.
   *
   * Increases the pledged amount. Only the pledgor can call this.
   * Returns the updated pledge and a new action requiring approval.
   *
   * @param pledgeId - The pledge ID to add collateral to
   * @param request - Collateral addition parameters
   * @returns Updated pledge and action requiring approval
   * @throws {@link ValidationError} If arguments are invalid
   * @throws {@link APIError} If API request fails
   */
  async addCollateral(pledgeId: string, request: AddPledgeCollateralRequest): Promise<AddCollateralResult> {
    if (!pledgeId || pledgeId.trim() === '') {
      throw new ValidationError('pledgeId is required');
    }
    if (!request.amount || request.amount.trim() === '') {
      throw new ValidationError('amount is required');
    }

    return this.execute(async () => {
      const response = await this.pledgeApi.taurusNetworkServiceAddPledgeCollateral({
        pledgeID: pledgeId,
        body: {
          amount: request.amount,
        },
      });

      return { pledgeActionId: response.pledgeActionID };
    });
  }

  /**
   * Withdraws from a pledge (pledgee operation).
   *
   * Allows the pledgee (target participant) to withdraw funds from the pledge.
   * Returns the withdrawal record and an action requiring approval.
   *
   * @param pledgeId - The pledge ID to withdraw from
   * @param request - Withdrawal parameters
   * @returns Withdrawal record and pledge action
   * @throws {@link ValidationError} If arguments are invalid
   * @throws {@link APIError} If API request fails
   */
  async withdrawPledge(pledgeId: string, request: WithdrawPledgeRequest): Promise<WithdrawPledgeResult> {
    if (!pledgeId || pledgeId.trim() === '') {
      throw new ValidationError('pledgeId is required');
    }
    if (!request.amount || request.amount.trim() === '') {
      throw new ValidationError('amount is required');
    }

    return this.execute(async () => {
      const response = await this.pledgeApi.taurusNetworkServiceWithdrawPledge({
        pledgeID: pledgeId,
        body: {
          amount: request.amount,
          destinationSharedAddressID: request.destinationSharedAddressId,
          destinationInternalAddressID: request.destinationInternalAddressId,
          externalReferenceID: request.externalReferenceId,
        },
      });

      return { pledgeWithdrawalId: response.pledgeWithdrawalID, pledgeActionId: response.pledgeActionID };
    });
  }

  /**
   * Initiates withdrawal from a pledge (pledgor operation).
   *
   * Allows the pledgor (owner participant) to initiate a withdrawal.
   * Returns the withdrawal record and an action requiring approval.
   *
   * @param pledgeId - The pledge ID to withdraw from
   * @param request - Withdrawal parameters
   * @returns Withdrawal record and pledge action
   * @throws {@link ValidationError} If arguments are invalid
   * @throws {@link APIError} If API request fails
   */
  async initiateWithdrawPledge(
    pledgeId: string,
    request: InitiateWithdrawPledgeRequest
  ): Promise<WithdrawPledgeResult> {
    if (!pledgeId || pledgeId.trim() === '') {
      throw new ValidationError('pledgeId is required');
    }
    if (!request.amount || request.amount.trim() === '') {
      throw new ValidationError('amount is required');
    }

    return this.execute(async () => {
      const response = await this.pledgeApi.taurusNetworkServiceInitiateWithdrawPledge({
        pledgeID: pledgeId,
        body: {
          amount: request.amount,
          destinationSharedAddressID: request.destinationSharedAddressId,
        },
      });

      return { pledgeWithdrawalId: response.pledgeWithdrawalID, pledgeActionId: response.pledgeActionID };
    });
  }

  /**
   * Unpledges all funds from a pledge.
   *
   * Releases all pledged funds back to the pledgor. Only the pledgor
   * can call this. Returns the updated pledge and an action requiring approval.
   *
   * @param pledgeId - The pledge ID to unpledge
   * @returns Updated pledge and pledge action requiring approval
   * @throws {@link ValidationError} If pledgeId is invalid
   * @throws {@link APIError} If API request fails
   */
  async unpledge(pledgeId: string): Promise<UnpledgeResult> {
    if (!pledgeId || pledgeId.trim() === '') {
      throw new ValidationError('pledgeId is required');
    }

    return this.execute(async () => {
      const response = await this.pledgeApi.taurusNetworkServiceUnpledge({
        pledgeID: pledgeId,
        body: {},
      });

      return { pledgeActionId: response.pledgeActionID };
    });
  }

  /**
   * Rejects a pledge.
   *
   * Rejects a pending pledge. Only the pledgee can reject an incoming pledge.
   *
   * @param pledgeId - The pledge ID to reject
   * @param request - Rejection parameters with comment
   * @throws {@link ValidationError} If arguments are invalid
   * @throws {@link APIError} If API request fails
   */
  async rejectPledge(pledgeId: string, request: RejectPledgeRequest): Promise<void> {
    if (!pledgeId || pledgeId.trim() === '') {
      throw new ValidationError('pledgeId is required');
    }
    if (!request.comment || request.comment.trim() === '') {
      throw new ValidationError('comment is required');
    }

    return this.execute(async () => {
      await this.pledgeApi.taurusNetworkServiceRejectPledge({
        pledgeID: pledgeId,
        body: {
          comment: request.comment,
        },
      });
    });
  }

  /**
   * Lists pledge actions with optional filtering.
   *
   * @param options - Optional filtering and pagination options
   * @returns Pledge actions list and pagination info
   * @throws {@link APIError} If API request fails
   */
  async listPledgeActions(
    options?: ListPledgeActionsOptions
  ): Promise<{ actions: PledgeAction[]; pagination?: CursorPagination }> {
    return this.execute(async () => {
      const response = await this.pledgeApi.taurusNetworkServiceGetPledgeActions({
        ids: options?.ids,
        pledgeID: options?.pledgeId,
        sortOrder: options?.sortOrder,
        cursorCurrentPage: options?.currentPage,
        cursorPageRequest: options?.pageRequest,
        cursorPageSize: options?.pageSize ? String(options.pageSize) : undefined,
      });

      const actions: PledgeAction[] = [];
      if (response.result) {
        for (const dto of response.result) {
          const action = pledgeActionFromDto(dto);
          if (action) {
            actions.push(action);
          }
        }
      }

      return {
        actions,
        pagination: extractCursorPagination(response.cursor),
      };
    });
  }

  /**
   * Lists pledge actions pending approval.
   *
   * Returns only actions that require approval from the current user.
   *
   * @param options - Optional filtering and pagination options
   * @returns Pledge actions list and pagination info
   * @throws {@link APIError} If API request fails
   */
  async listPledgeActionsForApproval(
    options?: ListPledgeActionsOptions
  ): Promise<{ actions: PledgeAction[]; pagination?: CursorPagination }> {
    return this.execute(async () => {
      const response = await this.pledgeApi.taurusNetworkServiceGetPledgeActionsForApproval({
        ids: options?.ids,
        types: options?.types,
        sortOrder: options?.sortOrder,
        cursorCurrentPage: options?.currentPage,
        cursorPageRequest: options?.pageRequest,
        cursorPageSize: options?.pageSize ? String(options.pageSize) : undefined,
      });

      const actions: PledgeAction[] = [];
      if (response.result) {
        for (const dto of response.result) {
          const action = pledgeActionFromDto(dto);
          if (action) {
            actions.push(action);
          }
        }
      }

      return {
        actions,
        pagination: extractCursorPagination(response.cursor),
      };
    });
  }

  /**
   * Approves multiple pledge actions with ECDSA signature.
   *
   * The actions are sorted by ID, and a signature is computed over
   * the concatenated hashes of their metadata using the provided private key.
   *
   * @param actionIds - List of pledge action IDs to approve
   * @param signature - Base64-encoded ECDSA signature
   * @param comment - Optional approval comment
   * @returns Number of actions successfully approved
   * @throws {@link ValidationError} If arguments are invalid
   * @throws {@link APIError} If API request fails
   */
  async approvePledgeActions(
    actionIds: string[],
    signature: string,
    comment: string = 'approving via taurus-protect-sdk-typescript'
  ): Promise<number> {
    if (!actionIds || actionIds.length === 0) {
      throw new ValidationError('actionIds list cannot be empty');
    }
    if (!signature || signature.trim() === '') {
      throw new ValidationError('signature is required');
    }

    return this.execute(async () => {
      const response = await this.pledgeApi.taurusNetworkServiceApprovePledgeActions({
        body: {
          ids: actionIds,
          signature,
          comment,
        },
      });

      return actionIds.length;
    });
  }

  /**
   * Rejects multiple pledge actions.
   *
   * @param request - Rejection request with action IDs and comment
   * @throws {@link ValidationError} If request is invalid
   * @throws {@link APIError} If API request fails
   */
  async rejectPledgeActions(request: RejectPledgeActionsRequest): Promise<void> {
    if (!request.ids || request.ids.length === 0) {
      throw new ValidationError('ids list cannot be empty');
    }
    if (!request.comment || request.comment.trim() === '') {
      throw new ValidationError('comment is required');
    }

    return this.execute(async () => {
      await this.pledgeApi.taurusNetworkServiceRejectPledgeActions({
        body: {
          ids: request.ids,
          comment: request.comment,
        },
      });
    });
  }

  /**
   * Lists pledge withdrawals with optional filtering.
   *
   * @param options - Optional filtering and pagination options
   * @returns Withdrawals list and pagination info
   * @throws {@link APIError} If API request fails
   */
  async listPledgeWithdrawals(
    options?: ListPledgeWithdrawalsOptions
  ): Promise<{ withdrawals: PledgeWithdrawal[]; pagination?: CursorPagination }> {
    return this.execute(async () => {
      const response = await this.pledgeApi.taurusNetworkServiceGetPledgesWithdrawals({
        pledgeID: options?.pledgeId,
        withdrawalStatus: options?.withdrawalStatus,
        sortOrder: options?.sortOrder,
        cursorCurrentPage: options?.currentPage,
        cursorPageRequest: options?.pageRequest,
        cursorPageSize: options?.pageSize ? String(options.pageSize) : undefined,
      });

      const withdrawals: PledgeWithdrawal[] = [];
      if (response.withdrawals) {
        for (const dto of response.withdrawals) {
          const withdrawal = pledgeWithdrawalFromDto(dto);
          if (withdrawal) {
            withdrawals.push(withdrawal);
          }
        }
      }

      return {
        withdrawals,
        pagination: extractCursorPagination(response.cursor),
      };
    });
  }
}
