/**
 * Unit tests for WhitelistedAssetService.
 *
 * Every read path must verify. These tests use GENUINELY SIGNED fixtures
 * (`buildFullAssetFixture`) rather than bare payloads, because a bare payload
 * cannot distinguish "the service verified and accepted" from "the service
 * never verified at all" — which is precisely the defect these tests missed
 * before: get/list/getEnvelope returned attacker-controllable payload data
 * labelled as verified.
 *
 * Covered here:
 * - get/list/getEnvelope all run the 5-step verification
 * - a tampered payload is rejected on every read path
 * - list excludes an unverifiable row and reports it
 * - list errors when rows were returned but none survived
 * - input validation
 */

import * as crypto from 'crypto';

import { verifySignature } from '../../../src/crypto';
import { WhitelistedAssetService } from '../../../src/services/whitelisted-asset-service';
import { IntegrityError, NotFoundError, ValidationError } from '../../../src/errors';
import type { ContractWhitelistingApi } from '../../../src/internal/openapi';
import type { WhitelistedAssetServiceConfig } from '../../../src/services/whitelisted-asset-service';
import type { DecodedRulesContainer } from '../../../src/models/governance-rules';
import type { RuleUserSignature } from '../../../src/models/governance-rules';
import { createEmptyRulesContainer } from '../../../src/models/governance-rules';
import {
  buildFullAssetFixture,
  fixtureToDto,
} from '../fixtures/whitelisted-asset-fixtures';

const TEST_SUPER_ADMIN_KEY_PEM = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEM2NtzaFhm7xIR3OvWq5chW3/GEvW
L+3uqoE6lEJ13eWbulxsP/5h36VCqYDIGN/0wDeWwLYdpu5HhSXWhxCsCA==
-----END PUBLIC KEY-----`;

const mockConfig: WhitelistedAssetServiceConfig = {
  superAdminKeysPem: [TEST_SUPER_ADMIN_KEY_PEM],
  minValidSignatures: 1,
  rulesContainerDecoder: (_base64: string): DecodedRulesContainer => createEmptyRulesContainer(),
  userSignaturesDecoder: (_base64: string): RuleUserSignature[] => [],
};

function createMockApi(): jest.Mocked<ContractWhitelistingApi> {
  return {
    whitelistServiceGetWhitelistedContract: jest.fn(),
    whitelistServiceGetWhitelistedContracts: jest.fn(),
    whitelistServiceCreateWhitelistedContract: jest.fn(),
    whitelistServiceApproveWhitelistedContract: jest.fn(),
    whitelistServiceRejectWhitelistedContract: jest.fn(),
    whitelistServiceUpdateWhitelistedContract: jest.fn(),
    whitelistServiceGetWhitelistedContractsForApproval: jest.fn(),
    whitelistServiceCreateWhitelistedContractAttributes: jest.fn(),
    whitelistServiceGetWhitelistedContractAttribute: jest.fn(),
    whitelistServiceDeleteWhitelistedContractAttribute: jest.fn(),
  } as unknown as jest.Mocked<ContractWhitelistingApi>;
}

/**
 * Builds a service wired to a genuinely signed fixture, so verification can
 * actually succeed. Returns the fixture too, for tampering.
 */
function signedSetup(): {
  api: jest.Mocked<ContractWhitelistingApi>;
  svc: WhitelistedAssetService;
  fixture: ReturnType<typeof buildFullAssetFixture>;
} {
  const fixture = buildFullAssetFixture();
  const api = createMockApi();
  const svc = new WhitelistedAssetService(api, {
    superAdminKeysPem: [fixture.saPem],
    minValidSignatures: 1,
    rulesContainerDecoder: fixture.rulesContainerDecoder,
    userSignaturesDecoder: fixture.userSignaturesDecoder,
  });
  return { api, svc, fixture };
}

describe('WhitelistedAssetService', () => {
  let mockApi: jest.Mocked<ContractWhitelistingApi>;
  let service: WhitelistedAssetService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new WhitelistedAssetService(mockApi, mockConfig);
  });

  describe('get', () => {
    it('should throw ValidationError when assetId is 0', async () => {
      await expect(service.get(0)).rejects.toThrow(ValidationError);
      await expect(service.get(0)).rejects.toThrow('assetId must be positive');
    });

    it('should throw ValidationError when assetId is negative', async () => {
      await expect(service.get(-1)).rejects.toThrow(ValidationError);
    });

    it('should throw NotFoundError when asset is not found', async () => {
      mockApi.whitelistServiceGetWhitelistedContract.mockResolvedValue({
        result: undefined,
      });

      await expect(service.get(999)).rejects.toThrow(NotFoundError);
    });

    it('should return the asset from a genuinely verified envelope', async () => {
      const { api, svc, fixture } = signedSetup();
      api.whitelistServiceGetWhitelistedContract.mockResolvedValue({
        result: fixtureToDto(fixture.envelope),
      });

      const asset = await svc.get(1);
      expect(asset).toBeDefined();
      expect(asset.blockchain).toBe('ETH');
      expect(asset.network).toBe('mainnet');
      expect(asset.symbol).toBe('USDC');
    });

    // THE REGRESSION TEST for the defect this change fixes. Before the fix,
    // get() never invoked the verifier, so a payload the attacker rewrote came
    // back as "verified" and this assertion failed.
    it('should reject a tampered payload rather than returning it as verified', async () => {
      const { api, svc, fixture } = signedSetup();
      const tamperedPayload = JSON.stringify({
        blockchain: 'ETH',
        network: 'mainnet',
        contractAddress: '0xATTACKER',
        name: 'USDC',
        symbol: 'USDC',
        decimals: 6,
      });
      api.whitelistServiceGetWhitelistedContract.mockResolvedValue({
        result: fixtureToDto(fixture.envelope, {
          payloadAsString: tamperedPayload,
        }),
      });

      await expect(svc.get(1)).rejects.toThrow(IntegrityError);
    });

    it('should throw IntegrityError when payload is missing', async () => {
      const { api, svc, fixture } = signedSetup();
      api.whitelistServiceGetWhitelistedContract.mockResolvedValue({
        result: fixtureToDto(fixture.envelope, { payloadAsString: '' }),
      });

      await expect(svc.get(1)).rejects.toThrow(IntegrityError);
    });

    it('should call API with correct parameters', async () => {
      const { api, svc, fixture } = signedSetup();
      api.whitelistServiceGetWhitelistedContract.mockResolvedValue({
        result: fixtureToDto(fixture.envelope),
      });

      await svc.get(1);
      expect(api.whitelistServiceGetWhitelistedContract).toHaveBeenCalledWith({
        id: '1',
      });
    });
  });

  describe('list', () => {
    it('should return verified rows and pagination', async () => {
      const { api, svc, fixture } = signedSetup();
      api.whitelistServiceGetWhitelistedContracts.mockResolvedValue({
        result: [
          fixtureToDto(fixture.envelope, { id: '1' }),
          fixtureToDto(fixture.envelope, { id: '2' }),
        ],
        totalItems: '10',
      });

      const result = await svc.list({ limit: 50 });
      expect(result.items).toHaveLength(2);
      expect(result.items[0].symbol).toBe('USDC');
      expect(result.pagination?.totalItems).toBe(10);
      expect(result.excludedUnverified).toHaveLength(0);
    });

    it('should exclude an unverifiable row and report it, keeping the rest', async () => {
      const { api, svc, fixture } = signedSetup();
      api.whitelistServiceGetWhitelistedContracts.mockResolvedValue({
        result: [
          fixtureToDto(fixture.envelope, { id: '1' }),
          fixtureToDto(fixture.envelope, {
            id: '2',
            payloadAsString: JSON.stringify({ blockchain: 'ETH', name: 'EVIL' }),
          }),
        ],
        totalItems: '2',
      });

      const result = await svc.list({ limit: 50 });
      expect(result.items).toHaveLength(1);
      expect(result.excludedUnverified).toHaveLength(1);
      expect(result.excludedUnverified[0].id).toBe(2);
      expect(result.excludedUnverified[0].reason).toBeTruthy();
    });

    // A short page and a filtered page must not look the same. Rows came back
    // and none survived: that is a systemic failure, not an empty whitelist.
    it('should error when rows were returned but none verified', async () => {
      const { api, svc, fixture } = signedSetup();
      api.whitelistServiceGetWhitelistedContracts.mockResolvedValue({
        result: [
          fixtureToDto(fixture.envelope, {
            id: '1',
            payloadAsString: JSON.stringify({ blockchain: 'ETH', name: 'EVIL' }),
          }),
        ],
        totalItems: '1',
      });

      await expect(svc.list({ limit: 50 })).rejects.toThrow(IntegrityError);
      await expect(svc.list({ limit: 50 })).rejects.toThrow(
        /all 1 whitelisted asset\(s\) failed verification/
      );
    });

    it('should throw ValidationError when limit is 0', async () => {
      await expect(service.list({ limit: 0 })).rejects.toThrow(ValidationError);
      await expect(service.list({ limit: 0 })).rejects.toThrow('limit must be positive');
    });

    it('should throw ValidationError when limit is negative', async () => {
      await expect(service.list({ limit: -1 })).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when offset is negative', async () => {
      await expect(service.list({ offset: -1 })).rejects.toThrow(ValidationError);
      await expect(service.list({ offset: -1 })).rejects.toThrow('offset cannot be negative');
    });

    it('should handle empty results', async () => {
      mockApi.whitelistServiceGetWhitelistedContracts.mockResolvedValue({
        result: [],
        totalItems: '0',
      });

      const result = await service.list();
      expect(result.items).toHaveLength(0);
    });

    it('should pass filter options to API', async () => {
      mockApi.whitelistServiceGetWhitelistedContracts.mockResolvedValue({
        result: [],
        totalItems: '0',
      });

      await service.list({
        blockchain: 'ETH',
        network: 'mainnet',
        query: 'USDC',
        kindTypes: ['token'],
      });

      expect(mockApi.whitelistServiceGetWhitelistedContracts).toHaveBeenCalledWith(
        expect.objectContaining({
          blockchain: 'ETH',
          network: 'mainnet',
          query: 'USDC',
          kindTypes: ['token'],
        })
      );
    });

    it('should pass ids to API', async () => {
      mockApi.whitelistServiceGetWhitelistedContracts.mockResolvedValue({
        result: [],
        totalItems: '0',
      });

      await service.list({ ids: ['7', '9'] });

      expect(mockApi.whitelistServiceGetWhitelistedContracts).toHaveBeenCalledWith(
        expect.objectContaining({ whitelistedContractAddressIds: ['7', '9'] })
      );
    });
  });

  // The for-approval endpoint had no verified reader, so approvers inspected rows
  // nothing had checked against governance.
  describe('listForApproval', () => {
    it('should return verified rows and pagination', async () => {
      const { api, svc, fixture } = signedSetup();
      api.whitelistServiceGetWhitelistedContractsForApproval.mockResolvedValue({
        result: [fixtureToDto(fixture.envelope, { id: '1' })],
        totalItems: '1',
      });

      const result = await svc.listForApproval({ limit: 50 });
      expect(result.items).toHaveLength(1);
      expect(result.items[0].symbol).toBe('USDC');
      expect(result.pagination?.totalItems).toBe(1);
      expect(result.excludedUnverified).toHaveLength(0);
    });

    it('should reject a tampered row rather than presenting it for approval', async () => {
      const { api, svc, fixture } = signedSetup();
      api.whitelistServiceGetWhitelistedContractsForApproval.mockResolvedValue({
        result: [
          fixtureToDto(fixture.envelope, { id: '1' }),
          fixtureToDto(fixture.envelope, {
            id: '2',
            payloadAsString: JSON.stringify({
              blockchain: 'ETH',
              contractAddress: '0xATTACKER',
            }),
          }),
        ],
        totalItems: '2',
      });

      const result = await svc.listForApproval();
      expect(result.items).toHaveLength(1);
      expect(result.items.map((a) => a.contractAddress)).not.toContain('0xATTACKER');
      expect(result.excludedUnverified).toHaveLength(1);
      expect(result.excludedUnverified[0].id).toBe(2);
    });

    it('should error when rows were returned but none verified', async () => {
      const { api, svc, fixture } = signedSetup();
      api.whitelistServiceGetWhitelistedContractsForApproval.mockResolvedValue({
        result: [
          fixtureToDto(fixture.envelope, {
            id: '1',
            payloadAsString: JSON.stringify({ blockchain: 'ETH', name: 'EVIL' }),
          }),
        ],
        totalItems: '1',
      });

      await expect(svc.listForApproval()).rejects.toThrow(IntegrityError);
    });

    it('should throw ValidationError when limit is 0', async () => {
      await expect(service.listForApproval({ limit: 0 })).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when offset is negative', async () => {
      await expect(service.listForApproval({ offset: -1 })).rejects.toThrow(ValidationError);
    });

    it('should pass limit, offset and ids to API', async () => {
      mockApi.whitelistServiceGetWhitelistedContractsForApproval.mockResolvedValue({
        result: [],
        totalItems: '0',
      });

      await service.listForApproval({ limit: 25, offset: 50, ids: ['3'] });

      expect(
        mockApi.whitelistServiceGetWhitelistedContractsForApproval
      ).toHaveBeenCalledWith({ limit: '25', offset: '50', ids: ['3'] });
    });

    it('should handle empty results', async () => {
      mockApi.whitelistServiceGetWhitelistedContractsForApproval.mockResolvedValue({
        result: [],
        totalItems: '0',
      });

      const result = await service.listForApproval();
      expect(result.items).toHaveLength(0);
      expect(result.excludedUnverified).toHaveLength(0);
    });
  });

  describe('getEnvelope', () => {
    it('should throw ValidationError when assetId is 0', async () => {
      await expect(service.getEnvelope(0)).rejects.toThrow(ValidationError);
    });

    it('should throw NotFoundError when asset is not found', async () => {
      mockApi.whitelistServiceGetWhitelistedContract.mockResolvedValue({
        result: undefined,
      });

      await expect(service.getEnvelope(999)).rejects.toThrow(NotFoundError);
    });

    it('should return the envelope once verified', async () => {
      const { api, svc, fixture } = signedSetup();
      api.whitelistServiceGetWhitelistedContract.mockResolvedValue({
        result: fixtureToDto(fixture.envelope),
      });

      const envelope = await svc.getEnvelope(1);
      expect(envelope).toBeDefined();
      expect(envelope.id).toBe(1);
      expect(envelope.metadata.hash).toBe(fixture.envelope.metadata.hash);
      expect(envelope.signedContractAddress.signatures).toHaveLength(1);
    });

    // getEnvelope used to hand back the raw envelope with nothing checked,
    // which made it a way around the verified read paths.
    it('should reject an unverifiable envelope instead of returning it raw', async () => {
      const { api, svc, fixture } = signedSetup();
      api.whitelistServiceGetWhitelistedContract.mockResolvedValue({
        result: fixtureToDto(fixture.envelope, {
          payloadAsString: JSON.stringify({ blockchain: 'ETH', name: 'EVIL' }),
        }),
      });

      await expect(svc.getEnvelope(1)).rejects.toThrow(IntegrityError);
    });
  });

  describe('getWithVerification', () => {
    it('should throw ValidationError when assetId is 0', async () => {
      await expect(service.getWithVerification(0)).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when assetId is negative', async () => {
      await expect(service.getWithVerification(-1)).rejects.toThrow(ValidationError);
    });

    it('should throw NotFoundError when asset is not found', async () => {
      mockApi.whitelistServiceGetWhitelistedContract.mockResolvedValue({
        result: undefined,
      });

      await expect(service.getWithVerification(999)).rejects.toThrow(NotFoundError);
    });
  });

  describe('withVerification factory', () => {
    it('should create a service with verification enabled', () => {
      const svc = WhitelistedAssetService.withVerification(mockApi, mockConfig);
      expect(svc).toBeInstanceOf(WhitelistedAssetService);
    });
  });

  // One signature covers every hash in the batch, so approval is all-or-nothing.
  // Signing only the rows that verified would tell the approver they approved less than
  // they did, and the API takes a single signature so there is no partial submission.
  describe('approve', () => {
    const approverKey = crypto.generateKeyPairSync('ec', {
      namedCurve: 'prime256v1',
    }).privateKey;

    it('signs nothing when one row of the batch fails verification', async () => {
      const { api, svc, fixture } = signedSetup();
      const tampered = JSON.stringify({
        blockchain: 'ETH',
        network: 'mainnet',
        contractAddress: '0xATTACKER',
        name: 'USDC',
        symbol: 'USDC',
        decimals: 6,
      });

      // Row 1 verifies, row 2 does not: the mixed batch is the case that separates
      // all-or-nothing from sign-the-survivors. The batch is re-read through the
      // verifying LIST path filtered by ids, so both rows arrive on one page.
      api.whitelistServiceGetWhitelistedContracts.mockResolvedValue({
        result: [
          fixtureToDto(fixture.envelope, { id: '1' }),
          fixtureToDto(fixture.envelope, { id: '2', payloadAsString: tampered }),
        ],
        totalItems: '2',
      });

      await expect(svc.approve([1, 2], approverKey, 'batch approval')).rejects.toThrow(
        /refusing to sign/
      );
      expect(api.whitelistServiceApproveWhitelistedContract).not.toHaveBeenCalled();
    });

    it('signs the verified hashes with the ids sorted numerically', async () => {
      const { api, svc, fixture } = signedSetup();
      // ONE id-filtered page through the verifying list path, not one GET per id.
      api.whitelistServiceGetWhitelistedContracts.mockResolvedValue({
        result: [
          fixtureToDto(fixture.envelope, { id: '3' }),
          fixtureToDto(fixture.envelope, { id: '7' }),
        ],
        totalItems: '2',
      });

      await svc.approve([7, 3], approverKey, 'batch approval');

      expect(api.whitelistServiceGetWhitelistedContracts).toHaveBeenCalledTimes(1);
      const listArgs = api.whitelistServiceGetWhitelistedContracts.mock
        .calls[0]?.[0] as { whitelistedContractAddressIds?: string[]; includeForApproval?: boolean };
      expect(listArgs.whitelistedContractAddressIds).toEqual(['3', '7']);
      expect(listArgs.includeForApproval).toBe(true);

      const body = api.whitelistServiceApproveWhitelistedContract.mock.calls[0]?.[0]
        ?.body as { ids: string[]; signature: string; comment: string };
      expect(body.ids).toEqual(['3', '7']);
      expect(body.comment).toBe('batch approval');

      // ECDSA is randomized, so verify rather than byte-compare. What is being pinned is
      // WHICH bytes were signed: the hashes the verifier produced, sorted with the ids.
      const hash = fixture.envelope.metadata.hash;
      const approverPub = crypto.createPublicKey(approverKey);
      expect(
        verifySignature(
          approverPub,
          Buffer.from(JSON.stringify([hash, hash]), 'utf-8'),
          body.signature
        )
      ).toBe(true);
      expect(
        verifySignature(
          approverPub,
          Buffer.from(JSON.stringify(['0xATTACKER']), 'utf-8'),
          body.signature
        )
      ).toBe(false);
    });

    it('rejects bad input without reaching the API', async () => {
      const { api, svc } = signedSetup();

      await expect(svc.approve([], approverKey, 'c')).rejects.toThrow('ids cannot be empty');
      await expect(
        svc.approve([1], undefined as unknown as crypto.KeyObject, 'c')
      ).rejects.toThrow('privateKey is required');
      await expect(svc.approve([1], approverKey, '')).rejects.toThrow('comment is required');
      await expect(svc.approve([0], approverKey, 'c')).rejects.toThrow(
        'must be a positive integer'
      );

      expect(api.whitelistServiceApproveWhitelistedContract).not.toHaveBeenCalled();
    });
  });
});
