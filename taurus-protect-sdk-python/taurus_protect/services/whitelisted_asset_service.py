"""Whitelisted asset service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from cryptography.hazmat.primitives.asymmetric.ec import EllipticCurvePublicKey

from taurus_protect._internal.openapi.exceptions import ApiException
from taurus_protect.models.pagination import Pagination
from taurus_protect.models.whitelisted_address import (
    SignedContractAddress,
    WhitelistedAsset,
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
            self._verify_asset(asset, dto=result)

            return asset
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

    def list(
        self,
        blockchain: Optional[str] = None,
        network: Optional[str] = None,
        limit: int = 50,
        offset: int = 0,
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
            )

            assets: List[WhitelistedAsset] = []
            if reply.result:
                for dto in reply.result:
                    asset = self._map_asset_from_dto(dto)
                    self._verify_asset(asset, dto=dto)
                    assets.append(asset)

            pagination = self._extract_pagination(
                getattr(reply, "total_items", None),
                offset,
                limit,
            )
            return assets, pagination
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

    def _verify_asset(
        self,
        asset: WhitelistedAsset,
        dto: Optional[Any] = None,
    ) -> None:
        """
        Perform the 5-step integrity verification on a whitelisted asset.

        Args:
            asset: The asset to verify.
            dto: The original DTO envelope (provides blockchain/network for rules lookup
                when the verified payload omits these fields).

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
            from taurus_protect.errors import IntegrityError

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

        self._verifier.verify_whitelisted_asset(
            asset,
            rules_container_from_base64,
            user_signatures_from_base64,
            dto_blockchain=dto_blockchain,
            dto_network=dto_network,
        )

    @staticmethod
    def _map_asset_from_dto(dto: Any) -> WhitelistedAsset:
        """Map OpenAPI DTO to asset model."""
        import json

        # Map metadata if present
        metadata = None
        dto_metadata = getattr(dto, "metadata", None)
        payload: dict = {}
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
            # SECURITY: Extract payload dict from verified payload_as_string ONLY
            # (not from dto_metadata.payload which is unverified)
            if payload_as_string:
                try:
                    parsed = json.loads(payload_as_string)
                    if isinstance(parsed, dict):
                        payload = parsed
                except json.JSONDecodeError:
                    pass  # payload remains empty dict

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
        # by manipulating DTO fields that differ from the signed payload.
        return WhitelistedAsset(
            id=str(getattr(dto, "id", "")),
            tenant_id=getattr(dto, "tenant_id", None) or getattr(dto, "tenantId", None),
            # Security-critical fields from verified payload only (no DTO fallback)
            name=payload.get("name"),
            symbol=payload.get("symbol"),
            blockchain=payload.get("blockchain") or payload.get("Blockchain"),
            network=payload.get("network") or payload.get("Network"),
            contract_address=payload.get("contract_address") or payload.get("contractAddress"),
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
