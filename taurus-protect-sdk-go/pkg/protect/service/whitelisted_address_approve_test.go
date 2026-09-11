package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// Whitelisted addresses had no approval workflow at all: both the for-approval read and
// the approve endpoint are generated in every SDK client and were wrapped by none, so an
// approver could not act on a whitelisted destination through the SDK — verified or not.
//
// The batch re-read is the part that needs guarding. It goes through the verifying list
// path filtered by ids, so a page that silently omits a row must abort rather than become
// an approval of fewer rows than the caller asked for.

func addressApproveService(t *testing.T, listBody string) (*WhitelistedAddressService, *bool) {
	t.Helper()
	approved := false

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(r.URL.Path, "/approve") {
			approved = true
			_, _ = fmt.Fprint(w, `{}`)
			return
		}
		_, _ = fmt.Fprint(w, listBody)
	}))
	t.Cleanup(srv.Close)

	cfg := openapi.NewConfiguration()
	cfg.Servers = []openapi.ServerConfiguration{{URL: srv.URL}}
	cfg.HTTPClient = srv.Client()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	svc := NewWhitelistedAddressServiceWithVerification(openapi.NewAPIClient(cfg),
		&WhitelistedAddressServiceConfig{
			SuperAdminKeys:     []*ecdsa.PublicKey{&key.PublicKey},
			MinValidSignatures: 1,
		})
	return svc, &approved
}

func addressApprovalKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return key
}

// reviewedSelection mints the content pin the way a caller does — from a verified read —
// without needing a real signed page: it builds the result the read would have produced.
//
// Going through model.WhitelistedAddressResult.Select rather than constructing the pin
// directly is deliberate: Select is the only producer, and a hand-built
// model.WhitelistedAddressApproval{} pins nothing (its field is unexported), which is what
// makes the witness forgeable-but-useless.
func reviewedSelection(t *testing.T, idToHash map[string]string) *model.WhitelistedAddressApproval {
	t.Helper()

	result := &model.WhitelistedAddressResult{}
	ids := make([]string, 0, len(idToHash))
	for id, hash := range idToHash {
		result.Addresses = append(result.Addresses, &model.WhitelistedAddress{
			ID:       id,
			Metadata: &model.WhitelistedAssetMetadata{Hash: hash},
		})
		ids = append(ids, id)
	}

	selection, err := result.Select(ids...)
	if err != nil {
		t.Fatalf("minting the reviewed selection: %v", err)
	}
	return selection
}

// The completeness guard: the verified read came back without one of the requested ids,
// so nothing may be signed.
func TestApproveWhitelistedAddressesAbortsOnAnOmittedRow(t *testing.T) {
	// An empty page: no row verified, so every requested id is missing.
	svc, approved := addressApproveService(t, `{"result":[],"totalItems":"0"}`)

	err := svc.ApproveWhitelistedAddresses(context.Background(),
		reviewedSelection(t, map[string]string{"1": "aaa", "2": "bbb"}),
		addressApprovalKey(t), "ok")
	if err == nil {
		t.Fatal("signed an approval for rows the verified read did not return")
	}
	if !strings.Contains(err.Error(), "was not returned by the verified read") {
		t.Errorf("error should name the completeness failure, got %q", err)
	}
	if *approved {
		t.Error("a refusal must never reach the wire")
	}
}

func TestApproveWhitelistedAddressesValidatesArguments(t *testing.T) {
	svc, approved := addressApproveService(t, `{"result":[],"totalItems":"0"}`)
	ctx := context.Background()
	key := addressApprovalKey(t)
	one := reviewedSelection(t, map[string]string{"1": "aaa"})

	for _, tc := range []struct {
		name string
		call func() error
	}{
		// A nil selection must be refused rather than treated as "approve nothing":
		// an empty pin would silently restore the unpinned behaviour.
		{"nil selection", func() error { return svc.ApproveWhitelistedAddresses(ctx, nil, key, "ok") }},
		// A hand-built witness pins nothing, because the field is unexported.
		{"forged empty selection", func() error {
			return svc.ApproveWhitelistedAddresses(ctx, &model.WhitelistedAddressApproval{}, key, "ok")
		}},
		{"nil key", func() error {
			return svc.ApproveWhitelistedAddresses(ctx, one, nil, "ok")
		}},
		{"no comment", func() error {
			return svc.ApproveWhitelistedAddresses(ctx, one, key, "")
		}},
		{"non-numeric id", func() error {
			return svc.ApproveWhitelistedAddresses(ctx,
				reviewedSelection(t, map[string]string{"abc": "aaa"}), key, "ok")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.call(); err == nil {
				t.Error("expected an error")
			}
		})
	}
	if *approved {
		t.Error("nothing should have reached the wire")
	}
}

// TestApproveWhitelistedAddressesRefusesWhenTheRowChangedSinceReview is the finding itself:
// the re-read returns a row under the requested id whose metadata hash is NOT the one the
// approver reviewed. Without the pin this signs the substituted hash.
func TestApproveWhitelistedAddressesRefusesASubstitutedRow(t *testing.T) {
	// A page whose row verifies as far as the completeness check is concerned, but whose
	// hash differs from the reviewed one. Verification of the row itself fails first here
	// (the fixture is not signed), which is why this test asserts on the abort rather than
	// on the pin message specifically — either refusal is correct, and both keep the
	// signature off the wire. The pin message is asserted directly in the model test below.
	svc, approved := addressApproveService(t,
		`{"result":[{"id":"1","metadata":{"hash":"substituted","payloadAsString":"{}"}}],"totalItems":"1"}`)

	err := svc.ApproveWhitelistedAddresses(context.Background(),
		reviewedSelection(t, map[string]string{"1": "the-hash-the-approver-reviewed"}),
		addressApprovalKey(t), "ok")
	if err == nil {
		t.Fatal("signed an approval over a row that changed since it was reviewed")
	}
	if *approved {
		t.Error("a refusal must never reach the wire")
	}
}

// TestWhitelistedAddressSelectRejectsUnreviewableRows pins the mint side: a caller cannot pin
// an id the verified read did not return, nor a row with no hash to pin to.
func TestWhitelistedAddressSelectRejectsUnreviewableRows(t *testing.T) {
	result := &model.WhitelistedAddressResult{
		Addresses: []*model.WhitelistedAddress{
			{ID: "1", Metadata: &model.WhitelistedAssetMetadata{Hash: "aaa"}},
			{ID: "2"}, // no metadata: nothing to pin
		},
	}

	if _, err := result.Select("3"); err == nil {
		t.Error("pinned an id this read never returned")
	}
	if _, err := result.Select("2"); err == nil {
		t.Error("pinned a row that carries no metadata hash")
	}
	if _, err := result.Select(); err == nil {
		t.Error("an empty selection must be refused, not treated as approve-nothing")
	}

	sel, err := result.Select("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hash, ok := sel.PinnedHash("1"); !ok || hash != "aaa" {
		t.Errorf("PinnedHash(1) = %q/%v, want aaa/true", hash, ok)
	}
	if sel.IsEmpty() {
		t.Error("a minted selection must not report empty")
	}

	// SelectAll pins what SURVIVED verification, which here excludes the hash-less row.
	if _, err := result.SelectAll(); err == nil {
		t.Error("SelectAll must fail when a returned row cannot be pinned")
	}
}

// The for-approval read must verify, not just return rows. An unsigned page has nothing
// that can verify, so it must not come back as a complete queue.
func TestListWhitelistedAddressesForApprovalVerifies(t *testing.T) {
	body := `{"result":[{"id":"1","metadata":{"hash":"deadbeef","payloadAsString":"{}"}}],"totalItems":"1"}`
	svc, _ := addressApproveService(t, body)

	_, err := svc.ListWhitelistedAddressesForApproval(context.Background(),
		&model.ListWhitelistedAddressesForApprovalOptions{Limit: 10})
	if err == nil {
		t.Fatal("returned an approval queue whose rows were never verified")
	}
}
