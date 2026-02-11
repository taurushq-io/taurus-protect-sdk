package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.StatisticsMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.PortfolioStatistics;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.StatisticsApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetPortfolioStatisticsReply;

import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for retrieving portfolio statistics in the Taurus Protect system.
 * <p>
 * This service provides access to aggregated portfolio statistics including
 * total balances, address counts, and wallet counts.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get portfolio statistics
 * PortfolioStatistics stats = client.getStatisticsService().getPortfolioStatistics();
 * System.out.println("Total wallets: " + stats.getWalletsCount());
 * System.out.println("Total addresses: " + stats.getAddressesCount());
 * System.out.println("Total balance (base currency): " + stats.getTotalBalanceBaseCurrency());
 * }</pre>
 *
 * @see PortfolioStatistics
 */
public class StatisticsService {

    private final StatisticsApi statisticsApi;
    private final ApiExceptionMapper apiExceptionMapper;
    private final StatisticsMapper statisticsMapper;

    /**
     * Instantiates a new Statistics service.
     *
     * @param openApiClient      the OpenAPI client
     * @param apiExceptionMapper the API exception mapper
     */
    public StatisticsService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.statisticsApi = new StatisticsApi(openApiClient);
        this.statisticsMapper = StatisticsMapper.INSTANCE;
    }

    /**
     * Retrieves aggregated portfolio statistics.
     * <p>
     * Returns summary statistics including total balance, number of wallets,
     * and number of addresses across the entire portfolio.
     *
     * @return the portfolio statistics
     * @throws ApiException if the API call fails
     */
    public PortfolioStatistics getPortfolioStatistics() throws ApiException {
        try {
            TgvalidatordGetPortfolioStatisticsReply reply = statisticsApi.statisticsServiceGetPortfolioStatistics();
            return statisticsMapper.fromDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
