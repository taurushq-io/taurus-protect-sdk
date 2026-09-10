package helper

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"os"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
)

// Regenerates scripts/resources/verification-signed-fixtures.json.
//
//	UPDATE_VERIFICATION_SIGNED_FIXTURES=1 go test ./pkg/protect/helper -run TestUpdateSignedFixtures
//
// Go is the sanctioned producer, as it is for the cell vectors. Fresh keys are
// generated on every run: the fixture records PUBLIC keys and signatures only, so
// there is no private key material to keep stable — or to commit. Regeneration
// therefore yields a different but equally valid fixture set, and the other three
// SDKs must still agree with every recorded outcome.
//
// Skipped unless the env var is set, so a normal run never rewrites the gate it is
// supposed to be checked against.
func TestUpdateSignedFixtures(t *testing.T) {
	if os.Getenv("UPDATE_VERIFICATION_SIGNED_FIXTURES") == "" {
		t.Skip("set UPDATE_VERIFICATION_SIGNED_FIXTURES=1 to regenerate")
	}

	// The bytes SuperAdmins sign. Content is irrelevant to the threshold logic; what
	// matters is that every signature below is over exactly these bytes.
	container := []byte(`{"rules":"fixture container"}`)

	keyA := mustKey(t)
	keyB := mustKey(t)
	keyC := mustKey(t) // configured nowhere: an unknown signer

	// keyC's PEM is deliberately never recorded: it is the signer configured NOWHERE,
	// so the vectors that use its signature must fail on that basis alone.
	pemA, pemB := mustPEM(t, keyA), mustPEM(t, keyB)

	// Two separate signatures from the SAME key. ECDSA is randomised, so these differ
	// byte-for-byte while attesting the same thing — which is the whole reason the
	// threshold counts distinct KEYS rather than entries.
	sigA1 := mustSign(t, keyA, container)
	sigA2 := mustSign(t, keyA, container)
	sigB1 := mustSign(t, keyB, container)
	sigC1 := mustSign(t, keyC, container)

	type sig struct {
		UserID    string `json:"user_id"`
		Signature string `json:"signature"`
	}
	type vector struct {
		Description        string   `json:"description"`
		Signatures         []sig    `json:"signatures"`
		SuperAdminKeysPEM  []string `json:"super_admin_keys_pem"`
		MinValidSignatures int      `json:"min_valid_signatures"`
		Expect             string   `json:"expect"`
	}

	vectors := []vector{
		{
			Description:        "two distinct keys satisfy a threshold of two",
			Signatures:         []sig{{"sa-a", sigA1}, {"sa-b", sigB1}},
			SuperAdminKeysPEM:  []string{pemA, pemB},
			MinValidSignatures: 2,
			Expect:             "ok",
		},
		{
			Description:        "two entries from ONE key do not satisfy a threshold of two",
			Signatures:         []sig{{"sa-a", sigA1}, {"sa-a-again", sigA2}},
			SuperAdminKeysPEM:  []string{pemA, pemB},
			MinValidSignatures: 2,
			Expect:             "error",
		},
		{
			Description:        "two entries from ONE key do satisfy a threshold of one",
			Signatures:         []sig{{"sa-a", sigA1}, {"sa-a-again", sigA2}},
			SuperAdminKeysPEM:  []string{pemA, pemB},
			MinValidSignatures: 1,
			Expect:             "ok",
		},
		{
			Description:        "the identical signature repeated counts once",
			Signatures:         []sig{{"sa-a", sigA1}, {"sa-a-copy", sigA1}},
			SuperAdminKeysPEM:  []string{pemA, pemB},
			MinValidSignatures: 2,
			Expect:             "error",
		},
		{
			Description:        "the same key configured twice counts once",
			Signatures:         []sig{{"sa-a", sigA1}, {"sa-b", sigB1}},
			SuperAdminKeysPEM:  []string{pemA, pemA, pemB},
			MinValidSignatures: 3,
			Expect:             "error",
		},
		{
			Description:        "a signature from an unconfigured key contributes nothing",
			Signatures:         []sig{{"sa-a", sigA1}, {"stranger", sigC1}},
			SuperAdminKeysPEM:  []string{pemA, pemB},
			MinValidSignatures: 2,
			Expect:             "error",
		},
		{
			Description:        "an unconfigured signer alongside enough real ones still passes",
			Signatures:         []sig{{"sa-a", sigA1}, {"sa-b", sigB1}, {"stranger", sigC1}},
			SuperAdminKeysPEM:  []string{pemA, pemB},
			MinValidSignatures: 2,
			Expect:             "ok",
		},
		{
			Description:        "a signature over different bytes does not count",
			Signatures:         []sig{{"sa-a", mustSign(t, keyA, []byte("other bytes"))}, {"sa-b", sigB1}},
			SuperAdminKeysPEM:  []string{pemA, pemB},
			MinValidSignatures: 2,
			Expect:             "error",
		},
		{
			Description:        "an empty signature entry is skipped, not an error on its own",
			Signatures:         []sig{{"sa-a", sigA1}, {"blank", ""}, {"sa-b", sigB1}},
			SuperAdminKeysPEM:  []string{pemA, pemB},
			MinValidSignatures: 2,
			Expect:             "ok",
		},
		{
			Description:        "a non-positive threshold is rejected outright",
			Signatures:         []sig{{"sa-a", sigA1}},
			SuperAdminKeysPEM:  []string{pemA},
			MinValidSignatures: 0,
			Expect:             "error",
		},
	}

	groupVectors := buildGroupThresholdVectors(t)

	doc := map[string]interface{}{
		"_comment": []string{
			"Signed fixtures for the two verification thresholds.",
			"",
			"PUBLIC keys only. The gate verifies, it never signs, so no private key",
			"material is needed — or committed. Generated by the Go SDK; regenerate with",
			"UPDATE_VERIFICATION_SIGNED_FIXTURES=1 go test ./pkg/protect/helper -run TestUpdateSignedFixtures",
			"",
			"superadmin_threshold covers step 2: minValidSignatures, tenant-wide.",
			"group_threshold covers step 5: GroupThreshold.minimumSignatures, per group.",
			"Different scopes, same counting rule — DISTINCT SIGNING KEYS, never signature",
			"entries and never the server-supplied userId. ECDSA is randomised, so one key",
			"can emit unlimited valid signatures over the same bytes; counting entries would",
			"make 2-of-N no stronger than 1-of-N. These vectors are what stops that",
			"regressing in any one SDK.",
			"",
			"A group signature is over JSON.stringify(hashes) — the array recorded on the",
			"entry, not the metadata hash alone.",
		},
		"rules_container_base64": encodeStd(container),
		"counts": map[string]int{
			"superadmin_threshold": len(vectors),
			"group_threshold":      len(groupVectors),
		},
		"superadmin_threshold": vectors,
		"group_threshold":      groupVectors,
	}

	out, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(signedFixturesRelPath, append(out, '\n'), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	t.Logf("wrote %d superadmin + %d group vectors to %s", len(vectors), len(groupVectors), signedFixturesRelPath)
}

// Step-5 vectors. The signature is over the JSON-encoded hashes array, and the group's
// members plus their public keys come from the SuperAdmin-verified container — which is
// why an entry's userId can never be the identity that gets counted.
func buildGroupThresholdVectors(t *testing.T) []groupVector {
	t.Helper()

	const metadataHash = "1f2e3d4c5b6a79889796a5b4c3d2e1f001122334455667788990aabbccddeeff"
	hashes := []string{metadataHash}
	signed := mustHashesJSON(t, hashes)

	keyOne := mustKey(t)
	keyTwo := mustKey(t)
	outsider := mustKey(t)

	pemOne, pemTwo, pemOutsider := mustPEM(t, keyOne), mustPEM(t, keyTwo), mustPEM(t, outsider)

	// Two signatures from keyOne over identical bytes. Randomised ECDSA, so they differ
	// byte-for-byte while attesting the same thing.
	sigOneA := mustSign(t, keyOne, signed)
	sigOneB := mustSign(t, keyOne, signed)
	sigTwo := mustSign(t, keyTwo, signed)
	sigOutsider := mustSign(t, outsider, signed)

	members := []groupUser{{"member-one", pemOne}, {"member-two", pemTwo}}
	memberIDs := []string{"member-one", "member-two"}

	return []groupVector{
		{
			Description:       "two distinct members satisfy a 2-of-N group",
			MetadataHash:      metadataHash,
			GroupID:           "approvers",
			MinimumSignatures: 2,
			Users:             members,
			GroupUserIDs:      memberIDs,
			Signatures:        []groupSig{{"member-one", sigOneA, hashes}, {"member-two", sigTwo, hashes}},
			Expect:            "ok",
		},
		{
			// The bug this file exists to pin: copy the entry one member really signed and
			// entry counting reports a 2-of-N group satisfied by one approver.
			Description:       "a duplicated entry from ONE member does not satisfy a 2-of-N group",
			MetadataHash:      metadataHash,
			GroupID:           "approvers",
			MinimumSignatures: 2,
			Users:             members,
			GroupUserIDs:      memberIDs,
			Signatures:        []groupSig{{"member-one", sigOneA, hashes}, {"member-one", sigOneA, hashes}},
			Expect:            "error",
		},
		{
			Description:       "two re-signed entries from ONE member do not satisfy a 2-of-N group",
			MetadataHash:      metadataHash,
			GroupID:           "approvers",
			MinimumSignatures: 2,
			Users:             members,
			GroupUserIDs:      memberIDs,
			Signatures:        []groupSig{{"member-one", sigOneA, hashes}, {"member-one", sigOneB, hashes}},
			Expect:            "error",
		},
		{
			// Same entries as above at a threshold of one: the counting rule must not
			// reject what really is satisfied.
			Description:       "one member satisfies a 1-of-N group",
			MetadataHash:      metadataHash,
			GroupID:           "approvers",
			MinimumSignatures: 1,
			Users:             members,
			GroupUserIDs:      memberIDs,
			Signatures:        []groupSig{{"member-one", sigOneA, hashes}, {"member-one", sigOneB, hashes}},
			Expect:            "ok",
		},
		{
			// One compromised secret, two identities. Deduping on userId would count this
			// as two approvers.
			Description:       "two member IDs sharing ONE key count once",
			MetadataHash:      metadataHash,
			GroupID:           "approvers",
			MinimumSignatures: 2,
			Users:             []groupUser{{"member-one", pemOne}, {"member-alias", pemOne}},
			GroupUserIDs:      []string{"member-one", "member-alias"},
			Signatures:        []groupSig{{"member-one", sigOneA, hashes}, {"member-alias", sigOneB, hashes}},
			Expect:            "error",
		},
		{
			Description:       "a signer outside the group contributes nothing",
			MetadataHash:      metadataHash,
			GroupID:           "approvers",
			MinimumSignatures: 2,
			Users:             append(append([]groupUser{}, members...), groupUser{"stranger", pemOutsider}),
			GroupUserIDs:      memberIDs,
			Signatures:        []groupSig{{"member-one", sigOneA, hashes}, {"stranger", sigOutsider, hashes}},
			Expect:            "error",
		},
		{
			// Step 4 rejects it before the signer is counted, so the covered-hashes list
			// is part of what the threshold depends on.
			Description:       "an entry that does not cover the metadata hash is not counted",
			MetadataHash:      metadataHash,
			GroupID:           "approvers",
			MinimumSignatures: 2,
			Users:             members,
			GroupUserIDs:      memberIDs,
			Signatures: []groupSig{
				{"member-one", sigOneA, hashes},
				{"member-two", mustSign(t, keyTwo, mustHashesJSON(t, []string{"deadbeef"})), []string{"deadbeef"}},
			},
			Expect: "error",
		},
		{
			// No post-loop threshold check exists in any SDK, so a zero would mean "one
			// signature suffices" and quietly turn a 2-of-N group into 1-of-N.
			Description:       "a populated group requiring zero signatures is malformed",
			MetadataHash:      metadataHash,
			GroupID:           "approvers",
			MinimumSignatures: 0,
			Users:             members,
			GroupUserIDs:      memberIDs,
			Signatures:        []groupSig{{"member-one", sigOneA, hashes}},
			Expect:            "error",
		},
	}
}

type groupUser struct {
	UserID       string `json:"user_id"`
	PublicKeyPEM string `json:"public_key_pem"`
}

type groupSig struct {
	UserID    string   `json:"user_id"`
	Signature string   `json:"signature"`
	Hashes    []string `json:"hashes"`
}

type groupVector struct {
	Description       string      `json:"description"`
	MetadataHash      string      `json:"metadata_hash"`
	GroupID           string      `json:"group_id"`
	MinimumSignatures int         `json:"minimum_signatures"`
	Users             []groupUser `json:"users"`
	GroupUserIDs      []string    `json:"group_user_ids"`
	Signatures        []groupSig  `json:"signatures"`
	Expect            string      `json:"expect"`
}

func mustHashesJSON(t *testing.T, hashes []string) []byte {
	t.Helper()
	b, err := json.Marshal(hashes)
	if err != nil {
		t.Fatalf("marshal hashes: %v", err)
	}
	return b
}

func mustKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	return k
}

func mustPEM(t *testing.T, k *ecdsa.PrivateKey) string {
	t.Helper()
	der, err := x509.MarshalPKIXPublicKey(&k.PublicKey)
	if err != nil {
		t.Fatalf("marshal public key: %v", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}))
}

func mustSign(t *testing.T, k *ecdsa.PrivateKey, data []byte) string {
	t.Helper()
	sig, err := crypto.SignData(k, data)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	return sig
}

func encodeStd(b []byte) string { return base64.StdEncoding.EncodeToString(b) }
