package testutil

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"
)

// SkipIfNotEnabled skips the test if integration/e2e tests are not enabled.
func SkipIfNotEnabled(t *testing.T) {
	t.Helper()
	config := LoadConfig()
	if !config.IsEnabled() {
		t.Skip("Skipping test. Set PROTECT_INTEGRATION_TEST=true or configure test.properties.")
	}
}

// GetTestClient creates a ProtectClient from the identity at the given 1-based index.
// The client is configured with SuperAdmin keys for integrity verification.
func GetTestClient(t *testing.T, identityIndex int) *protect.Client {
	t.Helper()
	config := LoadConfig()

	identity := config.GetIdentity(identityIndex)
	if identity == nil {
		t.Fatalf("Identity %d not configured", identityIndex)
	}
	if !identity.HasAPICredentials() {
		t.Fatalf("Identity %d (%s) has no API credentials", identityIndex, identity.Name)
	}
	if config.Host == "" {
		t.Fatal("API host not configured (set PROTECT_API_HOST or host in test.properties)")
	}

	opts := []protect.Option{
		protect.WithCredentials(identity.APIKey, identity.APISecret),
	}

	superAdminKeys := config.GetSuperAdminKeys()
	if len(superAdminKeys) > 0 {
		opts = append(opts,
			protect.WithSuperAdminKeysPEM(superAdminKeys),
			protect.WithMinValidSignatures(config.MinValidSignatures),
		)
	}

	client, err := protect.NewClient(config.Host, opts...)
	if err != nil {
		t.Fatalf("Failed to create client for identity %d (%s): %v", identityIndex, identity.Name, err)
	}

	return client
}

// GetPrivateKey returns the PEM-encoded private key for the identity at the given 1-based index.
// Returns empty string if the identity has no private key.
func GetPrivateKey(t *testing.T, identityIndex int) string {
	t.Helper()
	config := LoadConfig()

	identity := config.GetIdentity(identityIndex)
	if identity == nil {
		t.Fatalf("Identity %d not configured", identityIndex)
	}
	return identity.PrivateKey
}

// SkipIfInsufficientIdentities skips the test if fewer than required identities are configured.
func SkipIfInsufficientIdentities(t *testing.T, required int) {
	t.Helper()
	config := LoadConfig()
	if config.GetIdentityCount() < required {
		t.Skipf("Skipping: need %d identities but only %d configured.", required, config.GetIdentityCount())
	}
}
