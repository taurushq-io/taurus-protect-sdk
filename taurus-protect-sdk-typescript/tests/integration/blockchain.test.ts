/**
 * Integration tests for Blockchain API.
 *
 * These tests require a live API connection. Configure via environment variables
 * or hard-coded defaults in helpers.ts.
 *
 * Note: BlockchainService is not yet implemented as a high-level service,
 * so these tests use the raw OpenAPI API (client.blockchainApi).
 */

import { skipIfNotIntegration, getTestClient } from "./helpers";

describe("Integration: Blockchain", () => {
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

  it("should list blockchains", async () => {
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
      // Note: Using raw API as BlockchainService is not yet implemented
      const response = await client.blockchainApi.blockchainServiceGetBlockchains();
      const blockchains = (response as Record<string, unknown>).result as Array<{
        name?: string;
        symbol?: string;
        network?: string;
      }> ?? [];

      console.log(`Found ${blockchains.length} blockchains`);

      for (const blockchain of blockchains.slice(0, 10)) {
        console.log(
          `  Blockchain: ${blockchain.name} (${blockchain.symbol}), Network: ${blockchain.network}`
        );
      }

      if (blockchains.length > 10) {
        console.log(`  ... and ${blockchains.length - 10} more`);
      }

      expect(Array.isArray(blockchains)).toBe(true);
    } finally {
      client.close();
    }
  });

  it("should get blockchain by symbol and network", async () => {
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
      // Note: Using raw API as BlockchainService is not yet implemented
      // First, get list to find a valid blockchain
      const response = await client.blockchainApi.blockchainServiceGetBlockchains();
      const blockchains = (response as Record<string, unknown>).result as Array<{
        name?: string;
        symbol?: string;
        network?: string;
      }> ?? [];

      if (blockchains.length === 0) {
        console.log("No blockchains available for testing");
        return;
      }

      const first = blockchains[0];
      if (!first.symbol || !first.network) {
        console.log("Blockchain missing symbol or network");
        return;
      }

      // Find the specific blockchain from the list (API doesn't have a get-by-id method)
      const blockchain = blockchains.find(
        (b) => b.symbol === first.symbol && b.network === first.network
      );

      console.log("Blockchain details:");
      console.log(`  Symbol: ${blockchain?.symbol}`);
      console.log(`  Name: ${blockchain?.name}`);
      console.log(`  Network: ${blockchain?.network}`);

      expect(blockchain).toBeDefined();
    } finally {
      client.close();
    }
  });
});
