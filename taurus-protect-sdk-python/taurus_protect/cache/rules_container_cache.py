"""Thread-safe cache for decoded governance rules container."""

from __future__ import annotations

import threading
import time
from typing import TYPE_CHECKING, Optional

from taurus_protect.errors import APIError
from taurus_protect.models.governance_rules import DecodedRulesContainer

if TYPE_CHECKING:
    from taurus_protect.services.governance_rule_service import GovernanceRuleService


class RulesContainerCache:
    """
    Thread-safe cache for the decoded rules container with configurable TTL.

    This cache stores the decoded rules container and refreshes it automatically
    when the TTL expires. Used primarily for address signature verification
    against HSM slot public keys.

    The implementation uses a "fetching" flag with Condition coordination to avoid
    holding the lock during network I/O, which could cause deadlocks under load.

    Example:
        >>> cache = RulesContainerCache(governance_service, ttl_ms=300_000)
        >>> rules = cache.get_decoded_rules_container()
        >>> hsm_key = rules.get_hsm_public_key()

    Attributes:
        DEFAULT_CACHE_TTL_MS: Default TTL of 5 minutes (300,000 ms).
    """

    DEFAULT_CACHE_TTL_MS: int = 5 * 60 * 1000  # 5 minutes

    def __init__(
        self,
        governance_rule_service: "GovernanceRuleService",
        ttl_ms: int = DEFAULT_CACHE_TTL_MS,
    ) -> None:
        """
        Create a new rules container cache.

        Args:
            governance_rule_service: Service for fetching governance rules.
            ttl_ms: Cache time-to-live in milliseconds. Defaults to 5 minutes.

        Raises:
            ValueError: If governance_rule_service is None or ttl_ms <= 0.
        """
        if governance_rule_service is None:
            raise ValueError("governance_rule_service cannot be None")
        if ttl_ms <= 0:
            raise ValueError("ttl_ms must be positive")

        self._governance_rule_service = governance_rule_service
        self._ttl_ms = ttl_ms
        self._lock = threading.RLock()
        self._condition = threading.Condition(self._lock)
        self._cached_container: Optional[DecodedRulesContainer] = None
        # Use monotonic time to be immune to system clock changes
        self._cache_timestamp_mono: float = 0.0
        self._fetching: bool = False

    @property
    def ttl_ms(self) -> int:
        """Get the configured cache TTL in milliseconds."""
        return self._ttl_ms

    def _is_cache_expired(self) -> bool:
        """Check if cache is expired (must be called under lock)."""
        if self._cached_container is None:
            return True
        elapsed_ms = (time.monotonic() - self._cache_timestamp_mono) * 1000
        return elapsed_ms > self._ttl_ms

    def get_decoded_rules_container(self) -> DecodedRulesContainer:
        """
        Get the decoded rules container, fetching from API if cache is expired.

        This method is thread-safe and will only fetch once if multiple threads
        attempt to refresh simultaneously. The lock is released during network I/O
        to prevent deadlocks under concurrent load.

        Returns:
            The decoded rules container.

        Raises:
            APIError: If fetching the rules fails.
        """
        with self._condition:
            # If cache is valid, return immediately
            if not self._is_cache_expired():
                return self._cached_container  # type: ignore[return-value]

            # If another thread is already fetching, wait for it
            while self._fetching:
                self._condition.wait()
                # After waking, check if cache is now valid
                if not self._is_cache_expired():
                    return self._cached_container  # type: ignore[return-value]

            # We need to fetch - mark as fetching
            self._fetching = True

        # Network I/O happens OUTSIDE the lock to prevent deadlocks
        try:
            new_container = self._fetch_rules_container()
        finally:
            with self._condition:
                self._fetching = False
                self._condition.notify_all()

        # Update cache under lock
        with self._condition:
            self._cached_container = new_container
            self._cache_timestamp_mono = time.monotonic()
            return new_container

    def invalidate(self) -> None:
        """
        Force a cache refresh, fetching the latest rules from the API.

        Raises:
            APIError: If fetching the rules fails.
        """
        with self._condition:
            # If another thread is already fetching, wait for it
            while self._fetching:
                self._condition.wait()

            # Mark as fetching
            self._fetching = True

        # Network I/O happens OUTSIDE the lock
        try:
            new_container = self._fetch_rules_container()
        finally:
            with self._condition:
                self._fetching = False
                self._condition.notify_all()

        # Update cache under lock
        with self._condition:
            self._cached_container = new_container
            self._cache_timestamp_mono = time.monotonic()

    def is_cache_valid(self) -> bool:
        """
        Check if the cache is currently valid (not expired).

        Returns:
            True if the cache is valid, False if expired or empty.
        """
        with self._lock:
            return not self._is_cache_expired()

    def clear(self) -> None:
        """Clear the cache without refreshing."""
        with self._lock:
            self._cached_container = None
            self._cache_timestamp_mono = 0.0

    def _fetch_rules_container(self) -> DecodedRulesContainer:
        """
        Fetch governance rules from the API.

        This method performs network I/O and should be called WITHOUT holding the lock.

        Raises:
            APIError: If fetching fails or no rules are returned.
        """
        rules = self._governance_rule_service.get_rules()
        if rules is None:
            raise APIError("No governance rules available from API")

        # Get the decoded rules container (with signature verification)
        container = self._governance_rule_service.get_decoded_rules_container(rules)
        if container is None:
            raise APIError("Failed to decode governance rules container")
        return container
