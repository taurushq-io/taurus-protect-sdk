/**
 * Settlement service for Taurus Network in Taurus-PROTECT SDK.
 *
 * Provides methods for managing Taurus Network settlements between participants.
 */

import { NotFoundError, ValidationError } from '../../errors';
import type { TaurusNetworkSettlementApi } from '../../internal/openapi/apis/TaurusNetworkSettlementApi';
import type { TgvalidatordTnSettlement } from '../../internal/openapi/models/index';
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
 * Asset transfer in a settlement.
 */
export interface SettlementAssetTransfer {
  id?: string;
  sourceSharedAddressId?: string;
  destinationSharedAddressId?: string;
  currencyId?: string;
  blockchain?: string;
  network?: string;
  amount?: string;
  status?: string;
  pledgeId?: string;
  createdAt?: Date;
  updatedAt?: Date;
}

/**
 * Settlement clip transaction.
 */
export interface SettlementClipTransaction {
  requestId?: string;
  txHash?: string;
  status?: string;
}

/**
 * Settlement clip.
 */
export interface SettlementClip {
  id?: string;
  settlementId?: string;
  sourceSharedAddressId?: string;
  destinationSharedAddressId?: string;
  currencyId?: string;
  blockchain?: string;
  network?: string;
  amount?: string;
  status?: string;
  transactions?: SettlementClipTransaction[];
  createdAt?: Date;
  updatedAt?: Date;
}

/**
 * A Taurus Network settlement.
 */
export interface Settlement {
  id: string;
  creatorParticipantId?: string;
  targetParticipantId?: string;
  firstLegParticipantId?: string;
  status?: string;
  workflowId?: string;
  startExecutionDate?: Date;
  firstLegAssets?: SettlementAssetTransfer[];
  secondLegAssets?: SettlementAssetTransfer[];
  clips?: SettlementClip[];
  createdAt?: Date;
  updatedAt?: Date;
}

/**
 * Options for listing settlements.
 */
export interface ListSettlementsOptions {
  counterParticipantId?: string;
  statuses?: string[];
  sortOrder?: string;
  pageSize?: number;
  currentPage?: string;
  pageRequest?: string;
}

/**
 * Options for listing settlements for approval.
 */
export interface ListSettlementsForApprovalOptions {
  ids?: string[];
  sortOrder?: string;
  pageSize?: number;
  currentPage?: string;
  pageRequest?: string;
}

/**
 * Asset transfer definition for creating a settlement.
 */
export interface AssetTransferRequest {
  sourceSharedAddressId: string;
  destinationSharedAddressId: string;
  currencyId: string;
  amount: string;
  pledgeId?: string;
}

/**
 * Request to create a settlement.
 */
export interface CreateSettlementRequest {
  targetParticipantId: string;
  firstLegParticipantId: string;
  firstLegAssets: AssetTransferRequest[];
  secondLegAssets: AssetTransferRequest[];
  startExecutionDate?: Date;
}

/**
 * Request to replace (update) a settlement.
 */
export interface ReplaceSettlementRequest {
  firstLegAssets?: AssetTransferRequest[];
  secondLegAssets?: AssetTransferRequest[];
  startExecutionDate?: Date;
}

/**
 * Maps asset transfer DTOs.
 */
function assetTransferFromDto(
  dto: {
    id?: string;
    sourceSharedAddressID?: string;
    destinationSharedAddressID?: string;
    currencyID?: string;
    blockchain?: string;
    network?: string;
    amount?: string;
    status?: string;
    pledgeID?: string;
    createdAt?: Date;
    updatedAt?: Date;
  }
): SettlementAssetTransfer {
  return {
    id: dto.id,
    sourceSharedAddressId: dto.sourceSharedAddressID,
    destinationSharedAddressId: dto.destinationSharedAddressID,
    currencyId: dto.currencyID,
    blockchain: dto.blockchain,
    network: dto.network,
    amount: dto.amount,
    status: dto.status,
    pledgeId: dto.pledgeID,
    createdAt: dto.createdAt,
    updatedAt: dto.updatedAt,
  };
}

/**
 * Maps settlement clip DTOs.
 */
function settlementClipFromDto(
  dto: {
    id?: string;
    settlementID?: string;
    sourceSharedAddressID?: string;
    destinationSharedAddressID?: string;
    currencyID?: string;
    blockchain?: string;
    network?: string;
    amount?: string;
    status?: string;
    transactions?: Array<{
      requestID?: string;
      txHash?: string;
      status?: string;
    }>;
    createdAt?: Date;
    updatedAt?: Date;
  }
): SettlementClip {
  return {
    id: dto.id,
    settlementId: dto.settlementID,
    sourceSharedAddressId: dto.sourceSharedAddressID,
    destinationSharedAddressId: dto.destinationSharedAddressID,
    currencyId: dto.currencyID,
    blockchain: dto.blockchain,
    network: dto.network,
    amount: dto.amount,
    status: dto.status,
    transactions: dto.transactions?.map((t) => ({
      requestId: t.requestID,
      txHash: t.txHash,
      status: t.status,
    })),
    createdAt: dto.createdAt,
    updatedAt: dto.updatedAt,
  };
}

/**
 * Maps a DTO to a Settlement.
 */
function settlementFromDto(dto?: TgvalidatordTnSettlement): Settlement | undefined {
  if (!dto) {
    return undefined;
  }

  return {
    id: dto.id ?? '',
    creatorParticipantId: dto.creatorParticipantID,
    targetParticipantId: dto.targetParticipantID,
    firstLegParticipantId: dto.firstLegParticipantID,
    status: dto.status,
    workflowId: dto.workflowID,
    startExecutionDate: dto.startExecutionDate,
    firstLegAssets: dto.firstLegAssets?.map(assetTransferFromDto),
    secondLegAssets: dto.secondLegAssets?.map(assetTransferFromDto),
    clips: dto.clips?.map(settlementClipFromDto),
    createdAt: dto.createdAt,
    updatedAt: dto.updatedAt,
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
 * Service for Taurus Network settlement operations.
 *
 * Provides methods to create, manage, and monitor settlements between
 * Taurus Network participants. Settlements enable atomic swaps of assets
 * between participants.
 *
 * @example
 * ```typescript
 * // List settlements
 * const { settlements, pagination } = await settlementService.list({
 *   statuses: ['PENDING', 'COMPLETED'],
 * });
 *
 * // Get a single settlement
 * const settlement = await settlementService.get('settlement-123');
 * console.log(`Status: ${settlement.status}`);
 * ```
 */
export class SettlementService extends BaseService {
  private readonly settlementApi: TaurusNetworkSettlementApi;

  /**
   * Creates a new SettlementService instance.
   *
   * @param settlementApi - The TaurusNetworkSettlementApi instance from the OpenAPI client
   */
  constructor(settlementApi: TaurusNetworkSettlementApi) {
    super();
    this.settlementApi = settlementApi;
  }

  /**
   * Gets a settlement by ID.
   *
   * @param settlementId - The settlement ID to retrieve
   * @returns The settlement
   * @throws {@link ValidationError} If settlementId is empty
   * @throws {@link NotFoundError} If settlement not found
   * @throws {@link APIError} If API request fails
   */
  async get(settlementId: string): Promise<Settlement> {
    if (!settlementId || settlementId.trim() === '') {
      throw new ValidationError('settlementId is required');
    }

    return this.execute(async () => {
      const response = await this.settlementApi.taurusNetworkServiceGetSettlement({
        settlementID: settlementId,
      });

      const settlement = settlementFromDto(response.result);
      if (!settlement) {
        throw new NotFoundError(`Settlement ${settlementId} not found`);
      }

      return settlement;
    });
  }

  /**
   * Lists settlements with optional filtering.
   *
   * @param options - Optional filtering and pagination options
   * @returns Settlements list and pagination info
   * @throws {@link APIError} If API request fails
   */
  async list(
    options?: ListSettlementsOptions
  ): Promise<{ settlements: Settlement[]; pagination?: CursorPagination }> {
    return this.execute(async () => {
      const response = await this.settlementApi.taurusNetworkServiceGetSettlements({
        counterParticipantID: options?.counterParticipantId,
        statuses: options?.statuses,
        sortOrder: options?.sortOrder,
        cursorCurrentPage: options?.currentPage,
        cursorPageRequest: options?.pageRequest,
        cursorPageSize: options?.pageSize ? String(options.pageSize) : undefined,
      });

      const settlements: Settlement[] = [];
      if (response.result) {
        for (const dto of response.result) {
          const settlement = settlementFromDto(dto);
          if (settlement) {
            settlements.push(settlement);
          }
        }
      }

      return {
        settlements,
        pagination: extractCursorPagination(response.cursor),
      };
    });
  }

  /**
   * Lists settlements pending approval.
   *
   * @param options - Optional filtering and pagination options
   * @returns Settlements list and pagination info
   * @throws {@link APIError} If API request fails
   */
  async listForApproval(
    options?: ListSettlementsForApprovalOptions
  ): Promise<{ settlements: Settlement[]; pagination?: CursorPagination }> {
    return this.execute(async () => {
      const response = await this.settlementApi.taurusNetworkServiceGetSettlementsForApproval({
        ids: options?.ids,
        sortOrder: options?.sortOrder,
        cursorCurrentPage: options?.currentPage,
        cursorPageRequest: options?.pageRequest,
        cursorPageSize: options?.pageSize ? String(options.pageSize) : undefined,
      });

      const settlements: Settlement[] = [];
      if (response.result) {
        for (const dto of response.result) {
          const settlement = settlementFromDto(dto);
          if (settlement) {
            settlements.push(settlement);
          }
        }
      }

      return {
        settlements,
        pagination: extractCursorPagination(response.cursor),
      };
    });
  }

  /**
   * Creates a new settlement.
   *
   * Creates a settlement between two participants with defined asset
   * transfers in both legs. Once approved by both parties, the settlement
   * will execute atomically.
   *
   * @param request - Settlement creation parameters
   * @returns The created settlement ID
   * @throws {@link ValidationError} If required fields are missing
   * @throws {@link APIError} If API request fails
   */
  async create(request: CreateSettlementRequest): Promise<string> {
    if (!request.targetParticipantId || request.targetParticipantId.trim() === '') {
      throw new ValidationError('targetParticipantId is required');
    }
    if (!request.firstLegParticipantId || request.firstLegParticipantId.trim() === '') {
      throw new ValidationError('firstLegParticipantId is required');
    }
    if (!request.firstLegAssets || request.firstLegAssets.length === 0) {
      throw new ValidationError('firstLegAssets is required');
    }
    if (!request.secondLegAssets || request.secondLegAssets.length === 0) {
      throw new ValidationError('secondLegAssets is required');
    }

    return this.execute(async () => {
      const response = await this.settlementApi.taurusNetworkServiceCreateSettlement({
        body: {
          targetParticipantID: request.targetParticipantId,
          firstLegParticipantID: request.firstLegParticipantId,
          firstLegAssets: request.firstLegAssets.map((a) => ({
            sourceSharedAddressID: a.sourceSharedAddressId,
            destinationSharedAddressID: a.destinationSharedAddressId,
            currencyID: a.currencyId,
            amount: a.amount,
            pledgeID: a.pledgeId,
          })),
          secondLegAssets: request.secondLegAssets.map((a) => ({
            sourceSharedAddressID: a.sourceSharedAddressId,
            destinationSharedAddressID: a.destinationSharedAddressId,
            currencyID: a.currencyId,
            amount: a.amount,
            pledgeID: a.pledgeId,
          })),
          startExecutionDate: request.startExecutionDate,
        },
      });

      return response.settlementID ?? '';
    });
  }

  /**
   * Replaces (updates) a settlement.
   *
   * Updates the settlement's asset transfers before it has been approved.
   *
   * @param settlementId - The settlement ID to update
   * @param request - Update parameters
   * @throws {@link ValidationError} If settlementId is empty
   * @throws {@link APIError} If API request fails
   */
  async replace(settlementId: string, request: ReplaceSettlementRequest): Promise<void> {
    if (!settlementId || settlementId.trim() === '') {
      throw new ValidationError('settlementId is required');
    }

    return this.execute(async () => {
      await this.settlementApi.taurusNetworkServiceReplaceSettlement({
        settlementID: settlementId,
        body: {
          createSettlementRequest: {
            targetParticipantID: '', // Will be preserved from original
            firstLegParticipantID: '', // Will be preserved from original
            firstLegAssets: request.firstLegAssets?.map((a) => ({
              sourceSharedAddressID: a.sourceSharedAddressId,
              destinationSharedAddressID: a.destinationSharedAddressId,
              currencyID: a.currencyId,
              amount: a.amount,
              pledgeID: a.pledgeId,
            })) ?? [],
            secondLegAssets: request.secondLegAssets?.map((a) => ({
              sourceSharedAddressID: a.sourceSharedAddressId,
              destinationSharedAddressID: a.destinationSharedAddressId,
              currencyID: a.currencyId,
              amount: a.amount,
              pledgeID: a.pledgeId,
            })) ?? [],
            startExecutionDate: request.startExecutionDate,
          },
        },
      });
    });
  }

  /**
   * Cancels a settlement.
   *
   * Cancels a settlement before it has been executed. Can only be called
   * by the creator.
   *
   * @param settlementId - The settlement ID to cancel
   * @throws {@link ValidationError} If settlementId is empty
   * @throws {@link APIError} If API request fails
   */
  async cancel(settlementId: string): Promise<void> {
    if (!settlementId || settlementId.trim() === '') {
      throw new ValidationError('settlementId is required');
    }

    return this.execute(async () => {
      await this.settlementApi.taurusNetworkServiceCancelSettlement({
        settlementID: settlementId,
        body: {},
      });
    });
  }
}
