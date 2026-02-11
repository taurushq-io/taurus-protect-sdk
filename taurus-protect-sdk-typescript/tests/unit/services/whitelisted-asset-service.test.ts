/**
 * Unit tests for WhitelistedAssetService.
 *
 * Tests focus on the security-critical behavior:
 * - Verified payload data must be used for security-critical fields
 * - Missing payload must throw IntegrityError
 * - Validation of input parameters
 * - 5-step verification flow delegation
 */

import { WhitelistedAssetService } from '../../../src/services/whitelisted-asset-service';
import { IntegrityError, NotFoundError, ValidationError } from '../../../src/errors';
import type { ContractWhitelistingApi } from '../../../src/internal/openapi';
import type { WhitelistedAssetServiceConfig } from '../../../src/services/whitelisted-asset-service';
import type { DecodedRulesContainer } from '../../../src/models/governance-rules';
import type { RuleUserSignature } from '../../../src/models/governance-rules';
import { createEmptyRulesContainer } from '../../../src/models/governance-rules';

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

    it('should return asset from verified payload', async () => {
      mockApi.whitelistServiceGetWhitelistedContract.mockResolvedValue({
        result: {
          id: '123',
          metadata: {
            payloadAsString: JSON.stringify({
              blockchain: 'ETH',
              network: 'mainnet',
              address: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48',
              name: 'USD Coin',
              symbol: 'USDC',
              decimals: 6,
            }),
          },
          blockchain: 'UNVERIFIED',
          network: 'UNVERIFIED',
        },
      });

      const asset = await service.get(123);
      expect(asset).toBeDefined();
      expect(asset.id).toBe(123);
      // Security: fields must come from verified payload, not envelope
      expect(asset.blockchain).toBe('ETH');
      expect(asset.network).toBe('mainnet');
      expect(asset.name).toBe('USD Coin');
      expect(asset.symbol).toBe('USDC');
    });

    it('should throw IntegrityError when payload is missing', async () => {
      mockApi.whitelistServiceGetWhitelistedContract.mockResolvedValue({
        result: {
          id: '123',
          metadata: {},
          blockchain: 'ETH',
          network: 'mainnet',
        },
      });

      await expect(service.get(123)).rejects.toThrow(IntegrityError);
      await expect(service.get(123)).rejects.toThrow('payloadAsString is missing');
    });

    it('should call API with correct parameters', async () => {
      mockApi.whitelistServiceGetWhitelistedContract.mockResolvedValue({
        result: {
          id: '456',
          metadata: {
            payloadAsString: JSON.stringify({
              blockchain: 'BTC',
              network: 'mainnet',
              address: 'bc1q...',
              name: 'Bitcoin',
              symbol: 'BTC',
            }),
          },
        },
      });

      await service.get(456);
      expect(mockApi.whitelistServiceGetWhitelistedContract).toHaveBeenCalledWith({
        id: '456',
      });
    });
  });

  describe('list', () => {
    it('should return paginated results', async () => {
      mockApi.whitelistServiceGetWhitelistedContracts.mockResolvedValue({
        result: [
          {
            id: '1',
            metadata: {
              payloadAsString: JSON.stringify({
                blockchain: 'ETH',
                network: 'mainnet',
                name: 'Token A',
                symbol: 'TKA',
              }),
            },
          },
          {
            id: '2',
            metadata: {
              payloadAsString: JSON.stringify({
                blockchain: 'ETH',
                network: 'mainnet',
                name: 'Token B',
                symbol: 'TKB',
              }),
            },
          },
        ],
        totalItems: '10',
      });

      const result = await service.list({ limit: 50 });
      expect(result.items).toHaveLength(2);
      expect(result.items[0].name).toBe('Token A');
      expect(result.items[1].name).toBe('Token B');
      expect(result.pagination?.totalItems).toBe(10);
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

    it('should return envelope with mapped fields', async () => {
      mockApi.whitelistServiceGetWhitelistedContract.mockResolvedValue({
        result: {
          id: '123',
          metadata: {
            hash: 'abc123',
            payloadAsString: '{"name":"Test"}',
          },
          rulesContainer: 'cnVsZXM=',
          rulesSignatures: 'c2lncw==',
          signedContractAddress: {
            payload: 'signed-payload',
            signatures: [
              {
                signature: { userId: 'user1', signature: 'sig1', comment: 'ok' },
                hashes: ['hash1'],
              },
            ],
          },
          blockchain: 'ETH',
          network: 'mainnet',
        },
      });

      const envelope = await service.getEnvelope(123);
      expect(envelope).toBeDefined();
      expect(envelope.id).toBe(123);
      expect(envelope.metadata.hash).toBe('abc123');
      expect(envelope.rulesContainerBase64).toBe('cnVsZXM=');
      expect(envelope.signedContractAddress.signatures).toHaveLength(1);
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
});
