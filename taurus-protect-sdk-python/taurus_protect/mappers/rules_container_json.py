"""Canonical proto-JSON <-> base64 bridge for governance rules.

Mirrors the Go SDK's rules_container_json.go: operates purely on the generated
protobuf messages via google.protobuf.json_format, independent of the
DecodedRulesContainer convenience model, so callers can round-trip a rules
container or an individual rule-cell message through canonical protobuf JSON.
"""

from __future__ import annotations

import base64
from typing import Optional

from google.protobuf import json_format, symbol_database

from taurus_protect._strict_base64 import strict_b64decode
from taurus_protect._internal.proto import request_reply_pb2 as pb

# Proto package (e.g. "tgvalidatord") derived from the descriptor, matching Go's
# pb.File_request_reply_proto.Package().
_PROTO_PACKAGE = pb.RulesContainer.DESCRIPTOR.file.package
_SYM_DB = symbol_database.Default()


def rules_container_json_from_base64(base64_data: str) -> str:
    """Decode a base64 protobuf RulesContainer into canonical proto JSON."""
    container = pb.RulesContainer()
    container.ParseFromString(strict_b64decode(base64_data))
    return json_format.MessageToJson(container)


def rules_container_base64_from_json(json_data: str) -> str:
    """Encode canonical proto JSON into a base64 protobuf RulesContainer."""
    container = pb.RulesContainer()
    json_format.Parse(json_data, container)
    return base64.b64encode(container.SerializeToString(deterministic=True)).decode("ascii")


def rule_message_base64_from_json(message_type: str, json_data: str) -> str:
    """Encode a governance rule protobuf message (resolved by name) into base64.

    message_type is a concrete rule message name from request_reply.proto such as
    RuleSource, RuleFiatAmountRange, RulesContainer_Line, or
    RulesContainer_TransactionRules.
    """
    message = _new_governance_rule_message(message_type)
    json_format.Parse(json_data, message)
    return base64.b64encode(message.SerializeToString(deterministic=True)).decode("ascii")


def rule_message_json_from_base64(message_type: str, base64_data: str) -> Optional[str]:
    """Decode a base64 protobuf rule message (by name) into canonical proto JSON.

    Returns None for empty input (an empty cell means "match any").
    """
    if base64_data is None or base64_data.strip() == "":
        return None
    data = strict_b64decode(base64_data)
    if not data:
        return None
    message = _new_governance_rule_message(message_type)
    message.ParseFromString(data)
    return json_format.MessageToJson(message)


def _new_governance_rule_message(message_type: str):
    name = (message_type or "").strip()
    if not name:
        raise ValueError("message_type cannot be empty")
    full_name = f"{_PROTO_PACKAGE}." + name.replace("_", ".")
    try:
        message_class = _SYM_DB.GetSymbol(full_name)
    except KeyError as exc:
        raise ValueError(f"unsupported governance rule message type {message_type!r}") from exc
    return message_class()
