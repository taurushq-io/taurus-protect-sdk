package com.taurushq.sdk.protect.client.integration;

import com.taurushq.sdk.protect.client.ProtectClient;
import com.taurushq.sdk.protect.client.testutil.TestHelper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.AssetBalance;
import com.taurushq.sdk.protect.client.model.BalanceResult;
import com.taurushq.sdk.protect.client.model.NFTCollectionBalance;
import com.taurushq.sdk.protect.client.model.NFTCollectionBalanceResult;
import com.taurushq.sdk.protect.client.model.PageRequest;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.util.ArrayList;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Integration tests for BalanceService.
 */
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class BalanceIntegrationTest {

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
    void listBalances() throws ApiException {
        ApiRequestCursor cursor = new ApiRequestCursor(PageRequest.FIRST, 100);
        List<AssetBalance> allBalances = new ArrayList<>();
        int pageCount = 0;
        BalanceResult result;

        do {
            result = client.getBalanceService().getBalances(cursor);
            allBalances.addAll(result.getBalances());
            pageCount++;
            cursor = result.nextCursor(100);
        } while (result.hasNext() && pageCount < 5);

        System.out.println("Found " + allBalances.size() + " asset balances in " + pageCount + " pages");
        for (AssetBalance ab : allBalances.subList(0, Math.min(10, allBalances.size()))) {
            System.out.println("  " + ab);
        }

        assertNotNull(allBalances);
    }

    @Test
    void listBalancesByCurrency() throws ApiException {
        ApiRequestCursor cursor = new ApiRequestCursor(PageRequest.FIRST, 10);
        List<AssetBalance> allBalances = new ArrayList<>();
        int pageCount = 0;
        BalanceResult result;

        do {
            result = client.getBalanceService().getBalances("ETH", cursor);
            allBalances.addAll(result.getBalances());
            pageCount++;
            cursor = result.nextCursor(10);
        } while (result.hasNext() && pageCount < 3);

        System.out.println("Found " + allBalances.size() + " ETH balances in " + pageCount + " pages");
        for (AssetBalance ab : allBalances) {
            System.out.println("  " + ab);
        }

        assertNotNull(allBalances);
    }

    @Test
    void listNFTCollectionBalances() throws ApiException {
        ApiRequestCursor cursor = new ApiRequestCursor(PageRequest.FIRST, 10);
        List<NFTCollectionBalance> allBalances = new ArrayList<>();
        int pageCount = 0;
        NFTCollectionBalanceResult result;

        do {
            result = client.getBalanceService()
                    .getNFTCollectionBalances("ETH", "mainnet", cursor);
            allBalances.addAll(result.getBalances());
            pageCount++;
            cursor = result.nextCursor(10);
        } while (result.hasNext() && pageCount < 3);

        System.out.println("Found " + allBalances.size() + " NFT collection balances in " + pageCount + " pages");
        for (NFTCollectionBalance nft : allBalances) {
            System.out.println("  " + nft);
        }

        assertNotNull(allBalances);
    }
}
