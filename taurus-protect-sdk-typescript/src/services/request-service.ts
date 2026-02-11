/**
 * Request service for Taurus-PROTECT SDK.
 *
 * Provides operations for creating, approving, rejecting, and querying
 * transaction requests. Requests represent actions to be performed on
 * the blockchain, such as transfers between addresses.
 *
 * **CRITICAL SECURITY FEATURES:**
 * - All requests fetched via `get()` have their hash verified using constant-time comparison
 * - Request approval uses ECDSA signatures for cryptographic verification
 * - Hash verification prevents tampering with request data
 *
 * @example
 * ```typescript
 * // Get a single request with hash verification
 * const request = await requestService.get(123);
 *
 * // List requests pending approval
 * const { requests, cursor } = await requestService.listForApproval({ limit: 50 });
 *
 * // Approve requests with ECDSA signature
 * const signedCount = await requestService.approveRequests(requests, privateKey);
 * ```
 */

import type { KeyObject } from "crypto";
import type { RequestsApi } from "../internal/openapi/apis/RequestsApi";
import type {
  Request,
  ListRequestsOptions,
  ListRequestsForApprovalOptions,
  CreateInternalTransferOptions,
  CreateExternalTransferOptions,
  CreateInternalTransferFromWalletOptions,
  CreateExternalTransferFromWalletOptions,
  CreateIncomingRequestOptions,
} from "../models/request";
import type { CursorPagination } from "../models/pagination";
import { BaseService } from "./base";
import { IntegrityError, NotFoundError, ServerError } from "../errors";
import { calculateHexHash, constantTimeCompare, signData } from "../crypto";
import { requestFromDto, requestsFromDto } from "../mappers/request";

/**
 * Result of a list operation with cursor-based pagination.
 */
export interface ListRequestsResult {
  /** The requests in this page */
  requests: Request[];
  /** Cursor for pagination */
  cursor: CursorPagination;
}

/**
 * Service for managing transaction requests.
 *
 * Provides operations for creating, approving, rejecting, and querying
 * transaction requests. Requests represent actions to be performed on
 * the blockchain, such as transfers between addresses.
 *
 * **Important Security Features:**
 * - All requests fetched via `get()` have their hash verified using constant-time comparison
 * - Request approval uses ECDSA signatures for cryptographic verification
 */
export class RequestService extends BaseService {
  private readonly requestsApi: RequestsApi;

  /**
   * Creates a new RequestService instance.
   *
   * @param requestsApi - The OpenAPI RequestsApi instance
   */
  constructor(requestsApi: RequestsApi) {
    super();
    this.requestsApi = requestsApi;
  }

  /**
   * Get a request by ID with mandatory hash verification.
   *
   * **CRITICAL SECURITY:**
   * The hash of the request metadata payload is verified using
   * constant-time comparison to prevent timing attacks.
   *
   * @param requestId - The request ID to retrieve
   * @returns The verified request
   * @throws {IntegrityError} If hash verification fails
   * @throws {NotFoundError} If request not found
   * @throws {APIError} If API request fails
   *
   * @example
   * ```typescript
   * try {
   *   const request = await requestService.get(123);
   *   console.log(`Request ${request.id} status: ${request.status}`);
   * } catch (error) {
   *   if (error instanceof IntegrityError) {
   *     console.error('Security error: request hash mismatch');
   *   }
   * }
   * ```
   */
  async get(requestId: number): Promise<Request> {
    if (requestId <= 0) {
      throw new Error("requestId must be positive");
    }

    return this.execute(async () => {
      const response = await this.requestsApi.requestServiceGetRequest({
        id: String(requestId),
      });

      const result = response.result;
      if (!result) {
        throw new NotFoundError(`Request ${requestId} not found`);
      }

      const request = requestFromDto(result);
      if (!request) {
        throw new NotFoundError(`Request ${requestId} not found`);
      }

      // CRITICAL: Verify request hash using constant-time comparison
      this.verifyRequestHash(request);

      return request;
    });
  }

  /**
   * List requests with filtering and pagination.
   *
   * **Note:** Hash verification is NOT performed on list operations for performance.
   * Use `get()` to fetch individual requests with full verification.
   *
   * @param options - List options including filters and pagination
   * @returns Object containing requests array and cursor pagination
   * @throws {APIError} If API request fails
   *
   * @example
   * ```typescript
   * const result = await requestService.list({
   *   limit: 50,
   *   statuses: [RequestStatus.PENDING_APPROVAL],
   *   fromDate: new Date('2024-01-01'),
   * });
   * console.log(`Found ${result.requests.length} requests`);
   * ```
   */
  async list(options: ListRequestsOptions = {}): Promise<ListRequestsResult> {
    const limit = options.limit ?? 50;
    if (limit <= 0) {
      throw new Error("limit must be positive");
    }

    return this.execute(async () => {
      const response = await this.requestsApi.requestServiceGetRequestsV2({
        cursorPageSize: String(limit),
        cursorCurrentPage: options.currentPage,
        cursorPageRequest: options.pageRequest,
        from: options.fromDate,
        to: options.toDate,
        currencyID: options.currencyId,
        statuses: options.statuses?.map(String),
        types: options.types,
        ids: options.ids?.map(String),
        sortOrder: options.sortOrder,
        externalRequestIDs: options.externalRequestIds,
      });

      const requests = requestsFromDto(response.result);
      const cursor = response.cursor;

      return {
        requests,
        cursor: {
          nextCursor: cursor?.currentPage,
          hasMore: cursor?.currentPage !== undefined,
        },
      };
    });
  }

  /**
   * List requests pending approval.
   *
   * **Note:** Hash verification is NOT performed on list operations for performance.
   * Use `get()` to fetch individual requests with full verification.
   *
   * @param options - List options including filters and pagination
   * @returns Object containing requests array and cursor pagination
   * @throws {APIError} If API request fails
   *
   * @example
   * ```typescript
   * const result = await requestService.listForApproval({ limit: 50 });
   * console.log(`${result.requests.length} requests pending approval`);
   * ```
   */
  async listForApproval(
    options: ListRequestsForApprovalOptions = {}
  ): Promise<ListRequestsResult> {
    const limit = options.limit ?? 50;
    if (limit <= 0) {
      throw new Error("limit must be positive");
    }

    return this.execute(async () => {
      const response =
        await this.requestsApi.requestServiceGetRequestsForApprovalV2({
          cursorPageSize: String(limit),
          cursorCurrentPage: options.currentPage,
          cursorPageRequest: options.pageRequest,
          currencyID: options.currencyId,
          types: options.types,
          excludeTypes: options.excludeTypes,
          ids: options.ids?.map(String),
          sortOrder: options.sortOrder,
          externalRequestIDs: options.externalRequestIds,
        });

      const requests = requestsFromDto(response.result);
      const cursor = response.cursor;

      return {
        requests,
        cursor: {
          nextCursor: cursor?.currentPage,
          hasMore: cursor?.currentPage !== undefined,
        },
      };
    });
  }

  /**
   * Approve a single request with ECDSA signature.
   *
   * This is a convenience method that calls `approveRequests()` with a single request.
   *
   * @param request - The request to approve
   * @param privateKey - ECDSA private key for signing (P-256)
   * @param comment - Optional approval comment
   * @returns Number of requests successfully signed (0 or 1)
   * @throws {Error} If request is invalid
   * @throws {APIError} If API request fails
   *
   * @example
   * ```typescript
   * const count = await requestService.approveRequest(request, privateKey);
   * console.log(`Approved ${count} request`);
   * ```
   */
  async approveRequest(
    request: Request,
    privateKey: KeyObject,
    comment: string = "approving via taurus-protect-sdk-typescript"
  ): Promise<number> {
    return this.approveRequests([request], privateKey, comment);
  }

  /**
   * Approve multiple requests with ECDSA signature.
   *
   * **CRITICAL SECURITY:**
   * The requests are sorted by ID, and a JSON array of their hashes
   * is signed using the provided ECDSA private key. The signature
   * proves that the caller has the private key and authorizes
   * all the specified requests.
   *
   * @param requests - List of requests to approve
   * @param privateKey - ECDSA private key for signing (P-256)
   * @param comment - Optional approval comment
   * @returns Number of requests successfully signed
   * @throws {Error} If requests list is empty or invalid
   * @throws {APIError} If API request fails
   *
   * @example
   * ```typescript
   * // Load private key
   * const privateKey = decodePrivateKeyPem(pemString);
   *
   * // Approve multiple requests
   * const count = await requestService.approveRequests(
   *   pendingRequests,
   *   privateKey,
   *   'Batch approval for Q1 transfers'
   * );
   * console.log(`Approved ${count} requests`);
   * ```
   */
  async approveRequests(
    requests: Request[],
    privateKey: KeyObject,
    comment: string = "approving via taurus-protect-sdk-typescript"
  ): Promise<number> {
    if (!requests || requests.length === 0) {
      throw new Error("requests list cannot be empty");
    }
    if (!privateKey) {
      throw new Error("privateKey cannot be null or undefined");
    }

    // Validate all requests have metadata with hash
    for (const request of requests) {
      if (!request.metadata) {
        throw new Error(
          `Request ${request.id} metadata cannot be null or undefined`
        );
      }
      if (!request.metadata.hash) {
        throw new Error(
          `Request ${request.id} metadata hash cannot be null or empty`
        );
      }
    }

    return this.execute(async () => {
      // CRITICAL: Sort requests by ID (numeric sort)
      const sortedRequests = [...requests].sort((a, b) => a.id - b.id);

      // Build JSON array of hashes
      const hashes = sortedRequests.map((r) => r.metadata!.hash);
      const hashesJson = JSON.stringify(hashes);

      // Sign with ECDSA
      const signature = signData(privateKey, Buffer.from(hashesJson, "utf-8"));

      // Submit approval to API
      const response = await this.requestsApi.requestServiceApproveRequests({
        body: {
          ids: sortedRequests.map((r) => String(r.id)),
          signature,
          comment,
        },
      });

      // Parse signed count from response
      const signedRequests = response.signedRequests;
      if (signedRequests !== undefined) {
        const count = parseInt(signedRequests, 10);
        return isNaN(count) ? 0 : count;
      }
      return 0;
    });
  }

  /**
   * Reject a single request.
   *
   * @param requestId - The request ID to reject
   * @param comment - Rejection comment (required)
   * @throws {Error} If comment is empty
   * @throws {APIError} If API request fails
   *
   * @example
   * ```typescript
   * await requestService.rejectRequest(123, 'Amount exceeds daily limit');
   * ```
   */
  async rejectRequest(requestId: number, comment: string): Promise<void> {
    return this.rejectRequests([requestId], comment);
  }

  /**
   * Reject multiple requests.
   *
   * @param requestIds - List of request IDs to reject
   * @param comment - Rejection comment (required)
   * @throws {Error} If requestIds is empty or comment is empty
   * @throws {APIError} If API request fails
   *
   * @example
   * ```typescript
   * await requestService.rejectRequests([123, 456, 789], 'Rejected due to policy violation');
   * ```
   */
  async rejectRequests(requestIds: number[], comment: string): Promise<void> {
    if (!requestIds || requestIds.length === 0) {
      throw new Error("requestIds list cannot be empty");
    }
    if (!comment || comment.trim().length === 0) {
      throw new Error("comment is required and cannot be empty");
    }

    return this.execute(async () => {
      await this.requestsApi.requestServiceRejectRequests({
        body: {
          ids: requestIds.map(String),
          comment,
        },
      });
    });
  }

  /**
   * Create an internal transfer request between addresses.
   *
   * Internal transfers move funds between addresses within the same tenant.
   *
   * @param options - Transfer options
   * @returns The created request
   * @throws {Error} If arguments are invalid
   * @throws {APIError} If API request fails
   *
   * @example
   * ```typescript
   * const request = await requestService.createInternalTransferRequest({
   *   fromAddressId: 123,
   *   toAddressId: 456,
   *   amount: '1000000000000000000', // 1 ETH in wei
   * });
   * console.log(`Created request ${request.id}`);
   * ```
   */
  async createInternalTransferRequest(
    options: CreateInternalTransferOptions
  ): Promise<Request> {
    if (options.fromAddressId <= 0) {
      throw new Error("fromAddressId must be positive");
    }
    if (options.toAddressId <= 0) {
      throw new Error("toAddressId must be positive");
    }
    if (!options.amount || options.amount.trim().length === 0) {
      throw new Error("amount is required");
    }
    // Validate amount is a positive number
    const amountNum = parseFloat(options.amount);
    if (isNaN(amountNum) || amountNum <= 0) {
      throw new Error("amount must be a positive number");
    }

    return this.execute(async () => {
      const response =
        await this.requestsApi.requestServiceCreateOutgoingRequest({
          body: {
            amount: options.amount,
            fromAddressId: String(options.fromAddressId),
            toAddressId: String(options.toAddressId),
            comment: options.comment,
            externalRequestId: options.externalRequestId,
            gasLimit: options.gasLimit,
            feeLimit: options.feeLimit,
          },
        });

      const result = response.result;
      if (!result) {
        throw new ServerError("Failed to create request: no result returned");
      }

      const request = requestFromDto(result);
      if (!request) {
        throw new ServerError("Failed to create request: invalid response");
      }

      return request;
    });
  }

  /**
   * Create an external transfer request to a whitelisted address.
   *
   * External transfers move funds to addresses outside the tenant,
   * which must be whitelisted for security.
   *
   * @param options - Transfer options
   * @returns The created request
   * @throws {Error} If arguments are invalid
   * @throws {APIError} If API request fails
   *
   * @example
   * ```typescript
   * const request = await requestService.createExternalTransferRequest({
   *   fromAddressId: 123,
   *   toWhitelistedAddressId: 789,
   *   amount: '500000000000000000', // 0.5 ETH in wei
   * });
   * console.log(`Created request ${request.id}`);
   * ```
   */
  async createExternalTransferRequest(
    options: CreateExternalTransferOptions
  ): Promise<Request> {
    if (options.fromAddressId <= 0) {
      throw new Error("fromAddressId must be positive");
    }
    if (options.toWhitelistedAddressId <= 0) {
      throw new Error("toWhitelistedAddressId must be positive");
    }
    if (!options.amount || options.amount.trim().length === 0) {
      throw new Error("amount is required");
    }
    // Validate amount is a positive number
    const amountNum = parseFloat(options.amount);
    if (isNaN(amountNum) || amountNum <= 0) {
      throw new Error("amount must be a positive number");
    }

    return this.execute(async () => {
      const response =
        await this.requestsApi.requestServiceCreateOutgoingRequest({
          body: {
            amount: options.amount,
            fromAddressId: String(options.fromAddressId),
            toWhitelistedAddressId: String(options.toWhitelistedAddressId),
            comment: options.comment,
            externalRequestId: options.externalRequestId,
            gasLimit: options.gasLimit,
            feeLimit: options.feeLimit,
            destinationAddressMemo: options.destinationAddressMemo,
          },
        });

      const result = response.result;
      if (!result) {
        throw new ServerError("Failed to create request: no result returned");
      }

      const request = requestFromDto(result);
      if (!request) {
        throw new ServerError("Failed to create request: invalid response");
      }

      return request;
    });
  }

  /**
   * Create an internal transfer request from a wallet.
   *
   * Internal transfers from a wallet move funds from an omnibus wallet
   * to a specific address within the same tenant.
   *
   * @param options - Transfer options
   * @returns The created request
   * @throws {Error} If arguments are invalid
   * @throws {APIError} If API request fails
   *
   * @example
   * ```typescript
   * const request = await requestService.createInternalTransferFromWalletRequest({
   *   fromWalletId: 123,
   *   toAddressId: 456,
   *   amount: '1000000000000000000', // 1 ETH in wei
   * });
   * console.log(`Created request ${request.id}`);
   * ```
   */
  async createInternalTransferFromWalletRequest(
    options: CreateInternalTransferFromWalletOptions
  ): Promise<Request> {
    if (options.fromWalletId <= 0) {
      throw new Error("fromWalletId must be positive");
    }
    if (options.toAddressId <= 0) {
      throw new Error("toAddressId must be positive");
    }
    if (!options.amount || options.amount.trim().length === 0) {
      throw new Error("amount is required");
    }
    // Validate amount is a positive number
    const amountNum = parseFloat(options.amount);
    if (isNaN(amountNum) || amountNum <= 0) {
      throw new Error("amount must be a positive number");
    }

    return this.execute(async () => {
      const response =
        await this.requestsApi.requestServiceCreateOutgoingRequest({
          body: {
            amount: options.amount,
            fromWalletId: String(options.fromWalletId),
            toAddressId: String(options.toAddressId),
            comment: options.comment,
            externalRequestId: options.externalRequestId,
            gasLimit: options.gasLimit,
            feeLimit: options.feeLimit,
          },
        });

      const result = response.result;
      if (!result) {
        throw new ServerError("Failed to create request: no result returned");
      }

      const request = requestFromDto(result);
      if (!request) {
        throw new ServerError("Failed to create request: invalid response");
      }

      return request;
    });
  }

  /**
   * Create an external transfer request from a wallet to a whitelisted address.
   *
   * External transfers from a wallet move funds from an omnibus wallet
   * to addresses outside the tenant, which must be whitelisted for security.
   *
   * @param options - Transfer options
   * @returns The created request
   * @throws {Error} If arguments are invalid
   * @throws {APIError} If API request fails
   *
   * @example
   * ```typescript
   * const request = await requestService.createExternalTransferFromWalletRequest({
   *   fromWalletId: 123,
   *   toWhitelistedAddressId: 789,
   *   amount: '500000000000000000', // 0.5 ETH in wei
   * });
   * console.log(`Created request ${request.id}`);
   * ```
   */
  async createExternalTransferFromWalletRequest(
    options: CreateExternalTransferFromWalletOptions
  ): Promise<Request> {
    if (options.fromWalletId <= 0) {
      throw new Error("fromWalletId must be positive");
    }
    if (options.toWhitelistedAddressId <= 0) {
      throw new Error("toWhitelistedAddressId must be positive");
    }
    if (!options.amount || options.amount.trim().length === 0) {
      throw new Error("amount is required");
    }
    // Validate amount is a positive number
    const amountNum = parseFloat(options.amount);
    if (isNaN(amountNum) || amountNum <= 0) {
      throw new Error("amount must be a positive number");
    }

    return this.execute(async () => {
      const response =
        await this.requestsApi.requestServiceCreateOutgoingRequest({
          body: {
            amount: options.amount,
            fromWalletId: String(options.fromWalletId),
            toWhitelistedAddressId: String(options.toWhitelistedAddressId),
            comment: options.comment,
            externalRequestId: options.externalRequestId,
            gasLimit: options.gasLimit,
            feeLimit: options.feeLimit,
            destinationAddressMemo: options.destinationAddressMemo,
          },
        });

      const result = response.result;
      if (!result) {
        throw new ServerError("Failed to create request: no result returned");
      }

      const request = requestFromDto(result);
      if (!request) {
        throw new ServerError("Failed to create request: invalid response");
      }

      return request;
    });
  }

  /**
   * Create a cancel request for a pending transaction.
   *
   * Cancel requests attempt to cancel a pending transaction that has not
   * yet been broadcast to the blockchain by submitting a zero-value
   * transaction with the same nonce.
   *
   * @param addressId - The address ID of the pending transaction
   * @param nonce - The nonce of the transaction to cancel
   * @returns The created cancel request
   * @throws {Error} If arguments are invalid
   * @throws {APIError} If API request fails
   *
   * @example
   * ```typescript
   * const cancelRequest = await requestService.createCancelRequest(123, 42n);
   * console.log(`Created cancel request ${cancelRequest.id}`);
   * ```
   */
  async createCancelRequest(
    addressId: number,
    nonce: bigint | number
  ): Promise<Request> {
    if (addressId <= 0) {
      throw new Error("addressId must be positive");
    }
    const nonceValue = BigInt(nonce);
    if (nonceValue < 0n) {
      throw new Error("nonce cannot be negative");
    }

    return this.execute(async () => {
      const response =
        await this.requestsApi.requestServiceCreateOutgoingCancelRequest({
          body: {
            addressId: String(addressId),
            nonce: String(nonceValue),
          },
        });

      const result = response.result;
      if (!result) {
        throw new ServerError(
          "Failed to create cancel request: no result returned"
        );
      }

      const request = requestFromDto(result);
      if (!request) {
        throw new ServerError(
          "Failed to create cancel request: invalid response"
        );
      }

      return request;
    });
  }

  /**
   * Create an incoming request from an exchange.
   *
   * Incoming requests are used to transfer funds from an external exchange
   * account to an internal address.
   *
   * @param options - Incoming request options
   * @returns The created request
   * @throws {Error} If arguments are invalid
   * @throws {APIError} If API request fails
   *
   * @example
   * ```typescript
   * const request = await requestService.createIncomingRequest({
   *   fromExchangeId: 1,
   *   toAddressId: 456,
   *   amount: '1000000000000000000', // 1 ETH in wei
   * });
   * console.log(`Created request ${request.id}`);
   * ```
   */
  async createIncomingRequest(
    options: CreateIncomingRequestOptions
  ): Promise<Request> {
    if (options.fromExchangeId <= 0) {
      throw new Error("fromExchangeId must be positive");
    }
    if (options.toAddressId <= 0) {
      throw new Error("toAddressId must be positive");
    }
    if (!options.amount || options.amount.trim().length === 0) {
      throw new Error("amount is required");
    }
    // Validate amount is a positive number
    const amountNum = parseFloat(options.amount);
    if (isNaN(amountNum) || amountNum <= 0) {
      throw new Error("amount must be a positive number");
    }

    return this.execute(async () => {
      const response =
        await this.requestsApi.requestServiceCreateIncomingRequest({
          body: {
            amount: options.amount,
            fromExchangeId: String(options.fromExchangeId),
            toAddressId: String(options.toAddressId),
            comment: options.comment,
            externalRequestId: options.externalRequestId,
          },
        });

      const result = response.result;
      if (!result) {
        throw new ServerError("Failed to create request: no result returned");
      }

      const request = requestFromDto(result);
      if (!request) {
        throw new ServerError("Failed to create request: invalid response");
      }

      return request;
    });
  }

  /**
   * Verify the hash of a request using constant-time comparison.
   *
   * **CRITICAL SECURITY:**
   * This method computes the SHA-256 hash of the request's payload
   * and compares it to the provided hash using constant-time comparison
   * to prevent timing attacks.
   *
   * @param request - The request to verify
   * @throws {IntegrityError} If hash verification fails
   */
  private verifyRequestHash(request: Request): void {
    // If no metadata, nothing to verify
    if (!request.metadata) {
      return;
    }

    const providedHash = request.metadata.hash;
    const payloadAsString = request.metadata.payloadAsString;

    // If no hash and no payload, nothing to verify
    if (!providedHash && !payloadAsString) {
      // Note: This is acceptable - some requests may not have verification data
      return;
    }

    // SECURITY: If hash exists but payload is missing, we MUST fail
    // This prevents attackers from bypassing verification by stripping the payload
    if (!payloadAsString) {
      if (providedHash) {
        throw new IntegrityError(
          "Request hash verification failed: hash exists but payload is missing"
        );
      }
      return;
    }

    // Compute hash of the payload
    const computedHash = calculateHexHash(payloadAsString);

    // Explicit null check before constant-time comparison
    if (!providedHash) {
      throw new IntegrityError(
        "Request hash verification failed: provided hash is null"
      );
    }

    // CRITICAL: Use constant-time comparison to prevent timing attacks
    if (!constantTimeCompare(computedHash, providedHash)) {
      throw new IntegrityError(
        `Request hash mismatch: computed=${computedHash}, provided=${providedHash}`
      );
    }
  }
}
