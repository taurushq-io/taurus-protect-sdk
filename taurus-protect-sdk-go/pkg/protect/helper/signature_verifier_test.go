package helper

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestConstantTimeCompare(t *testing.T) {
	tests := []struct {
		name     string
		a, b     string
		expected bool
	}{
		{"equal strings", "hello", "hello", true},
		{"equal empty", "", "", true},
		{"different lengths", "hello", "hello!", false},
		{"different content", "hello", "world", false},
		{"case sensitive", "Hello", "hello", false},
		{"hex strings equal", "abc123", "abc123", true},
		{"hex strings different", "abc123", "def456", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConstantTimeCompare(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("ConstantTimeCompare(%q, %q) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestVerifyHashCoverage(t *testing.T) {
	tests := []struct {
		name       string
		hash       string
		signatures []model.WhitelistSignature
		expected   bool
	}{
		{
			name:       "empty signatures",
			hash:       "abc123",
			signatures: nil,
			expected:   false,
		},
		{
			name: "hash found in first signature",
			hash: "abc123",
			signatures: []model.WhitelistSignature{
				{Hashes: []string{"abc123", "def456"}},
			},
			expected: true,
		},
		{
			name: "hash found in second signature",
			hash: "def456",
			signatures: []model.WhitelistSignature{
				{Hashes: []string{"abc123"}},
				{Hashes: []string{"def456", "ghi789"}},
			},
			expected: true,
		},
		{
			name: "hash not found",
			hash: "xyz999",
			signatures: []model.WhitelistSignature{
				{Hashes: []string{"abc123"}},
				{Hashes: []string{"def456"}},
			},
			expected: false,
		},
		{
			name: "signature with nil hashes",
			hash: "abc123",
			signatures: []model.WhitelistSignature{
				{Hashes: nil},
				{Hashes: []string{"abc123"}},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := VerifyHashCoverage(tt.hash, tt.signatures)
			if result != tt.expected {
				t.Errorf("VerifyHashCoverage(%q, ...) = %v, want %v", tt.hash, result, tt.expected)
			}
		})
	}
}

func TestDecodeBase64(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{"valid base64", base64.StdEncoding.EncodeToString([]byte("hello")), false},
		{"empty string", "", false}, // Empty string decodes to empty bytes without error
		{"invalid base64", "not-valid-base64!!!", true},
		{"valid empty content", base64.StdEncoding.EncodeToString([]byte("")), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeBase64(tt.input)
			if tt.expectError && err == nil {
				t.Error("DecodeBase64() expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("DecodeBase64() unexpected error: %v", err)
			}
		})
	}
}

func TestIsValidSignature(t *testing.T) {
	// Generate a test key pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// Test data
	testData := []byte("test message")
	hash := sha256.Sum256(testData)

	// Create a valid signature
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		t.Fatalf("Failed to sign: %v", err)
	}

	// Encode signature in the expected format (r || s as raw bytes)
	sigBytes := append(r.Bytes(), s.Bytes()...)
	validSig := base64.StdEncoding.EncodeToString(sigBytes)

	t.Run("valid signature", func(t *testing.T) {
		result := IsValidSignature(testData, validSig, []*ecdsa.PublicKey{&privateKey.PublicKey})
		if !result {
			t.Error("IsValidSignature() = false for valid signature")
		}
	})

	t.Run("invalid signature", func(t *testing.T) {
		result := IsValidSignature(testData, "invalid-base64", []*ecdsa.PublicKey{&privateKey.PublicKey})
		if result {
			t.Error("IsValidSignature() = true for invalid signature")
		}
	})

	t.Run("wrong data", func(t *testing.T) {
		result := IsValidSignature([]byte("different data"), validSig, []*ecdsa.PublicKey{&privateKey.PublicKey})
		if result {
			t.Error("IsValidSignature() = true for wrong data")
		}
	})

	t.Run("empty public keys", func(t *testing.T) {
		result := IsValidSignature(testData, validSig, nil)
		if result {
			t.Error("IsValidSignature() = true with no public keys")
		}
	})
}

func TestVerifyGovernanceRulesSignatures(t *testing.T) {
	// Generate test key pairs
	key1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key1: %v", err)
	}
	key2, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key2: %v", err)
	}

	// Test data
	rulesData := []byte("test rules data")
	hash := sha256.Sum256(rulesData)

	// Create valid signatures
	r1, s1, _ := ecdsa.Sign(rand.Reader, key1, hash[:])
	r2, s2, _ := ecdsa.Sign(rand.Reader, key2, hash[:])

	sig1Bytes := append(padTo32Bytes(r1.Bytes()), padTo32Bytes(s1.Bytes())...)
	sig2Bytes := append(padTo32Bytes(r2.Bytes()), padTo32Bytes(s2.Bytes())...)

	validSig1 := base64.StdEncoding.EncodeToString(sig1Bytes)
	validSig2 := base64.StdEncoding.EncodeToString(sig2Bytes)

	signatures := []*model.RuleUserSignature{
		{UserID: "user1", Signature: validSig1},
		{UserID: "user2", Signature: validSig2},
	}

	publicKeys := []*ecdsa.PublicKey{&key1.PublicKey, &key2.PublicKey}

	t.Run("meets threshold", func(t *testing.T) {
		err := VerifyGovernanceRulesSignatures(rulesData, signatures, publicKeys, 2)
		if err != nil {
			t.Errorf("VerifyGovernanceRulesSignatures() error: %v", err)
		}
	})

	t.Run("exceeds available", func(t *testing.T) {
		err := VerifyGovernanceRulesSignatures(rulesData, signatures, publicKeys, 3)
		if err == nil {
			t.Error("VerifyGovernanceRulesSignatures() expected error for threshold > available")
		}
	})

	t.Run("zero threshold with signatures", func(t *testing.T) {
		// Even with zero threshold, function requires signatures to be provided
		err := VerifyGovernanceRulesSignatures(rulesData, signatures, publicKeys, 0)
		if err != nil {
			t.Errorf("VerifyGovernanceRulesSignatures() error for zero threshold: %v", err)
		}
	})

	t.Run("empty signatures", func(t *testing.T) {
		err := VerifyGovernanceRulesSignatures(rulesData, nil, publicKeys, 1)
		if err == nil {
			t.Error("VerifyGovernanceRulesSignatures() expected error for empty signatures")
		}
	})
}

// padTo32Bytes pads a byte slice to 32 bytes (for ECDSA P-256 signatures).
func padTo32Bytes(b []byte) []byte {
	if len(b) >= 32 {
		return b
	}
	padded := make([]byte, 32)
	copy(padded[32-len(b):], b)
	return padded
}

// Note: isValidECDSASignature is an unexported function, so we test it indirectly
// through the exported IsValidSignature function. The tests above cover the
// underlying ECDSA verification logic.
