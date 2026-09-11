"""Security tests for WhitelistedAssetService.

These tests verify that security-critical fields are sourced ONLY from the
cryptographically verified payload, not from unverified DTO attributes.
"""

from __future__ import annotations

import json
from typing import Any, Dict, Optional
from unittest.mock import MagicMock

import pytest

from taurus_protect.helpers.whitelisted_asset_verifier import AssetVerificationResult
from taurus_protect.services.whitelisted_asset_service import (
    WhitelistedAssetService,
)


class MockDTO:
    """Mock DTO object for testing."""

    def __init__(self, **kwargs: Any) -> None:
        for key, value in kwargs.items():
            setattr(self, key, value)


def create_mock_dto_with_payload(
    payload: Dict[str, Any],
    dto_overrides: Optional[Dict[str, Any]] = None,
) -> MockDTO:
    """Create a mock DTO with metadata containing the given payload."""
    metadata = MockDTO(
        hash="abc123",
        payload=payload,
        payload_as_string=json.dumps(payload, separators=(",", ":")),
        payloadAsString=json.dumps(payload, separators=(",", ":")),
    )

    dto_attrs: Dict[str, Any] = {
        "id": "asset-123",
        "tenant_id": "tenant-1",
        "tenantId": "tenant-1",
        "metadata": metadata,
        # _verified_asset refuses an envelope missing any of these three before it
        # verifies anything, so they have to be present for a field-sourcing test to
        # reach step 6 at all. The values are opaque: the verifier is stubbed.
        "signed_contract_address": MockDTO(payload="{}", signatures=[]),
        "signedContractAddress": None,
        "rules_container": "cnt",
        "rulesContainer": None,
        "rules_signatures": "sig",
        "rulesSignatures": None,
        "status": "APPROVED",
        "action": None,
        "rule": None,
        "created_at": None,
        "createdAt": None,
        "business_rule_enabled": False,
        "businessRuleEnabled": False,
        # Default DTO values (should be ignored for security fields)
        "name": None,
        "symbol": None,
        "blockchain": None,
        "network": None,
        "contract_address": None,
        "contractAddress": None,
    }

    if dto_overrides:
        dto_attrs.update(dto_overrides)

    return MockDTO(**dto_attrs)


def _verifier_reporting_delivered_payload() -> MagicMock:
    """A stubbed verifier whose step-4 answer is "the delivered payload matched".

    Step 6 now parses ``AssetVerificationResult.verified_payload`` rather than
    ``metadata.payload_as_string``, so a bare ``MagicMock()`` verifier hands the parse a
    mock object. These tests are about field SOURCING, not about which variant matched,
    so the stub reports the delivered payload -- the case where the two DIFFER is what
    ``TestLegacyVariantIsWhatStep6Parses`` covers, with real signatures.
    """
    verifier = MagicMock()

    def _verify(asset: Any, *_args: Any, **_kwargs: Any) -> Any:
        return AssetVerificationResult(
            rules_container=MagicMock(),
            verified_hash=asset.metadata.hash if asset.metadata else "",
            verified_payload=asset.metadata.payload_as_string if asset.metadata else "",
        )

    verifier.verify_whitelisted_asset.side_effect = _verify
    return verifier


def _create_service_with_mock_verifier() -> WhitelistedAssetService:
    """Create a WhitelistedAssetService with a stubbed verifier for unit tests.

    The verifier is always required, but for field sourcing tests we stub it
    to avoid needing real SuperAdmin keys.
    """
    api_client = MagicMock()
    assets_api = MagicMock()
    mock_keys = [MagicMock()]  # Mock key list (non-empty)
    service = WhitelistedAssetService(
        api_client=api_client,
        assets_api=assets_api,
        super_admin_keys=mock_keys,
        min_valid_signatures=1,
    )
    service._verifier = _verifier_reporting_delivered_payload()
    return service


def _verified(service: WhitelistedAssetService, dto: Any) -> Any:
    """Map a DTO and run the verification seam, i.e. what every read path does.

    Field sourcing is asserted THROUGH this rather than off ``_map_asset_from_dto``:
    the mapper no longer parses the payload at all, because at map time nothing has
    been verified. Step 6 is in ``_verified_asset``.
    """
    return service._verified_asset(WhitelistedAssetService._map_asset_from_dto(dto), dto=dto)


class TestWhitelistedAssetServiceSecurity:
    """Security tests for field sourcing in WhitelistedAssetService."""

    @pytest.fixture
    def service(self) -> WhitelistedAssetService:
        """Create a WhitelistedAssetService for testing with mocked verifier."""
        return _create_service_with_mock_verifier()

    def test_security_fields_from_payload_not_dto(
        self, service: WhitelistedAssetService
    ) -> None:
        """Test that name, symbol, contract_address come from payload, not DTO."""
        # Payload has verified values
        payload = {
            "name": "Verified Token Name",
            "symbol": "VTN",
            "contract_address": "0xverified_contract",
            "blockchain": "ETH",
            "network": "mainnet",
        }

        # DTO has DIFFERENT values (attacker-controlled, should be ignored)
        dto = create_mock_dto_with_payload(
            payload,
            dto_overrides={
                "name": "Malicious Name",
                "symbol": "FAKE",
                "contract_address": "0xmalicious_contract",
                "contractAddress": "0xmalicious_contract",
                "blockchain": "ATTACKER_CHAIN",
                "network": "attacker_net",
            },
        )

        asset = _verified(service, dto)

        # Security-critical fields MUST come from payload
        assert asset.name == "Verified Token Name"
        assert asset.symbol == "VTN"
        assert asset.contract_address == "0xverified_contract"
        assert asset.blockchain == "ETH"
        assert asset.network == "mainnet"

    def test_payload_missing_field_returns_none_not_dto_value(
        self, service: WhitelistedAssetService
    ) -> None:
        """Test that missing payload field returns None, not DTO value."""
        # Payload is missing 'symbol' (simulating incomplete payload)
        payload = {
            "name": "Token Name",
            # symbol is intentionally missing
            "contract_address": "0xcontract",
        }

        # DTO has symbol value (should NOT be used as fallback)
        dto = create_mock_dto_with_payload(
            payload,
            dto_overrides={
                "symbol": "SHOULD_NOT_USE",  # This MUST be ignored
            },
        )

        asset = _verified(service, dto)

        # Symbol should be None, not the DTO value
        assert asset.symbol is None
        # But name and contract_address should be from payload
        assert asset.name == "Token Name"
        assert asset.contract_address == "0xcontract"

    def test_payload_empty_returns_none_not_dto_value(
        self, service: WhitelistedAssetService
    ) -> None:
        """Test that empty payload results in None fields, not DTO values."""
        # Empty payload
        payload: Dict[str, Any] = {}

        # DTO has all values (should NOT be used)
        dto = create_mock_dto_with_payload(
            payload,
            dto_overrides={
                "name": "DTO Name",
                "symbol": "DTO",
                "contract_address": "0xdto_contract",
                "contractAddress": "0xdto_contract",
                "blockchain": "DTO_CHAIN",
                "network": "dto_net",
            },
        )

        asset = _verified(service, dto)

        # All security fields should be None (not from DTO)
        assert asset.name is None
        assert asset.symbol is None
        assert asset.contract_address is None
        assert asset.blockchain is None
        assert asset.network is None

    def test_contract_address_snake_case_in_payload(
        self, service: WhitelistedAssetService
    ) -> None:
        """Test contract_address is extracted with snake_case key."""
        payload = {
            "contract_address": "0xsnake_case_address",
        }

        dto = create_mock_dto_with_payload(payload)
        asset = _verified(service, dto)

        assert asset.contract_address == "0xsnake_case_address"

    def test_contract_address_camel_case_in_payload(
        self, service: WhitelistedAssetService
    ) -> None:
        """Test contract_address is extracted with camelCase key."""
        payload = {
            "contractAddress": "0xcamel_case_address",
        }

        dto = create_mock_dto_with_payload(payload)
        asset = _verified(service, dto)

        assert asset.contract_address == "0xcamel_case_address"

    def test_blockchain_case_variations_in_payload(
        self, service: WhitelistedAssetService
    ) -> None:
        """Test blockchain is extracted with different case variations."""
        # Test lowercase
        payload1 = {"blockchain": "eth"}
        dto1 = create_mock_dto_with_payload(payload1)
        asset1 = _verified(service, dto1)
        assert asset1.blockchain == "eth"

        # Test capitalized (Blockchain with capital B)
        payload2 = {"Blockchain": "ETH"}
        dto2 = create_mock_dto_with_payload(payload2)
        asset2 = _verified(service, dto2)
        assert asset2.blockchain == "ETH"

    def test_non_security_fields_can_come_from_dto(
        self, service: WhitelistedAssetService
    ) -> None:
        """Test that non-security fields (status, action, etc.) come from DTO."""
        payload = {"name": "Token"}

        dto = create_mock_dto_with_payload(
            payload,
            dto_overrides={
                "status": "APPROVED",
                "action": "ADD",
                "rule": "rule-1",
            },
        )

        asset = _verified(service, dto)

        # Non-security fields are fine from DTO
        assert asset.status == "APPROVED"
        assert asset.action == "ADD"
        assert asset.rule == "rule-1"

    def test_list_always_verifies_each_asset(self) -> None:
        """Every row list() returns has been through the verifier.

        Asserted on the VERIFIER, not on a service-private wrapper: the wrapper is the
        thing under test, so counting calls to it would pass against a wrapper that
        verifies nothing.
        """
        service = _create_service_with_mock_verifier()

        payload = {
            "name": "Test Token",
            "symbol": "TT",
            "contract_address": "0xtest",
        }
        dto = create_mock_dto_with_payload(payload)

        mock_reply = MagicMock()
        mock_reply.result = [dto, dto]
        mock_reply.total_items = "2"

        service._api.whitelist_service_get_whitelisted_contracts.return_value = mock_reply

        assets, pagination = service.list(limit=50, offset=0)

        # Should return both assets and verify each one
        assert len(assets) == 2
        assert all(a.name == "Test Token" for a in assets)
        assert service._verifier.verify_whitelisted_asset.call_count == 2

    def test_list_for_approval_always_verifies_each_asset(self) -> None:
        """The for-approval rows an approver inspects must be verified too."""
        service = _create_service_with_mock_verifier()

        dto = create_mock_dto_with_payload(
            {"name": "Test Token", "symbol": "TT", "contract_address": "0xtest"}
        )

        mock_reply = MagicMock()
        mock_reply.result = [dto, dto]
        mock_reply.total_items = "2"

        api = service._api
        api.whitelist_service_get_whitelisted_contracts_for_approval.return_value = (
            mock_reply
        )

        assets, pagination = service.list_for_approval(limit=50, offset=0)

        assert len(assets) == 2
        assert service._verifier.verify_whitelisted_asset.call_count == 2
        assert pagination is not None

    def test_list_for_approval_passes_filters(self) -> None:
        service = _create_service_with_mock_verifier()

        mock_reply = MagicMock()
        mock_reply.result = []
        mock_reply.total_items = "0"

        api = service._api
        api.whitelist_service_get_whitelisted_contracts_for_approval.return_value = (
            mock_reply
        )

        service.list_for_approval(ids=["3", "4"], limit=25, offset=50)

        api.whitelist_service_get_whitelisted_contracts_for_approval.assert_called_once_with(
            ids=["3", "4"], limit="25", offset="50"
        )

    def test_list_for_approval_rejects_bad_paging(self) -> None:
        service = _create_service_with_mock_verifier()

        with pytest.raises(ValueError):
            service.list_for_approval(limit=0)
        with pytest.raises(ValueError):
            service.list_for_approval(offset=-1)


def _reviewed_assets(**id_to_hash: str) -> "WhitelistedAssetApproval":
    """Mints the content pin the way a caller does -- from the ASSET OBJECTS a verified
    read returned, via ``WhitelistedAssetApproval.select``.

    ``_pinned`` is private, so a hand-built value pins nothing and the approval refuses
    it; ``test_rejects_bad_input`` exercises that directly.
    """
    from taurus_protect.models.whitelisted_address import (
        WhitelistedAsset,
        WhitelistedAssetApproval,
        WhitelistedAssetMetadata,
    )

    assets = [
        WhitelistedAsset(
            id=aid, metadata=WhitelistedAssetMetadata(hash=h, payload_as_string="{}")
        )
        for aid, h in id_to_hash.items()
    ]
    return WhitelistedAssetApproval.select(assets, *id_to_hash.keys())


class TestApproveWhitelistedAssets:
    """One signature covers every hash in the batch, so approval is all-or-nothing.

    Signing only the rows that verified would tell the approver they approved less than
    they did, and the API takes a single signature so there is no partial submission.
    """

    @staticmethod
    def _key() -> Any:
        from cryptography.hazmat.primitives.asymmetric import ec

        return ec.generate_private_key(ec.SECP256R1())

    def test_signs_nothing_when_one_row_fails_verification(self) -> None:
        from unittest.mock import patch

        from taurus_protect.errors import IntegrityError

        service = _create_service_with_mock_verifier()
        good = create_mock_dto_with_payload(
            {"name": "Test Token", "symbol": "TT", "contract_address": "0xtest"}
        )

        # Row 1 verifies, row 2 does not: the mixed batch is the case that separates
        # all-or-nothing from sign-the-survivors.
        def verify(asset: Any, dto: Any = None) -> Any:
            if asset.id == "asset-456":
                raise IntegrityError("hash mismatch")
            return asset

        api = service._api
        api.whitelist_service_get_whitelisted_contract.side_effect = [
            MagicMock(result=good),
            MagicMock(
                result=create_mock_dto_with_payload(
                    {"name": "Bad", "symbol": "B", "contract_address": "0xEVIL"},
                    {"id": "asset-456"},
                )
            ),
        ]

        with patch.object(service, "_verified_asset", side_effect=verify):
            with pytest.raises(IntegrityError, match="refusing to sign"):
                service.approve(
                    _reviewed_assets(**{"1": "abc123", "2": "abc123"}),
                    self._key(),
                    "batch approval",
                )

        api.whitelist_service_approve_whitelisted_contract.assert_not_called()

    def test_signs_the_verified_hashes_in_sorted_order(self) -> None:
        """
        The batch is re-read through the verifying LIST path filtered by ids -- one round
        trip and one rules-container fetch for the whole batch, where the per-id GET this
        replaced cost both per id.
        """
        from unittest.mock import patch

        service = _create_service_with_mock_verifier()
        api = service._api

        def row(asset_id: str) -> Any:
            return create_mock_dto_with_payload(
                {"name": "T", "symbol": "T", "contract_address": "0x1"},
                {"id": asset_id, "metadata": MockDTO(
                    hash=f"hash-{asset_id}",
                    payload={},
                    payload_as_string="{}",
                    payloadAsString="{}",
                )},
            )

        api.whitelist_service_get_whitelisted_contracts.return_value = MagicMock(
            result=[row("3"), row("7")], total_items="2"
        )

        with patch.object(service, "_verified_asset", side_effect=lambda asset, dto=None: asset):
            service.approve(
                _reviewed_assets(**{"7": "hash-7", "3": "hash-3"}),
                self._key(),
                "batch approval",
            )

        # ONE list call, not one GET per id.
        api.whitelist_service_get_whitelisted_contracts.assert_called_once()
        called = api.whitelist_service_get_whitelisted_contracts.call_args.kwargs
        assert called["whitelisted_contract_address_ids"] == ["3", "7"]
        assert called["include_for_approval"] is True

        body = api.whitelist_service_approve_whitelisted_contract.call_args.kwargs["body"]
        assert body.ids == ["3", "7"]
        assert body.comment == "batch approval"
        assert body.signature

    def test_aborts_when_the_verified_read_omits_a_row(self) -> None:
        """
        A page that silently omits a row must not become an approval of fewer rows than
        the caller asked for.
        """
        from unittest.mock import patch

        from taurus_protect.errors import IntegrityError

        service = _create_service_with_mock_verifier()
        api = service._api
        api.whitelist_service_get_whitelisted_contracts.return_value = MagicMock(
            result=[], total_items="0"
        )

        with patch.object(service, "_verified_asset", side_effect=lambda asset, dto=None: asset):
            with pytest.raises(IntegrityError, match="was not returned by the verified read"):
                service.approve(
                _reviewed_assets(**{"7": "hash-7", "3": "hash-3"}),
                self._key(),
                "batch approval",
            )

        api.whitelist_service_approve_whitelisted_contract.assert_not_called()

    def test_refuses_to_sign_a_row_that_changed_since_it_was_reviewed(self) -> None:
        """The content pin. Verification alone does not catch this: the substituted row is
        a genuine, validly-signed whitelist entry -- just not the one the approver read.

        Without the pin a response-controlling server answers the id-filtered re-read with
        a DIFFERENT row whose existing signatures already satisfy the container it
        presents, and harvests a real approver signature over content never reviewed.
        """
        from unittest.mock import patch

        from taurus_protect.errors import IntegrityError

        service = _create_service_with_mock_verifier()
        api = service._api

        substituted = create_mock_dto_with_payload(
            {"name": "T", "symbol": "T", "contract_address": "0xEVIL"},
            {"id": "3", "metadata": MockDTO(
                hash="hash-substituted",
                payload={},
                payload_as_string="{}",
                payloadAsString="{}",
            )},
        )
        api.whitelist_service_get_whitelisted_contracts.return_value = MagicMock(
            result=[substituted], total_items="1"
        )

        with patch.object(
            service, "_verified_asset", side_effect=lambda asset, dto=None: asset
        ):
            with pytest.raises(IntegrityError, match="changed since it was reviewed"):
                service.approve(
                    _reviewed_assets(**{"3": "hash-3"}), self._key(), "batch approval"
                )

        api.whitelist_service_approve_whitelisted_contract.assert_not_called()

    def test_rejects_bad_input(self) -> None:
        from taurus_protect.models.whitelisted_address import WhitelistedAssetApproval

        service = _create_service_with_mock_verifier()
        key = self._key()

        with pytest.raises(ValueError, match="selection cannot be empty"):
            # A hand-built value pins nothing, and must be refused rather than treated as
            # "approve nothing" -- an empty pin must not silently restore the unpinned
            # behaviour, the same rule approve_rules_proposal's container pin follows.
            service.approve(WhitelistedAssetApproval({}), key, "c")
        with pytest.raises(ValueError, match="cannot select an empty set of ids"):
            # And `select` refuses to mint one in the first place.
            _reviewed_assets()
        with pytest.raises(ValueError, match="private_key is required"):
            service.approve(_reviewed_assets(**{"1": "abc123"}), None, "c")
        with pytest.raises(ValueError, match="comment is required"):
            service.approve(_reviewed_assets(**{"1": "abc123"}), key, "")
        with pytest.raises(ValueError, match="positive integer"):
            service.approve(_reviewed_assets(**{"0": "abc123"}), key, "c")

        service._api.whitelist_service_approve_whitelisted_contract.assert_not_called()


class TestMapAssetFromDto:
    """Direct tests for the _map_asset_from_dto static method.

    The mapper is deliberately identity-BLIND now: at map time nothing has been
    verified, so ``name``/``symbol``/``blockchain``/``network``/``contract_address``/
    ``decimals``/``token_id`` are left unset and filled in by ``_verified_asset`` from
    the payload step 4 matched. Asserting them here would re-pin the defect.
    """

    @pytest.fixture
    def service(self) -> WhitelistedAssetService:
        return _create_service_with_mock_verifier()

    def test_the_mapper_never_sources_an_identity_field(self) -> None:
        """The mapper must not read the payload, even when one is present.

        A payload delivered alongside an envelope is unverified: the legacy-hash
        tolerance means the bytes a signature covered can differ from the bytes
        delivered, so a value read here is a value no signature covered.
        """
        dto = create_mock_dto_with_payload(
            {
                "name": "Delivered Name",
                "symbol": "DLV",
                "blockchain": "ETH",
                "network": "mainnet",
                "contractAddress": "0xdelivered",
                "decimals": 18,
                "tokenId": "7",
            }
        )

        asset = WhitelistedAssetService._map_asset_from_dto(dto)

        assert asset.name is None
        assert asset.symbol is None
        assert asset.blockchain is None
        assert asset.network is None
        assert asset.contract_address is None
        assert asset.decimals is None
        assert asset.token_id is None
        # Non-security fields and the verification material are still mapped.
        assert asset.id == "asset-123"
        assert asset.status == "APPROVED"
        assert asset.metadata is not None
        assert asset.metadata.hash == "abc123"

    def test_no_metadata_returns_none_security_fields(self) -> None:
        """Test that missing metadata results in None for security fields."""
        dto = MockDTO(
            id="asset-1",
            tenant_id="tenant-1",
            metadata=None,
            signed_contract_address=None,
            rules_container=None,
            rules_signatures=None,
            status="PENDING",
            action=None,
            rule=None,
            created_at=None,
            business_rule_enabled=False,
            # DTO has values but should NOT be used
            name="DTO Name",
            symbol="DTO",
            blockchain="DTO_CHAIN",
            network="dto_net",
            contract_address="0xdto",
        )

        asset = WhitelistedAssetService._map_asset_from_dto(dto)

        # Without metadata payload, security fields must be None
        assert asset.name is None
        assert asset.symbol is None
        assert asset.contract_address is None
        assert asset.blockchain is None
        assert asset.network is None

    def test_verified_payload_is_used_not_payload_dict(
        self, service: WhitelistedAssetService
    ) -> None:
        """The identity comes from the VERIFIED payload, never the raw payload dict.

        SECURITY: even when ``metadata.payload`` carries data, extraction reads the
        payload verification cleared -- the cryptographically committed source.
        """
        metadata = MockDTO(
            hash="abc",
            payload={"name": "attacker"},  # unverified object, must be ignored
            payload_as_string='{"name":"test"}',
            payloadAsString=None,
        )

        dto = MockDTO(
            id="asset-1",
            metadata=metadata,
            signed_contract_address=MockDTO(payload="{}", signatures=[]),
            signedContractAddress=None,
            rules_container="cnt",
            rulesContainer=None,
            rules_signatures="sig",
            rulesSignatures=None,
            status="PENDING",
            action=None,
            rule=None,
            created_at=None,
            tenant_id=None,
            business_rule_enabled=False,
        )

        asset = _verified(service, dto)

        assert asset.name == "test"

    def test_signed_contract_address_mapping(self) -> None:
        """Test that signed_contract_address is correctly mapped."""
        payload = {"name": "Token"}

        user_sig = MockDTO(
            user_id="user-1",
            userId="user-1",
            signature="sig123",
            comment="LGTM",
        )
        sig_entry = MockDTO(
            user_signature=user_sig,
            userSignature=user_sig,
            hashes=["hash1", "hash2"],
        )
        signed = MockDTO(
            payload="signed_payload",
            signatures=[sig_entry],
        )

        dto = create_mock_dto_with_payload(payload)
        dto.signed_contract_address = signed
        dto.signedContractAddress = signed

        asset = WhitelistedAssetService._map_asset_from_dto(dto)

        assert asset.signed_contract_address is not None
        assert len(asset.signed_contract_address.signatures) == 1
        assert asset.signed_contract_address.signatures[0].user_signature.user_id == "user-1"
        assert asset.signed_contract_address.signatures[0].hashes == ["hash1", "hash2"]
