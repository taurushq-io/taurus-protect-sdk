/**
 * Integration tests for Transactions API.
 *
 * These tests require a live API connection. Configure via environment variables
 * or hard-coded defaults in helpers.ts.
 *
 * Tests use the high-level TransactionService (client.transactions) rather than
 * the raw OpenAPI API to demonstrate proper SDK usage patterns.
 */

import { skipIfNotIntegration, getTestClient } from "./helpers";

describe("Integration: Transactions", () => {
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

  it("should list transactions", async () => {
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
      // Get transactions - note: API requires either currency or blockchain/network filter
      // We'll try with a common blockchain first
      const result = await client.transactions.list({
        limit: 10,
        blockchain: "ETH",
        network: "mainnet",
      });

      const transactions = result.items;
      console.log(`Found ${transactions.length} transactions`);

      if (result.pagination) {
        console.log(`Total items: ${result.pagination.totalItems}`);
      }

      for (const transaction of transactions) {
        console.log("Transaction - All Fields:");
        console.log("=".repeat(50));
        console.log(`  id: ${transaction.id}`);
        console.log(`  transactionId: ${transaction.transactionId ?? "(undefined)"}`);
        console.log(`  uniqueId: ${transaction.uniqueId ?? "(undefined)"}`);
        console.log(`  direction: ${transaction.direction ?? "(undefined)"}`);
        console.log(`  type: ${transaction.type ?? "(undefined)"}`);
        console.log(`  status: ${transaction.status ?? "(undefined)"}`);
        console.log(`  isConfirmed: ${transaction.isConfirmed}`);
        console.log(`  blockchain: ${transaction.blockchain ?? "(undefined)"}`);
        console.log(`  network: ${transaction.network ?? "(undefined)"}`);
        console.log(`  currency: ${transaction.currency ?? "(undefined)"}`);
        console.log(`  hash: ${transaction.hash ?? "(undefined)"}`);
        console.log(`  block: ${transaction.block ?? "(undefined)"}`);
        console.log(`  confirmationBlock: ${transaction.confirmationBlock ?? "(undefined)"}`);
        console.log(`  amount: ${transaction.amount ?? "(undefined)"}`);
        console.log(`  amountMainUnit: ${transaction.amountMainUnit ?? "(undefined)"}`);
        console.log(`  fee: ${transaction.fee ?? "(undefined)"}`);
        console.log(`  feeMainUnit: ${transaction.feeMainUnit ?? "(undefined)"}`);
        console.log(`  requestId: ${transaction.requestId ?? "(undefined)"}`);
        console.log(`  requestVisible: ${transaction.requestVisible ?? "(undefined)"}`);
        console.log(`  receptionDate: ${transaction.receptionDate?.toISOString() ?? "(undefined)"}`);
        console.log(`  confirmationDate: ${transaction.confirmationDate?.toISOString() ?? "(undefined)"}`);
        console.log(`  forkNumber: ${transaction.forkNumber ?? "(undefined)"}`);
        console.log(`  sources: [${transaction.sources.length} items]`);
        console.log(`  destinations: [${transaction.destinations.length} items]`);
        console.log(`  attributes: [${transaction.attributes.length} items]`);
        if (transaction.currencyInfo) {
          console.log(`  currencyInfo: { symbol: ${transaction.currencyInfo.symbol} }`);
        } else {
          console.log("  currencyInfo: (undefined)");
        }
      }

      expect(Array.isArray(transactions)).toBe(true);
    } finally {
      client.close();
    }
  });

  it("should list transactions by currency filter", async () => {
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
      // Filter by specific currency (ETH is commonly available)
      const result = await client.transactions.list({
        limit: 10,
        currency: "ETH",
      });

      const transactions = result.items;
      console.log(`Found ${transactions.length} ETH transactions`);

      if (result.pagination) {
        console.log(`Total ETH transactions: ${result.pagination.totalItems}`);
      }

      for (const transaction of transactions) {
        console.log(
          `Transaction: ID=${transaction.id}, Amount=${transaction.amount}, Direction=${transaction.direction}`
        );
        // Log sources if available (array of address info)
        if (transaction.sources && transaction.sources.length > 0) {
          console.log(
            `  Sources: ${transaction.sources.map((s) => s.address).join(", ")}`
          );
        }
        // Log destinations if available (array of address info)
        if (transaction.destinations && transaction.destinations.length > 0) {
          console.log(
            `  Destinations: ${transaction.destinations.map((d) => d.address).join(", ")}`
          );
        }
      }

      expect(Array.isArray(transactions)).toBe(true);
    } finally {
      client.close();
    }
  });

  it("should filter transactions by direction", async () => {
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
      // Filter for incoming transactions
      const result = await client.transactions.list({
        limit: 10,
        blockchain: "ETH",
        network: "mainnet",
        direction: "incoming",
      });

      const transactions = result.items;
      console.log(`Found ${transactions.length} incoming transactions`);

      for (const transaction of transactions) {
        console.log(
          `Transaction: ID=${transaction.id}, Hash=${transaction.hash}, Direction=${transaction.direction}`
        );
        // Verify all returned transactions have incoming direction
        if (transaction.direction) {
          expect(transaction.direction.toLowerCase()).toBe("incoming");
        }
      }

      expect(Array.isArray(transactions)).toBe(true);
    } finally {
      client.close();
    }
  });

  it("should paginate through transactions", async () => {
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
      const allTransactions: Array<{ id?: string; hash?: string }> = [];
      let offset = 0;

      // Fetch transactions using pagination
      while (true) {
        const result = await client.transactions.list({
          limit: pageSize,
          offset: offset,
          blockchain: "ETH",
          network: "mainnet",
        });

        const transactions = result.items;
        allTransactions.push(...transactions);
        console.log(
          `Fetched ${transactions.length} transactions (offset=${offset})`
        );

        for (const transaction of transactions) {
          console.log(`  ${transaction.id}: ${transaction.hash}`);
        }

        // Check if there are more pages
        const totalItems = result.pagination?.totalItems ?? 0;
        if (offset + pageSize >= totalItems || transactions.length === 0) {
          break;
        }

        offset += pageSize;

        // Safety limit for tests
        if (offset > 10) {
          console.log("Stopping pagination test at 10 items");
          break;
        }
      }

      console.log(
        `Total transactions fetched via pagination: ${allTransactions.length}`
      );
      expect(allTransactions.length).toBeGreaterThanOrEqual(0);
    } finally {
      client.close();
    }
  }, 60000); // Extended timeout for pagination with many API calls

  it("should get transaction by ID", async () => {
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
      // First list transactions to find a valid ID
      const result = await client.transactions.list({
        limit: 1,
        blockchain: "ETH",
        network: "mainnet",
      });

      const transactions = result.items;
      if (transactions.length === 0) {
        console.log("No transactions available to test get by ID");
        return;
      }

      const transactionId = transactions[0].id;
      console.log(`Getting transaction by ID: ${transactionId}`);

      // Get transaction by ID
      const transaction = await client.transactions.get(transactionId!);

      console.log("Transaction - All Fields:");
      console.log("=".repeat(50));
      console.log(`  id: ${transaction.id}`);
      console.log(`  transactionId: ${transaction.transactionId ?? "(undefined)"}`);
      console.log(`  uniqueId: ${transaction.uniqueId ?? "(undefined)"}`);
      console.log(`  direction: ${transaction.direction ?? "(undefined)"}`);
      console.log(`  type: ${transaction.type ?? "(undefined)"}`);
      console.log(`  status: ${transaction.status ?? "(undefined)"}`);
      console.log(`  isConfirmed: ${transaction.isConfirmed}`);
      console.log(`  blockchain: ${transaction.blockchain ?? "(undefined)"}`);
      console.log(`  network: ${transaction.network ?? "(undefined)"}`);
      console.log(`  currency: ${transaction.currency ?? "(undefined)"}`);
      console.log(`  hash: ${transaction.hash ?? "(undefined)"}`);
      console.log(`  block: ${transaction.block ?? "(undefined)"}`);
      console.log(`  confirmationBlock: ${transaction.confirmationBlock ?? "(undefined)"}`);
      console.log(`  amount: ${transaction.amount ?? "(undefined)"}`);
      console.log(`  amountMainUnit: ${transaction.amountMainUnit ?? "(undefined)"}`);
      console.log(`  fee: ${transaction.fee ?? "(undefined)"}`);
      console.log(`  feeMainUnit: ${transaction.feeMainUnit ?? "(undefined)"}`);
      console.log(`  requestId: ${transaction.requestId ?? "(undefined)"}`);
      console.log(`  requestVisible: ${transaction.requestVisible ?? "(undefined)"}`);
      console.log(`  receptionDate: ${transaction.receptionDate?.toISOString() ?? "(undefined)"}`);
      console.log(`  confirmationDate: ${transaction.confirmationDate?.toISOString() ?? "(undefined)"}`);
      console.log(`  forkNumber: ${transaction.forkNumber ?? "(undefined)"}`);
      console.log(`  sources: [${transaction.sources.length} items]`);
      for (const source of transaction.sources) {
        console.log(`    - address: ${source.address ?? "(undefined)"}`);
        console.log(`      label: ${source.label ?? "(undefined)"}`);
        console.log(`      container: ${source.container ?? "(undefined)"}`);
        console.log(`      customerId: ${source.customerId ?? "(undefined)"}`);
        console.log(`      amount: ${source.amount ?? "(undefined)"}`);
        console.log(`      amountMainUnit: ${source.amountMainUnit ?? "(undefined)"}`);
        console.log(`      type: ${source.type ?? "(undefined)"}`);
        console.log(`      idx: ${source.idx ?? "(undefined)"}`);
        console.log(`      internalAddressId: ${source.internalAddressId ?? "(undefined)"}`);
        console.log(`      whitelistedAddressId: ${source.whitelistedAddressId ?? "(undefined)"}`);
      }
      console.log(`  destinations: [${transaction.destinations.length} items]`);
      for (const dest of transaction.destinations) {
        console.log(`    - address: ${dest.address ?? "(undefined)"}`);
        console.log(`      label: ${dest.label ?? "(undefined)"}`);
        console.log(`      container: ${dest.container ?? "(undefined)"}`);
        console.log(`      customerId: ${dest.customerId ?? "(undefined)"}`);
        console.log(`      amount: ${dest.amount ?? "(undefined)"}`);
        console.log(`      amountMainUnit: ${dest.amountMainUnit ?? "(undefined)"}`);
        console.log(`      type: ${dest.type ?? "(undefined)"}`);
        console.log(`      idx: ${dest.idx ?? "(undefined)"}`);
        console.log(`      internalAddressId: ${dest.internalAddressId ?? "(undefined)"}`);
        console.log(`      whitelistedAddressId: ${dest.whitelistedAddressId ?? "(undefined)"}`);
      }
      console.log(`  attributes: [${transaction.attributes.length} items]`);
      for (const attr of transaction.attributes) {
        console.log(`    - ${attr.key}: ${attr.value}`);
      }
      if (transaction.currencyInfo) {
        console.log("  currencyInfo:");
        console.log(`    id: ${transaction.currencyInfo.id ?? "(undefined)"}`);
        console.log(`    symbol: ${transaction.currencyInfo.symbol ?? "(undefined)"}`);
        console.log(`    name: ${transaction.currencyInfo.name ?? "(undefined)"}`);
        console.log(`    decimals: ${transaction.currencyInfo.decimals ?? "(undefined)"}`);
        console.log(`    blockchain: ${transaction.currencyInfo.blockchain ?? "(undefined)"}`);
        console.log(`    network: ${transaction.currencyInfo.network ?? "(undefined)"}`);
      } else {
        console.log("  currencyInfo: (undefined)");
      }

      expect(transaction).toBeDefined();
      expect(transaction.id).toBe(transactionId);
    } finally {
      client.close();
    }
  });

  it("should get transaction by hash", async () => {
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
      // First list transactions to find a valid hash
      const result = await client.transactions.list({
        limit: 10,
        blockchain: "ETH",
        network: "mainnet",
      });

      const transactions = result.items;
      // Find a transaction with a hash
      const transactionWithHash = transactions.find((tx) => tx.hash);

      if (!transactionWithHash || !transactionWithHash.hash) {
        console.log("No transactions with hash available to test get by hash");
        return;
      }

      const txHash = transactionWithHash.hash;
      console.log(`Getting transaction by hash: ${txHash}`);

      // Get transaction by hash
      const transaction = await client.transactions.getByHash(txHash);

      console.log(`Transaction found:`);
      console.log(`  ID: ${transaction.id}`);
      console.log(`  Hash: ${transaction.hash}`);
      console.log(`  Type: ${transaction.type}`);
      console.log(`  Amount: ${transaction.amount}`);
      console.log(`  Direction: ${transaction.direction}`);

      expect(transaction).toBeDefined();
      expect(transaction.hash).toBe(txHash);
    } finally {
      client.close();
    }
  });

  it("should list transactions by address", async () => {
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
      // First list transactions to find a source address
      const result = await client.transactions.list({
        limit: 10,
        blockchain: "ETH",
        network: "mainnet",
      });

      const transactions = result.items;
      // Find a transaction with a source address
      let sourceAddress: string | undefined;
      for (const tx of transactions) {
        if (tx.sources && tx.sources.length > 0 && tx.sources[0].address) {
          sourceAddress = tx.sources[0].address;
          break;
        }
      }

      if (!sourceAddress) {
        console.log("No transactions with source address available to test");
        return;
      }

      console.log(`Listing transactions for address: ${sourceAddress}`);

      // Get transactions by address
      const addressResult = await client.transactions.listByAddress(
        sourceAddress,
        { limit: 10 }
      );

      const addressTransactions = addressResult.items;
      console.log(
        `Found ${addressTransactions.length} transactions for address`
      );

      if (addressResult.pagination) {
        console.log(`Total items: ${addressResult.pagination.totalItems}`);
      }

      for (const transaction of addressTransactions) {
        console.log(
          `Transaction: ID=${transaction.id}, Hash=${transaction.hash}, Direction=${transaction.direction}`
        );
      }

      expect(Array.isArray(addressTransactions)).toBe(true);
    } finally {
      client.close();
    }
  });
});
