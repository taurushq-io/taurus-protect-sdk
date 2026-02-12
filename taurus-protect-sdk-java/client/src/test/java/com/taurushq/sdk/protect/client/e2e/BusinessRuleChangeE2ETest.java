package com.taurushq.sdk.protect.client.e2e;

import com.taurushq.sdk.protect.client.ProtectClient;
import com.taurushq.sdk.protect.client.testutil.TestHelper;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.BusinessRule;
import com.taurushq.sdk.protect.client.model.BusinessRuleResult;
import com.taurushq.sdk.protect.client.model.CreateChangeRequest;
import com.taurushq.sdk.protect.client.model.PageRequest;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.util.ArrayList;
import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.fail;

/**
 * End-to-end integration test for the business rule change lifecycle.
 * <p>
 * Exercises the full flow: list rules, find a target rule, propose a change
 * with one admin, approve it with another admin, verify the update, then
 * restore the original value.
 * <p>
 * Requires at least 3 identities configured:
 * <ol>
 *   <li>Identity 1: default user (for reading rules)</li>
 *   <li>Identity 2: admin who proposes changes</li>
 *   <li>Identity 3: admin who approves changes</li>
 * </ol>
 */
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class BusinessRuleChangeE2ETest {

    private ProtectClient client;
    private ProtectClient admin1;
    private ProtectClient admin2;

    @BeforeAll
    void setup() throws Exception {
        TestHelper.skipIfNotEnabled();
        TestHelper.skipIfInsufficientIdentities(3);
        client = TestHelper.getTestClient(1);
        admin1 = TestHelper.getTestClient(2);
        admin2 = TestHelper.getTestClient(3);
    }

    @AfterAll
    void teardown() {
        if (client != null) {
            client.close();
        }
        if (admin1 != null) {
            admin1.close();
        }
        if (admin2 != null) {
            admin2.close();
        }
    }

    @Test
    void businessRuleChangeE2E() throws Exception {
        // ── Step 1: List all business rules ──────────────────────────────
        System.out.println("=== Step 1: Listing all business rules ===");
        List<BusinessRule> allRules = new ArrayList<>();
        ApiRequestCursor cursor = new ApiRequestCursor(PageRequest.FIRST, 50);
        do {
            BusinessRuleResult result = client.getBusinessRuleService().getBusinessRules(cursor);
            allRules.addAll(result.getRules());
            cursor = result.hasNext() ? result.nextCursor(50) : null;
        } while (cursor != null);
        System.out.println("Found " + allRules.size() + " business rules");

        // Print global rules (entityType=global)
        System.out.println("\n--- Global rules ---");
        for (BusinessRule rule : allRules) {
            if ("global".equalsIgnoreCase(rule.getEntityType())) {
                System.out.println(String.format("  %-45s = %-20s [group: %s]",
                        rule.getRuleKey(), rule.getRuleValue(), rule.getRuleGroup()));
            }
        }

        // Print distinct rule keys with their entity types
        Map<String, Integer> ruleKeyCounts = new LinkedHashMap<>();
        for (BusinessRule rule : allRules) {
            String key = rule.getRuleKey();
            ruleKeyCounts.put(key, ruleKeyCounts.getOrDefault(key, 0) + 1);
        }
        System.out.println("\n--- Distinct rule keys (count) ---");
        for (Map.Entry<String, Integer> entry : ruleKeyCounts.entrySet()) {
            System.out.println(String.format("  %-45s x %d", entry.getKey(), entry.getValue()));
        }

        // Print XLM-specific rules
        System.out.println("\n--- XLM rules ---");
        for (BusinessRule rule : allRules) {
            if ("XLM".equalsIgnoreCase(rule.getCurrency())) {
                System.out.println(String.format("  id=%-6s %-45s = %-20s [group: %s, entityType: %s]",
                        rule.getId(), rule.getRuleKey(), rule.getRuleValue(),
                        rule.getRuleGroup(), rule.getEntityType()));
            }
        }

        // ── Step 2: Find target rule (max transactions per day for XLM) ─
        System.out.println("\n=== Step 2: Finding target rule ===");

        BusinessRule targetRule = null;
        for (BusinessRule rule : allRules) {
            String ruleKey = rule.getRuleKey();
            if (ruleKey == null) {
                continue;
            }
            // Look for transaction-related rules involving XLM
            boolean isTransactionRule = ruleKey.toLowerCase().contains("transaction");
            boolean isXlm = "XLM".equalsIgnoreCase(rule.getCurrency());
            if (isTransactionRule && isXlm) {
                targetRule = rule;
                System.out.println("Found target rule: id=" + rule.getId()
                        + " key=" + ruleKey + " value=" + rule.getRuleValue());
                break;
            }
        }

        if (targetRule == null) {
            // Fallback: pick any rule with a numeric value that we can increment
            for (BusinessRule rule : allRules) {
                String value = rule.getRuleValue();
                if (value != null && !value.isEmpty()) {
                    try {
                        Long.parseLong(value);
                        targetRule = rule;
                        System.out.println("Fallback target rule: id=" + rule.getId()
                                + " key=" + rule.getRuleKey() + " value=" + value);
                        break;
                    } catch (NumberFormatException ignored) {
                        // not numeric, skip
                    }
                }
            }
        }

        if (targetRule == null) {
            // Print first 50 rules for diagnostics
            System.out.println("No XLM transaction rule found. First 50 rule keys:");
            for (int i = 0; i < Math.min(50, allRules.size()); i++) {
                BusinessRule rule = allRules.get(i);
                System.out.println("  id=" + rule.getId()
                        + " key=" + rule.getRuleKey()
                        + " value=" + rule.getRuleValue()
                        + " currency=" + rule.getCurrency());
            }
            fail("No suitable business rule found for testing.");
        }

        String originalValue = targetRule.getRuleValue();
        String targetRuleId = targetRule.getId();
        System.out.println("Target rule: id=" + targetRuleId
                + " key=" + targetRule.getRuleKey()
                + " originalValue=" + originalValue);

        // ── Step 3: Admin1 proposes a change ─────────────────────────────
        System.out.println("\n=== Step 3: Admin1 proposing change ===");
        String newValue = String.valueOf(Long.parseLong(originalValue) + 1);

        CreateChangeRequest request = new CreateChangeRequest();
        request.setAction("update");
        request.setEntity("businessrule");
        request.setEntityId(targetRuleId);
        request.setChanges(Collections.singletonMap("rulevalue", newValue));
        request.setComment("E2E test: temporarily change value from " + originalValue + " to " + newValue);

        String changeId = admin1.getChangeService().createChange(request);
        assertNotNull(changeId, "createChange should return a change ID");
        System.out.println("Created change: id=" + changeId + " (value " + originalValue + " -> " + newValue + ")");

        // ── Step 4: Admin2 approves the change ──────────────────────────
        System.out.println("\n=== Step 4: Admin2 approving change ===");
        admin2.getChangeService().approveChange(changeId);
        System.out.println("Change " + changeId + " approved by admin2");

        // ── Step 5: Verify the change took effect ───────────────────────
        System.out.println("\n=== Step 5: Verifying change ===");
        BusinessRule updatedRule = waitForRuleValue(targetRuleId, newValue);
        assertNotNull(updatedRule, "Should find the updated rule by ID");
        assertEquals(newValue, updatedRule.getRuleValue(),
                "Rule value should be updated to " + newValue);
        System.out.println("BEFORE: " + originalValue);
        System.out.println("AFTER:  " + updatedRule.getRuleValue());
        System.out.println("Verified: rule " + targetRuleId + " value changed successfully");

        // ── Step 6: Restore original value (cleanup) ────────────────────
        System.out.println("\n=== Step 6: Restoring original value ===");
        CreateChangeRequest restoreRequest = new CreateChangeRequest();
        restoreRequest.setAction("update");
        restoreRequest.setEntity("businessrule");
        restoreRequest.setEntityId(targetRuleId);
        restoreRequest.setChanges(Collections.singletonMap("rulevalue", originalValue));
        restoreRequest.setComment("E2E test: restore value from " + newValue + " to " + originalValue);

        String restoreChangeId = admin1.getChangeService().createChange(restoreRequest);
        assertNotNull(restoreChangeId, "createChange should return a change ID for restore");
        System.out.println("Created restore change: id=" + restoreChangeId);

        admin2.getChangeService().approveChange(restoreChangeId);
        System.out.println("Restore change " + restoreChangeId + " approved by admin2");

        BusinessRule restoredRule = waitForRuleValue(targetRuleId, originalValue);
        assertNotNull(restoredRule, "Should find the restored rule by ID");
        assertEquals(originalValue, restoredRule.getRuleValue(),
                "Rule value should be restored to " + originalValue);
        System.out.println("BEFORE: " + newValue);
        System.out.println("AFTER:  " + restoredRule.getRuleValue());
        System.out.println("Verified: rule " + targetRuleId + " value restored successfully");

        System.out.println("\n=== E2E PASSED ===");
    }

    /**
     * Finds a business rule by ID by iterating through paginated results.
     */
    private BusinessRule findRuleById(String ruleId) throws Exception {
        ApiRequestCursor cursor = new ApiRequestCursor(PageRequest.FIRST, 50);
        do {
            BusinessRuleResult result = client.getBusinessRuleService().getBusinessRules(cursor);
            for (BusinessRule rule : result.getRules()) {
                if (ruleId.equals(rule.getId())) {
                    return rule;
                }
            }
            cursor = result.hasNext() ? result.nextCursor(50) : null;
        } while (cursor != null);
        return null;
    }

    /**
     * Polls until the business rule has the expected value, or times out after 30 seconds.
     */
    private BusinessRule waitForRuleValue(String ruleId, String expectedValue) throws Exception {
        long deadline = System.currentTimeMillis() + 30_000L;
        BusinessRule rule = null;
        while (System.currentTimeMillis() < deadline) {
            rule = findRuleById(ruleId);
            if (rule != null && expectedValue.equals(rule.getRuleValue())) {
                return rule;
            }
            System.out.println("  Waiting for rule " + ruleId + " to have value " + expectedValue
                    + " (current: " + (rule != null ? rule.getRuleValue() : "null") + ")");
            Thread.sleep(2000L);
        }
        return rule;
    }
}
