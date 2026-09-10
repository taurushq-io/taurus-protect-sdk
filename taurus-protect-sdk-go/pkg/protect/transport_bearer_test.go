package protect

import (
	"context"
	"crypto/ecdsa"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// hostOfURL extracts host[:port] from a test server URL, which is what
// newBearerHTTPClient needs so the token stays on the configured origin.
func hostOfURL(t *testing.T, raw string) string {
	t.Helper()
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("parse %q: %v", raw, err)
	}
	return u.Host
}

func TestBearerTransportSetsAuthorizationHeader(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	hc := newBearerHTTPClient(staticBearer("session-token-xyz"), hostOfURL(t, srv.URL), &http.Client{})
	resp, err := hc.Get(srv.URL)
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(io.Discard, resp.Body)

	if want := "Bearer session-token-xyz"; gotAuth != want {
		t.Fatalf("Authorization = %q, want %q", gotAuth, want)
	}
}

// ctxTokenKey is a test context key carrying a per-request bearer token.
type ctxTokenKey struct{}

// TestBearerTransportResolvesPerRequestToken asserts one client serves different
// tokens per request via the provider reading the request context.
func TestBearerTransportResolvesPerRequestToken(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	provider := func(ctx context.Context) (string, error) {
		tok, _ := ctx.Value(ctxTokenKey{}).(string)
		return tok, nil
	}
	hc := newBearerHTTPClient(provider, hostOfURL(t, srv.URL), &http.Client{})

	for _, tok := range []string{"user-a-jwt", "user-b-jwt"} {
		req, err := http.NewRequestWithContext(
			context.WithValue(context.Background(), ctxTokenKey{}, tok),
			http.MethodGet, srv.URL, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := hc.Do(req)
		if err != nil {
			t.Fatalf("Do: %v", err)
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		if want := "Bearer " + tok; gotAuth != want {
			t.Fatalf("Authorization = %q, want %q", gotAuth, want)
		}
	}
}

// TestNewClientBearerSkipsTPV1 asserts that a bearer client uses the bearer
// transport (not TPV1) and creates no TPV1 auth — and that validate() accepts a
// bearer config with no apiKey/apiSecret/SuperAdmin keys.
func TestNewClientBearerSkipsTPV1(t *testing.T) {
	c, err := NewClient("https://api.example.com",
		WithCredentials(BearerTokenCredentials("tok")),
		WithSuperAdminKeys([]*ecdsa.PublicKey{testP256Key(t)}),
		WithMinValidSignatures(1))
	if err != nil {
		t.Fatalf("NewClient(bearer) error = %v", err)
	}
	defer func() { _ = c.Close() }()

	if _, ok := c.HTTPClient().Transport.(*bearerTransport); !ok {
		t.Fatalf("transport = %T, want *bearerTransport", c.HTTPClient().Transport)
	}
	if _, ok := c.HTTPClient().Transport.(*TPV1Transport); ok {
		t.Fatal("bearer client must not use the TPV1 transport")
	}
	if c.auth != nil {
		t.Fatal("bearer client must not create a TPV1 auth")
	}
}

// TestBearerTransportDoesNotMutateOriginalRequest asserts the token is added to a
// clone, leaving the caller's request headers untouched (mirrors TPV1Transport).
func TestBearerTransportDoesNotMutateOriginalRequest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	tr := &bearerTransport{provider: staticBearer("tok"), host: hostOfURL(t, srv.URL)}
	resp, err := tr.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(io.Discard, resp.Body)

	if got := req.Header.Get("Authorization"); got != "" {
		t.Fatalf("original request Authorization = %q, want empty (must not mutate)", got)
	}
}

// TestBearerTokenIsNotSentToARedirectTarget is the credential-leak gate.
//
// The token is attached inside the RoundTripper, BELOW net/http's redirect handling, so
// http.Client.do's cross-host Authorization stripping cannot see it — it strips only
// headers present on the request handed to Do. Left unguarded, the client follows a
// server-chosen Location and the transport mints a fresh token for the attacker's host,
// turning "can tamper with responses" into "holds a replayable credential".
//
// Two independent guards must hold: redirects are not followed by default, and the
// transport refuses to send the token to any host but the configured one.
func TestBearerTokenIsNotSentToARedirectTarget(t *testing.T) {
	var attackerSawAuth string
	attacker := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attackerSawAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer attacker.Close()

	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Redirect(w, &http.Request{}, attacker.URL+"/collect", http.StatusFound)
	}))
	defer api.Close()

	hc := newBearerHTTPClient(staticBearer("victim-jwt"), hostOfURL(t, api.URL), &http.Client{})
	resp, err := hc.Get(api.URL)
	if err == nil {
		defer func() { _ = resp.Body.Close() }()
		_, _ = io.Copy(io.Discard, resp.Body)
		// Not following the redirect is the expected outcome: the 302 is returned as-is.
		if resp.StatusCode != http.StatusFound {
			t.Errorf("expected the 302 to be returned unfollowed, got %d", resp.StatusCode)
		}
	}

	if attackerSawAuth != "" {
		t.Fatalf("bearer token leaked to the redirect target: %q", attackerSawAuth)
	}
}

// Even with a caller-supplied redirect-following policy, the transport's host check must
// keep the token on the configured origin. This is the guard that survives a consumer who
// deliberately re-enables redirects.
func TestBearerTokenIsNotSentCrossHostEvenWhenRedirectsAreFollowed(t *testing.T) {
	var attackerSawAuth string
	attacker := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attackerSawAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer attacker.Close()

	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Redirect(w, &http.Request{}, attacker.URL+"/collect", http.StatusFound)
	}))
	defer api.Close()

	follow := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return nil }}
	hc := newBearerHTTPClient(staticBearer("victim-jwt"), hostOfURL(t, api.URL), follow)

	resp, err := hc.Get(api.URL)
	if err == nil {
		defer func() { _ = resp.Body.Close() }()
		_, _ = io.Copy(io.Discard, resp.Body)
	}

	if attackerSawAuth != "" {
		t.Fatalf("bearer token leaked to the redirect target: %q", attackerSawAuth)
	}
}
