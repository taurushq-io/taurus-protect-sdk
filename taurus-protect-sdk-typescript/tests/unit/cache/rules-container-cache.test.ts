/**
 * Unit tests for RulesContainerCache.
 *
 * These tests verify cache behavior including TTL expiration,
 * concurrent refresh handling, and error cases.
 */

import {
  RulesContainerCache,
  type RulesContainerProvider,
} from "../../../src/cache/rules-container-cache";
import type { DecodedRulesContainer } from "../../../src/models/governance-rules";
import { createEmptyRulesContainer } from "../../../src/models/governance-rules";

/**
 * Create a mock provider that returns a rules container.
 */
function createMockProvider(
  container?: DecodedRulesContainer
): jest.Mock<Promise<DecodedRulesContainer>> {
  const result = container ?? {
    ...createEmptyRulesContainer(),
    timestamp: Date.now(),
  };
  return jest.fn().mockResolvedValue(result);
}

describe("RulesContainerCache", () => {
  describe("constructor", () => {
    it("should accept a valid provider and default TTL", () => {
      const provider = createMockProvider();
      const cache = new RulesContainerCache(provider);
      expect(cache.ttl).toBe(RulesContainerCache.DEFAULT_CACHE_TTL_MS);
    });

    it("should accept a custom TTL", () => {
      const provider = createMockProvider();
      const cache = new RulesContainerCache(provider, 60000);
      expect(cache.ttl).toBe(60000);
    });

    it("should throw if provider is null", () => {
      expect(
        () => new RulesContainerCache(null as unknown as RulesContainerProvider)
      ).toThrow("provider cannot be null or undefined");
    });

    it("should throw if provider is undefined", () => {
      expect(
        () =>
          new RulesContainerCache(
            undefined as unknown as RulesContainerProvider
          )
      ).toThrow("provider cannot be null or undefined");
    });

    it("should throw if TTL is zero", () => {
      const provider = createMockProvider();
      expect(() => new RulesContainerCache(provider, 0)).toThrow(
        "ttlMs must be positive"
      );
    });

    it("should throw if TTL is negative", () => {
      const provider = createMockProvider();
      expect(() => new RulesContainerCache(provider, -100)).toThrow(
        "ttlMs must be positive"
      );
    });
  });

  describe("get", () => {
    it("should trigger fetch from provider on first call (cache miss)", async () => {
      const provider = createMockProvider();
      const cache = new RulesContainerCache(provider);

      const result = await cache.get();

      expect(provider).toHaveBeenCalledTimes(1);
      expect(result).toBeDefined();
      expect(result.users).toEqual([]);
    });

    it("should return cached value on second call without re-fetching", async () => {
      const provider = createMockProvider();
      const cache = new RulesContainerCache(provider);

      const result1 = await cache.get();
      const result2 = await cache.get();

      expect(provider).toHaveBeenCalledTimes(1);
      expect(result1).toBe(result2);
    });

    it("should re-fetch after TTL expires", async () => {
      const container1 = {
        ...createEmptyRulesContainer(),
        timestamp: 1,
      };
      const container2 = {
        ...createEmptyRulesContainer(),
        timestamp: 2,
      };

      const provider = jest
        .fn<Promise<DecodedRulesContainer>, []>()
        .mockResolvedValueOnce(container1)
        .mockResolvedValueOnce(container2);

      // Use a very short TTL
      const cache = new RulesContainerCache(provider, 10);

      const result1 = await cache.get();
      expect(result1.timestamp).toBe(1);

      // Wait for TTL to expire
      await new Promise((resolve) => setTimeout(resolve, 20));

      const result2 = await cache.get();
      expect(result2.timestamp).toBe(2);
      expect(provider).toHaveBeenCalledTimes(2);
    });

    it("should throw when provider returns null", async () => {
      const provider = jest
        .fn()
        .mockResolvedValue(null) as jest.Mock<
        Promise<DecodedRulesContainer>
      >;
      const cache = new RulesContainerCache(provider);

      await expect(cache.get()).rejects.toThrow(
        "Provider returned null or undefined rules container"
      );
    });

    it("should throw when provider returns undefined", async () => {
      const provider = jest
        .fn()
        .mockResolvedValue(undefined) as jest.Mock<
        Promise<DecodedRulesContainer>
      >;
      const cache = new RulesContainerCache(provider);

      await expect(cache.get()).rejects.toThrow(
        "Provider returned null or undefined rules container"
      );
    });

    it("should propagate provider errors", async () => {
      const provider = jest
        .fn()
        .mockRejectedValue(new Error("API connection failed")) as jest.Mock<
        Promise<DecodedRulesContainer>
      >;
      const cache = new RulesContainerCache(provider);

      await expect(cache.get()).rejects.toThrow("API connection failed");
    });

    it("should share refresh promise for concurrent calls", async () => {
      let resolveProvider: ((value: DecodedRulesContainer) => void) | undefined;
      const provider = jest.fn().mockImplementation(
        () =>
          new Promise<DecodedRulesContainer>((resolve) => {
            resolveProvider = resolve;
          })
      );
      const cache = new RulesContainerCache(provider);

      // Start two concurrent gets
      const promise1 = cache.get();
      const promise2 = cache.get();

      // Both should be waiting on the same promise
      expect(provider).toHaveBeenCalledTimes(1);

      // Resolve the provider
      resolveProvider!({
        ...createEmptyRulesContainer(),
        timestamp: 42,
      });

      const [result1, result2] = await Promise.all([promise1, promise2]);
      expect(result1).toBe(result2);
      expect(provider).toHaveBeenCalledTimes(1);
    });
  });

  describe("refresh", () => {
    it("should force a fetch even when cache is valid", async () => {
      const provider = createMockProvider();
      const cache = new RulesContainerCache(provider);

      await cache.get();
      expect(provider).toHaveBeenCalledTimes(1);

      await cache.refresh();
      expect(provider).toHaveBeenCalledTimes(2);
    });
  });

  describe("clear", () => {
    it("should invalidate cache so next get triggers fetch", async () => {
      const provider = createMockProvider();
      const cache = new RulesContainerCache(provider);

      await cache.get();
      expect(provider).toHaveBeenCalledTimes(1);

      cache.clear();
      expect(cache.isValid()).toBe(false);

      await cache.get();
      expect(provider).toHaveBeenCalledTimes(2);
    });
  });

  describe("isValid", () => {
    it("should return false when cache is empty", () => {
      const provider = createMockProvider();
      const cache = new RulesContainerCache(provider);
      expect(cache.isValid()).toBe(false);
    });

    it("should return true after successful get", async () => {
      const provider = createMockProvider();
      const cache = new RulesContainerCache(provider);
      await cache.get();
      expect(cache.isValid()).toBe(true);
    });

    it("should return false after clear", async () => {
      const provider = createMockProvider();
      const cache = new RulesContainerCache(provider);
      await cache.get();
      cache.clear();
      expect(cache.isValid()).toBe(false);
    });
  });

  describe("getCached", () => {
    it("should return undefined when cache is empty", () => {
      const provider = createMockProvider();
      const cache = new RulesContainerCache(provider);
      expect(cache.getCached()).toBeUndefined();
    });

    it("should return the container when cache is valid", async () => {
      const container = {
        ...createEmptyRulesContainer(),
        timestamp: 99,
      };
      const provider = createMockProvider(container);
      const cache = new RulesContainerCache(provider);

      await cache.get();
      const cached = cache.getCached();

      expect(cached).toBeDefined();
      expect(cached!.timestamp).toBe(99);
    });

    it("should return undefined after clear", async () => {
      const provider = createMockProvider();
      const cache = new RulesContainerCache(provider);
      await cache.get();
      cache.clear();
      expect(cache.getCached()).toBeUndefined();
    });
  });

  describe("DEFAULT_CACHE_TTL_MS", () => {
    it("should be 5 minutes (300000 ms)", () => {
      expect(RulesContainerCache.DEFAULT_CACHE_TTL_MS).toBe(300000);
    });
  });
});
