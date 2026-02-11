/**
 * Unit tests for fee payer mapper functions.
 */

import {
  feePayerFromDto,
  feePayersFromDto,
  feePayerInfoFromDto,
  feePayerEthFromDto,
  feePayerEthLocalFromDto,
  feePayerEthRemoteFromDto,
} from '../../../src/mappers/fee-payer';

describe('feePayerEthLocalFromDto', () => {
  it('should map all fields', () => {
    const dto = {
      addressId: 'addr-1',
      forwarderAddressId: 'fwd-1',
      autoApprove: true,
      creatorAddressId: 'creator-1',
      forwarderKind: 'gsn',
      domainSeparator: '0xabc',
    };

    const result = feePayerEthLocalFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.addressId).toBe('addr-1');
    expect(result!.forwarderAddressId).toBe('fwd-1');
    expect(result!.autoApprove).toBe(true);
    expect(result!.creatorAddressId).toBe('creator-1');
    expect(result!.forwarderKind).toBe('gsn');
    expect(result!.domainSeparator).toBe('0xabc');
  });

  it('should handle snake_case fields', () => {
    const dto = {
      address_id: 'a1',
      forwarder_address_id: 'f1',
      auto_approve: false,
      creator_address_id: 'c1',
      forwarder_kind: 'eip2771',
      domain_separator: '0xdef',
    };

    const result = feePayerEthLocalFromDto(dto);

    expect(result!.addressId).toBe('a1');
    expect(result!.forwarderAddressId).toBe('f1');
    expect(result!.autoApprove).toBe(false);
  });

  it('should return undefined for null input', () => {
    expect(feePayerEthLocalFromDto(null)).toBeUndefined();
  });
});

describe('feePayerEthRemoteFromDto', () => {
  it('should map all fields', () => {
    const dto = {
      url: 'https://relay.example.com',
      username: 'admin',
      fromAddressId: 'from-1',
      forwarderAddress: '0xFWD',
      forwarderAddressId: 'fwd-1',
      creatorAddress: '0xCRT',
      creatorAddressId: 'crt-1',
      forwarderKind: 'gsn',
      domainSeparator: '0xsep',
    };

    const result = feePayerEthRemoteFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.url).toBe('https://relay.example.com');
    expect(result!.username).toBe('admin');
    expect(result!.fromAddressId).toBe('from-1');
    expect(result!.forwarderAddress).toBe('0xFWD');
  });

  it('should return undefined for null input', () => {
    expect(feePayerEthRemoteFromDto(null)).toBeUndefined();
  });
});

describe('feePayerEthFromDto', () => {
  it('should map eth fee payer with local and remote', () => {
    const dto = {
      kind: 'local',
      local: { addressId: 'a1', autoApprove: true },
      remote: { url: 'https://relay.example.com' },
      remoteEncrypted: 'encrypted-data',
    };

    const result = feePayerEthFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.kind).toBe('local');
    expect(result!.local).toBeDefined();
    expect(result!.remote).toBeDefined();
    expect(result!.remoteEncrypted).toBe('encrypted-data');
  });

  it('should return undefined for null input', () => {
    expect(feePayerEthFromDto(null)).toBeUndefined();
  });
});

describe('feePayerInfoFromDto', () => {
  it('should map fee payer info', () => {
    const dto = {
      blockchain: 'ETH',
      eth: { kind: 'local' },
    };

    const result = feePayerInfoFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.blockchain).toBe('ETH');
    expect(result!.eth).toBeDefined();
  });

  it('should return undefined for null input', () => {
    expect(feePayerInfoFromDto(null)).toBeUndefined();
  });
});

describe('feePayerFromDto', () => {
  it('should map full fee payer envelope', () => {
    const dto = {
      id: 'fp-1',
      tenantId: 't-1',
      blockchain: 'ETH',
      network: 'mainnet',
      name: 'Main Fee Payer',
      creationDate: new Date('2024-01-01'),
      feePayer: { blockchain: 'ETH', eth: { kind: 'local' } },
    };

    const result = feePayerFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.id).toBe('fp-1');
    expect(result!.tenantId).toBe('t-1');
    expect(result!.blockchain).toBe('ETH');
    expect(result!.network).toBe('mainnet');
    expect(result!.name).toBe('Main Fee Payer');
    expect(result!.creationDate).toBeInstanceOf(Date);
    expect(result!.feePayerInfo).toBeDefined();
  });

  it('should handle snake_case fields', () => {
    const dto = {
      id: 'fp-2',
      tenant_id: 't-2',
      creation_date: '2024-03-01T00:00:00Z',
      fee_payer: { blockchain: 'ETH' },
    };

    const result = feePayerFromDto(dto);

    expect(result!.tenantId).toBe('t-2');
    expect(result!.feePayerInfo).toBeDefined();
  });

  it('should return undefined for null input', () => {
    expect(feePayerFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(feePayerFromDto(undefined)).toBeUndefined();
  });
});

describe('feePayersFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { id: 'fp-1', name: 'FP1' },
      { id: 'fp-2', name: 'FP2' },
    ];

    const result = feePayersFromDto(dtos);

    expect(result).toHaveLength(2);
  });

  it('should return empty array for null input', () => {
    expect(feePayersFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(feePayersFromDto(undefined)).toEqual([]);
  });
});
