/**
 * Cursor and token lists, end to end through the generated client: rows are read from
 * the reply field the generated deserializer actually fills, `hasMore` comes from
 * `hasNext` (or the token's presence), the next cursor is the reply's and nothing else,
 * and a continuation carries it back with NEXT and the page size.
 */

import * as crypto from "crypto";
import { existsSync, readFileSync } from "fs";
import * as path from "path";
import type { RulesContainerCache } from "../../../src/cache";
import { encodePublicKeyPem } from "../../../src/crypto";
import {
  IntegrityError,
  NotFoundError,
  PaginationError,
  ValidationError,
} from "../../../src/errors";
import {
  createEmptyRulesContainer,
  type DecodedRulesContainer,
} from "../../../src/models/governance-rules";
import type { CursorPage } from "../../../src/models/pagination";
import { rulesContainerToBase64 } from "../../../src/mappers/protobuf-rules-container-encode";
import { ADAPTERS } from "./list-adapters";
import { RecordingTransport, pagedServices } from "./harness";

const LIST_REQUEST_VECTORS = path.resolve(
  __dirname,
  "../../../../scripts/resources/list-request-vectors.json"
);

/** The continuation token the shared vectors use (standard base64 with + / =). */
function vectorToken(): string {
  if (!existsSync(LIST_REQUEST_VECTORS)) {
    throw new Error(`cannot read the shared list-request vectors ${LIST_REQUEST_VECTORS}`);
  }
  const vectors = JSON.parse(readFileSync(LIST_REQUEST_VECTORS, "utf-8")) as {
    vectors: Array<{ options: Record<string, unknown> }>;
  };
  const token = vectors.vectors.find((v) => typeof v.options.cursor === "string")?.options
    .cursor as string | undefined;
  if (!token) {
    throw new Error("the shared vectors carry no continuation token");
  }
  return token;
}

/** A rules cache holding an empty (SuperAdmin-verified) container: no price signer. */
function emptyRulesCache(container: DecodedRulesContainer = createEmptyRulesContainer()): RulesContainerCache {
  return {
    get: () => Promise.resolve(container),
    clear: () => undefined,
  } as unknown as RulesContainerCache;
}

/** Where each cursor list's reply keeps its rows, and where its result returns them. */
interface CursorReplyShape {
  rowsField: string;
  resultKey: string;
  /** The reply's total field, for the families that return one. */
  total?: "total" | "totalItems";
  /** Token families carry the token in this field instead of a cursor object. */
  tokenField?: "next" | "cursor";
  /** Rows the reply carries (0 where a row would need a signature to survive). */
  rows?: number;
  options?: Record<string, unknown>;
  pathParams?: Record<string, string>;
}

const SHAPES: Record<string, CursorReplyShape> = {
  AssetServiceV2_ListAssetOperationsV2: {
    rowsField: "result",
    resultKey: "items",
    pathParams: { assetID: "asset-1" },
  },
  AssetServiceV2_QueryAssetAddressesV2: {
    rowsField: "result",
    resultKey: "items",
    pathParams: { assetID: "asset-1" },
  },
  AssetServiceV2_QueryAssetsV2: { rowsField: "result", resultKey: "items" },
  AuditService_GetAuditTrails: { rowsField: "result", resultKey: "items" },
  ChangeService_GetChanges: { rowsField: "result", resultKey: "changes" },
  ChangeService_GetChangesForApproval: { rowsField: "result", resultKey: "changes" },
  EarnService_GetRewards: { rowsField: "rewards", resultKey: "items" },
  ExchangeService_GetExchanges: { rowsField: "result", resultKey: "items" },
  FiatProviderService_GetFiatProviderAccounts: {
    rowsField: "result",
    resultKey: "accounts",
    options: { provider: "p", label: "l" },
  },
  FiatProviderService_GetFiatProviderCounterpartyAccounts: {
    rowsField: "result",
    resultKey: "accounts",
    options: { provider: "p", label: "l" },
  },
  FiatProviderService_GetFiatProviderEntities: { rowsField: "result", resultKey: "items" },
  FiatProviderService_GetFiatProviderOperations: { rowsField: "result", resultKey: "operations" },
  PriceService_QueryPricesV2: { rowsField: "result", resultKey: "items" },
  RequestService_GetRequestsV2: { rowsField: "result", resultKey: "requests" },
  RequestService_GetRequestsForApprovalV2: { rowsField: "result", resultKey: "requests" },
  RuleService_GetBusinessRulesV2: { rowsField: "result", resultKey: "rules" },
  // Unsigned rulesets are withheld (history is lenient), so the rows go to
  // excludedUnverified and reduce the total: see the dedicated test below.
  RuleService_GetRulesHistory: {
    rowsField: "result",
    resultKey: "items",
    total: "totalItems",
    tokenField: "cursor",
    rows: 0,
  },
  StakingService_GetStakeAccounts: { rowsField: "stakeAccounts", resultKey: "stakeAccounts" },
  TaurusNetworkService_GetLendingAgreements: {
    rowsField: "lendingAgreements",
    resultKey: "agreements",
  },
  TaurusNetworkService_GetLendingAgreementsForApproval: {
    rowsField: "result",
    resultKey: "agreements",
  },
  TaurusNetworkService_GetLendingOffers: { rowsField: "lendingOffers", resultKey: "offers" },
  TaurusNetworkService_GetPledgeActions: { rowsField: "result", resultKey: "actions" },
  TaurusNetworkService_GetPledgeActionsForApproval: { rowsField: "result", resultKey: "actions" },
  TaurusNetworkService_GetPledges: { rowsField: "pledges", resultKey: "pledges" },
  TaurusNetworkService_GetPledgesWithdrawals: {
    rowsField: "withdrawals",
    resultKey: "withdrawals",
  },
  TaurusNetworkService_GetSettlements: { rowsField: "result", resultKey: "settlements" },
  TaurusNetworkService_GetSettlementsForApproval: {
    rowsField: "result",
    resultKey: "settlements",
  },
  TaurusNetworkService_GetSharedAddresses: {
    rowsField: "sharedAddresses",
    resultKey: "sharedAddresses",
  },
  TaurusNetworkService_GetSharedAssets: { rowsField: "sharedAssets", resultKey: "sharedAssets" },
  // Every address row is signature-checked, so the rows test needs none here.
  WalletService_GetAssetAddresses: {
    rowsField: "addresses",
    resultKey: "items",
    total: "totalItems",
    rows: 0,
    options: { currency: "ETH" },
  },
  WalletService_GetAssetWallets: {
    rowsField: "wallets",
    resultKey: "items",
    total: "totalItems",
    options: { currency: "ETH" },
  },
  WalletService_GetBalances: { rowsField: "balances", resultKey: "items", total: "total" },
  WalletService_GetNFTCollectionBalances: { rowsField: "balances", resultKey: "items" },
  WalletService_GetReservations: { rowsField: "result", resultKey: "items" },
  WalletService_GetWalletTokens: {
    rowsField: "balances",
    resultKey: "items",
    total: "total",
    tokenField: "next",
    pathParams: { id: "7" },
  },
  WebhookService_GetWebhookCalls: { rowsField: "calls", resultKey: "calls" },
  WebhookService_GetWebhooks: { rowsField: "webhooks", resultKey: "items" },
};

const CURSOR_OPERATIONS = Object.entries(ADAPTERS)
  .filter(([, adapter]) => adapter.rule === "cursor" || adapter.rule === "token")
  .map(([operation]) => operation);

function reply(
  shape: CursorReplyShape,
  next: { token?: string; hasNext?: boolean } | undefined,
  total?: string
): Record<string, unknown> {
  const rowCount = shape.rows ?? 2;
  const body: Record<string, unknown> = {
    [shape.rowsField]: Array.from({ length: rowCount }, (_v, i) => ({ id: String(i + 1) })),
  };
  if (shape.tokenField !== undefined) {
    if (next?.token !== undefined) {
      body[shape.tokenField] = next.token;
    }
  } else if (next !== undefined) {
    body.cursor = { currentPage: next.token, hasNext: next.hasNext };
  }
  if (shape.total !== undefined && total !== undefined) {
    body[shape.total] = total;
  }
  return body;
}

async function run(
  operation: string,
  replies: unknown[],
  options: Record<string, unknown> = {}
): Promise<{ result: Record<string, unknown> & { pagination: CursorPage }; transport: RecordingTransport }> {
  const shape = SHAPES[operation];
  const transport = new RecordingTransport(...replies);
  const services = pagedServices(transport, emptyRulesCache());
  const result = (await ADAPTERS[operation].call(
    services,
    { ...shape.options, ...options },
    shape.pathParams ?? {}
  )) as Record<string, unknown> & { pagination: CursorPage };
  return { result, transport };
}

describe("cursor lists: reply -> pagination", () => {
  it("has a reply shape for every cursor and token operation", () => {
    expect([...CURSOR_OPERATIONS].sort()).toEqual(Object.keys(SHAPES).sort());
  });

  it.each(CURSOR_OPERATIONS)("%s: a middle page continues with the reply's cursor", async (operation) => {
    const token = vectorToken();
    const shape = SHAPES[operation];
    const hasTotal = shape.total !== undefined;
    const { result } = await run(
      operation,
      [reply(shape, { token, hasNext: true }, hasTotal ? "57" : undefined)],
      { pageSize: 7 }
    );
    const expected: CursorPage = hasTotal
      ? { pageSize: 7, nextCursor: token, hasMore: true, totalItems: 57 }
      : { pageSize: 7, nextCursor: token, hasMore: true };
    expect(result.pagination).toEqual(expected);
    expect([operation, (result[shape.resultKey] as unknown[]).length]).toEqual([
      operation,
      shape.rows ?? 2,
    ]);
  });

  it.each(CURSOR_OPERATIONS)("%s: the last page ends the walk", async (operation) => {
    const token = vectorToken();
    const shape = SHAPES[operation];
    // A v2 last page still carries currentPage; only hasNext decides.
    const last = shape.tokenField !== undefined ? undefined : { token, hasNext: undefined };
    const { result } = await run(operation, [reply(shape, last, shape.total ? "3" : undefined)]);
    expect([operation, result.pagination.hasMore, result.pagination.nextCursor]).toEqual([
      operation,
      false,
      "",
    ]);
  });

  it.each(CURSOR_OPERATIONS)("%s answers an empty page {} in full", async (operation) => {
    const shape = SHAPES[operation];
    const { result, transport } = await run(operation, [{}]);
    const expected: CursorPage =
      shape.total !== undefined
        ? { pageSize: 20, nextCursor: "", hasMore: false, totalItems: 0 }
        : { pageSize: 20, nextCursor: "", hasMore: false };
    expect(result.pagination).toEqual(expected);
    expect(result[shape.resultKey]).toEqual([]);
    expect(transport.requests).toHaveLength(1);
  });

  it.each(CURSOR_OPERATIONS.filter((op) => SHAPES[op].tokenField === undefined))(
    "%s refuses a cursor that reports a next page without naming it",
    async (operation) => {
      await expect(run(operation, [{ cursor: { hasNext: true } }])).rejects.toBeInstanceOf(
        PaginationError
      );
    }
  );

  it.each(CURSOR_OPERATIONS.filter((op) => SHAPES[op].tokenField === undefined))(
    "%s rejects cursor + currentPage and cursor + pageRequest before sending",
    async (operation) => {
      const token = vectorToken();
      for (const conflict of [{ currentPage: token }, { pageRequest: "PREVIOUS" }]) {
        const transport = new RecordingTransport({});
        const shape = SHAPES[operation];
        const call = ADAPTERS[operation].call(
          pagedServices(transport, emptyRulesCache()),
          { ...shape.options, cursor: token, ...conflict },
          shape.pathParams ?? {}
        );
        await expect(call).rejects.toBeInstanceOf(ValidationError);
        expect(transport.requests).toHaveLength(0);
      }
    }
  );
});

describe("walking cursor lists", () => {
  /** Two pages: the first names the next, the second is the last. */
  async function walk(
    operation: string
  ): Promise<{ transport: RecordingTransport; firstCursor: string }> {
    const token = vectorToken();
    const shape = SHAPES[operation];
    const first = reply(shape, { token, hasNext: true }, shape.total ? "4" : undefined);
    const last = reply(shape, shape.tokenField ? undefined : { token: "ignored", hasNext: false }, shape.total ? "4" : undefined);
    const transport = new RecordingTransport(first, last);
    const services = pagedServices(transport, emptyRulesCache());
    const call = ADAPTERS[operation].call;
    const page1 = (await call(services, { ...shape.options, pageSize: 2 }, shape.pathParams ?? {})) as {
      pagination: CursorPage;
    };
    expect(page1.pagination.hasMore).toBe(true);
    const page2 = (await call(
      services,
      { ...shape.options, pageSize: 2, cursor: page1.pagination.nextCursor },
      shape.pathParams ?? {}
    )) as { pagination: CursorPage };
    expect(page2.pagination.hasMore).toBe(false);
    return { transport, firstCursor: page1.pagination.nextCursor };
  }

  it("query cursor: the continuation carries currentPage, NEXT and the page size", async () => {
    const { transport, firstCursor } = await walk("RequestService_GetRequestsV2");
    expect(transport.requests[1].query.sort()).toEqual(
      [
        ["cursor.currentPage", firstCursor],
        ["cursor.pageRequest", "NEXT"],
        ["cursor.pageSize", "2"],
      ].sort()
    );
    expect(transport.requests[1].rawQuery).toContain(encodeURIComponent(firstCursor));
  });

  it("requestCursor query: balances continue through requestCursor.* only", async () => {
    const { transport, firstCursor } = await walk("WalletService_GetBalances");
    expect(transport.requests[1].query.sort()).toEqual(
      [
        ["requestCursor.currentPage", firstCursor],
        ["requestCursor.pageRequest", "NEXT"],
        ["requestCursor.pageSize", "2"],
      ].sort()
    );
  });

  it("body cursor: prices continue through the body", async () => {
    const { transport, firstCursor } = await walk("PriceService_QueryPricesV2");
    expect(transport.requests[1].body).toEqual({
      cursor: { currentPage: firstCursor, pageRequest: "NEXT", pageSize: "2" },
    });
  });

  it("token list: wallet tokens send the reply's next token back as cursor", async () => {
    const { transport, firstCursor } = await walk("WalletService_GetWalletTokens");
    expect(transport.requests[1].query.sort()).toEqual(
      [
        ["cursor", firstCursor],
        ["limit", "2"],
      ].sort()
    );
  });

  it("token list: rules history sends the reply's cursor token back", async () => {
    const { transport, firstCursor } = await walk("RuleService_GetRulesHistory");
    expect(transport.requests[1].query.sort()).toEqual(
      [
        ["cursor", firstCursor],
        ["limit", "2"],
      ].sort()
    );
  });
});

describe("rows are read from the fields the generated client fills", () => {
  it("NFT collections read `balances` and map the collection from currencyInfo", async () => {
    const transport = new RecordingTransport({
      balances: [
        {
          currencyInfo: {
            name: "Bored Apes",
            symbol: "BAYC",
            blockchain: "ETH",
            network: "mainnet",
            contractAddress: "0xbc4c",
            logo: "data:image/png;base64,AA==",
          },
          balance: { totalConfirmed: "3" },
        },
      ],
    });
    const page = await pagedServices(transport).balances.listNFTCollections({
      blockchain: "ETH",
      network: "mainnet",
    });
    expect(page.items).toEqual([
      {
        name: "Bored Apes",
        symbol: "BAYC",
        blockchain: "ETH",
        network: "mainnet",
        contractAddress: "0xbc4c",
        count: 3,
        logoUrl: "data:image/png;base64,AA==",
      },
    ]);
  });

  it("balances map the asset and the confirmed amount off the nested row", async () => {
    const transport = new RecordingTransport({
      balances: [
        {
          asset: {
            currency: "USDC",
            currencyInfo: {
              id: "c-usdc",
              symbol: "USDC",
              blockchain: "ETH",
              network: "mainnet",
              contractAddress: "0xa0b8",
            },
          },
          balance: { totalConfirmed: "1000" },
        },
      ],
      total: "1",
    });
    const page = await pagedServices(transport).balances.list({ currency: "USDC", tokenId: "t" });
    expect(page.items[0]).toMatchObject({
      currencyId: "c-usdc",
      currency: "USDC",
      blockchain: "ETH",
      network: "mainnet",
      contractAddress: "0xa0b8",
      balance: "1000",
    });
    expect(page.pagination.totalItems).toBe(1);
    expect(transport.single().query.sort()).toEqual(
      [
        ["currency", "USDC"],
        ["tokenId", "t"],
        ["requestCursor.pageSize", "20"],
      ].sort()
    );
  });

  it("rules history withholds unverifiable rulesets and reduces the total by them", async () => {
    const token = vectorToken();
    // Wire-valid containers carrying no SuperAdmin signature: each fails the threshold.
    const unsigned = {
      rulesContainer: rulesContainerToBase64(createEmptyRulesContainer()),
      rulesSignatures: [],
    };
    const transport = new RecordingTransport({
      result: [unsigned, unsigned],
      totalItems: "10",
      cursor: token,
    });
    const history = await pagedServices(transport).governanceRules.getRulesHistory({ pageSize: 2 });
    expect(history.items).toEqual([]);
    expect(history.excludedUnverified).toHaveLength(2);
    expect(history.pagination).toEqual({
      pageSize: 2,
      nextCursor: token,
      hasMore: true,
      totalItems: 8,
    });
  });

  it("the v2 price list verifies every row exactly as the v1 list did", async () => {
    const signer = crypto.generateKeyPairSync("ec", { namedCurve: "P-256" });
    const container = {
      ...createEmptyRulesContainer(),
      users: [
        {
          id: "price@bank.example",
          publicKeyPem: encodePublicKeyPem(signer.publicKey),
          roles: ["PRICEUPDATER"],
        },
      ],
    } as unknown as DecodedRulesContainer;
    const transport = new RecordingTransport({
      result: [{ currencyFrom: "ETH", currencyTo: "USD", rate: "2500", decimals: "18" }],
    });
    const services = pagedServices(transport, emptyRulesCache(container));
    await expect(services.prices.list()).rejects.toBeInstanceOf(IntegrityError);
  });

  it("the v2 price list sends the filters in the body", async () => {
    const transport = new RecordingTransport({});
    await pagedServices(transport, emptyRulesCache()).prices.list({
      onlyPrimary: true,
      sortOrder: "DESC",
      fromCurrencyId: "c1",
      toCurrencyIds: ["c2"],
    });
    expect(transport.single().body).toEqual({
      onlyPrimary: true,
      sortOrder: "DESC",
      cursor: { pageSize: "20" },
      fromTo: { currencyFromId: "c1", currencyToIds: ["c2"] },
    });
  });
});

describe("get-by-scan helpers walk pages", () => {
  it("webhooks.get walks to the page holding the id", async () => {
    const token = vectorToken();
    const transport = new RecordingTransport(
      { webhooks: [{ id: "a" }], cursor: { currentPage: token, hasNext: true } },
      { webhooks: [{ id: "target", url: "https://hooks.example.invalid" }] }
    );
    const webhook = await pagedServices(transport).webhooks.get("target");
    expect(webhook.id).toBe("target");
    expect(transport.requests).toHaveLength(2);
    expect(transport.requests[0].query).toEqual([["cursor.pageSize", "100"]]);
    expect(transport.requests[1].query.sort()).toEqual(
      [
        ["cursor.currentPage", token],
        ["cursor.pageRequest", "NEXT"],
        ["cursor.pageSize", "100"],
      ].sort()
    );
  });

  it("webhookCalls.get walks to the page holding the id", async () => {
    const token = vectorToken();
    const transport = new RecordingTransport(
      { calls: [{ id: "a" }], cursor: { currentPage: token, hasNext: true } },
      { calls: [{ id: "target" }] }
    );
    const call = await pagedServices(transport).webhookCalls.get("target");
    expect(call.id).toBe("target");
    expect(transport.requests).toHaveLength(2);
  });

  it("stops at the last page with NotFoundError", async () => {
    const transport = new RecordingTransport({ webhooks: [{ id: "a" }] });
    await expect(pagedServices(transport).webhooks.get("missing")).rejects.toBeInstanceOf(
      NotFoundError
    );
    expect(transport.requests).toHaveLength(1);
  });

  it("stops when the server hands back a cursor it already gave", async () => {
    const token = vectorToken();
    const looping = { calls: [{ id: "a" }], cursor: { currentPage: token, hasNext: true } };
    const transport = new RecordingTransport(looping);
    await expect(pagedServices(transport).webhookCalls.get("missing")).rejects.toBeInstanceOf(
      NotFoundError
    );
    expect(transport.requests).toHaveLength(2);
  });
});
