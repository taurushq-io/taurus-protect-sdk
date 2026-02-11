package integration

import (
	"crypto/ecdsa"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/test/testutil"
)

func skipIfNotIntegration(t *testing.T) {
	t.Helper()
	testutil.SkipIfNotEnabled(t)
}

// getTestClient returns a client configured with SuperAdmin keys for integrity verification.
// Uses identity 1 (default user) from the shared test config.
func getTestClient(t *testing.T) *protect.Client {
	t.Helper()
	return testutil.GetTestClient(t, 1)
}

// getTestClientWithVerification returns a client configured with SuperAdmin keys
// for governance signature verification. Uses identity 1 from the shared test config.
func getTestClientWithVerification(t *testing.T) *protect.Client {
	t.Helper()
	return testutil.GetTestClient(t, 1)
}

// getTeam1PrivateKey returns the test EC private key for request approval signing.
// Uses the private key from identity 1 in the shared test config.
func getTeam1PrivateKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()

	pem := testutil.GetPrivateKey(t, 1)
	if pem == "" {
		t.Fatal("Identity 1 has no private key configured")
	}

	key, err := crypto.DecodePrivateKeyPEM(pem)
	if err != nil {
		t.Fatalf("Failed to decode private key for identity 1: %v", err)
	}

	return key
}
