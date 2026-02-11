package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.FiatMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.FiatProvider;
import com.taurushq.sdk.protect.client.model.FiatProviderAccount;
import com.taurushq.sdk.protect.client.model.FiatProviderAccountResult;
import com.taurushq.sdk.protect.client.model.FiatProviderCounterpartyAccount;
import com.taurushq.sdk.protect.client.model.FiatProviderCounterpartyAccountResult;
import com.taurushq.sdk.protect.client.model.FiatProviderOperation;
import com.taurushq.sdk.protect.client.model.FiatProviderOperationResult;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.FiatApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFiatProviderAccountReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFiatProviderAccountsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFiatProviderCounterpartyAccountReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFiatProviderCounterpartyAccountsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFiatProviderOperationReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFiatProviderOperationsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFiatProvidersReply;

import java.util.List;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for managing fiat currency operations in the Taurus Protect system.
 * <p>
 * This service provides access to fiat provider accounts, counterparty accounts,
 * and operations for fiat currency management.
 * <p>
 * Example usage:
 * <pre>{@code
 * // List all fiat providers
 * List<FiatProvider> providers = client.getFiatService().getFiatProviders();
 *
 * // Get a specific account
 * FiatProviderAccount account = client.getFiatService()
 *     .getFiatProviderAccount("account-123");
 *
 * // List accounts with pagination
 * FiatProviderAccountResult result = client.getFiatService()
 *     .getFiatProviderAccounts(null, null, null, null, null);
 * }</pre>
 *
 * @see FiatProvider
 * @see FiatProviderAccount
 */
public class FiatService {

    /**
     * The underlying OpenAPI client for fiat operations.
     */
    private final FiatApi fiatApi;

    /**
     * Mapper for converting OpenAPI exceptions to SDK exceptions.
     */
    private final ApiExceptionMapper apiExceptionMapper;

    /**
     * Mapper for converting fiat DTOs to domain models.
     */
    private final FiatMapper mapper;

    /**
     * Instantiates a new Fiat service.
     *
     * @param openApiClient      the OpenAPI client
     * @param apiExceptionMapper the API exception mapper
     */
    public FiatService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.fiatApi = new FiatApi(openApiClient);
        this.mapper = FiatMapper.INSTANCE;
    }

    /**
     * Retrieves all configured fiat providers.
     *
     * @return the list of fiat providers
     * @throws ApiException if the API call fails
     */
    public List<FiatProvider> getFiatProviders() throws ApiException {
        try {
            TgvalidatordGetFiatProvidersReply reply = fiatApi.fiatProviderServiceGetFiatProviders();
            return mapper.fromProviderDTOList(reply.getFiatProviders());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves a fiat provider account by ID.
     *
     * @param id the account ID
     * @return the fiat provider account
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if id is null or empty
     */
    public FiatProviderAccount getFiatProviderAccount(final String id) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(id), "id cannot be null or empty");

        try {
            TgvalidatordGetFiatProviderAccountReply reply = fiatApi.fiatProviderServiceGetFiatProviderAccount(id);
            return mapper.fromAccountDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves fiat provider accounts with optional filtering.
     *
     * @param provider    filter by provider (optional)
     * @param label       filter by label (optional)
     * @param accountType filter by account type (optional)
     * @param sortOrder   sort order for results (optional, "ASC" or "DESC")
     * @param cursor      pagination cursor (optional, null for first page)
     * @return a paginated result containing fiat provider accounts
     * @throws ApiException if the API call fails
     */
    public FiatProviderAccountResult getFiatProviderAccounts(final String provider, final String label,
                                                              final String accountType, final String sortOrder,
                                                              final ApiRequestCursor cursor)
            throws ApiException {

        String cursorCurrentPage = null;
        String cursorPageRequest = null;
        String cursorPageSize = null;

        if (cursor != null) {
            cursorCurrentPage = cursor.getCurrentPage();
            cursorPageRequest = cursor.getPageRequest() != null ? cursor.getPageRequest().name() : null;
            cursorPageSize = String.valueOf(cursor.getPageSize());
        }

        try {
            TgvalidatordGetFiatProviderAccountsReply reply = fiatApi.fiatProviderServiceGetFiatProviderAccounts(
                    provider,
                    label,
                    sortOrder,
                    cursorCurrentPage,
                    cursorPageRequest,
                    cursorPageSize,
                    accountType
            );
            return mapper.fromAccountsReply(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves a fiat provider counterparty account by ID.
     *
     * @param id the counterparty account ID
     * @return the fiat provider counterparty account
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if id is null or empty
     */
    public FiatProviderCounterpartyAccount getFiatProviderCounterpartyAccount(final String id) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(id), "id cannot be null or empty");

        try {
            TgvalidatordGetFiatProviderCounterpartyAccountReply reply =
                    fiatApi.fiatProviderServiceGetFiatProviderCounterpartyAccount(id);
            return mapper.fromCounterpartyAccountDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves fiat provider counterparty accounts with optional filtering.
     *
     * @param provider       filter by provider (optional)
     * @param label          filter by label (optional)
     * @param counterpartyId filter by counterparty ID (optional)
     * @param sortOrder      sort order for results (optional, "ASC" or "DESC")
     * @param cursor         pagination cursor (optional, null for first page)
     * @return a paginated result containing fiat provider counterparty accounts
     * @throws ApiException if the API call fails
     */
    public FiatProviderCounterpartyAccountResult getFiatProviderCounterpartyAccounts(
            final String provider, final String label, final String counterpartyId,
            final String sortOrder, final ApiRequestCursor cursor)
            throws ApiException {

        String cursorCurrentPage = null;
        String cursorPageRequest = null;
        String cursorPageSize = null;

        if (cursor != null) {
            cursorCurrentPage = cursor.getCurrentPage();
            cursorPageRequest = cursor.getPageRequest() != null ? cursor.getPageRequest().name() : null;
            cursorPageSize = String.valueOf(cursor.getPageSize());
        }

        try {
            TgvalidatordGetFiatProviderCounterpartyAccountsReply reply =
                    fiatApi.fiatProviderServiceGetFiatProviderCounterpartyAccounts(
                            provider,
                            label,
                            counterpartyId,
                            sortOrder,
                            cursorCurrentPage,
                            cursorPageRequest,
                            cursorPageSize
                    );
            return mapper.fromCounterpartyAccountsReply(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves a fiat provider operation by ID.
     *
     * @param id the operation ID
     * @return the fiat provider operation
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if id is null or empty
     */
    public FiatProviderOperation getFiatProviderOperation(final String id) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(id), "id cannot be null or empty");

        try {
            TgvalidatordGetFiatProviderOperationReply reply =
                    fiatApi.fiatProviderServiceGetFiatProviderOperation(id);
            return mapper.fromOperationDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves fiat provider operations with optional filtering.
     *
     * @param provider  filter by provider (optional)
     * @param label     filter by label (optional)
     * @param sortOrder sort order for results (optional, "ASC" or "DESC")
     * @param cursor    pagination cursor (optional, null for first page)
     * @return a paginated result containing fiat provider operations
     * @throws ApiException if the API call fails
     */
    public FiatProviderOperationResult getFiatProviderOperations(final String provider, final String label,
                                                                  final String sortOrder,
                                                                  final ApiRequestCursor cursor)
            throws ApiException {

        String cursorCurrentPage = null;
        String cursorPageRequest = null;
        String cursorPageSize = null;

        if (cursor != null) {
            cursorCurrentPage = cursor.getCurrentPage();
            cursorPageRequest = cursor.getPageRequest() != null ? cursor.getPageRequest().name() : null;
            cursorPageSize = String.valueOf(cursor.getPageSize());
        }

        try {
            TgvalidatordGetFiatProviderOperationsReply reply = fiatApi.fiatProviderServiceGetFiatProviderOperations(
                    provider,
                    label,
                    sortOrder,
                    cursorCurrentPage,
                    cursorPageRequest,
                    cursorPageSize
            );
            return mapper.fromOperationsReply(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
