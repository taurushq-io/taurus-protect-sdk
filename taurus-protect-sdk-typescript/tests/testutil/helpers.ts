/**
 * Shared test helper functions for integration and E2E tests.
 */

import { ProtectClient } from '../../src/client';
import {
  getHost,
  getIdentity,
  getIdentityCount,
  getSuperAdminKeys,
  getMinValidSignatures,
  isEnabled,
  hasPrivateKey,
} from './config';

/**
 * Skips the current test if tests are not enabled.
 * Call this at the beginning of a beforeAll() block.
 *
 * @throws Error with message 'SKIP_TEST' if not enabled
 */
export function skipIfNotEnabled(): void {
  if (!isEnabled()) {
    console.log(
      'Skipping test. Set PROTECT_INTEGRATION_TEST=true or configure test.properties.'
    );
    throw new Error('SKIP_TEST');
  }
}

/**
 * Skips the current test if fewer than `required` identities are configured.
 *
 * @param required - minimum number of identities needed
 * @throws Error with message 'SKIP_TEST' if insufficient identities
 */
export function skipIfInsufficientIdentities(required: number): void {
  const count = getIdentityCount();
  if (count < required) {
    console.log(
      `Skipping: need ${required} identities but only ${count} configured.`
    );
    throw new Error('SKIP_TEST');
  }
}

/**
 * Creates a ProtectClient for the identity at the given 1-based index.
 * Configures SuperAdmin keys for signature verification.
 *
 * @param identityIndex - 1-based identity index
 * @returns configured ProtectClient
 */
export function getTestClient(identityIndex: number = 1): ProtectClient {
  const identity = getIdentity(identityIndex);
  const host = getHost();

  if (!host) {
    throw new Error(
      'API host not configured. Set PROTECT_API_HOST or add host= to test.properties.'
    );
  }
  if (!identity.apiKey || !identity.apiSecret) {
    throw new Error(
      `Identity ${identityIndex} has no API credentials. ` +
      `Set PROTECT_API_KEY_${identityIndex}/PROTECT_API_SECRET_${identityIndex} ` +
      `or configure in test.properties.`
    );
  }

  const superAdminKeys = getSuperAdminKeys();
  const minSigs = getMinValidSignatures();

  return ProtectClient.create({
    host,
    apiKey: identity.apiKey,
    apiSecret: identity.apiSecret,
    ...(superAdminKeys.length > 0
      ? { superAdminKeysPem: superAdminKeys, minValidSignatures: minSigs }
      : { superAdminKeysPem: [], minValidSignatures: 0 }),
  });
}

/**
 * Returns the PEM-encoded private key for the identity at the given 1-based index.
 * Returns undefined if the identity has no private key.
 *
 * @param identityIndex - 1-based identity index
 * @returns PEM string or undefined
 */
export function getPrivateKey(identityIndex: number): string | undefined {
  const identity = getIdentity(identityIndex);
  if (!hasPrivateKey(identity)) {
    return undefined;
  }
  return identity.privateKey;
}

/**
 * Conditionally runs a describe block only if tests are enabled.
 */
export const describeIntegration = isEnabled() ? describe : describe.skip;

/**
 * Conditionally runs a test only if tests are enabled.
 */
export const itIntegration = isEnabled() ? it : it.skip;
