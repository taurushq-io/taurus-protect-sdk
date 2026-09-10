/**
 * Cross-SDK signed fixtures for the two verification thresholds.
 *
 * Step 2 is the tenant-wide SuperAdmin threshold; step 5 is the per-group one. Different
 * scopes, same counting rule, so one file gates both.
 *
 * The behaviour vectors cover everything expressible without key material; this file
 * covers what needs real signatures. Until it existed, the repo's own CLAUDE.md
 * recorded that NO automated cross-SDK gate covered the distinct-key rule, and the
 * five cases below lived as five hand-maintained copies in four suites.
 *
 * The fixture carries PUBLIC keys and signatures only — the gate verifies, it never
 * signs — so there is no private key material in the repo.
 */

import * as crypto from "crypto";
import { readFileSync, existsSync } from "fs";
import * as path from "path";

import { decodePublicKeyPem, encodePublicKeyPem } from "../../../src/crypto";
import { ConfigurationError, IntegrityError } from "../../../src/errors";
import { verifyGovernanceRulesSignatures } from "../../../src/helpers/signature-verifier";
import { WhitelistedAddressVerifier } from "../../../src/helpers/whitelisted-address-verifier";
import type {
  DecodedRulesContainer,
  GroupThreshold,
} from "../../../src/models/governance-rules";
import type { WhitelistSignatureEntry } from "../../../src/models/whitelisted-address";

// tests/unit/helpers -> tests/unit -> tests -> <sdk> -> <repo root>
const FIXTURES_PATH = path.resolve(
  __dirname,
  "../../../../scripts/resources/verification-signed-fixtures.json"
);

interface ThresholdCase {
  description: string;
  signatures: Array<{ user_id: string; signature: string }>;
  super_admin_keys_pem: string[];
  min_valid_signatures: number;
  expect: "ok" | "error";
}

interface GroupCase {
  description: string;
  metadata_hash: string;
  group_id: string;
  minimum_signatures: number;
  users: Array<{ user_id: string; public_key_pem: string }>;
  group_user_ids: string[];
  signatures: Array<{ user_id: string; signature: string; hashes: string[] }>;
  expect: "ok" | "error";
}

interface Fixtures {
  rules_container_base64: string;
  counts: Record<string, number>;
  superadmin_threshold: ThresholdCase[];
  group_threshold: GroupCase[];
}

/** Loaded inside `it()`, never in a describe body — a throw during collection drops
 * the test count instead of showing a red test. */
function loadFixtures(): Fixtures {
  if (!existsSync(FIXTURES_PATH)) {
    throw new Error(`cannot read shared signed fixtures ${FIXTURES_PATH}`);
  }
  const fixtures = JSON.parse(readFileSync(FIXTURES_PATH, "utf-8")) as Fixtures;

  for (const [key, expected] of Object.entries(fixtures.counts)) {
    const actual = (fixtures as unknown as Record<string, unknown[]>)[key].length;
    expect([key, actual]).toEqual([key, expected]);
  }
  expect(fixtures.rules_container_base64).toBeTruthy();
  return fixtures;
}

describe("signed fixtures: SuperAdmin threshold", () => {
  it("agrees with every recorded outcome", () => {
    const fixtures = loadFixtures();
    const container = Buffer.from(fixtures.rules_container_base64, "base64");

    for (const c of fixtures.superadmin_threshold) {
      const keys = c.super_admin_keys_pem.map(decodePublicKeyPem);
      const signatures = c.signatures.map((s) => ({ signature: s.signature }));

      let threw: unknown;
      try {
        verifyGovernanceRulesSignatures(
          container,
          signatures,
          keys,
          c.min_valid_signatures
        );
      } catch (error: unknown) {
        threw = error;
      }

      if (c.expect === "error") {
        // Assert as a tuple so a failing loop iteration names itself.
        expect([c.description, threw instanceof Error]).toEqual([c.description, true]);
        expect(
          threw instanceof IntegrityError || threw instanceof ConfigurationError
        ).toBe(true);
      } else {
        expect([c.description, threw]).toEqual([c.description, undefined]);
      }
    }
  });

  // Without a case where the entry COUNT meets the threshold but the distinct-key
  // count does not, the whole file would pass against an implementation that counts
  // entries — the regression these fixtures exist to catch.
  it("contains a case that distinguishes distinct-key counting from entry counting", () => {
    const discriminating = loadFixtures().superadmin_threshold.some(
      (c) =>
        c.expect === "error" &&
        c.min_valid_signatures > 0 &&
        c.signatures.length >= c.min_valid_signatures
    );
    expect(discriminating).toBe(true);
  });
});

// Step 5 — the per-group threshold. A different threshold from the SuperAdmin one above
// (per group, not tenant-wide) but the same counting rule, so it is gated from the same
// file: the four SDKs cannot drift on one without drifting on the other.
describe("signed fixtures: group threshold", () => {
  it("agrees with every recorded outcome", () => {
    const fixtures = loadFixtures();

    for (const c of fixtures.group_threshold) {
      const rules = {
        users: c.users.map((u) => ({
          id: u.user_id,
          publicKeyPem: u.public_key_pem,
          roles: ["USER"],
        })),
        groups: [{ id: c.group_id, userIds: c.group_user_ids }],
      } as unknown as DecodedRulesContainer;

      const signatures = c.signatures.map(
        (s) =>
          ({
            hashes: s.hashes,
            userSignature: { userId: s.user_id, signature: s.signature },
          }) as WhitelistSignatureEntry
      );

      const precomputed = new Map<number, Buffer>();
      signatures.forEach((sig, idx) => {
        precomputed.set(idx, Buffer.from(JSON.stringify(sig.hashes), "utf-8"));
      });

      const verifier = new WhitelistedAddressVerifier({
        superAdminKeysPem: [encodePublicKeyPem(
          crypto.generateKeyPairSync("ec", { namedCurve: "prime256v1" }).publicKey
        )],
        minValidSignatures: 1,
      });

      // verifyGroupThreshold is private; the walk is what is under test.
      const failure = (
        verifier as unknown as {
          verifyGroupThreshold: (
            t: GroupThreshold,
            r: DecodedRulesContainer,
            s: WhitelistSignatureEntry[],
            h: string,
            p: Map<number, Buffer>,
            k: Map<string, crypto.KeyObject>
          ) => string | null;
        }
      ).verifyGroupThreshold(
        {
          groupId: c.group_id,
          minimumSignatures: c.minimum_signatures,
          threshold: 0,
        } as GroupThreshold,
        rules,
        signatures,
        c.metadata_hash,
        precomputed,
        new Map()
      );

      // Assert as a tuple so a failing loop iteration names itself.
      expect([c.description, failure === null]).toEqual([
        c.description,
        c.expect === "ok",
      ]);
    }
  });

  // Same non-vacuity guard as the SuperAdmin section, applied per group.
  it("contains a case that distinguishes distinct-signer counting from entry counting", () => {
    const discriminating = loadFixtures().group_threshold.some((c) => {
      if (c.expect !== "error" || c.minimum_signatures <= 0) {
        return false;
      }
      const members = new Set(c.group_user_ids);
      const inGroup = c.signatures.filter((s) => members.has(s.user_id)).length;
      return inGroup >= c.minimum_signatures;
    });
    expect(discriminating).toBe(true);
  });
});
