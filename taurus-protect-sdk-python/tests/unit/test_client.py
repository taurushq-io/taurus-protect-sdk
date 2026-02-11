"""Comprehensive unit tests for ProtectClient."""

from __future__ import annotations

import threading
from typing import TYPE_CHECKING, List

import pytest
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric import ec

from taurus_protect import ProtectClient
from taurus_protect.errors import ConfigurationError

if TYPE_CHECKING:
    pass


# =============================================================================
# TestProtectClientCreate - Tests for client creation (15 tests)
# =============================================================================


class TestProtectClientCreate:
    """Tests for ProtectClient.create() factory method."""

    def test_create_with_valid_config_succeeds(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that valid configuration creates a client."""
        client = ProtectClient.create(
            host=host,
            api_key=api_key,
            api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        assert client is not None
        assert client.is_closed is False
        client.close()

    def test_create_with_empty_host_raises(
        self, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that empty host raises ConfigurationError."""
        with pytest.raises(ConfigurationError, match="host cannot be empty"):
            ProtectClient.create(
                host="", api_key=api_key, api_secret=api_secret_hex,
                super_admin_keys_pem=super_admin_keys_pem,
            )

    def test_create_with_none_host_raises(
        self, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that None host raises ConfigurationError."""
        with pytest.raises(ConfigurationError, match="host cannot be empty"):
            ProtectClient.create(
                host=None, api_key=api_key, api_secret=api_secret_hex,  # type: ignore
                super_admin_keys_pem=super_admin_keys_pem,
            )

    def test_create_with_empty_api_key_raises(
        self, host: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that empty api_key raises ConfigurationError."""
        with pytest.raises(ConfigurationError, match="api_key cannot be empty"):
            ProtectClient.create(
                host=host, api_key="", api_secret=api_secret_hex,
                super_admin_keys_pem=super_admin_keys_pem,
            )

    def test_create_with_none_api_key_raises(
        self, host: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that None api_key raises ConfigurationError."""
        with pytest.raises(ConfigurationError, match="api_key cannot be empty"):
            ProtectClient.create(
                host=host, api_key=None, api_secret=api_secret_hex,  # type: ignore
                super_admin_keys_pem=super_admin_keys_pem,
            )

    def test_create_with_empty_api_secret_raises(
        self, host: str, api_key: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that empty api_secret raises ConfigurationError."""
        with pytest.raises(ConfigurationError, match="api_secret cannot be empty"):
            ProtectClient.create(
                host=host, api_key=api_key, api_secret="",
                super_admin_keys_pem=super_admin_keys_pem,
            )

    def test_create_with_none_api_secret_raises(
        self, host: str, api_key: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that None api_secret raises ConfigurationError."""
        with pytest.raises(ConfigurationError, match="api_secret cannot be empty"):
            ProtectClient.create(
                host=host, api_key=api_key, api_secret=None,  # type: ignore
                super_admin_keys_pem=super_admin_keys_pem,
            )

    def test_create_with_invalid_hex_api_secret_raises(
        self, host: str, api_key: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that invalid hex api_secret raises ConfigurationError."""
        # Non-hex characters
        with pytest.raises(ConfigurationError):
            ProtectClient.create(
                host=host, api_key=api_key, api_secret="not-valid-hex-string!",
                super_admin_keys_pem=super_admin_keys_pem,
            )

    def test_create_with_super_admin_keys_pem(
        self, host: str, api_key: str, api_secret_hex: str, ecdsa_public_key_pem: str
    ) -> None:
        """Test creating client with SuperAdmin keys."""
        client = ProtectClient.create(
            host=host,
            api_key=api_key,
            api_secret=api_secret_hex,
            super_admin_keys_pem=[ecdsa_public_key_pem],
            min_valid_signatures=1,
        )
        assert client is not None
        assert client.is_closed is False
        client.close()

    def test_create_with_invalid_pem_raises(
        self, host: str, api_key: str, api_secret_hex: str
    ) -> None:
        """Test that invalid PEM format raises ConfigurationError."""
        with pytest.raises(ConfigurationError, match="Invalid SuperAdmin key"):
            ProtectClient.create(
                host=host,
                api_key=api_key,
                api_secret=api_secret_hex,
                super_admin_keys_pem=["not-a-valid-pem-key"],
            )

    def test_create_with_negative_min_valid_signatures_raises(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that negative min_valid_signatures raises ConfigurationError."""
        with pytest.raises(ConfigurationError, match="min_valid_signatures cannot be negative"):
            ProtectClient.create(
                host=host,
                api_key=api_key,
                api_secret=api_secret_hex,
                super_admin_keys_pem=super_admin_keys_pem,
                min_valid_signatures=-1,
            )

    def test_create_with_min_valid_signatures_exceeding_keys_raises(
        self, host: str, api_key: str, api_secret_hex: str, ecdsa_public_key_pem: str
    ) -> None:
        """Test that min_valid_signatures > len(super_admin_keys) raises ConfigurationError."""
        with pytest.raises(ConfigurationError, match="cannot exceed"):
            ProtectClient.create(
                host=host,
                api_key=api_key,
                api_secret=api_secret_hex,
                super_admin_keys_pem=[ecdsa_public_key_pem],
                min_valid_signatures=5,
            )

    def test_create_normalizes_trailing_slash(
        self, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that trailing slash is removed from host."""
        client = ProtectClient.create(
            host="https://api.test.taurushq.com/",
            api_key=api_key,
            api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        assert client.host == "https://api.test.taurushq.com"
        client.close()

    def test_create_without_super_admin_keys_raises(
        self, host: str, api_key: str, api_secret_hex: str
    ) -> None:
        """Test that omitting super_admin_keys_pem raises ConfigurationError."""
        with pytest.raises(ConfigurationError, match="super_admin_keys_pem is required"):
            ProtectClient.create(
                host=host,
                api_key=api_key,
                api_secret=api_secret_hex,
                super_admin_keys_pem=None,  # type: ignore
            )

    def test_create_with_empty_super_admin_keys_raises(
        self, host: str, api_key: str, api_secret_hex: str
    ) -> None:
        """Test that empty super_admin_keys_pem list raises ConfigurationError."""
        with pytest.raises(ConfigurationError, match="super_admin_keys_pem is required"):
            ProtectClient.create(
                host=host,
                api_key=api_key,
                api_secret=api_secret_hex,
                super_admin_keys_pem=[],
            )


# =============================================================================
# TestProtectClientContextManager - Tests for context manager (4 tests)
# =============================================================================


class TestProtectClientContextManager:
    """Tests for ProtectClient context manager protocol."""

    def test_context_manager_returns_client(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that context manager returns the client instance."""
        with ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        ) as client:
            assert client is not None
            assert client.is_closed is False

    def test_context_manager_closes_on_exit(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that context manager closes client on normal exit."""
        with ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        ) as client:
            pass
        assert client.is_closed is True

    def test_context_manager_closes_on_exception(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that context manager closes client even when exception occurs."""
        client_ref = None
        try:
            with ProtectClient.create(
                host=host, api_key=api_key, api_secret=api_secret_hex,
                super_admin_keys_pem=super_admin_keys_pem,
            ) as client:
                client_ref = client
                raise ValueError("test exception")
        except ValueError:
            pass
        assert client_ref is not None
        assert client_ref.is_closed is True

    def test_nested_context_managers(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that multiple clients can be nested."""
        with ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        ) as client1:
            with ProtectClient.create(
                host=host, api_key=api_key, api_secret=api_secret_hex,
                super_admin_keys_pem=super_admin_keys_pem,
            ) as client2:
                assert client1 is not client2
                assert client1.is_closed is False
                assert client2.is_closed is False
            assert client2.is_closed is True
            assert client1.is_closed is False
        assert client1.is_closed is True


# =============================================================================
# TestProtectClientLifecycle - Tests for client lifecycle (5 tests)
# =============================================================================


class TestProtectClientLifecycle:
    """Tests for ProtectClient lifecycle management."""

    def test_close_sets_is_closed_true(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that close() sets is_closed to True."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        assert client.is_closed is False
        client.close()
        assert client.is_closed is True

    def test_close_is_idempotent(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that multiple close() calls do not raise."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        client.close()
        client.close()  # Second call should not raise
        client.close()  # Third call should not raise
        assert client.is_closed is True

    def test_closed_client_raises_on_service_access(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that accessing a service on a closed client raises RuntimeError."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        client.close()
        with pytest.raises(RuntimeError, match="has been closed"):
            _ = client.wallets

    def test_closed_client_raises_on_any_service(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that accessing any service on closed client raises RuntimeError."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        client.close()

        # Test a few different services
        with pytest.raises(RuntimeError, match="has been closed"):
            _ = client.addresses
        with pytest.raises(RuntimeError, match="has been closed"):
            _ = client.requests
        with pytest.raises(RuntimeError, match="has been closed"):
            _ = client.health

    def test_is_closed_property_is_readonly(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that is_closed is a read-only property."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        with pytest.raises(AttributeError):
            client.is_closed = True  # type: ignore
        client.close()


# =============================================================================
# TestProtectClientLazyServiceInit - Tests for lazy initialization (5 tests)
# =============================================================================


class TestProtectClientLazyServiceInit:
    """Tests for lazy service initialization."""

    def test_services_not_initialized_on_create(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that services are None after client creation."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        # Check internal state - services should be None
        assert client._wallet_service is None
        assert client._address_service is None
        assert client._request_service is None
        assert client._health_service is None
        client.close()

    def test_wallet_service_initialized_on_first_access(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that wallet service is initialized on first access."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        assert client._wallet_service is None
        _ = client.wallets
        assert client._wallet_service is not None
        client.close()

    def test_service_returns_same_instance_on_subsequent_access(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that accessing a service multiple times returns the same instance."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        first = client.wallets
        second = client.wallets
        third = client.wallets
        assert first is second
        assert second is third
        client.close()

    def test_different_services_have_different_instances(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that different service properties return different instances."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        wallets = client.wallets
        addresses = client.addresses
        assert wallets is not addresses
        client.close()

    def test_api_client_initialized_lazily(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that the OpenAPI client is initialized lazily."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        # API client should be None before any service access
        assert client._api_client is None
        # Access a service which should initialize the API client
        _ = client.wallets
        # Now API client should be set
        assert client._api_client is not None
        client.close()


# =============================================================================
# TestProtectClientServiceProperties - Tests for all 38+ services (39 tests)
# =============================================================================


class TestProtectClientServiceProperties:
    """Tests for all service properties returning non-None values."""

    # Services that work correctly (API class names match)
    @pytest.mark.parametrize(
        "service_name",
        [
            "wallets",
            "addresses",
            "requests",
            "transactions",
            "health",
            "jobs",
            "users",
            "groups",
            "visibility_groups",
            "tags",
            "webhooks",
            "webhook_calls",
            "audits",
            "governance_rules",
            "whitelisted_addresses",
            "whitelisted_assets",
            "contract_whitelisting",
            "staking",
            "reservations",
            "multi_factor_signature",
            "business_rules",
            "air_gap",
            "config",
            "assets",
            "changes",
            "statistics",
            "blockchains",
            "prices",
            "fee_payers",
            "scores",
            "token_metadata",
            "actions",
            "balances",
            "currencies",
            "exchanges",
            "fees",
            "fiat",
            "user_devices",
        ],
    )
    def test_service_property_returns_non_none(
        self, host: str, api_key: str, api_secret_hex: str,
        super_admin_keys_pem: List[str], service_name: str
    ) -> None:
        """Test that each service property returns a non-None value."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        service = getattr(client, service_name)
        assert service is not None
        client.close()

    def test_taurus_network_returns_non_none(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that taurus_network property returns a non-None value."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        tn = client.taurus_network
        assert tn is not None
        client.close()


# =============================================================================
# TestProtectClientTaurusNetwork - Tests for TaurusNetwork namespace (5 tests)
# =============================================================================


class TestProtectClientTaurusNetwork:
    """Tests for TaurusNetwork namespace services."""

    def test_taurus_network_participants_accessible(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that taurus_network.participants is accessible."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        participants = client.taurus_network.participants
        assert participants is not None
        client.close()

    def test_taurus_network_pledges_accessible(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that taurus_network.pledges is accessible."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        pledges = client.taurus_network.pledges
        assert pledges is not None
        client.close()

    def test_taurus_network_lending_accessible(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that taurus_network.lending is accessible."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        lending = client.taurus_network.lending
        assert lending is not None
        client.close()

    def test_taurus_network_settlements_accessible(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that taurus_network.settlements is accessible."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        settlements = client.taurus_network.settlements
        assert settlements is not None
        client.close()

    def test_taurus_network_sharing_accessible(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that taurus_network.sharing is accessible."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        sharing = client.taurus_network.sharing
        assert sharing is not None
        client.close()


# =============================================================================
# TestProtectClientThreadSafety - Tests for thread safety (3 tests)
# =============================================================================


class TestProtectClientThreadSafety:
    """Tests for thread-safe service access."""

    def test_concurrent_service_access(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that concurrent service access returns the same instance."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        results: list = []
        errors: list = []

        def access_wallets() -> None:
            try:
                results.append(client.wallets)
            except Exception as e:
                errors.append(e)

        threads = [threading.Thread(target=access_wallets) for _ in range(10)]
        for t in threads:
            t.start()
        for t in threads:
            t.join()

        assert len(errors) == 0, f"Errors occurred: {errors}"
        assert len(results) == 10
        # All should be the same instance
        assert all(r is results[0] for r in results)
        client.close()

    def test_concurrent_different_service_access(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that concurrent access to different services works correctly."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        wallet_results: list = []
        address_results: list = []
        health_results: list = []
        errors: list = []

        def access_wallets() -> None:
            try:
                wallet_results.append(client.wallets)
            except Exception as e:
                errors.append(e)

        def access_addresses() -> None:
            try:
                address_results.append(client.addresses)
            except Exception as e:
                errors.append(e)

        def access_health() -> None:
            try:
                health_results.append(client.health)
            except Exception as e:
                errors.append(e)

        threads = (
            [threading.Thread(target=access_wallets) for _ in range(5)]
            + [threading.Thread(target=access_addresses) for _ in range(5)]
            + [threading.Thread(target=access_health) for _ in range(5)]
        )
        for t in threads:
            t.start()
        for t in threads:
            t.join()

        assert len(errors) == 0, f"Errors occurred: {errors}"
        assert len(wallet_results) == 5
        assert len(address_results) == 5
        assert len(health_results) == 5
        # Each service type should have the same instance
        assert all(r is wallet_results[0] for r in wallet_results)
        assert all(r is address_results[0] for r in address_results)
        assert all(r is health_results[0] for r in health_results)
        client.close()

    def test_concurrent_close_does_not_raise(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that concurrent close calls do not raise exceptions."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        errors: list = []

        def close_client() -> None:
            try:
                client.close()
            except Exception as e:
                errors.append(e)

        threads = [threading.Thread(target=close_client) for _ in range(10)]
        for t in threads:
            t.start()
        for t in threads:
            t.join()

        assert len(errors) == 0, f"Errors occurred: {errors}"
        assert client.is_closed is True


# =============================================================================
# TestProtectClientConfiguration - Tests for configuration (6 tests)
# =============================================================================


class TestProtectClientConfiguration:
    """Tests for client configuration options."""

    def test_default_timeout(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that default timeout is applied."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        assert client._timeout == ProtectClient.DEFAULT_TIMEOUT
        client.close()

    def test_custom_timeout(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that custom timeout is applied."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
            timeout=60.0,
        )
        assert client._timeout == 60.0
        client.close()

    def test_default_rules_cache_ttl(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that default rules cache TTL is applied."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        assert client._rules_cache_ttl == ProtectClient.DEFAULT_RULES_CACHE_TTL
        client.close()

    def test_custom_rules_cache_ttl(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that custom rules cache TTL is applied."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
            rules_cache_ttl=600.0,
        )
        assert client._rules_cache_ttl == 600.0
        client.close()

    def test_host_property_returns_configured_host(
        self, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that host property returns the configured host."""
        client = ProtectClient.create(
            host="https://my-custom-host.com", api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        assert client.host == "https://my-custom-host.com"
        client.close()

    def test_multiple_super_admin_keys(
        self, host: str, api_key: str, api_secret_hex: str
    ) -> None:
        """Test creating client with multiple SuperAdmin keys."""
        # Generate two different key pairs
        key1 = ec.generate_private_key(ec.SECP256R1())
        key2 = ec.generate_private_key(ec.SECP256R1())

        pem1 = key1.public_key().public_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PublicFormat.SubjectPublicKeyInfo,
        ).decode("utf-8")

        pem2 = key2.public_key().public_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PublicFormat.SubjectPublicKeyInfo,
        ).decode("utf-8")

        client = ProtectClient.create(
            host=host,
            api_key=api_key,
            api_secret=api_secret_hex,
            super_admin_keys_pem=[pem1, pem2],
            min_valid_signatures=2,
        )
        assert len(client._super_admin_keys) == 2
        client.close()


# =============================================================================
# TestProtectClientServiceTypes - Tests for service type correctness (10 tests)
# =============================================================================


class TestProtectClientServiceTypes:
    """Tests verifying correct service types are returned."""

    def test_wallets_returns_wallet_service(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that wallets property returns WalletService."""
        from taurus_protect.services.wallet_service import WalletService

        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        assert isinstance(client.wallets, WalletService)
        client.close()

    def test_addresses_returns_address_service(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that addresses property returns AddressService."""
        from taurus_protect.services.address_service import AddressService

        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        assert isinstance(client.addresses, AddressService)
        client.close()

    def test_requests_returns_request_service(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that requests property returns RequestService."""
        from taurus_protect.services.request_service import RequestService

        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        assert isinstance(client.requests, RequestService)
        client.close()

    def test_transactions_returns_transaction_service(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that transactions property returns TransactionService."""
        from taurus_protect.services.transaction_service import TransactionService

        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        assert isinstance(client.transactions, TransactionService)
        client.close()

    def test_health_returns_health_service(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that health property returns HealthService."""
        from taurus_protect.services.health_service import HealthService

        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        assert isinstance(client.health, HealthService)
        client.close()

    def test_governance_rules_returns_governance_rule_service(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that governance_rules property returns GovernanceRuleService."""
        from taurus_protect.services.governance_rule_service import GovernanceRuleService

        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        assert isinstance(client.governance_rules, GovernanceRuleService)
        client.close()

    def test_users_returns_user_service(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that users property returns UserService."""
        from taurus_protect.services.user_service import UserService

        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        assert isinstance(client.users, UserService)
        client.close()

    def test_groups_returns_group_service(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that groups property returns GroupService."""
        from taurus_protect.services.group_service import GroupService

        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        assert isinstance(client.groups, GroupService)
        client.close()

    def test_balances_returns_balance_service(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that balances property returns BalanceService."""
        from taurus_protect.services.balance_service import BalanceService

        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        assert isinstance(client.balances, BalanceService)
        client.close()

    def test_currencies_returns_currency_service(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that currencies property returns CurrencyService."""
        from taurus_protect.services.currency_service import CurrencyService

        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        assert isinstance(client.currencies, CurrencyService)
        client.close()


# =============================================================================
# TestProtectClientRulesCache - Tests for rules cache (4 tests)
# =============================================================================


class TestProtectClientRulesCache:
    """Tests for rules container cache initialization."""

    def test_rules_cache_not_initialized_on_create(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that rules cache is None after client creation."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        assert client._rules_cache is None
        client.close()

    def test_rules_cache_initialized_on_address_service_access(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that rules cache is initialized when address service is accessed."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        # Access address service which requires rules cache
        _ = client.addresses
        # Rules cache should now be initialized
        assert client._rules_cache is not None
        client.close()

    def test_governance_rules_service_initialized_before_rules_cache(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that governance_rules service is set up when rules cache is created."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        # Access addresses which triggers rules cache creation
        _ = client.addresses
        # governance_rules service should be initialized (rules cache depends on it)
        assert client._governance_rule_service is not None
        client.close()

    def test_rules_cache_reused_across_services(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that the same rules cache is reused."""
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        # Access addresses first
        _ = client.addresses
        first_cache = client._rules_cache
        # Access again
        _ = client.addresses
        second_cache = client._rules_cache
        assert first_cache is second_cache
        client.close()


# =============================================================================
# TestProtectClientCreateFromPem - Tests for create_from_pem (3 tests)
# =============================================================================


class TestProtectClientCreateFromPem:
    """Tests for create_from_pem factory method."""

    def test_create_from_pem_with_valid_hex(
        self, host: str, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test create_from_pem works with hex-encoded secret."""
        client = ProtectClient.create_from_pem(
            host=host, api_key=api_key, api_secret_pem=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        assert client is not None
        assert client.is_closed is False
        client.close()

    def test_create_from_pem_with_super_admin_keys(
        self, host: str, api_key: str, api_secret_hex: str, ecdsa_public_key_pem: str
    ) -> None:
        """Test create_from_pem accepts super_admin_keys_pem."""
        client = ProtectClient.create_from_pem(
            host=host,
            api_key=api_key,
            api_secret_pem=api_secret_hex,
            super_admin_keys_pem=[ecdsa_public_key_pem],
        )
        assert client is not None
        assert len(client._super_admin_keys) == 1
        client.close()

    def test_create_from_pem_validates_inputs(
        self, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test create_from_pem validates inputs like create()."""
        with pytest.raises(ConfigurationError, match="host cannot be empty"):
            ProtectClient.create_from_pem(
                host="", api_key=api_key, api_secret_pem=api_secret_hex,
                super_admin_keys_pem=super_admin_keys_pem,
            )


# =============================================================================
# TestProtectClientEdgeCases - Edge case tests (4 tests)
# =============================================================================


class TestProtectClientEdgeCases:
    """Tests for edge cases and boundary conditions."""

    def test_whitespace_only_host_creates_client(
        self, api_key: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that whitespace-only host is normalized to whitespace (no strip).

        Note: The current implementation does not strip whitespace before validation.
        This test documents the current behavior - whitespace strings are accepted
        because they are truthy in Python. The host will be stripped of trailing slash
        only, resulting in whitespace host being preserved.
        """
        # Current behavior: whitespace is not stripped, so "   " is truthy and accepted
        client = ProtectClient.create(
            host="   ", api_key=api_key, api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        # The host will only have trailing slashes removed
        assert client.host == "   "
        client.close()

    def test_whitespace_only_api_key_creates_client(
        self, host: str, api_secret_hex: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that whitespace-only api_key is accepted.

        Note: The current implementation does not strip whitespace before validation.
        This test documents the current behavior.
        """
        client = ProtectClient.create(
            host=host, api_key="   ", api_secret=api_secret_hex,
            super_admin_keys_pem=super_admin_keys_pem,
        )
        # Should be created since "   " is truthy
        assert client is not None
        client.close()

    def test_whitespace_only_api_secret_creates_client(
        self, host: str, api_key: str, super_admin_keys_pem: List[str]
    ) -> None:
        """Test that whitespace-only api_secret is accepted.

        Note: The current implementation strips whitespace during hex decoding,
        so whitespace results in an empty bytes secret. This test documents
        the current behavior.
        """
        # Whitespace is stripped during hex decoding, resulting in empty secret
        client = ProtectClient.create(
            host=host, api_key=api_key, api_secret="   ",
            super_admin_keys_pem=super_admin_keys_pem,
        )
        assert client is not None
        client.close()

    def test_zero_min_valid_signatures_rejected_with_keys(
        self, host: str, api_key: str, api_secret_hex: str, ecdsa_public_key_pem: str
    ) -> None:
        """Test that min_valid_signatures=0 is rejected when SuperAdmin keys are provided."""
        with pytest.raises(ConfigurationError, match="min_valid_signatures must be greater than zero"):
            ProtectClient.create(
                host=host,
                api_key=api_key,
                api_secret=api_secret_hex,
                super_admin_keys_pem=[ecdsa_public_key_pem],
                min_valid_signatures=0,
            )
