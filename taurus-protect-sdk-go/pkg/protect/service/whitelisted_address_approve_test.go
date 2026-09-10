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

// The completeness guard: the verified read came back without one of the requested ids,
// so nothing may be signed.
func TestApproveWhitelistedAddressesAbortsOnAnOmittedRow(t *testing.T) {
	// An empty page: no row verified, so every requested id is missing.
	svc, approved := addressApproveService(t, `{"result":[],"totalItems":"0"}`)

	err := svc.ApproveWhitelistedAddresses(context.Background(),
		[]string{"1", "2"}, addressApprovalKey(t), "ok")
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

	for _, tc := range []struct {
		name string
		call func() error
	}{
		{"no ids", func() error { return svc.ApproveWhitelistedAddresses(ctx, nil, key, "ok") }},
		{"nil key", func() error {
			return svc.ApproveWhitelistedAddresses(ctx, []string{"1"}, nil, "ok")
		}},
		{"no comment", func() error {
			return svc.ApproveWhitelistedAddresses(ctx, []string{"1"}, key, "")
		}},
		{"non-numeric id", func() error {
			return svc.ApproveWhitelistedAddresses(ctx, []string{"abc"}, key, "ok")
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
