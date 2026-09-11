package com.taurushq.sdk.protect.client.cache;

import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.GovernanceRules;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.client.service.GovernanceRuleService;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Thread-safe cache for the decoded rules container with configurable TTL.
 * <p>
 * This cache stores the decoded rules container and refreshes it automatically
 * when the TTL expires. Used primarily for address signature verification
 * against HSM slot public keys.
 * <p>
 * This implementation avoids holding a lock during network I/O by using a
 * "fetching" flag pattern. When a refresh is needed, one thread will perform
 * the fetch while other threads wait for completion.
 */
public class RulesContainerCache {

    /**
     * Default cache TTL: 5 minutes.
     */
    public static final long DEFAULT_CACHE_TTL_MS = 5 * 60 * 1000L;

    private final GovernanceRuleService governanceRuleService;
    private final long cacheTtlMs;

    private final Object lock = new Object();
    private DecodedRulesContainer cachedContainer;
    private long cacheTimestamp;

    /**
     * Flag indicating whether a fetch is currently in progress.
     */
    private boolean fetching = false;

    /**
     * Exception captured during fetch, to be re-thrown to waiting threads.
     */
    private ApiException fetchException = null;

    /**
     * Creates a new rules container cache with default TTL (5 minutes).
     *
     * @param governanceRuleService the governance rule service for fetching rules
     */
    public RulesContainerCache(GovernanceRuleService governanceRuleService) {
        this(governanceRuleService, DEFAULT_CACHE_TTL_MS);
    }

    /**
     * Creates a new rules container cache with a custom TTL.
     *
     * @param governanceRuleService the governance rule service for fetching rules
     * @param cacheTtlMs            the cache time-to-live in milliseconds
     */
    public RulesContainerCache(GovernanceRuleService governanceRuleService, long cacheTtlMs) {
        checkNotNull(governanceRuleService, "governanceRuleService cannot be null");
        checkArgument(cacheTtlMs > 0, "cacheTtlMs must be positive");

        this.governanceRuleService = governanceRuleService;
        this.cacheTtlMs = cacheTtlMs;
    }

    /**
     * Gets the decoded rules container, fetching from the API if the cache is expired.
     * <p>
     * This method is thread-safe and will only fetch once if multiple threads
     * attempt to refresh simultaneously. The lock is released during network I/O
     * to prevent blocking other threads unnecessarily.
     *
     * @return the decoded rules container
     * @throws ApiException if fetching the rules fails
     */
    public DecodedRulesContainer getDecodedRulesContainer() throws ApiException {
        synchronized (lock) {
            long now = System.currentTimeMillis();

            // If cache is valid, return immediately
            if (cachedContainer != null && (now - cacheTimestamp) <= cacheTtlMs) {
                return cachedContainer;
            }

            // If another thread is already fetching, wait for it
            while (fetching) {
                try {
                    lock.wait();
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                    ApiException ex = new ApiException();
                    ex.setCode(503);
                    ex.setError("Interrupted");
                    ex.setMessage("Interrupted while waiting for rules container refresh");
                    ex.initCause(e);
                    throw ex;
                }
            }

            // Re-check cache after waking up (another thread may have refreshed it)
            now = System.currentTimeMillis();
            if (cachedContainer != null && (now - cacheTimestamp) <= cacheTtlMs) {
                return cachedContainer;
            }

            // Check if the previous fetch failed
            if (fetchException != null) {
                ApiException ex = fetchException;
                fetchException = null;
                throw ex;
            }

            // This thread will perform the fetch
            fetching = true;
            fetchException = null;
        }

        // Perform network I/O outside the lock
        DecodedRulesContainer newContainer = null;
        ApiException exception = null;
        boolean settled = false;
        try {
            newContainer = doFetch();
            settled = true;
        } catch (ApiException e) {
            exception = e;
            settled = true;
        } finally {
            if (!settled) {
                releaseTheFetchSlot();
            }
        }

        // Update cache and notify waiting threads
        synchronized (lock) {
            fetching = false;
            if (exception != null) {
                fetchException = exception;
                lock.notifyAll();
                throw exception;
            }
            cachedContainer = newContainer;
            cacheTimestamp = System.currentTimeMillis();
            lock.notifyAll();
            return cachedContainer;
        }
    }

    /**
     * Clears the single-flight flag and wakes every waiter, for the path where the fetch
     * threw something the {@code catch (ApiException)} above does not see.
     *
     * <p>This is not defensive tidying. {@code doFetch()} runs governance verification,
     * whose {@code IntegrityException} is <b>unchecked</b> ({@code extends
     * SecurityException}), and the error-body path can raise {@code StackOverflowError}.
     * Either one escaped before {@code fetching} was cleared, so every later caller parked
     * on the untimed {@code lock.wait()} and never woke: <b>one crafted {@code /rules}
     * response permanently wedged address, asset and price verification process-wide</b> --
     * a denial of service on the whole verification surface from a single response.
     *
     * <p>{@code fetchException} is deliberately NOT set here: it is typed
     * {@code ApiException} and the escaping throwable propagates to this caller on its
     * own. Waiters re-check the cache, find it stale, and take a fresh turn at the fetch
     * -- which is the correct outcome, since nothing was cached.
     */
    private void releaseTheFetchSlot() {
        synchronized (lock) {
            fetching = false;
            lock.notifyAll();
        }
    }

    /**
     * Forces a cache refresh, fetching the latest rules from the API.
     * <p>
     * This method is thread-safe and releases the lock during network I/O.
     *
     * @throws ApiException if fetching the rules fails
     */
    public void invalidate() throws ApiException {
        synchronized (lock) {
            // Wait if another thread is already fetching
            while (fetching) {
                try {
                    lock.wait();
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                    ApiException ex = new ApiException();
                    ex.setCode(503);
                    ex.setError("Interrupted");
                    ex.setMessage("Interrupted while waiting for rules container invalidation");
                    ex.initCause(e);
                    throw ex;
                }
            }

            // This thread will perform the fetch
            fetching = true;
            fetchException = null;
        }

        // Perform network I/O outside the lock
        DecodedRulesContainer newContainer = null;
        ApiException exception = null;
        boolean settled = false;
        try {
            newContainer = doFetch();
            settled = true;
        } catch (ApiException e) {
            exception = e;
            settled = true;
        } finally {
            if (!settled) {
                releaseTheFetchSlot();
            }
        }

        // Update cache and notify waiting threads
        synchronized (lock) {
            fetching = false;
            if (exception != null) {
                fetchException = exception;
                lock.notifyAll();
                throw exception;
            }
            cachedContainer = newContainer;
            cacheTimestamp = System.currentTimeMillis();
            lock.notifyAll();
        }
    }

    /**
     * Checks if the cache is currently valid (not expired).
     *
     * @return true if the cache is valid, false if expired or empty
     */
    public boolean isCacheValid() {
        synchronized (lock) {
            if (cachedContainer == null) {
                return false;
            }
            long now = System.currentTimeMillis();
            return (now - cacheTimestamp) <= cacheTtlMs;
        }
    }

    /**
     * Gets the configured cache TTL in milliseconds.
     *
     * @return the cache TTL
     */
    public long getCacheTtlMs() {
        return cacheTtlMs;
    }

    /**
     * Performs the actual fetch from the governance rule service.
     * This method should be called without holding the lock.
     *
     * @return the decoded rules container
     * @throws ApiException if fetching the rules fails
     */
    private DecodedRulesContainer doFetch() throws ApiException {
        GovernanceRules rules = governanceRuleService.getRules();
        if (rules == null) {
            ApiException ex = new ApiException();
            ex.setCode(503);  // Service Unavailable - rules service not responding
            ex.setError("ServiceUnavailable");
            ex.setMessage("No governance rules available from API");
            throw ex;
        }
        return governanceRuleService.getDecodedRulesContainer(rules);
    }
}
