"""Unit tests for RequestService."""

from __future__ import annotations

from unittest.mock import MagicMock, patch

import pytest

from taurus_protect.errors import APIError, IntegrityError, NotFoundError
from taurus_protect.services.request_service import RequestService


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
        with patch(
            "taurus_protect.services.request_service.request_from_dto",
            return_value=mock_request,
        ), patch(
            "taurus_protect.services.request_service.calculate_hex_hash",
            return_value="correct_hash",
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


class TestList:
    """Tests for RequestService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        requests_api = MagicMock()
        service = RequestService(api_client=api_client, requests_api=requests_api)
        return service, requests_api

    def test_list_returns_requests_and_pagination(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = [MagicMock(), MagicMock()]
        reply.total_items = "100"
        api.request_service_get_requests.return_value = reply

        with patch(
            "taurus_protect.services.request_service.requests_from_dto",
            return_value=[MagicMock(), MagicMock()],
        ):
            requests, pagination = service.list(limit=50, offset=0)

        assert len(requests) == 2
        api.request_service_get_requests.assert_called_once()

    def test_list_raises_value_error_for_invalid_limit(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="limit must be positive"):
            service.list(limit=0)

    def test_list_raises_value_error_for_negative_offset(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="offset cannot be negative"):
            service.list(offset=-1)

    def test_list_returns_empty_when_no_result(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = None
        reply.total_items = None
        api.request_service_get_requests.return_value = reply

        requests, pagination = service.list()
        assert requests == []


class TestGetForApproval:
    """Tests for RequestService.get_for_approval()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        requests_api = MagicMock()
        service = RequestService(api_client=api_client, requests_api=requests_api)
        return service, requests_api

    def test_get_for_approval_returns_requests(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = [MagicMock()]
        reply.cursor = MagicMock()
        reply.cursor.total_items = "10"
        api.request_service_get_requests_for_approval_v2.return_value = reply

        with patch(
            "taurus_protect.services.request_service.requests_from_dto",
            return_value=[MagicMock()],
        ):
            requests, pagination = service.get_for_approval(limit=10)

        assert len(requests) == 1

    def test_get_for_approval_validates_limit(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="limit must be positive"):
            service.get_for_approval(limit=0)

    def test_get_for_approval_validates_offset(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="offset cannot be negative"):
            service.get_for_approval(offset=-1)


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

        mock_req = MagicMock()
        mock_req.id = "1"
        mock_req.metadata = MagicMock()
        mock_req.metadata.hash = "hash1"

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

        mock_request = MagicMock()
        with patch(
            "taurus_protect.services.request_service.request_from_dto",
            return_value=mock_request,
        ):
            result = service.create_internal_transfer(1, 2, "100")

        assert result is mock_request


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
