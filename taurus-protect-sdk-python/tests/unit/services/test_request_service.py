"""Unit tests for RequestService."""

from __future__ import annotations

from datetime import datetime, timezone
from unittest.mock import MagicMock, patch

import pytest

from taurus_protect._internal.openapi import RequestsApi
from taurus_protect.crypto.hashing import calculate_hex_hash
from taurus_protect.errors import APIError, IntegrityError, NotFoundError
from taurus_protect.models.pagination import CursorPage
from taurus_protect.models.request import Request, RequestMetadata, RequestStatus
from taurus_protect.services.request_service import RequestService
from tests.unit.transport_stub import StubTransport, api_client


class TestGet:
    """Tests for RequestService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        requests_api = MagicMock()
        service = RequestService(api_client=api_client, requests_api=requests_api)
        return service, requests_api

    def test_get_returns_verified_request(self) -> None:
        service, api = self._make_service()

        result_dto = MagicMock()
        reply = MagicMock()
        reply.result = result_dto

        api.request_service_get_request.return_value = reply

        mock_request = MagicMock()
        mock_request.metadata = MagicMock()
        mock_request.metadata.hash = None
        mock_request.metadata.payload_as_string = None

        with patch(
            "taurus_protect.services.request_service.request_from_dto",
            return_value=mock_request,
        ):
            result = service.get(1)

        assert result is mock_request
        api.request_service_get_request.assert_called_once_with("1")

    def test_get_raises_value_error_for_non_positive_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="request_id must be positive"):
            service.get(0)

        with pytest.raises(ValueError, match="request_id must be positive"):
            service.get(-1)

    def test_get_raises_not_found_when_result_is_none(self) -> None:
        service, api = self._make_service()
        reply = MagicMock()
        reply.result = None
        api.request_service_get_request.return_value = reply

        with pytest.raises(NotFoundError):
            service.get(1)

    def test_get_verifies_request_hash(self) -> None:
        service, api = self._make_service()

        result_dto = MagicMock()
        reply = MagicMock()
        reply.result = result_dto
        api.request_service_get_request.return_value = reply

        mock_request = MagicMock()
        mock_request.metadata = MagicMock()
        mock_request.metadata.hash = "wrong_hash"
        mock_request.metadata.payload_as_string = '{"some":"payload"}'

        # IntegrityError from _verify_request_hash propagates directly
        with (
            patch(
                "taurus_protect.services.request_service.request_from_dto",
                return_value=mock_request,
            ),
            patch(
                "taurus_protect.services.request_service.calculate_hex_hash",
                return_value="correct_hash",
            ),
        ):
            with pytest.raises(IntegrityError):
                service.get(1)

    def test_get_raises_when_hash_exists_but_payload_missing(self) -> None:
        """F1: hash exists but payload stripped -> must raise, not silently skip."""
        service, api = self._make_service()

        result_dto = MagicMock()
        reply = MagicMock()
        reply.result = result_dto
        api.request_service_get_request.return_value = reply

        mock_request = MagicMock()
        mock_request.metadata = MagicMock()
        mock_request.metadata.hash = "abc123hash"
        mock_request.metadata.payload_as_string = None  # payload stripped

        with patch(
            "taurus_protect.services.request_service.request_from_dto",
            return_value=mock_request,
        ):
            with pytest.raises(IntegrityError):
                service.get(1)


def _verified_request(request_id: str = "1") -> Request:
    """
    A request whose metadata actually verifies.

    List paths now drop rows whose hash does not cover payload_as_string, so a bare
    MagicMock is excluded -- correctly. Fixtures have to describe a real artefact or
    they pin the behaviour the fix removed.
    """
    payload = '[{"key": "currency", "value": "BTC"}]'
    return Request(
        id=request_id,
        metadata=RequestMetadata(hash=calculate_hex_hash(payload), payload_as_string=payload),
    )


def _tampered_request(request_id: str = "2") -> Request:
    """
    A request whose hash no longer covers its payload -- the documented attack:
    alter the structured payload, leave the hashed string's hash in place.
    """
    good = '[{"key": "currency", "value": "BTC"}]'
    return Request(
        id=request_id,
        metadata=RequestMetadata(
            hash=calculate_hex_hash(good),
            payload_as_string='[{"key": "currency", "value": "ETH"}]',
        ),
    )


class TestList:
    """RequestService.list goes through GetRequestsV2 (cursor) instead of the deprecated v1."""

    def _service(self) -> RequestService:
        ac = api_client()
        return RequestService(ac, RequestsApi(ac))

    def test_filters_reach_the_v2_endpoint(self) -> None:
        start = datetime(2026, 1, 2, tzinfo=timezone.utc)
        with StubTransport() as transport:
            self._service().list(
                page_size=10,
                from_date=start,
                currency_id="c",
                statuses=[RequestStatus.CREATED, RequestStatus.APPROVED],
                types=["transfer"],
                ids=["1"],
                external_request_ids=["x"],
                sort_order="DESC",
            )

        assert transport.last.path == "/api/rest/v2/requests"
        query = dict(transport.last.query)
        assert query["from"].startswith("2026-01-02T00:00:00")
        assert query["currencyID"] == "c"
        assert query["types"] == "transfer"
        assert query["ids"] == "1"
        assert query["externalRequestIDs"] == "x"
        assert query["sortOrder"] == "DESC"
        assert query["cursor.pageSize"] == "10"
        assert [v for k, v in transport.last.query if k == "statuses"] == ["APPROVED", "CREATED"]

    def test_page_and_continuation(self) -> None:
        with StubTransport(
            {"result": [], "cursor": {"currentPage": "n", "hasNext": True}},
            {},
        ) as transport:
            _, page = self._service().list()
            assert page == CursorPage(page_size=20, next_cursor="n", has_more=True)
            _, page = self._service().list(cursor=page.next_cursor)

        assert page == CursorPage(page_size=20)
        assert transport.requests[1].param("cursor.currentPage") == "n"
        assert transport.requests[1].param("cursor.pageRequest") == "NEXT"

    def test_invalid_page_size_sends_nothing(self) -> None:
        with StubTransport() as transport:
            with pytest.raises(ValueError, match="page_size"):
                self._service().list(page_size=-1)

        assert transport.requests == []


class TestGetForApproval:
    """get_for_approval: valid arguments (it sent an invalid ``statuses`` and a numeric pageRequest)."""

    def _service(self) -> RequestService:
        ac = api_client()
        return RequestService(ac, RequestsApi(ac))

    def test_filters_reach_the_wire(self) -> None:
        with StubTransport() as transport:
            self._service().get_for_approval(
                page_size=5,
                currency_id="c",
                types=["transfer"],
                exclude_types=["stake"],
                ids=["1"],
                sort_order="ASC",
            )

        assert transport.last.path == "/api/rest/v2/requests/for-approval"
        assert transport.last.query == sorted(
            [
                ("currencyID", "c"),
                ("types", "transfer"),
                ("excludeTypes", "stake"),
                ("ids", "1"),
                ("sortOrder", "ASC"),
                ("cursor.pageSize", "5"),
            ]
        )

    def test_external_request_ids_is_refused_by_name(self) -> None:
        # validatord's approval paginator drops externalRequestIDs, so sending it would return an
        # unfiltered queue that reads as filtered.
        with StubTransport() as transport:
            with pytest.raises(TypeError, match="external_request_ids"):
                self._service().get_for_approval(external_request_ids=["x"])  # type: ignore[call-arg]
        assert transport.requests == []

    def test_continuation_sends_next_not_a_number(self) -> None:
        with StubTransport() as transport:
            self._service().get_for_approval(cursor="n")

        assert transport.last.param("cursor.pageRequest") == "NEXT"
        assert transport.last.param("cursor.currentPage") == "n"

    def test_statuses_cannot_filter_the_queue(self) -> None:
        with pytest.raises(TypeError, match="statuses"):
            self._service().get_for_approval(statuses=["CREATED"])  # type: ignore[call-arg]


class TestApproveRequests:
    """Tests for RequestService.approve_requests()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        requests_api = MagicMock()
        service = RequestService(api_client=api_client, requests_api=requests_api)
        return service, requests_api

    def test_approve_requests_raises_for_empty_list(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="requests list cannot be empty"):
            service.approve_requests([], MagicMock())

    def test_approve_requests_raises_for_none_key(self) -> None:
        service, _ = self._make_service()

        mock_req = MagicMock()
        mock_req.metadata = MagicMock()
        mock_req.metadata.hash = "abc"

        with pytest.raises(ValueError, match="private_key cannot be None"):
            service.approve_requests([mock_req], None)

    def test_approve_requests_raises_for_missing_hash(self) -> None:
        service, _ = self._make_service()

        mock_req = MagicMock()
        mock_req.metadata = MagicMock()
        mock_req.metadata.hash = ""

        with pytest.raises(ValueError, match="request metadata hash cannot be None or empty"):
            service.approve_requests([mock_req], MagicMock())

    def test_approve_requests_raises_for_none_metadata(self) -> None:
        service, _ = self._make_service()

        mock_req = MagicMock()
        mock_req.metadata = None

        with pytest.raises(ValueError, match="request metadata cannot be None"):
            service.approve_requests([mock_req], MagicMock())

    def test_approve_requests_returns_signed_count(self) -> None:
        service, api = self._make_service()

        # A real artefact. The previous fixture was a MagicMock whose
        # `hash_verified` was a truthy mock, so the refusal gate was vacuous and the
        # signing path was never actually exercised against a real hash.
        mock_req = _verified_request("1")

        reply = MagicMock()
        reply.signed_requests = "1"
        api.request_service_approve_requests.return_value = reply

        with patch(
            "taurus_protect.services.request_service.sign_data",
            return_value="signature_base64",
        ):
            count = service.approve_requests([mock_req], MagicMock())

        assert count == 1
        api.request_service_approve_requests.assert_called_once()

    def test_approve_requests_refuses_unverified_metadata(self) -> None:
        """The signature attests to the hash, so an unverified one must be refused."""
        from taurus_protect.errors import IntegrityError

        service, api = self._make_service()

        mock_req = _tampered_request("1")

        with pytest.raises(IntegrityError):
            service.approve_requests([mock_req], MagicMock())

        api.request_service_approve_requests.assert_not_called()

    def test_approve_requests_refuses_forged_hash_verified(self) -> None:
        """
        Why approve re-verifies instead of trusting ``hash_verified``: the flag is a
        plain model field, so this is what a Request rebuilt from JSON -- a queue, a
        webhook, a cached blob -- can look like. Trusting it gets an attacker-chosen
        hash signed by the approver's real key.
        """
        from taurus_protect.errors import IntegrityError

        service, api = self._make_service()

        forged = _tampered_request("9")
        forged = forged.model_copy(
            update={"metadata": forged.metadata.model_copy(update={"hash_verified": True})}
        )
        assert forged.metadata.hash_verified, "fixture is not exercising the forgery"

        with pytest.raises(IntegrityError, match="refusing to sign request 9"):
            service.approve_requests([forged], MagicMock())

        api.request_service_approve_requests.assert_not_called()

    def test_approve_requests_refuses_batch_when_one_row_unverified(self) -> None:
        from taurus_protect.errors import IntegrityError

        service, api = self._make_service()

        verified = _verified_request("1")
        unverified = _tampered_request("2")

        with pytest.raises(IntegrityError):
            service.approve_requests([verified, unverified], MagicMock())

        api.request_service_approve_requests.assert_not_called()


class TestRejectRequests:
    """Tests for RequestService.reject_requests()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        requests_api = MagicMock()
        service = RequestService(api_client=api_client, requests_api=requests_api)
        return service, requests_api

    def test_reject_requests_raises_for_empty_ids(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="request_ids list cannot be empty"):
            service.reject_requests([], "reason")

    def test_reject_requests_raises_for_empty_comment(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="comment"):
            service.reject_requests([1], "")

    def test_reject_requests_calls_api(self) -> None:
        service, api = self._make_service()

        service.reject_requests([1, 2], "not needed")

        api.request_service_reject_requests.assert_called_once()


class TestCreateInternalTransfer:
    """Tests for RequestService.create_internal_transfer()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        requests_api = MagicMock()
        service = RequestService(api_client=api_client, requests_api=requests_api)
        return service, requests_api

    def test_create_internal_transfer_validates_from_address(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="from_address_id must be positive"):
            service.create_internal_transfer(0, 1, "100")

    def test_create_internal_transfer_validates_to_address(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="to_address_id must be positive"):
            service.create_internal_transfer(1, 0, "100")

    def test_create_internal_transfer_validates_amount(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="amount must be positive"):
            service.create_internal_transfer(1, 2, "-1")

    def test_create_internal_transfer_returns_request(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = MagicMock()
        api.request_service_create_outgoing_request.return_value = reply

        mock_request = _verified_request("1")
        with patch(
            "taurus_protect.services.request_service.request_from_dto",
            return_value=mock_request,
        ):
            result = service.create_internal_transfer(1, 2, "100")

        # create now returns through the verifying seam, so the request comes back
        # marked -- an unverified create response is an error, not a silent pass.
        assert result.id == "1"
        assert result.metadata.hash_verified


class TestCreateExternalTransfer:
    """Tests for RequestService.create_external_transfer()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        requests_api = MagicMock()
        service = RequestService(api_client=api_client, requests_api=requests_api)
        return service, requests_api

    def test_create_external_transfer_validates_from_address(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="from_address_id must be positive"):
            service.create_external_transfer(0, 1, "100")

    def test_create_external_transfer_validates_whitelisted_address(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="to_whitelisted_address_id must be positive"):
            service.create_external_transfer(1, 0, "100")


class TestCreateCancelRequest:
    """Tests for RequestService.create_cancel_request()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        requests_api = MagicMock()
        service = RequestService(api_client=api_client, requests_api=requests_api)
        return service, requests_api

    def test_create_cancel_request_validates_address_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="address_id must be positive"):
            service.create_cancel_request(0, 1)

    def test_create_cancel_request_validates_nonce(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="nonce cannot be negative"):
            service.create_cancel_request(1, -1)


class TestListExcludesUnverifiedRows:
    """
    The defect: list() and get_for_approval() returned every row without verifying
    any of them, while get() verified. A payload altered in transit reached the
    caller with no error and no flag.
    """

    def _service(self) -> RequestService:
        ac = api_client()
        return RequestService(ac, RequestsApi(ac))

    def _tampered_request(self, request_id: str) -> Request:
        good = '[{"key": "currency", "value": "BTC"}]'
        return Request(
            id=request_id,
            metadata=RequestMetadata(
                hash=calculate_hex_hash(good),  # hash of the original
                payload_as_string='[{"key": "currency", "value": "ETH"}]',  # altered
            ),
        )

    def test_list_drops_the_tampered_row(self) -> None:
        with StubTransport({"result": [{"id": "1"}, {"id": "2"}]}):
            with patch(
                "taurus_protect.services.request_service.request_from_dto",
                side_effect=[_verified_request("1"), self._tampered_request("2")],
            ):
                requests, _ = self._service().list()

        assert [r.id for r in requests] == ["1"], "the tampered row must not be returned"
        assert requests[0].metadata.hash_verified

    def test_get_for_approval_drops_the_tampered_row(self) -> None:
        with StubTransport({"result": [{"id": "1"}, {"id": "2"}]}):
            with patch(
                "taurus_protect.services.request_service.request_from_dto",
                side_effect=[_verified_request("1"), self._tampered_request("2")],
            ):
                requests, _ = self._service().get_for_approval(page_size=10)

        assert [r.id for r in requests] == ["1"]

    def test_metadata_less_row_is_kept_unmarked(self) -> None:
        """An early-status request has no metadata yet. That is not a failure."""
        with StubTransport({"result": [{"id": "3"}]}):
            with patch(
                "taurus_protect.services.request_service.request_from_dto",
                side_effect=[Request(id="3", metadata=None)],
            ):
                requests, _ = self._service().list()

        assert [r.id for r in requests] == ["3"]
