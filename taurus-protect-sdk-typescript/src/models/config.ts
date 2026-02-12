/**
 * Tenant configuration models for Taurus-PROTECT SDK.
 *
 * Contains configuration settings for the tenant, including security requirements,
 * feature flags, and system parameters.
 */

/**
 * NFT minting configuration for the tenant.
 */
export interface NFTMintingConfig {
  /**
   * Whether NFT minting is enabled for the tenant.
   */
  enabled?: boolean;

  /**
   * The public base URL for NFT metadata.
   */
  publicBaseURL?: string;
}

/**
 * Represents the tenant configuration in the Taurus-PROTECT system.
 *
 * Contains various configuration settings for the tenant, including
 * security requirements, feature flags, and system parameters.
 *
 * @example
 * ```typescript
 * const config = await client.config.getTenantConfig();
 *
 * // Check if MFA is mandatory
 * if (config.mfaMandatory) {
 *   console.log('MFA is required for this tenant');
 * }
 *
 * // Get the base currency
 * console.log(`Base currency: ${config.baseCurrency}`);
 * ```
 */
export interface TenantConfig {
  /**
   * The unique identifier for the tenant.
   */
  tenantId?: string;

  /**
   * The minimum number of SuperAdmin signatures required for governance operations.
   */
  superAdminMinimumSignatures?: string;

  /**
   * The base currency for the tenant (e.g., 'USD', 'EUR', 'CHF').
   */
  baseCurrency?: string;

  /**
   * Whether multi-factor authentication is mandatory for the tenant.
   */
  mfaMandatory?: boolean;

  /**
   * Whether to exclude container information from responses.
   */
  excludeContainer?: boolean;

  /**
   * The fee limit factor applied to transaction fees.
   */
  feeLimitFactor?: number;

  /**
   * The version of the Protect Engine.
   */
  protectEngineVersion?: string;

  /**
   * Whether sources are restricted for whitelisted addresses.
   */
  restrictSourcesForWhitelistedAddresses?: boolean;

  /**
   * NFT minting configuration for the tenant.
   */
  nftMinting?: NFTMintingConfig;

  /**
   * Whether the Protect Engine is in cold mode.
   */
  protectEngineCold?: boolean;

  /**
   * Whether the cold Protect Engine is offline.
   */
  coldProtectEngineOffline?: boolean;

  /**
   * Whether physical air gap is enabled for the tenant.
   */
  physicalAirGapEnabled?: boolean;
}
