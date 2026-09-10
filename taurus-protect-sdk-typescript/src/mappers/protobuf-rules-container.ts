/**
 * Protobuf → model decoding for RulesContainer (mirrors the Go SDK's
 * rules_container.go decode path).
 *
 * The decode is lossless: every governance field is carried onto the model,
 * transaction-rule cells are decoded into the typed {@link RuleCell} union, and
 * per-node protobuf unknown fields are captured (ts-proto `_unknownFields`,
 * enabled via `--ts_proto_opt=unknownFields=true`) so a container from a newer
 * schema re-encodes without dropping data. Encoding lives in
 * `protobuf-rules-container-encode.ts`.
 */

import type {
  DecodedRulesContainer,
  RuleUser,
  RuleGroup,
  TransactionRules,
  TransactionRuleDetails,
  RuleColumn,
  RuleLine,
  AddressWhitelistingRules,
  AddressWhitelistingLine,
  ContractAddressWhitelistingRules,
  GroupThreshold,
  SequentialThresholds,
  RuleSource,
  UnknownFields,
} from '../models/governance-rules';
import { RuleSourceType } from '../models/governance-rules';
import { ruleCellFromBytes } from './rule-cell-codec';
import { ruleSourceToBytes } from './protobuf-rules-container-encode';
import {
  RulesContainer as ProtobufRulesContainer,
  Role,
  Blockchain,
  RuleSource as ProtobufRuleSource,
  RuleSource_RuleSourceType,
  RuleSourceInternalWallet as PbSourceInternalWallet,
  RuleSourceInternalAddress as PbSourceInternalAddress,
  RuleSourceExchange as PbSourceExchange,
  RuleSourceExternalAddress as PbSourceExternalAddress,
  RulesContainer_ColumnType,
  RulesContainer_TransactionRules_TransactionRuleDetails_RuleDomain as PbRuleDomain,
  RulesContainer_TransactionRules_TransactionRuleDetails_RuleSubDomain as PbRuleSubDomain,
  type User as ProtobufUser,
  type Group as ProtobufGroup,
  type RulesContainer_AddressWhitelistingRules as PbAddressWhitelistingRules,
  type RulesContainer_AddressWhitelistingRules_Line as PbAddressWhitelistingLine,
  type RulesContainer_ContractAddressWhitelistingRules as PbContractAddressWhitelistingRules,
  type RulesContainer_TransactionRules as PbTransactionRules,
  type RulesContainer_Column as PbColumn,
  type RulesContainer_Line as PbLine,
  type RulesContainer_TransactionRules_TransactionRuleDetails as PbDetails,
  type SequentialThresholds as PbSequentialThresholds,
  type GroupThreshold as PbGroupThreshold,
} from '../internal/proto/request_reply';

/** ts-proto stores unknown fields on `_unknownFields`; return it if non-empty. */
function uf(m: { _unknownFields?: UnknownFields }): UnknownFields | undefined {
  const u = m._unknownFields;
  return u && Object.keys(u).length > 0 ? u : undefined;
}

function mapOf(m: { [key: string]: Uint8Array } | undefined): { [key: string]: Uint8Array } | undefined {
  return m && Object.keys(m).length > 0 ? m : undefined;
}

/** Number → enum-value name, or the decimal string for values unknown to this SDK (passthrough). */
export function pbEnumName(enumObj: Record<string | number, string | number>, value: number): string {
  const name = enumObj[value];
  return typeof name === 'string' && name !== 'UNRECOGNIZED' ? name : String(value);
}

/**
 * Attempts to decode protobuf bytes to a DecodedRulesContainer.
 * Returns undefined if the bytes are not a protobuf RulesContainer.
 */
export function tryDecodeProtobufRulesContainer(bytes: Uint8Array): DecodedRulesContainer | undefined {
  try {
    return rulesContainerFromProtobuf(ProtobufRulesContainer.decode(bytes));
  } catch {
    return undefined;
  }
}

function rulesContainerFromProtobuf(pb: ReturnType<typeof ProtobufRulesContainer.decode>): DecodedRulesContainer {
  const users: RuleUser[] = pb.users.map((u: ProtobufUser) => ({
    id: u.id || undefined,
    name: undefined, // protobuf User has no name field
    publicKeyPem: u.publicKey || undefined,
    roles: u.roles.map((r: Role) => pbEnumName(Role, r)),
    properties: mapOf(u.properties),
    unknownFields: uf(u),
  }));

  const groups: RuleGroup[] = pb.groups.map((g: ProtobufGroup) => ({
    id: g.id || undefined,
    name: undefined,
    userIds: [...g.userIds],
    properties: mapOf(g.properties),
    unknownFields: uf(g),
  }));

  const transactionRules: TransactionRules[] = (pb.transactionRules || []).map(transactionRulesFromProtobuf);

  const addressWhitelistingRules: AddressWhitelistingRules[] = pb.addressWhitelistingRules.map(
    (r: PbAddressWhitelistingRules) => ({
      currency: r.currency || undefined,
      network: r.network || undefined,
      parallelThresholds: mapSequentialThresholds(r.parallelThresholds || []),
      lines: (r.lines || []).map(addressWhitelistingLineFromProtobuf),
      properties: mapOf(r.properties),
      unknownFields: uf(r),
    })
  );

  const contractAddressWhitelistingRules: ContractAddressWhitelistingRules[] =
    pb.contractAddressWhitelistingRules.map((r: PbContractAddressWhitelistingRules) => ({
      blockchain:
        r.blockchain !== undefined && r.blockchain !== Blockchain.None
          ? pbEnumName(Blockchain, r.blockchain)
          : undefined,
      network: r.network || undefined,
      parallelThresholds: mapSequentialThresholds(r.parallelThresholds || []),
      properties: mapOf(r.properties),
      unknownFields: uf(r),
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
    timestamp: typeof pb.timestamp === 'bigint' ? Number(pb.timestamp) : pb.timestamp || 0,
    hsmSlotId: pb.hsmSlotId || 0,
    minimumCommitmentSignatures: pb.minimumCommitmentSignatures || 0,
    engineIdentities: pb.engineIdentities ? [...pb.engineIdentities] : [],
    properties: mapOf(pb.properties),
    unknownFields: uf(pb),
  };
}

function transactionRulesFromProtobuf(r: PbTransactionRules): TransactionRules {
  const columns: RuleColumn[] = (r.columns || []).map((c: PbColumn) => ({
    type: pbEnumName(RulesContainer_ColumnType, c.type),
    name: c.name,
    metadataKey: c.metadataKey,
    unknownFields: uf(c),
  }));

  const lines: RuleLine[] = (r.lines || []).map((l: PbLine) => ({
    cells: (l.cells || []).map((cellBytes, i) =>
      ruleCellFromBytes(i < columns.length ? columns[i].type : '', cellBytes)
    ),
    parallelThresholds: mapSequentialThresholds(l.parallelThresholds || []),
    priority: l.priority || 0,
    properties: mapOf(l.properties),
    unknownFields: uf(l),
  }));

  let details: TransactionRuleDetails | undefined;
  if (r.details) {
    const d: PbDetails = r.details;
    details = {
      domain: pbEnumName(PbRuleDomain, d.domain),
      subDomain: pbEnumName(PbRuleSubDomain, d.subDomain),
      blockchain: d.blockchain,
      network: d.network,
      evmCallContract: d.evmCallContract
        ? { contractType: d.evmCallContract.contractType, methodSignature: d.evmCallContract.methodSignature, unknownFields: uf(d.evmCallContract) }
        : undefined,
      xtzCallContract: d.xtzCallContract
        ? { contractType: d.xtzCallContract.contractType, methodSignature: d.xtzCallContract.methodSignature, unknownFields: uf(d.xtzCallContract) }
        : undefined,
      cashSettlement: d.cashSettlement
        ? { provider: d.cashSettlement.provider, requestType: d.cashSettlement.requestType, unknownFields: uf(d.cashSettlement) }
        : undefined,
      cosmosDetails: d.cosmosDetails
        ? { methodSignatures: [...d.cosmosDetails.methodSignatures], unknownFields: uf(d.cosmosDetails) }
        : undefined,
      unknownFields: uf(d),
    };
  }

  return { key: r.key, columns, lines, details, unknownFields: uf(r) };
}

function addressWhitelistingLineFromProtobuf(line: PbAddressWhitelistingLine): AddressWhitelistingLine {
  return {
    cells: (line.cells || []).map(ruleSourceFromBytes),
    parallelThresholds: mapSequentialThresholds(line.parallelThresholds || []),
    properties: mapOf(line.properties),
    unknownFields: uf(line),
  };
}

/**
 * Decodes a whitelisting RuleSource cell, keeping anything this SDK cannot
 * reproduce byte-for-byte as a verbatim raw passthrough.
 *
 * Mirrors the cell codec's lossless guard. An unknown-field check alone is too weak:
 * a non-canonical encoding carries no unknown field yet still re-encodes differently.
 * An explicitly-present zero-length payload (`08 01 12 00`) decodes to the same typed
 * value as an absent one (`08 01`), so re-emitting the typed form would drop two
 * bytes from a container the SuperAdmins signed.
 */
export function ruleSourceFromBytes(data: Uint8Array): RuleSource {
  const typed = ruleSourceFromBytesTyped(data);
  if (typed === undefined) {
    return { type: RuleSourceType.Unknown, raw: data };
  }
  let reencoded: Uint8Array;
  try {
    reencoded = ruleSourceToBytes(typed);
  } catch {
    return { type: RuleSourceType.Unknown, raw: data };
  }
  if (!Buffer.from(reencoded).equals(Buffer.from(data))) {
    return { type: RuleSourceType.Unknown, raw: data };
  }
  return typed;
}

/** Decodes into the typed form, or undefined when this SDK cannot type it. */
function ruleSourceFromBytesTyped(data: Uint8Array): RuleSource | undefined {
  let pb: ProtobufRuleSource;
  try {
    pb = ProtobufRuleSource.decode(data);
  } catch {
    return undefined;
  }
  // A schema-newer field on the source itself cannot be represented by the typed model,
  // so keep the bytes verbatim rather than dropping it on the next encode.
  if (uf(pb)) return undefined;

  const src: { -readonly [K in keyof RuleSource]?: RuleSource[K] } = { type: pb.type as number as RuleSourceType };
  try {
    switch (pb.type) {
      case RuleSource_RuleSourceType.RuleSourceAny:
      case RuleSource_RuleSourceType.RuleSourceAnyExchange:
        // Payload-less arms: a present payload is data this SDK would discard.
        if (pb.payload.length > 0) return undefined;
        break;
      case RuleSource_RuleSourceType.RuleSourceInternalWallet:
        if (pb.payload.length > 0) {
          const i = PbSourceInternalWallet.decode(pb.payload);
          if (uf(i)) return undefined;
          src.internalWallet = { path: i.path };
        }
        break;
      case RuleSource_RuleSourceType.RuleSourceInternalAddress:
        if (pb.payload.length > 0) {
          const i = PbSourceInternalAddress.decode(pb.payload);
          if (uf(i)) return undefined;
          src.internalAddress = { address: i.address, path: i.path };
        }
        break;
      case RuleSource_RuleSourceType.RuleSourceExchange:
        if (pb.payload.length > 0) {
          const i = PbSourceExchange.decode(pb.payload);
          if (uf(i)) return undefined;
          src.exchange = { label: i.label };
        }
        break;
      case RuleSource_RuleSourceType.RuleSourceExternalAddress:
        if (pb.payload.length > 0) {
          const i = PbSourceExternalAddress.decode(pb.payload);
          if (uf(i)) return undefined;
          src.externalAddress = { address: i.address, memo: i.memo };
        }
        break;
      default:
        return undefined; // source type newer than this SDK
    }
  } catch {
    return undefined;
  }
  return src as RuleSource;
}

function mapSequentialThresholds(seq: PbSequentialThresholds[]): SequentialThresholds[] {
  return seq.map((s) => ({
    thresholds: (s.thresholds || []).map(mapGroupThreshold),
    unknownFields: uf(s),
  }));
}

function mapGroupThreshold(t: PbGroupThreshold): GroupThreshold {
  return {
    groupId: t.groupId || undefined,
    minimumSignatures: t.minimumSignatures || 0,
    threshold: 0, // no separate proto field; retained for backward compatibility
    unknownFields: uf(t),
  };
}
