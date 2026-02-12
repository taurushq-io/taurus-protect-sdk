package integration

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/test/testutil"
)

// DefaultSuperAdminKeysPEM provides backward compatibility for integration tests
// that reference SuperAdmin keys directly. Loaded from testutil.
var DefaultSuperAdminKeysPEM []string

// DefaultMinValidSignatures provides backward compatibility for integration tests
// that reference MinValidSignatures directly. Loaded from testutil.
var DefaultMinValidSignatures int

func init() {
	config := testutil.LoadConfig()
	DefaultSuperAdminKeysPEM = config.GetSuperAdminKeys()
	DefaultMinValidSignatures = config.MinValidSignatures
}

// GetConfig returns configuration from the shared testutil config (identity 1).
// Provided for backward compatibility with existing integration tests.
func GetConfig() (host, apiKey, apiSecret string) {
	config := testutil.LoadConfig()
	host = config.Host
	identity := config.GetIdentity(1)
	if identity != nil {
		apiKey = identity.APIKey
		apiSecret = identity.APISecret
	}
	return host, apiKey, apiSecret
}

// IsIntegrationEnabled returns true if integration tests should run.
func IsIntegrationEnabled() bool {
	return testutil.LoadConfig().IsEnabled()
}
