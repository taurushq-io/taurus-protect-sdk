package protect

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
)

// Default configuration values.
const (
	// DefaultRulesCacheTTL is the default TTL for the rules container cache.
	DefaultRulesCacheTTL = 5 * time.Minute
	// DefaultHTTPTimeout is the default timeout for HTTP requests.
	DefaultHTTPTimeout = 30 * time.Second
)

// clientConfig holds the configuration for the ProtectClient.
type clientConfig struct {
	host               string
	credentials        Credentials
	superAdminKeys     []*ecdsa.PublicKey
	minValidSignatures int
	rulesCacheTTL      time.Duration
	httpClient         *http.Client
	httpTimeout        time.Duration
	logger             Logger
}

// validate checks the configuration for required fields and valid values. The
// auth mechanism is a Credentials value (WithCredentials). SuperAdmin keys are
// mandatory regardless of the mechanism, since client-side rules verification is
// mandatory.
func (c *clientConfig) validate() error {
	if c.host == "" {
		return errors.New("host is required")
	}
	if c.credentials == nil {
		return errors.New("credentials are required: pass WithCredentials(...)")
	}
	if len(c.superAdminKeys) == 0 {
		return errors.New("superAdminKeys are required: at least one SuperAdmin public key must be provided for integrity verification")
	}
	if c.minValidSignatures <= 0 {
		return errors.New("minValidSignatures must be greater than zero")
	}
	if c.minValidSignatures > len(c.superAdminKeys) {
		return errors.New("minValidSignatures cannot exceed number of superAdminKeys")
	}
	return nil
}

// Option configures a ProtectClient.
type Option func(*clientConfig) error

// WithCredentials sets the client's authentication mechanism. Construct the
// Credentials with APIKeyCredentials (static TPV1-HMAC),
// BearerTokenProviderCredentials (per-request Bearer token), or
// BearerTokenCredentials (a static Bearer token).
func WithCredentials(credentials Credentials) Option {
	return func(c *clientConfig) error {
		if credentials == nil {
			return errors.New("credentials cannot be nil")
		}
		c.credentials = credentials
		return nil
	}
}

// WithSuperAdminKeysPEM sets the SuperAdmin public keys from PEM-encoded strings.
// These keys are used to verify governance rules signatures.
func WithSuperAdminKeysPEM(pemKeys []string) Option {
	return func(c *clientConfig) error {
		keys := make([]*ecdsa.PublicKey, 0, len(pemKeys))
		for i, pem := range pemKeys {
			key, err := crypto.DecodePublicKeyPEM(pem)
			if err != nil {
				return fmt.Errorf("invalid SuperAdmin key at index %d: %w", i, err)
			}
			keys = append(keys, key)
		}
		c.superAdminKeys = keys
		return nil
	}
}

// WithSuperAdminKeys sets the SuperAdmin public keys directly.
func WithSuperAdminKeys(keys []*ecdsa.PublicKey) Option {
	return func(c *clientConfig) error {
		c.superAdminKeys = keys
		return nil
	}
}

// WithMinValidSignatures sets the minimum number of valid SuperAdmin signatures
// required to verify governance rules.
func WithMinValidSignatures(n int) Option {
	return func(c *clientConfig) error {
		if n < 0 {
			return errors.New("minValidSignatures must be non-negative")
		}
		c.minValidSignatures = n
		return nil
	}
}

// WithRulesCacheTTL sets the TTL for the rules container cache.
// The default is 5 minutes.
func WithRulesCacheTTL(ttl time.Duration) Option {
	return func(c *clientConfig) error {
		if ttl < 0 {
			return errors.New("rulesCacheTTL must be non-negative")
		}
		c.rulesCacheTTL = ttl
		return nil
	}
}

// WithHTTPClient sets a custom HTTP client to use for requests.
// Note: The client's Transport will be wrapped with TPV1 authentication.
func WithHTTPClient(client *http.Client) Option {
	return func(c *clientConfig) error {
		c.httpClient = client
		return nil
	}
}

// WithLogger sets the logger the SDK reports diagnostic events to. The default
// discards them.
//
// The SDK logs metadata only — an identifier, a resource, a reason — never a
// payload, token or key. It reports, among other things, rows dropped from a list
// because their integrity could not be verified; without a logger those exclusions
// are silent, and a shortened list is indistinguishable from a complete one.
func WithLogger(l Logger) Option {
	return func(c *clientConfig) error {
		if l == nil {
			return errors.New("logger must not be nil; omit the option to disable logging")
		}
		c.logger = l
		return nil
	}
}

// WithHTTPTimeout sets the timeout for HTTP requests.
// The default is 30 seconds.
func WithHTTPTimeout(timeout time.Duration) Option {
	return func(c *clientConfig) error {
		if timeout < 0 {
			return errors.New("httpTimeout must be non-negative")
		}
		c.httpTimeout = timeout
		return nil
	}
}
