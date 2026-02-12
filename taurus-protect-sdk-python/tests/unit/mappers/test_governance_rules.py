"""Unit tests for governance rules mapper."""

from __future__ import annotations

import base64
import json

import pytest

from taurus_protect.mappers.governance_rules import (
    _parse_rules_container_from_dict,
    _parse_sequential_thresholds,
    rules_container_from_base64,
    user_signatures_from_base64,
)
from taurus_protect.models.governance_rules import (
    DecodedRulesContainer,
    GroupThreshold,
    RuleUserSignature,
    SequentialThresholds,
)


class TestRulesContainerFromBase64:
    """Tests for rules_container_from_base64."""

    def test_returns_empty_for_empty_input(self) -> None:
        result = rules_container_from_base64("")
        assert isinstance(result, DecodedRulesContainer)
        assert result.users == []
        assert result.groups == []

    def test_returns_empty_for_none_input(self) -> None:
        result = rules_container_from_base64(None)
        assert isinstance(result, DecodedRulesContainer)

    def test_parses_json_rules_container(self) -> None:
        data = {
            "users": [
                {
                    "id": "user-1",
                    "name": "Admin",
                    "publicKey": "MFkwEwYH...",
                    "roles": ["ADMIN"],
                }
            ],
            "groups": [
                {
                    "id": "group-1",
                    "name": "Approvers",
                    "userIds": ["user-1"],
                }
            ],
            "minimumDistinctUserSignatures": 2,
            "minimumDistinctGroupSignatures": 1,
            "addressWhitelistingRules": [],
            "contractAddressWhitelistingRules": [],
            "enforcedRulesHash": "abc123",
            "timestamp": 1700000000,
        }
        encoded = base64.b64encode(json.dumps(data).encode()).decode()

        result = rules_container_from_base64(encoded)

        assert isinstance(result, DecodedRulesContainer)
        assert len(result.users) == 1
        assert result.users[0].id == "user-1"
        assert result.users[0].name == "Admin"
        assert result.users[0].public_key_pem == "MFkwEwYH..."
        assert result.users[0].roles == ["ADMIN"]
        assert len(result.groups) == 1
        assert result.groups[0].id == "group-1"
        assert result.groups[0].name == "Approvers"
        assert result.groups[0].user_ids == ["user-1"]
        assert result.minimum_distinct_user_signatures == 2
        assert result.minimum_distinct_group_signatures == 1
        assert result.enforced_rules_hash == "abc123"
        assert result.timestamp == 1700000000

    def test_raises_on_invalid_base64(self) -> None:
        from taurus_protect.errors import IntegrityError

        with pytest.raises(IntegrityError, match="Failed to decode"):
            rules_container_from_base64("!!!not-base64!!!")

    def test_parses_address_whitelisting_rules(self) -> None:
        data = {
            "users": [],
            "groups": [],
            "addressWhitelistingRules": [
                {
                    "currency": "ETH",
                    "network": "mainnet",
                    "parallelThresholds": [
                        {
                            "thresholds": [
                                {"groupId": "g1", "minimumSignatures": 2}
                            ]
                        }
                    ],
                    "lines": [],
                }
            ],
            "contractAddressWhitelistingRules": [],
        }
        encoded = base64.b64encode(json.dumps(data).encode()).decode()

        result = rules_container_from_base64(encoded)

        assert len(result.address_whitelisting_rules) == 1
        rule = result.address_whitelisting_rules[0]
        assert rule.currency == "ETH"
        assert rule.network == "mainnet"
        assert len(rule.parallel_thresholds) == 1
        assert len(rule.parallel_thresholds[0].thresholds) == 1
        assert rule.parallel_thresholds[0].thresholds[0].group_id == "g1"
        assert rule.parallel_thresholds[0].thresholds[0].minimum_signatures == 2

    def test_parses_contract_address_whitelisting_rules(self) -> None:
        data = {
            "users": [],
            "groups": [],
            "addressWhitelistingRules": [],
            "contractAddressWhitelistingRules": [
                {
                    "blockchain": "ETH",
                    "network": "mainnet",
                    "parallelThresholds": [
                        {"groupId": "g1", "minimumSignatures": 1}
                    ],
                }
            ],
        }
        encoded = base64.b64encode(json.dumps(data).encode()).decode()

        result = rules_container_from_base64(encoded)

        assert len(result.contract_address_whitelisting_rules) == 1
        rule = result.contract_address_whitelisting_rules[0]
        assert rule.blockchain == "ETH"
        assert rule.network == "mainnet"
        assert len(rule.parallel_thresholds) == 1


class TestParseRulesContainerFromDict:
    """Tests for _parse_rules_container_from_dict."""

    def test_handles_empty_dict(self) -> None:
        result = _parse_rules_container_from_dict({})
        assert isinstance(result, DecodedRulesContainer)
        assert result.users == []
        assert result.groups == []

    def test_handles_snake_case_keys(self) -> None:
        data = {
            "users": [
                {
                    "id": "u1",
                    "name": "User",
                    "public_key_pem": "key123",
                    "roles": [],
                }
            ],
            "groups": [
                {
                    "id": "g1",
                    "name": "Group",
                    "user_ids": ["u1"],
                }
            ],
            "minimum_distinct_user_signatures": 3,
            "minimum_distinct_group_signatures": 2,
            "address_whitelisting_rules": [],
            "contract_address_whitelisting_rules": [],
            "enforced_rules_hash": "hash",
        }
        result = _parse_rules_container_from_dict(data)

        assert len(result.users) == 1
        assert result.users[0].public_key_pem == "key123"
        assert len(result.groups) == 1
        assert result.groups[0].user_ids == ["u1"]
        assert result.minimum_distinct_user_signatures == 3
        assert result.minimum_distinct_group_signatures == 2

    def test_handles_camelcase_public_key(self) -> None:
        data = {
            "users": [
                {
                    "id": "u1",
                    "publicKeyPem": "pemKey",
                    "roles": [],
                }
            ],
        }
        result = _parse_rules_container_from_dict(data)
        assert result.users[0].public_key_pem == "pemKey"

    def test_handles_address_whitelisting_lines(self) -> None:
        data = {
            "addressWhitelistingRules": [
                {
                    "currency": "BTC",
                    "network": "mainnet",
                    "parallelThresholds": [],
                    "lines": [
                        {
                            "cells": [
                                {"type": 1, "internalWallet": {"path": "m/44'/0'/0'"}}
                            ],
                            "parallelThresholds": [
                                {"groupId": "g1", "minimumSignatures": 1}
                            ],
                        }
                    ],
                }
            ],
        }
        result = _parse_rules_container_from_dict(data)

        assert len(result.address_whitelisting_rules) == 1
        rule = result.address_whitelisting_rules[0]
        assert len(rule.lines) == 1
        line = rule.lines[0]
        assert len(line.cells) == 1
        assert line.cells[0].type == 1
        assert line.cells[0].internal_wallet is not None
        assert line.cells[0].internal_wallet.path == "m/44'/0'/0'"


class TestParseSequentialThresholds:
    """Tests for _parse_sequential_thresholds."""

    def test_handles_nested_format(self) -> None:
        data = [
            {
                "thresholds": [
                    {"groupId": "g1", "minimumSignatures": 2},
                    {"groupId": "g2", "minimumSignatures": 1},
                ]
            }
        ]
        result = _parse_sequential_thresholds(data)

        assert len(result) == 1
        assert len(result[0].thresholds) == 2
        assert result[0].thresholds[0].group_id == "g1"
        assert result[0].thresholds[0].minimum_signatures == 2
        assert result[0].thresholds[1].group_id == "g2"

    def test_handles_flat_format(self) -> None:
        data = [
            {"groupId": "g1", "minimumSignatures": 3},
            {"groupId": "g2", "minimumSignatures": 1},
        ]
        result = _parse_sequential_thresholds(data)

        assert len(result) == 2
        assert len(result[0].thresholds) == 1
        assert result[0].thresholds[0].group_id == "g1"
        assert result[0].thresholds[0].minimum_signatures == 3

    def test_handles_snake_case_keys(self) -> None:
        data = [{"group_id": "g1", "minimum_signatures": 2}]
        result = _parse_sequential_thresholds(data)

        assert len(result) == 1
        assert result[0].thresholds[0].group_id == "g1"
        assert result[0].thresholds[0].minimum_signatures == 2

    def test_handles_unknown_format(self) -> None:
        data = [{"unknown": "field"}]
        result = _parse_sequential_thresholds(data)

        assert len(result) == 1
        assert result[0].thresholds == []

    def test_handles_empty_list(self) -> None:
        result = _parse_sequential_thresholds([])
        assert result == []


class TestUserSignaturesFromBase64:
    """Tests for user_signatures_from_base64."""

    def test_returns_empty_for_empty_input(self) -> None:
        result = user_signatures_from_base64("")
        assert result == []

    def test_returns_empty_for_none(self) -> None:
        result = user_signatures_from_base64(None)
        assert result == []

    def test_parses_json_signatures(self) -> None:
        data = {
            "signatures": [
                {"userId": "u1", "signature": "sig1base64"},
                {"userId": "u2", "signature": "sig2base64"},
            ]
        }
        encoded = base64.b64encode(json.dumps(data).encode()).decode()

        result = user_signatures_from_base64(encoded)

        assert len(result) == 2
        assert result[0].user_id == "u1"
        assert result[0].signature == "sig1base64"
        assert result[1].user_id == "u2"
        assert result[1].signature == "sig2base64"

    def test_parses_json_array_format(self) -> None:
        data = [
            {"userId": "u1", "signature": "sig1"},
        ]
        encoded = base64.b64encode(json.dumps(data).encode()).decode()

        result = user_signatures_from_base64(encoded)

        assert len(result) == 1
        assert result[0].user_id == "u1"

    def test_parses_snake_case_keys(self) -> None:
        data = {
            "signatures": [
                {"user_id": "u1", "signature": "sig1"},
            ]
        }
        encoded = base64.b64encode(json.dumps(data).encode()).decode()

        result = user_signatures_from_base64(encoded)

        assert len(result) == 1
        assert result[0].user_id == "u1"
