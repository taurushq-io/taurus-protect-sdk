/**
 * TPV1-HMAC-SHA256 middleware for OpenAPI client.
 *
 * This middleware automatically signs all outgoing requests using
 * the TPV1 authentication protocol.
 */

import type {
  Middleware,
  RequestContext,
  FetchParams,
} from "../internal/openapi/runtime";
import { TPV1Auth } from "../crypto/tpv1";

/**
 * Creates TPV1 authentication middleware.
 *
 * This middleware intercepts all outgoing requests and adds the
 * TPV1-HMAC-SHA256 Authorization header.
 *
 * @param apiKey - The API key for authentication
 * @param apiSecretHex - The API secret in hexadecimal format
 * @returns Middleware that signs requests with TPV1
 *
 * @example
 * ```typescript
 * import { Configuration } from './internal/openapi/runtime';
 * import { createTPV1Middleware } from './transport/tpv1-middleware';
 *
 * const middleware = createTPV1Middleware('my-api-key', 'hex-encoded-secret');
 * const config = new Configuration({
 *   basePath: 'https://api.protect.taurushq.com',
 *   middleware: [middleware],
 * });
 * ```
 */
export function createTPV1Middleware(
  apiKey: string,
  apiSecretHex: string
): Middleware {
  const tpv1 = new TPV1Auth(apiKey, apiSecretHex);

  return {
    async pre(context: RequestContext): Promise<FetchParams | void> {
      const { url, init } = context;

      // Parse URL to extract components
      const parsedUrl = new URL(url);
      const method = (init.method as string) ?? "GET";
      const host = parsedUrl.host;
      const path = parsedUrl.pathname;
      const query = parsedUrl.search ? parsedUrl.search.substring(1) : undefined;

      // Get headers - handle both Headers object and plain object
      let headers: Record<string, string> = {};
      if (init.headers) {
        if (init.headers instanceof Headers) {
          init.headers.forEach((value, key) => {
            headers[key] = value;
          });
        } else if (Array.isArray(init.headers)) {
          for (const header of init.headers) {
            const [key, value] = header;
            if (key !== undefined && value !== undefined) {
              headers[key] = value;
            }
          }
        } else {
          headers = { ...(init.headers as Record<string, string>) };
        }
      }

      // Get content type (case-insensitive lookup)
      const contentType =
        headers["Content-Type"] ?? headers["content-type"] ?? undefined;

      // Get body as string
      let body: string | undefined;
      if (init.body !== null && init.body !== undefined) {
        if (typeof init.body === "string") {
          body = init.body;
        } else if (init.body instanceof ArrayBuffer) {
          body = new TextDecoder().decode(init.body);
        } else if (ArrayBuffer.isView(init.body)) {
          body = new TextDecoder().decode(init.body);
        }
        // Note: FormData and Blob bodies are not supported for TPV1 signing
        // as they cannot be converted to a deterministic string representation
      }

      // Calculate signed header using TPV1Auth
      const authHeader = tpv1.signRequest(
        method,
        host,
        path,
        query,
        contentType,
        body
      );

      // Add Authorization header
      const newHeaders: Record<string, string> = {
        ...headers,
        Authorization: authHeader,
      };

      return {
        url,
        init: {
          ...init,
          headers: newHeaders,
        },
      };
    },
  };
}
