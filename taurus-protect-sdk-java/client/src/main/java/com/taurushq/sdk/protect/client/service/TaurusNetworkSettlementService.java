package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.TaurusNetworkMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.taurusnetwork.Settlement;
import com.taurushq.sdk.protect.client.model.taurusnetwork.SettlementResult;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.TaurusNetworkSettlementApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetSettlementReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetSettlementsReply;

import java.util.List;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for Taurus Network settlement operations.
 * <p>
 * This service provides access to settlements in the Taurus Network.
 * <p>
 * Example usage:
 * <pre>{@code
 * // List settlements
 * SettlementResult settlements = client.getTaurusNetworkSettlementService()
 *     .getSettlements(null, null, null, null);
 *
 * // Get a specific settlement
 * Settlement settlement = client.getTaurusNetworkSettlementService()
 *     .getSettlement("settlement-123");
 * }</pre>
 */
public class TaurusNetworkSettlementService {

    private final TaurusNetworkSettlementApi settlementApi;
    private final ApiExceptionMapper apiExceptionMapper;
    private final TaurusNetworkMapper mapper;

    /**
     * Instantiates a new Taurus Network Settlement service.
     *
     * @param openApiClient      the OpenAPI client
     * @param apiExceptionMapper the API exception mapper
     */
    public TaurusNetworkSettlementService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.settlementApi = new TaurusNetworkSettlementApi(openApiClient);
        this.mapper = TaurusNetworkMapper.INSTANCE;
    }

    /**
     * Retrieves a settlement by ID.
     *
     * @param settlementId the settlement ID
     * @return the settlement
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if settlementId is null or empty
     */
    public Settlement getSettlement(final String settlementId) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(settlementId), "settlementId cannot be null or empty");

        try {
            TgvalidatordGetSettlementReply reply = settlementApi.taurusNetworkServiceGetSettlement(settlementId);
            return mapper.fromSettlementDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves settlements with optional filtering.
     *
     * @param counterParticipantId filter by counter participant ID (optional)
     * @param statuses             filter by statuses (optional)
     * @param sortOrder            sort order for results (optional, "ASC" or "DESC")
     * @param cursor               pagination cursor (optional, null for first page)
     * @return a paginated result containing settlements
     * @throws ApiException if the API call fails
     */
    public SettlementResult getSettlements(final String counterParticipantId, final List<String> statuses,
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
            TgvalidatordGetSettlementsReply reply = settlementApi.taurusNetworkServiceGetSettlements(
                    counterParticipantId,
                    statuses,
                    sortOrder,
                    cursorCurrentPage,
                    cursorPageRequest,
                    cursorPageSize
            );
            return mapper.fromSettlementsReply(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
