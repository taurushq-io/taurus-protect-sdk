/**
 * Integration tests for Requests API.
 *
 * These tests require a live API connection. Configure via environment variables
 * or hard-coded defaults in helpers.ts.
 *
 * Tests use the high-level RequestService (client.requests) rather than
 * the raw OpenAPI API to demonstrate proper SDK usage patterns.
 */

import { skipIfNotIntegration, getTestClient } from "./helpers";
import { RequestStatus } from "../../src/models/request";

describe("Integration: Requests", () => {
  beforeEach(() => {
    try {
      skipIfNotIntegration();
    } catch (e) {
      if (e instanceof Error && e.message === "SKIP_INTEGRATION") {
        return;
      }
      throw e;
    }
  });

  it("should list requests", async () => {
    try {
      skipIfNotIntegration();
    } catch (e) {
      if (e instanceof Error && e.message === "SKIP_INTEGRATION") {
        return;
      }
      throw e;
    }

    const client = getTestClient();
    try {
      const result = await client.requests.list({ limit: 10 });

      const requests = result.requests;
      console.log(`Found ${requests.length} requests`);

      for (const request of requests) {
        console.log("Request - All Fields:");
        console.log("=".repeat(50));
        console.log(`  id: ${request.id}`);
        console.log(`  type: ${request.type}`);
        console.log(`  status: ${request.status}`);
        console.log(`  metadata.hash: ${request.metadata?.hash ?? "(undefined)"}`);
        console.log(`  metadata.payloadAsString: ${request.metadata?.payloadAsString ?? "(undefined)"}`);
        console.log(`  tenantId: ${request.tenantId ?? "(undefined)"}`);
        console.log(`  currency: ${request.currency ?? "(undefined)"}`);
        console.log(`  memo: ${request.memo ?? "(undefined)"}`);
        console.log(`  rule: ${request.rule ?? "(undefined)"}`);
        console.log(`  externalRequestId: ${request.externalRequestId ?? "(undefined)"}`);
        console.log(`  requestBundleId: ${request.requestBundleId ?? "(undefined)"}`);
        console.log(`  createdAt: ${request.createdAt?.toISOString() ?? "(undefined)"}`);
        console.log(`  updatedAt: ${request.updatedAt?.toISOString() ?? "(undefined)"}`);
        console.log(`  needsApprovalFrom: [${request.needsApprovalFrom.length} items]`);
        console.log(`  tags: [${request.tags.length} items]`);
        if (request.currencyInfo) {
          console.log("  currencyInfo:");
          console.log(`    id: ${request.currencyInfo.id ?? "(undefined)"}`);
          console.log(`    symbol: ${request.currencyInfo.symbol ?? "(undefined)"}`);
        } else {
          console.log("  currencyInfo: (undefined)");
        }
        if (request.metadata) {
          console.log("  metadata:");
          console.log(`    hash: ${request.metadata.hash}`);
        } else {
          console.log("  metadata: (undefined)");
        }
        if (request.approvers) {
          console.log(`  approvers: [${request.approvers.parallel.length} parallel groups]`);
        } else {
          console.log("  approvers: (undefined)");
        }
      }

      expect(Array.isArray(requests)).toBe(true);
    } finally {
      client.close();
    }
  });

  it("should get request by ID with metadata", async () => {
    try {
      skipIfNotIntegration();
    } catch (e) {
      if (e instanceof Error && e.message === "SKIP_INTEGRATION") {
        return;
      }
      throw e;
    }

    const client = getTestClient();
    try {
      // First, get a list to find a valid request ID
      const listResult = await client.requests.list({ limit: 1 });

      const requests = listResult.requests;
      if (requests.length === 0) {
        console.log("No requests available for testing");
        return;
      }

      const requestId = requests[0].id;
      if (!requestId) {
        console.log("Request ID is undefined");
        return;
      }

      const request = await client.requests.get(requestId);

      expect(request).toBeDefined();

      console.log("Request - All Fields:");
      console.log("=".repeat(50));
      console.log(`  id: ${request.id}`);
      console.log(`  type: ${request.type}`);
      console.log(`  status: ${request.status}`);
      console.log(`  metadata.hash: ${request.metadata?.hash ?? "(undefined)"}`);
      console.log(`  metadata.payloadAsString: ${request.metadata?.payloadAsString ?? "(undefined)"}`);
      console.log(`  tenantId: ${request.tenantId ?? "(undefined)"}`);
      console.log(`  currency: ${request.currency ?? "(undefined)"}`);
      console.log(`  memo: ${request.memo ?? "(undefined)"}`);
      console.log(`  rule: ${request.rule ?? "(undefined)"}`);
      console.log(`  externalRequestId: ${request.externalRequestId ?? "(undefined)"}`);
      console.log(`  requestBundleId: ${request.requestBundleId ?? "(undefined)"}`);
      console.log(`  createdAt: ${request.createdAt?.toISOString() ?? "(undefined)"}`);
      console.log(`  updatedAt: ${request.updatedAt?.toISOString() ?? "(undefined)"}`);
      console.log(`  needsApprovalFrom: [${request.needsApprovalFrom.length} items]`);
      if (request.needsApprovalFrom.length > 0) {
        for (const group of request.needsApprovalFrom) {
          console.log(`    - ${group}`);
        }
      }
      console.log(`  tags: [${request.tags.length} items]`);
      if (request.tags.length > 0) {
        for (const tag of request.tags) {
          console.log(`    - ${tag}`);
        }
      }
      if (request.currencyInfo) {
        console.log("  currencyInfo:");
        console.log(`    id: ${request.currencyInfo.id ?? "(undefined)"}`);
        console.log(`    symbol: ${request.currencyInfo.symbol ?? "(undefined)"}`);
        console.log(`    name: ${request.currencyInfo.name ?? "(undefined)"}`);
        console.log(`    decimals: ${request.currencyInfo.decimals ?? "(undefined)"}`);
        console.log(`    blockchain: ${request.currencyInfo.blockchain ?? "(undefined)"}`);
        console.log(`    network: ${request.currencyInfo.network ?? "(undefined)"}`);
      } else {
        console.log("  currencyInfo: (undefined)");
      }
      if (request.metadata) {
        console.log("  metadata:");
        console.log(`    hash: ${request.metadata.hash}`);
        console.log(`    payloadAsString: ${request.metadata.payloadAsString ?? "(undefined)"}`);
        // Note: payload field is intentionally omitted for security - use payloadAsString and JSON.parse
      } else {
        console.log("  metadata: (undefined)");
      }
      if (request.approvers) {
        console.log("  approvers:");
        console.log(`    parallel: [${request.approvers.parallel.length} groups]`);
        for (const parallelGroup of request.approvers.parallel) {
          console.log(`      sequential: [${parallelGroup.sequential.length} items]`);
          for (const seq of parallelGroup.sequential) {
            console.log(`        - externalGroupId: ${seq.externalGroupId ?? "(undefined)"}, minimumSignatures: ${seq.minimumSignatures ?? "(undefined)"}`);
          }
        }
      } else {
        console.log("  approvers: (undefined)");
      }
    } finally {
      client.close();
    }
  });

  it("should filter requests by status", async () => {
    try {
      skipIfNotIntegration();
    } catch (e) {
      if (e instanceof Error && e.message === "SKIP_INTEGRATION") {
        return;
      }
      throw e;
    }

    const client = getTestClient();
    try {
      // Query for confirmed requests (CONFIRMED is a valid API status)
      const result = await client.requests.list({
        limit: 10,
        statuses: [RequestStatus.CONFIRMED],
      });

      const requests = result.requests;
      console.log(`Found ${requests.length} confirmed requests`);

      for (const request of requests) {
        console.log(
          `Request: ID=${request.id}, Type=${request.type}, Status=${request.status}`
        );
        // Verify all returned requests have CONFIRMED status
        expect(request.status).toBe("CONFIRMED");
      }

      expect(Array.isArray(requests)).toBe(true);
    } finally {
      client.close();
    }
  });
});
