package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.ExchangeMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.Exchange;
import com.taurushq.sdk.protect.client.model.ExchangeCounterparty;
import com.taurushq.sdk.protect.client.model.ExchangeWithdrawalFee;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.ExchangeApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordExportExchangesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetExchangeCounterpartiesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetExchangeReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetExchangeWithdrawalFeeReply;

import java.util.List;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for managing exchange accounts in the Taurus Protect system.
 * <p>
 * This service provides access to exchange account information, including
 * balances, counterparty summaries, and withdrawal fee calculations.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get an exchange account
 * Exchange exchange = client.getExchangeService().getExchange("exchange-123");
 *
 * // Get all counterparties (summary by exchange)
 * List<ExchangeCounterparty> counterparties = client.getExchangeService()
 *     .getExchangeCounterparties();
 *
 * // Calculate withdrawal fee
 * ExchangeWithdrawalFee fee = client.getExchangeService()
 *     .getExchangeWithdrawalFee("exchange-123", "address-456", "1000000000");
 * }</pre>
 *
 * @see Exchange
 * @see ExchangeCounterparty
 * @see ExchangeWithdrawalFee
 */
public class ExchangeService {

    private final ExchangeApi exchangeApi;
    private final ApiExceptionMapper apiExceptionMapper;
    private final ExchangeMapper exchangeMapper;

    /**
     * Instantiates a new Exchange service.
     *
     * @param openApiClient      the OpenAPI client
     * @param apiExceptionMapper the API exception mapper
     */
    public ExchangeService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.exchangeApi = new ExchangeApi(openApiClient);
        this.exchangeMapper = ExchangeMapper.INSTANCE;
    }

    /**
     * Retrieves a specific exchange account by ID.
     *
     * @param id the exchange account ID
     * @return the exchange account details
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if id is null or empty
     */
    public Exchange getExchange(final String id) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(id), "id cannot be null or empty");

        try {
            TgvalidatordGetExchangeReply reply = exchangeApi.exchangeServiceGetExchange(id);
            return exchangeMapper.fromDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves all exchange counterparties.
     * <p>
     * Returns a summary of holdings grouped by exchange name.
     *
     * @return the list of exchange counterparties
     * @throws ApiException if the API call fails
     */
    public List<ExchangeCounterparty> getExchangeCounterparties() throws ApiException {
        try {
            TgvalidatordGetExchangeCounterpartiesReply reply =
                    exchangeApi.exchangeServiceGetExchangeCounterparties();
            return exchangeMapper.fromCounterpartyDTOList(reply.getExchanges());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Calculates the withdrawal fee for a transfer from an exchange.
     *
     * @param exchangeId  the exchange account ID
     * @param toAddressId the destination address ID
     * @param amount      the amount to withdraw (in smallest currency unit)
     * @return the withdrawal fee details
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if exchangeId is null or empty
     */
    public ExchangeWithdrawalFee getExchangeWithdrawalFee(final String exchangeId,
                                                          final String toAddressId,
                                                          final String amount) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(exchangeId), "exchangeId cannot be null or empty");

        try {
            TgvalidatordGetExchangeWithdrawalFeeReply reply =
                    exchangeApi.exchangeServiceGetExchangeWithdrawalFee(exchangeId, toAddressId, amount);
            return exchangeMapper.fromWithdrawalFeeReply(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Exports all exchanges to a specified format.
     *
     * @param format the export format (e.g., "csv", "json")
     * @return the exported data as a string
     * @throws ApiException if the API call fails
     */
    public String exportExchanges(final String format) throws ApiException {
        try {
            TgvalidatordExportExchangesReply reply = exchangeApi.exchangeServiceExportExchanges(format);
            return reply.getResult();
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
