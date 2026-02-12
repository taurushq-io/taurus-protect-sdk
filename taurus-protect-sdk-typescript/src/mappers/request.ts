/**
 * Request mapper for converting OpenAPI DTOs to domain models.
 */

import type { TgvalidatordRequest } from "../internal/openapi/models/TgvalidatordRequest";
import type { TgvalidatordMetadata } from "../internal/openapi/models/TgvalidatordMetadata";
import type { TgvalidatordCurrency } from "../internal/openapi/models/TgvalidatordCurrency";
import type { TgvalidatordApprovers } from "../internal/openapi/models/TgvalidatordApprovers";
import type { TgvalidatordParallelApproversGroups } from "../internal/openapi/models/TgvalidatordParallelApproversGroups";
import type { TgvalidatordApproversGroup } from "../internal/openapi/models/TgvalidatordApproversGroup";
import type { RequestSignedRequest } from "../internal/openapi/models/RequestSignedRequest";
import type {
  Request,
  RequestMetadata,
  SignedRequest,
  CurrencyInfo,
  Approvers,
  ApproversGroup,
  ParallelApproversGroups,
  RequestStatus,
  RequestType,
} from "../models/request";
import { safeIntDefault, safeStringDefault, safeString, safeInt, safeDate } from "./base";

/**
 * Maps a TgvalidatordRequest DTO to a Request domain model.
 *
 * @param dto - The OpenAPI DTO
 * @returns The domain model, or undefined if dto is null/undefined
 */
export function requestFromDto(dto: TgvalidatordRequest | null | undefined): Request | undefined {
  if (!dto) {
    return undefined;
  }

  const metadata = mapMetadata(dto.metadata);

  return {
    id: safeIntDefault(dto.id, 0),
    type: (dto.type ?? "UNKNOWN") as RequestType | string,
    status: (dto.status ?? "UNKNOWN") as RequestStatus | string,
    metadata,
    tenantId: safeString(dto.tenantId),
    currency: safeString(dto.currency),
    currencyInfo: mapCurrencyInfo(dto.currencyInfo),
    memo: undefined, // Not directly available in DTO
    rule: safeString(dto.rule),
    externalRequestId: safeString(dto.externalRequestId),
    requestBundleId: safeString(dto.requestBundleId),
    needsApprovalFrom: dto.needsApprovalFrom ?? [],
    approvers: mapApprovers(dto.approvers),
    createdAt: safeDate(dto.creationDate),
    updatedAt: safeDate(dto.updateDate),
    tags: [], // Tags not directly available in DTO
    signedRequests: signedRequestsFromDto(dto.signedRequests),
  };
}

/**
 * Maps an array of TgvalidatordRequest DTOs to Request domain models.
 *
 * @param dtos - The array of OpenAPI DTOs
 * @returns Array of domain models (null/undefined entries are filtered out)
 */
export function requestsFromDto(
  dtos: TgvalidatordRequest[] | null | undefined
): Request[] {
  if (!dtos) {
    return [];
  }
  const result: Request[] = [];
  for (const dto of dtos) {
    const mapped = requestFromDto(dto);
    if (mapped) {
      result.push(mapped);
    }
  }
  return result;
}

/**
 * Maps metadata DTO to domain model.
 *
 * SECURITY: payload field intentionally not mapped - use payloadAsString only.
 * The raw payload object could be tampered with while payloadAsString
 * remains unchanged (hash still verifies). By extracting from payloadAsString,
 * we ensure all data comes from the cryptographically verified source.
 */
function mapMetadata(
  dto: TgvalidatordMetadata | null | undefined
): RequestMetadata | undefined {
  if (!dto) {
    return undefined;
  }

  const hash = safeString(dto.hash);
  if (!hash) {
    return undefined;
  }

  return {
    hash,
    payloadAsString: safeString(dto.payloadAsString),
    // SECURITY: payload intentionally not mapped
  };
}

/**
 * Maps currency info DTO to domain model.
 */
function mapCurrencyInfo(
  dto: TgvalidatordCurrency | null | undefined
): CurrencyInfo | undefined {
  if (!dto) {
    return undefined;
  }

  return {
    id: safeString(dto.id),
    symbol: safeString(dto.symbol),
    name: safeString(dto.name),
    decimals: safeInt(dto.decimals),
    blockchain: safeString(dto.blockchain),
    network: safeString(dto.network),
  };
}

/**
 * Maps approvers DTO to domain model.
 */
function mapApprovers(
  dto: TgvalidatordApprovers | null | undefined
): Approvers | undefined {
  if (!dto) {
    return undefined;
  }

  return {
    parallel: mapParallelApproversGroups(dto.parallel),
  };
}

/**
 * Maps parallel approvers groups.
 */
function mapParallelApproversGroups(
  dtos: TgvalidatordParallelApproversGroups[] | null | undefined
): ParallelApproversGroups[] {
  if (!dtos) {
    return [];
  }
  return dtos.map((dto) => ({
    sequential: mapApproversGroups(dto.sequential),
  }));
}

/**
 * Maps approvers groups.
 */
function mapApproversGroups(
  dtos: TgvalidatordApproversGroup[] | null | undefined
): ApproversGroup[] {
  if (!dtos) {
    return [];
  }
  return dtos.map((dto) => ({
    externalGroupId: safeString(dto.externalGroupID),
    minimumSignatures: dto.minimumSignatures,
  }));
}

/**
 * Maps a RequestSignedRequest DTO to a SignedRequest domain model.
 */
function signedRequestFromDto(dto: RequestSignedRequest): SignedRequest {
  return {
    id: safeStringDefault(dto.id, ""),
    signedRequest: safeStringDefault(dto.signedRequest, ""),
    status: (dto.status ?? "UNKNOWN") as RequestStatus | string,
    hash: safeStringDefault(dto.hash, ""),
    block: safeIntDefault(dto.block, 0),
    details: safeStringDefault(dto.details, ""),
    creationDate: safeDate(dto.creationDate),
    updateDate: safeDate(dto.updateDate),
    broadcastDate: safeDate(dto.broadcastDate),
    confirmationDate: safeDate(dto.confirmationDate),
  };
}

/**
 * Maps an array of RequestSignedRequest DTOs to SignedRequest domain models.
 */
function signedRequestsFromDto(dtos: RequestSignedRequest[] | null | undefined): SignedRequest[] {
  if (!dtos || dtos.length === 0) {
    return [];
  }
  return dtos.map(signedRequestFromDto);
}
