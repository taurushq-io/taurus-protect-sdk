"""Blockchain service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional

from taurus_protect.mappers._base import safe_bool, safe_int, safe_string
from taurus_protect.models.blockchain import Blockchain
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


def blockchain_from_dto(dto: Any) -> Optional[Blockchain]:
    """
    Convert an OpenAPI blockchain DTO to domain model.

    Args:
        dto: The OpenAPI DTO object.

    Returns:
        Blockchain model or None if dto is None.
    """
    if dto is None:
        return None

    # Generate ID from blockchain and network
    blockchain_name = safe_string(getattr(dto, "blockchain", None) or getattr(dto, "name", None))
    network_name = safe_string(getattr(dto, "network", None))
    blockchain_id = f"{blockchain_name}_{network_name}" if network_name else blockchain_name

    return Blockchain(
        id=getattr(dto, "id", None) or blockchain_id,
        name=blockchain_name,
        network=network_name,
        display_name=getattr(dto, "display_name", None) or getattr(dto, "displayName", None),
        enabled=safe_bool(getattr(dto, "enabled", True)),
        native_currency=getattr(dto, "native_currency", None)
        or getattr(dto, "nativeCurrency", None),
        block_height=safe_int(
            getattr(dto, "block_height", None) or getattr(dto, "blockHeight", None)
        )
        or None,
        block_time=safe_int(getattr(dto, "block_time", None) or getattr(dto, "blockTime", None))
        or None,
        confirmations_required=safe_int(
            getattr(dto, "confirmations_required", None)
            or getattr(dto, "confirmationsRequired", None)
        )
        or None,
    )


def blockchains_from_dto(dtos: Any) -> List[Blockchain]:
    """
    Convert a list of OpenAPI blockchain DTOs to domain models.

    Args:
        dtos: List of OpenAPI DTO objects.

    Returns:
        List of Blockchain models.
    """
    if dtos is None:
        return []
    return [b for dto in dtos if (b := blockchain_from_dto(dto)) is not None]


class BlockchainService(BaseService):
    """
    Service for blockchain operations.

    Provides methods to list supported blockchains and retrieve
    blockchain-specific information.

    Example:
        >>> # List all blockchains
        >>> blockchains = client.blockchains.list()
        >>> for bc in blockchains:
        ...     print(f"{bc.name} ({bc.network}): {bc.native_currency}")
        >>>
        >>> # Get specific blockchain info
        >>> blockchain = client.blockchains.get("ETH", "mainnet")
        >>> print(f"Block height: {blockchain.block_height}")
    """

    def __init__(self, api_client: Any, blockchain_api: Any) -> None:
        """
        Initialize blockchain service.

        Args:
            api_client: The OpenAPI client instance.
            blockchain_api: The BlockchainAPI service from OpenAPI client.
        """
        super().__init__(api_client)
        self._blockchain_api = blockchain_api

    def list(
        self,
        blockchain: Optional[str] = None,
        network: Optional[str] = None,
        include_block_height: bool = False,
    ) -> List[Blockchain]:
        """
        List supported blockchains.

        Args:
            blockchain: Optional filter by blockchain name.
            network: Optional filter by network name.
            include_block_height: Whether to include current block height.
                                  Requires blockchain/network to be specified.

        Returns:
            List of supported blockchains.

        Raises:
            APIError: If API request fails.
        """
        try:
            resp = self._blockchain_api.blockchain_service_get_blockchains(
                blockchain=blockchain,
                network=network,
                include_block_height=include_block_height if include_block_height else None,
            )

            result = getattr(resp, "result", None) or getattr(resp, "blockchains", None)
            return blockchains_from_dto(result) if result else []
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def get(
        self,
        blockchain: str,
        network: str = "mainnet",
        include_block_height: bool = False,
    ) -> Blockchain:
        """
        Get blockchain information.

        Args:
            blockchain: The blockchain name (e.g., "ETH", "BTC").
            network: The network name (default: "mainnet").
            include_block_height: Whether to include current block height.

        Returns:
            The blockchain information.

        Raises:
            ValueError: If blockchain is invalid.
            NotFoundError: If blockchain not found.
            APIError: If API request fails.
        """
        self._validate_required(blockchain, "blockchain")
        self._validate_required(network, "network")

        try:
            resp = self._blockchain_api.blockchain_service_get_blockchains(
                blockchain=blockchain,
                network=network,
                include_block_height=include_block_height if include_block_height else None,
            )

            result = getattr(resp, "result", None) or getattr(resp, "blockchains", None)
            blockchains = blockchains_from_dto(result) if result else []

            if not blockchains:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Blockchain {blockchain}/{network} not found")

            # Return the first matching blockchain
            return blockchains[0]
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_by_id(self, blockchain_id: str) -> Blockchain:
        """
        Get blockchain by composite ID.

        The ID is expected in format "BLOCKCHAIN_NETWORK" (e.g., "ETH_mainnet").

        Args:
            blockchain_id: The blockchain ID.

        Returns:
            The blockchain information.

        Raises:
            ValueError: If blockchain_id is invalid.
            NotFoundError: If blockchain not found.
            APIError: If API request fails.
        """
        self._validate_required(blockchain_id, "blockchain_id")

        # Parse the ID
        parts = blockchain_id.split("_", 1)
        blockchain = parts[0]
        network = parts[1] if len(parts) > 1 else "mainnet"

        return self.get(blockchain=blockchain, network=network)
