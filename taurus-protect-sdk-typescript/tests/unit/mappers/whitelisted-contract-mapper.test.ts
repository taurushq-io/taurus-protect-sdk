/**
 * Unit tests for whitelisted contract mapper functions.
 *
 * This test file covers the same mapper as contract-whitelist-mapper.test.ts
 * but focuses on additional edge cases and the envelope-to-contract mapping.
 */

import {
  whitelistedContractFromDto,
  whitelistedContractsFromDto,
  whitelistedContractAttributeFromDto,
  whitelistedContractAttributesFromDto,
  whitelistedContractMetadataFromDto,
  signedWhitelistedContractEnvelopeFromDto,
  whitelistedContractFromEnvelope,
  whitelistedContractResultFromDto,
  whitelistedContractTrailFromDto,
  whitelistedContractTrailsFromDto,
  whitelistedContractApproversFromDto,
} from '../../../src/mappers/contract-whitelist';

describe('whitelistedContractFromDto', () => {
  it('should map a full envelope DTO to a WhitelistedContract', () => {
    const dto = {
      id: 'wc-1',
      tenantId: 'tenant-1',
      blockchain: 'ETH',
      network: 'mainnet',
      metadata: {
        hash: 'abc123',
        payloadAsString: JSON.stringify({
          contractAddress: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48',
          symbol: 'USDC',
          name: 'USD Coin',
          decimals: 6,
          kind: 'erc20',
        }),
      },
      status: 'ACTIVE',
      businessRuleEnabled: true,
      attributes: [
        { id: 'a1', key: 'category', value: 'stablecoin' },
      ],
    };

    const result = whitelistedContractFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('wc-1');
    expect(result!.blockchain).toBe('ETH');
    expect(result!.network).toBe('mainnet');
    expect(result!.contractAddress).toBe('0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48');
    expect(result!.symbol).toBe('USDC');
    expect(result!.name).toBe('USD Coin');
    expect(result!.decimals).toBe(6);
    expect(result!.kind).toBe('erc20');
    expect(result!.status).toBe('ACTIVE');
    expect(result!.businessRuleEnabled).toBe(true);
    expect(result!.attributes).toHaveLength(1);
  });

  it('should return undefined for null input', () => {
    expect(whitelistedContractFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(whitelistedContractFromDto(undefined)).toBeUndefined();
  });

  it('should handle DTO without metadata payload', () => {
    const dto = {
      id: 'wc-2',
      blockchain: 'ETH',
      network: 'mainnet',
    };

    const result = whitelistedContractFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('wc-2');
    expect(result!.contractAddress).toBeUndefined();
    expect(result!.symbol).toBeUndefined();
  });
});

describe('whitelistedContractsFromDto', () => {
  it('should map multiple DTOs', () => {
    const dtos = [
      { id: 'wc-1', metadata: { payloadAsString: '{"symbol":"USDC"}' } },
      { id: 'wc-2', metadata: { payloadAsString: '{"symbol":"WETH"}' } },
    ];

    const result = whitelistedContractsFromDto(dtos);
    expect(result).toHaveLength(2);
    expect(result[0].symbol).toBe('USDC');
    expect(result[1].symbol).toBe('WETH');
  });

  it('should return empty array for null', () => {
    expect(whitelistedContractsFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined', () => {
    expect(whitelistedContractsFromDto(undefined)).toEqual([]);
  });
});

describe('whitelistedContractAttributesFromDto', () => {
  it('should map array of attributes', () => {
    const dtos = [
      { id: 'a1', key: 'k1', value: 'v1' },
      { id: 'a2', key: 'k2', value: 'v2' },
    ];

    const result = whitelistedContractAttributesFromDto(dtos);
    expect(result).toHaveLength(2);
    expect(result[0].key).toBe('k1');
    expect(result[1].value).toBe('v2');
  });

  it('should return empty array for null', () => {
    expect(whitelistedContractAttributesFromDto(null)).toEqual([]);
  });

  it('should handle snake_case attribute fields', () => {
    const dto = {
      id: 'a1',
      key: 'logo',
      value: 'data:image/png;base64,...',
      content_type: 'image/png',
      sub_type: 'logo',
      is_file: true,
    };

    const result = whitelistedContractAttributeFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.contentType).toBe('image/png');
    expect(result!.subType).toBe('logo');
    expect(result!.isFile).toBe(true);
  });
});

describe('whitelistedContractMetadataFromDto - security', () => {
  it('should only map hash and payloadAsString (not raw payload)', () => {
    const dto = {
      hash: 'secure-hash',
      payloadAsString: '{"contractAddress":"0x1234"}',
      payload: { contractAddress: '0x9999-TAMPERED' },
    };

    const result = whitelistedContractMetadataFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.hash).toBe('secure-hash');
    expect(result!.payloadAsString).toBe('{"contractAddress":"0x1234"}');
    // payload field should not be mapped (security)
    expect((result as unknown as Record<string, unknown>)['payload']).toBeUndefined();
  });
});

describe('whitelistedContractTrailsFromDto', () => {
  it('should map array of trails', () => {
    const dtos = [
      { id: 't1', action: 'created', userId: 'u1', timestamp: '2024-01-01' },
      { id: 't2', action: 'approved', userId: 'u2', timestamp: '2024-01-02' },
    ];

    const result = whitelistedContractTrailsFromDto(dtos);
    expect(result).toHaveLength(2);
    expect(result[0].action).toBe('created');
    expect(result[1].action).toBe('approved');
  });

  it('should return empty array for null', () => {
    expect(whitelistedContractTrailsFromDto(null)).toEqual([]);
  });

  it('should handle snake_case trail fields', () => {
    const dto = {
      id: 't1',
      action: 'created',
      user_id: 'u1',
      comment: 'Initial',
      created_at: '2024-01-01',
    };

    const result = whitelistedContractTrailFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.userId).toBe('u1');
    expect(result!.timestamp).toBe('2024-01-01');
  });
});

describe('whitelistedContractApproversFromDto', () => {
  it('should map approvers with nested groups and users', () => {
    const dto = {
      required: 2,
      groups: [
        {
          id: 'g1',
          name: 'Admins',
          required: 1,
          users: [
            { userId: 'u1', userName: 'Alice', pending: false },
            { userId: 'u2', userName: 'Bob', pending: true },
          ],
        },
      ],
    };

    const result = whitelistedContractApproversFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.required).toBe(2);
    expect(result!.groups).toHaveLength(1);
    expect(result!.groups[0].name).toBe('Admins');
    expect(result!.groups[0].users).toHaveLength(2);
    expect(result!.groups[0].users[0].userName).toBe('Alice');
  });
});

describe('signedWhitelistedContractEnvelopeFromDto', () => {
  it('should handle snake_case field names', () => {
    const dto = {
      id: 'env-1',
      tenant_id: 'tenant-1',
      signed_contract_address: {
        payload: '{"name":"Token"}',
        signatures: [],
      },
      rules_container: 'base64-rules',
      rules_signatures: 'base64-sigs',
      business_rule_enabled: true,
    };

    const result = signedWhitelistedContractEnvelopeFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.tenantId).toBe('tenant-1');
    expect(result!.signedContractAddress).toBeDefined();
    expect(result!.rulesContainer).toBe('base64-rules');
    expect(result!.rulesSignatures).toBe('base64-sigs');
    expect(result!.businessRuleEnabled).toBe(true);
  });
});

describe('whitelistedContractFromEnvelope - payload extraction', () => {
  it('should prefer metadata payload over signed contract payload', () => {
    const envelope = {
      id: 'env-1',
      metadata: {
        hash: 'h1',
        payloadAsString: JSON.stringify({
          contractAddress: '0xFROM_METADATA',
          symbol: 'META',
          name: 'Metadata Token',
          decimals: 18,
        }),
      },
      signedContractAddress: {
        payload: JSON.stringify({
          contractAddress: '0xFROM_SIGNED',
          symbol: 'SIGNED',
          name: 'Signed Token',
          decimals: 6,
        }),
        signatures: [],
      },
      attributes: [],
    } as any;

    const result = whitelistedContractFromEnvelope(envelope);
    expect(result.contractAddress).toBe('0xFROM_METADATA');
    expect(result.symbol).toBe('META');
    expect(result.name).toBe('Metadata Token');
    expect(result.decimals).toBe(18);
  });

  it('should handle malformed metadata payload gracefully', () => {
    const envelope = {
      id: 'env-2',
      metadata: {
        hash: 'h1',
        payloadAsString: 'not-json',
      },
      signedContractAddress: {
        payload: JSON.stringify({ symbol: 'FB' }),
        signatures: [],
      },
      attributes: [],
    } as any;

    const result = whitelistedContractFromEnvelope(envelope);
    // Falls back to signed contract payload
    expect(result.symbol).toBe('FB');
  });
});

describe('whitelistedContractResultFromDto', () => {
  it('should map result with total_items snake_case', () => {
    const dto = {
      result: [{ id: 'wc-1' }],
      total_items: 25,
    };

    const result = whitelistedContractResultFromDto(dto);
    expect(result.contracts).toHaveLength(1);
    expect(result.totalItems).toBe(25);
  });

  it('should return defaults for empty result', () => {
    const result = whitelistedContractResultFromDto({});
    expect(result.contracts).toEqual([]);
    expect(result.totalItems).toBe(0);
  });
});
