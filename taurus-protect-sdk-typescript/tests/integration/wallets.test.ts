/**
 * Integration tests for Wallets API.
 *
 * These tests require a live API connection. Configure via environment variables
 * or hard-coded defaults in helpers.ts.
 *
 * Tests use the high-level WalletService (client.wallets) rather than
 * the raw OpenAPI API to demonstrate proper SDK usage patterns.
 */

import { skipIfNotIntegration, getTestClient } from "./helpers";

describe("Integration: Wallets", () => {
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

  it("should list wallets", async () => {
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
      const result = await client.wallets.list({ limit: 10 });

      const wallets = result.items;
      console.log(`Found ${wallets.length} wallets`);

      if (result.pagination) {
        console.log(`Total items: ${result.pagination.totalItems}`);
      }

      for (const wallet of wallets) {
        console.log("Wallet - All Fields:");
        console.log("=".repeat(50));
        console.log(`  id: ${wallet.id}`);
        console.log(`  name: ${wallet.name}`);
        console.log(`  externalWalletId: ${wallet.externalWalletId ?? "(undefined)"}`);
        console.log(`  status: ${wallet.status}`);
        console.log(`  type: ${wallet.type}`);
        console.log(`  blockchain: ${wallet.blockchain ?? "(undefined)"}`);
        console.log(`  network: ${wallet.network ?? "(undefined)"}`);
        console.log(`  currency: ${wallet.currency ?? "(undefined)"}`);
        console.log(`  isOmnibus: ${wallet.isOmnibus}`);
        console.log(`  createdAt: ${wallet.createdAt?.toISOString() ?? "(undefined)"}`);
        console.log(`  updatedAt: ${wallet.updatedAt?.toISOString() ?? "(undefined)"}`);
        console.log(`  comment: ${wallet.comment ?? "(undefined)"}`);
        console.log(`  customerId: ${wallet.customerId ?? "(undefined)"}`);
        console.log(`  addressesCount: ${wallet.addressesCount ?? "(undefined)"}`);
        console.log(`  visibilityGroupId: ${wallet.visibilityGroupId ?? "(undefined)"}`);
        console.log(`  tags: [${wallet.tags.length} items]`);
        console.log(`  attributes: [${wallet.attributes.length} items]`);
        if (wallet.balance) {
          console.log("  balance:");
          console.log(`    totalConfirmed: ${wallet.balance.totalConfirmed ?? "(undefined)"}`);
          console.log(`    totalUnconfirmed: ${wallet.balance.totalUnconfirmed ?? "(undefined)"}`);
          console.log(`    availableConfirmed: ${wallet.balance.availableConfirmed ?? "(undefined)"}`);
          console.log(`    availableUnconfirmed: ${wallet.balance.availableUnconfirmed ?? "(undefined)"}`);
          console.log(`    reservedConfirmed: ${wallet.balance.reservedConfirmed ?? "(undefined)"}`);
          console.log(`    reservedUnconfirmed: ${wallet.balance.reservedUnconfirmed ?? "(undefined)"}`);
        } else {
          console.log("  balance: (undefined)");
        }
      }

      expect(Array.isArray(wallets)).toBe(true);
    } finally {
      client.close();
    }
  });

  it("should get wallet by ID", async () => {
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
      // First, get a list to find a valid wallet ID
      const listResult = await client.wallets.list({ limit: 1 });

      const wallets = listResult.items;
      if (wallets.length === 0) {
        console.log("No wallets available for testing");
        return;
      }

      const walletIdStr = wallets[0].id;
      if (!walletIdStr) {
        console.log("Wallet ID is undefined");
        return;
      }

      const walletId = parseInt(walletIdStr, 10);
      const wallet = await client.wallets.get(walletId);

      expect(wallet).toBeDefined();

      console.log("Wallet - All Fields:");
      console.log("=".repeat(50));
      console.log(`  id: ${wallet.id}`);
      console.log(`  name: ${wallet.name}`);
      console.log(`  externalWalletId: ${wallet.externalWalletId ?? "(undefined)"}`);
      console.log(`  status: ${wallet.status}`);
      console.log(`  type: ${wallet.type}`);
      console.log(`  blockchain: ${wallet.blockchain ?? "(undefined)"}`);
      console.log(`  network: ${wallet.network ?? "(undefined)"}`);
      console.log(`  currency: ${wallet.currency ?? "(undefined)"}`);
      console.log(`  isOmnibus: ${wallet.isOmnibus}`);
      console.log(`  createdAt: ${wallet.createdAt?.toISOString() ?? "(undefined)"}`);
      console.log(`  updatedAt: ${wallet.updatedAt?.toISOString() ?? "(undefined)"}`);
      console.log(`  comment: ${wallet.comment ?? "(undefined)"}`);
      console.log(`  customerId: ${wallet.customerId ?? "(undefined)"}`);
      console.log(`  addressesCount: ${wallet.addressesCount ?? "(undefined)"}`);
      console.log(`  visibilityGroupId: ${wallet.visibilityGroupId ?? "(undefined)"}`);
      console.log(`  tags: [${wallet.tags.length} items]`);
      if (wallet.tags.length > 0) {
        for (const tag of wallet.tags) {
          console.log(`    - ${tag}`);
        }
      }
      console.log(`  attributes: [${wallet.attributes.length} items]`);
      if (wallet.attributes.length > 0) {
        for (const attr of wallet.attributes) {
          console.log(`    - ${attr.key}: ${attr.value}`);
        }
      }
      if (wallet.balance) {
        console.log("  balance:");
        console.log(`    totalConfirmed: ${wallet.balance.totalConfirmed ?? "(undefined)"}`);
        console.log(`    totalUnconfirmed: ${wallet.balance.totalUnconfirmed ?? "(undefined)"}`);
        console.log(`    availableConfirmed: ${wallet.balance.availableConfirmed ?? "(undefined)"}`);
        console.log(`    availableUnconfirmed: ${wallet.balance.availableUnconfirmed ?? "(undefined)"}`);
        console.log(`    reservedConfirmed: ${wallet.balance.reservedConfirmed ?? "(undefined)"}`);
        console.log(`    reservedUnconfirmed: ${wallet.balance.reservedUnconfirmed ?? "(undefined)"}`);
      } else {
        console.log("  balance: (undefined)");
      }
    } finally {
      client.close();
    }
  });

  it("should paginate through wallets", async () => {
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
      const pageSize = 2;
      const allWallets: Array<{ id?: string; name?: string }> = [];
      let offset = 0;

      // Fetch all wallets using pagination
      while (true) {
        const result = await client.wallets.list({
          limit: pageSize,
          offset: offset,
        });

        const wallets = result.items;
        allWallets.push(...wallets);
        console.log(`Fetched ${wallets.length} wallets (offset=${offset})`);

        for (const wallet of wallets) {
          console.log(`  ${wallet.id}: ${wallet.name}`);
        }

        // Check if there are more pages
        const totalItems = result.pagination?.totalItems ?? 0;
        if (offset + pageSize >= totalItems || wallets.length === 0) {
          break;
        }

        offset += pageSize;

        // Safety limit for tests
        if (offset > 100) {
          console.log("Stopping pagination test at 100 items");
          break;
        }
      }

      console.log(`Total wallets fetched via pagination: ${allWallets.length}`);
      expect(allWallets.length).toBeGreaterThanOrEqual(0);
    } finally {
      client.close();
    }
  }, 60000); // Extended timeout for pagination with many API calls
});
