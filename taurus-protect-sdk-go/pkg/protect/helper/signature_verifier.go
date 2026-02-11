// Package helper provides signature verification and validation utilities.
package helper

import (
	"crypto/ecdsa"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// VerifyGovernanceRulesSignatures verifies that the rules container is signed by the required
// number of SuperAdmin keys.
// Returns nil if verification succeeds, error otherwise.
func VerifyGovernanceRulesSignatures(
	rulesContainerData []byte,
	signatures []*model.RuleUserSignature,
	superAdminKeys []*ecdsa.PublicKey,
	minValidSignatures int,
) error {
	if len(rulesContainerData) == 0 {
		return fmt.Errorf("rules container data cannot be empty")
	}

	if len(superAdminKeys) == 0 {
		return fmt.Errorf("no SuperAdmin keys configured for verification")
	}

	if len(signatures) == 0 {
		return fmt.Errorf("no signatures provided")
	}

	validCount := 0
	for _, sig := range signatures {
		if sig == nil || sig.Signature == "" {
			continue
		}

		if IsValidSignature(rulesContainerData, sig.Signature, superAdminKeys) {
			validCount++
		}
	}

	if validCount < minValidSignatures {
		return fmt.Errorf("insufficient valid SuperAdmin signatures: got %d, need %d", validCount, minValidSignatures)
	}

	return nil
}

// IsValidSignature checks if a signature is valid against any of the provided public keys.
// The signature is expected to be base64-encoded.
func IsValidSignature(data []byte, base64Signature string, publicKeys []*ecdsa.PublicKey) bool {
	for _, key := range publicKeys {
		if key == nil {
			continue
		}

		valid, err := crypto.VerifySignature(key, data, base64Signature)
		if err == nil && valid {
			return true
		}
	}
	return false
}

// VerifySignatureWithKey verifies a signature against a specific public key.
// Returns true if the signature is valid, false otherwise.
func VerifySignatureWithKey(data []byte, base64Signature string, publicKey *ecdsa.PublicKey) bool {
	if publicKey == nil {
		return false
	}

	valid, err := crypto.VerifySignature(publicKey, data, base64Signature)
	return err == nil && valid
}

// ConstantTimeCompare compares two strings in constant time to prevent timing attacks.
// Returns true if the strings are equal.
func ConstantTimeCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// VerifyHashCoverage checks if a hash is covered by at least one of the signatures.
// Returns true if the hash is found in at least one signature's hashes list.
// This function iterates through all signatures to prevent timing-based information leaks.
func VerifyHashCoverage(hash string, signatures []model.WhitelistSignature) bool {
	found := false
	for _, sig := range signatures {
		for _, h := range sig.Hashes {
			if ConstantTimeCompare(h, hash) {
				found = true
				// Continue iterating to prevent timing leak
			}
		}
	}
	return found
}

// DecodeBase64 decodes a base64-encoded string.
// Returns the decoded bytes or an error.
func DecodeBase64(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}

// EncodeBase64 encodes bytes to base64 string.
func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
