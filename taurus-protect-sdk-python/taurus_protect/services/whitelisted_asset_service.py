"""Whitelisted asset service for Taurus-PROTECT SDK."""

from __future__ import annotations

import hmac
import json
from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from cryptography.hazmat.primitives.asymmetric.ec import EllipticCurvePublicKey

from taurus_protect.crypto.signing import sign_data
from taurus_protect.errors import APIError, IntegrityError, WhitelistError
from taurus_protect.helpers.whitelist_hash_helper import (
    parse_whitelisted_asset_identity_from_json,
)
from taurus_protect.models.pagination import Pagination
from taurus_protect.models.whitelisted_address import (
    SignedContractAddress,
    WhitelistedAsset,
    WhitelistedAssetApproval,
    WhitelistedAssetMetadata,
    WhitelistSignatureEntry,
    WhitelistUserSignature,
)
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    from taurus_protect.helpers.whitelisted_asset_verifier import (
        WhitelistedAssetVerifier,
    )


class WhitelistedAssetService(BaseService):
    """
    Service for managing whitelisted assets.

    Whitelisted assets are pre-approved tokens that can be transferred.
    This is typically used for ERC-20 and similar token standards.

    All retrieved assets are cryptographically verified using the 5-step
    verification flow:
    1. Verify metadata hash (SHA-256)
    2. Verify rules container signatures (SuperAdmin keys)
    3. Decode rules container
    4. Verify hash coverage
    5. Verify whitelist signatures meet governance thresholds

    Example:
        >>> client = ProtectClient.create(
        ...     host="...", api_key="...", api_secret="...",
        ...     super_admin_keys_pem=["-----BEGIN PUBLIC KEY-----..."],
        ...     min_valid_signatures=2,
        ... )
        >>> asset = client.whitelisted_assets.get(123)  # Verified automatically
        >>> assets, _ = client.whitelisted_assets.list(blockchain="ETH")
    """

    def __init__(
        self,
        api_client: Any,
        assets_api: Any,
        super_admin_keys: List[EllipticCurvePublicKey],
        min_valid_signatures: int,
    ) -> None:
        """
        Initialize the whitelisted asset service.

        Args:
            api_client: The OpenAPI client instance.
            assets_api: The whitelisted assets API instance.
            super_admin_keys: List of SuperAdmin public keys for verification.
            min_valid_signatures: Minimum valid signatures required.
        """
        super().__init__(api_client)
        self._api = assets_api

        from taurus_protect.helpers.whitelisted_asset_verifier import (
            WhitelistedAssetVerifier,
        )

        self._verifier: WhitelistedAssetVerifier = WhitelistedAssetVerifier(
            super_admin_keys=super_admin_keys,
            min_valid_signatures=min_valid_signatures,
        )

    def get(self, asset_id: int) -> WhitelistedAsset:
        """
        Get a whitelisted asset by ID.

        The asset integrity is cryptographically verified before being returned.

        Args:
            asset_id: The whitelisted asset ID.

        Returns:
            The whitelisted asset.

        Raises:
            NotFoundError: If the asset is not found.
            IntegrityError: If verification fails.
            WhitelistError: If signature thresholds are not met.
            APIError: If the API call fails.
        """
        if asset_id <= 0:
            raise ValueError("asset_id must be positive")

        try:
            reply = self._api.whitelist_service_get_whitelisted_contract(id=str(asset_id))
            result = reply.result
            if result is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Whitelisted asset {asset_id} not found")

            asset = self._map_asset_from_dto(result)
            return self._verified_asset(asset, dto=result)
        except Exception as e:
            # Funnel on the SDK error taxonomy, not on the raw generated ApiException.
            # Keying on ApiException left every other failure -- a urllib3 transport
            # error, say -- propagating unmapped, so a caller could not treat it as an
            # APIError at all. And the pass-through list must name IntegrityError and
            # WhitelistError explicitly: both are plain Exceptions, so _handle_error
            # would map them to ServerError(500), whose is_retryable() is True. That
            # inverts the documented "security error, DO NOT retry" contract.
            if isinstance(e, (APIError, IntegrityError, WhitelistError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list(
        self,
        blockchain: Optional[str] = None,
        network: Optional[str] = None,
        limit: int = 50,
        offset: int = 0,
        *,
        ids: Optional[List[str]] = None,
        include_for_approval: bool = False,
    ) -> Tuple[List[WhitelistedAsset], Optional[Pagination]]:
        """
        List whitelisted assets.

        Each asset's integrity is cryptographically verified.

        Args:
            blockchain: Filter by blockchain.
            network: Filter by network.
            limit: Maximum number of assets to return.
            offset: Offset for pagination.

        Returns:
            Tuple of (assets list, pagination info).

        Raises:
            IntegrityError: If verification fails for any asset.
            WhitelistError: If signature thresholds are not met.
            APIError: If the API call fails.
        """
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            reply = self._api.whitelist_service_get_whitelisted_contracts(
                blockchain=blockchain,
                network=network,
                limit=str(limit),
                offset=str(offset),
                whitelisted_contract_address_ids=ids or None,
                include_for_approval=include_for_approval or None,
            )

            assets: List[WhitelistedAsset] = []
            if reply.result:
                for dto in reply.result:
                    asset = self._map_asset_from_dto(dto)
                    assets.append(self._verified_asset(asset, dto=dto))

            pagination = self._extract_pagination(
                getattr(reply, "total_items", None),
                offset,
                limit,
            )
            return assets, pagination
        except Exception as e:
            # See above: funnel on the SDK taxonomy, and never let IntegrityError or
            # WhitelistError be remapped to a retryable ServerError(500).
            if isinstance(e, (APIError, IntegrityError, WhitelistError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list_for_approval(
        self,
        ids: Optional[List[str]] = None,
        limit: int = 50,
        offset: int = 0,
    ) -> Tuple[List[WhitelistedAsset], Optional[Pagination]]:
        """
        List whitelisted assets awaiting approval, verified as in list().

        Without this the only reader of the for-approval endpoint was the unverified
        contract-whitelisting service, so the rows an approver inspects before
        whitelisting a contract address were never checked against governance.

        Args:
            ids: Filter by specific whitelisted asset IDs.
            limit: Maximum number of assets to return.
            offset: Offset for pagination.

        Returns:
            Tuple of (assets list, pagination info).

        Raises:
            IntegrityError: If verification fails for any asset.
            WhitelistError: If signature thresholds are not met.
            APIError: If the API call fails.
        """
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            reply = self._api.whitelist_service_get_whitelisted_contracts_for_approval(
                ids=ids,
                limit=str(limit),
                offset=str(offset),
            )

            assets: List[WhitelistedAsset] = []
            if reply.result:
                for dto in reply.result:
                    asset = self._map_asset_from_dto(dto)
                    assets.append(self._verified_asset(asset, dto=dto))

            pagination = self._extract_pagination(
                getattr(reply, "total_items", None),
                offset,
                limit,
            )
            return assets, pagination
        except Exception as e:
            # See above: funnel on the SDK taxonomy, and never let IntegrityError or
            # WhitelistError be remapped to a retryable ServerError(500).
            if isinstance(e, (APIError, IntegrityError, WhitelistError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def approve(
        self,
        selection: WhitelistedAssetApproval,
        private_key: Any,
        comment: str,
    ) -> None:
        """
        Sign and submit an approval for the reviewed whitelisted assets, all-or-nothing.

        ``selection`` comes from a verified read --
        ``WhitelistedAssetApproval.select(assets, *ids)`` or ``.select_all(assets)`` --
        and carries the metadata hash each row had AT REVIEW TIME. Each asset is then
        re-read through the verified path and the re-read hash must equal the pin, so
        the approver's signature covers content the approver actually saw rather than
        whatever the server chooses to return under those ids. Without the pin a
        response-controlling server could substitute a row whose existing signatures
        already satisfy the container it presents and harvest a genuine approver
        signature over content the approver never reviewed. An empty selection RAISES
        rather than meaning "approve nothing", so the pin cannot be silently bypassed.

        ``ContractWhitelistingService.approve_whitelisted_contracts`` takes an opaque
        signature over hashes nothing verified.

        Any asset that is missing, fails verification, or whose hash no longer matches
        the pin aborts the whole call and nothing is signed: one signature covers every
        hash in the batch, so a partial approval would mean the caller believes they
        approved more than they did.

        Args:
            selection: The reviewed rows, pinned from a verified read.
            private_key: The approver's P-256 private key.
            comment: The approval comment.

        Raises:
            ValueError: If any argument is missing or malformed.
            IntegrityError: If any asset is missing, unverifiable, has no hash, or no
                longer matches the reviewed pin.
            APIError: If the API call fails.
        """
        if selection is None or selection.is_empty():
            raise ValueError(
                "selection cannot be empty: pin the rows with "
                "WhitelistedAssetApproval.select(assets, ids...) or "
                ".select_all(assets) from a verified read"
            )
        if private_key is None:
            raise ValueError("private_key is required")
        if not comment:
            raise ValueError("comment is required")

        ids: List[int] = []
        for raw_id in selection.ids():
            try:
                parsed = int(raw_id)
            except (TypeError, ValueError) as exc:
                raise ValueError(
                    f"whitelisted asset ID {raw_id!r} is not a valid numeric ID"
                ) from exc
            if parsed <= 0:
                raise ValueError(
                    f"whitelisted asset ID {raw_id!r} must be a positive integer"
                )
            ids.append(parsed)

        # Sorted numerically, as the request-approval path does, so the signed order is
        # independent of the order the caller passed.
        sorted_ids = sorted(ids)

        # ONE id-filtered page through the verifying list path, not one GET per id. The
        # list path verifies every row and fetches the rules container once per call, so
        # a 50-id approval costs one round trip and one container fetch instead of fifty
        # of each. include_for_approval is required: the rows being approved are pending,
        # so the default list does not return them.
        id_strings = [str(i) for i in sorted_ids]
        try:
            verified, _ = self.list(
                limit=len(id_strings),
                ids=id_strings,
                include_for_approval=True,
            )
        except (APIError, IntegrityError, WhitelistError):
            # The SDK taxonomy propagates unchanged, so a transport failure is not
            # reported as an integrity failure (and stays retryable). Go wraps with %w
            # and keeps the type reachable via errors.As; `raise ... from` does not.
            raise
        except Exception as e:
            raise IntegrityError(f"refusing to sign: the verified read failed: {e}") from e

        by_id = {str(asset.id): asset for asset in verified if asset is not None}

        hashes: List[str] = []
        for asset_id in id_strings:
            asset = by_id.get(asset_id)
            if asset is None:
                # A page that silently omits a row must not become an approval of fewer
                # rows than the caller asked for.
                raise IntegrityError(
                    f"refusing to sign: asset {asset_id} was not returned by the "
                    "verified read"
                )
            if asset.metadata is None or not asset.metadata.hash:
                raise IntegrityError(
                    f"refusing to sign: asset {asset_id} has no metadata hash"
                )

            # The pin. Constant-time because this compares hash material, matching
            # approve_rules_proposal's use of hmac.compare_digest on its container pin.
            pinned_hash = selection.pinned_hash(asset_id)
            if not pinned_hash:
                raise IntegrityError(
                    f"refusing to sign: asset {asset_id} is not in the reviewed selection"
                )
            if not hmac.compare_digest(pinned_hash, asset.metadata.hash):
                raise IntegrityError(
                    f"refusing to sign: whitelisted asset {asset_id} changed since it "
                    f"was reviewed: reviewed hash {pinned_hash}, re-read hash "
                    f"{asset.metadata.hash}. Re-read, re-review and re-approve"
                )

            hashes.append(asset.metadata.hash)

        to_sign = json.dumps(hashes, separators=(",", ":"))
        signature = sign_data(private_key, to_sign.encode("utf-8"))

        try:
            from taurus_protect._internal.openapi.models.tgvalidatord_approve_whitelisted_contract_address_request import (
                TgvalidatordApproveWhitelistedContractAddressRequest,
            )

            body = TgvalidatordApproveWhitelistedContractAddressRequest(
                ids=[str(i) for i in sorted_ids],
                signature=signature,
                comment=comment,
            )
            self._api.whitelist_service_approve_whitelisted_contract(body=body)
        except Exception as e:
            # See above: funnel on the SDK taxonomy, and never let IntegrityError or
            # WhitelistError be remapped to a retryable ServerError(500).
            if isinstance(e, (APIError, IntegrityError, WhitelistError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def _verified_asset(
        self,
        asset: WhitelistedAsset,
        dto: Optional[Any] = None,
    ) -> WhitelistedAsset:
        """
        Run the 5-step verification, then step 6: populate the identity fields from the
        payload a counted signature COVERED.

        This is the ONE seam every asset-returning path uses (``get``, ``list``,
        ``list_for_approval``), so an identity field cannot come from anywhere else.

        ``_map_asset_from_dto`` deliberately leaves ``name``/``symbol``/``blockchain``/
        ``network``/``contract_address``/``decimals``/``token_id`` unset: nothing has
        been verified at map time, and it used to read them out of the DELIVERED
        payload. That is the injection this closes -- for a legacy-signed row the strip
        recovers the signed bytes while the delivered payload can carry an appended
        member the parse would still read.

        Args:
            asset: The mapped envelope, with identity fields still unset.
            dto: The original DTO envelope (provides blockchain/network for rules lookup
                when the verified payload omits these fields).

        Returns:
            The same asset with its identity fields sourced from the verified payload.

        Raises:
            IntegrityError: If verification fails.
            WhitelistError: If signature thresholds are not met.
        """
        # Verification data is required — missing data is an error.
        # An attacker could strip verification data to bypass checks.
        if (
            asset.metadata is None
            or not asset.rules_container
            or asset.signed_contract_address is None
        ):
            raise IntegrityError("verification enabled but required data missing")

        from taurus_protect.mappers.governance_rules import (
            rules_container_from_base64,
            user_signatures_from_base64,
        )

        # Pass DTO blockchain/network for rules lookup.
        # This matches Java SDK behavior where the envelope's DTO fields
        # (not the verified payload fields) are used to locate governance rules.
        dto_blockchain = (
            getattr(dto, "blockchain", None) if dto else None
        )
        dto_network = (
            getattr(dto, "network", None) if dto else None
        )

        result = self._verifier.verify_whitelisted_asset(
            asset,
            rules_container_from_base64,
            user_signatures_from_base64,
            dto_blockchain=dto_blockchain,
            dto_network=dto_network,
        )

        # Step 6. result.verified_payload, NOT asset.metadata.payload_as_string: when
        # step 4 matched a legacy variant, those two differ by exactly the members no
        # signature covered.
        return asset.model_copy(
            update=parse_whitelisted_asset_identity_from_json(result.verified_payload)
        )

    @staticmethod
    def _map_asset_from_dto(dto: Any) -> WhitelistedAsset:
        """
        Map an OpenAPI DTO to the asset envelope, WITHOUT any identity field.

        ``name``/``symbol``/``blockchain``/``network``/``contract_address``/
        ``decimals``/``token_id`` are left ``None`` here and filled in by
        :meth:`_verified_asset` from the payload a counted signature covered. This used
        to parse ``payload_as_string`` at map time -- before any signature had been
        checked -- and the legacy-hash tolerance makes that exploitable: for a row
        signed before ``isNFT``/``kindType`` existed, the strip recovers the signed
        bytes while the delivered payload can carry members no signature covered.

        Leaving them unset rather than parsing early also means a caller who somehow
        reaches this mapper gets nulls instead of unverified values.
        """
        # Map metadata if present
        metadata = None
        dto_metadata = getattr(dto, "metadata", None)
        if dto_metadata:
            payload_as_string = getattr(dto_metadata, "payload_as_string", None) or getattr(
                dto_metadata, "payloadAsString", None
            )
            metadata = WhitelistedAssetMetadata(
                hash=getattr(dto_metadata, "hash", None),
                # SECURITY: payload intentionally not mapped - use payload_as_string only.
                # The raw payload object could be tampered with while payloadAsString
                # remains unchanged (hash still verifies).
                payload_as_string=payload_as_string,
            )

        # Map signed contract address if present
        signed_contract_address = None
        dto_signed = getattr(dto, "signed_contract_address", None) or getattr(
            dto, "signedContractAddress", None
        )
        if dto_signed:
            signatures = []
            dto_signatures = getattr(dto_signed, "signatures", []) or []
            for sig in dto_signatures:
                user_sig = None
                # The DTO field name varies:
                # - "user_signature" / "userSignature" in some schemas
                # - "signature" as a nested object in TgvalidatordWhitelistSignature
                dto_user_sig = getattr(sig, "user_signature", None) or getattr(
                    sig, "userSignature", None
                )
                if dto_user_sig is None:
                    # Check if 'signature' is a nested user signature object (not a string)
                    nested = getattr(sig, "signature", None)
                    if nested is not None and not isinstance(nested, (str, bytes)):
                        dto_user_sig = nested
                if dto_user_sig:
                    user_sig = WhitelistUserSignature(
                        user_id=getattr(dto_user_sig, "user_id", None)
                        or getattr(dto_user_sig, "userId", None),
                        signature=getattr(dto_user_sig, "signature", None),
                        comment=getattr(dto_user_sig, "comment", None),
                    )
                signatures.append(
                    WhitelistSignatureEntry(
                        user_signature=user_sig,
                        hashes=getattr(sig, "hashes", []) or [],
                    )
                )
            signed_contract_address = SignedContractAddress(
                payload=getattr(dto_signed, "payload", None),
                signatures=signatures,
            )

        # SECURITY: All security-critical fields MUST come from verified payload only.
        # No DTO fallbacks allowed - this prevents attackers from bypassing verification
        # by manipulating DTO fields that differ from the signed payload. They are set
        # by _verified_asset, from the payload step 4 matched -- not here.
        return WhitelistedAsset(
            id=str(getattr(dto, "id", "")),
            tenant_id=getattr(dto, "tenant_id", None) or getattr(dto, "tenantId", None),
            # Non-security fields can come from DTO
            status=getattr(dto, "status", None),
            action=getattr(dto, "action", None),
            rule=getattr(dto, "rule", None),
            created_at=getattr(dto, "created_at", None) or getattr(dto, "createdAt", None),
            metadata=metadata,
            rules_container=getattr(dto, "rules_container", None)
            or getattr(dto, "rulesContainer", None),
            rules_signatures=getattr(dto, "rules_signatures", None)
            or getattr(dto, "rulesSignatures", None),
            signed_contract_address=signed_contract_address,
            business_rule_enabled=getattr(dto, "business_rule_enabled", False)
            or getattr(dto, "businessRuleEnabled", False),
        )
