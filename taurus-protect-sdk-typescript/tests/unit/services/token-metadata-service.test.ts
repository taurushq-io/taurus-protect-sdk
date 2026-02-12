/**
 * Unit tests for TokenMetadataService.
 */

import { TokenMetadataService } from '../../../src/services/token-metadata-service';
import { ValidationError } from '../../../src/errors';
import type { TokenMetadataApi } from '../../../src/internal/openapi/apis/TokenMetadataApi';

function createMockApi(): jest.Mocked<TokenMetadataApi> {
  return {
    tokenMetadataServiceGetERCTokenMetadata: jest.fn(),
    tokenMetadataServiceGetEVMERCTokenMetadata: jest.fn(),
    tokenMetadataServiceGetFATokenMetadata: jest.fn(),
    tokenMetadataServiceGetCryptoPunksTokenMetadata: jest.fn(),
  } as unknown as jest.Mocked<TokenMetadataApi>;
}

describe('TokenMetadataService', () => {
  let mockApi: jest.Mocked<TokenMetadataApi>;
  let service: TokenMetadataService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new TokenMetadataService(mockApi);
  });

  describe('getERCTokenMetadata', () => {
    it('should throw ValidationError when network is empty', async () => {
      await expect(
        service.getERCTokenMetadata({ network: '', contract: '0x123' })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.getERCTokenMetadata({ network: '', contract: '0x123' })
      ).rejects.toThrow('network is required');
    });

    it('should throw ValidationError when contract is empty', async () => {
      await expect(
        service.getERCTokenMetadata({ network: 'mainnet', contract: '' })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.getERCTokenMetadata({ network: 'mainnet', contract: '' })
      ).rejects.toThrow('contract is required');
    });

    it('should return token metadata', async () => {
      mockApi.tokenMetadataServiceGetERCTokenMetadata.mockResolvedValue({
        result: {
          name: 'USD Coin',
          decimals: '6',
          uri: 'https://example.com',
        },
      } as never);

      const metadata = await service.getERCTokenMetadata({
        network: 'mainnet',
        contract: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48',
        blockchain: 'ETH',
      });

      expect(metadata).toBeDefined();
      expect(metadata.name).toBe('USD Coin');
      expect(metadata.decimals).toBe('6');
    });

    it('should handle missing result', async () => {
      mockApi.tokenMetadataServiceGetERCTokenMetadata.mockResolvedValue({
        result: undefined,
      } as never);

      const metadata = await service.getERCTokenMetadata({
        network: 'mainnet',
        contract: '0x123',
      });

      expect(metadata).toBeDefined();
    });
  });

  describe('getEVMERCTokenMetadata', () => {
    it('should throw ValidationError when network is empty', async () => {
      await expect(
        service.getEVMERCTokenMetadata({ network: '', contract: '0x123', blockchain: 'ETH' })
      ).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when contract is empty', async () => {
      await expect(
        service.getEVMERCTokenMetadata({ network: 'mainnet', contract: '', blockchain: 'ETH' })
      ).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when blockchain is empty', async () => {
      await expect(
        service.getEVMERCTokenMetadata({ network: 'mainnet', contract: '0x123', blockchain: '' })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.getEVMERCTokenMetadata({ network: 'mainnet', contract: '0x123', blockchain: '' })
      ).rejects.toThrow('blockchain is required');
    });

    it('should return token metadata', async () => {
      mockApi.tokenMetadataServiceGetEVMERCTokenMetadata.mockResolvedValue({
        result: {
          name: 'Bored Ape',
          description: 'An NFT',
          uri: 'ipfs://...',
        },
      } as never);

      const metadata = await service.getEVMERCTokenMetadata({
        network: 'mainnet',
        contract: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D',
        tokenId: '1234',
        withData: true,
        blockchain: 'ETH',
      });

      expect(metadata).toBeDefined();
      expect(metadata.name).toBe('Bored Ape');
    });

    it('should pass correct parameters to API', async () => {
      mockApi.tokenMetadataServiceGetEVMERCTokenMetadata.mockResolvedValue({
        result: {},
      } as never);

      await service.getEVMERCTokenMetadata({
        network: 'mainnet',
        contract: '0x123',
        tokenId: '42',
        withData: true,
        blockchain: 'MATIC',
      });

      expect(mockApi.tokenMetadataServiceGetEVMERCTokenMetadata).toHaveBeenCalledWith({
        network: 'mainnet',
        contract: '0x123',
        token: '42',
        withData: true,
        blockchain: 'MATIC',
      });
    });
  });

  describe('getFATokenMetadata', () => {
    it('should throw ValidationError when network is empty', async () => {
      await expect(
        service.getFATokenMetadata({ network: '', contract: 'KT1...' })
      ).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when contract is empty', async () => {
      await expect(
        service.getFATokenMetadata({ network: 'mainnet', contract: '' })
      ).rejects.toThrow(ValidationError);
    });

    it('should return FA token metadata', async () => {
      mockApi.tokenMetadataServiceGetFATokenMetadata.mockResolvedValue({
        result: {
          name: 'tzBTC',
          symbol: 'tzBTC',
          decimals: '8',
        },
      } as never);

      const metadata = await service.getFATokenMetadata({
        network: 'mainnet',
        contract: 'KT1PWx2mnDueood7fEmfbBDKx1D9BAnnXitn',
        tokenId: '0',
      });

      expect(metadata).toBeDefined();
      expect(metadata.name).toBe('tzBTC');
      expect(metadata.symbol).toBe('tzBTC');
    });
  });

  describe('getCryptoPunkMetadata', () => {
    it('should throw ValidationError when network is empty', async () => {
      await expect(
        service.getCryptoPunkMetadata({ network: '', contract: '0x...', punkId: '1' })
      ).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when contract is empty', async () => {
      await expect(
        service.getCryptoPunkMetadata({ network: 'mainnet', contract: '', punkId: '1' })
      ).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when punkId is empty', async () => {
      await expect(
        service.getCryptoPunkMetadata({ network: 'mainnet', contract: '0x...', punkId: '' })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.getCryptoPunkMetadata({ network: 'mainnet', contract: '0x...', punkId: '' })
      ).rejects.toThrow('punkId is required');
    });

    it('should return CryptoPunk metadata', async () => {
      mockApi.tokenMetadataServiceGetCryptoPunksTokenMetadata.mockResolvedValue({
        result: {
          punkId: '7804',
          punkAttributes: 'Alien, Small Shades, Cap Forward, Pipe',
          image: 'data:image/svg+xml;base64,...',
        },
      } as never);

      const metadata = await service.getCryptoPunkMetadata({
        network: 'mainnet',
        contract: '0xb47e3cd837dDF8e4c57F05d70Ab865de6e193BBB',
        punkId: '7804',
      });

      expect(metadata).toBeDefined();
      expect(metadata.punkId).toBe('7804');
      expect(metadata.punkAttributes).toBe('Alien, Small Shades, Cap Forward, Pipe');
    });

    it('should pass correct parameters to API', async () => {
      mockApi.tokenMetadataServiceGetCryptoPunksTokenMetadata.mockResolvedValue({
        result: {},
      } as never);

      await service.getCryptoPunkMetadata({
        network: 'mainnet',
        contract: '0xb47e3cd8',
        punkId: '100',
        blockchain: 'ETH',
      });

      expect(mockApi.tokenMetadataServiceGetCryptoPunksTokenMetadata).toHaveBeenCalledWith({
        network: 'mainnet',
        contract: '0xb47e3cd8',
        token: '100',
        blockchain: 'ETH',
      });
    });
  });
});
