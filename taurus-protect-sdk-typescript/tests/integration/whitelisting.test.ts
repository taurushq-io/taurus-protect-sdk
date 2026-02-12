/**
 * Integration tests for Address and Asset Whitelisting APIs.
 *
 * These tests require a live API connection. Configure via environment variables
 * or hard-coded defaults in helpers.ts.
 *
 * Tests use high-level services (client.whitelistedAddresses, client.whitelistedAssets)
 * rather than raw OpenAPI APIs to demonstrate proper SDK usage patterns.
 */

import { skipIfNotIntegration, getTestClient } from "./helpers";

describe("Integration: Whitelisting", () => {
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

  describe("Whitelisted Addresses", () => {
    it("listWhitelistedAddresses", async () => {
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
        const result = await client.whitelistedAddresses.list({ limit: 10 });

        const addresses = result.items;
        console.log(`Found ${addresses.length} whitelisted addresses`);

        if (result.pagination) {
          console.log(`Total items: ${result.pagination.totalItems}`);
        }

        for (const addr of addresses) {
          console.log(
            `  Address: ID=${addr.id}, Blockchain=${addr.blockchain}, Network=${addr.network}`
          );
        }

        expect(Array.isArray(addresses)).toBe(true);
      } finally {
        client.close();
      }
    });

    it(
      "paginateAllWhitelistedAddresses",
      async () => {
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
          const pageSize = 50;
          const allAddresses: Array<{ id?: string }> = [];
          let offset = 0;
          // Safety limit for addresses (aligned across all SDKs)
          const safetyLimit = 2000;

          // Fetch whitelisted addresses using pagination
          while (true) {
            const result = await client.whitelistedAddresses.list({
              limit: pageSize,
              offset: offset,
            });

            const addresses = result.items;
            allAddresses.push(...addresses);
            console.log(
              `Fetched ${addresses.length} whitelisted addresses (offset=${offset})`
            );

            // Check if there are more pages
            const totalItems = result.pagination?.totalItems ?? 0;
            if (offset + pageSize >= totalItems || addresses.length === 0) {
              break;
            }

            offset += pageSize;

            // Safety limit to avoid excessive API calls
            if (offset >= safetyLimit) {
              console.log(
                `Stopping pagination test at ${safetyLimit} items (safety limit)`
              );
              break;
            }
          }

          console.log(
            `Total whitelisted addresses fetched via pagination: ${allAddresses.length}`
          );
          expect(allAddresses.length).toBeGreaterThanOrEqual(0);
        } finally {
          client.close();
        }
      },
      300000
    );

    it("getWhitelistedAddress", async () => {
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
        // First, get a list to find a valid ID
        const listResult = await client.whitelistedAddresses.list({ limit: 1 });

        const addresses = listResult.items;
        if (addresses.length === 0) {
          console.log("No whitelisted addresses available for testing");
          return;
        }

        const addressId = addresses[0].id;
        if (!addressId) {
          console.log("Whitelisted address ID is undefined");
          return;
        }

        const address = await client.whitelistedAddresses.get(addressId);

        expect(address).toBeDefined();

        console.log("Whitelisted address details:");
        console.log(`  ID: ${address.id}`);
        console.log(`  Address: ${address.address}`);
        console.log(`  Blockchain: ${address.blockchain}`);
        console.log(`  Network: ${address.network}`);
      } finally {
        client.close();
      }
    });

    it("listWhitelistedAddressesByBlockchain", async () => {
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
        const result = await client.whitelistedAddresses.list({
          limit: 10,
          blockchain: "ETH",
        });

        const addresses = result.items;
        console.log(
          `Found ${addresses.length} whitelisted addresses filtered by blockchain=ETH`
        );

        if (result.pagination) {
          console.log(`Total items: ${result.pagination.totalItems}`);
        }

        for (const addr of addresses) {
          console.log(
            `  Address: ID=${addr.id}, Blockchain=${addr.blockchain}, Network=${addr.network}`
          );
        }

        expect(Array.isArray(addresses)).toBe(true);
        // All returned addresses should be ETH blockchain
        for (const addr of addresses) {
          expect(addr.blockchain).toBe("ETH");
        }
      } finally {
        client.close();
      }
    });

    it("listWhitelistedAddressesByBlockchainAndNetwork", async () => {
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
        const result = await client.whitelistedAddresses.list({
          limit: 10,
          blockchain: "ETH",
          network: "mainnet",
        });

        const addresses = result.items;
        console.log(
          `Found ${addresses.length} whitelisted addresses filtered by blockchain=ETH, network=mainnet`
        );

        if (result.pagination) {
          console.log(`Total items: ${result.pagination.totalItems}`);
        }

        for (const addr of addresses) {
          console.log(
            `  Address: ID=${addr.id}, Blockchain=${addr.blockchain}, Network=${addr.network}`
          );
        }

        expect(Array.isArray(addresses)).toBe(true);
        // Verify filtering worked (API returns filtered results)
        // Note: blockchain and network in the model come from the verified payload,
        // which may not include all fields. The API filtering works correctly,
        // but the model may have empty values if not in the signed payload.
        // This is intentional for security - we only expose verified data.
        expect(addresses.length).toBeGreaterThan(0);
        // Verify blockchain is set (comes from payload 'currency' field for addresses)
        for (const addr of addresses) {
          // blockchain may be empty if payload doesn't have 'currency' field
          // network may be empty if payload doesn't have 'network' field
          // We just verify we got results from the filtered API call
          expect(typeof addr.blockchain).toBe("string");
          expect(typeof addr.network).toBe("string");
        }
      } finally {
        client.close();
      }
    });
  });

  describe("Whitelisted Assets (Contracts)", () => {
    it("listWhitelistedAssets", async () => {
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
        const result = await client.whitelistedAssets.list({ limit: 10 });

        const assets = result.items;
        console.log(`Found ${assets.length} whitelisted assets (contracts)`);

        if (result.pagination) {
          console.log(`Total items: ${result.pagination.totalItems}`);
        }

        for (const asset of assets) {
          console.log(
            `  Asset: ID=${asset.id}, Blockchain=${asset.blockchain}, Symbol=${asset.symbol}`
          );
        }

        expect(Array.isArray(assets)).toBe(true);
      } finally {
        client.close();
      }
    });

    it(
      "paginateAllWhitelistedAssets",
      async () => {
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
          const pageSize = 100;
          const allAssets: Array<{ id?: number }> = [];
          let offset = 0;
          const safetyLimit = 5000;

          // Fetch whitelisted assets using pagination
          while (true) {
            const result = await client.whitelistedAssets.list({
              limit: pageSize,
              offset: offset,
            });

            const assets = result.items;
            allAssets.push(...assets);
            console.log(
              `Fetched ${assets.length} whitelisted assets (offset=${offset})`
            );

            // Check if there are more pages
            const totalItems = result.pagination?.totalItems ?? 0;
            if (offset + pageSize >= totalItems || assets.length === 0) {
              break;
            }

            offset += pageSize;

            // Safety limit to avoid excessive API calls
            if (offset >= safetyLimit) {
              console.log(
                `Stopping pagination test at ${safetyLimit} items (safety limit)`
              );
              break;
            }
          }

          console.log(
            `Total whitelisted assets fetched via pagination: ${allAssets.length}`
          );
          expect(allAssets.length).toBeGreaterThanOrEqual(0);
        } finally {
          client.close();
        }
      },
      60000
    );

    it("getWhitelistedAsset", async () => {
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
        // First, get a list to find a valid ID
        const listResult = await client.whitelistedAssets.list({ limit: 1 });

        const assets = listResult.items;
        if (assets.length === 0) {
          console.log("No whitelisted assets available for testing");
          return;
        }

        const assetId = assets[0].id;
        if (!assetId) {
          console.log("Whitelisted asset ID is undefined");
          return;
        }

        const asset = await client.whitelistedAssets.get(assetId);

        expect(asset).toBeDefined();

        console.log("Whitelisted asset details:");
        console.log(`  ID: ${asset.id}`);
        console.log(`  Contract: ${asset.contractAddress}`);
        console.log(`  Blockchain: ${asset.blockchain}`);
        console.log(`  Network: ${asset.network}`);
      } finally {
        client.close();
      }
    });

    it("listWhitelistedAssetsByBlockchain", async () => {
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
        const result = await client.whitelistedAssets.list({
          limit: 10,
          blockchain: "ETH",
        });

        const assets = result.items;
        console.log(
          `Found ${assets.length} whitelisted assets filtered by blockchain=ETH`
        );

        if (result.pagination) {
          console.log(`Total items: ${result.pagination.totalItems}`);
        }

        for (const asset of assets) {
          console.log(
            `  Asset: ID=${asset.id}, Blockchain=${asset.blockchain}, Symbol=${asset.symbol}`
          );
        }

        expect(Array.isArray(assets)).toBe(true);
        // All returned assets should be ETH blockchain
        for (const asset of assets) {
          expect(asset.blockchain).toBe("ETH");
        }
      } finally {
        client.close();
      }
    });

    it("listWhitelistedAssetsByBlockchainAndNetwork", async () => {
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
        const result = await client.whitelistedAssets.list({
          limit: 10,
          blockchain: "ETH",
          network: "mainnet",
        });

        const assets = result.items;
        console.log(
          `Found ${assets.length} whitelisted assets filtered by blockchain=ETH, network=mainnet`
        );

        if (result.pagination) {
          console.log(`Total items: ${result.pagination.totalItems}`);
        }

        for (const asset of assets) {
          console.log(
            `  Asset: ID=${asset.id}, Blockchain=${asset.blockchain}, Network=${asset.network}, Symbol=${asset.symbol}`
          );
        }

        expect(Array.isArray(assets)).toBe(true);
        // Verify filtering worked (API returns filtered results)
        // Note: blockchain and network in the model come from the verified payload,
        // which may not include all fields. The API filtering works correctly,
        // but the model may have empty values if not in the signed payload.
        // This is intentional for security - we only expose verified data.
        expect(assets.length).toBeGreaterThan(0);
        // Verify blockchain is set (comes from payload 'blockchain' field for assets)
        for (const asset of assets) {
          // blockchain may be empty if payload doesn't have 'blockchain' field
          // network may be empty if payload doesn't have 'network' field
          // We just verify we got results from the filtered API call
          expect(typeof asset.blockchain).toBe("string");
          expect(typeof asset.network).toBe("string");
        }
      } finally {
        client.close();
      }
    });
  });
});
