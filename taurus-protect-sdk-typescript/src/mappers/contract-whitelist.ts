/**
 * Contract whitelist mapper functions for converting OpenAPI DTOs to domain models.
 */

import type {
  WhitelistedContract,
  WhitelistedContractAttribute,
  WhitelistedContractMetadata,
  WhitelistedContractSignature,
  SignedWhitelistedContract,
  WhitelistedContractTrail,
  WhitelistedContractApprover,
  WhitelistedContractApprovers,
  WhitelistedContractApproverGroup,
  SignedWhitelistedContractEnvelope,
  WhitelistedContractResult,
} from '../models/contract-whitelist';
import { safeBool, safeInt, safeMap, safeString } from './base';

/**
 * Maps a whitelisted contract attribute from DTO.
 */
export function whitelistedContractAttributeFromDto(
  dto: unknown
): WhitelistedContractAttribute | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    id: safeString(d.id),
    key: safeString(d.key),
    value: safeString(d.value),
    contentType: safeString(d.contentType ?? d.content_type),
    type: safeString(d.type),
    subType: safeString(d.subtype ?? d.subType ?? d.sub_type),
    isFile: safeBool(d.isfile ?? d.isFile ?? d.is_file),
  };
}

/**
 * Maps an array of whitelisted contract attributes from DTOs.
 */
export function whitelistedContractAttributesFromDto(
  dtos: unknown[] | null | undefined
): WhitelistedContractAttribute[] {
  return safeMap(dtos, whitelistedContractAttributeFromDto);
}

/**
 * Maps whitelisted contract metadata from DTO.
 *
 * SECURITY: payload field intentionally not mapped - use payloadAsString only.
 * The raw payload object could be tampered with while payloadAsString
 * remains unchanged (hash still verifies).
 */
export function whitelistedContractMetadataFromDto(
  dto: unknown
): WhitelistedContractMetadata | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    hash: safeString(d.hash),
    payloadAsString: safeString(d.payloadAsString ?? d.payload_as_string),
    // SECURITY: payload intentionally not mapped
  };
}

/**
 * Maps a whitelisted contract signature from DTO.
 */
export function whitelistedContractSignatureFromDto(
  dto: unknown
): WhitelistedContractSignature | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  const hashes = Array.isArray(d.hashes)
    ? (d.hashes as unknown[]).map((h) => String(h))
    : [];

  return {
    signature: safeString(d.signature),
    comment: safeString(d.comment),
    hashes,
    userId: safeString(d.userId ?? d.user_id),
  };
}

/**
 * Maps an array of whitelisted contract signatures from DTOs.
 */
export function whitelistedContractSignaturesFromDto(
  dtos: unknown[] | null | undefined
): WhitelistedContractSignature[] {
  return safeMap(dtos, whitelistedContractSignatureFromDto);
}

/**
 * Maps signed whitelisted contract from DTO.
 */
export function signedWhitelistedContractFromDto(
  dto: unknown
): SignedWhitelistedContract | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    payload: safeString(d.payload),
    signatures: whitelistedContractSignaturesFromDto(d.signatures as unknown[]),
  };
}

/**
 * Maps a whitelisted contract trail from DTO.
 */
export function whitelistedContractTrailFromDto(
  dto: unknown
): WhitelistedContractTrail | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    id: safeString(d.id),
    action: safeString(d.action),
    userId: safeString(d.userId ?? d.user_id),
    comment: safeString(d.comment),
    timestamp: safeString(d.timestamp ?? d.createdAt ?? d.created_at),
  };
}

/**
 * Maps an array of whitelisted contract trails from DTOs.
 */
export function whitelistedContractTrailsFromDto(
  dtos: unknown[] | null | undefined
): WhitelistedContractTrail[] {
  return safeMap(dtos, whitelistedContractTrailFromDto);
}

/**
 * Maps a whitelisted contract approver from DTO.
 */
export function whitelistedContractApproverFromDto(
  dto: unknown
): WhitelistedContractApprover | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    userId: safeString(d.userId ?? d.user_id ?? d.id),
    userName: safeString(d.userName ?? d.user_name ?? d.name ?? d.email),
    pending: safeBool(d.pending),
  };
}

/**
 * Maps an array of whitelisted contract approvers from DTOs.
 */
export function whitelistedContractApproversListFromDto(
  dtos: unknown[] | null | undefined
): WhitelistedContractApprover[] {
  return safeMap(dtos, whitelistedContractApproverFromDto);
}

/**
 * Maps a whitelisted contract approver group from DTO.
 */
export function whitelistedContractApproverGroupFromDto(
  dto: unknown
): WhitelistedContractApproverGroup | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    id: safeString(d.id),
    name: safeString(d.name),
    required: safeInt(d.required),
    users: whitelistedContractApproversListFromDto(d.users as unknown[]),
  };
}

/**
 * Maps an array of whitelisted contract approver groups from DTOs.
 */
export function whitelistedContractApproverGroupsFromDto(
  dtos: unknown[] | null | undefined
): WhitelistedContractApproverGroup[] {
  return safeMap(dtos, whitelistedContractApproverGroupFromDto);
}

/**
 * Maps whitelisted contract approvers from DTO.
 */
export function whitelistedContractApproversFromDto(
  dto: unknown
): WhitelistedContractApprovers | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    required: safeInt(d.required),
    groups: whitelistedContractApproverGroupsFromDto(d.groups as unknown[]),
  };
}

/**
 * Extracts contract details from metadata payload string.
 */
function extractContractDetailsFromPayload(payloadAsString: string | undefined): {
  contractAddress?: string;
  symbol?: string;
  name?: string;
  decimals?: number;
  kind?: string;
  tokenId?: string;
} {
  if (!payloadAsString) {
    return {};
  }

  try {
    const payload = JSON.parse(payloadAsString);
    return {
      contractAddress: safeString(payload.contractAddress ?? payload.contract_address),
      symbol: safeString(payload.symbol),
      name: safeString(payload.name),
      decimals: safeInt(payload.decimals),
      kind: safeString(payload.kind),
      tokenId: safeString(payload.tokenId ?? payload.token_id),
    };
  } catch {
    return {};
  }
}

/**
 * Maps a signed whitelisted contract envelope from DTO.
 */
export function signedWhitelistedContractEnvelopeFromDto(
  dto: unknown
): SignedWhitelistedContractEnvelope | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    id: safeString(d.id) ?? '',
    tenantId: safeString(d.tenantId ?? d.tenant_id),
    blockchain: safeString(d.blockchain),
    network: safeString(d.network),
    metadata: whitelistedContractMetadataFromDto(d.metadata),
    signedContractAddress: signedWhitelistedContractFromDto(
      d.signedContractAddress ?? d.signed_contract_address
    ),
    action: safeString(d.action),
    trails: whitelistedContractTrailsFromDto(d.trails as unknown[]),
    rulesContainer: safeString(d.rulesContainer ?? d.rules_container),
    rule: safeString(d.rule),
    rulesSignatures: safeString(d.rulesSignatures ?? d.rules_signatures),
    approvers: whitelistedContractApproversFromDto(d.approvers),
    attributes: whitelistedContractAttributesFromDto(d.attributes as unknown[]),
    status: safeString(d.status),
    businessRuleEnabled: safeBool(d.businessRuleEnabled ?? d.business_rule_enabled),
  };
}

/**
 * Maps an array of signed whitelisted contract envelopes from DTOs.
 */
export function signedWhitelistedContractEnvelopesFromDto(
  dtos: unknown[] | null | undefined
): SignedWhitelistedContractEnvelope[] {
  return safeMap(dtos, signedWhitelistedContractEnvelopeFromDto);
}

/**
 * Maps a signed whitelisted contract envelope to a simplified WhitelistedContract.
 *
 * Extracts contract details from metadata payload and signed contract payload.
 */
export function whitelistedContractFromEnvelope(
  envelope: SignedWhitelistedContractEnvelope
): WhitelistedContract {
  // Extract details from metadata payload
  const metadataDetails = extractContractDetailsFromPayload(
    envelope.metadata?.payloadAsString
  );

  // Extract details from signed contract payload (fallback)
  const signedDetails = extractContractDetailsFromPayload(
    envelope.signedContractAddress?.payload
  );

  return {
    id: envelope.id,
    tenantId: envelope.tenantId,
    blockchain: envelope.blockchain,
    network: envelope.network,
    contractAddress: metadataDetails.contractAddress ?? signedDetails.contractAddress,
    symbol: metadataDetails.symbol ?? signedDetails.symbol,
    name: metadataDetails.name ?? signedDetails.name,
    decimals: metadataDetails.decimals ?? signedDetails.decimals,
    kind: metadataDetails.kind ?? signedDetails.kind,
    tokenId: metadataDetails.tokenId ?? signedDetails.tokenId,
    status: envelope.status,
    businessRuleEnabled: envelope.businessRuleEnabled,
    attributes: envelope.attributes,
  };
}

/**
 * Maps an envelope DTO directly to a WhitelistedContract.
 */
export function whitelistedContractFromDto(dto: unknown): WhitelistedContract | undefined {
  const envelope = signedWhitelistedContractEnvelopeFromDto(dto);
  if (!envelope) {
    return undefined;
  }
  return whitelistedContractFromEnvelope(envelope);
}

/**
 * Maps an array of envelope DTOs to WhitelistedContracts.
 */
export function whitelistedContractsFromDto(
  dtos: unknown[] | null | undefined
): WhitelistedContract[] {
  return safeMap(dtos, whitelistedContractFromDto);
}

/**
 * Maps a whitelisted contracts reply to a WhitelistedContractResult.
 */
export function whitelistedContractResultFromDto(
  dto: unknown
): WhitelistedContractResult {
  if (!dto || typeof dto !== 'object') {
    return { contracts: [], totalItems: 0 };
  }

  const d = dto as Record<string, unknown>;
  return {
    contracts: whitelistedContractsFromDto(d.result as unknown[]),
    totalItems: safeInt(d.totalItems ?? d.total_items) ?? 0,
  };
}
