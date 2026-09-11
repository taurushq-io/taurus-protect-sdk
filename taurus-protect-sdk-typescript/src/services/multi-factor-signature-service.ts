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
   * SECURITY — `payloadToSign` IS UNVERIFIED SERVER DATA. The caller must bind it to a
   * verified entity before signing it.
   *
   * This is the one read path in the SDK that returns bytes intended for a signing key
   * without verifying them, and it is deliberate-but-unresolved rather than an oversight
   * (see TODOS.md, "MultiFactorSignatureService.approve takes an opaque caller signature").
   * The reason it cannot be fixed here: the reply carries only
   * `{id, payloadToSign[], entityType}`. There is **no entity id**, singular or plural, and
   * `payloadToSign` is a bare string array with no per-element id or hash — so nothing in the
   * reply can be joined back to the entities the request was created for, and this method has
   * nothing to check against.
   *
   * What that means for a second-factor signer: a malicious or compromised API server can
   * answer this call with the metadata hash of an entity of its choosing — a withdrawal to an
   * attacker address, an attacker-controlled whitelisted address — under the expected kind.
   * Sign it and the server holds a valid MobileAppSigner approval over an entity nobody
   * reviewed. That signature is precisely the artefact a compromised server cannot forge on
   * its own. Contrast {@link RequestService.approveRequests}, which re-verifies every hash
   * against its `payloadAsString` and throws `IntegrityError` before the approver's key signs.
   *
   * The workable client-side binding, for when the decision lands: {@link create} takes the
   * entity IDs, so a caller who keeps them can re-read those entities through the verifying
   * reader for that `entityType` and require each `payloadToSign` element to equal a locally
   * recomputed, verified metadata hash. Until the SDK does that for you, do it yourself.
   *
   * @param id - The multi-factor signature ID.
   * @returns The multi-factor signature info, whose `payloadToSign` is UNVERIFIED.
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
   * SECURITY — THE SIGNATURE IS OPAQUE TO THIS SDK. It is forwarded after a non-empty check
   * and nothing here knows or checks what it covers, so this method cannot enforce that only
   * verified payloads are approved.
   *
   * That is why this site is absent from `scripts/signing-sites/manifest.json`: the manifest
   * enumerates sites where the *SDK* signs, and here the caller does. It is also why its
   * absence is not reassuring — the gate that would have forced a `verifies` /
   * `signs-own-bytes` classification simply does not see this path.
   *
   * The caller is responsible for having bound the bytes they signed to a verified entity.
   * See {@link get} for why the SDK cannot do it for them yet, and TODOS.md for the open
   * decision.
   *
   * @param request - The approval request with ID, signature, and comment. The signature is
   *   assumed to cover bytes the CALLER verified; this method does not check it.
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
