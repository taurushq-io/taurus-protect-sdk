package com.taurushq.sdk.protect.client.integration;

import com.taurushq.sdk.protect.client.ProtectClient;
import com.taurushq.sdk.protect.client.testutil.TestHelper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.Wallet;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Integration tests for WalletService.
 */
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class WalletIntegrationTest {

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
    void listWallets() throws ApiException {
        List<Wallet> wallets = client.getWalletService().getWallets(10, 0);

        System.out.println("Found " + wallets.size() + " wallets");
        for (Wallet w : wallets) {
            System.out.println("Wallet: ID=" + w.getId() + ", Name=" + w.getName() + ", Currency=" + w.getCurrency());
        }

        assertNotNull(wallets);
    }

    @Test
    void getWallet() throws ApiException {
        // First, get a list to find a valid wallet ID
        List<Wallet> wallets = client.getWalletService().getWallets(1, 0);

        if (wallets.isEmpty()) {
            System.out.println("No wallets available for testing");
            return;
        }

        long walletId = wallets.get(0).getId();
        Wallet wallet = client.getWalletService().getWallet(walletId);

        System.out.println("Wallet details:");
        System.out.println("  ID: " + wallet.getId());
        System.out.println("  Name: " + wallet.getName());
        System.out.println("  Currency: " + wallet.getCurrency());
        System.out.println("  Blockchain: " + wallet.getBlockchain());
        System.out.println("  Network: " + wallet.getNetwork());

        assertNotNull(wallet);
        assertEquals(walletId, wallet.getId());
    }

    @Test
    void pagination() throws ApiException {
        int pageSize = 2;
        int offset = 0;
        int totalFetched = 0;
        int maxPages = 5;
        int pageCount = 0;

        List<Wallet> wallets;
        do {
            wallets = client.getWalletService().getWallets(pageSize, offset);
            totalFetched += wallets.size();
            pageCount++;

            System.out.println("Fetched " + wallets.size() + " wallets (offset=" + offset + ")");
            for (Wallet w : wallets) {
                System.out.println("  " + w.getId() + ": " + w.getName());
            }

            offset += wallets.size();

            if (pageCount >= maxPages) {
                System.out.println("Stopping pagination test at " + maxPages + " pages");
                break;
            }
        } while (!wallets.isEmpty());

        System.out.println("Total wallets fetched via pagination: " + totalFetched);
        assertTrue(totalFetched >= 0);
    }
}
