package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"net/http"
	"testing"
)

// This file complements governance_rule_proposal_test.go with the approve/reject
// error paths (fetch failure, undecodable pending container, reject server error).

func TestApproveRulesProposal_FetchError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/rest/v1/rules/proposal", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom","message":"boom","code":13}`))
	})
	svc := newGovernanceServiceForTest(t, mux)

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	if err := svc.ApproveRulesProposal(context.Background(), key, "lgtm", "deadbeef"); err == nil {
		t.Fatal("expected error when the proposal fetch fails")
	}
}

func TestApproveRulesProposal_UndecodablePendingContainer(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/rest/v1/rules/proposal", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// A pending proposal whose rulesContainer is not valid base64.
		reply := map[string]any{"result": map[string]any{"rulesContainer": "!!!not-base64!!!"}}
		_ = json.NewEncoder(w).Encode(reply)
	})
	mux.HandleFunc("/api/rest/v1/rules/proposal/approve", func(http.ResponseWriter, *http.Request) {
		t.Fatal("approve must not be called when the pending container cannot be decoded")
	})
	svc := newGovernanceServiceForTest(t, mux)

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	if err := svc.ApproveRulesProposal(context.Background(), key, "lgtm", "deadbeef"); err == nil {
		t.Fatal("expected error for an undecodable pending container")
	}
}

func TestRejectRulesProposal_MapsServerError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/rest/v1/rules/proposal/reject", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom","message":"boom","code":13}`))
	})
	svc := newGovernanceServiceForTest(t, mux)
	if err := svc.RejectRulesProposal(context.Background(), "nope"); err == nil {
		t.Fatal("expected error from server failure")
	}
}
