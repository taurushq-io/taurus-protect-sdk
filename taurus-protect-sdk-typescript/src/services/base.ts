/**
 * Base service class for all Taurus-PROTECT services.
 *
 * Provides common functionality for error handling and response mapping.
 */

import { APIError, IntegrityError, ServerError, mapHttpError } from "../errors";
import { ResponseError } from "../internal/openapi/runtime";

/**
 * Error details extracted from API response.
 */
interface ErrorDetails {
  statusCode: number;
  message: string;
  errorCode?: string;
  body?: unknown;
  retryAfterMs?: number;
}

/**
 * Base class for all service implementations.
 *
 * Provides:
 * - Error handling and mapping from OpenAPI ResponseError to SDK errors
 * - Common utility methods for executing API calls
 *
 * @example
 * ```typescript
 * class WalletService extends BaseService {
 *   private api: WalletsApi;
 *
 *   constructor(api: WalletsApi) {
 *     super();
 *     this.api = api;
 *   }
 *
 *   async getWallet(walletId: string): Promise<Wallet> {
 *     return this.execute(async () => {
 *       const dto = await this.api.getWallet({ walletId });
 *       return mapWalletFromDto(dto);
 *     });
 *   }
 * }
 * ```
 */
export abstract class BaseService {
  /**
   * Handles API errors, converting OpenAPI ResponseError to appropriate SDK errors.
   *
   * @param error - The error thrown by the OpenAPI client
   * @throws APIError or appropriate subclass
   */
  protected handleError(error: unknown): never {
    if (error instanceof ResponseError) {
      throw this.mapResponseError(error);
    }
    // Pass through SDK errors without wrapping
    if (error instanceof APIError) {
      throw error;
    }
    if (error instanceof IntegrityError) {
      throw error;
    }
    if (error instanceof Error) {
      throw new ServerError(error.message, 500, undefined, undefined, error);
    }
    throw new ServerError("Unknown error occurred", 500);
  }

  /**
   * Maps a ResponseError to the appropriate APIError subclass.
   *
   * @param error - The ResponseError from the OpenAPI client
   * @returns The appropriate APIError subclass based on status code
   */
  private mapResponseError(error: ResponseError): APIError {
    const response = error.response;
    const statusCode = response.status;

    // Build error details from response
    const errorDetails: ErrorDetails = {
      statusCode,
      message: error.message ?? `Request failed with status ${statusCode}`,
    };

    // Check for Retry-After header (for rate limiting)
    const retryAfter = response.headers.get("Retry-After");
    if (retryAfter) {
      const seconds = parseInt(retryAfter, 10);
      if (!isNaN(seconds)) {
        errorDetails.retryAfterMs = seconds * 1000;
      }
    }

    return mapHttpError(
      errorDetails.statusCode,
      errorDetails.message,
      errorDetails.errorCode,
      errorDetails.body,
      errorDetails.retryAfterMs
    );
  }

  /**
   * Executes an API call with error handling.
   *
   * Wraps the API call in a try-catch block and converts any errors
   * to the appropriate SDK error types.
   *
   * @param apiCall - Async function that makes the API call
   * @returns The result of the API call
   * @throws APIError on failure
   *
   * @example
   * ```typescript
   * const wallet = await this.execute(async () => {
   *   const dto = await this.api.getWallet({ walletId });
   *   return mapWalletFromDto(dto);
   * });
   * ```
   */
  protected async execute<T>(apiCall: () => Promise<T>): Promise<T> {
    try {
      return await apiCall();
    } catch (error) {
      this.handleError(error);
    }
  }
}
