/**
 * Cross-SDK behaviour vectors for the verification primitives.
 *
 * These are the invariants that had drifted apart before: which (blockchain,
 * network) pair selects the governance rules, and whether a hash is covered by a
 * signature. Both are pure input → outcome mappings, so all four SDKs assert them
 * from one file rather than from four hand-maintained copies — the arrangement that
 * let Java's hash-coverage go non-constant-time while its peers did not.
 *
 * To add a case: append to the shared file, bump the matching count there, and
 * consume it in all four suites.
 */

import { readFileSync, existsSync } from "fs";
import * as path from "path";

import { IntegrityError } from "../../../src/errors";
import { calculateHexHash } from "../../../src/crypto/hashing";
import {
  computeLegacyPayloadVariants,
  resolveRuleKey,
  verifyHashCoverage,
} from "../../../src/helpers/whitelist-hash-helper";
import {
  createEmptyRulesContainer,
  findAddressWhitelistingRuleCandidates,
  findContractAddressWhitelistingRuleCandidates,
} from "../../../src/models/governance-rules";
import { rulesetVerificationKey } from "../../../src/services/governance-rule-service";
import type { GovernanceRules } from "../../../src/models/governance-rules";

// tests/unit/helpers -> tests/unit -> tests -> <sdk> -> <repo root>
const VECTORS_PATH = path.resolve(
  __dirname,
  "../../../../scripts/resources/verification-behaviour-vectors.json"
);

interface RuleKeyCase {
  description: string;
  payload_as_string: string;
  dto_blockchain: string;
  dto_network: string;
  expect: "ok" | "error";
  blockchain?: string;
  network?: string;
}

interface HashCoverageCase {
  description: string;
  signatures: string[][];
  hash: string;
  expect: boolean;
}

interface ContainsHashCase {
  description: string;
  hashes: string[];
  hash: string;
  expect: boolean;
}

interface MemoKeyRuleset {
  rules_container: string;
  signatures: { user_id: string; signature: string }[];
}

interface MemoKeyCase {
  description: string;
  a: MemoKeyRuleset;
  b: MemoKeyRuleset;
  expect: "distinct" | "same";
}

interface LegacyHashCase {
  description: string;
  signed_payload: string;
  delivered_payload: string;
  expect: "match" | "no_match";
  expect_matched_payload?: string;
}

interface RuleTier {
  blockchain: string;
  network: string;
}

interface RuleTierCandidatesCase {
  description: string;
  rules: RuleTier[];
  blockchain: string;
  expect_candidates: RuleTier[];
}

interface Vectors {
  counts: Record<string, number>;
  legacy_hash: LegacyHashCase[];
  rule_tier_candidates: RuleTierCandidatesCase[];
  rule_key: RuleKeyCase[];
  hash_coverage: HashCoverageCase[];
  contains_hash: ContainsHashCase[];
  memo_key: MemoKeyCase[];
}

/** Builds a GovernanceRules from a memo_key vector side. */
function toRuleset(spec: MemoKeyRuleset): GovernanceRules {
  return {
    rulesContainer: spec.rules_container,
    rulesSignatures: spec.signatures.map((sig) => ({
      userId: sig.user_id,
      signature: sig.signature,
    })),
  } as GovernanceRules;
}

/**
 * Loaded lazily inside `it()`, never in a describe body: a throw during collection
 * becomes "Test suite failed to run", which drops the test count instead of showing
 * a red test.
 */
function loadVectors(): Vectors {
  if (!existsSync(VECTORS_PATH)) {
    throw new Error(
      `cannot read shared verification behaviour vectors ${VECTORS_PATH}`
    );
  }
  const vectors = JSON.parse(readFileSync(VECTORS_PATH, "utf-8")) as Vectors;

  // Counts are asserted so a case added to the shared file without being consumed
  // here fails loudly rather than being silently ignored by this SDK.
  for (const [key, expected] of Object.entries(vectors.counts)) {
    const actual = (vectors as unknown as Record<string, unknown[]>)[key].length;
    expect([key, actual]).toEqual([key, expected]);
  }
  return vectors;
}

describe("verification behaviour vectors", () => {
  it("resolves the governance rule key for every vector", () => {
    for (const c of loadVectors().rule_key) {
      if (c.expect === "error") {
        expect(() =>
          resolveRuleKey(c.payload_as_string, c.dto_blockchain, c.dto_network)
        ).toThrow(IntegrityError);
        continue;
      }
      const got = resolveRuleKey(
        c.payload_as_string,
        c.dto_blockchain,
        c.dto_network
      );
      // Assert as a tuple so a failing loop iteration names itself: Jest has no
      // per-assertion message argument.
      expect([c.description, got.blockchain, got.network]).toEqual([
        c.description,
        c.blockchain,
        c.network,
      ]);
    }
  });

  // The governance verification memo key decides whether ECDSA runs at all, so it MUST
  // be injective. An unprefixed concatenation leaves the boundary between the container
  // and the signature list uncommitted, and a response-controlling attacker can then
  // shift bytes across it to make a MODIFIED container inherit a genuine one's "already
  // verified" status — skipping signature verification on the document that carries
  // every HSM key the SDK trusts.
  it("keeps the memo key injective for every vector", () => {
    const vectors = loadVectors();

    // Non-vacuity: without a "distinct" pair the section would pass against a key
    // function that returns a constant.
    expect(vectors.memo_key.some((c) => c.expect === "distinct")).toBe(true);

    for (const c of vectors.memo_key) {
      const keyA = rulesetVerificationKey(toRuleset(c.a));
      const keyB = rulesetVerificationKey(toRuleset(c.b));
      expect([c.description, keyA !== undefined, keyB !== undefined]).toEqual([
        c.description,
        true,
        true,
      ]);
      // Assert as a tuple so a failing loop iteration names itself: Jest has no
      // per-assertion message argument.
      expect([c.description, keyA === keyB]).toEqual([
        c.description,
        c.expect === "same",
      ]);
    }
  });

  it("agrees on hash coverage for every vector", () => {
    for (const c of loadVectors().hash_coverage) {
      const signatures = c.signatures.map((hashes) => ({ hashes }));
      expect([c.description, verifyHashCoverage(c.hash, signatures)]).toEqual([
        c.description,
        c.expect,
      ]);
    }
  });

  it("agrees on single-signature hash containment for every vector", () => {
    // containsHash is private to each verifier; verifyHashCoverage over a single
    // signature is the public expression of the same question.
    for (const c of loadVectors().contains_hash) {
      expect([
        c.description,
        verifyHashCoverage(c.hash, [{ hashes: c.hashes }]),
      ]).toEqual([c.description, c.expect]);
    }
  });

  // ------------------------------------------------------------------------------
  // legacy_hash and rule_tier_candidates
  // ------------------------------------------------------------------------------
  //
  // Until 2026-09-10 both sections were consumed by the GO suite alone, even though the
  // file's `counts` block declares them and every loader asserts the counts. Asserting a
  // section's LENGTH proves the file is well formed; it does not prove the behaviour is
  // checked. Python shipped the single-tier lookup `rule_tier_candidates` exists to forbid
  // while three SDKs had the fix, and nothing went red.

  it("parses the payload the signature COVERED, for every legacy_hash vector", () => {
    // Asserts the PARSED PAYLOAD, not a hash — the legacy-strip injection moves no hash at
    // all, so `crypto-test-vectors.json` is structurally blind to it.
    for (const c of loadVectors().legacy_hash) {
      const coveredHash = calculateHexHash(c.signed_payload);

      let matchedPayload: string | undefined;
      if (calculateHexHash(c.delivered_payload) === coveredHash) {
        matchedPayload = c.delivered_payload;
      } else {
        for (const variant of computeLegacyPayloadVariants(c.delivered_payload)) {
          if (variant.hash === coveredHash) {
            matchedPayload = variant.payload;
            break;
          }
        }
      }

      if (c.expect === "no_match") {
        expect([c.description, matchedPayload]).toEqual([c.description, undefined]);
        continue;
      }
      expect([c.description, matchedPayload]).toEqual([
        c.description,
        c.expect_matched_payload,
      ]);
    }
  });

  it("returns every reachable rule tier, for every rule_tier_candidates vector", () => {
    // When the signed payload omits `network` there is no authenticated way to learn
    // whether that was legitimate, so every tier the unsigned DTO value could have
    // selected must be enforced. Too few lets the server pick the quorum; too many
    // rejects rows governance would accept.
    for (const c of loadVectors().rule_tier_candidates) {
      const container = {
        ...createEmptyRulesContainer(),
        addressWhitelistingRules: c.rules.map((r) => ({
          currency: r.blockchain,
          network: r.network,
          parallelThresholds: [],
          lines: [],
        })),
        contractAddressWhitelistingRules: c.rules.map((r) => ({
          blockchain: r.blockchain,
          network: r.network,
          parallelThresholds: [],
        })),
      };
      const expected = c.expect_candidates.map((e) => [e.blockchain, e.network]);

      expect([
        c.description,
        findAddressWhitelistingRuleCandidates(container, c.blockchain).map((r) => [
          r.currency ?? "",
          r.network ?? "",
        ]),
      ]).toEqual([c.description, expected]);

      // The asset peer, from the same vectors: the two families name the chain field
      // differently (`currency` vs `blockchain`), so a walk reading one name treats every
      // contract rule as a wildcard global default — fail-OPEN, the broadest tier.
      expect([
        c.description,
        findContractAddressWhitelistingRuleCandidates(container, c.blockchain).map((r) => [
          r.blockchain ?? "",
          r.network ?? "",
        ]),
      ]).toEqual([c.description, expected]);
    }
  });
});
