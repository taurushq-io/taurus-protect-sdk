/**
 * Unit tests for governance rules mapper functions.
 */

import {
  rulesContainerFromBase64,
  userSignaturesFromBase64,
  governanceRulesFromDto,
  ruleUserSignatureFromDto,
  rulesTrailFromDto,
} from '../../../src/mappers/governance-rules';
import { IntegrityError } from '../../../src/errors';

describe('rulesContainerFromBase64', () => {
  it('should return empty container for empty string', () => {
    const result = rulesContainerFromBase64('');
    expect(result).toBeDefined();
    expect(result.users).toEqual([]);
    expect(result.groups).toEqual([]);
  });

  it('should decode JSON-encoded rules container', () => {
    const container = {
      users: [
        { id: 'u-1', name: 'Admin', publicKeyPem: '-----BEGIN PUBLIC KEY-----', roles: ['ADMIN'] },
      ],
      groups: [
        { id: 'g-1', name: 'Operators', userIds: ['u-1'] },
      ],
      minimumDistinctUserSignatures: 2,
      minimumDistinctGroupSignatures: 1,
      addressWhitelistingRules: [
        {
          currency: 'BTC',
          network: 'mainnet',
          parallelThresholds: [
            { thresholds: [{ groupId: 'g-1', minimumSignatures: 2 }] },
          ],
        },
      ],
      contractAddressWhitelistingRules: [
        {
          blockchain: 'ETH',
          network: 'mainnet',
          parallelThresholds: [],
        },
      ],
      enforcedRulesHash: 'hash123',
      timestamp: 1700000000,
    };
    const base64 = Buffer.from(JSON.stringify(container)).toString('base64');

    const result = rulesContainerFromBase64(base64);

    expect(result.users).toHaveLength(1);
    expect(result.users[0].id).toBe('u-1');
    expect(result.users[0].publicKeyPem).toBe('-----BEGIN PUBLIC KEY-----');
    expect(result.groups).toHaveLength(1);
    expect(result.groups[0].id).toBe('g-1');
    expect(result.minimumDistinctUserSignatures).toBe(2);
    expect(result.addressWhitelistingRules).toHaveLength(1);
    expect(result.addressWhitelistingRules[0].currency).toBe('BTC');
    expect(result.addressWhitelistingRules[0].parallelThresholds).toHaveLength(1);
    expect(result.contractAddressWhitelistingRules).toHaveLength(1);
    expect(result.enforcedRulesHash).toBe('hash123');
    expect(result.timestamp).toBe(1700000000);
  });

  it('should handle flat sequential thresholds format', () => {
    const container = {
      users: [],
      groups: [],
      addressWhitelistingRules: [
        {
          currency: 'ETH',
          network: 'mainnet',
          parallelThresholds: [
            { groupId: 'g-1', minimumSignatures: 3 },
          ],
        },
      ],
    };
    const base64 = Buffer.from(JSON.stringify(container)).toString('base64');

    const result = rulesContainerFromBase64(base64);
    expect(result.addressWhitelistingRules).toHaveLength(1);
    const thresholds = result.addressWhitelistingRules[0].parallelThresholds;
    expect(thresholds).toHaveLength(1);
    expect(thresholds[0].thresholds).toHaveLength(1);
    expect(thresholds[0].thresholds[0].groupId).toBe('g-1');
    expect(thresholds[0].thresholds[0].minimumSignatures).toBe(3);
  });

  it('should handle snake_case property names', () => {
    const container = {
      users: [
        { id: 'u-1', public_key_pem: 'PEM', roles: [] },
      ],
      groups: [
        { id: 'g-1', userIds: ['u-1'] },
      ],
      minimum_distinct_user_signatures: 1,
      address_whitelisting_rules: [],
    };
    const base64 = Buffer.from(JSON.stringify(container)).toString('base64');

    const result = rulesContainerFromBase64(base64);
    expect(result.users[0].publicKeyPem).toBe('PEM');
    expect(result.groups[0].userIds).toEqual(['u-1']);
    expect(result.minimumDistinctUserSignatures).toBe(1);
  });

  it('should throw IntegrityError for invalid data', () => {
    // Base64 of binary data that is neither valid protobuf nor JSON
    const invalidBase64 = Buffer.from(new Uint8Array([0xFF, 0xFE, 0xFD, 0xFC])).toString('base64');
    expect(() => rulesContainerFromBase64(invalidBase64)).toThrow(IntegrityError);
  });
});

describe('userSignaturesFromBase64', () => {
  it('should return empty array for empty string', () => {
    expect(userSignaturesFromBase64('')).toEqual([]);
  });

  it('should decode JSON array signatures', () => {
    const sigs = [
      { userId: 'u-1', signature: 'sig-base64-1' },
      { userId: 'u-2', signature: 'sig-base64-2' },
    ];
    const base64 = Buffer.from(JSON.stringify(sigs)).toString('base64');

    const result = userSignaturesFromBase64(base64);
    expect(result).toHaveLength(2);
    expect(result[0].userId).toBe('u-1');
    expect(result[0].signature).toBe('sig-base64-1');
  });

  it('should decode JSON object with signatures field', () => {
    const data = {
      signatures: [
        { userId: 'u-1', signature: 'sig1' },
      ],
    };
    const base64 = Buffer.from(JSON.stringify(data)).toString('base64');

    const result = userSignaturesFromBase64(base64);
    expect(result).toHaveLength(1);
  });

  it('should handle snake_case property names', () => {
    const sigs = [
      { user_id: 'u-1', signature: 'sig1' },
    ];
    const base64 = Buffer.from(JSON.stringify(sigs)).toString('base64');

    const result = userSignaturesFromBase64(base64);
    expect(result).toHaveLength(1);
    expect(result[0].userId).toBe('u-1');
  });

  it('should return empty array for invalid base64', () => {
    expect(userSignaturesFromBase64('not-valid!!!')).toEqual([]);
  });
});

describe('governanceRulesFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      rulesContainer: 'base64data',
      rulesSignatures: [
        { userId: 'u-1', signature: 'sig1' },
      ],
      locked: true,
      creationDate: new Date('2024-01-01'),
      updateDate: new Date('2024-06-01'),
      trails: [
        { userId: 'u-1', action: 'LOCKED', date: new Date('2024-01-15') },
      ],
    };

    const result = governanceRulesFromDto(dto as any);
    expect(result).toBeDefined();
    expect(result!.rulesContainer).toBe('base64data');
    expect(result!.locked).toBe(true);
    expect(result!.rulesSignatures).toHaveLength(1);
    expect(result!.trails).toHaveLength(1);
  });

  it('should return undefined for null input', () => {
    expect(governanceRulesFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(governanceRulesFromDto(undefined)).toBeUndefined();
  });
});

describe('ruleUserSignatureFromDto', () => {
  it('should map fields', () => {
    const dto = { userId: 'u-1', signature: 'sig-data' };
    const result = ruleUserSignatureFromDto(dto as any);
    expect(result).toBeDefined();
    expect(result!.userId).toBe('u-1');
    expect(result!.signature).toBe('sig-data');
  });

  it('should return undefined for null input', () => {
    expect(ruleUserSignatureFromDto(null)).toBeUndefined();
  });
});

describe('rulesTrailFromDto', () => {
  it('should map fields', () => {
    const dto = {
      userId: 'u-1',
      action: 'LOCKED',
      date: new Date('2024-01-15'),
    };
    const result = rulesTrailFromDto(dto as any);
    expect(result).toBeDefined();
    expect(result!.userId).toBe('u-1');
    expect(result!.action).toBe('LOCKED');
    expect(result!.timestamp).toBeDefined();
  });

  it('should return undefined for null input', () => {
    expect(rulesTrailFromDto(null)).toBeUndefined();
  });
});
