/**
 * Integration tests for Currencies API.
 *
 * These tests require a live API connection. Configure via environment variables
 * or hard-coded defaults in helpers.ts.
 *
 * Tests use the high-level CurrencyService (client.currencies) rather than
 * the raw OpenAPI API to demonstrate proper SDK usage patterns.
 */

import { skipIfNotIntegration, getTestClient } from "./helpers";

describe("Integration: Currencies", () => {
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

  it("should list currencies", async () => {
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
      const currencies = await client.currencies.list();

      console.log(`Found ${currencies.length} currencies`);

      // Log first few currencies for visibility
      for (const currency of currencies.slice(0, 5)) {
        console.log(
          `Currency: ID=${currency.id}, Symbol=${currency.symbol}, Blockchain=${currency.blockchain}`
        );
      }

      if (currencies.length > 5) {
        console.log(`... and ${currencies.length - 5} more`);
      }

      expect(Array.isArray(currencies)).toBe(true);
    } finally {
      client.close();
    }
  });

  it("should get currency by blockchain and network", async () => {
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
      // First, get a list to find a valid currency
      const currencies = await client.currencies.list();

      if (currencies.length === 0) {
        console.log("No currencies available for testing");
        return;
      }

      // Find a currency with blockchain and network (required for getByBlockchain)
      const currencyWithDetails = currencies.find(
        (c) => c.blockchain && c.network
      );
      if (!currencyWithDetails) {
        console.log("No currency with blockchain and network found");
        return;
      }

      const currency = await client.currencies.getByBlockchain({
        blockchain: currencyWithDetails.blockchain!,
        network: currencyWithDetails.network!,
        contractAddress: currencyWithDetails.contractAddress,
      });

      expect(currency).toBeDefined();

      console.log("Currency details:");
      console.log(`  ID: ${currency.id}`);
      console.log(`  Symbol: ${currency.symbol}`);
      console.log(`  Name: ${currency.name}`);
      console.log(`  Blockchain: ${currency.blockchain}`);
      console.log(`  Network: ${currency.network}`);
      console.log(`  Decimals: ${currency.decimals}`);
    } finally {
      client.close();
    }
  });

  it("should get base currency", async () => {
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
      const baseCurrency = await client.currencies.getBaseCurrency();

      expect(baseCurrency).toBeDefined();

      console.log(`Base currency: ${baseCurrency.symbol}`);
    } catch (e) {
      // Skip if base currency is not configured in the test environment
      if (e instanceof Error && e.message === "Base currency not configured") {
        console.log("Skipping: Base currency not configured in test environment");
        return;
      }
      throw e;
    } finally {
      client.close();
    }
  });
});
