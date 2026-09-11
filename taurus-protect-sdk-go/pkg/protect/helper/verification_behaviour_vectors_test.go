package helper

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
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
		RuleKey            int `json:"rule_key"`
		HashCoverage       int `json:"hash_coverage"`
		ContainsHash       int `json:"contains_hash"`
		LegacyHash         int `json:"legacy_hash"`
		RuleTierCandidates int `json:"rule_tier_candidates"`
	} `json:"counts"`
	LegacyHash []struct {
		Description          string `json:"description"`
		SignedPayload        string `json:"signed_payload"`
		DeliveredPayload     string `json:"delivered_payload"`
		Expect               string `json:"expect"`
		ExpectMatchedPayload string `json:"expect_matched_payload"`
	} `json:"legacy_hash"`
	RuleTierCandidates []struct {
		Description string `json:"description"`
		Rules       []struct {
			Blockchain string `json:"blockchain"`
			Network    string `json:"network"`
		} `json:"rules"`
		Blockchain       string `json:"blockchain"`
		ExpectCandidates []struct {
			Blockchain string `json:"blockchain"`
			Network    string `json:"network"`
		} `json:"expect_candidates"`
	} `json:"rule_tier_candidates"`
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
	if len(v.LegacyHash) != v.Counts.LegacyHash {
		t.Fatalf("legacy_hash: got %d vectors, file declares %d", len(v.LegacyHash), v.Counts.LegacyHash)
	}
	if len(v.RuleTierCandidates) != v.Counts.RuleTierCandidates {
		t.Fatalf("rule_tier_candidates: got %d vectors, file declares %d",
			len(v.RuleTierCandidates), v.Counts.RuleTierCandidates)
	}
	return v
}

// TestVerificationBehaviourVectors_LegacyHash pins WHICH PAYLOAD step 6 must parse.
//
// The vector gives the payload a signer covered and the payload the server delivered. The SDK
// must decide, from the delivered payload alone, which variant a signature over the signed
// payload covers — and hand back THAT payload. Asserting the hash alone (which is all
// crypto-test-vectors.json does) cannot catch the defect, because the attack does not move any
// hash: it changes which string the parse reads.
func TestVerificationBehaviourVectors_LegacyHash(t *testing.T) {
	for _, tc := range loadBehaviourVectors(t).LegacyHash {
		t.Run(tc.Description, func(t *testing.T) {
			coveredHash := crypto.CalculateHexHash(tc.SignedPayload)

			matchedPayload := ""
			matched := false
			if crypto.CalculateHexHash(tc.DeliveredPayload) == coveredHash {
				matchedPayload, matched = tc.DeliveredPayload, true
			} else {
				for _, variant := range ComputeLegacyPayloadVariants(tc.DeliveredPayload) {
					if variant.Hash == coveredHash {
						matchedPayload, matched = variant.Payload, true
						break
					}
				}
			}

			if tc.Expect == "no_match" {
				if matched {
					t.Fatalf("expected no variant to be covered, but %q matched", matchedPayload)
				}
				return
			}
			if !matched {
				t.Fatal("expected a covered variant, found none")
			}
			if matchedPayload != tc.ExpectMatchedPayload {
				t.Errorf("matched payload:\n got  %q\n want %q", matchedPayload, tc.ExpectMatchedPayload)
			}
		})
	}
}

// TestVerificationBehaviourVectors_RuleTierCandidates pins the set of rules that must be enforced
// when the signed payload omits `network`. Too few leaves the server free to pick the quorum; too
// many rejects rows governance would accept.
func TestVerificationBehaviourVectors_RuleTierCandidates(t *testing.T) {
	for _, tc := range loadBehaviourVectors(t).RuleTierCandidates {
		t.Run(tc.Description, func(t *testing.T) {
			container := &model.DecodedRulesContainer{}
			for _, r := range tc.Rules {
				container.AddressWhitelistingRules = append(container.AddressWhitelistingRules,
					&model.AddressWhitelistingRules{Currency: r.Blockchain, Network: r.Network})
			}

			got := container.FindAddressWhitelistingRuleCandidates(tc.Blockchain)
			if len(got) != len(tc.ExpectCandidates) {
				t.Fatalf("got %d candidates, want %d", len(got), len(tc.ExpectCandidates))
			}
			for i, want := range tc.ExpectCandidates {
				if got[i].Currency != want.Blockchain || got[i].Network != want.Network {
					t.Errorf("candidate %d = %q/%q, want %q/%q",
						i, got[i].Currency, got[i].Network, want.Blockchain, want.Network)
				}
			}
		})
	}
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
