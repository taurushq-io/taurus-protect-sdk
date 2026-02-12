/**
 * Unit tests for action mapper functions.
 */

import {
  actionEnvelopeFromDto,
  actionEnvelopesFromDto,
  actionFromDto,
  actionAttributeFromDto,
  actionTrailFromDto,
  actionTargetFromDto,
  actionTriggerFromDto,
  actionAmountFromDto,
  actionSourceFromDto,
  actionDestinationFromDto,
  taskTransferFromDto,
  taskNotificationFromDto,
  actionTaskFromDto,
} from '../../../src/mappers/action';

describe('actionEnvelopeFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      id: 'act-1',
      tenantId: 't-1',
      label: 'Auto-rebalance',
      status: 'ACTIVE',
      autoApprove: true,
      creationDate: new Date('2024-01-01'),
      updateDate: new Date('2024-06-01'),
      lastcheckeddate: new Date('2024-06-02'),
      action: {
        trigger: { kind: 'BALANCE' },
        tasks: [],
      },
      attributes: [
        { id: 'a1', key: 'priority', value: 'high' },
      ],
      trails: [
        { id: 'tr-1', action: 'created', date: new Date('2024-01-15') },
      ],
    };

    const result = actionEnvelopeFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('act-1');
    expect(result!.tenantId).toBe('t-1');
    expect(result!.label).toBe('Auto-rebalance');
    expect(result!.status).toBe('ACTIVE');
    expect(result!.autoApprove).toBe(true);
    expect(result!.attributes).toHaveLength(1);
    expect(result!.trails).toHaveLength(1);
  });

  it('should return undefined for null input', () => {
    expect(actionEnvelopeFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(actionEnvelopeFromDto(undefined)).toBeUndefined();
  });

  it('should handle missing optional fields', () => {
    const dto = { id: 'act-2' };
    const result = actionEnvelopeFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('act-2');
  });
});

describe('actionEnvelopesFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { id: '1', label: 'Action 1' },
      { id: '2', label: 'Action 2' },
    ];
    const result = actionEnvelopesFromDto(dtos);
    expect(result).toHaveLength(2);
  });

  it('should return empty array for null input', () => {
    expect(actionEnvelopesFromDto(null)).toEqual([]);
  });
});

describe('actionFromDto', () => {
  it('should map trigger and tasks', () => {
    const dto = {
      trigger: { kind: 'BALANCE', balance: { target: { kind: 'ADDRESS' } } },
      tasks: [{ kind: 'TRANSFER' }],
    };
    const result = actionFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.trigger).toBeDefined();
    expect(result!.tasks).toHaveLength(1);
  });

  it('should return undefined for null input', () => {
    expect(actionFromDto(null)).toBeUndefined();
  });
});

describe('actionAttributeFromDto', () => {
  it('should map fields', () => {
    const dto = { id: 'a1', tenantId: 't-1', key: 'priority', value: 'high', contentType: 'text' };
    const result = actionAttributeFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('a1');
    expect(result!.key).toBe('priority');
    expect(result!.value).toBe('high');
  });

  it('should return undefined for null input', () => {
    expect(actionAttributeFromDto(null)).toBeUndefined();
  });
});

describe('actionTrailFromDto', () => {
  it('should map fields', () => {
    const dto = {
      id: 'tr-1',
      action: 'approved',
      comment: 'looks good',
      date: new Date('2024-01-15'),
      actionStatus: 'OK',
    };
    const result = actionTrailFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('tr-1');
    expect(result!.action).toBe('approved');
    expect(result!.comment).toBe('looks good');
  });

  it('should return undefined for null input', () => {
    expect(actionTrailFromDto(null)).toBeUndefined();
  });
});

describe('actionTargetFromDto', () => {
  it('should map nested address and wallet', () => {
    const dto = {
      kind: 'ADDRESS',
      address: { kind: 'SPECIFIC', addressId: 'a-1' },
      wallet: { kind: 'ANY', walletId: 'w-1' },
    };
    const result = actionTargetFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.kind).toBe('ADDRESS');
    expect(result!.address).toBeDefined();
    expect(result!.wallet).toBeDefined();
  });

  it('should return undefined for null input', () => {
    expect(actionTargetFromDto(null)).toBeUndefined();
  });
});

describe('actionAmountFromDto', () => {
  it('should map fields', () => {
    const dto = { kind: 'FIXED', cryptoAmount: '1.5' };
    const result = actionAmountFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.kind).toBe('FIXED');
    expect(result!.cryptoAmount).toBe('1.5');
  });

  it('should return undefined for null input', () => {
    expect(actionAmountFromDto(null)).toBeUndefined();
  });
});

describe('actionSourceFromDto', () => {
  it('should map fields with camelCase', () => {
    const dto = { kind: 'ADDRESS', addressId: 'a-1', walletId: 'w-1' };
    const result = actionSourceFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.addressId).toBe('a-1');
    expect(result!.walletId).toBe('w-1');
  });

  it('should return undefined for null input', () => {
    expect(actionSourceFromDto(null)).toBeUndefined();
  });
});

describe('actionDestinationFromDto', () => {
  it('should map fields', () => {
    const dto = {
      kind: 'WHITELISTED',
      whitelistedAddressId: 'wl-1',
    };
    const result = actionDestinationFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.whitelistedAddressId).toBe('wl-1');
  });

  it('should return undefined for null input', () => {
    expect(actionDestinationFromDto(null)).toBeUndefined();
  });
});

describe('actionTriggerFromDto', () => {
  it('should map trigger with balance', () => {
    const dto = {
      kind: 'BALANCE',
      balance: {
        target: { kind: 'ADDRESS' },
        comparator: { kind: 'LESS_THAN' },
        amount: { kind: 'FIXED', cryptoAmount: '10' },
      },
    };
    const result = actionTriggerFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.kind).toBe('BALANCE');
    expect(result!.balance).toBeDefined();
    expect(result!.balance!.amount).toBeDefined();
    expect(result!.balance!.amount!.cryptoAmount).toBe('10');
  });

  it('should return undefined for null input', () => {
    expect(actionTriggerFromDto(null)).toBeUndefined();
  });
});

describe('taskTransferFromDto', () => {
  it('should map nested from/to/amount', () => {
    const dto = {
      from: { kind: 'ADDRESS', addressId: 'a-1' },
      to: { kind: 'WHITELISTED', whitelistedAddressId: 'wl-1' },
      amount: { kind: 'ALL' },
      topUp: true,
      useAllFunds: false,
    };
    const result = taskTransferFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.from).toBeDefined();
    expect(result!.to).toBeDefined();
    expect(result!.topUp).toBe(true);
    expect(result!.useAllFunds).toBe(false);
  });

  it('should return undefined for null input', () => {
    expect(taskTransferFromDto(null)).toBeUndefined();
  });
});

describe('taskNotificationFromDto', () => {
  it('should map fields', () => {
    const dto = {
      emailAddresses: ['admin@example.com'],
      notificationMessage: 'Balance low',
      numberOfReminders: '3',
    };
    const result = taskNotificationFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.emailAddresses).toEqual(['admin@example.com']);
    expect(result!.notificationMessage).toBe('Balance low');
  });

  it('should return undefined for null input', () => {
    expect(taskNotificationFromDto(null)).toBeUndefined();
  });
});

describe('actionTaskFromDto', () => {
  it('should map transfer task', () => {
    const dto = {
      kind: 'TRANSFER',
      transfer: { from: { kind: 'ADDRESS' } },
    };
    const result = actionTaskFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.kind).toBe('TRANSFER');
    expect(result!.transfer).toBeDefined();
  });

  it('should return undefined for null input', () => {
    expect(actionTaskFromDto(null)).toBeUndefined();
  });
});
