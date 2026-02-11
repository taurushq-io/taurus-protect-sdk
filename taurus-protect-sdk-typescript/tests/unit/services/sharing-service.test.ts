/**
 * Unit tests for SharingService (Taurus Network).
 *
 * Tests all public methods for shared addresses and assets.
 */

import { ValidationError } from "../../../src/errors";
import type { TaurusNetworkSharedAddressAssetApi } from "../../../src/internal/openapi/apis/TaurusNetworkSharedAddressAssetApi";
import { SharingService } from "../../../src/services/taurus-network/sharing-service";

function createMockApi(): jest.Mocked<TaurusNetworkSharedAddressAssetApi> {
  return {
    taurusNetworkServiceGetSharedAddresses: jest.fn(),
    taurusNetworkServiceShareAddress: jest.fn(),
    taurusNetworkServiceUnshareAddress: jest.fn(),
    taurusNetworkServiceGetSharedAssets: jest.fn(),
    taurusNetworkServiceShareWhitelistedAsset: jest.fn(),
    taurusNetworkServiceUnshareWhitelistedAsset: jest.fn(),
  } as unknown as jest.Mocked<TaurusNetworkSharedAddressAssetApi>;
}

describe("SharingService", () => {
  let mockApi: jest.Mocked<TaurusNetworkSharedAddressAssetApi>;
  let service: SharingService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new SharingService(mockApi);
  });

  // ===========================================================================
  // listSharedAddresses
  // ===========================================================================

  describe("listSharedAddresses", () => {
    it("should return shared addresses from API", async () => {
      mockApi.taurusNetworkServiceGetSharedAddresses.mockResolvedValue({
        sharedAddresses: [
          {
            id: "sa-1",
            ownerParticipantId: "part-owner",
            targetParticipantId: "part-target",
            blockchain: "ETH",
            network: "mainnet",
            address: "0xabc123",
            status: "ACCEPTED",
          },
        ],
        cursor: {
          currentPage: "1",
          hasNext: true,
          hasPrevious: false,
        },
      });

      const result = await service.listSharedAddresses({
        ownerParticipantId: "part-owner",
      });

      expect(result.sharedAddresses).toHaveLength(1);
      expect(result.sharedAddresses[0].id).toBe("sa-1");
      expect(result.sharedAddresses[0].blockchain).toBe("ETH");
      expect(result.sharedAddresses[0].address).toBe("0xabc123");
      expect(result.sharedAddresses[0].ownerParticipantId).toBe("part-owner");
      expect(result.pagination).toBeDefined();
      expect(result.pagination?.hasNext).toBe(true);
      expect(result.pagination?.hasPrevious).toBe(false);
    });

    it("should return empty list when no shared addresses", async () => {
      mockApi.taurusNetworkServiceGetSharedAddresses.mockResolvedValue({
        sharedAddresses: [],
      });

      const result = await service.listSharedAddresses();

      expect(result.sharedAddresses).toHaveLength(0);
      expect(result.pagination).toBeUndefined();
    });

    it("should pass filter options to API", async () => {
      mockApi.taurusNetworkServiceGetSharedAddresses.mockResolvedValue({
        sharedAddresses: [],
      });

      await service.listSharedAddresses({
        participantId: "part-1",
        blockchain: "ETH",
        network: "mainnet",
        statuses: ["ACCEPTED"],
        pageSize: 25,
        currentPage: "page-token",
        pageRequest: "NEXT",
      });

      expect(mockApi.taurusNetworkServiceGetSharedAddresses).toHaveBeenCalledWith({
        participantID: "part-1",
        ownerParticipantID: undefined,
        targetParticipantID: undefined,
        blockchain: "ETH",
        network: "mainnet",
        ids: undefined,
        statuses: ["ACCEPTED"],
        sortOrder: undefined,
        cursorCurrentPage: "page-token",
        cursorPageRequest: "NEXT",
        cursorPageSize: "25",
      });
    });

    it("should handle undefined response fields gracefully", async () => {
      mockApi.taurusNetworkServiceGetSharedAddresses.mockResolvedValue({
        sharedAddresses: [
          { id: "sa-1" },
        ],
      });

      const result = await service.listSharedAddresses();

      expect(result.sharedAddresses).toHaveLength(1);
      expect(result.sharedAddresses[0].id).toBe("sa-1");
      expect(result.sharedAddresses[0].blockchain).toBeUndefined();
      expect(result.sharedAddresses[0].address).toBeUndefined();
    });

    it("should map proof of ownership fields", async () => {
      mockApi.taurusNetworkServiceGetSharedAddresses.mockResolvedValue({
        sharedAddresses: [
          {
            id: "sa-poo",
            proofOfOwnership: {
              signedPayload: {
                payload: {
                  ownerParticipantID: "owner-1",
                  targetParticipantID: "target-1",
                  address: "0xpoo",
                  blockchain: "ETH",
                  network: "mainnet",
                },
                ownerParticipantSignature: "sig-data",
              },
              signedPayloadHash: "hash123",
              proofOfReserve: {
                curve: "Secp256r1",
                cipher: "ECDSA_SHA256",
                path: "m/44/60/0/0",
                address: "0xpoo",
                publicKey: "pk123",
                challenge: "ch",
                challengeResponse: "cr",
                type: "ADA",
                stakePublicKey: "spk",
                stakeChallengeResponse: "scr",
              },
              signedPayloadAsString: '{"data":"test"}',
            },
          },
        ],
      });

      const result = await service.listSharedAddresses();
      const addr = result.sharedAddresses[0];

      expect(addr.proofOfOwnership).toBeDefined();
      expect(addr.proofOfOwnership?.signedPayload?.payload?.ownerParticipantId).toBe("owner-1");
      expect(addr.proofOfOwnership?.signedPayload?.ownerParticipantSignature).toBe("sig-data");
      expect(addr.proofOfOwnership?.signedPayloadHash).toBe("hash123");
      expect(addr.proofOfOwnership?.proofOfReserve?.curve).toBe("Secp256r1");
      expect(addr.proofOfOwnership?.proofOfReserve?.stakePublicKey).toBe("spk");
      expect(addr.proofOfOwnership?.signedPayloadAsString).toBe('{"data":"test"}');
    });

    it("should map trail entries", async () => {
      const trailDate = new Date("2024-06-15T10:00:00Z");
      mockApi.taurusNetworkServiceGetSharedAddresses.mockResolvedValue({
        sharedAddresses: [
          {
            id: "sa-trails",
            trails: [
              {
                id: "trail-1",
                sharedAddressID: "sa-trails",
                addressStatus: "PENDING",
                comment: "Created",
                createdAt: trailDate,
              },
            ],
          },
        ],
      });

      const result = await service.listSharedAddresses();
      const addr = result.sharedAddresses[0];

      expect(addr.trails).toHaveLength(1);
      expect(addr.trails?.[0].id).toBe("trail-1");
      expect(addr.trails?.[0].sharedAddressId).toBe("sa-trails");
      expect(addr.trails?.[0].addressStatus).toBe("PENDING");
      expect(addr.trails?.[0].createdAt).toEqual(trailDate);
    });
  });

  // ===========================================================================
  // shareAddress
  // ===========================================================================

  describe("shareAddress", () => {
    it("should call API with correct parameters", async () => {
      mockApi.taurusNetworkServiceShareAddress.mockResolvedValue({});

      await service.shareAddress({
        addressId: "addr-123",
        toParticipantId: "part-456",
      });

      expect(mockApi.taurusNetworkServiceShareAddress).toHaveBeenCalledWith({
        body: {
          addressID: "addr-123",
          toParticipantID: "part-456",
          keyValueAttributes: undefined,
        },
      });
    });

    it("should pass key-value attributes", async () => {
      mockApi.taurusNetworkServiceShareAddress.mockResolvedValue({});

      await service.shareAddress({
        addressId: "addr-123",
        toParticipantId: "part-456",
        keyValueAttributes: [
          { key: "label", value: "My Address" },
          { key: "memo", value: "test" },
        ],
      });

      expect(mockApi.taurusNetworkServiceShareAddress).toHaveBeenCalledWith({
        body: {
          addressID: "addr-123",
          toParticipantID: "part-456",
          keyValueAttributes: [
            { key: "label", value: "My Address" },
            { key: "memo", value: "test" },
          ],
        },
      });
    });

    it("should throw ValidationError when addressId is empty", async () => {
      await expect(
        service.shareAddress({ addressId: "", toParticipantId: "part-456" })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.shareAddress({ addressId: "", toParticipantId: "part-456" })
      ).rejects.toThrow("addressId is required");
    });

    it("should throw ValidationError when addressId is whitespace only", async () => {
      await expect(
        service.shareAddress({ addressId: "   ", toParticipantId: "part-456" })
      ).rejects.toThrow(ValidationError);
    });

    it("should throw ValidationError when toParticipantId is empty", async () => {
      await expect(
        service.shareAddress({ addressId: "addr-123", toParticipantId: "" })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.shareAddress({ addressId: "addr-123", toParticipantId: "" })
      ).rejects.toThrow("toParticipantId is required");
    });
  });

  // ===========================================================================
  // unshareAddress
  // ===========================================================================

  describe("unshareAddress", () => {
    it("should call API with correct parameters", async () => {
      mockApi.taurusNetworkServiceUnshareAddress.mockResolvedValue({});

      await service.unshareAddress("sa-789");

      expect(mockApi.taurusNetworkServiceUnshareAddress).toHaveBeenCalledWith({
        tnSharedAddressID: "sa-789",
        body: {},
      });
    });

    it("should throw ValidationError when sharedAddressId is empty", async () => {
      await expect(service.unshareAddress("")).rejects.toThrow(ValidationError);
      await expect(service.unshareAddress("")).rejects.toThrow("sharedAddressId is required");
    });

    it("should throw ValidationError when sharedAddressId is whitespace only", async () => {
      await expect(service.unshareAddress("   ")).rejects.toThrow(ValidationError);
    });
  });

  // ===========================================================================
  // listSharedAssets
  // ===========================================================================

  describe("listSharedAssets", () => {
    it("should return shared assets from API", async () => {
      mockApi.taurusNetworkServiceGetSharedAssets.mockResolvedValue({
        sharedAssets: [
          {
            id: "asset-1",
            ownerParticipantId: "part-owner",
            targetParticipantId: "part-target",
            blockchain: "ETH",
            network: "mainnet",
            name: "USDC",
            symbol: "USDC",
            decimals: "6",
            contractAddress: "0xA0b8...",
            kind: "ERC20",
            status: "ACCEPTED",
          },
        ],
        cursor: {
          currentPage: "2",
          hasNext: false,
          hasPrevious: true,
        },
      });

      const result = await service.listSharedAssets({
        blockchain: "ETH",
      });

      expect(result.sharedAssets).toHaveLength(1);
      expect(result.sharedAssets[0].id).toBe("asset-1");
      expect(result.sharedAssets[0].name).toBe("USDC");
      expect(result.sharedAssets[0].symbol).toBe("USDC");
      expect(result.sharedAssets[0].decimals).toBe("6");
      expect(result.sharedAssets[0].contractAddress).toBe("0xA0b8...");
      expect(result.sharedAssets[0].kind).toBe("ERC20");
      expect(result.pagination?.hasNext).toBe(false);
      expect(result.pagination?.hasPrevious).toBe(true);
    });

    it("should return empty list when no shared assets", async () => {
      mockApi.taurusNetworkServiceGetSharedAssets.mockResolvedValue({
        sharedAssets: [],
      });

      const result = await service.listSharedAssets();

      expect(result.sharedAssets).toHaveLength(0);
    });

    it("should pass filter options to API", async () => {
      mockApi.taurusNetworkServiceGetSharedAssets.mockResolvedValue({
        sharedAssets: [],
      });

      await service.listSharedAssets({
        participantId: "part-1",
        ownerParticipantId: "owner-1",
        targetParticipantId: "target-1",
        blockchain: "ETH",
        network: "mainnet",
        ids: ["id-1", "id-2"],
        statuses: ["PENDING"],
        sortOrder: "ASC",
        pageSize: 10,
      });

      expect(mockApi.taurusNetworkServiceGetSharedAssets).toHaveBeenCalledWith({
        participantID: "part-1",
        ownerParticipantID: "owner-1",
        targetParticipantID: "target-1",
        blockchain: "ETH",
        network: "mainnet",
        ids: ["id-1", "id-2"],
        statuses: ["PENDING"],
        sortOrder: "ASC",
        cursorCurrentPage: undefined,
        cursorPageRequest: undefined,
        cursorPageSize: "10",
      });
    });

    it("should map shared asset trail entries", async () => {
      const trailDate = new Date("2024-07-01T12:00:00Z");
      mockApi.taurusNetworkServiceGetSharedAssets.mockResolvedValue({
        sharedAssets: [
          {
            id: "asset-trails",
            trails: [
              {
                id: "at-1",
                sharedAssetID: "asset-trails",
                assetStatus: "ACCEPTED",
                comment: "Accepted by target",
                createdAt: trailDate,
              },
            ],
          },
        ],
      });

      const result = await service.listSharedAssets();
      const asset = result.sharedAssets[0];

      expect(asset.trails).toHaveLength(1);
      expect(asset.trails?.[0].id).toBe("at-1");
      expect(asset.trails?.[0].sharedAssetId).toBe("asset-trails");
      expect(asset.trails?.[0].assetStatus).toBe("ACCEPTED");
      expect(asset.trails?.[0].createdAt).toEqual(trailDate);
    });

    it("should map date fields correctly", async () => {
      const created = new Date("2024-01-01T00:00:00Z");
      const updated = new Date("2024-06-01T00:00:00Z");
      const rejected = new Date("2024-03-01T00:00:00Z");

      mockApi.taurusNetworkServiceGetSharedAssets.mockResolvedValue({
        sharedAssets: [
          {
            id: "asset-dates",
            originCreationDate: created,
            targetRejectedAt: rejected,
            createdAt: created,
            updatedAt: updated,
          },
        ],
      });

      const result = await service.listSharedAssets();
      const asset = result.sharedAssets[0];

      expect(asset.originCreationDate).toEqual(created);
      expect(asset.targetRejectedAt).toEqual(rejected);
      expect(asset.createdAt).toEqual(created);
      expect(asset.updatedAt).toEqual(updated);
    });
  });

  // ===========================================================================
  // shareWhitelistedAsset
  // ===========================================================================

  describe("shareWhitelistedAsset", () => {
    it("should call API with correct parameters", async () => {
      mockApi.taurusNetworkServiceShareWhitelistedAsset.mockResolvedValue({});

      await service.shareWhitelistedAsset({
        whitelistedContractId: "wc-123",
        toParticipantId: "part-456",
      });

      expect(mockApi.taurusNetworkServiceShareWhitelistedAsset).toHaveBeenCalledWith({
        body: {
          whitelistedContractID: "wc-123",
          toParticipantID: "part-456",
        },
      });
    });

    it("should throw ValidationError when whitelistedContractId is empty", async () => {
      await expect(
        service.shareWhitelistedAsset({ whitelistedContractId: "", toParticipantId: "part-456" })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.shareWhitelistedAsset({ whitelistedContractId: "", toParticipantId: "part-456" })
      ).rejects.toThrow("whitelistedContractId is required");
    });

    it("should throw ValidationError when whitelistedContractId is whitespace only", async () => {
      await expect(
        service.shareWhitelistedAsset({ whitelistedContractId: "  ", toParticipantId: "part-456" })
      ).rejects.toThrow(ValidationError);
    });

    it("should throw ValidationError when toParticipantId is empty", async () => {
      await expect(
        service.shareWhitelistedAsset({ whitelistedContractId: "wc-123", toParticipantId: "" })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.shareWhitelistedAsset({ whitelistedContractId: "wc-123", toParticipantId: "" })
      ).rejects.toThrow("toParticipantId is required");
    });
  });

  // ===========================================================================
  // unshareWhitelistedAsset
  // ===========================================================================

  describe("unshareWhitelistedAsset", () => {
    it("should call API with correct parameters", async () => {
      mockApi.taurusNetworkServiceUnshareWhitelistedAsset.mockResolvedValue({});

      await service.unshareWhitelistedAsset("sa-999");

      expect(mockApi.taurusNetworkServiceUnshareWhitelistedAsset).toHaveBeenCalledWith({
        tnSharedAssetID: "sa-999",
        body: {},
      });
    });

    it("should throw ValidationError when sharedAssetId is empty", async () => {
      await expect(service.unshareWhitelistedAsset("")).rejects.toThrow(ValidationError);
      await expect(service.unshareWhitelistedAsset("")).rejects.toThrow("sharedAssetId is required");
    });

    it("should throw ValidationError when sharedAssetId is whitespace only", async () => {
      await expect(service.unshareWhitelistedAsset("   ")).rejects.toThrow(ValidationError);
    });
  });
});
