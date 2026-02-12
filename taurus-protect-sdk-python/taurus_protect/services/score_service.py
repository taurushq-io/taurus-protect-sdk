"""Score service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional

from taurus_protect.mappers.statistics import scores_from_dto
from taurus_protect.models.statistics import Score
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


class ScoreService(BaseService):
    """
    Service for address and transaction risk scoring.

    Provides methods to retrieve risk scores from compliance/AML providers
    for addresses and transactions.

    Example:
        >>> # Get risk score for an address
        >>> scores = client.scores.get_address_score(address_id="123")
        >>> for score in scores:
        ...     print(f"{score.provider}: {score.score}")
        >>>
        >>> # Refresh scores for a whitelisted address
        >>> scores = client.scores.refresh_whitelisted_address_score(
        ...     address_id="456",
        ...     provider="chainalysis",
        ... )
    """

    def __init__(self, api_client: Any, scores_api: Any) -> None:
        """
        Initialize score service.

        Args:
            api_client: The OpenAPI client instance.
            scores_api: The ScoresApi service from OpenAPI client.
        """
        super().__init__(api_client)
        self._scores_api = scores_api

    def get_address_score(
        self,
        address_id: str,
        provider: Optional[str] = None,
    ) -> List[Score]:
        """
        Get risk score for an address.

        Retrieves risk scores for a specific address from compliance/AML
        providers. If no provider is specified, returns scores from all
        available providers.

        Args:
            address_id: The address ID to get scores for.
            provider: Optional provider name to filter results
                     (e.g., "chainalysis", "elliptic").

        Returns:
            List of risk scores for the address.

        Raises:
            ValueError: If address_id is empty.
            APIError: If API request fails.

        Example:
            >>> scores = client.scores.get_address_score(address_id="123")
            >>> for score in scores:
            ...     print(f"{score.provider}: {score.score}")
        """
        self._validate_required(address_id, "address_id")

        try:
            # Build request body
            body = {}
            if provider:
                body["provider"] = provider

            resp = self._scores_api.score_service_refresh_address_score(
                address_id=address_id,
                body=body,
            )

            scores_list = getattr(resp, "scores", None)
            return scores_from_dto(scores_list) if scores_list else []
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_transaction_score(
        self,
        tx_hash: str,
        provider: Optional[str] = None,
    ) -> List[Score]:
        """
        Get risk score for a transaction.

        Retrieves risk scores for a specific transaction from compliance/AML
        providers. If no provider is specified, returns scores from all
        available providers.

        Note: This method may not be directly supported by the current API.
        Transaction scoring is typically done through address-based scoring
        where the transaction's sender/receiver addresses are evaluated.

        Args:
            tx_hash: The transaction hash to get scores for.
            provider: Optional provider name to filter results
                     (e.g., "chainalysis", "elliptic").

        Returns:
            List of risk scores for the transaction.

        Raises:
            ValueError: If tx_hash is empty.
            NotImplementedError: If transaction scoring is not supported.
            APIError: If API request fails.

        Example:
            >>> scores = client.scores.get_transaction_score(
            ...     tx_hash="0x123...",
            ...     provider="chainalysis",
            ... )
        """
        self._validate_required(tx_hash, "tx_hash")

        # The current API does not have a direct transaction scoring endpoint.
        # Transaction scoring is typically done by:
        # 1. Looking up the transaction
        # 2. Scoring the sender/receiver addresses
        #
        # This method is provided as a placeholder for future API support
        # or for custom implementations.

        raise NotImplementedError(
            "Transaction scoring is not directly supported by the API. "
            "Use get_address_score() to score the transaction's addresses instead."
        )

    def refresh_whitelisted_address_score(
        self,
        address_id: str,
        provider: Optional[str] = None,
    ) -> List[Score]:
        """
        Refresh risk score for a whitelisted address.

        Triggers a refresh of risk scores for a whitelisted address from
        compliance/AML providers. This is useful for getting updated
        scores for addresses on the whitelist.

        Args:
            address_id: The whitelisted address ID to refresh scores for.
            provider: Optional provider name to filter which provider to use
                     (e.g., "chainalysis", "elliptic").

        Returns:
            List of refreshed risk scores for the whitelisted address.

        Raises:
            ValueError: If address_id is empty.
            APIError: If API request fails.

        Example:
            >>> scores = client.scores.refresh_whitelisted_address_score(
            ...     address_id="456",
            ...     provider="chainalysis",
            ... )
            >>> for score in scores:
            ...     print(f"{score.provider}: {score.score} (updated: {score.updated_at})")
        """
        self._validate_required(address_id, "address_id")

        try:
            # Build request body
            body = {}
            if provider:
                body["provider"] = provider

            resp = self._scores_api.score_service_refresh_wla_score(
                address_id=address_id,
                body=body,
            )

            scores_list = getattr(resp, "scores", None)
            return scores_from_dto(scores_list) if scores_list else []
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
