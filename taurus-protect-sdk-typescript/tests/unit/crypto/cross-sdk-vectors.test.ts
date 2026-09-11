/**
 * Cross-SDK cryptographic test vectors.
 *
 * These tests verify that TypeScript SDK cryptographic functions produce
 * identical results to the Java, Go, and Python SDKs. All SDKs read the
 * same test vectors from docs/test-vectors/crypto-test-vectors.json.
 */

import * as fs from "fs";
import * as path from "path";
import { calculateSignedHeader } from "../../../src/crypto/tpv1";
import {
  calculateHexHash,
  calculateBase64Hmac,
  constantTimeCompare,
} from "../../../src/crypto/hashing";
import {
  computeLegacyHashes,
  computeAssetLegacyHashes,
} from "../../../src/helpers/whitelist-hash-helper";

const VECTORS_PATH = path.join(
  __dirname,
  "..",
  "..",
  "..",
  "..",
  "docs",
  "test-vectors",
  "crypto-test-vectors.json"
);

interface TestVectors {
  vectors: {
    hex_hash: Array<{
      description: string;
      input: string;
      expected: string;
    }>;
    hmac_sha256: Array<{
      description: string;
      key_hex: string;
      data: string;
      expected_base64: string;
    }>;
    constant_time_compare: Array<{
      description: string;
      a: string;
      b: string;
      expected: boolean;
    }>;
    legacy_hash_address: Array<{
      description: string;
      payload: string;
      original_hash: string;
      expected_without_contract_type?: string;
      expected_without_labels?: string;
      expected_without_both?: string;
      expected_legacy_count: number;
    }>;
    legacy_hash_asset: Array<{
      description: string;
      payload: string;
      original_hash: string;
      expected_without_is_nft?: string;
      expected_without_kind_type?: string;
      expected_without_both?: string;
      expected_legacy_count: number;
    }>;
    canonical_string: {
      secret_hex: string;
      count: number;
      cases: Array<{
        description: string;
        api_key: string;
        nonce: string;
        timestamp: number;
        method: string;
        host: string;
        path: string;
        query: string;
        content_type: string;
        body: string;
        expected_message: string;
        expected_signature: string;
      }>;
    };
  };
}

const testData: TestVectors = JSON.parse(
  fs.readFileSync(VECTORS_PATH, "utf8")
);
const vectors = testData.vectors;

describe("Cross-SDK SHA-256 Hex Hash", () => {
  it.each(vectors.hex_hash)(
    "should match for: $description",
    ({ input, expected }) => {
      expect(calculateHexHash(input)).toBe(expected);
    }
  );
});

describe("Cross-SDK HMAC-SHA256", () => {
  it.each(vectors.hmac_sha256)(
    "should match for: $description",
    ({ key_hex, data, expected_base64 }) => {
      const key = Buffer.from(key_hex, "hex");
      expect(calculateBase64Hmac(key, data)).toBe(expected_base64);
    }
  );
});

describe("Cross-SDK Constant-Time Compare", () => {
  it.each(vectors.constant_time_compare)(
    "should match for: $description",
    ({ a, b, expected }) => {
      expect(constantTimeCompare(a, b)).toBe(expected);
    }
  );
});

describe("Cross-SDK Legacy Address Hash", () => {
  it.each(vectors.legacy_hash_address)(
    "original hash should match for: $description",
    ({ payload, original_hash }) => {
      expect(calculateHexHash(payload)).toBe(original_hash);
    }
  );

  it.each(vectors.legacy_hash_address)(
    "legacy hashes should match for: $description",
    ({
      payload,
      expected_without_contract_type,
      expected_without_labels,
      expected_without_both,
      expected_legacy_count,
    }) => {
      const legacyHashes = computeLegacyHashes(payload);
      expect(legacyHashes.length).toBe(expected_legacy_count);

      if (expected_legacy_count > 0) {
        expect(legacyHashes).toContain(expected_without_contract_type);
        expect(legacyHashes).toContain(expected_without_labels);
        expect(legacyHashes).toContain(expected_without_both);
      }
    }
  );
});

describe("Cross-SDK Legacy Asset Hash", () => {
  it.each(vectors.legacy_hash_asset)(
    "original hash should match for: $description",
    ({ payload, original_hash }) => {
      expect(calculateHexHash(payload)).toBe(original_hash);
    }
  );

  it.each(vectors.legacy_hash_asset)(
    "legacy hashes should match for: $description",
    ({
      payload,
      expected_without_is_nft,
      expected_without_kind_type,
      expected_without_both,
      expected_legacy_count,
    }) => {
      const legacyHashes = computeAssetLegacyHashes(payload);
      expect(legacyHashes.length).toBe(expected_legacy_count);

      if (expected_legacy_count > 0) {
        expect(legacyHashes).toContain(expected_without_is_nft);
        expect(legacyHashes).toContain(expected_without_kind_type);
        expect(legacyHashes).toContain(expected_without_both);
      }
    }
  );
});

describe("Cross-SDK TPV1 Canonical String", () => {
  // Nothing pinned the canonical MESSAGE before 2026-09-10, and that is how a real interop
  // break shipped: the `hmac_sha256` group HMACs a hardcoded string that merely *looks*
  // like a canonical message, and no consumer routed through calculateSignedHeader -- so
  // Python and TypeScript upper-cased the HTTP method while Java and Go signed it verbatim,
  // leaving a caller who issued a lowercase `get` unable to authenticate against one of the
  // two families. This section was consumed by the Go suite alone until now.
  const group = vectors.canonical_string;

  it("declares the number of cases it carries", () => {
    // A case added and consumed by nobody must fail loudly.
    expect(group.cases.length).toBe(group.count);
  });

  it.each(group.cases)("signs the canonical message for: $description", (vec) => {
    // Asserting the SIGNATURE is what pins the MESSAGE: the secret is fixed, so a match
    // means the exact byte string was signed. Rebuilding the message in the test would
    // assert a copy of the implementation instead.
    const header = calculateSignedHeader(
      vec.api_key,
      Buffer.from(group.secret_hex, "hex"),
      vec.nonce,
      vec.timestamp,
      vec.method,
      vec.host,
      vec.path,
      vec.query || undefined,
      vec.content_type || undefined,
      vec.body || undefined
    );

    expect([vec.description, header.includes(`Signature=${vec.expected_signature}`)]).toEqual(
      [vec.description, true]
    );
  });
});
