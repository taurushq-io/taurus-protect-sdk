"""Multi-factor signature service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Union

from taurus_protect.errors import APIError, IntegrityError, WhitelistError
from taurus_protect.models.multi_factor_signature import (
    MultiFactorSignatureApprovalResult,
    MultiFactorSignatureEntityType,
    MultiFactorSignatureInfo,
    MultiFactorSignatureResult,
)
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    from taurus_protect._internal.openapi.api.multi_factor_signature_api import (
        MultiFactorSignatureApi,
    )


class MultiFactorSignatureService(BaseService):
    """
    Service for multi-factor signature operations.

    A multi-factor signature is a SECOND approval channel over the same entities this SDK
    otherwise protects: a request, a whitelisted address, or a whitelisted contract. The
    caller here is the second-factor signing device (validatord requires the
    ``RequestMobileAppSigner`` / ``WhitelistedAddressMobileAppSigner`` role), and the
    signature it submits is precisely the artefact a compromised server cannot forge on
    its own.

    Read
    :meth:`get_multi_factor_signature_info` before using this service: the payload it
    returns is unverified server data, and binding it to a verified entity is currently
    the caller's job.

    Four operations, matching Go, Java and TypeScript:
    ``get_multi_factor_signature_info`` / ``create_multi_factor_signatures`` /
    ``approve_multi_factor_signature`` / ``reject_multi_factor_signature``.
    """

    def __init__(
        self,
        api_client: Any,
        mfs_api: "MultiFactorSignatureApi",
    ) -> None:
        super().__init__(api_client)
        self._api = mfs_api

    def get_multi_factor_signature_info(self, id: str) -> MultiFactorSignatureInfo:
        """
        Retrieve information about a multi-factor signature request.

        SECURITY -- ``payload_to_sign`` IS UNVERIFIED SERVER DATA. The caller must bind it
        to a verified entity before signing it.

        This is the one read path in the SDK that returns bytes intended for a signing key
        without verifying them, and it is deliberate-but-unresolved rather than an
        oversight (see ``TODOS.md``). The reason it cannot be fixed here: the reply carries
        only ``{id, payloadToSign[], entityType}``. There is **no entity id**, singular or
        plural, and ``payloadToSign`` is a bare string array with no per-element id or
        hash -- so nothing in the reply can be joined back to the entities the request was
        created for, and this method has nothing to check against.

        What that means for a second-factor signer: a malicious or compromised API server
        can answer this call with the metadata hash of an entity of its choosing -- a
        withdrawal to an attacker address, an attacker-controlled whitelisted address --
        under the expected kind. Sign it and the server holds a valid MobileAppSigner
        approval over an entity nobody reviewed.

        The workable client-side binding, for when the decision lands:
        :meth:`create_multi_factor_signatures` takes the entity IDs, so a caller who keeps
        them can re-read those entities through the verifying reader for that
        ``entity_type`` (``client.requests``, ``client.whitelisted_addresses``,
        ``client.whitelisted_assets``) and require each ``payload_to_sign`` element to
        equal a locally recomputed, verified metadata hash. Until the SDK does that for
        you, do it yourself.

        Args:
            id: The multi-factor signature ID.

        Returns:
            The signature request info, whose ``payload_to_sign`` is UNVERIFIED.

        Raises:
            ValueError: If ``id`` is empty.
            APIError: If the API call fails.
        """
        self._validate_required(id, "id")

        try:
            reply = self._api.multi_factor_signature_service_get_multi_factor_signature_entities_info(
                id=id
            )
            return MultiFactorSignatureInfo(
                id=str(getattr(reply, "id", "") or ""),
                payload_to_sign=list(getattr(reply, "payload_to_sign", None) or []),
                entity_type=_entity_type_from_dto(getattr(reply, "entity_type", None)),
            )
        except Exception as e:
            if isinstance(e, (APIError, IntegrityError, WhitelistError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def create_multi_factor_signatures(
        self,
        entity_ids: List[str],
        entity_type: Union[MultiFactorSignatureEntityType, str],
    ) -> MultiFactorSignatureResult:
        """
        Create a batch of multi-factor signature requests.

        Keep the ``entity_ids`` you pass here. They are the only thing that can bind a
        later :meth:`get_multi_factor_signature_info` reply to real entities -- the reply
        itself carries no entity id.

        Args:
            entity_ids: The entity IDs to put through multi-factor approval.
            entity_type: The kind of entity those IDs refer to.

        Returns:
            The created batch, carrying its ID.

        Raises:
            ValueError: If ``entity_ids`` is empty or ``entity_type`` is missing.
            APIError: If the API call fails.
        """
        if not entity_ids:
            raise ValueError("entity_ids cannot be empty")
        if not entity_type:
            raise ValueError("entity_type cannot be empty")

        try:
            from taurus_protect._internal.openapi.models.tgvalidatord_create_multi_factor_signatures_request import (  # noqa: E501
                TgvalidatordCreateMultiFactorSignaturesRequest,
            )

            body = TgvalidatordCreateMultiFactorSignaturesRequest(
                entityIDs=[str(entity_id) for entity_id in entity_ids],
                entityType=_entity_type_to_dto(entity_type),
            )
            reply = self._api.multi_factor_signature_service_create_multi_factor_signature_batch(
                body=body
            )
            return MultiFactorSignatureResult(id=str(getattr(reply, "id", "") or ""))
        except Exception as e:
            if isinstance(e, (APIError, IntegrityError, WhitelistError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def approve_multi_factor_signature(
        self,
        id: str,
        signature: str,
        comment: str = "",
    ) -> MultiFactorSignatureApprovalResult:
        """
        Approve a multi-factor signature request.

        SECURITY -- THE SIGNATURE IS OPAQUE TO THIS SDK. It is forwarded after a non-empty
        check and nothing here knows or checks what it covers, so this method cannot
        enforce that only verified payloads are approved.

        That is why this site is absent from ``scripts/signing-sites/manifest.json``: the
        manifest enumerates sites where the *SDK* signs, and here the caller does. It is
        also why its absence is not reassuring -- the gate that would have forced a
        ``verifies`` / ``signs-own-bytes`` classification simply does not see this path.

        The caller is responsible for having bound the bytes they signed to a verified
        entity. See :meth:`get_multi_factor_signature_info` for why the SDK cannot do it
        for them yet, and ``TODOS.md`` for the open decision.

        Args:
            id: The multi-factor signature ID.
            signature: Base64 signature of the entities' metadata, produced by the caller.
            comment: An optional approval comment.

        Returns:
            The approval result, carrying the server's signature count.

        Raises:
            ValueError: If ``id`` or ``signature`` is empty.
            APIError: If the API call fails.
        """
        self._validate_required(id, "id")
        self._validate_required(signature, "signature")

        try:
            from taurus_protect._internal.openapi.models.multi_factor_signature_service_approve_multi_factor_signature_body import (  # noqa: E501
                MultiFactorSignatureServiceApproveMultiFactorSignatureBody,
            )

            body = MultiFactorSignatureServiceApproveMultiFactorSignatureBody(
                signature=signature,
                comment=comment or "",
            )
            reply = self._api.multi_factor_signature_service_approve_multi_factor_signature(
                id=id, body=body
            )
            return MultiFactorSignatureApprovalResult(
                signature_count=str(getattr(reply, "signatures", "") or "")
            )
        except Exception as e:
            if isinstance(e, (APIError, IntegrityError, WhitelistError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def reject_multi_factor_signature(self, id: str, comment: str = "") -> None:
        """
        Reject a multi-factor signature request.

        Args:
            id: The multi-factor signature ID.
            comment: An optional rejection comment.

        Raises:
            ValueError: If ``id`` is empty.
            APIError: If the API call fails.
        """
        self._validate_required(id, "id")

        try:
            from taurus_protect._internal.openapi.models.multi_factor_signature_service_reject_multi_factor_signature_body import (  # noqa: E501
                MultiFactorSignatureServiceRejectMultiFactorSignatureBody,
            )

            body = MultiFactorSignatureServiceRejectMultiFactorSignatureBody(
                comment=comment or ""
            )
            self._api.multi_factor_signature_service_reject_multi_factor_signature(
                id=id, body=body
            )
        except Exception as e:
            if isinstance(e, (APIError, IntegrityError, WhitelistError, ValueError)):
                raise
            raise self._handle_error(e) from e


def _entity_type_to_dto(
    entity_type: Union[MultiFactorSignatureEntityType, str],
) -> Any:
    """Map the domain entity type onto the generated enum.

    A value the generated enum does not know is a ``ValueError`` rather than a silent
    substitution: sending the wrong kind would put an entity through the wrong approval
    channel.
    """
    from taurus_protect._internal.openapi.models.tgvalidatord_multi_factor_signatures_entity_type import (  # noqa: E501
        TgvalidatordMultiFactorSignaturesEntityType,
    )

    value = (
        entity_type.value
        if isinstance(entity_type, MultiFactorSignatureEntityType)
        else str(entity_type)
    )
    try:
        return TgvalidatordMultiFactorSignaturesEntityType(value)
    except ValueError as exc:
        raise ValueError(f"unknown multi-factor signature entity type {value!r}") from exc


def _entity_type_from_dto(value: Any) -> MultiFactorSignatureEntityType:
    """Map the generated enum back onto the domain one.

    An unknown kind is refused rather than defaulted. Defaulting would tell the caller a
    payload covers a REQUEST when the server said something else, and the kind is what
    decides which verifying reader they must check the payload against.
    """
    if value is None:
        raise IntegrityError("multi-factor signature reply carries no entity type")
    raw = getattr(value, "value", value)
    try:
        return MultiFactorSignatureEntityType(str(raw))
    except ValueError as exc:
        raise IntegrityError(
            f"unknown multi-factor signature entity type {raw!r}: the kind decides which "
            "verifying reader the payload must be checked against, so it is not defaulted"
        ) from exc
