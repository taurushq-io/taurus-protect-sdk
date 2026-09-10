/**
 * Cross-SDK cell wire-format alignment.
 *
 * Every entry in the shared golden-vector file (produced by the Go SDK) must
 * encode to exactly the recorded bytes and decode back to the same typed cell.
 * The same file is consumed by the Go/Java/Python suites, so this pins
 * byte-for-byte parity across all four SDKs.
 */
import * as fs from 'fs';
import * as path from 'path';
import type { RuleCell } from '../../../src/models/rule-cell';
import { ruleCellFromBytes, ruleCellToBytes } from '../../../src/mappers/rule-cell-codec';

const VECTORS_PATH = path.resolve(
  __dirname,
  '../../../../scripts/resources/governance-cell-vectors.json'
);

interface CellVector {
  description: string;
  column_type: string;
  cell_type: string;
  typed_value_json: unknown;
  wire_base64: string;
}

const RAW_DESCRIPTION = 'raw cell (unknown cell type preserved verbatim)';

// Hardcoded cell cases keyed by the shared vectors' `description`, mirroring the
// Go SDK's allCellCases table. The vectors file is the byte oracle.
const CASES: Record<string, { columnType: string; cell: RuleCell }> = {
  'fiat amount any': { columnType: 'RuleFiatAmount', cell: { kind: 'FiatAmountAny' } },
  'fiat amount is zero': { columnType: 'RuleFiatAmount', cell: { kind: 'FiatAmountIsZero' } },
  'fiat amount range': { columnType: 'RuleFiatAmount', cell: { kind: 'FiatAmountRange', minAmount: '1000', maxAmount: '50000' } },

  'source any': { columnType: 'RuleSource', cell: { kind: 'SourceAny' } },
  'source internal wallet': { columnType: 'RuleSource', cell: { kind: 'SourceInternalWallet', path: "m/44'/60'/0'" } },
  'source internal address': { columnType: 'RuleSource', cell: { kind: 'SourceInternalAddress', address: '0xabc', path: "m/44'/60'/0'/0/0" } },
  'source any exchange': { columnType: 'RuleSource', cell: { kind: 'SourceAnyExchange' } },
  'source exchange': { columnType: 'RuleSource', cell: { kind: 'SourceExchange', label: 'kraken-main' } },
  'source external address': { columnType: 'RuleSource', cell: { kind: 'SourceExternalAddress', address: '0xdef', memo: 'memo-1' } },

  'destination any': { columnType: 'RuleDestination', cell: { kind: 'DestinationAny' } },
  'destination internal wallet': { columnType: 'RuleDestination', cell: { kind: 'DestinationInternalWallet', path: "m/44'/60'/1'" } },
  'destination internal address': { columnType: 'RuleDestination', cell: { kind: 'DestinationInternalAddress', address: '0x111', path: "m/44'/60'/1'/0/0" } },
  'destination external address': { columnType: 'RuleDestination', cell: { kind: 'DestinationExternalAddress', address: '0x222', memo: 'dest-memo' } },
  'destination any exchange': { columnType: 'RuleDestination', cell: { kind: 'DestinationAnyExchange' } },
  'destination exchange': { columnType: 'RuleDestination', cell: { kind: 'DestinationExchange', label: 'binance-desk', memo: 'x' } },
  'destination contract address': { columnType: 'RuleDestination', cell: { kind: 'DestinationContractAddress', address: '0x333', name: 'USDC', symbol: 'USDC', blockchain: 'ETH' } },
  'destination contract address (unknown blockchain)': { columnType: 'RuleDestination', cell: { kind: 'DestinationContractAddress', address: '0x1', name: '', symbol: '', blockchain: '4242' } },
  'destination any external address': { columnType: 'RuleDestination', cell: { kind: 'DestinationAnyExternalAddress' } },
  'destination any contract address': { columnType: 'RuleDestination', cell: { kind: 'DestinationAnyContractAddress' } },

  'string equal any': { columnType: 'RuleStringEqual', cell: { kind: 'StringEqualAny' } },
  'string equal empty': { columnType: 'RuleStringEqual', cell: { kind: 'StringEqualEmpty' } },
  'string equal value': { columnType: 'RuleStringEqual', cell: { kind: 'StringEqualValue', value: 'contract-id-42' } },

  'bytes equal any': { columnType: 'RuleBytesEqual', cell: { kind: 'BytesEqualAny' } },
  'bytes equal empty': { columnType: 'RuleBytesEqual', cell: { kind: 'BytesEqualEmpty' } },
  'bytes equal value': { columnType: 'RuleBytesEqual', cell: { kind: 'BytesEqualValue', value: new Uint8Array([0xde, 0xad, 0xbe, 0xef]) } },

  'string array equal any': { columnType: 'RuleStringArrayEqual', cell: { kind: 'StringArrayEqualAny' } },
  'string array equal empty': { columnType: 'RuleStringArrayEqual', cell: { kind: 'StringArrayEqualEmpty' } },
  'string array equal value': { columnType: 'RuleStringArrayEqual', cell: { kind: 'StringArrayEqualValue', values: ['a', 'b', 'c'] } },

  'integer greater any': { columnType: 'RuleIntegerGreater', cell: { kind: 'IntegerGreaterAny' } },
  'integer greater positive value': { columnType: 'RuleIntegerGreater', cell: { kind: 'IntegerGreaterValue', value: 50n } },
  'integer greater negative value': { columnType: 'RuleIntegerGreater', cell: { kind: 'IntegerGreaterValue', value: -50n } },
  'integer greater zero': { columnType: 'RuleIntegerGreater', cell: { kind: 'IntegerGreaterValue', value: 0n } },

  'uinteger greater any': { columnType: 'RuleUIntegerGreater', cell: { kind: 'UIntegerGreaterAny' } },
  'uinteger greater is zero': { columnType: 'RuleUIntegerGreater', cell: { kind: 'UIntegerGreaterIsZero' } },
  'uinteger greater value': { columnType: 'RuleUIntegerGreater', cell: { kind: 'UIntegerGreaterValue', value: 18446744073709551615n } },
  'uinteger greater is equal': { columnType: 'RuleUIntegerGreater', cell: { kind: 'UIntegerGreaterIsEqual', value: 10n } },

  'whitelisted contract any': { columnType: 'RuleWhitelistedContract', cell: { kind: 'WhitelistedContractAny' } },
  'whitelisted contract address': { columnType: 'RuleWhitelistedContract', cell: { kind: 'WhitelistedContractAddress', address: '0x444', name: 'DAI', symbol: 'DAI', blockchain: 'ETH' } },
};

function toB64(bytes: Uint8Array): string {
  return Buffer.from(bytes).toString('base64');
}
function fromB64(s: string): Uint8Array {
  return new Uint8Array(Buffer.from(s, 'base64'));
}

describe('governance cell golden vectors (cross-SDK parity)', () => {
  const vectors: CellVector[] = JSON.parse(fs.readFileSync(VECTORS_PATH, 'utf-8'));

  it('covers every hardcoded case (grammar completeness)', () => {
    const described = new Set(vectors.map((v) => v.description));
    for (const desc of Object.keys(CASES)) {
      expect(described.has(desc)).toBe(true);
    }
    // vectors = the typed cases + one raw case
    expect(vectors.length).toBe(Object.keys(CASES).length + 1);
  });

  for (const v of vectorsFor()) {
    it(`encodes + decodes: ${v.description}`, () => {
      if (v.description === RAW_DESCRIPTION) {
        const wire = fromB64(v.wire_base64);
        const decoded = ruleCellFromBytes(v.column_type, wire);
        expect(decoded.kind).toBe('RawCell');
        expect(toB64(ruleCellToBytes(v.column_type, decoded))).toBe(v.wire_base64);
        return;
      }
      const c = CASES[v.description];
      expect(c).toBeDefined();
      expect(c.columnType).toBe(v.column_type);
      // encode(typed) === recorded wire bytes (cross-SDK byte parity)
      expect(toB64(ruleCellToBytes(c.columnType, c.cell))).toBe(v.wire_base64);
      // decode(wire) === typed
      expect(ruleCellFromBytes(c.columnType, fromB64(v.wire_base64))).toEqual(c.cell);
    });
  }

  // Read vectors once at module load for the generated `it` blocks.
  function vectorsFor(): CellVector[] {
    return JSON.parse(fs.readFileSync(VECTORS_PATH, 'utf-8'));
  }
});
