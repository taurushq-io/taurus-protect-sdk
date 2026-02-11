/**
 * Exception classes for Taurus-PROTECT SDK.
 *
 * This module provides a comprehensive error hierarchy for handling API errors,
 * security/integrity errors, and configuration errors.
 */

/**
 * Base exception for all Taurus-PROTECT API errors.
 *
 * This exception captures HTTP status codes, error messages, and provides
 * helper methods for common error handling patterns.
 *
 * Specific subclasses provide more context for different error categories:
 * - {@link ValidationError} - 400 Bad Request errors
 * - {@link AuthenticationError} - 401 Unauthorized errors
 * - {@link AuthorizationError} - 403 Forbidden errors
 * - {@link NotFoundError} - 404 Not Found errors
 * - {@link RateLimitError} - 429 Too Many Requests errors
 * - {@link ServerError} - 5xx Server errors
 */
export class APIError extends Error {
  /**
   * HTTP status code of the error response.
   */
  public readonly statusCode: number;

  /**
   * Application-specific error code from the API.
   */
  public readonly errorCode: string | undefined;

  /**
   * Raw response body, if available.
   */
  public readonly body: unknown | undefined;

  /**
   * Suggested retry delay in milliseconds for rate-limited requests.
   * Only set for RateLimitError.
   */
  public readonly retryAfterMs: number | undefined;

  /**
   * The underlying error that caused this exception, if any.
   */
  public readonly cause: Error | undefined;

  /**
   * Constructs an APIError with the specified details.
   *
   * @param statusCode - HTTP status code
   * @param errorCode - Application-specific error code from the API
   * @param message - Human-readable error message
   * @param body - Raw response body
   * @param retryAfterMs - Suggested retry delay in milliseconds
   * @param cause - The underlying error that caused this exception
   */
  constructor(
    statusCode: number,
    errorCode: string | undefined,
    message: string,
    body?: unknown,
    retryAfterMs?: number,
    cause?: Error
  ) {
    super(message);
    this.name = "APIError";
    this.statusCode = statusCode;
    this.errorCode = errorCode;
    this.body = body;
    this.retryAfterMs = retryAfterMs;
    this.cause = cause;

    // Maintains proper stack trace for where our error was thrown (only available on V8)
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, this.constructor);
    }
  }

  /**
   * Determines if this error is potentially retryable.
   *
   * Returns true for:
   * - 429 Too Many Requests (rate limited)
   * - All 5xx Server errors (transient failures)
   *
   * @returns true if the request might succeed on retry
   */
  isRetryable(): boolean {
    return this.statusCode === 429 || this.statusCode >= 500;
  }

  /**
   * Determines if this is a client error (4xx).
   *
   * @returns true if this is a client error
   */
  isClientError(): boolean {
    return this.statusCode >= 400 && this.statusCode < 500;
  }

  /**
   * Determines if this is a server error (5xx).
   *
   * @returns true if this is a server error
   */
  isServerError(): boolean {
    return this.statusCode >= 500;
  }

  /**
   * Gets the suggested retry delay in milliseconds.
   *
   * For rate limit errors, this may return a specific delay from the server.
   * For server errors (5xx), returns a default backoff value.
   * For non-retryable errors, returns undefined.
   *
   * @returns suggested delay in milliseconds, or undefined if not retryable
   */
  suggestedRetryDelayMs(): number | undefined {
    if (this.statusCode === 429) {
      return this.retryAfterMs ?? 1000; // Default 1 second for rate limits
    }
    if (this.statusCode >= 500) {
      return 5000; // Default 5 seconds for server errors
    }
    return undefined;
  }

  /**
   * Returns a string representation of this error.
   */
  override toString(): string {
    const parts = [this.message];
    if (this.statusCode) {
      parts.push(`(HTTP ${this.statusCode})`);
    }
    if (this.errorCode) {
      parts.push(`[${this.errorCode}]`);
    }
    return parts.join(" ");
  }
}

/**
 * 400 Bad Request - Input validation failed.
 *
 * This exception indicates that the request was malformed or contained
 * invalid parameters that the server could not process.
 */
export class ValidationError extends APIError {
  /**
   * Constructs a ValidationError with the specified details.
   *
   * @param message - Human-readable error message
   * @param errorCode - Application-specific error code
   * @param body - Raw response body
   * @param cause - The underlying error
   */
  constructor(
    message: string,
    errorCode?: string,
    body?: unknown,
    cause?: Error
  ) {
    super(400, errorCode, message, body, undefined, cause);
    this.name = "ValidationError";
  }
}

/**
 * 401 Unauthorized - Invalid or missing credentials.
 *
 * This exception indicates that authentication failed, typically due to
 * invalid API credentials or an expired/invalid authentication token.
 */
export class AuthenticationError extends APIError {
  /**
   * Constructs an AuthenticationError with the specified details.
   *
   * @param message - Human-readable error message
   * @param errorCode - Application-specific error code
   * @param body - Raw response body
   * @param cause - The underlying error
   */
  constructor(
    message: string,
    errorCode?: string,
    body?: unknown,
    cause?: Error
  ) {
    super(401, errorCode, message, body, undefined, cause);
    this.name = "AuthenticationError";
  }
}

/**
 * 403 Forbidden - Insufficient permissions.
 *
 * This exception indicates that the authenticated user does not have
 * sufficient permissions to perform the requested operation.
 */
export class AuthorizationError extends APIError {
  /**
   * Constructs an AuthorizationError with the specified details.
   *
   * @param message - Human-readable error message
   * @param errorCode - Application-specific error code
   * @param body - Raw response body
   * @param cause - The underlying error
   */
  constructor(
    message: string,
    errorCode?: string,
    body?: unknown,
    cause?: Error
  ) {
    super(403, errorCode, message, body, undefined, cause);
    this.name = "AuthorizationError";
  }
}

/**
 * 404 Not Found - Resource does not exist.
 *
 * This exception indicates that the requested resource could not be found.
 */
export class NotFoundError extends APIError {
  /**
   * Constructs a NotFoundError with the specified details.
   *
   * @param message - Human-readable error message
   * @param errorCode - Application-specific error code
   * @param body - Raw response body
   * @param cause - The underlying error
   */
  constructor(
    message: string,
    errorCode?: string,
    body?: unknown,
    cause?: Error
  ) {
    super(404, errorCode, message, body, undefined, cause);
    this.name = "NotFoundError";
  }
}

/**
 * 429 Too Many Requests - Rate limit exceeded.
 *
 * This exception indicates that the client has exceeded the API rate limit
 * and should wait before retrying. Use {@link suggestedRetryDelayMs} to get
 * the recommended wait time.
 *
 * @example
 * ```typescript
 * try {
 *   await client.walletService.getWallets();
 * } catch (error) {
 *   if (error instanceof RateLimitError) {
 *     const delay = error.suggestedRetryDelayMs() ?? 1000;
 *     await new Promise(resolve => setTimeout(resolve, delay));
 *     // Retry the request
 *   }
 * }
 * ```
 */
export class RateLimitError extends APIError {
  /**
   * Constructs a RateLimitError with the specified details.
   *
   * @param message - Human-readable error message
   * @param retryAfterMs - Suggested retry delay in milliseconds
   * @param errorCode - Application-specific error code
   * @param body - Raw response body
   * @param cause - The underlying error
   */
  constructor(
    message: string,
    retryAfterMs?: number,
    errorCode?: string,
    body?: unknown,
    cause?: Error
  ) {
    super(429, errorCode, message, body, retryAfterMs ?? 1000, cause);
    this.name = "RateLimitError";
  }

  /**
   * Rate limit errors are always retryable.
   */
  override isRetryable(): boolean {
    return true;
  }

  /**
   * Returns the retry delay, always defined for rate limit errors.
   */
  override suggestedRetryDelayMs(): number {
    return this.retryAfterMs ?? 1000;
  }
}

/**
 * 5xx Server Error - Internal server error.
 *
 * This exception indicates a server-side error that may be transient.
 * Errors with status codes 503 and 504 are considered retryable.
 */
export class ServerError extends APIError {
  /**
   * Constructs a ServerError with the specified details.
   *
   * @param message - Human-readable error message
   * @param statusCode - HTTP status code (default: 500)
   * @param errorCode - Application-specific error code
   * @param body - Raw response body
   * @param cause - The underlying error
   */
  constructor(
    message: string,
    statusCode: number = 500,
    errorCode?: string,
    body?: unknown,
    cause?: Error
  ) {
    super(statusCode, errorCode, message, body, undefined, cause);
    this.name = "ServerError";
  }
}

/**
 * Cryptographic integrity verification failure.
 *
 * This exception indicates that hash or signature verification failed.
 * This is a security-critical error and should NEVER be retried.
 *
 * Common causes include:
 * - Request hash mismatch
 * - Invalid address signature
 * - Invalid whitelist signature
 * - SuperAdmin signature verification failure
 * - Tampered request or response data
 *
 * This exception typically indicates a serious security issue that should be
 * investigated and not simply retried.
 */
export class IntegrityError extends Error {
  /**
   * Constructs an IntegrityError with the specified message.
   *
   * @param message - The detail message describing the verification failure
   * @param cause - The underlying error that caused this exception
   */
  constructor(message: string, cause?: Error) {
    super(message);
    this.name = "IntegrityError";
    this.cause = cause;

    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, this.constructor);
    }
  }

  /**
   * Returns a string representation of this error.
   */
  override toString(): string {
    return `IntegrityError: ${this.message}`;
  }
}

/**
 * Whitelisted address/asset verification failure.
 *
 * This exception indicates that whitelist verification failed, which may be due to:
 * - Invalid hash in metadata
 * - Insufficient governance rule signatures
 * - Missing required fields in payload
 * - Invalid address format for the target blockchain
 * - Whitelist entry not found
 */
export class WhitelistError extends Error {
  /**
   * Constructs a WhitelistError with the specified message.
   *
   * @param message - The detail message describing the whitelist operation failure
   * @param cause - The underlying error that caused this exception
   */
  constructor(message: string, cause?: Error) {
    super(message);
    this.name = "WhitelistError";
    this.cause = cause;

    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, this.constructor);
    }
  }

  /**
   * Returns a string representation of this error.
   */
  override toString(): string {
    return `WhitelistError: ${this.message}`;
  }
}

/**
 * SDK configuration error.
 *
 * This exception indicates invalid SDK configuration such as:
 * - Missing required credentials
 * - Invalid API host URL
 * - Invalid SuperAdmin key format
 * - Invalid min_valid_signatures value
 */
export class ConfigurationError extends Error {
  /**
   * Constructs a ConfigurationError with the specified message.
   *
   * @param message - The detail message describing the configuration error
   * @param cause - The underlying error that caused this exception
   */
  constructor(message: string, cause?: Error) {
    super(message);
    this.name = "ConfigurationError";
    this.cause = cause;

    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, this.constructor);
    }
  }
}

/**
 * Request metadata error.
 *
 * This exception is thrown when request metadata cannot be parsed or extracted.
 * Common causes include:
 * - Missing required fields in the metadata payload
 * - Malformed JSON structure
 * - Type mismatches when parsing values
 */
export class RequestMetadataError extends Error {
  /**
   * Constructs a RequestMetadataError with the specified message.
   *
   * @param message - The detail message describing what metadata field
   *                  could not be extracted or parsed
   * @param cause - The underlying error that caused this exception
   */
  constructor(message: string, cause?: Error) {
    super(message);
    this.name = "RequestMetadataError";
    this.cause = cause;

    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, this.constructor);
    }
  }

  /**
   * Returns a string representation of this error.
   */
  override toString(): string {
    return `RequestMetadataError: ${this.message}`;
  }
}

/**
 * Maps an HTTP status code to the appropriate error class.
 *
 * @param statusCode - HTTP status code
 * @param message - Error message
 * @param errorCode - Optional application error code
 * @param body - Optional response body
 * @param retryAfterMs - Optional retry delay for rate limits (in milliseconds)
 * @param cause - Optional underlying error
 * @returns Appropriate APIError subclass
 *
 * @example
 * ```typescript
 * const error = mapHttpError(404, "Wallet not found", "WALLET_NOT_FOUND");
 * if (error instanceof NotFoundError) {
 *   console.log("Resource not found");
 * }
 * ```
 */
export function mapHttpError(
  statusCode: number,
  message: string,
  errorCode?: string,
  body?: unknown,
  retryAfterMs?: number,
  cause?: Error
): APIError {
  switch (statusCode) {
    case 400:
      return new ValidationError(message, errorCode, body, cause);
    case 401:
      return new AuthenticationError(message, errorCode, body, cause);
    case 403:
      return new AuthorizationError(message, errorCode, body, cause);
    case 404:
      return new NotFoundError(message, errorCode, body, cause);
    case 429:
      return new RateLimitError(message, retryAfterMs, errorCode, body, cause);
    default:
      if (statusCode >= 500) {
        return new ServerError(message, statusCode, errorCode, body, cause);
      }
      return new APIError(statusCode, errorCode, message, body, undefined, cause);
  }
}
