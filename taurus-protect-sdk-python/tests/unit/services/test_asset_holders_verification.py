"""
The v2 asset-holder rows carry no signature, so query_asset_addresses returns an address
as a Taurus-PROTECT address only after its verified counterpart confirms it: an INTERNAL
row through the HSM-verified managed-address read, a WHITELISTED row through the
whitelist read.

Every request goes through the transport stub. The HSM signatures are real; only the
whitelist's 6-step verifier (tested on its own) is replaced, by one that reads the
address from the envelope's payload.
"""

from __future__ import annotations

import base64
import json
from typing import Any, Dict, List, Optional, Tuple
from unittest.mock import MagicMock
from urllib.parse import parse_qsl

import pytest
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric import ec

from taurus_protect._internal.openapi import (
    AddressesApi,
    AddressWhitelistingApi,
    AssetsApi,
    AssetV2Api,
)
from taurus_protect.crypto.signing import sign_data
from taurus_protect.errors import APIError, ContainerIntegrityError, IntegrityError
from taurus_protect.helpers.whitelist_hash_helper import parse_whitelisted_address_from_json
from taurus_protect.helpers.whitelisted_address_verifier import AddressVerificationResult
from taurus_protect.models.asset import QueryAssetAddressesResult
from taurus_protect.models.governance_rules import DecodedRulesContainer, RuleUser
from taurus_protect.models.pagination import CursorPage
from taurus_protect.services.address_service import AddressService
from taurus_protect.services.asset_service import AssetService
from taurus_protect.services.whitelisted_address_service import WhitelistedAddressService
from tests.unit.transport_stub import RecordedRequest, Reply, StubTransport, api_client

INTERNAL = "ADDRESS_TYPE_V2_INTERNAL"
WHITELISTED = "ADDRESS_TYPE_V2_WHITELISTED"
EXTERNAL = "ADDRESS_TYPE_V2_EXTERNAL"
HOLDERS_PATH = "/api/rest/v2/assets/a1/addresses/query"
MANAGED_PATH = "/api/rest/v1/addresses"
WHITELIST_PATH = "/api/rest/v1/whitelists/addresses"
NEXT = base64.b64encode(b"holders-next").decode()


def holder(
    address: str,
    address_type: Optional[str] = None,
    address_id: Optional[str] = None,
    whitelisted_id: Optional[str] = None,
) -> Dict[str, Any]:
    row: Dict[str, Any] = {
        "address": address,
        "balance": "5",
        "kycStatus": "KYC_STATUS_V2_APPROVED",
    }
    if address_type is not None:
        row["addressType"] = address_type
    if address_id is not None:
        row["addressID"] = address_id
    if whitelisted_id is not None:
        row["whitelistedAddressID"] = whitelisted_id
    return row


def holders_page(*rows: Dict[str, Any]) -> Dict[str, Any]:
    """The holders page always claims a next page, so a moved cursor shows."""
    return {"result": list(rows), "cursor": {"currentPage": NEXT, "hasNext": True}}


def page(*rows: Dict[str, Any]) -> Dict[str, Any]:
    return {"result": list(rows), "totalItems": str(len(rows))}


def whitelisted(whitelisted_id: str, address: str) -> Dict[str, Any]:
    payload = json.dumps({"address": address, "currency": "ETH", "network": "mainnet"})
    return {
        "id": whitelisted_id,
        "metadata": {"hash": f"h{whitelisted_id}", "payloadAsString": payload},
    }


class SixStep:
    """The whitelist verifier's seam: the address comes from the payload; tampered ids fail."""

    def __init__(self) -> None:
        self.tampered: set = set()

    def verify_whitelisted_address(
        self,
        envelope: Any,
        rules_container_decoder: Any,
        user_signatures_decoder: Any,
        cached_rules_container: Any = None,
    ) -> AddressVerificationResult:
        # The envelope model carries no id; whitelisted() files each row's hash as h<id>.
        if envelope.metadata.hash in {f"h{i}" for i in self.tampered}:
            raise IntegrityError(f"signature threshold not met for {envelope.metadata.hash}")
        payload = envelope.metadata.payload_as_string
        return AddressVerificationResult(
            rules_container=DecodedRulesContainer(),
            verified_hash=envelope.metadata.hash,
            verified_payload=payload,
            verified_whitelisted_address=parse_whitelisted_address_from_json(payload),
        )


class Holders:
    """An AssetService over the two verifying services the client would hand it."""

    def __init__(self, *, hsm_key: bool = True) -> None:
        self.hsm = ec.generate_private_key(ec.SECP256R1())
        users = []
        if hsm_key:
            pem = self.hsm.public_key().public_bytes(
                serialization.Encoding.PEM, serialization.PublicFormat.SubjectPublicKeyInfo
            )
            users = [RuleUser(id="hsm", public_key_pem=pem.decode(), roles=["HSMSLOT"])]
        rules_cache = MagicMock()
        rules_cache.get_decoded_rules_container.return_value = DecodedRulesContainer(users=users)

        ac = api_client()
        self.addresses = AddressService(ac, AddressesApi(ac), rules_cache)
        self.whitelisted = WhitelistedAddressService(
            ac,
            AddressWhitelistingApi(ac),
            [ec.generate_private_key(ec.SECP256R1()).public_key()],
            1,
        )
        self.six_step = SixStep()
        self.whitelisted._verifier = self.six_step
        self.service = AssetService(
            ac,
            AssetsApi(ac),
            rules_cache,
            AssetV2Api(ac),
            address_service=self.addresses,
            whitelisted_address_service=self.whitelisted,
        )

    def managed(
        self, address_id: str, address: str, signed: Optional[str] = None
    ) -> Dict[str, Any]:
        """A managed-address row whose HSM signature covers ``signed`` (default: the address)."""
        signature = sign_data(self.hsm, (signed or address).encode("utf-8"))
        return {
            "id": address_id,
            "walletId": "1",
            "address": address,
            "signature": signature,
            "status": "confirmed",
        }

    def query(self, *replies: Any) -> Tuple[QueryAssetAddressesResult, StubTransport]:
        with StubTransport(*replies) as transport:
            result = self.service.query_asset_addresses("a1")
        return result, transport


def reads(transport: StubTransport, path: str) -> List[RecordedRequest]:
    return [r for r in transport.requests if r.path == path]


def ids(request: RecordedRequest, name: str) -> List[str]:
    """In the order sent (``RecordedRequest.query`` is sorted)."""
    return [v for k, v in parse_qsl(request.raw_query) if k == name]


def kept(result: QueryAssetAddressesResult) -> List[Tuple[Optional[str], bool]]:
    return [(a.address, a.verified) for a in result.addresses]


def spy(obj: Any, name: str) -> List[Any]:
    """Record what a verified reader returned, to see where a kept address came from."""
    returned: List[Any] = []
    real = getattr(obj, name)

    def wrapper(*args: Any, **kwargs: Any) -> Any:
        returned.append(real(*args, **kwargs))
        return returned[-1]

    setattr(obj, name, wrapper)
    return returned


UNTOUCHED = CursorPage(page_size=20, next_cursor=NEXT, has_more=True)


class TestInternalRows:
    # (case, second holder row, extra managed rows by id, excluded id, reason, ids read)
    CASES = [
        (
            "confirmed",
            holder("party::8", INTERNAL, "8"),
            {"8": ("party::8", None)},
            None,
            None,
            ["7", "8"],
        ),
        (
            "address differs",
            holder("party::8", INTERNAL, "8"),
            {"8": ("party::other", None)},
            "8",
            "differs",
            ["7", "8"],
        ),
        ("not returned", holder("party::9", INTERNAL, "9"), {}, "9", "did not return", ["7", "9"]),
        ("no addressID", holder("party::10", INTERNAL), {}, "party::10", "no addressID", ["7"]),
        (
            "zero addressID",
            holder("party::11", INTERNAL, "0"),
            {},
            "party::11",
            "no addressID",
            ["7"],
        ),
        (
            "HSM signature fails",
            holder("party::12", INTERNAL, "12"),
            {"12": ("party::12", "party::forged")},
            "12",
            "did not verify",
            ["7", "12"],
        ),
    ]

    @pytest.mark.parametrize(
        "row,managed,excluded_id,reason,read_ids", [c[1:] for c in CASES], ids=[c[0] for c in CASES]
    )
    def test_confirmed_through_the_verified_managed_address_read(
        self,
        row: Dict[str, Any],
        managed: Dict[str, Tuple[str, Optional[str]]],
        excluded_id: Optional[str],
        reason: Optional[str],
        read_ids: List[str],
    ) -> None:
        h = Holders()
        rows = [h.managed("7", "party::7")] + [
            h.managed(i, address, signed) for i, (address, signed) in managed.items()
        ]

        result, transport = h.query(
            holders_page(holder("party::7", INTERNAL, "7"), row), page(*rows)
        )

        want = [("party::7", True)] + ([] if excluded_id else [("party::8", True)])
        assert kept(result) == want
        assert [e.id for e in result.excluded_unverified] == ([excluded_id] if excluded_id else [])
        if reason:
            assert reason in result.excluded_unverified[0].reason
        # Exclusions never move the cursor.
        assert result.page == UNTOUCHED
        (read,) = reads(transport, MANAGED_PATH)
        assert ids(read, "addressIds") == read_ids
        assert read.param("limit") == str(len(read_ids))
        assert reads(transport, WHITELIST_PATH) == []

    def test_a_kept_row_carries_the_verified_address(self) -> None:
        h = Holders()
        returned = spy(h.addresses, "_verified_addresses_by_id")

        result, _ = h.query(
            holders_page(holder("party::7", INTERNAL, "7")), page(h.managed("7", "party::7"))
        )

        ((verified, _failed),) = returned
        assert result.addresses[0].address is verified["7"].address
        assert result.addresses[0].verified is True


class TestWhitelistedRows:
    # (case, second holder row, extra whitelist rows, tampered ids, excluded id, reason, ids read)
    CASES = [
        (
            "confirmed",
            holder("ADDR2", WHITELISTED, whitelisted_id="2"),
            [("2", "ADDR2")],
            set(),
            None,
            None,
            ["1", "2"],
        ),
        (
            "address differs",
            holder("ADDR2x", WHITELISTED, whitelisted_id="2"),
            [("2", "ADDR2")],
            set(),
            "2",
            "differs",
            ["1", "2"],
        ),
        (
            "not returned",
            holder("ADDR3", WHITELISTED, whitelisted_id="3"),
            [],
            set(),
            "3",
            "did not return",
            ["1", "3"],
        ),
        (
            "no whitelistedAddressID",
            holder("ADDR4", WHITELISTED),
            [],
            set(),
            "ADDR4",
            "no whitelistedAddressID",
            ["1"],
        ),
        (
            "zero whitelistedAddressID",
            holder("ADDR6", WHITELISTED, whitelisted_id="0"),
            [],
            set(),
            "ADDR6",
            "no whitelistedAddressID",
            ["1"],
        ),
        (
            "envelope fails",
            holder("ADDR5", WHITELISTED, whitelisted_id="5"),
            [("5", "ADDR5")],
            {"5"},
            "5",
            "did not verify",
            ["1", "5"],
        ),
    ]

    @pytest.mark.parametrize(
        "row,extra,tampered,excluded_id,reason,read_ids",
        [c[1:] for c in CASES],
        ids=[c[0] for c in CASES],
    )
    def test_confirmed_through_the_verified_whitelist_read(
        self,
        row: Dict[str, Any],
        extra: List[Tuple[str, str]],
        tampered: set,
        excluded_id: Optional[str],
        reason: Optional[str],
        read_ids: List[str],
    ) -> None:
        h = Holders()
        h.six_step.tampered = tampered
        rows = [whitelisted("1", "ADDR1")] + [whitelisted(i, a) for i, a in extra]

        result, transport = h.query(
            holders_page(holder("ADDR1", WHITELISTED, whitelisted_id="1"), row), page(*rows)
        )

        want = [("ADDR1", True)] + ([] if excluded_id else [("ADDR2", True)])
        assert kept(result) == want
        assert [e.id for e in result.excluded_unverified] == ([excluded_id] if excluded_id else [])
        if reason:
            assert reason in result.excluded_unverified[0].reason
        assert result.page == UNTOUCHED
        (read,) = reads(transport, WHITELIST_PATH)
        assert ids(read, "ids") == read_ids
        assert read.param("limit") == str(len(read_ids))
        assert read.param("rulesContainerNormalized") == "true"
        # A whitelisted address still awaiting approval does not confirm a holder.
        assert read.param("includeForApproval") is None
        assert reads(transport, MANAGED_PATH) == []

    def test_a_kept_row_carries_the_verified_address(self) -> None:
        h = Holders()
        returned = spy(h.whitelisted, "_verified_envelopes_by_id")

        result, _ = h.query(
            holders_page(holder("ADDR1", WHITELISTED, whitelisted_id="1")),
            page(whitelisted("1", "ADDR1")),
        )

        ((envelopes, _failed),) = returned
        assert result.addresses[0].address is envelopes["1"].verified_whitelisted_address.address
        assert result.addresses[0].verified is True


def test_other_rows_are_returned_unverified_without_a_read() -> None:
    """Not even an EXTERNAL row that carries a platform id is looked up."""
    h = Holders()

    result, transport = h.query(
        holders_page(
            holder("party::ext", EXTERNAL),
            holder("party::untyped"),
            holder("party::7", EXTERNAL, "7"),
        )
    )

    assert kept(result) == [("party::ext", False), ("party::untyped", False), ("party::7", False)]
    assert result.excluded_unverified == []
    assert [r.path for r in transport.requests] == [HOLDERS_PATH]


def test_an_unknown_type_is_returned_unverified_without_a_read() -> None:
    """A type this SDK does not know decodes with its raw value, and the decision keys on
    the two verified types only -- even for a row that carries a platform id."""
    h = Holders()

    result, transport = h.query(holders_page(holder("party::7", "ADDRESS_TYPE_V2_FUTURE", "7")))

    assert kept(result) == [("party::7", False)]
    assert result.addresses[0].address_type == "ADDRESS_TYPE_V2_FUTURE"
    assert result.excluded_unverified == []
    assert [r.path for r in transport.requests] == [HOLDERS_PATH]


def test_verified_reads_are_batched_at_their_id_caps() -> None:
    """The managed-address read takes at most 50 ids, the whitelist read at most 100."""
    h = Holders()
    internal = [str(i) for i in range(1, 52)]
    listed = [str(1000 + i) for i in range(1, 102)]
    rows = [holder(f"party::{i}", INTERNAL, i) for i in internal] + [
        holder(f"ADDR{i}", WHITELISTED, whitelisted_id=i) for i in listed
    ]

    result, transport = h.query(
        holders_page(*rows),
        page(*[h.managed(i, f"party::{i}") for i in internal[:50]]),
        page(*[h.managed(i, f"party::{i}") for i in internal[50:]]),
        page(*[whitelisted(i, f"ADDR{i}") for i in listed[:100]]),
        page(*[whitelisted(i, f"ADDR{i}") for i in listed[100:]]),
    )

    assert len(result.addresses) == 152
    assert all(a.verified for a in result.addresses)
    assert result.excluded_unverified == []
    managed_reads = reads(transport, MANAGED_PATH)
    assert [len(ids(r, "addressIds")) for r in managed_reads] == [50, 1]
    assert [r.param("limit") for r in managed_reads] == ["50", "1"]
    whitelist_reads = reads(transport, WHITELIST_PATH)
    assert [len(ids(r, "ids")) for r in whitelist_reads] == [100, 1]
    assert [r.param("limit") for r in whitelist_reads] == ["100", "1"]


def test_raises_when_no_row_survives() -> None:
    """Rows came back but none can be confirmed: an error, never an empty page."""
    h = Holders()

    with pytest.raises(IntegrityError, match="all 2 asset address"):
        h.query(
            holders_page(
                holder("party::1", INTERNAL), holder("ADDR2", WHITELISTED, whitelisted_id="2")
            ),
            page(),
        )


def test_an_empty_page_is_not_an_error() -> None:
    result, transport = Holders().query(holders_page())

    assert result.addresses == [] and result.excluded_unverified == []
    assert len(transport.requests) == 1


class TestVerifiedReaderFailuresAbort:
    """A verified reader that cannot answer fails the call; its failure never becomes
    exclusions that let the rest of the page through."""

    def test_managed_address_read_fails(self) -> None:
        h = Holders()
        with pytest.raises(APIError):
            h.query(
                holders_page(holder("party::ext", EXTERNAL), holder("party::7", INTERNAL, "7")),
                Reply({"code": 13, "message": "unavailable"}, status=500),
            )

    def test_rules_container_without_an_hsm_key(self) -> None:
        h = Holders(hsm_key=False)
        with StubTransport(
            holders_page(holder("party::ext", EXTERNAL), holder("party::7", INTERNAL, "7"))
        ) as transport:
            with pytest.raises(IntegrityError, match="HSMSLOT"):
                h.service.query_asset_addresses("a1")

        assert [r.path for r in transport.requests] == [HOLDERS_PATH]

    def test_whitelist_read_fails_for_the_whole_batch(self) -> None:
        h = Holders()
        h.six_step.tampered = {"5"}
        with pytest.raises(IntegrityError, match="whitelisted address"):
            h.query(
                holders_page(
                    holder("party::ext", EXTERNAL), holder("ADDR5", WHITELISTED, whitelisted_id="5")
                ),
                page(whitelisted("5", "ADDR5")),
            )

    def test_whitelist_rules_container_cannot_be_interpreted(self) -> None:
        h = Holders()
        reply = page(whitelisted("5", "ADDR5"))
        reply["rulesContainers"] = [{"hash": "not-its-hash", "rulesContainer": "e30="}]
        with pytest.raises(ContainerIntegrityError):
            h.query(
                holders_page(
                    holder("party::ext", EXTERNAL), holder("ADDR5", WHITELISTED, whitelisted_id="5")
                ),
                reply,
            )


def test_the_verifying_services_are_mandatory() -> None:
    ac = api_client()
    with pytest.raises(ValueError, match="asset holders are verified"):
        AssetService(ac, AssetsApi(ac), MagicMock(), AssetV2Api(ac))
