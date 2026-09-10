/**
 * Bearer-token middleware for the OpenAPI client.
 *
 * Adds an "Authorization: Bearer <token>" header to every request. The token is
 * obtained per request from a provider, so one client can serve rotating or
 * per-caller tokens (e.g. a provider reading async-local storage). Unlike the
 * TPV1 middleware it neither reads nor signs the body. Mirrors the Go SDK's
 * per-request bearer provider.
 */

import type {
  Middleware,
  RequestContext,
  FetchParams,
} from "../internal/openapi/runtime";
import { ConfigurationError } from "../errors";

/** Provides the bearer token for a request (resolved per request). */
export type BearerTokenProvider = () => string | Promise<string>;

/**
 * Creates bearer-token authentication middleware.
 *
 * @param tokenProvider - Resolves the bearer token for each request
 * @returns Middleware that sets the Authorization header
 */
export function createBearerMiddleware(
  tokenProvider: BearerTokenProvider
): Middleware {
  return {
    async pre(context: RequestContext): Promise<FetchParams | void> {
      const { url, init } = context;

      // Copy existing headers, tolerating Headers / array / object forms.
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

      // A refresh that silently yields nothing would send `Bearer undefined` and surface
      // as an opaque 401 rather than the real cause.
      const token = await tokenProvider();
      if (!token) {
        throw new ConfigurationError('bearer token provider returned an empty token');
      }
      headers.Authorization = `Bearer ${token}`;

      return {
        url,
        init: {
          ...init,
          headers,
        },
      };
    },
  };
}
