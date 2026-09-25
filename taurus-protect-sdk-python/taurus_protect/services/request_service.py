"""Request service for Taurus-PROTECT SDK."""

from __future__ import annotations

import hmac
import json
import logging
from datetime import datetime
from decimal import Decimal
from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from cryptography.hazmat.primitives.asymmetric.ec import EllipticCurvePrivateKey

from taurus_protect.crypto.hashing import calculate_hex_hash
from taurus_protect.crypto.signing import sign_data
from taurus_protect.errors import APIError, IntegrityError
from taurus_protect.mappers.request import request_from_dto
from taurus_protect.models.pagination import CursorPage, cursor_page, cursor_request
from taurus_protect.models.request import (
    CreateExternalTransferRequest,
    CreateInternalTransferRequest,
    Request,
    RequestStatus,
)
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


# Module logger: the Python idiom. Metadata only -- id and reason, never the payload.
_LOGGER = logging.getLogger(__name__)


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
        >>> requests, page = client.requests.get_for_approval(page_size=100)
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

            # Mandatory hash verification, and mark the result so a caller can
            # tell a verified request from one that merely came back.
            return self._verified_request(result)
        except Exception as e:

            # IntegrityError is a plain Exception, NOT an APIError, so a funnel that
            # omits it hands the failure to _handle_error, which maps an unknown
            # exception to ServerError(500) -- and ServerError.is_retryable() is True.
            # The documented contract of IntegrityError is "security-critical, NEVER
            # retry", so the omission inverts it.
            #
            # On the six create paths this is worse than a taxonomy slip: the request is
            # ALREADY created server-side, so _verified_request puts the id in the error
            # precisely so the caller can reconcile "rather than retrying and creating a
            # second request". Reporting it as retryable invites exactly that.
            if isinstance(e, (APIError, IntegrityError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list(
        self,
        page_size: Optional[int] = None,
        cursor: Optional[str] = None,
        from_date: Optional[datetime] = None,
        to_date: Optional[datetime] = None,
        currency_id: Optional[str] = None,
        statuses: Optional[List[RequestStatus]] = None,
        *,
        types: Optional[List[str]] = None,
        ids: Optional[List[str]] = None,
        external_request_ids: Optional[List[str]] = None,
        sort_order: Optional[str] = None,
        current_page: Optional[str] = None,
        page_request: Optional[str] = None,
    ) -> Tuple[List[Request], CursorPage]:
        """
        List requests with filtering, one page at a time.

        Args:
            page_size: Page size (default 20, max 100).
            cursor: ``page.next_cursor`` from the previous page, to continue.
            from_date: Filter requests created after this date.
            to_date: Filter requests created before this date.
            currency_id: Filter by currency ID.
            statuses: Filter by request statuses.
            types: Filter by request types.
            ids: Filter by request IDs.
            external_request_ids: Filter by external request IDs.
            sort_order: ASC or DESC.
            current_page: Low-level page token; not with ``cursor``.
            page_request: Low-level page direction (FIRST, PREVIOUS, NEXT, LAST).

        Returns:
            Tuple of (verified requests, page). Rows whose metadata fails verification
            are withheld and logged.

        Raises:
            ValueError: If the page size is invalid or cursor options conflict.
            APIError: If API request fails.
        """
        req = cursor_request(
            page_size, cursor, current_page=current_page, page_request=page_request
        )

        try:
            resp = self._requests_api.request_service_get_requests_v2(
                var_from=from_date,
                to=to_date,
                currency_id=currency_id,
                statuses=[s.value for s in statuses] if statuses else None,
                types=types,
                ids=ids,
                external_request_ids=external_request_ids,
                sort_order=sort_order,
                **req.query_params(),
            )

            requests = self._verified_requests(resp.result)
            return requests, cursor_page(req.page_size, resp.cursor)
        except Exception as e:

            # See get(): IntegrityError is not an APIError, so omitting it here turns a
            # security failure into a retryable ServerError(500).
            if isinstance(e, (APIError, IntegrityError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_for_approval(
        self,
        page_size: Optional[int] = None,
        cursor: Optional[str] = None,
        *,
        currency_id: Optional[str] = None,
        types: Optional[List[str]] = None,
        exclude_types: Optional[List[str]] = None,
        ids: Optional[List[str]] = None,
        sort_order: Optional[str] = None,
        current_page: Optional[str] = None,
        page_request: Optional[str] = None,
    ) -> Tuple[List[Request], CursorPage]:
        """
        List requests pending the caller's approval, one page at a time.

        The approval queue has no status filter (every row is awaiting approval) and no
        ``external_request_ids`` filter: validatord's approval paginator drops it before the query.

        Args:
            page_size: Page size (default 20, max 100).
            cursor: ``page.next_cursor`` from the previous page, to continue.
            currency_id: Filter by currency ID.
            types: Only these request types.
            exclude_types: Leave out these request types.
            ids: Filter by request IDs.
            sort_order: ASC or DESC.
            current_page: Low-level page token; not with ``cursor``.
            page_request: Low-level page direction (FIRST, PREVIOUS, NEXT, LAST).

        Returns:
            Tuple of (verified requests, page).

        Raises:
            ValueError: If the page size is invalid or cursor options conflict.
            APIError: If API request fails.
        """
        req = cursor_request(
            page_size, cursor, current_page=current_page, page_request=page_request
        )

        try:
            resp = self._requests_api.request_service_get_requests_for_approval_v2(
                currency_id=currency_id,
                types=types,
                exclude_types=exclude_types,
                ids=ids,
                sort_order=sort_order,
                **req.query_params(),
            )

            requests = self._verified_requests(resp.result)
            return requests, cursor_page(req.page_size, resp.cursor)
        except Exception as e:

            # See get(): IntegrityError is not an APIError, so omitting it here turns a
            # security failure into a retryable ServerError(500).
            if isinstance(e, (APIError, IntegrityError, ValueError)):
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
            # RE-VERIFY the hash this signature will attest to. ``hash_verified`` is
            # deliberately NOT consulted: it is a plain model field, so a Request
            # rebuilt from JSON (a queue, a webhook, a cached blob) can arrive
            # claiming True and get an attacker-chosen hash signed by the approver's
            # real key. Re-verification is a SHA-256 over a string already in hand,
            # and it is the same rule ``WhitelistedAssetService.approve`` gets by
            # re-reading -- so both signing paths rest on one rule, not two.
            # Ordered after the presence check so an early-status request still
            # reports the clearer absence message.
            try:
                self._verify_request_hash(r)
            except IntegrityError as exc:
                raise IntegrityError(f"refusing to sign request {r.id}: {exc}") from exc

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

            # See get(): IntegrityError is not an APIError, so omitting it here turns a
            # security failure into a retryable ServerError(500).
            if isinstance(e, (APIError, IntegrityError, ValueError)):
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

            # See get(): IntegrityError is not an APIError, so omitting it here turns a
            # security failure into a retryable ServerError(500).
            if isinstance(e, (APIError, IntegrityError, ValueError)):
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

                raise APIError(500, "Failed to create request: no result returned")

            return self._verified_request(result)
        except Exception as e:

            # See get(): IntegrityError is not an APIError, so omitting it here turns a
            # security failure into a retryable ServerError(500).
            if isinstance(e, (APIError, IntegrityError, ValueError)):
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

                raise APIError(500, "Failed to create request: no result returned")

            return self._verified_request(result)
        except Exception as e:

            # See get(): IntegrityError is not an APIError, so omitting it here turns a
            # security failure into a retryable ServerError(500).
            if isinstance(e, (APIError, IntegrityError, ValueError)):
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

                raise APIError(500, "Failed to create request: no result returned")

            return self._verified_request(result)
        except Exception as e:

            # See get(): IntegrityError is not an APIError, so omitting it here turns a
            # security failure into a retryable ServerError(500).
            if isinstance(e, (APIError, IntegrityError, ValueError)):
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

                raise APIError(500, "Failed to create request: no result returned")

            return self._verified_request(result)
        except Exception as e:

            # See get(): IntegrityError is not an APIError, so omitting it here turns a
            # security failure into a retryable ServerError(500).
            if isinstance(e, (APIError, IntegrityError, ValueError)):
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

                raise APIError(500, "Failed to create cancel request: no result")

            return self._verified_request(result)
        except Exception as e:

            # See get(): IntegrityError is not an APIError, so omitting it here turns a
            # security failure into a retryable ServerError(500).
            if isinstance(e, (APIError, IntegrityError, ValueError)):
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

                raise APIError(500, "Failed to create incoming request: no result")

            return self._verified_request(result)
        except Exception as e:

            # See get(): IntegrityError is not an APIError, so omitting it here turns a
            # security failure into a retryable ServerError(500).
            if isinstance(e, (APIError, IntegrityError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def _verified_request(self, dto: Any) -> Request:
        """
        Map one request DTO and verify its metadata. The ONLY way this module builds
        a ``Request``.

        get() ---------------+
        list/for_approval ---+
        create_* (x6) -------+--> _verified_request -> request_from_dto
                                                    -> _verify_and_mark
                                        ok <--------+--------> IntegrityError
                                (hash_verified set)          (names the request id)

        Verification lived in ``get()`` alone. The list paths skipped it in one pass
        and the six create paths in the next -- both times because "remember to
        verify" was a rule rather than the only available construction path. So do
        NOT call ``request_from_dto`` / ``requests_from_dto`` anywhere else here.

        The id goes into the error because a failing create has already succeeded
        server-side; the caller needs the id to reconcile rather than retrying and
        creating a second request.

        Raises:
            NotFoundError: If the DTO maps to nothing.
            IntegrityError: If the metadata hash does not cover the payload.
        """
        request = request_from_dto(dto)
        if request is None:
            from taurus_protect.errors import NotFoundError

            raise NotFoundError("request not found")
        try:
            return self._verify_and_mark(request)
        except IntegrityError as exc:
            # ``exc.message``, not ``str(exc)``: IntegrityError.__str__ prefixes its own
            # class name, so re-wrapping the string form yields
            # "IntegrityError: request 2: IntegrityError: ...".
            raise IntegrityError(f"request {request.id}: {exc.message}") from exc

    def _verified_requests(self, dtos: Optional[List[Any]]) -> List[Request]:
        """
        Keep only the rows whose metadata integrity verifies.

        rows -> _verified_request -+- ok    -> kept, marked hash_verified
                                   +- fails -> excluded + logged

        The list endpoints previously returned every row without verifying any of
        them, so a payload altered in transit reached the caller with no error and
        no flag -- while ``get()`` verified. Excluding rather than raising keeps one
        bad row from denying access to every good one; logging keeps a shortened
        list from passing as a complete one.

        Takes DTOs rather than models, so the mapper stays behind the seam.
        """
        kept: List[Request] = []
        for dto in dtos or []:
            if dto is None:
                continue
            try:
                kept.append(self._verified_request(dto))
            except IntegrityError as exc:
                _LOGGER.warning(
                    "request excluded: metadata integrity verification failed "
                    "(request_id=%s, reason=%s)",
                    getattr(dto, "id", None),
                    exc.message,
                )
                continue
        return kept

    def _verify_and_mark(self, request: Request) -> Request:
        """
        Verify a request's metadata hash and return it marked verified.

        The single place verification and the flag are set together. Keeping them
        apart is what let ``get()`` verify and then return a request whose
        ``hash_verified`` was still False, so a caller checking the flag could not
        tell a verified request from an unverified one.

        Metadata is frozen, so a verified row is rebuilt with ``model_copy`` rather
        than mutated.

        Args:
            request: The request to verify.

        Returns:
            The request, with ``metadata.hash_verified`` set when it carries metadata.

        Raises:
            IntegrityError: If hash verification fails.
        """
        self._verify_request_hash(request)

        if request.metadata is not None and (
            request.metadata.hash or request.metadata.payload_as_string
        ):
            request = request.model_copy(
                update={"metadata": request.metadata.model_copy(update={"hash_verified": True})}
            )
        return request

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
                raise IntegrityError(
                    "request hash verification failed: hash exists but payload is missing"
                )
            return

        # Compute hash of the payload
        computed_hash = calculate_hex_hash(payload)

        # Explicit null check before constant-time comparison
        if provided_hash is None:
            raise IntegrityError("request hash verification failed: provided hash is null")

        # Use constant-time comparison to prevent timing attacks
        if not hmac.compare_digest(computed_hash, provided_hash):
            raise IntegrityError(
                f"request hash verification failed: computed={computed_hash}, provided={provided_hash}"
            )
