/**
 * Multi-factor signature domain models for Taurus-PROTECT SDK.
 *
 * Multi-factor signatures are used for operations that require approval from
 * multiple parties, such as critical governance changes, high-value transactions,
 * or sensitive administrative actions.
 */

/**
 * Entity type for multi-factor signatures.
 *
 * Identifies what kind of entity is associated with the multi-factor signature request.
 */
export enum MultiFactorSignatureEntityType {
  /** Request entity type */
  REQUEST = 'REQUEST',
  /** Whitelisted address entity type */
  WHITELISTED_ADDRESS = 'WHITELISTED_ADDRESS',
  /** Whitelisted contract entity type */
  WHITELISTED_CONTRACT = 'WHITELISTED_CONTRACT',
}

/**
 * Represents information about a multi-factor signature request.
 *
 * Multi-factor signatures are used for operations that require multiple
 * approvals, such as critical governance changes or high-value transactions.
 */
export interface MultiFactorSignatureInfo {
  /** The multi-factor signature ID */
  readonly id: string;
  /** The payloads that need to be signed */
  readonly payloadToSign: string[];
  /** The type of entity associated with this signature request */
  readonly entityType: MultiFactorSignatureEntityType;
}

/**
 * Request to create multi-factor signatures.
 */
export interface CreateMultiFactorSignatureRequest {
  /** The type of entities */
  readonly entityType: MultiFactorSignatureEntityType;
  /** The list of entity IDs to create signatures for */
  readonly entityIds: string[];
}

/**
 * Request to approve a multi-factor signature.
 */
export interface ApproveMultiFactorSignatureRequest {
  /** The multi-factor signature ID */
  readonly id: string;
  /** The signature for approval (base64 encoded) */
  readonly signature: string;
  /** Optional comment for the approval */
  readonly comment?: string;
}

/**
 * Request to reject a multi-factor signature.
 */
export interface RejectMultiFactorSignatureRequest {
  /** The multi-factor signature ID */
  readonly id: string;
  /** Optional comment explaining the rejection */
  readonly comment?: string;
}
