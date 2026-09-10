package mapper

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// The lossless/parity scenarios below cannot be produced by this SDK's encoder — they
// are deliberately non-canonical or schema-newer wire bytes — so they are not in
// governance-cell-vectors.json, which Go generates. They used to live as base64
// literals hand-copied into all four SDKs' round-trip suites, where nothing compared
// the copies and they could drift silently. They now come from one shared file, the
// same pattern as the cell vectors and the authorization-error vectors.
//
// To add a scenario: build the bytes with the Python protobuf runtime (it can emit
// arbitrary/malformed wire forms most easily), append an entry here, and consume it in
// all four suites.
const losslessVectorsRelPath = "../../../../scripts/resources/governance-lossless-vectors.json"

type losslessVector struct {
	Description string `json:"description"`
	Scenario    string `json:"scenario"`
	WireBase64  string `json:"wire_base64"`
}

// losslessVectorCount is asserted so a vector added to the shared file without being
// consumed here fails loudly instead of being silently ignored by this SDK.
const losslessVectorCount = 9

func loadLosslessVectors(t *testing.T) []losslessVector {
	t.Helper()

	path := filepath.Clean(losslessVectorsRelPath)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read shared lossless vectors file %s: %v", path, err)
	}

	var vectors []losslessVector
	if err := json.Unmarshal(raw, &vectors); err != nil {
		t.Fatalf("cannot parse %s: %v", path, err)
	}
	if len(vectors) != losslessVectorCount {
		t.Fatalf("shared lossless vectors file has %d entries, expected %d — update every SDK's suite in lockstep",
			len(vectors), losslessVectorCount)
	}
	return vectors
}

// losslessVectorsByScenario groups the shared vectors, preserving file order.
func losslessVectorsByScenario(t *testing.T, scenario string) []losslessVector {
	t.Helper()

	var out []losslessVector
	for _, v := range loadLosslessVectors(t) {
		if v.Scenario == scenario {
			out = append(out, v)
		}
	}
	if len(out) == 0 {
		t.Fatalf("no shared lossless vector for scenario %q", scenario)
	}
	return out
}

// losslessVector returns the single vector for a scenario.
func losslessVectorFor(t *testing.T, scenario string) losslessVector {
	t.Helper()

	vs := losslessVectorsByScenario(t, scenario)
	if len(vs) != 1 {
		t.Fatalf("scenario %q has %d vectors, expected exactly 1", scenario, len(vs))
	}
	return vs[0]
}
