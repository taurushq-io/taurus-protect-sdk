"""User device service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect.mappers.user_device import (
    user_device_pairing_from_dto,
    user_device_pairing_info_from_dto,
)
from taurus_protect.models.pagination import Pagination
from taurus_protect.models.user_device import UserDevicePairing, UserDevicePairingInfo
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


class UserDeviceService(BaseService):
    """
    Service for user device management operations.

    Provides methods to manage user device pairing for multi-factor
    authentication, including creating, starting, and approving pairings.

    Note: The API does not provide traditional list/get operations.
    The pairing workflow consists of:
    1. Create pairing request (get pairing ID)
    2. Start pairing (device provides encryption key and nonce)
    3. Approve pairing (complete the pairing process)
    4. Get pairing status (check if pairing is complete)

    Example:
        >>> # Create a new pairing request
        >>> pairing = client.user_devices.create_pairing()
        >>> print(f"Pairing ID: {pairing.pairing_id}")
        >>>
        >>> # Get pairing status
        >>> info = client.user_devices.get_pairing_status(pairing_id, nonce)
        >>> print(f"Status: {info.status}")
    """

    def __init__(self, api_client: Any, user_device_api: Any) -> None:
        """
        Initialize user device service.

        Args:
            api_client: The OpenAPI client instance.
            user_device_api: The UserDeviceApi service from OpenAPI client.
        """
        super().__init__(api_client)
        self._user_device_api = user_device_api

    def list(
        self,
        limit: int = 50,
        offset: int = 0,
    ) -> Tuple[List[UserDevicePairing], Optional[Pagination]]:
        """
        List user device pairings.

        Note: The underlying API does not support listing pairings.
        This method is provided for interface consistency but will
        return an empty list.

        Args:
            limit: Maximum number of pairings to return (must be positive).
            offset: Number of pairings to skip (must be non-negative).

        Returns:
            Tuple of (empty pairings list, None pagination).

        Raises:
            ValueError: If limit or offset are invalid.
        """
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        # The API does not support listing user device pairings
        # Return empty list for interface consistency
        return [], None

    def get(self, device_id: str) -> UserDevicePairingInfo:
        """
        Get user device pairing by ID.

        Note: The underlying API requires a nonce to get pairing status.
        Use get_pairing_status() for the full API functionality.

        Args:
            device_id: The device/pairing ID to retrieve.

        Raises:
            NotFoundError: Always raised as the API requires a nonce.
        """
        self._validate_required(device_id, "device_id")

        from taurus_protect.errors import NotFoundError

        raise NotFoundError(
            f"Device {device_id} not found. Use get_pairing_status(pairing_id, nonce) instead."
        )

    def create_pairing(self) -> UserDevicePairing:
        """
        Create a new user device pairing request (Step 1).

        Creates a pairing request for the current user. Returns a pairing ID
        that should be used by the user's device to complete the pairing process.

        Returns:
            The created pairing request with pairing ID.

        Raises:
            APIError: If API request fails.
        """
        try:
            resp = self._user_device_api.user_device_service_create_user_device_pairing()

            pairing = user_device_pairing_from_dto(resp)
            if pairing is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to create pairing: invalid response")

            return pairing
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def start_pairing(
        self,
        pairing_id: str,
        nonce: str,
        encryption_key: str,
    ) -> None:
        """
        Start a user device pairing request (Step 2).

        Starts a specific pairing process. The nonce and user's device
        encryption key will be used to validate future requests.

        Args:
            pairing_id: The pairing ID from create_pairing().
            nonce: A 6-digit number provided by the device.
            encryption_key: The device's encryption public key.

        Raises:
            ValueError: If required arguments are invalid.
            APIError: If API request fails.
        """
        self._validate_required(pairing_id, "pairing_id")
        self._validate_required(nonce, "nonce")
        self._validate_required(encryption_key, "encryption_key")

        try:
            from taurus_protect._internal.openapi.models.user_device_service_start_user_device_pairing_body import (
                UserDeviceServiceStartUserDevicePairingBody,
            )

            body = UserDeviceServiceStartUserDevicePairingBody(
                nonce=nonce,
                encryption_key=encryption_key,
            )
            self._user_device_api.user_device_service_start_user_device_pairing(
                pairing_id=pairing_id,
                body=body,
            )
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def approve_pairing(
        self,
        pairing_id: str,
        nonce: str,
    ) -> None:
        """
        Approve a user device pairing request (Step 3).

        Final approval of the pairing request. Checks the validity of the
        nonce and completes the device pairing.

        Args:
            pairing_id: The pairing ID from create_pairing().
            nonce: The same 6-digit number used in start_pairing().

        Raises:
            ValueError: If required arguments are invalid.
            APIError: If API request fails.
        """
        self._validate_required(pairing_id, "pairing_id")
        self._validate_required(nonce, "nonce")

        try:
            from taurus_protect._internal.openapi.models.user_device_service_approve_user_device_pairing_body import (
                UserDeviceServiceApproveUserDevicePairingBody,
            )

            body = UserDeviceServiceApproveUserDevicePairingBody(
                nonce=nonce,
            )
            self._user_device_api.user_device_service_approve_user_device_pairing(
                pairing_id=pairing_id,
                body=body,
            )
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_pairing_status(
        self,
        pairing_id: str,
        nonce: str,
    ) -> UserDevicePairingInfo:
        """
        Get the status of a user device pairing request.

        Returns the status of a pairing request identified by its pairing ID
        and the nonce used to start the pairing.

        Args:
            pairing_id: The pairing ID from create_pairing().
            nonce: A 6-digit number used to identify the pairing session.

        Returns:
            The pairing status information.

        Raises:
            ValueError: If required arguments are invalid.
            NotFoundError: If pairing not found.
            APIError: If API request fails.
        """
        self._validate_required(pairing_id, "pairing_id")
        self._validate_required(nonce, "nonce")

        try:
            resp = self._user_device_api.user_device_service_get_user_device_pairing_status(
                pairing_id=pairing_id,
                nonce=nonce,
            )

            info = user_device_pairing_info_from_dto(resp)
            if info is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Pairing {pairing_id} not found")

            return info
        except Exception as e:
            from taurus_protect.errors import APIError, NotFoundError

            if isinstance(e, (APIError, NotFoundError, ValueError)):
                raise
            raise self._handle_error(e) from e
