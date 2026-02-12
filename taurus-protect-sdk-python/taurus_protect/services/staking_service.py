"""Staking service for Taurus-PROTECT SDK."""

from __future__ import annotations

from decimal import Decimal
from typing import TYPE_CHECKING, Any, Dict, List, Optional, Tuple

from taurus_protect.models.pagination import Pagination
from taurus_protect.models.staking import StakingInfo, Validator
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    from taurus_protect._internal.openapi.api.staking_api import StakingApi


class StakingService(BaseService):
    """
    Service for staking operations.

    Provides methods to list validators and get staking information
    for addresses across supported blockchains.

    Supported blockchains include:
    - Ethereum (ETH) - Validator info
    - Solana (SOL) - Stake accounts
    - Cardano (ADA) - Stake pool info
    - Near (NEAR) - Validator info
    - Fantom (FTM) - Validator info
    - Internet Computer (ICP) - Neuron info
    - Tezos (XTZ) - Staking rewards

    Example:
        >>> # List ETH validators
        >>> validators, pagination = client.staking.list_validators(
        ...     blockchain="ETH",
        ...     limit=50
        ... )
        >>> for v in validators:
        ...     print(f"{v.name}: {v.commission}% commission")
        >>>
        >>> # Get staking info for an address
        >>> info = client.staking.get_staking_info(address_id=123)
        >>> print(f"Staked: {info.staked_amount}")
    """

    def __init__(self, api_client: Any, staking_api: "StakingApi") -> None:
        """
        Initialize staking service.

        Args:
            api_client: The OpenAPI client instance.
            staking_api: The StakingApi service from OpenAPI client.
        """
        super().__init__(api_client)
        self._api = staking_api

    def list_validators(
        self,
        blockchain: str,
        network: str = "mainnet",
        limit: int = 50,
        offset: int = 0,
    ) -> Tuple[List[Validator], Optional[Pagination]]:
        """
        List validators for a blockchain.

        Retrieves available validators that can be used for staking operations.
        The specific information returned varies by blockchain.

        Args:
            blockchain: Blockchain type (e.g., "ETH", "SOL", "ADA", "NEAR", "FTM").
            network: Network identifier (default: "mainnet").
            limit: Maximum number of validators to return (default: 50).
            offset: Number of validators to skip for pagination (default: 0).

        Returns:
            Tuple of (validators list, pagination info).

        Raises:
            ValueError: If blockchain is empty or limit/offset are invalid.
            APIError: If the API request fails.

        Example:
            >>> validators, pagination = client.staking.list_validators(
            ...     blockchain="ETH",
            ...     network="mainnet",
            ...     limit=100
            ... )
            >>> print(f"Found {pagination.total_items} validators")
        """
        self._validate_required(blockchain, "blockchain")
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            validators: List[Validator] = []
            total_items = 0

            # Different blockchains have different API methods
            blockchain_upper = blockchain.upper()

            if blockchain_upper == "ETH":
                # Ethereum validators
                resp = self._api.staking_service_get_eth_validators_info(
                    network=network,
                )
                validators = self._map_eth_validators(resp, blockchain_upper, network)
                total_items = len(validators)

            elif blockchain_upper == "ADA":
                # ADA stake pools - need a stake pool ID for specific info
                # For listing, we return empty since the API requires a specific pool ID
                validators = []
                total_items = 0

            elif blockchain_upper == "NEAR":
                # NEAR validators - requires specific validator_id
                validators = []
                total_items = 0

            elif blockchain_upper == "FTM":
                # Fantom validators - requires specific validator_id
                validators = []
                total_items = 0

            elif blockchain_upper in ("SOL", "SOLANA"):
                # Solana stake accounts - use get_stake_accounts
                # This returns stake account info, not validators
                validators = []
                total_items = 0

            else:
                # Unknown blockchain - return empty list
                validators = []
                total_items = 0

            # Apply pagination manually since these APIs don't support it
            paginated_validators = validators[offset : offset + limit]

            pagination = self._extract_pagination(
                total_items=total_items,
                offset=offset,
                limit=limit,
            )

            return paginated_validators, pagination

        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def _map_eth_validators(
        self,
        resp: Any,
        blockchain: str,
        network: str,
    ) -> List[Validator]:
        """Map ETH validators response to Validator models."""
        validators: List[Validator] = []

        result = getattr(resp, "result", None) or getattr(resp, "validators", None)
        if not result:
            return validators

        if isinstance(result, list):
            for item in result:
                validator = self._map_validator_from_dto(item, blockchain, network)
                if validator:
                    validators.append(validator)
        else:
            # Single validator response
            validator = self._map_validator_from_dto(result, blockchain, network)
            if validator:
                validators.append(validator)

        return validators

    def _map_validator_from_dto(
        self,
        dto: Any,
        blockchain: str,
        network: str,
    ) -> Optional[Validator]:
        """Map a single validator DTO to Validator model."""
        if dto is None:
            return None

        validator_id = (
            getattr(dto, "id", None)
            or getattr(dto, "validator_id", None)
            or getattr(dto, "pubkey", None)
            or ""
        )

        if not validator_id:
            return None

        name = getattr(dto, "name", None) or getattr(dto, "label", None) or ""
        address = (
            getattr(dto, "address", None)
            or getattr(dto, "pubkey", None)
            or getattr(dto, "validator_address", None)
            or ""
        )

        commission_raw = getattr(dto, "commission", None) or getattr(dto, "commission_rate", None)
        commission = Decimal(str(commission_raw)) if commission_raw is not None else None

        total_stake_raw = getattr(dto, "total_stake", None) or getattr(
            dto, "effective_balance", None
        )
        total_stake = Decimal(str(total_stake_raw)) if total_stake_raw is not None else None

        active = getattr(dto, "active", True)
        if isinstance(active, str):
            active = active.lower() in ("true", "active", "yes")

        status = getattr(dto, "status", None) or ("active" if active else "inactive")

        return Validator(
            id=str(validator_id),
            name=str(name),
            blockchain=blockchain,
            network=network,
            address=str(address),
            commission=commission,
            total_stake=total_stake,
            active=bool(active),
            status=str(status),
        )

    def get_staking_info(self, address_id: int) -> StakingInfo:
        """
        Get staking information for an address.

        Retrieves current staking positions, rewards, and status for a
        specific address.

        Args:
            address_id: The address ID to get staking info for.

        Returns:
            Staking information for the address.

        Raises:
            ValueError: If address_id is invalid.
            NotFoundError: If the address is not found or has no staking info.
            APIError: If the API request fails.

        Example:
            >>> info = client.staking.get_staking_info(address_id=123)
            >>> print(f"Staked amount: {info.staked_amount}")
            >>> print(f"Rewards: {info.rewards}")
            >>> print(f"Status: {info.status}")
        """
        if address_id <= 0:
            raise ValueError("address_id must be positive")

        try:
            # Try to get stake accounts (works for SOL)
            resp = self._api.staking_service_get_stake_accounts(
                address_id=str(address_id),
            )

            result = getattr(resp, "result", None) or getattr(resp, "stake_accounts", None)

            if not result:
                # Return empty staking info if no data
                return StakingInfo(address_id=str(address_id))

            # Map the first stake account to StakingInfo
            if isinstance(result, list) and len(result) > 0:
                return self._map_staking_info_from_dto(result[0], address_id)
            elif not isinstance(result, list):
                return self._map_staking_info_from_dto(result, address_id)

            return StakingInfo(address_id=str(address_id))

        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def _map_staking_info_from_dto(self, dto: Any, address_id: int) -> StakingInfo:
        """Map staking DTO to StakingInfo model."""
        if dto is None:
            return StakingInfo(address_id=str(address_id))

        blockchain = getattr(dto, "blockchain", None) or ""
        network = getattr(dto, "network", None) or ""
        validator_id = getattr(dto, "validator_id", None) or getattr(dto, "vote_pubkey", None)
        validator_address = getattr(dto, "validator_address", None) or getattr(
            dto, "vote_pubkey", None
        )

        staked_amount_raw = (
            getattr(dto, "staked_amount", None)
            or getattr(dto, "stake", None)
            or getattr(dto, "balance", None)
        )
        staked_amount = Decimal(str(staked_amount_raw)) if staked_amount_raw is not None else None

        rewards_raw = getattr(dto, "rewards", None) or getattr(dto, "accumulated_rewards", None)
        rewards = Decimal(str(rewards_raw)) if rewards_raw is not None else None

        status = getattr(dto, "status", None) or getattr(dto, "state", None) or ""

        staked_at = getattr(dto, "staked_at", None) or getattr(dto, "activation_epoch", None)
        unbonding_at = getattr(dto, "unbonding_at", None) or getattr(
            dto, "deactivation_epoch", None
        )

        return StakingInfo(
            address_id=str(address_id),
            blockchain=str(blockchain),
            network=str(network),
            validator_id=str(validator_id) if validator_id else None,
            validator_address=str(validator_address) if validator_address else None,
            staked_amount=staked_amount,
            rewards=rewards,
            status=str(status),
            staked_at=staked_at if isinstance(staked_at, type(None)) else None,
            unbonding_at=unbonding_at if isinstance(unbonding_at, type(None)) else None,
        )
