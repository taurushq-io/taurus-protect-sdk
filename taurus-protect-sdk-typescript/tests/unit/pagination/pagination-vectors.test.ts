/**
 * Cross-SDK pagination vectors: the pure pagination helpers against
 * `scripts/resources/pagination-vectors.json`, which all four SDKs load.
 *
 * The vectors are loaded inside `it()`, never in a describe body: a throw during
 * collection becomes "Test suite failed to run" and drops the test count instead of
 * showing a red test. Jest has no per-assertion message, so loops assert tuples that
 * carry the vector's description.
 */

import { existsSync, readFileSync } from "fs";
import * as path from "path";
import { PaginationError, ValidationError } from "../../../src/errors";
import {
  DEFAULT_PAGE_SIZE,
  MAX_COUNT,
  MAX_PAGE_SIZE,
  MAX_PRICE_HISTORY_LIMIT,
  buildCursorPage,
  buildOffsetPagination,
  resolveOffset,
  resolvePageSize,
  type OffsetRule,
  type PageSizeKind,
  type ReplyCursor,
} from "../../../src/models/pagination";

// tests/unit/pagination -> tests/unit -> tests -> <sdk> -> <repo root>
const VECTORS_PATH = path.resolve(
  __dirname,
  "../../../../scripts/resources/pagination-vectors.json"
);

// The counts this loader was written against; regenerate and bump together.
const EXPECTED_COUNTS = {
  operations: 53,
  offset: 26,
  cursor: 12,
  page_size: 14,
  offset_input: 4,
};

interface OffsetVector {
  description: string;
  rule: OffsetRule;
  request: { limit: number; offset: number };
  served_rows: number;
  sdk_excluded: number;
  reply: { totalItems?: string; offset?: string };
  expect?: {
    limit: number;
    offset: number;
    total_items: number;
    next_offset: number;
    has_more: boolean;
  };
  expect_error?: boolean;
}

interface CursorVector {
  description: string;
  family: "cursor" | "token";
  page_size: number;
  has_total: boolean;
  reply_cursor?: ReplyCursor | null;
  reply_token?: string | null;
  reply_total: string | null;
  expect?: {
    page_size: number;
    next_cursor: string;
    has_more: boolean;
    total_items: number | null;
  };
  expect_error?: boolean;
}

interface PageSizeVector {
  description: string;
  kind: PageSizeKind;
  input: number | null;
  expect?: number;
  expect_error?: boolean;
}

interface OffsetInputVector {
  description: string;
  input: number | null;
  expect?: number;
  expect_error?: boolean;
}

interface PaginationVectors {
  constants: {
    default_page_size: number;
    max_page_size: number;
    price_history_max: number;
    max_count: number;
  };
  counts: Record<string, number>;
  operations: Record<string, { rule: string; total: boolean; size?: string }>;
  offset: OffsetVector[];
  cursor: CursorVector[];
  page_size: PageSizeVector[];
  offset_input: OffsetInputVector[];
}

function loadVectors(): PaginationVectors {
  if (!existsSync(VECTORS_PATH)) {
    throw new Error(`cannot read the shared pagination vectors ${VECTORS_PATH}`);
  }
  return JSON.parse(readFileSync(VECTORS_PATH, "utf-8")) as PaginationVectors;
}

/** Runs `fn`, returning its value or the class name of what it threw. */
function outcome<T>(fn: () => T): T | { threw: string } {
  try {
    return fn();
  } catch (err) {
    return { threw: err instanceof Error ? err.constructor.name : typeof err };
  }
}

describe("pagination vectors (scripts/resources/pagination-vectors.json)", () => {
  it("declares the counts this loader was written against", () => {
    const v = loadVectors();
    expect(v.counts).toEqual(EXPECTED_COUNTS);
    expect(["operations", Object.keys(v.operations).length]).toEqual([
      "operations",
      EXPECTED_COUNTS.operations,
    ]);
    for (const section of ["offset", "cursor", "page_size", "offset_input"] as const) {
      expect([section, v[section].length]).toEqual([section, EXPECTED_COUNTS[section]]);
    }
  });

  it("pins the same constants as src/models/pagination.ts", () => {
    const v = loadVectors();
    expect(v.constants).toEqual({
      default_page_size: DEFAULT_PAGE_SIZE,
      max_page_size: MAX_PAGE_SIZE,
      price_history_max: MAX_PRICE_HISTORY_LIMIT,
      max_count: MAX_COUNT,
    });
  });

  it("builds every offset page the vectors describe", () => {
    const v = loadVectors();
    for (const vector of v.offset) {
      const actual = outcome(() => {
        const p = buildOffsetPagination(
          vector.rule,
          vector.request,
          vector.reply,
          vector.served_rows,
          vector.sdk_excluded
        );
        return {
          limit: p.limit,
          offset: p.offset,
          total_items: p.totalItems,
          next_offset: p.nextOffset,
          has_more: p.hasMore,
        };
      });
      const expected = vector.expect_error ? { threw: PaginationError.name } : vector.expect;
      expect([vector.description, actual]).toEqual([vector.description, expected]);
    }
  });

  it("builds every cursor and token page the vectors describe", () => {
    const v = loadVectors();
    for (const vector of v.cursor) {
      const next =
        vector.family === "token"
          ? (vector.reply_token ?? undefined)
          : (vector.reply_cursor ?? undefined);
      const total = vector.has_total ? { total: vector.reply_total ?? undefined } : undefined;
      const actual = outcome(() => {
        const p = buildCursorPage(vector.page_size, next, total);
        return {
          page_size: p.pageSize,
          next_cursor: p.nextCursor,
          has_more: p.hasMore,
          total_items: p.totalItems === undefined ? null : p.totalItems,
        };
      });
      const expected = vector.expect_error ? { threw: PaginationError.name } : vector.expect;
      expect([vector.description, actual]).toEqual([vector.description, expected]);
    }
  });

  it("resolves every page size the vectors describe", () => {
    const v = loadVectors();
    for (const vector of v.page_size) {
      const actual = outcome(() =>
        resolvePageSize(vector.input ?? undefined, "pageSize", vector.kind)
      );
      const expected = vector.expect_error ? { threw: ValidationError.name } : vector.expect;
      expect([vector.description, vector.kind, actual]).toEqual([
        vector.description,
        vector.kind,
        expected,
      ]);
    }
  });

  it("resolves every offset the vectors describe", () => {
    const v = loadVectors();
    for (const vector of v.offset_input) {
      const actual = outcome(() => resolveOffset(vector.input ?? undefined));
      const expected = vector.expect_error ? { threw: ValidationError.name } : vector.expect;
      expect([vector.description, actual]).toEqual([vector.description, expected]);
    }
  });

  it("names the option and the bound in a page-size rejection", () => {
    expect(() => resolvePageSize(101, "pageSize")).toThrow("pageSize must be at most 100, got 101");
    expect(() => resolvePageSize(-1, "limit")).toThrow("limit must not be negative, got -1");
    expect(() => resolvePageSize(366, "limit", "price_history")).toThrow(
      "limit must be at most 365, got 366"
    );
    expect(() => resolvePageSize(1.5, "pageSize")).toThrow("pageSize must be a whole number");
    expect(() => resolveOffset(-1)).toThrow("offset must not be negative, got -1");
  });
});
