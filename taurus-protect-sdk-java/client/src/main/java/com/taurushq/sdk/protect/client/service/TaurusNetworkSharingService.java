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
 *     .listSharedAddresses(null, null, null, "ETH", "mainnet", null, null, Pagination.page(20, null));
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
     * Retrieves a page of shared addresses with optional filtering.
     * <p>
     * Seven filters leave no room for separate page-size and cursor parameters (Checkstyle
     * caps a method at eight), so this list takes the page as one request cursor: pass
     * {@code Pagination.page(pageSize, cursor)}, with {@code cursor} a previous page's
     * {@code getPage().getNextCursor()}.
     *
     * @param participantId       filter by participant ID (optional)
     * @param ownerParticipantId  filter by owner participant ID (optional)
     * @param targetParticipantId filter by target participant ID (optional)
     * @param blockchain          filter by blockchain (optional)
     * @param network             filter by network (optional)
     * @param ids                 filter by shared address IDs (optional)
     * @param sortOrder           sort order for results (optional, "ASC" or "DESC")
     * @param cursor              the page, e.g. {@code Pagination.page(20, next)}; null for the
     *                            first page with the default size
     * @return the shared addresses and their page
     * @throws ApiException if the API call fails
     */
    public SharedAddressResult listSharedAddresses(final String participantId, final String ownerParticipantId,
                                                   final String targetParticipantId, final String blockchain,
                                                   final String network, final List<String> ids,
                                                   final String sortOrder, final ApiRequestCursor cursor)
            throws ApiException {

        final CursorRequest page = CursorRequest.of(cursor);

        try {
            TgvalidatordGetSharedAddressesReply reply = sharedAddressApi.taurusNetworkServiceGetSharedAddresses(
                    participantId,
                    ownerParticipantId,
                    targetParticipantId,
                    blockchain,
                    network,
                    ids,
                    sortOrder,
                    page.currentPage(),
                    page.pageRequest(),
                    page.pageSizeParam(),
                    null  // statuses
            );
            return page.complete(mapper.fromSharedAddressesReply(reply), PagedOperation.SHARED_ADDRESSES);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
