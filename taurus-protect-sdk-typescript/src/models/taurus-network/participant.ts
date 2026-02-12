/**
 * Participant models for Taurus Network.
 *
 * This module provides domain models for Taurus Network participants,
 * which represent organizations or entities that participate in Taurus Network
 * for collateral management, lending, and settlement operations.
 */

/**
 * Attribute associated with a Taurus Network participant.
 */
export interface ParticipantAttribute {
  /** The attribute ID. */
  readonly id: string;
  /** The attribute key. */
  readonly key: string;
  /** The attribute value. */
  readonly value: string;
  /** The content type of the attribute value. */
  readonly contentType: string;
  /** The type of the attribute. */
  readonly attributeType: string;
  /** The subtype of the attribute. */
  readonly subtype: string;
  /** Whether the attribute is shared to other participants. */
  readonly shared: boolean;
  /** When the attribute was created. */
  readonly createdAt: Date | undefined;
  /** When the attribute was last updated. */
  readonly updatedAt: Date | undefined;
}

/**
 * A Taurus Network participant.
 *
 * Represents an organization or entity that participates in Taurus Network
 * for collateral management, lending, and settlement operations.
 */
export interface Participant {
  /** The participant ID. */
  readonly id: string;
  /** The participant name. */
  readonly name: string;
  /** The legal address. */
  readonly legalAddress: string;
  /** The country code. */
  readonly country: string;
  /** The participant's public key. */
  readonly publicKey: string;
  /** The shield identifier. */
  readonly shield: string;
  /** The participant status. */
  readonly status: string;
  /** Public subname for the participant. */
  readonly publicSubname: string;
  /** Legal entity identifier (LEI). */
  readonly legalEntityIdentifier: string;
  /** Count of addresses owned by this participant. */
  readonly ownedSharedAddressesCount: number;
  /** Count of addresses targeting this participant. */
  readonly targetedSharedAddressesCount: number;
  /** Total valuation of outgoing pledges. */
  readonly outgoingTotalPledgesValuation: string;
  /** Total valuation of incoming pledges. */
  readonly incomingTotalPledgesValuation: string;
  /** List of participant attributes. */
  readonly attributes: ParticipantAttribute[];
  /** Date of original registration. */
  readonly originRegistrationDate: Date | undefined;
  /** Date of deletion (if applicable). */
  readonly originDeletionDate: Date | undefined;
  /** When the participant was created. */
  readonly createdAt: Date | undefined;
  /** When the participant was last updated. */
  readonly updatedAt: Date | undefined;
}

/**
 * Settings for the current participant (my participant).
 */
export interface ParticipantSettings {
  /** The current status of participant settings. */
  readonly status: string;
  /** List of country codes allowed for interaction. */
  readonly interactingAllowedCountries: string[];
  /** When terms and conditions were accepted. */
  readonly termsAndConditionsAcceptedAt: Date | undefined;
}

/**
 * The current participant with associated settings.
 *
 * This represents the authenticated tenant's participant information
 * along with their Taurus Network settings.
 */
export interface MyParticipant {
  /** The participant details. */
  readonly participant: Participant | undefined;
  /** The participant settings. */
  readonly settings: ParticipantSettings | undefined;
}

/**
 * Options for retrieving a participant.
 */
export interface GetParticipantOptions {
  /** Include aggregated pledge valuations. */
  readonly includeTotalPledgesValuation?: boolean;
}

/**
 * Options for listing participants.
 */
export interface ListParticipantsOptions {
  /** Filter by specific participant IDs. */
  readonly participantIds?: string[];
  /** Include aggregated pledge valuations. */
  readonly includeTotalPledgesValuation?: boolean;
}

/**
 * Request to create a participant attribute.
 */
export interface CreateParticipantAttributeRequest {
  /** The attribute key. */
  readonly key: string;
  /** The attribute value. */
  readonly value: string;
  /** The content type of the value (optional). */
  readonly contentType?: string;
  /** The type of the attribute (optional). */
  readonly attributeType?: string;
  /** The subtype of the attribute (optional). */
  readonly subtype?: string;
  /** Whether to share to linked participant. */
  readonly shareToTaurusNetworkParticipant?: boolean;
}
