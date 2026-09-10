package mapper

import (
	"bytes"
	"math/big"
	"reflect"
	"strings"
	"testing"

	pb "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/proto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
	"google.golang.org/protobuf/proto"
)

func mustNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func wantEqual(t *testing.T, want, got any) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("mismatch:\nwant: %#v\ngot:  %#v", want, got)
	}
}

func wantErrContaining(t *testing.T, err error, substr string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error containing %q, got nil", substr)
	}
	if !strings.Contains(err.Error(), substr) {
		t.Fatalf("expected error containing %q, got: %v", substr, err)
	}
}

func TestRuleCellCodec_RoundTripAllTypes(t *testing.T) {
	for _, tc := range allCellCases {
		t.Run(tc.Description, func(t *testing.T) {
			data, err := ruleCellToBytes(tc.ColumnType, tc.Cell)
			mustNoErr(t, err)

			decoded := ruleCellFromBytes(tc.ColumnType, data, nil)
			wantEqual(t, tc.Cell, decoded)
		})
	}
}

// TestRuleCellCodec_MatchAnySemantics pins the *Any behavior: they are cells
// like any other, and since they are the protobuf zero values their wire form
// is the empty cell. Round-tripping goes through the same wrapper
// unmarshal/switch as every other cell type. A nil cell is accepted at encode
// as an empty cell; decode yields nil only for columns with no cell family.
func TestRuleCellCodec_MatchAnySemantics(t *testing.T) {
	// *Any serializes to the empty cell (protobuf zero value).
	data, err := ruleCellToBytes("RuleFiatAmount", model.FiatAmountAny{})
	mustNoErr(t, err)
	if len(data) != 0 {
		t.Fatalf("expected empty bytes for FiatAmountAny, got %v", data)
	}
	// nil is accepted at encode as an empty cell.
	data, err = ruleCellToBytes("RuleFiatAmount", nil)
	mustNoErr(t, err)
	if len(data) != 0 {
		t.Fatalf("expected empty bytes for nil cell, got %v", data)
	}

	// Decode goes through the regular wrapper switch and lands on the *Any arm.
	wantEqual(t, model.RuleCell(model.FiatAmountAny{}), ruleCellFromBytes("RuleFiatAmount", nil, nil))
	wantEqual(t, model.RuleCell(model.SourceAny{}), ruleCellFromBytes("RuleSource", nil, nil))

	// Columns without any cell family cannot type an empty cell: nil.
	if got := ruleCellFromBytes("RuleAny", nil, nil); got != nil {
		t.Fatalf("expected nil cell for empty bytes under RuleAny, got %#v", got)
	}
	if got := ruleCellFromBytes("", nil, nil); got != nil {
		t.Fatalf("expected nil cell for empty bytes without a column, got %#v", got)
	}
}

func TestRuleCellCodec_RawCellPassesThroughVerbatim(t *testing.T) {
	raw := model.RawCell{ColumnType: "RuleFiatAmount", Payload: []byte{0x08, 0x63}}
	data, err := ruleCellToBytes("RuleFiatAmount", raw)
	mustNoErr(t, err)
	if !bytes.Equal(raw.Payload, data) {
		t.Fatalf("raw cell not passed through verbatim: %v != %v", raw.Payload, data)
	}
}

func TestRuleCellCodec_ColumnTypeMismatchErrors(t *testing.T) {
	_, err := ruleCellToBytes("RuleSource", model.FiatAmountRange{MinAmount: "1", MaxAmount: "2"})
	wantErrContaining(t, err, "not valid for column type")
}

func TestRuleCellCodec_IntegerRequiresValue(t *testing.T) {
	if _, err := ruleCellToBytes("RuleIntegerGreater", model.IntegerGreaterValue{}); err == nil {
		t.Fatal("expected error for nil IntegerGreaterValue.Value")
	}
	if _, err := ruleCellToBytes("RuleUIntegerGreater", model.UIntegerGreaterValue{}); err == nil {
		t.Fatal("expected error for nil UIntegerGreaterValue.Value")
	}
	_, err := ruleCellToBytes("RuleUIntegerGreater", model.UIntegerGreaterValue{Value: big.NewInt(-1)})
	wantErrContaining(t, err, "non-negative")
}

// TestRuleCellCodec_IntegerWireFormat pins the ecosystem wire convention: the
// payload is the big-endian magnitude and the sign selects the enum arm
// (RuleIntegerGreaterValue vs RuleIntegerGreaterNegValue), as parsed by the
// protect-engine rule evaluator.
func TestRuleCellCodec_IntegerWireFormat(t *testing.T) {
	data, err := ruleCellToBytes("RuleIntegerGreater", model.IntegerGreaterValue{Value: big.NewInt(-50)})
	mustNoErr(t, err)

	var w pb.RuleIntegerGreater
	mustNoErr(t, proto.Unmarshal(data, &w))
	wantEqual(t, pb.RuleIntegerGreater_RuleIntegerGreaterNegValue, w.GetType())
	wantEqual(t, []byte{50}, w.GetPayload())

	data, err = ruleCellToBytes("RuleIntegerGreater", model.IntegerGreaterValue{Value: big.NewInt(50)})
	mustNoErr(t, err)
	mustNoErr(t, proto.Unmarshal(data, &w))
	wantEqual(t, pb.RuleIntegerGreater_RuleIntegerGreaterValue, w.GetType())
	wantEqual(t, []byte{50}, w.GetPayload())

	// Zero encodes as {0}, mirroring the ecosystem BigIntToBytes convention.
	data, err = ruleCellToBytes("RuleIntegerGreater", model.IntegerGreaterValue{Value: big.NewInt(0)})
	mustNoErr(t, err)
	mustNoErr(t, proto.Unmarshal(data, &w))
	wantEqual(t, []byte{0}, w.GetPayload())
}

// TestRuleCellCodec_StringEqualWireFormat pins that RuleStringEqual carries the
// raw string bytes in the payload (not a nested message).
func TestRuleCellCodec_StringEqualWireFormat(t *testing.T) {
	data, err := ruleCellToBytes("RuleStringEqual", model.StringEqualValue{Value: "abc"})
	mustNoErr(t, err)

	var w pb.RuleStringEqual
	mustNoErr(t, proto.Unmarshal(data, &w))
	wantEqual(t, pb.RuleStringEqual_RuleStringEqualValue, w.GetType())
	wantEqual(t, []byte("abc"), w.GetPayload())
}

func TestRuleCellCodec_UnknownCellTypeDecodesToRawCell(t *testing.T) {
	data, err := proto.Marshal(&pb.RuleFiatAmount{Type: pb.RuleFiatAmount_RuleFiatAmountType(999)})
	mustNoErr(t, err)

	stats := &decodeStats{}
	decoded := ruleCellFromBytes("RuleFiatAmount", data, stats)
	raw, ok := decoded.(model.RawCell)
	if !ok {
		t.Fatalf("expected RawCell, got %T", decoded)
	}
	wantEqual(t, data, raw.Payload)
	wantEqual(t, "RuleFiatAmount", raw.ColumnType)
	wantEqual(t, 1, stats.rawCells)

	// Re-encode is verbatim.
	out, err := ruleCellToBytes("RuleFiatAmount", raw)
	mustNoErr(t, err)
	wantEqual(t, data, out)
}

func TestRuleCellCodec_UnknownColumnTypeDecodesToRawCell(t *testing.T) {
	data, err := ruleCellToBytes("RuleFiatAmount", model.FiatAmountRange{MinAmount: "1", MaxAmount: "2"})
	mustNoErr(t, err)

	decoded := ruleCellFromBytes("999", data, nil)
	raw, ok := decoded.(model.RawCell)
	if !ok {
		t.Fatalf("expected RawCell, got %T", decoded)
	}
	wantEqual(t, "999", raw.ColumnType)
	wantEqual(t, data, raw.Payload)
}

// TestRuleCellCodec_AnyEnumIsWireEmpty pins the wire truth behind the *Any
// handling: the *Any enum values are proto zero values, so proto3 serializes
// an "explicit Any" wrapper to empty bytes — identical to an empty cell — and
// decode recovers the typed *Any from the column type.
func TestRuleCellCodec_AnyEnumIsWireEmpty(t *testing.T) {
	data, err := proto.Marshal(&pb.RuleFiatAmount{Type: pb.RuleFiatAmount_RuleFiatAmountAny})
	mustNoErr(t, err)
	if len(data) != 0 {
		t.Fatalf("proto3 emitted a zero enum: %v", data)
	}
	wantEqual(t, model.RuleCell(model.FiatAmountAny{}), ruleCellFromBytes("RuleFiatAmount", data, nil))
}

func TestRuleCellCodec_UnknownFieldsInCellDecodeToRawCell(t *testing.T) {
	w := &pb.RuleFiatAmount{Type: pb.RuleFiatAmount_RuleFiatAmountAny}
	w.ProtoReflect().SetUnknown([]byte{0xC0, 0x0C, 0x2A}) // field 200, varint 42
	data, err := proto.Marshal(w)
	mustNoErr(t, err)

	decoded := ruleCellFromBytes("RuleFiatAmount", data, nil)
	raw, ok := decoded.(model.RawCell)
	if !ok {
		t.Fatalf("expected RawCell, got %T", decoded)
	}

	out, err := ruleCellToBytes("RuleFiatAmount", raw)
	mustNoErr(t, err)
	wantEqual(t, data, out)
}
