"""Multi-factor signature service for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import datetime
from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect._internal.openapi.exceptions import ApiException
from taurus_protect.models.pagination import Pagination
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    from taurus_protect._internal.openapi.api.multi_factor_signature_api import (
        MultiFactorSignatureApi,
    )


class MultiFactorSignatureChallenge:
    """A multi-factor signature challenge."""

    def __init__(
        self,
        id: str,
        request_id: Optional[str] = None,
        user_id: Optional[str] = None,
        status: Optional[str] = None,
        challenge_type: Optional[str] = None,
        created_at: Optional[datetime] = None,
        expires_at: Optional[datetime] = None,
    ):
        self.id = id
        self.request_id = request_id
        self.user_id = user_id
        self.status = status
        self.challenge_type = challenge_type
        self.created_at = created_at
        self.expires_at = expires_at


class MultiFactorSignatureService(BaseService):
    """
    Service for multi-factor signature operations.

    Multi-factor signatures provide additional security for
    high-value or sensitive transactions by requiring multiple
    authentication factors.
    """

    def __init__(
        self,
        api_client: Any,
        mfs_api: "MultiFactorSignatureApi",
    ) -> None:
        super().__init__(api_client)
        self._api = mfs_api

    def get_challenge(self, challenge_id: str) -> MultiFactorSignatureChallenge:
        """Get a multi-factor signature challenge by ID."""
        self._validate_required(challenge_id, "challenge_id")

        try:
            reply = self._api.multi_factor_signature_service_get_challenge(challenge_id)
            result = reply.result
            if result is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Challenge {challenge_id} not found")
            return self._map_challenge_from_dto(result)
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

    def list_challenges(
        self,
        request_id: Optional[int] = None,
        limit: int = 50,
        offset: int = 0,
    ) -> Tuple[List[MultiFactorSignatureChallenge], Optional[Pagination]]:
        """List multi-factor signature challenges."""
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            reply = self._api.multi_factor_signature_service_get_challenges(
                request_id=str(request_id) if request_id else None,
                limit=str(limit),
                offset=str(offset),
            )

            challenges: List[MultiFactorSignatureChallenge] = []
            if reply.result:
                for dto in reply.result:
                    challenges.append(self._map_challenge_from_dto(dto))

            pagination = self._extract_pagination(
                getattr(reply, "total_items", None),
                offset,
                limit,
            )
            return challenges, pagination
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

    def create_challenge(
        self,
        request_id: int,
        challenge_type: str,
    ) -> str:
        """Create a new multi-factor signature challenge."""
        if request_id <= 0:
            raise ValueError("request_id must be positive")
        self._validate_required(challenge_type, "challenge_type")

        try:
            from taurus_protect._internal.openapi.models.tgvalidatord_create_mfs_challenge_request import (
                TgvalidatordCreateMfsChallengeRequest,
            )

            body = TgvalidatordCreateMfsChallengeRequest(
                request_id=str(request_id),
                challenge_type=challenge_type,
            )

            reply = self._api.multi_factor_signature_service_create_challenge(body=body)
            return reply.result or ""
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

    def verify_challenge(
        self,
        challenge_id: str,
        response: str,
    ) -> bool:
        """Verify a multi-factor signature challenge response."""
        self._validate_required(challenge_id, "challenge_id")
        self._validate_required(response, "response")

        try:
            from taurus_protect._internal.openapi.models.tgvalidatord_verify_mfs_challenge_request import (
                TgvalidatordVerifyMfsChallengeRequest,
            )

            body = TgvalidatordVerifyMfsChallengeRequest(
                challenge_id=challenge_id,
                response=response,
            )

            reply = self._api.multi_factor_signature_service_verify_challenge(body=body)
            return getattr(reply, "verified", False)
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

    @staticmethod
    def _map_challenge_from_dto(dto: Any) -> MultiFactorSignatureChallenge:
        return MultiFactorSignatureChallenge(
            id=str(getattr(dto, "id", "")),
            request_id=(
                str(getattr(dto, "request_id", "")) if getattr(dto, "request_id", None) else None
            ),
            user_id=str(getattr(dto, "user_id", "")) if getattr(dto, "user_id", None) else None,
            status=getattr(dto, "status", None),
            challenge_type=getattr(dto, "challenge_type", None)
            or getattr(dto, "challengeType", None),
            created_at=getattr(dto, "created_at", None) or getattr(dto, "createdAt", None),
            expires_at=getattr(dto, "expires_at", None) or getattr(dto, "expiresAt", None),
        )
