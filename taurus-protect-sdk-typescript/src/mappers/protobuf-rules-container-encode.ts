/**
 * Model → protobuf encoding for RulesContainer (mirrors the Go SDK's
 * rules_container_encode.go).
 *
 * Inverse of the decode in `protobuf-rules-container.ts`. Server-controlled
 * `enforcedRulesHash` and `timestamp` are stripped (the server recomputes them
 * and rejects submissions asserting stale values). Per-node unknown protobuf
 * fields captured at decode are re-attached (ts-proto `_unknownFields`), and
 * unknown enum values pass through numerically, so a container from a newer
 * schema re-encodes without dropping data. Enum names the caller authored that
 * this SDK does not recognize are a hard error.
 */

import type {
  DecodedRulesContainer,
  RuleUser,
  RuleGroup,
  TransactionRules,
  RuleLine,
  AddressWhitelistingRules,
  AddressWhitelistingLine,
  ContractAddressWhitelistingRules,
  SequentialThresholds,
  RuleSource,
  UnknownFields,
} from '../models/governance-rules';
import { RuleSourceType } from '../models/governance-rules';
import { ruleCellToBytes } from './rule-cell-codec';
import {
  RulesContainer,
  User,
  Group,
  Role,
  Blockchain,
  SequentialThresholds as PbSequentialThresholds,
  GroupThreshold as PbGroupThreshold,
  RuleSource as PbRuleSource,
  RuleSource_RuleSourceType,
  RuleSourceInternalWallet as PbSourceInternalWallet,
  RuleSourceInternalAddress as PbSourceInternalAddress,
  RuleSourceExchange as PbSourceExchange,
  RuleSourceExternalAddress as PbSourceExternalAddress,
  RulesContainer_Column,
  RulesContainer_ColumnType,
  RulesContainer_Line,
  RulesContainer_TransactionRules,
  RulesContainer_TransactionRules_TransactionRuleDetails as PbDetails,
  RulesContainer_TransactionRules_TransactionRuleDetails_RuleDomain as PbRuleDomain,
  RulesContainer_TransactionRules_TransactionRuleDetails_RuleSubDomain as PbRuleSubDomain,
  RulesContainer_TransactionRules_TransactionRuleDetails_EvmCallContract as PbEvm,
  RulesContainer_TransactionRules_TransactionRuleDetails_XtzCallContract as PbXtz,
  RulesContainer_TransactionRules_TransactionRuleDetails_CashSettlement as PbCash,
  RulesContainer_TransactionRules_TransactionRuleDetails_CosmosDetails as PbCosmos,
  RulesContainer_AddressWhitelistingRules,
  RulesContainer_AddressWhitelistingRules_Line,
  RulesContainer_ContractAddressWhitelistingRules,
} from '../internal/proto/request_reply';

/** Encodes a DecodedRulesContainer to base64 protobuf wire format. */
export function rulesContainerToBase64(container: DecodedRulesContainer): string {
  return Buffer.from(rulesContainerToBytes(container)).toString('base64');
}

/** Encodes a DecodedRulesContainer to raw protobuf bytes. */
export function rulesContainerToBytes(container: DecodedRulesContainer): Uint8Array {
  return RulesContainer.encode(rulesContainerToProto(container)).finish();
}


/** Name → enum number: decimal strings pass through; unknown names throw. */
function enumNumber(enumObj: Record<string, string | number>, name: string, what: string): number {
  if (name === '') return 0;
  if (/^-?\d+$/.test(name)) return parseInt(name, 10);
  const v = enumObj[name];
  if (typeof v === 'number') return v;
  throw new Error(`unknown ${what} "${name}"`);
}

function attach<T extends { _unknownFields?: UnknownFields }>(m: T, unknown: UnknownFields | undefined): T {
  if (unknown && Object.keys(unknown).length > 0) m._unknownFields = unknown;
  return m;
}

// Nested protobuf messages must be assigned AFTER the parent's fromPartial:
// ts-proto's fromPartial re-runs the child's fromPartial on each element, and the
// generated child fromPartial does not copy `_unknownFields` — so building a
// parent via `fromPartial({ children: [...built] })` would silently strip every
// child's preserved unknown fields. Assigning post-build keeps them intact.
function thresholdsToProto(list: SequentialThresholds[]): PbSequentialThresholds[] {
  return list.map((st) => {
    const pb = PbSequentialThresholds.fromPartial({});
    pb.thresholds = st.thresholds.map((t) =>
      attach(PbGroupThreshold.fromPartial({ groupId: t.groupId ?? '', minimumSignatures: t.minimumSignatures }), t.unknownFields)
    );
    return attach(pb, st.unknownFields);
  });
}

function userToProto(u: RuleUser): User {
  return attach(
    User.fromPartial({
      id: u.id ?? '',
      publicKey: u.publicKeyPem ?? '',
      roles: u.roles.map((r) => enumNumber(Role as never, r, 'user role')),
      properties: u.properties,
    }),
    u.unknownFields
  );
}

function groupToProto(g: RuleGroup): Group {
  return attach(
    Group.fromPartial({ id: g.id ?? '', userIds: g.userIds, properties: g.properties }),
    g.unknownFields
  );
}

function transactionRulesToProto(r: TransactionRules): RulesContainer_TransactionRules {
  const columns = r.columns.map((c) =>
    attach(
      RulesContainer_Column.fromPartial({
        type: enumNumber(RulesContainer_ColumnType as never, c.type, 'column type'),
        name: c.name,
        metadataKey: c.metadataKey,
      }),
      c.unknownFields
    )
  );

  const lines = r.lines.map((l: RuleLine) => {
    const pbLine = RulesContainer_Line.fromPartial({
      cells: l.cells.map((cell, i) => ruleCellToBytes(i < r.columns.length ? r.columns[i].type : '', cell)),
      priority: l.priority,
      properties: l.properties,
    });
    pbLine.parallelThresholds = thresholdsToProto(l.parallelThresholds);
    return attach(pbLine, l.unknownFields);
  });

  let details: PbDetails | undefined;
  if (r.details) {
    const d = r.details;
    const pbDetails = PbDetails.fromPartial({
      domain: enumNumber(PbRuleDomain as never, d.domain, 'rule domain'),
      subDomain: enumNumber(PbRuleSubDomain as never, d.subDomain, 'rule sub-domain'),
      blockchain: d.blockchain,
      network: d.network,
    });
    // Assign the nested scoping messages after fromPartial (see thresholdsToProto):
    // passing them in would re-run each child's fromPartial, which drops the
    // `_unknownFields` bag and silently narrows a rule's contract-call scoping.
    if (d.evmCallContract) {
      pbDetails.evmCallContract = attach(
        PbEvm.fromPartial({ contractType: d.evmCallContract.contractType, methodSignature: d.evmCallContract.methodSignature }),
        d.evmCallContract.unknownFields
      );
    }
    if (d.xtzCallContract) {
      pbDetails.xtzCallContract = attach(
        PbXtz.fromPartial({ contractType: d.xtzCallContract.contractType, methodSignature: d.xtzCallContract.methodSignature }),
        d.xtzCallContract.unknownFields
      );
    }
    if (d.cashSettlement) {
      pbDetails.cashSettlement = attach(
        PbCash.fromPartial({ provider: d.cashSettlement.provider, requestType: d.cashSettlement.requestType }),
        d.cashSettlement.unknownFields
      );
    }
    if (d.cosmosDetails) {
      pbDetails.cosmosDetails = attach(
        PbCosmos.fromPartial({ methodSignatures: d.cosmosDetails.methodSignatures }),
        d.cosmosDetails.unknownFields
      );
    }
    details = attach(pbDetails, d.unknownFields);
  }

  // Assign nested messages after fromPartial (see thresholdsToProto) so their
  // preserved _unknownFields are not stripped.
  const pbRule = RulesContainer_TransactionRules.fromPartial({ key: r.key });
  pbRule.columns = columns;
  pbRule.lines = lines;
  if (details) pbRule.details = details;
  return attach(pbRule, r.unknownFields);
}

/**
 * Serializes a whitelisting source cell. A preserved `raw` cell is re-emitted byte for
 * byte: decoding it first would re-run the parse whose failure caused it to be preserved,
 * and even when it parses, a decode/encode round trip normalizes field order — which
 * would change the container the SuperAdmin signatures were computed over.
 */
export function ruleSourceToBytes(s: RuleSource): Uint8Array {
  if (s.raw !== undefined && s.raw.length > 0) {
    return s.raw;
  }
  return PbRuleSource.encode(ruleSourceToProto(s)).finish();
}

function ruleSourceToProto(s: RuleSource): PbRuleSource {
  let payload: Uint8Array | undefined;
  switch (s.type) {
    case RuleSourceType.Unknown:
    case RuleSourceType.AnyExchange:
      break;
    case RuleSourceType.InternalWallet:
      if (s.internalWallet) payload = PbSourceInternalWallet.encode(PbSourceInternalWallet.fromPartial(s.internalWallet)).finish();
      break;
    case RuleSourceType.InternalAddress:
      if (s.internalAddress) payload = PbSourceInternalAddress.encode(PbSourceInternalAddress.fromPartial(s.internalAddress)).finish();
      break;
    case RuleSourceType.Exchange:
      if (s.exchange) payload = PbSourceExchange.encode(PbSourceExchange.fromPartial(s.exchange)).finish();
      break;
    case RuleSourceType.ExternalAddress:
      if (s.externalAddress) payload = PbSourceExternalAddress.encode(PbSourceExternalAddress.fromPartial(s.externalAddress)).finish();
      break;
    default:
      throw new Error(`unknown rule source type ${s.type}`);
  }
  return PbRuleSource.fromPartial({ type: s.type as number as RuleSource_RuleSourceType, payload });
}

function addressWhitelistingRulesToProto(r: AddressWhitelistingRules): RulesContainer_AddressWhitelistingRules {
  const lines = r.lines.map((l: AddressWhitelistingLine) => {
    const pbLine = RulesContainer_AddressWhitelistingRules_Line.fromPartial({
      cells: l.cells.map(ruleSourceToBytes),
      properties: l.properties,
    });
    pbLine.parallelThresholds = thresholdsToProto(l.parallelThresholds);
    return attach(pbLine, l.unknownFields);
  });
  // Assign nested messages after fromPartial (see thresholdsToProto).
  const pb = RulesContainer_AddressWhitelistingRules.fromPartial({
    currency: r.currency ?? '',
    network: r.network ?? '',
    properties: r.properties,
  });
  pb.parallelThresholds = thresholdsToProto(r.parallelThresholds);
  pb.lines = lines;
  return attach(pb, r.unknownFields);
}

function contractAddressWhitelistingRulesToProto(
  r: ContractAddressWhitelistingRules
): RulesContainer_ContractAddressWhitelistingRules {
  // Assign nested thresholds after fromPartial (see thresholdsToProto).
  const pb = RulesContainer_ContractAddressWhitelistingRules.fromPartial({
    blockchain: enumNumber(Blockchain as never, r.blockchain ?? '', 'blockchain'),
    network: r.network ?? '',
    properties: r.properties,
  });
  pb.parallelThresholds = thresholdsToProto(r.parallelThresholds);
  return attach(pb, r.unknownFields);
}

function rulesContainerToProto(c: DecodedRulesContainer): RulesContainer {
  const pb = RulesContainer.fromPartial({
    minimumDistinctUserSignatures: c.minimumDistinctUserSignatures,
    minimumDistinctGroupSignatures: c.minimumDistinctGroupSignatures,
    minimumCommitmentSignatures: c.minimumCommitmentSignatures,
    engineIdentities: c.engineIdentities,
    hsmSlotId: c.hsmSlotId,
    properties: c.properties,
    // enforcedRulesHash and timestamp intentionally omitted (server-controlled).
  });
  // Assign nested messages after fromPartial (see thresholdsToProto) so per-node
  // preserved _unknownFields are not stripped by the container's fromPartial.
  pb.users = c.users.map(userToProto);
  pb.groups = c.groups.map(groupToProto);
  pb.transactionRules = c.transactionRules.map(transactionRulesToProto);
  pb.addressWhitelistingRules = c.addressWhitelistingRules.map(addressWhitelistingRulesToProto);
  pb.contractAddressWhitelistingRules = c.contractAddressWhitelistingRules.map(contractAddressWhitelistingRulesToProto);
  return attach(pb, c.unknownFields);
}
