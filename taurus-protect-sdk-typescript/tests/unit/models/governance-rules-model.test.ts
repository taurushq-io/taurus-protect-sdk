/**
 * Coverage for the DecodedRulesContainer model helper functions — parity with
 * the Python/Go model tests (find-by-id, priority-based whitelisting lookups,
 * HSM key resolution, empty-container factory).
 */
import {
  createEmptyRulesContainer,
  findAddressWhitelistingRules,
  findContractAddressWhitelistingRules,
  findGroupById,
  findUserById,
  getHsmPublicKey,
  hasUnknownFields,
} from '../../../src/models/governance-rules';
import type {
  AddressWhitelistingRules,
  ContractAddressWhitelistingRules,
  DecodedRulesContainer,
} from '../../../src/models/governance-rules';

const user = (id: string, roles: string[] = [], publicKeyPem?: string) => ({ id, name: undefined, publicKeyPem, roles });

describe('createEmptyRulesContainer', () => {
  it('returns an empty, no-unknown-fields container', () => {
    const c = createEmptyRulesContainer();
    expect(c.users).toEqual([]);
    expect(c.transactionRules).toEqual([]);
    expect(c.timestamp).toBe(0);
    expect(hasUnknownFields(c)).toBe(false);
  });
});

describe('findUserById / findGroupById', () => {
  const c: DecodedRulesContainer = {
    ...createEmptyRulesContainer(),
    users: [user('u1'), user('u2')],
    groups: [{ id: 'g1', name: undefined, userIds: ['u1'] }],
  };
  it('finds present entities and returns undefined for absent ones', () => {
    expect(findUserById(c, 'u2')?.id).toBe('u2');
    expect(findUserById(c, 'nope')).toBeUndefined();
    expect(findGroupById(c, 'g1')?.id).toBe('g1');
    expect(findGroupById(c, 'nope')).toBeUndefined();
  });
});

describe('getHsmPublicKey', () => {
  it('returns the HSMSLOT user public key, or undefined', () => {
    const withHsm: DecodedRulesContainer = {
      ...createEmptyRulesContainer(),
      users: [user('auth', ['SUPERADMIN']), user('hsm', ['HSMSLOT'], 'PEM-KEY')],
    };
    expect(getHsmPublicKey(withHsm)).toBe('PEM-KEY');

    // HSMSLOT user without a key -> undefined
    const noKey: DecodedRulesContainer = { ...createEmptyRulesContainer(), users: [user('hsm', ['HSMSLOT'])] };
    expect(getHsmPublicKey(noKey)).toBeUndefined();

    // no HSMSLOT user -> undefined
    const noHsm: DecodedRulesContainer = { ...createEmptyRulesContainer(), users: [user('u', ['SUPERADMIN'], 'PEM')] };
    expect(getHsmPublicKey(noHsm)).toBeUndefined();
  });
});

describe('findAddressWhitelistingRules priority tiers', () => {
  const exact: AddressWhitelistingRules = { currency: 'ETH', network: 'mainnet', parallelThresholds: [], lines: [] };
  const blockchainOnly: AddressWhitelistingRules = { currency: 'ETH', network: 'Any', parallelThresholds: [], lines: [] };
  const global: AddressWhitelistingRules = { currency: '', network: 'mainnet', parallelThresholds: [], lines: [] };
  const c: DecodedRulesContainer = { ...createEmptyRulesContainer(), addressWhitelistingRules: [exact, blockchainOnly, global] };

  it('prefers exact > blockchain-only > global default', () => {
    expect(findAddressWhitelistingRules(c, 'ETH', 'mainnet')).toBe(exact);
    expect(findAddressWhitelistingRules(c, 'ETH', 'testnet')).toBe(blockchainOnly);
    expect(findAddressWhitelistingRules(c, 'BTC', 'mainnet')).toBe(global);
  });

  it('returns undefined when nothing matches and there is no global default', () => {
    const noGlobal: DecodedRulesContainer = { ...createEmptyRulesContainer(), addressWhitelistingRules: [exact] };
    expect(findAddressWhitelistingRules(noGlobal, 'SOL', 'x')).toBeUndefined();
  });
});

describe('findContractAddressWhitelistingRules priority tiers', () => {
  const exact: ContractAddressWhitelistingRules = { blockchain: 'ETH', network: 'mainnet', parallelThresholds: [] };
  const blockchainOnly: ContractAddressWhitelistingRules = { blockchain: 'ETH', network: 'Any', parallelThresholds: [] };
  const global: ContractAddressWhitelistingRules = { blockchain: '', network: 'mainnet', parallelThresholds: [] };
  const c: DecodedRulesContainer = { ...createEmptyRulesContainer(), contractAddressWhitelistingRules: [exact, blockchainOnly, global] };

  it('prefers exact > blockchain-only > global default', () => {
    expect(findContractAddressWhitelistingRules(c, 'ETH', 'mainnet')).toBe(exact);
    expect(findContractAddressWhitelistingRules(c, 'ETH', 'testnet')).toBe(blockchainOnly);
    expect(findContractAddressWhitelistingRules(c, 'BTC', 'mainnet')).toBe(global);
  });
});
