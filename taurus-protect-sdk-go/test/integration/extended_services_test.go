package integration

import (
	"context"
	"testing"
)

// =============================================================================
// BusinessRuleService
// =============================================================================

func TestIntegration_ListBusinessRules(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	result, err := client.BusinessRules().ListBusinessRules(ctx, nil)
	if err != nil {
		t.Fatalf("ListBusinessRules() error = %v", err)
	}

	t.Logf("Found %d business rules", len(result.BusinessRules))
	for _, rule := range result.BusinessRules {
		t.Logf("  Rule: ID=%s", rule.ID)
	}
}

// =============================================================================
// WebhookService
// =============================================================================

func TestIntegration_ListWebhooks(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	result, err := client.Webhooks().ListWebhooks(ctx, nil)
	if err != nil {
		t.Fatalf("ListWebhooks() error = %v", err)
	}

	t.Logf("Found %d webhooks", len(result.Webhooks))
	for _, wh := range result.Webhooks {
		t.Logf("  Webhook: ID=%s, URL=%s", wh.ID, wh.URL)
	}
}

// =============================================================================
// WebhookCallService
// =============================================================================

func TestIntegration_ListWebhookCalls(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	result, err := client.WebhookCalls().ListWebhookCalls(ctx, nil)
	if err != nil {
		t.Fatalf("ListWebhookCalls() error = %v", err)
	}

	t.Logf("Found %d webhook calls", len(result.WebhookCalls))
}

// =============================================================================
// FeePayerService
// =============================================================================

func TestIntegration_ListFeePayers(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	result, err := client.FeePayers().ListFeePayers(ctx, nil)
	if err != nil {
		t.Fatalf("ListFeePayers() error = %v", err)
	}

	t.Logf("Found %d fee payers", len(result.FeePayers))
	for _, fp := range result.FeePayers {
		t.Logf("  FeePayer: ID=%s", fp.ID)
	}
}

// =============================================================================
// ExchangeService
// =============================================================================

func TestIntegration_ListExchangeCounterparties(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	result, err := client.Exchanges().ListExchangeCounterparties(ctx)
	if err != nil {
		t.Fatalf("ListExchangeCounterparties() error = %v", err)
	}

	t.Logf("Found %d exchange counterparties", len(result.Exchanges))
	for _, cp := range result.Exchanges {
		t.Logf("  Exchange: %s", cp.Name)
	}
}

// =============================================================================
// FiatService
// =============================================================================

func TestIntegration_ListFiatProviders(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	result, err := client.Fiat().ListFiatProviders(ctx)
	if err != nil {
		t.Fatalf("ListFiatProviders() error = %v", err)
	}

	t.Logf("Found %d fiat providers", len(result.FiatProviders))
	for _, p := range result.FiatProviders {
		t.Logf("  Provider: %s", p.Provider)
	}
}

// =============================================================================
// JobService
// =============================================================================

func TestIntegration_ListJobs(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	jobs, err := client.Jobs().ListJobs(ctx)
	if err != nil {
		// Jobs API may require elevated permissions not available to the test API key
		t.Logf("ListJobs() not available: %v", err)
		t.Skip("Jobs API not available (may require elevated permissions)")
	}

	t.Logf("Found %d jobs", len(jobs))
	for _, job := range jobs {
		t.Logf("  Job: Name=%s", job.Name)
	}
}

// =============================================================================
// TaurusNetwork - Participants
// =============================================================================

func TestIntegration_ListTaurusNetworkParticipants(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	result, err := client.TaurusNetwork().Participants().ListParticipants(ctx, nil)
	if err != nil {
		// TaurusNetwork may not be enabled in all environments
		t.Logf("TaurusNetwork participants not available: %v", err)
		t.Skip("TaurusNetwork not available")
	}

	t.Logf("Found %d TaurusNetwork participants", len(result.Participants))
	for _, p := range result.Participants {
		t.Logf("  Participant: ID=%s, Name=%s", p.ID, p.Name)
	}
}

// =============================================================================
// TaurusNetwork - Pledges
// =============================================================================

func TestIntegration_ListTaurusNetworkPledges(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	pledges, _, err := client.TaurusNetwork().Pledges().ListPledges(ctx, nil)
	if err != nil {
		// TaurusNetwork may not be enabled in all environments
		t.Logf("TaurusNetwork pledges not available: %v", err)
		t.Skip("TaurusNetwork not available")
	}

	t.Logf("Found %d TaurusNetwork pledges", len(pledges))
	for _, pledge := range pledges {
		t.Logf("  Pledge: ID=%s, Status=%s", pledge.ID, pledge.Status)
	}
}

// =============================================================================
// TaurusNetwork - Lending
// =============================================================================

func TestIntegration_ListTaurusNetworkLendingOffers(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	result, err := client.TaurusNetwork().Lending().ListLendingOffers(ctx, nil)
	if err != nil {
		t.Logf("TaurusNetwork lending offers not available: %v", err)
		t.Skip("TaurusNetwork not available")
	}

	t.Logf("Found %d TaurusNetwork lending offers", len(result.LendingOffers))
}

func TestIntegration_ListTaurusNetworkLendingAgreements(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	result, err := client.TaurusNetwork().Lending().ListLendingAgreements(ctx, nil)
	if err != nil {
		t.Logf("TaurusNetwork lending agreements not available: %v", err)
		t.Skip("TaurusNetwork not available")
	}

	t.Logf("Found %d TaurusNetwork lending agreements", len(result.LendingAgreements))
}

// =============================================================================
// TaurusNetwork - Settlements
// =============================================================================

func TestIntegration_ListTaurusNetworkSettlements(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	result, err := client.TaurusNetwork().Settlements().ListSettlements(ctx, nil)
	if err != nil {
		t.Logf("TaurusNetwork settlements not available: %v", err)
		t.Skip("TaurusNetwork not available")
	}

	t.Logf("Found %d TaurusNetwork settlements", len(result.Settlements))
}

// =============================================================================
// TaurusNetwork - Sharing
// =============================================================================

func TestIntegration_ListTaurusNetworkSharedAddresses(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	result, err := client.TaurusNetwork().Sharing().ListSharedAddresses(ctx, nil)
	if err != nil {
		t.Logf("TaurusNetwork shared addresses not available: %v", err)
		t.Skip("TaurusNetwork not available")
	}

	t.Logf("Found %d TaurusNetwork shared addresses", len(result.SharedAddresses))
}

func TestIntegration_ListTaurusNetworkSharedAssets(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	result, err := client.TaurusNetwork().Sharing().ListSharedAssets(ctx, nil)
	if err != nil {
		t.Logf("TaurusNetwork shared assets not available: %v", err)
		t.Skip("TaurusNetwork not available")
	}

	t.Logf("Found %d TaurusNetwork shared assets", len(result.SharedAssets))
}

// =============================================================================
// ContractWhitelistingService
// =============================================================================

func TestIntegration_ListWhitelistedContracts(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	result, err := client.WhitelistedContracts().ListWhitelistedContracts(ctx, nil)
	if err != nil {
		t.Fatalf("ListWhitelistedContracts() error = %v", err)
	}

	t.Logf("Found %d whitelisted contracts", len(result.Contracts))
}
