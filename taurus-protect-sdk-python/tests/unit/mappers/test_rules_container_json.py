"""Tests for the canonical proto-JSON <-> base64 governance mappers."""

import base64

import pytest
from google.protobuf import json_format

from taurus_protect._internal.proto import request_reply_pb2 as pb
from taurus_protect.mappers.rules_container_json import (
    rule_message_base64_from_json,
    rule_message_json_from_base64,
    rules_container_base64_from_json,
    rules_container_json_from_base64,
)


def test_rules_container_json_round_trip():
    src = pb.RulesContainer()
    json_format.Parse('{"minimumDistinctUserSignatures": 3}', src)
    encoded = base64.b64encode(src.SerializeToString()).decode("ascii")

    json_str = rules_container_json_from_base64(encoded)
    assert "minimumDistinctUserSignatures" in json_str

    re_encoded = rules_container_base64_from_json(json_str)
    assert re_encoded == encoded


def test_json_bridge_is_order_independent_and_cross_sdk_stable():
    """The JSON bridge is the path the MCP governance tools drive: decode a container to
    JSON, let a caller edit it, re-encode, submit — and approvers sign the re-encoded
    bytes. An edited container arrives with its map keys in arbitrary order, so encoding
    MUST be order-independent: every runtime emits map<string, bytes> in unspecified
    order unless told otherwise. The input below carries the five ``properties`` keys in
    reverse-sorted order on purpose; feeding back the canonical (already sorted) JSON
    would pass even with a non-deterministic encoder and prove nothing. Expected bytes
    are the same vector asserted by test_encoding_is_deterministic_and_cross_sdk_stable
    for the typed encoder, and the same reversed input is asserted in all four SDKs.
    """
    expected = (
        "CmgKAnUxEgNQRU0aAQMiEAoGa0FscGhhEgZrQWxwaGEiEAoGa0JyYXZvEgZrQnJhdm8iFAoIa0No"
        "YXJsaWUSCGtDaGFybGllIhAKBmtEZWx0YRIGa0RlbHRhIg4KBWtFY2hvEgVrRWNob0oQCgZrQWxw"
        "aGESBmtBbHBoYUoQCgZrQnJhdm8SBmtCcmF2b0oUCghrQ2hhcmxpZRIIa0NoYXJsaWVKEAoGa0Rl"
        "bHRhEgZrRGVsdGFKDgoFa0VjaG8SBWtFY2hv"
    )
    reversed_key_json = (
        '{"users":[{"id":"u1","publicKey":"PEM","roles":["SUPERADMIN"],"properties":'
        '{"kEcho":"a0VjaG8=","kDelta":"a0RlbHRh","kCharlie":"a0NoYXJsaWU=",'
        '"kBravo":"a0JyYXZv","kAlpha":"a0FscGhh"}}],"properties":'
        '{"kEcho":"a0VjaG8=","kDelta":"a0RlbHRh","kCharlie":"a0NoYXJsaWU=",'
        '"kBravo":"a0JyYXZv","kAlpha":"a0FscGhh"}}'
    )

    for _ in range(20):
        assert rules_container_base64_from_json(reversed_key_json) == expected


def test_rules_container_base64_from_json_rejects_invalid_json():
    with pytest.raises(json_format.ParseError):
        rules_container_base64_from_json("not json")


def test_rule_message_round_trip_stable():
    encoded = rule_message_base64_from_json("RuleFiatAmountRange", '{"minAmount": "1000"}')
    decoded = rule_message_json_from_base64("RuleFiatAmountRange", encoded)
    assert decoded is not None
    assert "1000" in decoded
    # Re-encoding the decoded JSON is stable.
    assert rule_message_base64_from_json("RuleFiatAmountRange", decoded) == encoded


def test_rule_message_json_from_base64_empty_returns_none():
    assert rule_message_json_from_base64("RuleSource", "") is None
    assert rule_message_json_from_base64("RuleSource", "   ") is None


def test_rule_message_unknown_type_raises():
    with pytest.raises(ValueError):
        rule_message_base64_from_json("NoSuchRuleMessage", "{}")


def test_rule_message_empty_type_raises():
    with pytest.raises(ValueError):
        rule_message_base64_from_json("   ", "{}")
