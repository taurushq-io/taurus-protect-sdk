"""Transaction service for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import datetime
from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect._internal.openapi.exceptions import ApiException
from taurus_protect.errors import NotFoundError
from taurus_protect.mappers.transaction import map_transaction, map_transactions
from taurus_protect.models.pagination import (
    PLUS_ROWS,
    Pagination,
    offset_pagination,
    offset_query,
    parse_count,
    resolve_offset,
    resolve_page_size,
)
from taurus_protect.models.transaction import Transaction, TransactionExport
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
        >>> transactions, pagination = client.transactions.list(currency="ETH", limit=100)
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
        limit: Optional[int] = None,
        offset: Optional[int] = None,
        blockchain: Optional[str] = None,
        network: Optional[str] = None,
    ) -> Tuple[List[Transaction], Pagination]:
        """
        List transactions with filtering, one page at a time.

        Args:
            from_date: Filter transactions after this date.
            to_date: Filter transactions before this date.
            currency: Filter by currency ID or symbol.
            direction: Filter by direction ("incoming" or "outgoing").
            limit: Page size (default 20, max 100).
            offset: Number of transactions to skip; pass ``pagination.next_offset``.
            blockchain: Filter by blockchain (e.g. "ETH").
            network: Filter by network (e.g. "mainnet").

        Returns:
            Tuple of (transactions, pagination). The server's total can be an upper
            bound, so a trailing page may come back empty with ``has_more`` false.

        Raises:
            ValueError: If limit or offset are invalid.
            APIError: If the API call fails.
        """
        return self._list(
            limit,
            offset,
            currency=currency,
            direction=direction,
            var_from=from_date,
            to=to_date,
            blockchain=blockchain,
            network=network,
        )

    def list_by_address(
        self,
        address: str,
        limit: Optional[int] = None,
        offset: Optional[int] = None,
    ) -> Tuple[List[Transaction], Pagination]:
        """
        List transactions for a specific blockchain address, one page at a time.

        Args:
            address: The blockchain address.
            limit: Page size (default 20, max 100).
            offset: Number of transactions to skip; pass ``pagination.next_offset``.

        Returns:
            Tuple of (transactions, pagination).

        Raises:
            ValueError: If address is empty or limit/offset are invalid.
            APIError: If the API call fails.
        """
        self._validate_required(address, "address")
        return self._list(limit, offset, address=address)

    def _list(
        self, limit: Optional[int], offset: Optional[int], **filters: Any
    ) -> Tuple[List[Transaction], Pagination]:
        page_size = resolve_page_size(limit, "limit")
        start = resolve_offset(offset)

        try:
            reply = self._api.transaction_service_get_transactions(
                **offset_query(page_size, start), **filters
            )

            rows = reply.result or []
            transactions = map_transactions(rows)
            pagination = offset_pagination(
                PLUS_ROWS,
                limit=page_size,
                offset=start,
                served_rows=len(rows),
                total_items=reply.total_items,
                excluded=len(rows) - len(transactions),
            )
            return transactions, pagination
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

    def export(
        self,
        from_date: Optional[datetime] = None,
        to_date: Optional[datetime] = None,
        currency: Optional[str] = None,
        direction: Optional[str] = None,
        limit: Optional[int] = None,
        blockchain: Optional[str] = None,
        network: Optional[str] = None,
        format: Optional[str] = None,
    ) -> TransactionExport:
        """
        Export transactions in one reply.

        The export cannot page: the server always exports from the first matching row,
        so there is no offset. When ``total_items`` exceeds the rows exported, raise the
        limit or narrow the filters.

        Args:
            from_date: Filter transactions after this date.
            to_date: Filter transactions before this date.
            currency: Filter by currency ID or symbol.
            direction: Filter by direction ("incoming" or "outgoing").
            limit: Maximum rows to export (default 20, no SDK maximum).
            blockchain: Filter by blockchain (e.g. "ETH").
            network: Filter by network (e.g. "mainnet").
            format: Export format ("json", "csv", "csv_simple"); the server defaults
                to JSON.

        Returns:
            The exported text and the server's count of matching transactions.

        Raises:
            ValueError: If limit is negative.
            APIError: If the API call fails.
        """
        size = resolve_page_size(limit, "limit", maximum=None)

        try:
            reply = self._api.transaction_service_export_transactions(
                currency=currency,
                direction=direction,
                var_from=from_date,
                to=to_date,
                format=format,
                blockchain=blockchain,
                network=network,
                limit=str(size),
            )

            return TransactionExport(
                content=reply.result or "",
                total_items=parse_count(reply.total_items, "totalItems"),
            )
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
        limit: Optional[int] = None,
        blockchain: Optional[str] = None,
        network: Optional[str] = None,
    ) -> TransactionExport:
        """
        Export transactions as CSV: :meth:`export` with ``format="csv"``.

        Args:
            from_date: Filter transactions after this date.
            to_date: Filter transactions before this date.
            currency: Filter by currency ID or symbol.
            direction: Filter by direction ("incoming" or "outgoing").
            limit: Maximum rows to export (default 20, no SDK maximum).
            blockchain: Filter by blockchain (e.g. "ETH").
            network: Filter by network (e.g. "mainnet").

        Returns:
            The CSV text and the server's count of matching transactions.

        Raises:
            ValueError: If limit is negative.
            APIError: If the API call fails.
        """
        return self.export(
            from_date=from_date,
            to_date=to_date,
            currency=currency,
            direction=direction,
            limit=limit,
            blockchain=blockchain,
            network=network,
            format="csv",
        )
