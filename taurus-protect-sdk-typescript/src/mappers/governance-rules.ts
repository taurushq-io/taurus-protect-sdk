/**
 * Mappers for decoding governance rules from base64-encoded data.
 */

import type { TgvalidatordRules } from "../internal/openapi/models/TgvalidatordRules";
import type { TgvalidatordRulesTrail } from "../internal/openapi/models/TgvalidatordRulesTrail";
import type { TgvalidatordRuleUserSignature } from "../internal/openapi/models/TgvalidatordRuleUserSignature";
import {
  type DecodedRulesContainer,
  type RuleUser,
  type RuleGroup,
  type GroupThreshold,
  type SequentialThresholds,
  type TransactionRules,
  type AddressWhitelistingRules,
  type ContractAddressWhitelistingRules,
  type RuleUserSignature,
  type RulesTrail,
  type GovernanceRules,
  createEmptyRulesContainer,
} from "../models/governance-rules";
import { safeBoolDefault, safeDate, safeMap, safeString } from "./base";
import { IntegrityError } from "../errors";
import { tryDecodeProtobufRulesContainer } from "./protobuf-rules-container";

/**
 * Decodes a base64-encoded rules container.
 *
 * Supports both protobuf (primary format, matches Java SDK) and JSON formats.
 * Tries protobuf first, then falls back to JSON parsing.
 *
 * @param base64Data - Base64-encoded rules container
 * @returns Decoded rules container
 * @throws IntegrityError if decoding fails
 */
export function rulesContainerFromBase64(base64Data: string): DecodedRulesContainer {
  if (!base64Data) {
    return createEmptyRulesContainer();
  }

  // Decode base64 to bytes
  let decoded: Uint8Array;
  try {
    decoded = Buffer.from(base64Data, 'base64');
  } catch (error) {
    // Invalid base64 encoding - this is a security-critical failure
    throw new IntegrityError(
      `Failed to decode rules container: invalid base64 encoding - ${error instanceof Error ? error.message : "unknown error"}`
    );
  }

  // Try protobuf first (primary format, matches Java SDK)
  const protobufResult = tryDecodeProtobufRulesContainer(decoded);
  if (protobufResult) {
    return protobufResult;
  }

  // Fall back to JSON parsing
  try {
    const jsonString = new TextDecoder().decode(decoded);
    const data = JSON.parse(jsonString) as Record<string, unknown>;
    return parseRulesContainerFromDict(data);
  } catch {
    // Not valid protobuf or JSON - this is a security-critical failure
    throw new IntegrityError(
      "Failed to decode rules container: not valid protobuf or JSON"
    );
  }
}

/**
 * Parses rules container from a dictionary (JSON object).
 */
function parseRulesContainerFromDict(data: Record<string, unknown>): DecodedRulesContainer {
  // Parse users
  const users: RuleUser[] = [];
  const usersData = (data['users'] ?? []) as Record<string, unknown>[];
  for (const userData of usersData) {
    users.push({
      id: getString(userData, 'id'),
      name: getString(userData, 'name'),
      publicKeyPem: getString(userData, 'publicKeyPem') ?? getString(userData, 'public_key_pem') ?? getString(userData, 'publicKey'),
      roles: getStringArray(userData, 'roles'),
    });
  }

  // Parse groups
  const groups: RuleGroup[] = [];
  const groupsData = (data['groups'] ?? []) as Record<string, unknown>[];
  for (const groupData of groupsData) {
    groups.push({
      id: getString(groupData, 'id'),
      name: getString(groupData, 'name'),
      userIds: getStringArray(groupData, 'userIds') ?? getStringArray(groupData, 'user_ids'),
    });
  }

  // Parse transaction rules
  const transactionRules: TransactionRules[] = [];
  const transactionRulesData = (
    data['transactionRules'] ??
    data['transaction_rules'] ??
    []
  ) as Record<string, unknown>[];
  for (const ruleData of transactionRulesData) {
    const parallelThresholds = parseSequentialThresholds(
      getArray(ruleData, 'parallelThresholds') ?? getArray(ruleData, 'parallel_thresholds') ?? []
    );
    transactionRules.push({ parallelThresholds });
  }

  // Parse address whitelisting rules
  const addressWhitelistingRules: AddressWhitelistingRules[] = [];
  const addressRulesData = (
    data['addressWhitelistingRules'] ??
    data['address_whitelisting_rules'] ??
    []
  ) as Record<string, unknown>[];
  for (const ruleData of addressRulesData) {
    const parallelThresholds = parseSequentialThresholds(
      getArray(ruleData, 'parallelThresholds') ?? getArray(ruleData, 'parallel_thresholds') ?? []
    );
    addressWhitelistingRules.push({
      currency: getString(ruleData, 'currency'),
      network: getString(ruleData, 'network'),
      parallelThresholds,
    });
  }

  // Parse contract address whitelisting rules
  const contractAddressWhitelistingRules: ContractAddressWhitelistingRules[] = [];
  const contractRulesData = (
    data['contractAddressWhitelistingRules'] ??
    data['contract_address_whitelisting_rules'] ??
    []
  ) as Record<string, unknown>[];
  for (const ruleData of contractRulesData) {
    const parallelThresholds = parseSequentialThresholds(
      getArray(ruleData, 'parallelThresholds') ?? getArray(ruleData, 'parallel_thresholds') ?? []
    );
    contractAddressWhitelistingRules.push({
      blockchain: getString(ruleData, 'blockchain'),
      network: getString(ruleData, 'network'),
      parallelThresholds,
    });
  }

  return {
    users,
    groups,
    minimumDistinctUserSignatures:
      getNumber(data, 'minimumDistinctUserSignatures') ??
      getNumber(data, 'minimum_distinct_user_signatures') ??
      0,
    minimumDistinctGroupSignatures:
      getNumber(data, 'minimumDistinctGroupSignatures') ??
      getNumber(data, 'minimum_distinct_group_signatures') ??
      0,
    transactionRules,
    addressWhitelistingRules,
    contractAddressWhitelistingRules,
    enforcedRulesHash:
      getString(data, 'enforcedRulesHash') ?? getString(data, 'enforced_rules_hash'),
    timestamp: getNumber(data, 'timestamp') ?? 0,
    hsmSlotId:
      getNumber(data, 'hsmSlotId') ?? getNumber(data, 'hsm_slot_id') ?? 0,
    minimumCommitmentSignatures:
      getNumber(data, 'minimumCommitmentSignatures') ??
      getNumber(data, 'minimum_commitment_signatures') ??
      0,
    engineIdentities:
      data['engineIdentities'] !== undefined
        ? getStringArray(data, 'engineIdentities')
        : getStringArray(data, 'engine_identities'),
  };
}

/**
 * Parses group thresholds from list of dicts.
 */
function parseGroupThresholds(data: Record<string, unknown>[]): GroupThreshold[] {
  const thresholds: GroupThreshold[] = [];
  for (const item of data) {
    thresholds.push({
      groupId: getString(item, 'groupId') ?? getString(item, 'group_id'),
      minimumSignatures:
        getNumber(item, 'minimumSignatures') ?? getNumber(item, 'minimum_signatures') ?? 0,
      threshold: getNumber(item, 'threshold') ?? 0,
    });
  }
  return thresholds;
}

/**
 * Parses sequential thresholds from list of dicts.
 * Handles both nested format ({"thresholds": [...]}) and flat format
 * ({"groupId": ..., "minimumSignatures": ...}) for backward compatibility.
 */
function parseSequentialThresholds(data: Record<string, unknown>[]): SequentialThresholds[] {
  const result: SequentialThresholds[] = [];
  for (const item of data) {
    const thresholdsData = getArray(item, 'thresholds');
    if (thresholdsData) {
      // Nested format: {"thresholds": [{groupId, minimumSignatures}, ...]}
      const thresholds = parseGroupThresholds(thresholdsData);
      result.push({ thresholds });
    } else if (item['groupId'] !== undefined || item['group_id'] !== undefined || item['minimumSignatures'] !== undefined || item['minimum_signatures'] !== undefined) {
      // Flat format: {groupId, minimumSignatures} - wrap in SequentialThresholds
      const threshold = parseGroupThresholds([item]);
      result.push({ thresholds: threshold });
    }
  }
  return result;
}

/**
 * Decodes base64-encoded user signatures.
 *
 * Note: Full protobuf decoding is not yet implemented. This function
 * attempts to parse as JSON if the data is JSON-encoded.
 *
 * @param base64Data - Base64-encoded user signatures
 * @returns List of user signatures
 */
export function userSignaturesFromBase64(base64Data: string): RuleUserSignature[] {
  if (!base64Data) {
    return [];
  }

  const decodedBytes = Buffer.from(base64Data, 'base64');

  // Try protobuf first (primary format, matches Java/Go SDKs)
  try {
    const { UserSignatures: PbUserSignatures } = require('../internal/proto/request_reply') as typeof import('../internal/proto/request_reply');
    const pbSigs = PbUserSignatures.decode(decodedBytes);
    if (pbSigs.signatures && pbSigs.signatures.length > 0) {
      return pbSigs.signatures.map((sig) => ({
        userId: sig.userId || undefined,
        // Protobuf signature is raw bytes — re-encode as base64 (matches Go SDK)
        signature: sig.signature
          ? Buffer.from(sig.signature).toString('base64')
          : undefined,
      }));
    }
  } catch {
    // Not valid protobuf — fall through to JSON
  }

  // Fall back to JSON parsing
  try {
    const decoded = decodedBytes.toString('utf-8');
    const data = JSON.parse(decoded) as unknown;
    const signatures: RuleUserSignature[] = [];

    // Handle both array and object with signatures field
    const sigList = Array.isArray(data)
      ? data
      : ((data as Record<string, unknown>)['signatures'] ?? []);

    for (const sigData of sigList as Record<string, unknown>[]) {
      signatures.push({
        userId: getString(sigData, 'userId') ?? getString(sigData, 'user_id'),
        signature: getString(sigData, 'signature'),
      });
    }

    return signatures;
  } catch {
    return [];
  }
}

// Helper functions for safe type access

function getString(obj: Record<string, unknown>, key: string): string | undefined {
  const value = obj[key];
  return typeof value === 'string' ? value : undefined;
}

function getNumber(obj: Record<string, unknown>, key: string): number | undefined {
  const value = obj[key];
  return typeof value === 'number' ? value : undefined;
}

function getArray(obj: Record<string, unknown>, key: string): Record<string, unknown>[] | undefined {
  const value = obj[key];
  return Array.isArray(value) ? (value as Record<string, unknown>[]) : undefined;
}

function getStringArray(obj: Record<string, unknown>, key: string): string[] {
  const value = obj[key];
  if (!Array.isArray(value)) return [];
  return value.filter((v): v is string => typeof v === "string");
}

/**
 * Maps a TgvalidatordRuleUserSignature DTO to a RuleUserSignature domain model.
 *
 * @param dto - The rule user signature DTO from the OpenAPI response
 * @returns The RuleUserSignature domain model, or undefined if dto is null/undefined
 */
export function ruleUserSignatureFromDto(
  dto: TgvalidatordRuleUserSignature | null | undefined
): RuleUserSignature | undefined {
  if (!dto) {
    return undefined;
  }
  return {
    userId: safeString(dto.userId),
    signature: safeString(dto.signature),
  };
}

/**
 * Maps a TgvalidatordRulesTrail DTO to a RulesTrail domain model.
 *
 * @param dto - The rules trail DTO from the OpenAPI response
 * @returns The RulesTrail domain model, or undefined if dto is null/undefined
 */
export function rulesTrailFromDto(
  dto: TgvalidatordRulesTrail | null | undefined
): RulesTrail | undefined {
  if (!dto) {
    return undefined;
  }
  return {
    userId: safeString(dto.userId),
    action: safeString(dto.action),
    timestamp: safeDate(dto.date),
  };
}

/**
 * Maps a TgvalidatordRules DTO to a GovernanceRules domain model.
 *
 * @param dto - The rules DTO from the OpenAPI response
 * @returns The GovernanceRules domain model, or undefined if dto is null/undefined
 */
export function governanceRulesFromDto(
  dto: TgvalidatordRules | null | undefined
): GovernanceRules | undefined {
  if (!dto) {
    return undefined;
  }
  return {
    rulesContainer: safeString(dto.rulesContainer),
    rulesSignatures: safeMap(dto.rulesSignatures, ruleUserSignatureFromDto),
    locked: safeBoolDefault(dto.locked, false),
    creationDate: safeDate(dto.creationDate),
    updateDate: safeDate(dto.updateDate),
    trails: safeMap(dto.trails, rulesTrailFromDto),
  };
}

/**
 * Maps an array of TgvalidatordRules DTOs to GovernanceRules domain models.
 *
 * @param dtos - The array of rules DTOs
 * @returns Array of GovernanceRules domain models (undefined entries filtered out)
 */
export function governanceRulesArrayFromDto(
  dtos: TgvalidatordRules[] | null | undefined
): GovernanceRules[] {
  return safeMap(dtos, governanceRulesFromDto);
}
