package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.EarnMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.EarnRewardResult;
import com.taurushq.sdk.protect.client.model.Pagination;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.EarnApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetEarnRewardsReply;

import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for the rewards addresses earn (for example Merkl token rewards).
 * <p>
 * Example usage:
 * <pre>{@code
 * EarnRewardResult rewards = client.getEarnService().listRewards(null, 20, null);
 * while (rewards.getPage().hasMore()) {
 *     rewards = client.getEarnService().listRewards(null, 20, rewards.getPage().getNextCursor());
 * }
 * }</pre>
 *
 * @see com.taurushq.sdk.protect.client.model.EarnReward
 */
public class EarnService {

    private final EarnApi earnApi;
    private final ApiExceptionMapper apiExceptionMapper;

    /**
     * Instantiates a new Earn service.
     *
     * @param openApiClient      the OpenAPI client
     * @param apiExceptionMapper the API exception mapper
     */
    public EarnService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.earnApi = new EarnApi(openApiClient);
    }

    /**
     * Lists a page of earn rewards.
     *
     * @param recipientAddressId filter by the address the rewards were credited to (optional)
     * @param pageSize           the page size, null or 0 for the default
     * @param cursor             a previous page's {@code getPage().getNextCursor()}, null for
     *                           the first page
     * @return the rewards and their page
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if the page size is out of range
     */
    public EarnRewardResult listRewards(final String recipientAddressId, final Integer pageSize,
                                        final String cursor) throws ApiException {
        final CursorRequest page = CursorRequest.of(Pagination.page(pageSize, cursor));
        try {
            TgvalidatordGetEarnRewardsReply reply = earnApi.earnServiceGetRewards(
                    recipientAddressId,
                    page.currentPage(),
                    page.pageRequest(),
                    page.pageSizeParam()
            );
            EarnRewardResult result = new EarnRewardResult();
            result.setRewards(EarnMapper.INSTANCE.fromDTOList(reply.getRewards()));
            return page.complete(result, PagedOperation.EARN_REWARDS, reply.getCursor(), null);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
