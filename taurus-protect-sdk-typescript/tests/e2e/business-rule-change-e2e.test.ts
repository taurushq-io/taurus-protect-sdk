/**
 * End-to-end test for the business rule change lifecycle.
 *
 * Exercises the full flow: list rules, find a target rule, propose a change
 * with one admin, approve it with another admin, verify the update, then
 * restore the original value.
 *
 * Requires at least 3 identities configured:
 *   1. Identity 1: default user (for reading rules)
 *   2. Identity 2: admin who proposes changes
 *   3. Identity 3: admin who approves changes
 */

import { getTestClient, skipIfNotEnabled, skipIfInsufficientIdentities } from '../testutil';
import type { ProtectClient } from '../../src/client';
import type { BusinessRule, ListBusinessRulesResult } from '../../src/models/business-rule';

// ── Helpers ─────────────────────────────────────────────────────────────

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

/**
 * Collects all business rules via cursor pagination.
 */
async function listAllBusinessRules(client: ProtectClient): Promise<BusinessRule[]> {
  const allRules: BusinessRule[] = [];
  let currentPage: string | undefined;
  let isFirst = true;

  do {
    const result: ListBusinessRulesResult = await client.businessRules.list({
      pageSize: 50,
      currentPage,
      pageRequest: isFirst ? 'FIRST' : 'NEXT',
    });

    allRules.push(...result.rules);
    currentPage = result.nextCursor;
    isFirst = false;
  } while (currentPage);

  return allRules;
}

/**
 * Finds a business rule by ID by iterating through paginated results.
 */
async function findRuleById(client: ProtectClient, ruleId: string): Promise<BusinessRule | undefined> {
  let currentPage: string | undefined;
  let isFirst = true;

  do {
    const result: ListBusinessRulesResult = await client.businessRules.list({
      pageSize: 50,
      currentPage,
      pageRequest: isFirst ? 'FIRST' : 'NEXT',
    });

    for (const rule of result.rules) {
      if (rule.id === ruleId) {
        return rule;
      }
    }

    currentPage = result.nextCursor;
    isFirst = false;
  } while (currentPage);

  return undefined;
}

/**
 * Polls until the business rule has the expected value, or times out after 30 seconds.
 */
async function waitForRuleValue(
  client: ProtectClient,
  ruleId: string,
  expectedValue: string,
): Promise<BusinessRule | undefined> {
  const deadline = Date.now() + 30_000;
  let rule: BusinessRule | undefined;

  while (Date.now() < deadline) {
    rule = await findRuleById(client, ruleId);
    if (rule && rule.ruleValue === expectedValue) {
      return rule;
    }
    console.log(
      `  Waiting for rule ${ruleId} to have value ${expectedValue}` +
      ` (current: ${rule?.ruleValue ?? 'null'})`
    );
    await sleep(2000);
  }

  return rule;
}

// ── Test ─────────────────────────────────────────────────────────────────

describe('BusinessRuleChange E2E', () => {
  let client: ProtectClient;
  let admin1: ProtectClient;
  let admin2: ProtectClient;

  beforeAll(() => {
    try {
      skipIfNotEnabled();
      skipIfInsufficientIdentities(3);
    } catch (e) {
      if (e instanceof Error && e.message === 'SKIP_TEST') return;
      throw e;
    }

    client = getTestClient(1);
    admin1 = getTestClient(2);
    admin2 = getTestClient(3);
  });

  afterAll(() => {
    if (client) client.close();
    if (admin1) admin1.close();
    if (admin2) admin2.close();
  });

  it('should propose, approve, verify, and restore a business rule change', async () => {
    // Skip if clients were not created (test infrastructure not configured)
    if (!client || !admin1 || !admin2) return;

    // ── Step 1: List all business rules ──────────────────────────────
    console.log('=== Step 1: Listing all business rules ===');
    const allRules = await listAllBusinessRules(client);
    console.log(`Found ${allRules.length} business rules`);

    // Print global rules
    console.log('\n--- Global rules ---');
    for (const rule of allRules) {
      if (rule.entityType?.toLowerCase() === 'global') {
        console.log(`  ${(rule.ruleKey ?? '').padEnd(45)} = ${(rule.ruleValue ?? '').padEnd(20)} [group: ${rule.ruleGroup}]`);
      }
    }

    // Print XLM-specific rules
    console.log('\n--- XLM rules ---');
    for (const rule of allRules) {
      if (rule.currency?.toUpperCase() === 'XLM') {
        console.log(
          `  id=${(rule.id ?? '').padEnd(6)} ${(rule.ruleKey ?? '').padEnd(45)}` +
          ` = ${(rule.ruleValue ?? '').padEnd(20)} [group: ${rule.ruleGroup}, entityType: ${rule.entityType}]`
        );
      }
    }

    // ── Step 2: Find target rule (transaction rule for XLM) ─────────
    console.log('\n=== Step 2: Finding target rule ===');

    let targetRule: BusinessRule | undefined;

    // Look for XLM transaction-related rules
    for (const rule of allRules) {
      const ruleKey = rule.ruleKey;
      if (!ruleKey) continue;
      const isTransactionRule = ruleKey.toLowerCase().includes('transaction');
      const isXlm = rule.currency?.toUpperCase() === 'XLM';
      if (isTransactionRule && isXlm) {
        targetRule = rule;
        console.log(`Found target rule: id=${rule.id} key=${ruleKey} value=${rule.ruleValue}`);
        break;
      }
    }

    // Fallback: pick any rule with a numeric value
    if (!targetRule) {
      for (const rule of allRules) {
        const value = rule.ruleValue;
        if (value && value.length > 0) {
          const parsed = parseInt(value, 10);
          if (!isNaN(parsed) && String(parsed) === value) {
            targetRule = rule;
            console.log(`Fallback target rule: id=${rule.id} key=${rule.ruleKey} value=${value}`);
            break;
          }
        }
      }
    }

    if (!targetRule) {
      console.log('No suitable business rule found. First 50 rule keys:');
      for (let i = 0; i < Math.min(50, allRules.length); i++) {
        const rule = allRules[i];
        console.log(`  id=${rule.id} key=${rule.ruleKey} value=${rule.ruleValue} currency=${rule.currency}`);
      }
      throw new Error('No suitable business rule found for testing.');
    }

    const originalValue = targetRule.ruleValue!;
    const targetRuleId = targetRule.id!;
    console.log(`Target rule: id=${targetRuleId} key=${targetRule.ruleKey} originalValue=${originalValue}`);

    // ── Step 3: Admin1 proposes a change ─────────────────────────────
    console.log('\n=== Step 3: Admin1 proposing change ===');
    const newValue = String(parseInt(originalValue, 10) + 1);

    const changeId = await admin1.changes.create({
      action: 'update',
      entity: 'businessrule',
      entityId: targetRuleId,
      changes: { rulevalue: newValue },
      comment: `E2E test: temporarily change value from ${originalValue} to ${newValue}`,
    });

    expect(changeId).toBeTruthy();
    console.log(`Created change: id=${changeId} (value ${originalValue} -> ${newValue})`);

    // ── Step 4: Admin2 approves the change ──────────────────────────
    console.log('\n=== Step 4: Admin2 approving change ===');
    await admin2.changes.approve(changeId);
    console.log(`Change ${changeId} approved by admin2`);

    // ── Step 5: Verify the change took effect ───────────────────────
    console.log('\n=== Step 5: Verifying change ===');
    const updatedRule = await waitForRuleValue(client, targetRuleId, newValue);
    expect(updatedRule).toBeDefined();
    expect(updatedRule!.ruleValue).toBe(newValue);
    console.log(`BEFORE: ${originalValue}`);
    console.log(`AFTER:  ${updatedRule!.ruleValue}`);
    console.log(`Verified: rule ${targetRuleId} value changed successfully`);

    // ── Step 6: Restore original value (cleanup) ────────────────────
    console.log('\n=== Step 6: Restoring original value ===');

    const restoreChangeId = await admin1.changes.create({
      action: 'update',
      entity: 'businessrule',
      entityId: targetRuleId,
      changes: { rulevalue: originalValue },
      comment: `E2E test: restore value from ${newValue} to ${originalValue}`,
    });

    expect(restoreChangeId).toBeTruthy();
    console.log(`Created restore change: id=${restoreChangeId}`);

    await admin2.changes.approve(restoreChangeId);
    console.log(`Restore change ${restoreChangeId} approved by admin2`);

    const restoredRule = await waitForRuleValue(client, targetRuleId, originalValue);
    expect(restoredRule).toBeDefined();
    expect(restoredRule!.ruleValue).toBe(originalValue);
    console.log(`BEFORE: ${newValue}`);
    console.log(`AFTER:  ${restoredRule!.ruleValue}`);
    console.log(`Verified: rule ${targetRuleId} value restored successfully`);

    console.log('\n=== E2E PASSED ===');
  }, 120_000); // 2 minute timeout
});
