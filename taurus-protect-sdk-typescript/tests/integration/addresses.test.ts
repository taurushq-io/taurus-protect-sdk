/**
 * Integration tests for Addresses API.
 *
 * These tests require a live API connection. Configure via environment variables
 * or hard-coded defaults in helpers.ts.
 *
 * Tests use the high-level AddressService (client.addresses) rather than
 * the raw OpenAPI API to demonstrate proper SDK usage patterns.
 */

import { skipIfNotIntegration, getTestClient } from "./helpers";

describe("Integration: Addresses", () => {
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

  it("should list addresses", async () => {
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
      const result = await client.addresses.listWithOptions({ limit: 10 });

      const addresses = result.items;
      console.log(`Found ${addresses.length} addresses`);

      if (result.pagination) {
        console.log(`Total items: ${result.pagination.totalItems}`);
      }

      for (const address of addresses) {
        console.log("Address - All Fields:");
        console.log("=".repeat(50));
        console.log(`  id: ${address.id}`);
        console.log(`  walletId: ${address.walletId}`);
        console.log(`  address: ${address.address}`);
        console.log(`  alternateAddress: ${address.alternateAddress ?? "(undefined)"}`);
        console.log(`  label: ${address.label ?? "(undefined)"}`);
        console.log(`  comment: ${address.comment ?? "(undefined)"}`);
        console.log(`  customerId: ${address.customerId ?? "(undefined)"}`);
        console.log(`  externalAddressId: ${address.externalAddressId ?? "(undefined)"}`);
        console.log(`  currency: ${address.currency}`);
        console.log(`  addressPath: ${address.addressPath ?? "(undefined)"}`);
        console.log(`  addressIndex: ${address.addressIndex ?? "(undefined)"}`);
        console.log(`  nonce: ${address.nonce ?? "(undefined)"}`);
        console.log(`  status: ${address.status ?? "(undefined)"}`);
        console.log(`  signature: ${address.signature ?? "(undefined)"}`);
        console.log(`  disabled: ${address.disabled}`);
        console.log(`  canUseAllFunds: ${address.canUseAllFunds}`);
        console.log(`  createdAt: ${address.createdAt?.toISOString() ?? "(undefined)"}`);
        console.log(`  updatedAt: ${address.updatedAt?.toISOString() ?? "(undefined)"}`);
        console.log(`  attributes: [${address.attributes.length} items]`);
        console.log(`  linkedWhitelistedAddressIds: [${address.linkedWhitelistedAddressIds.length} items]`);
      }

      expect(Array.isArray(addresses)).toBe(true);
    } finally {
      client.close();
    }
  }, 30000); // Extended timeout: addresses endpoint is slow due to large dataset (7.5M+ addresses)

  it("should get address by ID", async () => {
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
      // First, get a list to find a valid address ID
      const listResult = await client.addresses.listWithOptions({ limit: 1 });

      const addresses = listResult.items;
      if (addresses.length === 0) {
        console.log("No addresses available for testing");
        return;
      }

      const addressIdStr = addresses[0].id;
      if (!addressIdStr) {
        console.log("Address ID is undefined");
        return;
      }

      const addressId = parseInt(addressIdStr, 10);
      const address = await client.addresses.get(addressId);

      expect(address).toBeDefined();

      console.log("Address - All Fields:");
      console.log("=".repeat(50));
      console.log(`  id: ${address.id}`);
      console.log(`  walletId: ${address.walletId}`);
      console.log(`  address: ${address.address}`);
      console.log(`  alternateAddress: ${address.alternateAddress ?? "(undefined)"}`);
      console.log(`  label: ${address.label ?? "(undefined)"}`);
      console.log(`  comment: ${address.comment ?? "(undefined)"}`);
      console.log(`  customerId: ${address.customerId ?? "(undefined)"}`);
      console.log(`  externalAddressId: ${address.externalAddressId ?? "(undefined)"}`);
      console.log(`  currency: ${address.currency}`);
      console.log(`  addressPath: ${address.addressPath ?? "(undefined)"}`);
      console.log(`  addressIndex: ${address.addressIndex ?? "(undefined)"}`);
      console.log(`  nonce: ${address.nonce ?? "(undefined)"}`);
      console.log(`  status: ${address.status ?? "(undefined)"}`);
      console.log(`  signature: ${address.signature ?? "(undefined)"}`);
      console.log(`  disabled: ${address.disabled}`);
      console.log(`  canUseAllFunds: ${address.canUseAllFunds}`);
      console.log(`  createdAt: ${address.createdAt?.toISOString() ?? "(undefined)"}`);
      console.log(`  updatedAt: ${address.updatedAt?.toISOString() ?? "(undefined)"}`);
      console.log(`  attributes: [${address.attributes.length} items]`);
      if (address.attributes.length > 0) {
        for (const attr of address.attributes) {
          console.log(`    - ${attr.key}: ${attr.value}`);
        }
      }
      console.log(`  linkedWhitelistedAddressIds: [${address.linkedWhitelistedAddressIds.length} items]`);
      if (address.linkedWhitelistedAddressIds.length > 0) {
        for (const linkedId of address.linkedWhitelistedAddressIds) {
          console.log(`    - ${linkedId}`);
        }
      }
    } finally {
      client.close();
    }
  }, 30000); // Extended timeout: addresses endpoint is slow due to large dataset
});
