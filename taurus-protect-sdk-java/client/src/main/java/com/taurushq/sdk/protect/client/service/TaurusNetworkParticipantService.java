package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.TaurusNetworkMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.taurusnetwork.Participant;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.TaurusNetworkParticipantApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetMyParticipantReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetParticipantReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetParticipantsReply;

import java.util.List;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for Taurus Network participant operations.
 * <p>
 * This service provides access to participant management functionality
 * including retrieving the current participant and other connected participants.
 * <p>
 * Access via: {@code client.taurusNetwork().participants()}
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get current participant
 * Participant me = client.taurusNetwork().participants().getMyParticipant();
 *
 * // Get specific participant
 * Participant participant = client.taurusNetwork().participants()
 *     .get("participant-id", true);
 * }</pre>
 */
public class TaurusNetworkParticipantService {

    private final TaurusNetworkParticipantApi participantApi;
    private final ApiExceptionMapper apiExceptionMapper;
    private final TaurusNetworkMapper mapper;

    /**
     * Instantiates a new Taurus Network participant service.
     *
     * @param openApiClient      the OpenAPI client
     * @param apiExceptionMapper the API exception mapper
     */
    public TaurusNetworkParticipantService(final ApiClient openApiClient,
                                           final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.participantApi = new TaurusNetworkParticipantApi(openApiClient);
        this.mapper = TaurusNetworkMapper.INSTANCE;
    }

    /**
     * Retrieves the current participant.
     *
     * @return the current participant
     * @throws ApiException if the API call fails
     */
    public Participant getMyParticipant() throws ApiException {
        try {
            TgvalidatordGetMyParticipantReply reply = participantApi.taurusNetworkServiceGetMyParticipant();
            return mapper.fromParticipantDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves a participant by ID.
     *
     * @param participantId                the participant ID
     * @param includeTotalPledgesValuation whether to include total pledges valuation
     * @return the participant
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if participantId is null or empty
     */
    public Participant get(final String participantId, final Boolean includeTotalPledgesValuation)
            throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(participantId), "participantId cannot be null or empty");

        try {
            TgvalidatordGetParticipantReply reply = participantApi.taurusNetworkServiceGetParticipant(
                    participantId, includeTotalPledgesValuation);
            return mapper.fromParticipantDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves multiple participants by IDs.
     *
     * @param participantIds               the list of participant IDs
     * @param includeTotalPledgesValuation whether to include total pledges valuation
     * @return the list of participants
     * @throws ApiException if the API call fails
     */
    public List<Participant> list(final List<String> participantIds,
                                  final Boolean includeTotalPledgesValuation)
            throws ApiException {
        try {
            TgvalidatordGetParticipantsReply reply = participantApi.taurusNetworkServiceGetParticipants(
                    participantIds, includeTotalPledgesValuation);
            return mapper.fromParticipantDTOList(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
