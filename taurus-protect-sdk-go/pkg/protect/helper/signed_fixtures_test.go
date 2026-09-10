package helper

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// Cross-SDK signed fixtures for verification step 2 — the SuperAdmin threshold.
//
// The behaviour vectors cover everything expressible without key material; this file
// covers what needs real signatures. Until it existed, the repo's own CLAUDE.md
// recorded that NO automated cross-SDK gate covered the distinct-key rule, and the
// five cases below lived as five hand-maintained copies in four suites.
//
// The fixture carries PUBLIC keys and signatures only — the gate verifies, it never
// signs — so there is no private key material in the repo.
const signedFixturesRelPath = "../../../../scripts/resources/verification-signed-fixtures.json"

type signedFixtureVector struct {
	Description string `json:"description"`
	Signatures  []struct {
		UserID    string `json:"user_id"`
		Signature string `json:"signature"`
	} `json:"signatures"`
	SuperAdminKeysPEM  []string `json:"super_admin_keys_pem"`
	MinValidSignatures int      `json:"min_valid_signatures"`
	Expect             string   `json:"expect"`
}

// Step-5 vector. The container's users and group membership are given, so the loader
// builds the same DecodedRulesContainer the verified container would have produced.
type groupFixtureVector struct {
	Description string `json:"description"`
	Users       []struct {
		UserID       string `json:"user_id"`
		PublicKeyPEM string `json:"public_key_pem"`
	} `json:"users"`
	Signatures []struct {
		UserID    string   `json:"user_id"`
		Signature string   `json:"signature"`
		Hashes    []string `json:"hashes"`
	} `json:"signatures"`
	MetadataHash      string   `json:"metadata_hash"`
	GroupID           string   `json:"group_id"`
	GroupUserIDs      []string `json:"group_user_ids"`
	MinimumSignatures int      `json:"minimum_signatures"`
	Expect            string   `json:"expect"`
}

type signedFixtures struct {
	RulesContainerBase64 string `json:"rules_container_base64"`
	Counts               struct {
		SuperAdminThreshold int `json:"superadmin_threshold"`
		GroupThreshold      int `json:"group_threshold"`
	} `json:"counts"`
	SuperAdminThreshold []signedFixtureVector `json:"superadmin_threshold"`
	GroupThreshold      []groupFixtureVector  `json:"group_threshold"`
}

func loadSignedFixtures(t *testing.T) signedFixtures {
	t.Helper()

	path := filepath.Clean(signedFixturesRelPath)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read shared signed fixtures %s: %v", path, err)
	}

	var f signedFixtures
	if err := json.Unmarshal(raw, &f); err != nil {
		t.Fatalf("cannot parse %s: %v", path, err)
	}
	if len(f.SuperAdminThreshold) != f.Counts.SuperAdminThreshold {
		t.Fatalf("superadmin_threshold: got %d vectors, file declares %d",
			len(f.SuperAdminThreshold), f.Counts.SuperAdminThreshold)
	}
	if len(f.GroupThreshold) != f.Counts.GroupThreshold {
		t.Fatalf("group_threshold: got %d vectors, file declares %d",
			len(f.GroupThreshold), f.Counts.GroupThreshold)
	}
	if f.RulesContainerBase64 == "" {
		t.Fatal("fixture carries no rules container to verify against")
	}
	return f
}

func TestSignedFixtures_SuperAdminThreshold(t *testing.T) {
	fixtures := loadSignedFixtures(t)

	container, err := base64.StdEncoding.DecodeString(fixtures.RulesContainerBase64)
	if err != nil {
		t.Fatalf("cannot decode the fixture rules container: %v", err)
	}

	for _, tc := range fixtures.SuperAdminThreshold {
		t.Run(tc.Description, func(t *testing.T) {
			keys := make([]*ecdsa.PublicKey, 0, len(tc.SuperAdminKeysPEM))
			for i, keyPEM := range tc.SuperAdminKeysPEM {
				key, err := crypto.DecodePublicKeyPEM(keyPEM)
				if err != nil {
					t.Fatalf("cannot decode fixture key %d: %v", i, err)
				}
				keys = append(keys, key)
			}

			signatures := make([]*model.RuleUserSignature, 0, len(tc.Signatures))
			for _, s := range tc.Signatures {
				signatures = append(signatures, &model.RuleUserSignature{
					UserID:    s.UserID,
					Signature: s.Signature,
				})
			}

			err := VerifyGovernanceRulesSignatures(container, signatures, keys, tc.MinValidSignatures)

			if tc.Expect == "error" {
				if err == nil {
					t.Fatal("expected verification to fail, it passed")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected verification to pass, got: %v", err)
			}
		})
	}
}

// The fixture must actually exercise the distinct-key rule: without a case where the
// signature COUNT meets the threshold but the distinct-key count does not, the whole
// file would pass against an implementation that counts entries.
func TestSignedFixtures_CoverTheDistinctKeyRule(t *testing.T) {
	fixtures := loadSignedFixtures(t)

	found := false
	for _, tc := range fixtures.SuperAdminThreshold {
		if tc.Expect == "error" && len(tc.Signatures) >= tc.MinValidSignatures && tc.MinValidSignatures > 0 {
			found = true
			break
		}
	}
	if !found {
		t.Error("no vector where entry count meets the threshold but distinct keys do not; " +
			"an implementation counting entries would pass this file")
	}
}

// Step 5 — the per-group threshold. A different threshold from the SuperAdmin one above
// (per group, not tenant-wide) but the same counting rule, so it is gated from the same
// file: the four SDKs cannot drift on one without drifting on the other.
func TestSignedFixtures_GroupThreshold(t *testing.T) {
	fixtures := loadSignedFixtures(t)

	for _, tc := range fixtures.GroupThreshold {
		t.Run(tc.Description, func(t *testing.T) {
			users := make([]*model.RuleUser, 0, len(tc.Users))
			for i, u := range tc.Users {
				key, err := crypto.DecodePublicKeyPEM(u.PublicKeyPEM)
				if err != nil {
					t.Fatalf("cannot decode fixture key for user %d: %v", i, err)
				}
				users = append(users, &model.RuleUser{ID: u.UserID, PublicKey: key})
			}

			container := &model.DecodedRulesContainer{
				Users:  users,
				Groups: []*model.RuleGroup{{ID: tc.GroupID, UserIDs: tc.GroupUserIDs}},
			}

			signatures := make([]model.WhitelistSignature, 0, len(tc.Signatures))
			for _, s := range tc.Signatures {
				signatures = append(signatures, model.WhitelistSignature{
					Hashes: s.Hashes,
					UserSignature: &model.WhitelistUserSignature{
						UserID:    s.UserID,
						Signature: s.Signature,
					},
				})
			}

			gt := &model.GroupThreshold{GroupID: tc.GroupID, MinimumSignatures: tc.MinimumSignatures}
			err := verifyGroupThreshold(gt, container, signatures, tc.MetadataHash,
				precomputeHashesJSON(signatures))

			if tc.Expect == "error" {
				if err == nil {
					t.Fatal("expected the threshold to fail, it passed")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected the threshold to pass, got: %v", err)
			}
		})
	}
}

// Same non-vacuity guard as the SuperAdmin section: without a vector where the entry
// count meets the threshold but the distinct-signer count does not, the group section
// would pass against an implementation that counts entries.
func TestSignedFixtures_GroupSectionCoversTheDistinctSignerRule(t *testing.T) {
	fixtures := loadSignedFixtures(t)

	found := false
	for _, tc := range fixtures.GroupThreshold {
		if tc.Expect != "error" || tc.MinimumSignatures <= 0 {
			continue
		}
		inGroup := 0
		members := make(map[string]bool, len(tc.GroupUserIDs))
		for _, id := range tc.GroupUserIDs {
			members[id] = true
		}
		for _, s := range tc.Signatures {
			if members[s.UserID] {
				inGroup++
			}
		}
		if inGroup >= tc.MinimumSignatures {
			found = true
			break
		}
	}
	if !found {
		t.Error("no group vector where in-group entry count meets the threshold but " +
			"distinct signers do not; an implementation counting entries would pass")
	}
}
