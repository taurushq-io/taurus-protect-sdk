package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.BusinessRuleMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.BusinessRuleResult;
import com.taurushq.sdk.protect.client.model.Pagination;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.BusinessRulesApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordBusinessRule;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetBusinessRulesV2Reply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordUpdateTransactionsEnabledBusinessRuleRequest;

import java.util.Collections;
import java.util.List;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for managing business rules in the Taurus Protect system.
 * <p>
 * Business rules define automated policies that apply to wallets, addresses, or currencies.
 * They can enforce constraints like spending limits, approval requirements, or allowed
 * transaction types.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get business rules, one page at a time
 * BusinessRuleResult result = client.getBusinessRuleService().getBusinessRules(20, null);
 * while (result.getPage().hasMore()) {
 *     result = client.getBusinessRuleService().getBusinessRules(20, result.getPage().getNextCursor());
 * }
 *
 * // Get rules for a specific wallet
 * BusinessRuleResult walletRules = client.getBusinessRuleService()
 *     .getBusinessRulesByWallet(walletId, 20, null);
 *
 * // Get rules for a specific currency
 * BusinessRuleResult currencyRules = client.getBusinessRuleService()
 *     .getBusinessRulesByCurrency("ETH", 20, null);
 * }</pre>
 *
 * @see BusinessRuleResult
 * @see GovernanceRuleService
 */
public class BusinessRuleService {

    /**
     * The underlying OpenAPI client for business rule operations.
     */
    private final BusinessRulesApi businessRulesApi;

    /**
     * Mapper for converting OpenAPI exceptions to SDK exceptions.
     */
    private final ApiExceptionMapper apiExceptionMapper;

    /**
     * Instantiates a new Business rule service.
     *
     * @param openApiClient      the open api client
     * @param apiExceptionMapper the api exception mapper
     */
    public BusinessRuleService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {

        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.businessRulesApi = new BusinessRulesApi(openApiClient);
    }


    /**
     * Gets a page of business rules.
     *
     * @param pageSize the page size, null or 0 for the default
     * @param cursor   a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the business rules and their page
     * @throws ApiException             the api exception
     * @throws IllegalArgumentException if the page size is out of range
     */
    public BusinessRuleResult getBusinessRules(final Integer pageSize, final String cursor) throws ApiException {
        return getBusinessRules(Pagination.page(pageSize, cursor));
    }

    /**
     * Gets a page of business rules, with a low-level request cursor.
     *
     * @param cursor the request cursor, null for the first page with the default size
     * @return the business rules and their page
     * @throws ApiException the api exception
     */
    public BusinessRuleResult getBusinessRules(final ApiRequestCursor cursor) throws ApiException {
        return listBusinessRules(null, null, cursor);
    }


    /**
     * Gets a page of the business rules of a wallet.
     *
     * @param walletId the wallet id
     * @param pageSize the page size, null or 0 for the default
     * @param cursor   a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the business rules and their page
     * @throws ApiException             the api exception
     * @throws IllegalArgumentException if walletId or the page size is out of range
     */
    public BusinessRuleResult getBusinessRulesByWallet(final long walletId, final Integer pageSize,
                                                       final String cursor) throws ApiException {
        return getBusinessRulesByWallet(walletId, Pagination.page(pageSize, cursor));
    }

    /**
     * Gets a page of the business rules of a wallet, with a low-level request cursor.
     *
     * @param walletId the wallet id
     * @param cursor   the request cursor, null for the first page with the default size
     * @return the business rules and their page
     * @throws ApiException the api exception
     */
    public BusinessRuleResult getBusinessRulesByWallet(final long walletId, final ApiRequestCursor cursor)
            throws ApiException {
        checkArgument(walletId > 0, "walletId must be positive");
        return listBusinessRules(Collections.singletonList(String.valueOf(walletId)), null, cursor);
    }


    /**
     * Gets a page of the business rules of a currency.
     *
     * @param currencyId the currency id
     * @param pageSize   the page size, null or 0 for the default
     * @param cursor     a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the business rules and their page
     * @throws ApiException             the api exception
     * @throws IllegalArgumentException if currencyId is empty or the page size is out of range
     */
    public BusinessRuleResult getBusinessRulesByCurrency(final String currencyId, final Integer pageSize,
                                                         final String cursor) throws ApiException {
        return getBusinessRulesByCurrency(currencyId, Pagination.page(pageSize, cursor));
    }

    /**
     * Gets a page of the business rules of a currency, with a low-level request cursor.
     *
     * @param currencyId the currency id
     * @param cursor     the request cursor, null for the first page with the default size
     * @return the business rules and their page
     * @throws ApiException the api exception
     */
    public BusinessRuleResult getBusinessRulesByCurrency(final String currencyId, final ApiRequestCursor cursor)
            throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(currencyId), "currencyId cannot be null or empty");
        return listBusinessRules(null, Collections.singletonList(currencyId), cursor);
    }

    private BusinessRuleResult listBusinessRules(final List<String> walletIds, final List<String> currencyIds,
                                                 final ApiRequestCursor cursor) throws ApiException {
        final CursorRequest page = CursorRequest.of(cursor);

        try {
            TgvalidatordGetBusinessRulesV2Reply reply = businessRulesApi.ruleServiceGetBusinessRulesV2(
                    null,                               // ids
                    null,                               // ruleKeys
                    null,                               // ruleGroups
                    walletIds,                          // walletIds
                    currencyIds,                        // currencyIds
                    null,                               // addressIds
                    null,                               // level
                    page.currentPage(),                 // cursorCurrentPage
                    page.pageRequest(),                 // cursorPageRequest
                    page.pageSizeParam(),               // cursorPageSize
                    null,                               // entityType
                    null                                // entityIDs
            );

            BusinessRuleResult result = new BusinessRuleResult();
            List<TgvalidatordBusinessRule> rules = reply.getResult();
            result.setRules(rules == null ? Collections.emptyList() : BusinessRuleMapper.INSTANCE.fromDTO(rules));
            return page.complete(result, PagedOperation.BUSINESS_RULES, reply.getCursor(), null);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Enables or disables transaction processing for the tenant (the
     * transactions-enabled business rule).
     *
     * <p>Peer of Go BusinessRuleService.UpdateTransactionsEnabled, Python
     * update_transactions_enabled and TS updateTransactionsEnabled. This is the kill
     * switch tg-protect-mcpd drives, so it must exist in every SDK.
     *
     * @param enabled true to allow transactions, false to halt them
     * @throws ApiException the api exception
     */
    public void updateTransactionsEnabled(final boolean enabled) throws ApiException {
        TgvalidatordUpdateTransactionsEnabledBusinessRuleRequest body =
                new TgvalidatordUpdateTransactionsEnabledBusinessRuleRequest();
        body.setEnabled(enabled);

        try {
            businessRulesApi.ruleServiceUpdateTransactionsEnabledBusinessRule(body);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
