"""Unit tests for request mapper."""

from types import SimpleNamespace

import pytest

from taurus_protect.mappers.request import (
    request_from_dto,
    request_metadata_from_dto,
    request_trail_from_dto,
    requests_from_dto,
    signed_request_from_dto,
    signed_requests_from_dto,
)
from taurus_protect.models.request import RequestStatus


class TestRequestFromDto:
    """Tests for request_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="123",
            tenant_id=1,
            currency="BTC",
            envelope="env-data",
            status="APPROVED",
            trails=[],
            creation_date="2024-01-15T10:30:00Z",
            update_date="2024-06-01T10:30:00Z",
            metadata=SimpleNamespace(
                hash="abc123",
                payload_as_string='{"key":"value"}',
                amount=None,
                fee=None,
                source_address=None,
                destination_address=None,
                memo=None,
            ),
            rule="default",
            signed_requests=None,
        )
        result = request_from_dto(dto)
        assert result is not None
        assert result.id == "123"
        assert result.tenant_id == 1
        assert result.currency == "BTC"
        assert result.envelope == "env-data"
        assert result.status == RequestStatus.APPROVED
        assert result.rule == "default"

    def test_returns_none_for_none(self) -> None:
        assert request_from_dto(None) is None

    def test_handles_missing_optional_fields(self) -> None:
        dto = SimpleNamespace(
            id="1",
            tenant_id=None,
            currency=None,
            envelope=None,
            status=None,
            trails=None,
            creation_date=None,
            update_date=None,
            metadata=None,
            rule=None,
            signed_requests=None,
        )
        result = request_from_dto(dto)
        assert result is not None
        assert result.id == "1"
        assert result.metadata is None

    def test_maps_trails(self) -> None:
        dto = SimpleNamespace(
            id="1",
            tenant_id=None,
            currency=None,
            envelope=None,
            status="PENDING",
            trails=[
                SimpleNamespace(
                    timestamp="2024-01-15T10:30:00Z",
                    action="created",
                    user_id="user-1",
                    comment="test",
                ),
            ],
            creation_date=None,
            update_date=None,
            metadata=None,
            rule=None,
            signed_requests=None,
        )
        result = request_from_dto(dto)
        assert result is not None
        assert len(result.trails) == 1
        assert result.trails[0].action == "created"
        assert result.trails[0].user_id == "user-1"


class TestRequestsFromDto:
    """Tests for requests_from_dto function."""

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(
                id="1", tenant_id=None, currency="BTC", envelope=None,
                status="PENDING", trails=None, creation_date=None,
                update_date=None, metadata=None, rule=None, signed_requests=None,
            ),
            SimpleNamespace(
                id="2", tenant_id=None, currency="ETH", envelope=None,
                status="APPROVED", trails=None, creation_date=None,
                update_date=None, metadata=None, rule=None, signed_requests=None,
            ),
        ]
        result = requests_from_dto(dtos)
        assert len(result) == 2
        assert result[0].id == "1"
        assert result[1].id == "2"

    def test_returns_empty_for_none(self) -> None:
        assert requests_from_dto(None) == []

    def test_returns_empty_for_empty_list(self) -> None:
        assert requests_from_dto([]) == []

    def test_filters_none_entries(self) -> None:
        result = requests_from_dto([None])
        assert result == []


class TestRequestMetadataFromDto:
    """Tests for request_metadata_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            hash="abc123",
            payload_as_string='{"test":true}',
        )
        result = request_metadata_from_dto(dto)
        assert result is not None
        assert result.hash == "abc123"
        assert result.payload_as_string == '{"test":true}'

    def test_returns_none_for_none(self) -> None:
        assert request_metadata_from_dto(None) is None


class TestRequestTrailFromDto:
    """Tests for request_trail_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            timestamp="2024-01-15T10:30:00Z",
            action="approved",
            user_id="user-42",
            comment="looks good",
        )
        result = request_trail_from_dto(dto)
        assert result.action == "approved"
        assert result.user_id == "user-42"
        assert result.comment == "looks good"
        assert result.timestamp is not None


class TestSignedRequestFromDto:
    """Tests for signed_request_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="sr-1",
            signed_request="0xdeadbeef",
            status="BROADCASTED",
            hash="txhash123",
            block=12345,
            details="confirmed",
            creation_date="2024-01-15T10:30:00Z",
            update_date="2024-06-01T10:30:00Z",
            broadcast_date="2024-01-16T10:30:00Z",
            confirmation_date="2024-01-17T10:30:00Z",
        )
        result = signed_request_from_dto(dto)
        assert result.id == "sr-1"
        assert result.signed_request == "0xdeadbeef"
        assert result.status == RequestStatus.BROADCASTED
        assert result.hash == "txhash123"
        assert result.block == 12345

    def test_handles_invalid_status(self) -> None:
        dto = SimpleNamespace(
            id="sr-1",
            signed_request="0x",
            status="UNKNOWN_STATUS",
            hash="h",
            block=0,
            details="",
            creation_date=None,
            update_date=None,
            broadcast_date=None,
            confirmation_date=None,
        )
        result = signed_request_from_dto(dto)
        assert result.status is None


class TestSignedRequestsFromDto:
    """Tests for signed_requests_from_dto function."""

    def test_returns_empty_for_none(self) -> None:
        assert signed_requests_from_dto(None) == []

    def test_returns_empty_for_empty(self) -> None:
        assert signed_requests_from_dto([]) == []

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(
                id="1", signed_request="sr1", status="BROADCASTED",
                hash="h1", block=1, details="d1",
                creation_date=None, update_date=None,
                broadcast_date=None, confirmation_date=None,
            ),
        ]
        result = signed_requests_from_dto(dtos)
        assert len(result) == 1
        assert result[0].id == "1"
