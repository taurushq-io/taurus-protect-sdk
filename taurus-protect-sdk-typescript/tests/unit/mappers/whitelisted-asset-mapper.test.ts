/**
 * Unit tests for whitelisted asset mapping logic.
 *
 * Tests parseWhitelistedAssetFromJson and createEmptyWhitelistedAsset
 * from the whitelisted-asset model module.
 */

import {
  parseWhitelistedAssetFromJson,
  createEmptyWhitelistedAsset,
} from '../../../src/models/whitelisted-asset';

describe('parseWhitelistedAssetFromJson', () => {
  it('should parse a complete whitelisted asset payload', () => {
    const payload = JSON.stringify({
      blockchain: 'ETH',
      network: 'mainnet',
      contractAddress: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48',
      name: 'USD Coin',
      symbol: 'USDC',
      decimals: 6,
    });

    const result = parseWhitelistedAssetFromJson(payload);

    expect(result.blockchain).toBe('ETH');
    expect(result.network).toBe('mainnet');
    expect(result.contractAddress).toBe('0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48');
    expect(result.name).toBe('USD Coin');
    expect(result.symbol).toBe('USDC');
    expect(result.decimals).toBe(6);
  });

  it('should set id to 0 (not in signed payload)', () => {
    const payload = JSON.stringify({
      blockchain: 'ETH',
      network: 'mainnet',
      contractAddress: '0x1234',
    });

    const result = parseWhitelistedAssetFromJson(payload);
    expect(result.id).toBe(0);
  });

  it('should handle missing optional fields', () => {
    const payload = JSON.stringify({
      blockchain: 'MATIC',
      network: 'mainnet',
      contractAddress: '0xabc',
    });

    const result = parseWhitelistedAssetFromJson(payload);
    expect(result.name).toBeUndefined();
    expect(result.symbol).toBeUndefined();
    expect(result.decimals).toBeUndefined();
    expect(result.createdAt).toBeUndefined();
  });

  it('should default missing required fields to empty string', () => {
    const payload = JSON.stringify({});

    const result = parseWhitelistedAssetFromJson(payload);
    expect(result.blockchain).toBe('');
    expect(result.network).toBe('');
    expect(result.contractAddress).toBe('');
  });

  it('should parse NFT contract (decimals 0)', () => {
    const payload = JSON.stringify({
      blockchain: 'ETH',
      network: 'mainnet',
      contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D',
      name: 'Bored Ape Yacht Club',
      symbol: 'BAYC',
      decimals: 0,
    });

    const result = parseWhitelistedAssetFromJson(payload);
    expect(result.name).toBe('Bored Ape Yacht Club');
    expect(result.symbol).toBe('BAYC');
    expect(result.decimals).toBe(0);
  });

  it('should parse Tezos FA token', () => {
    const payload = JSON.stringify({
      blockchain: 'XTZ',
      network: 'mainnet',
      contractAddress: 'KT1PWx2mnDueood7fEmfbBDKx1D9BAnnXitn',
      name: 'tzBTC',
      symbol: 'tzBTC',
      decimals: 8,
    });

    const result = parseWhitelistedAssetFromJson(payload);
    expect(result.blockchain).toBe('XTZ');
    expect(result.contractAddress).toBe('KT1PWx2mnDueood7fEmfbBDKx1D9BAnnXitn');
    expect(result.decimals).toBe(8);
  });

  it('should throw on empty payload', () => {
    expect(() => parseWhitelistedAssetFromJson('')).toThrow('JSON payload cannot be empty');
  });

  it('should throw on invalid JSON', () => {
    expect(() => parseWhitelistedAssetFromJson('not-json')).toThrow('Failed to parse whitelist payload');
  });

  it('should handle payload with extra fields gracefully', () => {
    const payload = JSON.stringify({
      blockchain: 'ETH',
      network: 'mainnet',
      contractAddress: '0x1234',
      name: 'Token',
      symbol: 'TKN',
      decimals: 18,
      isNFT: false,
      kindType: 'erc20',
      tokenId: 'tok-1',
    });

    const result = parseWhitelistedAssetFromJson(payload);
    expect(result.blockchain).toBe('ETH');
    expect(result.name).toBe('Token');
    expect(result.decimals).toBe(18);
  });
});

describe('createEmptyWhitelistedAsset', () => {
  it('should create an asset with all required fields as defaults', () => {
    const empty = createEmptyWhitelistedAsset();

    expect(empty.id).toBe(0);
    expect(empty.contractAddress).toBe('');
    expect(empty.blockchain).toBe('');
    expect(empty.network).toBe('');
    expect(empty.name).toBeUndefined();
    expect(empty.symbol).toBeUndefined();
    expect(empty.decimals).toBeUndefined();
    expect(empty.createdAt).toBeUndefined();
  });
});
