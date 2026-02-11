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
	apiKey             string
	apiSecret          string
	superAdminKeys     []*ecdsa.PublicKey
	minValidSignatures int
	rulesCacheTTL      time.Duration
	httpClient         *http.Client
	httpTimeout        time.Duration
}

// validate checks the configuration for required fields and valid values.
func (c *clientConfig) validate() error {
	if c.host == "" {
		return errors.New("host is required")
	}
	if c.apiKey == "" {
		return errors.New("apiKey is required")
	}
	if c.apiSecret == "" {
		return errors.New("apiSecret is required")
	}
	if len(c.superAdminKeys) == 0 {
		return errors.New("superAdminKeys are required: at least one SuperAdmin public key must be provided for integrity verification")
	}
	if c.minValidSignatures < 0 {
		return errors.New("minValidSignatures must be non-negative")
	}
	if len(c.superAdminKeys) > 0 && c.minValidSignatures > len(c.superAdminKeys) {
		return errors.New("minValidSignatures cannot exceed number of superAdminKeys")
	}
	// Reject minValidSignatures=0 when SuperAdmin keys are provided.
	// This matches Java SDK behavior (IllegalArgumentException).
	if len(c.superAdminKeys) > 0 && c.minValidSignatures == 0 {
		return errors.New("minValidSignatures must be greater than zero when SuperAdmin keys are provided")
	}
	return nil
}

// Option configures a ProtectClient.
type Option func(*clientConfig) error

// WithCredentials sets the API key and secret for authentication.
// The apiSecret should be hex-encoded.
func WithCredentials(apiKey, apiSecret string) Option {
	return func(c *clientConfig) error {
		if apiKey == "" {
			return errors.New("apiKey cannot be empty")
		}
		if apiSecret == "" {
			return errors.New("apiSecret cannot be empty")
		}
		c.apiKey = apiKey
		c.apiSecret = apiSecret
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
