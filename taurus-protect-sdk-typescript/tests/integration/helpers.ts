/**
 * Integration test helpers.
 *
 * Provides utilities for running integration tests with Jest.
 */

import { ProtectClient } from "../../src/client";
import {
  getConfig,
  isIntegrationEnabled,
  DEFAULT_SUPER_ADMIN_KEYS_PEM,
  DEFAULT_MIN_VALID_SIGNATURES,
} from "./config";

/**
 * Skips the current test if integration tests are not enabled.
 *
 * Call this at the beginning of integration tests to conditionally skip.
 *
 * @example
 * ```typescript
 * beforeAll(() => {
 *   skipIfNotIntegration();
 * });
 *
 * it('should list wallets', async () => {
 *   const client = getTestClient();
 *   // ...
 * });
 * ```
 */
export function skipIfNotIntegration(): void {
  if (!isIntegrationEnabled()) {
    console.log(
      "Skipping integration test. Set PROTECT_INTEGRATION_TEST=true or configure defaults in config.ts"
    );
    // Throwing will cause beforeAll to fail, effectively skipping tests
    throw new Error("SKIP_INTEGRATION");
  }
}

/**
 * Conditionally runs a describe block only if integration tests are enabled.
 *
 * @example
 * ```typescript
 * describeIntegration('WalletsApi Integration', () => {
 *   it('should list wallets', async () => {
 *     const client = getTestClient();
 *     // ...
 *   });
 * });
 * ```
 */
export const describeIntegration = isIntegrationEnabled()
  ? describe
  : describe.skip;

/**
 * Conditionally runs a test only if integration tests are enabled.
 *
 * @example
 * ```typescript
 * itIntegration('should list wallets', async () => {
 *   const client = getTestClient();
 *   // ...
 * });
 * ```
 */
export const itIntegration = isIntegrationEnabled() ? it : it.skip;

/**
 * Creates and returns a ProtectClient instance configured for integration tests.
 *
 * The client is created using configuration from environment variables or defaults.
 * SuperAdmin keys are always provided since integrity verification is mandatory.
 *
 * @returns A configured ProtectClient instance
 *
 * @example
 * ```typescript
 * const client = getTestClient();
 * try {
 *   // Use high-level services, not raw OpenAPI APIs
 *   const result = await client.wallets.list({ limit: 10 });
 *   for (const wallet of result.items) {
 *     console.log(wallet.name);
 *   }
 * } finally {
 *   client.close();
 * }
 * ```
 */
export function getTestClient(): ProtectClient {
  const config = getConfig();

  if (!config.host) {
    throw new Error(
      "API host not configured. Set PROTECT_API_HOST or DEFAULT_API_HOST in config.ts"
    );
  }
  if (!config.apiKey) {
    throw new Error(
      "API key not configured. Set PROTECT_API_KEY or DEFAULT_API_KEY in config.ts"
    );
  }
  if (!config.apiSecret) {
    throw new Error(
      "API secret not configured. Set PROTECT_API_SECRET or DEFAULT_API_SECRET in config.ts"
    );
  }

  return ProtectClient.create({
    host: config.host,
    apiKey: config.apiKey,
    apiSecret: config.apiSecret,
    superAdminKeysPem: DEFAULT_SUPER_ADMIN_KEYS_PEM,
    minValidSignatures: DEFAULT_MIN_VALID_SIGNATURES,
  });
}

/**
 * Creates and returns a ProtectClient instance configured with SuperAdmin keys
 * for governance signature verification.
 *
 * Use this client for tests that verify governance rules, whitelisted addresses,
 * or whitelisted assets signatures.
 *
 * @returns A configured ProtectClient instance with signature verification enabled
 *
 * @example
 * ```typescript
 * const client = getTestClientWithVerification();
 * try {
 *   const rules = await client.governanceRules.getRules();
 *   // GetDecodedRulesContainer will verify signatures
 *   const decoded = await client.governanceRules.getDecodedRulesContainer(rules);
 * } finally {
 *   client.close();
 * }
 * ```
 */
export function getTestClientWithVerification(): ProtectClient {
  const config = getConfig();

  if (!config.host) {
    throw new Error(
      "API host not configured. Set PROTECT_API_HOST or DEFAULT_API_HOST in config.ts"
    );
  }
  if (!config.apiKey) {
    throw new Error(
      "API key not configured. Set PROTECT_API_KEY or DEFAULT_API_KEY in config.ts"
    );
  }
  if (!config.apiSecret) {
    throw new Error(
      "API secret not configured. Set PROTECT_API_SECRET or DEFAULT_API_SECRET in config.ts"
    );
  }

  return ProtectClient.create({
    host: config.host,
    apiKey: config.apiKey,
    apiSecret: config.apiSecret,
    superAdminKeysPem: DEFAULT_SUPER_ADMIN_KEYS_PEM,
    minValidSignatures: DEFAULT_MIN_VALID_SIGNATURES,
  });
}

// Re-export config functions for convenience
export { getConfig, isIntegrationEnabled } from "./config";
export {
  DEFAULT_SUPER_ADMIN_KEYS_PEM,
  DEFAULT_MIN_VALID_SIGNATURES,
  TEAM1_PRIVATE_KEY_PEM,
} from "./config";
