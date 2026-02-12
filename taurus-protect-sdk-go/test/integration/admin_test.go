package integration

import (
	"context"
	"testing"
)

func TestIntegration_ListGroups(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	result, err := client.Groups().ListGroups(ctx, nil)
	if err != nil {
		t.Fatalf("ListGroups() error = %v", err)
	}

	t.Logf("Found %d groups", len(result.Groups))
	for _, g := range result.Groups {
		t.Logf("Group: ID=%s, Name=%s", g.ID, g.Name)
	}
}

func TestIntegration_ListVisibilityGroups(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	groups, err := client.VisibilityGroups().ListVisibilityGroups(ctx)
	if err != nil {
		t.Fatalf("ListVisibilityGroups() error = %v", err)
	}

	t.Logf("Found %d visibility groups", len(groups))
	for _, g := range groups {
		t.Logf("VisibilityGroup: ID=%s, Name=%s", g.ID, g.Name)
	}
}

func TestIntegration_GetTenantConfig(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	config, err := client.Config().GetTenantConfig(ctx)
	if err != nil {
		t.Fatalf("GetTenantConfig() error = %v", err)
	}

	t.Logf("Tenant config:")
	t.Logf("  TenantID: %s", config.TenantID)
	t.Logf("  BaseCurrency: %s", config.BaseCurrency)
	t.Logf("  IsMFAMandatory: %v", config.IsMFAMandatory)
}

func TestIntegration_ListAuditTrails(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	result, err := client.Audits().ListAuditTrails(ctx, nil)
	if err != nil {
		t.Fatalf("ListAuditTrails() error = %v", err)
	}

	t.Logf("Found %d audit trails", len(result.AuditTrails))
	for _, a := range result.AuditTrails {
		userID := ""
		if a.User != nil {
			userID = a.User.ID
		}
		t.Logf("AuditTrail: ID=%s, Action=%s, Entity=%s, UserID=%s", a.ID, a.Action, a.Entity, userID)
	}
}
