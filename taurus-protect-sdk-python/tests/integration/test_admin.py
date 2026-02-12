"""Integration tests for admin operations.

These tests verify administrative operations against a live Taurus-PROTECT API.
Tests include: Groups, Visibility Groups, Tenant Config, and Audit Trails.
"""

from __future__ import annotations

import pytest

from taurus_protect.client import ProtectClient


@pytest.mark.integration
def test_list_groups(client: ProtectClient) -> None:
    """Test listing groups."""
    groups, pagination = client.groups.list()

    print(f"Found {len(groups)} groups")
    if pagination:
        print(f"Total items: {pagination.total_items}, HasMore: {pagination.has_more}")

    for group in groups:
        print(f"Group: ID={group.id}, Name={group.name}")


@pytest.mark.integration
def test_list_visibility_groups(client: ProtectClient) -> None:
    """Test listing visibility groups."""
    groups, pagination = client.visibility_groups.list()

    print(f"Found {len(groups)} visibility groups")
    if pagination:
        print(f"Total items: {pagination.total_items}, HasMore: {pagination.has_more}")

    for group in groups:
        print(f"VisibilityGroup: ID={group.id}, Name={group.name}")


@pytest.mark.integration
def test_get_tenant_config(client: ProtectClient) -> None:
    """Test getting tenant configuration."""
    config = client.config.get()

    print("Tenant config:")
    print(f"  TenantID: {config.tenant_id}")
    print(f"  BaseCurrency: {config.base_currency}")
    print(f"  IsMFAMandatory: {config.is_mfa_mandatory}")


@pytest.mark.integration
def test_list_audit_trails(client: ProtectClient) -> None:
    """Test listing audit trails."""
    audits, pagination = client.audits.list()

    print(f"Found {len(audits)} audit trails")
    if pagination:
        print(f"Total items: {pagination.total_items}, HasMore: {pagination.has_more}")

    for audit in audits:
        print(f"Audit: ID={audit.id}, Type={audit.type}, Description={audit.description}")
