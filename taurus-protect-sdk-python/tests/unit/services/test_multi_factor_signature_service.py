"""Unit tests for MultiFactorSignatureService.

This whole service used to be dead code: it called
``multi_factor_signature_service_get_challenge`` / ``_get_challenges`` /
``_create_challenge`` / ``_verify_challenge``, none of which exist on the generated
``MultiFactorSignatureApi`` (which has only approve / create-batch /
get-entities-info / reject), and imported two request models that do not exist either.
Every call raised ``AttributeError`` or ``ModuleNotFoundError`` at runtime.

**The tests passed anyway, because ``mfs_api`` was a bare ``MagicMock()``** -- which
answers any attribute with another mock, so four non-existent operations reported green.
So the API stub here is built from the real generated class with
``spec=MultiFactorSignatureApi``: a call to a method the API does not have is now an
``AttributeError`` in the test rather than a silent pass.
"""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect._internal.openapi.api.multi_factor_signature_api import (
    MultiFactorSignatureApi,
)
from taurus_protect._internal.openapi.models.tgvalidatord_multi_factor_signatures_entity_type import (  # noqa: E501
    TgvalidatordMultiFactorSignaturesEntityType,
)
from taurus_protect.models.multi_factor_signature import MultiFactorSignatureEntityType
from taurus_protect.services.multi_factor_signature_service import (
    MultiFactorSignatureService,
)


def _make_service() -> tuple:
    # spec= is the point of this fixture: it binds the mock to the real generated
    # surface, so the four-non-existent-methods defect cannot recur silently.
    mfs_api = MagicMock(spec=MultiFactorSignatureApi)
    service = MultiFactorSignatureService(api_client=MagicMock(), mfs_api=mfs_api)
    return service, mfs_api


class TestTheApiSurfaceExists:
    """The generated API really does have these four operations, and only these four."""

    @pytest.mark.parametrize(
        "operation",
        [
            "multi_factor_signature_service_approve_multi_factor_signature",
            "multi_factor_signature_service_create_multi_factor_signature_batch",
            "multi_factor_signature_service_get_multi_factor_signature_entities_info",
            "multi_factor_signature_service_reject_multi_factor_signature",
        ],
    )
    def test_operation_exists(self, operation: str) -> None:
        assert hasattr(MultiFactorSignatureApi, operation)

    @pytest.mark.parametrize(
        "phantom",
        [
            "multi_factor_signature_service_get_challenge",
            "multi_factor_signature_service_get_challenges",
            "multi_factor_signature_service_create_challenge",
            "multi_factor_signature_service_verify_challenge",
        ],
    )
    def test_the_operations_the_old_service_called_do_not_exist(self, phantom: str) -> None:
        """Pinned so the challenge-shaped API cannot be reintroduced from memory."""
        assert not hasattr(MultiFactorSignatureApi, phantom)


class TestGetMultiFactorSignatureInfo:
    def test_raises_on_empty_id(self) -> None:
        service, _ = _make_service()
        with pytest.raises(ValueError):
            service.get_multi_factor_signature_info("")

    def test_returns_the_payload_and_the_kind(self) -> None:
        service, api = _make_service()
        reply = MagicMock()
        reply.id = "mfs-1"
        reply.payload_to_sign = ["hash-a", "hash-b"]
        reply.entity_type = TgvalidatordMultiFactorSignaturesEntityType.WHITELISTED_ADDRESS
        api.multi_factor_signature_service_get_multi_factor_signature_entities_info.return_value = (
            reply
        )

        info = service.get_multi_factor_signature_info("mfs-1")

        assert info.id == "mfs-1"
        assert info.payload_to_sign == ["hash-a", "hash-b"]
        assert info.entity_type is MultiFactorSignatureEntityType.WHITELISTED_ADDRESS
        api.multi_factor_signature_service_get_multi_factor_signature_entities_info.assert_called_once_with(
            id="mfs-1"
        )

    def test_the_payload_is_not_claimed_verified(self) -> None:
        """The threat model's whole point: nothing here checks these bytes.

        The reply carries no entity id, so there is nothing to check them against. This
        test exists to keep that contract stated in code as well as in the docstring --
        if a future change adds verification, this assertion is what says so out loud.
        """
        from taurus_protect.models.multi_factor_signature import MultiFactorSignatureInfo

        field = MultiFactorSignatureInfo.model_fields["payload_to_sign"]
        assert "UNVERIFIED SERVER DATA" in (field.description or "")


class TestCreateMultiFactorSignatures:
    def test_raises_on_empty_entity_ids(self) -> None:
        service, _ = _make_service()
        with pytest.raises(ValueError, match="entity_ids cannot be empty"):
            service.create_multi_factor_signatures([], MultiFactorSignatureEntityType.REQUEST)

    def test_raises_on_missing_entity_type(self) -> None:
        service, _ = _make_service()
        with pytest.raises(ValueError, match="entity_type cannot be empty"):
            service.create_multi_factor_signatures(["1"], "")

    def test_sends_the_ids_and_the_kind(self) -> None:
        service, api = _make_service()
        reply = MagicMock()
        reply.id = "batch-7"
        api.multi_factor_signature_service_create_multi_factor_signature_batch.return_value = reply

        result = service.create_multi_factor_signatures(
            ["10", "11"], MultiFactorSignatureEntityType.REQUEST
        )

        assert result.id == "batch-7"
        body = api.multi_factor_signature_service_create_multi_factor_signature_batch.call_args.kwargs[
            "body"
        ]
        assert body.entity_ids == ["10", "11"]
        assert body.entity_type == TgvalidatordMultiFactorSignaturesEntityType.REQUEST

    def test_refuses_an_unknown_kind(self) -> None:
        service, api = _make_service()
        with pytest.raises(ValueError, match="unknown multi-factor signature entity type"):
            service.create_multi_factor_signatures(["10"], "WALLET")
        api.multi_factor_signature_service_create_multi_factor_signature_batch.assert_not_called()


class TestApproveMultiFactorSignature:
    def test_raises_on_empty_id(self) -> None:
        service, _ = _make_service()
        with pytest.raises(ValueError):
            service.approve_multi_factor_signature("", "sig", "c")

    def test_raises_on_empty_signature(self) -> None:
        service, api = _make_service()
        with pytest.raises(ValueError):
            service.approve_multi_factor_signature("mfs-1", "", "c")
        api.multi_factor_signature_service_approve_multi_factor_signature.assert_not_called()

    def test_forwards_the_signature_and_returns_the_count(self) -> None:
        service, api = _make_service()
        reply = MagicMock()
        reply.signatures = "2"
        api.multi_factor_signature_service_approve_multi_factor_signature.return_value = reply

        result = service.approve_multi_factor_signature("mfs-1", "c2ln", "looks fine")

        assert result.signature_count == "2"
        call = api.multi_factor_signature_service_approve_multi_factor_signature.call_args
        assert call.kwargs["id"] == "mfs-1"
        assert call.kwargs["body"].signature == "c2ln"
        assert call.kwargs["body"].comment == "looks fine"


class TestRejectMultiFactorSignature:
    def test_raises_on_empty_id(self) -> None:
        service, _ = _make_service()
        with pytest.raises(ValueError):
            service.reject_multi_factor_signature("", "no")

    def test_sends_the_comment(self) -> None:
        service, api = _make_service()
        service.reject_multi_factor_signature("mfs-1", "not mine")
        call = api.multi_factor_signature_service_reject_multi_factor_signature.call_args
        assert call.kwargs["id"] == "mfs-1"
        assert call.kwargs["body"].comment == "not mine"
