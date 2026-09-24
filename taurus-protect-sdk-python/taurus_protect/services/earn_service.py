"""Earn service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import Any, List, Optional, Tuple

from taurus_protect.errors import APIError
from taurus_protect.mappers._base import safe_string
from taurus_protect.models.earn import EarnReward
from taurus_protect.models.pagination import CursorPage, cursor_page, cursor_request
from taurus_protect.services._base import BaseService


def earn_reward_from_dto(dto: Any) -> Optional[EarnReward]:
    """Convert a generated ``TgvalidatordEarnReward`` to the domain model."""
    if dto is None:
        return None

    merkl = dto.merkl_token_reward
    token = merkl.token if merkl is not None else None
    reward_type = dto.reward_type
    return EarnReward(
        id=safe_string(dto.id),
        recipient_address_id=dto.recipient_address_id,
        recipient_address=dto.recipient_address,
        reward_type=reward_type.value if reward_type is not None else None,
        amount=merkl.amount if merkl is not None else None,
        claimed=merkl.claimed if merkl is not None else None,
        pending=merkl.pending if merkl is not None else None,
        token_address=token.address if token is not None else None,
        token_symbol=token.symbol if token is not None else None,
        token_asset_id=token.asset_id if token is not None else None,
    )


class EarnService(BaseService):
    """
    Service for earn rewards.

    Example:
        >>> cursor = None
        >>> while True:
        ...     rewards, page = client.earn.list_rewards(page_size=100, cursor=cursor)
        ...     for reward in rewards:
        ...         print(f"{reward.recipient_address}: {reward.amount} {reward.token_symbol}")
        ...     if not page.has_more:
        ...         break
        ...     cursor = page.next_cursor
    """

    def __init__(self, api_client: Any, earn_api: Any) -> None:
        """
        Initialize earn service.

        Args:
            api_client: The OpenAPI client instance.
            earn_api: The EarnApi service from OpenAPI client.
        """
        super().__init__(api_client)
        self._earn_api = earn_api

    def list_rewards(
        self,
        recipient_address_id: Optional[str] = None,
        page_size: Optional[int] = None,
        cursor: Optional[str] = None,
        *,
        current_page: Optional[str] = None,
        page_request: Optional[str] = None,
    ) -> Tuple[List[EarnReward], CursorPage]:
        """
        List earn rewards, one page at a time.

        Args:
            recipient_address_id: Filter by the receiving address ID.
            page_size: Page size (default 20, max 100).
            cursor: ``page.next_cursor`` from the previous page, to continue.
            current_page: Low-level page token; not with ``cursor``.
            page_request: Low-level page direction (FIRST, PREVIOUS, NEXT, LAST).

        Returns:
            Tuple of (rewards, page).

        Raises:
            ValueError: If the page size is invalid or cursor options conflict.
            APIError: If API request fails.
        """
        req = cursor_request(
            page_size, cursor, current_page=current_page, page_request=page_request
        )

        try:
            resp = self._earn_api.earn_service_get_rewards(
                recipient_address_id=recipient_address_id,
                **req.query_params(),
            )

            rewards = [
                r for dto in resp.rewards or [] if (r := earn_reward_from_dto(dto)) is not None
            ]
            return rewards, cursor_page(req.page_size, resp.cursor)
        except Exception as e:
            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
