package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.TaurusNetworkMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.Pagination;
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
 * SettlementResult settlements = client.taurusNetwork().settlements()
 *     .getSettlements(null, null, null, 20, null);
 *
 * // Get a specific settlement
 * Settlement settlement = client.taurusNetwork().settlements()
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
     * Retrieves a page of settlements with optional filtering.
     *
     * @param counterParticipantId filter by counter participant ID (optional)
     * @param statuses             filter by statuses (optional)
     * @param sortOrder            sort order for results (optional, "ASC" or "DESC")
     * @param pageSize             the page size, null or 0 for the default
     * @param cursor               a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the settlements and their page
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if the page size is out of range
     */
    public SettlementResult getSettlements(final String counterParticipantId, final List<String> statuses,
                                           final String sortOrder, final Integer pageSize, final String cursor)
            throws ApiException {
        return getSettlements(counterParticipantId, statuses, sortOrder, Pagination.page(pageSize, cursor));
    }

    /**
     * Retrieves a page of settlements with optional filtering, with a low-level request cursor.
     *
     * @param counterParticipantId filter by counter participant ID (optional)
     * @param statuses             filter by statuses (optional)
     * @param sortOrder            sort order for results (optional, "ASC" or "DESC")
     * @param cursor               the request cursor, null for the first page with the default size
     * @return the settlements and their page
     * @throws ApiException if the API call fails
     */
    public SettlementResult getSettlements(final String counterParticipantId, final List<String> statuses,
                                           final String sortOrder, final ApiRequestCursor cursor)
            throws ApiException {

        final CursorRequest page = CursorRequest.of(cursor);

        try {
            TgvalidatordGetSettlementsReply reply = settlementApi.taurusNetworkServiceGetSettlements(
                    counterParticipantId,
                    statuses,
                    sortOrder,
                    page.currentPage(),
                    page.pageRequest(),
                    page.pageSizeParam()
            );
            return page.complete(mapper.fromSettlementsReply(reply), PagedOperation.SETTLEMENTS);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
