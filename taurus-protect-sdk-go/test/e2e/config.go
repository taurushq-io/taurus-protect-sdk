package e2e

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/test/testutil"
)

// IsIntegrationEnabled returns true if E2E tests should run.
func IsIntegrationEnabled() bool {
	return testutil.LoadConfig().IsEnabled()
}
