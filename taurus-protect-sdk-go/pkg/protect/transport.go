package protect

import (
	"bytes"
	"io"
	"net/http"

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
func newHTTPClient(auth *crypto.TPV1Auth, base *http.Client) *http.Client {
	var baseTransport http.RoundTripper
	if base != nil && base.Transport != nil {
		baseTransport = base.Transport
	}

	return &http.Client{
		Transport: &TPV1Transport{
			Base: baseTransport,
			Auth: auth,
		},
		Timeout: base.Timeout,
	}
}
