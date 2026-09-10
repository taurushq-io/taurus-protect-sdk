package protect

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
)

// TPV1Transport is an http.RoundTripper that signs requests with TPV1-HMAC-SHA256.
//
// Important: This transport reads and clones the request body for signing.
// If you use middleware that wraps this transport and expects to retry requests,
// be aware that the original request body is consumed during signing.
// The transport properly resets the body on the cloned request, but middleware
// should not rely on re-reading the original request's body.
//
// Middleware ordering: Place retry middleware ABOVE this transport (closer to
// the http.Client), so retries create fresh requests rather than reusing
// requests with consumed bodies.
type TPV1Transport struct {
	// Base is the underlying transport. If nil, http.DefaultTransport is used.
	Base http.RoundTripper
	// Auth provides TPV1 credentials for signing.
	Auth *crypto.TPV1Auth
}

// RoundTrip executes a single HTTP transaction, signing the request with TPV1.
func (t *TPV1Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Read the body from the original request first
	var body []byte
	if req.Body != nil {
		var err error
		body, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		// Reset the original body so it can be re-read by upstream middleware (e.g., for retries)
		req.Body = io.NopCloser(bytes.NewReader(body))
	}

	// Clone the request to avoid mutating the original headers
	req2 := req.Clone(req.Context())

	// Set the body on the clone for the actual request
	if body != nil {
		req2.Body = io.NopCloser(bytes.NewReader(body))
	}

	// Sign the cloned request
	if err := t.Auth.SignRequest(req2, body); err != nil {
		return nil, err
	}

	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}

	return base.RoundTrip(req2)
}

// newHTTPClient creates an http.Client with TPV1 authentication.
//
// CheckRedirect and Jar are carried over from the caller-supplied client. They used to be
// dropped, which silently discarded a redirect policy a consumer had deliberately set.
func newHTTPClient(auth *crypto.TPV1Auth, base *http.Client) *http.Client {
	var baseTransport http.RoundTripper
	if base != nil && base.Transport != nil {
		baseTransport = base.Transport
	}

	client := &http.Client{
		Transport: &TPV1Transport{
			Base: baseTransport,
			Auth: auth,
		},
	}
	if base != nil {
		client.Timeout = base.Timeout
		client.CheckRedirect = base.CheckRedirect
		client.Jar = base.Jar
	}
	return client
}

// BearerTokenProvider yields the bearer token to use for a request. It is invoked
// once per request with that request's context, so a single client can serve many
// callers/tokens — e.g. a multi-user proxy that carries the per-user token in the
// context. The token is dynamic (unlike static TPV1 credentials).
type BearerTokenProvider func(ctx context.Context) (string, error)

// bearerTransport is an http.RoundTripper that authenticates each request with a
// bearer token obtained per-request from a provider. Unlike TPV1Transport it
// neither reads nor signs the request body.
type bearerTransport struct {
	// Base is the underlying transport. If nil, http.DefaultTransport is used.
	Base http.RoundTripper
	// provider supplies the bearer token for each request.
	provider BearerTokenProvider
	// host is the configured API origin (host[:port]). The token is attached ONLY to
	// requests for this host. See RoundTrip for why the check has to live here.
	host string
}

// RoundTrip executes a single HTTP transaction, adding the per-request bearer
// token to a clone so the caller's original request headers are left untouched.
//
// The token is attached ONLY when the request is for the configured host, and that check
// cannot be delegated upwards. net/http strips Authorization on a cross-host redirect in
// http.Client.do, but it only strips headers present on the request handed to Do — a header
// set here, inside the RoundTripper, is invisible to it, so there is nothing for it to
// strip and it would re-attach a freshly minted token to whatever host the server named in
// a Location. (Same caveat as golang.org/x/oauth2's Transport.) newBearerHTTPClient also
// refuses to follow redirects by default; this check is what holds if a caller supplies a
// redirect-following policy of their own.
func (t *bearerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL == nil || !strings.EqualFold(req.URL.Host, t.host) {
		return nil, fmt.Errorf(
			"refusing to send the bearer token to %q: the client is configured for %q",
			hostOf(req), t.host)
	}

	token, err := t.provider(req.Context())
	if err != nil {
		return nil, fmt.Errorf("bearer token provider: %w", err)
	}
	// A refresh that silently yields nothing would send an empty credential and
	// surface as an opaque 401 rather than the real cause.
	if token == "" {
		return nil, errors.New("bearer token provider returned an empty token")
	}
	req2 := req.Clone(req.Context())
	req2.Header.Set("Authorization", "Bearer "+token)

	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(req2)
}

// hostOf names the request host for an error message without dereferencing a nil URL.
func hostOf(req *http.Request) string {
	if req.URL == nil {
		return ""
	}
	return req.URL.Host
}

// newBearerHTTPClient creates an http.Client that authenticates with a per-request
// bearer token instead of TPV1-HMAC signing.
//
// host is the configured API origin, used to keep the token from leaving it.
//
// Redirects are NOT followed by default: a JSON API client has no reason to, and following
// one is how a bearer token reaches a host the caller never configured. A caller who
// genuinely needs redirects can set CheckRedirect on the client passed to WithHTTPClient —
// it is carried over here, and bearerTransport's own host check still refuses to send the
// token anywhere but the configured origin.
func newBearerHTTPClient(provider BearerTokenProvider, host string, base *http.Client) *http.Client {
	var baseTransport http.RoundTripper
	if base != nil && base.Transport != nil {
		baseTransport = base.Transport
	}

	client := &http.Client{
		Transport: &bearerTransport{
			Base:     baseTransport,
			provider: provider,
			host:     host,
		},
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	if base != nil {
		client.Timeout = base.Timeout
		client.Jar = base.Jar
		if base.CheckRedirect != nil {
			client.CheckRedirect = base.CheckRedirect
		}
	}
	return client
}
