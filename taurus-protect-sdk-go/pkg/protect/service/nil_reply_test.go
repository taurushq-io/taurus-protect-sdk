package service

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

// Every read/list service follows the same shape: `resp, httpResp, err := req.Execute()`, then
// `resp.<Field>` as soon as err is nil. The generated decode path used to report SUCCESS
// without populating the typed reply for two bodies a server fully controls — an empty one, and
// the JSON literal `null` — so Execute returned (nil, resp, nil) and the very next line was a
// nil-pointer dereference.
//
// That is a whole-process kill, not a failed call: nothing in the SDK recovers, so in a
// consumer like tg-protect-mcpd (one client per tenant in one daemon) a single empty 200 aimed
// at one tenant takes down every tenant, repeatably on every retry. A server that merely errors
// or returns malformed JSON cannot do that — it only fails the caller's own request.
//
// ~110 unchecked dereferences across 38 service files, so the guard lives in the shared decode
// path (assertReplyDecoded in internal/openapi/client.go, and in the vendored template at
// scripts/resources/templates/go/client.mustache so it survives regeneration) rather than in
// each service.

func nilReplyClient(t *testing.T, status int, body string) *openapi.APIClient {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != "" {
			_, _ = fmt.Fprint(w, body)
		}
	}))
	t.Cleanup(srv.Close)

	cfg := openapi.NewConfiguration()
	cfg.Servers = []openapi.ServerConfiguration{{URL: srv.URL}}
	cfg.HTTPClient = srv.Client()
	return openapi.NewAPIClient(cfg)
}

// The two bodies that used to panic. Each subtest would have crashed the test binary before the
// fix — a panic, not a failure — which is why this is written as a table over several services
// rather than a single case: the defect was in the shared path, so it showed up everywhere.
func TestServicesDoNotPanicOnANilReply(t *testing.T) {
	// Both of these reported SUCCESS with a nil reply before the fix. A whitespace-only or
	// malformed body is NOT in this table: it already failed closed through the JSON parser,
	// so including it would blur what this test pins.
	bodies := map[string]string{
		"empty body":     "",
		"JSON null body": "null",
	}

	for name, body := range bodies {
		t.Run(name, func(t *testing.T) {
			client := nilReplyClient(t, http.StatusOK, body)
			ctx := context.Background()

			calls := map[string]func() error{
				"GetAction": func() error {
					_, err := NewActionService(client).GetAction(ctx, "1")
					return err
				},
				"ListActions": func() error {
					_, err := NewActionService(client).ListActions(ctx, nil)
					return err
				},
				"GetChange": func() error {
					_, err := NewChangeService(client).GetChange(ctx, "1")
					return err
				},
				"ListWallets": func() error {
					_, _, err := NewWalletService(client).ListWallets(ctx, nil)
					return err
				},
				"ListBlockchains": func() error {
					_, err := NewBlockchainService(client).ListBlockchains(ctx, nil)
					return err
				},
				"GetCurrencies": func() error {
					_, err := NewCurrencyService(client).GetCurrencies(ctx, nil)
					return err
				},
			}

			for label, call := range calls {
				t.Run(label, func(t *testing.T) {
					// The assertion is that this RETURNS. A panic here fails the whole
					// binary, which is precisely the pre-fix behaviour.
					err := call()
					if err == nil {
						t.Fatalf("a 2xx with a %s produced no error; the reply cannot have "+
							"been populated, so the caller would dereference nil", name)
					}
					if !strings.Contains(err.Error(), "no usable") {
						t.Errorf("error should explain the empty reply, got %q", err)
					}
				})
			}
		})
	}
}

// A legitimately empty body must still be accepted where the caller is not expecting a typed
// reply. Without this, "reject every empty body" would pass the test above while breaking every
// endpoint that returns google.protobuf.Empty.
func TestEmptyBodyIsStillFineWhenNoTypedReplyIsExpected(t *testing.T) {
	client := nilReplyClient(t, http.StatusOK, "")

	// A void endpoint: nothing to populate, so an empty body is the normal answer.
	err := NewTagService(client).DeleteTag(context.Background(), "1")
	if err != nil {
		t.Errorf("an empty body is the expected reply for a void endpoint: %v", err)
	}
}

// The guard must not swallow a real payload.
func TestPopulatedReplyIsUnaffected(t *testing.T) {
	client := nilReplyClient(t, http.StatusOK, `{"blockchains":[{"id":"1","name":"ETH"}]}`)

	blockchains, err := NewBlockchainService(client).ListBlockchains(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(blockchains) != 1 {
		t.Errorf("got %d blockchains, want 1", len(blockchains))
	}
}
