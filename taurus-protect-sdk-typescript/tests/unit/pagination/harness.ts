/**
 * Transport-level harness for the pagination tests.
 *
 * Every service is built on the REAL generated `*Api`, whose `Configuration` carries a
 * recording `fetchApi` (`src/internal/openapi/runtime.ts`). The generated request
 * serialization and reply deserializers therefore run exactly as in production — the
 * layer the older, above-the-deserializer mocks skipped, which is how lists came to read
 * reply fields the deserializer drops (NFT collections, exchange/fiat/webhook-call
 * cursors, business-rules `hasMore`).
 */

import * as crypto from "crypto";
import type { RulesContainerCache } from "../../../src/cache";
import {
  ActionsApi,
  AddressWhitelistingApi,
  AddressesApi,
  AssetV2Api,
  AssetsApi,
  AuditApi,
  BalancesApi,
  BusinessRulesApi,
  ChangesApi,
  ContractWhitelistingApi,
  EarnApi,
  ExchangeApi,
  FeePayersApi,
  FiatApi,
  GovernanceRulesApi,
  GroupsApi,
  PricesApi,
  RequestsApi,
  ReservationsApi,
  StakingApi,
  TaurusNetworkLendingApi,
  TaurusNetworkPledgeApi,
  TaurusNetworkSettlementApi,
  TaurusNetworkSharedAddressAssetApi,
  TransactionsApi,
  UsersApi,
  WalletsApi,
  WebhookCallsApi,
  WebhooksApi,
} from "../../../src/internal/openapi/apis";
import { Configuration } from "../../../src/internal/openapi/runtime";
import {
  createEmptyRulesContainer,
  type DecodedRulesContainer,
  type RuleUserSignature,
} from "../../../src/models/governance-rules";
import { ActionService } from "../../../src/services/action-service";
import { AddressService } from "../../../src/services/address-service";
import { AssetService } from "../../../src/services/asset-service";
import { AuditService } from "../../../src/services/audit-service";
import { BalanceService } from "../../../src/services/balance-service";
import { BusinessRuleService } from "../../../src/services/business-rule-service";
import { ChangeService } from "../../../src/services/change-service";
import { EarnService } from "../../../src/services/earn-service";
import { ExchangeService } from "../../../src/services/exchange-service";
import { FeePayerService } from "../../../src/services/fee-payer-service";
import { FiatService } from "../../../src/services/fiat-service";
import { GovernanceRuleService } from "../../../src/services/governance-rule-service";
import { GroupService } from "../../../src/services/group-service";
import { PriceService } from "../../../src/services/price-service";
import { RequestService } from "../../../src/services/request-service";
import { ReservationService } from "../../../src/services/reservation-service";
import { StakingService } from "../../../src/services/staking-service";
import { TransactionService } from "../../../src/services/transaction-service";
import { UserService } from "../../../src/services/user-service";
import { WalletService } from "../../../src/services/wallet-service";
import { WebhookCallService } from "../../../src/services/webhook-call-service";
import { WebhookService } from "../../../src/services/webhook-service";
import { WhitelistedAddressService } from "../../../src/services/whitelisted-address-service";
import { WhitelistedAssetService } from "../../../src/services/whitelisted-asset-service";
import { LendingService } from "../../../src/services/taurus-network/lending-service";
import { PledgeService } from "../../../src/services/taurus-network/pledge-service";
import { SettlementService } from "../../../src/services/taurus-network/settlement-service";
import { SharingService } from "../../../src/services/taurus-network/sharing-service";

/** Placeholder host: nothing is ever sent over the network. */
export const BASE_PATH = "https://protect.example.invalid";

/** One request as it reached the transport. */
export interface RecordedRequest {
  readonly method: string;
  /** Path relative to the host, still percent-encoded. */
  readonly path: string;
  /** The query string exactly as sent, without the leading "?". */
  readonly rawQuery: string;
  /** Query pairs in wire order, URL-decoded. */
  readonly query: Array<[string, string]>;
  /** The parsed JSON body, when there is one. */
  readonly body: unknown;
}

/**
 * A fetch stand-in that records every request and answers with canned JSON replies, one
 * per request in order (the last one repeats).
 */
export class RecordingTransport {
  readonly requests: RecordedRequest[] = [];
  readonly configuration: Configuration;
  private readonly replies: unknown[];

  constructor(...replies: unknown[]) {
    this.replies = replies.length > 0 ? [...replies] : [{}];
    this.configuration = new Configuration({
      basePath: BASE_PATH,
      fetchApi: (input, init) => this.handle(input, init),
    });
  }

  /** The only request sent; fails the test when there were none or several. */
  single(): RecordedRequest {
    expect(this.requests).toHaveLength(1);
    return this.requests[0];
  }

  private async handle(input: RequestInfo | URL, init?: RequestInit): Promise<Response> {
    const href =
      typeof input === "string" ? input : input instanceof URL ? input.href : input.url;
    const url = new URL(href);
    const query: Array<[string, string]> = [];
    url.searchParams.forEach((value, key) => query.push([key, value]));
    const rawBody = init?.body;
    this.requests.push({
      method: init?.method ?? "GET",
      path: url.pathname,
      rawQuery: url.search.replace(/^\?/, ""),
      query,
      body: typeof rawBody === "string" && rawBody.length > 0 ? JSON.parse(rawBody) : undefined,
    });
    const reply = this.replies.length > 1 ? this.replies.shift() : this.replies[0];
    return new Response(JSON.stringify(reply ?? {}), {
      status: 200,
      headers: { "Content-Type": "application/json" },
    });
  }
}

/** Sorts query pairs, for the multiset comparison the vectors specify. */
export function sortedPairs(pairs: Array<[string, string]>): Array<[string, string]> {
  return [...pairs].sort((a, b) =>
    a[0] === b[0] ? (a[1] < b[1] ? -1 : a[1] > b[1] ? 1 : 0) : a[0] < b[0] ? -1 : 1
  );
}

/**
 * A rules cache that fails the test if a list path consults it: the replies these tests
 * use carry no signed rows, so nothing should need the container.
 */
export function unusedRulesCache(): RulesContainerCache {
  return {
    get: () => Promise.reject(new Error("the rules container must not be fetched")),
    clear: () => undefined,
  } as unknown as RulesContainerCache;
}

// A throwaway P-256 public key: the whitelist verifiers need one to construct, and the
// replies these tests use carry no rows for them to verify.
const THROWAWAY_KEY_PEM = crypto
  .generateKeyPairSync("ec", { namedCurve: "P-256" })
  .publicKey.export({ type: "spki", format: "pem" })
  .toString();

function whitelistConfig(): {
  superAdminKeysPem: string[];
  minValidSignatures: number;
  rulesContainerDecoder: (base64: string) => DecodedRulesContainer;
  userSignaturesDecoder: (base64: string) => RuleUserSignature[];
} {
  return {
    superAdminKeysPem: [THROWAWAY_KEY_PEM],
    minValidSignatures: 1,
    rulesContainerDecoder: () => createEmptyRulesContainer(),
    userSignaturesDecoder: () => [],
  };
}

/** Every paged service, built on the real generated APIs over one transport. */
export interface PagedServices {
  actions: ActionService;
  addresses: AddressService;
  assets: AssetService;
  audits: AuditService;
  balances: BalanceService;
  businessRules: BusinessRuleService;
  changes: ChangeService;
  earn: EarnService;
  exchanges: ExchangeService;
  feePayers: FeePayerService;
  fiat: FiatService;
  governanceRules: GovernanceRuleService;
  groups: GroupService;
  lending: LendingService;
  pledges: PledgeService;
  prices: PriceService;
  requests: RequestService;
  reservations: ReservationService;
  settlements: SettlementService;
  sharing: SharingService;
  staking: StakingService;
  transactions: TransactionService;
  users: UserService;
  wallets: WalletService;
  webhookCalls: WebhookCallService;
  webhooks: WebhookService;
  whitelistedAddresses: WhitelistedAddressService;
  whitelistedAssets: WhitelistedAssetService;
}

/**
 * Builds every paged service over `transport`.
 *
 * @param transport - The recording transport the generated APIs send through
 * @param rulesCache - The rules cache for the services that verify addresses and prices
 */
export function pagedServices(
  transport: RecordingTransport,
  rulesCache: RulesContainerCache = unusedRulesCache()
): PagedServices {
  const c = transport.configuration;
  const addresses = new AddressService(new AddressesApi(c), rulesCache);
  const whitelistedAddresses = new WhitelistedAddressService(
    new AddressWhitelistingApi(c),
    whitelistConfig()
  );
  return {
    actions: new ActionService(new ActionsApi(c)),
    addresses,
    assets: new AssetService(new AssetsApi(c), rulesCache, new AssetV2Api(c), {
      addresses: () => addresses,
      whitelistedAddresses: () => whitelistedAddresses,
    }),
    audits: new AuditService(new AuditApi(c)),
    balances: new BalanceService(new BalancesApi(c)),
    businessRules: new BusinessRuleService(new BusinessRulesApi(c)),
    changes: new ChangeService(new ChangesApi(c)),
    earn: new EarnService(new EarnApi(c)),
    exchanges: new ExchangeService(new ExchangeApi(c)),
    feePayers: new FeePayerService(new FeePayersApi(c)),
    fiat: new FiatService(new FiatApi(c)),
    governanceRules: new GovernanceRuleService(new GovernanceRulesApi(c), {
      superAdminKeys: [crypto.createPublicKey(THROWAWAY_KEY_PEM)],
      minValidSignatures: 1,
    }),
    groups: new GroupService(new GroupsApi(c)),
    lending: new LendingService(new TaurusNetworkLendingApi(c)),
    pledges: new PledgeService(new TaurusNetworkPledgeApi(c)),
    prices: new PriceService(new PricesApi(c), rulesCache),
    requests: new RequestService(new RequestsApi(c)),
    reservations: new ReservationService(new ReservationsApi(c)),
    settlements: new SettlementService(new TaurusNetworkSettlementApi(c)),
    sharing: new SharingService(new TaurusNetworkSharedAddressAssetApi(c)),
    staking: new StakingService(new StakingApi(c)),
    transactions: new TransactionService(new TransactionsApi(c)),
    users: new UserService(new UsersApi(c)),
    wallets: new WalletService(new WalletsApi(c)),
    webhookCalls: new WebhookCallService(new WebhookCallsApi(c)),
    webhooks: new WebhookService(new WebhooksApi(c)),
    whitelistedAddresses,
    whitelistedAssets: new WhitelistedAssetService(
      new ContractWhitelistingApi(c),
      whitelistConfig()
    ),
  };
}
