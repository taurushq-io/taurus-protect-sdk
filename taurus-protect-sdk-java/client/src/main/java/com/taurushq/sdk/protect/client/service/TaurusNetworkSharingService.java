package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.TaurusNetworkMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.taurusnetwork.SharedAddressResult;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.TaurusNetworkSharedAddressAssetApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetSharedAddressesReply;

import java.util.List;

import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for Taurus Network sharing operations.
 * <p>
 * This service provides access to address and asset sharing functionality
 * between participants in the Taurus Network.
 * <p>
 * Access via: {@code client.taurusNetwork().sharing()}
 * <p>
 * Example usage:
 * <pre>{@code
 * // List shared addresses
 * SharedAddressResult result = client.taurusNetwork().sharing()
 *     .listSharedAddresses(null, null, null, "ETH", "mainnet", null, null, null);
 * }</pre>
 */
public class TaurusNetworkSharingService {

    private final TaurusNetworkSharedAddressAssetApi sharedAddressApi;
    private final ApiExceptionMapper apiExceptionMapper;
    private final TaurusNetworkMapper mapper;

    /**
     * Instantiates a new Taurus Network sharing service.
     *
     * @param openApiClient      the OpenAPI client
     * @param apiExceptionMapper the API exception mapper
     */
    public TaurusNetworkSharingService(final ApiClient openApiClient,
                                       final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.sharedAddressApi = new TaurusNetworkSharedAddressAssetApi(openApiClient);
        this.mapper = TaurusNetworkMapper.INSTANCE;
    }

    /**
     * Retrieves shared addresses with optional filtering.
     *
     * @param participantId       filter by participant ID (optional)
     * @param ownerParticipantId  filter by owner participant ID (optional)
     * @param targetParticipantId filter by target participant ID (optional)
     * @param blockchain          filter by blockchain (optional)
     * @param network             filter by network (optional)
     * @param ids                 filter by shared address IDs (optional)
     * @param sortOrder           sort order for results (optional, "ASC" or "DESC")
     * @param cursor              pagination cursor (optional, null for first page)
     * @return a paginated result containing shared addresses
     * @throws ApiException if the API call fails
     */
    public SharedAddressResult listSharedAddresses(final String participantId, final String ownerParticipantId,
                                                   final String targetParticipantId, final String blockchain,
                                                   final String network, final List<String> ids,
                                                   final String sortOrder, final ApiRequestCursor cursor)
            throws ApiException {

        String cursorCurrentPage = null;
        String cursorPageRequest = null;
        String cursorPageSize = null;

        if (cursor != null) {
            cursorCurrentPage = cursor.getCurrentPage();
            cursorPageRequest = cursor.getPageRequest() != null ? cursor.getPageRequest().name() : null;
            cursorPageSize = String.valueOf(cursor.getPageSize());
        }

        try {
            TgvalidatordGetSharedAddressesReply reply = sharedAddressApi.taurusNetworkServiceGetSharedAddresses(
                    participantId,
                    ownerParticipantId,
                    targetParticipantId,
                    blockchain,
                    network,
                    ids,
                    sortOrder,
                    cursorCurrentPage,
                    cursorPageRequest,
                    cursorPageSize,
                    null  // statuses
            );
            return mapper.fromSharedAddressesReply(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
