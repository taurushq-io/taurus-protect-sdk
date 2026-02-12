package com.taurushq.sdk.protect.client.integration;

import com.taurushq.sdk.protect.client.ProtectClient;
import com.taurushq.sdk.protect.client.testutil.TestHelper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.AuditTrail;
import com.taurushq.sdk.protect.client.model.AuditTrailResult;
import com.taurushq.sdk.protect.client.model.Group;
import com.taurushq.sdk.protect.client.model.TenantConfig;
import com.taurushq.sdk.protect.client.model.VisibilityGroup;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Integration tests for administrative services (Groups, VisibilityGroups, Config, Audits).
 */
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class AdminIntegrationTest {

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
    void listGroups() throws ApiException {
        List<Group> groups = client.getGroupService().getGroups("10", "0", null, null, null);

        System.out.println("Found " + groups.size() + " groups");
        for (Group g : groups) {
            System.out.println("  " + g.getId() + ": " + g.getName());
        }

        assertNotNull(groups);
    }

    @Test
    void listVisibilityGroups() throws ApiException {
        List<VisibilityGroup> visibilityGroups = client.getVisibilityGroupService().getVisibilityGroups();

        System.out.println("Found " + visibilityGroups.size() + " visibility groups");
        for (VisibilityGroup vg : visibilityGroups) {
            System.out.println("  " + vg.getId() + ": " + vg.getName());
        }

        assertNotNull(visibilityGroups);
    }

    @Test
    void getTenantConfig() throws ApiException {
        TenantConfig config = client.getConfigService().getTenantConfig();

        System.out.println("Tenant configuration:");
        System.out.println("  Tenant ID: " + config.getTenantId());
        System.out.println("  Base Currency: " + config.getBaseCurrency());

        assertNotNull(config);
    }

    @Test
    void listAuditTrails() throws ApiException {
        AuditTrailResult result = client.getAuditService().getAuditTrails(
                null, null, null, null, null, null);

        List<AuditTrail> audits = result.getAuditTrails();
        System.out.println("Found " + audits.size() + " audit trails");
        for (AuditTrail a : audits.subList(0, Math.min(10, audits.size()))) {
            System.out.println("  " + a.getId() + ": " + a.getEntity() + "/" + a.getAction());
        }

        assertNotNull(audits);
    }
}
