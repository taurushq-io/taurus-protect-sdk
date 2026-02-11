/**
 * Unit tests for currency mapper functions.
 */

import { currencyFromDto, currenciesFromDto } from '../../../src/mappers/currency';

describe('currencyFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      id: 'BTC',
      symbol: 'BTC',
      name: 'Bitcoin',
      displayName: 'Bitcoin (BTC)',
      type: 'NATIVE',
      blockchain: 'BTC',
      network: 'mainnet',
      decimals: 8,
      coinTypeIndex: '0',
      contractAddress: null,
      tokenId: null,
      wlcaId: 1,
      logo: 'https://example.com/btc.png',
      isToken: false,
      isERC20: false,
      isFA12: false,
      isFA20: false,
      isNFT: false,
      isUTXOBased: true,
      isAccountBased: false,
      isFiat: false,
      hasStaking: false,
      enabled: true,
    };

    const result = currencyFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('BTC');
    expect(result!.symbol).toBe('BTC');
    expect(result!.name).toBe('Bitcoin');
    expect(result!.decimals).toBe(8);
    expect(result!.isUTXOBased).toBe(true);
    expect(result!.isToken).toBe(false);
    expect(result!.enabled).toBe(true);
  });

  it('should handle snake_case field names', () => {
    const dto = {
      id: 'ETH',
      symbol: 'ETH',
      name: 'Ethereum',
      display_name: 'Ethereum',
      is_token: false,
      is_erc20: false,
      is_utxo_based: false,
      is_account_based: true,
      coin_type_index: '60',
      contract_address: null,
      has_staking: true,
    };

    const result = currencyFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.displayName).toBeDefined();
    expect(result!.isAccountBased).toBe(true);
    expect(result!.hasStaking).toBe(true);
  });

  it('should return undefined for null input', () => {
    expect(currencyFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(currencyFromDto(undefined)).toBeUndefined();
  });

  it('should default enabled to true', () => {
    const dto = { id: 'BTC' };
    const result = currencyFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.enabled).toBe(true);
  });

  it('should handle token contract address fallback', () => {
    const dto = {
      id: 'USDT',
      tokenContractAddress: '0xdAC17F958D2ee523a2206206994597C13D831ec7',
    };
    const result = currencyFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.contractAddress).toBe('0xdAC17F958D2ee523a2206206994597C13D831ec7');
  });
});

describe('currenciesFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { id: 'BTC', symbol: 'BTC' },
      { id: 'ETH', symbol: 'ETH' },
    ];
    const result = currenciesFromDto(dtos);
    expect(result).toHaveLength(2);
  });

  it('should return empty array for null input', () => {
    expect(currenciesFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(currenciesFromDto(undefined)).toEqual([]);
  });
});
