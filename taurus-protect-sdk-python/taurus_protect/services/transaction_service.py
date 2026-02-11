"""Transaction service for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import datetime
from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect._internal.openapi.exceptions import ApiException
from taurus_protect.errors import NotFoundError
from taurus_protect.mappers.transaction import map_transaction, map_transactions
from taurus_protect.models.pagination import Pagination
from taurus_protect.models.transaction import Transaction
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    from taurus_protect._internal.openapi.api.transactions_api import TransactionsApi


class TransactionService(BaseService):
    """
    Service for retrieving and exporting blockchain transactions.

    Transactions represent the movement of cryptocurrency on the blockchain
    and can be either incoming (received) or outgoing (sent).

    Example:
        >>> # Get recent transactions
        >>> transactions, _ = client.transactions.list(currency="ETH", limit=50)
        >>> for tx in transactions:
        ...     print(f"{tx.tx_hash}: {tx.amount} {tx.currency}")
        ...
        >>> # Get a transaction by hash
        >>> tx = client.transactions.get_by_hash("0x1234...")
    """

    def __init__(
        self,
        api_client: Any,
        transactions_api: "TransactionsApi",
    ) -> None:
        """
        Initialize the transaction service.

        Args:
            api_client: The OpenAPI client instance.
            transactions_api: The transactions API instance.
        """
        super().__init__(api_client)
        self._api = transactions_api

    def get(self, transaction_id: int) -> Transaction:
        """
        Get a single transaction by ID.

        Args:
            transaction_id: The transaction ID.

        Returns:
            The transaction.

        Raises:
            NotFoundError: If the transaction is not found.
            APIError: If the API call fails.
        """
        if transaction_id <= 0:
            raise ValueError("transaction_id must be positive")

        try:
            reply = self._api.transaction_service_get_transactions(
                currency=None,
                direction=None,
                query=None,
                limit="1",
                offset="0",
                var_from=None,
                to=None,
                transaction_ids=None,
                type=None,
                source=None,
                destination=None,
                ids=[str(transaction_id)],
                blockchain=None,
                network=None,
                from_block_number=None,
                to_block_number=None,
                hashes=None,
                address=None,
                amount_above=None,
                exclude_unknown_source_destination=None,
                customer_id=None,
            )

            result = reply.result
            if not result:
                raise NotFoundError(f"Transaction with id '{transaction_id}' not found")

            return map_transaction(result[0])
        except NotFoundError:
            raise
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

    def get_by_hash(self, tx_hash: str) -> Transaction:
        """
        Get a transaction by its blockchain hash.

        Args:
            tx_hash: The transaction hash.

        Returns:
            The transaction.

        Raises:
            NotFoundError: If the transaction is not found.
            APIError: If the API call fails.
        """
        self._validate_required(tx_hash, "tx_hash")

        try:
            reply = self._api.transaction_service_get_transactions(
                currency=None,
                direction=None,
                query=None,
                limit="1",
                offset="0",
                var_from=None,
                to=None,
                transaction_ids=None,
                type=None,
                source=None,
                destination=None,
                ids=None,
                blockchain=None,
                network=None,
                from_block_number=None,
                to_block_number=None,
                hashes=[tx_hash],
                address=None,
                amount_above=None,
                exclude_unknown_source_destination=None,
                customer_id=None,
            )

            result = reply.result
            if not result:
                raise NotFoundError(f"Transaction with hash '{tx_hash}' not found")

            return map_transaction(result[0])
        except NotFoundError:
            raise
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

    def list(
        self,
        from_date: Optional[datetime] = None,
        to_date: Optional[datetime] = None,
        currency: Optional[str] = None,
        direction: Optional[str] = None,
        limit: int = 50,
        offset: int = 0,
    ) -> Tuple[List[Transaction], Optional[Pagination]]:
        """
        List transactions with filtering.

        Args:
            from_date: Filter transactions after this date.
            to_date: Filter transactions before this date.
            currency: Filter by currency ID or symbol.
            direction: Filter by direction ("incoming" or "outgoing").
            limit: Maximum number of transactions to return.
            offset: Offset for pagination.

        Returns:
            Tuple of (transactions list, pagination info).

        Raises:
            APIError: If the API call fails.
        """
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            reply = self._api.transaction_service_get_transactions(
                currency=currency,
                direction=direction,
                query=None,
                limit=str(limit),
                offset=str(offset),
                var_from=from_date,
                to=to_date,
                transaction_ids=None,
                type=None,
                source=None,
                destination=None,
                ids=None,
                blockchain=None,
                network=None,
                from_block_number=None,
                to_block_number=None,
                hashes=None,
                address=None,
                amount_above=None,
                exclude_unknown_source_destination=None,
                customer_id=None,
            )

            transactions = map_transactions(reply.result)
            pagination = self._extract_pagination(
                getattr(reply, "total_items", None),
                offset,
                limit,
            )
            return transactions, pagination
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

    def list_by_address(
        self,
        address: str,
        limit: int = 50,
        offset: int = 0,
    ) -> Tuple[List[Transaction], Optional[Pagination]]:
        """
        List transactions for a specific blockchain address.

        Args:
            address: The blockchain address.
            limit: Maximum number of transactions to return.
            offset: Offset for pagination.

        Returns:
            Tuple of (transactions list, pagination info).

        Raises:
            APIError: If the API call fails.
        """
        self._validate_required(address, "address")
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            reply = self._api.transaction_service_get_transactions(
                currency=None,
                direction=None,
                query=None,
                limit=str(limit),
                offset=str(offset),
                var_from=None,
                to=None,
                transaction_ids=None,
                type=None,
                source=None,
                destination=None,
                ids=None,
                blockchain=None,
                network=None,
                from_block_number=None,
                to_block_number=None,
                hashes=None,
                address=address,
                amount_above=None,
                exclude_unknown_source_destination=None,
                customer_id=None,
            )

            transactions = map_transactions(reply.result)
            pagination = self._extract_pagination(
                getattr(reply, "total_items", None),
                offset,
                limit,
            )
            return transactions, pagination
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

    def export_csv(
        self,
        from_date: Optional[datetime] = None,
        to_date: Optional[datetime] = None,
        currency: Optional[str] = None,
        direction: Optional[str] = None,
        limit: int = 1000,
        offset: int = 0,
    ) -> str:
        """
        Export transactions to CSV format.

        Args:
            from_date: Filter transactions after this date.
            to_date: Filter transactions before this date.
            currency: Filter by currency ID or symbol.
            direction: Filter by direction ("incoming" or "outgoing").
            limit: Maximum number of transactions to export.
            offset: Offset for pagination.

        Returns:
            CSV content as a string.

        Raises:
            APIError: If the API call fails.
        """
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            reply = self._api.transaction_service_export_transactions(
                currency=currency,
                direction=direction,
                query=None,
                limit=str(limit),
                offset=str(offset),
                var_from=from_date,
                to=to_date,
                transaction_ids=None,
                format="csv",
                type=None,
                source=None,
                destination=None,
                ids=None,
                blockchain=None,
                network=None,
                from_block_number=None,
                to_block_number=None,
                amount_above=None,
                exclude_unknown_source_destination=None,
                hashes=None,
                address=None,
            )

            return reply.result or ""
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise
