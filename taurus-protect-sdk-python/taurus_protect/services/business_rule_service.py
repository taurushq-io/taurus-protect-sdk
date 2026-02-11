"""Business rules service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional

from taurus_protect.mappers.business_rule import business_rule_from_dto, business_rules_from_dto
from taurus_protect.models.business_rule import BusinessRule, BusinessRuleResult
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    from taurus_protect._internal.openapi.api.business_rules_api import BusinessRulesApi


def _extract_cursor(cursor: Any) -> tuple:
    """Extract current_page and has_next from a response cursor."""
    if cursor is None:
        return None, False
    cp = getattr(cursor, "current_page", None)
    if cp is None:
        cp = getattr(cursor, "currentPage", None)
    if cp is not None and not isinstance(cp, str):
        cp = None
    hn = getattr(cursor, "has_next", None)
    if hn is None:
        hn = getattr(cursor, "hasNext", None)
    return cp, bool(hn) if hn is not None else False


class BusinessRuleService(BaseService):
    """Service for managing business rules.

    Business rules define operational constraints and configurations at various
    scopes (tenant, wallet, address, currency).
    """

    def __init__(
        self,
        api_client: Any,
        business_rules_api: "BusinessRulesApi",
    ) -> None:
        super().__init__(api_client)
        self._api = business_rules_api

    def list(
        self,
        page_size: Optional[int] = None,
        current_page: Optional[str] = None,
        page_request: Optional[str] = None,
        rule_keys: Optional[List[str]] = None,
        wallet_ids: Optional[List[str]] = None,
        currency_ids: Optional[List[str]] = None,
        entity_type: Optional[str] = None,
        entity_ids: Optional[List[str]] = None,
    ) -> BusinessRuleResult:
        """List business rules with cursor-based pagination (v2 API).

        Args:
            page_size: Number of rules per page.
            current_page: Cursor for current page.
            page_request: Page request type (FIRST, NEXT, PREVIOUS, LAST).
            rule_keys: Filter by rule keys.
            wallet_ids: Filter by wallet IDs.
            currency_ids: Filter by currency IDs.
            entity_type: Filter by entity type.
            entity_ids: Filter by entity IDs.

        Returns:
            BusinessRuleResult with rules list and cursor info.

        Raises:
            APIError: If API request fails.
        """
        try:
            reply = self._api.rule_service_get_business_rules_v2(
                ids=None,
                rule_keys=rule_keys,
                rule_groups=None,
                wallet_ids=wallet_ids,
                currency_ids=currency_ids,
                address_ids=None,
                level=None,
                cursor_current_page=current_page,
                cursor_page_request=page_request or "FIRST",
                cursor_page_size=str(page_size) if page_size else None,
                entity_type=entity_type,
                entity_ids=entity_ids,
            )

            rules_list = getattr(reply, "result", None)
            rules = business_rules_from_dto(rules_list) if rules_list else []

            cursor = getattr(reply, "cursor", None)
            result_current_page, has_next = _extract_cursor(cursor)

            return BusinessRuleResult(
                rules=rules,
                current_page=result_current_page,
                has_next=has_next,
            )
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list_by_wallet(
        self,
        wallet_id: int,
        page_size: Optional[int] = None,
        current_page: Optional[str] = None,
        page_request: Optional[str] = None,
    ) -> BusinessRuleResult:
        """List business rules for a specific wallet.

        Args:
            wallet_id: The wallet ID (must be positive).
            page_size: Number of rules per page.
            current_page: Cursor for current page.
            page_request: Page request type.

        Returns:
            BusinessRuleResult with rules list and cursor info.

        Raises:
            ValueError: If wallet_id is not positive.
            APIError: If API request fails.
        """
        if wallet_id <= 0:
            raise ValueError("wallet_id must be positive")
        return self.list(
            wallet_ids=[str(wallet_id)],
            page_size=page_size,
            current_page=current_page,
            page_request=page_request,
        )

    def list_by_currency(
        self,
        currency_id: str,
        page_size: Optional[int] = None,
        current_page: Optional[str] = None,
        page_request: Optional[str] = None,
    ) -> BusinessRuleResult:
        """List business rules for a specific currency.

        Args:
            currency_id: The currency ID (must not be empty).
            page_size: Number of rules per page.
            current_page: Cursor for current page.
            page_request: Page request type.

        Returns:
            BusinessRuleResult with rules list and cursor info.

        Raises:
            ValueError: If currency_id is empty.
            APIError: If API request fails.
        """
        self._validate_required(currency_id, "currency_id")
        return self.list(
            currency_ids=[currency_id],
            page_size=page_size,
            current_page=current_page,
            page_request=page_request,
        )
