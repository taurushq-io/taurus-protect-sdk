/**
 * Unit tests for contract whitelist mapper functions.
 */

import {
  whitelistedContractAttributeFromDto,
  whitelistedContractMetadataFromDto,
  whitelistedContractSignatureFromDto,
  signedWhitelistedContractFromDto,
  whitelistedContractTrailFromDto,
  whitelistedContractApproverFromDto,
  whitelistedContractApproverGroupFromDto,
  whitelistedContractApproversFromDto,
  signedWhitelistedContractEnvelopeFromDto,
  whitelistedContractFromEnvelope,
  whitelistedContractFromDto,
  whitelistedContractsFromDto,
  whitelistedContractResultFromDto,
} from '../../../src/mappers/contract-whitelist';

describe('whitelistedContractAttributeFromDto', () => {
  it('should map all fields', () => {
    const dto = {
      id: 'attr-1',
      key: 'description',
      value: 'Test contract',
      contentType: 'text/plain',
      type: 'custom',
      subtype: 'info',
      isfile: false,
    };

    const result = whitelistedContractAttributeFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.id).toBe('attr-1');
    expect(result!.key).toBe('description');
    expect(result!.value).toBe('Test contract');
    expect(result!.contentType).toBe('text/plain');
    expect(result!.type).toBe('custom');
    expect(result!.subType).toBe('info');
    expect(result!.isFile).toBe(false);
  });

  it('should return undefined for null input', () => {
    expect(whitelistedContractAttributeFromDto(null)).toBeUndefined();
  });
});

describe('whitelistedContractMetadataFromDto', () => {
  it('should map hash and payloadAsString', () => {
    const dto = {
      hash: 'abc123',
      payloadAsString: '{"contractAddress":"0x1234"}',
    };

    const result = whitelistedContractMetadataFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.hash).toBe('abc123');
    expect(result!.payloadAsString).toBe('{"contractAddress":"0x1234"}');
  });

  it('should handle snake_case payload_as_string', () => {
    const dto = {
      hash: 'xyz',
      payload_as_string: '{"name":"Token"}',
    };

    const result = whitelistedContractMetadataFromDto(dto);

    expect(result!.payloadAsString).toBe('{"name":"Token"}');
  });

  it('should return undefined for null input', () => {
    expect(whitelistedContractMetadataFromDto(null)).toBeUndefined();
  });
});

describe('whitelistedContractSignatureFromDto', () => {
  it('should map signature fields', () => {
    const dto = {
      signature: 'sig-value',
      comment: 'Approved',
      hashes: ['hash1', 'hash2'],
      userId: 'user-1',
    };

    const result = whitelistedContractSignatureFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.signature).toBe('sig-value');
    expect(result!.comment).toBe('Approved');
    expect(result!.hashes).toEqual(['hash1', 'hash2']);
    expect(result!.userId).toBe('user-1');
  });

  it('should handle missing hashes', () => {
    const dto = { signature: 'sig' };

    const result = whitelistedContractSignatureFromDto(dto);

    expect(result!.hashes).toEqual([]);
  });

  it('should return undefined for null input', () => {
    expect(whitelistedContractSignatureFromDto(null)).toBeUndefined();
  });
});

describe('signedWhitelistedContractFromDto', () => {
  it('should map payload and signatures', () => {
    const dto = {
      payload: '{"contractAddress":"0x1234"}',
      signatures: [
        { signature: 'sig1', hashes: ['h1'] },
      ],
    };

    const result = signedWhitelistedContractFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.payload).toBe('{"contractAddress":"0x1234"}');
    expect(result!.signatures).toHaveLength(1);
  });

  it('should return undefined for null input', () => {
    expect(signedWhitelistedContractFromDto(null)).toBeUndefined();
  });
});

describe('whitelistedContractTrailFromDto', () => {
  it('should map trail fields', () => {
    const dto = {
      id: 'trail-1',
      action: 'created',
      userId: 'user-1',
      comment: 'Initial creation',
      timestamp: '2024-01-01T00:00:00Z',
    };

    const result = whitelistedContractTrailFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.id).toBe('trail-1');
    expect(result!.action).toBe('created');
    expect(result!.userId).toBe('user-1');
    expect(result!.comment).toBe('Initial creation');
    expect(result!.timestamp).toBe('2024-01-01T00:00:00Z');
  });

  it('should return undefined for null input', () => {
    expect(whitelistedContractTrailFromDto(null)).toBeUndefined();
  });
});

describe('whitelistedContractApproverFromDto', () => {
  it('should map approver fields', () => {
    const dto = {
      userId: 'u-1',
      userName: 'John Doe',
      pending: true,
    };

    const result = whitelistedContractApproverFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.userId).toBe('u-1');
    expect(result!.userName).toBe('John Doe');
    expect(result!.pending).toBe(true);
  });

  it('should return undefined for null input', () => {
    expect(whitelistedContractApproverFromDto(null)).toBeUndefined();
  });
});

describe('whitelistedContractApproverGroupFromDto', () => {
  it('should map group with users', () => {
    const dto = {
      id: 'g-1',
      name: 'Admins',
      required: 2,
      users: [
        { userId: 'u-1', userName: 'Alice', pending: false },
        { userId: 'u-2', userName: 'Bob', pending: true },
      ],
    };

    const result = whitelistedContractApproverGroupFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.id).toBe('g-1');
    expect(result!.name).toBe('Admins');
    expect(result!.required).toBe(2);
    expect(result!.users).toHaveLength(2);
  });

  it('should return undefined for null input', () => {
    expect(whitelistedContractApproverGroupFromDto(null)).toBeUndefined();
  });
});

describe('whitelistedContractApproversFromDto', () => {
  it('should map approvers with groups', () => {
    const dto = {
      required: 1,
      groups: [
        { id: 'g-1', name: 'Group1', required: 1, users: [] },
      ],
    };

    const result = whitelistedContractApproversFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.required).toBe(1);
    expect(result!.groups).toHaveLength(1);
  });

  it('should return undefined for null input', () => {
    expect(whitelistedContractApproversFromDto(null)).toBeUndefined();
  });
});

describe('signedWhitelistedContractEnvelopeFromDto', () => {
  it('should map full envelope DTO', () => {
    const dto = {
      id: 'env-1',
      tenantId: 't-1',
      blockchain: 'ETH',
      network: 'mainnet',
      metadata: { hash: 'h1', payloadAsString: '{"name":"Token"}' },
      action: 'CREATED',
      status: 'ACTIVE',
      businessRuleEnabled: true,
    };

    const result = signedWhitelistedContractEnvelopeFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.id).toBe('env-1');
    expect(result!.tenantId).toBe('t-1');
    expect(result!.blockchain).toBe('ETH');
    expect(result!.network).toBe('mainnet');
    expect(result!.metadata).toBeDefined();
    expect(result!.action).toBe('CREATED');
    expect(result!.status).toBe('ACTIVE');
    expect(result!.businessRuleEnabled).toBe(true);
  });

  it('should return undefined for null input', () => {
    expect(signedWhitelistedContractEnvelopeFromDto(null)).toBeUndefined();
  });
});

describe('whitelistedContractFromEnvelope', () => {
  it('should extract contract details from metadata payload', () => {
    const envelope = {
      id: 'env-1',
      tenantId: 't-1',
      blockchain: 'ETH',
      network: 'mainnet',
      metadata: {
        hash: 'h1',
        payloadAsString: JSON.stringify({
          contractAddress: '0xABC',
          symbol: 'TKN',
          name: 'Token',
          decimals: 18,
          kind: 'ERC20',
          tokenId: 'tid-1',
        }),
      },
      status: 'ACTIVE',
      businessRuleEnabled: false,
      attributes: [],
    } as any;

    const result = whitelistedContractFromEnvelope(envelope);

    expect(result.id).toBe('env-1');
    expect(result.contractAddress).toBe('0xABC');
    expect(result.symbol).toBe('TKN');
    expect(result.name).toBe('Token');
    expect(result.decimals).toBe(18);
    expect(result.kind).toBe('ERC20');
    expect(result.tokenId).toBe('tid-1');
  });

  it('should fall back to signed contract payload', () => {
    const envelope = {
      id: 'env-2',
      metadata: { hash: 'h1' },
      signedContractAddress: {
        payload: JSON.stringify({
          contractAddress: '0xDEF',
          symbol: 'FB',
        }),
        signatures: [],
      },
      attributes: [],
    } as any;

    const result = whitelistedContractFromEnvelope(envelope);

    expect(result.contractAddress).toBe('0xDEF');
    expect(result.symbol).toBe('FB');
  });
});

describe('whitelistedContractsFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      {
        id: 'e1',
        metadata: { payloadAsString: '{"contractAddress":"0x1"}' },
      },
      {
        id: 'e2',
        metadata: { payloadAsString: '{"contractAddress":"0x2"}' },
      },
    ];

    const result = whitelistedContractsFromDto(dtos);

    expect(result).toHaveLength(2);
  });

  it('should return empty array for null input', () => {
    expect(whitelistedContractsFromDto(null)).toEqual([]);
  });
});

describe('whitelistedContractResultFromDto', () => {
  it('should map result with total items', () => {
    const dto = {
      result: [{ id: 'e1' }],
      totalItems: 10,
    };

    const result = whitelistedContractResultFromDto(dto);

    expect(result.contracts).toHaveLength(1);
    expect(result.totalItems).toBe(10);
  });

  it('should return default for null input', () => {
    const result = whitelistedContractResultFromDto(null);

    expect(result.contracts).toEqual([]);
    expect(result.totalItems).toBe(0);
  });
});
