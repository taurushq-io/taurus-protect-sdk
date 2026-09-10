/**
 * Branch/edge coverage for the cell codec that the golden vectors exercise only
 * transitively: error paths, the RawCell fallback, the magnitude encoding, and
 * cellFamily.
 */
import type { RuleCell } from '../../../src/models/rule-cell';
import { cellFamily, ruleCellFromBytes, ruleCellToBytes } from '../../../src/mappers/rule-cell-codec';
import {
  Blockchain,
  RuleDestination,
  RuleDestination_RuleDestinationType,
  RuleDestinationContractAddress,
  RuleFiatAmount,
  RuleFiatAmount_RuleFiatAmountType,
  RuleIntegerGreater,
  RuleIntegerGreater_RuleIntegerGreaterType,
  RuleUIntegerGreater,
  RuleUIntegerGreater_RuleUIntegerGreaterType,
} from '../../../src/internal/proto/request_reply';

const bytes = (u: Uint8Array): number[] => Array.from(u);

// --- Error branches ---

describe('rule cell codec error branches', () => {
  it('throws when a cell is placed in the wrong column family', () => {
    expect(() => ruleCellToBytes('RuleSource', { kind: 'FiatAmountAny' })).toThrow(/not valid for column type/);
  });

  it('throws for a negative unsigned integer', () => {
    expect(() => ruleCellToBytes('RuleUIntegerGreater', { kind: 'UIntegerGreaterValue', value: -1n })).toThrow(/non-negative/);
    expect(() => ruleCellToBytes('RuleUIntegerGreater', { kind: 'UIntegerGreaterIsEqual', value: -1n })).toThrow(/non-negative/);
  });
});

// --- RawCell fallback branches ---

describe('rule cell codec RawCell fallback', () => {
  it('preserves an unknown column type verbatim (empty and non-empty)', () => {
    const nonEmpty = ruleCellFromBytes('RuleUnknownColumn', new Uint8Array([0x08, 0x01]));
    expect(nonEmpty.kind).toBe('RawCell');
    expect(bytes(ruleCellToBytes('RuleUnknownColumn', nonEmpty))).toEqual([0x08, 0x01]);

    const empty = ruleCellFromBytes('RuleUnknownColumn', new Uint8Array());
    expect(empty.kind).toBe('RawCell');
    expect(bytes(ruleCellToBytes('RuleUnknownColumn', empty))).toEqual([]);
  });

  it('preserves a wrapper with an unknown sub-field verbatim (lossless guard)', () => {
    const data = new Uint8Array([0xc0, 0x0c, 0x2a]); // field 200, varint 42, on an otherwise-empty RuleFiatAmount
    const cell = ruleCellFromBytes('RuleFiatAmount', data);
    expect(cell.kind).toBe('RawCell');
    expect(bytes(ruleCellToBytes('RuleFiatAmount', cell))).toEqual([0xc0, 0x0c, 0x2a]);
  });

  it('preserves a payload on a payload-free cell type verbatim', () => {
    const data = RuleFiatAmount.encode(
      RuleFiatAmount.fromPartial({ type: RuleFiatAmount_RuleFiatAmountType.RuleFiatAmountAny, payload: new Uint8Array([1, 2]) })
    ).finish();
    const cell = ruleCellFromBytes('RuleFiatAmount', data);
    expect(cell.kind).toBe('RawCell');
    expect(bytes(ruleCellToBytes('RuleFiatAmount', cell))).toEqual(bytes(data));
  });

  it('preserves an unknown cell-type enum across every family', () => {
    // Cell-type enum value newer than this SDK (cast to inject an out-of-range value).
    const cases: Array<[string, Uint8Array]> = [
      ['RuleFiatAmount', RuleFiatAmount.encode(RuleFiatAmount.fromPartial({ type: 902 as RuleFiatAmount_RuleFiatAmountType })).finish()],
      ['RuleIntegerGreater', RuleIntegerGreater.encode(RuleIntegerGreater.fromPartial({ type: 902 as RuleIntegerGreater_RuleIntegerGreaterType })).finish()],
      ['RuleUIntegerGreater', RuleUIntegerGreater.encode(RuleUIntegerGreater.fromPartial({ type: 902 as RuleUIntegerGreater_RuleUIntegerGreaterType })).finish()],
    ];
    for (const [col, data] of cases) {
      const cell = ruleCellFromBytes(col, data);
      expect(cell.kind).toBe('RawCell');
      expect(bytes(ruleCellToBytes(col, cell))).toEqual(bytes(data));
    }
  });

});

// --- Blockchain numeric passthrough ---

describe('rule cell codec blockchain numeric passthrough', () => {
  it('round-trips a contract-address cell with an unknown blockchain as a typed cell', () => {
    // ts-proto collapses an unknown blockchain enum to UNRECOGNIZED; the codec's
    // numeric passthrough keeps the raw number as a decimal string so the cell
    // stays typed and re-encodes byte-identically.
    const inner = RuleDestinationContractAddress.encode(
      RuleDestinationContractAddress.fromPartial({ address: '0x1', blockchain: 4242 as Blockchain })
    ).finish();
    const data = RuleDestination.encode(
      RuleDestination.fromPartial({ type: RuleDestination_RuleDestinationType.RuleDestinationContractAddress, payload: inner })
    ).finish();
    const cell = ruleCellFromBytes('RuleDestination', data);
    expect(cell).toEqual({ kind: 'DestinationContractAddress', address: '0x1', name: '', symbol: '', blockchain: '4242' });
    expect(bytes(ruleCellToBytes('RuleDestination', cell))).toEqual(bytes(data));
  });

  it('throws when encoding a contract-address cell with an unparseable blockchain name', () => {
    expect(() =>
      ruleCellToBytes('RuleDestination', { kind: 'DestinationContractAddress', address: '0x1', name: '', symbol: '', blockchain: 'NotAChain' })
    ).toThrow(/unknown blockchain/);
  });
});

// --- Magnitude encoding (via the public codec) ---

describe('rule cell codec magnitude encoding', () => {
  it('encodes the big-endian magnitude and round-trips zero', () => {
    const zero = ruleCellToBytes('RuleUIntegerGreater', { kind: 'UIntegerGreaterValue', value: 0n });
    expect(ruleCellFromBytes('RuleUIntegerGreater', zero)).toEqual({ kind: 'UIntegerGreaterValue', value: 0n });

    const big = RuleUIntegerGreater.decode(ruleCellToBytes('RuleUIntegerGreater', { kind: 'UIntegerGreaterValue', value: 256n }));
    expect(bytes(big.payload)).toEqual([0x01, 0x00]);

    // odd-length hex is left-padded to whole bytes
    const small = RuleIntegerGreater.decode(ruleCellToBytes('RuleIntegerGreater', { kind: 'IntegerGreaterValue', value: 8n }));
    expect(bytes(small.payload)).toEqual([0x08]);
  });

  it('selects the NegValue arm for negatives (sign carried by the arm, magnitude in the payload)', () => {
    const neg = RuleIntegerGreater.decode(ruleCellToBytes('RuleIntegerGreater', { kind: 'IntegerGreaterValue', value: -50n }));
    expect(neg.payload && bytes(neg.payload)).toEqual([50]);
    expect(ruleCellFromBytes('RuleIntegerGreater', ruleCellToBytes('RuleIntegerGreater', { kind: 'IntegerGreaterValue', value: -50n })))
      .toEqual({ kind: 'IntegerGreaterValue', value: -50n });
  });
});

// --- cellFamily (direct) ---

describe('cellFamily', () => {
  it('maps each cell kind to its column family', () => {
    const cases: Array<[RuleCell, string]> = [
      [{ kind: 'FiatAmountRange', minAmount: '1', maxAmount: '2' }, 'RuleFiatAmount'],
      [{ kind: 'SourceInternalWallet', path: 'p' }, 'RuleSource'],
      [{ kind: 'DestinationAny' }, 'RuleDestination'],
      [{ kind: 'StringEqualValue', value: 'v' }, 'RuleStringEqual'],
      [{ kind: 'BytesEqualValue', value: new Uint8Array() }, 'RuleBytesEqual'],
      [{ kind: 'StringArrayEqualValue', values: [] }, 'RuleStringArrayEqual'], // distinct from RuleStringEqual
      [{ kind: 'IntegerGreaterValue', value: 1n }, 'RuleIntegerGreater'],
      [{ kind: 'UIntegerGreaterValue', value: 1n }, 'RuleUIntegerGreater'], // distinct from RuleIntegerGreater
      [{ kind: 'WhitelistedContractAny' }, 'RuleWhitelistedContract'],
    ];
    for (const [cell, want] of cases) {
      expect(cellFamily(cell)).toBe(want);
    }
    // RawCell reports its own column type.
    expect(cellFamily({ kind: 'RawCell', columnType: 'RuleFiatAmount', payload: new Uint8Array() })).toBe('RuleFiatAmount');
  });
});
