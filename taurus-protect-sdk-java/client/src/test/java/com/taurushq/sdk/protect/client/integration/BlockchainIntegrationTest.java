package com.taurushq.sdk.protect.client.integration;

import com.taurushq.sdk.protect.client.ProtectClient;
import com.taurushq.sdk.protect.client.testutil.TestHelper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.BlockchainInfo;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Integration tests for BlockchainService.
 */
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class BlockchainIntegrationTest {

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
    void listBlockchains() throws ApiException {
        List<BlockchainInfo> blockchains = client.getBlockchainService().getBlockchains();

        System.out.println("Found " + blockchains.size() + " blockchains");
        for (BlockchainInfo b : blockchains) {
            System.out.println("  " + b.getSymbol() + " (" + b.getNetwork() + "): " + b.getName());
        }

        assertNotNull(blockchains);
        assertFalse(blockchains.isEmpty(), "Should have at least one blockchain");
    }
}
