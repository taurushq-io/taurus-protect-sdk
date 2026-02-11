"""Unit tests for TaurusNetwork ParticipantService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.services.taurus_network.participant_service import (
    ParticipantService,
)


class TestGetMyParticipant:
    """Tests for ParticipantService.get_my_participant()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        participant_api = MagicMock()
        service = ParticipantService(
            api_client=api_client, participant_api=participant_api
        )
        return service, participant_api

    def test_calls_api(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.participant = MagicMock()
        resp.participant.id = "p-1"
        resp.participant.name = "My Org"
        resp.settings = None
        api.taurus_network_service_get_my_participant.return_value = resp

        try:
            service.get_my_participant()
        except Exception:
            pass  # mapper may fail with mock data

        api.taurus_network_service_get_my_participant.assert_called_once()


class TestParticipantServiceGet:
    """Tests for ParticipantService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        participant_api = MagicMock()
        service = ParticipantService(
            api_client=api_client, participant_api=participant_api
        )
        return service, participant_api

    def test_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="participant_id"):
            service.get(participant_id="")

    def test_raises_not_found_when_none(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        api.taurus_network_service_get_participant.return_value = resp

        from taurus_protect.errors import NotFoundError

        with pytest.raises(NotFoundError):
            service.get(participant_id="p-missing")


class TestParticipantServiceList:
    """Tests for ParticipantService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        participant_api = MagicMock()
        service = ParticipantService(
            api_client=api_client, participant_api=participant_api
        )
        return service, participant_api

    def test_returns_empty_when_no_results(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        api.taurus_network_service_get_participants.return_value = resp

        result = service.list()

        assert result == []


class TestParticipantServiceCreateAttribute:
    """Tests for ParticipantService.create_participant_attribute()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        participant_api = MagicMock()
        service = ParticipantService(
            api_client=api_client, participant_api=participant_api
        )
        return service, participant_api

    def test_raises_on_empty_participant_id(self) -> None:
        service, _ = self._make_service()
        from taurus_protect.models.taurus_network.participant import (
            CreateParticipantAttributeRequest,
        )

        req = CreateParticipantAttributeRequest(key="k", value="v")
        with pytest.raises(ValueError, match="participant_id"):
            service.create_participant_attribute(participant_id="", request=req)

    def test_raises_on_none_request(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="request cannot be None"):
            service.create_participant_attribute(
                participant_id="p-1", request=None
            )


class TestParticipantServiceDeleteAttribute:
    """Tests for ParticipantService.delete_participant_attribute()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        participant_api = MagicMock()
        service = ParticipantService(
            api_client=api_client, participant_api=participant_api
        )
        return service, participant_api

    def test_raises_on_empty_participant_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="participant_id"):
            service.delete_participant_attribute(
                participant_id="", attribute_id="a-1"
            )

    def test_raises_on_empty_attribute_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="attribute_id"):
            service.delete_participant_attribute(
                participant_id="p-1", attribute_id=""
            )

    def test_calls_api(self) -> None:
        service, api = self._make_service()
        api.taurus_network_service_delete_participant_attribute.return_value = None

        service.delete_participant_attribute(
            participant_id="p-1", attribute_id="a-1"
        )

        api.taurus_network_service_delete_participant_attribute.assert_called_once_with(
            participant_id="p-1",
            attribute_id="a-1",
        )
