"""Participant service for Taurus Network in Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional

from taurus_protect.mappers.taurus_network.participant import (
    my_participant_from_dto,
    participant_from_dto,
    participants_from_dto,
)
from taurus_protect.models.taurus_network.participant import (
    CreateParticipantAttributeRequest,
    GetParticipantOptions,
    ListParticipantsOptions,
    MyParticipant,
    Participant,
)
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    from taurus_protect._internal.openapi.exceptions import ApiException


class ParticipantService(BaseService):
    """
    Service for Taurus Network participant operations.

    Provides methods to retrieve and manage Taurus Network participants,
    including the current tenant's participant information and attributes.

    Example:
        >>> # Get my participant info
        >>> my_participant = client.taurus_network.participants.get_my_participant()
        >>> print(f"My participant: {my_participant.participant.name}")
        >>>
        >>> # List all visible participants
        >>> participants = client.taurus_network.participants.list()
        >>> for p in participants:
        ...     print(f"{p.name}: {p.country}")
        >>>
        >>> # Get specific participant
        >>> opts = GetParticipantOptions(include_total_pledges_valuation=True)
        >>> participant = client.taurus_network.participants.get("participant-id", opts)
    """

    def __init__(
        self,
        api_client: Any,
        participant_api: Any,
    ) -> None:
        """
        Initialize participant service.

        Args:
            api_client: The OpenAPI client instance.
            participant_api: The TaurusNetworkParticipantApi service from OpenAPI client.
        """
        super().__init__(api_client)
        self._participant_api = participant_api

    def get_my_participant(self) -> MyParticipant:
        """
        Get the current participant with settings.

        Returns the participant linked to the current tenant along with
        their Taurus Network settings.

        Returns:
            MyParticipant containing the participant and settings.

        Raises:
            NotFoundError: If the current tenant is not a participant.
            APIError: If API request fails.
        """
        try:
            resp = self._participant_api.taurus_network_service_get_my_participant()

            result = my_participant_from_dto(resp)
            if result is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError("My participant not found")

            return result
        except Exception as e:
            from taurus_protect._internal.openapi.exceptions import ApiException
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            if isinstance(e, ApiException):
                raise self._handle_error(e) from e
            raise self._handle_error(e) from e

    def get(
        self,
        participant_id: str,
        options: Optional[GetParticipantOptions] = None,
    ) -> Participant:
        """
        Get a participant by ID.

        Args:
            participant_id: The participant ID to retrieve.
            options: Optional retrieval options.

        Returns:
            The participant.

        Raises:
            ValueError: If participant_id is empty.
            NotFoundError: If participant not found.
            APIError: If API request fails.
        """
        self._validate_required(participant_id, "participant_id")
        opts = options or GetParticipantOptions()

        try:
            resp = self._participant_api.taurus_network_service_get_participant(
                participant_id=participant_id,
                include_total_pledges_valuation=(
                    opts.include_total_pledges_valuation
                    if opts.include_total_pledges_valuation
                    else None
                ),
            )

            result = getattr(resp, "result", None)
            if result is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Participant {participant_id} not found")

            participant = participant_from_dto(result)
            if participant is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Participant {participant_id} not found")

            return participant
        except Exception as e:
            from taurus_protect._internal.openapi.exceptions import ApiException
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            if isinstance(e, ApiException):
                raise self._handle_error(e) from e
            raise self._handle_error(e) from e

    def list(
        self,
        options: Optional[ListParticipantsOptions] = None,
    ) -> List[Participant]:
        """
        List visible Taurus Network participants.

        Returns participants that are visible to the current tenant's
        participant based on allowed interactions.

        Args:
            options: Optional filtering options.

        Returns:
            List of visible participants.

        Raises:
            APIError: If API request fails.
        """
        opts = options or ListParticipantsOptions()

        try:
            resp = self._participant_api.taurus_network_service_get_participants(
                participant_ids=opts.participant_ids,
                include_total_pledges_valuation=(
                    opts.include_total_pledges_valuation
                    if opts.include_total_pledges_valuation
                    else None
                ),
            )

            result = getattr(resp, "result", None)
            return participants_from_dto(result) if result else []
        except Exception as e:
            from taurus_protect._internal.openapi.exceptions import ApiException
            from taurus_protect.errors import APIError

            if isinstance(e, APIError):
                raise
            if isinstance(e, ApiException):
                raise self._handle_error(e) from e
            raise self._handle_error(e) from e

    def create_participant_attribute(
        self,
        participant_id: str,
        request: CreateParticipantAttributeRequest,
    ) -> None:
        """
        Create an attribute for a participant.

        Creates a new attribute for the specified participant. The attribute
        can optionally be shared to the linked Taurus Network participant.

        Args:
            participant_id: The participant ID.
            request: The attribute creation request.

        Raises:
            ValueError: If participant_id or request is invalid.
            NotFoundError: If participant not found.
            APIError: If API request fails.

        Note:
            Required role: Admin.
        """
        self._validate_required(participant_id, "participant_id")
        if request is None:
            raise ValueError("request cannot be None")
        self._validate_required(request.key, "key")
        self._validate_required(request.value, "value")

        try:
            # Import the body model
            from taurus_protect._internal.openapi.models.taurus_network_service_create_participant_attribute_body import (
                TaurusNetworkServiceCreateParticipantAttributeBody,
            )
            from taurus_protect._internal.openapi.models.tgvalidatord_participant_attribute_data import (
                TgvalidatordParticipantAttributeData,
            )

            # Build the attribute data
            attribute_data = TgvalidatordParticipantAttributeData(
                key=request.key,
                value=request.value,
                content_type=request.content_type if request.content_type else None,
                type=request.attribute_type if request.attribute_type else None,
                subtype=request.subtype if request.subtype else None,
            )

            # Build the request body
            body = TaurusNetworkServiceCreateParticipantAttributeBody(
                attribute_data=attribute_data,
                share_to_taurus_network_participant=request.share_to_taurus_network_participant,
            )

            self._participant_api.taurus_network_service_create_participant_attribute(
                participant_id=participant_id,
                body=body,
            )
        except Exception as e:
            from taurus_protect._internal.openapi.exceptions import ApiException
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            if isinstance(e, ApiException):
                raise self._handle_error(e) from e
            raise self._handle_error(e) from e

    def delete_participant_attribute(
        self,
        participant_id: str,
        attribute_id: str,
    ) -> None:
        """
        Delete an attribute for a participant.

        Deletes the specified attribute from the participant. If the attribute
        was shared to the linked Taurus Network participant, a notification
        will be sent.

        Args:
            participant_id: The participant ID.
            attribute_id: The attribute ID to delete.

        Raises:
            ValueError: If participant_id or attribute_id is empty.
            NotFoundError: If participant or attribute not found.
            APIError: If API request fails.

        Note:
            Required role: Admin.
        """
        self._validate_required(participant_id, "participant_id")
        self._validate_required(attribute_id, "attribute_id")

        try:
            self._participant_api.taurus_network_service_delete_participant_attribute(
                participant_id=participant_id,
                attribute_id=attribute_id,
            )
        except Exception as e:
            from taurus_protect._internal.openapi.exceptions import ApiException
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            if isinstance(e, ApiException):
                raise self._handle_error(e) from e
            raise self._handle_error(e) from e
