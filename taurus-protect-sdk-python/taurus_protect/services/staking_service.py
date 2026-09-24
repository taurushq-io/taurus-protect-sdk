"""Staking service for Taurus-PROTECT SDK."""

from __future__ import annotations

from decimal import Decimal
from typing import TYPE_CHECKING, Any, List, Optional

from taurus_protect.models.pagination import cursor_request
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
        >>> validators = client.staking.list_validators(blockchain="ETH")
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
        *,
        ids: Optional[List[str]] = None,
    ) -> List[Validator]:
        """
        List validators for a blockchain: every validator the endpoint returns.

        Only ETH has a validators list, and it does not page; every other chain returns
        an empty list.

        Args:
            blockchain: Blockchain type (e.g., "ETH").
            network: Network identifier (default: "mainnet").
            ids: Keep only these ETH validator IDs.

        Returns:
            Every matching validator.

        Raises:
            ValueError: If blockchain is empty.
            APIError: If the API request fails.

        Example:
            >>> validators = client.staking.list_validators(blockchain="ETH", network="mainnet")
            >>> print(f"Found {len(validators)} validators")
        """
        self._validate_required(blockchain, "blockchain")

        blockchain_upper = blockchain.upper()
        if blockchain_upper != "ETH":
            # ADA, NEAR, FTM and ICP are looked up one pool, validator or neuron at a
            # time, and SOL has stake accounts; none has a list of validators.
            return []

        try:
            resp = self._api.staking_service_get_eth_validators_info(
                network=network, ids=ids or None
            )
            return self._map_eth_validators(resp, blockchain_upper, network)
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
        for item in resp.validators or []:
            validator = self._map_validator_from_dto(item, blockchain, network)
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
            # Stake accounts (SOL); only the first one is read, so ask for one.
            resp = self._api.staking_service_get_stake_accounts(
                address_id=str(address_id),
                **cursor_request(1).query_params(),
            )

            accounts = resp.stake_accounts or []
            if not accounts:
                return StakingInfo(address_id=str(address_id))
            return self._map_staking_info_from_dto(accounts[0], address_id)

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
