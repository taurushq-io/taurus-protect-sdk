"""Whitelisted address service for Taurus-PROTECT SDK."""

from __future__ import annotations

import base64
import binascii
import json
from typing import TYPE_CHECKING, Any, Dict, List, Optional, Tuple

from cryptography.hazmat.primitives.asymmetric.ec import EllipticCurvePublicKey

from taurus_protect._internal.openapi.exceptions import ApiException
from taurus_protect.errors import APIError, IntegrityError
from taurus_protect.helpers.signature_verifier import is_valid_signature
from taurus_protect.helpers.whitelisted_address_verifier import WhitelistedAddressVerifier
from taurus_protect.mappers.governance_rules import (
    rules_container_from_base64,
    user_signatures_from_base64,
)
from taurus_protect.models.governance_rules import DecodedRulesContainer
from taurus_protect.models.pagination import Pagination
from taurus_protect.models.whitelisted_address import (
    InternalWallet,
    SignedWhitelistedAddress,
    SignedWhitelistedAddressEnvelope,
    WhitelistedAddress,
    WhitelistMetadata,
    WhitelistSignature,
    WhitelistSignatureEntry,
    WhitelistUserSignature,
)
from taurus_protect.services._base import BaseService

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
    ) -> Tuple[List[WhitelistedAddress], Optional[Pagination]]:
        """
        List whitelisted addresses with cryptographic verification.

        Each address is verified using the 6-step verification flow.

        Args:
            currency: Filter by currency.
            limit: Maximum number of addresses to return.
            offset: Offset for pagination.

        Returns:
            Tuple of (verified addresses list, pagination info).

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
            )

            # Build rules container cache from normalized containers
            rules_container_cache = self._build_rules_container_cache(reply)

            addresses: List[WhitelistedAddress] = []
            if reply.result:
                for dto in reply.result:
                    # Map to envelope and verify (strict mode - fail on first error)
                    envelope = self._map_envelope_from_dto(dto)

                    # Look up cached rules container by hash
                    cached = None
                    if envelope.rules_container_hash:
                        cached = rules_container_cache.get(envelope.rules_container_hash)

                    self._verify_and_populate_envelope(envelope, dto, cached_rules_container=cached)
                    if envelope.verified_whitelisted_address:
                        addresses.append(envelope.verified_whitelisted_address)

            pagination = self._extract_pagination(
                getattr(reply, "total_items", None),
                offset,
                limit,
            )
            return addresses, pagination
        except IntegrityError:
            raise
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

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

        # Verify signatures
        valid_count = 0
        for sig in signatures:
            if sig.signature and is_valid_signature(
                rules_data, sig.signature, self._verifier._super_admin_keys
            ):
                valid_count += 1

        if valid_count < self._verifier._min_valid_signatures:
            raise IntegrityError(
                f"rules container signature verification failed: only {valid_count} valid "
                f"signatures, minimum {self._verifier._min_valid_signatures} required"
            )

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
