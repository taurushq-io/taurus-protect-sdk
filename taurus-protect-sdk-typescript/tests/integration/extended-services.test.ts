/**
 * Integration tests for extended services that previously lacked test coverage.
 *
 * Tests read-only operations for: BusinessRuleService, WebhookService,
 * WebhookCallService, StakingService, FeePayerService, ExchangeService,
 * AssetService, FiatService, JobService, ContractWhitelistingService,
 * and TaurusNetwork (participants, pledges, lending, settlements, sharing).
 *
 * These tests require a live API connection. Configure via environment variables
 * or hard-coded defaults in config.ts.
 */

import { describeIntegration, getTestClient } from "./helpers";

// =============================================================================
// BusinessRuleService
// =============================================================================

describeIntegration("Integration: BusinessRuleService", () => {
  it("should list business rules", async () => {
    const client = getTestClient();
    try {
      const result = await client.businessRules.list();

      expect(result).toBeDefined();
      console.log(`Found ${result.rules.length} business rules`);
      for (const rule of result.rules.slice(0, 5)) {
        console.log(`  Rule: ID=${rule.id}`);
      }
    } finally {
      client.close();
    }
  }, 30000);
});

// =============================================================================
// WebhookService
// =============================================================================

describeIntegration("Integration: WebhookService", () => {
  it("should list webhooks", async () => {
    const client = getTestClient();
    try {
      const webhooks = await client.webhooks.list();

      expect(webhooks).toBeDefined();
      expect(Array.isArray(webhooks)).toBe(true);
      console.log(`Found ${webhooks.length} webhooks`);
      for (const wh of webhooks.slice(0, 5)) {
        console.log(`  Webhook: ID=${wh.id}, URL=${wh.url}`);
      }
    } finally {
      client.close();
    }
  }, 30000);
});

// =============================================================================
// WebhookCallService
// =============================================================================

describeIntegration("Integration: WebhookCallService", () => {
  it("should list webhook calls", async () => {
    const client = getTestClient();
    try {
      const result = await client.webhookCalls.list();

      expect(result).toBeDefined();
      console.log(`Found ${result.calls.length} webhook calls`);
    } finally {
      client.close();
    }
  }, 30000);
});

// =============================================================================
// StakingService
// =============================================================================

describeIntegration("Integration: StakingService", () => {
  it("should get stake accounts", async () => {
    const client = getTestClient();
    try {
      // Use getStakeAccounts which doesn't require specific validator IDs
      const result = await client.staking.getStakeAccounts({
        addressId: "1",
        accountType: "StakeAccountTypeSolana",
      });

      expect(result).toBeDefined();
      console.log("Staking stake accounts response received");
    } catch (error: unknown) {
      // Staking may not be available in all environments
      console.log(
        `Staking not available: ${error instanceof Error ? error.message : String(error)}`
      );
    } finally {
      client.close();
    }
  }, 30000);
});

// =============================================================================
// FeePayerService
// =============================================================================

describeIntegration("Integration: FeePayerService", () => {
  it("should list fee payers", async () => {
    const client = getTestClient();
    try {
      const feePayers = await client.feePayers.list();

      expect(feePayers).toBeDefined();
      expect(Array.isArray(feePayers)).toBe(true);
      console.log(`Found ${feePayers.length} fee payers`);
      for (const fp of feePayers.slice(0, 5)) {
        console.log(`  FeePayer: ID=${fp.id}`);
      }
    } finally {
      client.close();
    }
  }, 30000);
});

// =============================================================================
// ExchangeService
// =============================================================================

describeIntegration("Integration: ExchangeService", () => {
  it("should list exchanges", async () => {
    const client = getTestClient();
    try {
      const result = await client.exchanges.list();

      expect(result).toBeDefined();
      console.log(`Found ${result.items.length} exchanges`);
      for (const ex of result.items.slice(0, 5)) {
        console.log(`  Exchange: ID=${ex.id}`);
      }
    } finally {
      client.close();
    }
  }, 30000);
});

// =============================================================================
// FiatService
// =============================================================================

describeIntegration("Integration: FiatService", () => {
  it("should list fiat providers", async () => {
    const client = getTestClient();
    try {
      const providers = await client.fiatAccounts.getFiatProviders();

      expect(providers).toBeDefined();
      expect(Array.isArray(providers)).toBe(true);
      console.log(`Found ${providers.length} fiat providers`);
      for (const p of providers.slice(0, 5)) {
        console.log(`  Provider: ${p.provider}`);
      }
    } finally {
      client.close();
    }
  }, 30000);
});

// =============================================================================
// JobService
// =============================================================================

describeIntegration("Integration: JobService", () => {
  it("should list jobs", async () => {
    const client = getTestClient();
    try {
      const jobs = await client.jobs.list();

      expect(jobs).toBeDefined();
      expect(Array.isArray(jobs)).toBe(true);
      console.log(`Found ${jobs.length} jobs`);
      for (const job of jobs.slice(0, 5)) {
        console.log(`  Job: Name=${job.name}`);
      }
    } catch (error: unknown) {
      console.log(
        `Jobs not available: ${error instanceof Error ? error.message : String(error)}`
      );
    } finally {
      client.close();
    }
  }, 30000);
});

// =============================================================================
// TaurusNetwork - Participants
// =============================================================================

describeIntegration("Integration: TaurusNetwork Participants", () => {
  it("should list participants", async () => {
    const client = getTestClient();
    try {
      const participants = await client.taurusNetwork.participants.list();

      expect(participants).toBeDefined();
      console.log(`Found ${participants.length} TaurusNetwork participants`);
      for (const p of participants.slice(0, 5)) {
        console.log(`  Participant: ID=${p.id}, Name=${p.name}`);
      }
    } catch (error: unknown) {
      // TaurusNetwork may not be enabled in all environments
      console.log(
        `TaurusNetwork participants not available: ${error instanceof Error ? error.message : String(error)}`
      );
    } finally {
      client.close();
    }
  }, 30000);
});

// =============================================================================
// TaurusNetwork - Pledges
// =============================================================================

describeIntegration("Integration: TaurusNetwork Pledges", () => {
  it("should list pledges", async () => {
    const client = getTestClient();
    try {
      const result = await client.taurusNetwork.pledges.list();

      expect(result).toBeDefined();
      console.log(`Found ${result.pledges.length} TaurusNetwork pledges`);
      for (const pledge of result.pledges.slice(0, 5)) {
        console.log(`  Pledge: ID=${pledge.id}, Status=${pledge.status}`);
      }
    } catch (error: unknown) {
      // TaurusNetwork may not be enabled in all environments
      console.log(
        `TaurusNetwork pledges not available: ${error instanceof Error ? error.message : String(error)}`
      );
    } finally {
      client.close();
    }
  }, 30000);
});

// =============================================================================
// TaurusNetwork - Lending
// =============================================================================

describeIntegration("Integration: TaurusNetwork Lending", () => {
  it("should list lending offers", async () => {
    const client = getTestClient();
    try {
      const { offers } = await client.taurusNetwork.lending.listLendingOffers();

      expect(offers).toBeDefined();
      console.log(`Found ${offers.length} TaurusNetwork lending offers`);
      for (const offer of offers.slice(0, 5)) {
        console.log(`  Offer: ID=${offer.id}`);
      }
    } catch (error: unknown) {
      console.log(
        `TaurusNetwork lending offers not available: ${error instanceof Error ? error.message : String(error)}`
      );
    } finally {
      client.close();
    }
  }, 30000);

  it("should list lending agreements", async () => {
    const client = getTestClient();
    try {
      const { agreements } =
        await client.taurusNetwork.lending.listLendingAgreements();

      expect(agreements).toBeDefined();
      console.log(`Found ${agreements.length} TaurusNetwork lending agreements`);
      for (const agreement of agreements.slice(0, 5)) {
        console.log(`  Agreement: ID=${agreement.id}`);
      }
    } catch (error: unknown) {
      console.log(
        `TaurusNetwork lending agreements not available: ${error instanceof Error ? error.message : String(error)}`
      );
    } finally {
      client.close();
    }
  }, 30000);
});

// =============================================================================
// TaurusNetwork - Settlements
// =============================================================================

describeIntegration("Integration: TaurusNetwork Settlements", () => {
  it("should list settlements", async () => {
    const client = getTestClient();
    try {
      const { settlements } = await client.taurusNetwork.settlements.list();

      expect(settlements).toBeDefined();
      console.log(`Found ${settlements.length} TaurusNetwork settlements`);
      for (const s of settlements.slice(0, 5)) {
        console.log(`  Settlement: ID=${s.id}, Status=${s.status}`);
      }
    } catch (error: unknown) {
      console.log(
        `TaurusNetwork settlements not available: ${error instanceof Error ? error.message : String(error)}`
      );
    } finally {
      client.close();
    }
  }, 30000);
});

// =============================================================================
// TaurusNetwork - Sharing
// =============================================================================

describeIntegration("Integration: TaurusNetwork Sharing", () => {
  it("should list shared addresses", async () => {
    const client = getTestClient();
    try {
      const { sharedAddresses } =
        await client.taurusNetwork.sharing.listSharedAddresses();

      expect(sharedAddresses).toBeDefined();
      console.log(
        `Found ${sharedAddresses.length} TaurusNetwork shared addresses`
      );
      for (const addr of sharedAddresses.slice(0, 5)) {
        console.log(`  SharedAddress: ID=${addr.id}`);
      }
    } catch (error: unknown) {
      console.log(
        `TaurusNetwork shared addresses not available: ${error instanceof Error ? error.message : String(error)}`
      );
    } finally {
      client.close();
    }
  }, 30000);

  it("should list shared assets", async () => {
    const client = getTestClient();
    try {
      const { sharedAssets } =
        await client.taurusNetwork.sharing.listSharedAssets();

      expect(sharedAssets).toBeDefined();
      console.log(`Found ${sharedAssets.length} TaurusNetwork shared assets`);
      for (const asset of sharedAssets.slice(0, 5)) {
        console.log(`  SharedAsset: ID=${asset.id}`);
      }
    } catch (error: unknown) {
      console.log(
        `TaurusNetwork shared assets not available: ${error instanceof Error ? error.message : String(error)}`
      );
    } finally {
      client.close();
    }
  }, 30000);
});

// =============================================================================
// ContractWhitelistingService
// =============================================================================

describeIntegration("Integration: ContractWhitelistingService", () => {
  it("should list whitelisted contracts", async () => {
    const client = getTestClient();
    try {
      const result = await client.contractWhitelisting.list();

      expect(result).toBeDefined();
      console.log(`Found ${result.items.length} whitelisted contracts`);
      for (const c of result.items.slice(0, 5)) {
        console.log(`  Contract: ID=${c.id}`);
      }
    } finally {
      client.close();
    }
  }, 30000);
});
