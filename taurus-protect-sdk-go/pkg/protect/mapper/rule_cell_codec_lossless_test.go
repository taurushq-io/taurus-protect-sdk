package mapper

import (
	"bytes"
	"encoding/base64"
	"math/big"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// Non-canonical wire encodings carry no unknown protobuf field, so an unknown-field
// check alone accepts them and silently rewrites the cell on re-encode. Each case
// below must demote to a verbatim RawCell.
func TestRuleCellFromBytes_NonCanonicalDemotesToRawCell(t *testing.T) {
	cases := []struct {
		name    string
		colType string
		data    []byte
		why     string
	}{
		{
			name:    "negative zero flips the sign arm",
			colType: "RuleIntegerGreater",
			// {type: RuleIntegerGreaterNegValue, payload: 0x00}
			data: []byte{0x08, 0x02, 0x12, 0x01, 0x00},
			why:  "re-encodes as RuleIntegerGreaterValue: the comparison operator changes",
		},
		{
			name:    "non-minimal magnitude",
			colType: "RuleIntegerGreater",
			// {type: RuleIntegerGreaterValue, payload: 0x00 0x32}
			data: []byte{0x08, 0x01, 0x12, 0x02, 0x00, 0x32},
			why:  "re-encodes minimally as 0x32",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stats := &decodeStats{}
			got := ruleCellFromBytes(tc.colType, tc.data, stats)

			raw, ok := got.(model.RawCell)
			if !ok {
				t.Fatalf("expected RawCell (%s), got %T: %#v", tc.why, got, got)
			}
			if !bytes.Equal(raw.Payload, tc.data) {
				t.Fatalf("RawCell must keep the exact wire bytes: got %x want %x", raw.Payload, tc.data)
			}
			if stats.rawCells != 1 {
				t.Fatalf("expected the cell to be reported as raw, got rawCells=%d", stats.rawCells)
			}

			// The unguarded decode is what the byte-compare protects against: it
			// produces a typed cell whose re-encoding differs from the input.
			untyped := ruleCellFromBytesTyped(tc.colType, tc.data)
			if untyped == nil {
				t.Fatalf("precondition: expected the typed decode to succeed, making the guard load-bearing")
			}
			reencoded, err := ruleCellToBytes(tc.colType, untyped)
			if err == nil && bytes.Equal(reencoded, tc.data) {
				t.Fatalf("precondition: bytes round-trip, so this case does not exercise the guard")
			}
		})
	}
}

// A canonical cell must stay typed — the guard must not demote valid input.
func TestRuleCellFromBytes_CanonicalStaysTyped(t *testing.T) {
	cell := model.IntegerGreaterValue{Value: big.NewInt(50)}
	data, err := ruleCellToBytes("RuleIntegerGreater", cell)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

	stats := &decodeStats{}
	got := ruleCellFromBytes("RuleIntegerGreater", data, stats)
	decoded, ok := got.(model.IntegerGreaterValue)
	if !ok {
		t.Fatalf("expected IntegerGreaterValue, got %T", got)
	}
	if decoded.Value.Cmp(big.NewInt(50)) != 0 {
		t.Fatalf("value = %s, want 50", decoded.Value)
	}
	if stats.rawCells != 0 {
		t.Fatalf("canonical cell must not be reported raw, got rawCells=%d", stats.rawCells)
	}
}

// A whitelisting RuleSource this SDK cannot fully represent must keep its exact wire
// bytes. The same three base64 vectors are asserted in all four SDKs.
func TestRuleSourceFromBytes_PreservesUnrepresentable(t *testing.T) {
	for _, v := range losslessVectorsByScenario(t, "rule_source_lossless") {
		t.Run(v.Description, func(t *testing.T) {
			data, err := base64.StdEncoding.DecodeString(v.WireBase64)
			if err != nil {
				t.Fatal(err)
			}
			src := ruleSourceFromBytes(data, nil)
			if len(src.Raw) == 0 {
				t.Fatalf("expected Raw preserved, got typed %+v", src)
			}
			out, err := ruleSourceToBytes(src)
			if err != nil {
				t.Fatalf("encode: %v", err)
			}
			if !bytes.Equal(out, data) {
				t.Fatalf("not re-emitted verbatim: got %x want %x", out, data)
			}
		})
	}
}

// Enum values newer than this SDK keep their numbers across a decode/encode round trip:
// collapsing them to the zero value would rewrite a column's family or widen a rule's
// sub-domain. The same base64 vector is asserted in all four SDKs.
func TestRulesContainer_UnknownEnumsPassThroughNumerically(t *testing.T) {
	vector := losslessVectorFor(t, "unknown_enum_passthrough").WireBase64

	c, err := RulesContainerFromBase64(vector)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	for _, tc := range []struct{ what, got, want string }{
		{"role", c.Users[0].Roles[0], "201"},
		{"column type", c.TransactionRules[0].Columns[0].Type, "77"},
		{"sub-domain", c.TransactionRules[0].Details.SubDomain, "202"},
		{"blockchain", c.ContractAddressWhitelistingRules[0].Blockchain, "203"},
	} {
		if tc.got != tc.want {
			t.Errorf("%s = %q, want %q", tc.what, tc.got, tc.want)
		}
	}

	out, err := RulesContainerToBase64(c)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if out != vector {
		t.Fatalf("round trip not lossless:\n got %s\nwant %s", out, vector)
	}
}

// The container carries map<string, bytes> properties at five levels. Map iteration
// order is unspecified, so without deterministic marshaling the same reviewed container
// encodes to different bytes across runs and across SDKs — which would make any
// client-side hash of the encoded container meaningless. The expected value below is
// asserted byte-for-byte in all four SDKs.
func TestRulesContainerToBase64_IsDeterministicAndCrossSDKStable(t *testing.T) {
	expected := losslessVectorFor(t, "deterministic_encoding").WireBase64

	// Keys are inserted in reverse-sorted order on purpose.
	props := map[string][]byte{}
	for _, k := range []string{"kEcho", "kDelta", "kCharlie", "kBravo", "kAlpha"} {
		props[k] = []byte(k)
	}

	for i := 0; i < 20; i++ {
		c := &model.DecodedRulesContainer{
			Users: []*model.RuleUser{{
				ID:           "u1",
				PublicKeyPEM: "PEM",
				Roles:        []string{"SUPERADMIN"},
				Properties:   props,
			}},
			Properties: props,
		}
		out, err := RulesContainerToBase64(c)
		if err != nil {
			t.Fatalf("encode: %v", err)
		}
		if out != expected {
			t.Fatalf("run %d not stable/cross-SDK-equal:\n got %s\nwant %s", i, out, expected)
		}
	}
}

// A malformed cell must degrade on its own and never abort the container: every rule for
// that tenant would go down with it, taking whitelisted-address verification with them.
// A cell payload is protobuf bytes, so a non-UTF-8 string cell and a truncated wrapper are
// both legal on the wire. The typed result is language-dependent (Go strings hold arbitrary
// bytes, so the string cell stays typed; Python/TS/Java demote it), but the container must
// decode and re-encode byte-identically in all four. Same vector in all four SDKs.
func TestRulesContainer_MalformedCellDoesNotAbortDecode(t *testing.T) {
	vector := losslessVectorFor(t, "malformed_cell_degrades_alone").WireBase64

	c, err := RulesContainerFromBase64(vector)
	if err != nil {
		t.Fatalf("a malformed cell must not fail the container decode: %v", err)
	}
	if len(c.TransactionRules) != 1 {
		t.Fatalf("transaction rules = %d, want 1", len(c.TransactionRules))
	}
	if got := len(c.TransactionRules[0].Lines[0].Cells); got != 2 {
		t.Fatalf("cells = %d, want 2", got)
	}

	out, err := RulesContainerToBase64(c)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if out != vector {
		t.Fatalf("malformed cells not preserved verbatim:\n got %s\nwant %s", out, vector)
	}
}

// Unknown protobuf fields inside the nested contract-call scoping sub-messages must
// survive a round trip: dropping them silently narrows which contract calls a rule
// covers. The 12 sibling node types already preserved theirs; these four did not.
// The same base64 vector is asserted in all four SDKs.
func TestRulesContainer_NestedRuleDetailUnknownFieldsSurvive(t *testing.T) {
	vector := losslessVectorFor(t, "nested_detail_unknown_fields").WireBase64

	c, err := RulesContainerFromBase64(vector)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	d := c.TransactionRules[0].Details
	for _, tc := range []struct {
		what string
		got  []byte
	}{
		{"evmCallContract", d.EvmCallContract.UnknownFields},
		{"xtzCallContract", d.XtzCallContract.UnknownFields},
		{"cashSettlement", d.CashSettlement.UnknownFields},
		{"cosmosDetails", d.CosmosDetails.UnknownFields},
	} {
		if len(tc.got) == 0 {
			t.Errorf("%s: unknown fields were dropped", tc.what)
		}
	}
	// The container-level report must see the nested nodes, not just the details node.
	if !c.HasUnknownFields() {
		t.Error("HasUnknownFields must report nested unknown fields")
	}

	out, err := RulesContainerToBase64(c)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if out != vector {
		t.Fatalf("nested unknown fields not re-emitted:\n got %s\nwant %s", out, vector)
	}
}

// An empty payload on a payload-carrying arm leaves the typed variant unset: parsing it
// would materialize a default sub-message, so a caller checking the variant for nil
// would see an empty wallet instead of nothing. Asserted in all four SDKs.
func TestRuleSourceFromBytes_EmptyPayloadLeavesVariantUnset(t *testing.T) {
	data, err := base64.StdEncoding.DecodeString(losslessVectorFor(t, "empty_payload_arm").WireBase64)
	if err != nil {
		t.Fatal(err)
	}
	src := ruleSourceFromBytes(data, nil)
	if src.Type != model.RuleSourceTypeInternalWallet {
		t.Fatalf("type = %v, want InternalWallet", src.Type)
	}
	if src.InternalWallet != nil {
		t.Fatalf("InternalWallet must stay nil for an empty payload, got %+v", src.InternalWallet)
	}
}

// The canonical form above (0801) and this one (08011200) decode to the SAME typed
// value: an explicitly-present zero-length payload is legal on the wire, and reading
// it leaves the variant unset either way. Re-encoding therefore emits 0801 and drops
// two bytes — from a container the SuperAdmins signed. So the non-canonical form must
// be kept verbatim as a raw source instead of being typed.
//
// An unknown-field check cannot catch this: there is no unknown field. Only the
// decode -> re-encode -> byte-compare guard sees it. Java has had this guard on this
// path; Go, Python and TypeScript did not.
func TestRuleSourceFromBytes_NonCanonicalEmptyPayloadStaysRaw(t *testing.T) {
	vec := losslessVectorFor(t, "explicit_empty_payload_noncanonical")
	data, err := base64.StdEncoding.DecodeString(vec.WireBase64)
	if err != nil {
		t.Fatal(err)
	}

	src := ruleSourceFromBytes(data, nil)
	if !bytes.Equal(src.Raw, data) {
		t.Fatalf("Raw = %x, want the input preserved verbatim %x", src.Raw, data)
	}

	// And the whole point: it must survive a re-encode unchanged.
	out, err := ruleSourceToBytes(src)
	if err != nil {
		t.Fatalf("ruleSourceToBytes: %v", err)
	}
	if !bytes.Equal(out, data) {
		t.Fatalf("re-encoded to %x, want %x — two bytes of signed container lost", out, data)
	}
}
