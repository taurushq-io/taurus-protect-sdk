package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.TaurusNetworkMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.Pagination;
import com.taurushq.sdk.protect.client.model.taurusnetwork.LendingAgreement;
import com.taurushq.sdk.protect.client.model.taurusnetwork.LendingAgreementResult;
import com.taurushq.sdk.protect.client.model.taurusnetwork.LendingOffer;
import com.taurushq.sdk.protect.client.model.taurusnetwork.LendingOfferResult;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.TaurusNetworkLendingApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetLendingAgreementReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetLendingAgreementsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetLendingOfferReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetLendingOffersReply;

import java.util.List;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for Taurus Network lending operations.
 * <p>
 * This service provides access to lending offers and agreements
 * in the Taurus Network.
 * <p>
 * Example usage:
 * <pre>{@code
 * // List lending offers, first page
 * LendingOfferResult offers = client.taurusNetwork().lending()
 *     .getLendingOffers(null, null, null, null, 20, null);
 *
 * // Get a specific lending agreement
 * LendingAgreement agreement = client.taurusNetwork().lending()
 *     .getLendingAgreement("agreement-123");
 * }</pre>
 */
public class TaurusNetworkLendingService {

    private final TaurusNetworkLendingApi lendingApi;
    private final ApiExceptionMapper apiExceptionMapper;
    private final TaurusNetworkMapper mapper;

    /**
     * Instantiates a new Taurus Network Lending service.
     *
     * @param openApiClient      the OpenAPI client
     * @param apiExceptionMapper the API exception mapper
     */
    public TaurusNetworkLendingService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.lendingApi = new TaurusNetworkLendingApi(openApiClient);
        this.mapper = TaurusNetworkMapper.INSTANCE;
    }

    /**
     * Retrieves a lending offer by ID.
     *
     * @param offerId the lending offer ID
     * @return the lending offer
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if offerId is null or empty
     */
    public LendingOffer getLendingOffer(final String offerId) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(offerId), "offerId cannot be null or empty");

        try {
            TgvalidatordGetLendingOfferReply reply = lendingApi.taurusNetworkServiceGetLendingOffer(offerId);
            return mapper.fromLendingOfferDTO(reply.getLendingOffer());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves a page of lending offers with optional filtering.
     *
     * @param currencyIds   filter by currency IDs (optional)
     * @param participantId filter by participant ID (optional)
     * @param duration      filter by duration (optional)
     * @param sortOrder     sort order for results (optional, "ASC" or "DESC")
     * @param pageSize      the page size, null or 0 for the default
     * @param cursor        a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the lending offers and their page
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if the page size is out of range
     */
    public LendingOfferResult getLendingOffers(final List<String> currencyIds, final String participantId,
                                               final String duration, final String sortOrder,
                                               final Integer pageSize, final String cursor)
            throws ApiException {
        return getLendingOffers(currencyIds, participantId, duration, sortOrder, Pagination.page(pageSize, cursor));
    }

    /**
     * Retrieves a page of lending offers, with a low-level request cursor.
     *
     * @param currencyIds   filter by currency IDs (optional)
     * @param participantId filter by participant ID (optional)
     * @param duration      filter by duration (optional)
     * @param sortOrder     sort order for results (optional, "ASC" or "DESC")
     * @param cursor        the request cursor, null for the first page with the default size
     * @return the lending offers and their page
     * @throws ApiException if the API call fails
     */
    public LendingOfferResult getLendingOffers(final List<String> currencyIds, final String participantId,
                                               final String duration, final String sortOrder,
                                               final ApiRequestCursor cursor)
            throws ApiException {

        final CursorRequest page = CursorRequest.of(cursor);

        try {
            TgvalidatordGetLendingOffersReply reply = lendingApi.taurusNetworkServiceGetLendingOffers(
                    sortOrder,
                    page.currentPage(),
                    page.pageRequest(),
                    page.pageSizeParam(),
                    currencyIds,
                    participantId,
                    duration
            );
            return page.complete(mapper.fromLendingOffersReply(reply), PagedOperation.LENDING_OFFERS);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves a lending agreement by ID.
     *
     * @param agreementId the lending agreement ID
     * @return the lending agreement
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if agreementId is null or empty
     */
    public LendingAgreement getLendingAgreement(final String agreementId) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(agreementId), "agreementId cannot be null or empty");

        try {
            TgvalidatordGetLendingAgreementReply reply = lendingApi.taurusNetworkServiceGetLendingAgreement(agreementId);
            return mapper.fromLendingAgreementDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves a page of lending agreements.
     *
     * @param sortOrder sort order for results (optional, "ASC" or "DESC")
     * @param pageSize  the page size, null or 0 for the default
     * @param cursor    a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the lending agreements and their page
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if the page size is out of range
     */
    public LendingAgreementResult getLendingAgreements(final String sortOrder, final Integer pageSize,
                                                       final String cursor) throws ApiException {
        return getLendingAgreements(sortOrder, Pagination.page(pageSize, cursor));
    }

    /**
     * Retrieves a page of lending agreements, with a low-level request cursor.
     *
     * @param sortOrder sort order for results (optional, "ASC" or "DESC")
     * @param cursor    the request cursor, null for the first page with the default size
     * @return the lending agreements and their page
     * @throws ApiException if the API call fails
     */
    public LendingAgreementResult getLendingAgreements(final String sortOrder, final ApiRequestCursor cursor)
            throws ApiException {

        final CursorRequest page = CursorRequest.of(cursor);

        try {
            TgvalidatordGetLendingAgreementsReply reply = lendingApi.taurusNetworkServiceGetLendingAgreements(
                    sortOrder,
                    page.currentPage(),
                    page.pageRequest(),
                    page.pageSizeParam()
            );
            return page.complete(mapper.fromLendingAgreementsReply(reply), PagedOperation.LENDING_AGREEMENTS);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
