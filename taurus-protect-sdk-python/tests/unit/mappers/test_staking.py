"""Unit tests for staking service mapper functions."""

from decimal import Decimal
from types import SimpleNamespace

from taurus_protect.services.staking_service import StakingService


class TestMapValidatorFromDto:
    """Tests for StakingService._map_validator_from_dto."""

    def _make_service(self) -> StakingService:
        """Create a StakingService instance without calling __init__."""
        return StakingService.__new__(StakingService)

    def test_maps_all_fields(self) -> None:
        svc = self._make_service()
        dto = SimpleNamespace(
            id="val-1",
            validator_id=None,
            pubkey=None,
            name="Validator One",
            label=None,
            address="0xabc",
            validator_address=None,
            commission="5.0",
            commission_rate=None,
            total_stake="1000000",
            effective_balance=None,
            active=True,
            status="active",
        )
        result = svc._map_validator_from_dto(dto, "ETH", "mainnet")
        assert result is not None
        assert result.id == "val-1"
        assert result.name == "Validator One"
        assert result.blockchain == "ETH"
        assert result.network == "mainnet"
        assert result.address == "0xabc"
        assert result.commission == Decimal("5.0")
        assert result.total_stake == Decimal("1000000")
        assert result.active is True
        assert result.status == "active"

    def test_returns_none_for_none(self) -> None:
        svc = self._make_service()
        assert svc._map_validator_from_dto(None, "ETH", "mainnet") is None

    def test_returns_none_for_empty_id(self) -> None:
        svc = self._make_service()
        dto = SimpleNamespace(
            id=None, validator_id=None, pubkey=None,
            name=None, label=None, address=None, validator_address=None,
            commission=None, commission_rate=None,
            total_stake=None, effective_balance=None,
            active=True, status=None,
        )
        result = svc._map_validator_from_dto(dto, "ETH", "mainnet")
        assert result is None

    def test_id_fallback_to_pubkey(self) -> None:
        svc = self._make_service()
        dto = SimpleNamespace(
            id=None, validator_id=None, pubkey="0xpubkey",
            name=None, label="My Val", address=None, validator_address=None,
            commission=None, commission_rate="1.5",
            total_stake=None, effective_balance="500",
            active="true", status=None,
        )
        result = svc._map_validator_from_dto(dto, "SOL", "testnet")
        assert result is not None
        assert result.id == "0xpubkey"
        assert result.name == "My Val"
        assert result.commission == Decimal("1.5")
        assert result.total_stake == Decimal("500")
        assert result.active is True

    def test_active_string_parsing(self) -> None:
        svc = self._make_service()
        dto = SimpleNamespace(
            id="v-1", validator_id=None, pubkey=None,
            name="V", label=None, address="addr", validator_address=None,
            commission=None, commission_rate=None,
            total_stake=None, effective_balance=None,
            active="inactive", status=None,
        )
        result = svc._map_validator_from_dto(dto, "ETH", "mainnet")
        assert result is not None
        assert result.active is False
        assert result.status == "inactive"


class TestMapStakingInfoFromDto:
    """Tests for StakingService._map_staking_info_from_dto."""

    def _make_service(self) -> StakingService:
        return StakingService.__new__(StakingService)

    def test_maps_basic_fields(self) -> None:
        svc = self._make_service()
        dto = SimpleNamespace(
            blockchain="SOL",
            network="mainnet",
            validator_id="vote-123",
            vote_pubkey=None,
            validator_address="val-addr",
            staked_amount="1000",
            stake=None,
            balance=None,
            rewards="50",
            accumulated_rewards=None,
            status="staked",
            state=None,
            staked_at=None,
            activation_epoch=None,
            unbonding_at=None,
            deactivation_epoch=None,
        )
        result = svc._map_staking_info_from_dto(dto, 42)
        assert result.address_id == "42"
        assert result.blockchain == "SOL"
        assert result.network == "mainnet"
        assert result.validator_id == "vote-123"
        assert result.staked_amount == Decimal("1000")
        assert result.rewards == Decimal("50")
        assert result.status == "staked"

    def test_returns_empty_for_none_dto(self) -> None:
        svc = self._make_service()
        result = svc._map_staking_info_from_dto(None, 99)
        assert result.address_id == "99"
        assert result.staked_amount is None
