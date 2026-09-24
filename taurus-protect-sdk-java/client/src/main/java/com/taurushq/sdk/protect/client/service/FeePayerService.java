package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.FeePayerMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.FeePayer;
import com.taurushq.sdk.protect.client.model.FeePayerResult;
import com.taurushq.sdk.protect.client.model.Pagination;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.FeePayersApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFeePayerReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordFeePayerEnvelope;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFeePayersReply;

import java.util.Collections;
import java.util.List;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for managing fee payers in the Taurus Protect system.
 * <p>
 * Fee payers are accounts used to pay transaction fees on behalf of other
 * addresses, commonly used for sponsored transactions on EVM-compatible
 * blockchains like Ethereum.
 * <p>
 * Example usage:
 * <pre>{@code
 * // First page of fee payers (default page size)
 * FeePayerResult feePayers = client.getFeePayerService().getFeePayers();
 *
 * // Fee payers for a specific blockchain
 * FeePayerResult ethFeePayers = client.getFeePayerService()
 *     .getFeePayers(0, 0, null, "ETH", "mainnet");
 *
 * // Get a specific fee payer
 * FeePayer feePayer = client.getFeePayerService().getFeePayer("fp-123");
 * }</pre>
 *
 * @see FeePayer
 */
public class FeePayerService {

    private final FeePayersApi feePayersApi;
    private final ApiExceptionMapper apiExceptionMapper;
    private final FeePayerMapper feePayerMapper;

    /**
     * Creates a new FeePayerService.
     *
     * @param apiClient          the API client for making HTTP requests
     * @param apiExceptionMapper the mapper for converting API exceptions
     * @throws NullPointerException if any parameter is null
     */
    public FeePayerService(final ApiClient apiClient, final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(apiClient, "apiClient must not be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper must not be null");
        this.feePayersApi = new FeePayersApi(apiClient);
        this.apiExceptionMapper = apiExceptionMapper;
        this.feePayerMapper = FeePayerMapper.INSTANCE;
    }

    /**
     * Retrieves the first page of fee payers, with the default page size.
     *
     * @return the fee payers and their pagination
     * @throws ApiException if the API call fails
     */
    public FeePayerResult getFeePayers() throws ApiException {
        return getFeePayers(0, 0, null, null, null);
    }

    /**
     * Retrieves a page of fee payers with optional filters.
     *
     * @param limit      the page size, 0 for the default ({@link Pagination#DEFAULT_PAGE_SIZE})
     * @param offset     the offset, 0 for the first page
     * @param ids        list of specific IDs to filter by (optional)
     * @param blockchain blockchain to filter by (optional)
     * @param network    network to filter by (optional)
     * @return the fee payers and their pagination
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if limit or offset is out of range
     */
    public FeePayerResult getFeePayers(
            final int limit,
            final long offset,
            final List<String> ids,
            final String blockchain,
            final String network
    ) throws ApiException {
        final int size = PagedOperation.FEE_PAYERS.resolveSize("limit", limit);
        final long from = Pagination.resolveOffset("offset", offset);
        try {
            TgvalidatordGetFeePayersReply reply = feePayersApi.feePayerServiceGetFeePayers(
                    String.valueOf(size),
                    from == 0 ? null : String.valueOf(from),
                    ids,
                    blockchain,
                    network
            );
            List<TgvalidatordFeePayerEnvelope> rows = reply.getResult() == null
                    ? Collections.emptyList() : reply.getResult();
            return new FeePayerResult(feePayerMapper.fromDTOList(rows),
                    PagedOperation.FEE_PAYERS.offsetPage(size, from, rows.size(), 0,
                            reply.getTotalItems(), null));
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves a fee payer by ID.
     *
     * @param id the fee payer ID
     * @return the fee payer
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if id is null or empty
     */
    public FeePayer getFeePayer(final String id) throws ApiException {
        checkArgument(id != null && !id.isEmpty(), "id must not be null or empty");
        try {
            TgvalidatordGetFeePayerReply reply = feePayersApi.feePayerServiceGetFeePayer(id);
            return feePayerMapper.fromDTO(reply.getFeepayer());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
