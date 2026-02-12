/**
 * Unit tests for wallet mapper functions.
 */

import {
  walletFromDto,
  walletFromCreateDto,
  walletsFromDto,
  balanceHistoryFromDto,
  walletBalanceFromDto,
  walletAttributeFromDto,
} from '../../../src/mappers/wallet';
import { WalletStatus } from '../../../src/models/wallet';

describe('walletFromDto', () => {
  it('should map all fields correctly from WalletInfo DTO', () => {
    const dto = {
      id: '123',
      name: 'My Wallet',
      blockchain: 'ETH',
      network: 'mainnet',
      currency: 'ETH',
      isOmnibus: false,
      disabled: false,
      balance: {
        totalConfirmed: '1000',
        totalUnconfirmed: '500',
      },
      creationDate: new Date('2024-01-01'),
      updateDate: new Date('2024-06-01'),
      comment: 'Test wallet',
      customerId: 'cust-1',
      addressesCount: '5',
      attributes: [
        { id: 'attr1', key: 'department', value: 'treasury' },
      ],
      visibilityGroupID: 'vg-1',
      externalWalletId: 'ext-123',
    };

    const wallet = walletFromDto(dto as any);

    expect(wallet).toBeDefined();
    expect(wallet!.id).toBe('123');
    expect(wallet!.name).toBe('My Wallet');
    expect(wallet!.blockchain).toBe('ETH');
    expect(wallet!.network).toBe('mainnet');
    expect(wallet!.currency).toBe('ETH');
    expect(wallet!.isOmnibus).toBe(false);
    expect(wallet!.status).toBe(WalletStatus.ACTIVE);
    expect(wallet!.balance?.totalConfirmed).toBe('1000');
    expect(wallet!.comment).toBe('Test wallet');
    expect(wallet!.customerId).toBe('cust-1');
    expect(wallet!.addressesCount).toBe(5);
    expect(wallet!.attributes).toHaveLength(1);
    expect(wallet!.visibilityGroupId).toBe('vg-1');
    expect(wallet!.externalWalletId).toBe('ext-123');
  });

  it('should return undefined for null input', () => {
    expect(walletFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(walletFromDto(undefined)).toBeUndefined();
  });

  it('should set status to DISABLED when disabled is true', () => {
    const dto = { id: '1', name: 'Test', disabled: true };
    const wallet = walletFromDto(dto as any);
    expect(wallet!.status).toBe(WalletStatus.DISABLED);
  });

  it('should set status to ACTIVE when disabled is false', () => {
    const dto = { id: '1', name: 'Test', disabled: false };
    const wallet = walletFromDto(dto as any);
    expect(wallet!.status).toBe(WalletStatus.ACTIVE);
  });
});

describe('walletFromCreateDto', () => {
  it('should map fields from TgvalidatordWallet DTO', () => {
    const dto = {
      id: '456',
      name: 'New Wallet',
      blockchain: 'BTC',
      currency: 'BTC',
      isOmnibus: true,
    };

    const wallet = walletFromCreateDto(dto as any);

    expect(wallet).toBeDefined();
    expect(wallet!.id).toBe('456');
    expect(wallet!.name).toBe('New Wallet');
    expect(wallet!.blockchain).toBe('BTC');
    expect(wallet!.isOmnibus).toBe(true);
  });

  it('should return undefined for null input', () => {
    expect(walletFromCreateDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(walletFromCreateDto(undefined)).toBeUndefined();
  });
});

describe('walletsFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { id: '1', name: 'Wallet 1' },
      { id: '2', name: 'Wallet 2' },
    ];

    const wallets = walletsFromDto(dtos as any);

    expect(wallets).toHaveLength(2);
    expect(wallets[0].name).toBe('Wallet 1');
    expect(wallets[1].name).toBe('Wallet 2');
  });

  it('should return empty array for null input', () => {
    expect(walletsFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(walletsFromDto(undefined)).toEqual([]);
  });

  it('should return empty array for empty array', () => {
    expect(walletsFromDto([])).toEqual([]);
  });
});

describe('balanceHistoryFromDto', () => {
  it('should map balance history points', () => {
    const dtos = [
      {
        pointDate: new Date('2024-01-01'),
        balance: { totalConfirmed: '1000' },
      },
      {
        pointDate: new Date('2024-01-02'),
        balance: { totalConfirmed: '1500' },
      },
    ];

    const history = balanceHistoryFromDto(dtos as any);

    expect(history).toHaveLength(2);
    expect(history[0].balance?.totalConfirmed).toBe('1000');
    expect(history[1].balance?.totalConfirmed).toBe('1500');
  });

  it('should return empty array for null input', () => {
    expect(balanceHistoryFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(balanceHistoryFromDto(undefined)).toEqual([]);
  });
});

describe('walletBalanceFromDto', () => {
  it('should map balance fields', () => {
    const dto = {
      totalConfirmed: '100',
      totalUnconfirmed: '50',
      availableConfirmed: '80',
      availableUnconfirmed: '30',
      reservedConfirmed: '20',
      reservedUnconfirmed: '20',
    };

    const balance = walletBalanceFromDto(dto as any);

    expect(balance).toBeDefined();
    expect(balance!.totalConfirmed).toBe('100');
    expect(balance!.availableConfirmed).toBe('80');
    expect(balance!.reservedConfirmed).toBe('20');
  });

  it('should return undefined for null input', () => {
    expect(walletBalanceFromDto(null)).toBeUndefined();
  });
});

describe('walletAttributeFromDto', () => {
  it('should map attribute fields', () => {
    const dto = {
      id: 'attr-1',
      key: 'department',
      value: 'treasury',
      contentType: 'text/plain',
      owner: 'admin',
      type: 'custom',
    };

    const attr = walletAttributeFromDto(dto as any);

    expect(attr).toBeDefined();
    expect(attr!.id).toBe('attr-1');
    expect(attr!.key).toBe('department');
    expect(attr!.value).toBe('treasury');
  });

  it('should return undefined for null input', () => {
    expect(walletAttributeFromDto(null)).toBeUndefined();
  });
});
