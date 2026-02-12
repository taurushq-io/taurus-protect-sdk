package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.CurrencyMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.Currency;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.CurrenciesApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCurrency;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetBaseCurrencyReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetCurrenciesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetCurrencyReply;

import java.util.Collections;
import java.util.List;
import java.util.stream.Collectors;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for retrieving currency information in the Taurus Protect system.
 * <p>
 * This service provides operations for querying supported cryptocurrencies and their
 * metadata including blockchain, network, decimals, and logos. Currencies can be
 * retrieved by ID, blockchain/network combination, or as a full list.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get all supported currencies
 * List<Currency> currencies = client.getCurrencyService().getCurrencies();
 *
 * // Get a specific currency by ID
 * Currency eth = client.getCurrencyService().getCurrency("ETH");
 *
 * // Get currency by blockchain and network
 * Currency mainnetEth = client.getCurrencyService()
 *     .getCurrencyByBlockchain("ETH", "mainnet");
 *
 * // Get the base currency (for fiat conversion)
 * Currency baseCurrency = client.getCurrencyService().getBaseCurrency();
 * }</pre>
 *
 * @see Currency
 * @see PriceService
 */
public class CurrencyService {

    /**
     * The underlying OpenAPI client for currency operations.
     */
    private final CurrenciesApi currenciesApi;

    /**
     * Mapper for converting OpenAPI exceptions to SDK exceptions.
     */
    private final ApiExceptionMapper apiExceptionMapper;

    /**
     * Instantiates a new Currency service.
     *
     * @param openApiClient      the open api client
     * @param apiExceptionMapper the api exception mapper
     */
    public CurrencyService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {

        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.currenciesApi = new CurrenciesApi(openApiClient);
    }


    /**
     * Gets all currencies.
     *
     * @return the list of currencies
     * @throws ApiException the api exception
     */
    public List<Currency> getCurrencies() throws ApiException {
        return getCurrencies(false, false);
    }


    /**
     * Gets all currencies with options.
     *
     * @param showDisabled include disabled currencies
     * @param includeLogo  include logo in response
     * @return the list of currencies
     * @throws ApiException the api exception
     */
    public List<Currency> getCurrencies(final boolean showDisabled, final boolean includeLogo) throws ApiException {
        try {
            TgvalidatordGetCurrenciesReply reply = currenciesApi.walletServiceGetCurrencies(showDisabled, includeLogo);

            List<TgvalidatordCurrency> result = reply.getResult();
            if (result == null) {
                return Collections.emptyList();
            }
            return result.stream()
                    .map(CurrencyMapper.INSTANCE::fromDTO)
                    .collect(Collectors.toList());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets currency by ID.
     * <p>
     * Note: The API requires blockchain and network to look up a currency.
     * This method finds the currency by ID in the full list of currencies.
     *
     * @param currencyId the currency id
     * @return the currency, or null if not found
     * @throws ApiException the api exception
     */
    public Currency getCurrency(final String currencyId) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(currencyId), "currencyId cannot be null or empty");

        // The API requires blockchain and network, so we need to find the currency
        // in the list first to get those values
        List<Currency> currencies = getCurrencies(true, false);
        return currencies.stream()
                .filter(c -> currencyId.equals(c.getId()))
                .findFirst()
                .orElse(null);
    }


    /**
     * Gets currency by blockchain and network.
     *
     * @param blockchain the blockchain
     * @param network    the network
     * @return the currency
     * @throws ApiException the api exception
     */
    public Currency getCurrencyByBlockchain(final String blockchain, final String network) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(blockchain), "blockchain cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(network), "network cannot be null or empty");

        try {
            TgvalidatordGetCurrencyReply reply = currenciesApi.walletServiceGetCurrency(
                    blockchain,     // uniqueCurrencyFilterBlockchain
                    network,        // uniqueCurrencyFilterNetwork
                    false,          // showDisabled
                    null,           // uniqueCurrencyFilterTokenContractAddress
                    null,           // uniqueCurrencyFilterTokenID
                    null,           // currencyID
                    false           // includeLogo
            );
            return CurrencyMapper.INSTANCE.fromDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets the base currency.
     *
     * @return the base currency
     * @throws ApiException the api exception
     */
    public Currency getBaseCurrency() throws ApiException {
        try {
            TgvalidatordGetBaseCurrencyReply reply = currenciesApi.walletServiceGetBaseCurrency();
            String baseCurrencyId = reply.getResult();
            if (baseCurrencyId == null) {
                return null;
            }
            return getCurrency(baseCurrencyId);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}