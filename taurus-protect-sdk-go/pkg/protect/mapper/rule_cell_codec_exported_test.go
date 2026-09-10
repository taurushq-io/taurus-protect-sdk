package mapper_test

// The codec entry points must be reachable from OUTSIDE the mapper package: Java,
// Python and TS all expose encode/decode/cellFamily, while this SDK kept only the
// lowercase internals, so a Go caller could not encode or decode a model.RuleCell at
// all. This file deliberately lives in package mapper_test (not mapper) — an
// in-package test would pass just as well against the unexported functions and prove
// nothing about the public surface.

import (
	"bytes"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestExportedCodecRoundTripsFromAnotherPackage(t *testing.T) {
	cases := []struct {
		name    string
		colType string
		cell    model.RuleCell
	}{
		{"source wallet path", "RuleSource", model.SourceInternalWallet{Path: "m/44'/60'/0'"}},
		{"destination external address", "RuleDestination", model.DestinationExternalAddress{
			Address: "0xabc", Memo: "memo",
		}},
		{"fiat amount range", "RuleFiatAmount", model.FiatAmountRange{MinAmount: "1", MaxAmount: "2"}},
		{"string equal value", "RuleStringEqual", model.StringEqualValue{Value: "ETH"}},
		{"any arm encodes to the empty cell", "RuleSource", model.SourceAny{}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := mapper.RuleCellToBytes(tc.colType, tc.cell)
			if err != nil {
				t.Fatalf("RuleCellToBytes: %v", err)
			}

			decoded := mapper.RuleCellFromBytes(tc.colType, encoded)
			if decoded != tc.cell {
				t.Fatalf("round trip = %#v, want %#v", decoded, tc.cell)
			}

			reencoded, err := mapper.RuleCellToBytes(tc.colType, decoded)
			if err != nil {
				t.Fatalf("re-encode: %v", err)
			}
			if !bytes.Equal(reencoded, encoded) {
				t.Fatalf("re-encode changed the bytes: %x -> %x", encoded, reencoded)
			}

			if got := mapper.CellFamily(tc.cell); got != tc.colType {
				t.Fatalf("CellFamily(%T) = %q, want %q", tc.cell, got, tc.colType)
			}
		})
	}
}

func TestExportedCodecPreservesUnknownCellVerbatim(t *testing.T) {
	// A cell under a column type this SDK does not know must survive a decode/encode
	// round trip byte-for-byte rather than being dropped or rewritten.
	unknown := []byte{0x08, 0x7f, 0x12, 0x03, 'a', 'b', 'c'}

	decoded := mapper.RuleCellFromBytes("RuleSomethingNewer", unknown)
	rawCell, ok := decoded.(model.RawCell)
	if !ok {
		t.Fatalf("decoded %T, want model.RawCell", decoded)
	}
	if got := mapper.CellFamily(rawCell); got != "RuleSomethingNewer" {
		t.Fatalf("CellFamily(RawCell) = %q, want the column it was read under", got)
	}

	reencoded, err := mapper.RuleCellToBytes("RuleSomethingNewer", decoded)
	if err != nil {
		t.Fatalf("re-encode: %v", err)
	}
	if !bytes.Equal(reencoded, unknown) {
		t.Fatalf("unknown cell was rewritten: %x -> %x", unknown, reencoded)
	}
}

func TestExportedCodecEmptyUnknownCellIsNil(t *testing.T) {
	// Deliberate Go-specific behaviour: an EMPTY cell under an unknown column has
	// nothing to preserve, so it decodes to nil rather than an empty RawCell. Keep it
	// — the container decoder relies on it to mean "match any".
	if got := mapper.RuleCellFromBytes("RuleSomethingNewer", nil); got != nil {
		t.Fatalf("RuleCellFromBytes(unknown column, empty) = %#v, want nil", got)
	}
}
