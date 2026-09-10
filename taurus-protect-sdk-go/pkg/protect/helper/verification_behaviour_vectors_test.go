package helper

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// Cross-SDK behaviour vectors for the verification primitives.
//
// These are the invariants that had drifted apart before: which (blockchain,
// network) pair selects the governance rules, and whether a hash is covered by a
// signature. Both are pure input -> outcome mappings, so all four SDKs assert them
// from one file rather than from four hand-maintained copies — the arrangement that
// let Java's hash-coverage go non-constant-time while its peers did not.
//
// To add a case: append to the shared file, bump the matching count there, and
// consume it in all four suites.
const behaviourVectorsRelPath = "../../../../scripts/resources/verification-behaviour-vectors.json"

type behaviourVectors struct {
	Counts struct {
		RuleKey      int `json:"rule_key"`
		HashCoverage int `json:"hash_coverage"`
		ContainsHash int `json:"contains_hash"`
	} `json:"counts"`
	RuleKey []struct {
		Description     string `json:"description"`
		PayloadAsString string `json:"payload_as_string"`
		DTOBlockchain   string `json:"dto_blockchain"`
		DTONetwork      string `json:"dto_network"`
		Expect          string `json:"expect"`
		Blockchain      string `json:"blockchain"`
		Network         string `json:"network"`
	} `json:"rule_key"`
	HashCoverage []struct {
		Description string     `json:"description"`
		Signatures  [][]string `json:"signatures"`
		Hash        string     `json:"hash"`
		Expect      bool       `json:"expect"`
	} `json:"hash_coverage"`
	ContainsHash []struct {
		Description string   `json:"description"`
		Hashes      []string `json:"hashes"`
		Hash        string   `json:"hash"`
		Expect      bool     `json:"expect"`
	} `json:"contains_hash"`
}

func loadBehaviourVectors(t *testing.T) behaviourVectors {
	t.Helper()

	path := filepath.Clean(behaviourVectorsRelPath)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read shared verification behaviour vectors %s: %v", path, err)
	}

	var v behaviourVectors
	if err := json.Unmarshal(raw, &v); err != nil {
		t.Fatalf("cannot parse %s: %v", path, err)
	}

	// Counts are asserted so a case added to the shared file without being consumed
	// here fails loudly rather than being silently ignored by this SDK.
	if len(v.RuleKey) != v.Counts.RuleKey {
		t.Fatalf("rule_key: got %d vectors, file declares %d", len(v.RuleKey), v.Counts.RuleKey)
	}
	if len(v.HashCoverage) != v.Counts.HashCoverage {
		t.Fatalf("hash_coverage: got %d vectors, file declares %d", len(v.HashCoverage), v.Counts.HashCoverage)
	}
	if len(v.ContainsHash) != v.Counts.ContainsHash {
		t.Fatalf("contains_hash: got %d vectors, file declares %d", len(v.ContainsHash), v.Counts.ContainsHash)
	}
	return v
}

func TestVerificationBehaviourVectors_RuleKey(t *testing.T) {
	for _, tc := range loadBehaviourVectors(t).RuleKey {
		t.Run(tc.Description, func(t *testing.T) {
			blockchain, network, err := ResolveRuleKey(tc.PayloadAsString, tc.DTOBlockchain, tc.DTONetwork)

			if tc.Expect == "error" {
				if err == nil {
					t.Fatalf("expected an error, got %q/%q", blockchain, network)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if blockchain != tc.Blockchain || network != tc.Network {
				t.Errorf("got %q/%q, want %q/%q", blockchain, network, tc.Blockchain, tc.Network)
			}
		})
	}
}

func TestVerificationBehaviourVectors_HashCoverage(t *testing.T) {
	for _, tc := range loadBehaviourVectors(t).HashCoverage {
		t.Run(tc.Description, func(t *testing.T) {
			signatures := make([]model.WhitelistSignature, 0, len(tc.Signatures))
			for _, hashes := range tc.Signatures {
				signatures = append(signatures, model.WhitelistSignature{Hashes: hashes})
			}

			if got := VerifyHashCoverage(tc.Hash, signatures); got != tc.Expect {
				t.Errorf("VerifyHashCoverage(%q) = %v, want %v", tc.Hash, got, tc.Expect)
			}
		})
	}
}

func TestVerificationBehaviourVectors_ContainsHash(t *testing.T) {
	for _, tc := range loadBehaviourVectors(t).ContainsHash {
		t.Run(tc.Description, func(t *testing.T) {
			if got := containsHash(tc.Hashes, tc.Hash); got != tc.Expect {
				t.Errorf("containsHash(%v, %q) = %v, want %v", tc.Hashes, tc.Hash, got, tc.Expect)
			}
		})
	}
}
