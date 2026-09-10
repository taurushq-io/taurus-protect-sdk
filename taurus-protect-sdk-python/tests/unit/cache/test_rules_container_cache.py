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

        # Exactly once. This used to allow <= 2 "if timing is tight", which is the
        # tolerance that hid a real stampede: waiters were woken before the fetched
        # container was published, so they saw an expired cache with no fetch in
        # flight and each started their own.
        assert mock_service.get_rules.call_count == 1

    @staticmethod
    def _notify_observations(cache: RulesContainerCache) -> list:
        """Record (fetching, container_published) at every notify_all() call.

        Asserting the ORDERING INVARIANT directly rather than racing for the
        symptom. A timing-based test does not work here: CPython hands the lock
        straight back to the thread that released it, so the fetcher almost always
        wins the window and a stampede test passes even with the bug present.

        The invariant: a notify that clears `_fetching` must find the container
        already published, or the waiters it wakes will see an expired cache with
        no fetch in flight and start their own.
        """
        observations: list = []
        condition = cache._condition
        original = condition.notify_all

        def spy() -> None:
            observations.append(
                (cache._fetching, cache._cached_container is not None)
            )
            original()

        condition.notify_all = spy  # type: ignore[method-assign]
        return observations

    def test_publishes_container_before_waking_waiters(self) -> None:
        """get_decoded_rules_container must publish before it notifies."""
        mock_service = MagicMock()
        mock_container = DecodedRulesContainer()
        mock_service.get_rules.return_value = MagicMock()
        mock_service.get_decoded_rules_container.return_value = mock_container

        cache = RulesContainerCache(mock_service)
        observations = self._notify_observations(cache)

        assert cache.get_decoded_rules_container() is mock_container

        assert observations, "notify_all was never called"
        fetching, published = observations[-1]
        assert fetching is False, "the notify should be the one releasing the slot"
        assert published, (
            "container was still unpublished when waiters were woken: they will "
            "see an expired cache with _fetching already false and refetch"
        )

    def test_invalidate_publishes_container_before_waking_waiters(self) -> None:
        """`invalidate` carries the identical ordering; both entry points are fixed."""
        mock_service = MagicMock()
        mock_container = DecodedRulesContainer()
        mock_service.get_rules.return_value = MagicMock()
        mock_service.get_decoded_rules_container.return_value = mock_container

        cache = RulesContainerCache(mock_service)
        observations = self._notify_observations(cache)

        cache.invalidate()

        assert observations, "notify_all was never called"
        fetching, published = observations[-1]
        assert fetching is False
        assert published, "invalidate woke waiters before publishing the container"

    def test_failed_fetch_still_releases_the_slot(self) -> None:
        """A failing fetch must not leave waiters parked forever."""
        mock_service = MagicMock()
        mock_service.get_rules.side_effect = APIError("boom")

        cache = RulesContainerCache(mock_service)
        observations = self._notify_observations(cache)

        with pytest.raises(APIError):
            cache.get_decoded_rules_container()

        assert observations, "waiters were never woken after a failed fetch"
        fetching, published = observations[-1]
        assert fetching is False, "the single-flight slot was not released"
        assert not published, "nothing should be published when the fetch failed"

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


class TestRulesContainerCacheVerification:
    """The cache must never hand out a container nothing verified.

    This cache supplies the HSM public key that address signature verification
    trusts, so an unverified container makes that check pass against whatever key
    the container carries — attacker-chosen addresses verify clean.

    Cross-SDK invariant, and TypeScript had drifted off it: its cache fetched
    through the generated API and the raw mapper, skipping verification entirely.
    Nothing in any SDK tested the path, which is why it survived. Guard for this one.
    """

    @staticmethod
    def _wire_valid_container() -> str:
        """A container that DECODES cleanly, so only the signatures are at fault.

        A malformed blob would raise IntegrityError from the decoder instead, which
        makes the test pass whether or not verification runs at all.
        """
        from taurus_protect.mappers import rules_container_to_base64

        return rules_container_to_base64(
            DecodedRulesContainer(minimum_distinct_user_signatures=2)
        )

    @staticmethod
    def _service_returning(container: str, signatures: list) -> object:
        """A real service with real keys over a stubbed API."""
        from cryptography.hazmat.primitives.asymmetric import ec

        from taurus_protect.services.governance_rule_service import (
            GovernanceRuleService,
        )

        key = ec.generate_private_key(ec.SECP256R1()).public_key()

        dto = MagicMock()
        dto.rules_container = container
        dto.rules_signatures = signatures
        dto.locked = False
        reply = MagicMock()
        reply.result = dto

        api = MagicMock()
        api.rule_service_get_rules.return_value = reply

        return GovernanceRuleService(MagicMock(), api, [key], 1)

    def test_the_container_decodes_so_only_signatures_can_fail(self) -> None:
        """Guards the two tests below from passing on a decode error instead."""
        from taurus_protect.mappers import rules_container_from_base64

        decoded = rules_container_from_base64(self._wire_valid_container())
        assert decoded.minimum_distinct_user_signatures == 2

    def test_get_refuses_a_container_with_no_signatures(self) -> None:
        """An unsigned container must not reach the caller."""
        from taurus_protect.errors import IntegrityError

        service = self._service_returning(self._wire_valid_container(), [])
        cache = RulesContainerCache(service)

        with pytest.raises(IntegrityError):
            cache.get_decoded_rules_container()

    def test_get_refuses_signatures_from_an_unconfigured_key(self) -> None:
        """Signatures present, but from a key this client does not trust."""
        from taurus_protect.errors import IntegrityError

        signature = MagicMock()
        signature.user_id = "attacker"
        signature.signature = "YmFkLXNpZy1ub3QtZWNkc2E="

        service = self._service_returning(self._wire_valid_container(), [signature])
        cache = RulesContainerCache(service)

        with pytest.raises(IntegrityError):
            cache.get_decoded_rules_container()
