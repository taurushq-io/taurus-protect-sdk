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
 * The flat error payload Taurus-PROTECT returns on any status >= 400. "error" is a
 * string (the gRPC status mnemonic), not a nested object; error_code is empty on
 * some paths.
 */
interface ErrorBody {
  error?: string;
  message?: string;
  code?: number;
  error_code?: string;
}

/**
 * Reads and parses the error body, tolerating a non-JSON or already-consumed body.
 *
 * @param response - The error response
 * @returns The parsed body, or undefined when it is unavailable or not an object
 */
async function readErrorBody(response: Response): Promise<ErrorBody | undefined> {
  try {
    const text = await response.clone().text();
    if (text.length === 0) {
      return undefined;
    }
    const parsed: unknown = JSON.parse(text);
    if (typeof parsed !== "object" || parsed === null) {
      return undefined;
    }
    return parsed as ErrorBody;
  } catch {
    return undefined;
  }
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
  protected async handleError(error: unknown): Promise<never> {
    if (error instanceof ResponseError) {
      throw await this.mapResponseError(error);
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
  private async mapResponseError(error: ResponseError): Promise<APIError> {
    const response = error.response;
    const statusCode = response.status;

    // Build error details from response
    const errorDetails: ErrorDetails = {
      statusCode,
      message: error.message ?? `Request failed with status ${statusCode}`,
    };

    // The server's own message and error code live in the body; without reading it
    // the caller only ever sees "Response returned an error code".
    const body = await readErrorBody(response);
    if (body !== undefined) {
      errorDetails.body = body;
      if (typeof body.message === "string" && body.message.length > 0) {
        errorDetails.message = body.message;
      }
      if (typeof body.error_code === "string" && body.error_code.length > 0) {
        errorDetails.errorCode = body.error_code;
      }
    }

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
      // `return` rather than a bare call: handleError is async now, so an
      // un-awaited call would leave a floating promise and return undefined.
      return await this.handleError(error);
    }
  }
}
