package helper

import (
	"crypto/ecdsa"
	"fmt"
	"strings"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// WhitelistedAssetVerifier provides methods to verify whitelisted asset (contract) integrity.
type WhitelistedAssetVerifier struct {
	superAdminKeys     []*ecdsa.PublicKey
	minValidSignatures int
}

// NewWhitelistedAssetVerifier creates a new verifier with the given configuration.
//
// Accepts an empty key set on purpose — see NewWhitelistedAddressVerifier.
func NewWhitelistedAssetVerifier(superAdminKeys []*ecdsa.PublicKey, minValidSignatures int) *WhitelistedAssetVerifier {
	return &WhitelistedAssetVerifier{
		superAdminKeys:     superAdminKeys,
		minValidSignatures: minValidSignatures,
	}
}

// AssetVerificationResult contains the result of verification and the decoded rules container.
type AssetVerificationResult struct {
	// RulesContainer is the decoded and verified rules container.
	RulesContainer *model.DecodedRulesContainer
	// VerifiedHash is the hash that was matched during verification.
	// This may differ from the input hash if a legacy hash format was used.
	VerifiedHash string
	// VerifiedPayload is the payload VerifiedHash covers. Step 6 for assets runs in the
	// service, so it must parse THIS rather than Metadata.PayloadAsString: when a legacy
	// variant matched, the delivered payload carries members no signature covered.
	VerifiedPayload string
}

// 5-step verification for a whitelisted asset, plus the parse that makes it usable.
//
//	envelope ──▶ 1 hash ──▶ 2 SuperAdmin sigs ──▶ 3 decode rules
//	                                                   │
//	             6 parse VERIFIED payload ◀── 5 thresholds ◀── 4 coverage
//	                        │                       │                │
//	                        ▼                       │      returns the hash it
//	                 returned to caller             │      matched, which is
//	                                                │      what step 5 checks
//	                                                ▼
//	                             per group: DISTINCT signers, keyed on the
//	                             container's public key, never the entry's userId
//
// Steps 1-5 prove the envelope is authentic; step 6 is what stops an unsigned
// value reaching the caller. Assets use ContractAddressWhitelistingRules and the
// ASSET legacy hashes (isNFT / kindType) — not the address ones.
//
// Steps 1-5 run here; step 6 runs in service/whitelisted_asset.go, which calls
// ParseWhitelistedAssetFromJSON on the payload this verifier just cleared.
// VerifiedAsset is proof that an asset went through VerifyWhitelistedAsset.
//
// Its fields are unexported, so a value built outside this package carries nothing:
// Asset() returns nil. That is the same guarantee RequestMetadata.entries already
// gives — unverified input yields zero values — with the difference that a function
// taking a VerifiedAsset now SAYS so in its signature instead of a comment.
//
// Be precise about the limit: Go permits `helper.VerifiedAsset{}` from another
// package. What it does not permit is populating it. A forged value is therefore
// useless rather than unforgeable, which is the same security property.
//
// It proves verification RAN, not that it ran against the right keys or thresholds —
// a structural guard against forgetting to call the verifier, not a cryptographic one.
type VerifiedAsset struct {
	asset *model.WhitelistedAsset
	rules *model.DecodedRulesContainer
	hash  string
}

// Asset returns the asset parsed from the verified payload, or nil for a value this
// package did not produce.
func (v VerifiedAsset) Asset() *model.WhitelistedAsset { return v.asset }

// RulesContainer returns the decoded, SuperAdmin-verified rules container.
func (v VerifiedAsset) RulesContainer() *model.DecodedRulesContainer { return v.rules }

// VerifiedHash returns the hash that was matched, which may be a legacy variant.
func (v VerifiedAsset) VerifiedHash() string { return v.hash }

// IsVerified reports whether this value came from the verifier.
func (v VerifiedAsset) IsVerified() bool { return v.asset != nil }

// VerifyWhitelistedAsset performs the complete 5-step verification of a whitelisted asset.
// This implements the same verification flow as the Java SDK.
//
// Steps:
// 1. Verify metadata hash (SHA-256 of payloadAsString == metadata.hash)
// 2. Verify rules container signatures (SuperAdmin signatures)
// 3. Decode rules container (base64 -> protobuf -> model)
// 4. Verify hash coverage (metadata.hash in at least one signature.hashes)
// 5. Verify whitelist signatures (user signatures meet governance thresholds)
//
// The function does not mutate the input asset. If a legacy hash was matched during
// verification, it is returned in the result's VerifiedHash field.
// If cachedRulesContainer is non-nil, steps 2-3 are skipped — already done when the
// page's cache was built.
func (v *WhitelistedAssetVerifier) VerifyWhitelistedAsset(
	asset *model.WhitelistedAsset,
	rulesContainerDecoder func(base64Data string) (*model.DecodedRulesContainer, error),
	userSignaturesDecoder func(base64Data string) ([]*model.RuleUserSignature, error),
	cachedRulesContainer ...*model.DecodedRulesContainer,
) (*AssetVerificationResult, error) {
	if asset == nil {
		return nil, fmt.Errorf("whitelisted asset cannot be nil")
	}
	if asset.Metadata == nil {
		return nil, fmt.Errorf("metadata cannot be nil")
	}

	// Step 1: Verify metadata hash
	if err := v.verifyMetadataHash(asset); err != nil {
		return nil, err
	}

	var rulesContainer *model.DecodedRulesContainer
	if len(cachedRulesContainer) > 0 && cachedRulesContainer[0] != nil {
		// Steps 2-3 already done during cache building
		rulesContainer = cachedRulesContainer[0]
	} else {
		// Step 2: Verify rules container signatures
		if err := v.verifyRulesContainerSignatures(asset, userSignaturesDecoder); err != nil {
			return nil, err
		}

		// Step 3: Decode rules container
		var err error
		rulesContainer, err = v.decodeRulesContainer(asset, rulesContainerDecoder)
		if err != nil {
			return nil, err
		}
	}

	// Step 4: Verify hash coverage
	// verifiedHash may differ from asset.Metadata.Hash if a legacy hash format was matched
	verifiedHash, verifiedPayload, err := v.verifyHashInSignedHashes(asset)
	if err != nil {
		return nil, err
	}

	// Step 5: Verify whitelist signatures using the verified hash
	if err := v.verifyWhitelistSignatures(asset, rulesContainer, verifiedHash); err != nil {
		return nil, err
	}

	return &AssetVerificationResult{
		RulesContainer:  rulesContainer,
		VerifiedHash:    verifiedHash,
		VerifiedPayload: verifiedPayload,
	}, nil
}

// VerifyAndDecodeRulesContainer performs steps 2-3 for one container, so the service can
// build its per-page cache instead of re-verifying per row.
func (v *WhitelistedAssetVerifier) VerifyAndDecodeRulesContainer(
	rulesContainerBase64 string,
	rulesSignaturesBase64 string,
	rulesContainerDecoder func(base64Data string) (*model.DecodedRulesContainer, error),
	userSignaturesDecoder func(base64Data string) ([]*model.RuleUserSignature, error),
) (*model.DecodedRulesContainer, error) {
	return VerifyAndDecodeRulesContainer(
		rulesContainerBase64,
		rulesSignaturesBase64,
		v.superAdminKeys,
		v.minValidSignatures,
		rulesContainerDecoder,
		userSignaturesDecoder,
	)
}

// verifyMetadataHash verifies that the computed hash matches the provided hash.
// Step 1 of the verification flow.
func (v *WhitelistedAssetVerifier) verifyMetadataHash(asset *model.WhitelistedAsset) error {
	if asset.Metadata.PayloadAsString == "" {
		return &model.IntegrityError{Message: "payloadAsString is empty"}
	}
	if asset.Metadata.Hash == "" {
		return &model.IntegrityError{Message: "metadata hash is empty"}
	}

	computedHash := crypto.CalculateHexHash(asset.Metadata.PayloadAsString)
	if !ConstantTimeCompare(computedHash, asset.Metadata.Hash) {
		return &model.IntegrityError{
			Message: "metadata hash verification failed",
		}
	}

	return nil
}

// verifyRulesContainerSignatures verifies SuperAdmin signatures on the rules container.
// Step 2 of the verification flow.
func (v *WhitelistedAssetVerifier) verifyRulesContainerSignatures(
	asset *model.WhitelistedAsset,
	userSignaturesDecoder func(base64Data string) ([]*model.RuleUserSignature, error),
) error {
	if len(v.superAdminKeys) == 0 {
		return &model.IntegrityError{Message: "no SuperAdmin keys configured for verification"}
	}

	if asset.RulesContainer == "" {
		return &model.IntegrityError{Message: "rulesContainer is empty"}
	}
	if asset.RulesSignatures == "" {
		return &model.IntegrityError{Message: "rulesSignatures is empty"}
	}

	// Decode rules signatures (protobuf UserSignatures)
	signatures, err := userSignaturesDecoder(asset.RulesSignatures)
	if err != nil {
		return &model.IntegrityError{
			Message: fmt.Sprintf("failed to decode rules signatures: %v", err),
		}
	}

	// Decode rules container data
	rulesData, err := DecodeBase64(asset.RulesContainer)
	if err != nil {
		return &model.IntegrityError{
			Message: fmt.Sprintf("failed to decode rules container: %v", err),
		}
	}

	// Verify signatures
	if err := VerifyGovernanceRulesSignatures(rulesData, signatures, v.superAdminKeys, v.minValidSignatures); err != nil {
		return &model.IntegrityError{
			Message: fmt.Sprintf("rules container signature verification failed: %v", err),
		}
	}

	return nil
}

// decodeRulesContainer decodes the base64 protobuf rules container.
// Step 3 of the verification flow.
func (v *WhitelistedAssetVerifier) decodeRulesContainer(
	asset *model.WhitelistedAsset,
	rulesContainerDecoder func(base64Data string) (*model.DecodedRulesContainer, error),
) (*model.DecodedRulesContainer, error) {
	if rulesContainerDecoder == nil {
		return nil, fmt.Errorf("rulesContainerDecoder is required")
	}

	container, err := rulesContainerDecoder(asset.RulesContainer)
	if err != nil {
		return nil, &model.IntegrityError{
			Message: fmt.Sprintf("failed to decode rules container: %v", err),
		}
	}

	return container, nil
}

// verifyHashInSignedHashes verifies that the metadata hash is covered by at least one signature.
// Step 4 of the verification flow.
//
// Returns BOTH the hash that was found (which may be a legacy hash) and the payload that hash
// covers. Step 6 runs in the service here rather than in this verifier, so the payload travels
// on AssetVerificationResult — see the address verifier for why parsing the delivered payload
// instead is a verification bypass.
func (v *WhitelistedAssetVerifier) verifyHashInSignedHashes(
	asset *model.WhitelistedAsset,
) (matchedHash string, matchedPayload string, err error) {
	if asset.SignedContractAddress == nil {
		return "", "", &model.IntegrityError{Message: "signedContractAddress is nil"}
	}

	signatures := asset.SignedContractAddress.Signatures
	if len(signatures) == 0 {
		return "", "", &model.IntegrityError{Message: "no signatures in signedContractAddress"}
	}

	// Try the provided hash first
	providedHash := asset.Metadata.Hash
	if VerifyHashCoverage(providedHash, signatures) {
		return providedHash, asset.Metadata.PayloadAsString, nil
	}

	// Try legacy variants for backward compatibility. This handles assets signed before schema
	// changes (e.g., before isNFT or kindType was added); the variant's PAYLOAD travels with
	// its hash so step 6 parses the bytes the signature actually covered.
	for _, variant := range ComputeAssetLegacyPayloadVariants(asset.Metadata.PayloadAsString) {
		if VerifyHashCoverage(variant.Hash, signatures) {
			return variant.Hash, variant.Payload, nil
		}
	}

	return "", "", &model.IntegrityError{
		Message: "metadata hash is not covered by any signature",
	}
}

// verifyWhitelistSignatures verifies user signatures meet governance threshold requirements.
// Step 5 of the verification flow.
func (v *WhitelistedAssetVerifier) verifyWhitelistSignatures(
	asset *model.WhitelistedAsset,
	rulesContainer *model.DecodedRulesContainer,
	metadataHash string,
) error {

	// Keyed off the SIGNED payload, not the response. See the address verifier.
	blockchain, network, networkFromPayload, err := resolveRuleKeyFor(
		asset.Metadata, asset.Blockchain, asset.Network)
	if err != nil {
		return err
	}

	// When the payload omits `network` the key is only partly signed, so enforce every tier
	// the unsigned DTO value could have selected rather than the one it named. See the address
	// verifier and model.FindContractAddressWhitelistingRuleCandidates.
	var applicableRules []*model.ContractAddressWhitelistingRules
	if networkFromPayload {
		if r := rulesContainer.FindContractAddressWhitelistingRules(blockchain, network); r != nil {
			applicableRules = []*model.ContractAddressWhitelistingRules{r}
		}
	} else {
		applicableRules = rulesContainer.FindContractAddressWhitelistingRuleCandidates(blockchain)
	}
	if len(applicableRules) == 0 {
		return &model.WhitelistError{
			Message: fmt.Sprintf("no contract address whitelisting rules found for blockchain=%s network=%s",
				blockchain, network),
		}
	}

	for _, whitelistRules := range applicableRules {
		// Contract whitelisting uses parallelThresholds directly (no rule lines matching)
		parallelThresholds := whitelistRules.ParallelThresholds
		if len(parallelThresholds) == 0 {
			return &model.WhitelistError{Message: "no threshold rules defined"}
		}

		// Try to verify all paths (OR logic - only one needs to succeed)
		pathFailures := tryVerifyAllPaths(parallelThresholds, rulesContainer, asset.SignedContractAddress.Signatures, metadataHash)
		if len(pathFailures) > 0 {
			scope := fmt.Sprintf("blockchain=%s network=%s", whitelistRules.Blockchain, whitelistRules.Network)
			if !networkFromPayload && len(applicableRules) > 1 {
				scope += " (enforced because the signed payload carries no network, so the " +
					"response could otherwise choose which quorum applies)"
			}
			return &model.WhitelistError{
				Message: fmt.Sprintf("signature verification failed for whitelisted asset (ID: %s) "+
					"against %s: no approval path satisfied the threshold requirements. %s",
					asset.ID, scope, strings.Join(pathFailures, "; ")),
			}
		}
	}

	return nil
}

// The step-5 threshold walk lives in group_threshold.go, shared with the address verifier.

// VerifyAsset is the witness-returning form of VerifyWhitelistedAsset.
//
// Prefer it in new code: a caller holding a VerifiedAsset cannot have skipped
// verification, whereas one holding a *model.WhitelistedAsset might have. The
// *AssetVerificationResult form is kept because existing callers return it.
func (v *WhitelistedAssetVerifier) VerifyAsset(
	asset *model.WhitelistedAsset,
	rulesContainerDecoder func(base64Data string) (*model.DecodedRulesContainer, error),
	userSignaturesDecoder func(base64Data string) ([]*model.RuleUserSignature, error),
) (VerifiedAsset, error) {
	result, err := v.VerifyWhitelistedAsset(asset, rulesContainerDecoder, userSignaturesDecoder)
	if err != nil {
		return VerifiedAsset{}, err
	}
	return VerifiedAsset{
		asset: asset,
		rules: result.RulesContainer,
		hash:  result.VerifiedHash,
	}, nil
}
