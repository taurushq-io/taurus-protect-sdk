package mapper

import (
	"bytes"
	"math/big"
	"testing"

	pb "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/proto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
	"google.golang.org/protobuf/proto"
)

// This file complements rule_cell_codec_test.go with branch/edge coverage that
// the shared allCellCases table exercises only transitively: per-family wire
// shapes, every RawCell fallback branch, the scalar magnitude helpers, and
// cellFamily.

// --- Error branches ---

func TestRuleCellCodec_UIntegerIsEqualRequiresNonNegative(t *testing.T) {
	if _, err := ruleCellToBytes("RuleUIntegerGreater", model.UIntegerGreaterIsEqual{}); err == nil {
		t.Fatal("expected error for nil UIntegerGreaterIsEqual.Value")
	}
	_, err := ruleCellToBytes("RuleUIntegerGreater", model.UIntegerGreaterIsEqual{Value: big.NewInt(-1)})
	wantErrContaining(t, err, "non-negative")
}

// --- RawCell fallback branches ---

func TestRuleCellCodec_UnknownColumnEmptyBytesIsNil(t *testing.T) {
	// An unknown column type with an empty cell cannot be typed as *Any, but
	// there is nothing to preserve either: the result is nil, not a RawCell.
	if got := ruleCellFromBytes("999", nil, nil); got != nil {
		t.Fatalf("expected nil for empty cell under unknown column, got %#v", got)
	}
}

func TestRuleCellCodec_RuleAnyNonEmptyDecodesToRawCell(t *testing.T) {
	stats := &decodeStats{}
	decoded := ruleCellFromBytes("RuleAny", []byte{0x08, 0x2a}, stats)
	raw, ok := decoded.(model.RawCell)
	if !ok {
		t.Fatalf("expected RawCell, got %T", decoded)
	}
	wantEqual(t, "RuleAny", raw.ColumnType)
	wantEqual(t, 1, stats.rawCells)
}

func TestRuleCellCodec_PayloadOnPayloadFreeTypeDecodesToRawCell(t *testing.T) {
	// A payload-less cell type (RuleFiatAmountAny) that unexpectedly carries a
	// payload must be preserved verbatim, not silently dropped.
	data, err := proto.Marshal(&pb.RuleFiatAmount{
		Type:    pb.RuleFiatAmount_RuleFiatAmountAny,
		Payload: []byte("unexpected"),
	})
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

// familyWrappers pairs each column family with its wrapper's *Any (zero) value,
// so the same skew scenario can be asserted for every family, not just
// RuleFiatAmount.
var familyWrappers = []struct {
	col string
	any proto.Message
}{
	{"RuleFiatAmount", &pb.RuleFiatAmount{}},
	{"RuleSource", &pb.RuleSource{}},
	{"RuleDestination", &pb.RuleDestination{}},
	{"RuleWhitelistedContract", &pb.RuleWhitelistedContract{}},
	{"RuleStringEqual", &pb.RuleStringEqual{}},
	{"RuleBytesEqual", &pb.RuleBytesEqual{}},
	{"RuleStringArrayEqual", &pb.RuleStringArrayEqual{}},
	{"RuleIntegerGreater", &pb.RuleIntegerGreater{}},
	{"RuleUIntegerGreater", &pb.RuleUIntegerGreater{}},
}

func TestRuleCellCodec_UnknownFieldsInCellDecodeToRawCell_AllFamilies(t *testing.T) {
	for _, fw := range familyWrappers {
		t.Run(fw.col, func(t *testing.T) {
			msg := proto.Clone(fw.any)
			msg.ProtoReflect().SetUnknown([]byte{0xC0, 0x0C, 0x2A}) // field 200, varint 42
			data, err := proto.Marshal(msg)
			mustNoErr(t, err)

			decoded := ruleCellFromBytes(fw.col, data, nil)
			raw, ok := decoded.(model.RawCell)
			if !ok {
				t.Fatalf("%s: expected RawCell, got %T", fw.col, decoded)
			}
			out, err := ruleCellToBytes(fw.col, raw)
			mustNoErr(t, err)
			wantEqual(t, data, out)
		})
	}
}

func TestRuleCellCodec_UnknownCellEnumDecodesToRawCell_AllFamilies(t *testing.T) {
	// Marshal each wrapper with a cell-type enum value newer than this SDK.
	cases := []struct {
		col  string
		data []byte
	}{
		{"RuleFiatAmount", mustMarshal(t, &pb.RuleFiatAmount{Type: pb.RuleFiatAmount_RuleFiatAmountType(999)})},
		{"RuleSource", mustMarshal(t, &pb.RuleSource{Type: pb.RuleSource_RuleSourceType(999)})},
		{"RuleDestination", mustMarshal(t, &pb.RuleDestination{Type: pb.RuleDestination_RuleDestinationType(999)})},
		{"RuleWhitelistedContract", mustMarshal(t, &pb.RuleWhitelistedContract{Type: pb.RuleWhitelistedContract_RuleWhitelistedContractType(999)})},
		{"RuleStringEqual", mustMarshal(t, &pb.RuleStringEqual{Type: pb.RuleStringEqual_RuleStringEqualType(999)})},
		{"RuleBytesEqual", mustMarshal(t, &pb.RuleBytesEqual{Type: pb.RuleBytesEqual_RuleBytesEqualType(999)})},
		{"RuleStringArrayEqual", mustMarshal(t, &pb.RuleStringArrayEqual{Type: pb.RuleStringArrayEqual_RuleStringArrayEqualType(999)})},
		{"RuleIntegerGreater", mustMarshal(t, &pb.RuleIntegerGreater{Type: pb.RuleIntegerGreater_RuleIntegerGreaterType(999)})},
		{"RuleUIntegerGreater", mustMarshal(t, &pb.RuleUIntegerGreater{Type: pb.RuleUIntegerGreater_RuleUIntegerGreaterType(999)})},
	}
	for _, tc := range cases {
		t.Run(tc.col, func(t *testing.T) {
			decoded := ruleCellFromBytes(tc.col, tc.data, nil)
			raw, ok := decoded.(model.RawCell)
			if !ok {
				t.Fatalf("%s: expected RawCell, got %T", tc.col, decoded)
			}
			out, err := ruleCellToBytes(tc.col, raw)
			mustNoErr(t, err)
			wantEqual(t, tc.data, out)
		})
	}
}

func TestRuleCellCodec_UndecodableNestedPayloadDecodesToRawCell(t *testing.T) {
	// A RuleFiatAmountRange arm whose payload is not a valid RuleFiatAmountRange
	// message must be preserved verbatim (inner unmarshal guard).
	data, err := proto.Marshal(&pb.RuleFiatAmount{
		Type:    pb.RuleFiatAmount_RuleFiatAmountRange,
		Payload: []byte{0xff, 0xff, 0xff, 0xff}, // not a valid sub-message
	})
	mustNoErr(t, err)

	decoded := ruleCellFromBytes("RuleFiatAmount", data, nil)
	if _, ok := decoded.(model.RawCell); !ok {
		t.Fatalf("expected RawCell for undecodable nested payload, got %T", decoded)
	}
}

// --- Per-family wire-shape pins (complementing the Integer/String pins) ---

func TestRuleCellCodec_BytesEqualWireFormat(t *testing.T) {
	data, err := ruleCellToBytes("RuleBytesEqual", model.BytesEqualValue{Value: []byte{0xde, 0xad}})
	mustNoErr(t, err)
	var w pb.RuleBytesEqual
	mustNoErr(t, proto.Unmarshal(data, &w))
	wantEqual(t, pb.RuleBytesEqual_RuleBytesEqualValue, w.GetType())
	wantEqual(t, []byte{0xde, 0xad}, w.GetPayload()) // raw bytes, not a nested message
}

func TestRuleCellCodec_StringArrayEqualWireFormat(t *testing.T) {
	data, err := ruleCellToBytes("RuleStringArrayEqual", model.StringArrayEqualValue{Values: []string{"a", "b"}})
	mustNoErr(t, err)
	var w pb.RuleStringArrayEqual
	mustNoErr(t, proto.Unmarshal(data, &w))
	wantEqual(t, pb.RuleStringArrayEqual_RuleStringArrayEqualValue, w.GetType())
	var inner pb.RuleStringArrayEqualValue // payload is a nested message
	mustNoErr(t, proto.Unmarshal(w.GetPayload(), &inner))
	wantEqual(t, []string{"a", "b"}, inner.GetValues())
}

func TestRuleCellCodec_UIntegerWireFormat(t *testing.T) {
	data, err := ruleCellToBytes("RuleUIntegerGreater", model.UIntegerGreaterValue{Value: big.NewInt(256)})
	mustNoErr(t, err)
	var w pb.RuleUIntegerGreater
	mustNoErr(t, proto.Unmarshal(data, &w))
	wantEqual(t, pb.RuleUIntegerGreater_RuleUIntegerGreaterValue, w.GetType())
	wantEqual(t, []byte{0x01, 0x00}, w.GetPayload()) // big-endian magnitude

	// The UInteger zero path (only the signed-int zero was covered before).
	data, err = ruleCellToBytes("RuleUIntegerGreater", model.UIntegerGreaterValue{Value: new(big.Int)})
	mustNoErr(t, err)
	mustNoErr(t, proto.Unmarshal(data, &w))
	wantEqual(t, []byte{0}, w.GetPayload())
	wantEqual(t, model.RuleCell(model.UIntegerGreaterValue{Value: new(big.Int)}), ruleCellFromBytes("RuleUIntegerGreater", data, nil))
}

func TestRuleCellCodec_ContractAddressBlockchainWireFormat(t *testing.T) {
	data, err := ruleCellToBytes("RuleDestination",
		model.DestinationContractAddress{Address: "0x1", Name: "USDC", Symbol: "USDC", Blockchain: "ETH"})
	mustNoErr(t, err)
	var w pb.RuleDestination
	mustNoErr(t, proto.Unmarshal(data, &w))
	var inner pb.RuleDestinationContractAddress
	mustNoErr(t, proto.Unmarshal(w.GetPayload(), &inner))
	wantEqual(t, pb.Blockchain_ETH, inner.GetBlockchain()) // blockchain carried as a nested enum
}

// --- Blockchain enum passthrough (vectors only use "ETH") ---

func TestRuleCellCodec_BlockchainNumericPassthrough(t *testing.T) {
	// A blockchain value newer than this SDK is authored numerically; it must
	// encode to that number and decode back to the same numeric string.
	cell := model.DestinationContractAddress{Address: "0x1", Blockchain: "4242"}
	data, err := ruleCellToBytes("RuleDestination", cell)
	mustNoErr(t, err)

	decoded := ruleCellFromBytes("RuleDestination", data, nil)
	got, ok := decoded.(model.DestinationContractAddress)
	if !ok {
		t.Fatalf("expected DestinationContractAddress, got %T", decoded)
	}
	wantEqual(t, "4242", got.Blockchain)
}

func TestBlockchainToProto(t *testing.T) {
	n, err := blockchainToProto("ETH")
	mustNoErr(t, err)
	wantEqual(t, pb.Blockchain_ETH, n)

	n, err = blockchainToProto("4242") // numeric passthrough
	mustNoErr(t, err)
	wantEqual(t, pb.Blockchain(4242), n)

	_, err = blockchainToProto("NotAChain") // caller-authored typo
	wantErrContaining(t, err, "unknown blockchain")
}

// --- Magnitude scalar helpers (direct) ---

func TestBigIntMagnitudeBytes(t *testing.T) {
	wantEqual(t, []byte{0}, bigIntMagnitudeBytes(new(big.Int)))       // zero -> {0}
	wantEqual(t, []byte{0x80}, bigIntMagnitudeBytes(big.NewInt(128))) // top bit set, no sign byte
	wantEqual(t, []byte{0x01, 0x00}, bigIntMagnitudeBytes(big.NewInt(256)))
	wantEqual(t, []byte{50}, bigIntMagnitudeBytes(big.NewInt(-50))) // magnitude of a negative
}

func TestBigFromMagnitude(t *testing.T) {
	wantEqual(t, new(big.Int), bigFromMagnitude(nil))       // empty -> canonical zero
	wantEqual(t, new(big.Int), bigFromMagnitude([]byte{0})) // {0} -> canonical zero
	wantEqual(t, big.NewInt(128), bigFromMagnitude([]byte{0x80}))
	wantEqual(t, big.NewInt(256), bigFromMagnitude([]byte{0x01, 0x00}))
}

// --- cellFamily (direct) ---

func TestCellFamily(t *testing.T) {
	// A slice, not a map: some cells hold []byte/[]string and are unhashable.
	cases := []struct {
		cell model.RuleCell
		want string
	}{
		{model.FiatAmountRange{}, "RuleFiatAmount"},
		{model.SourceInternalWallet{}, "RuleSource"},
		{model.DestinationContractAddress{}, "RuleDestination"},
		{model.StringEqualValue{}, "RuleStringEqual"},
		{model.BytesEqualValue{}, "RuleBytesEqual"},
		{model.StringArrayEqualValue{}, "RuleStringArrayEqual"}, // distinct from RuleStringEqual
		{model.IntegerGreaterValue{}, "RuleIntegerGreater"},
		{model.UIntegerGreaterValue{}, "RuleUIntegerGreater"}, // distinct from RuleIntegerGreater
		{model.WhitelistedContractAny{}, "RuleWhitelistedContract"},
	}
	for _, tc := range cases {
		if got := cellFamily(tc.cell); got != tc.want {
			t.Fatalf("cellFamily(%T) = %q, want %q", tc.cell, got, tc.want)
		}
	}
	// A RawCell reports the column it was read under, matching Java, Python and TS.
	// Returning "" here (as this SDK used to) discards the one piece of grammar
	// information a raw cell still carries.
	wantEqual(t, "RuleFiatAmount", cellFamily(model.RawCell{ColumnType: "RuleFiatAmount"}))
	wantEqual(t, "", cellFamily(model.RawCell{}))
}

// sanity: encoding a RawCell whose bytes came from a genuinely empty cell.
func TestRuleCellCodec_RawCellEmptyBytes(t *testing.T) {
	out, err := ruleCellToBytes("RuleFiatAmount", model.RawCell{ColumnType: "RuleFiatAmount", Payload: nil})
	mustNoErr(t, err)
	if !bytes.Equal(out, nil) && len(out) != 0 {
		t.Fatalf("expected empty bytes, got %v", out)
	}
}
