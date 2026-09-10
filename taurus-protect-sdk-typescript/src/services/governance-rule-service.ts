/**
 * Governance rule service for Taurus-PROTECT SDK.
 *
 * Provides methods for retrieving governance rules and managing rule proposals.
 */

import type { KeyObject } from "crypto";
import { ConfigurationError, IntegrityError, ValidationError } from "../errors";
import { signData } from "../crypto/signing";
import { verifyGovernanceRulesSignatures } from "../helpers/signature-verifier";
import type { GovernanceRulesApi } from "../internal/openapi/apis/GovernanceRulesApi";
import {
  governanceRulesFromDto,
  governanceRulesArrayFromDto,
  rulesContainerFromBase64,
  superAdminPublicKeysFromDto,
} from "../mappers/governance-rules";
import { rulesContainerToBase64 } from "../mappers/protobuf-rules-container-encode";
import { createHash, timingSafeEqual, type Hash } from "crypto";
import { strictBase64Decode } from "../helpers/strict-base64";
import type {
  DecodedRulesContainer,
  ExcludedRuleset,
  GovernanceRules,
  GovernanceRulesHistoryResult,
  ListGovernanceRulesHistoryOptions,
  SuperAdminPublicKey,
} from "../models/governance-rules";
import { BaseService } from "./base";

/**
 * Configuration options for GovernanceRuleService.
 */
/** Bounds the verification memo. */
const MAX_VERIFIED_RULESETS = 256;

/**
 * Writes an 8-byte big-endian length, then the bytes, so a sequence of fields cannot be
 * re-partitioned into a different sequence that encodes identically.
 */
function updateLengthPrefixed(digest: Hash, data: Buffer): void {
  const length = Buffer.alloc(8);
  length.writeBigUInt64BE(BigInt(data.length));
  digest.update(length);
  digest.update(data);
}

/**
 * Identifies the exact document verification cleared: the DECODED container bytes -- the
 * same bytes {@link verifyGovernanceRulesSignatures} checks the signatures over -- plus
 * every signature over them. A container whose signature set changed is a memo MISS rather
 * than a stale hit; signature order is server-controlled, so it is sorted out of the key.
 *
 * Every field is LENGTH-PREFIXED, and that is load-bearing. A plain concatenation is not
 * injective: with `sha256(container || sig || ...)` the boundary between the container and
 * the signature list is not committed, so a response-controlling attacker can shift bytes
 * across it and make a MODIFIED container collide with a genuine one's key -- inheriting
 * its "already verified" status and skipping ECDSA entirely. Interleaving a NUL separator
 * does NOT fix this; only length prefixes do.
 *
 * `userId` is deliberately excluded: verification never reads it (it consumes only
 * `signature`), so a server-controlled field that verification ignores must not be able to
 * influence a verification-SKIP decision.
 *
 * Keying on the decoded bytes rather than the base64 text also removes base64 leniency as a
 * lever, and guarantees the key cannot identify something other than what was verified.
 *
 * @returns the hex digest, or undefined when the container does not decode -- in which case
 *   there is nothing worth memoising and the caller must fall through to verification.
 *
 * Exported for the cross-SDK memo-key vector gate only. NOT re-exported from the package
 * barrel, so it is not public API -- same convention as `attestVerified`.
 */
export function rulesetVerificationKey(rules: GovernanceRules): string | undefined {
  let data: Buffer;
  try {
    data = strictBase64Decode(rules.rulesContainer ?? "");
  } catch {
    // A container we cannot decode has no stable identity to memoise on.
    return undefined;
  }

  const digest = createHash("sha256");
  updateLengthPrefixed(digest, data);

  const signatures = (rules.rulesSignatures ?? []).map((sig) => sig.signature ?? "").sort();
  const count = Buffer.alloc(8);
  count.writeBigUInt64BE(BigInt(signatures.length));
  digest.update(count);
  for (const signature of signatures) {
    updateLengthPrefixed(digest, Buffer.from(signature, "utf8"));
  }
  return digest.digest("hex");
}

export interface GovernanceRuleServiceConfig {
  /**
   * SuperAdmin public keys for signature verification (KeyObject instances).
   * REQUIRED: verification is mandatory on every governance read.
   */
  readonly superAdminKeys: KeyObject[];
  /** Minimum number of valid signatures required. Must be positive. */
  readonly minValidSignatures: number;
}

/**
 * Service for managing governance rules.
 *
 * Governance rules define the approval workflows and policies for
 * transaction requests and address whitelisting.
 *
 * @example
 * ```typescript
 * // Get current rules
 * const rules = await governanceRuleService.getRules();
 * if (rules) {
 *   console.log(`Rules locked: ${rules.locked}`);
 * }
 *
 * // Get rules by ID
 * const historical = await governanceRuleService.getRulesById("rules-123");
 *
 * // Get rules history
 * const history = await governanceRuleService.getRulesHistory({ limit: 10 });
 * for (const rules of history.items) {
 *   console.log(`Created: ${rules.creationDate}`);
 * }
 *
 * // Get decoded rules container
 * const decoded = await governanceRuleService.getDecodedRulesContainer();
 * console.log(`Users: ${decoded.users.length}`);
 * ```
 */
export class GovernanceRuleService extends BaseService {
  private readonly governanceRulesApi: GovernanceRulesApi;
  private readonly superAdminKeys: KeyObject[];
  private readonly minValidSignatures: number;
  /** Rulesets whose signatures already checked out; see verifiedRuleset. */
  private readonly verified = new Set<string>();

  /**
   * Creates a new GovernanceRuleService instance.
   *
   * @param governanceRulesApi - The GovernanceRulesApi instance from the OpenAPI client
   * @param config - Optional configuration for signature verification
   */
  constructor(
    governanceRulesApi: GovernanceRulesApi,
    config: GovernanceRuleServiceConfig
  ) {
    super();
    this.governanceRulesApi = governanceRulesApi;
    this.superAdminKeys = config.superAdminKeys ?? [];
    this.minValidSignatures = config.minValidSignatures ?? 0;

    // There is no "explicitly unverified service" any more. `config` used to be
    // optional and defaulted to zero keys / zero threshold, which every read then
    // gated on — so the class exported from the package barrel could be constructed
    // to return the UNVERIFIED trust root, silently and with no error. Java
    // (checkArgument), Python (ValueError) and Go (no keyless constructor) all refuse
    // it; this now does too.
    if (this.superAdminKeys.length === 0) {
      throw new ConfigurationError(
        "superAdminKeys are required: governance verification is mandatory"
      );
    }
    if (this.minValidSignatures <= 0) {
      throw new ConfigurationError(
        `minValidSignatures must be positive, got ${this.minValidSignatures}`
      );
    }
  }

  /**
   * Verifies a ruleset's SuperAdmin signatures, memoising the outcome so the same
   * document is not re-verified on every read.
   *
   *   getRules / getRulesById / getRulesHistory -> verifiedRuleset
   *                                                 |- memo hit -> return
   *                                                 +- miss -> ECDSA -> memoise
   *
   * The container itself is already cached for the address/asset/price paths by
   * `RulesContainerCache`, which fetches through `getDecodedRulesContainer`. This memo
   * covers only the reads that return the RAW ruleset. Do not add a third cache of the
   * same bytes.
   *
   * Only SUCCESSES are memoised: a failure must resurface its error on every call.
   */
  private verifiedRuleset(rules: GovernanceRules): GovernanceRules {
    const key = rulesetVerificationKey(rules);
    if (key !== undefined && this.verified.has(key)) {
      return rules;
    }
    this.verifyGovernanceRules(rules);
    // An undecodable container has no stable identity to memoise on; verification above
    // has already had the final say.
    if (key === undefined) {
      return rules;
    }
    // Bounded so walking a long history cannot grow it without limit. Containers
    // change rarely, so a plain reset beats LRU bookkeeping.
    if (this.verified.size >= MAX_VERIFIED_RULESETS) {
      this.verified.clear();
    }
    this.verified.add(key);
    return rules;
  }

  /**
   * Gets the currently enforced governance rules.
   *
   * @returns The governance rules, or undefined if not available
   * @throws {@link IntegrityError} If signature verification fails
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const rules = await governanceRuleService.getRules();
   * if (rules) {
   *   console.log(`Rules locked: ${rules.locked}`);
   *   console.log(`Signatures: ${rules.rulesSignatures.length}`);
   * }
   * ```
   */
  async getRules(): Promise<GovernanceRules | undefined> {
    return this.execute(async () => {
      const response = await this.governanceRulesApi.ruleServiceGetRules();
      const rules = governanceRulesFromDto(response.result);

      if (rules) {
        this.verifiedRuleset(rules);
      }

      return rules;
    });
  }

  /**
   * Gets a governance ruleset by its ID.
   *
   * @param rulesId - The ruleset ID
   * @returns The governance rules, or undefined if not found
   * @throws {@link ValidationError} If rulesId is invalid
   * @throws {@link IntegrityError} If signature verification fails
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const rules = await governanceRuleService.getRulesById("rules-123");
   * if (rules) {
   *   console.log(`Created: ${rules.creationDate}`);
   * }
   * ```
   */
  async getRulesById(rulesId: string): Promise<GovernanceRules | undefined> {
    if (!rulesId || rulesId.trim() === "") {
      throw new ValidationError("rulesId is required");
    }

    return this.execute(async () => {
      const response = await this.governanceRulesApi.ruleServiceGetRulesByID({
        id: rulesId,
      });
      const rules = governanceRulesFromDto(response.result);

      if (rules) {
        this.verifiedRuleset(rules);
      }

      return rules;
    });
  }

  /**
   * Gets the proposed governance rules.
   *
   * Requires SuperAdmin or SuperAdminReadOnly role.
   *
   * @returns The proposed governance rules, or undefined if not available
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const proposal = await governanceRuleService.getRulesProposal();
   * if (proposal) {
   *   console.log(`Proposal locked: ${proposal.locked}`);
   * }
   * ```
   */
  async getRulesProposal(): Promise<GovernanceRules | undefined> {
    return this.execute(async () => {
      const response =
        await this.governanceRulesApi.ruleServiceGetRulesProposal();
      // Proposal rules are not verified (they're not yet enforced)
      return governanceRulesFromDto(response.result);
    });
  }

  /**
   * Submits a rules container as a governance proposal (SuperAdmin only).
   *
   * The container is encoded to the wire format internally; the server-controlled
   * `enforcedRulesHash` and `timestamp` fields are stripped (the server recomputes
   * them and rejects submissions asserting stale values). The endpoint returns no
   * body — callers that need the persisted proposal should call {@link getProposal}
   * (note: rules reads are cached server-side, so an immediate read-back may be stale).
   *
   * @param container - The rules container to submit
   * @throws {@link APIError} If the API request fails
   */
  async updateRulesProposal(container: DecodedRulesContainer): Promise<void> {
    const rulesContainer = rulesContainerToBase64(container);
    await this.execute(async () => {
      await this.governanceRulesApi.ruleServiceUpdateRulesProposal({ body: { rulesContainer } });
    });
  }

  /**
   * Signs the pending proposal's rules container with a SuperAdmin private key
   * and submits the approval (SuperAdmin only).
   *
   * `expectedContainerHash` PINS the content being approved: it must be the
   * {@link proposalContainerHash} of the proposal the caller actually reviewed. The
   * proposal is re-fetched here, and if its container does not match that digest the
   * call aborts WITHOUT signing.
   *
   * That pin is the security control. Without it a server able to shape responses could
   * serve the benign proposal to the review call and a different container to the
   * re-fetch inside this method, obtaining a GENUINE SuperAdmin signature over bytes of
   * its choosing — and such a container then verifies clean everywhere, including in
   * independent and air-gapped verifiers that never trusted that server. Nothing else
   * binds the two calls: the wire format carries no proposal id, hash or version, and a
   * pending proposal has no signatures to verify against, so re-verification cannot
   * substitute for pinning.
   *
   * The signature is computed over the decoded bytes of the pending proposal's rules
   * container (SHA-256 + P-256 ECDSA, base64 raw r||s).
   *
   * @param privateKey - SuperAdmin P-256 private key
   * @param comment - Approval comment
   * @param expectedContainerHash - {@link proposalContainerHash} of the reviewed proposal
   * @throws {@link ValidationError} If no proposal is pending, or the pin is missing
   * @throws {@link IntegrityError} If the pending container differs from the pin
   * @throws {@link APIError} If the API request fails
   */
  async approveRulesProposal(
    privateKey: KeyObject,
    comment: string,
    expectedContainerHash: string
  ): Promise<void> {
    if (!expectedContainerHash) {
      throw new ValidationError(
        "expectedContainerHash is required: it pins the approval to the container you " +
          "reviewed (see proposalContainerHash)"
      );
    }

    const proposal = await this.getRulesProposal();
    if (!proposal || !proposal.rulesContainer) {
      throw new ValidationError("No pending rules proposal to approve");
    }

    const actualHash = this.proposalContainerHash(proposal);
    // Constant-time, matching how every other hash comparison in this SDK is done.
    const expected = Buffer.from(expectedContainerHash, "utf8");
    const actual = Buffer.from(actualHash, "utf8");
    const matches =
      expected.length === actual.length && timingSafeEqual(expected, actual);
    if (!matches) {
      throw new IntegrityError(
        `Refusing to sign: the pending rules proposal changed since it was reviewed ` +
          `(reviewed ${expectedContainerHash}, pending ${actualHash})`
      );
    }

    const data = new Uint8Array(strictBase64Decode(proposal.rulesContainer));
    const signature = signData(privateKey, data);
    await this.execute(async () => {
      await this.governanceRulesApi.ruleServiceApproveRulesProposal({ body: { signature, comment } });
    });
  }

  /**
   * Rejects the pending rules proposal with a comment (SuperAdmin only).
   *
   * @param comment - Rejection comment
   * @throws {@link APIError} If the API request fails
   */
  async rejectRulesProposal(comment: string): Promise<void> {
    await this.execute(async () => {
      await this.governanceRulesApi.ruleServiceRejectRulesProposal({ body: { comment } });
    });
  }

  /**
   * Gets the history of governance rules.
   *
   * Returns a paginated list of governance rule sets in reverse chronological
   * order (most recent first).
   *
   * @param options - Optional pagination options
   * @returns History result with rules and next cursor
   * @throws {@link ValidationError} If options are invalid
   * @throws {@link IntegrityError} If signature verification fails
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // Get first page
   * const page1 = await governanceRuleService.getRulesHistory({ limit: 10 });
   * console.log(`Found ${page1.items.length} rules`);
   *
   * // Get next page
   * if (page1.nextCursor) {
   *   const page2 = await governanceRuleService.getRulesHistory({
   *     limit: 10,
   *     cursor: page1.nextCursor,
   *   });
   * }
   * ```
   */
  async getRulesHistory(
    options?: ListGovernanceRulesHistoryOptions
  ): Promise<GovernanceRulesHistoryResult> {
    const limit = options?.limit ?? 50;

    if (limit <= 0) {
      throw new ValidationError("limit must be positive");
    }

    return this.execute(async () => {
      const response = await this.governanceRulesApi.ruleServiceGetRulesHistory(
        {
          limit: String(limit),
          cursor: options?.cursor,
        }
      );

      // Every historical ruleset is a past ENFORCED document, so each is verified. The
      // memo makes this one ECDSA pass per distinct container, not one per row.
      //
      //   entries -> verifiedRuleset -+- ok    -> kept
      //                               +- error -> excluded + named in the result
      //
      // LENIENT, unlike the whitelist lists, and it does NOT throw when nothing
      // survives: a SuperAdmin key rotation makes every pre-rotation ruleset
      // unverifiable, so aborting would deny the whole audit trail from then on.
      const kept: GovernanceRules[] = [];
      const excludedUnverified: ExcludedRuleset[] = [];
      for (const rule of governanceRulesArrayFromDto(response.result)) {
        try {
          kept.push(this.verifiedRuleset(rule));
        } catch (err) {
          excludedUnverified.push({
            creationDate: rule.creationDate,
            reason: err instanceof Error ? err.message : String(err),
          });
        }
      }

      return {
        items: kept,
        excludedUnverified,
        nextCursor: response.cursor,
      };
    });
  }

  /**
   * Gets the SuperAdmin public keys the server has configured.
   *
   * Useful for confirming that the keys this client verifies against are the ones the
   * server actually enforces — a mismatch is why signature verification fails with an
   * otherwise valid container.
   *
   * @returns The configured SuperAdmin public keys
   * @throws {@link APIError} If the API request fails
   *
   * @example
   * ```typescript
   * const keys = await governanceRuleService.getPublicKeys();
   * for (const key of keys) {
   *   console.log(`${key.userId}: ${key.publicKey}`);
   * }
   * ```
   */
  async getPublicKeys(): Promise<SuperAdminPublicKey[]> {
    return this.execute(async () => {
      const response = await this.governanceRulesApi.ruleServiceGetPublicKeys();
      return superAdminPublicKeysFromDto(response.publicKeys);
    });
  }

  /**
   * Gets the decoded rules container from the current governance rules.
   *
   * This method fetches the current rules and decodes the rules container.
   * The decoded container contains users, groups, and whitelisting rules.
   *
   * @returns The decoded rules container
   * @throws {@link IntegrityError} If signature verification fails or decoding fails
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const container = await governanceRuleService.getDecodedRulesContainer();
   * console.log(`Users: ${container.users.length}`);
   * console.log(`Groups: ${container.groups.length}`);
   * ```
   */
  async getDecodedRulesContainer(): Promise<DecodedRulesContainer> {
    const rules = await this.getRules();

    if (!rules) {
      throw new IntegrityError("No governance rules available");
    }

    if (!rules.rulesContainer) {
      throw new IntegrityError("Rules container is empty");
    }

    return rulesContainerFromBase64(rules.rulesContainer);
  }

  /**
   * Decodes a PENDING rules proposal so a SuperAdmin can inspect it before approving.
   *
   * UNVERIFIED BY DESIGN. A pending proposal legitimately carries 0..N signatures —
   * signing IS the approval step — so there is no threshold to check yet, and
   * {@link getDecodedRulesContainer} (which verifies unconditionally) always fails on
   * one. That is why this is a separate, explicitly-named entry point rather than a
   * flag: the absence of verification has to be visible at the call site.
   *
   * Pair it with {@link proposalContainerHash} and pass that digest to
   * {@link approveRulesProposal}, so the bytes signed are the bytes reviewed.
   *
   * @param rules - The pending proposal, from {@link getRulesProposal}
   * @returns The decoded rules container, NOT verified
   * @throws {@link IntegrityError} If decoding fails
   *
   * @example
   * ```typescript
   * const proposal = await governanceRuleService.getRulesProposal();
   * if (proposal) {
   *   const container = governanceRuleService.decodeProposalForReview(proposal);
   *   const pin = governanceRuleService.proposalContainerHash(proposal);
   *   await governanceRuleService.approveRulesProposal(key, "lgtm", pin);
   * }
   * ```
   */
  decodeProposalForReview(rules: GovernanceRules): DecodedRulesContainer {
    if (!rules.rulesContainer) {
      throw new IntegrityError("Rules container is empty");
    }

    return rulesContainerFromBase64(rules.rulesContainer);
  }

  /**
   * Returns the canonical SHA-256 hex digest of a ruleset's decoded rules container.
   *
   * This is the value to pass as `expectedContainerHash` to
   * {@link approveRulesProposal}, and it is what pins the approval to reviewed content.
   *
   * It digests the DECODED bytes, not the base64 text: base64 has multiple encodings of
   * the same bytes, so hashing the text would let a re-encoded but byte-identical
   * container read as a mismatch. Those decoded bytes are also exactly what gets signed.
   *
   * @param rules - The ruleset whose container to digest
   * @returns The hex digest
   * @throws {@link ValidationError} If the ruleset carries no container
   * @throws {@link IntegrityError} If the container is not valid base64
   */
  proposalContainerHash(rules: GovernanceRules): string {
    if (!rules.rulesContainer) {
      throw new ValidationError("Ruleset carries no rules container");
    }
    let data: Buffer;
    try {
      data = strictBase64Decode(rules.rulesContainer);
    } catch (error) {
      throw new IntegrityError(
        `Failed to decode rules container: ${error instanceof Error ? error.message : "unknown error"}`
      );
    }
    return createHash("sha256").update(data).digest("hex");
  }

  /**
   * Verifies that governance rules have enough valid SuperAdmin signatures.
   *
   * This method is called automatically when minValidSignatures > 0.
   * For manual verification, call this method directly.
   *
   * @param rules - The governance rules to verify
   * @returns The verified rules
   * @throws {@link IntegrityError} If verification fails
   *
   * @example
   * ```typescript
   * const rules = await governanceRuleService.getRules();
   * if (rules) {
   *   governanceRuleService.verifyGovernanceRules(rules);
   *   console.log("Rules verified successfully");
   * }
   * ```
   */
  verifyGovernanceRules(rules: GovernanceRules): GovernanceRules {
    // No `minValidSignatures <= 0` escape hatch: a caller who reaches this method
    // is asking for verification, and returning the rules unverified would answer
    // "verified" to a question that was never asked. The constructor rejects the
    // keys-with-zero-threshold combination, so reaching here with a zero threshold
    // means no keys either, which the guard below catches.
    if (!rules.rulesContainer) {
      throw new IntegrityError("Rules container is empty, cannot verify");
    }

    if (!rules.rulesSignatures || rules.rulesSignatures.length === 0) {
      throw new IntegrityError("No signatures found on rules");
    }

    if (this.superAdminKeys.length === 0) {
      throw new IntegrityError(
        "No SuperAdmin keys configured for verification"
      );
    }

    // Decode the rules container to get the signed data (raw bytes)
    const rulesData = Buffer.from(rules.rulesContainer, "base64");

    verifyGovernanceRulesSignatures(
      rulesData,
      rules.rulesSignatures,
      this.superAdminKeys,
      this.minValidSignatures
    );

    // Returning the verified rules matches Java and Python, so a caller can chain
    // verification into an assignment instead of relying on the throw alone.
    return rules;
  }

  /**
   * Gets the configured SuperAdmin public keys.
   */
  get configuredSuperAdminKeys(): readonly KeyObject[] {
    return this.superAdminKeys;
  }

  /**
   * Gets the configured minimum valid signatures requirement.
   */
  get configuredMinValidSignatures(): number {
    return this.minValidSignatures;
  }
}
