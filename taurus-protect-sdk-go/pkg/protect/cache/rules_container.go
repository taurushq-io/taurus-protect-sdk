// Package cache provides thread-safe caching for the Taurus-PROTECT SDK.
package cache

import (
	"context"
	"sync"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// RulesContainerFetcher is a function that fetches the rules container.
type RulesContainerFetcher func(ctx context.Context) (*model.DecodedRulesContainer, error)

// RulesContainerCache provides a thread-safe cache for the decoded rules container.
// It automatically refreshes the cache when the TTL expires.
type RulesContainerCache struct {
	mu       sync.RWMutex
	rules    *model.DecodedRulesContainer
	expiry   time.Time
	ttl      time.Duration
	fetcher  RulesContainerFetcher
	fetching bool          // True when a fetch is in progress
	fetchCh  chan struct{} // Closed when fetch completes
	fetchErr error         // Error from the last fetch (propagated to waiters)
}

// NewRulesContainerCache creates a new cache with the given TTL and fetcher function.
// The fetcher function is called to refresh the cache when it expires, and is the only
// way a container enters the cache (see below).
//
// The fetcher MUST verify the container's SuperAdmin signatures before returning it —
// this cache supplies the HSM public key that address signature verification trusts.
// protect.NewClient wires it to GovernanceRuleService.GetDecodedRulesContainer, which
// does. A nil fetcher leaves the cache permanently empty; Get then returns nil, which
// the address verifier rejects rather than treating as "nothing to check".
func NewRulesContainerCache(ttl time.Duration, fetcher RulesContainerFetcher) *RulesContainerCache {
	return &RulesContainerCache{
		ttl:     ttl,
		fetcher: fetcher,
	}
}

// Get returns the cached rules container, refreshing if expired.
// If the cache is empty and no fetcher is configured, returns nil.
func (c *RulesContainerCache) Get(ctx context.Context) (*model.DecodedRulesContainer, error) {
	// Fast path: check if cache is valid
	c.mu.RLock()
	if c.rules != nil && time.Now().Before(c.expiry) {
		rules := c.rules
		c.mu.RUnlock()
		return rules, nil
	}
	c.mu.RUnlock()

	// Slow path: need to refresh cache
	c.mu.Lock()

	// Double-check after acquiring write lock
	if c.rules != nil && time.Now().Before(c.expiry) {
		rules := c.rules
		c.mu.Unlock()
		return rules, nil
	}

	// If no fetcher, return current value (may be nil)
	if c.fetcher == nil {
		rules := c.rules
		c.mu.Unlock()
		return rules, nil
	}

	// Check if another goroutine is already fetching
	if c.fetching {
		// Wait for the ongoing fetch to complete
		waitCh := c.fetchCh
		c.mu.Unlock()

		// Wait for fetch to complete or context cancellation
		select {
		case <-waitCh:
			// Fetch completed, get the result (and any error)
			c.mu.RLock()
			rules := c.rules
			fetchErr := c.fetchErr
			c.mu.RUnlock()
			if fetchErr != nil {
				return nil, fetchErr
			}
			return rules, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	// Mark that we are fetching and create a channel to signal completion
	c.fetching = true
	c.fetchCh = make(chan struct{})
	fetcher := c.fetcher
	c.mu.Unlock()

	// Fetch new rules without holding the lock
	// This prevents deadlocks when network calls block under concurrent load
	rules, err := fetcher(ctx)

	// Update cache and signal completion
	c.mu.Lock()
	c.fetching = false
	if err != nil {
		c.fetchErr = err
		close(c.fetchCh)
		c.mu.Unlock()
		return nil, err
	}

	c.fetchErr = nil // Clear any previous error on successful fetch
	c.rules = rules
	c.expiry = time.Now().Add(c.ttl)
	close(c.fetchCh)
	c.mu.Unlock()

	return rules, nil
}

// The constructor's fetcher is the ONLY way a container enters this cache.
//
// There used to be a Set(container) and a SetFetcher(fn) here, either of which seated
// a container that nothing had verified — and this cache supplies the HSM public key
// that address signature verification trusts, so an unverified one makes that check
// pass against whatever key the container carries. Neither had a caller outside tests,
// and no sibling SDK offers an equivalent: Java and Python inject the governance
// service, TypeScript injects a provider closure. Do not reintroduce them.

// Invalidate clears the cache, forcing a refresh on the next Get().
func (c *RulesContainerCache) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.rules = nil
	c.expiry = time.Time{}
}

// IsValid returns true if the cache contains a non-expired value.
func (c *RulesContainerCache) IsValid() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.rules != nil && time.Now().Before(c.expiry)
}
