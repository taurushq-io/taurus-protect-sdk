"""Cross-SDK enforced-in-rules alignment.

validatord computes the enforced-in-rules flags on some endpoints only, and omits a false
bool from its JSON. Every SDK reads a flag the same way: the wire value when present,
False when absent on an endpoint that computes it, None otherwise. The shared vector file
is consumed by the Go/Java/Python/TypeScript suites alike.
"""

from __future__ import annotations

import json
from pathlib import Path
from typing import Any, Dict, List

import pytest

from taurus_protect._internal.openapi import GroupsApi, UsersApi
from taurus_protect.models.user import Group, User
from taurus_protect.services.group_service import GroupService
from taurus_protect.services.user_service import UserService
from tests.unit.transport_stub import StubTransport, api_client

VECTORS_PATH = (
    Path(__file__).resolve().parents[3]
    / "scripts"
    / "resources"
    / "enforced-in-rules-vectors.json"
)


def load_vectors() -> list:
    with VECTORS_PATH.open(encoding="utf-8") as handle:
        return json.load(handle)


VECTORS = load_vectors()


def _user_service() -> UserService:
    ac = api_client()
    return UserService(ac, UsersApi(ac))


def _group_service() -> GroupService:
    ac = api_client()
    return GroupService(ac, GroupsApi(ac))


def _user_flags(user: User) -> Dict[str, Any]:
    return {
        "enforcedInRules": user.enforced_in_rules,
        "publicKeyEnforcedInRules": user.public_key_enforced_in_rules,
        "groupsEnforcedInRules": [g.enforced_in_rules for g in user.groups],
    }


def _group_flags(group: Group) -> Dict[str, Any]:
    return {
        "enforcedInRules": group.enforced_in_rules,
        "usersEnforcedInRules": [u.enforced_in_rules for u in group.users],
    }


def test_vector_count() -> None:
    assert len(VECTORS) == 7


@pytest.mark.parametrize("vector", VECTORS, ids=[v["name"] for v in VECTORS])
def test_vector(vector: dict) -> None:
    operation = vector["operation"]
    expected = vector["expected"]
    assert set(expected) <= {"query", "users", "groups"}, vector["name"]

    users: List[User] = []
    groups: List[Group] = []
    with StubTransport(vector["reply"]) as transport:
        if operation == "getMe":
            users = [_user_service().get_current()]
        elif operation == "getUser":
            users = [_user_service().get(vector["request"]["id"])]
        elif operation == "listUsers":
            users, _ = _user_service().list()
        elif operation == "listGroups":
            groups, _ = _group_service().list()
        else:
            pytest.fail(f"unknown operation {operation!r}")

    for name, value in expected.get("query", {}).items():
        assert transport.last.param(name) == value, name
    if "users" in expected:
        assert [_user_flags(u) for u in users] == expected["users"]
    if "groups" in expected:
        assert [_group_flags(g) for g in groups] == expected["groups"]
