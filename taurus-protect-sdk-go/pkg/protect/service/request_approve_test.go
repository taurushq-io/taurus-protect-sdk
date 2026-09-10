package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// These tests cover the create→approve flow, which had NO test of any kind: the
// package carried Create*, GetRequest_* and VerifyRequestHash_* tests and nothing that
// called ApproveRequests. That gap is why "the approve guard requires verification, the
// create paths never perform it" was invisible.

const approvePath = "/api/rest/v1/requests/approve"

// Local rather than shared: the approval-key helpers in the whitelist test files are
// being reworked alongside this change, and a test fixture should not break because a
// sibling file moved.
func testApprovalKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return key
}

// requestRow builds the DTO shape the server returns for a single request.
func requestRow(id, payload, hash string) map[string]any {
	return map[string]any{
		"id":     id,
		"status": "APPROVING",
		"metadata": map[string]any{
			"hash":            hash,
			"payloadAsString": payload,
		},
	}
}

// createAndApproveService stubs both endpoints a create→approve flow touches, and
// records whether the approval was actually submitted — a refusal must never reach the
// wire.
func createAndApproveService(t *testing.T, row map[string]any) (*RequestService, *bool) {
	t.Helper()
	approved := false

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, approvePath) {
			approved = true
			_, _ = fmt.Fprint(w, `{"signedRequests":"1"}`)
			return
		}
		b, err := json.Marshal(map[string]any{"result": row})
		if err != nil {
			t.Error(err)
			return
		}
		_, _ = fmt.Fprint(w, string(b))
	}))
	t.Cleanup(srv.Close)

	cfg := openapi.NewConfiguration()
	cfg.Servers = []openapi.ServerConfiguration{{URL: srv.URL}}
	cfg.HTTPClient = srv.Client()
	return NewRequestService(openapi.NewAPIClient(cfg), WithServiceLogger(&recordingLogger{})), &approved
}

// The flow the Java service's own javadoc walks through, and the one that could not
// succeed: every create path returned an unverified request while ApproveRequests
// refused anything unverified.
func TestCreatePathsProduceAnApprovableRequest(t *testing.T) {
	payload := `[{"key":"currency","value":"BTC"}]`
	row := requestRow("7", payload, crypto.CalculateHexHash(payload))

	for _, tc := range []struct {
		name string
		call func(context.Context, *RequestService) (*model.Request, error)
	}{
		{"CreateOutgoingRequest", func(ctx context.Context, s *RequestService) (*model.Request, error) {
			return s.CreateOutgoingRequest(ctx, &model.CreateOutgoingRequest{
				Amount: "1", FromAddressID: "1", ToAddressID: "2",
			})
		}},
		{"CreateInternalTransferRequest", func(ctx context.Context, s *RequestService) (*model.Request, error) {
			return s.CreateInternalTransferRequest(ctx, "1", "2", "1")
		}},
		{"CreateInternalTransferFromWalletRequest", func(ctx context.Context, s *RequestService) (*model.Request, error) {
			return s.CreateInternalTransferFromWalletRequest(ctx, "1", "2", "1")
		}},
		{"CreateExternalTransferRequest", func(ctx context.Context, s *RequestService) (*model.Request, error) {
			return s.CreateExternalTransferRequest(ctx, "1", "2", "1")
		}},
		{"CreateExternalTransferFromWalletRequest", func(ctx context.Context, s *RequestService) (*model.Request, error) {
			return s.CreateExternalTransferFromWalletRequest(ctx, "1", "2", "1")
		}},
		{"CreateCancelRequest", func(ctx context.Context, s *RequestService) (*model.Request, error) {
			return s.CreateCancelRequest(ctx, "1", "")
		}},
		{"CreateIncomingRequest", func(ctx context.Context, s *RequestService) (*model.Request, error) {
			return s.CreateIncomingRequest(ctx, &model.CreateIncomingRequest{
				FromExchangeID: "1", ToAddressID: "2", Amount: "1",
			})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			svc, approved := createAndApproveService(t, row)
			ctx := context.Background()

			r, err := tc.call(ctx, svc)
			if err != nil {
				t.Fatalf("create: %v", err)
			}
			// The create response carried a payload, so the seam must have cleared it.
			if !r.Metadata.HashVerified {
				t.Error("created request is not marked verified — the create path skipped the seam")
			}
			// And the accessors must work, which they cannot on unmaterialised metadata.
			if got, err := r.Metadata.GetMetadataCurrency(); err != nil || got != "BTC" {
				t.Errorf("GetMetadataCurrency() = %q, %v; want \"BTC\", nil", got, err)
			}

			if _, err := svc.ApproveRequest(ctx, r, testApprovalKey(t), "ok"); err != nil {
				t.Fatalf("approve after create: %v", err)
			}
			if !*approved {
				t.Error("approval never reached the wire")
			}
		})
	}
}

// A create whose response metadata does not verify must fail loudly and name the id:
// the server-side write already happened, so the caller needs the id to reconcile
// rather than retrying and creating a second request.
func TestCreateRejectsUnverifiableMetadataAndNamesTheID(t *testing.T) {
	good := `[{"key":"currency","value":"BTC"}]`
	tampered := `[{"key":"currency","value":"ETH"}]`
	row := requestRow("42", tampered, crypto.CalculateHexHash(good))

	svc, approved := createAndApproveService(t, row)
	r, err := svc.CreateExternalTransferRequest(context.Background(), "1", "2", "1")
	if err == nil {
		t.Fatal("create returned a request whose metadata does not verify")
	}
	if r != nil {
		t.Error("create returned a non-nil request alongside the error")
	}
	var intErr *model.IntegrityError
	if !errors.As(err, &intErr) {
		t.Fatalf("error should unwrap to IntegrityError, got %T", err)
	}
	if !strings.Contains(err.Error(), "42") {
		t.Errorf("error must name the created request id so the caller can reconcile, got %q", err)
	}
	if *approved {
		t.Error("nothing should have been approved")
	}
}

// The reason ApproveRequests re-verifies instead of trusting HashVerified: the flag is
// a serialized field, so this is what a request decoded from a queue, webhook or cached
// blob looks like. Trusting it gets an attacker-chosen hash signed by the approver's
// real key.
func TestApproveRefusesForgedHashVerifiedFromJSON(t *testing.T) {
	good := `[{"key":"currency","value":"BTC"}]`
	forged := fmt.Sprintf(
		`{"id":"9","metadata":{"hash":%q,"payload_as_string":%q,"hash_verified":true}}`,
		crypto.CalculateHexHash(good), `[{"key":"currency","value":"ETH"}]`,
	)

	var r model.Request
	if err := json.Unmarshal([]byte(forged), &r); err != nil {
		t.Fatal(err)
	}
	// The forgery worked at the model level — that is the point.
	if !r.Metadata.HashVerified {
		t.Fatal("fixture is not exercising the forgery: HashVerified did not survive the decode")
	}

	svc, approved := createAndApproveService(t, requestRow("9", good, crypto.CalculateHexHash(good)))
	if _, err := svc.ApproveRequests(context.Background(), []*model.Request{&r}, testApprovalKey(t), "ok"); err == nil {
		t.Fatal("signed a hash that does not cover its payload, on the strength of a forged flag")
	} else if !strings.Contains(err.Error(), "refusing to sign request 9") {
		t.Errorf("unexpected error: %v", err)
	}
	if *approved {
		t.Error("a refusal must never reach the wire")
	}
}

// A hash with no payload is the shape there is nothing to check against; approving it
// would mean signing whatever the response claimed.
func TestApproveRefusesHashWithoutPayload(t *testing.T) {
	r := &model.Request{
		ID:       "5",
		Metadata: &model.RequestMetadata{Hash: crypto.CalculateHexHash("x"), HashVerified: true},
	}

	svc, approved := createAndApproveService(t, requestRow("5", "x", crypto.CalculateHexHash("x")))
	if _, err := svc.ApproveRequests(context.Background(), []*model.Request{r}, testApprovalKey(t), "ok"); err == nil {
		t.Fatal("approved metadata carrying a hash with no payload")
	}
	if *approved {
		t.Error("a refusal must never reach the wire")
	}
}

// One signature covers every hash in the batch, so a partial approval would mean the
// caller believes they approved more than they did.
func TestApproveIsAllOrNothing(t *testing.T) {
	good := `[{"key":"currency","value":"BTC"}]`
	ok := &model.Request{
		ID:       "1",
		Metadata: &model.RequestMetadata{Hash: crypto.CalculateHexHash(good), PayloadAsString: good},
	}
	bad := &model.Request{
		ID:       "2",
		Metadata: &model.RequestMetadata{Hash: crypto.CalculateHexHash(good), PayloadAsString: `[{"key":"currency","value":"ETH"}]`},
	}

	svc, approved := createAndApproveService(t, requestRow("1", good, crypto.CalculateHexHash(good)))
	if _, err := svc.ApproveRequests(context.Background(), []*model.Request{ok, bad}, testApprovalKey(t), "ok"); err == nil {
		t.Fatal("approved a batch containing an unverifiable request")
	}
	if *approved {
		t.Error("nothing should have been signed")
	}
}
