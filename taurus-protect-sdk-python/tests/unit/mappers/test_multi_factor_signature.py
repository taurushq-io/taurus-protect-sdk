"""Unit tests for the multi-factor signature entity-type mapping.

The kind is the only thing in a multi-factor signature reply that tells a caller which
verifying reader to check ``payload_to_sign`` against, so neither direction of the
mapping is allowed to default an unknown value to something plausible.
"""

from __future__ import annotations

import pytest

from taurus_protect._internal.openapi.models.tgvalidatord_multi_factor_signatures_entity_type import (  # noqa: E501
    TgvalidatordMultiFactorSignaturesEntityType,
)
from taurus_protect.errors import IntegrityError
from taurus_protect.models.multi_factor_signature import MultiFactorSignatureEntityType
from taurus_protect.services.multi_factor_signature_service import (
    _entity_type_from_dto,
    _entity_type_to_dto,
)


class TestEntityTypeToDto:
    """Domain enum -> generated enum."""

    @pytest.mark.parametrize(
        "kind",
        [
            MultiFactorSignatureEntityType.REQUEST,
            MultiFactorSignatureEntityType.WHITELISTED_ADDRESS,
            MultiFactorSignatureEntityType.WHITELISTED_CONTRACT,
        ],
    )
    def test_maps_every_kind(self, kind: MultiFactorSignatureEntityType) -> None:
        assert _entity_type_to_dto(kind) == TgvalidatordMultiFactorSignaturesEntityType(
            kind.value
        )

    def test_accepts_the_plain_string_form(self) -> None:
        assert (
            _entity_type_to_dto("WHITELISTED_ADDRESS")
            == TgvalidatordMultiFactorSignaturesEntityType.WHITELISTED_ADDRESS
        )

    def test_refuses_an_unknown_kind(self) -> None:
        """Sending the wrong kind would put an entity through the wrong approval channel."""
        with pytest.raises(ValueError, match="unknown multi-factor signature entity type"):
            _entity_type_to_dto("WALLET")


class TestEntityTypeFromDto:
    """Generated enum -> domain enum."""

    def test_maps_the_generated_enum(self) -> None:
        assert (
            _entity_type_from_dto(TgvalidatordMultiFactorSignaturesEntityType.REQUEST)
            is MultiFactorSignatureEntityType.REQUEST
        )

    def test_maps_a_bare_string(self) -> None:
        assert (
            _entity_type_from_dto("WHITELISTED_CONTRACT")
            is MultiFactorSignatureEntityType.WHITELISTED_CONTRACT
        )

    def test_an_unknown_kind_is_an_integrity_error_not_a_default(self) -> None:
        """Defaulting would tell the caller a payload covers a REQUEST when it does not."""
        with pytest.raises(IntegrityError, match="unknown multi-factor signature entity type"):
            _entity_type_from_dto("SOMETHING_NEW")

    def test_a_missing_kind_is_an_integrity_error(self) -> None:
        with pytest.raises(IntegrityError, match="carries no entity type"):
            _entity_type_from_dto(None)
