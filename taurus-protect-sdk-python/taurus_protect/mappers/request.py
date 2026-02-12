"""Request mapper for converting OpenAPI DTOs to domain models."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional

from taurus_protect.mappers._base import safe_datetime, safe_int, safe_string
from taurus_protect.mappers.currency import currency_from_dto
from taurus_protect.models.request import (
    Attribute,
    Request,
    RequestApprovers,
    RequestApproversGroup,
    RequestMetadata,
    RequestParallelApproversGroups,
    RequestStatus,
    RequestTrail,
    SignedRequest,
)

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


def request_from_dto(dto: Any) -> Optional[Request]:
    """
    Convert OpenAPI TgvalidatordRequest to domain Request.

    Args:
        dto: OpenAPI request DTO (TgvalidatordRequest).

    Returns:
        Domain Request model or None if dto is None.
    """
    if dto is None:
        return None

    # Parse status
    status_str = safe_string(getattr(dto, "status", None))
    status = RequestStatus.from_string(status_str) if status_str else RequestStatus.PENDING

    # Extract metadata if present
    metadata = None
    dto_metadata = getattr(dto, "metadata", None)
    if dto_metadata is not None:
        metadata = request_metadata_from_dto(dto_metadata)

    # Extract trails if present
    trails: List[RequestTrail] = []
    dto_trails = getattr(dto, "trails", None)
    if dto_trails is not None:
        trails = [request_trail_from_dto(t) for t in dto_trails]

    # Extract signed requests if present
    signed_requests = signed_requests_from_dto(getattr(dto, "signed_requests", None))

    # Extract approvers if present
    approvers = _approvers_from_dto(getattr(dto, "approvers", None))

    # Extract currency info if present
    currency_info = currency_from_dto(getattr(dto, "currency_info", None))

    # Extract needs_approval_from
    needs_approval_from = getattr(dto, "needs_approval_from", None) or []

    # Extract attributes if present
    attributes = _attributes_from_dto(getattr(dto, "attributes", None))

    return Request(
        id=safe_string(getattr(dto, "id", None)),
        tenant_id=safe_int(getattr(dto, "tenant_id", None)),
        currency=getattr(dto, "currency", None),
        envelope=getattr(dto, "envelope", None),
        status=status,
        trails=trails,
        created_at=safe_datetime(getattr(dto, "creation_date", None)),
        updated_at=safe_datetime(getattr(dto, "update_date", None)),
        metadata=metadata,
        rule=getattr(dto, "rule", None),
        signed_requests=signed_requests,
        type=getattr(dto, "type", None),
        approvers=approvers,
        currency_info=currency_info,
        needs_approval_from=needs_approval_from,
        request_bundle_id=getattr(dto, "request_bundle_id", None),
        external_request_id=getattr(dto, "external_request_id", None),
        attributes=attributes,
    )


def requests_from_dto(dtos: Optional[List[Any]]) -> List[Request]:
    """
    Convert list of OpenAPI request DTOs to domain Requests.

    Args:
        dtos: List of OpenAPI request DTOs.

    Returns:
        List of domain Request models.
    """
    if dtos is None:
        return []
    return [r for dto in dtos if (r := request_from_dto(dto)) is not None]


def request_metadata_from_dto(dto: Any) -> Optional[RequestMetadata]:
    """
    Convert OpenAPI request metadata DTO to domain RequestMetadata.

    Args:
        dto: OpenAPI request metadata DTO.

    Returns:
        Domain RequestMetadata model or None if dto is None.
    """
    if dto is None:
        return None

    return RequestMetadata(
        hash=getattr(dto, "hash", None),
        payload_as_string=getattr(dto, "payload_as_string", None),
    )


def request_trail_from_dto(dto: Any) -> RequestTrail:
    """
    Convert OpenAPI request trail DTO to domain RequestTrail.

    Args:
        dto: OpenAPI request trail DTO.

    Returns:
        Domain RequestTrail model.
    """
    return RequestTrail(
        timestamp=safe_datetime(getattr(dto, "timestamp", None)),
        action=getattr(dto, "action", None),
        user_id=getattr(dto, "user_id", None),
        comment=getattr(dto, "comment", None),
    )


def signed_request_from_dto(dto: Any) -> SignedRequest:
    """
    Map OpenAPI RequestSignedRequest DTO to domain SignedRequest.

    Args:
        dto: OpenAPI RequestSignedRequest DTO.

    Returns:
        Domain SignedRequest model.
    """
    status = None
    if dto.status:
        try:
            status = RequestStatus(dto.status)
        except ValueError:
            status = None

    return SignedRequest(
        id=safe_string(dto.id) or None,
        signed_request=safe_string(dto.signed_request) or None,
        status=status,
        hash=safe_string(dto.hash) or None,
        block=safe_int(dto.block) or None,
        details=safe_string(dto.details) or None,
        creation_date=safe_datetime(dto.creation_date),
        update_date=safe_datetime(dto.update_date),
        broadcast_date=safe_datetime(dto.broadcast_date),
        confirmation_date=safe_datetime(dto.confirmation_date),
    )


def signed_requests_from_dto(dtos: Optional[List[Any]]) -> List[SignedRequest]:
    """
    Map list of OpenAPI RequestSignedRequest DTOs to domain SignedRequests.

    Args:
        dtos: List of OpenAPI RequestSignedRequest DTOs.

    Returns:
        List of domain SignedRequest models.
    """
    if not dtos:
        return []
    return [signed_request_from_dto(dto) for dto in dtos]


def _approvers_from_dto(dto: Any) -> Optional[RequestApprovers]:
    """Convert OpenAPI approvers DTO to domain RequestApprovers."""
    if dto is None:
        return None

    parallel_list = getattr(dto, "parallel", None) or []
    parallel = []
    for pg_dto in parallel_list:
        sequential_list = getattr(pg_dto, "sequential", None) or []
        sequential = []
        for ag_dto in sequential_list:
            sequential.append(
                RequestApproversGroup(
                    external_group_id=getattr(ag_dto, "external_group_id", None)
                    or getattr(ag_dto, "external_group_i_d", None),
                    minimum_signatures=safe_int(
                        getattr(ag_dto, "minimum_signatures", None)
                    ),
                )
            )
        parallel.append(RequestParallelApproversGroups(sequential=sequential))

    return RequestApprovers(parallel=parallel)


def _attributes_from_dto(dtos: Optional[List[Any]]) -> List[Attribute]:
    """Convert OpenAPI attribute DTOs to domain Attributes."""
    if not dtos:
        return []
    result = []
    for dto in dtos:
        result.append(
            Attribute(
                id=safe_int(getattr(dto, "id", None))
                if getattr(dto, "id", None) is not None
                else None,
                key=getattr(dto, "key", None),
                value=getattr(dto, "value", None),
                content_type=getattr(dto, "content_type", None),
                owner=getattr(dto, "owner", None),
                type=getattr(dto, "type", None),
                sub_type=getattr(dto, "sub_type", None),
                is_file=bool(getattr(dto, "is_file", False)),
            )
        )
    return result
