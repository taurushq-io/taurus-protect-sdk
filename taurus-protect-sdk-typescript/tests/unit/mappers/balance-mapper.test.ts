/**
 * Unit tests for balance mapper functions.
 */

import {
  assetBalanceFromDto,
  assetBalancesFromDto,
  nftCollectionBalanceFromDto,
  nftCollectionBalancesFromDto,
} from '../../../src/mappers/balance';

describe('assetBalanceFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      currencyId: 'BTC',
      currency: 'BTC',
      blockchain: 'BTC',
      network: 'mainnet',
      contractAddress: '0x123',
      tokenId: 'tok-1',
      balance: '1000',
      fiatValue: '50000',
      fiatCurrency: 'USD',
    };

    const result = assetBalanceFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.currencyId).toBeDefined();
    expect(result!.currency).toBeDefined();
    expect(result!.blockchain).toBeDefined();
    expect(result!.balance).toBe('1000');
    expect(result!.fiatValue).toBeDefined();
  });

  it('should handle snake_case field names', () => {
    const dto = {
      currency_id: 'ETH',
      contract_address: '0x456',
      token_id: 'tok-2',
      total_confirmed: '500',
      fiat_value: '25000',
      fiat_currency: 'EUR',
    };

    const result = assetBalanceFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.currencyId).toBeDefined();
    expect(result!.contractAddress).toBeDefined();
    expect(result!.balance).toBe('500');
  });

  it('should return undefined for null input', () => {
    expect(assetBalanceFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(assetBalanceFromDto(undefined)).toBeUndefined();
  });

  it('should handle empty object', () => {
    const result = assetBalanceFromDto({});
    expect(result).toBeDefined();
  });
});

describe('assetBalancesFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { currency: 'BTC', balance: '100' },
      { currency: 'ETH', balance: '200' },
    ];
    const result = assetBalancesFromDto(dtos);
    expect(result).toHaveLength(2);
  });

  it('should return empty array for null input', () => {
    expect(assetBalancesFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(assetBalancesFromDto(undefined)).toEqual([]);
  });
});

describe('nftCollectionBalanceFromDto', () => {
  it('should map all fields', () => {
    const dto = {
      name: 'Bored Apes',
      symbol: 'BAYC',
      blockchain: 'ETH',
      network: 'mainnet',
      contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D',
      count: 10,
      logoUrl: 'https://example.com/logo.png',
    };

    const result = nftCollectionBalanceFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.name).toBe('Bored Apes');
    expect(result!.count).toBe(10);
    expect(result!.contractAddress).toBeDefined();
  });

  it('should handle snake_case fallbacks', () => {
    const dto = {
      name: 'CryptoPunks',
      contract_address: '0x1234',
      logo_url: 'https://example.com/punks.png',
      balance: 5,
    };

    const result = nftCollectionBalanceFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.contractAddress).toBe('0x1234');
    expect(result!.count).toBe(5);
  });

  it('should return undefined for null input', () => {
    expect(nftCollectionBalanceFromDto(null)).toBeUndefined();
  });
});

describe('nftCollectionBalancesFromDto', () => {
  it('should return empty array for null input', () => {
    expect(nftCollectionBalancesFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(nftCollectionBalancesFromDto(undefined)).toEqual([]);
  });
});
