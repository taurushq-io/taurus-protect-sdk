/**
 * Cross-SDK cryptographic test vectors.
 *
 * These tests verify that TypeScript SDK cryptographic functions produce
 * identical results to the Java, Go, and Python SDKs. All SDKs read the
 * same test vectors from docs/test-vectors/crypto-test-vectors.json.
 */

import * as fs from "fs";
import * as path from "path";
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
