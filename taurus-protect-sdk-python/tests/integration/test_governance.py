"""Integration tests for GovernanceRuleService.

These tests verify governance rules operations against a live Taurus-PROTECT API.
"""

from __future__ import annotations

import pytest

from taurus_protect.client import ProtectClient


@pytest.mark.integration
def test_get_governance_rules(client: ProtectClient) -> None:
    """Test getting governance rules."""
    rules = client.governance_rules.get_rules()

    if rules is None:
        pytest.skip("No governance rules available")

    print("Governance rules:")
    print(f"  Locked: {rules.locked}")
    print(f"  CreatedAt: {rules.creation_date}")
    print(f"  UpdatedAt: {rules.update_date}")
    print(f"  RulesContainer length: {len(rules.rules_container) if rules.rules_container else 0}")
