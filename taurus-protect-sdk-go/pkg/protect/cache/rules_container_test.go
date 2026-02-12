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

func TestRulesContainerCache_Set(t *testing.T) {
	cache := NewRulesContainerCache(time.Minute, nil)

	rules := &model.DecodedRulesContainer{}
	cache.Set(rules)

	got, err := cache.Get(context.Background())
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}
	if got != rules {
		t.Errorf("Get() = %v, want %v", got, rules)
	}
}

func TestRulesContainerCache_Set_ResetsExpiry(t *testing.T) {
	fetchCount := 0
	fetcher := func(ctx context.Context) (*model.DecodedRulesContainer, error) {
		fetchCount++
		return &model.DecodedRulesContainer{}, nil
	}

	cache := NewRulesContainerCache(50*time.Millisecond, fetcher)

	// Initial fetch
	cache.Get(context.Background())

	// Wait partial TTL
	time.Sleep(30 * time.Millisecond)

	// Manual set should reset expiry
	newRules := &model.DecodedRulesContainer{}
	cache.Set(newRules)

	// Wait another partial TTL (total 60ms from initial, but only 30ms from Set)
	time.Sleep(30 * time.Millisecond)

	// Should still use cached value (not expired from Set)
	got, _ := cache.Get(context.Background())
	if got != newRules {
		t.Errorf("Get() = %v, want manually set rules", got)
	}
	if fetchCount != 1 {
		t.Errorf("fetchCount = %v, want 1 (should not refetch)", fetchCount)
	}
}

func TestRulesContainerCache_Invalidate(t *testing.T) {
	rules := &model.DecodedRulesContainer{}
	cache := NewRulesContainerCache(time.Minute, nil)

	// Set value
	cache.Set(rules)

	// Verify it's cached
	if !cache.IsValid() {
		t.Error("Cache should be valid after Set()")
	}

	// Invalidate
	cache.Invalidate()

	// Verify it's cleared
	if cache.IsValid() {
		t.Error("Cache should be invalid after Invalidate()")
	}

	// Get should return nil (no fetcher)
	got, _ := cache.Get(context.Background())
	if got != nil {
		t.Errorf("Get() after Invalidate() = %v, want nil", got)
	}
}

func TestRulesContainerCache_IsValid(t *testing.T) {
	cache := NewRulesContainerCache(time.Minute, nil)

	// Initially invalid
	if cache.IsValid() {
		t.Error("New cache should be invalid")
	}

	// After Set, should be valid
	cache.Set(&model.DecodedRulesContainer{})
	if !cache.IsValid() {
		t.Error("Cache should be valid after Set()")
	}

	// After Invalidate, should be invalid
	cache.Invalidate()
	if cache.IsValid() {
		t.Error("Cache should be invalid after Invalidate()")
	}
}

func TestRulesContainerCache_IsValid_Expiration(t *testing.T) {
	cache := NewRulesContainerCache(10*time.Millisecond, nil)
	cache.Set(&model.DecodedRulesContainer{})

	// Initially valid
	if !cache.IsValid() {
		t.Error("Cache should be valid after Set()")
	}

	// Wait for expiration
	time.Sleep(20 * time.Millisecond)

	// Should be invalid after expiration
	if cache.IsValid() {
		t.Error("Cache should be invalid after expiration")
	}
}

func TestRulesContainerCache_SetFetcher(t *testing.T) {
	cache := NewRulesContainerCache(time.Minute, nil)

	// Without fetcher, returns nil
	got, _ := cache.Get(context.Background())
	if got != nil {
		t.Error("Get() should return nil without fetcher")
	}

	// Set fetcher
	rules := &model.DecodedRulesContainer{}
	cache.SetFetcher(func(ctx context.Context) (*model.DecodedRulesContainer, error) {
		return rules, nil
	})

	// Now Get should fetch
	got, _ = cache.Get(context.Background())
	if got != rules {
		t.Errorf("Get() = %v, want %v", got, rules)
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

func TestRulesContainerCache_ConcurrentGetAndSet(t *testing.T) {
	cache := NewRulesContainerCache(time.Minute, nil)

	var wg sync.WaitGroup

	// Concurrent Sets
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cache.Set(&model.DecodedRulesContainer{})
		}(i)
	}

	// Concurrent Gets
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.Get(context.Background())
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
