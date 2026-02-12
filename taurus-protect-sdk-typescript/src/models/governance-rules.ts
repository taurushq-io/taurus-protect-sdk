/**
 * A user in the governance rules system.
 */
export interface RuleUser {
  readonly id: string | undefined;
  readonly name: string | undefined;
  readonly publicKeyPem: string | undefined;
  readonly roles: string[];
}

/**
 * A group in the governance rules system.
 */
export interface RuleGroup {
  readonly id: string | undefined;
  readonly name: string | undefined;
  readonly userIds: string[];
}

/**
 * Group threshold for approval.
 */
export interface GroupThreshold {
  readonly groupId: string | undefined;
  readonly minimumSignatures: number;
  readonly threshold: number;
}

/**
 * Sequential thresholds containing group thresholds.
 */
export interface SequentialThresholds {
  readonly thresholds: GroupThreshold[];
}

/**
 * Transaction approval rules.
 */
export interface TransactionRules {
  readonly parallelThresholds: SequentialThresholds[];
}

/**
 * Address whitelisting rules for a blockchain/network.
 */
export interface AddressWhitelistingRules {
  readonly currency: string | undefined;
  readonly network: string | undefined;
  readonly parallelThresholds: SequentialThresholds[];
}

/**
 * Contract address whitelisting rules.
 */
export interface ContractAddressWhitelistingRules {
  readonly blockchain: string | undefined;
  readonly network: string | undefined;
  readonly parallelThresholds: SequentialThresholds[];
}

/**
 * User signature in governance rules.
 */
export interface RuleUserSignature {
  readonly userId: string | undefined;
  readonly signature: string | undefined;
}

/**
 * Rules trail entry showing history of rule changes.
 */
export interface RulesTrail {
  readonly userId: string | undefined;
  readonly action: string | undefined;
  readonly timestamp: Date | undefined;
}

/**
 * Governance rules as returned by the API.
 *
 * Contains the encoded rules container and associated metadata including
 * SuperAdmin signatures for verification.
 */
export interface GovernanceRules {
  /** Base64-encoded rules container */
  readonly rulesContainer: string | undefined;
  /** List of SuperAdmin signatures on the rules */
  readonly rulesSignatures: RuleUserSignature[];
  /** Whether the rules are locked */
  readonly locked: boolean;
  /** Creation timestamp */
  readonly creationDate: Date | undefined;
  /** Last update timestamp */
  readonly updateDate: Date | undefined;
  /** Audit trail of rule changes */
  readonly trails: RulesTrail[];
}

/**
 * Options for listing governance rules history.
 */
export interface ListGovernanceRulesHistoryOptions {
  /** Maximum number of rules to return (default: 50) */
  readonly limit?: number;
  /** Pagination cursor from a previous call */
  readonly cursor?: string;
}

/**
 * Result from listing governance rules history.
 */
export interface GovernanceRulesHistoryResult {
  /** List of historical governance rules */
  readonly items: GovernanceRules[];
  /** Cursor for fetching the next page, undefined if no more pages */
  readonly nextCursor: string | undefined;
}

/**
 * Decoded rules container containing all governance rules.
 */
export interface DecodedRulesContainer {
  readonly users: RuleUser[];
  readonly groups: RuleGroup[];
  readonly minimumDistinctUserSignatures: number;
  readonly minimumDistinctGroupSignatures: number;
  readonly transactionRules: TransactionRules[];
  readonly addressWhitelistingRules: AddressWhitelistingRules[];
  readonly contractAddressWhitelistingRules: ContractAddressWhitelistingRules[];
  readonly enforcedRulesHash: string | undefined;
  readonly timestamp: number;
  readonly hsmSlotId: number;
  readonly minimumCommitmentSignatures: number;
  readonly engineIdentities: string[];
}

/**
 * Creates an empty DecodedRulesContainer.
 */
export function createEmptyRulesContainer(): DecodedRulesContainer {
  return {
    users: [],
    groups: [],
    minimumDistinctUserSignatures: 0,
    minimumDistinctGroupSignatures: 0,
    transactionRules: [],
    addressWhitelistingRules: [],
    contractAddressWhitelistingRules: [],
    enforcedRulesHash: undefined,
    timestamp: 0,
    hsmSlotId: 0,
    minimumCommitmentSignatures: 0,
    engineIdentities: [],
  };
}

/**
 * Finds the HSM public key from the rules container.
 * Looks for a user with the HSMSLOT role.
 */
export function getHsmPublicKey(container: DecodedRulesContainer): string | undefined {
  for (const user of container.users) {
    if (user.roles.includes('HSMSLOT') && user.publicKeyPem) {
      return user.publicKeyPem;
    }
  }
  return undefined;
}

/**
 * Finds a user by ID in the rules container.
 */
export function findUserById(container: DecodedRulesContainer, userId: string): RuleUser | undefined {
  return container.users.find(u => u.id === userId);
}

/**
 * Finds a group by ID in the rules container.
 */
export function findGroupById(container: DecodedRulesContainer, groupId: string): RuleGroup | undefined {
  return container.groups.find(g => g.id === groupId);
}

/**
 * Checks if a value represents a wildcard (undefined, null, empty, or "Any").
 * Matches Java's DecodedRulesContainer.isWildcard() behavior.
 */
function isWildcard(value: string | undefined): boolean {
  return !value || value === '' || value.toLowerCase() === 'any';
}

/**
 * Priority-based matching for address whitelisting rules.
 * Returns: exact match > blockchain-only match > global default (wildcard).
 *
 * Wildcard values are: undefined, null, empty string, or "Any" (case-insensitive).
 */
export function findAddressWhitelistingRules(
  container: DecodedRulesContainer,
  blockchain: string,
  network: string
): AddressWhitelistingRules | undefined {
  let blockchainMatch: AddressWhitelistingRules | undefined;
  let globalMatch: AddressWhitelistingRules | undefined;

  for (const rules of container.addressWhitelistingRules) {
    const currency = rules.currency ?? '';
    const ruleNetwork = rules.network ?? '';
    const ruleIsGlobalDefault = isWildcard(currency);
    const blockchainMatches = !ruleIsGlobalDefault && currency === blockchain;
    const networkMatches = ruleNetwork === network;
    const ruleHasWildcardNetwork = isWildcard(ruleNetwork);

    // Priority 1: Exact match (blockchain + network)
    if (blockchainMatches && networkMatches) {
      return rules;
    }

    // Priority 2: Blockchain match with wildcard network
    if (blockchainMatches && ruleHasWildcardNetwork && !blockchainMatch) {
      blockchainMatch = rules;
    }

    // Priority 3: Global default (wildcard blockchain)
    if (ruleIsGlobalDefault && !globalMatch) {
      globalMatch = rules;
    }
  }

  return blockchainMatch ?? globalMatch;
}

/**
 * Priority-based matching for contract address whitelisting rules.
 *
 * Wildcard values are: undefined, null, empty string, or "Any" (case-insensitive).
 */
export function findContractAddressWhitelistingRules(
  container: DecodedRulesContainer,
  blockchain: string,
  network: string
): ContractAddressWhitelistingRules | undefined {
  let blockchainMatch: ContractAddressWhitelistingRules | undefined;
  let globalMatch: ContractAddressWhitelistingRules | undefined;

  for (const rules of container.contractAddressWhitelistingRules) {
    const ruleBlockchain = rules.blockchain ?? '';
    const ruleNetwork = rules.network ?? '';
    const ruleIsGlobalDefault = isWildcard(ruleBlockchain);
    const blockchainMatches = !ruleIsGlobalDefault && ruleBlockchain === blockchain;
    const networkMatches = ruleNetwork === network;
    const ruleHasWildcardNetwork = isWildcard(ruleNetwork);

    // Priority 1: Exact match (blockchain + network)
    if (blockchainMatches && networkMatches) {
      return rules;
    }

    // Priority 2: Blockchain match with wildcard network
    if (blockchainMatches && ruleHasWildcardNetwork && !blockchainMatch) {
      blockchainMatch = rules;
    }

    // Priority 3: Global default (wildcard blockchain)
    if (ruleIsGlobalDefault && !globalMatch) {
      globalMatch = rules;
    }
  }

  return blockchainMatch ?? globalMatch;
}
