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
// number of DISTINCT SuperAdmin keys.
// Returns nil if verification succeeds, error otherwise.
//
// minValidSignatures counts distinct signing keys, not signature entries: ECDSA is
// randomized, so counting entries would let a single key produce as many valid
// signatures as any threshold demands, making a 2-of-N quorum no stronger than 1-of-N.
func VerifyGovernanceRulesSignatures(
	rulesContainerData []byte,
	signatures []*model.RuleUserSignature,
	superAdminKeys []*ecdsa.PublicKey,
	minValidSignatures int,
) error {
	// A non-positive threshold makes the comparison below vacuously true, so the
	// function would report success having matched zero signers. Java, Python and
	// TypeScript all reject it here; this SDK relied on its callers to.
	if minValidSignatures <= 0 {
		return fmt.Errorf("minValidSignatures must be positive")
	}

	if len(rulesContainerData) == 0 {
		return fmt.Errorf("rules container data cannot be empty")
	}

	if len(superAdminKeys) == 0 {
		return fmt.Errorf("no SuperAdmin keys configured for verification")
	}

	if len(signatures) == 0 {
		return fmt.Errorf("no signatures provided")
	}

	signers := make(map[string]struct{}, len(superAdminKeys))
	for _, sig := range signatures {
		if sig == nil || sig.Signature == "" {
			continue
		}

		if key := matchingKey(rulesContainerData, sig.Signature, superAdminKeys); key != nil {
			fingerprint, err := keyFingerprint(key)
			if err != nil {
				// An unfingerprintable key cannot be counted as a distinct signer.
				continue
			}
			signers[fingerprint] = struct{}{}
		}
	}

	if len(signers) < minValidSignatures {
		return fmt.Errorf("insufficient distinct valid SuperAdmin signers: got %d, need %d", len(signers), minValidSignatures)
	}

	return nil
}

// matchingKey returns the first configured key that verifies the signature, or nil.
func matchingKey(data []byte, base64Signature string, publicKeys []*ecdsa.PublicKey) *ecdsa.PublicKey {
	for _, key := range publicKeys {
		if key == nil {
			continue
		}

		valid, err := crypto.VerifySignature(key, data, base64Signature)
		if err == nil && valid {
			return key
		}
	}
	return nil
}

// keyFingerprint identifies a key by its encoded bytes, so the same key configured
// twice counts as one signer.
func keyFingerprint(key *ecdsa.PublicKey) (string, error) {
	return model.KeyFingerprint(key)
}

// IsValidSignature checks if a signature is valid against any of the provided public keys.
// The signature is expected to be base64-encoded.
func IsValidSignature(data []byte, base64Signature string, publicKeys []*ecdsa.PublicKey) bool {
	return matchingKey(data, base64Signature, publicKeys) != nil
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

// VerifyHashCoverage reports whether a hash is covered by at least one signature.
//
// Asks across every signature what containsHash asks within one, so it delegates
// rather than re-rolling the inner loop — these two are the ONLY places this SDK
// compares hash strings, and an inlined second loop is how copies start to drift.
//
// SECURITY: constant-time, and the loop does not break early.
func VerifyHashCoverage(hash string, signatures []model.WhitelistSignature) bool {
	found := false
	for _, sig := range signatures {
		if containsHash(sig.Hashes, hash) {
			found = true
			// No break: returning on a match would leak its position through timing.
		}
	}
	return found
}

// containsHash reports whether one signature's hashes list covers hash.
//
// The per-signature half of the pair. SECURITY: constant-time, and the loop does not
// break early — it returned on the first match here while Java, Python and TypeScript
// all scanned the full list, so Go was the outlier against the invariant this package's
// CLAUDE.md states. Behaviour vectors cannot catch that: early return and full scan
// agree on every output, differing only in timing.
func containsHash(hashes []string, hash string) bool {
	found := false
	for _, h := range hashes {
		if ConstantTimeCompare(h, hash) {
			found = true
			// No break, as above.
		}
	}
	return found
}

// VerifyAndDecodeRulesContainer verifies a rules container's SuperAdmin signatures and
// decodes it — steps 2-3 of both whitelist flows, and independent of the entity being
// verified. Lets a service verify each distinct container once per page.
func VerifyAndDecodeRulesContainer(
	rulesContainerBase64 string,
	rulesSignaturesBase64 string,
	superAdminKeys []*ecdsa.PublicKey,
	minValidSignatures int,
	rulesContainerDecoder func(base64Data string) (*model.DecodedRulesContainer, error),
	userSignaturesDecoder func(base64Data string) ([]*model.RuleUserSignature, error),
) (*model.DecodedRulesContainer, error) {
	if rulesContainerBase64 == "" {
		return nil, &model.IntegrityError{Message: "rulesContainer is empty"}
	}
	if rulesSignaturesBase64 == "" {
		return nil, &model.IntegrityError{Message: "rulesSignatures is empty"}
	}

	// Decode and verify signatures
	signatures, err := userSignaturesDecoder(rulesSignaturesBase64)
	if err != nil {
		return nil, &model.IntegrityError{
			Message: fmt.Sprintf("failed to decode rules signatures: %v", err),
		}
	}

	rulesData, err := DecodeBase64(rulesContainerBase64)
	if err != nil {
		return nil, &model.IntegrityError{
			Message: fmt.Sprintf("failed to decode rules container: %v", err),
		}
	}

	if err := VerifyGovernanceRulesSignatures(rulesData, signatures, superAdminKeys, minValidSignatures); err != nil {
		return nil, &model.IntegrityError{
			Message: fmt.Sprintf("rules container signature verification failed: %v", err),
		}
	}

	// Decode rules container
	container, err := rulesContainerDecoder(rulesContainerBase64)
	if err != nil {
		return nil, &model.IntegrityError{
			Message: fmt.Sprintf("failed to decode rules container: %v", err),
		}
	}

	return container, nil
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

// VerifyGovernanceRules verifies a GovernanceRuleset directly, decoding its base64 rules
// container and checking the SuperAdmin signatures on it.
//
// This is the helper-level convenience the other three SDKs expose (Java
// SignatureVerifier.verifyGovernanceRules, Python verify_governance_rules, TS
// verifyGovernanceRules). Only VerifyGovernanceRulesSignatures existed here, so a caller
// outside the service layer had to decode the container itself — and this package's own
// CLAUDE.md documented the missing name as if it were present.
//
// minValidSignatures counts distinct signing keys; see VerifyGovernanceRulesSignatures.
func VerifyGovernanceRules(
	rules *model.GovernanceRuleset,
	minValidSignatures int,
	superAdminKeys []*ecdsa.PublicKey,
) error {
	if rules == nil {
		return fmt.Errorf("rules cannot be nil")
	}
	if minValidSignatures <= 0 {
		return fmt.Errorf("minValidSignatures must be positive")
	}
	if len(superAdminKeys) == 0 {
		return fmt.Errorf("superAdminKeys cannot be empty")
	}
	if rules.RulesContainer == "" {
		return &model.IntegrityError{
			Message: "governance rules verification failed: rules container is empty",
		}
	}
	if len(rules.Signatures) == 0 {
		return &model.IntegrityError{
			Message: "governance rules verification failed: no signatures present",
		}
	}

	rulesData, err := base64.StdEncoding.DecodeString(rules.RulesContainer)
	if err != nil {
		return &model.IntegrityError{
			Message: fmt.Sprintf("governance rules verification failed: cannot decode rules container: %v", err),
		}
	}

	signatures := make([]*model.RuleUserSignature, len(rules.Signatures))
	for i := range rules.Signatures {
		signatures[i] = &model.RuleUserSignature{
			UserID:    rules.Signatures[i].UserID,
			Signature: rules.Signatures[i].Signature,
		}
	}

	if err := VerifyGovernanceRulesSignatures(rulesData, signatures, superAdminKeys, minValidSignatures); err != nil {
		return &model.IntegrityError{
			Message: fmt.Sprintf("governance rules verification failed: %v", err),
		}
	}

	return nil
}
