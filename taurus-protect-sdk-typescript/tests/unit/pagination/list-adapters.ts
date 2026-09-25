/**
 * The list-method adapter table: every paged operation in the shared vector files, mapped
 * to the TypeScript method that wraps it and to the TypeScript name of each canonical
 * option. Shared by the wire-vector loader and the reply tests, so one table decides both
 * what is sent and which next-page rule is expected.
 *
 * An operation missing from ADAPTERS and from NOT_WRAPPED, or a canonical option missing
 * from an adapter's `options`, fails the loaders: nothing is skipped silently.
 */

import type { PagedServices } from "./harness";

/** How one paged operation is reached through this SDK. */
export interface ListAdapter {
  /** The next-page rule the SDK applies (must equal the shared operations map). */
  readonly rule: string;
  /** Canonical (snake_case) option name → this SDK's option field. */
  readonly options: Readonly<Record<string, string>>;
  /** Calls the method with the mapped options and the vector's path parameters. */
  readonly call: (
    services: PagedServices,
    options: Record<string, unknown>,
    pathParams: Record<string, string>
  ) => Promise<unknown>;
}

const OFFSET = { limit: "limit", offset: "offset" } as const;
const CURSOR = { page_size: "pageSize", cursor: "cursor" } as const;

/* eslint-disable @typescript-eslint/no-explicit-any -- options arrive as vector JSON */
export const ADAPTERS: Readonly<Record<string, ListAdapter>> = {
  ActionService_GetActions: {
    rule: "plus_rows",
    options: OFFSET,
    call: (s, o) => s.actions.list(o as any),
  },
  FeePayerService_GetFeePayers: {
    rule: "plus_rows",
    options: OFFSET,
    call: (s, o) => s.feePayers.list(o as any),
  },
  TransactionService_GetTransactions: {
    rule: "plus_rows",
    options: OFFSET,
    call: (s, o) => s.transactions.list(o as any),
  },
  TransactionService_ExportTransactions: {
    rule: "limit_only",
    options: { limit: "limit" },
    call: (s, o) => s.transactions.exportTransactions(o as any),
  },
  UserService_GetUsers: {
    rule: "plus_min_rows_limit",
    options: OFFSET,
    call: (s, o) => s.users.list(o as any),
  },
  UserService_GetGroups: {
    rule: "plus_min_rows_limit",
    options: OFFSET,
    call: (s, o) => s.groups.list(o as any),
  },
  WalletService_GetWalletsV2: {
    rule: "reply_offset",
    options: { ...OFFSET, exclude_disabled: "excludeDisabled" },
    call: (s, o) => s.wallets.list(o as any),
  },
  WalletService_GetAddresses: {
    rule: "reply_offset",
    options: { ...OFFSET, exclude_disabled: "excludeDisabled" },
    call: (s, o) => s.addresses.listWithOptions(o as any),
  },
  WhitelistService_GetWhitelistedAddresses: {
    rule: "plus_server_rows",
    options: OFFSET,
    call: (s, o) => s.whitelistedAddresses.list(o as any),
  },
  WhitelistService_GetWhitelistedAddressesForApproval: {
    rule: "plus_server_rows",
    options: OFFSET,
    call: (s, o) => s.whitelistedAddresses.listForApproval(o as any),
  },
  WhitelistService_GetWhitelistedContracts: {
    rule: "plus_limit",
    options: OFFSET,
    call: (s, o) => s.whitelistedAssets.list(o as any),
  },
  WhitelistService_GetWhitelistedContractsForApproval: {
    rule: "plus_limit",
    options: OFFSET,
    call: (s, o) => s.whitelistedAssets.listForApproval(o as any),
  },
  PriceService_GetPricesHistory: {
    rule: "limit_only",
    options: { limit: "limit" },
    call: (s, o, p) => s.prices.getHistory({ base: p.base, quote: p.quote, ...(o as any) }),
  },
  AssetServiceV2_ListAssetOperationsV2: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o, p) => s.assets.listAssetOperations(p.assetID, o as any),
  },
  AssetServiceV2_QueryAssetAddressesV2: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o, p) => s.assets.queryAssetAddresses(p.assetID, o as any),
  },
  AssetServiceV2_QueryAssetsV2: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.assets.queryAssets(o as any),
  },
  AuditService_GetAuditTrails: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.audits.list(o as any),
  },
  ChangeService_GetChanges: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.changes.list(o as any),
  },
  ChangeService_GetChangesForApproval: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.changes.listForApproval(o as any),
  },
  EarnService_GetRewards: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.earn.listRewards(o as any),
  },
  ExchangeService_GetExchanges: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.exchanges.list(o as any),
  },
  FiatProviderService_GetFiatProviderAccounts: {
    rule: "cursor",
    options: { ...CURSOR, provider: "provider", label: "label" },
    call: (s, o) => s.fiat.getFiatProviderAccounts(o as any),
  },
  FiatProviderService_GetFiatProviderCounterpartyAccounts: {
    rule: "cursor",
    options: { ...CURSOR, provider: "provider", label: "label" },
    call: (s, o) => s.fiat.getFiatProviderCounterpartyAccounts(o as any),
  },
  FiatProviderService_GetFiatProviderEntities: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.fiat.listFiatProviderEntities(o as any),
  },
  FiatProviderService_GetFiatProviderOperations: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.fiat.getFiatProviderOperations(o as any),
  },
  PriceService_QueryPricesV2: {
    rule: "cursor",
    options: {
      ...CURSOR,
      from_currency_id: "fromCurrencyId",
      to_currency_ids: "toCurrencyIds",
    },
    call: (s, o) => s.prices.list(o as any),
  },
  RequestService_GetRequestsV2: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.requests.list(o as any),
  },
  RequestService_GetRequestsForApprovalV2: {
    rule: "cursor",
    // `statuses` and `externalRequestIds` are not part of the options type; the method rejects
    // both by name.
    options: { ...CURSOR, statuses: "statuses", external_request_ids: "externalRequestIds" },
    call: (s, o) => s.requests.listForApproval(o as any),
  },
  RuleService_GetBusinessRulesV2: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.businessRules.list(o as any),
  },
  RuleService_GetRulesHistory: {
    rule: "token",
    options: CURSOR,
    call: (s, o) => s.governanceRules.getRulesHistory(o as any),
  },
  StakingService_GetStakeAccounts: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.staking.getStakeAccounts(o as any),
  },
  TaurusNetworkService_GetLendingAgreements: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.lending.listLendingAgreements(o as any),
  },
  TaurusNetworkService_GetLendingAgreementsForApproval: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.lending.listLendingAgreementsForApproval(o as any),
  },
  TaurusNetworkService_GetLendingOffers: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.lending.listLendingOffers(o as any),
  },
  TaurusNetworkService_GetPledgeActions: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.pledges.listPledgeActions(o as any),
  },
  TaurusNetworkService_GetPledgeActionsForApproval: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.pledges.listPledgeActionsForApproval(o as any),
  },
  TaurusNetworkService_GetPledges: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.pledges.list(o as any),
  },
  TaurusNetworkService_GetPledgesWithdrawals: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.pledges.listPledgeWithdrawals(o as any),
  },
  TaurusNetworkService_GetSettlements: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.settlements.list(o as any),
  },
  TaurusNetworkService_GetSettlementsForApproval: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.settlements.listForApproval(o as any),
  },
  TaurusNetworkService_GetSharedAddresses: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.sharing.listSharedAddresses(o as any),
  },
  TaurusNetworkService_GetSharedAssets: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.sharing.listSharedAssets(o as any),
  },
  WalletService_GetAssetAddresses: {
    rule: "cursor",
    options: { ...CURSOR, currency: "currency" },
    call: (s, o) => s.assets.getAssetAddresses(o as any),
  },
  WalletService_GetAssetWallets: {
    rule: "cursor",
    options: { ...CURSOR, currency: "currency" },
    call: (s, o) => s.assets.getAssetWallets(o as any),
  },
  WalletService_GetBalances: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.balances.list(o as any),
  },
  WalletService_GetNFTCollectionBalances: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.balances.listNFTCollections(o as any),
  },
  WalletService_GetReservations: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.reservations.list(o as any),
  },
  WalletService_GetWalletTokens: {
    rule: "token",
    options: CURSOR,
    call: (s, o, p) => s.wallets.getWalletTokens(Number(p.id), o as any),
  },
  WebhookService_GetWebhookCalls: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.webhookCalls.list(o as any),
  },
  WebhookService_GetWebhooks: {
    rule: "cursor",
    options: CURSOR,
    call: (s, o) => s.webhooks.list(o as any),
  },
};
/* eslint-enable @typescript-eslint/no-explicit-any */

/**
 * Operations this SDK genuinely does not wrap, with the reason. Only the Go SDK has
 * methods for them; adding one is a new public method, recorded as a known SDK gap.
 */
export const NOT_WRAPPED: Readonly<Record<string, string>> = {
  PriceService_ExportPricesHistory: "no TypeScript method (Go-only ExportPriceHistory)",
  StatisticsService_GetAggregatedTagStats: "no TypeScript method (Go-only ListTagStatistics)",
  StatisticsService_GetPortfolioStatisticsHistory:
    "no TypeScript method (Go-only GetPortfolioStatisticsHistory)",
};

/**
 * Maps a vector's canonical options onto this SDK's option names. An option the adapter
 * does not map throws, so a new canonical option cannot be dropped on the floor.
 */
export function mapOptions(
  operation: string,
  adapter: ListAdapter,
  canonical: Record<string, unknown>
): Record<string, unknown> {
  const mapped: Record<string, unknown> = {};
  for (const [name, value] of Object.entries(canonical)) {
    const field = adapter.options[name];
    if (field === undefined) {
      throw new Error(`${operation}: canonical option "${name}" has no TypeScript mapping`);
    }
    mapped[field] = value;
  }
  return mapped;
}
