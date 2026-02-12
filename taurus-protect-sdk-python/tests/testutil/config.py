"""Test configuration with multi-identity support.

Mirrors Java TestConfig.java — loads from test.properties with env var overrides.
"""

from __future__ import annotations

import os
from dataclasses import dataclass
from typing import List, Optional

from tests.testutil.properties import load_properties


@dataclass
class Identity:
    """A user identity with optional API credentials, private key, and public key."""

    index: int
    name: str
    api_key: str
    api_secret: str
    private_key: str
    public_key: str

    def has_api_credentials(self) -> bool:
        return bool(self.api_key and self.api_secret)

    def has_private_key(self) -> bool:
        return bool(self.private_key)

    def has_public_key(self) -> bool:
        return bool(self.public_key)

    def __str__(self) -> str:
        return f"{self.name} (identity {self.index})"


class TestConfig:
    """Configuration for tests with multi-identity support.

    Loads from test.properties in the project root, with environment
    variable overrides for CI/CD pipelines.
    """

    _ENV_INTEGRATION_TEST = "PROTECT_INTEGRATION_TEST"
    _ENV_API_HOST = "PROTECT_API_HOST"

    def __init__(self) -> None:
        self._properties = self._load_properties()
        self._identities = self._load_identities()

    def _load_properties(self) -> dict:
        search_paths = [
            "tests/testutil/test.properties",
            "test.properties",
            "taurus-protect-sdk-python/tests/testutil/test.properties",
            "taurus-protect-sdk-python/test.properties",
        ]
        for path in search_paths:
            props = load_properties(path)
            if props:
                return props
        return {}

    def _resolve(self, env_var: Optional[str], prop_key: str) -> str:
        if env_var:
            val = os.environ.get(env_var, "")
            if val:
                return val
        return self._properties.get(prop_key, "")

    def _load_identities(self) -> List[Identity]:
        identities: List[Identity] = []
        for i in range(1, 100):
            name = self._resolve(None, f"identity.{i}.name")
            api_key = self._resolve(f"PROTECT_API_KEY_{i}", f"identity.{i}.apiKey")
            api_secret = self._resolve(f"PROTECT_API_SECRET_{i}", f"identity.{i}.apiSecret")
            private_key = self._resolve(f"PROTECT_PRIVATE_KEY_{i}", f"identity.{i}.privateKey")
            public_key = self._resolve(f"PROTECT_PUBLIC_KEY_{i}", f"identity.{i}.publicKey")

            # Skip gaps (e.g., identity.1 then identity.4)
            if not any([name, api_key, api_secret, private_key, public_key]):
                continue

            if not name:
                name = f"identity-{i}"

            identities.append(Identity(
                index=i,
                name=name,
                api_key=api_key,
                api_secret=api_secret,
                private_key=private_key,
                public_key=public_key,
            ))
        return identities

    @property
    def host(self) -> str:
        env = os.environ.get(self._ENV_API_HOST, "")
        if env:
            return env
        return self._properties.get("host", "")

    def get_identity(self, index: int) -> Identity:
        """Get identity by 1-based index (supports gaps in identity numbering)."""
        for identity in self._identities:
            if identity.index == index:
                return identity
        raise IndexError(
            f"Identity {index} not found. Available: "
            f"{[i.index for i in self._identities]}"
        )

    @property
    def identity_count(self) -> int:
        return len(self._identities)

    @property
    def identities(self) -> List[Identity]:
        return list(self._identities)

    def get_super_admin_keys(self) -> List[str]:
        """Return PEM-encoded public keys from identities that have one."""
        return [i.public_key for i in self._identities if i.has_public_key()]

    def get_min_valid_signatures(self) -> int:
        val = self._properties.get("minValidSignatures", "2")
        try:
            return int(val)
        except ValueError:
            return 2

    def is_enabled(self) -> bool:
        """Check if tests should run."""
        env = os.environ.get(self._ENV_INTEGRATION_TEST, "")
        if env.lower() == "true":
            return True
        return any(i.has_api_credentials() for i in self._identities)


# Module-level singleton
_config: Optional[TestConfig] = None


def get_config() -> TestConfig:
    """Get or create the singleton TestConfig."""
    global _config
    if _config is None:
        _config = TestConfig()
    return _config
