"""Config service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List

from taurus_protect.mappers.webhook import tenant_config_from_dto
from taurus_protect.models.webhook import Feature, TenantConfig
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


class ConfigService(BaseService):
    """
    Service for tenant configuration operations.

    Provides methods to retrieve tenant configuration and enabled features.

    Example:
        >>> # Get tenant configuration
        >>> config = client.config.get()
        >>> print(f"Tenant: {config.tenant_id}")
        >>> print(f"Base currency: {config.base_currency}")
        >>>
        >>> # Get enabled features
        >>> features = client.config.get_features()
        >>> for feature in features:
        ...     print(f"{feature.name}: {feature.enabled}")
    """

    def __init__(self, api_client: Any, config_api: Any) -> None:
        """
        Initialize config service.

        Args:
            api_client: The OpenAPI client instance.
            config_api: The ConfigAPI service from OpenAPI client.
        """
        super().__init__(api_client)
        self._config_api = config_api

    def get(self) -> TenantConfig:
        """
        Get tenant configuration.

        Returns:
            The tenant configuration.

        Raises:
            APIError: If API request fails.
        """
        try:
            resp = self._config_api.status_service_get_config_tenant()

            config_dto = getattr(resp, "config", None)
            if config_dto is None:
                # Return empty config if none returned
                return TenantConfig()

            config = tenant_config_from_dto(config_dto)
            if config is None:
                return TenantConfig()

            return config
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def get_features(self) -> List[Feature]:
        """
        Get enabled features for the tenant.

        This method retrieves the tenant configuration and extracts
        feature flags from it.

        Returns:
            List of enabled features.

        Raises:
            APIError: If API request fails.
        """
        try:
            config = self.get()

            # Build feature list from config booleans
            features: List[Feature] = []

            if config.is_mfa_mandatory:
                features.append(Feature(name="mfa_mandatory", enabled=True))

            if config.is_protect_engine_cold:
                features.append(Feature(name="protect_engine_cold", enabled=True))

            if config.is_cold_protect_engine_offline:
                features.append(Feature(name="cold_protect_engine_offline", enabled=True))

            if config.is_physical_air_gap_enabled:
                features.append(Feature(name="physical_air_gap", enabled=True))

            if config.restrict_sources_for_whitelisted_addresses:
                features.append(
                    Feature(name="restrict_sources_for_whitelisted_addresses", enabled=True)
                )

            if config.exclude_container:
                features.append(Feature(name="exclude_container", enabled=True))

            return features
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e
