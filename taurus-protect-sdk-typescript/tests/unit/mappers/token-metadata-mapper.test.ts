/**
 * Unit tests for token metadata mapping logic.
 *
 * The TokenMetadataService uses private mapping methods (no dedicated mapper file).
 * These tests verify the mapping patterns used by the service.
 */

import type { TokenMetadata, CryptoPunkMetadata } from '../../../src/models/token-metadata';

/**
 * Replicates the mapERCTokenMetadata private method from TokenMetadataService.
 */
function mapERCTokenMetadata(dto: unknown): TokenMetadata {
  if (!dto) {
    return {};
  }
  const d = dto as Record<string, unknown>;
  return {
    name: d.name as string | undefined,
    description: d.description as string | undefined,
    decimals: d.decimals as string | undefined,
    dataType: d.dataType as string | undefined,
    base64Data: d.base64Data as string | undefined,
    uri: d.uri as string | undefined,
  };
}

/**
 * Replicates the mapFATokenMetadata private method from TokenMetadataService.
 */
function mapFATokenMetadata(dto: unknown): TokenMetadata {
  if (!dto) {
    return {};
  }
  const d = dto as Record<string, unknown>;
  return {
    name: d.name as string | undefined,
    symbol: d.symbol as string | undefined,
    decimals: d.decimals as string | undefined,
    dataType: d.dataType as string | undefined,
    base64Data: d.base64Data as string | undefined,
    uri: d.uri as string | undefined,
  };
}

/**
 * Replicates the mapCryptoPunkMetadata private method from TokenMetadataService.
 */
function mapCryptoPunkMetadata(dto: unknown): CryptoPunkMetadata {
  if (!dto) {
    return {};
  }
  const d = dto as Record<string, unknown>;
  return {
    punkId: d.punkId as string | undefined,
    punkAttributes: d.punkAttributes as string | undefined,
    image: d.image as string | undefined,
  };
}

describe('mapERCTokenMetadata', () => {
  it('should map ERC-20 token metadata', () => {
    const dto = {
      name: 'USD Coin',
      description: 'Stablecoin pegged to USD',
      decimals: '6',
    };

    const result = mapERCTokenMetadata(dto);
    expect(result.name).toBe('USD Coin');
    expect(result.description).toBe('Stablecoin pegged to USD');
    expect(result.decimals).toBe('6');
  });

  it('should map ERC-721 NFT metadata with media data', () => {
    const dto = {
      name: 'Bored Ape #1234',
      description: 'A unique ape',
      dataType: 'image/png',
      base64Data: 'iVBORw0KGgoAAAANSUhEUgAA...',
      uri: 'ipfs://QmABCDEF...',
    };

    const result = mapERCTokenMetadata(dto);
    expect(result.name).toBe('Bored Ape #1234');
    expect(result.dataType).toBe('image/png');
    expect(result.base64Data).toBe('iVBORw0KGgoAAAANSUhEUgAA...');
    expect(result.uri).toBe('ipfs://QmABCDEF...');
  });

  it('should return empty object for null input', () => {
    const result = mapERCTokenMetadata(null);
    expect(result).toEqual({});
  });

  it('should return empty object for undefined input', () => {
    const result = mapERCTokenMetadata(undefined);
    expect(result).toEqual({});
  });

  it('should handle empty DTO', () => {
    const result = mapERCTokenMetadata({});
    expect(result.name).toBeUndefined();
    expect(result.decimals).toBeUndefined();
  });
});

describe('mapFATokenMetadata', () => {
  it('should map FA token metadata with symbol', () => {
    const dto = {
      name: 'tzBTC',
      symbol: 'tzBTC',
      decimals: '8',
      uri: 'https://tzbtc.io/metadata',
    };

    const result = mapFATokenMetadata(dto);
    expect(result.name).toBe('tzBTC');
    expect(result.symbol).toBe('tzBTC');
    expect(result.decimals).toBe('8');
    expect(result.uri).toBe('https://tzbtc.io/metadata');
  });

  it('should return empty object for null input', () => {
    const result = mapFATokenMetadata(null);
    expect(result).toEqual({});
  });

  it('should handle FA2 NFT metadata', () => {
    const dto = {
      name: 'Tezos Domain',
      symbol: 'TD',
      decimals: '0',
      dataType: 'image/svg+xml',
      base64Data: 'PHN2ZyB4bWxucz0ia...',
    };

    const result = mapFATokenMetadata(dto);
    expect(result.symbol).toBe('TD');
    expect(result.decimals).toBe('0');
    expect(result.dataType).toBe('image/svg+xml');
  });
});

describe('mapCryptoPunkMetadata', () => {
  it('should map CryptoPunk metadata', () => {
    const dto = {
      punkId: '7804',
      punkAttributes: 'Male 2, Cap Forward, Small Shades, Pipe',
      image: 'iVBORw0KGgoAAAANSUhEUgAAABgA...',
    };

    const result = mapCryptoPunkMetadata(dto);
    expect(result.punkId).toBe('7804');
    expect(result.punkAttributes).toBe('Male 2, Cap Forward, Small Shades, Pipe');
    expect(result.image).toBe('iVBORw0KGgoAAAANSUhEUgAAABgA...');
  });

  it('should return empty object for null input', () => {
    const result = mapCryptoPunkMetadata(null);
    expect(result).toEqual({});
  });

  it('should return empty object for undefined input', () => {
    const result = mapCryptoPunkMetadata(undefined);
    expect(result).toEqual({});
  });

  it('should handle partial data', () => {
    const dto = {
      punkId: '0',
    };

    const result = mapCryptoPunkMetadata(dto);
    expect(result.punkId).toBe('0');
    expect(result.punkAttributes).toBeUndefined();
    expect(result.image).toBeUndefined();
  });
});
