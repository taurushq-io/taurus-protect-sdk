package mapper

import (
	"bytes"
	"fmt"
	"math/big"

	pb "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/proto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
	"google.golang.org/protobuf/proto"
)

// Transaction-rule cells are wrapper messages selected by the column's type
// ({type: <cell-type enum>, payload: <bytes>}). Message-family payloads nest a
// serialized sub-message; RuleStringEqual/RuleBytesEqual carry raw scalar
// bytes; the integer families carry big-endian magnitude bytes with the sign
// expressed by the enum arm (Value vs NegValue). The *Any cell types are the
// protobuf zero values, so their serialized wrapper is the empty cell.
//
// Decode never fails: anything this SDK version cannot fully represent
// (unknown column type, unknown cell type, undecodable payload, or unknown
// protobuf fields inside the cell) becomes a model.RawCell whose bytes are
// re-emitted verbatim on encode.

// RuleCellToBytes encodes a typed governance rule cell to its protobuf wire form under
// the given column type. Peer of Java RuleCellCodec.encode, Python rule_cell_to_bytes
// and TS ruleCellToBytes.
func RuleCellToBytes(colType string, cell model.RuleCell) ([]byte, error) {
	return ruleCellToBytes(colType, cell)
}

// RuleCellFromBytes decodes one governance rule cell under the given column type.
//
// It never returns an error: anything this SDK version cannot represent losslessly comes
// back as a model.RawCell that re-encodes verbatim, so a caller can round-trip a
// schema-newer container without rewriting it. An *empty* cell under a column type this
// SDK does not know returns nil — there is nothing to preserve. Peer of Java
// RuleCellCodec.decode, Python rule_cell_from_bytes and TS ruleCellFromBytes.
func RuleCellFromBytes(colType string, data []byte) model.RuleCell {
	return ruleCellFromBytes(colType, data, nil)
}

// CellFamily returns the column type a typed cell belongs to, or a RawCell's own column
// type. Peer of Java RuleCellCodec.cellFamily, Python cell_family and TS cellFamily.
func CellFamily(cell model.RuleCell) string {
	return cellFamily(cell)
}

// ruleCellFromBytes decodes one transaction-rule cell and guarantees the decode is
// lossless: a typed value that does not re-encode to the exact input bytes is demoted
// to a verbatim RawCell. Checking for unknown protobuf *fields* alone is not enough —
// a non-canonical encoding (a non-minimal integer magnitude, a payload on a value arm,
// non-canonical field order in a nested wrapper) carries no unknown field yet still
// re-encodes differently, which would silently rewrite signed governance bytes.
func ruleCellFromBytes(colType string, data []byte, stats *decodeStats) model.RuleCell {
	raw := func() model.RuleCell {
		if stats != nil {
			stats.rawCells++
		}
		return model.RawCell{ColumnType: colType, Payload: bytes.Clone(data)}
	}

	typed := ruleCellFromBytesTyped(colType, data)
	if typed == nil {
		// Untypeable with nothing to preserve: an empty cell under a column type
		// this SDK does not know.
		return nil
	}
	if _, alreadyRaw := typed.(model.RawCell); alreadyRaw {
		return raw()
	}

	return losslessOrRaw(data, typed, func(c model.RuleCell) ([]byte, error) {
		return ruleCellToBytes(colType, c)
	}, raw)
}

// losslessOrRaw is the guard both governance codecs share: a decoded value may only
// be trusted when re-encoding it reproduces the input byte for byte.
//
//	data ──▶ decode ──▶ typed ──▶ re-encode ──▶ bytes
//	 │                                            │
//	 └──────────────── compare ───────────────────┘
//	          equal → typed        differ → raw (verbatim)
//
// An unknown-field check alone is too weak. A non-canonical encoding carries no
// unknown field yet still re-encodes differently — an explicitly-present zero-length
// payload is legal on the wire and decodes to the same typed value as an absent one,
// so re-emitting it drops bytes from a container the SuperAdmins signed.
func losslessOrRaw[T any](data []byte, typed T, encode func(T) ([]byte, error), raw func() T) T {
	reencoded, err := encode(typed)
	if err != nil || !bytes.Equal(reencoded, data) {
		return raw()
	}
	return typed
}

// ruleCellFromBytesTyped decodes one transaction-rule cell using the column type it
// is aligned with, or returns nil when it cannot be represented as a typed cell. The
// *Any cell types take the same path as every other cell type: they are the proto zero
// values, so their wire form is the empty cell, which unmarshals to the zero wrapper
// and selects the *Any arm.
func ruleCellFromBytesTyped(colType string, data []byte) model.RuleCell {
	raw := func() model.RuleCell {
		return model.RawCell{ColumnType: colType, Payload: bytes.Clone(data)}
	}

	switch colType {
	case "RuleFiatAmount":
		var w pb.RuleFiatAmount
		if proto.Unmarshal(data, &w) != nil || messageHasUnknown(&w) {
			return raw()
		}
		switch w.GetType() {
		case pb.RuleFiatAmount_RuleFiatAmountAny:
			return payloadFree(model.FiatAmountAny{}, w.GetPayload(), raw)
		case pb.RuleFiatAmount_RuleFiatAmountIsZero:
			return payloadFree(model.FiatAmountIsZero{}, w.GetPayload(), raw)
		case pb.RuleFiatAmount_RuleFiatAmountRange:
			var inner pb.RuleFiatAmountRange
			if proto.Unmarshal(w.GetPayload(), &inner) != nil || messageHasUnknown(&inner) {
				return raw()
			}
			return model.FiatAmountRange{MinAmount: inner.GetMinAmount(), MaxAmount: inner.GetMaxAmount()}
		default:
			return raw()
		}

	case "RuleSource":
		var w pb.RuleSource
		if proto.Unmarshal(data, &w) != nil || messageHasUnknown(&w) {
			return raw()
		}
		switch w.GetType() {
		case pb.RuleSource_RuleSourceAny:
			return payloadFree(model.SourceAny{}, w.GetPayload(), raw)
		case pb.RuleSource_RuleSourceAnyExchange:
			return payloadFree(model.SourceAnyExchange{}, w.GetPayload(), raw)
		case pb.RuleSource_RuleSourceInternalWallet:
			var inner pb.RuleSourceInternalWallet
			if proto.Unmarshal(w.GetPayload(), &inner) != nil || messageHasUnknown(&inner) {
				return raw()
			}
			return model.SourceInternalWallet{Path: inner.GetPath()}
		case pb.RuleSource_RuleSourceInternalAddress:
			var inner pb.RuleSourceInternalAddress
			if proto.Unmarshal(w.GetPayload(), &inner) != nil || messageHasUnknown(&inner) {
				return raw()
			}
			return model.SourceInternalAddress{Address: inner.GetAddress(), Path: inner.GetPath()}
		case pb.RuleSource_RuleSourceExchange:
			var inner pb.RuleSourceExchange
			if proto.Unmarshal(w.GetPayload(), &inner) != nil || messageHasUnknown(&inner) {
				return raw()
			}
			return model.SourceExchange{Label: inner.GetLabel()}
		case pb.RuleSource_RuleSourceExternalAddress:
			var inner pb.RuleSourceExternalAddress
			if proto.Unmarshal(w.GetPayload(), &inner) != nil || messageHasUnknown(&inner) {
				return raw()
			}
			return model.SourceExternalAddress{Address: inner.GetAddress(), Memo: inner.GetMemo()}
		default:
			return raw()
		}

	case "RuleDestination":
		var w pb.RuleDestination
		if proto.Unmarshal(data, &w) != nil || messageHasUnknown(&w) {
			return raw()
		}
		switch w.GetType() {
		case pb.RuleDestination_RuleDestinationAny:
			return payloadFree(model.DestinationAny{}, w.GetPayload(), raw)
		case pb.RuleDestination_RuleDestinationAnyExchange:
			return payloadFree(model.DestinationAnyExchange{}, w.GetPayload(), raw)
		case pb.RuleDestination_RuleDestinationAnyExternalAddress:
			return payloadFree(model.DestinationAnyExternalAddress{}, w.GetPayload(), raw)
		case pb.RuleDestination_RuleDestinationAnyContractAddress:
			return payloadFree(model.DestinationAnyContractAddress{}, w.GetPayload(), raw)
		case pb.RuleDestination_RuleDestinationInternalWallet:
			var inner pb.RuleDestinationInternalWallet
			if proto.Unmarshal(w.GetPayload(), &inner) != nil || messageHasUnknown(&inner) {
				return raw()
			}
			return model.DestinationInternalWallet{Path: inner.GetPath()}
		case pb.RuleDestination_RuleDestinationInternalAddress:
			var inner pb.RuleDestinationInternalAddress
			if proto.Unmarshal(w.GetPayload(), &inner) != nil || messageHasUnknown(&inner) {
				return raw()
			}
			return model.DestinationInternalAddress{Address: inner.GetAddress(), Path: inner.GetPath()}
		case pb.RuleDestination_RuleDestinationExternalAddress:
			var inner pb.RuleDestinationExternalAddress
			if proto.Unmarshal(w.GetPayload(), &inner) != nil || messageHasUnknown(&inner) {
				return raw()
			}
			return model.DestinationExternalAddress{Address: inner.GetAddress(), Memo: inner.GetMemo()}
		case pb.RuleDestination_RuleDestinationExchange:
			var inner pb.RuleDestinationExchange
			if proto.Unmarshal(w.GetPayload(), &inner) != nil || messageHasUnknown(&inner) {
				return raw()
			}
			return model.DestinationExchange{Label: inner.GetLabel(), Memo: inner.GetMemo()}
		case pb.RuleDestination_RuleDestinationContractAddress:
			var inner pb.RuleDestinationContractAddress
			if proto.Unmarshal(w.GetPayload(), &inner) != nil || messageHasUnknown(&inner) {
				return raw()
			}
			return model.DestinationContractAddress{
				Address:    inner.GetAddress(),
				Name:       inner.GetName(),
				Symbol:     inner.GetSymbol(),
				Blockchain: inner.GetBlockchain().String(),
			}
		default:
			return raw()
		}

	case "RuleWhitelistedContract":
		var w pb.RuleWhitelistedContract
		if proto.Unmarshal(data, &w) != nil || messageHasUnknown(&w) {
			return raw()
		}
		switch w.GetType() {
		case pb.RuleWhitelistedContract_RuleWhitelistedContractAny:
			return payloadFree(model.WhitelistedContractAny{}, w.GetPayload(), raw)
		case pb.RuleWhitelistedContract_RuleWhitelistedContract_RuleDestinationContractAddress:
			var inner pb.RuleDestinationContractAddress
			if proto.Unmarshal(w.GetPayload(), &inner) != nil || messageHasUnknown(&inner) {
				return raw()
			}
			return model.WhitelistedContractAddress{
				Address:    inner.GetAddress(),
				Name:       inner.GetName(),
				Symbol:     inner.GetSymbol(),
				Blockchain: inner.GetBlockchain().String(),
			}
		default:
			return raw()
		}

	case "RuleStringEqual":
		var w pb.RuleStringEqual
		if proto.Unmarshal(data, &w) != nil || messageHasUnknown(&w) {
			return raw()
		}
		switch w.GetType() {
		case pb.RuleStringEqual_RuleStringEqualAny:
			return payloadFree(model.StringEqualAny{}, w.GetPayload(), raw)
		case pb.RuleStringEqual_RuleStringEqualEmpty:
			return payloadFree(model.StringEqualEmpty{}, w.GetPayload(), raw)
		case pb.RuleStringEqual_RuleStringEqualValue:
			return model.StringEqualValue{Value: string(w.GetPayload())}
		default:
			return raw()
		}

	case "RuleBytesEqual":
		var w pb.RuleBytesEqual
		if proto.Unmarshal(data, &w) != nil || messageHasUnknown(&w) {
			return raw()
		}
		switch w.GetType() {
		case pb.RuleBytesEqual_RuleBytesEqualAny:
			return payloadFree(model.BytesEqualAny{}, w.GetPayload(), raw)
		case pb.RuleBytesEqual_RuleBytesEqualEmpty:
			return payloadFree(model.BytesEqualEmpty{}, w.GetPayload(), raw)
		case pb.RuleBytesEqual_RuleBytesEqualValue:
			return model.BytesEqualValue{Value: bytes.Clone(w.GetPayload())}
		default:
			return raw()
		}

	case "RuleStringArrayEqual":
		var w pb.RuleStringArrayEqual
		if proto.Unmarshal(data, &w) != nil || messageHasUnknown(&w) {
			return raw()
		}
		switch w.GetType() {
		case pb.RuleStringArrayEqual_RuleStringArrayEqualAny:
			return payloadFree(model.StringArrayEqualAny{}, w.GetPayload(), raw)
		case pb.RuleStringArrayEqual_RuleStringArrayEqualEmpty:
			return payloadFree(model.StringArrayEqualEmpty{}, w.GetPayload(), raw)
		case pb.RuleStringArrayEqual_RuleStringArrayEqualValue:
			var inner pb.RuleStringArrayEqualValue
			if proto.Unmarshal(w.GetPayload(), &inner) != nil || messageHasUnknown(&inner) {
				return raw()
			}
			return model.StringArrayEqualValue{Values: inner.GetValues()}
		default:
			return raw()
		}

	case "RuleIntegerGreater":
		var w pb.RuleIntegerGreater
		if proto.Unmarshal(data, &w) != nil || messageHasUnknown(&w) {
			return raw()
		}
		switch w.GetType() {
		case pb.RuleIntegerGreater_RuleIntegerGreaterAny:
			return payloadFree(model.IntegerGreaterAny{}, w.GetPayload(), raw)
		case pb.RuleIntegerGreater_RuleIntegerGreaterValue:
			return model.IntegerGreaterValue{Value: bigFromMagnitude(w.GetPayload())}
		case pb.RuleIntegerGreater_RuleIntegerGreaterNegValue:
			return model.IntegerGreaterValue{Value: new(big.Int).Neg(bigFromMagnitude(w.GetPayload()))}
		default:
			return raw()
		}

	case "RuleUIntegerGreater":
		var w pb.RuleUIntegerGreater
		if proto.Unmarshal(data, &w) != nil || messageHasUnknown(&w) {
			return raw()
		}
		switch w.GetType() {
		case pb.RuleUIntegerGreater_RuleUIntegerGreaterAny:
			return payloadFree(model.UIntegerGreaterAny{}, w.GetPayload(), raw)
		case pb.RuleUIntegerGreater_RuleUIntegerGreaterIsZero:
			return payloadFree(model.UIntegerGreaterIsZero{}, w.GetPayload(), raw)
		case pb.RuleUIntegerGreater_RuleUIntegerGreaterValue:
			return model.UIntegerGreaterValue{Value: bigFromMagnitude(w.GetPayload())}
		case pb.RuleUIntegerGreater_RuleUIntegerGreaterIsEqual:
			return model.UIntegerGreaterIsEqual{Value: bigFromMagnitude(w.GetPayload())}
		default:
			return raw()
		}

	case "RuleAny":
		// RuleAny columns have no cell family; their cells are empty.
		if len(data) == 0 {
			return nil
		}
		return raw()

	default:
		// Unknown column type (including numeric passthrough of enum values
		// newer than this SDK): nothing to type an empty cell with; non-empty
		// cells are preserved verbatim.
		if len(data) == 0 {
			return nil
		}
		return raw()
	}
}

// payloadFree returns the typed cell when the wire payload is empty, as
// expected for payload-less cell types; unexpected payload bytes fall back to
// RawCell so they are preserved.
func payloadFree(cell model.RuleCell, payload []byte, raw func() model.RuleCell) model.RuleCell {
	if len(payload) > 0 {
		return raw()
	}
	return cell
}

// bigFromMagnitude decodes big-endian magnitude bytes, normalizing zero to the
// canonical zero-value big.Int so decoded values compare equal to caller-built
// ones.
func bigFromMagnitude(payload []byte) *big.Int {
	v := new(big.Int).SetBytes(payload)
	if v.Sign() == 0 {
		return new(big.Int)
	}
	return v
}

// ruleCellToBytes encodes one typed transaction-rule cell into the wrapped
// protobuf bytes stored in the rule line, where colType is the type of the
// column the cell is aligned with. A nil cell encodes to empty bytes ("match
// any"); a RawCell is emitted verbatim.
func ruleCellToBytes(colType string, cell model.RuleCell) ([]byte, error) {
	if cell == nil {
		return nil, nil
	}
	if rc, ok := cell.(model.RawCell); ok {
		return rc.Payload, nil
	}

	family := cellFamily(cell)
	if colType != "" && colType != family {
		return nil, fmt.Errorf("cell %T is not valid for column type %q", cell, colType)
	}

	switch c := cell.(type) {
	// --- RuleFiatAmount ---
	case model.FiatAmountAny:
		return marshalCell(&pb.RuleFiatAmount{Type: pb.RuleFiatAmount_RuleFiatAmountAny})
	case model.FiatAmountIsZero:
		return marshalCell(&pb.RuleFiatAmount{Type: pb.RuleFiatAmount_RuleFiatAmountIsZero})
	case model.FiatAmountRange:
		inner, err := marshalCell(&pb.RuleFiatAmountRange{MinAmount: c.MinAmount, MaxAmount: c.MaxAmount})
		if err != nil {
			return nil, err
		}
		return marshalCell(&pb.RuleFiatAmount{Type: pb.RuleFiatAmount_RuleFiatAmountRange, Payload: inner})

	// --- RuleSource ---
	case model.SourceAny:
		return marshalCell(&pb.RuleSource{Type: pb.RuleSource_RuleSourceAny})
	case model.SourceAnyExchange:
		return marshalCell(&pb.RuleSource{Type: pb.RuleSource_RuleSourceAnyExchange})
	case model.SourceInternalWallet:
		inner, err := marshalCell(&pb.RuleSourceInternalWallet{Path: c.Path})
		if err != nil {
			return nil, err
		}
		return marshalCell(&pb.RuleSource{Type: pb.RuleSource_RuleSourceInternalWallet, Payload: inner})
	case model.SourceInternalAddress:
		inner, err := marshalCell(&pb.RuleSourceInternalAddress{Address: c.Address, Path: c.Path})
		if err != nil {
			return nil, err
		}
		return marshalCell(&pb.RuleSource{Type: pb.RuleSource_RuleSourceInternalAddress, Payload: inner})
	case model.SourceExchange:
		inner, err := marshalCell(&pb.RuleSourceExchange{Label: c.Label})
		if err != nil {
			return nil, err
		}
		return marshalCell(&pb.RuleSource{Type: pb.RuleSource_RuleSourceExchange, Payload: inner})
	case model.SourceExternalAddress:
		inner, err := marshalCell(&pb.RuleSourceExternalAddress{Address: c.Address, Memo: c.Memo})
		if err != nil {
			return nil, err
		}
		return marshalCell(&pb.RuleSource{Type: pb.RuleSource_RuleSourceExternalAddress, Payload: inner})

	// --- RuleDestination ---
	case model.DestinationAny:
		return marshalCell(&pb.RuleDestination{Type: pb.RuleDestination_RuleDestinationAny})
	case model.DestinationAnyExchange:
		return marshalCell(&pb.RuleDestination{Type: pb.RuleDestination_RuleDestinationAnyExchange})
	case model.DestinationAnyExternalAddress:
		return marshalCell(&pb.RuleDestination{Type: pb.RuleDestination_RuleDestinationAnyExternalAddress})
	case model.DestinationAnyContractAddress:
		return marshalCell(&pb.RuleDestination{Type: pb.RuleDestination_RuleDestinationAnyContractAddress})
	case model.DestinationInternalWallet:
		inner, err := marshalCell(&pb.RuleDestinationInternalWallet{Path: c.Path})
		if err != nil {
			return nil, err
		}
		return marshalCell(&pb.RuleDestination{Type: pb.RuleDestination_RuleDestinationInternalWallet, Payload: inner})
	case model.DestinationInternalAddress:
		inner, err := marshalCell(&pb.RuleDestinationInternalAddress{Address: c.Address, Path: c.Path})
		if err != nil {
			return nil, err
		}
		return marshalCell(&pb.RuleDestination{Type: pb.RuleDestination_RuleDestinationInternalAddress, Payload: inner})
	case model.DestinationExternalAddress:
		inner, err := marshalCell(&pb.RuleDestinationExternalAddress{Address: c.Address, Memo: c.Memo})
		if err != nil {
			return nil, err
		}
		return marshalCell(&pb.RuleDestination{Type: pb.RuleDestination_RuleDestinationExternalAddress, Payload: inner})
	case model.DestinationExchange:
		inner, err := marshalCell(&pb.RuleDestinationExchange{Label: c.Label, Memo: c.Memo})
		if err != nil {
			return nil, err
		}
		return marshalCell(&pb.RuleDestination{Type: pb.RuleDestination_RuleDestinationExchange, Payload: inner})
	case model.DestinationContractAddress:
		blockchain, err := blockchainToProto(c.Blockchain)
		if err != nil {
			return nil, err
		}
		inner, err := marshalCell(&pb.RuleDestinationContractAddress{
			Address:    c.Address,
			Name:       c.Name,
			Symbol:     c.Symbol,
			Blockchain: blockchain,
		})
		if err != nil {
			return nil, err
		}
		return marshalCell(&pb.RuleDestination{Type: pb.RuleDestination_RuleDestinationContractAddress, Payload: inner})

	// --- RuleWhitelistedContract ---
	case model.WhitelistedContractAny:
		return marshalCell(&pb.RuleWhitelistedContract{Type: pb.RuleWhitelistedContract_RuleWhitelistedContractAny})
	case model.WhitelistedContractAddress:
		blockchain, err := blockchainToProto(c.Blockchain)
		if err != nil {
			return nil, err
		}
		inner, err := marshalCell(&pb.RuleDestinationContractAddress{
			Address:    c.Address,
			Name:       c.Name,
			Symbol:     c.Symbol,
			Blockchain: blockchain,
		})
		if err != nil {
			return nil, err
		}
		return marshalCell(&pb.RuleWhitelistedContract{
			Type:    pb.RuleWhitelistedContract_RuleWhitelistedContract_RuleDestinationContractAddress,
			Payload: inner,
		})

	// --- RuleStringEqual (payload is the raw string bytes) ---
	case model.StringEqualAny:
		return marshalCell(&pb.RuleStringEqual{Type: pb.RuleStringEqual_RuleStringEqualAny})
	case model.StringEqualEmpty:
		return marshalCell(&pb.RuleStringEqual{Type: pb.RuleStringEqual_RuleStringEqualEmpty})
	case model.StringEqualValue:
		return marshalCell(&pb.RuleStringEqual{Type: pb.RuleStringEqual_RuleStringEqualValue, Payload: []byte(c.Value)})

	// --- RuleBytesEqual (payload is the raw bytes) ---
	case model.BytesEqualAny:
		return marshalCell(&pb.RuleBytesEqual{Type: pb.RuleBytesEqual_RuleBytesEqualAny})
	case model.BytesEqualEmpty:
		return marshalCell(&pb.RuleBytesEqual{Type: pb.RuleBytesEqual_RuleBytesEqualEmpty})
	case model.BytesEqualValue:
		return marshalCell(&pb.RuleBytesEqual{Type: pb.RuleBytesEqual_RuleBytesEqualValue, Payload: bytes.Clone(c.Value)})

	// --- RuleStringArrayEqual ---
	case model.StringArrayEqualAny:
		return marshalCell(&pb.RuleStringArrayEqual{Type: pb.RuleStringArrayEqual_RuleStringArrayEqualAny})
	case model.StringArrayEqualEmpty:
		return marshalCell(&pb.RuleStringArrayEqual{Type: pb.RuleStringArrayEqual_RuleStringArrayEqualEmpty})
	case model.StringArrayEqualValue:
		inner, err := marshalCell(&pb.RuleStringArrayEqualValue{Values: c.Values})
		if err != nil {
			return nil, err
		}
		return marshalCell(&pb.RuleStringArrayEqual{Type: pb.RuleStringArrayEqual_RuleStringArrayEqualValue, Payload: inner})

	// --- RuleIntegerGreater (magnitude payload, sign selects the enum arm) ---
	case model.IntegerGreaterAny:
		return marshalCell(&pb.RuleIntegerGreater{Type: pb.RuleIntegerGreater_RuleIntegerGreaterAny})
	case model.IntegerGreaterValue:
		if c.Value == nil {
			return nil, fmt.Errorf("IntegerGreaterValue requires a non-nil Value")
		}
		arm := pb.RuleIntegerGreater_RuleIntegerGreaterValue
		if c.Value.Sign() < 0 {
			arm = pb.RuleIntegerGreater_RuleIntegerGreaterNegValue
		}
		return marshalCell(&pb.RuleIntegerGreater{Type: arm, Payload: bigIntMagnitudeBytes(c.Value)})

	// --- RuleUIntegerGreater ---
	case model.UIntegerGreaterAny:
		return marshalCell(&pb.RuleUIntegerGreater{Type: pb.RuleUIntegerGreater_RuleUIntegerGreaterAny})
	case model.UIntegerGreaterIsZero:
		return marshalCell(&pb.RuleUIntegerGreater{Type: pb.RuleUIntegerGreater_RuleUIntegerGreaterIsZero})
	case model.UIntegerGreaterValue:
		payload, err := unsignedMagnitudeBytes(c.Value, "UIntegerGreaterValue")
		if err != nil {
			return nil, err
		}
		return marshalCell(&pb.RuleUIntegerGreater{Type: pb.RuleUIntegerGreater_RuleUIntegerGreaterValue, Payload: payload})
	case model.UIntegerGreaterIsEqual:
		payload, err := unsignedMagnitudeBytes(c.Value, "UIntegerGreaterIsEqual")
		if err != nil {
			return nil, err
		}
		return marshalCell(&pb.RuleUIntegerGreater{Type: pb.RuleUIntegerGreater_RuleUIntegerGreaterIsEqual, Payload: payload})

	default:
		return nil, fmt.Errorf("unsupported rule cell type %T", cell)
	}
}

// cellFamily returns the column type name a typed cell belongs to.
func cellFamily(cell model.RuleCell) string {
	switch c := cell.(type) {
	case model.FiatAmountAny, model.FiatAmountIsZero, model.FiatAmountRange:
		return "RuleFiatAmount"
	case model.SourceAny, model.SourceInternalWallet, model.SourceInternalAddress,
		model.SourceAnyExchange, model.SourceExchange, model.SourceExternalAddress:
		return "RuleSource"
	case model.DestinationAny, model.DestinationInternalWallet, model.DestinationInternalAddress,
		model.DestinationExternalAddress, model.DestinationAnyExchange, model.DestinationExchange,
		model.DestinationContractAddress, model.DestinationAnyExternalAddress, model.DestinationAnyContractAddress:
		return "RuleDestination"
	case model.StringEqualAny, model.StringEqualEmpty, model.StringEqualValue:
		return "RuleStringEqual"
	case model.BytesEqualAny, model.BytesEqualEmpty, model.BytesEqualValue:
		return "RuleBytesEqual"
	case model.StringArrayEqualAny, model.StringArrayEqualEmpty, model.StringArrayEqualValue:
		return "RuleStringArrayEqual"
	case model.IntegerGreaterAny, model.IntegerGreaterValue:
		return "RuleIntegerGreater"
	case model.UIntegerGreaterAny, model.UIntegerGreaterIsZero, model.UIntegerGreaterValue, model.UIntegerGreaterIsEqual:
		return "RuleUIntegerGreater"
	case model.WhitelistedContractAny, model.WhitelistedContractAddress:
		return "RuleWhitelistedContract"
	case model.RawCell:
		// A RawCell already carries the column it was read under; reporting "" here
		// would throw that away. Java, Python and TS all return it.
		return c.ColumnType
	default:
		return ""
	}
}

// bigIntMagnitudeBytes returns the big-endian magnitude bytes of v, using
// {0} for zero (mirroring the ecosystem's BigIntToBytes convention).
func bigIntMagnitudeBytes(v *big.Int) []byte {
	b := new(big.Int).Abs(v).Bytes()
	if len(b) == 0 {
		return []byte{0}
	}
	return b
}

func unsignedMagnitudeBytes(v *big.Int, what string) ([]byte, error) {
	if v == nil {
		return nil, fmt.Errorf("%s requires a non-nil Value", what)
	}
	if v.Sign() < 0 {
		return nil, fmt.Errorf("%s requires a non-negative Value, got %s", what, v)
	}
	return bigIntMagnitudeBytes(v), nil
}

func marshalCell(m proto.Message) ([]byte, error) {
	data, err := deterministicMarshal.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %T: %w", m, err)
	}
	return data, nil
}
