/**
 * Offset lists, end to end through the generated client: the reply a real validatord
 * sends (zero values omitted) becomes a pagination value that is never undefined, whose
 * next offset follows the operation's rule, and whose total excludes withheld rows.
 *
 * The expected next offset is computed by `referenceNextOffset` — an independent copy of
 * the contract's rule table — so these tests do not trust `buildOffsetPagination`.
 */

import { IntegrityError, PaginationError, ValidationError } from "../../../src/errors";
import { createEmptyRulesContainer } from "../../../src/models/governance-rules";
import type { WhitelistedAddress } from "../../../src/models/whitelisted-address";
import type { WhitelistedAsset } from "../../../src/models/whitelisted-asset";
import type { Pagination } from "../../../src/models/pagination";
import { ADAPTERS } from "./list-adapters";
import { RecordingTransport, pagedServices, type PagedServices } from "./harness";

const OFFSET_RULES = [
  "reply_offset",
  "plus_rows",
  "plus_min_rows_limit",
  "plus_server_rows",
  "plus_limit",
];

/** The offset operations, by adapter rule. */
const OFFSET_OPERATIONS = Object.entries(ADAPTERS)
  .filter(([, adapter]) => OFFSET_RULES.includes(adapter.rule))
  .map(([operation]) => operation);

/** The contract's rule table, restated independently of the SDK. */
function referenceNextOffset(
  rule: string,
  limit: number,
  offset: number,
  servedRows: number,
  replyOffset: number | undefined
): number {
  switch (rule) {
    case "reply_offset":
      return replyOffset ?? offset + servedRows;
    case "plus_rows":
    case "plus_server_rows":
      return offset + servedRows;
    case "plus_min_rows_limit":
      return offset + Math.min(servedRows, limit);
    case "plus_limit":
      return offset + limit;
    default:
      throw new Error(`not an offset rule: ${rule}`);
  }
}

/** Makes every row pass the whitelist verifiers, so the tests see the paging only. */
function acceptEveryRow(services: PagedServices, failFirst = 0): void {
  let seen = 0;
  const reject = (): boolean => {
    seen += 1;
    return seen <= failFirst;
  };
  (
    services.whitelistedAddresses as unknown as { verifier: { verify: () => unknown } }
  ).verifier.verify = () => {
    if (reject()) {
      throw new IntegrityError("row failed verification");
    }
    return {
      verifiedWhitelistedAddress: { id: "1", address: "0xabc" } as unknown as WhitelistedAddress,
      rulesContainer: createEmptyRulesContainer(),
      verifiedHash: "h",
    };
  };
  (
    services.whitelistedAssets as unknown as { verifier: { verify: () => unknown } }
  ).verifier.verify = () => {
    if (reject()) {
      throw new IntegrityError("row failed verification");
    }
    return {
      verifiedAsset: { id: 1 } as unknown as WhitelistedAsset,
      rulesContainer: createEmptyRulesContainer(),
      verifiedHash: "h",
    };
  };
}

/** `n` rows every offset list's mapper accepts, carrying nothing that needs verifying. */
function rows(n: number): Array<Record<string, unknown>> {
  return Array.from({ length: n }, (_v, i) => ({
    id: String(i + 1),
    metadata: { hash: `h${i + 1}`, payloadAsString: "{}" },
  }));
}

async function run(
  operation: string,
  reply: unknown,
  options: Record<string, unknown>,
  failFirst = 0
): Promise<{ result: { items: unknown[]; pagination: Pagination }; transport: RecordingTransport }> {
  const transport = new RecordingTransport(reply);
  const services = pagedServices(transport);
  acceptEveryRow(services, failFirst);
  const result = (await ADAPTERS[operation].call(services, options, {})) as {
    items: unknown[];
    pagination: Pagination;
  };
  return { result, transport };
}

describe("offset lists: reply -> pagination", () => {
  it("covers every offset operation of the adapter table", () => {
    expect([...OFFSET_OPERATIONS].sort()).toEqual(
      [
        "ActionService_GetActions",
        "FeePayerService_GetFeePayers",
        "TransactionService_GetTransactions",
        "UserService_GetGroups",
        "UserService_GetUsers",
        "WalletService_GetAddresses",
        "WalletService_GetWalletsV2",
        "WhitelistService_GetWhitelistedAddresses",
        "WhitelistService_GetWhitelistedAddressesForApproval",
        "WhitelistService_GetWhitelistedContracts",
        "WhitelistService_GetWhitelistedContractsForApproval",
      ].sort()
    );
  });

  // Two replies that tell every rule apart: 21 rows (one beyond the limit, as with a
  // synthetic daemon user) with a reply offset, then a short page without one.
  const scenarios = [
    { name: "over-full page with a reply offset", served: 21, replyOffset: 70 },
    { name: "short page, reply offset omitted", served: 3, replyOffset: undefined },
  ];

  it.each(OFFSET_OPERATIONS)("%s follows its next-offset rule", async (operation) => {
    const rule = ADAPTERS[operation].rule;
    for (const scenario of scenarios) {
      const reply: Record<string, unknown> = {
        result: rows(scenario.served),
        totalItems: "200",
      };
      if (scenario.replyOffset !== undefined) {
        reply.offset = String(scenario.replyOffset);
      }
      const { result } = await run(operation, reply, { limit: 20, offset: 40 });
      const next = referenceNextOffset(rule, 20, 40, scenario.served, scenario.replyOffset);
      expect([scenario.name, result.pagination]).toEqual([
        scenario.name,
        { limit: 20, offset: 40, totalItems: 200, nextOffset: next, hasMore: next > 40 && next < 200 },
      ]);
    }
  });

  it.each(OFFSET_OPERATIONS)("%s answers an empty page {} in full", async (operation) => {
    const rule = ADAPTERS[operation].rule;
    const { result, transport } = await run(operation, {}, {});
    expect(result.items).toEqual([]);
    expect(result.pagination).toEqual({
      limit: 20,
      offset: 0,
      totalItems: 0,
      nextOffset: referenceNextOffset(rule, 20, 0, 0, undefined),
      hasMore: false,
    });
    // The default page size is always sent; an offset of 0 is not.
    const query = transport.single().query;
    expect(query).toContainEqual(["limit", "20"]);
    expect(query.find(([key]) => key === "offset")).toBeUndefined();
  });

  it.each(OFFSET_OPERATIONS)("%s ends the walk on the last page", async (operation) => {
    const rule = ADAPTERS[operation].rule;
    const served = 5;
    const reply: Record<string, unknown> = { result: rows(served), totalItems: "45" };
    if (rule === "reply_offset") {
      reply.offset = "45";
    }
    const { result } = await run(operation, reply, { limit: 20, offset: 40 });
    expect(result.pagination.hasMore).toBe(false);
    expect(result.pagination.offset).toBe(40);
  });

  it.each(OFFSET_OPERATIONS)("%s refuses a count it cannot read", async (operation) => {
    for (const bad of ["abc", "-1", "020", "9007199254740992"]) {
      await expect(
        run(operation, { result: rows(1), totalItems: bad }, {})
      ).rejects.toBeInstanceOf(PaginationError);
    }
  });

  it("refuses a non-numeric reply offset on the reply-offset lists", async () => {
    for (const operation of ["WalletService_GetWalletsV2", "WalletService_GetAddresses"]) {
      await expect(
        run(operation, { result: rows(1), totalItems: "5", offset: "x" }, {})
      ).rejects.toBeInstanceOf(PaginationError);
    }
  });

  it.each([
    "WhitelistService_GetWhitelistedAddresses",
    "WhitelistService_GetWhitelistedAddressesForApproval",
    "WhitelistService_GetWhitelistedContracts",
    "WhitelistService_GetWhitelistedContractsForApproval",
  ])("%s: withheld rows reduce the total, never the next offset", async (operation) => {
    const rule = ADAPTERS[operation].rule;
    const { result } = await run(
      operation,
      { result: rows(20), totalItems: "30" },
      { limit: 20 },
      2
    );
    const next = referenceNextOffset(rule, 20, 0, 20, undefined);
    expect(result.items).toHaveLength(18);
    expect(result.pagination).toEqual({
      limit: 20,
      offset: 0,
      totalItems: 28,
      nextOffset: next,
      hasMore: true,
    });

    // hasMore compares against the SERVER total: with 21 rows on the server and 2
    // withheld, the reduced total (19) is below the next offset (20), yet page 2 exists.
    const { result: edge } = await run(
      operation,
      { result: rows(20), totalItems: "21" },
      { limit: 20 },
      2
    );
    expect(edge.pagination).toEqual({
      limit: 20,
      offset: 0,
      totalItems: 19,
      nextOffset: next,
      hasMore: true,
    });
  });

  it("wallets: the reply offset is the next page, so the last page is not dropped", async () => {
    // The old reading took the reply offset as the CURRENT one: offset=20 with a total of
    // 25 read as "a page at 20 exists" on page 1, and as nothing more on page 2.
    const { result } = await run(
      "WalletService_GetWalletsV2",
      { result: rows(20), totalItems: "25", offset: "20" },
      {}
    );
    expect(result.pagination).toEqual({
      limit: 20,
      offset: 0,
      totalItems: 25,
      nextOffset: 20,
      hasMore: true,
    });
  });

  it("rejects out-of-bounds pages by name before sending", async () => {
    for (const operation of OFFSET_OPERATIONS) {
      for (const [options, message] of [
        [{ limit: 101 }, "limit must be at most 100, got 101"],
        [{ limit: -1 }, "limit must not be negative, got -1"],
        [{ offset: -1 }, "offset must not be negative, got -1"],
      ] as const) {
        const transport = new RecordingTransport({});
        const call = ADAPTERS[operation].call(pagedServices(transport), { ...options }, {});
        await expect(call).rejects.toThrow(ValidationError);
        await expect(call).rejects.toThrow(message);
        expect([operation, transport.requests.length]).toEqual([operation, 0]);
      }
    }
  });
});

describe("offset lists beyond the adapter table", () => {
  it("transactions by request and by address page like the main list", async () => {
    const transport = new RecordingTransport({ result: rows(2), totalItems: "12" });
    const services = pagedServices(transport);

    const byRequest = await services.transactions.listByRequest("r-1", { limit: 2, offset: 10 });
    expect(byRequest.pagination).toEqual({
      limit: 2,
      offset: 10,
      totalItems: 12,
      nextOffset: 12,
      hasMore: false,
    });

    const byAddress = await services.transactions.listByAddress("addr-1", { limit: 2 });
    expect(byAddress.pagination).toEqual({
      limit: 2,
      offset: 0,
      totalItems: 12,
      nextOffset: 2,
      hasMore: true,
    });
    expect(transport.requests[0].query).toEqual(
      expect.arrayContaining([
        ["transactionIds", "r-1"],
        ["limit", "2"],
        ["offset", "10"],
      ])
    );
    expect(transport.requests[1].query).toEqual(
      expect.arrayContaining([
        ["address", "addr-1"],
        ["limit", "2"],
      ])
    );
  });

  it("addresses send every filter, and excludeDisabled as includeDisabledAddresses=exclude", async () => {
    const transport = new RecordingTransport({});
    await pagedServices(transport).addresses.listWithOptions({
      walletId: "7",
      query: "q",
      blockchain: "ETH",
      network: "mainnet",
      addressIds: ["1", "2"],
      addresses: ["0xa"],
      tagIds: ["t"],
      onlyPositiveBalance: true,
      balanceAbove: "10",
      balanceBelow: "20",
      sortBy: "balance",
      sortOrder: "DESC",
      excludeDisabled: true,
      limit: 5,
      offset: 10,
    });
    expect(transport.single().query.sort()).toEqual(
      [
        ["walletId", "7"],
        ["query", "q"],
        ["blockchain", "ETH"],
        ["network", "mainnet"],
        ["addressIds", "1"],
        ["addressIds", "2"],
        ["addresses", "0xa"],
        ["tagIDs", "t"],
        ["onlyPositiveBalance", "true"],
        ["balanceAbove", "10"],
        ["balanceBelow", "20"],
        ["sortBy", "balance"],
        ["sortOrder", "DESC"],
        ["includeDisabledAddresses", "exclude"],
        ["limit", "5"],
        ["offset", "10"],
      ].sort()
    );
  });

  it("the wallet-scoped address list pages like the filtered one", async () => {
    const transport = new RecordingTransport({ result: rows(1), totalItems: "1", offset: "1" });
    const page = await pagedServices(transport).addresses.list(42, { limit: 1 });
    expect(page.pagination).toEqual({
      limit: 1,
      offset: 0,
      totalItems: 1,
      nextOffset: 1,
      hasMore: false,
    });
    expect(transport.single().query.sort()).toEqual(
      [
        ["walletId", "42"],
        ["limit", "1"],
      ].sort()
    );
  });
});

describe("transaction export (cannot page)", () => {
  it("returns the text and the server's total, sending only the limit", async () => {
    const transport = new RecordingTransport({ result: "id,amount\n1,5\n", totalItems: "37" });
    const exported = await pagedServices(transport).transactions.exportTransactions({
      limit: 1000,
    });
    expect(exported).toEqual({ data: "id,amount\n1,5\n", totalItems: 37 });
    expect(transport.single().query).toEqual([["limit", "1000"]]);
  });

  it("answers {} with an empty export", async () => {
    const transport = new RecordingTransport({});
    const exported = await pagedServices(transport).transactions.exportTransactions();
    expect(exported).toEqual({ data: "", totalItems: 0 });
    expect(transport.single().query).toEqual([["limit", "20"]]);
  });

  it("sends the format only when asked", async () => {
    const transport = new RecordingTransport({});
    await pagedServices(transport).transactions.exportTransactions({ format: "csv" });
    expect(transport.single().query.sort()).toEqual(
      [
        ["format", "csv"],
        ["limit", "20"],
      ].sort()
    );
  });

  it("refuses a total it cannot read", async () => {
    const transport = new RecordingTransport({ result: "", totalItems: "many" });
    await expect(
      pagedServices(transport).transactions.exportTransactions()
    ).rejects.toBeInstanceOf(PaginationError);
  });
});
