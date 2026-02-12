/**
 * Configuration for integration tests.
 *
 * Delegates to the shared testutil module for configuration loading.
 * Credentials are loaded from test.properties (git-ignored) or environment variables.
 *
 * Environment variables:
 * - PROTECT_INTEGRATION_TEST: Set to "true" to enable integration tests
 * - PROTECT_API_HOST: API host URL
 * - PROTECT_API_KEY_1: API key (identity 1)
 * - PROTECT_API_SECRET_1: API secret (identity 1)
 */

import {
  getHost,
  getIdentity,
  getIdentityCount,
  getSuperAdminKeys,
  getMinValidSignatures,
  isEnabled,
} from '../testutil/config';

/**
 * SuperAdmin public keys for governance signature verification.
 * Loaded from test.properties or environment variables.
 */
export const DEFAULT_SUPER_ADMIN_KEYS_PEM: string[] = getSuperAdminKeys();

/**
 * Minimum number of valid SuperAdmin signatures required for verification.
 */
export const DEFAULT_MIN_VALID_SIGNATURES = getMinValidSignatures();

/**
 * Test EC private key for request approval signing.
 * Loaded from identity 1's private key in test.properties.
 */
export const TEAM1_PRIVATE_KEY_PEM: string = (() => {
  try {
    if (getIdentityCount() > 0) {
      return getIdentity(1).privateKey;
    }
  } catch {
    // Config not loaded
  }
  return '';
})();

/**
 * Configuration interface for integration tests.
 */
export interface IntegrationConfig {
  host: string;
  apiKey: string;
  apiSecret: string;
}

/**
 * Returns configuration from the first identity, preferring ENV vars over properties file.
 */
export function getConfig(): IntegrationConfig {
  let host = '';
  let apiKey = '';
  let apiSecret = '';

  try {
    host = getHost();
    if (getIdentityCount() > 0) {
      const identity = getIdentity(1);
      apiKey = identity.apiKey;
      apiSecret = identity.apiSecret;
    }
  } catch {
    // Config not loaded
  }

  return { host, apiKey, apiSecret };
}

/**
 * Returns true if integration tests should run.
 */
export function isIntegrationEnabled(): boolean {
  return isEnabled();
}
