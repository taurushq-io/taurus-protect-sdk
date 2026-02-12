/**
 * Thread-safe cache for decoded governance rules containers.
 *
 * This cache stores decoded rules containers with a configurable TTL.
 * It's used by AddressService for HSM key lookup during signature verification.
 */

import type { DecodedRulesContainer } from '../models/governance-rules';

/**
 * Provider function that fetches the rules container.
 */
export type RulesContainerProvider = () => Promise<DecodedRulesContainer>;

/**
 * Thread-safe cache for governance rules containers.
 *
 * The cache automatically refreshes from the provider when expired.
 * Uses a simple time-based expiration strategy.
 *
 * Note: JavaScript is single-threaded, but async operations can interleave.
 * This implementation handles concurrent async calls by sharing the refresh promise.
 */
export class RulesContainerCache {
  /** Default TTL of 5 minutes (300,000 ms). */
  static readonly DEFAULT_CACHE_TTL_MS = 300000;

  private container: DecodedRulesContainer | undefined;
  private expiresAt: number = 0;
  private refreshPromise: Promise<DecodedRulesContainer> | undefined;

  /**
   * Creates a new RulesContainerCache.
   *
   * @param provider - Function that fetches the rules container
   * @param ttlMs - Time-to-live in milliseconds (default: 5 minutes)
   * @throws Error if provider is not provided or ttlMs is not positive
   */
  constructor(
    private readonly provider: RulesContainerProvider,
    private readonly ttlMs: number = RulesContainerCache.DEFAULT_CACHE_TTL_MS
  ) {
    if (!provider) {
      throw new Error('provider cannot be null or undefined');
    }
    if (ttlMs <= 0) {
      throw new Error('ttlMs must be positive');
    }
  }

  /**
   * Gets the configured cache TTL in milliseconds.
   */
  get ttl(): number {
    return this.ttlMs;
  }

  /**
   * Gets the cached rules container, refreshing if expired.
   *
   * This method is safe for concurrent async calls - if a refresh is already
   * in progress, concurrent calls will share the same refresh promise.
   *
   * @returns The cached or freshly fetched rules container
   * @throws Error if the provider fails to return a rules container
   */
  async get(): Promise<DecodedRulesContainer> {
    const now = Date.now();

    // Check if cache is still valid
    if (this.container && now < this.expiresAt) {
      return this.container;
    }

    // If already refreshing, wait for that
    if (this.refreshPromise) {
      return this.refreshPromise;
    }

    // Start refresh
    this.refreshPromise = this.refresh();

    try {
      return await this.refreshPromise;
    } finally {
      this.refreshPromise = undefined;
    }
  }

  /**
   * Forces a cache refresh, fetching the latest rules from the provider.
   *
   * @returns The freshly fetched rules container
   * @throws Error if the provider fails
   */
  async refresh(): Promise<DecodedRulesContainer> {
    const container = await this.provider();
    if (!container) {
      throw new Error('Provider returned null or undefined rules container');
    }
    this.container = container;
    this.expiresAt = Date.now() + this.ttlMs;
    return container;
  }

  /**
   * Clears the cache without refreshing.
   */
  clear(): void {
    this.container = undefined;
    this.expiresAt = 0;
  }

  /**
   * Checks if the cache has a valid (non-expired) entry.
   *
   * @returns True if the cache is valid, false if expired or empty
   */
  isValid(): boolean {
    return this.container !== undefined && Date.now() < this.expiresAt;
  }

  /**
   * Gets the cached container without refreshing.
   * Returns undefined if the cache is empty or expired.
   *
   * @returns The cached container or undefined
   */
  getCached(): DecodedRulesContainer | undefined {
    if (this.isValid()) {
      return this.container;
    }
    return undefined;
  }
}
