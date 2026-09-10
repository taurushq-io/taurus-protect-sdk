package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func wlaListBody(t *testing.T) string {
	t.Helper()
	good := `{"address":"0xabc","label":"good"}`
	rows := []map[string]any{
		{"id": "1", "metadata": map[string]any{
			"hash": crypto.CalculateHexHash(good), "payloadAsString": good}},
		// hash does not cover this payload — the shape that 500s the live endpoint
		{"id": "2", "metadata": map[string]any{
			"hash": crypto.CalculateHexHash(good), "payloadAsString": `{"address":"0xEVIL"}`}},
	}
	b, err := json.Marshal(map[string]any{"result": rows, "totalItems": "2"})
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func wlaService(t *testing.T, body string, withVerifier bool, logger Logger) *WhitelistedAddressService {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, body)
	}))
	t.Cleanup(srv.Close)

	cfg := openapi.NewConfiguration()
	cfg.Servers = []openapi.ServerConfiguration{{URL: srv.URL}}
	cfg.HTTPClient = srv.Client()

	svcCfg := &WhitelistedAddressServiceConfig{}
	if withVerifier {
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		svcCfg.SuperAdminKeys = []*ecdsa.PublicKey{&key.PublicKey}
		svcCfg.MinValidSignatures = 1
	}
	return NewWhitelistedAddressServiceWithVerification(openapi.NewAPIClient(cfg), svcCfg, WithServiceLogger(logger))
}

// Neither row in this fixture survives: row 2's hash does not cover its payload, and
// row 1 carries no signedAddress, so full verification rejects it too.
//
// This test used to assert only that row 2 was absent from the results — which was
// vacuously true, because the result was EMPTY. It passed while the call returned an
// empty whitelist and err == nil, which is exactly the failure the "rows came back but
// none survived" rule exists to stop: a filtered page must never read as an empty
// whitelist.
func TestListWhitelistedAddressesErrorsWhenNoRowSurvives(t *testing.T) {
	logger := &recordingLogger{}
	result, err := wlaService(t, wlaListBody(t), true, logger).
		ListWhitelistedAddresses(context.Background(), nil)

	if err == nil {
		t.Fatalf("all rows failed verification; returning %+v with no error hides that", result)
	}
	if result != nil {
		t.Errorf("no result may accompany the error, got %+v", result)
	}
	var integrity *model.IntegrityError
	if !asIntegrityError(err, &integrity) {
		t.Fatalf("want IntegrityError, got %T", err)
	}
	if !strings.Contains(err.Error(), "all 2 whitelisted address(es) failed verification") {
		t.Errorf("error should say how many rows were dropped, got: %v", err)
	}
	if len(logger.warns) == 0 {
		t.Error("each exclusion must still be logged so an operator can find the bad row")
	}
}

// The server counts rows it returned; the caller receives only those that verified, so
// reporting the server's total lets a filtered page pass for a complete one. HasMore is
// a separate question — whether another page exists — and is answered from the server's
// own total, not from the reduced one.
func TestAdjustedPaginationSubtractsExcludedRows(t *testing.T) {
	// excludedCount covers THIS page; TotalItems is the server's global count. So the
	// reduction is right for TotalItems (the caller cannot read the excluded rows) and
	// wrong for HasMore (pagination position is a server-side fact). Earlier cases here
	// asserted the coupled behaviour using an excluded count larger than the page limit,
	// which cannot happen: you cannot exclude 30 rows from a 20-row page.
	cases := []struct {
		name          string
		total         string
		limit, offset int64
		excluded      int
		wantTotal     int64
		wantMore      bool
	}{
		{"nothing excluded", "50", 20, 0, 0, 50, true},
		{"one excluded", "50", 20, 0, 1, 49, true},
		{"most of a page excluded, more pages remain", "50", 20, 0, 15, 35, true},
		{"last page with exclusions ends pagination", "50", 20, 40, 5, 45, false},
		// The regression: a page-local exclusion count subtracted from a global total
		// ended pagination early — 100+100 < 190 is false, yet 50 rows sit at offset 200.
		{"mid-result page with many exclusions still pages on", "250", 100, 100, 60, 190, true},
		{"more excluded than reported never goes negative", "20", 20, 0, 99, 0, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opts := &model.ListWhitelistedAddressesOptions{Limit: tc.limit, Offset: tc.offset}
			got := adjustedPagination(&tc.total, tc.excluded, opts)
			if got == nil {
				t.Fatal("pagination missing")
			}
			if got.TotalItems != tc.wantTotal {
				t.Errorf("TotalItems = %d, want %d", got.TotalItems, tc.wantTotal)
			}
			if got.HasMore != tc.wantMore {
				t.Errorf("HasMore = %v, want %v", got.HasMore, tc.wantMore)
			}
		})
	}

	opts := &model.ListWhitelistedAddressesOptions{Limit: 20, Offset: 0}
	if adjustedPagination(nil, 0, opts) != nil {
		t.Error("a server that reports no total must not be given a fabricated one")
	}
}

// A missing verifier is CONFIGURATION, not one bad row: it fails identically for
// every address. Treating it per-row would return an empty whitelist presented as
// complete — strictly worse than the 500 being fixed. It must fail the call.
//
// This and the test above must both hold; passing the first while failing this one
// is the dangerous outcome.
func TestListWhitelistedAddressesFailsWithoutVerifier(t *testing.T) {
	result, err := wlaService(t, wlaListBody(t), false, nil).
		ListWhitelistedAddresses(context.Background(), nil)
	if err == nil {
		t.Fatal("a missing verifier must fail the call, not silently empty the list")
	}
	if result != nil {
		t.Errorf("no result may be returned when nothing can be verified, got %+v", result)
	}
	var integrity *model.IntegrityError
	if !asIntegrityError(err, &integrity) {
		t.Errorf("want IntegrityError, got %T", err)
	}
}

func asIntegrityError(err error, target **model.IntegrityError) bool {
	for err != nil {
		if e, ok := err.(*model.IntegrityError); ok {
			*target = e
			return true
		}
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		err = u.Unwrap()
	}
	return false
}
