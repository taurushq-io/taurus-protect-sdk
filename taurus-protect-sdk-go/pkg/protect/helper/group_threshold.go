package helper

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// Step 5 for both whitelist flows. The ONE place the per-group threshold is evaluated —
// shared because none of it reads entity state. Do not add a per-verifier copy.
//
// Selecting WHICH thresholds apply stays per-entity: addresses resolve rule lines,
// contract rules go straight to ParallelThresholds.
//
//	parallelThresholds ──┬─▶ path 1 ─▶ ALL group thresholds met? ─┬─ yes ─▶ PASS (OR: first wins)
//	                     ├─▶ path 2 ─▶ ...                        └─ no ──▶ record failure
//	                     └─▶ path N                                          │
//	                                                                         ▼
//	                                            every path failed ─▶ WhitelistError(all reasons)

// precomputeHashesJSON pre-computes JSON serialization of each signature's hashes array.
// This avoids redundant json.Marshal calls when the same signature is checked
// across multiple group thresholds in the verification loops.
// Returns a map from signature index to marshaled JSON bytes.
// If marshaling fails for a signature, that index is absent from the map.
func precomputeHashesJSON(signatures []model.WhitelistSignature) map[int][]byte {
	result := make(map[int][]byte, len(signatures))
	for i, sig := range signatures {
		hashesJSON, err := json.Marshal(sig.Hashes)
		if err == nil {
			result[i] = hashesJSON
		}
	}
	return result
}

// tryVerifyAllPaths tries to verify all parallel threshold paths.
// Returns empty slice if verification passed, or list of failure messages if all paths failed.
func tryVerifyAllPaths(
	parallelThresholds []*model.SequentialThresholds,
	rulesContainer *model.DecodedRulesContainer,
	signatures []model.WhitelistSignature,
	metadataHash string,
) []string {
	// Pre-compute JSON serialization of each signature's hashes array once,
	// so it can be reused across all group threshold checks.
	hashesJSONMap := precomputeHashesJSON(signatures)

	var pathFailures []string

	for i, seqThreshold := range parallelThresholds {
		err := verifySequentialThresholds(seqThreshold, rulesContainer, signatures, metadataHash, hashesJSONMap)
		if err == nil {
			return nil // Verification passed
		}
		pathFailures = append(pathFailures, fmt.Sprintf("Path %d: %s", i+1, err.Error()))
	}

	return pathFailures
}

// verifySequentialThresholds verifies all group thresholds in a sequential threshold path.
func verifySequentialThresholds(
	seqThreshold *model.SequentialThresholds,
	rulesContainer *model.DecodedRulesContainer,
	signatures []model.WhitelistSignature,
	metadataHash string,
	hashesJSONMap map[int][]byte,
) error {
	if seqThreshold == nil || len(seqThreshold.Thresholds) == 0 {
		return &model.IntegrityError{Message: "no group thresholds defined"}
	}

	// ALL group thresholds must be satisfied (AND logic)
	for _, groupThreshold := range seqThreshold.Thresholds {
		if err := verifyGroupThreshold(groupThreshold, rulesContainer, signatures, metadataHash, hashesJSONMap); err != nil {
			return err
		}
	}

	return nil
}

// verifyGroupThreshold verifies that a group threshold is met.
//
// The unit is a SIGNER, not an entry. Every skip below is silent by design — a page
// carries the whole batch's approvals, so most entries belong to other groups — which
// is why the failure message carries the skip reasons.
//
//	entry ─▶ in this group? ──no──▶ skip (not an error)
//	            │yes
//	            ▼
//	         covers the metadata hash? ──no──▶ skip, reason recorded
//	            │yes
//	            ▼
//	         user + key in the container? ──no──▶ skip, reason recorded
//	            │yes
//	            ▼
//	         signature over hashes-JSON valid? ──no──▶ skip, reason recorded
//	            │yes
//	            ▼
//	         signers[fingerprint(key)] ──▶ |signers| >= minSigs? ──yes──▶ PASS
//	                                            │no
//	                                            ▼
//	                                     next entry; exhausted ─▶ IntegrityError
func verifyGroupThreshold(
	groupThreshold *model.GroupThreshold,
	rulesContainer *model.DecodedRulesContainer,
	signatures []model.WhitelistSignature,
	metadataHash string,
	hashesJSONMap map[int][]byte,
) error {
	groupID := groupThreshold.GroupID
	minSigs := groupThreshold.MinimumSignatures

	group := rulesContainer.FindGroupByID(groupID)
	if group == nil {
		return &model.IntegrityError{
			Message: fmt.Sprintf("group '%s' not found in rules container", groupID),
		}
	}

	if len(group.UserIDs) == 0 {
		if minSigs > 0 {
			return &model.IntegrityError{
				Message: fmt.Sprintf("group '%s' has no users but requires %d signature(s)", groupID, minSigs),
			}
		}
		return nil // minSignatures == 0, so empty group is OK
	}

	// A populated group with a zero threshold is a malformed container, not a group
	// that anyone may satisfy: the only success exit below is inside the loop, once the
	// signer set reaches the threshold, so a zero would mean "one signature suffices"
	// and turn a 2-of-N group into 1-of-N. Fail closed rather than guess.
	if minSigs <= 0 {
		return &model.IntegrityError{
			Message: fmt.Sprintf("group '%s' has %d user(s) but requires 0 signature(s): "+
				"minimumSignatures must be positive", groupID, len(group.UserIDs)),
		}
	}

	// Build set for faster lookup
	groupUserIDSet := make(map[string]bool)
	for _, uid := range group.UserIDs {
		groupUserIDSet[uid] = true
	}

	// Count DISTINCT signers, not signature entries.
	//
	// The entries come from the server-supplied userSignatures blob, so counting them
	// let a duplicated entry from one group member satisfy an N-of-M group: copy the
	// entry a single member really signed and an under-approved whitelist entry reads
	// as approved. Keyed on the container-resolved public key rather than the
	// server-supplied userId, so one compromised key shared by two IDs counts once —
	// the same rule the SuperAdmin threshold already applies.
	signers := make(map[string]struct{}, len(signatures))
	var skippedReasons []string

	for i, sig := range signatures {
		if sig.UserSignature == nil {
			skippedReasons = append(skippedReasons, "signature has nil userSig")
			continue
		}

		sigUserID := sig.UserSignature.UserID
		if !groupUserIDSet[sigUserID] {
			continue // Signer not in this group - not an error, just not relevant
		}

		// Check that metadata hash is covered by this signature
		if !containsHash(sig.Hashes, metadataHash) {
			skippedReasons = append(skippedReasons, fmt.Sprintf(
				"user '%s' signature does not cover metadata hash '%s' (signed hashes=%v)",
				sigUserID, metadataHash, sig.Hashes))
			continue
		}

		user := rulesContainer.FindUserByID(sigUserID)
		if user == nil {
			skippedReasons = append(skippedReasons, fmt.Sprintf("user '%s' not found in rules container", sigUserID))
			continue
		}
		if user.PublicKey == nil {
			skippedReasons = append(skippedReasons, fmt.Sprintf("user '%s' has no public key", sigUserID))
			continue
		}

		// Use pre-computed JSON-encoded hashes array
		hashesJSON, ok := hashesJSONMap[i]
		if !ok {
			skippedReasons = append(skippedReasons, fmt.Sprintf("failed to marshal hashes for user '%s'", sigUserID))
			continue
		}

		valid, err := crypto.VerifySignature(user.PublicKey, hashesJSON, sig.UserSignature.Signature)
		if err != nil || !valid {
			skippedReasons = append(skippedReasons, fmt.Sprintf("user '%s' signature verification failed", sigUserID))
			continue
		}

		fingerprint, err := user.KeyFingerprint()
		if err != nil {
			// An unfingerprintable key cannot be counted as a distinct signer.
			skippedReasons = append(skippedReasons, fmt.Sprintf("user '%s' key cannot be fingerprinted", sigUserID))
			continue
		}
		signers[fingerprint] = struct{}{}
		if len(signers) >= minSigs {
			return nil // Threshold met
		}
	}

	// Threshold not met
	message := fmt.Sprintf("group '%s' requires %d distinct signer(s) but only %d valid",
		groupID, minSigs, len(signers))
	if len(skippedReasons) > 0 {
		message += " [" + strings.Join(skippedReasons, "; ") + "]"
	}
	return &model.IntegrityError{Message: message}
}
