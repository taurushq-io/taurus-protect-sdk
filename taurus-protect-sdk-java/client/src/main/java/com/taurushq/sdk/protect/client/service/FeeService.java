package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.FeeMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.Fee;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.FeeApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFeesReply;

import java.util.List;

import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for retrieving network fee information in the Taurus Protect system.
 * <p>
 * This service provides access to current network fees for various blockchains,
 * which can be used to estimate transaction costs.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get all current network fees
 * List<Fee> fees = client.getFeeService().getFees();
 * for (Fee fee : fees) {
 *     System.out.println(fee.getKey() + ": " + fee.getValue());
 * }
 * }</pre>
 *
 * @see Fee
 */
public class FeeService {

    private final FeeApi feeApi;
    private final ApiExceptionMapper apiExceptionMapper;
    private final FeeMapper feeMapper;

    /**
     * Instantiates a new Fee service.
     *
     * @param openApiClient      the OpenAPI client
     * @param apiExceptionMapper the API exception mapper
     */
    public FeeService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.feeApi = new FeeApi(openApiClient);
        this.feeMapper = FeeMapper.INSTANCE;
    }

    /**
     * Retrieves current network fees for all supported blockchains.
     *
     * @return the list of fees
     * @throws ApiException if the API call fails
     */
    public List<Fee> getFees() throws ApiException {
        try {
            TgvalidatordGetFeesReply reply = feeApi.feeServiceGetFees();
            return feeMapper.fromDTOList(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
