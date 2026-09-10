package helper

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

const groupTestHash = "abc123"

// signedEntry produces a genuinely valid signature over the hashes array, as a group
// member's approval arrives on the wire.
func signedEntry(t *testing.T, userID string, key *ecdsa.PrivateKey) model.WhitelistSignature {
	t.Helper()
	hashes := []string{groupTestHash}
	hashesJSON, err := json.Marshal(hashes)
	if err != nil {
		t.Fatalf("marshal hashes: %v", err)
	}
	signature, err := crypto.SignData(key, hashesJSON)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	return model.WhitelistSignature{
		Hashes: hashes,
		UserSignature: &model.WhitelistUserSignature{
			UserID:    userID,
			Signature: signature,
		},
	}
}

func twoMemberGroup(t *testing.T, minSigs int) (*model.DecodedRulesContainer, *ecdsa.PrivateKey, *ecdsa.PrivateKey, *model.GroupThreshold) {
	t.Helper()
	key1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key1: %v", err)
	}
	key2, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key2: %v", err)
	}

	container := &model.DecodedRulesContainer{
		Users: []*model.RuleUser{
			{ID: "user1@bank.com", PublicKey: &key1.PublicKey, Roles: []string{"USER"}},
			{ID: "user2@bank.com", PublicKey: &key2.PublicKey, Roles: []string{"USER"}},
		},
		Groups: []*model.RuleGroup{
			{ID: "approvers", UserIDs: []string{"user1@bank.com", "user2@bank.com"}},
		},
	}
	return container, key1, key2, &model.GroupThreshold{GroupID: "approvers", MinimumSignatures: minSigs}
}

// The entries come from the server-supplied userSignatures blob, so counting them let a
// duplicated entry from ONE member satisfy a 2-of-N group — promoting an under-approved
// whitelist entry to approved.
func TestGroupThreshold_DuplicateEntryFromOneMemberDoesNotMeetTwoOfN(t *testing.T) {
	container, key1, _, gt := twoMemberGroup(t, 2)

	entry := signedEntry(t, "user1@bank.com", key1)
	signatures := []model.WhitelistSignature{entry, entry}

	err := verifyGroupThreshold(gt, container, signatures, groupTestHash, precomputeHashesJSON(signatures))
	if err == nil {
		t.Fatal("one member's entry duplicated must NOT satisfy a 2-of-N group")
	}
}

// Re-signing rather than copying produces a different signature for the same key, since
// ECDSA is randomized. It must still count as one signer.
func TestGroupThreshold_ReSignedEntryFromOneMemberDoesNotMeetTwoOfN(t *testing.T) {
	container, key1, _, gt := twoMemberGroup(t, 2)

	signatures := []model.WhitelistSignature{
		signedEntry(t, "user1@bank.com", key1),
		signedEntry(t, "user1@bank.com", key1),
	}
	if signatures[0].UserSignature.Signature == signatures[1].UserSignature.Signature {
		t.Fatal("test setup: expected two distinct signatures from the same key")
	}

	err := verifyGroupThreshold(gt, container, signatures, groupTestHash, precomputeHashesJSON(signatures))
	if err == nil {
		t.Fatal("two signatures from ONE key must NOT satisfy a 2-of-N group")
	}
}

func TestGroupThreshold_TwoDistinctMembersMeetTwoOfN(t *testing.T) {
	container, key1, key2, gt := twoMemberGroup(t, 2)

	signatures := []model.WhitelistSignature{
		signedEntry(t, "user1@bank.com", key1),
		signedEntry(t, "user2@bank.com", key2),
	}

	if err := verifyGroupThreshold(gt, container, signatures, groupTestHash, precomputeHashesJSON(signatures)); err != nil {
		t.Fatalf("two distinct members must satisfy a 2-of-N group, got %v", err)
	}
}

func TestGroupThreshold_OneMemberMeetsOneOfN(t *testing.T) {
	container, key1, _, gt := twoMemberGroup(t, 1)

	signatures := []model.WhitelistSignature{signedEntry(t, "user1@bank.com", key1)}

	if err := verifyGroupThreshold(gt, container, signatures, groupTestHash, precomputeHashesJSON(signatures)); err != nil {
		t.Fatalf("one member must satisfy a 1-of-N group, got %v", err)
	}
}

// Two user IDs sharing one key is one compromised secret, so it must not satisfy 2-of-N.
// This is why the set keys on the key and not the server-supplied userId.
func TestGroupThreshold_TwoUserIDsSharingOneKeyCountOnce(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	container := &model.DecodedRulesContainer{
		Users: []*model.RuleUser{
			{ID: "user1@bank.com", PublicKey: &key.PublicKey, Roles: []string{"USER"}},
			{ID: "user2@bank.com", PublicKey: &key.PublicKey, Roles: []string{"USER"}},
		},
		Groups: []*model.RuleGroup{
			{ID: "approvers", UserIDs: []string{"user1@bank.com", "user2@bank.com"}},
		},
	}
	gt := &model.GroupThreshold{GroupID: "approvers", MinimumSignatures: 2}

	signatures := []model.WhitelistSignature{
		signedEntry(t, "user1@bank.com", key),
		signedEntry(t, "user2@bank.com", key),
	}

	err = verifyGroupThreshold(gt, container, signatures, groupTestHash, precomputeHashesJSON(signatures))
	if err == nil {
		t.Fatal("two user IDs backed by ONE key must NOT satisfy a 2-of-N group")
	}
}

func TestGroupThreshold_SignerOutsideGroupContributesNothing(t *testing.T) {
	container, key1, _, gt := twoMemberGroup(t, 2)

	outsiderKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate outsider key: %v", err)
	}
	container.Users = append(container.Users, &model.RuleUser{
		ID: "outsider@bank.com", PublicKey: &outsiderKey.PublicKey, Roles: []string{"USER"},
	})

	signatures := []model.WhitelistSignature{
		signedEntry(t, "user1@bank.com", key1),
		signedEntry(t, "outsider@bank.com", outsiderKey),
	}

	err = verifyGroupThreshold(gt, container, signatures, groupTestHash, precomputeHashesJSON(signatures))
	if err == nil {
		t.Fatal("a valid signer outside the group must not count toward its threshold")
	}
}
