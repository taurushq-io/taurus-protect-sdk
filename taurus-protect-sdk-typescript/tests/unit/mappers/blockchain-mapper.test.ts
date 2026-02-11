/**
 * Unit tests for blockchain mapper functions.
 */

import { blockchainFromDto, blockchainsFromDto } from '../../../src/mappers/blockchain';

describe('blockchainFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      symbol: 'ETH',
      name: 'Ethereum',
      network: 'mainnet',
      chainId: '1',
      confirmations: '12',
      blockHeight: '18000000',
      blackholeAddress: '0x0000000000000000000000000000000000000000',
      isLayer2Chain: false,
      layer1Network: undefined,
      baseCurrency: { id: 'c1', symbol: 'ETH', name: 'Ethereum' },
      dotInfo: null,
      ethInfo: { chainId: '1' },
      xtzInfo: null,
    };

    const result = blockchainFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.symbol).toBe('ETH');
    expect(result!.name).toBe('Ethereum');
    expect(result!.network).toBe('mainnet');
    expect(result!.chainId).toBe('1');
    expect(result!.confirmations).toBe('12');
    expect(result!.blockHeight).toBe('18000000');
    expect(result!.blackholeAddress).toBe('0x0000000000000000000000000000000000000000');
    expect(result!.isLayer2Chain).toBe(false);
    expect(result!.baseCurrency).toBeDefined();
    expect(result!.ethInfo).toBeDefined();
    expect(result!.ethInfo!.chainId).toBe('1');
  });

  it('should handle snake_case field names', () => {
    const dto = {
      symbol: 'DOT',
      chain_id: '42',
      block_height: '999',
      blackhole_address: '0x0',
      is_layer2_chain: true,
      layer1_network: 'polkadot',
      base_currency: { id: 'c2', symbol: 'DOT' },
      dot_info: { ss_58_format: 42 },
    };

    const result = blockchainFromDto(dto);

    expect(result!.chainId).toBe('42');
    expect(result!.blockHeight).toBe('999');
    expect(result!.isLayer2Chain).toBe(true);
    expect(result!.layer1Network).toBe('polkadot');
    expect(result!.dotInfo).toBeDefined();
    expect(result!.dotInfo!.ss58Format).toBe(42);
  });

  it('should handle XTZ blockchain info', () => {
    const dto = {
      symbol: 'XTZ',
      xtzInfo: { protocolHash: 'PtNairobiyssHuh87hEhfVBGCVrK3WnS8Z2FT4ymB5tAa4r1nQf' },
    };

    const result = blockchainFromDto(dto);

    expect(result!.xtzInfo).toBeDefined();
    expect(result!.xtzInfo!.protocolHash).toBe('PtNairobiyssHuh87hEhfVBGCVrK3WnS8Z2FT4ymB5tAa4r1nQf');
  });

  it('should return undefined for null input', () => {
    expect(blockchainFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(blockchainFromDto(undefined)).toBeUndefined();
  });

  it('should handle empty object', () => {
    const result = blockchainFromDto({});
    expect(result).toBeDefined();
    expect(result!.symbol).toBeUndefined();
  });
});

describe('blockchainsFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { symbol: 'ETH', name: 'Ethereum' },
      { symbol: 'BTC', name: 'Bitcoin' },
    ];

    const result = blockchainsFromDto(dtos);

    expect(result).toHaveLength(2);
    expect(result[0].symbol).toBe('ETH');
    expect(result[1].symbol).toBe('BTC');
  });

  it('should return empty array for null input', () => {
    expect(blockchainsFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(blockchainsFromDto(undefined)).toEqual([]);
  });

  it('should return empty array for empty array', () => {
    expect(blockchainsFromDto([])).toEqual([]);
  });
});
