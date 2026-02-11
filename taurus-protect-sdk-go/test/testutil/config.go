package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
)

// Identity represents a user identity with optional API credentials, private key, and public key.
type Identity struct {
	Index      int
	Name       string
	APIKey     string
	APISecret  string
	PrivateKey string // PEM string
	PublicKey  string // PEM string
}

// HasAPICredentials returns true if this identity has both API key and secret.
func (id *Identity) HasAPICredentials() bool {
	return id.APIKey != "" && id.APISecret != ""
}

// HasPrivateKey returns true if this identity has a private key.
func (id *Identity) HasPrivateKey() bool {
	return id.PrivateKey != ""
}

// HasPublicKey returns true if this identity has a public key.
func (id *Identity) HasPublicKey() bool {
	return id.PublicKey != ""
}

// TestConfig holds multi-identity test configuration.
type TestConfig struct {
	Host               string
	MinValidSignatures int
	Identities         []Identity
}

var (
	globalConfig     *TestConfig
	globalConfigOnce sync.Once
)

// LoadConfig loads configuration from test.properties file with environment variable overrides.
// The properties file is searched in multiple locations relative to the project root.
// Environment variables: PROTECT_API_HOST, PROTECT_MIN_VALID_SIGNATURES,
// PROTECT_API_KEY_N, PROTECT_API_SECRET_N, PROTECT_PRIVATE_KEY_N, PROTECT_PUBLIC_KEY_N.
func LoadConfig() *TestConfig {
	globalConfigOnce.Do(func() {
		globalConfig = loadConfigInternal()
	})
	return globalConfig
}

func loadConfigInternal() *TestConfig {
	props := loadPropertiesFile()

	config := &TestConfig{
		Host:               resolve("PROTECT_API_HOST", "host", props),
		MinValidSignatures: 2,
	}

	// Parse minValidSignatures
	minSigStr := resolve("PROTECT_MIN_VALID_SIGNATURES", "minValidSignatures", props)
	if minSigStr != "" {
		if n, err := strconv.Atoi(minSigStr); err == nil {
			config.MinValidSignatures = n
		}
	}

	// Load identities (supports gaps, e.g., identity.1 then identity.4)
	for i := 1; i < 100; i++ {
		name := resolve("", fmt.Sprintf("identity.%d.name", i), props)
		apiKey := resolve(fmt.Sprintf("PROTECT_API_KEY_%d", i), fmt.Sprintf("identity.%d.apiKey", i), props)
		apiSecret := resolve(fmt.Sprintf("PROTECT_API_SECRET_%d", i), fmt.Sprintf("identity.%d.apiSecret", i), props)
		privateKey := resolve(fmt.Sprintf("PROTECT_PRIVATE_KEY_%d", i), fmt.Sprintf("identity.%d.privateKey", i), props)
		publicKey := resolve(fmt.Sprintf("PROTECT_PUBLIC_KEY_%d", i), fmt.Sprintf("identity.%d.publicKey", i), props)

		// Skip gaps in identity indices
		if name == "" && apiKey == "" && apiSecret == "" && privateKey == "" && publicKey == "" {
			continue
		}

		if name == "" {
			name = fmt.Sprintf("identity-%d", i)
		}

		config.Identities = append(config.Identities, Identity{
			Index:      i,
			Name:       name,
			APIKey:     apiKey,
			APISecret:  apiSecret,
			PrivateKey: privateKey,
			PublicKey:  publicKey,
		})
	}

	return config
}

// GetIdentity returns the identity with the given 1-based index (supports gaps).
func (c *TestConfig) GetIdentity(index int) *Identity {
	for i := range c.Identities {
		if c.Identities[i].Index == index {
			return &c.Identities[i]
		}
	}
	return nil
}

// GetIdentityCount returns the number of configured identities.
func (c *TestConfig) GetIdentityCount() int {
	return len(c.Identities)
}

// IsEnabled returns true if tests should run.
// Tests are enabled if PROTECT_INTEGRATION_TEST=true or any identity has API credentials.
func (c *TestConfig) IsEnabled() bool {
	if os.Getenv("PROTECT_INTEGRATION_TEST") == "true" {
		return true
	}
	for i := range c.Identities {
		if c.Identities[i].HasAPICredentials() {
			return true
		}
	}
	return false
}

// GetSuperAdminKeys returns all public keys from identities that have them.
func (c *TestConfig) GetSuperAdminKeys() []string {
	var keys []string
	for i := range c.Identities {
		if c.Identities[i].HasPublicKey() {
			keys = append(keys, c.Identities[i].PublicKey)
		}
	}
	return keys
}

// resolve returns the first non-empty value from: environment variable, properties file, or empty string.
func resolve(envVar, propKey string, props map[string]string) string {
	if envVar != "" {
		if val := os.Getenv(envVar); val != "" {
			return val
		}
	}
	if props != nil {
		if val, ok := props[propKey]; ok && val != "" {
			return val
		}
	}
	return ""
}

// loadPropertiesFile searches for test.properties in several locations.
func loadPropertiesFile() map[string]string {
	// Get the directory of this source file to find test.properties relative to project root
	_, thisFile, _, _ := runtime.Caller(0)
	testutilDir := filepath.Dir(thisFile)
	goSDKRoot := filepath.Dir(filepath.Dir(testutilDir)) // test/testutil -> test -> go-sdk-root

	searchPaths := []string{
		filepath.Join(goSDKRoot, "test", "testutil", "test.properties"),
		filepath.Join(goSDKRoot, "test.properties"),
		"test/testutil/test.properties",
		"test.properties",
	}

	for _, path := range searchPaths {
		props, err := ParseProperties(path)
		if err == nil {
			return props
		}
	}
	return nil
}
