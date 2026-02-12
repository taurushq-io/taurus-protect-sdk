"""ProtectClient - Main entry point for Taurus-PROTECT SDK."""

from __future__ import annotations

import threading
from typing import TYPE_CHECKING, Any, List, Optional

from cryptography.hazmat.primitives.asymmetric.ec import EllipticCurvePublicKey

from taurus_protect.crypto.keys import decode_public_keys_pem
from taurus_protect.crypto.tpv1 import TPV1Auth
from taurus_protect.errors import ConfigurationError

if TYPE_CHECKING:
    from taurus_protect.services.action_service import ActionService
    from taurus_protect.services.address_service import AddressService
    from taurus_protect.services.air_gap_service import AirGapService
    from taurus_protect.services.asset_service import AssetService
    from taurus_protect.services.audit_service import AuditService
    from taurus_protect.services.balance_service import BalanceService
    from taurus_protect.services.blockchain_service import BlockchainService
    from taurus_protect.services.business_rule_service import BusinessRuleService
    from taurus_protect.services.change_service import ChangeService
    from taurus_protect.services.config_service import ConfigService
    from taurus_protect.services.contract_whitelisting_service import (
        ContractWhitelistingService,
    )
    from taurus_protect.services.currency_service import CurrencyService
    from taurus_protect.services.exchange_service import ExchangeService
    from taurus_protect.services.fee_payer_service import FeePayerService
    from taurus_protect.services.fee_service import FeeService
    from taurus_protect.services.fiat_service import FiatService
    from taurus_protect.services.governance_rule_service import GovernanceRuleService
    from taurus_protect.services.group_service import GroupService
    from taurus_protect.services.health_service import HealthService
    from taurus_protect.services.job_service import JobService
    from taurus_protect.services.multi_factor_signature_service import (
        MultiFactorSignatureService,
    )
    from taurus_protect.services.price_service import PriceService
    from taurus_protect.services.request_service import RequestService
    from taurus_protect.services.reservation_service import ReservationService
    from taurus_protect.services.score_service import ScoreService
    from taurus_protect.services.staking_service import StakingService
    from taurus_protect.services.statistics_service import StatisticsService
    from taurus_protect.services.tag_service import TagService
    from taurus_protect.services.taurus_network._client import TaurusNetworkClient
    from taurus_protect.services.token_metadata_service import TokenMetadataService
    from taurus_protect.services.transaction_service import TransactionService
    from taurus_protect.services.user_device_service import UserDeviceService
    from taurus_protect.services.user_service import UserService
    from taurus_protect.services.visibility_group_service import VisibilityGroupService
    from taurus_protect.services.wallet_service import WalletService
    from taurus_protect.services.webhook_call_service import WebhookCallService
    from taurus_protect.services.webhook_service import WebhookService
    from taurus_protect.services.whitelisted_address_service import (
        WhitelistedAddressService,
    )
    from taurus_protect.services.whitelisted_asset_service import WhitelistedAssetService


class ProtectClient:
    """
    Main entry point for the Taurus-PROTECT SDK.

    Provides access to all API services through lazy-initialized properties.
    Implements context manager protocol for proper resource cleanup.

    Example:
        >>> with ProtectClient.create(
        ...     host="https://api.protect.taurushq.com",
        ...     api_key="your-api-key",
        ...     api_secret="your-api-secret-hex",
        ...     super_admin_keys_pem=["-----BEGIN PUBLIC KEY-----..."],
        ...     min_valid_signatures=2,
        ... ) as client:
        ...     # List wallets
        ...     wallets, _ = client.wallets.list()
        ...     for wallet in wallets:
        ...         print(f"{wallet.name}: {wallet.currency}")
        ...
        ...     # Get an address
        ...     address = client.addresses.get(123)
        ...     print(f"Address: {address.address}")
        ...
        ...     # Approve a request
        ...     request = client.requests.get(456)
        ...     client.requests.approve_request(request, private_key)

    Attributes:
        host: API host URL.
    """

    DEFAULT_TIMEOUT: float = 30.0
    DEFAULT_RULES_CACHE_TTL: float = 300.0  # 5 minutes

    def __init__(
        self,
        host: str,
        auth: TPV1Auth,
        super_admin_keys: List[EllipticCurvePublicKey],
        min_valid_signatures: int,
        rules_cache_ttl: float,
        timeout: float,
    ) -> None:
        """
        Initialize ProtectClient.

        Use the `create()` or `create_from_pem()` factory methods instead
        of calling this constructor directly.
        """
        self.host = host
        self._auth = auth
        self._super_admin_keys = super_admin_keys
        self._min_valid_signatures = min_valid_signatures
        self._rules_cache_ttl = rules_cache_ttl
        self._timeout = timeout
        self._lock = threading.RLock()
        self._closed = False

        # OpenAPI client will be initialized lazily
        self._api_client: Any = None

        # Lazy-initialized service instances
        self._action_service: Optional["ActionService"] = None
        self._address_service: Optional["AddressService"] = None
        self._air_gap_service: Optional["AirGapService"] = None
        self._asset_service: Optional["AssetService"] = None
        self._audit_service: Optional["AuditService"] = None
        self._balance_service: Optional["BalanceService"] = None
        self._blockchain_service: Optional["BlockchainService"] = None
        self._business_rules_service: Optional["BusinessRuleService"] = None
        self._change_service: Optional["ChangeService"] = None
        self._config_service: Optional["ConfigService"] = None
        self._contract_whitelisting_service: Optional["ContractWhitelistingService"] = None
        self._currency_service: Optional["CurrencyService"] = None
        self._exchange_service: Optional["ExchangeService"] = None
        self._fee_payer_service: Optional["FeePayerService"] = None
        self._fee_service: Optional["FeeService"] = None
        self._fiat_service: Optional["FiatService"] = None
        self._governance_rule_service: Optional["GovernanceRuleService"] = None
        self._group_service: Optional["GroupService"] = None
        self._health_service: Optional["HealthService"] = None
        self._job_service: Optional["JobService"] = None
        self._multi_factor_signature_service: Optional["MultiFactorSignatureService"] = None
        self._price_service: Optional["PriceService"] = None
        self._request_service: Optional["RequestService"] = None
        self._reservation_service: Optional["ReservationService"] = None
        self._score_service: Optional["ScoreService"] = None
        self._staking_service: Optional["StakingService"] = None
        self._statistics_service: Optional["StatisticsService"] = None
        self._tag_service: Optional["TagService"] = None
        self._token_metadata_service: Optional["TokenMetadataService"] = None
        self._transaction_service: Optional["TransactionService"] = None
        self._user_device_service: Optional["UserDeviceService"] = None
        self._user_service: Optional["UserService"] = None
        self._visibility_group_service: Optional["VisibilityGroupService"] = None
        self._wallet_service: Optional["WalletService"] = None
        self._webhook_call_service: Optional["WebhookCallService"] = None
        self._webhook_service: Optional["WebhookService"] = None
        self._whitelisted_address_service: Optional["WhitelistedAddressService"] = None
        self._whitelisted_asset_service: Optional["WhitelistedAssetService"] = None
        self._taurus_network_client: Optional["TaurusNetworkClient"] = None

        # Rules container cache for mandatory address signature verification
        self._rules_cache: Any = None

    @classmethod
    def create(
        cls,
        host: str,
        api_key: str,
        api_secret: str,
        super_admin_keys_pem: List[str],
        min_valid_signatures: int = 1,
        rules_cache_ttl: float = DEFAULT_RULES_CACHE_TTL,
        timeout: float = DEFAULT_TIMEOUT,
    ) -> "ProtectClient":
        """
        Create a new ProtectClient instance.

        Args:
            host: API host URL (e.g., "https://api.protect.taurushq.com").
            api_key: API key for authentication.
            api_secret: API secret as hex-encoded string.
            super_admin_keys_pem: List of PEM-encoded SuperAdmin public keys
                for integrity verification. At least one key is required.
            min_valid_signatures: Minimum number of valid SuperAdmin signatures
                required for verification. Cannot exceed number of keys.
            rules_cache_ttl: TTL in seconds for rules container cache.
            timeout: HTTP request timeout in seconds.

        Returns:
            Configured ProtectClient instance.

        Raises:
            ConfigurationError: If configuration is invalid.
            ValueError: If credentials are empty.
        """
        # Validate required fields
        if not host:
            raise ConfigurationError("host cannot be empty")
        if not api_key:
            raise ConfigurationError("api_key cannot be empty")
        if not api_secret:
            raise ConfigurationError("api_secret cannot be empty")

        # Validate SuperAdmin keys are provided
        if not super_admin_keys_pem:
            raise ConfigurationError(
                "super_admin_keys_pem is required: at least one SuperAdmin public key "
                "must be provided for integrity verification"
            )

        # Normalize host URL
        host = host.rstrip("/")

        # Create auth handler
        try:
            auth = TPV1Auth(api_key, api_secret)
        except ValueError as e:
            raise ConfigurationError(str(e)) from e

        # Decode SuperAdmin keys
        try:
            super_admin_keys = decode_public_keys_pem(super_admin_keys_pem)
        except ValueError as e:
            raise ConfigurationError(f"Invalid SuperAdmin key: {e}") from e

        # Validate min_valid_signatures
        if min_valid_signatures < 0:
            raise ConfigurationError("min_valid_signatures cannot be negative")
        if min_valid_signatures > len(super_admin_keys):
            raise ConfigurationError(
                f"min_valid_signatures ({min_valid_signatures}) cannot exceed "
                f"number of SuperAdmin keys ({len(super_admin_keys)})"
            )
        if min_valid_signatures < 1:
            raise ConfigurationError(
                "min_valid_signatures must be greater than zero"
            )

        return cls(
            host=host,
            auth=auth,
            super_admin_keys=super_admin_keys,
            min_valid_signatures=min_valid_signatures,
            rules_cache_ttl=rules_cache_ttl,
            timeout=timeout,
        )

    @classmethod
    def create_from_pem(
        cls,
        host: str,
        api_key: str,
        api_secret_pem: str,
        **kwargs: Any,
    ) -> "ProtectClient":
        """
        Create a ProtectClient from PEM-encoded API secret.

        This is a convenience method for cases where the API secret
        is stored in PEM format rather than hex-encoded.

        Args:
            host: API host URL.
            api_key: API key for authentication.
            api_secret_pem: PEM-encoded API secret.
            **kwargs: Additional arguments passed to create().

        Returns:
            Configured ProtectClient instance.
        """
        # For now, assume PEM secret is just the hex value
        # In practice, you might need to parse PEM format
        return cls.create(host, api_key, api_secret_pem, **kwargs)

    def __enter__(self) -> "ProtectClient":
        """Enter context manager."""
        return self

    def __exit__(
        self,
        exc_type: Any,
        exc_val: Any,
        exc_tb: Any,
    ) -> None:
        """Exit context manager and clean up resources."""
        self.close()

    def close(self) -> None:
        """
        Close the client and release resources.

        Securely wipes credentials from memory.
        """
        with self._lock:
            if not self._closed:
                if self._auth:
                    self._auth.close()
                self._closed = True

    @property
    def is_closed(self) -> bool:
        """Check if the client has been closed."""
        return self._closed

    def _check_not_closed(self) -> None:
        """Raise if client is closed."""
        if self._closed:
            raise RuntimeError("ProtectClient has been closed")

    def _get_api_client(self) -> Any:
        """
        Get or create the OpenAPI client.

        Returns lazily initialized OpenAPI client configured with TPV1 auth.
        The REST client is replaced with an authenticated version that signs
        all requests with TPV1-HMAC-SHA256 (similar to Go SDK's TPV1Transport).
        """
        if self._api_client is None:
            from taurus_protect._internal.openapi import ApiClient, Configuration
            from taurus_protect.crypto.authenticated_rest import AuthenticatedRESTClient

            config = Configuration(host=self.host)

            # Fix datetime format for Go server compatibility.
            # Python's strftime %z produces +0000 but Go expects +00:00 or Z suffix.
            # Use Z suffix format which Go parses correctly.
            # Note: This assumes datetime objects are in UTC when converted.
            config.datetime_format = "%Y-%m-%dT%H:%M:%S.%fZ"

            # Create the ApiClient
            api_client = ApiClient(configuration=config)

            # Replace rest_client with authenticated version (like Go's TPV1Transport)
            # This intercepts all requests and adds TPV1 Authorization header
            api_client.rest_client = AuthenticatedRESTClient(config, self._auth)

            self._api_client = api_client

        return self._api_client

    def _get_rules_cache(self) -> Any:
        """
        Get or create the rules container cache for address signature verification.

        Returns lazily initialized RulesContainerCache. This is required for
        mandatory address signature verification.
        """
        if self._rules_cache is None:
            from taurus_protect.cache.rules_container_cache import RulesContainerCache

            # Get or create the governance rules service first
            gov_service = self.governance_rules
            ttl_ms = int(self._rules_cache_ttl * 1000)  # Convert seconds to ms
            self._rules_cache = RulesContainerCache(gov_service, ttl_ms=ttl_ms)

        return self._rules_cache

    # Service property accessors with lazy initialization

    @property
    def wallets(self) -> "WalletService":
        """
        Access wallet operations.

        Returns:
            WalletService instance for wallet management.
        """
        self._check_not_closed()
        with self._lock:
            if self._wallet_service is None:
                from taurus_protect._internal.openapi import WalletsApi
                from taurus_protect.services.wallet_service import WalletService

                api_client = self._get_api_client()
                wallets_api = WalletsApi(api_client)
                self._wallet_service = WalletService(api_client, wallets_api)
            return self._wallet_service

    @property
    def addresses(self) -> "AddressService":
        """
        Access address operations with mandatory signature verification.

        Returns:
            AddressService instance for address management.
        """
        self._check_not_closed()
        with self._lock:
            if self._address_service is None:
                from taurus_protect._internal.openapi import AddressesApi
                from taurus_protect.services.address_service import AddressService

                api_client = self._get_api_client()
                addresses_api = AddressesApi(api_client)
                rules_cache = self._get_rules_cache()
                self._address_service = AddressService(
                    api_client,
                    addresses_api,
                    rules_cache=rules_cache,
                )
            return self._address_service

    @property
    def requests(self) -> "RequestService":
        """
        Access request operations.

        Returns:
            RequestService instance for transaction request management.
        """
        self._check_not_closed()
        with self._lock:
            if self._request_service is None:
                from taurus_protect._internal.openapi import RequestsApi
                from taurus_protect.services.request_service import RequestService

                api_client = self._get_api_client()
                requests_api = RequestsApi(api_client)
                self._request_service = RequestService(api_client, requests_api)
            return self._request_service

    @property
    def health(self) -> "HealthService":
        """Access health check operations."""
        self._check_not_closed()
        with self._lock:
            if self._health_service is None:
                from taurus_protect._internal.openapi import HealthApi
                from taurus_protect.services.health_service import HealthService

                api_client = self._get_api_client()
                health_api = HealthApi(api_client)
                self._health_service = HealthService(api_client, health_api)
            return self._health_service

    @property
    def actions(self) -> "ActionService":
        """Access action operations."""
        self._check_not_closed()
        with self._lock:
            if self._action_service is None:
                from taurus_protect._internal.openapi import ActionsApi
                from taurus_protect.services.action_service import ActionService

                api_client = self._get_api_client()
                actions_api = ActionsApi(api_client)
                self._action_service = ActionService(api_client, actions_api)
            return self._action_service

    @property
    def air_gap(self) -> "AirGapService":
        """Access air-gap operations."""
        self._check_not_closed()
        with self._lock:
            if self._air_gap_service is None:
                from taurus_protect._internal.openapi import AirGapApi
                from taurus_protect.services.air_gap_service import AirGapService

                api_client = self._get_api_client()
                air_gap_api = AirGapApi(api_client)
                self._air_gap_service = AirGapService(api_client, air_gap_api)
            return self._air_gap_service

    @property
    def assets(self) -> "AssetService":
        """Access asset operations."""
        self._check_not_closed()
        with self._lock:
            if self._asset_service is None:
                from taurus_protect._internal.openapi import AssetsApi
                from taurus_protect.services.asset_service import AssetService

                api_client = self._get_api_client()
                assets_api = AssetsApi(api_client)
                self._asset_service = AssetService(api_client, assets_api)
            return self._asset_service

    @property
    def audits(self) -> "AuditService":
        """Access audit log operations."""
        self._check_not_closed()
        with self._lock:
            if self._audit_service is None:
                from taurus_protect._internal.openapi import AuditApi
                from taurus_protect.services.audit_service import AuditService

                api_client = self._get_api_client()
                audit_api = AuditApi(api_client)
                self._audit_service = AuditService(api_client, audit_api)
            return self._audit_service

    @property
    def balances(self) -> "BalanceService":
        """Access balance operations."""
        self._check_not_closed()
        with self._lock:
            if self._balance_service is None:
                from taurus_protect._internal.openapi import BalancesApi
                from taurus_protect.services.balance_service import BalanceService

                api_client = self._get_api_client()
                balances_api = BalancesApi(api_client)
                self._balance_service = BalanceService(api_client, balances_api)
            return self._balance_service

    @property
    def blockchains(self) -> "BlockchainService":
        """Access blockchain operations."""
        self._check_not_closed()
        with self._lock:
            if self._blockchain_service is None:
                from taurus_protect._internal.openapi import BlockchainApi
                from taurus_protect.services.blockchain_service import BlockchainService

                api_client = self._get_api_client()
                blockchain_api = BlockchainApi(api_client)
                self._blockchain_service = BlockchainService(api_client, blockchain_api)
            return self._blockchain_service

    @property
    def business_rules(self) -> "BusinessRuleService":
        """Access business rules operations."""
        self._check_not_closed()
        with self._lock:
            if self._business_rules_service is None:
                from taurus_protect._internal.openapi import BusinessRulesApi
                from taurus_protect.services.business_rule_service import BusinessRuleService

                api_client = self._get_api_client()
                business_rules_api = BusinessRulesApi(api_client)
                self._business_rules_service = BusinessRuleService(api_client, business_rules_api)
            return self._business_rules_service

    @property
    def changes(self) -> "ChangeService":
        """Access change tracking operations."""
        self._check_not_closed()
        with self._lock:
            if self._change_service is None:
                from taurus_protect._internal.openapi import ChangesApi
                from taurus_protect.services.change_service import ChangeService

                api_client = self._get_api_client()
                changes_api = ChangesApi(api_client)
                self._change_service = ChangeService(api_client, changes_api)
            return self._change_service

    @property
    def config(self) -> "ConfigService":
        """Access configuration operations."""
        self._check_not_closed()
        with self._lock:
            if self._config_service is None:
                from taurus_protect._internal.openapi import ConfigApi
                from taurus_protect.services.config_service import ConfigService

                api_client = self._get_api_client()
                config_api = ConfigApi(api_client)
                self._config_service = ConfigService(api_client, config_api)
            return self._config_service

    @property
    def contract_whitelisting(self) -> "ContractWhitelistingService":
        """Access contract whitelisting operations."""
        self._check_not_closed()
        with self._lock:
            if self._contract_whitelisting_service is None:
                from taurus_protect._internal.openapi import ContractWhitelistingApi
                from taurus_protect.services.contract_whitelisting_service import (
                    ContractWhitelistingService,
                )

                api_client = self._get_api_client()
                contract_whitelisting_api = ContractWhitelistingApi(api_client)
                self._contract_whitelisting_service = ContractWhitelistingService(
                    api_client, contract_whitelisting_api
                )
            return self._contract_whitelisting_service

    @property
    def currencies(self) -> "CurrencyService":
        """Access currency operations."""
        self._check_not_closed()
        with self._lock:
            if self._currency_service is None:
                from taurus_protect._internal.openapi import CurrenciesApi
                from taurus_protect.services.currency_service import CurrencyService

                api_client = self._get_api_client()
                currencies_api = CurrenciesApi(api_client)
                self._currency_service = CurrencyService(api_client, currencies_api)
            return self._currency_service

    @property
    def exchanges(self) -> "ExchangeService":
        """Access exchange operations."""
        self._check_not_closed()
        with self._lock:
            if self._exchange_service is None:
                from taurus_protect._internal.openapi import ExchangeApi
                from taurus_protect.services.exchange_service import ExchangeService

                api_client = self._get_api_client()
                exchange_api = ExchangeApi(api_client)
                self._exchange_service = ExchangeService(api_client, exchange_api)
            return self._exchange_service

    @property
    def fee_payers(self) -> "FeePayerService":
        """Access fee payer operations."""
        self._check_not_closed()
        with self._lock:
            if self._fee_payer_service is None:
                from taurus_protect._internal.openapi import FeePayersApi
                from taurus_protect.services.fee_payer_service import FeePayerService

                api_client = self._get_api_client()
                fee_payers_api = FeePayersApi(api_client)
                self._fee_payer_service = FeePayerService(api_client, fee_payers_api)
            return self._fee_payer_service

    @property
    def fees(self) -> "FeeService":
        """Access fee operations."""
        self._check_not_closed()
        with self._lock:
            if self._fee_service is None:
                from taurus_protect._internal.openapi import FeeApi
                from taurus_protect.services.fee_service import FeeService

                api_client = self._get_api_client()
                fee_api = FeeApi(api_client)
                self._fee_service = FeeService(api_client, fee_api)
            return self._fee_service

    @property
    def fiat(self) -> "FiatService":
        """Access fiat currency operations."""
        self._check_not_closed()
        with self._lock:
            if self._fiat_service is None:
                from taurus_protect._internal.openapi import CurrenciesApi, FiatApi
                from taurus_protect.services.fiat_service import FiatService

                api_client = self._get_api_client()
                fiat_api = FiatApi(api_client)
                currencies_api = CurrenciesApi(api_client)
                self._fiat_service = FiatService(api_client, fiat_api, currencies_api)
            return self._fiat_service

    @property
    def governance_rules(self) -> "GovernanceRuleService":
        """Access governance rules operations."""
        self._check_not_closed()
        with self._lock:
            if self._governance_rule_service is None:
                from taurus_protect._internal.openapi import GovernanceRulesApi
                from taurus_protect.services.governance_rule_service import GovernanceRuleService

                api_client = self._get_api_client()
                governance_rules_api = GovernanceRulesApi(api_client)
                self._governance_rule_service = GovernanceRuleService(
                    api_client,
                    governance_rules_api,
                    super_admin_keys=self._super_admin_keys,
                    min_valid_signatures=self._min_valid_signatures,
                )
            return self._governance_rule_service

    @property
    def groups(self) -> "GroupService":
        """Access group operations."""
        self._check_not_closed()
        with self._lock:
            if self._group_service is None:
                from taurus_protect._internal.openapi import GroupsApi
                from taurus_protect.services.group_service import GroupService

                api_client = self._get_api_client()
                groups_api = GroupsApi(api_client)
                self._group_service = GroupService(api_client, groups_api)
            return self._group_service

    @property
    def jobs(self) -> "JobService":
        """Access job operations."""
        self._check_not_closed()
        with self._lock:
            if self._job_service is None:
                from taurus_protect._internal.openapi import JobsApi
                from taurus_protect.services.job_service import JobService

                api_client = self._get_api_client()
                jobs_api = JobsApi(api_client)
                self._job_service = JobService(api_client, jobs_api)
            return self._job_service

    @property
    def multi_factor_signature(self) -> "MultiFactorSignatureService":
        """Access multi-factor signature operations."""
        self._check_not_closed()
        with self._lock:
            if self._multi_factor_signature_service is None:
                from taurus_protect._internal.openapi import MultiFactorSignatureApi
                from taurus_protect.services.multi_factor_signature_service import (
                    MultiFactorSignatureService,
                )

                api_client = self._get_api_client()
                mfs_api = MultiFactorSignatureApi(api_client)
                self._multi_factor_signature_service = MultiFactorSignatureService(
                    api_client, mfs_api
                )
            return self._multi_factor_signature_service

    @property
    def prices(self) -> "PriceService":
        """Access price operations."""
        self._check_not_closed()
        with self._lock:
            if self._price_service is None:
                from taurus_protect._internal.openapi import PricesApi
                from taurus_protect.services.price_service import PriceService

                api_client = self._get_api_client()
                prices_api = PricesApi(api_client)
                self._price_service = PriceService(api_client, prices_api)
            return self._price_service

    @property
    def reservations(self) -> "ReservationService":
        """Access reservation operations."""
        self._check_not_closed()
        with self._lock:
            if self._reservation_service is None:
                from taurus_protect._internal.openapi import ReservationsApi
                from taurus_protect.services.reservation_service import ReservationService

                api_client = self._get_api_client()
                reservations_api = ReservationsApi(api_client)
                self._reservation_service = ReservationService(api_client, reservations_api)
            return self._reservation_service

    @property
    def scores(self) -> "ScoreService":
        """Access score operations."""
        self._check_not_closed()
        with self._lock:
            if self._score_service is None:
                from taurus_protect._internal.openapi import ScoresApi
                from taurus_protect.services.score_service import ScoreService

                api_client = self._get_api_client()
                scores_api = ScoresApi(api_client)
                self._score_service = ScoreService(api_client, scores_api)
            return self._score_service

    @property
    def staking(self) -> "StakingService":
        """Access staking operations."""
        self._check_not_closed()
        with self._lock:
            if self._staking_service is None:
                from taurus_protect._internal.openapi import StakingApi
                from taurus_protect.services.staking_service import StakingService

                api_client = self._get_api_client()
                staking_api = StakingApi(api_client)
                self._staking_service = StakingService(api_client, staking_api)
            return self._staking_service

    @property
    def statistics(self) -> "StatisticsService":
        """Access statistics operations."""
        self._check_not_closed()
        with self._lock:
            if self._statistics_service is None:
                from taurus_protect._internal.openapi import StatisticsApi
                from taurus_protect.services.statistics_service import StatisticsService

                api_client = self._get_api_client()
                statistics_api = StatisticsApi(api_client)
                self._statistics_service = StatisticsService(api_client, statistics_api)
            return self._statistics_service

    @property
    def tags(self) -> "TagService":
        """Access tag operations."""
        self._check_not_closed()
        with self._lock:
            if self._tag_service is None:
                from taurus_protect._internal.openapi import TagsApi
                from taurus_protect.services.tag_service import TagService

                api_client = self._get_api_client()
                tags_api = TagsApi(api_client)
                self._tag_service = TagService(api_client, tags_api)
            return self._tag_service

    @property
    def token_metadata(self) -> "TokenMetadataService":
        """Access token metadata operations."""
        self._check_not_closed()
        with self._lock:
            if self._token_metadata_service is None:
                from taurus_protect._internal.openapi import TokenMetadataApi
                from taurus_protect.services.token_metadata_service import TokenMetadataService

                api_client = self._get_api_client()
                token_metadata_api = TokenMetadataApi(api_client)
                self._token_metadata_service = TokenMetadataService(api_client, token_metadata_api)
            return self._token_metadata_service

    @property
    def transactions(self) -> "TransactionService":
        """Access transaction operations."""
        self._check_not_closed()
        with self._lock:
            if self._transaction_service is None:
                from taurus_protect._internal.openapi import TransactionsApi
                from taurus_protect.services.transaction_service import TransactionService

                api_client = self._get_api_client()
                transactions_api = TransactionsApi(api_client)
                self._transaction_service = TransactionService(api_client, transactions_api)
            return self._transaction_service

    @property
    def user_devices(self) -> "UserDeviceService":
        """Access user device operations."""
        self._check_not_closed()
        with self._lock:
            if self._user_device_service is None:
                from taurus_protect._internal.openapi import UserDeviceApi
                from taurus_protect.services.user_device_service import UserDeviceService

                api_client = self._get_api_client()
                user_device_api = UserDeviceApi(api_client)
                self._user_device_service = UserDeviceService(api_client, user_device_api)
            return self._user_device_service

    @property
    def users(self) -> "UserService":
        """Access user operations."""
        self._check_not_closed()
        with self._lock:
            if self._user_service is None:
                from taurus_protect._internal.openapi import UsersApi
                from taurus_protect.services.user_service import UserService

                api_client = self._get_api_client()
                users_api = UsersApi(api_client)
                self._user_service = UserService(api_client, users_api)
            return self._user_service

    @property
    def visibility_groups(self) -> "VisibilityGroupService":
        """Access visibility group operations."""
        self._check_not_closed()
        with self._lock:
            if self._visibility_group_service is None:
                from taurus_protect._internal.openapi import RestrictedVisibilityGroupsApi
                from taurus_protect.services.visibility_group_service import VisibilityGroupService

                api_client = self._get_api_client()
                visibility_groups_api = RestrictedVisibilityGroupsApi(api_client)
                self._visibility_group_service = VisibilityGroupService(
                    api_client, visibility_groups_api
                )
            return self._visibility_group_service

    @property
    def webhook_calls(self) -> "WebhookCallService":
        """Access webhook call operations."""
        self._check_not_closed()
        with self._lock:
            if self._webhook_call_service is None:
                from taurus_protect._internal.openapi import WebhookCallsApi
                from taurus_protect.services.webhook_call_service import WebhookCallService

                api_client = self._get_api_client()
                webhook_calls_api = WebhookCallsApi(api_client)
                self._webhook_call_service = WebhookCallService(api_client, webhook_calls_api)
            return self._webhook_call_service

    @property
    def webhooks(self) -> "WebhookService":
        """Access webhook operations."""
        self._check_not_closed()
        with self._lock:
            if self._webhook_service is None:
                from taurus_protect._internal.openapi import WebhooksApi
                from taurus_protect.services.webhook_service import WebhookService

                api_client = self._get_api_client()
                webhooks_api = WebhooksApi(api_client)
                self._webhook_service = WebhookService(api_client, webhooks_api)
            return self._webhook_service

    @property
    def whitelisted_addresses(self) -> "WhitelistedAddressService":
        """Access whitelisted address operations."""
        self._check_not_closed()
        with self._lock:
            if self._whitelisted_address_service is None:
                from taurus_protect._internal.openapi import AddressWhitelistingApi
                from taurus_protect.services.whitelisted_address_service import (
                    WhitelistedAddressService,
                )

                api_client = self._get_api_client()
                whitelisted_addresses_api = AddressWhitelistingApi(api_client)
                self._whitelisted_address_service = WhitelistedAddressService(
                    api_client,
                    whitelisted_addresses_api,
                    super_admin_keys=self._super_admin_keys,
                    min_valid_signatures=self._min_valid_signatures,
                )
            return self._whitelisted_address_service

    @property
    def whitelisted_assets(self) -> "WhitelistedAssetService":
        """Access whitelisted asset operations."""
        self._check_not_closed()
        with self._lock:
            if self._whitelisted_asset_service is None:
                from taurus_protect._internal.openapi import ContractWhitelistingApi
                from taurus_protect.services.whitelisted_asset_service import (
                    WhitelistedAssetService,
                )

                api_client = self._get_api_client()
                whitelisted_assets_api = ContractWhitelistingApi(api_client)
                self._whitelisted_asset_service = WhitelistedAssetService(
                    api_client,
                    whitelisted_assets_api,
                    super_admin_keys=self._super_admin_keys,
                    min_valid_signatures=self._min_valid_signatures,
                )
            return self._whitelisted_asset_service

    @property
    def taurus_network(self) -> "TaurusNetworkClient":
        """
        Access Taurus-NETWORK services.

        Provides namespace-based access to Taurus-NETWORK specific services:
        - participants: Participant management
        - pledges: Pledge lifecycle operations
        - lending: Offers and lending agreements
        - settlements: Settlement operations
        - sharing: Address and asset sharing

        Example:
            >>> client.taurus_network.participants.get_my_participant()
            >>> client.taurus_network.pledges.list_pledges()
            >>> client.taurus_network.sharing.list_shared_addresses()
        """
        self._check_not_closed()
        with self._lock:
            if self._taurus_network_client is None:
                from taurus_protect.services.taurus_network._client import (
                    TaurusNetworkClient,
                )

                api_client = self._get_api_client()
                self._taurus_network_client = TaurusNetworkClient(api_client)
            return self._taurus_network_client
