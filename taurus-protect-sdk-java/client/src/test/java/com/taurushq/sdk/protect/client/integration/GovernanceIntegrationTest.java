package com.taurushq.sdk.protect.client.integration;

import com.taurushq.sdk.protect.client.ProtectClient;
import com.taurushq.sdk.protect.client.testutil.TestHelper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.GovernanceRules;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Integration tests for GovernanceRuleService.
 */
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class GovernanceIntegrationTest {

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
    void getGovernanceRulesAndVerifySignatures() throws ApiException {
        GovernanceRules rules = client.getGovernanceRuleService().getRules();

        System.out.println("Governance Rules retrieved");

        // Get whitelisting rules
        System.out.println("Whitelisting rules:");
        client.getGovernanceRuleService().getDecodedRulesContainer(rules)
                .getAddressWhitelistingRules()
                .forEach(rule -> System.out.println("  " + rule));

        assertNotNull(rules);
    }
}
