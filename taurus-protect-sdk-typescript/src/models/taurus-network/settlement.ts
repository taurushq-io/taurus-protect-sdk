/**
 * Settlement models for Taurus Network.
 *
 * This module provides domain models for settlements, which represent
 * atomic exchanges of assets between two participants.
 */

import type { CursorNavigationOptions } from '../pagination';

/**
 * Settlement status enum.
 */
export enum SettlementStatus {
  PENDING = "PENDING",
  PENDING_APPROVAL = "PENDING_APPROVAL",
  APPROVED = "APPROVED",
  EXECUTING = "EXECUTING",
  COMPLETED = "COMPLETED",
  REJECTED = "REJECTED",
  CANCELED = "CANCELED",
  FAILED = "FAILED",
}

/**
 * Asset transfer within a settlement.
 *
 * Represents one leg of an asset transfer in a settlement.
 */
export interface SettlementAssetTransfer {
  /** Currency being transferred. */
  readonly currencyId: string | undefined;
  /** Amount to transfer. */
  readonly amount: string | undefined;
  /** Source shared address ID. */
  readonly sourceSharedAddressId: string | undefined;
  /** Destination shared address ID. */
  readonly destinationSharedAddressId: string | undefined;
}

/**
 * Transaction within a settlement clip.
 */
export interface SettlementClipTransaction {
  /** Transaction ID. */
  readonly id: string | undefined;
  /** Blockchain transaction hash. */
  readonly txHash: string | undefined;
  /** Transaction status. */
  readonly status: string | undefined;
  /** Transaction amount. */
  readonly amount: string | undefined;
  /** Currency ID. */
  readonly currencyId: string | undefined;
}

/**
 * Settlement clip representing a portion of the settlement.
 */
export interface SettlementClip {
  /** Clip ID. */
  readonly id: string | undefined;
  /** Clip status. */
  readonly status: string | undefined;
  /** Clip transactions. */
  readonly transactions: SettlementClipTransaction[];
}

/**
 * Audit trail entry for settlement changes.
 */
export interface SettlementTrail {
  /** Trail entry ID. */
  readonly id: string | undefined;
  /** Trail timestamp. */
  readonly timestamp: Date | undefined;
  /** Action performed. */
  readonly action: string | undefined;
  /** Actor who performed action. */
  readonly actor: string | undefined;
  /** Optional comment. */
  readonly comment: string | undefined;
}

/**
 * Settlement between participants.
 *
 * Represents an atomic exchange of assets between two participants.
 */
export interface Settlement {
  /** Unique settlement identifier. */
  readonly id: string | undefined;
  /** Participant who created the settlement. */
  readonly creatorParticipantId: string | undefined;
  /** Counter-party participant. */
  readonly targetParticipantId: string | undefined;
  /** Who executes the first leg. */
  readonly firstLegParticipantId: string | undefined;
  /** Assets in first leg. */
  readonly firstLegAssets: SettlementAssetTransfer[];
  /** Assets in second leg. */
  readonly secondLegAssets: SettlementAssetTransfer[];
  /** Settlement clips. */
  readonly clips: SettlementClip[];
  /** When execution started. */
  readonly startExecutionDate: Date | undefined;
  /** Settlement status. */
  readonly status: string | undefined;
  /** Creation timestamp. */
  readonly createdAt: Date | undefined;
  /** Update timestamp. */
  readonly updatedAt: Date | undefined;
  /** Associated workflow ID. */
  readonly workflowId: string | undefined;
  /** Audit trail. */
  readonly trails: SettlementTrail[];
}

// Request models

/**
 * Asset transfer specification for creating a settlement.
 */
export interface SettlementAssetTransferRequest {
  /** Currency ID. */
  readonly currencyId: string;
  /** Amount to transfer. */
  readonly amount: string;
  /** Source shared address ID. */
  readonly sourceSharedAddressId: string;
  /** Destination shared address ID. */
  readonly destinationSharedAddressId: string;
}

/**
 * Request to create a new settlement.
 */
export interface CreateSettlementRequest {
  /** Target participant ID. */
  readonly targetParticipantId: string;
  /** Who executes the first leg. */
  readonly firstLegParticipantId: string;
  /** First leg asset transfers. */
  readonly firstLegAssets: SettlementAssetTransferRequest[];
  /** Second leg asset transfers. */
  readonly secondLegAssets: SettlementAssetTransferRequest[];
}

/**
 * Request to accept a settlement.
 */
export interface AcceptSettlementRequest {
  /** Acceptance comment. */
  readonly comment?: string;
}

/**
 * Request to reject a settlement.
 */
export interface RejectSettlementRequest {
  /** Rejection reason. */
  readonly comment: string;
}

// Filter options

/**
 * Options for listing settlements. A cursor list: `pageSize` 1-100 (default 20) and the
 * `cursor` of a previous page.
 */
export interface ListSettlementsOptions extends CursorNavigationOptions {
  /** Filter by counterparty participant ID */
  readonly counterParticipantId?: string;
  /** Filter by statuses */
  readonly statuses?: string[];
  /** Sort order ("ASC" or "DESC") */
  readonly sortOrder?: string;
}

/**
 * Options for listing the settlements awaiting approval. A cursor list: `pageSize` 1-100
 * (default 20) and the `cursor` of a previous page.
 */
export interface ListSettlementsForApprovalOptions extends CursorNavigationOptions {
  /** Filter by settlement IDs */
  readonly ids?: string[];
  /** Sort order ("ASC" or "DESC") */
  readonly sortOrder?: string;
}
