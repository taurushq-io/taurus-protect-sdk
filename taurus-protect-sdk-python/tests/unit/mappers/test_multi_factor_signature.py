"""Unit tests for multi-factor signature mapper functions."""

from datetime import datetime, timezone
from types import SimpleNamespace

from taurus_protect.services.multi_factor_signature_service import (
    MultiFactorSignatureChallenge,
    MultiFactorSignatureService,
)


class TestMapChallengeFromDto:
    """Tests for MultiFactorSignatureService._map_challenge_from_dto."""

    def test_maps_all_fields(self) -> None:
        created = datetime(2024, 6, 15, 12, 0, 0, tzinfo=timezone.utc)
        expires = datetime(2024, 6, 15, 12, 5, 0, tzinfo=timezone.utc)
        dto = SimpleNamespace(
            id="ch-1",
            request_id="req-100",
            user_id="u-1",
            status="PENDING",
            challenge_type="TOTP",
            challengeType=None,
            created_at=created,
            createdAt=None,
            expires_at=expires,
            expiresAt=None,
        )
        result = MultiFactorSignatureService._map_challenge_from_dto(dto)
        assert result.id == "ch-1"
        assert result.request_id == "req-100"
        assert result.user_id == "u-1"
        assert result.status == "PENDING"
        assert result.challenge_type == "TOTP"
        assert result.created_at == created
        assert result.expires_at == expires

    def test_handles_camelcase_fields(self) -> None:
        dto = SimpleNamespace(
            id="ch-2",
            request_id=None,
            user_id=None,
            status="VERIFIED",
            challenge_type=None,
            challengeType="SMS",
            created_at=None,
            createdAt="2024-01-01",
            expires_at=None,
            expiresAt="2024-01-02",
        )
        result = MultiFactorSignatureService._map_challenge_from_dto(dto)
        assert result.id == "ch-2"
        assert result.challenge_type == "SMS"
        assert result.created_at == "2024-01-01"
        assert result.expires_at == "2024-01-02"

    def test_handles_none_optional_fields(self) -> None:
        dto = SimpleNamespace(
            id="ch-3",
            request_id=None,
            user_id=None,
            status=None,
            challenge_type=None,
            challengeType=None,
            created_at=None,
            createdAt=None,
            expires_at=None,
            expiresAt=None,
        )
        result = MultiFactorSignatureService._map_challenge_from_dto(dto)
        assert result.id == "ch-3"
        assert result.request_id is None
        assert result.user_id is None
        assert result.status is None
        assert result.challenge_type is None

    def test_result_is_challenge_instance(self) -> None:
        dto = SimpleNamespace(
            id="ch-4",
            request_id="42",
            user_id="u-5",
            status="EXPIRED",
            challenge_type="EMAIL",
            challengeType=None,
            created_at=None,
            createdAt=None,
            expires_at=None,
            expiresAt=None,
        )
        result = MultiFactorSignatureService._map_challenge_from_dto(dto)
        assert isinstance(result, MultiFactorSignatureChallenge)
