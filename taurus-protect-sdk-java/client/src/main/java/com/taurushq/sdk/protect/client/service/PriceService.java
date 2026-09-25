package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.cache.RulesContainerCache;
import com.taurushq.sdk.protect.client.helper.PriceVerifier;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.PriceMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ConversionResult;
import com.taurushq.sdk.protect.client.model.Pagination;
import com.taurushq.sdk.protect.client.model.Price;
import com.taurushq.sdk.protect.client.model.PriceHistoryPoint;
import com.taurushq.sdk.protect.client.model.PriceResult;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.PricesApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordConversionReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordConversionValue;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCurrencyFromFilter;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCurrencyFromToFilter;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCurrencyPrice;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCurrencyToFilter;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetPricesHistoryReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordPricesHistoryPoint;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordQueryPricesV2Reply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordQueryPricesV2Request;

import java.util.ArrayList;
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
 * // Get current prices, one page at a time
 * PriceResult prices = client.getPriceService().getPrices(null, null, null, null, 20, null);
 * while (prices.getPage().hasMore()) {
 *     prices = client.getPriceService().getPrices(null, null, null, null, 20, prices.getPage().getNextCursor());
 * }
 *
 * // Get price history for a currency pair
 * List<PriceHistoryPoint> history = client.getPriceService()
 *     .getPriceHistory("ETH", "USD", 30);
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
     * Cache supplying the PRICEUPDATER keys from a verified rules container.
     */
    private final RulesContainerCache rulesContainerCache;

    /**
     * Instantiates a new Price service.
     * <p>
     * Price signature verification is mandatory: rate and decimals feed amount
     * conversion, so an unverified price is a wrong number a caller acts on. Whether
     * prices must be signed is decided by the SuperAdmin-verified rules container.
     *
     * @param openApiClient       the open api client
     * @param apiExceptionMapper  the api exception mapper
     * @param rulesContainerCache the cache supplying the PRICEUPDATER keys, required
     */
    public PriceService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper,
                        final RulesContainerCache rulesContainerCache) {

        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");
        checkNotNull(rulesContainerCache,
                "rulesContainerCache cannot be null - price signature verification is mandatory");

        this.rulesContainerCache = rulesContainerCache;
        this.apiExceptionMapper = apiExceptionMapper;
        this.pricesApi = new PricesApi(openApiClient);
    }


    /**
     * Gets the first page of prices, with the default page size.
     *
     * @return the verified prices and their page
     * @throws ApiException the api exception
     */
    public PriceResult getPrices() throws ApiException {
        return getPrices(null, null, null, null, null, null);
    }


    /**
     * Gets a page of prices, each one's signature verified.
     * <p>
     * The currency filter: {@code fromCurrencyId} alone lists that currency's prices;
     * {@code toCurrencyIds} alone lists the prices into those currencies; both together list
     * the prices of {@code fromCurrencyId} into those currencies; neither lists every price.
     *
     * @param fromCurrencyId the source currency id, or null
     * @param toCurrencyIds  the target currency ids, or null
     * @param onlyPrimary    true to list only the primary price of each pair, or null
     * @param sortOrder      "ASC" or "DESC", or null for the server default
     * @param pageSize       the page size, null or 0 for the default
     * @param cursor         a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the verified prices and their page
     * @throws ApiException             the api exception
     * @throws IllegalArgumentException if the page size is out of range
     */
    public PriceResult getPrices(final String fromCurrencyId, final List<String> toCurrencyIds,
                                 final Boolean onlyPrimary, final String sortOrder,
                                 final Integer pageSize, final String cursor) throws ApiException {
        final CursorRequest page = CursorRequest.of(Pagination.page(pageSize, cursor));

        TgvalidatordQueryPricesV2Request request = new TgvalidatordQueryPricesV2Request();
        request.setCursor(page.toDTO());
        request.setOnlyPrimary(onlyPrimary);
        request.setSortOrder(sortOrder);
        boolean hasFrom = !Strings.isNullOrEmpty(fromCurrencyId);
        boolean hasTo = toCurrencyIds != null && !toCurrencyIds.isEmpty();
        if (hasFrom && hasTo) {
            TgvalidatordCurrencyFromToFilter fromTo = new TgvalidatordCurrencyFromToFilter();
            fromTo.setCurrencyFromId(fromCurrencyId);
            fromTo.setCurrencyToIds(new ArrayList<>(toCurrencyIds));
            request.setFromTo(fromTo);
        } else if (hasFrom) {
            TgvalidatordCurrencyFromFilter from = new TgvalidatordCurrencyFromFilter();
            from.setCurrencyFromId(fromCurrencyId);
            request.setFrom(from);
        } else if (hasTo) {
            TgvalidatordCurrencyToFilter to = new TgvalidatordCurrencyToFilter();
            to.setCurrencyToIds(new ArrayList<>(toCurrencyIds));
            request.setTo(to);
        }

        try {
            TgvalidatordQueryPricesV2Reply reply = pricesApi.priceServiceQueryPricesV2(request);

            List<TgvalidatordCurrencyPrice> rows = reply.getResult();
            PriceResult result = new PriceResult();
            result.setBaseCurrency(reply.getBaseCurrency());
            result.setPrices(rows == null || rows.isEmpty()
                    ? Collections.emptyList() : verifiedPrices(PriceMapper.INSTANCE.fromDTO(rows)));
            return page.complete(result, PagedOperation.PRICES, reply.getCursor(), null);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets the price history of a currency pair, newest first. This endpoint cannot page.
     *
     * @param base  the base currency
     * @param quote the quote currency
     * @param limit how many daily points, 0 for the default ({@link Pagination#DEFAULT_PAGE_SIZE}),
     *              at most {@link Pagination#MAX_PRICE_HISTORY_LIMIT}
     * @return the list of price history points
     * @throws ApiException             the api exception
     * @throws IllegalArgumentException if base or quote is empty, or limit is out of range
     */
    public List<PriceHistoryPoint> getPriceHistory(final String base, final String quote, final int limit) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(base), "base cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(quote), "quote cannot be null or empty");
        final int size = PagedOperation.PRICE_HISTORY.resolveSize("limit", limit);

        try {
            TgvalidatordGetPricesHistoryReply reply = pricesApi.priceServiceGetPricesHistory(
                    base,
                    quote,
                    String.valueOf(size)
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

    /**
     * Verifies every price against the container's PRICEUPDATER keys.
     *
     * @param prices the mapped prices
     * @return the same prices, once verified
     * @throws ApiException if the rules container cannot be fetched
     */
    private List<Price> verifiedPrices(final List<Price> prices) throws ApiException {
        if (prices == null || prices.isEmpty()) {
            return prices;
        }
        PriceVerifier.verifyPrices(prices, rulesContainerCache.getDecodedRulesContainer());
        return prices;
    }
}
