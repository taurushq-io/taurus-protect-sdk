"""Fee payer service for Taurus-PROTECT SDK."""

from __future__ import annotations

from decimal import Decimal
from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect.models.pagination import Pagination
from taurus_protect.models.staking import FeePayer
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    from taurus_protect._internal.openapi.api.fee_payers_api import FeePayersApi


class FeePayerService(BaseService):
    """
    Service for fee payer management.

    Fee payers are accounts used to pay transaction fees on behalf of
    other accounts. This is commonly used for:

    - Gas station networks (meta-transactions)
    - Sponsored transactions
    - Account abstraction patterns

    Example:
        >>> # List all fee payers
        >>> fee_payers, pagination = client.fee_payers.list(limit=50)
        >>> for fp in fee_payers:
        ...     print(f"{fp.id}: {fp.address} ({fp.balance})")
        >>>
        >>> # Get a specific fee payer
        >>> fee_payer = client.fee_payers.get(fee_payer_id="fp-123")
        >>> print(f"Balance: {fee_payer.balance}")
    """

    def __init__(self, api_client: Any, fee_payers_api: "FeePayersApi") -> None:
        """
        Initialize fee payer service.

        Args:
            api_client: The OpenAPI client instance.
            fee_payers_api: The FeePayersApi service from OpenAPI client.
        """
        super().__init__(api_client)
        self._api = fee_payers_api

    def list(
        self,
        limit: int = 50,
        offset: int = 0,
        blockchain: Optional[str] = None,
        network: Optional[str] = None,
    ) -> Tuple[List[FeePayer], Optional[Pagination]]:
        """
        List fee payers with pagination.

        Args:
            limit: Maximum number of fee payers to return (default: 50).
            offset: Number of fee payers to skip for pagination (default: 0).
            blockchain: Optional filter by blockchain type.
            network: Optional filter by network identifier.

        Returns:
            Tuple of (fee payers list, pagination info).

        Raises:
            ValueError: If limit or offset are invalid.
            APIError: If the API request fails.

        Example:
            >>> # List all fee payers
            >>> fee_payers, pagination = client.fee_payers.list()
            >>> print(f"Total: {pagination.total_items}")
            >>>
            >>> # Filter by blockchain
            >>> fee_payers, _ = client.fee_payers.list(blockchain="SOL")
        """
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            resp = self._api.fee_payer_service_get_fee_payers(
                limit=str(limit),
                offset=str(offset),
                ids=None,
                blockchain=blockchain,
                network=network,
            )

            result = getattr(resp, "result", None) or getattr(resp, "fee_payers", None)
            if not result:
                return [], None

            fee_payers: List[FeePayer] = []
            if isinstance(result, list):
                for dto in result:
                    fee_payer = self._map_fee_payer(dto)
                    if fee_payer:
                        fee_payers.append(fee_payer)
            else:
                fee_payer = self._map_fee_payer(result)
                if fee_payer:
                    fee_payers.append(fee_payer)

            pagination = self._extract_pagination(
                total_items=getattr(resp, "total_items", None),
                offset=getattr(resp, "offset", None),
                limit=limit,
            )

            return fee_payers, pagination

        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get(self, fee_payer_id: str) -> FeePayer:
        """
        Get a fee payer by ID.

        Args:
            fee_payer_id: The fee payer ID to retrieve.

        Returns:
            The fee payer.

        Raises:
            ValueError: If fee_payer_id is empty.
            NotFoundError: If the fee payer is not found.
            APIError: If the API request fails.

        Example:
            >>> fee_payer = client.fee_payers.get(fee_payer_id="fp-123")
            >>> print(f"Address: {fee_payer.address}")
            >>> print(f"Balance: {fee_payer.balance}")
        """
        self._validate_required(fee_payer_id, "fee_payer_id")

        try:
            resp = self._api.fee_payer_service_get_fee_payer(id=fee_payer_id)

            result = getattr(resp, "result", None) or getattr(resp, "fee_payer", resp)

            if result is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Fee payer {fee_payer_id} not found")

            fee_payer = self._map_fee_payer(result)
            if fee_payer is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Fee payer {fee_payer_id} not found")

            return fee_payer

        except Exception as e:
            from taurus_protect.errors import APIError, NotFoundError

            if isinstance(e, (APIError, ValueError, NotFoundError)):
                raise
            raise self._handle_error(e) from e

    def _map_fee_payer(self, dto: Any) -> Optional[FeePayer]:
        """Map fee payer DTO to FeePayer model."""
        if dto is None:
            return None

        fee_payer_id = getattr(dto, "id", None)
        if not fee_payer_id:
            return None

        blockchain = getattr(dto, "blockchain", None) or ""
        network = getattr(dto, "network", None) or ""
        address = (
            getattr(dto, "address", None)
            or getattr(dto, "public_key", None)
            or getattr(dto, "pubkey", None)
            or ""
        )

        balance_raw = getattr(dto, "balance", None)
        balance = Decimal(str(balance_raw)) if balance_raw is not None else None

        status = getattr(dto, "status", None) or "active"
        if isinstance(status, bool):
            status = "active" if status else "disabled"

        created_at = getattr(dto, "created_at", None)
        updated_at = getattr(dto, "updated_at", None)

        return FeePayer(
            id=str(fee_payer_id),
            blockchain=str(blockchain),
            network=str(network),
            address=str(address),
            balance=balance,
            status=str(status),
            created_at=created_at,
            updated_at=updated_at,
        )
