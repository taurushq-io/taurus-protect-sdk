/**
 * Status of a transaction request.
 *
 * Request statuses track the progress of a request from creation through
 * approval, signing, broadcasting, and confirmation on the blockchain.
 *
 * Common status transitions:
 * - CREATED -> PENDING -> APPROVED -> BROADCASTED -> CONFIRMED
 * - PENDING -> REJECTED (if rejected by an approver)
 * - BROADCASTED -> PERMANENT_FAILURE (if transaction fails)
 */
export enum RequestStatus {
  /** Request has received secondary approval */
  APPROVED_2 = "APPROVED_2",
  /** Request has been approved and is ready for signing */
  APPROVED = "APPROVED",
  /** Request is in the process of being approved */
  APPROVING = "APPROVING",
  /** Auto prepared secondary status */
  AUTO_PREPARED_2 = "AUTO_PREPARED_2",
  /** Auto prepared status */
  AUTO_PREPARED = "AUTO_PREPARED",
  /** Broadcasting secondary status */
  BROADCASTING_2 = "BROADCASTING_2",
  /** Broadcasting status */
  BROADCASTING = "BROADCASTING",
  /** Broadcasted to the blockchain */
  BROADCASTED = "BROADCASTED",
  /** Bundle has been approved */
  BUNDLE_APPROVED = "BUNDLE_APPROVED",
  /** Bundle is being broadcast */
  BUNDLE_BROADCASTING = "BUNDLE_BROADCASTING",
  /** Bundle is ready */
  BUNDLE_READY = "BUNDLE_READY",
  /** Request was canceled */
  CANCELED = "CANCELED",
  /** Transaction confirmed on blockchain */
  CONFIRMED = "CONFIRMED",
  /** Request has been created */
  CREATED = "CREATED",
  /** Diem burn MBS approved */
  DIEM_BURN_MBS_APPROVED = "DIEM_BURN_MBS_APPROVED",
  /** Diem burn MBS pending */
  DIEM_BURN_MBS_PENDING = "DIEM_BURN_MBS_PENDING",
  /** Diem mint MBS approved */
  DIEM_MINT_MBS_APPROVED = "DIEM_MINT_MBS_APPROVED",
  /** Diem mint MBS completed */
  DIEM_MINT_MBS_COMPLETED = "DIEM_MINT_MBS_COMPLETED",
  /** Diem mint MBS pending */
  DIEM_MINT_MBS_PENDING = "DIEM_MINT_MBS_PENDING",
  /** Request has expired */
  EXPIRED = "EXPIRED",
  /** Fast approved secondary status */
  FAST_APPROVED_2 = "FAST_APPROVED_2",
  /** HSM signing failed (secondary) */
  HSM_FAILED_2 = "HSM_FAILED_2",
  /** HSM signing failed */
  HSM_FAILED = "HSM_FAILED",
  /** HSM ready for signing (secondary) */
  HSM_READY_2 = "HSM_READY_2",
  /** HSM ready for signing */
  HSM_READY = "HSM_READY",
  /** HSM signed (secondary) */
  HSM_SIGNED_2 = "HSM_SIGNED_2",
  /** HSM has signed the transaction */
  HSM_SIGNED = "HSM_SIGNED",
  /** Manual broadcast required */
  MANUAL_BROADCAST = "MANUAL_BROADCAST",
  /** Transaction has been mined */
  MINED = "MINED",
  /** Transaction partially confirmed */
  PARTIALLY_CONFIRMED = "PARTIALLY_CONFIRMED",
  /** Request is pending approval */
  PENDING = "PENDING",
  /** Transaction permanently failed */
  PERMANENT_FAILURE = "PERMANENT_FAILURE",
  /** Request is ready */
  READY = "READY",
  /** Request was rejected */
  REJECTED = "REJECTED",
  /** Transaction has been sent */
  SENT = "SENT",
  /** Signet transaction completed */
  SIGNET_COMPLETED = "SIGNET_COMPLETED",
  /** Signet transaction pending */
  SIGNET_PENDING = "SIGNET_PENDING",
  /** Unknown status */
  UNKNOWN = "UNKNOWN",
}

/**
 * Request type enum.
 */
export enum RequestType {
  UNKNOWN = "UNKNOWN",
  TRANSFER = "TRANSFER",
  SIGN = "SIGN",
  CONTRACT_CALL = "CONTRACT_CALL",
  PAYMENT = "payment",
  FUNCTION_CALL = "function_call",
}

/**
 * Request metadata containing cryptographic hash information.
 *
 * The metadata is used for integrity verification:
 * - hash: SHA-256 hash of the payloadAsString
 * - payloadAsString: The canonical string representation used for hash computation
 *
 * SECURITY: payload field intentionally omitted - use payloadAsString only.
 * The raw payload object could be tampered with by an attacker while
 * payloadAsString remains unchanged (hash still verifies). By not having
 * this field, we enforce that all data extraction uses the verified source.
 */
export interface RequestMetadata {
  /** SHA-256 hash of the request payload (hex-encoded) */
  readonly hash: string;
  /** The raw payload used for hash computation */
  readonly payloadAsString: string | undefined;
  // SECURITY: payload field intentionally omitted - use payloadAsString only.
  // Use JSON.parse(payloadAsString) for secure data extraction.
}

/** A signed blockchain transaction associated with a request. */
export interface SignedRequest {
  readonly id: string;
  readonly signedRequest: string;
  readonly status: RequestStatus | string;
  readonly hash: string;
  readonly block: number;
  readonly details: string;
  readonly creationDate: Date | undefined;
  readonly updateDate: Date | undefined;
  readonly broadcastDate: Date | undefined;
  readonly confirmationDate: Date | undefined;
}

/** Amount information extracted from request metadata payload. */
export interface RequestMetadataAmount {
  readonly valueFrom: string;
  readonly valueTo: string;
  readonly rate: string;
  readonly decimals: number;
  readonly currencyFrom: string;
  readonly currencyTo: string;
}

/**
 * Represents a transaction request in Taurus-PROTECT.
 *
 * Requests are cryptographically protected:
 * - The metadata.hash field contains a SHA-256 hash of the request payload
 * - Hash verification should use constant-time comparison
 * - Approval requires ECDSA signing of the hash
 *
 * **Security Features:**
 * - All requests fetched via RequestService.get() have their hash verified
 * - Request approval uses ECDSA signatures for cryptographic verification
 */
export interface Request {
  /** Unique identifier */
  readonly id: number;
  /** Request type */
  readonly type: RequestType | string;
  /** Current status */
  readonly status: RequestStatus | string;
  /** Metadata containing hash and payload for verification */
  readonly metadata: RequestMetadata | undefined;
  /** Tenant ID */
  readonly tenantId: string | undefined;
  /** Currency/asset identifier */
  readonly currency: string | undefined;
  /** Currency information */
  readonly currencyInfo: CurrencyInfo | undefined;
  /** Request memo/reference */
  readonly memo: string | undefined;
  /** Request rule applied */
  readonly rule: string | undefined;
  /** External request ID for idempotency */
  readonly externalRequestId: string | undefined;
  /** Request bundle ID */
  readonly requestBundleId: string | undefined;
  /** Groups that need to approve this request */
  readonly needsApprovalFrom: string[];
  /** Approvers information */
  readonly approvers: Approvers | undefined;
  /** Creation timestamp */
  readonly createdAt: Date | undefined;
  /** Last update timestamp */
  readonly updatedAt: Date | undefined;
  /** Associated tags */
  readonly tags: string[];
  /** Signed blockchain transactions associated with this request */
  readonly signedRequests: SignedRequest[];
}

/**
 * Currency information.
 */
export interface CurrencyInfo {
  readonly id: string | undefined;
  readonly symbol: string | undefined;
  readonly name: string | undefined;
  readonly decimals: number | undefined;
  readonly blockchain: string | undefined;
  readonly network: string | undefined;
}

/**
 * Approvers group specifying a group and required signature count.
 */
export interface ApproversGroup {
  /** External group ID */
  readonly externalGroupId: string | undefined;
  /** Minimum signatures required from this group */
  readonly minimumSignatures: number | undefined;
}

/**
 * Sequential approvers groups that must approve in order.
 */
export interface ParallelApproversGroups {
  /** Sequential groups */
  readonly sequential: ApproversGroup[];
}

/**
 * Approvers information for a request.
 *
 * The approval structure supports parallel and sequential approval groups,
 * allowing complex multi-signature approval workflows.
 */
export interface Approvers {
  /** Parallel approval groups */
  readonly parallel: ParallelApproversGroups[];
}

/**
 * Options for listing requests.
 */
export interface ListRequestsOptions {
  /** Maximum number of requests to return */
  limit?: number;
  /** Page to request (for cursor-based pagination) */
  pageRequest?: "FIRST" | "PREVIOUS" | "NEXT" | "LAST";
  /** Current page cursor */
  currentPage?: string;
  /** Filter by request statuses */
  statuses?: RequestStatus[];
  /** Filter by request types */
  types?: string[];
  /** Filter by request IDs */
  ids?: number[];
  /** Filter by currency ID */
  currencyId?: string;
  /** Filter requests created after this date */
  fromDate?: Date;
  /** Filter requests created before this date */
  toDate?: Date;
  /** Sort order */
  sortOrder?: "ASC" | "DESC";
  /** Filter by external request IDs */
  externalRequestIds?: string[];
}

/**
 * Options for listing requests pending approval.
 */
export interface ListRequestsForApprovalOptions {
  /** Maximum number of requests to return */
  limit?: number;
  /** Page to request (for cursor-based pagination) */
  pageRequest?: "FIRST" | "PREVIOUS" | "NEXT" | "LAST";
  /** Current page cursor */
  currentPage?: string;
  /** Filter by request types */
  types?: string[];
  /** Exclude request types */
  excludeTypes?: string[];
  /** Filter by request IDs */
  ids?: number[];
  /** Filter by currency ID */
  currencyId?: string;
  /** Sort order */
  sortOrder?: "ASC" | "DESC";
  /** Filter by external request IDs */
  externalRequestIds?: string[];
}

/**
 * Options for creating an internal transfer request.
 */
export interface CreateInternalTransferOptions {
  /** Source address ID */
  fromAddressId: number;
  /** Destination address ID */
  toAddressId: number;
  /** Transfer amount as string (to preserve precision) */
  amount: string;
  /** Optional comment */
  comment?: string;
  /** Optional external request ID for idempotency */
  externalRequestId?: string;
  /** Optional gas limit */
  gasLimit?: string;
  /** Optional fee limit */
  feeLimit?: string;
}

/**
 * Options for creating an external transfer request.
 */
export interface CreateExternalTransferOptions {
  /** Source address ID */
  fromAddressId: number;
  /** Destination whitelisted address ID */
  toWhitelistedAddressId: number;
  /** Transfer amount as string (to preserve precision) */
  amount: string;
  /** Optional comment */
  comment?: string;
  /** Optional external request ID for idempotency */
  externalRequestId?: string;
  /** Optional gas limit */
  gasLimit?: string;
  /** Optional fee limit */
  feeLimit?: string;
  /** Optional destination address memo */
  destinationAddressMemo?: string;
}

/**
 * Options for creating an internal transfer from a wallet.
 */
export interface CreateInternalTransferFromWalletOptions {
  /** Source wallet ID (must be an omnibus wallet) */
  fromWalletId: number;
  /** Destination address ID */
  toAddressId: number;
  /** Transfer amount as string (to preserve precision) */
  amount: string;
  /** Optional comment */
  comment?: string;
  /** Optional external request ID for idempotency */
  externalRequestId?: string;
  /** Optional gas limit */
  gasLimit?: string;
  /** Optional fee limit */
  feeLimit?: string;
}

/**
 * Options for creating an external transfer from a wallet.
 */
export interface CreateExternalTransferFromWalletOptions {
  /** Source wallet ID (must be an omnibus wallet) */
  fromWalletId: number;
  /** Destination whitelisted address ID */
  toWhitelistedAddressId: number;
  /** Transfer amount as string (to preserve precision) */
  amount: string;
  /** Optional comment */
  comment?: string;
  /** Optional external request ID for idempotency */
  externalRequestId?: string;
  /** Optional gas limit */
  gasLimit?: string;
  /** Optional fee limit */
  feeLimit?: string;
  /** Optional destination address memo */
  destinationAddressMemo?: string;
}

/**
 * Options for creating an incoming request from an exchange.
 */
export interface CreateIncomingRequestOptions {
  /** Source exchange ID */
  fromExchangeId: number;
  /** Destination address ID */
  toAddressId: number;
  /** Transfer amount as string (to preserve precision) */
  amount: string;
  /** Optional comment */
  comment?: string;
  /** Optional external request ID for idempotency */
  externalRequestId?: string;
}
