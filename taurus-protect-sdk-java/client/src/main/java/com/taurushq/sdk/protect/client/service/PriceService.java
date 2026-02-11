package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.PriceMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ConversionResult;
import com.taurushq.sdk.protect.client.model.Price;
import com.taurushq.sdk.protect.client.model.PriceHistoryPoint;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.PricesApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordConversionReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordConversionValue;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCurrencyPrice;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetPricesHistoryReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetPricesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordPricesHistoryPoint;

import java.util.Collections;
import java.util.List;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for retrieving cryptocurrency prices and performing conversions.
 * <p>
 * This service provides operations for querying current and historical prices
 * of supported cryptocurrencies, as well as converting amounts between currencies.
 * Prices are provided against the configured base currency (typically USD).
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get all current prices
 * List<Price> prices = client.getPriceService().getPrices();
 *
 * // Get price history for a currency pair
 * List<PriceHistoryPoint> history = client.getPriceService()
 *     .getPriceHistory("ETH", "USD", 100);
 *
 * // Convert an amount to target currencies
 * List<ConversionResult> converted = client.getPriceService()
 *     .convert("ETH", "1000000000000000000", Arrays.asList("USD", "BTC"));
 * }</pre>
 *
 * @see Price
 * @see PriceHistoryPoint
 * @see ConversionResult
 * @see CurrencyService
 */
public class PriceService {

    /**
     * The underlying OpenAPI client for price operations.
     */
    private final PricesApi pricesApi;

    /**
     * Mapper for converting OpenAPI exceptions to SDK exceptions.
     */
    private final ApiExceptionMapper apiExceptionMapper;

    /**
     * Instantiates a new Price service.
     *
     * @param openApiClient      the open api client
     * @param apiExceptionMapper the api exception mapper
     */
    public PriceService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {

        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.pricesApi = new PricesApi(openApiClient);
    }


    /**
     * Gets all prices.
     *
     * @return the list of prices
     * @throws ApiException the api exception
     */
    public List<Price> getPrices() throws ApiException {
        try {
            TgvalidatordGetPricesReply reply = pricesApi.priceServiceGetPrices();

            List<TgvalidatordCurrencyPrice> result = reply.getResult();
            if (result == null) {
                return Collections.emptyList();
            }
            return PriceMapper.INSTANCE.fromDTO(result);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets price history.
     *
     * @param base  the base currency
     * @param quote the quote currency
     * @param limit the limit
     * @return the list of price history points
     * @throws ApiException the api exception
     */
    public List<PriceHistoryPoint> getPriceHistory(final String base, final String quote, final int limit) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(base), "base cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(quote), "quote cannot be null or empty");
        checkArgument(limit > 0, "limit must be positive");

        try {
            TgvalidatordGetPricesHistoryReply reply = pricesApi.priceServiceGetPricesHistory(
                    base,
                    quote,
                    String.valueOf(limit)
            );

            List<TgvalidatordPricesHistoryPoint> result = reply.getResult();
            if (result == null) {
                return Collections.emptyList();
            }
            return PriceMapper.INSTANCE.fromDTOHistory(result);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Converts an amount from one currency to target currencies.
     *
     * @param currency           the source currency
     * @param amount             the amount to convert
     * @param targetCurrencyIds  the target currency ids
     * @return the list of conversion results
     * @throws ApiException the api exception
     */
    public List<ConversionResult> convert(final String currency, final String amount, final List<String> targetCurrencyIds) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(currency), "currency cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(amount), "amount cannot be null or empty");

        try {
            TgvalidatordConversionReply reply = pricesApi.priceServiceConvert(
                    currency,
                    amount,
                    null,               // symbols
                    targetCurrencyIds   // targetCurrencyIds
            );

            List<TgvalidatordConversionValue> result = reply.getResult();
            if (result == null) {
                return Collections.emptyList();
            }
            return PriceMapper.INSTANCE.fromDTOConversion(result);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
