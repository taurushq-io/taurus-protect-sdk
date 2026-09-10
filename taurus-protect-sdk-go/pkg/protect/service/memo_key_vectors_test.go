package service

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// Cross-SDK vectors for the governance verification memo key.
//
// The key decides whether ECDSA runs at all, so it MUST be injective. An unprefixed
// concatenation leaves the boundary between the container and the signature list
// uncommitted, and a response-controlling attacker can then shift bytes across it to make
// a MODIFIED container inherit a genuine one's "already verified" status — skipping
// signature verification on the document that carries every HSM key the SDK trusts.
//
// The pairs live in the shared file so all four SDKs assert the same property from one
// place. This loader is in package service rather than beside the other behaviour vectors
// because rulesetVerificationKey is unexported here.
//
// To add a case: append to the shared file, bump counts.memo_key there, and consume it in
// all four suites.
const memoKeyVectorsRelPath = "../../../../scripts/resources/verification-behaviour-vectors.json"

type memoKeyVectorFile struct {
	Counts struct {
		MemoKey int `json:"memo_key"`
	} `json:"counts"`
	MemoKey []memoKeyVector `json:"memo_key"`
}

type memoKeyVector struct {
	Description string         `json:"description"`
	A           memoKeyRuleset `json:"a"`
	B           memoKeyRuleset `json:"b"`
	Expect      string         `json:"expect"`
}

type memoKeyRuleset struct {
	RulesContainer string `json:"rules_container"`
	Signatures     []struct {
		UserID    string `json:"user_id"`
		Signature string `json:"signature"`
	} `json:"signatures"`
}

func (r memoKeyRuleset) toModel() *model.GovernanceRuleset {
	sigs := make([]model.RuleUserSignature, 0, len(r.Signatures))
	for _, s := range r.Signatures {
		sigs = append(sigs, model.RuleUserSignature{UserID: s.UserID, Signature: s.Signature})
	}
	return &model.GovernanceRuleset{RulesContainer: r.RulesContainer, Signatures: sigs}
}

func TestMemoKeyVectors(t *testing.T) {
	path := filepath.Clean(memoKeyVectorsRelPath)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read shared verification behaviour vectors %s: %v", path, err)
	}

	var file memoKeyVectorFile
	if err := json.Unmarshal(raw, &file); err != nil {
		t.Fatalf("cannot parse %s: %v", path, err)
	}

	// Asserted so a case added to the shared file without being consumed here fails
	// loudly rather than being silently ignored by this SDK.
	if len(file.MemoKey) != file.Counts.MemoKey {
		t.Fatalf("memo_key: got %d vectors, file declares %d", len(file.MemoKey), file.Counts.MemoKey)
	}
	if len(file.MemoKey) == 0 {
		t.Fatal("memo_key: no vectors loaded")
	}

	// Non-vacuity: without at least one "distinct" pair the whole section would pass
	// against a key function that returns a constant.
	distinct := 0
	for _, v := range file.MemoKey {
		if v.Expect == "distinct" {
			distinct++
		}
	}
	if distinct == 0 {
		t.Fatal("memo_key: no 'distinct' vector, so the section cannot catch a colliding key")
	}

	for _, v := range file.MemoKey {
		t.Run(v.Description, func(t *testing.T) {
			keyA, okA := rulesetVerificationKey(v.A.toModel())
			keyB, okB := rulesetVerificationKey(v.B.toModel())
			if !okA || !okB {
				t.Fatalf("both containers must be decodable: okA=%v okB=%v", okA, okB)
			}

			switch v.Expect {
			case "distinct":
				if keyA == keyB {
					t.Error("distinct (container, signatures) pairs must not share a memo key")
				}
			case "same":
				if keyA != keyB {
					t.Error("these inputs describe the same verified document, so the key must match")
				}
			default:
				t.Fatalf("unknown expect %q", v.Expect)
			}
		})
	}
}
