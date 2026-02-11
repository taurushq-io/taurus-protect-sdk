/**
 * Configuration for integration and E2E tests with multi-identity support.
 *
 * Loads configuration from `test.properties` (searched in multiple paths),
 * with environment variable overrides for CI/CD pipelines.
 *
 * Each identity may have API credentials (for making API calls), a private key
 * (for signing operations), and/or a public key (for SuperAdmin verification).
 *
 * Environment variables:
 * - PROTECT_INTEGRATION_TEST: Set to "true" to enable tests
 * - PROTECT_API_HOST: API host URL
 * - PROTECT_API_KEY_N: API key for identity N (1-based)
 * - PROTECT_API_SECRET_N: API secret for identity N
 * - PROTECT_PRIVATE_KEY_N: Private key (PEM) for identity N
 * - PROTECT_PUBLIC_KEY_N: Public key (PEM) for identity N
 */

import * as path from 'path';
import { loadProperties } from './properties';

// ── Identity ────────────────────────────────────────────────────────────

/**
 * Represents a user identity in the test system.
 */
export interface Identity {
  /** 1-based index */
  readonly index: number;
  /** Human-readable name */
  readonly name: string;
  /** API key (empty string if not configured) */
  readonly apiKey: string;
  /** API secret (empty string if not configured) */
  readonly apiSecret: string;
  /** PEM-encoded EC private key (empty string if not configured) */
  readonly privateKey: string;
  /** PEM-encoded EC public key (empty string if not configured) */
  readonly publicKey: string;
}

function hasApiCredentials(identity: Identity): boolean {
  return identity.apiKey.length > 0 && identity.apiSecret.length > 0;
}

function hasPrivateKey(identity: Identity): boolean {
  return identity.privateKey.length > 0;
}

function hasPublicKey(identity: Identity): boolean {
  return identity.publicKey.length > 0;
}

// ── Properties loading ──────────────────────────────────────────────────

const searchPaths = [
  'tests/testutil/test.properties',
  'test.properties',
  '../tests/testutil/test.properties',
];

let _properties: Map<string, string> | null = null;
let _propertiesLoaded = false;

function getProperties(): Map<string, string> | null {
  if (_propertiesLoaded) return _properties;
  _propertiesLoaded = true;

  for (const searchPath of searchPaths) {
    const fullPath = path.resolve(searchPath);
    const props = loadProperties(fullPath);
    if (props !== null) {
      _properties = props;
      return _properties;
    }
  }
  return null;
}

/**
 * Resolves a setting value: env var takes priority, then properties file, then empty string.
 */
function resolve(envVar: string | null, propertiesKey: string): string {
  if (envVar) {
    const env = process.env[envVar];
    if (env && env.length > 0) {
      return env;
    }
  }
  const props = getProperties();
  if (props) {
    const prop = props.get(propertiesKey);
    if (prop && prop.length > 0) {
      return prop;
    }
  }
  return '';
}

// ── Identity loading ────────────────────────────────────────────────────

let _identities: Identity[] | null = null;

function loadIdentities(): Identity[] {
  const list: Identity[] = [];

  for (let i = 1; i < 100; i++) {
    const name = resolve(null, `identity.${i}.name`);
    const apiKey = resolve(`PROTECT_API_KEY_${i}`, `identity.${i}.apiKey`);
    const apiSecret = resolve(`PROTECT_API_SECRET_${i}`, `identity.${i}.apiSecret`);
    const privateKey = resolve(`PROTECT_PRIVATE_KEY_${i}`, `identity.${i}.privateKey`);
    const publicKey = resolve(`PROTECT_PUBLIC_KEY_${i}`, `identity.${i}.publicKey`);

    // Skip gaps (e.g., identity.1 then identity.4)
    if (!name && !apiKey && !apiSecret && !privateKey && !publicKey) {
      continue;
    }

    list.push({
      index: i,
      name: name || `identity-${i}`,
      apiKey,
      apiSecret,
      privateKey,
      publicKey,
    });
  }

  return list;
}

function getIdentities(): Identity[] {
  if (_identities === null) {
    _identities = loadIdentities();
  }
  return _identities;
}

// ── Public API ──────────────────────────────────────────────────────────

/**
 * Returns the API host, preferring ENV var over properties file.
 */
export function getHost(): string {
  return resolve('PROTECT_API_HOST', 'host');
}

/**
 * Returns the identity at the given 1-based index.
 */
export function getIdentity(index: number): Identity {
  const identities = getIdentities();
  const identity = identities.find((id) => id.index === index);
  if (!identity) {
    throw new RangeError(
      `Identity ${index} not found. Available: [${identities.map((id) => id.index).join(', ')}]`
    );
  }
  return identity;
}

/**
 * Returns how many identities are configured.
 */
export function getIdentityCount(): number {
  return getIdentities().length;
}

/**
 * Returns all configured identities.
 */
export function getAllIdentities(): Identity[] {
  return [...getIdentities()];
}

/**
 * Returns PEM-encoded SuperAdmin public keys from identities that have a public key.
 */
export function getSuperAdminKeys(): string[] {
  return getIdentities()
    .filter(hasPublicKey)
    .map((id) => id.publicKey);
}

/**
 * Returns the minimum number of valid signatures required.
 */
export function getMinValidSignatures(): number {
  const props = getProperties();
  if (props) {
    const prop = props.get('minValidSignatures');
    if (prop) {
      const parsed = parseInt(prop, 10);
      if (!isNaN(parsed)) return parsed;
    }
  }
  return 2;
}

/**
 * Returns true if tests should run.
 * Enabled if PROTECT_INTEGRATION_TEST=true OR at least one identity has API credentials.
 */
export function isEnabled(): boolean {
  if (process.env.PROTECT_INTEGRATION_TEST === 'true') {
    return true;
  }
  return getIdentities().some(hasApiCredentials);
}

// Re-export utility functions for identity checks
export { hasApiCredentials, hasPrivateKey, hasPublicKey };
