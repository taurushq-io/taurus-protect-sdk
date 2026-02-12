package com.taurushq.sdk.protect.client.integration;

import com.taurushq.sdk.protect.client.ProtectClient;
import com.taurushq.sdk.protect.client.testutil.TestHelper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.HealthCheck;
import com.taurushq.sdk.protect.client.model.PortfolioStatistics;
import com.taurushq.sdk.protect.client.model.Tag;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Integration tests for miscellaneous services (Tags, Statistics, Client lifecycle).
 */
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class MiscIntegrationTest {

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
    void listTags() throws ApiException {
        List<Tag> tags = client.getTagService().getTags();

        System.out.println("Found " + tags.size() + " tags");
        for (Tag t : tags) {
            System.out.println("  " + t.getId() + ": " + t.getValue());
        }

        assertNotNull(tags);
    }

    @Test
    void getPortfolioStatistics() throws ApiException {
        PortfolioStatistics stats = client.getStatisticsService().getPortfolioStatistics();

        System.out.println("Portfolio statistics:");
        System.out.println("  Total Balance: " + stats.getTotalBalance());
        System.out.println("  Total Balance (Base Currency): " + stats.getTotalBalanceBaseCurrency());
        System.out.println("  Wallets Count: " + stats.getWalletsCount());
        System.out.println("  Addresses Count: " + stats.getAddressesCount());

        assertNotNull(stats);
    }

    @Test
    void clientLifecycle() throws Exception {
        // Test that we can create and close multiple clients
        ProtectClient client1 = TestHelper.getTestClient(1);
        assertNotNull(client1);

        // Verify client is functional
        HealthCheck health = client1.getHealthService().getAllHealthChecks();
        assertNotNull(health);

        // Close client
        client1.close();

        // Create another client to verify closing doesn't affect new clients
        ProtectClient client2 = TestHelper.getTestClient(1);
        assertNotNull(client2);

        health = client2.getHealthService().getAllHealthChecks();
        assertNotNull(health);

        client2.close();

        System.out.println("Client lifecycle test passed");
    }
}
