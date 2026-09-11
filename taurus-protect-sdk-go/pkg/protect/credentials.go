package protect

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
)

// Credentials is a client authentication mechanism: either static TPV1-HMAC
// (API key + secret) or a Bearer token resolved per request. Construct one with
// APIKeyCredentials, BearerTokenProviderCredentials, or BearerTokenCredentials
// and pass it via WithCredentials. The interface is closed — only these three
// implementations exist — so the client applies it without null-checking.
// SuperAdmin keys are required regardless of the mechanism (client-side rules
// verification is mandatory).
type Credentials interface {
	// apply builds the authenticating http.Client from base. For TPV1-HMAC it also
	// returns the TPV1Auth to wipe on Close; bearer returns a nil auth.
	apply(host string, base *http.Client) (*http.Client, *crypto.TPV1Auth, error)
}

// APIKeyCredentials authenticates with static TPV1-HMAC credentials. The
// apiSecret must be hex-encoded.
func APIKeyCredentials(apiKey, apiSecret string) Credentials {
	return apiKeyCredentials{apiKey: apiKey, apiSecret: apiSecret}
}

// BearerTokenProviderCredentials authenticates with a Bearer token resolved per
// request from provider(ctx) — so one client can serve many callers/tokens.
func BearerTokenProviderCredentials(provider BearerTokenProvider) Credentials {
	return bearerCredentials{provider: provider}
}

// BearerTokenCredentials is a convenience for a single static Bearer token (a
// constant provider), for single-user/CLI use.
func BearerTokenCredentials(token string) Credentials {
	return bearerCredentials{
		token:    token,
		static:   true,
		provider: func(context.Context) (string, error) { return token, nil },
	}
}

type apiKeyCredentials struct {
	apiKey    string
	apiSecret string
}

func (c apiKeyCredentials) apply(host string, base *http.Client) (*http.Client, *crypto.TPV1Auth, error) {
	if c.apiKey == "" {
		return nil, nil, errors.New("apiKey is required")
	}
	if c.apiSecret == "" {
		return nil, nil, errors.New("apiSecret is required")
	}
	auth, err := crypto.NewTPV1Auth(c.apiKey, c.apiSecret)
	if err != nil {
		return nil, nil, err
	}
	// The signature must not be minted for any origin but the configured one, so the
	// transport needs to know it. This used to discard the host (`_ string`), which left
	// the api-key client free to sign for whatever host a redirect named — the bearer half
	// of this file has parsed and forwarded it for exactly that reason since 2026-09.
	// An unparseable host fails closed rather than disabling the check silently.
	parsed, err := url.Parse(host)
	if err != nil || parsed.Host == "" {
		return nil, nil, fmt.Errorf("cannot determine the API host from %q: api-key "+
			"credentials require it so a request is never signed for another origin", host)
	}
	return newHTTPClient(auth, parsed.Host, base), auth, nil
}

type bearerCredentials struct {
	// token is set only by BearerTokenCredentials, so a static empty token is rejected
	// at construction rather than surfacing as an opaque 401 on every request. A
	// provider-supplied token is validated per request instead, in bearerTransport.
	token    string
	static   bool
	provider BearerTokenProvider
}

func (c bearerCredentials) apply(host string, base *http.Client) (*http.Client, *crypto.TPV1Auth, error) {
	if c.provider == nil {
		return nil, nil, errors.New("bearer token provider is required")
	}
	if c.static && c.token == "" {
		return nil, nil, errors.New("bearerToken is required")
	}
	// The token must not leave the configured origin, so the transport needs to know it.
	// An unparseable host fails closed here rather than disabling the check silently.
	parsed, err := url.Parse(host)
	if err != nil || parsed.Host == "" {
		return nil, nil, fmt.Errorf("cannot determine the API host from %q: bearer "+
			"credentials require it so the token is never sent elsewhere", host)
	}
	return newBearerHTTPClient(c.provider, parsed.Host, base), nil, nil
}
