"""Business rules service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional

from taurus_protect.mappers.business_rule import business_rules_from_dto
from taurus_protect.models.business_rule import BusinessRule, BusinessRuleResult
from taurus_protect.models.pagination import cursor_page, cursor_request
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    from taurus_protect._internal.openapi.api.business_rules_api import BusinessRulesApi


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
        *,
        cursor: Optional[str] = None,
        ids: Optional[List[str]] = None,
        rule_groups: Optional[List[str]] = None,
        address_ids: Optional[List[str]] = None,
        level: Optional[str] = None,
    ) -> BusinessRuleResult:
        """List business rules, one page at a time (v2 API).

        Args:
            page_size: Page size (default 20, max 100).
            current_page: Low-level page token; not with ``cursor``.
            page_request: Low-level page direction (FIRST, NEXT, PREVIOUS, LAST).
            rule_keys: Filter by rule keys.
            wallet_ids: Filter by wallet IDs.
            currency_ids: Filter by currency IDs.
            entity_type: Filter by entity type.
            entity_ids: Filter by entity IDs.
            cursor: ``result.page.next_cursor`` from the previous page, to continue.
            ids: Filter by rule IDs.
            rule_groups: Filter by rule groups.
            address_ids: Filter by address IDs.
            level: Filter by rule level.

        Returns:
            BusinessRuleResult with the rules and the page.

        Raises:
            ValueError: If the page size is invalid or cursor options conflict.
            APIError: If API request fails.
        """
        req = cursor_request(
            page_size, cursor, current_page=current_page, page_request=page_request
        )

        try:
            reply = self._api.rule_service_get_business_rules_v2(
                ids=ids,
                rule_keys=rule_keys,
                rule_groups=rule_groups,
                wallet_ids=wallet_ids,
                currency_ids=currency_ids,
                address_ids=address_ids,
                level=level,
                entity_type=entity_type,
                entity_ids=entity_ids,
                **req.query_params(),
            )

            return BusinessRuleResult(
                rules=business_rules_from_dto(reply.result or []),
                page=cursor_page(req.page_size, reply.cursor),
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
        *,
        cursor: Optional[str] = None,
    ) -> BusinessRuleResult:
        """List business rules for a specific wallet, one page at a time.

        Args:
            wallet_id: The wallet ID (must be positive).
            page_size: Page size (default 20, max 100).
            current_page: Low-level page token; not with ``cursor``.
            page_request: Low-level page direction.
            cursor: ``result.page.next_cursor`` from the previous page, to continue.

        Returns:
            BusinessRuleResult with the rules and the page.

        Raises:
            ValueError: If wallet_id is not positive or paging options are invalid.
            APIError: If API request fails.
        """
        if wallet_id <= 0:
            raise ValueError("wallet_id must be positive")
        return self.list(
            wallet_ids=[str(wallet_id)],
            page_size=page_size,
            current_page=current_page,
            page_request=page_request,
            cursor=cursor,
        )

    def list_by_currency(
        self,
        currency_id: str,
        page_size: Optional[int] = None,
        current_page: Optional[str] = None,
        page_request: Optional[str] = None,
        *,
        cursor: Optional[str] = None,
    ) -> BusinessRuleResult:
        """List business rules for a specific currency, one page at a time.

        Args:
            currency_id: The currency ID (must not be empty).
            page_size: Page size (default 20, max 100).
            current_page: Low-level page token; not with ``cursor``.
            page_request: Low-level page direction.
            cursor: ``result.page.next_cursor`` from the previous page, to continue.

        Returns:
            BusinessRuleResult with the rules and the page.

        Raises:
            ValueError: If currency_id is empty or paging options are invalid.
            APIError: If API request fails.
        """
        self._validate_required(currency_id, "currency_id")
        return self.list(
            currency_ids=[currency_id],
            page_size=page_size,
            current_page=current_page,
            page_request=page_request,
            cursor=cursor,
        )

    def update_transactions_enabled(self, enabled: bool) -> None:
        """Enable or disable transaction processing for the tenant.

        This is the transactions-enabled business rule — the kill switch that
        tg-protect-mcpd drives — so it must exist in every SDK. Peer of Go
        BusinessRuleService.UpdateTransactionsEnabled, Java updateTransactionsEnabled
        and TS updateTransactionsEnabled.

        Args:
            enabled: True to allow transactions, False to halt them.

        Raises:
            APIError: If the API request fails.
        """
        from taurus_protect._internal.openapi.models.tgvalidatord_update_transactions_enabled_business_rule_request import (  # noqa: E501
            TgvalidatordUpdateTransactionsEnabledBusinessRuleRequest,
        )

        body = TgvalidatordUpdateTransactionsEnabledBusinessRuleRequest(enabled=enabled)

        try:
            self._api.rule_service_update_transactions_enabled_business_rule(body=body)
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
