/**
 * TPV1-HMAC-SHA256 authentication for Taurus-PROTECT API.
 *
 * Generates Authorization headers for API requests.
 */

import * as crypto from "crypto";

import { calculateBase64Hmac } from "./hashing";

/**
 * TPV1-HMAC-SHA256 authentication handler.
 *
 * Generates Authorization headers for Taurus-PROTECT API requests.
 *
 * @example
 * ```typescript
 * const auth = new TPV1Auth("api-key", "hex-encoded-secret");
 * const header = auth.signRequest("POST", "api.example.com", "/v1/wallets");
 * // header = "TPV1-HMAC-SHA256 ApiKey=... Nonce=... Timestamp=... Signature=..."
 *
 * // When done, securely wipe the secret
 * auth.close();
 * ```
 */
export class TPV1Auth {
  private readonly apiKey: string;
  private secret: Uint8Array;
  private closed: boolean = false;

  /**
   * Initialize TPV1 authentication.
   *
   * @param apiKey - The API key
   * @param apiSecretHex - The API secret as a hex-encoded string
   * @throws Error if apiKey or apiSecretHex is empty
   */
  constructor(apiKey: string, apiSecretHex: string) {
    if (!apiKey) {
      throw new Error("apiKey cannot be empty");
    }
    if (!apiSecretHex) {
      throw new Error("apiSecret cannot be empty");
    }

    this.apiKey = apiKey;
    this.secret = Buffer.from(apiSecretHex, "hex");
  }

  /**
   * Securely wipe the API secret from memory.
   *
   * After calling this method, the auth instance cannot be used to sign requests.
   *
   * **SECURITY NOTE - Memory Wiping Limitations in JavaScript:**
   * Due to JavaScript's garbage collector, copies of the secret may still exist
   * in memory even after this method is called. The V8 engine and other JavaScript
   * runtimes may create intermediate string copies, cache values, or defer garbage
   * collection indefinitely. This wiping is best-effort and provides defense-in-depth,
   * but cannot guarantee complete removal of sensitive data from memory.
   *
   * For environments with strict memory security requirements, consider:
   * - Using short-lived process isolation for sensitive operations
   * - Implementing HSM or secure enclave solutions for key management
   * - Avoiding long-lived processes that handle secrets
   */
  close(): void {
    if (!this.closed) {
      // Best-effort: Overwrite secret buffer with zeros
      // Note: This cannot prevent V8/GC from retaining copies elsewhere
      if (this.secret instanceof Buffer) {
        this.secret.fill(0);
      } else {
        for (let i = 0; i < this.secret.length; i++) {
          this.secret[i] = 0;
        }
      }
      this.closed = true;
    }
  }

  /**
   * Generate TPV1 Authorization header for a request.
   *
   * Message format: "TPV1 ApiKey Nonce Timestamp METHOD host path [query] [content-type] [body]"
   * Header format: "TPV1-HMAC-SHA256 ApiKey=X Nonce=Y Timestamp=Z Signature=W"
   *
   * @param method - HTTP method (GET, POST, etc.)
   * @param host - Request host (e.g., "api.protect.taurushq.com")
   * @param path - Request path (e.g., "/v1/wallets")
   * @param query - Query string without leading "?" (optional)
   * @param contentType - Content-Type header value (optional)
   * @param body - Request body as string (optional)
   * @returns The Authorization header value
   * @throws Error if auth has been closed
   *
   * @example
   * ```typescript
   * // Simple GET request
   * const header = auth.signRequest("GET", "api.protect.taurushq.com", "/v1/wallets");
   *
   * // POST with body
   * const header = auth.signRequest(
   *   "POST",
   *   "api.protect.taurushq.com",
   *   "/v1/wallets",
   *   undefined,
   *   "application/json",
   *   '{"name":"my-wallet"}'
   * );
   * ```
   */
  signRequest(
    method: string,
    host: string,
    path: string,
    query?: string,
    contentType?: string,
    body?: string
  ): string {
    if (this.closed) {
      throw new Error("TPV1Auth has been closed");
    }

    const nonce = this.generateNonce();
    const timestamp = Date.now();

    // Build message from non-empty parts
    // Format: "TPV1 ApiKey Nonce Timestamp METHOD host path [query] [content-type] [body]"
    const parts: string[] = [
      "TPV1",
      this.apiKey,
      nonce,
      timestamp.toString(),
      method.toUpperCase(),
      host,
      path,
    ];

    // Add optional parts only if they have values
    if (query) {
      parts.push(query);
    }
    if (contentType) {
      parts.push(contentType);
    }
    if (body) {
      parts.push(body);
    }

    const message = parts.join(" ");

    // Calculate HMAC-SHA256 signature
    const signature = calculateBase64Hmac(this.secret, message);

    // Format: "TPV1-HMAC-SHA256 ApiKey=X Nonce=Y Timestamp=Z Signature=W"
    return `TPV1-HMAC-SHA256 ApiKey=${this.apiKey} Nonce=${nonce} Timestamp=${timestamp} Signature=${signature}`;
  }

  /**
   * Generate a unique nonce for request signing.
   *
   * Uses Node.js built-in crypto.randomUUID() (available in Node.js 14.17+).
   *
   * @returns UUID v4 string
   */
  private generateNonce(): string {
    return crypto.randomUUID();
  }

  /**
   * Parse URL into host, path, and query components.
   *
   * @param url - Full URL to parse
   * @returns Object containing host, path, and optional query
   *
   * @example
   * ```typescript
   * const { host, path, query } = TPV1Auth.parseUrl("https://api.example.com/v1/wallets?limit=10");
   * // host = "api.example.com"
   * // path = "/v1/wallets"
   * // query = "limit=10"
   * ```
   */
  static parseUrl(url: string): { host: string; path: string; query?: string } {
    const parsed = new URL(url);
    return {
      host: parsed.host,
      path: parsed.pathname || "/",
      query: parsed.search ? parsed.search.substring(1) : undefined,
    };
  }
}

/**
 * Calculate the TPV1 signed header value.
 *
 * This is a standalone function for use cases where you don't want to maintain
 * a TPV1Auth instance.
 *
 * @param apiKey - The API key
 * @param apiSecret - The API secret as bytes
 * @param nonce - Unique nonce for this request
 * @param timestamp - Request timestamp in milliseconds
 * @param method - HTTP method
 * @param host - Request host
 * @param path - Request path
 * @param query - Query string (optional)
 * @param contentType - Content-Type header (optional)
 * @param body - Request body (optional)
 * @returns The Authorization header value
 */
export function calculateSignedHeader(
  apiKey: string,
  apiSecret: Uint8Array,
  nonce: string,
  timestamp: number,
  method: string,
  host: string,
  path: string,
  query?: string,
  contentType?: string,
  body?: string
): string {
  // Build message from non-empty parts
  const parts: string[] = [
    "TPV1",
    apiKey,
    nonce,
    timestamp.toString(),
    method.toUpperCase(),
    host,
    path,
  ];

  if (query) {
    parts.push(query);
  }
  if (contentType) {
    parts.push(contentType);
  }
  if (body) {
    parts.push(body);
  }

  const message = parts.join(" ");
  const signature = calculateBase64Hmac(apiSecret, message);

  return `TPV1-HMAC-SHA256 ApiKey=${apiKey} Nonce=${nonce} Timestamp=${timestamp} Signature=${signature}`;
}
