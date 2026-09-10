"""Whitelisted address service for Taurus-PROTECT SDK."""

from __future__ import annotations

import logging

import base64
import binascii
import hashlib
import json
from typing import Any, Dict, List, Optional, TYPE_CHECKING, Tuple

from cryptography.hazmat.primitives.asymmetric.ec import EllipticCurvePublicKey

from taurus_protect.crypto.signing import sign_data
from taurus_protect._internal.openapi.exceptions import ApiException
from taurus_protect.errors import ContainerIntegrityError, APIError, IntegrityError, WhitelistError
from taurus_protect.helpers.signature_verifier import verify_governance_rules_signatures
from taurus_protect.helpers.whitelisted_address_verifier import WhitelistedAddressVerifier
from taurus_protect.mappers.governance_rules import (
    rules_container_from_base64,
    user_signatures_from_base64,
)
from taurus_protect.models.governance_rules import DecodedRulesContainer
from taurus_protect.models.whitelisted_address import (
    ExcludedWhitelistedAddress,
    InternalWallet,
    SignedWhitelistedAddress,
    SignedWhitelistedAddressEnvelope,
    WhitelistedAddress,
    WhitelistedAddressListResult,
    WhitelistMetadata,
    WhitelistSignature,
    WhitelistSignatureEntry,
    WhitelistUserSignature,
)
from taurus_protect.services._base import BaseService

_LOGGER = logging.getLogger(__name__)


def _container_hash_label(container_base64: str) -> str:
    """
    Recompute the label validatord files a normalized rules container under.

    A row's ``rules_container_hash`` then resolves only to bytes that really hash to it.
    The convention is validatord's, in its whitelist controller::

        h := base64.StdEncoding.EncodeToString(crypto.Sha256([]byte(e.GetRulesContainer())))

    Note what is hashed. ``GetRulesContainer()`` is ALREADY a base64 string there, so the
    digest is over the base64 TEXT, not over the decoded protobuf, and the output is
    base64 rather than hex. Decoding the container first and hashing the protobuf yields
    a different value and would reject every container, breaking all list calls.

    Do not confuse it with ``enforcedRulesHash``, which is ``base64(SHA256(raw
    protobuf))`` and is a backlink to a ruleset's predecessor rather than its own
    identity. Both are 44-character base64 SHA-256 digests, so mixing them up is silent.

    Args:
        container_base64: the container exactly as the response carried it.

    Returns:
        The label those bytes should be filed under.
    """
    digest = hashlib.sha256(container_base64.encode("utf-8")).digest()
    return base64.b64encode(digest).decode("ascii")


if TYPE_CHECKING:
    from taurus_protect._internal.openapi.api.address_whitelisting_api import (
        AddressWhitelistingApi,
    )


class WhitelistedAddressService(BaseService):
    """
    Service for managing whitelisted addresses.

    Whitelisted addresses are pre-approved external destinations for withdrawals.
    This service performs cryptographic verification of whitelisted addresses
    using the 6-step verification flow:
    1. Metadata hash verification
    2. Rules container signature verification (SuperAdmin keys)
    3. Decode rules container
    4. Verify hash coverage
    5. Whitelist signature verification (per governance rules)
    6. Parse WhitelistedAddress from verified payload

    Example:
        >>> # Get a verified whitelisted address
        >>> address = client.whitelisted_addresses.get(123)
        >>> print(f"{address.label}: {address.address}")
        ...
        >>> # List whitelisted addresses
        >>> addresses, _ = client.whitelisted_addresses.list(limit=50)
    """

    def __init__(
        self,
        api_client: Any,
        whitelisting_api: "AddressWhitelistingApi",
        super_admin_keys: List[EllipticCurvePublicKey],
        min_valid_signatures: int,
    ) -> None:
        """
        Initialize the whitelisted address service.

        Args:
            api_client: The OpenAPI client instance.
            whitelisting_api: The address whitelisting API instance.
            super_admin_keys: List of SuperAdmin public keys for verification.
            min_valid_signatures: Minimum valid signatures required.
        """
        super().__init__(api_client)
        self._api = whitelisting_api
        self._verifier = WhitelistedAddressVerifier(
            super_admin_keys=super_admin_keys,
            min_valid_signatures=min_valid_signatures,
        )

    def get(self, whitelisted_address_id: int) -> WhitelistedAddress:
        """
        Get a whitelisted address by ID with verification.

        Performs cryptographic verification of the address envelope
        before returning the address.

        Args:
            whitelisted_address_id: The whitelisted address ID.

        Returns:
            The verified whitelisted address.

        Raises:
            IntegrityError: If verification fails.
            APIError: If the API call fails.
        """
        if whitelisted_address_id <= 0:
            raise ValueError("whitelisted_address_id must be positive")

        envelope = self.get_envelope(whitelisted_address_id)
        return envelope.verified_whitelisted_address or WhitelistedAddress(id=str(whitelisted_address_id))

    def get_envelope(self, whitelisted_address_id: int) -> SignedWhitelistedAddressEnvelope:
        """
        Get the signed envelope for a whitelisted address.

        Performs full 6-step verification.

        Args:
            whitelisted_address_id: The whitelisted address ID.

        Returns:
            The verified signed envelope.

        Raises:
            IntegrityError: If verification fails.
            APIError: If the API call fails.
        """
        if whitelisted_address_id <= 0:
            raise ValueError("whitelisted_address_id must be positive")

        try:
            reply = self._api.whitelist_service_get_whitelisted_address(str(whitelisted_address_id))
            result = reply.result
            if result is None:
                raise APIError(f"Whitelisted address {whitelisted_address_id} not found")

            envelope = self._map_envelope_from_dto(result)
            self._verify_and_populate_envelope(envelope, result)
            return envelope
        except IntegrityError:
            raise
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

    def list(
        self,
        currency: Optional[str] = None,
        limit: int = 50,
        offset: int = 0,
        *,
        ids: Optional[List[str]] = None,
        include_for_approval: bool = False,
    ) -> WhitelistedAddressListResult:
        """
        List whitelisted addresses with cryptographic verification.

        Each address is verified using the 6-step verification flow.

        Args:
            currency: Filter by currency.
            limit: Maximum number of addresses to return.
            offset: Offset for pagination.

        Returns:
            The verified addresses, the page window with total_items reduced by the
            number of excluded rows, and the excluded rows themselves.

        Raises:
            IntegrityError: If verification fails for any address.
            APIError: If the API call fails.
        """
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            reply = self._api.whitelist_service_get_whitelisted_addresses(
                currency=currency,
                limit=str(limit),
                offset=str(offset),
                rules_container_normalized=True,
                ids=ids or None,
                include_for_approval=include_for_approval or None,
            )

            # Build rules container cache from normalized containers
            rules_container_cache = self._build_rules_container_cache(reply)

            # rows -> verify -+- ok    -> returned
            #                  +- fails -> excluded + logged, list survives
            #
            # One unverifiable row used to fail the whole call, which took down the
            # whitelist for every consumer -- and listing is how an operator would
            # find the bad row, so the failure hid its own cause. Excluding stays
            # fail-closed: an omitted destination cannot be selected.
            rows = list(reply.result or [])
            envelopes, excluded = self._verified_addresses(rows, rules_container_cache)
            addresses = [e.verified_whitelisted_address for e in envelopes]

            # The server counts rows it returned; the caller receives only those that
            # verified. Reporting the server's total lets a filtered page pass for a
            # complete one, and makes has_more promise a page never fully readable.
            pagination = self._extract_pagination(
                getattr(reply, "total_items", None),
                offset,
                limit,
            )
            if pagination is not None and pagination.total_items is not None:
                pagination = pagination.model_copy(
                    update={"total_items": max(0, pagination.total_items - len(excluded))}
                )

            return WhitelistedAddressListResult(
                addresses=addresses,
                pagination=pagination,
                excluded_unverified=excluded,
            )
        except IntegrityError:
            raise
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

    def list_for_approval(
        self,
        limit: int = 50,
        offset: int = 0,
        *,
        ids: Optional[List[str]] = None,
        include_already_signed_by_user: bool = False,
    ) -> WhitelistedAddressListResult:
        """
        List whitelisted addresses awaiting approval, verified exactly as ``list`` is.

        Both this endpoint and the approval endpoint below are generated in all four SDK
        clients and were wrapped by none of them, so the rows an approver reads before
        whitelisting a destination were not reachable through the SDK at all -- verified
        or not.

        Args:
            limit: Maximum number of rows to return.
            offset: Offset for pagination.
            ids: Filter by specific whitelisted address IDs.
            include_already_signed_by_user: Include rows the calling user has already
                signed, which is how an approver tells "waiting for me" from "waiting
                for someone else".

        Returns:
            The verified rows, the page window, and the rows withheld.

        Raises:
            IntegrityError: If the container cannot be interpreted, or no row survived.
            APIError: If the API call fails.
        """
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            reply = self._api.whitelist_service_get_whitelisted_addresses_for_approval(
                limit=str(limit),
                offset=str(offset),
                ids=ids or None,
                include_already_signed_by_user=include_already_signed_by_user or None,
            )

            # This endpoint has no normalized-container mode, so the per-row containers
            # are used and the cache starts empty.
            rows = list(reply.result or [])
            envelopes, excluded = self._verified_addresses(rows, {})
            addresses = [e.verified_whitelisted_address for e in envelopes]

            pagination = self._extract_pagination(
                getattr(reply, "total_items", None), offset, limit
            )
            if pagination is not None and pagination.total_items is not None:
                pagination = pagination.model_copy(
                    update={"total_items": max(0, pagination.total_items - len(excluded))}
                )

            return WhitelistedAddressListResult(
                addresses=addresses,
                pagination=pagination,
                excluded_unverified=excluded,
            )
        except IntegrityError:
            raise
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

    def approve(
        self,
        ids: List[int],
        private_key: Any,
        comment: str,
    ) -> None:
        """
        Sign and submit an approval for the given whitelisted addresses, all-or-nothing.

        The batch is re-read through the verifying path and the hashes THOSE rows carry
        are what gets signed, so the approver's signature covers metadata this SDK
        checked rather than whatever a caller was handed. Same shape as
        ``WhitelistedAssetService.approve``.

        ids -> sort numerically -> ONE filtered verified page -> completeness check
                                                                        |
                                                                        v
                                                    sign(JSON(hashes)) -> POST once

        Any address that is missing or fails verification aborts the whole call and
        nothing is signed: one signature covers every hash in the batch, so a partial
        approval would mean the caller believes they approved more than they did.

        Args:
            ids: The whitelisted address IDs to approve.
            private_key: The approver's P-256 private key.
            comment: The approval comment.

        Raises:
            ValueError: If any argument is missing or malformed.
            IntegrityError: If any address is missing, unverifiable, or has no hash.
            APIError: If the API call fails.
        """
        if not ids:
            raise ValueError("ids cannot be empty")
        if private_key is None:
            raise ValueError("private_key is required")
        if not comment:
            raise ValueError("comment is required")
        for address_id in ids:
            if not isinstance(address_id, int) or isinstance(address_id, bool) or address_id <= 0:
                raise ValueError(
                    f"whitelisted address ID {address_id!r} must be a positive integer"
                )

        # The endpoint requires ascending order, and sorting here also makes the signed
        # order independent of the order the caller passed.
        id_strings = [str(i) for i in sorted(ids)]

        # ONE id-filtered page through the verifying path, not one GET per id.
        try:
            reply = self._api.whitelist_service_get_whitelisted_addresses(
                limit=str(len(id_strings)),
                offset="0",
                rules_container_normalized=True,
                ids=id_strings,
                include_for_approval=True,
            )
            cache = self._build_rules_container_cache(reply)
            envelopes, _ = self._verified_addresses(list(reply.result or []), cache)
        except IntegrityError:
            raise
        except Exception as e:
            raise IntegrityError(f"refusing to sign: the verified read failed: {e}") from e

        by_id = {
            str(e.verified_whitelisted_address.id): e
            for e in envelopes
            if e.verified_whitelisted_address is not None
        }

        hashes: List[str] = []
        for address_id in id_strings:
            addr = by_id.get(address_id)
            if addr is None:
                # Excluded, or simply absent. Either way a page that omits a row must
                # not become an approval of fewer rows than the caller asked for.
                raise IntegrityError(
                    f"refusing to sign: address {address_id} was not returned by the "
                    "verified read"
                )
            if addr.metadata is None or not addr.metadata.hash:
                raise IntegrityError(
                    f"refusing to sign: address {address_id} has no metadata hash"
                )
            hashes.append(addr.metadata.hash)

        to_sign = json.dumps(hashes, separators=(",", ":"))
        signature = sign_data(private_key, to_sign.encode("utf-8"))

        try:
            from taurus_protect._internal.openapi.models.tgvalidatord_approve_whitelisted_address_request import (
                TgvalidatordApproveWhitelistedAddressRequest,
            )

            body = TgvalidatordApproveWhitelistedAddressRequest(
                ids=id_strings, signature=signature, comment=comment
            )
            self._api.whitelist_service_approve_whitelisted_address(body=body)
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

    def _verified_addresses(
        self,
        rows: List[Any],
        rules_container_cache: Dict[str, DecodedRulesContainer],
    ) -> Tuple[List[SignedWhitelistedAddressEnvelope], List[ExcludedWhitelistedAddress]]:
        """
        Verify a page of address envelopes; return the survivors plus the rows withheld.

        Returns ENVELOPES, not addresses: this SDK's ``WhitelistedAddress`` carries no
        metadata (Go's does), so the hash an approver signs is only reachable here.

        rows -> verify -+- ok    -> returned
                        +- fails -> excluded + logged, list survives

        LENIENT per row, but two things abort the whole call: a ContainerIntegrityError
        (it invalidates every row judged against that container) and
        rows-returned-but-none-surviving.

        Shared by ``list`` and ``list_for_approval``. The for-approval endpoint had no
        verified reader at all, and duplicating this loop into a second one is precisely
        how the list paths drifted from ``get()`` in the first place.
        """
        verified: List[SignedWhitelistedAddressEnvelope] = []
        excluded: List[ExcludedWhitelistedAddress] = []

        for dto in rows:
            dto_id = getattr(dto, "id", None)
            try:
                envelope = self._map_envelope_from_dto(dto)

                cached = None
                if envelope.rules_container_hash:
                    cached = rules_container_cache.get(envelope.rules_container_hash)

                self._verify_and_populate_envelope(
                    envelope, dto, cached_rules_container=cached
                )
                if envelope.verified_whitelisted_address:
                    verified.append(envelope)
            except ContainerIntegrityError:
                # Not a property of this row: the container is uninterpretable, so
                # every row judged against it is equally unverifiable. Excluding
                # them one by one would empty the whitelist and report success.
                raise
            except (IntegrityError, WhitelistError) as exc:
                excluded.append(
                    ExcludedWhitelistedAddress(
                        id=str(dto_id) if dto_id is not None else None,
                        reason=str(exc),
                    )
                )
                _LOGGER.warning(
                    "whitelisted address excluded: verification failed "
                    "(id=%s, reason=%s)",
                    dto_id,
                    exc,
                )

        # Rows came back but none survived: a systemic failure, not an empty
        # whitelist, and returning [] would be indistinguishable from one.
        if rows and not verified:
            raise IntegrityError(
                f"all {len(rows)} whitelisted address(es) failed verification; "
                f"first failure: {excluded[0].reason if excluded else 'unknown'}"
            )

        return verified, excluded

    def _verify_and_populate_envelope(
        self,
        envelope: SignedWhitelistedAddressEnvelope,
        dto: Optional[Any] = None,
        cached_rules_container: Optional[DecodedRulesContainer] = None,
    ) -> None:
        """
        Verify the envelope using the 6-step verifier and populate the address from payload.

        Args:
            envelope: The envelope to verify.
            dto: Optional DTO to supplement non-security fields not in the payload.
            cached_rules_container: Pre-verified and decoded rules container.
                When provided, steps 2-3 are skipped (already done during cache building).
        """
        # Run the 6-step verification via WhitelistedAddressVerifier
        # Step 6 parses the WhitelistedAddress from the verified payload
        result = self._verifier.verify_whitelisted_address(
            envelope,
            rules_container_decoder=rules_container_from_base64,
            user_signatures_decoder=user_signatures_from_base64,
            cached_rules_container=cached_rules_container,
        )

        # Merge payload fields from verification result with DTO-sourced non-security fields
        verified_addr = result.verified_whitelisted_address

        # Extract created_at from DTO trails (find 'created' action)
        created_at = None
        if dto:
            trails = getattr(dto, "trails", None) or []
            for trail in trails:
                if getattr(trail, "action", None) == "created":
                    created_at = getattr(trail, "var_date", None)
                    break

        # Extract attributes from DTO as dict
        attributes_dict: Dict[str, Any] = {}
        if dto:
            dto_attrs = getattr(dto, "attributes", None) or []
            for attr in dto_attrs:
                key = getattr(attr, "key", None)
                value = getattr(attr, "value", None)
                if key:
                    attributes_dict[key] = value

        envelope.verified_whitelisted_address = WhitelistedAddress(
            # Non-security: ID from payload, fallback to DTO
            id=verified_addr.id or str(getattr(dto, "id", "") if dto else ""),
            # Security-critical: from verified payload only
            address=verified_addr.address,
            label=verified_addr.label,
            currency=verified_addr.currency,
            contract_type=verified_addr.contract_type,
            memo=verified_addr.memo,
            customer_id=verified_addr.customer_id,
            address_type=verified_addr.address_type,
            tn_participant_id=verified_addr.tn_participant_id,
            exchange_account_id=verified_addr.exchange_account_id,
            linked_internal_addresses=verified_addr.linked_internal_addresses,
            linked_wallets=verified_addr.linked_wallets,
            # Non-security: from payload, fallback to DTO
            network=verified_addr.network
            or (getattr(dto, "network", None) if dto else None),
            status=verified_addr.status
            or (getattr(dto, "status", None) if dto else None),
            created_at=created_at,
            attributes=attributes_dict,
        )
        envelope.verified_rules_container = result.rules_container

    def _build_rules_container_cache(self, reply: Any) -> Dict[str, DecodedRulesContainer]:
        """
        Build a cache of verified rules containers from the normalized response.

        When rulesContainerNormalized=True, the API returns deduplicated rules containers
        in reply.rules_containers. Each container is verified once and cached by hash,
        so per-address verification can skip steps 2-3 for cache hits.

        Args:
            reply: The API response containing rules_containers.

        Returns:
            Dict mapping rules container hash to decoded rules container.
        """
        cache: Dict[str, DecodedRulesContainer] = {}
        rules_containers = getattr(reply, "rules_containers", None)
        if not rules_containers:
            return cache

        # Deduplicate by base64 container string to avoid re-verifying identical containers
        verified_containers: Dict[str, DecodedRulesContainer] = {}

        for hash_container in rules_containers:
            container_hash = getattr(hash_container, "hash", None)
            container_base64 = getattr(hash_container, "rules_container", None)
            signatures_base64 = getattr(hash_container, "rules_signatures", None)

            if not container_hash or not container_base64:
                continue

            # The hash is the LABEL a row uses to pick its container, and it arrives in
            # the same response as the container. Recompute it: otherwise a server can
            # file container A under container B's label and steer any row to any other
            # validly-signed container -- an older ruleset with a weaker group threshold,
            # say. Both pass the SuperAdmin check, so signature verification alone does
            # not catch it.
            computed_label = _container_hash_label(container_base64)
            if computed_label != container_hash:
                raise ContainerIntegrityError(
                    f"rules container hash mismatch: response labelled it "
                    f"{container_hash} but its bytes hash to {computed_label}"
                )

            # Check if we already verified this container (dedup by content)
            decoded = verified_containers.get(container_base64)
            if decoded is None:
                decoded = self._verify_and_decode_rules_container(
                    container_base64, signatures_base64
                )
                verified_containers[container_base64] = decoded

            cache[container_hash] = decoded

        return cache

    def _verify_and_decode_rules_container(
        self, rules_container_base64: str, rules_signatures_base64: Optional[str]
    ) -> DecodedRulesContainer:
        """
        Verify SuperAdmin signatures on a rules container and decode it.

        This performs steps 2-3 of the verification flow for a single rules container.

        Args:
            rules_container_base64: Base64-encoded rules container.
            rules_signatures_base64: Base64-encoded rules signatures.

        Returns:
            The decoded rules container.

        Raises:
            IntegrityError: If signature verification or decoding fails.
        """
        if not rules_signatures_base64:
            raise IntegrityError("rules signatures is empty for normalized container")
        if not rules_container_base64:
            raise IntegrityError("rules container is empty for normalized container")

        # Decode signatures
        try:
            signatures = user_signatures_from_base64(rules_signatures_base64)
        except (ValueError, binascii.Error, KeyError) as e:
            raise IntegrityError(f"failed to decode rules signatures: {e}") from e

        # Decode rules container data
        try:
            rules_data = base64.b64decode(rules_container_base64)
        except (binascii.Error, ValueError) as e:
            raise IntegrityError(f"failed to decode rules container: {e}") from e

        try:
            verify_governance_rules_signatures(
                rules_data,
                signatures,
                self._verifier._super_admin_keys,
                self._verifier._min_valid_signatures,
            )
        except IntegrityError as e:
            raise IntegrityError(f"rules container signature verification failed: {e}") from e

        # Decode rules container
        try:
            return rules_container_from_base64(rules_container_base64)
        except (ValueError, KeyError, binascii.Error) as e:
            raise IntegrityError(f"failed to decode rules container: {e}") from e

    @staticmethod
    def _map_envelope_from_dto(dto: Any) -> SignedWhitelistedAddressEnvelope:
        """Map OpenAPI DTO to envelope model."""
        metadata = None
        if hasattr(dto, "metadata") and dto.metadata:
            metadata = WhitelistMetadata(
                hash=getattr(dto.metadata, "hash", None),
                payload_as_string=getattr(dto.metadata, "payload_as_string", None)
                or getattr(dto.metadata, "payloadAsString", None),
            )

        # Map signature entries for WhitelistSignatureEntry (used by verifier)
        sig_entries: List[WhitelistSignatureEntry] = []
        # Also keep flat WhitelistSignature list for backward compatibility
        flat_signatures: List[WhitelistSignature] = []

        signed_address_dto = getattr(dto, "signed_address", None) or getattr(
            dto, "signedAddress", None
        )
        if signed_address_dto and hasattr(signed_address_dto, "signatures"):
            for sig_dto in signed_address_dto.signatures or []:
                # sig_dto has nested structure:
                # - signature: TgvalidatordWhitelistUserSignature (user_id, signature, comment)
                # - hashes: List[str]
                user_sig_dto = getattr(sig_dto, "signature", None)
                hashes = getattr(sig_dto, "hashes", None) or []

                user_sig = None
                if user_sig_dto:
                    user_sig = WhitelistUserSignature(
                        user_id=getattr(user_sig_dto, "user_id", None),
                        signature=getattr(user_sig_dto, "signature", None),
                        comment=getattr(user_sig_dto, "comment", None),
                    )

                sig_entries.append(
                    WhitelistSignatureEntry(
                        user_signature=user_sig,
                        hashes=list(hashes),
                    )
                )

                flat_signatures.append(
                    WhitelistSignature(
                        user_id=getattr(user_sig_dto, "user_id", None) if user_sig_dto else None,
                        signature=getattr(user_sig_dto, "signature", None)
                        if user_sig_dto
                        else None,
                        hash=hashes[0] if hashes else None,
                        hashes=list(hashes),
                    )
                )

        signed_address = SignedWhitelistedAddress(signatures=sig_entries) if sig_entries else None

        # Parse linked wallets from payload (if available)
        linked_wallets: List[InternalWallet] = []
        if metadata and metadata.payload_as_string:
            try:
                payload = json.loads(metadata.payload_as_string)
                raw_wallets = payload.get("linkedWallets", [])
                for w in raw_wallets:
                    if isinstance(w, dict):
                        linked_wallets.append(
                            InternalWallet(
                                id=int(w.get("id", 0)),
                                path=w.get("path"),
                            )
                        )
            except (json.JSONDecodeError, ValueError):
                pass

        return SignedWhitelistedAddressEnvelope(
            metadata=metadata,
            blockchain=getattr(dto, "blockchain", None),
            network=getattr(dto, "network", None),
            rules_container=getattr(dto, "rules_container", None)
            or getattr(dto, "rulesContainer", None),
            rules_signatures=getattr(dto, "rules_signatures", None)
            or getattr(dto, "rulesSignatures", None),
            rules_container_hash=getattr(dto, "rules_container_hash", None)
            or getattr(dto, "rulesContainerHash", None),
            signatures=flat_signatures,
            signed_address=signed_address,
            linked_wallets=linked_wallets,
        )
