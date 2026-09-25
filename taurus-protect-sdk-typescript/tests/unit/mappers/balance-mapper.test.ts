/**
 * Unit tests for balance mapper functions.
 *
 * The rows are shaped as the generated client delivers them: an asset balance is
 * `{asset: {currency, currencyInfo, nft}, balance: {totalConfirmed, ...}}` and an NFT
 * collection balance is `{currencyInfo, balance}`. The flat `currency` / `balance` keys
 * earlier versions of these tests used do not exist on the wire.
 */

import {
  assetBalanceFromDto,
  assetBalancesFromDto,
  nftCollectionBalanceFromDto,
  nftCollectionBalancesFromDto,
} from '../../../src/mappers/balance';

describe('assetBalanceFromDto', () => {
  it('should map the currency from asset.currencyInfo and the amount from balance', () => {
    const result = assetBalanceFromDto({
      asset: {
        currency: 'USDC',
        currencyInfo: {
          id: 'c-usdc',
          symbol: 'USDC',
          blockchain: 'ETH',
          network: 'mainnet',
          contractAddress: '0xa0b8',
          tokenID: 'tok-info',
        },
        nft: { tokenid: 'tok-1' },
      },
      balance: { totalConfirmed: '1000', availableConfirmed: '900' },
    });
    expect(result).toEqual({
      currencyId: 'c-usdc',
      currency: 'USDC',
      blockchain: 'ETH',
      network: 'mainnet',
      contractAddress: '0xa0b8',
      tokenId: 'tok-1',
      balance: '1000',
      fiatValue: undefined,
      fiatCurrency: undefined,
    });
  });

  it('should fall back to asset.currency and currencyInfo.tokenID', () => {
    const result = assetBalanceFromDto({
      asset: { currency: 'ALGO', currencyInfo: { tokenID: 'tok-info' } },
    });
    expect(result?.currency).toBe('ALGO');
    expect(result?.tokenId).toBe('tok-info');
    expect(result?.balance).toBeUndefined();
  });

  it('should return undefined for null input', () => {
    expect(assetBalanceFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(assetBalanceFromDto(undefined)).toBeUndefined();
  });
});

describe('assetBalancesFromDto', () => {
  it('should map array of rows', () => {
    const result = assetBalancesFromDto([
      { asset: { currency: 'BTC' }, balance: { totalConfirmed: '100' } },
      { asset: { currency: 'ETH' }, balance: { totalConfirmed: '200' } },
    ]);
    expect(result.map((b) => [b.currency, b.balance])).toEqual([
      ['BTC', '100'],
      ['ETH', '200'],
    ]);
  });

  it('should return empty array for null input', () => {
    expect(assetBalancesFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(assetBalancesFromDto(undefined)).toEqual([]);
  });
});

describe('nftCollectionBalanceFromDto', () => {
  it('should map the collection from currencyInfo and the count from the balance', () => {
    const result = nftCollectionBalanceFromDto({
      currencyInfo: {
        name: 'Bored Apes',
        symbol: 'BAYC',
        blockchain: 'ETH',
        network: 'mainnet',
        contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D',
        logo: 'data:image/png;base64,AA==',
      },
      balance: { totalConfirmed: '10' },
    });
    expect(result).toEqual({
      name: 'Bored Apes',
      symbol: 'BAYC',
      blockchain: 'ETH',
      network: 'mainnet',
      contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D',
      count: 10,
      logoUrl: 'data:image/png;base64,AA==',
    });
  });

  it('should leave the count undefined without a balance', () => {
    expect(nftCollectionBalanceFromDto({ currencyInfo: { name: 'X' } })?.count).toBeUndefined();
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
