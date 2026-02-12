package com.taurushq.sdk.protect.client.integration;

import com.taurushq.sdk.protect.client.ProtectClient;
import com.taurushq.sdk.protect.client.testutil.TestHelper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.HealthCheck;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Integration tests for HealthService.
 */
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class HealthIntegrationTest {

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
    void healthCheck() throws ApiException {
        HealthCheck health = client.getHealthService().getAllHealthChecks();

        System.out.println("Health check completed");
        if (health.getComponents() != null) {
            health.getComponents().forEach((name, component) -> {
                System.out.println("  Component: " + name);
                if (component.getGroups() != null) {
                    component.getGroups().forEach((groupName, group) ->
                            System.out.println("    Group: " + groupName));
                }
            });
        }

        assertNotNull(health, "Health check should return result");
    }
}
