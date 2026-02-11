/**
 * Integration tests for miscellaneous APIs (Tags, Statistics) and client lifecycle.
 *
 * These tests require a live API connection. Configure via environment variables
 * or hard-coded defaults in helpers.ts.
 *
 * Tests use high-level services (client.tags, client.statistics) rather than
 * raw OpenAPI APIs to demonstrate proper SDK usage patterns.
 */

import { skipIfNotIntegration, getTestClient } from "./helpers";
import { ProtectClient } from "../../src/client";

describe("Integration: Tags", () => {
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

  it("should list tags", async () => {
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
      const tags = await client.tags.list();

      console.log(`Found ${tags.length} tags`);

      for (const tag of tags.slice(0, 10)) {
        console.log(`  Tag: ${tag.name}, ID: ${tag.id}`);
      }

      if (tags.length > 10) {
        console.log(`  ... and ${tags.length - 10} more`);
      }

      expect(Array.isArray(tags)).toBe(true);
    } finally {
      client.close();
    }
  });
});

describe("Integration: Statistics", () => {
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

  it("should get portfolio statistics", async () => {
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
      const stats = await client.statistics.getPortfolioStatistics();

      console.log("Portfolio statistics response received");
      expect(stats).toBeDefined();

      console.log(`  Total balance: ${stats.totalBalance}`);
      console.log(
        `  Total balance (base currency): ${stats.totalBalanceBaseCurrency}`
      );
      if (stats.walletsCount !== undefined) {
        console.log(`  Wallets count: ${stats.walletsCount}`);
      }
      if (stats.addressesCount !== undefined) {
        console.log(`  Addresses count: ${stats.addressesCount}`);
      }
    } finally {
      client.close();
    }
  });
});

describe("Integration: Client Lifecycle", () => {
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

  it("should lazily initialize services", async () => {
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
      // Services should be lazily initialized - access them and verify they work
      // First access to wallets service should create the instance
      const walletsService1 = client.wallets;
      expect(walletsService1).toBeDefined();

      // Second access should return the same instance
      const walletsService2 = client.wallets;
      expect(walletsService2).toBe(walletsService1);

      console.log("Lazy initialization verified: same instance returned");

      // Verify the service is functional by making a request
      const result = await client.wallets.list({ limit: 1 });
      expect(result).toBeDefined();
      console.log("Service is functional after lazy initialization");
    } finally {
      client.close();
    }
  });

  it("should handle close() idempotently", async () => {
    try {
      skipIfNotIntegration();
    } catch (e) {
      if (e instanceof Error && e.message === "SKIP_INTEGRATION") {
        return;
      }
      throw e;
    }

    const client = getTestClient();

    // Make a request to ensure client is working
    const result = await client.wallets.list({ limit: 1 });
    expect(result).toBeDefined();
    console.log("Client is working before close");

    // Close should work without error
    client.close();
    console.log("First close() succeeded");

    // Calling close() again should not throw
    expect(() => client.close()).not.toThrow();
    console.log("Second close() succeeded (idempotent)");

    // Calling close() a third time should also not throw
    expect(() => client.close()).not.toThrow();
    console.log("Third close() succeeded (idempotent)");
  });

  it("should create client with factory method", () => {
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
      // Verify client is an instance of ProtectClient
      expect(client).toBeInstanceOf(ProtectClient);
      console.log("Client created via factory method is ProtectClient instance");

      // Verify client has expected high-level services
      expect(client.wallets).toBeDefined();
      expect(client.health).toBeDefined();
      // Note: blockchainApi is available, but BlockchainService is not yet implemented
      expect(client.blockchainApi).toBeDefined();
      expect(client.tags).toBeDefined();
      expect(client.statistics).toBeDefined();
      console.log("Client has all expected high-level service properties");
    } finally {
      client.close();
    }
  });
});
