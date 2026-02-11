package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.ApiResponseCursorMapper;
import com.taurushq.sdk.protect.client.mapper.BusinessRuleMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.BusinessRuleResult;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.BusinessRulesApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordBusinessRule;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetBusinessRulesV2Reply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRequestCursor;

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
 * // Get all business rules with pagination
 * ApiRequestCursor cursor = Pagination.first(50);
 * BusinessRuleResult result = client.getBusinessRuleService().getBusinessRules(cursor);
 *
 * // Get rules for a specific wallet
 * BusinessRuleResult walletRules = client.getBusinessRuleService()
 *     .getBusinessRulesByWallet(walletId, cursor);
 *
 * // Get rules for a specific currency
 * BusinessRuleResult currencyRules = client.getBusinessRuleService()
 *     .getBusinessRulesByCurrency("ETH", cursor);
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
     * Gets all business rules with pagination.
     *
     * @param cursor the request cursor for pagination
     * @return the business rule result with list and response cursor
     * @throws ApiException the api exception
     */
    public BusinessRuleResult getBusinessRules(final ApiRequestCursor cursor) throws ApiException {
        checkNotNull(cursor, "cursor cannot be null");

        TgvalidatordRequestCursor requestCursor = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        try {
            TgvalidatordGetBusinessRulesV2Reply reply = businessRulesApi.ruleServiceGetBusinessRulesV2(
                    null,                               // ids
                    null,                               // ruleKeys
                    null,                               // ruleGroups
                    null,                               // walletIds
                    null,                               // currencyIds
                    null,                               // addressIds
                    null,                               // level
                    requestCursor.getCurrentPage(),     // cursorCurrentPage
                    requestCursor.getPageRequest(),     // cursorPageRequest
                    requestCursor.getPageSize(),        // cursorPageSize
                    null,                               // entityType
                    null                                // entityIDs
            );

            BusinessRuleResult result = new BusinessRuleResult();

            List<TgvalidatordBusinessRule> rules = reply.getResult();
            if (rules == null) {
                result.setRules(Collections.emptyList());
            } else {
                result.setRules(BusinessRuleMapper.INSTANCE.fromDTO(rules));
            }

            result.setCursor(ApiResponseCursorMapper.INSTANCE.fromDTO(reply.getCursor()));

            return result;
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets business rules by wallet id.
     *
     * @param walletId the wallet id
     * @param cursor   the request cursor for pagination
     * @return the business rule result with list and response cursor
     * @throws ApiException the api exception
     */
    public BusinessRuleResult getBusinessRulesByWallet(final long walletId, final ApiRequestCursor cursor) throws ApiException {
        checkArgument(walletId > 0, "walletId must be positive");
        checkNotNull(cursor, "cursor cannot be null");

        TgvalidatordRequestCursor requestCursor = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        try {
            TgvalidatordGetBusinessRulesV2Reply reply = businessRulesApi.ruleServiceGetBusinessRulesV2(
                    null,                               // ids
                    null,                               // ruleKeys
                    null,                               // ruleGroups
                    Collections.singletonList(String.valueOf(walletId)),  // walletIds
                    null,                               // currencyIds
                    null,                               // addressIds
                    null,                               // level
                    requestCursor.getCurrentPage(),     // cursorCurrentPage
                    requestCursor.getPageRequest(),     // cursorPageRequest
                    requestCursor.getPageSize(),        // cursorPageSize
                    null,                               // entityType
                    null                                // entityIDs
            );

            BusinessRuleResult result = new BusinessRuleResult();

            List<TgvalidatordBusinessRule> rules = reply.getResult();
            if (rules == null) {
                result.setRules(Collections.emptyList());
            } else {
                result.setRules(BusinessRuleMapper.INSTANCE.fromDTO(rules));
            }

            result.setCursor(ApiResponseCursorMapper.INSTANCE.fromDTO(reply.getCursor()));

            return result;
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets business rules by currency id.
     *
     * @param currencyId the currency id
     * @param cursor     the request cursor for pagination
     * @return the business rule result with list and response cursor
     * @throws ApiException the api exception
     */
    public BusinessRuleResult getBusinessRulesByCurrency(final String currencyId, final ApiRequestCursor cursor) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(currencyId), "currencyId cannot be null or empty");
        checkNotNull(cursor, "cursor cannot be null");

        TgvalidatordRequestCursor requestCursor = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        try {
            TgvalidatordGetBusinessRulesV2Reply reply = businessRulesApi.ruleServiceGetBusinessRulesV2(
                    null,                               // ids
                    null,                               // ruleKeys
                    null,                               // ruleGroups
                    null,                               // walletIds
                    Collections.singletonList(currencyId),  // currencyIds
                    null,                               // addressIds
                    null,                               // level
                    requestCursor.getCurrentPage(),     // cursorCurrentPage
                    requestCursor.getPageRequest(),     // cursorPageRequest
                    requestCursor.getPageSize(),        // cursorPageSize
                    null,                               // entityType
                    null                                // entityIDs
            );

            BusinessRuleResult result = new BusinessRuleResult();

            List<TgvalidatordBusinessRule> rules = reply.getResult();
            if (rules == null) {
                result.setRules(Collections.emptyList());
            } else {
                result.setRules(BusinessRuleMapper.INSTANCE.fromDTO(rules));
            }

            result.setCursor(ApiResponseCursorMapper.INSTANCE.fromDTO(reply.getCursor()));

            return result;
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
