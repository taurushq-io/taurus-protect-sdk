/**
 * Transaction service for Taurus-PROTECT SDK.
 *
 * Provides methods for retrieving and filtering blockchain transactions.
 */

import { NotFoundError, ValidationError } from "../errors";
import type { TransactionsApi } from "../internal/openapi/apis/TransactionsApi";
import { transactionsFromDto } from "../mappers/transaction";
import {
  buildOffsetPagination,
  offsetRequest,
  parseReplyCount,
  resolvePageSize,
  type OffsetPageOptions,
  type PaginatedResult,
} from "../models/pagination";
import type {
  ExportTransactionsOptions,
  ExportTransactionsResult,
  ListTransactionsOptions,
  Transaction,
} from "../models/transaction";
import { BaseService } from "./base";
import { offsetQuery } from "./paging";

/**
 * Service for retrieving and filtering blockchain transactions.
 *
 * Transactions represent the movement of cryptocurrency on the blockchain
 * and can be either incoming (received) or outgoing (sent).
 *
 * @example
 * ```typescript
 * // Get a single transaction by ID
 * const tx = await transactionService.get("12345");
 * console.log(`${tx.hash}: ${tx.amount} ${tx.currency}`);
 *
 * // List transactions with filters
 * const result = await transactionService.list({
 *   currency: "ETH",
 *   direction: "incoming",
 *   limit: 50,
 * });
 * for (const tx of result.items) {
 *   console.log(`${tx.hash}: ${tx.amount}`);
 * }
 *
 * // List transactions by request ID
 * const txs = await transactionService.listByRequest("req-123");
 * ```
 */
export class TransactionService extends BaseService {
  private readonly transactionsApi: TransactionsApi;

  /**
   * Creates a new TransactionService instance.
   *
   * @param transactionsApi - The TransactionsApi instance from the OpenAPI client
   */
  constructor(transactionsApi: TransactionsApi) {
    super();
    this.transactionsApi = transactionsApi;
  }

  /**
   * Gets a transaction by ID.
   *
   * @param transactionId - The transaction ID to retrieve
   * @returns The transaction
   * @throws {@link ValidationError} If transactionId is invalid
   * @throws {@link NotFoundError} If transaction not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const tx = await transactionService.get("12345");
   * console.log(`Hash: ${tx.hash}, Amount: ${tx.amount}`);
   * ```
   */
  async get(transactionId: string): Promise<Transaction> {
    if (!transactionId || transactionId.trim() === "") {
      throw new ValidationError("transactionId is required");
    }

    return this.execute(async () => {
      const response = await this.transactionsApi.transactionServiceGetTransactions({
        ids: [transactionId],
        limit: "1",
      });

      const transactions = transactionsFromDto(response.result);
      if (transactions.length === 0) {
        throw new NotFoundError(`Transaction ${transactionId} not found`);
      }

      return transactions[0];
    });
  }

  /**
   * Gets a transaction by its blockchain hash.
   *
   * @param txHash - The blockchain transaction hash
   * @returns The transaction
   * @throws {@link ValidationError} If txHash is invalid
   * @throws {@link NotFoundError} If transaction not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const tx = await transactionService.getByHash("0x1234...");
   * console.log(`ID: ${tx.id}, Amount: ${tx.amount}`);
   * ```
   */
  async getByHash(txHash: string): Promise<Transaction> {
    if (!txHash || txHash.trim() === "") {
      throw new ValidationError("txHash is required");
    }

    return this.execute(async () => {
      const response = await this.transactionsApi.transactionServiceGetTransactions({
        hashes: [txHash],
        limit: "1",
      });

      const transactions = transactionsFromDto(response.result);
      if (transactions.length === 0) {
        throw new NotFoundError(`Transaction with hash ${txHash} not found`);
      }

      return transactions[0];
    });
  }

  /**
   * Lists transactions with pagination and optional filtering.
   *
   * @param options - Optional filtering and pagination options
   * @returns Paginated result containing transactions and pagination info
   * @throws {@link ValidationError} If limit or offset are invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // List first 50 transactions
   * const result = await transactionService.list({ limit: 50, offset: 0 });
   * console.log(`Found ${result.pagination.totalItems} transactions`);
   *
   * // List with filters
   * const filtered = await transactionService.list({
   *   currency: "ETH",
   *   direction: "incoming",
   *   blockchain: "ETH",
   *   network: "mainnet",
   * });
   * ```
   */
  async list(
    options?: ListTransactionsOptions
  ): Promise<PaginatedResult<Transaction>> {
    const page = offsetRequest(options);

    return this.execute(async () => {
      const response = await this.transactionsApi.transactionServiceGetTransactions({
        currency: options?.currency,
        direction: options?.direction,
        query: options?.query,
        ...offsetQuery(page),
        from: options?.fromDate,
        to: options?.toDate,
        type: options?.type,
        source: options?.source,
        destination: options?.destination,
        ids: options?.ids,
        blockchain: options?.blockchain,
        network: options?.network,
        fromBlockNumber: options?.fromBlockNumber,
        toBlockNumber: options?.toBlockNumber,
        hashes: options?.hashes,
        address: options?.address,
        amountAbove: options?.amountAbove,
        excludeUnknownSourceDestination: options?.excludeUnknownSourceDestination,
        customerId: options?.customerId,
      });

      const rows = response.result ?? [];
      return {
        items: transactionsFromDto(rows),
        pagination: buildOffsetPagination("plus_rows", page, response, rows.length),
      };
    });
  }

  /**
   * Lists transactions associated with a specific request ID.
   *
   * @param requestId - The request ID to filter by
   * @param options - Optional pagination options
   * @returns Paginated result containing transactions
   * @throws {@link ValidationError} If requestId is invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const result = await transactionService.listByRequest("req-123");
   * for (const tx of result.items) {
   *   console.log(`${tx.hash}: ${tx.status}`);
   * }
   * ```
   */
  async listByRequest(
    requestId: string,
    options?: OffsetPageOptions
  ): Promise<PaginatedResult<Transaction>> {
    if (!requestId || requestId.trim() === "") {
      throw new ValidationError("requestId is required");
    }
    const page = offsetRequest(options);

    return this.execute(async () => {
      // Note: The API doesn't have a direct requestId filter, so we use transactionIds
      // which maps to the internal transaction IDs associated with requests
      const response = await this.transactionsApi.transactionServiceGetTransactions({
        transactionIds: [requestId],
        ...offsetQuery(page),
      });

      const rows = response.result ?? [];
      return {
        items: transactionsFromDto(rows),
        pagination: buildOffsetPagination("plus_rows", page, response, rows.length),
      };
    });
  }

  /**
   * Lists transactions for a specific blockchain address.
   *
   * @param address - The blockchain address
   * @param options - Optional pagination options
   * @returns Paginated result containing transactions
   * @throws {@link ValidationError} If address is invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const result = await transactionService.listByAddress("0x1234...");
   * for (const tx of result.items) {
   *   console.log(`${tx.hash}: ${tx.amount}`);
   * }
   * ```
   */
  async listByAddress(
    address: string,
    options?: OffsetPageOptions
  ): Promise<PaginatedResult<Transaction>> {
    if (!address || address.trim() === "") {
      throw new ValidationError("address is required");
    }
    const page = offsetRequest(options);

    return this.execute(async () => {
      const response = await this.transactionsApi.transactionServiceGetTransactions({
        address,
        ...offsetQuery(page),
      });

      const rows = response.result ?? [];
      return {
        items: transactionsFromDto(rows),
        pagination: buildOffsetPagination("plus_rows", page, response, rows.length),
      };
    });
  }

  /**
   * Exports transactions as a formatted string (JSON or CSV).
   *
   * The export cannot page — the server always starts at the first matching row — so it
   * takes a `limit` (default 20, no SDK maximum) and no offset. The result's `totalItems`
   * tells a truncated export apart from a complete one.
   *
   * @param options - Filters, format and `limit`
   * @returns The export payload and the server's total
   * @throws {@link ValidationError} If the limit is negative
   * @throws {@link PaginationError} If the reply's total is not a canonical count
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const exported = await transactionService.exportTransactions({
   *   currency: 'ETH',
   *   fromDate: new Date('2024-01-01'),
   *   format: 'csv',
   *   limit: 1000,
   * });
   * console.log(exported.data);
   * console.log(`${exported.totalItems} transactions match`);
   * ```
   */
  async exportTransactions(
    options?: ExportTransactionsOptions
  ): Promise<ExportTransactionsResult> {
    const limit = resolvePageSize(options?.limit, "limit", "export");

    return this.execute(async () => {
      const response = await this.transactionsApi.transactionServiceExportTransactions({
        from: options?.fromDate,
        to: options?.toDate,
        currency: options?.currency,
        direction: options?.direction,
        limit: String(limit),
        format: options?.format,
        blockchain: options?.blockchain,
        network: options?.network,
        query: options?.query,
        address: options?.address,
      });
      return {
        data: response.result ?? "",
        totalItems: parseReplyCount(response.totalItems, "totalItems"),
      };
    });
  }
}
