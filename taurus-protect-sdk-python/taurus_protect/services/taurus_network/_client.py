"""TaurusNetwork client namespace for Taurus-PROTECT SDK."""

from __future__ import annotations

import threading
from typing import TYPE_CHECKING, Any, Optional

if TYPE_CHECKING:
    from taurus_protect.services.taurus_network.lending_service import LendingService
    from taurus_protect.services.taurus_network.participant_service import (
        ParticipantService,
    )
    from taurus_protect.services.taurus_network.pledge_service import PledgeService
    from taurus_protect.services.taurus_network.settlement_service import (
        SettlementService,
    )
    from taurus_protect.services.taurus_network.sharing_service import SharingService


class TaurusNetworkClient:
    """
    Namespace client for Taurus-NETWORK services.

    Provides access to TaurusNetwork-specific services through lazy-initialized
    properties. Access via the main ProtectClient:

    Example:
        >>> with ProtectClient.create(...) as client:
        ...     # Get my participant info
        ...     me = client.taurus_network.participants.get_my_participant()
        ...
        ...     # List pledges
        ...     pledges, _ = client.taurus_network.pledges.list_pledges()
        ...
        ...     # List shared addresses
        ...     addresses, _ = client.taurus_network.sharing.list_shared_addresses()
    """

    def __init__(self, api_client: Any) -> None:
        """
        Initialize TaurusNetworkClient.

        Args:
            api_client: The underlying OpenAPI client instance.
        """
        self._api_client = api_client
        self._lock = threading.RLock()

        # Lazy-initialized service instances
        self._participant_service: Optional["ParticipantService"] = None
        self._pledge_service: Optional["PledgeService"] = None
        self._lending_service: Optional["LendingService"] = None
        self._settlement_service: Optional["SettlementService"] = None
        self._sharing_service: Optional["SharingService"] = None

    @property
    def participants(self) -> "ParticipantService":
        """Access participant operations."""
        with self._lock:
            if self._participant_service is None:
                from taurus_protect._internal.openapi import TaurusNetworkParticipantApi
                from taurus_protect.services.taurus_network.participant_service import (
                    ParticipantService,
                )

                participant_api = TaurusNetworkParticipantApi(self._api_client)
                self._participant_service = ParticipantService(self._api_client, participant_api)
            return self._participant_service

    @property
    def pledges(self) -> "PledgeService":
        """Access pledge operations."""
        with self._lock:
            if self._pledge_service is None:
                from taurus_protect._internal.openapi import TaurusNetworkPledgeApi
                from taurus_protect.services.taurus_network.pledge_service import (
                    PledgeService,
                )

                pledge_api = TaurusNetworkPledgeApi(self._api_client)
                self._pledge_service = PledgeService(self._api_client, pledge_api)
            return self._pledge_service

    @property
    def lending(self) -> "LendingService":
        """Access lending operations (offers and agreements)."""
        with self._lock:
            if self._lending_service is None:
                from taurus_protect._internal.openapi import TaurusNetworkLendingApi
                from taurus_protect.services.taurus_network.lending_service import (
                    LendingService,
                )

                lending_api = TaurusNetworkLendingApi(self._api_client)
                self._lending_service = LendingService(self._api_client, lending_api)
            return self._lending_service

    @property
    def settlements(self) -> "SettlementService":
        """Access settlement operations."""
        with self._lock:
            if self._settlement_service is None:
                from taurus_protect._internal.openapi import TaurusNetworkSettlementApi
                from taurus_protect.services.taurus_network.settlement_service import (
                    SettlementService,
                )

                settlement_api = TaurusNetworkSettlementApi(self._api_client)
                self._settlement_service = SettlementService(self._api_client, settlement_api)
            return self._settlement_service

    @property
    def sharing(self) -> "SharingService":
        """Access shared address and asset operations."""
        with self._lock:
            if self._sharing_service is None:
                from taurus_protect._internal.openapi import (
                    TaurusNetworkSharedAddressAssetApi,
                )
                from taurus_protect.services.taurus_network.sharing_service import (
                    SharingService,
                )

                shared_api = TaurusNetworkSharedAddressAssetApi(self._api_client)
                self._sharing_service = SharingService(self._api_client, shared_api)
            return self._sharing_service
