package integration

import (
	"context"
	"testing"
)

func TestIntegration_HealthCheck(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	result, err := client.Health().GetAllHealthChecks(ctx, nil)
	if err != nil {
		t.Fatalf("GetAllHealthChecks() error = %v", err)
	}

	t.Logf("Health checks: %d components", len(result.Components))
	for name, component := range result.Components {
		t.Logf("  %s: %d groups", name, len(component.Groups))
	}
}
