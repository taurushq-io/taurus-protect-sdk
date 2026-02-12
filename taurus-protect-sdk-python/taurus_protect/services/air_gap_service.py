"""Air-gap service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any

from taurus_protect.models.staking import UnsignedPayload
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    from taurus_protect._internal.openapi.api.air_gap_api import AirGapApi


class AirGapService(BaseService):
    """
    Service for air-gap (cold storage) signing operations.

    Air-gap signing is used in high-security environments where signing keys
    are stored on offline devices. This service provides methods to:

    1. Get unsigned transaction payloads for offline signing
    2. Submit signed payloads back to the system

    Example:
        >>> # Get unsigned payload for a request
        >>> unsigned = client.air_gap.get_unsigned_payload(request_id=123)
        >>> print(f"Payload hash: {unsigned.hash}")
        >>>
        >>> # After signing offline, submit the signed payload
        >>> client.air_gap.submit_signed_payload(
        ...     request_id=123,
        ...     signed_payload="0x1234..."
        ... )
    """

    def __init__(self, api_client: Any, air_gap_api: "AirGapApi") -> None:
        """
        Initialize air-gap service.

        Args:
            api_client: The OpenAPI client instance.
            air_gap_api: The AirGapApi service from OpenAPI client.
        """
        super().__init__(api_client)
        self._api = air_gap_api

    def get_unsigned_payload(self, request_id: int) -> bytes:
        """
        Get unsigned transaction payload for offline signing.

        Retrieves the raw transaction data that needs to be signed by a cold HSM
        or offline signing device.

        Args:
            request_id: The request ID to get the unsigned payload for.

        Returns:
            The unsigned payload as bytes.

        Raises:
            ValueError: If request_id is invalid.
            NotFoundError: If the request is not found.
            APIError: If the API request fails.

        Example:
            >>> payload = client.air_gap.get_unsigned_payload(request_id=123)
            >>> # Transfer payload to offline signing device
            >>> print(f"Payload size: {len(payload)} bytes")
        """
        if request_id <= 0:
            raise ValueError("request_id must be positive")

        try:
            # Build request body
            body = {"request_ids": [str(request_id)]}

            # The API returns raw bytes for the cold HSM
            resp = self._api.air_gap_service_get_outgoing_air_gap(body=body)

            # Response is a bytearray
            if resp is None:
                return bytes()

            if isinstance(resp, bytearray):
                return bytes(resp)
            if isinstance(resp, bytes):
                return resp

            # Try to get bytes from response object
            return bytes(resp) if resp else bytes()
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def submit_signed_payload(self, request_id: int, signed_payload: bytes) -> None:
        """
        Submit a signed payload back to the system.

        After signing the transaction offline (e.g., with a cold HSM),
        use this method to submit the signed data back to Taurus-PROTECT.

        Args:
            request_id: The request ID the signed payload belongs to.
            signed_payload: The signed payload bytes from the offline signer.

        Raises:
            ValueError: If request_id is invalid or signed_payload is empty.
            NotFoundError: If the request is not found.
            ValidationError: If the signature is invalid.
            APIError: If the API request fails.

        Example:
            >>> # After receiving signed payload from cold HSM
            >>> signed_data = bytes.fromhex("1234abcd...")
            >>> client.air_gap.submit_signed_payload(
            ...     request_id=123,
            ...     signed_payload=signed_data
            ... )
        """
        if request_id <= 0:
            raise ValueError("request_id must be positive")
        if not signed_payload:
            raise ValueError("signed_payload cannot be empty")

        try:
            # Build request body with the signed payload
            body = {"payload": signed_payload}

            self._api.air_gap_service_submit_incoming_air_gap(body=body)
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
