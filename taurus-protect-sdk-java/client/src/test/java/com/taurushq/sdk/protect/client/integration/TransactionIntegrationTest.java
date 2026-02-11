package com.taurushq.sdk.protect.client.integration;

import com.taurushq.sdk.protect.client.ProtectClient;
import com.taurushq.sdk.protect.client.testutil.TestHelper;
import com.taurushq.sdk.protect.client.model.AddressInfo;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.Transaction;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.time.LocalDate;
import java.time.OffsetDateTime;
import java.time.ZoneOffset;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Integration tests for TransactionService.
 */
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class TransactionIntegrationTest {

    private ProtectClient client;

    @BeforeAll
    void setup() throws Exception {
        TestHelper.skipIfNotEnabled();
        client = TestHelper.getTestClient(1);
    }

    @AfterAll
    void teardown() {
        if (client != null) {
            client.close();
        }
    }

    @Test
    void getTransaction() throws ApiException {
        // First get a list to find a valid transaction ID
        List<Transaction> transactions = client.getTransactionService()
                .getTransactions(null, null, null, null, 1, 0);

        if (transactions.isEmpty()) {
            System.out.println("No transactions available for testing");
            return;
        }

        long txId = transactions.get(0).getId();
        Transaction tx = client.getTransactionService().getTransactionById(txId);

        System.out.println("Transaction details:");
        System.out.println("  ID: " + tx.getId());
        System.out.println("  Hash: " + tx.getHash());
        System.out.println("  Currency: " + tx.getCurrency());
        System.out.println("  Amount: " + tx.getAmount());

        assertNotNull(tx);
        assertEquals(txId, tx.getId());
    }

    @Test
    void listTransactions() throws ApiException {
        List<Transaction> transactions = client.getTransactionService()
                .getTransactions(null, null, null, null, 10, 0);

        System.out.println("Found " + transactions.size() + " transactions");
        for (Transaction tx : transactions) {
            System.out.println("  " + tx.getId() + ": " + tx.getCurrency() + " " + tx.getAmount());
        }

        assertNotNull(transactions);
    }

    @Test
    void listTransactionsByCurrency() throws ApiException {
        List<Transaction> transactions = client.getTransactionService()
                .getTransactions(null, null, "ETH", null, 10, 0);

        System.out.println("Found " + transactions.size() + " ETH transactions");
        for (Transaction tx : transactions) {
            System.out.println("  " + tx.getId() + ": " + tx.getHash());
        }

        assertNotNull(transactions);
    }

    @Test
    void listTransactionsByAddress() throws ApiException {
        // First get a transaction to find an address
        List<Transaction> transactions = client.getTransactionService()
                .getTransactions(null, null, null, null, 1, 0);

        if (transactions.isEmpty()) {
            System.out.println("No transactions available for testing");
            return;
        }

        // Get transactions by address (using the first source address from first transaction)
        List<AddressInfo> sources = transactions.get(0).getSources();
        if (sources == null || sources.isEmpty()) {
            System.out.println("No source addresses in transaction for testing");
            return;
        }

        String address = sources.get(0).getAddress();
        if (address == null || address.isEmpty()) {
            System.out.println("No source address in transaction for testing");
            return;
        }

        List<Transaction> txByAddress = client.getTransactionService()
                .getTransactionsByAddress(address, 10, 0);

        System.out.println("Found " + txByAddress.size() + " transactions for address " + address);

        assertNotNull(txByAddress);
    }

    @Test
    void getTransactionByHash() throws ApiException {
        // First get a transaction to find a hash
        List<Transaction> transactions = client.getTransactionService()
                .getTransactions(null, null, null, null, 1, 0);

        if (transactions.isEmpty()) {
            System.out.println("No transactions available for testing");
            return;
        }

        String hash = transactions.get(0).getHash();
        if (hash == null || hash.isEmpty()) {
            System.out.println("No hash in transaction for testing");
            return;
        }

        Transaction tx = client.getTransactionService().getTransactionByHash(hash);

        System.out.println("Transaction by hash:");
        System.out.println("  Hash: " + tx.getHash());
        System.out.println("  Currency: " + tx.getCurrency());

        assertNotNull(tx);
        assertEquals(hash, tx.getHash());
    }

    @Test
    void exportTransactions() throws ApiException {
        OffsetDateTime start = LocalDate.now().minusDays(30).atStartOfDay().atOffset(ZoneOffset.UTC);
        int limit = 10;

        List<Transaction> transactions = client.getTransactionService()
                .getTransactions(start, null, null, null, limit, 0);
        System.out.println("Found " + transactions.size() + " transactions since " + start);

        String csv = client.getTransactionService()
                .exportTransactions(start, null, null, null, limit, 0);

        int csvLines = csv.split("\r?\n", -1).length - 1; // Subtract header line
        System.out.println("Exported CSV (found " + csvLines + " out of " + transactions.size() + " expected transactions):");
        System.out.println(csv.substring(0, Math.min(500, csv.length())) + "...");

        assertNotNull(csv);
        assertTrue(csv.length() > 0);
    }
}
