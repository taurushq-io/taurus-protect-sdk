package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.model.ApiResponseCursor;
import com.taurushq.sdk.protect.client.model.CursorPage;
import com.taurushq.sdk.protect.client.model.OffsetPagination;
import com.taurushq.sdk.protect.client.model.OffsetRule;
import com.taurushq.sdk.protect.client.model.Pagination;

/**
 * Every paged operation this SDK calls, with the pagination rule its reply follows.
 * <p>
 * The services build every pagination value through these constants, so this table is the
 * one place a rule is decided; {@code PaginationVectorsTest} holds it equal to the shared
 * {@code pagination-vectors.json#operations}.
 */
enum PagedOperation {

    WALLETS("WalletService_GetWalletsV2", OffsetRule.REPLY_OFFSET),
    ADDRESSES("WalletService_GetAddresses", OffsetRule.REPLY_OFFSET),
    TRANSACTIONS("TransactionService_GetTransactions", OffsetRule.PLUS_ROWS),
    FEE_PAYERS("FeePayerService_GetFeePayers", OffsetRule.PLUS_ROWS),
    ACTIONS("ActionService_GetActions", OffsetRule.PLUS_ROWS),
    USERS("UserService_GetUsers", OffsetRule.PLUS_MIN_ROWS_LIMIT),
    GROUPS("UserService_GetGroups", OffsetRule.PLUS_MIN_ROWS_LIMIT),
    WHITELISTED_ADDRESSES("WhitelistService_GetWhitelistedAddresses", OffsetRule.PLUS_SERVER_ROWS),
    WHITELISTED_ADDRESSES_FOR_APPROVAL("WhitelistService_GetWhitelistedAddressesForApproval",
            OffsetRule.PLUS_SERVER_ROWS),
    WHITELISTED_CONTRACTS("WhitelistService_GetWhitelistedContracts", OffsetRule.PLUS_LIMIT),
    WHITELISTED_CONTRACTS_FOR_APPROVAL("WhitelistService_GetWhitelistedContractsForApproval",
            OffsetRule.PLUS_LIMIT),

    REQUESTS("RequestService_GetRequestsV2", Kind.CURSOR, false),
    REQUESTS_FOR_APPROVAL("RequestService_GetRequestsForApprovalV2", Kind.CURSOR, false),
    CHANGES("ChangeService_GetChanges", Kind.CURSOR, false),
    CHANGES_FOR_APPROVAL("ChangeService_GetChangesForApproval", Kind.CURSOR, false),
    AUDIT_TRAILS("AuditService_GetAuditTrails", Kind.CURSOR, false),
    BUSINESS_RULES("RuleService_GetBusinessRulesV2", Kind.CURSOR, false),
    STAKE_ACCOUNTS("StakingService_GetStakeAccounts", Kind.CURSOR, false),
    SHARED_ADDRESSES("TaurusNetworkService_GetSharedAddresses", Kind.CURSOR, false),
    LENDING_AGREEMENTS("TaurusNetworkService_GetLendingAgreements", Kind.CURSOR, false),
    LENDING_OFFERS("TaurusNetworkService_GetLendingOffers", Kind.CURSOR, false),
    PLEDGES("TaurusNetworkService_GetPledges", Kind.CURSOR, false),
    PLEDGE_WITHDRAWALS("TaurusNetworkService_GetPledgesWithdrawals", Kind.CURSOR, false),
    SETTLEMENTS("TaurusNetworkService_GetSettlements", Kind.CURSOR, false),
    BALANCES("WalletService_GetBalances", Kind.CURSOR, true),
    ASSET_ADDRESSES("WalletService_GetAssetAddresses", Kind.CURSOR, true),
    ASSET_WALLETS("WalletService_GetAssetWallets", Kind.CURSOR, true),
    NFT_COLLECTION_BALANCES("WalletService_GetNFTCollectionBalances", Kind.CURSOR, false),
    RESERVATIONS("WalletService_GetReservations", Kind.CURSOR, false),
    FIAT_PROVIDER_ACCOUNTS("FiatProviderService_GetFiatProviderAccounts", Kind.CURSOR, false),
    FIAT_PROVIDER_COUNTERPARTY_ACCOUNTS("FiatProviderService_GetFiatProviderCounterpartyAccounts",
            Kind.CURSOR, false),
    FIAT_PROVIDER_OPERATIONS("FiatProviderService_GetFiatProviderOperations", Kind.CURSOR, false),
    FIAT_PROVIDER_ENTITIES("FiatProviderService_GetFiatProviderEntities", Kind.CURSOR, false),
    WEBHOOKS("WebhookService_GetWebhooks", Kind.CURSOR, false),
    WEBHOOK_CALLS("WebhookService_GetWebhookCalls", Kind.CURSOR, false),
    PRICES("PriceService_QueryPricesV2", Kind.CURSOR, false),
    ASSETS_V2("AssetServiceV2_QueryAssetsV2", Kind.CURSOR, false),
    ASSET_ADDRESSES_V2("AssetServiceV2_QueryAssetAddressesV2", Kind.CURSOR, false),
    ASSET_OPERATIONS_V2("AssetServiceV2_ListAssetOperationsV2", Kind.CURSOR, false),
    EARN_REWARDS("EarnService_GetRewards", Kind.CURSOR, false),

    WALLET_TOKENS("WalletService_GetWalletTokens", Kind.TOKEN, true),
    RULES_HISTORY("RuleService_GetRulesHistory", Kind.TOKEN, true),

    PRICE_HISTORY("PriceService_GetPricesHistory", Kind.LIMIT_ONLY, null, false,
            Pagination.MAX_PRICE_HISTORY_LIMIT),
    // The export cannot page: the server ignores the offset, so it has no SDK maximum.
    TRANSACTION_EXPORT("TransactionService_ExportTransactions", Kind.LIMIT_ONLY, null, true, null);

    /**
     * How an operation pages.
     */
    enum Kind {
        /** limit + offset, next offset by an {@link OffsetRule}. */
        OFFSET,
        /** page size + a reply cursor with currentPage/hasNext. */
        CURSOR,
        /** limit + a reply token that is the next cursor. */
        TOKEN,
        /** a limit only: the endpoint cannot page. */
        LIMIT_ONLY
    }

    private final String operationId;
    private final Kind kind;
    private final OffsetRule offsetRule;
    private final boolean hasTotal;
    private final Integer maxSize;

    PagedOperation(final String operationId, final OffsetRule offsetRule) {
        this(operationId, Kind.OFFSET, offsetRule, true, Pagination.MAX_PAGE_SIZE);
    }

    PagedOperation(final String operationId, final Kind kind, final boolean hasTotal) {
        this(operationId, kind, null, hasTotal, Pagination.MAX_PAGE_SIZE);
    }

    PagedOperation(final String operationId, final Kind kind, final OffsetRule offsetRule,
                   final boolean hasTotal, final Integer maxSize) {
        this.operationId = operationId;
        this.kind = kind;
        this.offsetRule = offsetRule;
        this.hasTotal = hasTotal;
        this.maxSize = maxSize;
    }

    /**
     * Returns the swagger operation id.
     *
     * @return the operation id
     */
    String operationId() {
        return operationId;
    }

    /**
     * Returns how the operation pages.
     *
     * @return the kind
     */
    Kind kind() {
        return kind;
    }

    /**
     * Returns the offset rule, for an offset operation.
     *
     * @return the rule, or null for other kinds
     */
    OffsetRule offsetRule() {
        return offsetRule;
    }

    /**
     * Reports whether the reply carries a count.
     *
     * @return true when the reply reports a total
     */
    boolean hasTotal() {
        return hasTotal;
    }

    /**
     * Returns the largest page size or limit the SDK accepts.
     *
     * @return the maximum, or null when there is none
     */
    Integer maxSize() {
        return maxSize;
    }

    /**
     * Resolves the caller's page size or limit for this operation.
     *
     * @param option the option name used in the error message
     * @param value  the caller's value, null or 0 for the default
     * @return the value to send
     * @throws IllegalArgumentException if it is negative or above this operation's maximum
     */
    int resolveSize(final String option, final Integer value) {
        return Pagination.resolveSize(option, value, maxSize);
    }

    /**
     * Builds an offset list's pagination.
     *
     * @param limit        the page size sent
     * @param offset       the offset sent
     * @param serverRows   the rows the server returned
     * @param excludedRows the rows the SDK withheld
     * @param totalItems   the reply's total as received
     * @param replyOffset  the reply's offset as received
     * @return the pagination
     */
    OffsetPagination offsetPage(final int limit, final long offset, final int serverRows,
                                final int excludedRows, final String totalItems,
                                final String replyOffset) {
        checkKind(Kind.OFFSET);
        return OffsetPagination.of(offsetRule, limit, offset, serverRows, excludedRows,
                totalItems, replyOffset);
    }

    /**
     * Builds a cursor list's page.
     *
     * @param pageSize    the page size sent
     * @param replyCursor the reply's cursor
     * @param totalItems  the reply's total as received, ignored when the reply has none
     * @return the page
     */
    CursorPage cursorPage(final int pageSize, final ApiResponseCursor replyCursor,
                          final String totalItems) {
        checkKind(Kind.CURSOR);
        return CursorPage.fromCursor(pageSize, replyCursor, totalItems, hasTotal);
    }

    /**
     * Builds a token list's page.
     *
     * @param pageSize   the page size sent
     * @param replyToken the reply's token as text
     * @param totalItems the reply's total as received
     * @return the page
     */
    CursorPage tokenPage(final int pageSize, final String replyToken, final String totalItems) {
        checkKind(Kind.TOKEN);
        return CursorPage.fromToken(pageSize, replyToken, totalItems, hasTotal);
    }

    private void checkKind(final Kind expected) {
        if (kind != expected) {
            throw new IllegalStateException(operationId + " is a " + kind + " operation, not " + expected);
        }
    }
}
