"""Transaction mapper for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import Any, List, Optional

from taurus_protect.mappers._base import safe_datetime, safe_int, safe_string
from taurus_protect.models.transaction import Transaction


def map_transaction(dto: Any) -> Transaction:
    """Map OpenAPI transaction DTO to domain model."""
    return Transaction(
        id=safe_string(getattr(dto, "id", None)) or "",
        request_id=safe_string(getattr(dto, "request_id", None)),
        wallet_id=safe_string(getattr(dto, "wallet_id", None)),
        address_id=safe_string(getattr(dto, "address_id", None)),
        currency=safe_string(getattr(dto, "currency", None)),
        blockchain=safe_string(getattr(dto, "blockchain", None)),
        tx_hash=safe_string(getattr(dto, "tx_hash", None)) or safe_string(getattr(dto, "hash", None)),
        block_height=safe_int(getattr(dto, "block_height", None)) or safe_int(getattr(dto, "block_number", None)),
        block_hash=safe_string(getattr(dto, "block_hash", None)),
        amount=safe_string(getattr(dto, "amount", None)),
        fee=safe_string(getattr(dto, "fee", None)),
        direction=safe_string(getattr(dto, "direction", None)),
        status=safe_string(getattr(dto, "status", None)),
        confirmations=safe_int(getattr(dto, "confirmations", None)) or 0,
        created_at=safe_datetime(getattr(dto, "created_at", None)) or safe_datetime(getattr(dto, "creation_date", None)),
        confirmed_at=safe_datetime(getattr(dto, "confirmed_at", None)) or safe_datetime(getattr(dto, "confirmation_date", None)),
    )


def map_transactions(dtos: Optional[List[Any]]) -> List[Transaction]:
    """Map list of OpenAPI transaction DTOs to domain models."""
    if not dtos:
        return []
    return [map_transaction(dto) for dto in dtos]
