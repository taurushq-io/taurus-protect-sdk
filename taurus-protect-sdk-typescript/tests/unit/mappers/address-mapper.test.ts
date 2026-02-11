/**
 * Unit tests for address mapper functions.
 */

import {
  addressFromDto,
  addressesFromDto,
  addressAttributeFromDto,
} from '../../../src/mappers/address';

describe('addressFromDto', () => {
  it('should map all fields correctly', () => {
    const dto = {
      id: '456',
      walletId: '123',
      address: '0xabc123def456',
      alternateAddress: '0xalt',
      label: 'My Address',
      comment: 'Test address',
      currency: 'ETH',
      customerId: 'cust-1',
      externalAddressId: 'ext-456',
      addressPath: "m/44'/60'/0'/0/0",
      addressIndex: '0',
      nonce: '42',
      status: 'active',
      signature: 'sig-data-base64',
      disabled: false,
      canUseAllFunds: true,
      creationDate: new Date('2024-01-01'),
      updateDate: new Date('2024-06-01'),
      attributes: [
        { id: 'a1', key: 'tag', value: 'vip' },
      ],
      linkedWhitelistedAddressIds: ['wl-1', 'wl-2'],
    };

    const address = addressFromDto(dto as any);

    expect(address).toBeDefined();
    expect(address!.id).toBe('456');
    expect(address!.walletId).toBe('123');
    expect(address!.address).toBe('0xabc123def456');
    expect(address!.alternateAddress).toBe('0xalt');
    expect(address!.label).toBe('My Address');
    expect(address!.comment).toBe('Test address');
    expect(address!.currency).toBe('ETH');
    expect(address!.customerId).toBe('cust-1');
    expect(address!.signature).toBe('sig-data-base64');
    expect(address!.disabled).toBe(false);
    expect(address!.canUseAllFunds).toBe(true);
    expect(address!.attributes).toHaveLength(1);
    expect(address!.linkedWhitelistedAddressIds).toEqual(['wl-1', 'wl-2']);
  });

  it('should return undefined for null input', () => {
    expect(addressFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(addressFromDto(undefined)).toBeUndefined();
  });

  it('should handle missing optional fields', () => {
    const dto = {
      id: '1',
      walletId: '10',
      address: '0xminimal',
      currency: 'ETH',
    };

    const address = addressFromDto(dto as any);

    expect(address).toBeDefined();
    expect(address!.id).toBe('1');
    expect(address!.address).toBe('0xminimal');
    expect(address!.label).toBeUndefined();
    expect(address!.comment).toBeUndefined();
    expect(address!.signature).toBeUndefined();
  });
});

describe('addressesFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { id: '1', walletId: '10', address: '0xaaa', currency: 'ETH' },
      { id: '2', walletId: '10', address: '0xbbb', currency: 'ETH' },
    ];

    const addresses = addressesFromDto(dtos as any);

    expect(addresses).toHaveLength(2);
    expect(addresses[0].address).toBe('0xaaa');
    expect(addresses[1].address).toBe('0xbbb');
  });

  it('should return empty array for null input', () => {
    expect(addressesFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(addressesFromDto(undefined)).toEqual([]);
  });

  it('should return empty array for empty input', () => {
    expect(addressesFromDto([])).toEqual([]);
  });
});

describe('addressAttributeFromDto', () => {
  it('should map attribute fields', () => {
    const dto = {
      id: 'a1',
      key: 'tag',
      value: 'important',
    };

    const attr = addressAttributeFromDto(dto as any);

    expect(attr).toBeDefined();
    expect(attr.id).toBe('a1');
    expect(attr.key).toBe('tag');
    expect(attr.value).toBe('important');
  });
});
