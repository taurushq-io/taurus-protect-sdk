package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.ScoreMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.Score;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.ScoresApi;
import com.taurushq.sdk.protect.openapi.model.ScoreServiceRefreshAddressScoreBody;
import com.taurushq.sdk.protect.openapi.model.ScoreServiceRefreshWLAScoreBody;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRefreshScoreReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordScore;

import java.util.Collections;
import java.util.List;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for managing blockchain analytics scores in the Taurus Protect system.
 * <p>
 * Scores represent risk assessments provided by third-party blockchain analytics
 * providers (e.g., Chainalysis, Elliptic). This service allows refreshing scores
 * for internal addresses and whitelisted external addresses.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Refresh score for an internal address
 * List<Score> scores = client.getScoreService()
 *     .refreshAddressScore(addressId, "chainalysis");
 *
 * // Refresh score for a whitelisted external address
 * List<Score> wlaScores = client.getScoreService()
 *     .refreshWhitelistedAddressScore(whitelistedAddressId, "elliptic");
 * }</pre>
 *
 * @see Score
 * @see AddressService
 * @see WhitelistedAddressService
 */
public class ScoreService {

    /**
     * The underlying OpenAPI client for score operations.
     */
    private final ScoresApi scoresApi;

    /**
     * Mapper for converting OpenAPI exceptions to SDK exceptions.
     */
    private final ApiExceptionMapper apiExceptionMapper;

    /**
     * Instantiates a new Score service.
     *
     * @param openApiClient      the open api client
     * @param apiExceptionMapper the api exception mapper
     */
    public ScoreService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {

        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.scoresApi = new ScoresApi(openApiClient);
    }


    /**
     * Refresh address score.
     *
     * @param addressId     the address id
     * @param scoreProvider the score provider
     * @return the list of refreshed scores
     * @throws ApiException the api exception
     */
    public List<Score> refreshAddressScore(final long addressId, final String scoreProvider) throws ApiException {
        checkArgument(addressId > 0, "addressId must be positive");
        checkNotNull(scoreProvider, "scoreProvider cannot be null");

        try {
            ScoreServiceRefreshAddressScoreBody body = new ScoreServiceRefreshAddressScoreBody();
            body.setScoreProvider(scoreProvider);

            TgvalidatordRefreshScoreReply reply = scoresApi.scoreServiceRefreshAddressScore(
                    String.valueOf(addressId),
                    body
            );

            List<TgvalidatordScore> result = reply.getScores();
            if (result == null) {
                return Collections.emptyList();
            }
            return ScoreMapper.INSTANCE.fromDTO(result);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Refresh whitelisted address score.
     *
     * @param addressId     the whitelisted address id
     * @param scoreProvider the score provider
     * @return the list of refreshed scores
     * @throws ApiException the api exception
     */
    public List<Score> refreshWhitelistedAddressScore(final long addressId, final String scoreProvider) throws ApiException {
        checkArgument(addressId > 0, "addressId must be positive");
        checkNotNull(scoreProvider, "scoreProvider cannot be null");

        try {
            ScoreServiceRefreshWLAScoreBody body = new ScoreServiceRefreshWLAScoreBody();
            body.setScoreProvider(scoreProvider);

            TgvalidatordRefreshScoreReply reply = scoresApi.scoreServiceRefreshWLAScore(
                    String.valueOf(addressId),
                    body
            );

            List<TgvalidatordScore> result = reply.getScores();
            if (result == null) {
                return Collections.emptyList();
            }
            return ScoreMapper.INSTANCE.fromDTO(result);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
