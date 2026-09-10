package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// The governance container carries the HSM public key that address verification
// trusts, so an unverified one is the worst object this SDK can hand back. GetRules,
// GetRulesByID and GetRulesHistory returned it unverified, and
// GetDecodedRulesContainer skipped verification entirely when no keys were configured
// — silently, with no error. Nothing tested any of that.

// wireValidUnsignedContainer is a container that DECODES cleanly but carries no
// signatures. It has to decode: a malformed blob throws from the decoder instead, and
// the test then passes whether verification ran or not. This trap has bitten this repo
// twice.
func wireValidUnsignedContainer(t *testing.T) string {
	t.Helper()
	encoded, err := mapper.RulesContainerToBase64(&model.DecodedRulesContainer{
		MinimumDistinctUserSignatures: 1,
	})
	if err != nil {
		t.Fatalf("building the fixture container: %v", err)
	}
	return encoded
}

// The guard that keeps the tests below non-vacuous.
func TestWireValidContainerDecodesCleanly(t *testing.T) {
	if _, err := mapper.RulesContainerFromBase64(wireValidUnsignedContainer(t)); err != nil {
		t.Fatalf("fixture must decode, or the verification tests pass for the wrong reason: %v", err)
	}
}

// governanceService serves one canned rules reply on every governance path and counts
// the requests, so a test can tell a memo hit from a re-fetch.
func governanceService(t *testing.T, body string, keys []*ecdsa.PublicKey) *GovernanceRuleService {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, body)
	}))
	t.Cleanup(srv.Close)

	cfg := openapi.NewConfiguration()
	cfg.Servers = []openapi.ServerConfiguration{{URL: srv.URL}}
	cfg.HTTPClient = srv.Client()
	return NewGovernanceRuleServiceWithVerification(openapi.NewAPIClient(cfg),
		&GovernanceRuleServiceConfig{SuperAdminKeys: keys, MinValidSignatures: 1})
}

func testSuperAdminKey(t *testing.T) *ecdsa.PublicKey {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return &key.PublicKey
}

func rulesReply(t *testing.T, container string) string {
	t.Helper()
	b, err := json.Marshal(map[string]any{
		"result": map[string]any{"id": "1", "rulesContainer": container, "rulesSignatures": []any{}},
	})
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func historyReply(t *testing.T, container string) string {
	t.Helper()
	b, err := json.Marshal(map[string]any{
		"result": []map[string]any{
			{"id": "1", "rulesContainer": container, "rulesSignatures": []any{}},
			{"id": "2", "rulesContainer": container, "rulesSignatures": []any{}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// Every governance read path must refuse an unsigned container.
func TestGovernanceReadsRefuseUnsignedContainer(t *testing.T) {
	container := wireValidUnsignedContainer(t)
	keys := []*ecdsa.PublicKey{testSuperAdminKey(t)}

	for _, tc := range []struct {
		name string
		body string
		call func(context.Context, *GovernanceRuleService) error
	}{
		{"GetRules", rulesReply(t, container), func(ctx context.Context, s *GovernanceRuleService) error {
			_, err := s.GetRules(ctx)
			return err
		}},
		{"GetRulesByID", rulesReply(t, container), func(ctx context.Context, s *GovernanceRuleService) error {
			_, err := s.GetRulesByID(ctx, "1")
			return err
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			svc := governanceService(t, tc.body, keys)
			if err := tc.call(context.Background(), svc); err == nil {
				t.Fatal("returned an unsigned governance container with no error")
			}
		})
	}
}

// The fail-open case: a service with no SuperAdmin keys used to decode and return the
// container anyway. Verification is mandatory now, so it must error instead.
func TestGetDecodedRulesContainerFailsClosedWithoutKeys(t *testing.T) {
	container := wireValidUnsignedContainer(t)
	svc := governanceService(t, rulesReply(t, container), nil)

	_, err := svc.GetDecodedRulesContainer(&model.GovernanceRuleset{RulesContainer: container})
	if err == nil {
		t.Fatal("decoded and returned an unverified container because no keys were configured")
	}
	if !strings.Contains(err.Error(), "signature") && !strings.Contains(err.Error(), "key") {
		t.Errorf("error should name the verification failure, got %q", err)
	}
}

// memoKey returns the memo key, failing the test if the ruleset is not memoisable.
func memoKey(t *testing.T, rules *model.GovernanceRuleset) string {
	t.Helper()
	key, ok := rulesetVerificationKey(rules)
	if !ok {
		t.Fatalf("container should be memoisable, got ok=false")
	}
	return key
}

// The memo must key on the signature set, not just the container bytes: otherwise
// re-signing the same container would be served from a stale hit.
func TestVerificationMemoKeysOnTheSignatureSet(t *testing.T) {
	container := wireValidUnsignedContainer(t)
	base := &model.GovernanceRuleset{RulesContainer: container}
	resigned := &model.GovernanceRuleset{
		RulesContainer: container,
		Signatures:     []model.RuleUserSignature{{UserID: "u1", Signature: "sig"}},
	}

	if memoKey(t, base) == memoKey(t, resigned) {
		t.Fatal("a container with a different signature set must not share a memo key")
	}

	// Same content, signatures listed in the other order: one key, so the memo hits.
	a := &model.GovernanceRuleset{RulesContainer: container, Signatures: []model.RuleUserSignature{
		{UserID: "u1", Signature: "s1"}, {UserID: "u2", Signature: "s2"},
	}}
	b := &model.GovernanceRuleset{RulesContainer: container, Signatures: []model.RuleUserSignature{
		{UserID: "u2", Signature: "s2"}, {UserID: "u1", Signature: "s1"},
	}}
	if memoKey(t, a) != memoKey(t, b) {
		t.Error("signature order is server-controlled, so it must not change the memo key")
	}
}

// userId is not read by verification, so it must not participate in the memo key —
// otherwise a field the server controls freely and verification ignores can influence a
// verification-SKIP decision.
func TestVerificationMemoIgnoresUserID(t *testing.T) {
	container := wireValidUnsignedContainer(t)
	a := &model.GovernanceRuleset{RulesContainer: container, Signatures: []model.RuleUserSignature{
		{UserID: "u1", Signature: "s1"},
	}}
	b := &model.GovernanceRuleset{RulesContainer: container, Signatures: []model.RuleUserSignature{
		{UserID: "attacker-chosen", Signature: "s1"},
	}}
	if memoKey(t, a) != memoKey(t, b) {
		t.Error("userId is ignored by verification, so it must not change the memo key")
	}
}

// The memo key must be INJECTIVE. A plain concatenation leaves the boundary between the
// container and the signature list uncommitted, so a response-controlling attacker can
// shift bytes across it and make a MODIFIED container collide with a genuine one's key —
// inheriting its "already verified" status and skipping ECDSA entirely.
//
// Each pair below collides under an unprefixed key and must not collide here. Reverting
// rulesetVerificationKey to a concatenation takes this test red.
func TestVerificationMemoKeyIsInjective(t *testing.T) {
	container := wireValidUnsignedContainer(t)
	raw, err := base64.StdEncoding.DecodeString(container)
	if err != nil {
		t.Fatalf("fixture container must decode: %v", err)
	}

	signed := &model.GovernanceRuleset{RulesContainer: container, Signatures: []model.RuleUserSignature{
		{Signature: "s1"},
	}}

	// Fold the signature list into the container: genuine (C, [sig]) vs (C‖tail, []).
	// Three tail shapes, one per plausible separator-free encoding, so this test is red
	// against ALL of them rather than only the one that happens to be in place:
	//	"s1"         → plain concatenation
	//	"\x00s1"     → a 0x00 written before each entry
	//	"\x00\x00s1" → 0x00 separator plus an empty userId field
	for _, tail := range []string{"s1", "\x00s1", "\x00\x00s1"} {
		t.Run(fmt.Sprintf("signature folded into container tail %q", tail), func(t *testing.T) {
			folded := base64.StdEncoding.EncodeToString(append(append([]byte{}, raw...), tail...))
			if memoKey(t, signed) == memoKey(t, &model.GovernanceRuleset{RulesContainer: folded}) {
				t.Fatal("distinct (container, signatures) pairs must not share a memo key")
			}
		})
	}

	t.Run("one signature split into two", func(t *testing.T) {
		a := &model.GovernanceRuleset{RulesContainer: container, Signatures: []model.RuleUserSignature{
			{Signature: "abcd"},
		}}
		b := &model.GovernanceRuleset{RulesContainer: container, Signatures: []model.RuleUserSignature{
			{Signature: "ab"}, {Signature: "cd"},
		}}
		if memoKey(t, a) == memoKey(t, b) {
			t.Fatal("distinct (container, signatures) pairs must not share a memo key")
		}
	})
}

// An undecodable container has no stable identity, so it must not be memoised — and
// verification must still get the final say.
func TestVerificationMemoSkipsUndecodableContainer(t *testing.T) {
	if _, ok := rulesetVerificationKey(&model.GovernanceRuleset{RulesContainer: "!!!not base64!!!"}); ok {
		t.Error("an undecodable container must not produce a memo key")
	}
}

// A failure must resurface on every call rather than being remembered as a verdict.
func TestVerificationMemoDoesNotCacheFailures(t *testing.T) {
	container := wireValidUnsignedContainer(t)
	svc := governanceService(t, rulesReply(t, container), []*ecdsa.PublicKey{testSuperAdminKey(t)})
	ctx := context.Background()

	first := func() error { _, err := svc.GetRules(ctx); return err }
	if err := first(); err == nil {
		t.Fatal("expected the first read to fail")
	}
	if err := first(); err == nil {
		t.Fatal("a memoised FAILURE would make the second read succeed; failures must not be cached")
	}
}

// History is LENIENT where the single-ruleset reads are strict, and that is deliberate:
// a SuperAdmin key rotation makes every pre-rotation ruleset unverifiable, so a strict
// page would deny access to the whole audit trail from the rotation onwards. It must
// still NAME what it dropped, or a shortened page reads as a complete one.
func TestGetRulesHistoryExcludesAndNamesUnverifiedEntries(t *testing.T) {
	container := wireValidUnsignedContainer(t)
	svc := governanceService(t, historyReply(t, container), []*ecdsa.PublicKey{testSuperAdminKey(t)})

	result, err := svc.GetRulesHistory(context.Background(), nil)
	if err != nil {
		t.Fatalf("history must not fail outright on unverifiable entries: %v", err)
	}
	if len(result.Rules) != 0 {
		t.Errorf("no entry verified, so none may be returned; got %d", len(result.Rules))
	}
	if len(result.ExcludedUnverified) != 2 {
		t.Fatalf("both entries must be named as excluded; got %d", len(result.ExcludedUnverified))
	}
	for _, ex := range result.ExcludedUnverified {
		if ex.Reason == "" {
			t.Error("an exclusion must carry a reason")
		}
	}
	// The server counted 2; the caller can read 0.
	if result.TotalItems != 0 {
		t.Errorf("TotalItems must be reduced by the exclusions, got %d", result.TotalItems)
	}
}
