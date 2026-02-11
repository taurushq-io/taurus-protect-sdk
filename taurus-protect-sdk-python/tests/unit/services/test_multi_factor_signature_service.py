"""Unit tests for MultiFactorSignatureService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.services.multi_factor_signature_service import (
    MultiFactorSignatureService,
)


class TestGetChallenge:
    """Tests for MultiFactorSignatureService.get_challenge()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        mfs_api = MagicMock()
        service = MultiFactorSignatureService(api_client=api_client, mfs_api=mfs_api)
        return service, mfs_api

    def test_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="challenge_id"):
            service.get_challenge(challenge_id="")

    def test_returns_challenge(self) -> None:
        service, api = self._make_service()
        dto = MagicMock()
        dto.id = "ch-1"
        dto.request_id = "42"
        dto.user_id = "u-1"
        dto.status = "PENDING"
        dto.challenge_type = "OTP"
        dto.challengeType = None
        dto.created_at = None
        dto.createdAt = None
        dto.expires_at = None
        dto.expiresAt = None
        reply = MagicMock()
        reply.result = dto
        api.multi_factor_signature_service_get_challenge.return_value = reply

        challenge = service.get_challenge("ch-1")

        assert challenge.id == "ch-1"
        assert challenge.status == "PENDING"

    def test_raises_not_found_when_none(self) -> None:
        service, api = self._make_service()
        reply = MagicMock()
        reply.result = None
        api.multi_factor_signature_service_get_challenge.return_value = reply

        from taurus_protect.errors import NotFoundError

        with pytest.raises(NotFoundError):
            service.get_challenge("ch-missing")


class TestListChallenges:
    """Tests for MultiFactorSignatureService.list_challenges()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        mfs_api = MagicMock()
        service = MultiFactorSignatureService(api_client=api_client, mfs_api=mfs_api)
        return service, mfs_api

    def test_raises_on_invalid_limit(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="limit must be positive"):
            service.list_challenges(limit=0)

    def test_raises_on_negative_offset(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="offset cannot be negative"):
            service.list_challenges(offset=-1)

    def test_returns_empty_when_no_results(self) -> None:
        service, api = self._make_service()
        reply = MagicMock()
        reply.result = None
        reply.total_items = None
        api.multi_factor_signature_service_get_challenges.return_value = reply

        challenges, pagination = service.list_challenges()

        assert challenges == []

    def test_passes_request_id_filter(self) -> None:
        service, api = self._make_service()
        reply = MagicMock()
        reply.result = None
        reply.total_items = None
        api.multi_factor_signature_service_get_challenges.return_value = reply

        service.list_challenges(request_id=42)

        api.multi_factor_signature_service_get_challenges.assert_called_once_with(
            request_id="42",
            limit="50",
            offset="0",
        )


class TestCreateChallenge:
    """Tests for MultiFactorSignatureService.create_challenge()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        mfs_api = MagicMock()
        service = MultiFactorSignatureService(api_client=api_client, mfs_api=mfs_api)
        return service, mfs_api

    def test_raises_on_invalid_request_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="request_id must be positive"):
            service.create_challenge(request_id=0, challenge_type="OTP")

    def test_raises_on_empty_type(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="challenge_type"):
            service.create_challenge(request_id=1, challenge_type="")


class TestVerifyChallenge:
    """Tests for MultiFactorSignatureService.verify_challenge()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        mfs_api = MagicMock()
        service = MultiFactorSignatureService(api_client=api_client, mfs_api=mfs_api)
        return service, mfs_api

    def test_raises_on_empty_challenge_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="challenge_id"):
            service.verify_challenge(challenge_id="", response="123456")

    def test_raises_on_empty_response(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="response"):
            service.verify_challenge(challenge_id="ch-1", response="")
