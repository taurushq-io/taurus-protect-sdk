package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.TaurusNetworkMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.Pagination;
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
 *     .list(null, null, null, null, null, 20, null);
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
     * Retrieves a page of pledges with optional filtering.
     *
     * @param ownerParticipantId  filter by owner participant ID (optional)
     * @param targetParticipantId filter by target participant ID (optional)
     * @param sharedAddressIds    filter by shared address IDs (optional)
     * @param currencyId          filter by currency ID (optional)
     * @param sortOrder           sort order for results (optional, "ASC" or "DESC")
     * @param pageSize            the page size, null or 0 for the default
     * @param cursor              a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the pledges and their page
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if the page size is out of range
     */
    public PledgeResult list(final String ownerParticipantId, final String targetParticipantId,
                             final List<String> sharedAddressIds, final String currencyId,
                             final String sortOrder, final Integer pageSize, final String cursor)
            throws ApiException {
        return list(ownerParticipantId, targetParticipantId, sharedAddressIds, currencyId, sortOrder,
                Pagination.page(pageSize, cursor));
    }

    /**
     * Retrieves a page of pledges with optional filtering, with a low-level request cursor.
     *
     * @param ownerParticipantId  filter by owner participant ID (optional)
     * @param targetParticipantId filter by target participant ID (optional)
     * @param sharedAddressIds    filter by shared address IDs (optional)
     * @param currencyId          filter by currency ID (optional)
     * @param sortOrder           sort order for results (optional, "ASC" or "DESC")
     * @param cursor              the request cursor, null for the first page with the default size
     * @return the pledges and their page
     * @throws ApiException if the API call fails
     */
    public PledgeResult list(final String ownerParticipantId, final String targetParticipantId,
                             final List<String> sharedAddressIds, final String currencyId,
                             final String sortOrder, final ApiRequestCursor cursor)
            throws ApiException {

        final CursorRequest page = CursorRequest.of(cursor);

        try {
            TgvalidatordGetPledgesReply reply = pledgeApi.taurusNetworkServiceGetPledges(
                    ownerParticipantId,
                    targetParticipantId,
                    sharedAddressIds,
                    currencyId,
                    sortOrder,
                    page.currentPage(),
                    page.pageRequest(),
                    page.pageSizeParam(),
                    null,  // attributeFiltersJson
                    null,  // statuses
                    null   // attributeFiltersOperator
            );
            return page.complete(mapper.fromPledgesReply(reply), PagedOperation.PLEDGES);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves a page of pledge withdrawals.
     *
     * @param pledgeId         filter by pledge ID (optional)
     * @param withdrawalStatus filter by withdrawal status (optional)
     * @param sortOrder        sort order for results (optional, "ASC" or "DESC")
     * @param pageSize         the page size, null or 0 for the default
     * @param cursor           a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the withdrawals and their page
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if the page size is out of range
     */
    public PledgeWithdrawalResult listWithdrawals(final String pledgeId, final String withdrawalStatus,
                                                  final String sortOrder, final Integer pageSize,
                                                  final String cursor) throws ApiException {
        return listWithdrawals(pledgeId, withdrawalStatus, sortOrder, Pagination.page(pageSize, cursor));
    }

    /**
     * Retrieves a page of pledge withdrawals, with a low-level request cursor.
     *
     * @param pledgeId         filter by pledge ID (optional)
     * @param withdrawalStatus filter by withdrawal status (optional)
     * @param sortOrder        sort order for results (optional, "ASC" or "DESC")
     * @param cursor           the request cursor, null for the first page with the default size
     * @return the withdrawals and their page
     * @throws ApiException if the API call fails
     */
    public PledgeWithdrawalResult listWithdrawals(final String pledgeId, final String withdrawalStatus,
                                                  final String sortOrder, final ApiRequestCursor cursor)
            throws ApiException {
        final CursorRequest page = CursorRequest.of(cursor);

        try {
            TgvalidatordGetPledgesWithdrawalsReply reply = pledgeApi.taurusNetworkServiceGetPledgesWithdrawals(
                    pledgeId,
                    sortOrder,
                    page.currentPage(),
                    page.pageRequest(),
                    page.pageSizeParam(),
                    withdrawalStatus
            );
            return page.complete(mapper.fromPledgeWithdrawalsReply(reply), PagedOperation.PLEDGE_WITHDRAWALS);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
