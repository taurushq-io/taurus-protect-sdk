/**
 * Tests for rules container decoding from base64.
 *
 * These tests verify the JSON-based rules container decoding functionality
 * using the fixture data from whitelisted-address-raw-response.json.
 */

import * as fs from 'fs';
import * as path from 'path';
import {
  rulesContainerFromBase64,
  userSignaturesFromBase64,
} from '../../../src/mappers/governance-rules';
import { IntegrityError } from '../../../src/errors';

// Load fixture data
const fixtureData = JSON.parse(
  fs.readFileSync(
    path.join(__dirname, '../fixtures/whitelisted-address-raw-response.json'),
    'utf-8'
  )
);

// Encode the rulesContainerJson as base64 JSON for testing
const rulesContainerBase64 = Buffer.from(
  JSON.stringify(fixtureData.rulesContainerJson)
).toString('base64');

// The rulesSignatures from the fixture (protobuf format)
const rulesSignaturesBase64 = fixtureData.rulesSignatures;

describe('rulesContainerFromBase64', () => {
  describe('Category A: Rules container decoding tests', () => {
    it('should decode rules container from base64 successfully', () => {
      const decoded = rulesContainerFromBase64(rulesContainerBase64);

      expect(decoded).toBeDefined();
      expect(decoded.users).toBeDefined();
      expect(decoded.groups).toBeDefined();
      expect(decoded.addressWhitelistingRules).toBeDefined();
    });

    it('should decode container with 4 users', () => {
      const decoded = rulesContainerFromBase64(rulesContainerBase64);

      expect(decoded.users).toHaveLength(4);

      // Verify user IDs
      const userIds = decoded.users.map((u) => u.id);
      expect(userIds).toContain('superadmin1@bank.com');
      expect(userIds).toContain('superadmin2@bank.com');
      expect(userIds).toContain('team1@bank.com');
      expect(userIds).toContain('hsmslot@bank.com');
    });

    it('should decode container with 2 groups', () => {
      const decoded = rulesContainerFromBase64(rulesContainerBase64);

      expect(decoded.groups).toHaveLength(2);

      // Verify group IDs
      const groupIds = decoded.groups.map((g) => g.id);
      expect(groupIds).toContain('team1');
      expect(groupIds).toContain('superadmins');

      // Verify group memberships
      const team1Group = decoded.groups.find((g) => g.id === 'team1');
      expect(team1Group?.userIds).toEqual(['team1@bank.com']);

      const superadminsGroup = decoded.groups.find((g) => g.id === 'superadmins');
      expect(superadminsGroup?.userIds).toContain('superadmin1@bank.com');
      expect(superadminsGroup?.userIds).toContain('superadmin2@bank.com');
    });

    it('should decode users with PEM public keys', () => {
      const decoded = rulesContainerFromBase64(rulesContainerBase64);

      for (const user of decoded.users) {
        // The fixture uses 'publicKey' field, but mapper looks for 'publicKeyPem' or 'publicKey'
        // Check that we extracted the key (may come from the 'publicKey' field)
        expect(user.publicKeyPem).toBeDefined();
        expect(user.publicKeyPem).toContain('-----BEGIN PUBLIC KEY-----');
        expect(user.publicKeyPem).toContain('-----END PUBLIC KEY-----');
      }
    });

    it('should decode users with roles', () => {
      const decoded = rulesContainerFromBase64(rulesContainerBase64);

      // Find superadmin1 and check roles
      const superadmin1 = decoded.users.find((u) => u.id === 'superadmin1@bank.com');
      expect(superadmin1?.roles).toContain('SUPERADMIN');

      // Find team1 user and check roles
      const team1User = decoded.users.find((u) => u.id === 'team1@bank.com');
      expect(team1User?.roles).toContain('USER');
      expect(team1User?.roles).toContain('OPERATOR');
    });

    it('should find 2 users with SUPERADMIN role', () => {
      const decoded = rulesContainerFromBase64(rulesContainerBase64);

      const superadminUsers = decoded.users.filter((u) =>
        u.roles.includes('SUPERADMIN')
      );
      expect(superadminUsers).toHaveLength(2);

      const superadminIds = superadminUsers.map((u) => u.id);
      expect(superadminIds).toContain('superadmin1@bank.com');
      expect(superadminIds).toContain('superadmin2@bank.com');
    });

    it('should find 1 user with HSMSLOT role', () => {
      const decoded = rulesContainerFromBase64(rulesContainerBase64);

      const hsmslotUsers = decoded.users.filter((u) =>
        u.roles.includes('HSMSLOT')
      );
      expect(hsmslotUsers).toHaveLength(1);
      expect(hsmslotUsers[0].id).toBe('hsmslot@bank.com');
    });

    it('should decode 1 address whitelisting rule for ALGO/mainnet', () => {
      const decoded = rulesContainerFromBase64(rulesContainerBase64);

      expect(decoded.addressWhitelistingRules).toHaveLength(1);

      const algoRule = decoded.addressWhitelistingRules[0];
      expect(algoRule.currency).toBe('ALGO');
      expect(algoRule.network).toBe('mainnet');

      // Verify parallel thresholds (SequentialThresholds[] wrapping GroupThreshold[])
      expect(algoRule.parallelThresholds).toHaveLength(1);
      expect(algoRule.parallelThresholds[0].thresholds).toHaveLength(1);
      expect(algoRule.parallelThresholds[0].thresholds[0].groupId).toBe('team1');
      expect(algoRule.parallelThresholds[0].thresholds[0].minimumSignatures).toBe(1);
    });

    it('should throw IntegrityError for invalid base64', () => {
      // Create truly invalid base64 that can't be decoded at all
      // Note: Buffer.from is lenient with base64, so we need something that
      // decodes to invalid data that is neither valid protobuf nor JSON
      const malformedJson = '{"users": [invalid json';
      const base64 = Buffer.from(malformedJson).toString('base64');

      expect(() => rulesContainerFromBase64(base64)).toThrow(IntegrityError);
    });
  });

  describe('Additional edge cases', () => {
    it('should return empty container for empty string', () => {
      const decoded = rulesContainerFromBase64('');

      expect(decoded.users).toEqual([]);
      expect(decoded.groups).toEqual([]);
      expect(decoded.addressWhitelistingRules).toEqual([]);
    });

    it('should decode timestamp correctly', () => {
      const decoded = rulesContainerFromBase64(rulesContainerBase64);

      expect(decoded.timestamp).toBe(1706194800);
    });

    it('should decode minimumDistinctUserSignatures and minimumDistinctGroupSignatures', () => {
      const decoded = rulesContainerFromBase64(rulesContainerBase64);

      expect(decoded.minimumDistinctUserSignatures).toBe(0);
      expect(decoded.minimumDistinctGroupSignatures).toBe(0);
    });

    it('should decode empty contractAddressWhitelistingRules', () => {
      const decoded = rulesContainerFromBase64(rulesContainerBase64);

      expect(decoded.contractAddressWhitelistingRules).toEqual([]);
    });
  });
});

describe('userSignaturesFromBase64', () => {
  describe('Category B: User signatures decoding tests', () => {
    // Create JSON-formatted signatures for testing (since protobuf parsing returns empty)
    const jsonSignatures = [
      { userId: 'superadmin1@bank.com', signature: 'sig1base64==' },
      { userId: 'superadmin2@bank.com', signature: 'sig2base64==' },
    ];
    const jsonSignaturesBase64 = Buffer.from(
      JSON.stringify(jsonSignatures)
    ).toString('base64');

    it('should decode user signatures from base64 successfully', () => {
      const signatures = userSignaturesFromBase64(jsonSignaturesBase64);

      expect(signatures).toBeDefined();
      expect(Array.isArray(signatures)).toBe(true);
      expect(signatures.length).toBeGreaterThan(0);
    });

    it('should extract user IDs from signatures', () => {
      const signatures = userSignaturesFromBase64(jsonSignaturesBase64);

      const userIds = signatures.map((s) => s.userId);
      expect(userIds).toContain('superadmin1@bank.com');
      expect(userIds).toContain('superadmin2@bank.com');
    });

    it('should extract signature bytes from signatures', () => {
      const signatures = userSignaturesFromBase64(jsonSignaturesBase64);

      for (const sig of signatures) {
        expect(sig.signature).toBeDefined();
        expect(typeof sig.signature).toBe('string');
        expect(sig.signature!.length).toBeGreaterThan(0);
      }
    });

    it('should have expected signature count of 2', () => {
      const signatures = userSignaturesFromBase64(jsonSignaturesBase64);

      expect(signatures).toHaveLength(2);
    });

    it('should return empty array for invalid base64 signatures', () => {
      // The function returns empty array on parse failure, doesn't throw
      const signatures = userSignaturesFromBase64('not-valid-base64!!!');

      expect(signatures).toEqual([]);
    });
  });

  describe('Additional edge cases', () => {
    it('should return empty array for empty string', () => {
      const signatures = userSignaturesFromBase64('');

      expect(signatures).toEqual([]);
    });

    it('should handle object with signatures field', () => {
      // The function handles both array and object with 'signatures' field
      const objWithSignatures = {
        signatures: [
          { userId: 'user1', signature: 'sig1' },
          { userId: 'user2', signature: 'sig2' },
        ],
      };
      const base64 = Buffer.from(JSON.stringify(objWithSignatures)).toString('base64');

      const signatures = userSignaturesFromBase64(base64);

      expect(signatures).toHaveLength(2);
      expect(signatures[0].userId).toBe('user1');
      expect(signatures[1].userId).toBe('user2');
    });

    it('should decode protobuf-encoded signatures', () => {
      // The fixture rulesSignatures is in protobuf format — now correctly decoded
      const signatures = userSignaturesFromBase64(rulesSignaturesBase64);

      // Protobuf format is now decoded correctly (matches Go SDK behavior)
      expect(signatures.length).toBeGreaterThan(0);
      // Each signature should have userId and base64-encoded signature
      for (const sig of signatures) {
        expect(sig.userId).toBeTruthy();
        expect(sig.signature).toBeTruthy();
      }
    });

    it('should handle snake_case user_id field', () => {
      const snakeCaseSignatures = [
        { user_id: 'user1', signature: 'sig1' },
        { user_id: 'user2', signature: 'sig2' },
      ];
      const base64 = Buffer.from(JSON.stringify(snakeCaseSignatures)).toString('base64');

      const signatures = userSignaturesFromBase64(base64);

      expect(signatures).toHaveLength(2);
      expect(signatures[0].userId).toBe('user1');
      expect(signatures[1].userId).toBe('user2');
    });
  });
});
