"""Token metadata service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, Optional

from taurus_protect.mappers.token_metadata import (
    crypto_punk_metadata_from_dto,
    fa_token_metadata_from_dto,
    token_metadata_from_dto,
)
from taurus_protect.models.token_metadata import (
    CryptoPunkMetadata,
    FATokenMetadata,
    TokenMetadata,
)
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


class TokenMetadataService(BaseService):
    """
    Service for token metadata operations.

    Provides methods to get metadata for various token types including
    ERC tokens (ERC721, ERC1155), FA tokens (Tezos), and CryptoPunks.

    Example:
        >>> # Get ERC token metadata
        >>> metadata = client.token_metadata.get("ETH", "0x1234...")
        >>> print(f"Token: {metadata.name}")
        >>>
        >>> # Get with token ID for NFTs
        >>> metadata = client.token_metadata.get_erc(
        ...     network="mainnet",
        ...     contract_address="0x1234...",
        ...     token_id="42",
        ... )
    """

    def __init__(self, api_client: Any, token_metadata_api: Any) -> None:
        """
        Initialize token metadata service.

        Args:
            api_client: The OpenAPI client instance.
            token_metadata_api: The TokenMetadataApi service from OpenAPI client.
        """
        super().__init__(api_client)
        self._token_metadata_api = token_metadata_api

    def get(
        self,
        blockchain: str,
        contract_address: str,
        token_id: str = "0",
        network: str = "mainnet",
        with_data: bool = False,
    ) -> Optional[TokenMetadata]:
        """
        Get token metadata for an ERC token.

        Args:
            blockchain: The blockchain (e.g., "ETH").
            contract_address: The token contract address.
            token_id: The token ID (default "0" for fungible tokens).
            network: The network (default "mainnet").
            with_data: Whether to include base64 encoded data.

        Returns:
            The token metadata or None if not found.

        Raises:
            ValueError: If required arguments are invalid.
            APIError: If API request fails.
        """
        self._validate_required(blockchain, "blockchain")
        self._validate_required(contract_address, "contract_address")

        try:
            resp = self._token_metadata_api.token_metadata_service_get_evmerc_token_metadata(
                network=network,
                contract=contract_address,
                token=token_id,
                with_data=with_data if with_data else None,
                blockchain=blockchain,
            )

            result = getattr(resp, "result", None)
            return token_metadata_from_dto(result)
        except Exception as e:
            from taurus_protect.errors import APIError, NotFoundError

            if isinstance(e, (APIError, NotFoundError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_erc(
        self,
        network: str,
        contract_address: str,
        token_id: str,
        blockchain: Optional[str] = None,
        with_data: bool = False,
    ) -> Optional[TokenMetadata]:
        """
        Get ERC token metadata (ERC721 or ERC1155).

        Args:
            network: The network (e.g., "mainnet", "goerli").
            contract_address: The ERC721 or ERC1155 contract address.
            token_id: The token ID.
            blockchain: Optional blockchain override.
            with_data: Whether to include base64 encoded data and type.

        Returns:
            The token metadata or None if not found.

        Raises:
            ValueError: If required arguments are invalid.
            APIError: If API request fails.
        """
        self._validate_required(network, "network")
        self._validate_required(contract_address, "contract_address")
        self._validate_required(token_id, "token_id")

        try:
            resp = self._token_metadata_api.token_metadata_service_get_evmerc_token_metadata(
                network=network,
                contract=contract_address,
                token=token_id,
                with_data=with_data if with_data else None,
                blockchain=blockchain,
            )

            result = getattr(resp, "result", None)
            return token_metadata_from_dto(result)
        except Exception as e:
            from taurus_protect.errors import APIError, NotFoundError

            if isinstance(e, (APIError, NotFoundError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_fa(
        self,
        network: str,
        contract_address: str,
        token_id: str = "0",
        with_data: bool = False,
    ) -> Optional[FATokenMetadata]:
        """
        Get FA token metadata (Tezos FA1.2 or FA2).

        Args:
            network: The Tezos network (e.g., "mainnet").
            contract_address: The FA1.2 or FA2 contract address.
            token_id: The token ID. Must be "0" for FA1.2, any existing token for FA2.
            with_data: Whether to include base64 encoded data and type.

        Returns:
            The FA token metadata or None if not found.

        Raises:
            ValueError: If required arguments are invalid.
            APIError: If API request fails.
        """
        self._validate_required(network, "network")
        self._validate_required(contract_address, "contract_address")

        try:
            resp = self._token_metadata_api.token_metadata_service_get_fa_token_metadata(
                network=network,
                contract=contract_address,
                token=token_id,
                with_data=with_data if with_data else None,
            )

            result = getattr(resp, "result", None)
            return fa_token_metadata_from_dto(result)
        except Exception as e:
            from taurus_protect.errors import APIError, NotFoundError

            if isinstance(e, (APIError, NotFoundError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_crypto_punk(
        self,
        network: str,
        contract_address: str,
        punk_id: str,
        blockchain: Optional[str] = None,
    ) -> Optional[CryptoPunkMetadata]:
        """
        Get CryptoPunk token metadata.

        Args:
            network: The network (e.g., "mainnet").
            contract_address: The CryptoPunks contract address
                             (0xb47e3cd837dDF8e4c57F05d70Ab865de6e193BBB for ETH/mainnet).
            punk_id: The punk ID (0-9999).
            blockchain: Optional blockchain override.

        Returns:
            The CryptoPunk metadata or None if not found.

        Raises:
            ValueError: If required arguments are invalid.
            APIError: If API request fails.
        """
        self._validate_required(network, "network")
        self._validate_required(contract_address, "contract_address")
        self._validate_required(punk_id, "punk_id")

        try:
            resp = self._token_metadata_api.token_metadata_service_get_crypto_punks_token_metadata(
                network=network,
                contract=contract_address,
                token=punk_id,
                blockchain=blockchain,
            )

            result = getattr(resp, "result", None)
            return crypto_punk_metadata_from_dto(result)
        except Exception as e:
            from taurus_protect.errors import APIError, NotFoundError

            if isinstance(e, (APIError, NotFoundError, ValueError)):
                raise
            raise self._handle_error(e) from e
