package service

import (
	"strings"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// stubContainerVerifier verifies whatever is in `good` and rejects everything else. Under
// test is how the cache records the outcome, not the cryptography behind it.
type stubContainerVerifier struct {
	good  map[string]*model.DecodedRulesContainer
	calls int
}

func (s *stubContainerVerifier) VerifyAndDecodeRulesContainer(
	rulesContainerBase64 string,
	_ string,
	_ func(string) (*model.DecodedRulesContainer, error),
	_ func(string) ([]*model.RuleUserSignature, error),
) (*model.DecodedRulesContainer, error) {
	s.calls++
	if container, ok := s.good[rulesContainerBase64]; ok {
		return container, nil
	}
	return nil, &model.IntegrityError{
		Message: "rules container signature verification failed: insufficient distinct valid SuperAdmin signers",
	}
}

func hashContainer(hash, container, signatures string) openapi.TgvalidatordHashRulesContainer {
	return openapi.TgvalidatordHashRulesContainer{
		Hash:            &hash,
		RulesContainer:  &container,
		RulesSignatures: &signatures,
	}
}

// labelled builds an entry whose hash is the one validatord would compute for those
// bytes. Fixtures must use it: the cache now recomputes the label, so an invented one is
// rejected before verification is even attempted.
func labelled(container, signatures string) openapi.TgvalidatordHashRulesContainer {
	return hashContainer(containerHashLabel(container), container, signatures)
}

// A failed container and an absent one must not read alike: both are cache misses, but
// only one means the response carried a container that did not verify.
func TestBuildRulesContainerCache_DistinguishesFailedFromAbsent(t *testing.T) {
	good := &model.DecodedRulesContainer{}
	verifier := &stubContainerVerifier{good: map[string]*model.DecodedRulesContainer{"good-b64": good}}

	goodHash := containerHashLabel("good-b64")
	badHash := containerHashLabel("bad-b64")

	containers := buildRulesContainerCache(verifier, []openapi.TgvalidatordHashRulesContainer{
		labelled("good-b64", "sigs"),
		labelled("bad-b64", "sigs"),
	})

	if got := containers.verified[goodHash]; got != good {
		t.Errorf("a container that verified must be cached by hash, got %v", got)
	}
	if _, failed := containers.reasonFor(goodHash); failed {
		t.Error("a container that verified must not be recorded as failed")
	}

	if _, ok := containers.verified[badHash]; ok {
		t.Error("a container that failed verification must NOT be cached")
	}
	reason, failed := containers.reasonFor(badHash)
	if !failed {
		t.Fatal("a container that failed verification must be recorded, not silently skipped")
	}
	if !strings.Contains(reason, "signature verification failed") {
		t.Errorf("the recorded reason must carry the verifier's own message, got %q", reason)
	}

	if _, failed := containers.reasonFor("hash-never-sent"); failed {
		t.Error("a hash absent from the response must not read as a failed container")
	}
}

// A repeated entry verifies once. Note the label is a pure function of the bytes, so two
// DIFFERENT labels for one container is not a legitimate response shape any more — that
// case is covered by TestBuildRulesContainerCache_RejectsMislabelledContainer below.
func TestBuildRulesContainerCache_DedupesByContent(t *testing.T) {
	good := &model.DecodedRulesContainer{}
	verifier := &stubContainerVerifier{good: map[string]*model.DecodedRulesContainer{"good-b64": good}}

	containers := buildRulesContainerCache(verifier, []openapi.TgvalidatordHashRulesContainer{
		labelled("good-b64", "sigs"),
		labelled("good-b64", "sigs"),
		labelled("bad-b64", "sigs"),
		labelled("bad-b64", "sigs"),
	})

	if containers.verified[containerHashLabel("good-b64")] != good {
		t.Error("the good container must resolve from its own label")
	}
	if _, failed := containers.reasonFor(containerHashLabel("bad-b64")); !failed {
		t.Error("the failing container must be recorded against its label")
	}

	// Two distinct container payloads, so two verifications — not four.
	if verifier.calls != 2 {
		t.Errorf("expected 2 verifications for 2 distinct containers, got %d", verifier.calls)
	}
}

// The hash is the LABEL a row uses to select its container, and it arrives in the same
// response as the container. Without recomputing it, a server can file container A under
// container B's label and steer any row to any other validly-signed container — an older
// ruleset with a weaker group threshold, say. Both containers pass step 2, so SuperAdmin
// signature verification alone does not catch this.
func TestBuildRulesContainerCache_RejectsMislabelledContainer(t *testing.T) {
	good := &model.DecodedRulesContainer{}
	stale := &model.DecodedRulesContainer{}
	verifier := &stubContainerVerifier{good: map[string]*model.DecodedRulesContainer{
		"current-b64": good,
		"stale-b64":   stale, // also validly signed, just older
	}}

	// The stale container filed under the current one's label.
	currentLabel := containerHashLabel("current-b64")
	containers := buildRulesContainerCache(verifier, []openapi.TgvalidatordHashRulesContainer{
		hashContainer(currentLabel, "stale-b64", "sigs"),
	})

	if got, ok := containers.verified[currentLabel]; ok {
		t.Errorf("a mislabelled container must never be cached, got %v", got)
	}
	if got := containers.verified[currentLabel]; got == stale {
		t.Fatal("a row asking for the current container was handed the stale one")
	}
	reason, failed := containers.reasonFor(currentLabel)
	if !failed {
		t.Fatal("a mislabelled container must be recorded as failed, not silently skipped")
	}
	if !strings.Contains(reason, "hash mismatch") {
		t.Errorf("the reason must name the mismatch, got %q", reason)
	}
	if verifier.calls != 0 {
		t.Errorf("a mislabelled container is rejected before verification, got %d calls", verifier.calls)
	}
}

// The label is base64(SHA256(base64 container text)) — over the base64 TEXT, not the
// decoded protobuf, and base64 out, not hex. Hashing the decoded bytes instead would
// reject every container and break all list calls, so pin the convention itself.
func TestContainerHashLabel_MatchesValidatordConvention(t *testing.T) {
	// echo -n "abc" | sha256sum -> ba7816bf..., base64 of those raw bytes:
	const containerBase64 = "abc"
	const want = "ungWv48Bz+pBQUDeXa4iI7ADYaOWF3qctBD/YfIAFa0="

	if got := containerHashLabel(containerBase64); got != want {
		t.Errorf("containerHashLabel(%q) = %q, want %q\n"+
			"if this changed to hex, or to hashing the decoded protobuf, every container "+
			"would be rejected and every list call would fail", containerBase64, got, want)
	}
}

func TestBuildRulesContainerCache_SkipsIncompleteEntries(t *testing.T) {
	verifier := &stubContainerVerifier{good: map[string]*model.DecodedRulesContainer{}}

	containers := buildRulesContainerCache(verifier, []openapi.TgvalidatordHashRulesContainer{
		hashContainer("", "some-b64", "sigs"), // no hash to key on
		hashContainer("hash-x", "", "sigs"),   // nothing to verify
	})

	if len(containers.verified) != 0 || len(containers.failed) != 0 {
		t.Errorf("entries missing a hash or a container are not verifiable either way: %+v", containers)
	}
	if verifier.calls != 0 {
		t.Errorf("expected no verification attempts, got %d", verifier.calls)
	}
}

// The contracts endpoint carries the container on every row, so without memoization a
// page of N assets sharing one container pays N SuperAdmin verifications.
func TestInlineContainerCache_VerifiesEachContainerOnce(t *testing.T) {
	good := &model.DecodedRulesContainer{}
	verifier := &stubContainerVerifier{good: map[string]*model.DecodedRulesContainer{"good-b64": good}}
	cache := newInlineContainerCache(verifier)

	for i := 0; i < 5; i++ {
		container, err := cache.get("good-b64", "sigs")
		if err != nil {
			t.Fatalf("row %d: unexpected error: %v", i, err)
		}
		if container != good {
			t.Fatalf("row %d: expected the cached container", i)
		}
	}
	if verifier.calls != 1 {
		t.Errorf("expected 1 verification across 5 rows, got %d", verifier.calls)
	}
}

// A failing container must not be re-verified per row either, and every row must still
// see the failure.
func TestInlineContainerCache_MemoizesFailures(t *testing.T) {
	verifier := &stubContainerVerifier{good: map[string]*model.DecodedRulesContainer{}}
	cache := newInlineContainerCache(verifier)

	for i := 0; i < 3; i++ {
		if _, err := cache.get("bad-b64", "sigs"); err == nil {
			t.Fatalf("row %d: expected the container failure to surface", i)
		}
	}
	if verifier.calls != 1 {
		t.Errorf("expected 1 verification attempt across 3 rows, got %d", verifier.calls)
	}
}

// Same container bytes under different signatures is a different claim.
func TestInlineContainerCache_KeysOnSignaturesToo(t *testing.T) {
	good := &model.DecodedRulesContainer{}
	verifier := &stubContainerVerifier{good: map[string]*model.DecodedRulesContainer{"good-b64": good}}
	cache := newInlineContainerCache(verifier)

	if _, err := cache.get("good-b64", "sigs-a"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := cache.get("good-b64", "sigs-b"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if verifier.calls != 2 {
		t.Errorf("a different signature set must be verified again, got %d calls", verifier.calls)
	}
}

func TestBuildRulesContainerCache_NilVerifierYieldsEmpty(t *testing.T) {
	containers := buildRulesContainerCache(nil, []openapi.TgvalidatordHashRulesContainer{
		hashContainer("hash-good", "good-b64", "sigs"),
	})
	if len(containers.verified) != 0 || len(containers.failed) != 0 {
		t.Error("a nil verifier must yield an empty cache; the caller refuses the call separately")
	}
}
