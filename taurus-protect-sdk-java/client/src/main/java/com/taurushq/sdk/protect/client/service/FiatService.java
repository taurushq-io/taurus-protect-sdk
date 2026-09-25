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
import com.taurushq.sdk.protect.client.model.FiatProviderEntityResult;
import com.taurushq.sdk.protect.client.model.FiatProviderOperation;
import com.taurushq.sdk.protect.client.model.FiatProviderOperationResult;
import com.taurushq.sdk.protect.client.model.Pagination;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.FiatApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFiatProviderAccountReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFiatProviderAccountsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFiatProviderCounterpartyAccountReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFiatProviderCounterpartyAccountsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFiatProviderEntitiesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFiatProviderOperationReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFiatProviderOperationsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFiatProvidersReply;

import java.util.Collections;
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
 * // List a provider's accounts, first page
 * FiatProviderAccountResult result = client.getFiatService()
 *     .getFiatProviderAccounts("provider", "label", null, null, 20, null);
 * // next page: pass result.getPage().getNextCursor() while result.getPage().hasMore()
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
     * Retrieves a page of a provider's accounts.
     *
     * @param provider    the provider, required
     * @param label       the provider configuration label, required
     * @param accountType filter by account type (optional)
     * @param sortOrder   sort order for results (optional, "ASC" or "DESC")
     * @param pageSize    the page size, null or 0 for the default
     * @param cursor      a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the accounts and their page
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if provider or label is empty, or the page size is out of range
     */
    public FiatProviderAccountResult getFiatProviderAccounts(final String provider, final String label,
                                                              final String accountType, final String sortOrder,
                                                              final Integer pageSize, final String cursor)
            throws ApiException {
        return getFiatProviderAccounts(provider, label, accountType, sortOrder, Pagination.page(pageSize, cursor));
    }

    /**
     * Retrieves a page of a provider's accounts, with a low-level request cursor.
     *
     * @param provider    the provider, required
     * @param label       the provider configuration label, required
     * @param accountType filter by account type (optional)
     * @param sortOrder   sort order for results (optional, "ASC" or "DESC")
     * @param cursor      the request cursor, null for the first page with the default size
     * @return the accounts and their page
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if provider or label is empty
     */
    public FiatProviderAccountResult getFiatProviderAccounts(final String provider, final String label,
                                                              final String accountType, final String sortOrder,
                                                              final ApiRequestCursor cursor)
            throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(provider), "provider cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(label), "label cannot be null or empty");
        final CursorRequest page = CursorRequest.of(cursor);

        try {
            TgvalidatordGetFiatProviderAccountsReply reply = fiatApi.fiatProviderServiceGetFiatProviderAccounts(
                    provider,
                    label,
                    sortOrder,
                    page.currentPage(),
                    page.pageRequest(),
                    page.pageSizeParam(),
                    accountType
            );
            return page.complete(mapper.fromAccountsReply(reply), PagedOperation.FIAT_PROVIDER_ACCOUNTS);
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
     * Retrieves a page of a provider's counterparty accounts.
     *
     * @param provider       the provider, required
     * @param label          the provider configuration label, required
     * @param counterpartyId filter by counterparty ID (optional)
     * @param sortOrder      sort order for results (optional, "ASC" or "DESC")
     * @param pageSize       the page size, null or 0 for the default
     * @param cursor         a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the counterparty accounts and their page
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if provider or label is empty, or the page size is out of range
     */
    public FiatProviderCounterpartyAccountResult getFiatProviderCounterpartyAccounts(
            final String provider, final String label, final String counterpartyId,
            final String sortOrder, final Integer pageSize, final String cursor)
            throws ApiException {
        return getFiatProviderCounterpartyAccounts(provider, label, counterpartyId, sortOrder,
                Pagination.page(pageSize, cursor));
    }

    /**
     * Retrieves a page of a provider's counterparty accounts, with a low-level request cursor.
     *
     * @param provider       the provider, required
     * @param label          the provider configuration label, required
     * @param counterpartyId filter by counterparty ID (optional)
     * @param sortOrder      sort order for results (optional, "ASC" or "DESC")
     * @param cursor         the request cursor, null for the first page with the default size
     * @return the counterparty accounts and their page
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if provider or label is empty
     */
    public FiatProviderCounterpartyAccountResult getFiatProviderCounterpartyAccounts(
            final String provider, final String label, final String counterpartyId,
            final String sortOrder, final ApiRequestCursor cursor)
            throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(provider), "provider cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(label), "label cannot be null or empty");
        final CursorRequest page = CursorRequest.of(cursor);

        try {
            TgvalidatordGetFiatProviderCounterpartyAccountsReply reply =
                    fiatApi.fiatProviderServiceGetFiatProviderCounterpartyAccounts(
                            provider,
                            label,
                            counterpartyId,
                            sortOrder,
                            page.currentPage(),
                            page.pageRequest(),
                            page.pageSizeParam()
                    );
            return page.complete(mapper.fromCounterpartyAccountsReply(reply),
                    PagedOperation.FIAT_PROVIDER_COUNTERPARTY_ACCOUNTS);
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
     * Retrieves a page of fiat provider operations.
     *
     * @param provider  filter by provider (optional)
     * @param label     filter by label (optional)
     * @param sortOrder sort order for results (optional, "ASC" or "DESC")
     * @param pageSize  the page size, null or 0 for the default
     * @param cursor    a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the operations and their page
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if the page size is out of range
     */
    public FiatProviderOperationResult getFiatProviderOperations(final String provider, final String label,
                                                                  final String sortOrder, final Integer pageSize,
                                                                  final String cursor)
            throws ApiException {
        return getFiatProviderOperations(provider, label, sortOrder, Pagination.page(pageSize, cursor));
    }

    /**
     * Retrieves a page of fiat provider operations, with a low-level request cursor.
     *
     * @param provider  filter by provider (optional)
     * @param label     filter by label (optional)
     * @param sortOrder sort order for results (optional, "ASC" or "DESC")
     * @param cursor    the request cursor, null for the first page with the default size
     * @return the operations and their page
     * @throws ApiException if the API call fails
     */
    public FiatProviderOperationResult getFiatProviderOperations(final String provider, final String label,
                                                                  final String sortOrder,
                                                                  final ApiRequestCursor cursor)
            throws ApiException {
        final CursorRequest page = CursorRequest.of(cursor);

        try {
            TgvalidatordGetFiatProviderOperationsReply reply = fiatApi.fiatProviderServiceGetFiatProviderOperations(
                    provider,
                    label,
                    sortOrder,
                    page.currentPage(),
                    page.pageRequest(),
                    page.pageSizeParam()
            );
            return page.complete(mapper.fromOperationsReply(reply), PagedOperation.FIAT_PROVIDER_OPERATIONS);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves a page of the entities registered with fiat providers.
     *
     * @param provider  filter by provider (optional)
     * @param label     filter by label (optional)
     * @param sortOrder sort order for results (optional, "ASC" or "DESC")
     * @param pageSize  the page size, null or 0 for the default
     * @param cursor    a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the entities and their page
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if the page size is out of range
     */
    public FiatProviderEntityResult listFiatProviderEntities(final String provider, final String label,
                                                             final String sortOrder, final Integer pageSize,
                                                             final String cursor) throws ApiException {
        final CursorRequest page = CursorRequest.of(Pagination.page(pageSize, cursor));

        try {
            TgvalidatordGetFiatProviderEntitiesReply reply = fiatApi.fiatProviderServiceGetFiatProviderEntities(
                    provider,
                    label,
                    sortOrder,
                    page.currentPage(),
                    page.pageRequest(),
                    page.pageSizeParam()
            );
            FiatProviderEntityResult result = new FiatProviderEntityResult();
            result.setEntities(reply.getResult() == null
                    ? Collections.emptyList() : mapper.fromEntityDTOList(reply.getResult()));
            return page.complete(result, PagedOperation.FIAT_PROVIDER_ENTITIES, reply.getCursor(), null);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
