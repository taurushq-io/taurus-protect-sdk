/**
 * Participant service for Taurus Network in Taurus-PROTECT SDK.
 *
 * Provides methods for managing Taurus Network participants,
 * including the current tenant's participant information and attributes.
 */

import { NotFoundError, ValidationError } from '../../errors';
import type { TaurusNetworkParticipantApi } from '../../internal/openapi/apis/TaurusNetworkParticipantApi';
import type {
  TgvalidatordTnParticipant,
  TgvalidatordTnParticipantSettings,
} from '../../internal/openapi/models/index';
import { BaseService } from '../base';

/**
 * Participant information from Taurus Network.
 */
export interface Participant {
  id: string;
  name: string;
  legalAddress?: string;
  country?: string;
  logoBase64?: string;
  publicKey?: string;
  shield?: string;
  publicSubname?: string;
  legalEntityIdentifier?: string;
  status?: string;
  originRegistrationDate?: Date;
  originDeletionDate?: Date;
  createdAt?: Date;
  updatedAt?: Date;
  ownedSharedAddressesCount?: string;
  targetedSharedAddressesCount?: string;
  outgoingTotalPledgesValuationBaseCurrency?: string;
  incomingTotalPledgesValuationBaseCurrency?: string;
  attributes?: ParticipantAttribute[];
}

/**
 * Participant attribute.
 */
export interface ParticipantAttribute {
  id?: string;
  key: string;
  value: string;
  contentType?: string;
  type?: string;
  subtype?: string;
  owner?: string;
  isTaurusNetworkShared?: boolean;
}

/**
 * Allowed participant for interactions.
 */
export interface AllowedParticipant {
  id?: string;
  name?: string;
}

/**
 * Participant settings for the current tenant.
 */
export interface ParticipantSettings {
  interactingAllowedCountries?: string[];
  status?: string;
  interactingAllowedParticipants?: AllowedParticipant[];
  termsAndConditionsAcceptedAt?: Date;
}

/**
 * My participant with settings.
 */
export interface MyParticipant {
  participant: Participant;
  settings?: ParticipantSettings;
}

/**
 * Options for getting a participant.
 */
export interface GetParticipantOptions {
  includeTotalPledgesValuation?: boolean;
}

/**
 * Options for listing participants.
 */
export interface ListParticipantsOptions {
  participantIds?: string[];
  includeTotalPledgesValuation?: boolean;
}

/**
 * Request to create a participant attribute.
 */
export interface CreateParticipantAttributeRequest {
  key: string;
  value: string;
  contentType?: string;
  attributeType?: string;
  subtype?: string;
  shareToTaurusNetworkParticipant?: boolean;
}

/**
 * Maps a DTO to a Participant.
 */
function participantFromDto(dto?: TgvalidatordTnParticipant): Participant | undefined {
  if (!dto) {
    return undefined;
  }

  return {
    id: dto.id ?? '',
    name: dto.name ?? '',
    legalAddress: dto.legalAddress,
    country: dto.country,
    logoBase64: dto.logoBase64,
    publicKey: dto.publicKey,
    shield: dto.shield,
    publicSubname: dto.publicSubname,
    legalEntityIdentifier: dto.legalEntityIdentifier,
    status: dto.status,
    originRegistrationDate: dto.originRegistrationDate,
    originDeletionDate: dto.originDeletionDate,
    createdAt: dto.createdAt,
    updatedAt: dto.updatedAt,
    ownedSharedAddressesCount: dto.ownedSharedAddressesCount,
    targetedSharedAddressesCount: dto.targetedSharedAddressesCount,
    outgoingTotalPledgesValuationBaseCurrency: dto.outgoingTotalPledgesValuationBaseCurrency,
    incomingTotalPledgesValuationBaseCurrency: dto.incomingTotalPledgesValuationBaseCurrency,
    attributes: dto.attributes?.map((attr) => ({
      id: attr.id,
      key: attr.key ?? '',
      value: attr.value ?? '',
      contentType: attr.contentType,
      type: attr.type,
      subtype: attr.subtype,
      owner: attr.owner,
      isTaurusNetworkShared: attr.isTaurusNetworkShared,
    })),
  };
}

/**
 * Maps settings DTO to ParticipantSettings.
 */
function settingsFromDto(dto?: TgvalidatordTnParticipantSettings): ParticipantSettings | undefined {
  if (!dto) {
    return undefined;
  }

  return {
    interactingAllowedCountries: dto.interactingAllowedCountries,
    status: dto.status,
    interactingAllowedParticipants: dto.interactingAllowedParticipants?.map((p) => ({
      id: p.id,
      name: p.name,
    })),
    termsAndConditionsAcceptedAt: dto.termsAndConditionsAcceptedAt,
  };
}

/**
 * Service for Taurus Network participant operations.
 *
 * Provides methods to retrieve and manage Taurus Network participants,
 * including the current tenant's participant information and attributes.
 *
 * @example
 * ```typescript
 * // Get my participant info
 * const myParticipant = await participantService.getMyParticipant();
 * console.log(`My participant: ${myParticipant.participant.name}`);
 *
 * // List all visible participants
 * const participants = await participantService.list();
 * for (const p of participants) {
 *   console.log(`${p.name}: ${p.country}`);
 * }
 *
 * // Get specific participant
 * const participant = await participantService.get('participant-id', {
 *   includeTotalPledgesValuation: true,
 * });
 * ```
 */
export class ParticipantService extends BaseService {
  private readonly participantApi: TaurusNetworkParticipantApi;

  /**
   * Creates a new ParticipantService instance.
   *
   * @param participantApi - The TaurusNetworkParticipantApi instance from the OpenAPI client
   */
  constructor(participantApi: TaurusNetworkParticipantApi) {
    super();
    this.participantApi = participantApi;
  }

  /**
   * Gets the current participant with settings.
   *
   * Returns the participant linked to the current tenant along with
   * their Taurus Network settings.
   *
   * @returns MyParticipant containing the participant and settings
   * @throws {@link NotFoundError} If the current tenant is not a participant
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const myParticipant = await participantService.getMyParticipant();
   * console.log(`My participant: ${myParticipant.participant.name}`);
   * console.log(`Settings: ${JSON.stringify(myParticipant.settings)}`);
   * ```
   */
  async getMyParticipant(): Promise<MyParticipant> {
    return this.execute(async () => {
      const response = await this.participantApi.taurusNetworkServiceGetMyParticipant();

      const participant = participantFromDto(response.result);
      if (!participant) {
        throw new NotFoundError('My participant not found');
      }

      return {
        participant,
        settings: settingsFromDto(response.settings),
      };
    });
  }

  /**
   * Gets a participant by ID.
   *
   * @param participantId - The participant ID to retrieve
   * @param options - Optional retrieval options
   * @returns The participant
   * @throws {@link ValidationError} If participantId is empty
   * @throws {@link NotFoundError} If participant not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const participant = await participantService.get('participant-123', {
   *   includeTotalPledgesValuation: true,
   * });
   * console.log(`Participant: ${participant.name}`);
   * ```
   */
  async get(participantId: string, options?: GetParticipantOptions): Promise<Participant> {
    if (!participantId || participantId.trim() === '') {
      throw new ValidationError('participantId is required');
    }

    return this.execute(async () => {
      const response = await this.participantApi.taurusNetworkServiceGetParticipant({
        participantID: participantId,
        includeTotalPledgesValuation: options?.includeTotalPledgesValuation,
      });

      const participant = participantFromDto(response.result);
      if (!participant) {
        throw new NotFoundError(`Participant ${participantId} not found`);
      }

      return participant;
    });
  }

  /**
   * Lists visible Taurus Network participants.
   *
   * Returns participants that are visible to the current tenant's
   * participant based on allowed interactions.
   *
   * @param options - Optional filtering options
   * @returns List of visible participants
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // List all participants
   * const participants = await participantService.list();
   *
   * // List specific participants with valuation
   * const filtered = await participantService.list({
   *   participantIds: ['id1', 'id2'],
   *   includeTotalPledgesValuation: true,
   * });
   * ```
   */
  async list(options?: ListParticipantsOptions): Promise<Participant[]> {
    return this.execute(async () => {
      const response = await this.participantApi.taurusNetworkServiceGetParticipants({
        participantIDs: options?.participantIds,
        includeTotalPledgesValuation: options?.includeTotalPledgesValuation,
      });

      const participants: Participant[] = [];
      if (response.result) {
        for (const dto of response.result) {
          const participant = participantFromDto(dto);
          if (participant) {
            participants.push(participant);
          }
        }
      }

      return participants;
    });
  }

  /**
   * Creates an attribute for a participant.
   *
   * Creates a new attribute for the specified participant. The attribute
   * can optionally be shared to the linked Taurus Network participant.
   *
   * @param participantId - The participant ID
   * @param request - The attribute creation request
   * @throws {@link ValidationError} If participantId or request is invalid
   * @throws {@link NotFoundError} If participant not found
   * @throws {@link APIError} If API request fails
   *
   * Note: Required role is Admin.
   *
   * @example
   * ```typescript
   * await participantService.createParticipantAttribute('participant-123', {
   *   key: 'department',
   *   value: 'treasury',
   *   shareToTaurusNetworkParticipant: true,
   * });
   * ```
   */
  async createParticipantAttribute(
    participantId: string,
    request: CreateParticipantAttributeRequest
  ): Promise<void> {
    if (!participantId || participantId.trim() === '') {
      throw new ValidationError('participantId is required');
    }
    if (!request.key || request.key.trim() === '') {
      throw new ValidationError('key is required');
    }
    if (!request.value || request.value.trim() === '') {
      throw new ValidationError('value is required');
    }

    return this.execute(async () => {
      await this.participantApi.taurusNetworkServiceCreateParticipantAttribute({
        participantID: participantId,
        body: {
          attributeData: {
            key: request.key,
            value: request.value,
            contentType: request.contentType,
            type: request.attributeType,
            subtype: request.subtype,
          },
          shareToTaurusNetworkParticipant: request.shareToTaurusNetworkParticipant,
        },
      });
    });
  }

  /**
   * Deletes an attribute for a participant.
   *
   * Deletes the specified attribute from the participant. If the attribute
   * was shared to the linked Taurus Network participant, a notification
   * will be sent.
   *
   * @param participantId - The participant ID
   * @param attributeId - The attribute ID to delete
   * @throws {@link ValidationError} If participantId or attributeId is empty
   * @throws {@link NotFoundError} If participant or attribute not found
   * @throws {@link APIError} If API request fails
   *
   * Note: Required role is Admin.
   *
   * @example
   * ```typescript
   * await participantService.deleteParticipantAttribute('participant-123', 'attr-456');
   * ```
   */
  async deleteParticipantAttribute(participantId: string, attributeId: string): Promise<void> {
    if (!participantId || participantId.trim() === '') {
      throw new ValidationError('participantId is required');
    }
    if (!attributeId || attributeId.trim() === '') {
      throw new ValidationError('attributeId is required');
    }

    return this.execute(async () => {
      await this.participantApi.taurusNetworkServiceDeleteParticipantAttribute({
        participantID: participantId,
        attributeID: attributeId,
      });
    });
  }
}
