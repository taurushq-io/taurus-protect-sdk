package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.TaurusNetworkMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.taurusnetwork.Pledge;
import com.taurushq.sdk.protect.client.model.taurusnetwork.PledgeResult;
import com.taurushq.sdk.protect.client.model.taurusnetwork.PledgeWithdrawalResult;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.TaurusNetworkPledgeApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetPledgeReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetPledgesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetPledgesWithdrawalsReply;

import java.util.List;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for Taurus Network pledge operations.
 * <p>
 * This service provides access to pledge management functionality
 * including creating, retrieving, and managing pledges and withdrawals.
 * <p>
 * Access via: {@code client.taurusNetwork().pledges()}
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get a pledge
 * Pledge pledge = client.taurusNetwork().pledges().get("pledge-id");
 *
 * // List pledges
 * PledgeResult pledges = client.taurusNetwork().pledges()
 *     .list(null, null, null, null, null, null);
 * }</pre>
 */
public class TaurusNetworkPledgeService {

    private final TaurusNetworkPledgeApi pledgeApi;
    private final ApiExceptionMapper apiExceptionMapper;
    private final TaurusNetworkMapper mapper;

    /**
     * Instantiates a new Taurus Network pledge service.
     *
     * @param openApiClient      the OpenAPI client
     * @param apiExceptionMapper the API exception mapper
     */
    public TaurusNetworkPledgeService(final ApiClient openApiClient,
                                      final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.pledgeApi = new TaurusNetworkPledgeApi(openApiClient);
        this.mapper = TaurusNetworkMapper.INSTANCE;
    }

    /**
     * Retrieves a pledge by ID.
     *
     * @param pledgeId the pledge ID
     * @return the pledge
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if pledgeId is null or empty
     */
    public Pledge get(final String pledgeId) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(pledgeId), "pledgeId cannot be null or empty");

        try {
            TgvalidatordGetPledgeReply reply = pledgeApi.taurusNetworkServiceGetPledge(pledgeId);
            return mapper.fromPledgeDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves pledges with optional filtering.
     *
     * @param ownerParticipantId  filter by owner participant ID (optional)
     * @param targetParticipantId filter by target participant ID (optional)
     * @param sharedAddressIds    filter by shared address IDs (optional)
     * @param currencyId          filter by currency ID (optional)
     * @param sortOrder           sort order for results (optional, "ASC" or "DESC")
     * @param cursor              pagination cursor (optional, null for first page)
     * @return a paginated result containing pledges
     * @throws ApiException if the API call fails
     */
    public PledgeResult list(final String ownerParticipantId, final String targetParticipantId,
                             final List<String> sharedAddressIds, final String currencyId,
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
            TgvalidatordGetPledgesReply reply = pledgeApi.taurusNetworkServiceGetPledges(
                    ownerParticipantId,
                    targetParticipantId,
                    sharedAddressIds,
                    currencyId,
                    sortOrder,
                    cursorCurrentPage,
                    cursorPageRequest,
                    cursorPageSize,
                    null,  // attributeFiltersJson
                    null,  // statuses
                    null   // attributeFiltersOperator
            );
            return mapper.fromPledgesReply(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves pledge withdrawals for a specific pledge.
     *
     * @param pledgeId         the pledge ID
     * @param withdrawalStatus filter by withdrawal status (optional)
     * @param sortOrder        sort order for results (optional, "ASC" or "DESC")
     * @param cursor           pagination cursor (optional, null for first page)
     * @return a paginated result containing pledge withdrawals
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if pledgeId is null or empty
     */
    public PledgeWithdrawalResult listWithdrawals(final String pledgeId, final String withdrawalStatus,
                                                  final String sortOrder, final ApiRequestCursor cursor)
            throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(pledgeId), "pledgeId cannot be null or empty");

        String cursorCurrentPage = null;
        String cursorPageRequest = null;
        String cursorPageSize = null;

        if (cursor != null) {
            cursorCurrentPage = cursor.getCurrentPage();
            cursorPageRequest = cursor.getPageRequest() != null ? cursor.getPageRequest().name() : null;
            cursorPageSize = String.valueOf(cursor.getPageSize());
        }

        try {
            TgvalidatordGetPledgesWithdrawalsReply reply = pledgeApi.taurusNetworkServiceGetPledgesWithdrawals(
                    pledgeId,
                    sortOrder,
                    cursorCurrentPage,
                    cursorPageRequest,
                    cursorPageSize,
                    withdrawalStatus
            );
            return mapper.fromPledgeWithdrawalsReply(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
