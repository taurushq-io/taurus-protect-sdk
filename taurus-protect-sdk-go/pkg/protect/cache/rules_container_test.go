package cache

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewRulesContainerCache(t *testing.T) {
	tests := []struct {
		name    string
		ttl     time.Duration
		fetcher RulesContainerFetcher
	}{
		{
			name:    "with TTL and fetcher",
			ttl:     5 * time.Minute,
			fetcher: func(ctx context.Context) (*model.DecodedRulesContainer, error) { return nil, nil },
		},
		{
			name:    "with zero TTL",
			ttl:     0,
			fetcher: nil,
		},
		{
			name:    "with nil fetcher",
			ttl:     time.Hour,
			fetcher: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewRulesContainerCache(tt.ttl, tt.fetcher)
			if cache == nil {
				t.Fatal("NewRulesContainerCache() returned nil")
			}
			if cache.ttl != tt.ttl {
				t.Errorf("TTL = %v, want %v", cache.ttl, tt.ttl)
			}
		})
	}
}

func TestRulesContainerCache_Get_WithoutFetcher(t *testing.T) {
	cache := NewRulesContainerCache(time.Minute, nil)

	// Without fetcher, Get should return nil when cache is empty
	got, err := cache.Get(context.Background())
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}
	if got != nil {
		t.Errorf("Get() = %v, want nil", got)
	}
}

func TestRulesContainerCache_Get_WithFetcher(t *testing.T) {
	expectedRules := &model.DecodedRulesContainer{}
	fetchCount := 0

	fetcher := func(ctx context.Context) (*model.DecodedRulesContainer, error) {
		fetchCount++
		return expectedRules, nil
	}

	cache := NewRulesContainerCache(time.Minute, fetcher)

	// First call should fetch
	got, err := cache.Get(context.Background())
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}
	if got != expectedRules {
		t.Errorf("Get() = %v, want %v", got, expectedRules)
	}
	if fetchCount != 1 {
		t.Errorf("fetchCount = %v, want 1", fetchCount)
	}

	// Second call should use cache
	got, err = cache.Get(context.Background())
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}
	if got != expectedRules {
		t.Errorf("Get() = %v, want %v", got, expectedRules)
	}
	if fetchCount != 1 {
		t.Errorf("fetchCount = %v, want 1 (should use cache)", fetchCount)
	}
}

func TestRulesContainerCache_Get_FetcherError(t *testing.T) {
	expectedErr := errors.New("fetch failed")

	fetcher := func(ctx context.Context) (*model.DecodedRulesContainer, error) {
		return nil, expectedErr
	}

	cache := NewRulesContainerCache(time.Minute, fetcher)

	got, err := cache.Get(context.Background())
	if err != expectedErr {
		t.Errorf("Get() error = %v, want %v", err, expectedErr)
	}
	if got != nil {
		t.Errorf("Get() = %v, want nil", got)
	}
}

func TestRulesContainerCache_Get_Expiration(t *testing.T) {
	fetchCount := 0
	rules1 := &model.DecodedRulesContainer{}
	rules2 := &model.DecodedRulesContainer{}

	fetcher := func(ctx context.Context) (*model.DecodedRulesContainer, error) {
		fetchCount++
		if fetchCount == 1 {
			return rules1, nil
		}
		return rules2, nil
	}

	// Very short TTL for testing
	cache := NewRulesContainerCache(10*time.Millisecond, fetcher)

	// First fetch
	got, _ := cache.Get(context.Background())
	if got != rules1 {
		t.Errorf("First Get() = %v, want %v", got, rules1)
	}

	// Wait for expiration
	time.Sleep(20 * time.Millisecond)

	// Should refetch
	got, _ = cache.Get(context.Background())
	if got != rules2 {
		t.Errorf("Second Get() = %v, want %v", got, rules2)
	}
	if fetchCount != 2 {
		t.Errorf("fetchCount = %v, want 2", fetchCount)
	}
}

// seededCache returns a cache holding rules, populated the only way a container can
// now enter: through the constructor's fetcher. It also reports the fetch count, so a
// caller can tell a cache hit from a refetch.
func seededCache(t *testing.T, ttl time.Duration, rules *model.DecodedRulesContainer) (*RulesContainerCache, func() int) {
	t.Helper()
	fetches := 0
	cache := NewRulesContainerCache(ttl, func(context.Context) (*model.DecodedRulesContainer, error) {
		fetches++
		return rules, nil
	})
	if _, err := cache.Get(context.Background()); err != nil {
		t.Fatalf("seeding Get() error = %v", err)
	}
	return cache, func() int { return fetches }
}

func TestRulesContainerCache_Invalidate(t *testing.T) {
	rules := &model.DecodedRulesContainer{}
	cache, fetches := seededCache(t, time.Minute, rules)

	if !cache.IsValid() {
		t.Error("Cache should be valid after the seeding fetch")
	}

	cache.Invalidate()

	if cache.IsValid() {
		t.Error("Cache should be invalid after Invalidate()")
	}

	// The next Get must go back to the fetcher rather than serve the cleared value.
	got, err := cache.Get(context.Background())
	if err != nil {
		t.Errorf("Get() after Invalidate() error = %v", err)
	}
	if got != rules {
		t.Errorf("Get() after Invalidate() = %v, want a refetch", got)
	}
	if fetches() != 2 {
		t.Errorf("fetchCount = %d, want 2 (Invalidate must force a refetch)", fetches())
	}
}

func TestRulesContainerCache_IsValid(t *testing.T) {
	cache := NewRulesContainerCache(time.Minute, func(context.Context) (*model.DecodedRulesContainer, error) {
		return &model.DecodedRulesContainer{}, nil
	})

	// Initially invalid
	if cache.IsValid() {
		t.Error("New cache should be invalid")
	}

	// After a fetch, should be valid
	if _, err := cache.Get(context.Background()); err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if !cache.IsValid() {
		t.Error("Cache should be valid after a fetch")
	}

	// After Invalidate, should be invalid
	cache.Invalidate()
	if cache.IsValid() {
		t.Error("Cache should be invalid after Invalidate()")
	}
}

func TestRulesContainerCache_IsValid_Expiration(t *testing.T) {
	cache, _ := seededCache(t, 10*time.Millisecond, &model.DecodedRulesContainer{})

	// Initially valid
	if !cache.IsValid() {
		t.Error("Cache should be valid after the seeding fetch")
	}

	// Wait for expiration
	time.Sleep(20 * time.Millisecond)

	// Should be invalid after expiration
	if cache.IsValid() {
		t.Error("Cache should be invalid after expiration")
	}
}

func TestRulesContainerCache_ConcurrentAccess(t *testing.T) {
	var fetchCount int64

	fetcher := func(ctx context.Context) (*model.DecodedRulesContainer, error) {
		atomic.AddInt64(&fetchCount, 1)
		time.Sleep(10 * time.Millisecond) // Simulate slow fetch
		return &model.DecodedRulesContainer{}, nil
	}

	cache := NewRulesContainerCache(time.Minute, fetcher)

	var wg sync.WaitGroup
	numGoroutines := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := cache.Get(context.Background())
			if err != nil {
				t.Errorf("Get() error = %v", err)
			}
		}()
	}

	wg.Wait()

	// Due to double-check locking, only a few fetches should occur
	count := atomic.LoadInt64(&fetchCount)
	if count > 5 {
		t.Errorf("fetchCount = %v, expected fewer due to double-check locking", count)
	}
}

// Concurrent Get + Invalidate, which is the remaining write path now that a container
// can only enter through the fetcher. Run with -race.
func TestRulesContainerCache_ConcurrentGetAndInvalidate(t *testing.T) {
	cache := NewRulesContainerCache(time.Minute, func(context.Context) (*model.DecodedRulesContainer, error) {
		return &model.DecodedRulesContainer{}, nil
	})

	var wg sync.WaitGroup

	// Concurrent Gets
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = cache.Get(context.Background())
		}()
	}

	// Concurrent Invalidates
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.Invalidate()
		}()
	}

	wg.Wait()
	// Test passes if no race conditions detected
}

func TestRulesContainerCache_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	fetcher := func(ctx context.Context) (*model.DecodedRulesContainer, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(100 * time.Millisecond):
			return &model.DecodedRulesContainer{}, nil
		}
	}

	cache := NewRulesContainerCache(time.Minute, fetcher)

	// Cancel context immediately
	cancel()

	_, err := cache.Get(ctx)
	if err == nil {
		t.Error("Get() should return error when context is cancelled")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Get() error = %v, want context.Canceled", err)
	}
}
