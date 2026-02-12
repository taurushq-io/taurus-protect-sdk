/**
 * Sharing service for Taurus Network in Taurus-PROTECT SDK.
 *
 * Provides methods for managing shared addresses and assets between
 * Taurus Network participants.
 */

import { ValidationError } from '../../errors';
import type { TaurusNetworkSharedAddressAssetApi } from '../../internal/openapi/apis/TaurusNetworkSharedAddressAssetApi';
import type {
  TgvalidatordTnSharedAddress,
  TgvalidatordTnSharedAsset,
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
 * Trail entry for shared address.
 */
export interface SharedAddressTrail {
  id?: string;
  sharedAddressId?: string;
  addressStatus?: string;
  comment?: string;
  createdAt?: Date;
}

/**
 * Proof of ownership payload.
 */
export interface ProofOfOwnershipPayload {
  ownerParticipantId?: string;
  targetParticipantId?: string;
  address?: string;
  blockchain?: string;
  network?: string;
}

/**
 * Signed proof of ownership payload.
 */
export interface SignedProofOfOwnershipPayload {
  payload?: ProofOfOwnershipPayload;
  ownerParticipantSignature?: string;
}

/**
 * Proof of reserve for a shared address.
 */
export interface ProofOfReserve {
  curve?: string;
  cipher?: string;
  path?: string;
  address?: string;
  publicKey?: string;
  challenge?: string;
  challengeResponse?: string;
  type?: string;
  stakePublicKey?: string;
  stakeChallengeResponse?: string;
}

/**
 * Proof of ownership for a shared address.
 */
export interface ProofOfOwnership {
  signedPayload?: SignedProofOfOwnershipPayload;
  signedPayloadHash?: string;
  proofOfReserve?: ProofOfReserve;
  signedPayloadAsString?: string;
}

/**
 * A shared address in the Taurus Network.
 */
export interface SharedAddress {
  id: string;
  internalAddressId?: string;
  whitelistedAddressId?: string;
  ownerParticipantId?: string;
  targetParticipantId?: string;
  blockchain?: string;
  network?: string;
  address?: string;
  originLabel?: string;
  status?: string;
  pledgesCount?: string;
  proofOfOwnership?: ProofOfOwnership;
  trails?: SharedAddressTrail[];
  originCreationDate?: Date;
  originDeletionDate?: Date;
  targetAcceptedAt?: Date;
  createdAt?: Date;
  updatedAt?: Date;
}

/**
 * Trail entry for shared asset.
 */
export interface SharedAssetTrail {
  id?: string;
  sharedAssetId?: string;
  assetStatus?: string;
  comment?: string;
  createdAt?: Date;
}

/**
 * A shared asset (whitelisted contract) in the Taurus Network.
 */
export interface SharedAsset {
  id: string;
  whitelistedContractAddressId?: string;
  ownerParticipantId?: string;
  targetParticipantId?: string;
  blockchain?: string;
  network?: string;
  name?: string;
  symbol?: string;
  decimals?: string;
  contractAddress?: string;
  tokenId?: string;
  kind?: string;
  status?: string;
  trails?: SharedAssetTrail[];
  originCreationDate?: Date;
  originDeletionDate?: Date;
  targetAcceptedAt?: Date;
  targetRejectedAt?: Date;
  createdAt?: Date;
  updatedAt?: Date;
}

/**
 * Options for listing shared addresses.
 */
export interface ListSharedAddressesOptions {
  participantId?: string;
  ownerParticipantId?: string;
  targetParticipantId?: string;
  blockchain?: string;
  network?: string;
  ids?: string[];
  statuses?: string[];
  sortOrder?: string;
  pageSize?: number;
  currentPage?: string;
  pageRequest?: string;
}

/**
 * Options for listing shared assets.
 */
export interface ListSharedAssetsOptions {
  participantId?: string;
  ownerParticipantId?: string;
  targetParticipantId?: string;
  blockchain?: string;
  network?: string;
  ids?: string[];
  statuses?: string[];
  sortOrder?: string;
  pageSize?: number;
  currentPage?: string;
  pageRequest?: string;
}

/**
 * Key-value attribute for sharing.
 */
export interface KeyValueAttribute {
  key: string;
  value: string;
}

/**
 * Request to share an address.
 */
export interface ShareAddressRequest {
  addressId: string;
  toParticipantId: string;
  keyValueAttributes?: KeyValueAttribute[];
}

/**
 * Request to share a whitelisted asset.
 */
export interface ShareWhitelistedAssetRequest {
  whitelistedContractId: string;
  toParticipantId: string;
}

/**
 * Maps a DTO to a SharedAddress.
 */
function sharedAddressFromDto(dto?: TgvalidatordTnSharedAddress): SharedAddress | undefined {
  if (!dto) {
    return undefined;
  }

  return {
    id: dto.id ?? '',
    internalAddressId: dto.internalAddressID,
    whitelistedAddressId: dto.wladdressID,
    ownerParticipantId: dto.ownerParticipantId,
    targetParticipantId: dto.targetParticipantId,
    blockchain: dto.blockchain,
    network: dto.network,
    address: dto.address,
    originLabel: dto.originLabel,
    status: dto.status,
    pledgesCount: dto.pledgesCount,
    proofOfOwnership: dto.proofOfOwnership
      ? {
          signedPayload: dto.proofOfOwnership.signedPayload
            ? {
                payload: dto.proofOfOwnership.signedPayload.payload
                  ? {
                      ownerParticipantId: dto.proofOfOwnership.signedPayload.payload.ownerParticipantID,
                      targetParticipantId: dto.proofOfOwnership.signedPayload.payload.targetParticipantID,
                      address: dto.proofOfOwnership.signedPayload.payload.address,
                      blockchain: dto.proofOfOwnership.signedPayload.payload.blockchain,
                      network: dto.proofOfOwnership.signedPayload.payload.network,
                    }
                  : undefined,
                ownerParticipantSignature: dto.proofOfOwnership.signedPayload.ownerParticipantSignature,
              }
            : undefined,
          signedPayloadHash: dto.proofOfOwnership.signedPayloadHash,
          proofOfReserve: dto.proofOfOwnership.proofOfReserve
            ? {
                curve: dto.proofOfOwnership.proofOfReserve.curve,
                cipher: dto.proofOfOwnership.proofOfReserve.cipher,
                path: dto.proofOfOwnership.proofOfReserve.path,
                address: dto.proofOfOwnership.proofOfReserve.address,
                publicKey: dto.proofOfOwnership.proofOfReserve.publicKey,
                challenge: dto.proofOfOwnership.proofOfReserve.challenge,
                challengeResponse: dto.proofOfOwnership.proofOfReserve.challengeResponse,
                type: dto.proofOfOwnership.proofOfReserve.type,
                stakePublicKey: dto.proofOfOwnership.proofOfReserve.stakePublicKey,
                stakeChallengeResponse: dto.proofOfOwnership.proofOfReserve.stakeChallengeResponse,
              }
            : undefined,
          signedPayloadAsString: dto.proofOfOwnership.signedPayloadAsString,
        }
      : undefined,
    trails: dto.trails?.map((trail) => ({
      id: trail.id,
      sharedAddressId: trail.sharedAddressID,
      addressStatus: trail.addressStatus,
      comment: trail.comment,
      createdAt: trail.createdAt,
    })),
    originCreationDate: dto.originCreationDate,
    originDeletionDate: dto.originDeletionDate,
    targetAcceptedAt: dto.targetAcceptedAt,
    createdAt: dto.createdAt,
    updatedAt: dto.updatedAt,
  };
}

/**
 * Maps a DTO to a SharedAsset.
 */
function sharedAssetFromDto(dto?: TgvalidatordTnSharedAsset): SharedAsset | undefined {
  if (!dto) {
    return undefined;
  }

  return {
    id: dto.id ?? '',
    whitelistedContractAddressId: dto.wlContractAddressID,
    ownerParticipantId: dto.ownerParticipantId,
    targetParticipantId: dto.targetParticipantId,
    blockchain: dto.blockchain,
    network: dto.network,
    name: dto.name,
    symbol: dto.symbol,
    decimals: dto.decimals,
    contractAddress: dto.contractAddress,
    tokenId: dto.tokenId,
    kind: dto.kind,
    status: dto.status,
    trails: dto.trails?.map((trail) => ({
      id: trail.id,
      sharedAssetId: trail.sharedAssetID,
      assetStatus: trail.assetStatus,
      comment: trail.comment,
      createdAt: trail.createdAt,
    })),
    originCreationDate: dto.originCreationDate,
    originDeletionDate: dto.originDeletionDate,
    targetAcceptedAt: dto.targetAcceptedAt,
    targetRejectedAt: dto.targetRejectedAt,
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
 * Service for Taurus Network address and asset sharing operations.
 *
 * Provides methods to share and unshare addresses and assets between
 * Taurus Network participants.
 *
 * @example
 * ```typescript
 * // List shared addresses
 * const { sharedAddresses, pagination } = await sharingService.listSharedAddresses({
 *   ownerParticipantId: myParticipantId,
 * });
 *
 * // Share an address with another participant
 * await sharingService.shareAddress({
 *   internalAddressId: 'addr-123',
 *   targetParticipantId: 'part-456',
 * });
 * ```
 */
export class SharingService extends BaseService {
  private readonly sharedAddressAssetApi: TaurusNetworkSharedAddressAssetApi;

  /**
   * Creates a new SharingService instance.
   *
   * @param sharedAddressAssetApi - The TaurusNetworkSharedAddressAssetApi instance from the OpenAPI client
   */
  constructor(sharedAddressAssetApi: TaurusNetworkSharedAddressAssetApi) {
    super();
    this.sharedAddressAssetApi = sharedAddressAssetApi;
  }

  // =========================================================================
  // Shared Address Methods
  // =========================================================================

  /**
   * Lists shared addresses with optional filtering.
   *
   * Returns addresses that are shared to or from the current participant.
   *
   * @param options - Optional filtering and pagination options
   * @returns Shared addresses list and pagination info
   * @throws {@link APIError} If API request fails
   */
  async listSharedAddresses(
    options?: ListSharedAddressesOptions
  ): Promise<{ sharedAddresses: SharedAddress[]; pagination?: CursorPagination }> {
    return this.execute(async () => {
      const response = await this.sharedAddressAssetApi.taurusNetworkServiceGetSharedAddresses({
        participantID: options?.participantId,
        ownerParticipantID: options?.ownerParticipantId,
        targetParticipantID: options?.targetParticipantId,
        blockchain: options?.blockchain,
        network: options?.network,
        ids: options?.ids,
        statuses: options?.statuses,
        sortOrder: options?.sortOrder,
        cursorCurrentPage: options?.currentPage,
        cursorPageRequest: options?.pageRequest,
        cursorPageSize: options?.pageSize ? String(options.pageSize) : undefined,
      });

      const sharedAddresses: SharedAddress[] = [];
      if (response.sharedAddresses) {
        for (const dto of response.sharedAddresses) {
          const sharedAddress = sharedAddressFromDto(dto);
          if (sharedAddress) {
            sharedAddresses.push(sharedAddress);
          }
        }
      }

      return {
        sharedAddresses,
        pagination: extractCursorPagination(response.cursor),
      };
    });
  }

  /**
   * Shares an internal address with a Taurus Network participant.
   *
   * This will automatically create a whitelisted address on the target
   * participant's side for them to approve/reject.
   *
   * @param request - Address sharing parameters
   * @throws {@link ValidationError} If required fields are missing
   * @throws {@link APIError} If API request fails
   */
  async shareAddress(request: ShareAddressRequest): Promise<void> {
    if (!request.addressId || request.addressId.trim() === '') {
      throw new ValidationError('addressId is required');
    }
    if (!request.toParticipantId || request.toParticipantId.trim() === '') {
      throw new ValidationError('toParticipantId is required');
    }

    return this.execute(async () => {
      await this.sharedAddressAssetApi.taurusNetworkServiceShareAddress({
        body: {
          addressID: request.addressId,
          toParticipantID: request.toParticipantId,
          keyValueAttributes: request.keyValueAttributes?.map((attr) => ({
            key: attr.key,
            value: attr.value,
          })),
        },
      });
    });
  }

  /**
   * Unshares an address from a Taurus Network participant.
   *
   * This updates the shared address status to unshared. It does not
   * delete the shared address from the registry.
   *
   * @param sharedAddressId - The shared address ID to unshare
   * @throws {@link ValidationError} If sharedAddressId is empty
   * @throws {@link APIError} If API request fails
   */
  async unshareAddress(sharedAddressId: string): Promise<void> {
    if (!sharedAddressId || sharedAddressId.trim() === '') {
      throw new ValidationError('sharedAddressId is required');
    }

    return this.execute(async () => {
      await this.sharedAddressAssetApi.taurusNetworkServiceUnshareAddress({
        tnSharedAddressID: sharedAddressId,
        body: {},
      });
    });
  }

  // =========================================================================
  // Shared Asset Methods
  // =========================================================================

  /**
   * Lists shared assets with optional filtering.
   *
   * Returns whitelisted assets that are shared to or from the current participant.
   *
   * @param options - Optional filtering and pagination options
   * @returns Shared assets list and pagination info
   * @throws {@link APIError} If API request fails
   */
  async listSharedAssets(
    options?: ListSharedAssetsOptions
  ): Promise<{ sharedAssets: SharedAsset[]; pagination?: CursorPagination }> {
    return this.execute(async () => {
      const response = await this.sharedAddressAssetApi.taurusNetworkServiceGetSharedAssets({
        participantID: options?.participantId,
        ownerParticipantID: options?.ownerParticipantId,
        targetParticipantID: options?.targetParticipantId,
        blockchain: options?.blockchain,
        network: options?.network,
        ids: options?.ids,
        statuses: options?.statuses,
        sortOrder: options?.sortOrder,
        cursorCurrentPage: options?.currentPage,
        cursorPageRequest: options?.pageRequest,
        cursorPageSize: options?.pageSize ? String(options.pageSize) : undefined,
      });

      const sharedAssets: SharedAsset[] = [];
      if (response.sharedAssets) {
        for (const dto of response.sharedAssets) {
          const sharedAsset = sharedAssetFromDto(dto);
          if (sharedAsset) {
            sharedAssets.push(sharedAsset);
          }
        }
      }

      return {
        sharedAssets,
        pagination: extractCursorPagination(response.cursor),
      };
    });
  }

  /**
   * Shares a whitelisted asset with a Taurus Network participant.
   *
   * This will automatically create a whitelisted asset on the target
   * participant's side for them to approve/reject.
   *
   * @param request - Asset sharing parameters
   * @throws {@link ValidationError} If required fields are missing
   * @throws {@link APIError} If API request fails
   */
  async shareWhitelistedAsset(request: ShareWhitelistedAssetRequest): Promise<void> {
    if (!request.whitelistedContractId || request.whitelistedContractId.trim() === '') {
      throw new ValidationError('whitelistedContractId is required');
    }
    if (!request.toParticipantId || request.toParticipantId.trim() === '') {
      throw new ValidationError('toParticipantId is required');
    }

    return this.execute(async () => {
      await this.sharedAddressAssetApi.taurusNetworkServiceShareWhitelistedAsset({
        body: {
          whitelistedContractID: request.whitelistedContractId,
          toParticipantID: request.toParticipantId,
        },
      });
    });
  }

  /**
   * Unshares an asset from a Taurus Network participant.
   *
   * This updates the shared asset status to unshared. It does not
   * delete the shared asset from the registry.
   *
   * @param sharedAssetId - The shared asset ID to unshare
   * @throws {@link ValidationError} If sharedAssetId is empty
   * @throws {@link APIError} If API request fails
   */
  async unshareWhitelistedAsset(sharedAssetId: string): Promise<void> {
    if (!sharedAssetId || sharedAssetId.trim() === '') {
      throw new ValidationError('sharedAssetId is required');
    }

    return this.execute(async () => {
      await this.sharedAddressAssetApi.taurusNetworkServiceUnshareWhitelistedAsset({
        tnSharedAssetID: sharedAssetId,
        body: {},
      });
    });
  }
}
