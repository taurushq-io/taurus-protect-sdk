"""Address service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, Dict, List, Optional, Tuple

from taurus_protect.errors import APIError, ContainerIntegrityError, IntegrityError, NotFoundError
from taurus_protect.mappers.address import address_from_dto
from taurus_protect.models.address import Address, CreateAddressRequest, ListAddressesOptions
from taurus_protect.models.pagination import (
    REPLY_OFFSET,
    Pagination,
    offset_pagination,
    offset_query,
    resolve_offset,
    resolve_page_size,
)
from taurus_protect.services._base import BaseService

# validatord caps an addressIds filter at 50.
MAX_ADDRESS_IDS = 50

if TYPE_CHECKING:
    from taurus_protect.cache.rules_container_cache import RulesContainerCache
    from taurus_protect.helpers.address_signature_verifier import AddressSignatureVerifier


class AddressService(BaseService):
    """
    Service for managing blockchain addresses.

    Provides operations for creating, retrieving, and managing addresses
    within wallets. All addresses retrieved through this service are
    automatically verified for cryptographic integrity using the rules
    container public keys. Signature verification is mandatory.

    Example:
        >>> # Create a new address
        >>> request = CreateAddressRequest(
        ...     wallet_id="123",
        ...     label="Customer Deposit",
        ...     comment="Primary deposit address",
        ... )
        >>> address = client.addresses.create(request)
        >>>
        >>> # Get an address (signature verified automatically)
        >>> address = client.addresses.get(456)
        >>> print(f"Address: {address.address}")
        >>>
        >>> # List addresses for a wallet (all signatures verified)
        >>> addresses, pagination = client.addresses.list(wallet_id=123, limit=100)
    """

    def __init__(
        self,
        api_client: Any,
        addresses_api: Any,
        rules_cache: "RulesContainerCache",
    ) -> None:
        """
        Initialize address service with mandatory signature verification.

        Args:
            api_client: The OpenAPI client instance.
            addresses_api: The AddressesAPI service from OpenAPI client.
            rules_cache: Rules container cache for signature verification.
                This is required - address signature verification is mandatory.

        Raises:
            ValueError: If rules_cache is None.
        """
        if rules_cache is None:
            raise ValueError(
                "rules_cache cannot be None - address signature verification is mandatory"
            )
        super().__init__(api_client)
        self._addresses_api = addresses_api
        self._rules_cache = rules_cache

    def get(self, address_id: int) -> Address:
        """
        Get an address by ID with mandatory signature verification.

        Args:
            address_id: The address ID to retrieve.

        Returns:
            The verified address.

        Raises:
            ValueError: If address_id is invalid.
            NotFoundError: If address not found.
            IntegrityError: If signature verification fails.
            APIError: If API request fails.
        """
        if address_id <= 0:
            raise ValueError("address_id must be positive")

        try:
            resp = self._addresses_api.wallet_service_get_address(str(address_id))

            result = getattr(resp, "result", None)
            if result is None:

                raise NotFoundError(f"Address {address_id} not found")

            address = self._verified_address(result)
            if address is None:

                raise NotFoundError(f"Address {address_id} not found")

            return address
        except Exception as e:

            # IntegrityError is NOT an APIError, so without naming it here a failed HSM
            # signature check would be remapped to a retryable ServerError(500) and a
            # caller following is_retryable() would retry a suspected forgery.
            if isinstance(e, (APIError, IntegrityError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def _verified_address(
        self, dto: Any, rules_container: Optional[Any] = None
    ) -> Optional[Address]:
        """
        The ONE construction seam for an :class:`Address`.

        Every path that returns an Address goes through here -- ``get``, ``list``,
        ``list_with_options``, ``create_address`` and ``AssetService.get_addresses`` --
        because "remember to verify" was a rule rather than the only available
        construction path, and ``create_address`` is what that cost: ``get_addresses``
        had already been fixed for exactly this, and the create path was missed.

        Asynchronous creation is the one case that is NOT simply verify-and-throw. The
        DTO carries a ``status`` of ``created``/``creating``/``signed``/``observed``/
        ``confirmed``, so a reply can legitimately arrive before the HSM has signed the
        address. The rule that holds either way: **never hand back a non-empty
        ``Address.address`` that has not been verified.** So a signature present must
        verify, and a signature absent means the address string is withheld rather than
        returned unchecked -- the caller re-reads through this same seam once the status
        advances.

        The branch is on the address STRING, not on ``status``: status is
        server-controlled, so keying the decision on it would let a response claim
        ``creating`` while handing over an attacker-chosen destination.

        Args:
            dto: The generated address DTO.
            rules_container: Optional pre-fetched container, to avoid an N+1 cache
                lookup when verifying a page.

        Returns:
            The verified address, or None when the DTO maps to nothing.

        Raises:
            IntegrityError: If the signature does not verify, or an address string
                arrived with no signature to check it against.
        """
        from taurus_protect.helpers.address_signature_verifier import verified_address

        address = address_from_dto(dto)
        if address is None:
            return None

        # Only pay for a container fetch when there is something to verify -- an
        # address still being generated carries no destination. The DECISION itself
        # stays in the helper, so this cannot drift from AssetService.get_addresses.
        if rules_container is None and address.address:
            rules_container = self._rules_cache.get_decoded_rules_container()
        return verified_address(address, rules_container)

    def list(
        self,
        wallet_id: int,
        limit: Optional[int] = None,
        offset: Optional[int] = None,
        exclude_disabled: Optional[bool] = None,
    ) -> Tuple[List[Address], Pagination]:
        """
        List a wallet's addresses, one page at a time, with mandatory signature verification.

        Args:
            wallet_id: The wallet ID to list addresses for.
            limit: Page size (default 20, max 100).
            offset: Number of addresses to skip; pass ``pagination.next_offset`` to continue.
            exclude_disabled: True hides disabled addresses, False includes them.

        Returns:
            Tuple of (addresses, pagination).

        Raises:
            ValueError: If wallet_id is invalid or limit/offset are invalid.
            IntegrityError: If signature verification fails for any address.
            APIError: If API request fails.
        """
        if wallet_id <= 0:
            raise ValueError("wallet_id must be positive")
        return self.list_with_options(
            ListAddressesOptions(
                wallet_id=str(wallet_id),
                limit=limit,
                offset=offset,
                exclude_disabled=exclude_disabled,
            )
        )

    def list_with_options(
        self,
        options: Optional[ListAddressesOptions] = None,
    ) -> Tuple[List[Address], Pagination]:
        """
        List addresses with filters, one page at a time, with mandatory signature verification.

        Args:
            options: Filters and page window; every field reaches the wire.

        Returns:
            Tuple of (addresses, pagination).

        Raises:
            ValueError: If limit or offset are invalid.
            IntegrityError: If signature verification fails for any address.
            APIError: If API request fails.
        """
        opts = options or ListAddressesOptions()
        limit = resolve_page_size(opts.limit, "limit")
        offset = resolve_offset(opts.offset)

        include_disabled = None
        if opts.exclude_disabled is not None:
            include_disabled = "exclude" if opts.exclude_disabled else "include"

        try:
            resp = self._addresses_api.wallet_service_get_addresses(
                query=opts.query,
                wallet_id=opts.wallet_id,
                include_disabled_addresses=include_disabled,
                **offset_query(limit, offset),
            )

            rows = resp.result or []

            # Every row through the one seam. Pre-fetch the rules container once to
            # avoid an N+1 cache lookup.
            rules_container = None
            if rows:
                rules_container = self._rules_cache.get_decoded_rules_container()
            addresses = []
            for dto in rows:
                address = self._verified_address(dto, rules_container)
                if address is not None:
                    addresses.append(address)

            pagination = offset_pagination(
                REPLY_OFFSET,
                limit=limit,
                offset=offset,
                served_rows=len(rows),
                total_items=resp.total_items,
                reply_offset=resp.offset,
                excluded=len(rows) - len(addresses),
            )
            return addresses, pagination
        except Exception as e:

            # IntegrityError is NOT an APIError, so without naming it here a failed HSM
            # signature check would be remapped to a retryable ServerError(500) and a
            # caller following is_retryable() would retry a suspected forgery.
            if isinstance(e, (APIError, IntegrityError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def _verified_addresses_by_id(
        self, address_ids: List[str]
    ) -> Tuple[Dict[str, Address], Dict[str, str]]:
        """
        Re-read managed addresses by id, every returned row through :meth:`_verified_address`.

        One request per ``MAX_ADDRESS_IDS`` ids. A row whose signature is missing or does
        not verify is reported under its id, not raised, so the caller decides without
        it. A rules container that cannot verify any address is not a property of one
        row, so it aborts the call, as does any request error.

        Args:
            address_ids: The address ids to read.

        Returns:
            The verified addresses by id, and the failure reason by id.

        Raises:
            IntegrityError: If the rules container has no HSMSLOT key.
            APIError: If an API request fails.
        """
        verified: Dict[str, Address] = {}
        failed: Dict[str, str] = {}
        try:
            rules_container = self._rules_cache.get_decoded_rules_container()
            if rules_container is None or rules_container.get_hsm_public_key() is None:
                raise IntegrityError(
                    "the rules container has no HSMSLOT key, so no address can be verified"
                )
            for start in range(0, len(address_ids), MAX_ADDRESS_IDS):
                chunk = address_ids[start : start + MAX_ADDRESS_IDS]
                resp = self._addresses_api.wallet_service_get_addresses(
                    address_ids=chunk, limit=str(len(chunk))
                )
                for dto in resp.result or []:
                    try:
                        address = self._verified_address(dto, rules_container)
                    except ContainerIntegrityError:
                        raise
                    except IntegrityError as exc:
                        failed[str(dto.id)] = exc.message
                        continue
                    if address is not None and address.id:
                        verified[str(address.id)] = address
        except Exception as e:
            if isinstance(e, (APIError, IntegrityError, ValueError)):
                raise
            raise self._handle_error(e) from e
        return verified, failed

    def create(self, request: CreateAddressRequest) -> Address:
        """
        Create a new address.

        Args:
            request: Address creation parameters.

        Returns:
            The created address.

        Raises:
            ValueError: If required fields are missing.
            ValidationError: If request is invalid.
            APIError: If API request fails.
        """
        if request is None:
            raise ValueError("request cannot be None")
        self._validate_required(request.wallet_id, "wallet_id")
        self._validate_required(request.label, "label")

        return self.create_address(
            wallet_id=int(request.wallet_id),
            label=request.label,
            comment=request.comment or "",
            customer_id=request.customer_id or "",
        )

    def create_address(
        self,
        wallet_id: int,
        label: str,
        comment: str = "",
        customer_id: str = "",
    ) -> Address:
        """
        Create a new address with explicit parameters.

        Args:
            wallet_id: The wallet ID to create the address in.
            label: Human-readable label for the address.
            comment: Optional description.
            customer_id: Optional customer identifier.

        Returns:
            The created address.

        Raises:
            ValueError: If required fields are missing.
            APIError: If API request fails.
        """
        if wallet_id <= 0:
            raise ValueError("wallet_id must be positive")
        self._validate_required(label, "label")
        self._validate_required(comment, "comment")

        try:
            body = {
                "wallet_id": str(wallet_id),
                "label": label,
                "comment": comment,
            }
            if customer_id:
                body["customer_id"] = customer_id

            resp = self._addresses_api.wallet_service_create_address(body=body)

            result = getattr(resp, "result", None)
            if result is None:

                raise APIError(500, "Failed to create address: no result returned")

            # Through the same seam as every read. The create reply carries the same
            # generated address DTO the read paths return, signature included, so there
            # was never a reason for this path to be the unverified one -- and callers
            # of create are precisely the ones about to publish or fund a fresh deposit
            # address.
            address = self._verified_address(result)
            if address is None:

                raise APIError(500, "Failed to create address: invalid response")

            return address
        except Exception as e:

            # See list(): an unverifiable freshly-created address must not be reported
            # as a retryable server error -- retrying would create a second address.
            if isinstance(e, (APIError, IntegrityError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def create_attribute(
        self,
        address_id: int,
        key: str,
        value: str,
    ) -> None:
        """
        Create an attribute for an address.

        Args:
            address_id: The address ID.
            key: The attribute key.
            value: The attribute value.

        Raises:
            ValueError: If any argument is invalid.
            APIError: If API request fails.
        """
        if address_id <= 0:
            raise ValueError("address_id must be positive")
        self._validate_required(key, "key")
        self._validate_required(value, "value")

        try:
            body = {
                "attributes": [{"key": key, "value": value}],
            }
            self._addresses_api.wallet_service_create_address_attributes(str(address_id), body=body)
        except Exception as e:

            if isinstance(e, (APIError, IntegrityError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def delete_attribute(
        self,
        address_id: int,
        attribute_id: int,
    ) -> None:
        """
        Delete an attribute from an address.

        Args:
            address_id: The address ID.
            attribute_id: The attribute ID to delete.

        Raises:
            ValueError: If any argument is invalid.
            APIError: If API request fails.
        """
        if address_id <= 0:
            raise ValueError("address_id must be positive")
        if attribute_id <= 0:
            raise ValueError("attribute_id must be positive")

        try:
            self._addresses_api.wallet_service_delete_address_attribute(
                str(address_id), str(attribute_id)
            )
        except Exception as e:

            if isinstance(e, (APIError, IntegrityError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_proof_of_reserve(
        self,
        address_id: int,
        challenge: Optional[str] = None,
    ) -> Any:
        """
        Get the proof of reserve for an address.

        Args:
            address_id: The address ID.
            challenge: Optional challenge string.

        Returns:
            The proof of reserve response.

        Raises:
            ValueError: If address_id is invalid.
            APIError: If API request fails.
        """
        if address_id <= 0:
            raise ValueError("address_id must be positive")

        try:
            resp = self._addresses_api.wallet_service_get_address_proof_of_reserve(
                str(address_id), challenge
            )
            return getattr(resp, "result", None)
        except Exception as e:

            if isinstance(e, (APIError, IntegrityError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def _verify_address_signature(
        self, address: Address, rules_container: Optional[Any] = None
    ) -> None:
        """
        Verify the signature of an address using the rules container.

        Address signature verification is mandatory and always performed.

        Args:
            address: The address to verify.
            rules_container: Optional pre-fetched rules container. If None,
                will be fetched from cache. Pass this when verifying multiple
                addresses to avoid N+1 cache lookups.

        Raises:
            IntegrityError: If signature verification fails.
        """
        from taurus_protect.helpers.address_signature_verifier import (
            verify_address_signature,
        )

        if rules_container is None:
            rules_container = self._rules_cache.get_decoded_rules_container()
        verify_address_signature(address, rules_container)
