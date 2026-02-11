/**
 * Shared Address and Shared Asset models for Taurus Network.
 *
 * This module provides domain models for shared addresses and shared assets,
 * which allow participants to share addresses and tokens for pledging,
 * transfers, and settlements.
 */

/**
 * Shared address status enum.
 */
export enum SharedAddressStatus {
  PENDING = "PENDING",
  ACTIVE = "ACTIVE",
  REJECTED = "REJECTED",
  DELETED = "DELETED",
}

/**
 * Shared asset status enum.
 */
export enum SharedAssetStatus {
  PENDING = "PENDING",
  ACTIVE = "ACTIVE",
  REJECTED = "REJECTED",
  DELETED = "DELETED",
}

/**
 * Audit trail entry for shared address changes.
 */
export interface SharedAddressTrail {
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
 * Proof of ownership for a shared address.
 */
export interface SharedAddressProofOfOwnership {
  /** Ownership signature. */
  readonly signature: string | undefined;
  /** Signed message. */
  readonly message: string | undefined;
}

/**
 * Shared address between participants.
 *
 * A shared address allows one participant (owner) to share an address
 * with another participant (target) for pledging and transfers.
 */
export interface SharedAddress {
  /** Unique shared address identifier (UUID). */
  readonly id: string | undefined;
  /** Internal address ID (for owner). */
  readonly internalAddressId: string | undefined;
  /** Whitelisted address ID (for target). */
  readonly wladdressId: string | undefined;
  /** Owner participant ID. */
  readonly ownerParticipantId: string | undefined;
  /** Target participant ID. */
  readonly targetParticipantId: string | undefined;
  /** Blockchain name. */
  readonly blockchain: string | undefined;
  /** Network name. */
  readonly network: string | undefined;
  /** Blockchain address string. */
  readonly address: string | undefined;
  /** Original label. */
  readonly originLabel: string | undefined;
  /** Network creation date. */
  readonly originCreationDate: Date | undefined;
  /** Network deletion date. */
  readonly originDeletionDate: Date | undefined;
  /** Local creation timestamp. */
  readonly createdAt: Date | undefined;
  /** Local update timestamp. */
  readonly updatedAt: Date | undefined;
  /** When target accepted. */
  readonly targetAcceptedAt: Date | undefined;
  /** Shared address status. */
  readonly status: string | undefined;
  /** Proof of ownership. */
  readonly proofOfOwnership: SharedAddressProofOfOwnership | undefined;
  /** Number of active pledges. */
  readonly pledgesCount: string | undefined;
  /** Audit trail. */
  readonly trails: SharedAddressTrail[];
}

/**
 * Audit trail entry for shared asset changes.
 */
export interface SharedAssetTrail {
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
 * Shared asset (token) between participants.
 *
 * A shared asset allows participants to share information about
 * tokens/contracts for use in settlements and transfers.
 */
export interface SharedAsset {
  /** Unique shared asset identifier. */
  readonly id: string | undefined;
  /** Whitelisted contract address ID. */
  readonly wlContractAddressId: string | undefined;
  /** Owner participant ID. */
  readonly ownerParticipantId: string | undefined;
  /** Target participant ID. */
  readonly targetParticipantId: string | undefined;
  /** Blockchain name. */
  readonly blockchain: string | undefined;
  /** Network name. */
  readonly network: string | undefined;
  /** Token/asset name. */
  readonly name: string | undefined;
  /** Token symbol. */
  readonly symbol: string | undefined;
  /** Token decimals. */
  readonly decimals: string | undefined;
  /** Smart contract address. */
  readonly contractAddress: string | undefined;
  /** Token ID (for NFTs). */
  readonly tokenId: string | undefined;
  /** Asset kind (ERC20, ERC721, etc.). */
  readonly kind: string | undefined;
  /** Network creation date. */
  readonly originCreationDate: Date | undefined;
  /** Network deletion date. */
  readonly originDeletionDate: Date | undefined;
  /** Local creation timestamp. */
  readonly createdAt: Date | undefined;
  /** Local update timestamp. */
  readonly updatedAt: Date | undefined;
  /** Target acceptance timestamp. */
  readonly targetAcceptedAt: Date | undefined;
  /** Target rejection timestamp. */
  readonly targetRejectedAt: Date | undefined;
  /** Shared asset status. */
  readonly status: string | undefined;
  /** Audit trail. */
  readonly trails: SharedAssetTrail[];
}

// Request models

/**
 * Request to create a new shared address.
 */
export interface CreateSharedAddressRequest {
  /** Internal address ID to share. */
  readonly internalAddressId: string;
  /** Target participant ID. */
  readonly targetParticipantId: string;
  /** Label for the shared address. */
  readonly label?: string;
  /** Proof of ownership signature. */
  readonly proofOfOwnershipSignature?: string;
  /** Proof of ownership message. */
  readonly proofOfOwnershipMessage?: string;
}

/**
 * Request to accept a shared address.
 */
export interface AcceptSharedAddressRequest {
  /** Acceptance comment. */
  readonly comment?: string;
}

/**
 * Request to reject a shared address.
 */
export interface RejectSharedAddressRequest {
  /** Rejection reason. */
  readonly comment: string;
}

/**
 * Request to create a new shared asset.
 */
export interface CreateSharedAssetRequest {
  /** Whitelisted contract address ID. */
  readonly wlContractAddressId: string;
  /** Target participant ID. */
  readonly targetParticipantId: string;
}

/**
 * Request to accept a shared asset.
 */
export interface AcceptSharedAssetRequest {
  /** Acceptance comment. */
  readonly comment?: string;
}

/**
 * Request to reject a shared asset.
 */
export interface RejectSharedAssetRequest {
  /** Rejection reason. */
  readonly comment: string;
}

// Filter options

/**
 * Options for listing shared addresses.
 */
export interface ListSharedAddressesOptions {
  /** Maximum items to return (default: 50, max: 1000). */
  readonly limit?: number;
  /** Number of items to skip. */
  readonly offset?: number;
  /** Filter by statuses. */
  readonly statuses?: string[];
  /** Filter by participant ID (owner or target). */
  readonly participantId?: string;
  /** Filter by blockchain. */
  readonly blockchain?: string;
  /** Filter by network. */
  readonly network?: string;
}

/**
 * Options for listing shared assets.
 */
export interface ListSharedAssetsOptions {
  /** Maximum items to return (default: 50, max: 1000). */
  readonly limit?: number;
  /** Number of items to skip. */
  readonly offset?: number;
  /** Filter by statuses. */
  readonly statuses?: string[];
  /** Filter by participant ID (owner or target). */
  readonly participantId?: string;
  /** Filter by blockchain. */
  readonly blockchain?: string;
  /** Filter by network. */
  readonly network?: string;
}
