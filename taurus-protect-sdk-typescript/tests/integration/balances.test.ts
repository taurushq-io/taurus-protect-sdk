/**
 * Integration tests for Balances API.
 *
 * These tests require a live API connection. Configure via environment variables
 * or hard-coded defaults in helpers.ts.
 *
 * Tests use the high-level BalanceService (client.balances) rather than
 * the raw OpenAPI API to demonstrate proper SDK usage patterns.
 */

import { skipIfNotIntegration, getTestClient } from "./helpers";

describe("Integration: Balances", () => {
  beforeEach(() => {
    try {
      skipIfNotIntegration();
    } catch (e) {
      if (e instanceof Error && e.message === "SKIP_INTEGRATION") {
        return;
      }
      throw e;
    }
  });

  it("should list balances", async () => {
    try {
      skipIfNotIntegration();
    } catch (e) {
      if (e instanceof Error && e.message === "SKIP_INTEGRATION") {
        return;
      }
      throw e;
    }

    const client = getTestClient();
    try {
      const balances = await client.balances.list({ limit: 10 });

      console.log(`Found ${balances.length} balances`);

      for (const assetBalance of balances) {
        console.log(
          `Balance: Currency=${assetBalance.currency}, ` +
            `Balance=${assetBalance.balance}, ` +
            `FiatValue=${assetBalance.fiatValue}`
        );
      }

      expect(Array.isArray(balances)).toBe(true);
    } finally {
      client.close();
    }
  });

  it("should list balances with currency filter", async () => {
    try {
      skipIfNotIntegration();
    } catch (e) {
      if (e instanceof Error && e.message === "SKIP_INTEGRATION") {
        return;
      }
      throw e;
    }

    const client = getTestClient();
    try {
      // First, get all balances to find a valid currency
      const allBalances = await client.balances.list({ limit: 1 });

      if (allBalances.length === 0) {
        console.log("No balances available for testing currency filter");
        return;
      }

      const currencyToFilter = allBalances[0].currency;
      if (!currencyToFilter) {
        console.log("Balance has no currency for filtering");
        return;
      }

      console.log(`Filtering balances by currency: ${currencyToFilter}`);

      const balances = await client.balances.list({
        currency: currencyToFilter,
        limit: 10,
      });

      console.log(
        `Found ${balances.length} balances for currency ${currencyToFilter}`
      );

      for (const assetBalance of balances) {
        console.log(
          `Balance: Currency=${assetBalance.currency}, ` +
            `Balance=${assetBalance.balance}, ` +
            `FiatValue=${assetBalance.fiatValue}`
        );
      }

      expect(Array.isArray(balances)).toBe(true);
    } finally {
      client.close();
    }
  });
});
