/**
 * Shared test utilities for integration and E2E tests.
 */

export { loadProperties } from './properties';

export {
  getHost,
  getIdentity,
  getIdentityCount,
  getAllIdentities,
  getSuperAdminKeys,
  getMinValidSignatures,
  isEnabled,
  hasApiCredentials,
  hasPrivateKey,
  hasPublicKey,
} from './config';

export type { Identity } from './config';

export {
  skipIfNotEnabled,
  skipIfInsufficientIdentities,
  getTestClient,
  getPrivateKey,
  describeIntegration,
  itIntegration,
} from './helpers';
