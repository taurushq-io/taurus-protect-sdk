/**
 * Mappers for decoding governance rules from base64-encoded data.
 */

import type { TgvalidatordRules } from "../internal/openapi/models/TgvalidatordRules";
import type { TgvalidatordRulesTrail } from "../internal/openapi/models/TgvalidatordRulesTrail";
import type { GetPublicKeysReplyPublicKey } from "../internal/openapi/models/GetPublicKeysReplyPublicKey";
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
  type SuperAdminPublicKey,
  createEmptyRulesContainer,
} from "../models/governance-rules";
import { safeBoolDefault, safeDate, safeMap, safeString } from "./base";
import { IntegrityError } from "../errors";
import { tryDecodeProtobufRulesContainer } from "./protobuf-rules-container";
import { strictBase64Decode } from "../helpers/strict-base64";
import { MAX_RULES_CONTAINER_BYTES } from "../helpers/signed-payload-guard";

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

  // Bound the input BEFORE decoding it. Nothing capped this container, and it is the
  // document every HSM and PRICEUPDATER public key is read from, so the decode runs on
  // the address, asset and price read paths for any response a server chooses to send.
  //
  // The check is on the base64 TEXT first, deliberately: `strictBase64Decode`'s
  // well-formedness regex recurses per four-character group, so a multi-megabyte string
  // exhausts the stack there (it fails closed, but with a misleading "invalid base64"
  // message and after the damage). Four base64 characters carry three bytes, so this
  // bounds the decoded size without allocating anything.
  const maxBase64Length = Math.ceil((MAX_RULES_CONTAINER_BYTES * 4) / 3) + 4;
  if (base64Data.length > maxBase64Length) {
    throw new IntegrityError(
      `Failed to decode rules container: encoded container exceeds ${MAX_RULES_CONTAINER_BYTES} bytes`
    );
  }

  // Decode base64 to bytes
  let decoded: Uint8Array;
  try {
    decoded = strictBase64Decode(base64Data);
  } catch (error) {
    // Invalid base64 encoding - this is a security-critical failure
    throw new IntegrityError(
      `Failed to decode rules container: invalid base64 encoding - ${error instanceof Error ? error.message : "unknown error"}`
    );
  }

  // Whitespace inside the base64 means the text bound above is not tight, so bound the
  // decoded bytes too rather than relying on the estimate.
  if (decoded.length > MAX_RULES_CONTAINER_BYTES) {
    throw new IntegrityError(
      `Failed to decode rules container: container exceeds ${MAX_RULES_CONTAINER_BYTES} bytes`
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
    // JSON fallback path (used only when protobuf decode fails): thresholds live
    // on lines in the typed model, so carry any legacy rule-level thresholds on a
    // single cell-less line. The protobuf path is the lossless one.
    const parallelThresholds = parseSequentialThresholds(
      getArray(ruleData, 'parallelThresholds') ?? getArray(ruleData, 'parallel_thresholds') ?? []
    );
    transactionRules.push({
      key: getString(ruleData, 'key') ?? getString(ruleData, 'rulesKey') ?? '',
      columns: [],
      lines: parallelThresholds.length > 0 ? [{ cells: [], parallelThresholds, priority: 0 }] : [],
    });
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
      lines: [],
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

  // Strict decoding, but this function's contract is to report a parse failure as "no
  // signatures" rather than throwing. That stays fail-closed: an empty list satisfies no
  // threshold, so ambiguous input can never contribute a signature towards one.
  let decodedBytes: Buffer;
  try {
    decodedBytes = strictBase64Decode(base64Data);
  } catch {
    return [];
  }

  // Try protobuf first (primary format, matches Java/Go SDKs)
  try {
    // Loaded lazily on purpose: protobuf decoding is the primary path but must not be a
    // hard static dependency, so the JSON fallback still works without the generated module.
    // eslint-disable-next-line @typescript-eslint/no-var-requires
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

/**
 * Maps a SuperAdmin public key DTO to its domain model.
 *
 * @param dto - The public key DTO
 * @returns The SuperAdminPublicKey, or undefined when the DTO is absent
 */
export function superAdminPublicKeyFromDto(
  dto: GetPublicKeysReplyPublicKey | null | undefined
): SuperAdminPublicKey | undefined {
  if (!dto) {
    return undefined;
  }
  return {
    userId: safeString(dto.userID) ?? "",
    publicKey: safeString(dto.publicKey) ?? "",
  };
}

/**
 * Maps an array of SuperAdmin public key DTOs to domain models.
 *
 * @param dtos - The array of public key DTOs
 * @returns Array of SuperAdminPublicKey models (undefined entries filtered out)
 */
export function superAdminPublicKeysFromDto(
  dtos: GetPublicKeysReplyPublicKey[] | null | undefined
): SuperAdminPublicKey[] {
  return safeMap(dtos, superAdminPublicKeyFromDto);
}
