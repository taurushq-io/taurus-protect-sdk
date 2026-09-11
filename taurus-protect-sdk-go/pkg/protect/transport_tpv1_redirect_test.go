package protect

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
)

// The api-key (TPV1) transport had neither of the two guards the bearer transport has, even
// though it has the same structural problem and a worse consequence.
//
// Because the Authorization header is created inside RoundTrip — below net/http's redirect
// handling — the follow-up request net/http synthesises from a server-supplied Location carries
// no Authorization for net/http to strip. The transport then mints a BRAND-NEW, fully valid
// TPV1 signature (fresh nonce and timestamp) over the attacker-chosen method, host, path, query
// and body. That is not a stale credential being replayed; it is a signing oracle. A 303 turns
// any signed call into an authenticated GET of the intermediary's choosing, with the response
// flowing back through them.
//
// Two guards, both required, mirroring the bearer pair:
//   - newHTTPClient refuses to follow redirects by default
//   - TPV1Transport refuses to sign for any host but the configured one, which is what holds
//     when a caller supplies a redirect-following policy of their own

func tpv1Auth(t *testing.T) *crypto.TPV1Auth {
	t.Helper()
	auth, err := crypto.NewTPV1Auth("api-key-id", "0badc0de0badc0de0badc0de0badc0de")
	if err != nil {
		t.Fatalf("building TPV1 auth: %v", err)
	}
	return auth
}

func TestTPV1SignatureIsNotSentToARedirectTarget(t *testing.T) {
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

	hc := newHTTPClient(tpv1Auth(t), hostOfURL(t, api.URL), &http.Client{})
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
		t.Fatalf("a freshly minted TPV1 signature leaked to the redirect target: %q", attackerSawAuth)
	}
}

// Even with a caller-supplied redirect-following policy, the transport's host check must keep
// the signature on the configured origin.
func TestTPV1SignatureIsNotSentCrossHostEvenWhenRedirectsAreFollowed(t *testing.T) {
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
	hc := newHTTPClient(tpv1Auth(t), hostOfURL(t, api.URL), follow)

	resp, err := hc.Get(api.URL)
	if err == nil {
		defer func() { _ = resp.Body.Close() }()
		_, _ = io.Copy(io.Discard, resp.Body)
	}

	if attackerSawAuth != "" {
		t.Fatalf("a freshly minted TPV1 signature leaked to the redirect target: %q", attackerSawAuth)
	}
}

// The same-host case must keep working: pinning is about the origin, not about refusing to
// sign at all. Without this, "refuse everything" would pass the two tests above.
func TestTPV1SignatureIsStillProducedForTheConfiguredHost(t *testing.T) {
	var sawAuth string
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer api.Close()

	hc := newHTTPClient(tpv1Auth(t), hostOfURL(t, api.URL), &http.Client{})
	resp, err := hc.Get(api.URL + "/api/rest/v1/wallets")
	if err != nil {
		t.Fatalf("a request for the configured host must be signed and sent: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(io.Discard, resp.Body)

	if !strings.HasPrefix(sawAuth, "TPV1-HMAC-SHA256 ") {
		t.Errorf("Authorization = %q, want a TPV1-HMAC-SHA256 header", sawAuth)
	}
}

// A cross-host request that did NOT arrive via a redirect must be refused too, with an error
// naming both hosts — the caller otherwise sees an opaque failure.
func TestTPV1TransportRefusesAForeignHostDirectly(t *testing.T) {
	transport := &TPV1Transport{Auth: tpv1Auth(t), Host: "api.example.test"}

	req, err := http.NewRequest(http.MethodGet, "https://attacker.test/collect", nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := transport.RoundTrip(req); err == nil {
		t.Fatal("signed a request for a host the client is not configured for")
	} else if !strings.Contains(err.Error(), "attacker.test") ||
		!strings.Contains(err.Error(), "api.example.test") {
		t.Errorf("error should name both hosts, got %q", err)
	}
}

// api-key credentials must fail closed on a host they cannot parse, matching the bearer half.
// Silently leaving the pin empty is how the guard gets disabled without anyone noticing.
func TestAPIKeyCredentialsRequireAParseableHost(t *testing.T) {
	_, _, err := apiKeyCredentials{apiKey: "k", apiSecret: "0badc0de0badc0de0badc0de0badc0de"}.apply("://not-a-url", nil)
	if err == nil {
		t.Fatal("expected api-key credentials to refuse an unparseable host")
	}
	if !strings.Contains(err.Error(), "another origin") {
		t.Errorf("error should explain why the host is required, got %q", err)
	}
}
