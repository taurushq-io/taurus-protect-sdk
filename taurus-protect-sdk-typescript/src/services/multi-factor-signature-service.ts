/**
 * Multi-factor signature service for Taurus-PROTECT SDK.
 *
 * Provides methods for managing multi-factor signature operations.
 * Multi-factor signatures are used for operations that require approval from
 * multiple parties, such as critical governance changes, high-value transactions,
 * or sensitive administrative actions.
 */

import { NotFoundError, ValidationError } from '../errors';
import type { MultiFactorSignatureApi } from '../internal/openapi/apis/MultiFactorSignatureApi';
import { TgvalidatordMultiFactorSignaturesEntityType } from '../internal/openapi/models/TgvalidatordMultiFactorSignaturesEntityType';
import { multiFactorSignatureInfoFromDto } from '../mappers/multi-factor-signature';
import type {
  ApproveMultiFactorSignatureRequest,
  CreateMultiFactorSignatureRequest,
  MultiFactorSignatureEntityType,
  MultiFactorSignatureInfo,
  RejectMultiFactorSignatureRequest,
} from '../models/multi-factor-signature';
import { BaseService } from './base';

/**
 * Converts SDK entity type to OpenAPI entity type.
 */
function toOpenApiEntityType(
  entityType: MultiFactorSignatureEntityType
): TgvalidatordMultiFactorSignaturesEntityType {
  switch (entityType) {
    case 'REQUEST':
      return TgvalidatordMultiFactorSignaturesEntityType.Request;
    case 'WHITELISTED_ADDRESS':
      return TgvalidatordMultiFactorSignaturesEntityType.WhitelistedAddress;
    case 'WHITELISTED_CONTRACT':
      return TgvalidatordMultiFactorSignaturesEntityType.WhitelistedContract;
    default:
      return TgvalidatordMultiFactorSignaturesEntityType.Request;
  }
}

/**
 * Service for multi-factor signature operations.
 *
 * Multi-factor signatures provide additional security for high-value or
 * sensitive transactions by requiring multiple authentication factors.
 * This service allows creating, approving, and rejecting multi-factor
 * signature workflows for requests, whitelisted addresses, and contracts.
 *
 * @example
 * ```typescript
 * // Create a multi-factor signature for requests
 * const mfsId = await multiFactorSignatureService.create({
 *   entityType: MultiFactorSignatureEntityType.REQUEST,
 *   entityIds: ['123', '456'],
 * });
 *
 * // Get the signature info with payloads to sign
 * const info = await multiFactorSignatureService.get(mfsId);
 *
 * // Approve with signature
 * await multiFactorSignatureService.approve({
 *   id: mfsId,
 *   signature: 'base64EncodedSignature',
 *   comment: 'Approved via SDK',
 * });
 * ```
 */
export class MultiFactorSignatureService extends BaseService {
  private readonly multiFactorSignatureApi: MultiFactorSignatureApi;

  /**
   * Creates a new MultiFactorSignatureService instance.
   *
   * @param multiFactorSignatureApi - The MultiFactorSignatureApi instance from the OpenAPI client
   */
  constructor(multiFactorSignatureApi: MultiFactorSignatureApi) {
    super();
    this.multiFactorSignatureApi = multiFactorSignatureApi;
  }

  /**
   * Get multi-factor signature entity info by ID.
   *
   * Returns the signature info including payloads to sign for approval.
   *
   * @param id - The multi-factor signature ID.
   * @returns The multi-factor signature info.
   * @throws {@link ValidationError} If ID is empty.
   * @throws {@link NotFoundError} If not found.
   * @throws {@link APIError} If API request fails.
   */
  async get(id: string): Promise<MultiFactorSignatureInfo> {
    if (!id || id.trim() === '') {
      throw new ValidationError('id is required');
    }

    return this.execute(async () => {
      const response =
        await this.multiFactorSignatureApi.multiFactorSignatureServiceGetMultiFactorSignatureEntitiesInfo(
          { id }
        );
      const result = multiFactorSignatureInfoFromDto(response);
      if (!result) {
        throw new NotFoundError(`Multi-factor signature ${id} not found`);
      }
      return result;
    });
  }

  /**
   * Create a multi-factor signature batch.
   *
   * Creates a signature process requiring multi-factor approvals for the
   * specified entities. Returns an ID that must be used by another factor
   * device to sign.
   *
   * @param request - The create request with entity type and IDs.
   * @returns The created multi-factor signature ID.
   * @throws {@link ValidationError} If request is invalid.
   * @throws {@link APIError} If API request fails.
   */
  async create(request: CreateMultiFactorSignatureRequest): Promise<string> {
    if (!request.entityIds || request.entityIds.length === 0) {
      throw new ValidationError('entityIds cannot be empty');
    }
    if (!request.entityType) {
      throw new ValidationError('entityType is required');
    }

    return this.execute(async () => {
      const response =
        await this.multiFactorSignatureApi.multiFactorSignatureServiceCreateMultiFactorSignatureBatch(
          {
            body: {
              entityType: toOpenApiEntityType(request.entityType),
              entityIDs: request.entityIds,
            },
          }
        );
      return response.id ?? '';
    });
  }

  /**
   * Approve a multi-factor signature.
   *
   * Approves the entities previously associated using a multi-factor approach.
   * Requires either 'RequestMobileAppSigner' or 'WhitelistedAddressMobileAppSigner'
   * role depending on the targeted resource.
   *
   * @param request - The approval request with ID, signature, and comment.
   * @throws {@link ValidationError} If request is invalid.
   * @throws {@link APIError} If API request fails.
   */
  async approve(request: ApproveMultiFactorSignatureRequest): Promise<void> {
    if (!request.id || request.id.trim() === '') {
      throw new ValidationError('id is required');
    }
    if (!request.signature || request.signature.trim() === '') {
      throw new ValidationError('signature is required');
    }

    await this.execute(async () => {
      await this.multiFactorSignatureApi.multiFactorSignatureServiceApproveMultiFactorSignature(
        {
          id: request.id,
          body: {
            signature: request.signature,
            comment: request.comment ?? '',
          },
        }
      );
    });
  }

  /**
   * Reject a multi-factor signature.
   *
   * Rejects the entities previously associated. Requires either
   * 'RequestMobileAppSigner' or 'WhitelistedAddressMobileAppSigner' role
   * depending on the targeted resource.
   *
   * @param request - The rejection request with ID and comment.
   * @throws {@link ValidationError} If request is invalid.
   * @throws {@link APIError} If API request fails.
   */
  async reject(request: RejectMultiFactorSignatureRequest): Promise<void> {
    if (!request.id || request.id.trim() === '') {
      throw new ValidationError('id is required');
    }

    await this.execute(async () => {
      await this.multiFactorSignatureApi.multiFactorSignatureServiceRejectMultiFactorSignature(
        {
          id: request.id,
          body: {
            comment: request.comment ?? '',
          },
        }
      );
    });
  }
}
