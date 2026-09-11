"""Multi-factor signature domain models."""

from __future__ import annotations

from enum import Enum
from typing import List

from pydantic import BaseModel, Field


class MultiFactorSignatureEntityType(str, Enum):
    """The kind of entity a multi-factor signature request covers.

    These are the three kinds validatord accepts, and each has a verifying reader in this
    SDK -- which is what a caller needs in order to bind a ``payload_to_sign`` element to
    an entity they have actually checked. See
    :meth:`~taurus_protect.services.multi_factor_signature_service.MultiFactorSignatureService.get_multi_factor_signature_info`.
    """

    REQUEST = "REQUEST"
    WHITELISTED_ADDRESS = "WHITELISTED_ADDRESS"
    WHITELISTED_CONTRACT = "WHITELISTED_CONTRACT"


class MultiFactorSignatureInfo(BaseModel):
    """Information about a multi-factor signature request."""

    id: str = Field(description="The multi-factor signature ID")
    payload_to_sign: List[str] = Field(
        default_factory=list,
        description=(
            "UNVERIFIED SERVER DATA. Unlike every other field this SDK hands back, these "
            "bytes have not been checked against anything -- the reply carries no entity "
            "id to check them against. Bind them to a verified entity before signing; see "
            "MultiFactorSignatureService.get_multi_factor_signature_info."
        ),
    )
    entity_type: MultiFactorSignatureEntityType = Field(
        description="The kind of entity this signature request covers"
    )

    model_config = {"frozen": True}


class MultiFactorSignatureResult(BaseModel):
    """The result of creating a batch of multi-factor signature requests."""

    id: str = Field(default="", description="ID of the created multi-factor signature batch")

    model_config = {"frozen": True}


class MultiFactorSignatureApprovalResult(BaseModel):
    """The result of approving a multi-factor signature request."""

    signature_count: str = Field(
        default="", description="Number of signatures applied, as returned by the server"
    )

    model_config = {"frozen": True}
