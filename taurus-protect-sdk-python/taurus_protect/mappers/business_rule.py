"""Business rule mapper for converting OpenAPI DTOs to domain models."""

from __future__ import annotations

from typing import Any, List, Optional

from taurus_protect.mappers._base import safe_int, safe_string
from taurus_protect.mappers.currency import currency_from_dto
from taurus_protect.models.business_rule import BusinessRule


def _resolve(dto: Any, snake_name: str, camel_name: str) -> Any:
    """Resolve a field value trying snake_case first, then camelCase."""
    val = getattr(dto, snake_name, None)
    if val is not None:
        return val
    return getattr(dto, camel_name, None)


def business_rule_from_dto(dto: Any) -> Optional[BusinessRule]:
    """Convert OpenAPI TgvalidatordBusinessRule DTO to domain BusinessRule.

    Maps tenantId string -> int (matching Java BusinessRuleMapper).
    """
    if dto is None:
        return None

    raw_id = getattr(dto, "id", None)
    return BusinessRule(
        id=str(raw_id) if raw_id is not None else "",
        tenant_id=safe_int(_resolve(dto, "tenant_id", "tenantId")),
        currency=getattr(dto, "currency", None),
        wallet_id=_resolve(dto, "wallet_id", "walletId"),
        address_id=_resolve(dto, "address_id", "addressId"),
        rule_key=_resolve(dto, "rule_key", "ruleKey"),
        rule_value=_resolve(dto, "rule_value", "ruleValue"),
        rule_group=_resolve(dto, "rule_group", "ruleGroup"),
        rule_description=_resolve(dto, "rule_description", "ruleDescription"),
        rule_validation=_resolve(dto, "rule_validation", "ruleValidation"),
        entity_type=_resolve(dto, "entity_type", "entityType"),
        entity_id=_resolve(dto, "entity_id", "entityID"),
        currency_info=currency_from_dto(_resolve(dto, "currency_info", "currencyInfo")),
    )


def business_rules_from_dto(dtos: Optional[List[Any]]) -> List[BusinessRule]:
    """Convert list of OpenAPI business rule DTOs to domain BusinessRules."""
    if dtos is None:
        return []
    return [r for dto in dtos if (r := business_rule_from_dto(dto)) is not None]
