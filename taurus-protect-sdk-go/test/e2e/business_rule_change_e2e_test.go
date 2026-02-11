package e2e

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/test/testutil"
)

// TestBusinessRuleChangeE2E exercises the full business rule change lifecycle:
// list rules, find a target, propose a change with one admin, approve with another,
// verify the update, then restore the original value.
//
// Requires at least 3 identities:
//   - Identity 1: default user (reader)
//   - Identity 2: admin who proposes changes
//   - Identity 3: admin who approves changes
func TestBusinessRuleChangeE2E(t *testing.T) {
	testutil.SkipIfNotEnabled(t)
	testutil.SkipIfInsufficientIdentities(t, 3)

	reader := testutil.GetTestClient(t, 1)
	defer reader.Close()
	proposer := testutil.GetTestClient(t, 2)
	defer proposer.Close()
	approver := testutil.GetTestClient(t, 3)
	defer approver.Close()

	ctx := context.Background()

	// Step 1: List all business rules with cursor pagination
	t.Log("=== Step 1: Listing all business rules ===")
	allRules := listAllBusinessRules(t, ctx, reader)
	t.Logf("Found %d business rules", len(allRules))

	// Print global rules
	t.Log("\n--- Global rules ---")
	for _, rule := range allRules {
		if strings.EqualFold(rule.EntityType, "global") {
			t.Logf("  %-45s = %-20s [group: %s]", rule.RuleKey, rule.RuleValue, rule.RuleGroup)
		}
	}

	// Print XLM-specific rules
	t.Log("\n--- XLM rules ---")
	for _, rule := range allRules {
		if strings.EqualFold(rule.Currency, "XLM") {
			t.Logf("  id=%-6s %-45s = %-20s [group: %s, entityType: %s]",
				rule.ID, rule.RuleKey, rule.RuleValue, rule.RuleGroup, rule.EntityType)
		}
	}

	// Step 2: Find target rule (transaction rule for XLM with numeric value)
	t.Log("\n=== Step 2: Finding target rule ===")
	targetRule := findTargetRule(t, allRules)
	if targetRule == nil {
		t.Fatal("No suitable business rule found for testing")
	}

	originalValue := targetRule.RuleValue
	targetRuleID := targetRule.ID
	t.Logf("Target rule: id=%s key=%s originalValue=%s", targetRuleID, targetRule.RuleKey, originalValue)

	// Step 3: Proposer creates a change (increment value by 1)
	t.Log("\n=== Step 3: Proposer creating change ===")
	origNum, err := strconv.ParseInt(originalValue, 10, 64)
	if err != nil {
		t.Fatalf("Failed to parse original value %q as integer: %v", originalValue, err)
	}
	newValue := strconv.FormatInt(origNum+1, 10)

	changeResult, err := proposer.Changes().CreateChange(ctx, &model.CreateChangeRequest{
		Action:   "update",
		Entity:   "businessrule",
		EntityID: targetRuleID,
		Changes:  map[string]string{"rulevalue": newValue},
		Comment:  fmt.Sprintf("E2E test: temporarily change value from %s to %s", originalValue, newValue),
	})
	if err != nil {
		t.Fatalf("CreateChange failed: %v", err)
	}
	if changeResult == nil || changeResult.ID == "" {
		t.Fatal("CreateChange returned nil or empty ID")
	}
	changeID := changeResult.ID
	t.Logf("Created change: id=%s (value %s -> %s)", changeID, originalValue, newValue)

	// Step 4: Approver approves the change
	t.Log("\n=== Step 4: Approver approving change ===")
	err = approver.Changes().ApproveChange(ctx, changeID)
	if err != nil {
		t.Fatalf("ApproveChange failed: %v", err)
	}
	t.Logf("Change %s approved by approver", changeID)

	// Step 5: Verify the change took effect
	t.Log("\n=== Step 5: Verifying change ===")
	updatedRule := waitForRuleValue(t, ctx, reader, targetRuleID, newValue)
	if updatedRule == nil {
		t.Fatalf("Rule %s did not reach expected value %s within timeout", targetRuleID, newValue)
	}
	if updatedRule.RuleValue != newValue {
		t.Fatalf("Expected rule value %s, got %s", newValue, updatedRule.RuleValue)
	}
	t.Logf("BEFORE: %s", originalValue)
	t.Logf("AFTER:  %s", updatedRule.RuleValue)
	t.Logf("Verified: rule %s value changed successfully", targetRuleID)

	// Step 6: Restore original value (cleanup)
	t.Log("\n=== Step 6: Restoring original value ===")
	restoreResult, err := proposer.Changes().CreateChange(ctx, &model.CreateChangeRequest{
		Action:   "update",
		Entity:   "businessrule",
		EntityID: targetRuleID,
		Changes:  map[string]string{"rulevalue": originalValue},
		Comment:  fmt.Sprintf("E2E test: restore value from %s to %s", newValue, originalValue),
	})
	if err != nil {
		t.Fatalf("CreateChange (restore) failed: %v", err)
	}
	restoreChangeID := restoreResult.ID
	t.Logf("Created restore change: id=%s", restoreChangeID)

	err = approver.Changes().ApproveChange(ctx, restoreChangeID)
	if err != nil {
		t.Fatalf("ApproveChange (restore) failed: %v", err)
	}
	t.Logf("Restore change %s approved by approver", restoreChangeID)

	restoredRule := waitForRuleValue(t, ctx, reader, targetRuleID, originalValue)
	if restoredRule == nil {
		t.Fatalf("Rule %s did not restore to value %s within timeout", targetRuleID, originalValue)
	}
	if restoredRule.RuleValue != originalValue {
		t.Fatalf("Expected restored value %s, got %s", originalValue, restoredRule.RuleValue)
	}
	t.Logf("BEFORE: %s", newValue)
	t.Logf("AFTER:  %s", restoredRule.RuleValue)
	t.Logf("Verified: rule %s value restored successfully", targetRuleID)

	t.Log("\n=== E2E PASSED ===")
}

// listAllBusinessRules retrieves all business rules using cursor pagination.
func listAllBusinessRules(t *testing.T, ctx context.Context, client *protect.Client) []*model.BusinessRule {
	t.Helper()

	var allRules []*model.BusinessRule
	opts := &model.ListBusinessRulesOptions{
		PageSize: 50,
	}

	// First page
	result, err := client.BusinessRules().ListBusinessRules(ctx, opts)
	if err != nil {
		t.Fatalf("ListBusinessRules failed: %v", err)
	}
	allRules = append(allRules, result.BusinessRules...)

	// Subsequent pages
	for result.HasNext {
		opts.CurrentPage = result.CurrentPage
		opts.PageRequest = "NEXT"
		result, err = client.BusinessRules().ListBusinessRules(ctx, opts)
		if err != nil {
			t.Fatalf("ListBusinessRules (page) failed: %v", err)
		}
		allRules = append(allRules, result.BusinessRules...)
	}

	return allRules
}

// findTargetRule finds a suitable business rule for testing.
// Prefers XLM transaction-related rules with numeric values, falls back to any numeric rule.
func findTargetRule(t *testing.T, rules []*model.BusinessRule) *model.BusinessRule {
	t.Helper()

	// Prefer XLM transaction rule
	for _, rule := range rules {
		if rule.RuleKey == "" {
			continue
		}
		isTransactionRule := strings.Contains(strings.ToLower(rule.RuleKey), "transaction")
		isXLM := strings.EqualFold(rule.Currency, "XLM")
		if isTransactionRule && isXLM {
			if _, err := strconv.ParseInt(rule.RuleValue, 10, 64); err == nil {
				t.Logf("Found target rule: id=%s key=%s value=%s", rule.ID, rule.RuleKey, rule.RuleValue)
				return rule
			}
		}
	}

	// Fallback: any rule with a numeric value
	for _, rule := range rules {
		if rule.RuleValue == "" {
			continue
		}
		if _, err := strconv.ParseInt(rule.RuleValue, 10, 64); err == nil {
			t.Logf("Fallback target rule: id=%s key=%s value=%s", rule.ID, rule.RuleKey, rule.RuleValue)
			return rule
		}
	}

	// Print diagnostics
	t.Log("No suitable rule found. First 50 rule keys:")
	for i, rule := range rules {
		if i >= 50 {
			break
		}
		t.Logf("  id=%s key=%s value=%s currency=%s", rule.ID, rule.RuleKey, rule.RuleValue, rule.Currency)
	}

	return nil
}

// findRuleByID searches all business rules for one with the given ID.
func findRuleByID(ctx context.Context, client *protect.Client, ruleID string) (*model.BusinessRule, error) {
	opts := &model.ListBusinessRulesOptions{
		PageSize: 50,
	}

	result, err := client.BusinessRules().ListBusinessRules(ctx, opts)
	if err != nil {
		return nil, err
	}
	for _, rule := range result.BusinessRules {
		if rule.ID == ruleID {
			return rule, nil
		}
	}

	for result.HasNext {
		opts.CurrentPage = result.CurrentPage
		opts.PageRequest = "NEXT"
		result, err = client.BusinessRules().ListBusinessRules(ctx, opts)
		if err != nil {
			return nil, err
		}
		for _, rule := range result.BusinessRules {
			if rule.ID == ruleID {
				return rule, nil
			}
		}
	}

	return nil, nil
}

// waitForRuleValue polls until the business rule has the expected value, or times out after 30 seconds.
func waitForRuleValue(t *testing.T, ctx context.Context, client *protect.Client, ruleID, expectedValue string) *model.BusinessRule {
	t.Helper()

	deadline := time.Now().Add(30 * time.Second)
	var lastRule *model.BusinessRule

	for time.Now().Before(deadline) {
		rule, err := findRuleByID(ctx, client, ruleID)
		if err != nil {
			t.Logf("  Warning: findRuleByID error: %v", err)
		} else if rule != nil {
			lastRule = rule
			if rule.RuleValue == expectedValue {
				return rule
			}
			t.Logf("  Waiting for rule %s to have value %s (current: %s)", ruleID, expectedValue, rule.RuleValue)
		}
		time.Sleep(2 * time.Second)
	}

	return lastRule
}
