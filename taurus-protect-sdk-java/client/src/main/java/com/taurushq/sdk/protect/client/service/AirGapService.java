package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.AirGapApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetOutgoingAirGapRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetOutgoingAirGapRequestRequests;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSubmitIncomingAirGapRequest;

import java.io.File;
import java.util.List;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for air gap operations in the Taurus Protect system.
 * <p>
 * This service provides operations for transferring data to and from a cold HSM
 * (Hardware Security Module) in an air-gapped environment.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Export requests for cold HSM signing
 * File payload = client.getAirGapService()
 *     .getOutgoingAirGap(Arrays.asList("request-1", "request-2"));
 *
 * // Submit signed responses from cold HSM
 * client.getAirGapService().submitIncomingAirGap(signedPayload);
 * }</pre>
 *
 * @see AirGapApi
 */
public class AirGapService {

    /**
     * The underlying OpenAPI client for air gap operations.
     */
    private final AirGapApi airGapApi;

    /**
     * Mapper for converting OpenAPI exceptions to SDK exceptions.
     */
    private final ApiExceptionMapper apiExceptionMapper;

    /**
     * Instantiates a new Air gap service.
     *
     * @param openApiClient      the OpenAPI client
     * @param apiExceptionMapper the API exception mapper
     */
    public AirGapService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.airGapApi = new AirGapApi(openApiClient);
    }

    /**
     * Exports HSM-ready requests for cold HSM signing.
     * <p>
     * This endpoint returns the payload to be transmitted to the cold HSM
     * for offline signing.
     *
     * @param requestIds the list of request IDs to export
     * @return a file containing the payload for the cold HSM
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if requestIds is null or empty
     */
    public File getOutgoingAirGap(final List<String> requestIds) throws ApiException {
        checkNotNull(requestIds, "requestIds cannot be null");
        checkArgument(!requestIds.isEmpty(), "requestIds cannot be empty");

        try {
            TgvalidatordGetOutgoingAirGapRequestRequests requests = new TgvalidatordGetOutgoingAirGapRequestRequests();
            requests.setIds(requestIds);

            TgvalidatordGetOutgoingAirGapRequest request = new TgvalidatordGetOutgoingAirGapRequest();
            request.setRequests(requests);
            return airGapApi.airGapServiceGetOutgoingAirGap(request);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Imports signed requests from the cold HSM.
     * <p>
     * This endpoint accepts an envelope of signed requests from the cold HSM
     * after offline signing.
     *
     * @param payload the signed payload from the cold HSM (base64-encoded)
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if payload is null or empty
     */
    public void submitIncomingAirGap(final String payload) throws ApiException {
        checkNotNull(payload, "payload cannot be null");
        checkArgument(!payload.isEmpty(), "payload cannot be empty");

        try {
            TgvalidatordSubmitIncomingAirGapRequest request = new TgvalidatordSubmitIncomingAirGapRequest();
            request.setPayload(payload);
            airGapApi.airGapServiceSubmitIncomingAirGap(request);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
