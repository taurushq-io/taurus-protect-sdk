/**
 * Unit tests for fiat provider mapper functions.
 */

import {
  fiatProviderFromDto,
  fiatProvidersFromDto,
  fiatProviderAccountFromDto,
  fiatProviderCounterpartyAccountFromDto,
  fiatProviderOperationFromDto,
  fiatResponseCursorFromDto,
  fiatProviderAccountResultFromDto,
  fiatProviderCounterpartyAccountResultFromDto,
  fiatProviderOperationResultFromDto,
} from '../../../src/mappers/fiat';

describe('fiatProviderFromDto', () => {
  it('should map all fields', () => {
    const dto = {
      provider: 'bank-abc',
      label: 'Bank ABC',
      baseCurrencyValuation: '100000',
    };

    const result = fiatProviderFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.provider).toBe('bank-abc');
    expect(result!.label).toBe('Bank ABC');
    expect(result!.baseCurrencyValuation).toBe('100000');
  });

  it('should return undefined for null input', () => {
    expect(fiatProviderFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(fiatProviderFromDto(undefined)).toBeUndefined();
  });
});

describe('fiatProvidersFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { provider: 'bank-a', label: 'A' },
      { provider: 'bank-b', label: 'B' },
    ];

    const result = fiatProvidersFromDto(dtos);

    expect(result).toHaveLength(2);
  });

  it('should return empty array for null input', () => {
    expect(fiatProvidersFromDto(null)).toEqual([]);
  });
});

describe('fiatProviderAccountFromDto', () => {
  it('should map all fields', () => {
    const dto = {
      id: 'acc-1',
      provider: 'bank-abc',
      label: 'Main Account',
      accountType: 'checking',
      accountIdentifier: 'CH1234',
      accountName: 'Treasury',
      totalBalance: '50000',
      currencyID: 'USD',
      baseCurrencyValuation: '50000',
      creationDate: new Date('2024-01-01'),
      updateDate: new Date('2024-06-01'),
    };

    const result = fiatProviderAccountFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.id).toBe('acc-1');
    expect(result!.provider).toBe('bank-abc');
    expect(result!.label).toBe('Main Account');
    expect(result!.accountType).toBe('checking');
    expect(result!.accountIdentifier).toBe('CH1234');
    expect(result!.accountName).toBe('Treasury');
    expect(result!.totalBalance).toBe('50000');
    expect(result!.currencyId).toBe('USD');
    expect(result!.baseCurrencyValuation).toBe('50000');
    expect(result!.creationDate).toBeInstanceOf(Date);
    expect(result!.updateDate).toBeInstanceOf(Date);
  });

  it('should handle currencyId alias', () => {
    const dto = { id: 'acc-2', currencyId: 'EUR' };

    const result = fiatProviderAccountFromDto(dto);

    expect(result!.currencyId).toBe('EUR');
  });

  it('should return undefined for null input', () => {
    expect(fiatProviderAccountFromDto(null)).toBeUndefined();
  });
});

describe('fiatProviderCounterpartyAccountFromDto', () => {
  it('should map all fields', () => {
    const dto = {
      id: 'cp-1',
      provider: 'bank-abc',
      label: 'Counterparty Acc',
      accountType: 'savings',
      accountIdentifier: 'DE5678',
      accountName: 'Partner',
      counterpartyID: 'partner-1',
      counterpartyName: 'Partner Corp',
      currencyID: 'EUR',
      creationDate: new Date('2024-02-01'),
      updateDate: new Date('2024-05-01'),
    };

    const result = fiatProviderCounterpartyAccountFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.id).toBe('cp-1');
    expect(result!.counterpartyId).toBe('partner-1');
    expect(result!.counterpartyName).toBe('Partner Corp');
    expect(result!.currencyId).toBe('EUR');
  });

  it('should return undefined for null input', () => {
    expect(fiatProviderCounterpartyAccountFromDto(null)).toBeUndefined();
  });
});

describe('fiatProviderOperationFromDto', () => {
  it('should map all fields', () => {
    const dto = {
      id: 'op-1',
      provider: 'bank-abc',
      label: 'Transfer',
      operationType: 'SEPA',
      operationIdentifier: 'TX-123',
      operationDirection: 'outgoing',
      status: 'completed',
      amount: '1000.50',
      currencyID: 'EUR',
      fromAccountID: 'acc-1',
      toAccountID: 'acc-2',
      fromDetails: 'Treasury',
      toDetails: 'Partner',
      comment: 'Monthly payment',
      operationDetails: '{"ref":"INV-001"}',
      creationDate: new Date('2024-03-01'),
      updateDate: new Date('2024-03-02'),
    };

    const result = fiatProviderOperationFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.id).toBe('op-1');
    expect(result!.operationType).toBe('SEPA');
    expect(result!.operationDirection).toBe('outgoing');
    expect(result!.status).toBe('completed');
    expect(result!.amount).toBe('1000.50');
    expect(result!.currencyId).toBe('EUR');
    expect(result!.fromAccountId).toBe('acc-1');
    expect(result!.toAccountId).toBe('acc-2');
    expect(result!.comment).toBe('Monthly payment');
  });

  it('should return undefined for null input', () => {
    expect(fiatProviderOperationFromDto(null)).toBeUndefined();
  });
});

describe('fiatResponseCursorFromDto', () => {
  it('should map cursor fields', () => {
    const dto = {
      currentPage: 'page-1',
      nextPage: 'page-2',
    };

    const result = fiatResponseCursorFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.currentPage).toBe('page-1');
    expect(result!.nextPage).toBe('page-2');
    expect(result!.hasMore).toBe(true);
  });

  it('should set hasMore to false when no nextPage', () => {
    const dto = { currentPage: 'page-1' };

    const result = fiatResponseCursorFromDto(dto);

    expect(result!.hasMore).toBe(false);
  });

  it('should return undefined for null input', () => {
    expect(fiatResponseCursorFromDto(null)).toBeUndefined();
  });
});

describe('fiatProviderAccountResultFromDto', () => {
  it('should map result with accounts and cursor', () => {
    const dto = {
      result: [{ id: 'acc-1', provider: 'bank-a' }],
      cursor: { currentPage: 'p1', nextPage: 'p2' },
    };

    const result = fiatProviderAccountResultFromDto(dto);

    expect(result.accounts).toHaveLength(1);
    expect(result.cursor).toBeDefined();
  });

  it('should return default for null input', () => {
    const result = fiatProviderAccountResultFromDto(null);

    expect(result.accounts).toEqual([]);
  });
});

describe('fiatProviderCounterpartyAccountResultFromDto', () => {
  it('should map result with accounts', () => {
    const dto = {
      result: [{ id: 'cp-1' }],
    };

    const result = fiatProviderCounterpartyAccountResultFromDto(dto);

    expect(result.accounts).toHaveLength(1);
  });

  it('should return default for null input', () => {
    const result = fiatProviderCounterpartyAccountResultFromDto(null);

    expect(result.accounts).toEqual([]);
  });
});

describe('fiatProviderOperationResultFromDto', () => {
  it('should map result with operations', () => {
    const dto = {
      result: [{ id: 'op-1', status: 'completed' }],
      cursor: { currentPage: 'p1' },
    };

    const result = fiatProviderOperationResultFromDto(dto);

    expect(result.operations).toHaveLength(1);
    expect(result.cursor).toBeDefined();
  });

  it('should return default for null input', () => {
    const result = fiatProviderOperationResultFromDto(null);

    expect(result.operations).toEqual([]);
  });
});
