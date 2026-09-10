/**
 * Authentication mechanism for a {@link ProtectClient}: either static TPV1-HMAC
 * (api key + secret) or a Bearer token (static, or resolved per request for
 * rotating/per-caller tokens).
 *
 * Build one with {@link Credentials.apiKey}, {@link Credentials.bearerToken}, or
 * {@link Credentials.bearerTokenProvider} and pass it as `credentials` to
 * {@link ProtectClient.create}. SuperAdmin keys are required regardless of the
 * mechanism.
 */

import type { Middleware } from "./internal/openapi/runtime";
import { ConfigurationError } from "./errors";
import {
  createBearerMiddleware,
  createTPV1Middleware,
  type BearerTokenProvider,
} from "./transport";

export class Credentials {
  // Builds the OpenAPI auth middleware for the chosen mechanism.
  private constructor(private readonly build: () => Middleware) {}

  /**
   * TPV1-HMAC credentials (a static shared-service identity).
   *
   * @param apiKey - The API key.
   * @param apiSecret - The API secret in hexadecimal format.
   */
  static apiKey(apiKey: string, apiSecret: string): Credentials {
    if (!apiKey) {
      throw new ConfigurationError("apiKey is required");
    }
    if (!apiSecret) {
      throw new ConfigurationError("apiSecret is required");
    }
    if (!/^[0-9a-fA-F]+$/.test(apiSecret)) {
      throw new ConfigurationError("apiSecret must be a valid hexadecimal string");
    }
    return new Credentials(() => createTPV1Middleware(apiKey, apiSecret));
  }

  /**
   * A single static Bearer token (the "Authorization: Bearer" header).
   */
  static bearerToken(token: string): Credentials {
    if (!token) {
      throw new ConfigurationError("bearerToken is required");
    }
    return new Credentials(() => createBearerMiddleware(() => token));
  }

  /**
   * A Bearer token resolved per request from `provider()` (for
   * rotating/per-caller tokens).
   */
  static bearerTokenProvider(provider: BearerTokenProvider): Credentials {
    if (!provider) {
      throw new ConfigurationError("bearerTokenProvider is required");
    }
    return new Credentials(() => createBearerMiddleware(provider));
  }

  /** @internal Builds the authenticating OpenAPI middleware. */
  toMiddleware(): Middleware {
    return this.build();
  }
}
