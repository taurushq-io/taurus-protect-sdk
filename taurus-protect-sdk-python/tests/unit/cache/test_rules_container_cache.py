"""Tests for RulesContainerCache."""

import threading
import time
from unittest.mock import MagicMock, patch

import pytest

from taurus_protect.cache.rules_container_cache import RulesContainerCache
from taurus_protect.errors import APIError
from taurus_protect.models.governance_rules import DecodedRulesContainer, RuleUser


class TestRulesContainerCacheInit:
    """Tests for RulesContainerCache initialization."""

    def test_init_with_valid_service(self) -> None:
        """Test initialization with valid service."""
        mock_service = MagicMock()
        cache = RulesContainerCache(mock_service)

        assert cache.ttl_ms == RulesContainerCache.DEFAULT_CACHE_TTL_MS
        assert cache._governance_rule_service is mock_service

    def test_init_with_custom_ttl(self) -> None:
        """Test initialization with custom TTL."""
        mock_service = MagicMock()
        cache = RulesContainerCache(mock_service, ttl_ms=60000)

        assert cache.ttl_ms == 60000

    def test_init_with_none_service_raises(self) -> None:
        """Test that None service raises ValueError."""
        with pytest.raises(ValueError, match="governance_rule_service cannot be None"):
            RulesContainerCache(None)  # type: ignore

    def test_init_with_zero_ttl_raises(self) -> None:
        """Test that zero TTL raises ValueError."""
        mock_service = MagicMock()
        with pytest.raises(ValueError, match="ttl_ms must be positive"):
            RulesContainerCache(mock_service, ttl_ms=0)

    def test_init_with_negative_ttl_raises(self) -> None:
        """Test that negative TTL raises ValueError."""
        mock_service = MagicMock()
        with pytest.raises(ValueError, match="ttl_ms must be positive"):
            RulesContainerCache(mock_service, ttl_ms=-1000)


class TestGetDecodedRulesContainer:
    """Tests for get_decoded_rules_container method."""

    def test_fetches_on_first_call(self) -> None:
        """Test that first call fetches from API."""
        mock_service = MagicMock()
        mock_rules = MagicMock()
        mock_container = DecodedRulesContainer()

        mock_service.get_rules.return_value = mock_rules
        mock_service.get_decoded_rules_container.return_value = mock_container

        cache = RulesContainerCache(mock_service)
        result = cache.get_decoded_rules_container()

        assert result is mock_container
        mock_service.get_rules.assert_called_once()
        mock_service.get_decoded_rules_container.assert_called_once_with(mock_rules)

    def test_returns_cached_on_subsequent_calls(self) -> None:
        """Test that subsequent calls return cached value."""
        mock_service = MagicMock()
        mock_rules = MagicMock()
        mock_container = DecodedRulesContainer()

        mock_service.get_rules.return_value = mock_rules
        mock_service.get_decoded_rules_container.return_value = mock_container

        cache = RulesContainerCache(mock_service)

        # First call
        result1 = cache.get_decoded_rules_container()
        # Second call
        result2 = cache.get_decoded_rules_container()

        assert result1 is result2
        # Should only have fetched once
        assert mock_service.get_rules.call_count == 1

    def test_refreshes_after_ttl_expires(self) -> None:
        """Test that cache refreshes after TTL expires."""
        mock_service = MagicMock()
        mock_rules = MagicMock()
        container1 = DecodedRulesContainer(timestamp=1)
        container2 = DecodedRulesContainer(timestamp=2)

        mock_service.get_rules.return_value = mock_rules
        mock_service.get_decoded_rules_container.side_effect = [container1, container2]

        # Use very short TTL for testing
        cache = RulesContainerCache(mock_service, ttl_ms=10)

        # First call
        result1 = cache.get_decoded_rules_container()
        assert result1.timestamp == 1

        # Wait for TTL to expire
        time.sleep(0.02)  # 20ms > 10ms TTL

        # Second call should refresh
        result2 = cache.get_decoded_rules_container()
        assert result2.timestamp == 2

        assert mock_service.get_rules.call_count == 2

    def test_raises_if_service_returns_none_rules(self) -> None:
        """Test that APIError is raised if service returns None rules."""
        mock_service = MagicMock()
        mock_service.get_rules.return_value = None

        cache = RulesContainerCache(mock_service)

        with pytest.raises(APIError, match="No governance rules available from API"):
            cache.get_decoded_rules_container()

    def test_raises_if_decoded_container_is_none(self) -> None:
        """Test APIError if service returns None decoded container."""
        mock_service = MagicMock()
        mock_rules = MagicMock()
        mock_service.get_rules.return_value = mock_rules
        mock_service.get_decoded_rules_container.return_value = None

        cache = RulesContainerCache(mock_service)

        # Service returns valid rules but decoded container is None
        with pytest.raises(APIError, match="Failed to decode governance rules container"):
            cache.get_decoded_rules_container()


class TestInvalidate:
    """Tests for invalidate method."""

    def test_invalidate_forces_refresh(self) -> None:
        """Test that invalidate forces a refresh."""
        mock_service = MagicMock()
        mock_rules = MagicMock()
        container1 = DecodedRulesContainer(timestamp=1)
        container2 = DecodedRulesContainer(timestamp=2)

        mock_service.get_rules.return_value = mock_rules
        mock_service.get_decoded_rules_container.side_effect = [container1, container2]

        # Use long TTL so it won't expire naturally
        cache = RulesContainerCache(mock_service, ttl_ms=300000)

        # First call
        result1 = cache.get_decoded_rules_container()
        assert result1.timestamp == 1

        # Force refresh
        cache.invalidate()

        # Next call should get new value
        result2 = cache.get_decoded_rules_container()
        assert result2.timestamp == 2


class TestIsCacheValid:
    """Tests for is_cache_valid method."""

    def test_returns_false_when_empty(self) -> None:
        """Test that empty cache returns False."""
        mock_service = MagicMock()
        cache = RulesContainerCache(mock_service)

        assert cache.is_cache_valid() is False

    def test_returns_true_when_valid(self) -> None:
        """Test that valid cache returns True."""
        mock_service = MagicMock()
        mock_rules = MagicMock()
        mock_container = DecodedRulesContainer()

        mock_service.get_rules.return_value = mock_rules
        mock_service.get_decoded_rules_container.return_value = mock_container

        cache = RulesContainerCache(mock_service, ttl_ms=300000)
        cache.get_decoded_rules_container()

        assert cache.is_cache_valid() is True

    def test_returns_false_when_expired(self) -> None:
        """Test that expired cache returns False."""
        mock_service = MagicMock()
        mock_rules = MagicMock()
        mock_container = DecodedRulesContainer()

        mock_service.get_rules.return_value = mock_rules
        mock_service.get_decoded_rules_container.return_value = mock_container

        cache = RulesContainerCache(mock_service, ttl_ms=10)
        cache.get_decoded_rules_container()

        # Wait for expiry
        time.sleep(0.02)

        assert cache.is_cache_valid() is False


class TestClear:
    """Tests for clear method."""

    def test_clear_empties_cache(self) -> None:
        """Test that clear empties the cache."""
        mock_service = MagicMock()
        mock_rules = MagicMock()
        mock_container = DecodedRulesContainer()

        mock_service.get_rules.return_value = mock_rules
        mock_service.get_decoded_rules_container.return_value = mock_container

        cache = RulesContainerCache(mock_service)
        cache.get_decoded_rules_container()

        assert cache.is_cache_valid() is True

        cache.clear()

        assert cache.is_cache_valid() is False

    def test_clear_does_not_refresh(self) -> None:
        """Test that clear does not trigger a refresh."""
        mock_service = MagicMock()
        mock_rules = MagicMock()
        mock_container = DecodedRulesContainer()

        mock_service.get_rules.return_value = mock_rules
        mock_service.get_decoded_rules_container.return_value = mock_container

        cache = RulesContainerCache(mock_service)
        cache.get_decoded_rules_container()

        call_count_before = mock_service.get_rules.call_count

        cache.clear()

        # Should not have made additional API calls
        assert mock_service.get_rules.call_count == call_count_before


class TestThreadSafety:
    """Tests for thread safety of the cache."""

    def test_concurrent_access(self) -> None:
        """Test that concurrent access is thread-safe."""
        mock_service = MagicMock()
        mock_rules = MagicMock()
        mock_container = DecodedRulesContainer()

        # Simulate slow API call
        def slow_get_rules():
            time.sleep(0.01)
            return mock_rules

        mock_service.get_rules.side_effect = slow_get_rules
        mock_service.get_decoded_rules_container.return_value = mock_container

        cache = RulesContainerCache(mock_service)

        results = []
        errors = []

        def get_container():
            try:
                result = cache.get_decoded_rules_container()
                results.append(result)
            except Exception as e:
                errors.append(e)

        # Start multiple threads
        threads = [threading.Thread(target=get_container) for _ in range(10)]
        for t in threads:
            t.start()
        for t in threads:
            t.join()

        # All should succeed
        assert len(errors) == 0
        assert len(results) == 10

        # All should get the same instance
        assert all(r is mock_container for r in results)

        # Should only have fetched once (or possibly twice if timing is tight)
        assert mock_service.get_rules.call_count <= 2

    def test_concurrent_invalidate_and_get(self) -> None:
        """Test concurrent invalidate and get operations."""
        mock_service = MagicMock()
        mock_rules = MagicMock()
        call_count = [0]

        def make_container(rules):
            call_count[0] += 1
            return DecodedRulesContainer(timestamp=call_count[0])

        mock_service.get_rules.return_value = mock_rules
        mock_service.get_decoded_rules_container.side_effect = make_container

        cache = RulesContainerCache(mock_service)

        errors = []

        def invalidate_loop():
            try:
                for _ in range(5):
                    cache.invalidate()
                    time.sleep(0.001)
            except Exception as e:
                errors.append(e)

        def get_loop():
            try:
                for _ in range(5):
                    cache.get_decoded_rules_container()
                    time.sleep(0.001)
            except Exception as e:
                errors.append(e)

        threads = [
            threading.Thread(target=invalidate_loop),
            threading.Thread(target=get_loop),
            threading.Thread(target=get_loop),
        ]

        for t in threads:
            t.start()
        for t in threads:
            t.join()

        # Should complete without errors
        assert len(errors) == 0
