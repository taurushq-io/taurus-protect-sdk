"""Request service for Taurus-PROTECT SDK."""

from __future__ import annotations

import hmac
import json
from datetime import datetime
from decimal import Decimal
from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from cryptography.hazmat.primitives.asymmetric.ec import EllipticCurvePrivateKey

from taurus_protect.crypto.hashing import calculate_hex_hash
from taurus_protect.crypto.signing import sign_data
from taurus_protect.mappers.request import request_from_dto, requests_from_dto
from taurus_protect.models.pagination import Pagination
from taurus_protect.models.request import (
    CreateExternalTransferRequest,
    CreateInternalTransferRequest,
    ListRequestsOptions,
    Request,
    RequestStatus,
)
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


class RequestService(BaseService):
    """
    Service for managing transaction requests.

    Provides operations for creating, approving, rejecting, and querying
    transaction requests. Requests represent actions to be performed on
    the blockchain, such as transfers between addresses.

    **Important Security Features:**
    - All requests fetched via `get()` have their hash verified using constant-time comparison
    - Request approval uses ECDSA signatures for cryptographic verification

    Example:
        >>> # Create an internal transfer
        >>> request = client.requests.create_internal_transfer(
        ...     from_address_id=123,
        ...     to_address_id=456,
        ...     amount="1000000000000000000",  # 1 ETH in wei
        ... )
        >>>
        >>> # Approve a request with a private key
        >>> signed_count = client.requests.approve_requests([request], private_key)
        >>>
        >>> # Get requests pending approval
        >>> requests, cursor = client.requests.get_for_approval(limit=50)
    """

    def __init__(self, api_client: Any, requests_api: Any) -> None:
        """
        Initialize request service.

        Args:
            api_client: The OpenAPI client instance.
            requests_api: The RequestsAPI service from OpenAPI client.
        """
        super().__init__(api_client)
        self._requests_api = requests_api

    def get(self, request_id: int) -> Request:
        """
        Get a request by ID with mandatory hash verification.

        The hash of the request metadata payload is verified using
        constant-time comparison to prevent timing attacks.

        Args:
            request_id: The request ID to retrieve.

        Returns:
            The verified request.

        Raises:
            ValueError: If request_id is invalid.
            NotFoundError: If request not found.
            IntegrityError: If hash verification fails.
            APIError: If API request fails.
        """
        if request_id <= 0:
            raise ValueError("request_id must be positive")

        try:
            resp = self._requests_api.request_service_get_request(str(request_id))

            result = getattr(resp, "result", None)
            if result is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Request {request_id} not found")

            request = request_from_dto(result)
            if request is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Request {request_id} not found")

            # Mandatory hash verification
            self._verify_request_hash(request)

            return request
        except Exception as e:
            from taurus_protect.errors import APIError, IntegrityError

            if isinstance(e, (APIError, IntegrityError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list(
        self,
        limit: int = 50,
        offset: int = 0,
        from_date: Optional[datetime] = None,
        to_date: Optional[datetime] = None,
        currency_id: Optional[str] = None,
        statuses: Optional[List[RequestStatus]] = None,
    ) -> Tuple[List[Request], Optional[Pagination]]:
        """
        List requests with filtering and pagination.

        Args:
            limit: Maximum number of requests to return.
            offset: Number of requests to skip.
            from_date: Filter requests created after this date.
            to_date: Filter requests created before this date.
            currency_id: Filter by currency ID.
            statuses: Filter by request statuses.

        Returns:
            Tuple of (requests list, pagination info).

        Raises:
            ValueError: If limit or offset are invalid.
            APIError: If API request fails.
        """
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            status_strings = [s.value for s in statuses] if statuses else None
            currencies = [currency_id] if currency_id else None

            # Use non-V2 endpoint which supports offset-based pagination
            resp = self._requests_api.request_service_get_requests(
                limit=str(limit),
                offset=str(offset),
                var_from=from_date,
                to=to_date,
                currencies=currencies,
                statuses=status_strings,
                types=None,
                ids=None,
            )

            result = getattr(resp, "result", None)
            requests = requests_from_dto(result) if result else []

            # Extract pagination from total_items
            total_items = getattr(resp, "total_items", None)
            pagination = None
            if total_items is not None:
                pagination = self._extract_pagination(total_items, offset, limit)

            return requests, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_for_approval(
        self,
        limit: int = 50,
        offset: int = 0,
    ) -> Tuple[List[Request], Optional[Pagination]]:
        """
        Get requests pending approval.

        Args:
            limit: Maximum number of requests to return.
            offset: Number of requests to skip.

        Returns:
            Tuple of (requests list, pagination info).

        Raises:
            ValueError: If limit or offset are invalid.
            APIError: If API request fails.
        """
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            resp = self._requests_api.request_service_get_requests_for_approval_v2(
                currency_id=None,
                types=None,
                ids=None,
                cursor_current_page=None,
                cursor_page_request=str(offset // limit) if limit > 0 else None,
                cursor_page_size=str(limit),
                sort_order=None,
                exclude_types=None,
                statuses=None,
            )

            result = getattr(resp, "result", None)
            requests = requests_from_dto(result) if result else []

            cursor = getattr(resp, "cursor", None)
            pagination = None
            if cursor is not None:
                total = getattr(cursor, "total_items", None)
                pagination = self._extract_pagination(total, offset, limit)

            return requests, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def approve_requests(
        self,
        requests: List[Request],
        private_key: EllipticCurvePrivateKey,
        comment: str = "approving via taurus-protect-sdk-python",
    ) -> int:
        """
        Approve multiple requests with ECDSA signature.

        The requests are sorted by ID, and a JSON array of their hashes
        is signed using the provided private key.

        Args:
            requests: List of requests to approve.
            private_key: ECDSA private key for signing.
            comment: Optional approval comment.

        Returns:
            Number of requests successfully signed.

        Raises:
            ValueError: If requests list is empty or invalid.
            APIError: If API request fails or signing fails.
        """
        if not requests:
            raise ValueError("requests list cannot be empty")
        if private_key is None:
            raise ValueError("private_key cannot be None")

        # Validate all requests have metadata with hash
        for r in requests:
            if r.metadata is None:
                raise ValueError("request metadata cannot be None")
            if r.metadata.hash is None or r.metadata.hash == "":
                raise ValueError("request metadata hash cannot be None or empty")

        try:
            # Sort requests by ID (numeric sort)
            sorted_requests = sorted(requests, key=lambda r: int(r.id))

            # Build JSON array of hashes
            hashes = [r.metadata.hash for r in sorted_requests]
            to_sign = json.dumps(hashes, separators=(",", ":"))

            # Sign with ECDSA
            signature = sign_data(private_key, to_sign.encode("utf-8"))

            # Build API request
            body = {
                "ids": [r.id for r in sorted_requests],
                "comment": comment,
                "signature": signature,
            }

            resp = self._requests_api.request_service_approve_requests(body=body)

            signed_requests = getattr(resp, "signed_requests", None)
            if signed_requests is not None:
                return int(signed_requests)
            return 0
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def approve_request(
        self,
        request: Request,
        private_key: EllipticCurvePrivateKey,
        comment: str = "approving via taurus-protect-sdk-python",
    ) -> int:
        """
        Approve a single request with ECDSA signature.

        Args:
            request: The request to approve.
            private_key: ECDSA private key for signing.
            comment: Optional approval comment.

        Returns:
            Number of requests successfully signed (0 or 1).

        Raises:
            ValueError: If request is invalid.
            APIError: If API request fails or signing fails.
        """
        return self.approve_requests([request], private_key, comment)

    def reject_requests(
        self,
        request_ids: List[int],
        comment: str,
    ) -> None:
        """
        Reject multiple requests.

        Args:
            request_ids: List of request IDs to reject.
            comment: Rejection comment (required).

        Raises:
            ValueError: If request_ids is empty or comment is missing.
            APIError: If API request fails.
        """
        if not request_ids:
            raise ValueError("request_ids list cannot be empty")
        self._validate_required(comment, "comment")

        try:
            body = {
                "ids": [str(rid) for rid in request_ids],
                "comment": comment,
            }
            self._requests_api.request_service_reject_requests(body=body)
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def reject_request(
        self,
        request_id: int,
        comment: str,
    ) -> None:
        """
        Reject a single request.

        Args:
            request_id: The request ID to reject.
            comment: Rejection comment (required).

        Raises:
            ValueError: If request_id is invalid or comment is missing.
            APIError: If API request fails.
        """
        self.reject_requests([request_id], comment)

    def create_internal_transfer(
        self,
        from_address_id: int,
        to_address_id: int,
        amount: str,
    ) -> Request:
        """
        Create an internal transfer request between addresses.

        Args:
            from_address_id: Source address ID.
            to_address_id: Destination address ID.
            amount: Transfer amount as string (to preserve precision).

        Returns:
            The created request.

        Raises:
            ValueError: If arguments are invalid.
            APIError: If API request fails.
        """
        if from_address_id <= 0:
            raise ValueError("from_address_id must be positive")
        if to_address_id <= 0:
            raise ValueError("to_address_id must be positive")
        self._validate_required(amount, "amount")
        if Decimal(amount) <= 0:
            raise ValueError("amount must be positive")

        try:
            body = {
                "from_address_id": str(from_address_id),
                "to_address_id": str(to_address_id),
                "amount": amount,
            }
            resp = self._requests_api.request_service_create_outgoing_request(body=body)

            result = getattr(resp, "result", None)
            if result is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to create request: no result returned")

            request = request_from_dto(result)
            if request is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to create request: invalid response")

            return request
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def create_internal_transfer_from_wallet(
        self,
        from_wallet_id: int,
        to_address_id: int,
        amount: str,
    ) -> Request:
        """
        Create an internal transfer from an omnibus wallet.

        Args:
            from_wallet_id: Source omnibus wallet ID.
            to_address_id: Destination address ID.
            amount: Transfer amount as string.

        Returns:
            The created request.

        Raises:
            ValueError: If arguments are invalid.
            APIError: If API request fails.
        """
        if from_wallet_id <= 0:
            raise ValueError("from_wallet_id must be positive")
        if to_address_id <= 0:
            raise ValueError("to_address_id must be positive")
        self._validate_required(amount, "amount")
        if Decimal(amount) <= 0:
            raise ValueError("amount must be positive")

        try:
            body = {
                "from_wallet_id": str(from_wallet_id),
                "to_address_id": str(to_address_id),
                "amount": amount,
            }
            resp = self._requests_api.request_service_create_outgoing_request(body=body)

            result = getattr(resp, "result", None)
            if result is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to create request: no result returned")

            return request_from_dto(result)
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def create_external_transfer(
        self,
        from_address_id: int,
        to_whitelisted_address_id: int,
        amount: str,
    ) -> Request:
        """
        Create an external transfer to a whitelisted address.

        Args:
            from_address_id: Source address ID.
            to_whitelisted_address_id: Destination whitelisted address ID.
            amount: Transfer amount as string.

        Returns:
            The created request.

        Raises:
            ValueError: If arguments are invalid.
            APIError: If API request fails.
        """
        if from_address_id <= 0:
            raise ValueError("from_address_id must be positive")
        if to_whitelisted_address_id <= 0:
            raise ValueError("to_whitelisted_address_id must be positive")
        self._validate_required(amount, "amount")
        if Decimal(amount) <= 0:
            raise ValueError("amount must be positive")

        try:
            body = {
                "from_address_id": str(from_address_id),
                "to_whitelisted_address_id": str(to_whitelisted_address_id),
                "amount": amount,
            }
            resp = self._requests_api.request_service_create_outgoing_request(body=body)

            result = getattr(resp, "result", None)
            if result is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to create request: no result returned")

            return request_from_dto(result)
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def create_external_transfer_from_wallet(
        self,
        from_wallet_id: int,
        to_whitelisted_address_id: int,
        amount: str,
    ) -> Request:
        """
        Create an external transfer from an omnibus wallet.

        Args:
            from_wallet_id: Source omnibus wallet ID.
            to_whitelisted_address_id: Destination whitelisted address ID.
            amount: Transfer amount as string.

        Returns:
            The created request.

        Raises:
            ValueError: If arguments are invalid.
            APIError: If API request fails.
        """
        if from_wallet_id <= 0:
            raise ValueError("from_wallet_id must be positive")
        if to_whitelisted_address_id <= 0:
            raise ValueError("to_whitelisted_address_id must be positive")
        self._validate_required(amount, "amount")
        if Decimal(amount) <= 0:
            raise ValueError("amount must be positive")

        try:
            body = {
                "from_wallet_id": str(from_wallet_id),
                "to_whitelisted_address_id": str(to_whitelisted_address_id),
                "amount": amount,
            }
            resp = self._requests_api.request_service_create_outgoing_request(body=body)

            result = getattr(resp, "result", None)
            if result is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to create request: no result returned")

            return request_from_dto(result)
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def create_cancel_request(
        self,
        address_id: int,
        nonce: int,
    ) -> Request:
        """
        Create a cancel request for a pending transaction.

        Args:
            address_id: The address ID.
            nonce: The nonce of the transaction to cancel.

        Returns:
            The created cancel request.

        Raises:
            ValueError: If arguments are invalid.
            APIError: If API request fails.
        """
        if address_id <= 0:
            raise ValueError("address_id must be positive")
        if nonce < 0:
            raise ValueError("nonce cannot be negative")

        try:
            body = {
                "address_id": str(address_id),
                "nonce": str(nonce),
            }
            resp = self._requests_api.request_service_create_outgoing_cancel_request(body=body)

            result = getattr(resp, "result", None)
            if result is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to create cancel request: no result")

            return request_from_dto(result)
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def create_incoming_request(
        self,
        from_exchange_id: int,
        to_address_id: int,
        amount: str,
    ) -> Request:
        """
        Create an incoming request from an exchange.

        Args:
            from_exchange_id: Source exchange ID.
            to_address_id: Destination address ID.
            amount: Transfer amount as string.

        Returns:
            The created incoming request.

        Raises:
            ValueError: If arguments are invalid.
            APIError: If API request fails.
        """
        if from_exchange_id <= 0:
            raise ValueError("from_exchange_id must be positive")
        if to_address_id <= 0:
            raise ValueError("to_address_id must be positive")
        self._validate_required(amount, "amount")
        if Decimal(amount) <= 0:
            raise ValueError("amount must be positive")

        try:
            body = {
                "from_exchange_id": str(from_exchange_id),
                "to_address_id": str(to_address_id),
                "amount": amount,
            }
            resp = self._requests_api.request_service_create_incoming_request(body=body)

            result = getattr(resp, "result", None)
            if result is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to create incoming request: no result")

            return request_from_dto(result)
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def _verify_request_hash(self, request: Request) -> None:
        """
        Verify the hash of a request using constant-time comparison.

        Args:
            request: The request to verify.

        Raises:
            IntegrityError: If hash verification fails.
        """
        if request.metadata is None:
            return

        provided_hash = request.metadata.hash
        payload = request.metadata.payload_as_string

        if not provided_hash and not payload:
            return

        if not payload:
            if provided_hash:
                from taurus_protect.errors import IntegrityError

                raise IntegrityError(
                    "request hash verification failed: hash exists but payload is missing"
                )
            return

        # Compute hash of the payload
        computed_hash = calculate_hex_hash(payload)

        # Explicit null check before constant-time comparison
        if provided_hash is None:
            from taurus_protect.errors import IntegrityError

            raise IntegrityError("request hash verification failed: provided hash is null")

        # Use constant-time comparison to prevent timing attacks
        if not hmac.compare_digest(computed_hash, provided_hash):
            from taurus_protect.errors import IntegrityError

            raise IntegrityError(
                f"request hash verification failed: computed={computed_hash}, provided={provided_hash}"
            )
