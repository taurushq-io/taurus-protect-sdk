package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// newGovernanceServiceForTest wires a GovernanceRuleService against a local
// test server.
func newGovernanceServiceForTest(t *testing.T, handler http.Handler) *GovernanceRuleService {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	cfg := openapi.NewConfiguration()
	cfg.Servers = openapi.ServerConfigurations{{URL: ts.URL}}
	// No keyless constructor exists any more: verification is mandatory on every
	// governance read. These proposal tests drive the write/sign paths, which read
	// GetRulesProposal — a proposal is not yet signed, so it is not verified — so an
	// empty key set is the right fixture and never reaches a verification call.
	return NewGovernanceRuleServiceWithVerification(openapi.NewAPIClient(cfg), nil)
}

// proposalTestContainer exercises the typed encode path end to end.
func proposalTestContainer() *model.DecodedRulesContainer {
	return &model.DecodedRulesContainer{
		MinimumDistinctUserSignatures: 1,
		Groups:                        []*model.RuleGroup{{ID: "approvers", UserIDs: []string{"u1"}}},
		TransactionRules: []*model.TransactionRules{{
			Key:     "ETH/transfer",
			Columns: []*model.RuleColumn{{Type: "RuleFiatAmount", Name: "amount", MetadataKey: "amount"}},
			Lines: []*model.RuleLine{{
				Cells: []model.RuleCell{model.FiatAmountRange{MinAmount: "1000", MaxAmount: "50000"}},
				ParallelThresholds: []*model.SequentialThresholds{{
					Thresholds: []*model.GroupThreshold{{GroupID: "approvers", MinimumSignatures: 1}},
				}},
			}},
		}},
	}
}

func TestUpdateRulesProposal_NilContainer(t *testing.T) {
	svc := newGovernanceServiceForTest(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("no HTTP call expected")
	}))
	if err := svc.UpdateRulesProposal(context.Background(), nil); err == nil {
		t.Fatal("expected error for nil container")
	}
}

func TestUpdateRulesProposal_EncodeErrorDoesNotCallServer(t *testing.T) {
	svc := newGovernanceServiceForTest(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("no HTTP call expected on encode failure")
	}))
	bad := &model.DecodedRulesContainer{
		TransactionRules: []*model.TransactionRules{{Key: "k", Columns: []*model.RuleColumn{{Type: "NotAColumnType"}}}},
	}
	if err := svc.UpdateRulesProposal(context.Background(), bad); err == nil {
		t.Fatal("expected encode error")
	}
}

func TestUpdateRulesProposal_SubmitsEncodedContainer(t *testing.T) {
	var gotContainer string
	mux := http.NewServeMux()
	mux.HandleFunc("/api/rest/v1/rules/proposal", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		var body struct {
			RulesContainer string `json:"rulesContainer"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		gotContainer = body.RulesContainer
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	})

	svc := newGovernanceServiceForTest(t, mux)
	if err := svc.UpdateRulesProposal(context.Background(), proposalTestContainer()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotContainer == "" {
		t.Fatal("server did not receive a rules container")
	}
	decoded, err := mapper.RulesContainerFromBase64(gotContainer)
	if err != nil {
		t.Fatalf("submitted container is not decodable: %v", err)
	}
	if len(decoded.TransactionRules) != 1 || decoded.TransactionRules[0].Key != "ETH/transfer" {
		t.Fatalf("submitted container lost content: %#v", decoded.TransactionRules)
	}
}

func TestUpdateRulesProposal_MapsServerError(t *testing.T) {
	svc := newGovernanceServiceForTest(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom","message":"boom","code":13}`))
	}))
	if err := svc.UpdateRulesProposal(context.Background(), proposalTestContainer()); err == nil {
		t.Fatal("expected error from server failure")
	}
}

func TestApproveRulesProposal_NilKey(t *testing.T) {
	svc := newGovernanceServiceForTest(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("no HTTP call expected")
	}))
	if err := svc.ApproveRulesProposal(context.Background(), nil, "lgtm", "deadbeef"); err == nil {
		t.Fatal("expected error for nil private key")
	}
}

func TestApproveRulesProposal_NoPendingProposal(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/rest/v1/rules/proposal", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`)) // no pending proposal
	})
	svc := newGovernanceServiceForTest(t, mux)

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	err = svc.ApproveRulesProposal(context.Background(), key, "lgtm", "deadbeef")
	if err == nil {
		t.Fatal("expected error when no proposal is pending")
	}
}

// TestApproveRulesProposal_RefusesWhenPendingContainerChanged is the content-binding gate.
// A server able to shape responses can serve the benign proposal to the review call and a
// DIFFERENT container to the re-fetch inside ApproveRulesProposal. Without the pin that
// yields a genuine SuperAdmin signature over attacker-chosen bytes — which then verifies
// clean everywhere, including in verifiers that never trusted that server. The approval
// must abort, and must not POST anything.
func TestApproveRulesProposal_RefusesWhenPendingContainerChanged(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	reviewed, err := mapper.RulesContainerToBase64(proposalTestContainer())
	if err != nil {
		t.Fatalf("failed to encode reviewed container: %v", err)
	}
	reviewedBytes, err := base64.StdEncoding.DecodeString(reviewed)
	if err != nil {
		t.Fatalf("failed to decode reviewed container: %v", err)
	}
	digest := sha256.Sum256(reviewedBytes)
	pin := hex.EncodeToString(digest[:])

	// What the server actually serves to the re-fetch: the reviewed bytes with an extra
	// field appended, the shape a merge-on-concatenation protobuf attack produces.
	substituted := base64.StdEncoding.EncodeToString(append(append([]byte{}, reviewedBytes...), 0x40, 0x01))

	approveCalled := false
	mux := http.NewServeMux()
	mux.HandleFunc("/api/rest/v1/rules/proposal", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		reply := map[string]any{"result": map[string]any{"rulesContainer": substituted}}
		_ = json.NewEncoder(w).Encode(reply)
	})
	mux.HandleFunc("/api/rest/v1/rules/proposal/approve", func(w http.ResponseWriter, r *http.Request) {
		approveCalled = true
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	})

	svc := newGovernanceServiceForTest(t, mux)
	err = svc.ApproveRulesProposal(context.Background(), key, "lgtm", pin)
	if err == nil {
		t.Fatal("signed a pending container that differs from the reviewed one")
	}
	if !strings.Contains(err.Error(), "changed since it was reviewed") {
		t.Errorf("error should name the substitution, got %q", err)
	}
	if approveCalled {
		t.Fatal("approval was submitted despite the content mismatch")
	}
}

// The pin is mandatory: an empty one would restore the unpinned behaviour silently.
func TestApproveRulesProposal_RequiresPin(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	svc := newGovernanceServiceForTest(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("no HTTP call expected: the pin must be validated before any fetch")
	}))
	if err := svc.ApproveRulesProposal(context.Background(), key, "lgtm", ""); err == nil {
		t.Fatal("expected error for an empty expectedContainerHash")
	}
}

// TestApproveRulesProposal_SignsPendingContainer pins the approval signature
// contract: SHA-256 + P-256 ECDSA (base64 raw r||s) over the decoded bytes of
// the PENDING proposal's rules container.
func TestApproveRulesProposal_SignsPendingContainer(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	pendingB64, err := mapper.RulesContainerToBase64(proposalTestContainer())
	if err != nil {
		t.Fatalf("failed to encode pending container: %v", err)
	}
	pendingBytes, err := base64.StdEncoding.DecodeString(pendingB64)
	if err != nil {
		t.Fatalf("failed to decode pending container: %v", err)
	}

	signatureVerified := false
	var gotComment string

	mux := http.NewServeMux()
	mux.HandleFunc("/api/rest/v1/rules/proposal", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		reply := map[string]any{"result": map[string]any{"rulesContainer": pendingB64}}
		_ = json.NewEncoder(w).Encode(reply)
	})
	mux.HandleFunc("/api/rest/v1/rules/proposal/approve", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		var body struct {
			Signature string `json:"signature"`
			Comment   string `json:"comment"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode approve body: %v", err)
		}
		gotComment = body.Comment

		sig, err := base64.StdEncoding.DecodeString(body.Signature)
		if err != nil {
			t.Fatalf("signature is not base64: %v", err)
		}
		if len(sig) != 64 {
			t.Fatalf("expected 64-byte r||s signature, got %d bytes", len(sig))
		}
		digest := sha256.Sum256(pendingBytes)
		rInt := new(big.Int).SetBytes(sig[:32])
		sInt := new(big.Int).SetBytes(sig[32:])
		signatureVerified = ecdsa.Verify(&key.PublicKey, digest[:], rInt, sInt)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	})

	svc := newGovernanceServiceForTest(t, mux)
	pin := hex.EncodeToString(func() []byte { d := sha256.Sum256(pendingBytes); return d[:] }())
	if err := svc.ApproveRulesProposal(context.Background(), key, "looks good", pin); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !signatureVerified {
		t.Fatal("approve signature did not verify against the pending container bytes")
	}
	if gotComment != "looks good" {
		t.Fatalf("unexpected comment: %q", gotComment)
	}
}

func TestRejectRulesProposal_SubmitsComment(t *testing.T) {
	var gotComment string
	mux := http.NewServeMux()
	mux.HandleFunc("/api/rest/v1/rules/proposal/reject", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		var body struct {
			Comment string `json:"comment"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode reject body: %v", err)
		}
		gotComment = body.Comment
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	})

	svc := newGovernanceServiceForTest(t, mux)
	if err := svc.RejectRulesProposal(context.Background(), "not compliant"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotComment != "not compliant" {
		t.Fatalf("unexpected comment: %q", gotComment)
	}
}
