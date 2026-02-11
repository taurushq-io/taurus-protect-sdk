"""End-to-end test for the business rule change lifecycle.

Exercises the full flow: list rules, find a target rule, propose a change
with one admin, approve it with another admin, verify the update, then
restore the original value.

Requires at least 3 identities configured:
  1. Identity 1: default user (for reading rules)
  2. Identity 2: admin who proposes changes
  3. Identity 3: admin who approves changes
"""

from __future__ import annotations

import time

import pytest

from taurus_protect.client import ProtectClient
from taurus_protect.models.audit import CreateChangeRequest
from taurus_protect.models.business_rule import BusinessRule

from tests.testutil import (
    get_test_client,
    skip_if_insufficient_identities,
    skip_if_not_enabled,
)


@pytest.mark.e2e
class TestBusinessRuleChangeE2E:
    """E2E test for business rule change proposal/approval lifecycle."""

    client: ProtectClient
    admin1: ProtectClient
    admin2: ProtectClient

    @pytest.fixture(autouse=True)
    def setup_clients(self) -> None:
        skip_if_not_enabled()
        skip_if_insufficient_identities(3)

        self.client = get_test_client(1)
        self.admin1 = get_test_client(2)
        self.admin2 = get_test_client(3)

        yield

        self.client.close()
        self.admin1.close()
        self.admin2.close()

    def _list_all_rules(self) -> list:
        """List all business rules using cursor pagination."""
        all_rules = []
        page_request = "FIRST"
        current_page = None

        while True:
            result = self.client.business_rules.list(
                page_size=50,
                page_request=page_request,
                current_page=current_page,
            )
            all_rules.extend(result.rules)
            if not result.has_next:
                break
            current_page = result.current_page
            page_request = "NEXT"

        return all_rules

    def _find_rule_by_id(self, rule_id: str) -> BusinessRule:
        """Find a business rule by ID through paginated listing."""
        page_request = "FIRST"
        current_page = None

        while True:
            result = self.client.business_rules.list(
                page_size=50,
                page_request=page_request,
                current_page=current_page,
            )
            for rule in result.rules:
                if rule.id == rule_id:
                    return rule
            if not result.has_next:
                break
            current_page = result.current_page
            page_request = "NEXT"

        return None

    def _wait_for_rule_value(
        self, rule_id: str, expected_value: str, timeout_seconds: int = 30
    ) -> BusinessRule:
        """Poll until a business rule has the expected value, or timeout."""
        deadline = time.monotonic() + timeout_seconds
        rule = None
        while time.monotonic() < deadline:
            rule = self._find_rule_by_id(rule_id)
            if rule is not None and rule.rule_value == expected_value:
                return rule
            current = rule.rule_value if rule else "null"
            print(
                f"  Waiting for rule {rule_id} to have value {expected_value}"
                f" (current: {current})"
            )
            time.sleep(2)
        return rule

    def test_business_rule_change_e2e(self) -> None:
        # Step 1: List all business rules
        print("=== Step 1: Listing all business rules ===")
        all_rules = self._list_all_rules()
        print(f"Found {len(all_rules)} business rules")

        # Print global rules
        print("\n--- Global rules ---")
        for rule in all_rules:
            if rule.entity_type and rule.entity_type.lower() == "global":
                print(
                    f"  {rule.rule_key:<45} = {rule.rule_value or '':<20}"
                    f" [group: {rule.rule_group}]"
                )

        # Print XLM rules
        print("\n--- XLM rules ---")
        for rule in all_rules:
            if rule.currency and rule.currency.upper() == "XLM":
                print(
                    f"  id={rule.id:<6} {rule.rule_key or '':<45}"
                    f" = {rule.rule_value or '':<20}"
                    f" [group: {rule.rule_group}, entityType: {rule.entity_type}]"
                )

        # Step 2: Find target rule (transaction-related XLM rule)
        print("\n=== Step 2: Finding target rule ===")
        target_rule = None
        for rule in all_rules:
            rule_key = rule.rule_key
            if not rule_key:
                continue
            is_transaction = "transaction" in rule_key.lower()
            is_xlm = rule.currency and rule.currency.upper() == "XLM"
            if is_transaction and is_xlm:
                target_rule = rule
                print(
                    f"Found target rule: id={rule.id}"
                    f" key={rule_key} value={rule.rule_value}"
                )
                break

        if target_rule is None:
            # Fallback: any rule with a numeric value
            for rule in all_rules:
                value = rule.rule_value
                if value:
                    try:
                        int(value)
                        target_rule = rule
                        print(
                            f"Fallback target rule: id={rule.id}"
                            f" key={rule.rule_key} value={value}"
                        )
                        break
                    except ValueError:
                        continue

        if target_rule is None:
            print("No suitable rule found. First 50 rule keys:")
            for rule in all_rules[:50]:
                print(
                    f"  id={rule.id} key={rule.rule_key}"
                    f" value={rule.rule_value} currency={rule.currency}"
                )
            pytest.fail("No suitable business rule found for testing.")

        original_value = target_rule.rule_value
        target_rule_id = target_rule.id
        print(
            f"Target rule: id={target_rule_id}"
            f" key={target_rule.rule_key}"
            f" originalValue={original_value}"
        )

        # Step 3: Admin1 proposes a change
        print("\n=== Step 3: Admin1 proposing change ===")
        new_value = str(int(original_value) + 1)

        request = CreateChangeRequest(
            action="update",
            entity="businessrule",
            entity_id=target_rule_id,
            changes={"rulevalue": new_value},
            comment=f"E2E test: temporarily change value from {original_value} to {new_value}",
        )
        change_id = self.admin1.changes.create_change(request)
        assert change_id, "createChange should return a change ID"
        print(
            f"Created change: id={change_id}"
            f" (value {original_value} -> {new_value})"
        )

        # Step 4: Admin2 approves the change
        print("\n=== Step 4: Admin2 approving change ===")
        self.admin2.changes.approve_change(change_id)
        print(f"Change {change_id} approved by admin2")

        # Step 5: Verify the change took effect
        print("\n=== Step 5: Verifying change ===")
        updated_rule = self._wait_for_rule_value(target_rule_id, new_value)
        assert updated_rule is not None, "Should find the updated rule by ID"
        assert updated_rule.rule_value == new_value, (
            f"Rule value should be updated to {new_value}"
        )
        print(f"BEFORE: {original_value}")
        print(f"AFTER:  {updated_rule.rule_value}")
        print(f"Verified: rule {target_rule_id} value changed successfully")

        # Step 6: Restore original value (cleanup)
        print("\n=== Step 6: Restoring original value ===")
        restore_request = CreateChangeRequest(
            action="update",
            entity="businessrule",
            entity_id=target_rule_id,
            changes={"rulevalue": original_value},
            comment=f"E2E test: restore value from {new_value} to {original_value}",
        )
        restore_change_id = self.admin1.changes.create_change(restore_request)
        assert restore_change_id, "createChange should return a change ID for restore"
        print(f"Created restore change: id={restore_change_id}")

        self.admin2.changes.approve_change(restore_change_id)
        print(f"Restore change {restore_change_id} approved by admin2")

        restored_rule = self._wait_for_rule_value(target_rule_id, original_value)
        assert restored_rule is not None, "Should find the restored rule by ID"
        assert restored_rule.rule_value == original_value, (
            f"Rule value should be restored to {original_value}"
        )
        print(f"BEFORE: {new_value}")
        print(f"AFTER:  {restored_rule.rule_value}")
        print(f"Verified: rule {target_rule_id} value restored successfully")

        print("\n=== E2E PASSED ===")
