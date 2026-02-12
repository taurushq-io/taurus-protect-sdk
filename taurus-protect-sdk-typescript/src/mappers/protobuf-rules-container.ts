/**
 * Protobuf to model conversion for RulesContainer.
 *
 * This module provides functions to decode protobuf-encoded rules containers
 * to the SDK's domain models. The protobuf format is the primary format used
 * by the Taurus-PROTECT API (aligning with the Java SDK).
 */

import type {
  DecodedRulesContainer,
  RuleUser,
  RuleGroup,
  TransactionRules,
  AddressWhitelistingRules,
  ContractAddressWhitelistingRules,
  GroupThreshold,
  SequentialThresholds,
} from '../models/governance-rules';
import {
  RulesContainer as ProtobufRulesContainer,
  Role,
  Blockchain,
  RuleSource as ProtobufRuleSource,
  RuleSourceInternalWallet as ProtobufRuleSourceInternalWallet,
  RuleSource_RuleSourceType,
  type User as ProtobufUser,
  type Group as ProtobufGroup,
  type RulesContainer_AddressWhitelistingRules as ProtobufAddressWhitelistingRules,
  type RulesContainer_AddressWhitelistingRules_Line as ProtobufAddressWhitelistingLine,
  type RulesContainer_ContractAddressWhitelistingRules as ProtobufContractAddressWhitelistingRules,
  type RulesContainer_TransactionRules as ProtobufTransactionRules,
  type SequentialThresholds as ProtobufSequentialThresholds,
  type GroupThreshold as ProtobufGroupThreshold,
} from '../internal/proto/request_reply';

/**
 * Attempts to decode protobuf bytes to DecodedRulesContainer.
 * Returns undefined if decoding fails (e.g., if data is JSON, not protobuf).
 *
 * @param bytes - The raw bytes to decode
 * @returns Decoded rules container, or undefined if decoding fails
 */
export function tryDecodeProtobufRulesContainer(
  bytes: Uint8Array
): DecodedRulesContainer | undefined {
  try {
    const pb = ProtobufRulesContainer.decode(bytes);
    return rulesContainerFromProtobuf(pb);
  } catch {
    return undefined;
  }
}

/**
 * Converts a protobuf RulesContainer to DecodedRulesContainer domain model.
 */
function rulesContainerFromProtobuf(
  pb: ReturnType<typeof ProtobufRulesContainer.decode>
): DecodedRulesContainer {
  // Map users
  const users: RuleUser[] = pb.users.map((u: ProtobufUser) => ({
    id: u.id || undefined,
    name: undefined, // protobuf User doesn't have name field, only id and publicKey
    publicKeyPem: u.publicKey || undefined,
    roles: u.roles.map((r: Role) => Role[r] ?? String(r)),
  }));

  // Map groups
  const groups: RuleGroup[] = pb.groups.map((g: ProtobufGroup) => ({
    id: g.id || undefined,
    name: undefined, // protobuf Group doesn't have name field
    userIds: [...g.userIds],
  }));

  // Map transaction rules
  const transactionRules: TransactionRules[] =
    (pb.transactionRules || []).map((r: ProtobufTransactionRules) => ({
      parallelThresholds: mapSequentialThresholds(
        (r.lines || []).flatMap((line) => line.parallelThresholds || [])
      ),
    }));

  // Map address whitelisting rules (including lines for per-wallet-path thresholds)
  const addressWhitelistingRules = pb.addressWhitelistingRules.map(
    (r: ProtobufAddressWhitelistingRules) => ({
      currency: r.currency || undefined,
      network: r.network || undefined,
      parallelThresholds: mapSequentialThresholds(r.parallelThresholds || []),
      lines: mapAddressWhitelistingLines(r.lines || []),
    })
  );

  // Map contract address whitelisting rules
  const contractAddressWhitelistingRules: ContractAddressWhitelistingRules[] =
    pb.contractAddressWhitelistingRules.map((r: ProtobufContractAddressWhitelistingRules) => ({
      blockchain: r.blockchain !== undefined && r.blockchain !== Blockchain.None
        ? Blockchain[r.blockchain] ?? String(r.blockchain)
        : undefined,
      network: r.network || undefined,
      parallelThresholds: mapSequentialThresholds(r.parallelThresholds || []),
    }));

  return {
    users,
    groups,
    minimumDistinctUserSignatures: pb.minimumDistinctUserSignatures || 0,
    minimumDistinctGroupSignatures: pb.minimumDistinctGroupSignatures || 0,
    transactionRules,
    addressWhitelistingRules,
    contractAddressWhitelistingRules,
    enforcedRulesHash: pb.enforcedRulesHash || undefined,
    timestamp: pb.timestamp || 0,
    hsmSlotId: pb.hsmSlotId || 0,
    minimumCommitmentSignatures: pb.minimumCommitmentSignatures || 0,
    engineIdentities: pb.engineIdentities ? [...pb.engineIdentities] : [],
  };
}

/**
 * Maps protobuf SequentialThresholds[] to domain SequentialThresholds[].
 */
function mapSequentialThresholds(
  seqThresholds: ProtobufSequentialThresholds[]
): SequentialThresholds[] {
  return seqThresholds.map((seq) => ({
    thresholds: (seq.thresholds || []).map(mapGroupThreshold),
  }));
}

/**
 * Maps a single protobuf GroupThreshold to domain GroupThreshold.
 */
function mapGroupThreshold(t: ProtobufGroupThreshold): GroupThreshold {
  return {
    groupId: t.groupId || undefined,
    minimumSignatures: t.minimumSignatures || 0,
    threshold: 0, // protobuf doesn't have a separate threshold field
  };
}

/**
 * Maps protobuf AddressWhitelistingRules.Line[] to domain model.
 * Each line contains cells (serialized RuleSource protobuf) and parallelThresholds.
 */
function mapAddressWhitelistingLines(
  lines: ProtobufAddressWhitelistingLine[]
): Array<{
  cells: Array<{ type: string; internalWallet?: { path?: string } }>;
  parallelThresholds: SequentialThresholds[];
}> {
  return lines.map((line) => ({
    cells: (line.cells || [])
      .map(ruleSourceFromBytes)
      .filter((s): s is NonNullable<typeof s> => s !== undefined),
    parallelThresholds: mapSequentialThresholds(line.parallelThresholds || []),
  }));
}

/**
 * Maps protobuf RuleSource_RuleSourceType to the type strings expected by the verifier.
 */
const RULE_SOURCE_TYPE_MAP: Record<number, string> = {
  [RuleSource_RuleSourceType.RuleSourceAny]: "ANY",
  [RuleSource_RuleSourceType.RuleSourceInternalWallet]: "INTERNAL_WALLET",
  [RuleSource_RuleSourceType.RuleSourceInternalAddress]: "INTERNAL_ADDRESS",
  [RuleSource_RuleSourceType.RuleSourceAnyExchange]: "ANY_EXCHANGE",
  [RuleSource_RuleSourceType.RuleSourceExchange]: "EXCHANGE",
  [RuleSource_RuleSourceType.RuleSourceExternalAddress]: "EXTERNAL_ADDRESS",
};

/**
 * Decodes a RuleSource from serialized protobuf bytes.
 * Matches Go SDK's ruleSourceFromBytes().
 */
function ruleSourceFromBytes(
  data: Uint8Array
): { type: string; internalWallet?: { path?: string } } | undefined {
  try {
    const pbSource = ProtobufRuleSource.decode(data);

    const result: { type: string; internalWallet?: { path?: string } } = {
      type: RULE_SOURCE_TYPE_MAP[pbSource.type] ?? String(pbSource.type),
    };

    // Decode payload for InternalWallet type
    if (
      pbSource.type === RuleSource_RuleSourceType.RuleSourceInternalWallet &&
      pbSource.payload.length > 0
    ) {
      try {
        const pbWallet = ProtobufRuleSourceInternalWallet.decode(pbSource.payload);
        result.internalWallet = { path: pbWallet.path || undefined };
      } catch {
        // Payload decode failure — skip wallet info
      }
    }

    return result;
  } catch {
    return undefined;
  }
}
