package integration

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"
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
