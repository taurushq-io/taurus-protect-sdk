/**
 * Codec between the typed {@link RuleCell} union and the wire bytes stored in a
 * transaction rule line's cells (mirrors the Go SDK's rule_cell_codec.go).
 *
 * A cell is a wrapper message selected by the column's type,
 * `{ type: <cell-type enum>, payload: <bytes> }`. Message-family payloads nest
 * a serialized sub-message; RuleStringEqual/RuleBytesEqual carry raw scalar
 * bytes; the integer families carry big-endian magnitude bytes with the sign
 * expressed by the enum arm. The `*Any` cells are the protobuf zero values, so
 * their serialized form is the empty cell.
 *
 * Decode never fails: any cell whose type/content this SDK version cannot
 * represent losslessly (unknown column type, unknown cell type, unknown
 * protobuf sub-fields) round-trips as a {@link RuleCell} of kind `RawCell`,
 * detected by a decode → re-encode → byte-compare check.
 */
import type { RuleCell } from '../models/rule-cell';
import {
  Blockchain,
  blockchainFromJSON,
  blockchainToJSON,
  RuleBytesEqual,
  RuleBytesEqual_RuleBytesEqualType,
  RuleDestination,
  RuleDestination_RuleDestinationType,
  RuleDestinationContractAddress,
  RuleDestinationExchange,
  RuleDestinationExternalAddress,
  RuleDestinationInternalAddress,
  RuleDestinationInternalWallet,
  RuleFiatAmount,
  RuleFiatAmount_RuleFiatAmountType,
  RuleFiatAmountRange,
  RuleIntegerGreater,
  RuleIntegerGreater_RuleIntegerGreaterType,
  RuleSource,
  RuleSource_RuleSourceType,
  RuleSourceExchange,
  RuleSourceExternalAddress,
  RuleSourceInternalAddress,
  RuleSourceInternalWallet,
  RuleStringArrayEqual,
  RuleStringArrayEqual_RuleStringArrayEqualType,
  RuleStringArrayEqualValue,
  RuleStringEqual,
  RuleStringEqual_RuleStringEqualType,
  RuleUIntegerGreater,
  RuleUIntegerGreater_RuleUIntegerGreaterType,
  RuleWhitelistedContract,
  RuleWhitelistedContract_RuleWhitelistedContractType,
} from '../internal/proto/request_reply';

const utf8 = new TextEncoder();
const utf8dec = new TextDecoder();

function bytesEqual(a: Uint8Array, b: Uint8Array): boolean {
  if (a.length !== b.length) return false;
  for (let i = 0; i < a.length; i++) if (a[i] !== b[i]) return false;
  return true;
}

/** Big-endian magnitude bytes of a non-negative bigint; zero → [0] (matches the ecosystem BigIntToBytes convention). */
function bigintToMagnitude(v: bigint): Uint8Array {
  if (v < 0n) v = -v;
  if (v === 0n) return new Uint8Array([0]);
  let hex = v.toString(16);
  if (hex.length % 2 === 1) hex = '0' + hex;
  const out = new Uint8Array(hex.length / 2);
  for (let i = 0; i < out.length; i++) out[i] = parseInt(hex.slice(i * 2, i * 2 + 2), 16);
  return out;
}

function magnitudeToBigint(b: Uint8Array): bigint {
  let v = 0n;
  for (const byte of b) v = (v << 8n) | BigInt(byte);
  return v;
}

/** The column type a typed cell belongs to. */
export function cellFamily(cell: RuleCell): string {
  switch (cell.kind) {
    case 'FiatAmountAny':
    case 'FiatAmountIsZero':
    case 'FiatAmountRange':
      return 'RuleFiatAmount';
    case 'SourceAny':
    case 'SourceInternalWallet':
    case 'SourceInternalAddress':
    case 'SourceAnyExchange':
    case 'SourceExchange':
    case 'SourceExternalAddress':
      return 'RuleSource';
    case 'DestinationAny':
    case 'DestinationInternalWallet':
    case 'DestinationInternalAddress':
    case 'DestinationExternalAddress':
    case 'DestinationAnyExchange':
    case 'DestinationExchange':
    case 'DestinationContractAddress':
    case 'DestinationAnyExternalAddress':
    case 'DestinationAnyContractAddress':
      return 'RuleDestination';
    case 'StringEqualAny':
    case 'StringEqualEmpty':
    case 'StringEqualValue':
      return 'RuleStringEqual';
    case 'BytesEqualAny':
    case 'BytesEqualEmpty':
    case 'BytesEqualValue':
      return 'RuleBytesEqual';
    case 'StringArrayEqualAny':
    case 'StringArrayEqualEmpty':
    case 'StringArrayEqualValue':
      return 'RuleStringArrayEqual';
    case 'IntegerGreaterAny':
    case 'IntegerGreaterValue':
      return 'RuleIntegerGreater';
    case 'UIntegerGreaterAny':
    case 'UIntegerGreaterIsZero':
    case 'UIntegerGreaterValue':
    case 'UIntegerGreaterIsEqual':
      return 'RuleUIntegerGreater';
    case 'WhitelistedContractAny':
    case 'WhitelistedContractAddress':
      return 'RuleWhitelistedContract';
    case 'RawCell':
      return cell.columnType;
  }
}

function blockchainToStr(b: Blockchain): string {
  // ts-proto's blockchainToJSON collapses every unknown enum to "UNRECOGNIZED",
  // which fails the cell's re-encode guard and demotes it to RawCell. Pass the
  // raw number through as a decimal string so an unknown blockchain stays typed.
  const name = blockchainToJSON(b);
  return name === 'UNRECOGNIZED' ? String(b) : name;
}
function blockchainFromStr(name: string): Blockchain {
  const b = blockchainFromJSON(name);
  if (b !== Blockchain.UNRECOGNIZED) return b;
  // Decimal passthrough for values newer than this SDK; a non-numeric unknown
  // name is a caller error.
  if (/^-?\d+$/.test(name)) return parseInt(name, 10) as Blockchain;
  throw new Error(`unknown blockchain "${name}"`);
}

/**
 * Encodes a typed cell into the wrapped protobuf bytes stored in a rule line,
 * where `colType` is the type of the column the cell is aligned with.
 */
export function ruleCellToBytes(colType: string, cell: RuleCell): Uint8Array {
  if (cell.kind === 'RawCell') {
    return cell.payload;
  }
  const family = cellFamily(cell);
  if (colType !== '' && colType !== family) {
    throw new Error(`cell ${cell.kind} is not valid for column type ${colType}`);
  }

  switch (cell.kind) {
    // RuleFiatAmount
    case 'FiatAmountAny':
      return RuleFiatAmount.encode(RuleFiatAmount.fromPartial({ type: RuleFiatAmount_RuleFiatAmountType.RuleFiatAmountAny })).finish();
    case 'FiatAmountIsZero':
      return RuleFiatAmount.encode(RuleFiatAmount.fromPartial({ type: RuleFiatAmount_RuleFiatAmountType.RuleFiatAmountIsZero })).finish();
    case 'FiatAmountRange': {
      const inner = RuleFiatAmountRange.encode(RuleFiatAmountRange.fromPartial({ minAmount: cell.minAmount, maxAmount: cell.maxAmount })).finish();
      return RuleFiatAmount.encode(RuleFiatAmount.fromPartial({ type: RuleFiatAmount_RuleFiatAmountType.RuleFiatAmountRange, payload: inner })).finish();
    }

    // RuleSource
    case 'SourceAny':
      return RuleSource.encode(RuleSource.fromPartial({ type: RuleSource_RuleSourceType.RuleSourceAny })).finish();
    case 'SourceAnyExchange':
      return RuleSource.encode(RuleSource.fromPartial({ type: RuleSource_RuleSourceType.RuleSourceAnyExchange })).finish();
    case 'SourceInternalWallet': {
      const inner = RuleSourceInternalWallet.encode(RuleSourceInternalWallet.fromPartial({ path: cell.path })).finish();
      return RuleSource.encode(RuleSource.fromPartial({ type: RuleSource_RuleSourceType.RuleSourceInternalWallet, payload: inner })).finish();
    }
    case 'SourceInternalAddress': {
      const inner = RuleSourceInternalAddress.encode(RuleSourceInternalAddress.fromPartial({ address: cell.address, path: cell.path })).finish();
      return RuleSource.encode(RuleSource.fromPartial({ type: RuleSource_RuleSourceType.RuleSourceInternalAddress, payload: inner })).finish();
    }
    case 'SourceExchange': {
      const inner = RuleSourceExchange.encode(RuleSourceExchange.fromPartial({ label: cell.label })).finish();
      return RuleSource.encode(RuleSource.fromPartial({ type: RuleSource_RuleSourceType.RuleSourceExchange, payload: inner })).finish();
    }
    case 'SourceExternalAddress': {
      const inner = RuleSourceExternalAddress.encode(RuleSourceExternalAddress.fromPartial({ address: cell.address, memo: cell.memo })).finish();
      return RuleSource.encode(RuleSource.fromPartial({ type: RuleSource_RuleSourceType.RuleSourceExternalAddress, payload: inner })).finish();
    }

    // RuleDestination
    case 'DestinationAny':
      return RuleDestination.encode(RuleDestination.fromPartial({ type: RuleDestination_RuleDestinationType.RuleDestinationAny })).finish();
    case 'DestinationAnyExchange':
      return RuleDestination.encode(RuleDestination.fromPartial({ type: RuleDestination_RuleDestinationType.RuleDestinationAnyExchange })).finish();
    case 'DestinationAnyExternalAddress':
      return RuleDestination.encode(RuleDestination.fromPartial({ type: RuleDestination_RuleDestinationType.RuleDestinationAnyExternalAddress })).finish();
    case 'DestinationAnyContractAddress':
      return RuleDestination.encode(RuleDestination.fromPartial({ type: RuleDestination_RuleDestinationType.RuleDestinationAnyContractAddress })).finish();
    case 'DestinationInternalWallet': {
      const inner = RuleDestinationInternalWallet.encode(RuleDestinationInternalWallet.fromPartial({ path: cell.path })).finish();
      return RuleDestination.encode(RuleDestination.fromPartial({ type: RuleDestination_RuleDestinationType.RuleDestinationInternalWallet, payload: inner })).finish();
    }
    case 'DestinationInternalAddress': {
      const inner = RuleDestinationInternalAddress.encode(RuleDestinationInternalAddress.fromPartial({ address: cell.address, path: cell.path })).finish();
      return RuleDestination.encode(RuleDestination.fromPartial({ type: RuleDestination_RuleDestinationType.RuleDestinationInternalAddress, payload: inner })).finish();
    }
    case 'DestinationExternalAddress': {
      const inner = RuleDestinationExternalAddress.encode(RuleDestinationExternalAddress.fromPartial({ address: cell.address, memo: cell.memo })).finish();
      return RuleDestination.encode(RuleDestination.fromPartial({ type: RuleDestination_RuleDestinationType.RuleDestinationExternalAddress, payload: inner })).finish();
    }
    case 'DestinationExchange': {
      const inner = RuleDestinationExchange.encode(RuleDestinationExchange.fromPartial({ label: cell.label, memo: cell.memo })).finish();
      return RuleDestination.encode(RuleDestination.fromPartial({ type: RuleDestination_RuleDestinationType.RuleDestinationExchange, payload: inner })).finish();
    }
    case 'DestinationContractAddress': {
      const inner = RuleDestinationContractAddress.encode(
        RuleDestinationContractAddress.fromPartial({ address: cell.address, name: cell.name, symbol: cell.symbol, blockchain: blockchainFromStr(cell.blockchain) })
      ).finish();
      return RuleDestination.encode(RuleDestination.fromPartial({ type: RuleDestination_RuleDestinationType.RuleDestinationContractAddress, payload: inner })).finish();
    }

    // RuleWhitelistedContract
    case 'WhitelistedContractAny':
      return RuleWhitelistedContract.encode(RuleWhitelistedContract.fromPartial({ type: RuleWhitelistedContract_RuleWhitelistedContractType.RuleWhitelistedContractAny })).finish();
    case 'WhitelistedContractAddress': {
      const inner = RuleDestinationContractAddress.encode(
        RuleDestinationContractAddress.fromPartial({ address: cell.address, name: cell.name, symbol: cell.symbol, blockchain: blockchainFromStr(cell.blockchain) })
      ).finish();
      return RuleWhitelistedContract.encode(
        RuleWhitelistedContract.fromPartial({ type: RuleWhitelistedContract_RuleWhitelistedContractType.RuleWhitelistedContract_RuleDestinationContractAddress, payload: inner })
      ).finish();
    }

    // RuleStringEqual (payload = raw string bytes)
    case 'StringEqualAny':
      return RuleStringEqual.encode(RuleStringEqual.fromPartial({ type: RuleStringEqual_RuleStringEqualType.RuleStringEqualAny })).finish();
    case 'StringEqualEmpty':
      return RuleStringEqual.encode(RuleStringEqual.fromPartial({ type: RuleStringEqual_RuleStringEqualType.RuleStringEqualEmpty })).finish();
    case 'StringEqualValue':
      return RuleStringEqual.encode(RuleStringEqual.fromPartial({ type: RuleStringEqual_RuleStringEqualType.RuleStringEqualValue, payload: utf8.encode(cell.value) })).finish();

    // RuleBytesEqual (payload = raw bytes)
    case 'BytesEqualAny':
      return RuleBytesEqual.encode(RuleBytesEqual.fromPartial({ type: RuleBytesEqual_RuleBytesEqualType.RuleBytesEqualAny })).finish();
    case 'BytesEqualEmpty':
      return RuleBytesEqual.encode(RuleBytesEqual.fromPartial({ type: RuleBytesEqual_RuleBytesEqualType.RuleBytesEqualEmpty })).finish();
    case 'BytesEqualValue':
      return RuleBytesEqual.encode(RuleBytesEqual.fromPartial({ type: RuleBytesEqual_RuleBytesEqualType.RuleBytesEqualValue, payload: cell.value })).finish();

    // RuleStringArrayEqual
    case 'StringArrayEqualAny':
      return RuleStringArrayEqual.encode(RuleStringArrayEqual.fromPartial({ type: RuleStringArrayEqual_RuleStringArrayEqualType.RuleStringArrayEqualAny })).finish();
    case 'StringArrayEqualEmpty':
      return RuleStringArrayEqual.encode(RuleStringArrayEqual.fromPartial({ type: RuleStringArrayEqual_RuleStringArrayEqualType.RuleStringArrayEqualEmpty })).finish();
    case 'StringArrayEqualValue': {
      const inner = RuleStringArrayEqualValue.encode(RuleStringArrayEqualValue.fromPartial({ values: cell.values })).finish();
      return RuleStringArrayEqual.encode(RuleStringArrayEqual.fromPartial({ type: RuleStringArrayEqual_RuleStringArrayEqualType.RuleStringArrayEqualValue, payload: inner })).finish();
    }

    // RuleIntegerGreater (magnitude payload, sign selects the enum arm)
    case 'IntegerGreaterAny':
      return RuleIntegerGreater.encode(RuleIntegerGreater.fromPartial({ type: RuleIntegerGreater_RuleIntegerGreaterType.RuleIntegerGreaterAny })).finish();
    case 'IntegerGreaterValue': {
      const arm = cell.value < 0n
        ? RuleIntegerGreater_RuleIntegerGreaterType.RuleIntegerGreaterNegValue
        : RuleIntegerGreater_RuleIntegerGreaterType.RuleIntegerGreaterValue;
      return RuleIntegerGreater.encode(RuleIntegerGreater.fromPartial({ type: arm, payload: bigintToMagnitude(cell.value) })).finish();
    }

    // RuleUIntegerGreater
    case 'UIntegerGreaterAny':
      return RuleUIntegerGreater.encode(RuleUIntegerGreater.fromPartial({ type: RuleUIntegerGreater_RuleUIntegerGreaterType.RuleUIntegerGreaterAny })).finish();
    case 'UIntegerGreaterIsZero':
      return RuleUIntegerGreater.encode(RuleUIntegerGreater.fromPartial({ type: RuleUIntegerGreater_RuleUIntegerGreaterType.RuleUIntegerGreaterIsZero })).finish();
    case 'UIntegerGreaterValue': {
      if (cell.value < 0n) throw new Error('UIntegerGreaterValue requires a non-negative value');
      return RuleUIntegerGreater.encode(RuleUIntegerGreater.fromPartial({ type: RuleUIntegerGreater_RuleUIntegerGreaterType.RuleUIntegerGreaterValue, payload: bigintToMagnitude(cell.value) })).finish();
    }
    case 'UIntegerGreaterIsEqual': {
      if (cell.value < 0n) throw new Error('UIntegerGreaterIsEqual requires a non-negative value');
      return RuleUIntegerGreater.encode(RuleUIntegerGreater.fromPartial({ type: RuleUIntegerGreater_RuleUIntegerGreaterType.RuleUIntegerGreaterIsEqual, payload: bigintToMagnitude(cell.value) })).finish();
    }
  }
}

/**
 * Decodes one transaction-rule cell using the column type it aligns with. Any
 * cell that does not round-trip byte-identically through the typed layer is
 * returned as a `RawCell` so nothing is silently dropped.
 */
export function ruleCellFromBytes(colType: string, data: Uint8Array): RuleCell {
  const raw = (): RuleCell => ({ kind: 'RawCell', columnType: colType, payload: data });
  // Decoding never fails: one malformed cell must degrade to a RawCell rather than
  // abort the whole container, which would take every rule for that tenant down with
  // it. A cell payload is protobuf `bytes`, so a non-UTF-8 string cell and a truncated
  // wrapper are both legal on the wire.
  try {
    const typed = ruleCellFromBytesTyped(colType, data);
    if (typed === undefined) return raw();
    // Lossless guard: if the typed value does not re-encode to the exact input
    // bytes (unknown sub-fields, unknown enum value, non-canonical encoding),
    // preserve the original verbatim.
    if (!bytesEqual(ruleCellToBytes(colType, typed), data)) return raw();
    return typed;
  } catch {
    return raw();
  }
}

function ruleCellFromBytesTyped(colType: string, data: Uint8Array): RuleCell | undefined {
  switch (colType) {
    case 'RuleFiatAmount': {
      const w = RuleFiatAmount.decode(data);
      switch (w.type) {
        case RuleFiatAmount_RuleFiatAmountType.RuleFiatAmountAny: return { kind: 'FiatAmountAny' };
        case RuleFiatAmount_RuleFiatAmountType.RuleFiatAmountIsZero: return { kind: 'FiatAmountIsZero' };
        case RuleFiatAmount_RuleFiatAmountType.RuleFiatAmountRange: {
          const i = RuleFiatAmountRange.decode(w.payload);
          return { kind: 'FiatAmountRange', minAmount: i.minAmount, maxAmount: i.maxAmount };
        }
        default: return undefined;
      }
    }
    case 'RuleSource': {
      const w = RuleSource.decode(data);
      switch (w.type) {
        case RuleSource_RuleSourceType.RuleSourceAny: return { kind: 'SourceAny' };
        case RuleSource_RuleSourceType.RuleSourceAnyExchange: return { kind: 'SourceAnyExchange' };
        case RuleSource_RuleSourceType.RuleSourceInternalWallet: {
          const i = RuleSourceInternalWallet.decode(w.payload);
          return { kind: 'SourceInternalWallet', path: i.path };
        }
        case RuleSource_RuleSourceType.RuleSourceInternalAddress: {
          const i = RuleSourceInternalAddress.decode(w.payload);
          return { kind: 'SourceInternalAddress', address: i.address, path: i.path };
        }
        case RuleSource_RuleSourceType.RuleSourceExchange: {
          const i = RuleSourceExchange.decode(w.payload);
          return { kind: 'SourceExchange', label: i.label };
        }
        case RuleSource_RuleSourceType.RuleSourceExternalAddress: {
          const i = RuleSourceExternalAddress.decode(w.payload);
          return { kind: 'SourceExternalAddress', address: i.address, memo: i.memo };
        }
        default: return undefined;
      }
    }
    case 'RuleDestination': {
      const w = RuleDestination.decode(data);
      switch (w.type) {
        case RuleDestination_RuleDestinationType.RuleDestinationAny: return { kind: 'DestinationAny' };
        case RuleDestination_RuleDestinationType.RuleDestinationAnyExchange: return { kind: 'DestinationAnyExchange' };
        case RuleDestination_RuleDestinationType.RuleDestinationAnyExternalAddress: return { kind: 'DestinationAnyExternalAddress' };
        case RuleDestination_RuleDestinationType.RuleDestinationAnyContractAddress: return { kind: 'DestinationAnyContractAddress' };
        case RuleDestination_RuleDestinationType.RuleDestinationInternalWallet: {
          const i = RuleDestinationInternalWallet.decode(w.payload);
          return { kind: 'DestinationInternalWallet', path: i.path };
        }
        case RuleDestination_RuleDestinationType.RuleDestinationInternalAddress: {
          const i = RuleDestinationInternalAddress.decode(w.payload);
          return { kind: 'DestinationInternalAddress', address: i.address, path: i.path };
        }
        case RuleDestination_RuleDestinationType.RuleDestinationExternalAddress: {
          const i = RuleDestinationExternalAddress.decode(w.payload);
          return { kind: 'DestinationExternalAddress', address: i.address, memo: i.memo };
        }
        case RuleDestination_RuleDestinationType.RuleDestinationExchange: {
          const i = RuleDestinationExchange.decode(w.payload);
          return { kind: 'DestinationExchange', label: i.label, memo: i.memo };
        }
        case RuleDestination_RuleDestinationType.RuleDestinationContractAddress: {
          const i = RuleDestinationContractAddress.decode(w.payload);
          return { kind: 'DestinationContractAddress', address: i.address, name: i.name, symbol: i.symbol, blockchain: blockchainToStr(i.blockchain) };
        }
        default: return undefined;
      }
    }
    case 'RuleWhitelistedContract': {
      const w = RuleWhitelistedContract.decode(data);
      switch (w.type) {
        case RuleWhitelistedContract_RuleWhitelistedContractType.RuleWhitelistedContractAny: return { kind: 'WhitelistedContractAny' };
        case RuleWhitelistedContract_RuleWhitelistedContractType.RuleWhitelistedContract_RuleDestinationContractAddress: {
          const i = RuleDestinationContractAddress.decode(w.payload);
          return { kind: 'WhitelistedContractAddress', address: i.address, name: i.name, symbol: i.symbol, blockchain: blockchainToStr(i.blockchain) };
        }
        default: return undefined;
      }
    }
    case 'RuleStringEqual': {
      const w = RuleStringEqual.decode(data);
      switch (w.type) {
        case RuleStringEqual_RuleStringEqualType.RuleStringEqualAny: return { kind: 'StringEqualAny' };
        case RuleStringEqual_RuleStringEqualType.RuleStringEqualEmpty: return { kind: 'StringEqualEmpty' };
        case RuleStringEqual_RuleStringEqualType.RuleStringEqualValue: return { kind: 'StringEqualValue', value: utf8dec.decode(w.payload) };
        default: return undefined;
      }
    }
    case 'RuleBytesEqual': {
      const w = RuleBytesEqual.decode(data);
      switch (w.type) {
        case RuleBytesEqual_RuleBytesEqualType.RuleBytesEqualAny: return { kind: 'BytesEqualAny' };
        case RuleBytesEqual_RuleBytesEqualType.RuleBytesEqualEmpty: return { kind: 'BytesEqualEmpty' };
        case RuleBytesEqual_RuleBytesEqualType.RuleBytesEqualValue: return { kind: 'BytesEqualValue', value: w.payload };
        default: return undefined;
      }
    }
    case 'RuleStringArrayEqual': {
      const w = RuleStringArrayEqual.decode(data);
      switch (w.type) {
        case RuleStringArrayEqual_RuleStringArrayEqualType.RuleStringArrayEqualAny: return { kind: 'StringArrayEqualAny' };
        case RuleStringArrayEqual_RuleStringArrayEqualType.RuleStringArrayEqualEmpty: return { kind: 'StringArrayEqualEmpty' };
        case RuleStringArrayEqual_RuleStringArrayEqualType.RuleStringArrayEqualValue: {
          const i = RuleStringArrayEqualValue.decode(w.payload);
          return { kind: 'StringArrayEqualValue', values: i.values };
        }
        default: return undefined;
      }
    }
    case 'RuleIntegerGreater': {
      const w = RuleIntegerGreater.decode(data);
      switch (w.type) {
        case RuleIntegerGreater_RuleIntegerGreaterType.RuleIntegerGreaterAny: return { kind: 'IntegerGreaterAny' };
        case RuleIntegerGreater_RuleIntegerGreaterType.RuleIntegerGreaterValue: return { kind: 'IntegerGreaterValue', value: magnitudeToBigint(w.payload) };
        case RuleIntegerGreater_RuleIntegerGreaterType.RuleIntegerGreaterNegValue: return { kind: 'IntegerGreaterValue', value: -magnitudeToBigint(w.payload) };
        default: return undefined;
      }
    }
    case 'RuleUIntegerGreater': {
      const w = RuleUIntegerGreater.decode(data);
      switch (w.type) {
        case RuleUIntegerGreater_RuleUIntegerGreaterType.RuleUIntegerGreaterAny: return { kind: 'UIntegerGreaterAny' };
        case RuleUIntegerGreater_RuleUIntegerGreaterType.RuleUIntegerGreaterIsZero: return { kind: 'UIntegerGreaterIsZero' };
        case RuleUIntegerGreater_RuleUIntegerGreaterType.RuleUIntegerGreaterValue: return { kind: 'UIntegerGreaterValue', value: magnitudeToBigint(w.payload) };
        case RuleUIntegerGreater_RuleUIntegerGreaterType.RuleUIntegerGreaterIsEqual: return { kind: 'UIntegerGreaterIsEqual', value: magnitudeToBigint(w.payload) };
        default: return undefined;
      }
    }
    default:
      // Unknown column type (incl. RuleAny): no cell family to type it.
      return undefined;
  }
}
