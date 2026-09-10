package mapper

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	pb "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/proto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// vectorsRelPath points at the monorepo-shared vectors file consumed by all
// four SDK test suites (same shared-resources pattern as the swagger/proto
// sources), so cross-SDK cell wire-format parity is CI-enforced.
const vectorsRelPath = "../../../../scripts/resources/governance-cell-vectors.json"

type cellVector struct {
	Description string          `json:"description"`
	ColumnType  string          `json:"column_type"`
	CellType    string          `json:"cell_type"`
	TypedValue  json.RawMessage `json:"typed_value_json"`
	WireBase64  string          `json:"wire_base64"`
}

// TestGovernanceCellVectors validates the shared cross-SDK golden vectors: for
// every entry the Go codec must produce exactly the recorded wire bytes and
// decode them back to the same typed value.
//
// Regenerate after a deliberate grammar change with:
//
//	UPDATE_GOVERNANCE_CELL_VECTORS=1 go test ./pkg/protect/mapper -run TestGovernanceCellVectors
func TestGovernanceCellVectors(t *testing.T) {
	path := filepath.Clean(vectorsRelPath)

	if os.Getenv("UPDATE_GOVERNANCE_CELL_VECTORS") != "" {
		writeCellVectors(t, path)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read shared vectors file %s (regenerate with UPDATE_GOVERNANCE_CELL_VECTORS=1): %v", path, err)
	}

	var vectors []cellVector
	mustNoErr(t, json.Unmarshal(raw, &vectors))
	if want := len(allCellCases) + 1; len(vectors) != want {
		t.Fatalf("vectors file has %d entries, expected %d (grammar cases + raw) — regenerate the file", len(vectors), want)
	}

	byDescription := make(map[string]cellCase, len(allCellCases))
	for _, c := range allCellCases {
		byDescription[c.Description] = c
	}

	for _, v := range vectors {
		t.Run(v.Description, func(t *testing.T) {
			if v.Description == rawCellVectorDescription {
				// Unknown cell types decode to RawCell and re-encode verbatim.
				wire, err := base64.StdEncoding.DecodeString(v.WireBase64)
				mustNoErr(t, err)
				decoded := ruleCellFromBytes(v.ColumnType, wire, nil)
				wantEqual(t, model.RuleCell(model.RawCell{ColumnType: v.ColumnType, Payload: wire}), decoded)
				out, err := ruleCellToBytes(v.ColumnType, decoded)
				mustNoErr(t, err)
				wantEqual(t, wire, out)
				return
			}

			c, ok := byDescription[v.Description]
			if !ok {
				t.Fatalf("vector %q has no matching Go cell case", v.Description)
			}
			wantEqual(t, c.ColumnType, v.ColumnType)
			wantEqual(t, c.CellType, v.CellType)

			data, err := ruleCellToBytes(c.ColumnType, c.Cell)
			mustNoErr(t, err)
			wantEqual(t, v.WireBase64, base64.StdEncoding.EncodeToString(data))

			wire, err := base64.StdEncoding.DecodeString(v.WireBase64)
			mustNoErr(t, err)
			wantEqual(t, c.Cell, ruleCellFromBytes(c.ColumnType, wire, nil))
		})
	}
}

const rawCellVectorDescription = "raw cell (unknown cell type preserved verbatim)"

// rawVectorWireBytes builds a cell with a cell-type enum value no SDK knows.
func rawVectorWireBytes(t *testing.T) []byte {
	t.Helper()
	data, err := deterministicMarshal.Marshal(&pb.RuleFiatAmount{
		Type:    pb.RuleFiatAmount_RuleFiatAmountType(902),
		Payload: []byte("future-cell-payload"),
	})
	mustNoErr(t, err)
	return data
}

func writeCellVectors(t *testing.T, path string) {
	t.Helper()
	vectors := make([]cellVector, 0, len(allCellCases)+1)
	for _, c := range allCellCases {
		data, err := ruleCellToBytes(c.ColumnType, c.Cell)
		mustNoErr(t, err)
		vectors = append(vectors, cellVector{
			Description: c.Description,
			ColumnType:  c.ColumnType,
			CellType:    c.CellType,
			TypedValue:  json.RawMessage(c.TypedJSON),
			WireBase64:  base64.StdEncoding.EncodeToString(data),
		})
	}
	vectors = append(vectors,
		cellVector{
			Description: rawCellVectorDescription,
			ColumnType:  "RuleFiatAmount",
			CellType:    "",
			TypedValue:  json.RawMessage(`{"preserved_verbatim":true}`),
			WireBase64:  base64.StdEncoding.EncodeToString(rawVectorWireBytes(t)),
		},
	)
	out, err := json.MarshalIndent(vectors, "", "  ")
	mustNoErr(t, err)
	out = append(out, '\n')
	mustNoErr(t, os.WriteFile(path, out, 0o644))
	t.Logf("wrote %d vectors to %s", len(vectors), path)
}
