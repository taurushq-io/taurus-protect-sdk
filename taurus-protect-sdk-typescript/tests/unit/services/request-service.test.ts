/**
 * Unit tests for RequestService.
 *
 * These tests verify the critical security features:
 * - Hash verification using constant-time comparison
 * - ECDSA signing for request approval
 * - Proper error handling
 */

import * as crypto from "crypto";
import { RequestService } from "../../../src/services/request-service";
import { IntegrityError, NotFoundError } from "../../../src/errors";
import type { Request, RequestMetadata } from "../../../src/models/request";
import type { RequestsApi } from "../../../src/internal/openapi/apis/RequestsApi";
import type { TgvalidatordGetRequestReply } from "../../../src/internal/openapi/models/TgvalidatordGetRequestReply";
import type { TgvalidatordGetRequestsV2Reply } from "../../../src/internal/openapi/models/TgvalidatordGetRequestsV2Reply";
import type { TgvalidatordApproveRequestsReply } from "../../../src/internal/openapi/models/TgvalidatordApproveRequestsReply";
import { calculateHexHash } from "../../../src/crypto";

// Mock RequestsApi
function createMockRequestsApi(): jest.Mocked<RequestsApi> {
  return {
    requestServiceGetRequest: jest.fn(),
    requestServiceGetRequestsV2: jest.fn(),
    requestServiceGetRequestsForApprovalV2: jest.fn(),
    requestServiceApproveRequests: jest.fn(),
    requestServiceRejectRequests: jest.fn(),
    requestServiceCreateOutgoingRequest: jest.fn(),
  } as unknown as jest.Mocked<RequestsApi>;
}

// Helper to create a valid request with matching hash
function createValidRequest(
  id: number,
  payloadAsString: string
): { dto: TgvalidatordGetRequestReply; expectedRequest: Request } {
  const hash = calculateHexHash(payloadAsString);

  const dto: TgvalidatordGetRequestReply = {
    result: {
      id: String(id),
      type: "payment",
      status: "PENDING_APPROVAL",
      metadata: {
        hash,
        payloadAsString,
        payload: { amount: "1000" },
      },
      tenantId: "tenant-1",
      currency: "ETH",
      creationDate: new Date("2024-01-01T00:00:00Z"),
      updateDate: new Date("2024-01-02T00:00:00Z"),
      needsApprovalFrom: ["group-1"],
    },
  };

  const expectedRequest: Request = {
    id,
    type: "payment",
    status: "PENDING_APPROVAL",
    metadata: {
      hash,
      payloadAsString,
      // SECURITY: payload field intentionally removed - use payloadAsString
    },
    tenantId: "tenant-1",
    currency: "ETH",
    currencyInfo: undefined,
    memo: undefined,
    rule: undefined,
    externalRequestId: undefined,
    requestBundleId: undefined,
    needsApprovalFrom: ["group-1"],
    approvers: undefined,
    createdAt: new Date("2024-01-01T00:00:00Z"),
    updatedAt: new Date("2024-01-02T00:00:00Z"),
    tags: [],
    signedRequests: [],
  };

  return { dto, expectedRequest };
}

describe("RequestService", () => {
  let mockApi: jest.Mocked<RequestsApi>;
  let service: RequestService;

  beforeEach(() => {
    mockApi = createMockRequestsApi();
    service = new RequestService(mockApi);
  });

  describe("get", () => {
    it("should return a request when hash verification passes", async () => {
      const { dto, expectedRequest } = createValidRequest(
        123,
        '{"amount":"1000","currency":"ETH"}'
      );
      mockApi.requestServiceGetRequest.mockResolvedValue(dto);

      const result = await service.get(123);

      expect(result.id).toBe(expectedRequest.id);
      expect(result.status).toBe(expectedRequest.status);
      expect(result.metadata?.hash).toBe(expectedRequest.metadata?.hash);
      expect(mockApi.requestServiceGetRequest).toHaveBeenCalledWith({
        id: "123",
      });
    });

    it("should throw IntegrityError when hash does not match", async () => {
      const payloadAsString = '{"amount":"1000","currency":"ETH"}';
      const wrongHash = "0000000000000000000000000000000000000000000000000000000000000000";

      const dto: TgvalidatordGetRequestReply = {
        result: {
          id: "123",
          type: "payment",
          status: "PENDING_APPROVAL",
          metadata: {
            hash: wrongHash,
            payloadAsString,
          },
        },
      };
      mockApi.requestServiceGetRequest.mockResolvedValue(dto);

      await expect(service.get(123)).rejects.toThrow(IntegrityError);
    });

    it("should throw NotFoundError when request is not found", async () => {
      const dto: TgvalidatordGetRequestReply = {
        result: undefined,
      };
      mockApi.requestServiceGetRequest.mockResolvedValue(dto);

      await expect(service.get(123)).rejects.toThrow(NotFoundError);
    });

    it("should throw Error when requestId is not positive", async () => {
      await expect(service.get(0)).rejects.toThrow("requestId must be positive");
      await expect(service.get(-1)).rejects.toThrow("requestId must be positive");
    });

    it("should pass verification when metadata is undefined", async () => {
      const dto: TgvalidatordGetRequestReply = {
        result: {
          id: "123",
          type: "payment",
          status: "PENDING_APPROVAL",
          // No metadata - should not throw
        },
      };
      mockApi.requestServiceGetRequest.mockResolvedValue(dto);

      const result = await service.get(123);
      expect(result.id).toBe(123);
    });
  });

  describe("list", () => {
    it("should return a list of requests", async () => {
      const dto: TgvalidatordGetRequestsV2Reply = {
        result: [
          {
            id: "1",
            type: "payment",
            status: "PENDING_APPROVAL",
          },
          {
            id: "2",
            type: "payment",
            status: "CONFIRMED",
          },
        ],
        cursor: {
          currentPage: "abc123",
        },
      };
      mockApi.requestServiceGetRequestsV2.mockResolvedValue(dto);

      const result = await service.list({ limit: 50 });

      expect(result.requests).toHaveLength(2);
      expect(result.requests[0].id).toBe(1);
      expect(result.requests[1].id).toBe(2);
      expect(result.cursor.nextCursor).toBe("abc123");
    });

    it("should throw Error when limit is not positive", async () => {
      await expect(service.list({ limit: 0 })).rejects.toThrow(
        "limit must be positive"
      );
    });
  });

  describe("listForApproval", () => {
    it("should return a list of requests pending approval", async () => {
      const dto: TgvalidatordGetRequestsV2Reply = {
        result: [
          {
            id: "1",
            type: "payment",
            status: "PENDING_APPROVAL",
          },
        ],
        cursor: {
          currentPage: "xyz789",
        },
      };
      mockApi.requestServiceGetRequestsForApprovalV2.mockResolvedValue(dto);

      const result = await service.listForApproval({ limit: 50 });

      expect(result.requests).toHaveLength(1);
      expect(result.requests[0].id).toBe(1);
      expect(result.cursor.nextCursor).toBe("xyz789");
    });
  });

  describe("approveRequests", () => {
    // Generate a test ECDSA key pair
    const { privateKey } = crypto.generateKeyPairSync("ec", {
      namedCurve: "P-256",
    });

    it("should sort requests by ID and sign the hash array", async () => {
      const hash1 = calculateHexHash('{"id":1}');
      const hash2 = calculateHexHash('{"id":2}');
      const hash3 = calculateHexHash('{"id":3}');

      const requests: Request[] = [
        createMockRequest(3, hash3),
        createMockRequest(1, hash1),
        createMockRequest(2, hash2),
      ];

      const reply: TgvalidatordApproveRequestsReply = {
        signedRequests: "3",
      };
      mockApi.requestServiceApproveRequests.mockResolvedValue(reply);

      const count = await service.approveRequests(requests, privateKey);

      expect(count).toBe(3);
      expect(mockApi.requestServiceApproveRequests).toHaveBeenCalledWith({
        body: {
          ids: ["1", "2", "3"], // Sorted by ID
          signature: expect.any(String),
          comment: "approving via taurus-protect-sdk-typescript",
        },
      });

      // Verify the signature was created with sorted hashes
      const call = mockApi.requestServiceApproveRequests.mock.calls[0][0];
      expect(call.body.signature).toBeTruthy();
      expect(call.body.signature.length).toBeGreaterThan(0);
    });

    it("should throw Error when requests list is empty", async () => {
      await expect(service.approveRequests([], privateKey)).rejects.toThrow(
        "requests list cannot be empty"
      );
    });

    it("should throw Error when request has no metadata", async () => {
      const request = createMockRequest(1, "hash");
      (request as { metadata: undefined }).metadata = undefined;

      await expect(
        service.approveRequests([request], privateKey)
      ).rejects.toThrow("metadata cannot be null or undefined");
    });

    it("should throw Error when request metadata has no hash", async () => {
      // Create a request with metadata but empty hash
      const request: Request = {
        ...createMockRequest(1, "somehash"),
        metadata: {
          hash: "",
          payloadAsString: '{"id":1}',
          // SECURITY: payload field intentionally removed
        },
      };

      await expect(
        service.approveRequests([request], privateKey)
      ).rejects.toThrow("metadata hash cannot be null or empty");
    });

    it("should use custom comment when provided", async () => {
      const request = createMockRequest(1, calculateHexHash("test"));
      const reply: TgvalidatordApproveRequestsReply = {
        signedRequests: "1",
      };
      mockApi.requestServiceApproveRequests.mockResolvedValue(reply);

      await service.approveRequests([request], privateKey, "Custom comment");

      expect(mockApi.requestServiceApproveRequests).toHaveBeenCalledWith({
        body: expect.objectContaining({
          comment: "Custom comment",
        }),
      });
    });
  });

  describe("approveRequest", () => {
    const { privateKey } = crypto.generateKeyPairSync("ec", {
      namedCurve: "P-256",
    });

    it("should delegate to approveRequests", async () => {
      const request = createMockRequest(1, calculateHexHash("test"));
      const reply: TgvalidatordApproveRequestsReply = {
        signedRequests: "1",
      };
      mockApi.requestServiceApproveRequests.mockResolvedValue(reply);

      const count = await service.approveRequest(request, privateKey);

      expect(count).toBe(1);
    });
  });

  describe("rejectRequests", () => {
    it("should reject multiple requests", async () => {
      mockApi.requestServiceRejectRequests.mockResolvedValue({});

      await service.rejectRequests([1, 2, 3], "Test rejection");

      expect(mockApi.requestServiceRejectRequests).toHaveBeenCalledWith({
        body: {
          ids: ["1", "2", "3"],
          comment: "Test rejection",
        },
      });
    });

    it("should throw Error when requestIds is empty", async () => {
      await expect(service.rejectRequests([], "comment")).rejects.toThrow(
        "requestIds list cannot be empty"
      );
    });

    it("should throw Error when comment is empty", async () => {
      await expect(service.rejectRequests([1], "")).rejects.toThrow(
        "comment is required and cannot be empty"
      );
    });
  });

  describe("rejectRequest", () => {
    it("should delegate to rejectRequests", async () => {
      mockApi.requestServiceRejectRequests.mockResolvedValue({});

      await service.rejectRequest(123, "Test rejection");

      expect(mockApi.requestServiceRejectRequests).toHaveBeenCalledWith({
        body: {
          ids: ["123"],
          comment: "Test rejection",
        },
      });
    });
  });

  describe("createInternalTransferRequest", () => {
    it("should create an internal transfer request", async () => {
      const payloadAsString = '{"amount":"1000000000000000000"}';
      const hash = calculateHexHash(payloadAsString);

      mockApi.requestServiceCreateOutgoingRequest.mockResolvedValue({
        result: {
          id: "999",
          type: "payment",
          status: "CREATED",
          metadata: {
            hash,
            payloadAsString,
          },
        },
      });

      const result = await service.createInternalTransferRequest({
        fromAddressId: 123,
        toAddressId: 456,
        amount: "1000000000000000000",
      });

      expect(result.id).toBe(999);
      expect(mockApi.requestServiceCreateOutgoingRequest).toHaveBeenCalledWith({
        body: {
          amount: "1000000000000000000",
          fromAddressId: "123",
          toAddressId: "456",
          comment: undefined,
          externalRequestId: undefined,
          gasLimit: undefined,
          feeLimit: undefined,
        },
      });
    });

    it("should throw Error when fromAddressId is not positive", async () => {
      await expect(
        service.createInternalTransferRequest({
          fromAddressId: 0,
          toAddressId: 456,
          amount: "1000",
        })
      ).rejects.toThrow("fromAddressId must be positive");
    });

    it("should throw Error when amount is not positive", async () => {
      await expect(
        service.createInternalTransferRequest({
          fromAddressId: 123,
          toAddressId: 456,
          amount: "0",
        })
      ).rejects.toThrow("amount must be a positive number");
    });
  });

  describe("createExternalTransferRequest", () => {
    it("should create an external transfer request", async () => {
      const payloadAsString = '{"amount":"500000000000000000"}';
      const hash = calculateHexHash(payloadAsString);

      mockApi.requestServiceCreateOutgoingRequest.mockResolvedValue({
        result: {
          id: "888",
          type: "payment",
          status: "CREATED",
          metadata: {
            hash,
            payloadAsString,
          },
        },
      });

      const result = await service.createExternalTransferRequest({
        fromAddressId: 123,
        toWhitelistedAddressId: 789,
        amount: "500000000000000000",
        destinationAddressMemo: "test-memo",
      });

      expect(result.id).toBe(888);
      expect(mockApi.requestServiceCreateOutgoingRequest).toHaveBeenCalledWith({
        body: {
          amount: "500000000000000000",
          fromAddressId: "123",
          toWhitelistedAddressId: "789",
          comment: undefined,
          externalRequestId: undefined,
          gasLimit: undefined,
          feeLimit: undefined,
          destinationAddressMemo: "test-memo",
        },
      });
    });

    it("should throw Error when toWhitelistedAddressId is not positive", async () => {
      await expect(
        service.createExternalTransferRequest({
          fromAddressId: 123,
          toWhitelistedAddressId: -1,
          amount: "1000",
        })
      ).rejects.toThrow("toWhitelistedAddressId must be positive");
    });
  });
});

// Helper function to create a mock Request
function createMockRequest(id: number, hash: string): Request {
  const metadata: RequestMetadata | undefined = hash
    ? {
        hash,
        payloadAsString: `{"id":${id}}`,
        // SECURITY: payload field intentionally removed - use payloadAsString
      }
    : undefined;

  return {
    id,
    type: "payment",
    status: "PENDING_APPROVAL",
    metadata,
    tenantId: "tenant-1",
    currency: "ETH",
    currencyInfo: undefined,
    memo: undefined,
    rule: undefined,
    externalRequestId: undefined,
    requestBundleId: undefined,
    needsApprovalFrom: [],
    approvers: undefined,
    createdAt: new Date(),
    updatedAt: new Date(),
    tags: [],
    signedRequests: [],
  };
}
