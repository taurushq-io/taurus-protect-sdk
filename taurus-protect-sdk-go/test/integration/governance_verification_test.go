package integration

import (
	"context"
	"strings"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"
)

func TestIntegration_VerifyGovernanceRulesSignatures(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClientWithVerification(t)
	defer client.Close()

	ctx := context.Background()
	rules, err := client.GovernanceRules().GetRules(ctx)
	if err != nil {
		t.Fatalf("GetRules() error = %v", err)
	}

	if rules == nil {
		t.Skip("No governance rules available")
	}

	t.Logf("Governance rules retrieved:")
	t.Logf("  Locked: %v", rules.Locked)
	t.Logf("  CreatedAt: %v", rules.CreatedAt)
	t.Logf("  RulesContainer length: %d bytes", len(rules.RulesContainer))
	t.Logf("  Signatures count: %d", len(rules.Signatures))

	// Log signature user IDs
	for i, sig := range rules.Signatures {
		t.Logf("  Signature[%d] userID: %s", i, sig.UserID)
	}

	// Verify that GetDecodedRulesContainer performs signature verification
	// (this is the main test - verification happens internally)
	decoded, err := client.GovernanceRules().GetDecodedRulesContainer(rules)
	if err != nil {
		t.Fatalf("GetDecodedRulesContainer() with verification error = %v", err)
	}

	t.Logf("Signature verification PASSED")
	t.Logf("Decoded rules container:")
	t.Logf("  Users count: %d", len(decoded.Users))
	t.Logf("  Groups count: %d", len(decoded.Groups))
	t.Logf("  AddressWhitelistingRules count: %d", len(decoded.AddressWhitelistingRules))
	t.Logf("  ContractAddressWhitelistingRules count: %d", len(decoded.ContractAddressWhitelistingRules))
	t.Logf("  TransactionRules count: %d", len(decoded.TransactionRules))
}

func TestIntegration_VerifyGovernanceRulesWithInvalidKeys(t *testing.T) {
	skipIfNotIntegration(t)

	host, apiKey, apiSecret := GetConfig()

	// Create client with invalid SuperAdmin keys (all zeros is not a valid EC point)
	// The SDK correctly rejects invalid keys at parse time
	invalidKeys := []string{
		`-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==
-----END PUBLIC KEY-----`,
	}

	_, err := protect.NewClient(host,
		protect.WithCredentials(apiKey, apiSecret),
		protect.WithSuperAdminKeysPEM(invalidKeys),
		protect.WithMinValidSignatures(1),
	)

	// Should fail because the invalid key cannot be parsed (not a valid EC point)
	if err == nil {
		t.Fatal("Expected client creation to fail with invalid keys, but it succeeded")
	}

	// Verify the error message indicates key parsing failure
	if !strings.Contains(err.Error(), "invalid SuperAdmin key") {
		t.Fatalf("Expected error about invalid SuperAdmin key, got: %v", err)
	}

	t.Logf("Correctly rejected invalid SuperAdmin keys at parse time: %v", err)
}

func TestIntegration_GetDecodedRulesContainer(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClientWithVerification(t)
	defer client.Close()

	ctx := context.Background()
	rules, err := client.GovernanceRules().GetRules(ctx)
	if err != nil {
		t.Fatalf("GetRules() error = %v", err)
	}

	if rules == nil {
		t.Skip("No governance rules available")
	}

	decoded, err := client.GovernanceRules().GetDecodedRulesContainer(rules)
	if err != nil {
		t.Fatalf("GetDecodedRulesContainer() error = %v", err)
	}

	// Log address whitelisting rules
	t.Logf("Address Whitelisting Rules (%d):", len(decoded.AddressWhitelistingRules))
	for i, rule := range decoded.AddressWhitelistingRules {
		if i >= 3 {
			t.Logf("  ... and %d more", len(decoded.AddressWhitelistingRules)-i)
			break
		}
		t.Logf("  [%d] Currency: %s, Network: %s, ParallelThresholds: %d", i, rule.Currency, rule.Network, len(rule.ParallelThresholds))
	}

	// Log contract address whitelisting rules
	t.Logf("Contract Address Whitelisting Rules (%d):", len(decoded.ContractAddressWhitelistingRules))
	for i, rule := range decoded.ContractAddressWhitelistingRules {
		if i >= 3 {
			t.Logf("  ... and %d more", len(decoded.ContractAddressWhitelistingRules)-i)
			break
		}
		t.Logf("  [%d] Blockchain: %s, Network: %s, ParallelThresholds: %d", i, rule.Blockchain, rule.Network, len(rule.ParallelThresholds))
	}

	// Log transaction rules
	t.Logf("Transaction Rules (%d):", len(decoded.TransactionRules))
	for i, rule := range decoded.TransactionRules {
		if i >= 3 {
			t.Logf("  ... and %d more", len(decoded.TransactionRules)-i)
			break
		}
		t.Logf("  [%d] Key: %s", i, rule.Key)
	}
}

func TestIntegration_ClientSuperAdminKeysConfigured(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClientWithVerification(t)
	defer client.Close()

	// Verify client has SuperAdmin keys configured
	keys := client.SuperAdminKeys()
	if len(keys) == 0 {
		t.Fatal("Client should have SuperAdmin keys configured")
	}

	t.Logf("Client configured with %d SuperAdmin keys", len(keys))
	t.Logf("MinValidSignatures: %d", client.MinValidSignatures())

	// Verify keys match expected count
	if len(keys) != len(DefaultSuperAdminKeysPEM) {
		t.Errorf("Expected %d SuperAdmin keys, got %d", len(DefaultSuperAdminKeysPEM), len(keys))
	}

	// Verify min signatures matches config
	if client.MinValidSignatures() != DefaultMinValidSignatures {
		t.Errorf("Expected MinValidSignatures=%d, got %d", DefaultMinValidSignatures, client.MinValidSignatures())
	}
}
