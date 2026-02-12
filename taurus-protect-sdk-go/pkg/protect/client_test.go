package protect

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"net/http"
	"sync"
	"testing"
	"time"
)

// testKey is a package-level test key generated once at init time for unit tests.
var testKey *ecdsa.PrivateKey

func init() {
	var err error
	testKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic("failed to generate test key: " + err.Error())
	}
}

// testSuperAdminKeyOpts returns options that provide a test SuperAdmin key and
// minValidSignatures=1, satisfying the mandatory integrity verification requirement.
func testSuperAdminKeyOpts() []Option {
	return []Option{
		WithSuperAdminKeys([]*ecdsa.PublicKey{&testKey.PublicKey}),
		WithMinValidSignatures(1),
	}
}

// newTestClient creates a client with credentials and test SuperAdmin keys for unit testing.
func newTestClient(t *testing.T) *Client {
	t.Helper()
	client, err := NewClient("https://api.example.com",
		WithCredentials("key", "deadbeef"),
		WithSuperAdminKeys([]*ecdsa.PublicKey{&testKey.PublicKey}),
		WithMinValidSignatures(1),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	return client
}

func TestNewClient(t *testing.T) {
	saOpts := testSuperAdminKeyOpts()

	tests := []struct {
		name    string
		host    string
		opts    []Option
		wantErr bool
	}{
		{
			name: "valid basic config",
			host: "https://api.example.com",
			opts: append([]Option{
				WithCredentials("test-key", "deadbeef"),
			}, saOpts...),
			wantErr: false,
		},
		{
			name: "with trailing slash",
			host: "https://api.example.com/",
			opts: append([]Option{
				WithCredentials("test-key", "deadbeef"),
			}, saOpts...),
			wantErr: false,
		},
		{
			name:    "missing credentials",
			host:    "https://api.example.com",
			opts:    []Option{},
			wantErr: true,
		},
		{
			name:    "empty host",
			host:    "",
			opts:    []Option{WithCredentials("test-key", "deadbeef")},
			wantErr: true,
		},
		{
			name: "invalid api secret",
			host: "https://api.example.com",
			opts: []Option{
				WithCredentials("test-key", "not-hex"),
			},
			wantErr: true,
		},
		{
			name: "missing super admin keys",
			host: "https://api.example.com",
			opts: []Option{
				WithCredentials("test-key", "deadbeef"),
			},
			wantErr: true,
		},
		{
			name: "with custom timeout",
			host: "https://api.example.com",
			opts: append([]Option{
				WithCredentials("test-key", "deadbeef"),
				WithHTTPTimeout(60 * time.Second),
			}, saOpts...),
			wantErr: false,
		},
		{
			name: "with custom http client",
			host: "https://api.example.com",
			opts: append([]Option{
				WithCredentials("test-key", "deadbeef"),
				WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
			}, saOpts...),
			wantErr: false,
		},
		{
			name: "with rules cache ttl",
			host: "https://api.example.com",
			opts: append([]Option{
				WithCredentials("test-key", "deadbeef"),
				WithRulesCacheTTL(10 * time.Minute),
			}, saOpts...),
			wantErr: false,
		},
		{
			name: "with min valid signatures but no keys",
			host: "https://api.example.com",
			opts: []Option{
				WithCredentials("test-key", "deadbeef"),
				WithMinValidSignatures(2),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.host, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if client != nil {
				defer client.Close()

				if !tt.wantErr {
					if client.HTTPClient() == nil {
						t.Error("Client should have HTTP client")
					}
				}
			}
		})
	}
}

func TestClient_BaseURL(t *testing.T) {
	opts := append([]Option{WithCredentials("key", "deadbeef")}, testSuperAdminKeyOpts()...)
	client, err := NewClient("https://api.example.com/", opts...)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Trailing slash should be removed
	if got := client.BaseURL(); got != "https://api.example.com" {
		t.Errorf("BaseURL() = %v, want %v", got, "https://api.example.com")
	}
}

func TestClient_Close(t *testing.T) {
	client := newTestClient(t)

	// Close should not error
	if err := client.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Close should be idempotent
	if err := client.Close(); err != nil {
		t.Errorf("Close() second call error = %v", err)
	}
}

func TestWithCredentials_Validation(t *testing.T) {
	tests := []struct {
		name      string
		apiKey    string
		apiSecret string
		wantErr   bool
	}{
		{"valid", "key", "deadbeef", false},
		{"empty key", "", "deadbeef", true},
		{"empty secret", "key", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithCredentials(tt.apiKey, tt.apiSecret)
			config := &clientConfig{}
			err := opt(config)
			if (err != nil) != tt.wantErr {
				t.Errorf("WithCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWithMinValidSignatures_Validation(t *testing.T) {
	tests := []struct {
		name    string
		n       int
		wantErr bool
	}{
		{"valid", 2, false},
		{"zero", 0, false},
		{"negative", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithMinValidSignatures(tt.n)
			config := &clientConfig{}
			err := opt(config)
			if (err != nil) != tt.wantErr {
				t.Errorf("WithMinValidSignatures(%d) error = %v, wantErr %v", tt.n, err, tt.wantErr)
			}
		})
	}
}

func TestWithRulesCacheTTL_Validation(t *testing.T) {
	tests := []struct {
		name    string
		ttl     time.Duration
		wantErr bool
	}{
		{"valid", 5 * time.Minute, false},
		{"zero", 0, false},
		{"negative", -1 * time.Second, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithRulesCacheTTL(tt.ttl)
			config := &clientConfig{}
			err := opt(config)
			if (err != nil) != tt.wantErr {
				t.Errorf("WithRulesCacheTTL(%v) error = %v, wantErr %v", tt.ttl, err, tt.wantErr)
			}
		})
	}
}

func TestWithHTTPTimeout_Validation(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		wantErr bool
	}{
		{"valid", 30 * time.Second, false},
		{"zero", 0, false},
		{"negative", -1 * time.Second, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithHTTPTimeout(tt.timeout)
			config := &clientConfig{}
			err := opt(config)
			if (err != nil) != tt.wantErr {
				t.Errorf("WithHTTPTimeout(%v) error = %v, wantErr %v", tt.timeout, err, tt.wantErr)
			}
		})
	}
}

// TestClient_ServiceGetters_ReturnsSameInstance tests that all service getters return non-nil
// and return the same instance on subsequent calls (singleton pattern).
func TestClient_ServiceGetters_ReturnsSameInstance(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	t.Run("Wallets", func(t *testing.T) {
		first := client.Wallets()
		second := client.Wallets()
		if first == nil {
			t.Error("Wallets() should not return nil")
		}
		if first != second {
			t.Error("Wallets() should return the same instance")
		}
	})

	t.Run("Addresses", func(t *testing.T) {
		first := client.Addresses()
		second := client.Addresses()
		if first == nil {
			t.Error("Addresses() should not return nil")
		}
		if first != second {
			t.Error("Addresses() should return the same instance")
		}
	})

	t.Run("Requests", func(t *testing.T) {
		first := client.Requests()
		second := client.Requests()
		if first == nil {
			t.Error("Requests() should not return nil")
		}
		if first != second {
			t.Error("Requests() should return the same instance")
		}
	})

	t.Run("Transactions", func(t *testing.T) {
		first := client.Transactions()
		second := client.Transactions()
		if first == nil {
			t.Error("Transactions() should not return nil")
		}
		if first != second {
			t.Error("Transactions() should return the same instance")
		}
	})

	t.Run("GovernanceRules", func(t *testing.T) {
		first := client.GovernanceRules()
		second := client.GovernanceRules()
		if first == nil {
			t.Error("GovernanceRules() should not return nil")
		}
		if first != second {
			t.Error("GovernanceRules() should return the same instance")
		}
	})

	t.Run("Balances", func(t *testing.T) {
		first := client.Balances()
		second := client.Balances()
		if first == nil {
			t.Error("Balances() should not return nil")
		}
		if first != second {
			t.Error("Balances() should return the same instance")
		}
	})

	t.Run("Currencies", func(t *testing.T) {
		first := client.Currencies()
		second := client.Currencies()
		if first == nil {
			t.Error("Currencies() should not return nil")
		}
		if first != second {
			t.Error("Currencies() should return the same instance")
		}
	})

	t.Run("WhitelistedAddresses", func(t *testing.T) {
		first := client.WhitelistedAddresses()
		second := client.WhitelistedAddresses()
		if first == nil {
			t.Error("WhitelistedAddresses() should not return nil")
		}
		if first != second {
			t.Error("WhitelistedAddresses() should return the same instance")
		}
	})

	t.Run("WhitelistedAssets", func(t *testing.T) {
		first := client.WhitelistedAssets()
		second := client.WhitelistedAssets()
		if first == nil {
			t.Error("WhitelistedAssets() should not return nil")
		}
		if first != second {
			t.Error("WhitelistedAssets() should return the same instance")
		}
	})

	t.Run("Fees", func(t *testing.T) {
		first := client.Fees()
		second := client.Fees()
		if first == nil {
			t.Error("Fees() should not return nil")
		}
		if first != second {
			t.Error("Fees() should return the same instance")
		}
	})

	t.Run("Audits", func(t *testing.T) {
		first := client.Audits()
		second := client.Audits()
		if first == nil {
			t.Error("Audits() should not return nil")
		}
		if first != second {
			t.Error("Audits() should return the same instance")
		}
	})

	t.Run("Changes", func(t *testing.T) {
		first := client.Changes()
		second := client.Changes()
		if first == nil {
			t.Error("Changes() should not return nil")
		}
		if first != second {
			t.Error("Changes() should return the same instance")
		}
	})

	t.Run("Prices", func(t *testing.T) {
		first := client.Prices()
		second := client.Prices()
		if first == nil {
			t.Error("Prices() should not return nil")
		}
		if first != second {
			t.Error("Prices() should return the same instance")
		}
	})

	t.Run("AirGap", func(t *testing.T) {
		first := client.AirGap()
		second := client.AirGap()
		if first == nil {
			t.Error("AirGap() should not return nil")
		}
		if first != second {
			t.Error("AirGap() should return the same instance")
		}
	})

	t.Run("Staking", func(t *testing.T) {
		first := client.Staking()
		second := client.Staking()
		if first == nil {
			t.Error("Staking() should not return nil")
		}
		if first != second {
			t.Error("Staking() should return the same instance")
		}
	})

	t.Run("WhitelistedContracts", func(t *testing.T) {
		first := client.WhitelistedContracts()
		second := client.WhitelistedContracts()
		if first == nil {
			t.Error("WhitelistedContracts() should not return nil")
		}
		if first != second {
			t.Error("WhitelistedContracts() should return the same instance")
		}
	})

	t.Run("BusinessRules", func(t *testing.T) {
		first := client.BusinessRules()
		second := client.BusinessRules()
		if first == nil {
			t.Error("BusinessRules() should not return nil")
		}
		if first != second {
			t.Error("BusinessRules() should return the same instance")
		}
	})

	t.Run("Reservations", func(t *testing.T) {
		first := client.Reservations()
		second := client.Reservations()
		if first == nil {
			t.Error("Reservations() should not return nil")
		}
		if first != second {
			t.Error("Reservations() should return the same instance")
		}
	})

	t.Run("Users", func(t *testing.T) {
		first := client.Users()
		second := client.Users()
		if first == nil {
			t.Error("Users() should not return nil")
		}
		if first != second {
			t.Error("Users() should return the same instance")
		}
	})

	t.Run("Groups", func(t *testing.T) {
		first := client.Groups()
		second := client.Groups()
		if first == nil {
			t.Error("Groups() should not return nil")
		}
		if first != second {
			t.Error("Groups() should return the same instance")
		}
	})

	t.Run("VisibilityGroups", func(t *testing.T) {
		first := client.VisibilityGroups()
		second := client.VisibilityGroups()
		if first == nil {
			t.Error("VisibilityGroups() should not return nil")
		}
		if first != second {
			t.Error("VisibilityGroups() should return the same instance")
		}
	})

	t.Run("Config", func(t *testing.T) {
		first := client.Config()
		second := client.Config()
		if first == nil {
			t.Error("Config() should not return nil")
		}
		if first != second {
			t.Error("Config() should return the same instance")
		}
	})

	t.Run("Webhooks", func(t *testing.T) {
		first := client.Webhooks()
		second := client.Webhooks()
		if first == nil {
			t.Error("Webhooks() should not return nil")
		}
		if first != second {
			t.Error("Webhooks() should return the same instance")
		}
	})

	t.Run("WebhookCalls", func(t *testing.T) {
		first := client.WebhookCalls()
		second := client.WebhookCalls()
		if first == nil {
			t.Error("WebhookCalls() should not return nil")
		}
		if first != second {
			t.Error("WebhookCalls() should return the same instance")
		}
	})

	t.Run("Tags", func(t *testing.T) {
		first := client.Tags()
		second := client.Tags()
		if first == nil {
			t.Error("Tags() should not return nil")
		}
		if first != second {
			t.Error("Tags() should return the same instance")
		}
	})

	t.Run("Assets", func(t *testing.T) {
		first := client.Assets()
		second := client.Assets()
		if first == nil {
			t.Error("Assets() should not return nil")
		}
		if first != second {
			t.Error("Assets() should return the same instance")
		}
	})

	t.Run("Actions", func(t *testing.T) {
		first := client.Actions()
		second := client.Actions()
		if first == nil {
			t.Error("Actions() should not return nil")
		}
		if first != second {
			t.Error("Actions() should return the same instance")
		}
	})

	t.Run("Blockchains", func(t *testing.T) {
		first := client.Blockchains()
		second := client.Blockchains()
		if first == nil {
			t.Error("Blockchains() should not return nil")
		}
		if first != second {
			t.Error("Blockchains() should return the same instance")
		}
	})

	t.Run("Exchanges", func(t *testing.T) {
		first := client.Exchanges()
		second := client.Exchanges()
		if first == nil {
			t.Error("Exchanges() should not return nil")
		}
		if first != second {
			t.Error("Exchanges() should return the same instance")
		}
	})

	t.Run("Fiat", func(t *testing.T) {
		first := client.Fiat()
		second := client.Fiat()
		if first == nil {
			t.Error("Fiat() should not return nil")
		}
		if first != second {
			t.Error("Fiat() should return the same instance")
		}
	})

	t.Run("FeePayers", func(t *testing.T) {
		first := client.FeePayers()
		second := client.FeePayers()
		if first == nil {
			t.Error("FeePayers() should not return nil")
		}
		if first != second {
			t.Error("FeePayers() should return the same instance")
		}
	})

	t.Run("Health", func(t *testing.T) {
		first := client.Health()
		second := client.Health()
		if first == nil {
			t.Error("Health() should not return nil")
		}
		if first != second {
			t.Error("Health() should return the same instance")
		}
	})

	t.Run("Jobs", func(t *testing.T) {
		first := client.Jobs()
		second := client.Jobs()
		if first == nil {
			t.Error("Jobs() should not return nil")
		}
		if first != second {
			t.Error("Jobs() should return the same instance")
		}
	})

	t.Run("Scores", func(t *testing.T) {
		first := client.Scores()
		second := client.Scores()
		if first == nil {
			t.Error("Scores() should not return nil")
		}
		if first != second {
			t.Error("Scores() should return the same instance")
		}
	})

	t.Run("Statistics", func(t *testing.T) {
		first := client.Statistics()
		second := client.Statistics()
		if first == nil {
			t.Error("Statistics() should not return nil")
		}
		if first != second {
			t.Error("Statistics() should return the same instance")
		}
	})

	t.Run("TokenMetadata", func(t *testing.T) {
		first := client.TokenMetadata()
		second := client.TokenMetadata()
		if first == nil {
			t.Error("TokenMetadata() should not return nil")
		}
		if first != second {
			t.Error("TokenMetadata() should return the same instance")
		}
	})

	t.Run("UserDevices", func(t *testing.T) {
		first := client.UserDevices()
		second := client.UserDevices()
		if first == nil {
			t.Error("UserDevices() should not return nil")
		}
		if first != second {
			t.Error("UserDevices() should return the same instance")
		}
	})

	t.Run("MultiFactorSignature", func(t *testing.T) {
		first := client.MultiFactorSignature()
		second := client.MultiFactorSignature()
		if first == nil {
			t.Error("MultiFactorSignature() should not return nil")
		}
		if first != second {
			t.Error("MultiFactorSignature() should return the same instance")
		}
	})

	t.Run("TaurusNetwork", func(t *testing.T) {
		first := client.TaurusNetwork()
		second := client.TaurusNetwork()
		if first == nil {
			t.Error("TaurusNetwork() should not return nil")
		}
		if first != second {
			t.Error("TaurusNetwork() should return the same instance")
		}
	})
}

// TestTaurusNetworkClient_SubServiceGetters tests that TaurusNetwork sub-service getters
// return non-nil and return the same instance on subsequent calls.
func TestTaurusNetworkClient_SubServiceGetters(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	tn := client.TaurusNetwork()
	if tn == nil {
		t.Fatal("TaurusNetwork() should not return nil")
	}

	t.Run("Participants", func(t *testing.T) {
		first := tn.Participants()
		second := tn.Participants()
		if first == nil {
			t.Error("Participants() should not return nil")
		}
		if first != second {
			t.Error("Participants() should return the same instance")
		}
	})

	t.Run("Pledges", func(t *testing.T) {
		first := tn.Pledges()
		second := tn.Pledges()
		if first == nil {
			t.Error("Pledges() should not return nil")
		}
		if first != second {
			t.Error("Pledges() should return the same instance")
		}
	})

	t.Run("Lending", func(t *testing.T) {
		first := tn.Lending()
		second := tn.Lending()
		if first == nil {
			t.Error("Lending() should not return nil")
		}
		if first != second {
			t.Error("Lending() should return the same instance")
		}
	})

	t.Run("Settlements", func(t *testing.T) {
		first := tn.Settlements()
		second := tn.Settlements()
		if first == nil {
			t.Error("Settlements() should not return nil")
		}
		if first != second {
			t.Error("Settlements() should return the same instance")
		}
	})

	t.Run("Sharing", func(t *testing.T) {
		first := tn.Sharing()
		second := tn.Sharing()
		if first == nil {
			t.Error("Sharing() should not return nil")
		}
		if first != second {
			t.Error("Sharing() should return the same instance")
		}
	})
}

// TestClient_ServiceGetters_ConcurrentAccess tests thread safety of lazy initialization.
func TestClient_ServiceGetters_ConcurrentAccess(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	t.Run("Wallets_ConcurrentAccess", func(t *testing.T) {
		const numGoroutines = 100
		results := make(chan interface{}, numGoroutines)

		var wg sync.WaitGroup
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				results <- client.Wallets()
			}()
		}

		wg.Wait()
		close(results)

		var first interface{}
		for svc := range results {
			if svc == nil {
				t.Error("Wallets() returned nil during concurrent access")
			}
			if first == nil {
				first = svc
			} else if svc != first {
				t.Error("Concurrent Wallets() calls returned different instances")
			}
		}
	})

	t.Run("TaurusNetwork_ConcurrentAccess", func(t *testing.T) {
		const numGoroutines = 100
		results := make(chan *TaurusNetworkClient, numGoroutines)

		var wg sync.WaitGroup
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				results <- client.TaurusNetwork()
			}()
		}

		wg.Wait()
		close(results)

		var first *TaurusNetworkClient
		for svc := range results {
			if svc == nil {
				t.Error("TaurusNetwork() returned nil during concurrent access")
			}
			if first == nil {
				first = svc
			} else if svc != first {
				t.Error("Concurrent TaurusNetwork() calls returned different instances")
			}
		}
	})

	t.Run("MultipleServices_ConcurrentAccess", func(t *testing.T) {
		const numGoroutines = 50

		var wg sync.WaitGroup

		// Results for different services
		walletResults := make(chan interface{}, numGoroutines)
		addressResults := make(chan interface{}, numGoroutines)
		requestResults := make(chan interface{}, numGoroutines)

		// Spawn goroutines for different services concurrently
		for i := 0; i < numGoroutines; i++ {
			wg.Add(3)
			go func() {
				defer wg.Done()
				walletResults <- client.Wallets()
			}()
			go func() {
				defer wg.Done()
				addressResults <- client.Addresses()
			}()
			go func() {
				defer wg.Done()
				requestResults <- client.Requests()
			}()
		}

		wg.Wait()
		close(walletResults)
		close(addressResults)
		close(requestResults)

		// Verify each service returned consistent instances
		verifyConsistentInstances(t, "Wallets", walletResults)
		verifyConsistentInstances(t, "Addresses", addressResults)
		verifyConsistentInstances(t, "Requests", requestResults)
	})
}

// verifyConsistentInstances is a helper for concurrent access tests.
func verifyConsistentInstances(t *testing.T, serviceName string, results <-chan interface{}) {
	t.Helper()
	var first interface{}
	for svc := range results {
		if svc == nil {
			t.Errorf("%s() returned nil during concurrent access", serviceName)
		}
		if first == nil {
			first = svc
		} else if svc != first {
			t.Errorf("Concurrent %s() calls returned different instances", serviceName)
		}
	}
}

// TestClient_Close_ClearsAuth verifies that auth is nil after close.
func TestClient_Close_ClearsAuth(t *testing.T) {
	client := newTestClient(t)

	// Auth should be set before close
	if client.auth == nil {
		t.Error("auth should not be nil before Close()")
	}

	if err := client.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Auth should be nil after close
	if client.auth != nil {
		t.Error("auth should be nil after Close()")
	}
}

// TestClient_Close_MultipleClosesSafe verifies that Close() is idempotent.
func TestClient_Close_MultipleClosesSafe(t *testing.T) {
	client := newTestClient(t)

	// First close should succeed
	if err := client.Close(); err != nil {
		t.Errorf("First Close() error = %v", err)
	}

	// Second close should also succeed (idempotent)
	if err := client.Close(); err != nil {
		t.Errorf("Second Close() error = %v", err)
	}

	// Third close should also succeed
	if err := client.Close(); err != nil {
		t.Errorf("Third Close() error = %v", err)
	}
}

// TestClient_AccessAfterClose verifies that services accessed before close still work reference-wise.
func TestClient_AccessAfterClose(t *testing.T) {
	client := newTestClient(t)

	// Access services before close
	wallets := client.Wallets()
	addresses := client.Addresses()
	requests := client.Requests()

	// Verify services are non-nil
	if wallets == nil {
		t.Error("Wallets() should not return nil before Close()")
	}
	if addresses == nil {
		t.Error("Addresses() should not return nil before Close()")
	}
	if requests == nil {
		t.Error("Requests() should not return nil before Close()")
	}

	// Close the client
	if err := client.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Services obtained before close should still be valid references
	// (they don't need the auth to exist for reference equality)
	if wallets == nil {
		t.Error("Wallets reference should still be non-nil after Close()")
	}
	if addresses == nil {
		t.Error("Addresses reference should still be non-nil after Close()")
	}
	if requests == nil {
		t.Error("Requests reference should still be non-nil after Close()")
	}

	// Verify they're still the same instances if accessed again
	// Note: The client can still initialize services after close,
	// but the auth won't be available for API calls
	walletsAfter := client.Wallets()
	if wallets != walletsAfter {
		t.Error("Wallets() should return the same instance after Close()")
	}
}

// TestClient_ConfigurationAccessors tests that configuration accessors return correct values.
func TestClient_ConfigurationAccessors(t *testing.T) {
	saOpts := testSuperAdminKeyOpts()

	t.Run("DefaultConfiguration", func(t *testing.T) {
		client := newTestClient(t)
		defer client.Close()

		// MinValidSignatures should be 1 (from test helper)
		if client.MinValidSignatures() != 1 {
			t.Errorf("MinValidSignatures() = %d, want 1", client.MinValidSignatures())
		}

		// RulesCache should always be initialized
		if client.RulesCache() == nil {
			t.Error("RulesCache() should not return nil")
		}

		// HTTPClient should be initialized
		if client.HTTPClient() == nil {
			t.Error("HTTPClient() should not return nil")
		}

		// SuperAdminKeys should be non-empty (mandatory)
		if len(client.SuperAdminKeys()) == 0 {
			t.Error("SuperAdminKeys() should not be empty")
		}
	})

	t.Run("MissingSuperAdminKeys", func(t *testing.T) {
		_, err := NewClient("https://api.example.com",
			WithCredentials("key", "deadbeef"),
		)
		if err == nil {
			t.Fatal("NewClient() should fail without SuperAdmin keys")
		}
	})

	t.Run("CustomMinValidSignatures", func(t *testing.T) {
		client, err := NewClient("https://api.example.com",
			WithCredentials("key", "deadbeef"),
			WithSuperAdminKeys([]*ecdsa.PublicKey{&testKey.PublicKey}),
			WithMinValidSignatures(1),
		)
		if err != nil {
			t.Fatalf("NewClient() error = %v", err)
		}
		defer client.Close()

		if client.MinValidSignatures() != 1 {
			t.Errorf("MinValidSignatures() = %d, want 1", client.MinValidSignatures())
		}
	})

	t.Run("CustomHTTPClient", func(t *testing.T) {
		customClient := &http.Client{Timeout: 60 * time.Second}
		opts := append([]Option{
			WithCredentials("key", "deadbeef"),
			WithHTTPClient(customClient),
		}, saOpts...)
		client, err := NewClient("https://api.example.com", opts...)
		if err != nil {
			t.Fatalf("NewClient() error = %v", err)
		}
		defer client.Close()

		// HTTPClient should be non-nil (note: it wraps the custom client)
		if client.HTTPClient() == nil {
			t.Error("HTTPClient() should not return nil")
		}
	})

	t.Run("CustomRulesCacheTTL", func(t *testing.T) {
		opts := append([]Option{
			WithCredentials("key", "deadbeef"),
			WithRulesCacheTTL(10 * time.Minute),
		}, saOpts...)
		client, err := NewClient("https://api.example.com", opts...)
		if err != nil {
			t.Fatalf("NewClient() error = %v", err)
		}
		defer client.Close()

		// RulesCache should be initialized even with custom TTL
		if client.RulesCache() == nil {
			t.Error("RulesCache() should not return nil with custom TTL")
		}
	})

	t.Run("BaseURL", func(t *testing.T) {
		opts := append([]Option{
			WithCredentials("key", "deadbeef"),
		}, saOpts...)
		client, err := NewClient("https://api.example.com/", opts...)
		if err != nil {
			t.Fatalf("NewClient() error = %v", err)
		}
		defer client.Close()

		// Trailing slash should be trimmed
		expected := "https://api.example.com"
		if client.BaseURL() != expected {
			t.Errorf("BaseURL() = %q, want %q", client.BaseURL(), expected)
		}
	})
}
