package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.FeePayerMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.FeePayer;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.FeePayersApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFeePayerReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFeePayersReply;

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
 * // Get all fee payers
 * List<FeePayer> feePayers = client.getFeePayerService().getFeePayers();
 *
 * // Get fee payers for a specific blockchain
 * List<FeePayer> ethFeePayers = client.getFeePayerService()
 *     .getFeePayers(null, null, null, "ETH", "mainnet");
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
     * Retrieves all fee payers.
     *
     * @return the list of fee payers
     * @throws ApiException if the API call fails
     */
    public List<FeePayer> getFeePayers() throws ApiException {
        return getFeePayers(null, null, null, null, null);
    }

    /**
     * Retrieves fee payers with optional filters.
     *
     * @param limit      maximum number of results to return (optional)
     * @param offset     number of results to skip for pagination (optional)
     * @param ids        list of specific IDs to filter by (optional)
     * @param blockchain blockchain to filter by (optional)
     * @param network    network to filter by (optional)
     * @return the list of fee payers matching the filters
     * @throws ApiException if the API call fails
     */
    public List<FeePayer> getFeePayers(
            final Integer limit,
            final Integer offset,
            final List<String> ids,
            final String blockchain,
            final String network
    ) throws ApiException {
        try {
            TgvalidatordGetFeePayersReply reply = feePayersApi.feePayerServiceGetFeePayers(
                    limit != null ? limit.toString() : null,
                    offset != null ? offset.toString() : null,
                    ids,
                    blockchain,
                    network
            );
            return feePayerMapper.fromDTOList(reply.getResult());
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
