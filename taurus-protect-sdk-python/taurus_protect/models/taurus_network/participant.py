"""Participant models for Taurus Network."""

from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime
from typing import List, Optional


@dataclass
class ParticipantAttribute:
    """
    Attribute associated with a Taurus Network participant.

    Attributes:
        id: The attribute ID.
        key: The attribute key.
        value: The attribute value.
        content_type: The content type of the attribute value.
        attribute_type: The type of the attribute.
        subtype: The subtype of the attribute.
        shared: Whether the attribute is shared to other participants.
        created_at: When the attribute was created.
        updated_at: When the attribute was last updated.
    """

    id: str = ""
    key: str = ""
    value: str = ""
    content_type: str = ""
    attribute_type: str = ""
    subtype: str = ""
    shared: bool = False
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None


@dataclass
class Participant:
    """
    A Taurus Network participant.

    Represents an organization or entity that participates in Taurus Network
    for collateral management, lending, and settlement operations.

    Attributes:
        id: The participant ID.
        name: The participant name.
        legal_address: The legal address.
        country: The country code.
        public_key: The participant's public key.
        shield: The shield identifier.
        status: The participant status.
        public_subname: Public subname for the participant.
        legal_entity_identifier: Legal entity identifier (LEI).
        owned_shared_addresses_count: Count of addresses owned by this participant.
        targeted_shared_addresses_count: Count of addresses targeting this participant.
        outgoing_total_pledges_valuation: Total valuation of outgoing pledges.
        incoming_total_pledges_valuation: Total valuation of incoming pledges.
        attributes: List of participant attributes.
        origin_registration_date: Date of original registration.
        origin_deletion_date: Date of deletion (if applicable).
        created_at: When the participant was created.
        updated_at: When the participant was last updated.
    """

    id: str = ""
    name: str = ""
    legal_address: str = ""
    country: str = ""
    public_key: str = ""
    shield: str = ""
    status: str = ""
    public_subname: str = ""
    legal_entity_identifier: str = ""
    owned_shared_addresses_count: int = 0
    targeted_shared_addresses_count: int = 0
    outgoing_total_pledges_valuation: str = ""
    incoming_total_pledges_valuation: str = ""
    attributes: List[ParticipantAttribute] = field(default_factory=list)
    origin_registration_date: Optional[datetime] = None
    origin_deletion_date: Optional[datetime] = None
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None


@dataclass
class ParticipantSettings:
    """
    Settings for the current participant (my participant).

    Attributes:
        status: The current status of participant settings.
        interacting_allowed_countries: List of country codes allowed for interaction.
        terms_and_conditions_accepted_at: When terms and conditions were accepted.
    """

    status: str = ""
    interacting_allowed_countries: List[str] = field(default_factory=list)
    terms_and_conditions_accepted_at: Optional[datetime] = None


@dataclass
class MyParticipant:
    """
    The current participant with associated settings.

    This represents the authenticated tenant's participant information
    along with their Taurus Network settings.

    Attributes:
        participant: The participant details.
        settings: The participant settings.
    """

    participant: Optional[Participant] = None
    settings: Optional[ParticipantSettings] = None


@dataclass
class GetParticipantOptions:
    """
    Options for retrieving a participant.

    Attributes:
        include_total_pledges_valuation: Include aggregated pledge valuations.
    """

    include_total_pledges_valuation: bool = False


@dataclass
class ListParticipantsOptions:
    """
    Options for listing participants.

    Attributes:
        participant_ids: Filter by specific participant IDs.
        include_total_pledges_valuation: Include aggregated pledge valuations.
    """

    participant_ids: Optional[List[str]] = None
    include_total_pledges_valuation: bool = False


@dataclass
class CreateParticipantAttributeRequest:
    """
    Request to create a participant attribute.

    Attributes:
        key: The attribute key.
        value: The attribute value.
        content_type: The content type of the value (optional).
        attribute_type: The type of the attribute (optional).
        subtype: The subtype of the attribute (optional).
        share_to_taurus_network_participant: Whether to share to linked participant.
    """

    key: str
    value: str
    content_type: str = ""
    attribute_type: str = ""
    subtype: str = ""
    share_to_taurus_network_participant: bool = False
