/**
 * Unit tests for config mapper functions.
 */

import {
  tenantConfigFromDto,
  nftMintingConfigFromDto,
} from '../../../src/mappers/config';

describe('nftMintingConfigFromDto', () => {
  it('should map all fields', () => {
    const dto = {
      enabled: true,
      publicBaseURL: 'https://nft.example.com',
    };

    const result = nftMintingConfigFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.enabled).toBe(true);
    expect(result!.publicBaseURL).toBe('https://nft.example.com');
  });

  it('should handle snake_case fields', () => {
    const dto = {
      enabled: false,
      public_base_url: 'https://other.example.com',
    };

    const result = nftMintingConfigFromDto(dto);

    expect(result!.publicBaseURL).toBe('https://other.example.com');
  });

  it('should return undefined for null input', () => {
    expect(nftMintingConfigFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(nftMintingConfigFromDto(undefined)).toBeUndefined();
  });
});

describe('tenantConfigFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      tenantId: 't-1',
      superAdminMinimumSignatures: '2',
      baseCurrency: 'USD',
      isMFAMandatory: true,
      excludeContainer: false,
      feeLimitFactor: 1.5,
      protectEngineVersion: 'v2.0',
      restrictSourcesForWhitelistedAddresses: true,
      nftMinting: { enabled: true, publicBaseURL: 'https://nft.example.com' },
      isProtectEngineCold: false,
      isColdProtectEngineOffline: false,
      isPhysicalAirGapEnabled: true,
    };

    const result = tenantConfigFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.tenantId).toBe('t-1');
    expect(result!.superAdminMinimumSignatures).toBe('2');
    expect(result!.baseCurrency).toBe('USD');
    expect(result!.mfaMandatory).toBe(true);
    expect(result!.excludeContainer).toBe(false);
    expect(result!.feeLimitFactor).toBe(1.5);
    expect(result!.protectEngineVersion).toBe('v2.0');
    expect(result!.restrictSourcesForWhitelistedAddresses).toBe(true);
    expect(result!.nftMinting).toBeDefined();
    expect(result!.nftMinting!.enabled).toBe(true);
    expect(result!.protectEngineCold).toBe(false);
    expect(result!.coldProtectEngineOffline).toBe(false);
    expect(result!.physicalAirGapEnabled).toBe(true);
  });

  it('should handle snake_case field names', () => {
    const dto = {
      tenant_id: 't-2',
      super_admin_minimum_signatures: '3',
      base_currency: 'EUR',
      is_mfa_mandatory: false,
      exclude_container: true,
      fee_limit_factor: 2.0,
      protect_engine_version: 'v3.0',
      restrict_sources_for_whitelisted_addresses: false,
      nft_minting: { enabled: false },
      is_protect_engine_cold: true,
      is_cold_protect_engine_offline: true,
      is_physical_air_gap_enabled: false,
    };

    const result = tenantConfigFromDto(dto);

    expect(result!.tenantId).toBe('t-2');
    expect(result!.superAdminMinimumSignatures).toBe('3');
    expect(result!.baseCurrency).toBe('EUR');
    expect(result!.mfaMandatory).toBe(false);
    expect(result!.protectEngineCold).toBe(true);
    expect(result!.coldProtectEngineOffline).toBe(true);
    expect(result!.physicalAirGapEnabled).toBe(false);
  });

  it('should return undefined for null input', () => {
    expect(tenantConfigFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(tenantConfigFromDto(undefined)).toBeUndefined();
  });

  it('should handle empty object', () => {
    const result = tenantConfigFromDto({});
    expect(result).toBeDefined();
    expect(result!.tenantId).toBeUndefined();
  });
});
