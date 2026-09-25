"""Cross-SDK decode-tolerance alignment.

An additive server change must not break a decode, and every SDK must survive it the
same way: a field the client does not know lands in ``additional_properties`` and is
written back, an enum value it does not know keeps its raw string, and a value the
caller supplies is sent verbatim. The shared vector file is consumed by the
Go/Java/Python/TypeScript suites alike, so a divergence is caught in one place for all
four.
"""

from __future__ import annotations

import json
from pathlib import Path
from typing import Any, Dict
from unittest.mock import MagicMock

import pytest
from cryptography.hazmat.primitives.asymmetric import ec

from taurus_protect._internal.openapi import (
    AddressesApi,
    AddressWhitelistingApi,
    AssetsApi,
    AssetV2Api,
)
from taurus_protect._internal.openapi.api.multi_factor_signature_api import (
    MultiFactorSignatureApi,
)
from taurus_protect._internal.openapi.models.tgvalidatord_get_balances_reply import (
    TgvalidatordGetBalancesReply,
)
from taurus_protect._internal.openapi.models.tgvalidatord_get_multi_factor_signature_entities_info_reply import (  # noqa: E501
    TgvalidatordGetMultiFactorSignatureEntitiesInfoReply,
)
from taurus_protect._internal.openapi.models.tgvalidatord_multi_factor_signatures_entity_type import (  # noqa: E501
    TgvalidatordMultiFactorSignaturesEntityType,
)
from taurus_protect._internal.openapi.models.tgvalidatord_stake_account import (
    TgvalidatordStakeAccount,
)
from taurus_protect._internal.openapi.models.tgvalidatord_stake_account_type import (
    TgvalidatordStakeAccountType,
)
from taurus_protect._internal.openapi.models.tgvalidatord_token_type import (
    TgvalidatordTokenType,
)
from taurus_protect.services.address_service import AddressService
from taurus_protect.services.asset_service import AssetService
from taurus_protect.services.multi_factor_signature_service import (
    MultiFactorSignatureService,
)
from taurus_protect.services.whitelisted_address_service import WhitelistedAddressService
from tests.unit.transport_stub import StubTransport, api_client

VECTORS_PATH = (
    Path(__file__).resolve().parents[3]
    / "scripts"
    / "resources"
    / "decode-tolerance-vectors.json"
)


def load_vectors() -> list:
    with VECTORS_PATH.open(encoding="utf-8") as handle:
        return json.load(handle)


VECTORS = load_vectors()

MODELS: Dict[str, Any] = {
    "tgvalidatordGetBalancesReply": TgvalidatordGetBalancesReply,
    "tgvalidatordStakeAccount": TgvalidatordStakeAccount,
    "tgvalidatordGetMultiFactorSignatureEntitiesInfoReply": (
        TgvalidatordGetMultiFactorSignatureEntitiesInfoReply
    ),
}

ENUMS: Dict[str, Any] = {
    "tgvalidatordTokenType": TgvalidatordTokenType,
    "tgvalidatordMultiFactorSignaturesEntityType": TgvalidatordMultiFactorSignaturesEntityType,
    "tgvalidatordStakeAccountType": TgvalidatordStakeAccountType,
}

OPERATIONS = {"getMultiFactorSignatureInfo", "createMultiFactorSignatures", "queryAssetAddresses"}


def _of_kind(kind: str) -> list:
    return [v for v in VECTORS if v["kind"] == kind]


def _ids(vectors: list) -> list:
    return [v["name"] for v in vectors]


def _mfs_service() -> MultiFactorSignatureService:
    ac = api_client()
    return MultiFactorSignatureService(api_client=ac, mfs_api=MultiFactorSignatureApi(ac))


def _asset_service() -> AssetService:
    ac = api_client()
    rules_cache = MagicMock()
    return AssetService(
        ac,
        AssetsApi(ac),
        rules_cache,
        AssetV2Api(ac),
        address_service=AddressService(ac, AddressesApi(ac), rules_cache),
        whitelisted_address_service=WhitelistedAddressService(
            ac,
            AddressWhitelistingApi(ac),
            [ec.generate_private_key(ec.SECP256R1()).public_key()],
            1,
        ),
    )


def _assert_contains(body: Any, expected: Dict[str, Any]) -> None:
    assert isinstance(body, dict), f"request body is {type(body).__name__}, not a JSON object"
    for key, value in expected.items():
        assert key in body, f"request body has no {key!r}: {body}"
        assert body[key] == value, f"request body {key!r} is {body[key]!r}, want {value!r}"


def test_vectors_file_is_not_empty() -> None:
    assert VECTORS


def test_every_vector_is_consumed_here() -> None:
    """A vector this suite cannot consume fails loudly instead of being skipped."""
    for vector in VECTORS:
        kind = vector["kind"]
        assert kind in {"model", "enum", "service"}, vector["name"]
        if kind == "model":
            assert vector["model"] in MODELS, vector["name"]
        elif kind == "enum":
            assert vector["enum"] in ENUMS, vector["name"]
        else:
            assert vector["operation"] in OPERATIONS, vector["name"]


@pytest.mark.parametrize("vector", _of_kind("model"), ids=_ids(_of_kind("model")))
def test_model_vector(vector: dict) -> None:
    model = MODELS[vector["model"]].from_dict(vector["wire"])
    expected = vector["expected"]

    assert model.additional_properties == expected["rootAdditionalProperties"]
    if expected.get("lossless"):
        assert json.loads(model.to_json()) == vector["wire"]


@pytest.mark.parametrize("vector", _of_kind("enum"), ids=_ids(_of_kind("enum")))
def test_enum_vector(vector: dict) -> None:
    enum_type = ENUMS[vector["enum"]]
    expected = vector["expected"]

    for decoded in (enum_type(vector["wire"]), enum_type.from_json(json.dumps(vector["wire"]))):
        assert decoded.value == expected["value"]
        assert any(decoded is member for member in enum_type) is expected["known"]


@pytest.mark.parametrize("vector", _of_kind("service"), ids=_ids(_of_kind("service")))
def test_service_vector(vector: dict) -> None:
    operation = vector["operation"]
    request = vector["request"]
    expected = vector["expected"]

    if operation == "getMultiFactorSignatureInfo":
        with StubTransport(vector["reply"]):
            info = _mfs_service().get_multi_factor_signature_info(request["id"])

        assert info.entity_type is not None
        assert info.entity_type.value == expected["entityType"]

    elif operation == "createMultiFactorSignatures":
        with StubTransport(vector["reply"]) as transport:
            _mfs_service().create_multi_factor_signatures(
                request["entityIds"], request["entityType"]
            )

        _assert_contains(transport.last.body, expected["requestBody"])

    elif operation == "queryAssetAddresses":
        with StubTransport(vector["reply"]) as transport:
            result = _asset_service().query_asset_addresses(
                request["assetId"], address_type=request["addressType"]
            )

        # Only INTERNAL and WHITELISTED rows are confirmed through a verified reader.
        assert len(transport.requests) == 1
        _assert_contains(transport.last.body, expected["requestBody"])
        assert [(a.address, a.address_type, a.verified) for a in result.addresses] == [
            (row["address"], row["addressType"], row["verified"]) for row in expected["rows"]
        ]
        assert len(result.excluded_unverified) == expected["excludedCount"]

    else:
        pytest.fail(f"unknown operation {operation!r}")
