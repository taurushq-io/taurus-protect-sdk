/**
 * Contract whitelist models for Taurus-PROTECT SDK.
 *
 * This module provides domain models for whitelisted contract addresses
 * such as ERC20 tokens, NFT collections (ERC721/ERC1155), FA2 tokens on Tezos,
 * and other smart contract-based assets.
 */

/**
 * Represents a whitelisted contract address (token, NFT, etc.).
 *
 * Whitelisted contract addresses allow transactions involving specific tokens or NFTs
 * to be processed according to governance rules.
 */
export interface WhitelistedContract {
  /** Unique identifier for the whitelisted contract. */
  readonly id: string;
  /** Tenant identifier. */
  readonly tenantId: string | undefined;
  /** Blockchain symbol (e.g., "ETH", "MATIC", "XTZ"). */
  readonly blockchain: string | undefined;
  /** Network identifier (e.g., "mainnet", "goerli"). */
  readonly network: string | undefined;
  /** The smart contract address. */
  readonly contractAddress: string | undefined;
  /** Token symbol (e.g., "USDC", "WETH"). */
  readonly symbol: string | undefined;
  /** Human-readable name. */
  readonly name: string | undefined;
  /** Number of decimals (0 for NFTs). */
  readonly decimals: number | undefined;
  /** Contract kind (e.g., "erc20", "erc721", "fa2"). */
  readonly kind: string | undefined;
  /** Token ID for NFTs (null for fungible tokens). */
  readonly tokenId: string | undefined;
  /** Current status of the whitelisted contract. */
  readonly status: string | undefined;
  /** Whether the business rule is enabled for this contract. */
  readonly businessRuleEnabled: boolean | undefined;
  /** Contract attributes. */
  readonly attributes: WhitelistedContractAttribute[];
}

/**
 * Represents an attribute on a whitelisted contract.
 *
 * Attributes are key-value metadata that can be attached to whitelisted contracts
 * for additional categorization or information.
 */
export interface WhitelistedContractAttribute {
  /** Attribute ID. */
  readonly id: string | undefined;
  /** Attribute key. */
  readonly key: string | undefined;
  /** Attribute value. */
  readonly value: string | undefined;
  /** Content type. */
  readonly contentType: string | undefined;
  /** Attribute type. */
  readonly type: string | undefined;
  /** Attribute subtype. */
  readonly subType: string | undefined;
  /** Whether this attribute is a file. */
  readonly isFile: boolean | undefined;
}

/**
 * Metadata containing the hash and signed payload for verification.
 *
 * SECURITY: payload field intentionally omitted - use payloadAsString only.
 * The raw payload object could be tampered with by an attacker while
 * payloadAsString remains unchanged (hash still verifies).
 */
export interface WhitelistedContractMetadata {
  /** SHA-256 hash of the payloadAsString. */
  readonly hash: string | undefined;
  /** JSON string of the signed payload. */
  readonly payloadAsString: string | undefined;
  // SECURITY: payload field intentionally omitted - use payloadAsString only.
}

/**
 * Signature data for a whitelisted contract.
 */
export interface WhitelistedContractSignature {
  /** Base64-encoded cryptographic signature. */
  readonly signature: string | undefined;
  /** Optional comment from the signer. */
  readonly comment: string | undefined;
  /** List of hashes covered by this signature. */
  readonly hashes: string[];
  /** User ID of the signer. */
  readonly userId: string | undefined;
}

/**
 * Signed whitelisted contract address data containing signatures.
 */
export interface SignedWhitelistedContract {
  /** Base64-encoded signed payload. */
  readonly payload: string | undefined;
  /** List of signatures on this contract. */
  readonly signatures: WhitelistedContractSignature[];
}

/**
 * Audit trail entry for a whitelisted contract.
 */
export interface WhitelistedContractTrail {
  /** Trail ID. */
  readonly id: string | undefined;
  /** Action taken. */
  readonly action: string | undefined;
  /** User ID who performed the action. */
  readonly userId: string | undefined;
  /** Comment on the action. */
  readonly comment: string | undefined;
  /** Timestamp of the action. */
  readonly timestamp: string | undefined;
}

/**
 * Approver information for a whitelisted contract.
 */
export interface WhitelistedContractApprover {
  /** User ID of the approver. */
  readonly userId: string | undefined;
  /** User name or email. */
  readonly userName: string | undefined;
  /** Whether approval is pending from this approver. */
  readonly pending: boolean | undefined;
}

/**
 * Approvers configuration for a whitelisted contract.
 */
export interface WhitelistedContractApprovers {
  /** Required number of approvals. */
  readonly required: number | undefined;
  /** List of approver groups. */
  readonly groups: WhitelistedContractApproverGroup[];
}

/**
 * Approver group for a whitelisted contract.
 */
export interface WhitelistedContractApproverGroup {
  /** Group ID. */
  readonly id: string | undefined;
  /** Group name. */
  readonly name: string | undefined;
  /** Required approvals from this group. */
  readonly required: number | undefined;
  /** Users in this group. */
  readonly users: WhitelistedContractApprover[];
}

/**
 * Envelope containing a whitelisted contract with all data needed for verification.
 *
 * This envelope contains:
 * - Metadata with hash and payload
 * - Signed contract address with user signatures
 * - Audit trails and approvers configuration
 * - Attributes attached to the contract
 */
export interface SignedWhitelistedContractEnvelope {
  /** Unique identifier of the whitelisted contract. */
  readonly id: string;
  /** Tenant identifier. */
  readonly tenantId: string | undefined;
  /** Blockchain identifier. */
  readonly blockchain: string | undefined;
  /** Network identifier. */
  readonly network: string | undefined;
  /** Metadata with hash and payload. */
  readonly metadata: WhitelistedContractMetadata | undefined;
  /** Signed contract data with user signatures. */
  readonly signedContractAddress: SignedWhitelistedContract | undefined;
  /** The action type for this entry (e.g., "create", "update", "delete"). */
  readonly action: string | undefined;
  /** Audit trail of actions taken on this entry. */
  readonly trails: WhitelistedContractTrail[];
  /** Rules container identifier. */
  readonly rulesContainer: string | undefined;
  /** Rule identifier. */
  readonly rule: string | undefined;
  /** Rules signatures for verification. */
  readonly rulesSignatures: string | undefined;
  /** Approvers configuration for this entry. */
  readonly approvers: WhitelistedContractApprovers | undefined;
  /** Custom attributes attached to this entry. */
  readonly attributes: WhitelistedContractAttribute[];
  /** Current status of the whitelist entry. */
  readonly status: string | undefined;
  /** Whether the business rule is enabled for this contract. */
  readonly businessRuleEnabled: boolean | undefined;
}

/**
 * Result of listing whitelisted contracts with pagination.
 */
export interface WhitelistedContractResult {
  /** List of whitelisted contracts. */
  readonly contracts: WhitelistedContract[];
  /** Total number of items matching the query. */
  readonly totalItems: number;
}

/**
 * Options for listing whitelisted contracts.
 */
export interface ListWhitelistedContractsOptions {
  /** Maximum number of items to return (max 100). */
  limit?: number;
  /** Offset for pagination. */
  offset?: number;
  /** Search query. */
  query?: string;
  /** Filter by blockchain. */
  blockchain?: string;
  /** Filter by network. */
  network?: string;
  /** Filter for NFT contracts only (deprecated: use kindTypes instead). */
  isNFT?: boolean;
  /** Filter by contract kind types (e.g., "nft", "token"). */
  kindTypes?: string[];
  /** Filter by specific contract IDs. */
  contractIds?: string[];
}

/**
 * Options for listing whitelisted contracts pending approval.
 */
export interface ListForApprovalOptions {
  /** Maximum number of items to return (max 100). */
  limit?: number;
  /** Offset for pagination. */
  offset?: number;
  /** Filter by specific IDs. */
  ids?: string[];
}

/**
 * Request to create a whitelisted contract.
 */
export interface CreateWhitelistedContractRequest {
  /** Blockchain identifier (e.g., "ETH", "MATIC", "XTZ"). */
  blockchain: string;
  /** Network identifier (e.g., "mainnet", "goerli"). */
  network: string;
  /** The smart contract address. */
  contractAddress?: string;
  /** Token symbol (e.g., "USDC", "WETH"). */
  symbol: string;
  /** Human-readable name. */
  name: string;
  /** Number of decimals (0 for NFTs). */
  decimals: number;
  /** Contract kind (e.g., "erc20", "erc721", "fa2"). */
  kind: string;
  /** Token ID for NFTs (null for fungible tokens). */
  tokenId?: string;
}

/**
 * Request to update a whitelisted contract.
 */
export interface UpdateWhitelistedContractRequest {
  /** New symbol. */
  symbol: string;
  /** New name. */
  name: string;
  /** New decimals value. */
  decimals: number;
}
