"""Statistics, Price, and Score mappers for converting OpenAPI DTOs to domain models."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional

from taurus_protect.mappers._base import safe_datetime, safe_string
from taurus_protect.models.statistics import (
    PortfolioStatistics,
    Price,
    PriceHistoryPoint,
    Score,
)

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


def price_from_dto(dto: Any) -> Optional[Price]:
    """
    Convert OpenAPI TgvalidatordCurrencyPrice to domain Price.

    Args:
        dto: OpenAPI currency price DTO.

    Returns:
        Domain Price model or None if dto is None.
    """
    if dto is None:
        return None

    return Price(
        currency_from=safe_string(getattr(dto, "currency_from", None)),
        currency_to=safe_string(getattr(dto, "currency_to", None)),
        rate=safe_string(getattr(dto, "rate", None)),
        blockchain=getattr(dto, "blockchain", None),
        decimals=getattr(dto, "decimals", None),
        change_percent_24h=getattr(dto, "change_percent24_hour", None),
        source=getattr(dto, "source", None),
        created_at=safe_datetime(getattr(dto, "creation_date", None)),
        updated_at=safe_datetime(getattr(dto, "update_date", None)),
    )


def prices_from_dto(dtos: Optional[List[Any]]) -> List[Price]:
    """
    Convert list of OpenAPI currency price DTOs to domain Prices.

    Args:
        dtos: List of OpenAPI currency price DTOs.

    Returns:
        List of domain Price models.
    """
    if dtos is None:
        return []
    return [p for dto in dtos if (p := price_from_dto(dto)) is not None]


def price_history_point_from_dto(dto: Any) -> Optional[PriceHistoryPoint]:
    """
    Convert OpenAPI TgvalidatordPricesHistoryPoint to domain PriceHistoryPoint.

    Args:
        dto: OpenAPI price history point DTO.

    Returns:
        Domain PriceHistoryPoint model or None if dto is None.
    """
    if dto is None:
        return None

    return PriceHistoryPoint(
        timestamp=safe_datetime(getattr(dto, "timestamp", None)),
        rate=safe_string(getattr(dto, "rate", None)),
    )


def price_history_from_dto(dtos: Optional[List[Any]]) -> List[PriceHistoryPoint]:
    """
    Convert list of OpenAPI price history DTOs to domain PriceHistoryPoints.

    Args:
        dtos: List of OpenAPI price history point DTOs.

    Returns:
        List of domain PriceHistoryPoint models.
    """
    if dtos is None:
        return []
    return [p for dto in dtos if (p := price_history_point_from_dto(dto)) is not None]


def portfolio_statistics_from_dto(dto: Any) -> Optional[PortfolioStatistics]:
    """
    Convert OpenAPI TgvalidatordAggregatedStatsData to domain PortfolioStatistics.

    Args:
        dto: OpenAPI aggregated stats DTO.

    Returns:
        Domain PortfolioStatistics model or None if dto is None.
    """
    if dto is None:
        return None

    return PortfolioStatistics(
        addresses_count=safe_string(getattr(dto, "addresses_count", None)),
        wallets_count=safe_string(getattr(dto, "wallets_count", None)),
        total_balance=safe_string(getattr(dto, "total_balance", None)),
        total_balance_base_currency=safe_string(getattr(dto, "total_balance_base_currency", None)),
        avg_balance_per_address=safe_string(getattr(dto, "avg_balance_per_address", None)),
    )


def score_from_dto(dto: Any) -> Optional[Score]:
    """
    Convert OpenAPI TgvalidatordScore to domain Score.

    Args:
        dto: OpenAPI score DTO.

    Returns:
        Domain Score model or None if dto is None.
    """
    if dto is None:
        return None

    return Score(
        id=safe_string(getattr(dto, "id", None)),
        provider=safe_string(getattr(dto, "provider", None)),
        score_type=safe_string(getattr(dto, "type", None)),
        score=safe_string(getattr(dto, "score", None)),
        updated_at=safe_datetime(getattr(dto, "update_date", None)),
    )


def scores_from_dto(dtos: Optional[List[Any]]) -> List[Score]:
    """
    Convert list of OpenAPI score DTOs to domain Scores.

    Args:
        dtos: List of OpenAPI score DTOs.

    Returns:
        List of domain Score models.
    """
    if dtos is None:
        return []
    return [s for dto in dtos if (s := score_from_dto(dto)) is not None]
