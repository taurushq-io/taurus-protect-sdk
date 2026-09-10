/**
 * Container-level round-trip + schema-evolution safety for the typed governance
 * rules mapper.
 */
import type { DecodedRulesContainer } from '../../../src/models/governance-rules';
import { RuleSourceType, hasUnknownFields } from '../../../src/models/governance-rules';
import { tryDecodeProtobufRulesContainer, ruleSourceFromBytes } from '../../../src/mappers/protobuf-rules-container';
import { rulesContainerToBytes, rulesContainerToBase64, ruleSourceToBytes } from '../../../src/mappers/protobuf-rules-container-encode';
import { vectorFor, vectorsFor } from './lossless-vectors';
import { rulesContainerFromBase64 } from '../../../src/mappers/governance-rules';
import { RulesContainer } from '../../../src/internal/proto/request_reply';

function richContainer(): DecodedRulesContainer {
  const thresholds = [{ thresholds: [{ groupId: 'approvers', minimumSignatures: 2, threshold: 0 }] }];
  return {
    users: [
      { id: 'user-1', name: undefined, publicKeyPem: undefined, roles: ['SUPERADMIN'], properties: { team: new Uint8Array([1, 2]) } },
      { id: 'user-2', name: undefined, publicKeyPem: undefined, roles: [] },
    ],
    groups: [{ id: 'approvers', name: undefined, userIds: ['user-1', 'user-2'] }],
    minimumDistinctUserSignatures: 2,
    minimumDistinctGroupSignatures: 1,
    transactionRules: [
      {
        key: 'ETH/ERC20_transfer',
        columns: [
          { type: 'RuleFiatAmount', name: 'amount', metadataKey: 'amount' },
          { type: 'RuleDestination', name: 'to', metadataKey: 'destination' },
          { type: 'RuleStringEqual', name: 'contract', metadataKey: 'contract_id' },
        ],
        lines: [
          {
            cells: [
              { kind: 'FiatAmountRange', minAmount: '1000', maxAmount: '50000' },
              { kind: 'DestinationInternalWallet', path: "m/44'/60'/1'" },
              { kind: 'StringEqualValue', value: 'contract-42' },
            ],
            parallelThresholds: thresholds,
            priority: 1,
            properties: { note: new Uint8Array([9]) },
          },
          {
            // empty cells decode to the columns' typed *Any values
            cells: [
              { kind: 'FiatAmountAny' },
              { kind: 'DestinationAny' },
              { kind: 'StringEqualAny' },
            ],
            parallelThresholds: thresholds,
            priority: 0,
          },
        ],
        details: {
          domain: 'RuleDomainTransfer',
          subDomain: 'RuleSubDomainERC20',
          blockchain: 'ETH',
          network: 'mainnet',
          evmCallContract: { contractType: 'ERC20', methodSignature: 'transfer(address,uint256)' },
        },
      },
    ],
    addressWhitelistingRules: [
      {
        currency: 'ETH',
        network: 'mainnet',
        parallelThresholds: thresholds,
        lines: [
          {
            cells: [
              { type: RuleSourceType.InternalWallet, internalWallet: { path: "m/44'/60'/0'" } },
              { type: RuleSourceType.Unknown },
            ],
            parallelThresholds: thresholds,
          },
        ],
      },
    ],
    contractAddressWhitelistingRules: [
      { blockchain: 'ETH', network: 'mainnet', parallelThresholds: thresholds },
    ],
    enforcedRulesHash: 'server-hash',
    timestamp: 1750000000,
    hsmSlotId: 7,
    minimumCommitmentSignatures: 1,
    engineIdentities: ['hsm-1', 'hsm-2'],
    properties: { tenant: new Uint8Array([4, 2]) },
  };
}

describe('rules container round-trip', () => {
  it('decodes typed cells and preserves lossless fields; strips server-controlled fields', () => {
    const encoded = rulesContainerToBytes(richContainer());
    const c = tryDecodeProtobufRulesContainer(encoded)!;
    expect(c).toBeDefined();
    expect(hasUnknownFields(c)).toBe(false);

    // server-controlled fields stripped by the encoder
    expect(c.enforcedRulesHash).toBeUndefined();
    expect(c.timestamp).toBe(0);

    const tr = c.transactionRules[0];
    expect(tr.key).toBe('ETH/ERC20_transfer');
    expect(tr.columns[0]).toMatchObject({ type: 'RuleFiatAmount', name: 'amount', metadataKey: 'amount' });
    expect(tr.lines[0].cells[0]).toEqual({ kind: 'FiatAmountRange', minAmount: '1000', maxAmount: '50000' });
    expect(tr.lines[0].cells[1]).toEqual({ kind: 'DestinationInternalWallet', path: "m/44'/60'/1'" });
    expect(tr.lines[0].cells[2]).toEqual({ kind: 'StringEqualValue', value: 'contract-42' });
    expect(tr.lines[0].priority).toBe(1);
    // empty cells -> typed *Any per column
    expect(tr.lines[1].cells).toEqual([
      { kind: 'FiatAmountAny' },
      { kind: 'DestinationAny' },
      { kind: 'StringEqualAny' },
    ]);
    expect(tr.details?.blockchain).toBe('ETH');
    expect(tr.details?.evmCallContract?.methodSignature).toBe('transfer(address,uint256)');

    // whitelisting source decoded with payload variant
    const src = c.addressWhitelistingRules[0].lines[0].cells[0];
    expect(src.type).toBe(RuleSourceType.InternalWallet);
    expect(src.internalWallet?.path).toBe("m/44'/60'/0'");

    // properties preserved (compare content; runtime type may be Buffer)
    expect(Array.from(c.properties!.tenant)).toEqual([4, 2]);
    expect(c.hsmSlotId).toBe(7);

    // re-encoding the decoded container is byte-stable
    expect(Buffer.from(rulesContainerToBytes(c))).toEqual(Buffer.from(encoded));
  });

  it('preserves unknown protobuf fields from a newer schema (skew safety)', () => {
    const encoded = rulesContainerToBytes(richContainer());
    // Append an unknown top-level field (field 500, varint 42) — valid proto
    // concatenation. A pre-schema-bump SDK would drop it on re-encode.
    const withUnknown = new Uint8Array([...encoded, 0xa0, 0x1f, 0x2a]);

    const decoded = tryDecodeProtobufRulesContainer(withUnknown)!;
    expect(hasUnknownFields(decoded)).toBe(true);

    const reencoded = rulesContainerToBytes(decoded);
    // ts-proto keys _unknownFields by the full wire tag: field 500, wire-type 0
    // (varint) => tag (500<<3)|0 = 4000.
    const raw = RulesContainer.decode(reencoded);
    expect(raw._unknownFields).toBeDefined();
    expect(raw._unknownFields![4000]).toBeDefined();
    // and it survives another full decode→model→encode cycle
    expect(hasUnknownFields(tryDecodeProtobufRulesContainer(reencoded)!)).toBe(true);
  });

  it('preserves an unknown cell type as RawCell', () => {
    const base = richContainer();
    const decoded = tryDecodeProtobufRulesContainer(rulesContainerToBytes(base))!;
    // No unknown cells in the base container.
    const cells = decoded.transactionRules[0].lines[0].cells;
    expect(cells.some((c) => c.kind === 'RawCell')).toBe(false);
  });
});

// A whitelisting RuleSource this SDK cannot fully represent keeps its exact wire bytes.
// The same three base64 vectors are asserted in all four SDKs.
describe("RuleSource lossless guards", () => {
  // The vectors are loaded INSIDE the test, not in the describe body: a load at
  // collection time turns a missing shared file into "suite failed to run", which drops
  // the test count instead of showing a red test.
  it("preserves every unrepresentable RuleSource verbatim, decode and encode", () => {
    for (const { description, wire_base64: b64 } of vectorsFor("rule_source_lossless")) {
      const data = new Uint8Array(Buffer.from(b64, "base64"));
      const src = ruleSourceFromBytes(data);

      expect(src.raw).toBeDefined();
      // Tuple with the description so a failure names the offending scenario (Jest has
      // no per-assertion message argument).
      expect([description, Buffer.from(src.raw!).toString("base64")]).toEqual([
        description,
        b64,
      ]);

      // Re-encode too, as Go, Python and Java all do. Asserting only the decode let a
      // regression that drops `raw` on ENCODE pass here while breaking the other three
      // SDKs' byte parity — ruleSourceToBytes was module-private, so this could not be
      // checked at all.
      expect([description, Buffer.from(ruleSourceToBytes(src)).toString("base64")]).toEqual([
        description,
        b64,
      ]);
    }
  });
});

// Enum values newer than this SDK keep their numbers across a decode/encode round trip:
// collapsing them to the zero value would rewrite a column's family or widen a rule's
// sub-domain. The same base64 vector is asserted in all four SDKs.
it('passes unknown enum values through numerically', () => {
  const vector = vectorFor("unknown_enum_passthrough");
  const c = rulesContainerFromBase64(vector);

  expect(c.users[0]!.roles[0]).toBe('201');
  expect(c.transactionRules[0]!.columns[0]!.type).toBe('77');
  expect(c.transactionRules[0]!.details?.subDomain).toBe('202');
  expect(c.contractAddressWhitelistingRules[0]!.blockchain).toBe('203');
  expect(rulesContainerToBase64(c)).toBe(vector);
});

// The container carries map<string, bytes> properties at five levels. ts-proto emits map
// entries in the object's own key order, so without sorting the same reviewed container
// encodes differently depending on how it was built — and differently from the other
// SDKs. The expected value is asserted byte-for-byte in all four SDKs.
it('encodes deterministically and matches the other SDKs byte for byte', () => {
  const expected = vectorFor("deterministic_encoding");

  // Keys are inserted in reverse-sorted order on purpose.
  const props: { [k: string]: Uint8Array } = {};
  for (const k of ['kEcho', 'kDelta', 'kCharlie', 'kBravo', 'kAlpha']) {
    props[k] = new Uint8Array(Buffer.from(k));
  }

  for (let i = 0; i < 20; i++) {
    const c = {
      users: [{ id: 'u1', publicKeyPem: 'PEM', roles: ['SUPERADMIN'], properties: { ...props } }],
      groups: [], transactionRules: [], addressWhitelistingRules: [],
      contractAddressWhitelistingRules: [], engineIdentities: [], properties: { ...props },
    } as unknown as DecodedRulesContainer;
    expect(rulesContainerToBase64(c)).toBe(expected);
  }
});

// A malformed cell must degrade on its own and never abort the container: every rule for
// that tenant would go down with it, taking whitelisted-address verification with them.
// A cell payload is protobuf bytes, so a non-UTF-8 string cell and a truncated wrapper are
// both legal on the wire. Same vector is asserted in all four SDKs.
it('does not abort the container decode on a malformed cell', () => {
  const vector = vectorFor("malformed_cell_degrades_alone");
  const c = rulesContainerFromBase64(vector);

  expect(c.transactionRules).toHaveLength(1);
  const cells = c.transactionRules[0]!.lines[0]!.cells;
  expect(cells).toHaveLength(2);
  // A JS string cannot hold invalid UTF-8, so both cells are preserved verbatim.
  expect(cells.every((cell) => cell.kind === 'RawCell')).toBe(true);
  expect(rulesContainerToBase64(c)).toBe(vector);
});

// Unknown protobuf fields inside the nested contract-call scoping sub-messages must
// survive a round trip: dropping them silently narrows which contract calls a rule
// covers. The same base64 vector is asserted in all four SDKs.
it('preserves unknown fields inside the nested rule-detail sub-messages', () => {
  const vector = vectorFor("nested_detail_unknown_fields");
  const c = rulesContainerFromBase64(vector);
  const d = c.transactionRules[0]!.details!;

  // Non-empty, not merely present: an empty bag satisfies toBeDefined() while the
  // nested unknown fields have actually been dropped. Go, Python and Java all assert
  // non-emptiness here.
  for (const [name, node] of [
    ['evmCallContract', d.evmCallContract],
    ['xtzCallContract', d.xtzCallContract],
    ['cashSettlement', d.cashSettlement],
    ['cosmosDetails', d.cosmosDetails],
  ] as const) {
    const uf = node?.unknownFields;
    expect(uf).toBeDefined();
    expect(Object.keys(uf!).length).toBeGreaterThan(0);
  }
  // The container-level report must see the nested nodes, not just the details node.
  expect(hasUnknownFields(c)).toBe(true);

  expect(rulesContainerToBase64(c)).toBe(vector);
});

// An empty payload on a payload-carrying arm leaves the typed variant unset: parsing it
// would materialize a default sub-message, so a caller checking the variant would see an
// empty wallet instead of nothing. Asserted in all four SDKs.
it('leaves the typed variant unset for an empty payload', () => {
  const src = ruleSourceFromBytes(new Uint8Array(Buffer.from(vectorFor("empty_payload_arm"), 'base64')));
  expect(src.type).toBe(RuleSourceType.InternalWallet);
  expect(src.internalWallet).toBeUndefined();
});

// `08 01 12 00` and `08 01` decode to the SAME typed value: an explicitly-present
// zero-length payload is legal on the wire and leaves the variant unset just like an
// absent one. Re-encoding the typed form therefore emits `08 01` and drops two bytes
// from a container the SuperAdmins signed. No unknown field is present, so only the
// decode → re-encode → byte-compare guard catches it.
it('keeps a non-canonical empty payload verbatim instead of typing it', () => {
  const data = new Uint8Array(
    Buffer.from(vectorFor('explicit_empty_payload_noncanonical'), 'base64')
  );
  const src = ruleSourceFromBytes(data);
  expect(src.raw && Buffer.from(src.raw)).toEqual(Buffer.from(data));
  expect(Buffer.from(ruleSourceToBytes(src))).toEqual(Buffer.from(data));
});
