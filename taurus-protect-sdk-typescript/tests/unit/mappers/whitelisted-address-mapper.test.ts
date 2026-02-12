/**
 * Unit tests for whitelisted address mapping logic.
 *
 * Tests parseWhitelistedAddressFromJson from the whitelist hash helper module,
 * and createEmptyWhitelistedAddress from the model module.
 */

import { parseWhitelistedAddressFromJson } from '../../../src/helpers/whitelist-hash-helper';
import { createEmptyWhitelistedAddress } from '../../../src/models/whitelisted-address';

describe('parseWhitelistedAddressFromJson', () => {
  it('should parse a complete whitelisted address payload', () => {
    const payload = JSON.stringify({
      currency: 'ETH',
      network: 'mainnet',
      address: '0x1234567890abcdef1234567890abcdef12345678',
      memo: 'test-memo',
      label: 'My Wallet',
      customerId: 'cust-123',
      contractType: 'ERC20',
      addressType: 'individual',
      tnParticipantID: 'tn-part-1',
      exchangeAccountId: '42',
      linkedInternalAddresses: [
        { id: 1, label: 'Internal Addr 1' },
        { id: 2, label: 'Internal Addr 2' },
      ],
      linkedWallets: [
        { id: 10, name: 'Wallet A', path: 'ETH/wallet-a' },
      ],
    });

    const result = parseWhitelistedAddressFromJson(payload);

    expect(result.blockchain).toBe('ETH');
    expect(result.network).toBe('mainnet');
    expect(result.address).toBe('0x1234567890abcdef1234567890abcdef12345678');
    expect(result.memo).toBe('test-memo');
    expect(result.label).toBe('My Wallet');
    expect(result.customerId).toBe('cust-123');
    expect(result.contractType).toBe('ERC20');
    expect(result.addressType).toBe('individual');
    expect(result.tnParticipantId).toBe('tn-part-1');
    expect(result.exchangeAccountId).toBe(42);
    expect(result.linkedInternalAddresses).toHaveLength(2);
    expect(result.linkedInternalAddresses[0].id).toBe(1);
    expect(result.linkedInternalAddresses[0].label).toBe('Internal Addr 1');
    expect(result.linkedWallets).toHaveLength(1);
    expect(result.linkedWallets[0].id).toBe(10);
    expect(result.linkedWallets[0].label).toBe('Wallet A');
    expect(result.linkedWallets[0].path).toBe('ETH/wallet-a');
  });

  it('should use currency field as blockchain', () => {
    const payload = JSON.stringify({
      currency: 'BTC',
      network: 'testnet',
      address: 'tb1q...',
    });

    const result = parseWhitelistedAddressFromJson(payload);
    expect(result.blockchain).toBe('BTC');
  });

  it('should set id to empty string (not in signed payload)', () => {
    const payload = JSON.stringify({
      currency: 'ETH',
      network: 'mainnet',
      address: '0xabc',
    });

    const result = parseWhitelistedAddressFromJson(payload);
    expect(result.id).toBe('');
  });

  it('should handle missing optional fields', () => {
    const payload = JSON.stringify({
      currency: 'ETH',
      network: 'mainnet',
      address: '0x1234',
    });

    const result = parseWhitelistedAddressFromJson(payload);
    expect(result.memo).toBeUndefined();
    expect(result.label).toBeUndefined();
    expect(result.customerId).toBeUndefined();
    expect(result.contractType).toBeUndefined();
    expect(result.addressType).toBeUndefined();
    expect(result.tnParticipantId).toBeUndefined();
    expect(result.exchangeAccountId).toBeUndefined();
    expect(result.linkedInternalAddresses).toEqual([]);
    expect(result.linkedWallets).toEqual([]);
    expect(result.createdAt).toBeUndefined();
    expect(result.attributes).toEqual({});
  });

  it('should handle empty linked arrays', () => {
    const payload = JSON.stringify({
      currency: 'XTZ',
      network: 'mainnet',
      address: 'tz1...',
      linkedInternalAddresses: [],
      linkedWallets: [],
    });

    const result = parseWhitelistedAddressFromJson(payload);
    expect(result.linkedInternalAddresses).toEqual([]);
    expect(result.linkedWallets).toEqual([]);
  });

  it('should handle non-numeric exchangeAccountId gracefully', () => {
    const payload = JSON.stringify({
      currency: 'ETH',
      network: 'mainnet',
      address: '0x1234',
      exchangeAccountId: 'not-a-number',
    });

    const result = parseWhitelistedAddressFromJson(payload);
    expect(result.exchangeAccountId).toBeUndefined();
  });

  it('should handle linked addresses without labels', () => {
    const payload = JSON.stringify({
      currency: 'ETH',
      network: 'mainnet',
      address: '0x1234',
      linkedInternalAddresses: [
        { id: 1 },
        { id: 2 },
      ],
    });

    const result = parseWhitelistedAddressFromJson(payload);
    expect(result.linkedInternalAddresses).toHaveLength(2);
    expect(result.linkedInternalAddresses[0].label).toBeUndefined();
    expect(result.linkedInternalAddresses[1].label).toBeUndefined();
  });

  it('should handle linked wallets without paths', () => {
    const payload = JSON.stringify({
      currency: 'ETH',
      network: 'mainnet',
      address: '0x1234',
      linkedWallets: [
        { id: 10, name: 'Wallet' },
      ],
    });

    const result = parseWhitelistedAddressFromJson(payload);
    expect(result.linkedWallets[0].path).toBeUndefined();
    expect(result.linkedWallets[0].label).toBe('Wallet');
  });

  it('should default missing currency/network/address to empty string', () => {
    const payload = JSON.stringify({});

    const result = parseWhitelistedAddressFromJson(payload);
    expect(result.blockchain).toBe('');
    expect(result.network).toBe('');
    expect(result.address).toBe('');
  });

  it('should throw on empty payload', () => {
    expect(() => parseWhitelistedAddressFromJson('')).toThrow('JSON payload cannot be empty');
  });

  it('should throw on invalid JSON', () => {
    expect(() => parseWhitelistedAddressFromJson('not-json')).toThrow('Failed to parse whitelist payload');
  });
});

describe('createEmptyWhitelistedAddress', () => {
  it('should create an address with all required fields as defaults', () => {
    const empty = createEmptyWhitelistedAddress();

    expect(empty.id).toBe('');
    expect(empty.address).toBe('');
    expect(empty.blockchain).toBe('');
    expect(empty.network).toBe('');
    expect(empty.label).toBeUndefined();
    expect(empty.memo).toBeUndefined();
    expect(empty.customerId).toBeUndefined();
    expect(empty.contractType).toBeUndefined();
    expect(empty.addressType).toBeUndefined();
    expect(empty.tnParticipantId).toBeUndefined();
    expect(empty.exchangeAccountId).toBeUndefined();
    expect(empty.linkedInternalAddresses).toEqual([]);
    expect(empty.linkedWallets).toEqual([]);
    expect(empty.createdAt).toBeUndefined();
    expect(empty.attributes).toEqual({});
  });
});
