/**
 * Governance rule service for Taurus-PROTECT SDK.
 *
 * Provides methods for retrieving governance rules and managing rule proposals.
 */

import type { KeyObject } from "crypto";
import { IntegrityError, ValidationError } from "../errors";
import { isValidSignature } from "../helpers/signature-verifier";
import type { GovernanceRulesApi } from "../internal/openapi/apis/GovernanceRulesApi";
import {
  governanceRulesFromDto,
  governanceRulesArrayFromDto,
  rulesContainerFromBase64,
} from "../mappers/governance-rules";
import type {
  DecodedRulesContainer,
  GovernanceRules,
  GovernanceRulesHistoryResult,
  ListGovernanceRulesHistoryOptions,
} from "../models/governance-rules";
import { BaseService } from "./base";

/**
 * Configuration options for GovernanceRuleService.
 */
export interface GovernanceRuleServiceConfig {
  /** SuperAdmin public keys for signature verification (KeyObject instances) */
  readonly superAdminKeys?: KeyObject[];
  /** Minimum number of valid signatures required (default: 0, disables verification) */
  readonly minValidSignatures?: number;
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
 * const rules = await governanceRuleService.get();
 * if (rules) {
 *   console.log(`Rules locked: ${rules.locked}`);
 * }
 *
 * // Get rules by ID
 * const historical = await governanceRuleService.getById("rules-123");
 *
 * // Get rules history
 * const history = await governanceRuleService.getHistory({ limit: 10 });
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

  /**
   * Creates a new GovernanceRuleService instance.
   *
   * @param governanceRulesApi - The GovernanceRulesApi instance from the OpenAPI client
   * @param config - Optional configuration for signature verification
   */
  constructor(
    governanceRulesApi: GovernanceRulesApi,
    config?: GovernanceRuleServiceConfig
  ) {
    super();
    this.governanceRulesApi = governanceRulesApi;
    this.superAdminKeys = config?.superAdminKeys ?? [];
    this.minValidSignatures = config?.minValidSignatures ?? 0;
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
   * const rules = await governanceRuleService.get();
   * if (rules) {
   *   console.log(`Rules locked: ${rules.locked}`);
   *   console.log(`Signatures: ${rules.rulesSignatures.length}`);
   * }
   * ```
   */
  async get(): Promise<GovernanceRules | undefined> {
    return this.execute(async () => {
      const response = await this.governanceRulesApi.ruleServiceGetRules();
      const rules = governanceRulesFromDto(response.result);

      if (rules && this.minValidSignatures > 0) {
        this.verifyGovernanceRules(rules);
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
   * const rules = await governanceRuleService.getById("rules-123");
   * if (rules) {
   *   console.log(`Created: ${rules.creationDate}`);
   * }
   * ```
   */
  async getById(rulesId: string): Promise<GovernanceRules | undefined> {
    if (!rulesId || rulesId.trim() === "") {
      throw new ValidationError("rulesId is required");
    }

    return this.execute(async () => {
      const response = await this.governanceRulesApi.ruleServiceGetRulesByID({
        id: rulesId,
      });
      const rules = governanceRulesFromDto(response.result);

      if (rules && this.minValidSignatures > 0) {
        this.verifyGovernanceRules(rules);
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
   * const proposal = await governanceRuleService.getProposal();
   * if (proposal) {
   *   console.log(`Proposal locked: ${proposal.locked}`);
   * }
   * ```
   */
  async getProposal(): Promise<GovernanceRules | undefined> {
    return this.execute(async () => {
      const response =
        await this.governanceRulesApi.ruleServiceGetRulesProposal();
      // Proposal rules are not verified (they're not yet enforced)
      return governanceRulesFromDto(response.result);
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
   * const page1 = await governanceRuleService.getHistory({ limit: 10 });
   * console.log(`Found ${page1.items.length} rules`);
   *
   * // Get next page
   * if (page1.nextCursor) {
   *   const page2 = await governanceRuleService.getHistory({
   *     limit: 10,
   *     cursor: page1.nextCursor,
   *   });
   * }
   * ```
   */
  async getHistory(
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

      const rules = governanceRulesArrayFromDto(response.result);

      // Verify each rule set if verification is enabled
      if (this.minValidSignatures > 0) {
        for (const rule of rules) {
          this.verifyGovernanceRules(rule);
        }
      }

      return {
        items: rules,
        nextCursor: response.cursor,
      };
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
    const rules = await this.get();

    if (!rules) {
      throw new IntegrityError("No governance rules available");
    }

    if (!rules.rulesContainer) {
      throw new IntegrityError("Rules container is empty");
    }

    return rulesContainerFromBase64(rules.rulesContainer);
  }

  /**
   * Decodes a rules container from a GovernanceRules object.
   *
   * @param rules - The governance rules containing the encoded container
   * @returns The decoded rules container
   * @throws {@link IntegrityError} If decoding fails
   *
   * @example
   * ```typescript
   * const rules = await governanceRuleService.get();
   * if (rules) {
   *   const container = governanceRuleService.decodeRulesContainer(rules);
   *   console.log(`Users: ${container.users.length}`);
   * }
   * ```
   */
  decodeRulesContainer(rules: GovernanceRules): DecodedRulesContainer {
    if (!rules.rulesContainer) {
      throw new IntegrityError("Rules container is empty");
    }

    return rulesContainerFromBase64(rules.rulesContainer);
  }

  /**
   * Verifies that governance rules have enough valid SuperAdmin signatures.
   *
   * This method is called automatically when minValidSignatures > 0.
   * For manual verification, call this method directly.
   *
   * @param rules - The governance rules to verify
   * @throws {@link IntegrityError} If verification fails
   *
   * @example
   * ```typescript
   * const rules = await governanceRuleService.get();
   * if (rules) {
   *   governanceRuleService.verifyGovernanceRules(rules);
   *   console.log("Rules verified successfully");
   * }
   * ```
   */
  verifyGovernanceRules(rules: GovernanceRules): void {
    if (this.minValidSignatures <= 0) {
      // Verification disabled
      return;
    }

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

    // Verify each signature cryptographically against SuperAdmin keys
    let validCount = 0;
    const seenUserIds = new Set<string>();

    for (const sig of rules.rulesSignatures) {
      if (!sig.signature) {
        continue;
      }

      // Ensure distinct user signatures (prevent same user signing multiple times)
      if (sig.userId && seenUserIds.has(sig.userId)) {
        continue;
      }

      // Verify signature using ECDSA against SuperAdmin public keys
      if (isValidSignature(rulesData, sig.signature, this.superAdminKeys)) {
        validCount++;
        if (sig.userId) {
          seenUserIds.add(sig.userId);
        }

        // Early return once threshold is met
        if (validCount >= this.minValidSignatures) {
          return;
        }
      }
    }

    if (validCount < this.minValidSignatures) {
      throw new IntegrityError(
        `Insufficient valid signatures: found ${validCount}, required ${this.minValidSignatures}`
      );
    }
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
