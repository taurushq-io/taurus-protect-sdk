"""
Whitelisted addresses had no approval workflow at all: both the for-approval read and
the approve endpoint are generated in every SDK client and were wrapped by none, so an
approver could not act on a whitelisted destination through the SDK -- verified or not.

The batch re-read is the part that needs guarding: it goes through the verifying list
path filtered by ids, so a page that omits a row must abort rather than become an
approval of fewer rows than the caller asked for.
"""

from __future__ import annotations

from unittest.mock import MagicMock, patch

import pytest
from cryptography.hazmat.primitives.asymmetric import ec

from taurus_protect.errors import IntegrityError
from taurus_protect.services.whitelisted_address_service import WhitelistedAddressService


def _service() -> WhitelistedAddressService:
    return WhitelistedAddressService(
        api_client=MagicMock(),
        whitelisting_api=MagicMock(),
        super_admin_keys=[ec.generate_private_key(ec.SECP256R1()).public_key()],
        min_valid_signatures=1,
    )


def _key():
    return ec.generate_private_key(ec.SECP256R1())


def _envelope(address_id: str, hash_value: str):
    """
    A verified envelope. Python's WhitelistedAddress carries no metadata (Go's does), so
    the hash the approver signs is only reachable on the envelope.
    """
    from taurus_protect.models.whitelisted_address import (
        SignedWhitelistedAddressEnvelope,
        WhitelistedAddress,
        WhitelistMetadata,
    )

    return SignedWhitelistedAddressEnvelope(
        metadata=WhitelistMetadata(hash=hash_value, payload_as_string="{}"),
        verified_whitelisted_address=WhitelistedAddress(id=address_id),
    )


class TestApproveWhitelistedAddresses:
    def test_signs_the_verified_hashes_in_sorted_order(self) -> None:
        service = _service()
        api = service._api
        api.whitelist_service_get_whitelisted_addresses.return_value = MagicMock(
            result=[MagicMock(), MagicMock()], total_items="2"
        )

        with patch.object(
            service,
            "_verified_addresses",
            return_value=([_envelope("3", "hash-3"), _envelope("7", "hash-7")], []),
        ):
            service.approve([7, 3], _key(), "batch approval")

        # ONE list call filtered by ids, not one GET per id.
        api.whitelist_service_get_whitelisted_addresses.assert_called_once()
        called = api.whitelist_service_get_whitelisted_addresses.call_args.kwargs
        assert called["ids"] == ["3", "7"]
        assert called["include_for_approval"] is True

        body = api.whitelist_service_approve_whitelisted_address.call_args.kwargs["body"]
        assert body.ids == ["3", "7"], "the endpoint requires ascending order"
        assert body.comment == "batch approval"
        assert body.signature

    def test_aborts_when_the_verified_read_omits_a_row(self) -> None:
        service = _service()
        api = service._api
        api.whitelist_service_get_whitelisted_addresses.return_value = MagicMock(
            result=[MagicMock()], total_items="1"
        )

        # Only one of the two requested rows survived verification.
        with patch.object(
            service, "_verified_addresses", return_value=([_envelope("3", "hash-3")], [])
        ):
            with pytest.raises(IntegrityError, match="was not returned by the verified read"):
                service.approve([7, 3], _key(), "batch approval")

        api.whitelist_service_approve_whitelisted_address.assert_not_called()

    def test_aborts_when_a_row_has_no_hash(self) -> None:
        service = _service()
        api = service._api
        api.whitelist_service_get_whitelisted_addresses.return_value = MagicMock(
            result=[MagicMock()], total_items="1"
        )

        with patch.object(
            service, "_verified_addresses", return_value=([_envelope("3", "")], [])
        ):
            with pytest.raises(IntegrityError, match="has no metadata hash"):
                service.approve([3], _key(), "batch approval")

        api.whitelist_service_approve_whitelisted_address.assert_not_called()

    @pytest.mark.parametrize(
        "ids,key,comment",
        [
            ([], "key", "ok"),
            ([1], None, "ok"),
            ([1], "key", ""),
            ([0], "key", "ok"),
            (["1"], "key", "ok"),
        ],
    )
    def test_rejects_bad_input(self, ids, key, comment) -> None:
        service = _service()
        with pytest.raises(ValueError):
            service.approve(ids, _key() if key else None, comment)
