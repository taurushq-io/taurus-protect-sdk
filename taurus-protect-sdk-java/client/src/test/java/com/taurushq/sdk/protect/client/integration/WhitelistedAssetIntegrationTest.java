package com.taurushq.sdk.protect.client.integration;

import com.taurushq.sdk.protect.client.ProtectClient;
import com.taurushq.sdk.protect.client.testutil.TestHelper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.SignedWhitelistedAssetEnvelope;
import com.taurushq.sdk.protect.client.model.WhitelistException;
import com.taurushq.sdk.protect.client.model.WhitelistedAsset;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Integration tests for WhitelistedAssetService.
 */
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class WhitelistedAssetIntegrationTest {

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
    void getWhitelistedAsset() throws ApiException, WhitelistException {
        // First get a list to find a valid ID
        List<SignedWhitelistedAssetEnvelope> assets = client.getWhitelistedAssetService()
                .getWhitelistedAssets(1, 0);

        if (assets.isEmpty()) {
            System.out.println("No whitelisted assets available for testing");
            return;
        }

        long id = assets.get(0).getId();
        WhitelistedAsset asset = client.getWhitelistedAssetService().getWhitelistedAsset(id);

        System.out.println("Whitelisted asset:");
        System.out.println("  Blockchain: " + asset.getBlockchain());
        System.out.println("  Symbol: " + asset.getSymbol());
        System.out.println("  Contract Address: " + asset.getContractAddress());
        System.out.println("  Name: " + asset.getName());
        System.out.println("  Decimals: " + asset.getDecimals());
        System.out.println("  Kind: " + asset.getKind());
        System.out.println("  Network: " + asset.getNetwork());

        assertNotNull(asset);
    }

    @Test
    void listWhitelistedAssets() throws ApiException, WhitelistException {
        List<SignedWhitelistedAssetEnvelope> assets = client.getWhitelistedAssetService()
                .getWhitelistedAssets(10, 0);

        System.out.println("Found " + assets.size() + " whitelisted assets");
        for (SignedWhitelistedAssetEnvelope envelope : assets) {
            WhitelistedAsset asset = envelope.getWhitelistedAsset();
            System.out.printf("ID: %d, Blockchain: %s, Symbol: %s, Contract: %s%n",
                    envelope.getId(),
                    asset.getBlockchain(),
                    asset.getSymbol(),
                    asset.getContractAddress());
        }

        assertNotNull(assets);
    }

    @Test
    void listWhitelistedAssetsByBlockchain() throws ApiException, WhitelistException {
        List<SignedWhitelistedAssetEnvelope> ethAssets = client.getWhitelistedAssetService()
                .getWhitelistedAssets(10, 0, "ETH");

        System.out.println("Found " + ethAssets.size() + " ETH whitelisted assets");
        for (SignedWhitelistedAssetEnvelope envelope : ethAssets) {
            WhitelistedAsset asset = envelope.getWhitelistedAsset();
            System.out.printf("Symbol: %s, Contract: %s%n",
                    asset.getSymbol(),
                    asset.getContractAddress());
        }

        assertNotNull(ethAssets);
    }

    @Test
    void listWhitelistedAssetsByBlockchainAndNetwork() throws ApiException, WhitelistException {
        List<SignedWhitelistedAssetEnvelope> mainnetAssets = client.getWhitelistedAssetService()
                .getWhitelistedAssets(10, 0, "ETH", "mainnet");

        System.out.println("Found " + mainnetAssets.size() + " ETH mainnet whitelisted assets");
        for (SignedWhitelistedAssetEnvelope envelope : mainnetAssets) {
            WhitelistedAsset asset = envelope.getWhitelistedAsset();
            System.out.printf("Symbol: %s, Network: %s, Contract: %s%n",
                    asset.getSymbol(),
                    asset.getNetwork(),
                    asset.getContractAddress());
        }

        assertNotNull(mainnetAssets);
    }

    @Test
    void paginateAllWhitelistedAssets() throws ApiException, WhitelistException {
        int limit = 10;
        int offset = 0;
        int totalCount = 0;
        int maxIterations = 500;
        int iteration = 0;

        List<SignedWhitelistedAssetEnvelope> assets;
        do {
            assets = client.getWhitelistedAssetService()
                    .getWhitelistedAssets(limit, offset);

            for (SignedWhitelistedAssetEnvelope envelope : assets) {
                WhitelistedAsset asset = envelope.getWhitelistedAsset();
                System.out.printf("ID: %d, Blockchain: %s, Symbol: %s, Contract: %s, Network: %s%n",
                        envelope.getId(),
                        asset.getBlockchain(),
                        asset.getSymbol(),
                        asset.getContractAddress(),
                        asset.getNetwork());
                totalCount++;
            }

            offset += assets.size();
            iteration++;
        } while (assets.size() == limit && iteration < maxIterations);

        System.out.println("Total whitelisted assets scanned: " + totalCount);
        assertTrue(totalCount >= 0);
    }
}
