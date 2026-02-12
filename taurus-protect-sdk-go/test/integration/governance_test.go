package integration

import (
	"context"
	"testing"
)

func TestIntegration_GetGovernanceRules(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	rules, err := client.GovernanceRules().GetRules(ctx)
	if err != nil {
		t.Fatalf("GetRules() error = %v", err)
	}

	t.Logf("Governance rules:")
	t.Logf("  Locked: %v", rules.Locked)
	t.Logf("  CreatedAt: %v", rules.CreatedAt)
	t.Logf("  UpdatedAt: %v", rules.UpdatedAt)
	t.Logf("  RulesContainer length: %d", len(rules.RulesContainer))
}
