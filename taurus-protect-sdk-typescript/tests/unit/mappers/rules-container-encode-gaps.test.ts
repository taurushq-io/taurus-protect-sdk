/**
 * Per-node unknown-field preservation, enum passthrough / named-error, the
 * RuleSource variant codec, and the details sub-messages — the encode/decode
 * branches the baseline round-trip test leaves at the container top level only.
 */
import type { DecodedRulesContainer, UnknownFields } from '../../../src/models/governance-rules';
import { RuleSourceType, hasUnknownFields, createEmptyRulesContainer } from '../../../src/models/governance-rules';
import { tryDecodeProtobufRulesContainer } from '../../../src/mappers/protobuf-rules-container';
import { rulesContainerToBytes } from '../../../src/mappers/protobuf-rules-container-encode';
import {
  RulesContainer,
  RuleSource as PbRuleSource,
  RuleSource_RuleSourceType,
} from '../../../src/internal/proto/request_reply';

const UF: UnknownFields = { 4000: [new Uint8Array([0x2a])] }; // field 500 (tag 4000), varint 42
const th = () => [{ thresholds: [{ groupId: 'g', minimumSignatures: 1, threshold: 0 }] }];
const roundtrip = (c: DecodedRulesContainer) => tryDecodeProtobufRulesContainer(rulesContainerToBytes(c))!;
const bytesOf = (u: Uint8Array | undefined): number[] => (u ? Array.from(u) : []);

describe('per-node unknown-field preservation', () => {
  it('captures and re-attaches unknown fields at every node', () => {
    const c: DecodedRulesContainer = {
      ...createEmptyRulesContainer(),
      unknownFields: UF,
      users: [{ id: 'u', name: undefined, publicKeyPem: undefined, roles: [], unknownFields: UF }],
      groups: [{ id: 'g', name: undefined, userIds: [], unknownFields: UF }],
      transactionRules: [
        {
          key: 'k',
          columns: [{ type: 'RuleFiatAmount', name: 'a', metadataKey: 'a', unknownFields: UF }],
          lines: [{ cells: [], parallelThresholds: [{ thresholds: [{ groupId: 'g', minimumSignatures: 1, threshold: 0, unknownFields: UF }], unknownFields: UF }], priority: 0, unknownFields: UF }],
          details: { domain: 'RuleDomainTransfer', subDomain: '', blockchain: '', network: '', unknownFields: UF },
          unknownFields: UF,
        },
      ],
      addressWhitelistingRules: [
        { currency: 'ETH', network: '', parallelThresholds: [], lines: [{ cells: [], parallelThresholds: [], unknownFields: UF }], unknownFields: UF },
      ],
      contractAddressWhitelistingRules: [{ blockchain: 'ETH', network: '', parallelThresholds: [], unknownFields: UF }],
    };

    const out = roundtrip(c);
    expect(hasUnknownFields(out)).toBe(true);

    const nodes: Record<string, UnknownFields | undefined> = {
      container: out.unknownFields,
      user: out.users[0].unknownFields,
      group: out.groups[0].unknownFields,
      transactionRule: out.transactionRules[0].unknownFields,
      column: out.transactionRules[0].columns[0].unknownFields,
      line: out.transactionRules[0].lines[0].unknownFields,
      details: out.transactionRules[0].details!.unknownFields,
      sequentialThresholds: out.transactionRules[0].lines[0].parallelThresholds[0].unknownFields,
      groupThreshold: out.transactionRules[0].lines[0].parallelThresholds[0].thresholds[0].unknownFields,
      awr: out.addressWhitelistingRules[0].unknownFields,
      awrLine: out.addressWhitelistingRules[0].lines[0].unknownFields,
      cawr: out.contractAddressWhitelistingRules[0].unknownFields,
    };
    const missing = Object.entries(nodes)
      .filter(([, uf]) => !uf || Object.keys(uf).length === 0)
      .map(([label]) => label);
    expect(missing).toEqual([]);
    // byte-stable re-encode
    expect(Buffer.from(rulesContainerToBytes(out))).toEqual(Buffer.from(rulesContainerToBytes(c)));
  });
});

describe('enum passthrough and named errors', () => {
  it('passes wire-decoded numeric enum values through numerically', () => {
    const c: DecodedRulesContainer = {
      ...createEmptyRulesContainer(),
      users: [{ id: 'u', name: undefined, publicKeyPem: undefined, roles: ['992'] }],
      transactionRules: [{ key: 'k', columns: [], lines: [], details: { domain: '993', subDomain: '994', blockchain: '', network: '' } }],
      contractAddressWhitelistingRules: [{ blockchain: '995', network: '', parallelThresholds: [] }],
    };
    const out = roundtrip(c);
    expect(out.users[0].roles).toEqual(['992']);
    expect(out.transactionRules[0].details!.domain).toBe('993');
    expect(out.transactionRules[0].details!.subDomain).toBe('994');
    expect(out.contractAddressWhitelistingRules[0].blockchain).toBe('995');
  });

  it('throws on caller-authored unknown enum names', () => {
    const role: DecodedRulesContainer = { ...createEmptyRulesContainer(), users: [{ id: 'u', name: undefined, publicKeyPem: undefined, roles: ['NotARole'] }] };
    expect(() => rulesContainerToBytes(role)).toThrow(/unknown user role/);

    const domain: DecodedRulesContainer = { ...createEmptyRulesContainer(), transactionRules: [{ key: 'k', columns: [], lines: [], details: { domain: 'Nope', subDomain: '', blockchain: '', network: '' } }] };
    expect(() => rulesContainerToBytes(domain)).toThrow(/unknown rule domain/);

    const chain: DecodedRulesContainer = { ...createEmptyRulesContainer(), contractAddressWhitelistingRules: [{ blockchain: 'NopeChain', network: '', parallelThresholds: [] }] };
    expect(() => rulesContainerToBytes(chain)).toThrow(/unknown blockchain/);
  });
});

describe('RuleSource variant codec', () => {
  it('round-trips the InternalAddress / Exchange / ExternalAddress / AnyExchange variants', () => {
    const c: DecodedRulesContainer = {
      ...createEmptyRulesContainer(),
      addressWhitelistingRules: [
        {
          currency: 'ETH',
          network: '',
          parallelThresholds: [],
          lines: [
            {
              cells: [
                { type: RuleSourceType.InternalAddress, internalAddress: { address: '0xabc', path: "m/44'/60'/0'/0/0" } },
                { type: RuleSourceType.Exchange, exchange: { label: 'kraken' } },
                { type: RuleSourceType.ExternalAddress, externalAddress: { address: '0xdef', memo: 'm' } },
                { type: RuleSourceType.AnyExchange },
              ],
              parallelThresholds: [],
            },
          ],
        },
      ],
    };
    const cells = roundtrip(c).addressWhitelistingRules[0].lines[0].cells;
    expect(cells.map((s) => s.type)).toEqual([
      RuleSourceType.InternalAddress,
      RuleSourceType.Exchange,
      RuleSourceType.ExternalAddress,
      RuleSourceType.AnyExchange,
    ]);
    expect(cells[0].internalAddress).toEqual({ address: '0xabc', path: "m/44'/60'/0'/0/0" });
    expect(cells[1].exchange).toEqual({ label: 'kraken' });
    expect(cells[2].externalAddress).toEqual({ address: '0xdef', memo: 'm' });
  });

  it('throws for an unknown source type without raw bytes', () => {
    const c: DecodedRulesContainer = {
      ...createEmptyRulesContainer(),
      addressWhitelistingRules: [
        { currency: 'ETH', network: '', parallelThresholds: [], lines: [{ cells: [{ type: 99 as RuleSourceType }], parallelThresholds: [] }] },
      ],
    };
    expect(() => rulesContainerToBytes(c)).toThrow(/unknown rule source type/);
  });

  it('preserves a malformed source payload as raw', () => {
    const bad = PbRuleSource.encode(
      PbRuleSource.fromPartial({ type: RuleSource_RuleSourceType.RuleSourceInternalWallet, payload: new Uint8Array([0xff, 0xff, 0xff, 0xff]) })
    ).finish();
    const pbC = RulesContainer.fromPartial({ addressWhitelistingRules: [{ currency: 'ETH', lines: [{ cells: [bad] }] }] });
    const out = tryDecodeProtobufRulesContainer(RulesContainer.encode(pbC).finish())!;
    const src = out.addressWhitelistingRules[0].lines[0].cells[0];
    expect(bytesOf(src.raw)).toEqual(Array.from(bad));
    // survives re-encode
    expect(bytesOf(roundtrip(out).addressWhitelistingRules[0].lines[0].cells[0].raw)).toEqual(Array.from(bad));
  });

  it('preserves a source type newer than this SDK as raw', () => {
    const future = PbRuleSource.encode(PbRuleSource.fromPartial({ type: 99 as RuleSource_RuleSourceType, payload: new Uint8Array([1, 2, 3]) })).finish();
    const pbC = RulesContainer.fromPartial({ addressWhitelistingRules: [{ currency: 'ETH', lines: [{ cells: [future] }] }] });
    const out = tryDecodeProtobufRulesContainer(RulesContainer.encode(pbC).finish())!;
    expect(hasUnknownFields(out)).toBe(true);
    expect(bytesOf(out.addressWhitelistingRules[0].lines[0].cells[0].raw)).toEqual(Array.from(future));
  });
});

describe('details sub-messages', () => {
  it('round-trips xtz / cash / cosmos details', () => {
    const c: DecodedRulesContainer = {
      ...createEmptyRulesContainer(),
      transactionRules: [
        {
          key: 'k',
          columns: [],
          lines: [],
          details: {
            domain: 'RuleDomainCallContract',
            subDomain: '',
            blockchain: '',
            network: '',
            xtzCallContract: { contractType: 'FA2', methodSignature: 'transfer' },
            cashSettlement: { provider: 'prov', requestType: 'settle' },
            cosmosDetails: { methodSignatures: ['/cosmos.bank.v1beta1.MsgSend'] },
          },
        },
      ],
    };
    const d = roundtrip(c).transactionRules[0].details!;
    expect(d.xtzCallContract).toEqual({ contractType: 'FA2', methodSignature: 'transfer' });
    expect(d.cashSettlement).toEqual({ provider: 'prov', requestType: 'settle' });
    expect(d.cosmosDetails).toEqual({ methodSignatures: ['/cosmos.bank.v1beta1.MsgSend'] });
  });
});
