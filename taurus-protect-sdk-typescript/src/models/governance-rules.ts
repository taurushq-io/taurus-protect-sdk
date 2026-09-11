import type { RuleCell } from './rule-cell';

export type { RuleCell } from './rule-cell';

/**
 * Unknown protobuf fields captured at decode, keyed by field number, in the
 * ts-proto `_unknownFields` shape. Re-attached verbatim on encode so data from
 * a newer schema than this SDK version is never silently dropped.
 */
export type UnknownFields = { [key: number]: Uint8Array[] };

/** A protobuf `map<string, bytes>` properties bag. */
export type Properties = { [key: string]: Uint8Array };

/**
 * A user in the governance rules system.
 */
export interface RuleUser {
  readonly id: string | undefined;
  readonly name: string | undefined;
  readonly publicKeyPem: string | undefined;
  readonly roles: string[];
  readonly properties?: Properties;
  readonly unknownFields?: UnknownFields;
}

/**
 * A group in the governance rules system.
 */
export interface RuleGroup {
  readonly id: string | undefined;
  readonly name: string | undefined;
  readonly userIds: string[];
  readonly properties?: Properties;
  readonly unknownFields?: UnknownFields;
}

/**
 * Group threshold for approval.
 */
export interface GroupThreshold {
  readonly groupId: string | undefined;
  readonly minimumSignatures: number;
  readonly threshold: number;
  readonly unknownFields?: UnknownFields;
}

/**
 * Sequential thresholds containing group thresholds.
 */
export interface SequentialThresholds {
  readonly thresholds: GroupThreshold[];
  readonly unknownFields?: UnknownFields;
}

/**
 * A column in a transaction rule.
 */
export interface RuleColumn {
  readonly type: string;
  readonly name: string;
  readonly metadataKey: string;
  readonly unknownFields?: UnknownFields;
}

/**
 * A line/row in a transaction rule. Cells align positionally with the rule's
 * columns.
 */
export interface RuleLine {
  readonly cells: RuleCell[];
  readonly parallelThresholds: SequentialThresholds[];
  readonly priority: number;
  readonly properties?: Properties;
  readonly unknownFields?: UnknownFields;
}

/** EVM contract-call scoping. */
export interface EvmCallContract {
  readonly contractType: string;
  readonly methodSignature: string;
  /** Protobuf fields from a newer schema, re-attached on encode. */
  readonly unknownFields?: UnknownFields;
}

/** Tezos contract-call scoping. */
export interface XtzCallContract {
  readonly contractType: string;
  readonly methodSignature: string;
  /** Protobuf fields from a newer schema, re-attached on encode. */
  readonly unknownFields?: UnknownFields;
}

/** Cash-settlement scoping. */
export interface CashSettlement {
  readonly provider: string;
  readonly requestType: string;
  /** Protobuf fields from a newer schema, re-attached on encode. */
  readonly unknownFields?: UnknownFields;
}

/** Cosmos scoping. */
export interface CosmosDetails {
  readonly methodSignatures: string[];
  /** Protobuf fields from a newer schema, re-attached on encode. */
  readonly unknownFields?: UnknownFields;
}

/**
 * Additional transaction-rule configuration.
 */
export interface TransactionRuleDetails {
  readonly domain: string;
  readonly subDomain: string;
  readonly blockchain: string;
  readonly network: string;
  readonly evmCallContract?: EvmCallContract;
  readonly xtzCallContract?: XtzCallContract;
  readonly cashSettlement?: CashSettlement;
  readonly cosmosDetails?: CosmosDetails;
  readonly unknownFields?: UnknownFields;
}

/**
 * Transaction approval rules.
 */
export interface TransactionRules {
  readonly key: string;
  readonly columns: RuleColumn[];
  readonly lines: RuleLine[];
  readonly details?: TransactionRuleDetails;
  readonly unknownFields?: UnknownFields;
}

/** The type of a rule source (whitelisting cell). */
export enum RuleSourceType {
  Unknown = 0,
  InternalWallet = 1,
  InternalAddress = 2,
  AnyExchange = 3,
  Exchange = 4,
  ExternalAddress = 5,
}

export interface RuleSourceInternalWallet {
  readonly path: string;
}
export interface RuleSourceInternalAddress {
  readonly address: string;
  readonly path: string;
}
export interface RuleSourceExchange {
  readonly label: string;
}
export interface RuleSourceExternalAddress {
  readonly address: string;
  readonly memo: string;
}

/**
 * A source specification in an address-whitelisting rule line. Sources whose
 * type/content is unknown to this SDK version keep their exact wire bytes in
 * `raw` and are re-emitted verbatim on encode.
 */
export interface RuleSource {
  readonly type: RuleSourceType;
  readonly internalWallet?: RuleSourceInternalWallet;
  readonly internalAddress?: RuleSourceInternalAddress;
  readonly exchange?: RuleSourceExchange;
  readonly externalAddress?: RuleSourceExternalAddress;
  readonly raw?: Uint8Array;
}

/**
 * A source-specific address-whitelisting rule line.
 */
export interface AddressWhitelistingLine {
  readonly cells: RuleSource[];
  readonly parallelThresholds: SequentialThresholds[];
  readonly properties?: Properties;
  readonly unknownFields?: UnknownFields;
}

/**
 * Address whitelisting rules for a blockchain/network.
 */
export interface AddressWhitelistingRules {
  readonly currency: string | undefined;
  readonly network: string | undefined;
  readonly parallelThresholds: SequentialThresholds[];
  readonly lines: AddressWhitelistingLine[];
  readonly properties?: Properties;
  readonly unknownFields?: UnknownFields;
}

/**
 * Contract address whitelisting rules.
 */
export interface ContractAddressWhitelistingRules {
  readonly blockchain: string | undefined;
  readonly network: string | undefined;
  readonly parallelThresholds: SequentialThresholds[];
  readonly properties?: Properties;
  readonly unknownFields?: UnknownFields;
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
/**
 * A SuperAdmin public key as configured on the server, used to verify the signatures on
 * a governance rules container.
 */
export interface SuperAdminPublicKey {
  /** ID of the SuperAdmin user this key belongs to */
  readonly userId: string;
  /** PEM-encoded public key */
  readonly publicKey: string;
}

/** A history entry withheld because its SuperAdmin signatures did not verify. */
export interface ExcludedRuleset {
  /**
   * When the excluded ruleset was created — the only stable identifier a history entry
   * carries.
   */
  readonly creationDate: Date | undefined;
  /** Why it was excluded. */
  readonly reason: string;
}

export interface GovernanceRulesHistoryResult {
  /** Historical governance rules whose SuperAdmin signatures verified */
  readonly items: GovernanceRules[];
  /**
   * Entries withheld because their signatures did not verify, so a shortened page
   * cannot read as a complete one. History is LENIENT and does not throw when nothing
   * survives: a SuperAdmin key rotation makes every pre-rotation ruleset unverifiable,
   * and aborting would deny the whole audit trail.
   */
  readonly excludedUnverified: ExcludedRuleset[];
  /** Cursor for fetching the next page, undefined if no more pages */
  readonly nextCursor: string | undefined;
}

/**
 * Decoded rules container containing all governance rules.
 *
 * The model is lossless: it carries every governance protobuf field (including
 * per-node `properties` and `unknownFields`), so a container decoded from the
 * wire re-encodes without dropping data — even data introduced by a newer
 * schema than this SDK version.
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
  readonly properties?: Properties;
  readonly unknownFields?: UnknownFields;
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

function hasUf(uf: UnknownFields | undefined): boolean {
  return uf !== undefined && Object.keys(uf).length > 0;
}

/**
 * Reports whether any part of the container carries data unknown to this SDK
 * version (unknown protobuf fields, raw cells, or raw whitelisting sources).
 * Such data is preserved verbatim on re-encode; a true result is a hint to
 * upgrade the SDK before editing the container.
 */
export function hasUnknownFields(container: DecodedRulesContainer): boolean {
  if (hasUf(container.unknownFields)) return true;
  if (container.users.some((u) => hasUf(u.unknownFields))) return true;
  if (container.groups.some((g) => hasUf(g.unknownFields))) return true;
  for (const tr of container.transactionRules) {
    if (hasUf(tr.unknownFields)) return true;
    if (tr.columns.some((c) => hasUf(c.unknownFields))) return true;
    if (ruleDetailsHaveUnknown(tr.details)) return true;
    for (const l of tr.lines) {
      if (hasUf(l.unknownFields)) return true;
      if (thresholdsHaveUnknown(l.parallelThresholds)) return true;
      if (l.cells.some((c) => c.kind === 'RawCell')) return true;
    }
  }
  for (const awr of container.addressWhitelistingRules) {
    if (hasUf(awr.unknownFields) || thresholdsHaveUnknown(awr.parallelThresholds)) return true;
    for (const l of awr.lines) {
      if (hasUf(l.unknownFields) || thresholdsHaveUnknown(l.parallelThresholds)) return true;
      if (l.cells.some((s) => s.raw !== undefined && s.raw.length > 0)) return true;
    }
  }
  for (const cawr of container.contractAddressWhitelistingRules) {
    if (hasUf(cawr.unknownFields) || thresholdsHaveUnknown(cawr.parallelThresholds)) return true;
  }
  return false;
}

/**
 * Reports unknown fields on the details node and each nested scoping sub-message.
 * Checking only the details node would report a container clean while a nested
 * contract-call scoping node carried data from a newer schema.
 */
function ruleDetailsHaveUnknown(details: TransactionRuleDetails | undefined): boolean {
  if (!details) return false;
  if (hasUf(details.unknownFields)) return true;
  return [
    details.evmCallContract,
    details.xtzCallContract,
    details.cashSettlement,
    details.cosmosDetails,
  ].some((n) => n !== undefined && hasUf(n.unknownFields));
}

function thresholdsHaveUnknown(thresholds: SequentialThresholds[]): boolean {
  return thresholds.some(
    (st) => hasUf(st.unknownFields) || st.thresholds.some((t) => hasUf(t.unknownFields))
  );
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
 * Every rule the tier walk in {@link findAddressWhitelistingRules} could select for this
 * blockchain, across all possible network values.
 *
 * This exists because the network half of the rule key is NOT always signed. Governance
 * carries a per-rule `includeNetworkInPayload` flag, and when it is off the signed
 * payload has no `network` member at all — the common case in captured production data
 * (`tests/unit/fixtures/whitelisted-address-raw-response.json`). The key then falls back
 * to the network on the unsigned response DTO, which hands a response-controlling server
 * the choice of WHICH rule judges the row, and therefore which group quorum it must
 * meet. That is not closeable by reading the flag: `includeNetworkInPayload` has no
 * proto backing in any of the four SDKs (`request_reply.proto`'s
 * `AddressWhitelistingRules` carries only currency, parallelThresholds, properties,
 * network, lines), so it is never part of the SuperAdmin-signed container — and this SDK
 * does not even model the field.
 *
 * So when the network is unsigned the caller must satisfy EVERY rule this returns, not
 * the one the DTO named. Where a chain has a single reachable tier — again the common
 * case — the set has one element and behaviour is unchanged.
 *
 * Reachability, mirroring the tier walk: every rule for this chain is reachable by naming
 * its network; the chain's wildcard-network rule is reachable by naming a network no
 * exact rule covers; and the global default is reachable ONLY when the chain has no
 * wildcard-network rule, because priority 2 would otherwise win.
 *
 * @param container - the verified rules container
 * @param blockchain - the chain from the signed payload
 * @returns the reachable rules, chain-specific ones first
 */
export function findAddressWhitelistingRuleCandidates(
  container: DecodedRulesContainer,
  blockchain: string
): AddressWhitelistingRules[] {
  const candidates: AddressWhitelistingRules[] = [];
  let globalDefault: AddressWhitelistingRules | undefined;
  let chainHasWildcardNetwork = false;

  for (const rules of container.addressWhitelistingRules) {
    const currency = rules.currency ?? '';
    if (isWildcard(currency)) {
      globalDefault ??= rules;
      continue;
    }
    if (currency !== blockchain) {
      continue;
    }
    candidates.push(rules);
    if (isWildcard(rules.network ?? '')) {
      chainHasWildcardNetwork = true;
    }
  }

  if (!chainHasWildcardNetwork && globalDefault) {
    candidates.push(globalDefault);
  }
  return candidates;
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

/**
 * The asset peer of {@link findAddressWhitelistingRuleCandidates}.
 *
 * Same reasoning: when the signed payload omits `network`, the unsigned response DTO
 * would otherwise choose which quorum judges the asset.
 *
 * @param container - the verified rules container
 * @param blockchain - the chain from the signed payload
 * @returns the reachable rules, chain-specific ones first
 */
export function findContractAddressWhitelistingRuleCandidates(
  container: DecodedRulesContainer,
  blockchain: string
): ContractAddressWhitelistingRules[] {
  const candidates: ContractAddressWhitelistingRules[] = [];
  let globalDefault: ContractAddressWhitelistingRules | undefined;
  let chainHasWildcardNetwork = false;

  for (const rules of container.contractAddressWhitelistingRules) {
    const ruleBlockchain = rules.blockchain ?? '';
    if (isWildcard(ruleBlockchain)) {
      globalDefault ??= rules;
      continue;
    }
    if (ruleBlockchain !== blockchain) {
      continue;
    }
    candidates.push(rules);
    if (isWildcard(rules.network ?? '')) {
      chainHasWildcardNetwork = true;
    }
  }

  if (!chainHasWildcardNetwork && globalDefault) {
    candidates.push(globalDefault);
  }
  return candidates;
}
