/**
 * Config mapper functions for converting OpenAPI DTOs to domain models.
 */

import type { NFTMintingConfig, TenantConfig } from '../models/config';
import { safeBool, safeFloat, safeString } from './base';

/**
 * Maps an NFTMinting DTO to an NFTMintingConfig domain model.
 */
export function nftMintingConfigFromDto(dto: unknown): NFTMintingConfig | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    enabled: safeBool(d.enabled),
    publicBaseURL: safeString(d.publicBaseURL ?? d.public_base_url),
  };
}

/**
 * Maps a TenantConfig DTO to a TenantConfig domain model.
 *
 * This mapper handles the conversion from the OpenAPI-generated
 * TgvalidatordTenantConfig to the clean domain model.
 */
export function tenantConfigFromDto(dto: unknown): TenantConfig | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    tenantId: safeString(d.tenantId ?? d.tenant_id),
    superAdminMinimumSignatures: safeString(
      d.superAdminMinimumSignatures ?? d.super_admin_minimum_signatures
    ),
    baseCurrency: safeString(d.baseCurrency ?? d.base_currency),
    mfaMandatory: safeBool(d.isMFAMandatory ?? d.is_mfa_mandatory ?? d.mfaMandatory),
    excludeContainer: safeBool(d.excludeContainer ?? d.exclude_container),
    feeLimitFactor: safeFloat(d.feeLimitFactor ?? d.fee_limit_factor),
    protectEngineVersion: safeString(d.protectEngineVersion ?? d.protect_engine_version),
    restrictSourcesForWhitelistedAddresses: safeBool(
      d.restrictSourcesForWhitelistedAddresses ?? d.restrict_sources_for_whitelisted_addresses
    ),
    nftMinting: nftMintingConfigFromDto(d.nftMinting ?? d.nft_minting),
    protectEngineCold: safeBool(d.isProtectEngineCold ?? d.is_protect_engine_cold ?? d.protectEngineCold),
    coldProtectEngineOffline: safeBool(
      d.isColdProtectEngineOffline ?? d.is_cold_protect_engine_offline ?? d.coldProtectEngineOffline
    ),
    physicalAirGapEnabled: safeBool(
      d.isPhysicalAirGapEnabled ?? d.is_physical_air_gap_enabled ?? d.physicalAirGapEnabled
    ),
  };
}
