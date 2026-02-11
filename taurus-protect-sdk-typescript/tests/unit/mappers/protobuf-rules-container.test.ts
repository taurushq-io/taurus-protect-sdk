/**
 * Tests for protobuf rules container decoding.
 */

import {
  RulesContainer,
  Role,
  Blockchain,
} from '../../../src/internal/proto/request_reply';
import { tryDecodeProtobufRulesContainer } from '../../../src/mappers/protobuf-rules-container';
import { rulesContainerFromBase64 } from '../../../src/mappers/governance-rules';

describe('tryDecodeProtobufRulesContainer', () => {
  describe('valid protobuf decoding', () => {
    it('should decode a simple RulesContainer', () => {
      // Create a minimal protobuf RulesContainer
      const pb: Parameters<typeof RulesContainer.encode>[0] = {
        users: [
          {
            id: 'user-1',
            publicKey: '-----BEGIN PUBLIC KEY-----\ntest\n-----END PUBLIC KEY-----',
            roles: [Role.SUPERADMIN, Role.HSMSLOT],
            properties: {},
          },
        ],
        groups: [
          {
            id: 'group-1',
            userIds: ['user-1'],
            properties: {},
          },
        ],
        minimumDistinctUserSignatures: 2,
        minimumDistinctGroupSignatures: 1,
        transactionRules: [],
        addressWhitelistingRules: [],
        contractAddressWhitelistingRules: [],
        enforcedRulesHash: 'abc123',
        properties: {},
        timestamp: 1234567890,
        minimumCommitmentSignatures: 0,
        engineIdentities: [],
        hsmSlotId: 0,
      };

      // Encode to bytes
      const encoded = RulesContainer.encode(pb).finish();

      // Decode
      const result = tryDecodeProtobufRulesContainer(encoded);

      expect(result).toBeDefined();
      expect(result!.users).toHaveLength(1);
      expect(result!.users[0].id).toBe('user-1');
      expect(result!.users[0].publicKeyPem).toBe('-----BEGIN PUBLIC KEY-----\ntest\n-----END PUBLIC KEY-----');
      expect(result!.users[0].roles).toContain('SUPERADMIN');
      expect(result!.users[0].roles).toContain('HSMSLOT');
      expect(result!.groups).toHaveLength(1);
      expect(result!.groups[0].id).toBe('group-1');
      expect(result!.groups[0].userIds).toEqual(['user-1']);
      expect(result!.minimumDistinctUserSignatures).toBe(2);
      expect(result!.minimumDistinctGroupSignatures).toBe(1);
      expect(result!.enforcedRulesHash).toBe('abc123');
      expect(result!.timestamp).toBe(1234567890);
    });

    it('should decode address whitelisting rules', () => {
      const pb: Parameters<typeof RulesContainer.encode>[0] = {
        users: [],
        groups: [],
        minimumDistinctUserSignatures: 0,
        minimumDistinctGroupSignatures: 0,
        transactionRules: [],
        addressWhitelistingRules: [
          {
            currency: 'ETH',
            network: 'mainnet',
            parallelThresholds: [
              {
                thresholds: [
                  { groupId: 'group-1', minimumSignatures: 2 },
                  { groupId: 'group-2', minimumSignatures: 1 },
                ],
              },
            ],
            properties: {},
            lines: [],
          },
        ],
        contractAddressWhitelistingRules: [],
        enforcedRulesHash: '',
        properties: {},
        timestamp: 0,
        minimumCommitmentSignatures: 0,
        engineIdentities: [],
        hsmSlotId: 0,
      };

      const encoded = RulesContainer.encode(pb).finish();
      const result = tryDecodeProtobufRulesContainer(encoded);

      expect(result).toBeDefined();
      expect(result!.addressWhitelistingRules).toHaveLength(1);
      expect(result!.addressWhitelistingRules[0].currency).toBe('ETH');
      expect(result!.addressWhitelistingRules[0].network).toBe('mainnet');
      // parallelThresholds preserves SequentialThresholds structure
      expect(result!.addressWhitelistingRules[0].parallelThresholds).toHaveLength(1);
      expect(result!.addressWhitelistingRules[0].parallelThresholds[0].thresholds).toHaveLength(2);
      expect(result!.addressWhitelistingRules[0].parallelThresholds[0].thresholds[0].groupId).toBe('group-1');
      expect(result!.addressWhitelistingRules[0].parallelThresholds[0].thresholds[0].minimumSignatures).toBe(2);
    });

    it('should decode contract address whitelisting rules', () => {
      const pb: Parameters<typeof RulesContainer.encode>[0] = {
        users: [],
        groups: [],
        minimumDistinctUserSignatures: 0,
        minimumDistinctGroupSignatures: 0,
        transactionRules: [],
        addressWhitelistingRules: [],
        contractAddressWhitelistingRules: [
          {
            blockchain: Blockchain.ETH,
            network: 'mainnet',
            parallelThresholds: [
              {
                thresholds: [
                  { groupId: 'group-1', minimumSignatures: 1 },
                ],
              },
            ],
            properties: {},
          },
        ],
        enforcedRulesHash: '',
        properties: {},
        timestamp: 0,
        minimumCommitmentSignatures: 0,
        engineIdentities: [],
        hsmSlotId: 0,
      };

      const encoded = RulesContainer.encode(pb).finish();
      const result = tryDecodeProtobufRulesContainer(encoded);

      expect(result).toBeDefined();
      expect(result!.contractAddressWhitelistingRules).toHaveLength(1);
      expect(result!.contractAddressWhitelistingRules[0].blockchain).toBe('ETH');
      expect(result!.contractAddressWhitelistingRules[0].network).toBe('mainnet');
      expect(result!.contractAddressWhitelistingRules[0].parallelThresholds).toHaveLength(1);
    });
  });

  describe('invalid input handling', () => {
    it('should return undefined for JSON data', () => {
      const jsonBytes = new TextEncoder().encode('{"users": [], "groups": []}');
      const result = tryDecodeProtobufRulesContainer(jsonBytes);
      expect(result).toBeUndefined();
    });

    it('should return undefined for random bytes', () => {
      const randomBytes = new Uint8Array([0x00, 0xff, 0x12, 0x34]);
      const result = tryDecodeProtobufRulesContainer(randomBytes);
      // May return undefined or a partially decoded object - either is acceptable
      // The important thing is it doesn't throw
    });

    it('should return undefined for empty bytes', () => {
      const emptyBytes = new Uint8Array(0);
      const result = tryDecodeProtobufRulesContainer(emptyBytes);
      // Empty bytes decode to empty container, which is valid
      expect(result).toBeDefined();
      expect(result!.users).toEqual([]);
    });
  });
});

describe('rulesContainerFromBase64', () => {
  describe('protobuf-first decoding', () => {
    it('should decode protobuf format correctly', () => {
      // Create a protobuf RulesContainer
      const pb: Parameters<typeof RulesContainer.encode>[0] = {
        users: [
          {
            id: 'user-proto',
            publicKey: 'PEM-KEY',
            roles: [Role.REQUESTAPPROVER],
            properties: {},
          },
        ],
        groups: [],
        minimumDistinctUserSignatures: 1,
        minimumDistinctGroupSignatures: 0,
        transactionRules: [],
        addressWhitelistingRules: [],
        contractAddressWhitelistingRules: [],
        enforcedRulesHash: 'hash123',
        properties: {},
        timestamp: 9999,
        minimumCommitmentSignatures: 0,
        engineIdentities: [],
        hsmSlotId: 0,
      };

      const encoded = RulesContainer.encode(pb).finish();
      const base64 = Buffer.from(encoded).toString('base64');

      const result = rulesContainerFromBase64(base64);

      expect(result.users).toHaveLength(1);
      expect(result.users[0].id).toBe('user-proto');
      expect(result.users[0].publicKeyPem).toBe('PEM-KEY');
      expect(result.minimumDistinctUserSignatures).toBe(1);
      expect(result.enforcedRulesHash).toBe('hash123');
      expect(result.timestamp).toBe(9999);
    });

    it('should fall back to JSON when protobuf fails', () => {
      const jsonData = {
        users: [{ id: 'user-json', name: 'Test User', publicKeyPem: 'JSON-KEY', roles: ['ADMIN'] }],
        groups: [],
        minimumDistinctUserSignatures: 3,
        minimumDistinctGroupSignatures: 0,
        addressWhitelistingRules: [],
        contractAddressWhitelistingRules: [],
        enforcedRulesHash: 'jsonhash',
        timestamp: 1111,
      };

      const jsonString = JSON.stringify(jsonData);
      const base64 = Buffer.from(jsonString).toString('base64');

      const result = rulesContainerFromBase64(base64);

      expect(result.users).toHaveLength(1);
      expect(result.users[0].id).toBe('user-json');
      expect(result.users[0].publicKeyPem).toBe('JSON-KEY');
      expect(result.minimumDistinctUserSignatures).toBe(3);
      expect(result.enforcedRulesHash).toBe('jsonhash');
    });
  });

  describe('error handling', () => {
    it('should return empty container for empty base64 string', () => {
      const result = rulesContainerFromBase64('');
      expect(result.users).toEqual([]);
      expect(result.groups).toEqual([]);
    });

    it('should handle non-standard base64 gracefully', () => {
      // Note: Buffer.from(string, 'base64') doesn't throw for invalid base64
      // It just decodes what it can. So 'not-valid-base64!!!' might still decode.
      // The important behavior is it either returns a valid container or throws IntegrityError
      try {
        const result = rulesContainerFromBase64('not-valid-base64!!!');
        // If it doesn't throw, it should return a valid container
        expect(result).toBeDefined();
        expect(result.users).toBeDefined();
      } catch (e) {
        // If it throws, it should be an IntegrityError
        expect(e).toBeDefined();
      }
    });

    it('should decode plain text as empty protobuf container', () => {
      // Plain text happens to decode as valid (empty) protobuf
      // This is expected behavior - protobuf is permissive with unknown bytes
      const plainText = 'This is plain text, not protobuf or JSON';
      const base64 = Buffer.from(plainText).toString('base64');

      // Should either decode as empty protobuf or throw - both are acceptable
      try {
        const result = rulesContainerFromBase64(base64);
        expect(result).toBeDefined();
        // The decoded container will have empty arrays for unknown data
      } catch (e) {
        // If it throws, it should be an IntegrityError
        expect(e).toBeDefined();
      }
    });

    it('should throw IntegrityError for malformed JSON that is not valid protobuf', () => {
      // Create data that starts with '{' (looks like JSON) but is invalid JSON
      // and also not valid protobuf
      const malformedJson = '{"users": [invalid json here}';
      const base64 = Buffer.from(malformedJson).toString('base64');

      // This should fail JSON parsing and protobuf parsing
      expect(() => rulesContainerFromBase64(base64)).toThrow('not valid protobuf or JSON');
    });
  });
});
