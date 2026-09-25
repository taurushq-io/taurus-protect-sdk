"""Business rule models for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import List, Optional

from pydantic import BaseModel, Field

from taurus_protect.models.currency import Currency
from taurus_protect.models.pagination import CursorPage


class BusinessRule(BaseModel):
    """Business rule for transaction validation and operational constraints."""

    id: str = Field(default="", description="Unique rule identifier")
    tenant_id: int = Field(default=0, description="Tenant ID")
    currency: Optional[str] = Field(default=None, description="Currency code")
    wallet_id: Optional[str] = Field(default=None, description="Wallet ID")
    address_id: Optional[str] = Field(default=None, description="Address ID")
    rule_key: Optional[str] = Field(default=None, description="Rule key/name")
    rule_value: Optional[str] = Field(default=None, description="Rule value/setting")
    rule_group: Optional[str] = Field(default=None, description="Rule group/category")
    rule_description: Optional[str] = Field(default=None, description="Human-readable description")
    rule_validation: Optional[str] = Field(default=None, description="Validation pattern")
    entity_type: Optional[str] = Field(
        default=None, description="Entity type (global, currency, wallet, address)"
    )
    entity_id: Optional[str] = Field(default=None, description="Entity ID")
    currency_info: Optional[Currency] = Field(default=None, description="Currency metadata")

    model_config = {"frozen": True}


class BusinessRuleResult(BaseModel):
    """
    One page of business rules.

    Continue with ``cursor=result.page.next_cursor`` while ``result.page.has_more``.
    """

    rules: List[BusinessRule] = Field(default_factory=list)
    page: CursorPage = Field(default_factory=CursorPage)

    model_config = {"frozen": True}
