/**
 * Cross-SDK list-request vectors: every paged operation, called through the real
 * generated client with a recording transport and a `{}` reply, must put exactly the
 * vector's query / body on the wire — or reject the options with a ValidationError
 * before sending anything.
 *
 * Vectors are loaded inside `it()`; loops assert tuples carrying the vector description.
 */

import { existsSync, readFileSync } from "fs";
import * as path from "path";
import { ValidationError } from "../../../src/errors";
import { ADAPTERS, NOT_WRAPPED, mapOptions } from "./list-adapters";
import { RecordingTransport, pagedServices, sortedPairs } from "./harness";

// tests/unit/pagination -> tests/unit -> tests -> <sdk> -> <repo root>
const VECTORS_PATH = path.resolve(
  __dirname,
  "../../../../scripts/resources/list-request-vectors.json"
);
const PAGINATION_VECTORS_PATH = path.resolve(
  __dirname,
  "../../../../scripts/resources/pagination-vectors.json"
);

// The counts this loader was written against; regenerate and bump together.
const EXPECTED_COUNTS = { methods: 53, vectors: 278 };

interface MethodSpec {
  method: string;
  path: string;
  paging: { style: string; prefix?: string; field?: string };
  rule: string;
}

interface RequestVector {
  operation: string;
  path_params: Record<string, string>;
  description: string;
  options: Record<string, unknown>;
  expect?: { query?: Array<[string, string]>; body?: unknown };
  expect_error?: boolean;
}

interface RequestVectors {
  counts: { methods: number; vectors: number };
  methods: Record<string, MethodSpec>;
  vectors: RequestVector[];
}

function loadVectors(): RequestVectors {
  if (!existsSync(VECTORS_PATH)) {
    throw new Error(`cannot read the shared list-request vectors ${VECTORS_PATH}`);
  }
  return JSON.parse(readFileSync(VECTORS_PATH, "utf-8")) as RequestVectors;
}

function loadOperationRules(): Record<string, { rule: string }> {
  if (!existsSync(PAGINATION_VECTORS_PATH)) {
    throw new Error(`cannot read the shared pagination vectors ${PAGINATION_VECTORS_PATH}`);
  }
  return (
    JSON.parse(readFileSync(PAGINATION_VECTORS_PATH, "utf-8")) as {
      operations: Record<string, { rule: string }>;
    }
  ).operations;
}

/** The path the vector names, its parameters substituted as the generated client does. */
function expectedPath(template: string, params: Record<string, string>): string {
  return template.replace(/\{(\w+)\}/g, (_m, name: string) => {
    const value = params[name];
    if (value === undefined) {
      throw new Error(`path parameter ${name} missing from the vector`);
    }
    return encodeURIComponent(value);
  });
}

describe("list-request vectors (scripts/resources/list-request-vectors.json)", () => {
  it("declares the counts this loader was written against", () => {
    const v = loadVectors();
    expect(v.counts).toEqual(EXPECTED_COUNTS);
    expect(["methods", Object.keys(v.methods).length]).toEqual([
      "methods",
      EXPECTED_COUNTS.methods,
    ]);
    expect(["vectors", v.vectors.length]).toEqual(["vectors", EXPECTED_COUNTS.vectors]);
  });

  it("maps every operation to a method, or lists it as not wrapped", () => {
    const v = loadVectors();
    for (const operation of Object.keys(v.methods)) {
      const covered = operation in ADAPTERS || operation in NOT_WRAPPED;
      expect([operation, covered]).toEqual([operation, true]);
      expect([operation, operation in ADAPTERS && operation in NOT_WRAPPED]).toEqual([
        operation,
        false,
      ]);
    }
    // No adapter for an operation the vectors do not know: the table cannot drift.
    for (const operation of [...Object.keys(ADAPTERS), ...Object.keys(NOT_WRAPPED)]) {
      expect([operation, operation in v.methods]).toEqual([operation, true]);
    }
  });

  it("applies the next-page rule the shared operations map names", () => {
    const methods = loadVectors().methods;
    const rules = loadOperationRules();
    for (const [operation, adapter] of Object.entries(ADAPTERS)) {
      expect([operation, adapter.rule]).toEqual([operation, rules[operation]?.rule]);
      expect([operation, adapter.rule]).toEqual([operation, methods[operation]?.rule]);
    }
  });

  it("sends exactly the vector's query or body, or rejects before sending", async () => {
    const v = loadVectors();
    for (const vector of v.vectors) {
      if (vector.operation in NOT_WRAPPED) {
        continue;
      }
      const spec = v.methods[vector.operation];
      const adapter = ADAPTERS[vector.operation];
      const transport = new RecordingTransport({});
      const services = pagedServices(transport);
      const options = mapOptions(vector.operation, adapter, vector.options);
      const label = `${vector.operation}: ${vector.description}`;

      let error: unknown;
      try {
        await adapter.call(services, options, vector.path_params);
      } catch (err) {
        error = err;
      }

      if (vector.expect_error) {
        expect([label, error instanceof ValidationError, transport.requests.length]).toEqual([
          label,
          true,
          0,
        ]);
        continue;
      }

      expect([label, error]).toEqual([label, undefined]);
      expect([label, transport.requests.length]).toEqual([label, 1]);
      const request = transport.requests[0];
      expect([label, request.method, request.path]).toEqual([
        label,
        spec.method,
        expectedPath(spec.path, vector.path_params),
      ]);
      if (vector.expect?.query !== undefined) {
        expect([label, sortedPairs(request.query)]).toEqual([
          label,
          sortedPairs(vector.expect.query),
        ]);
        // Decoding hides an unescaped "/" or "=", so check the raw wire form too.
        for (const [, value] of vector.expect.query) {
          expect([label, request.rawQuery.includes(encodeURIComponent(value))]).toEqual([
            label,
            true,
          ]);
        }
      }
      if (vector.expect?.body !== undefined) {
        expect([label, request.query]).toEqual([label, []]);
        expect([label, request.body]).toEqual([label, vector.expect.body]);
      }
    }
  });
});
