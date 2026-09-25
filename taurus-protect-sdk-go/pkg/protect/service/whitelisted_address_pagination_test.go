package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	pb "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/proto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// signedWhitelist builds whitelisted-address rows that pass the full verification chain: a
// SuperAdmin-signed rules container with an ALGO/mainnet rule needing one team1 signature, and
// a team1 signature over each row's metadata hash.
type signedWhitelist struct {
	superAdmin *ecdsa.PrivateKey
	team       *ecdsa.PrivateKey
	container  string
	signatures string
}

func newSignedWhitelist(t *testing.T) *signedWhitelist {
	t.Helper()
	newKey := func() *ecdsa.PrivateKey {
		k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		return k
	}
	f := &signedWhitelist{superAdmin: newKey(), team: newKey()}

	der, err := x509.MarshalPKIXPublicKey(&f.team.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	teamPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}))

	container, err := mapper.RulesContainerToBase64(&model.DecodedRulesContainer{
		Users:  []*model.RuleUser{{ID: "team1@bank.com", PublicKeyPEM: teamPEM, Roles: []string{"WHITELISTEDADDRESSAPPROVER"}}},
		Groups: []*model.RuleGroup{{ID: "team1", UserIDs: []string{"team1@bank.com"}}},
		AddressWhitelistingRules: []*model.AddressWhitelistingRules{{
			Currency: "ALGO",
			Network:  "mainnet",
			ParallelThresholds: []*model.SequentialThresholds{{
				Thresholds: []*model.GroupThreshold{{GroupID: "team1", MinimumSignatures: 1}},
			}},
		}},
	})
	if err != nil {
		t.Fatalf("encoding the rules container: %v", err)
	}
	raw, err := base64.StdEncoding.DecodeString(container)
	if err != nil {
		t.Fatal(err)
	}
	signature, err := crypto.SignData(f.superAdmin, raw)
	if err != nil {
		t.Fatal(err)
	}
	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		t.Fatal(err)
	}
	signatures, err := proto.Marshal(&pb.UserSignatures{Signatures: []*pb.UserSignature{
		{UserId: "superadmin@bank.com", Signature: sigBytes},
	}})
	if err != nil {
		t.Fatal(err)
	}
	f.container = container
	f.signatures = base64.StdEncoding.EncodeToString(signatures)
	return f
}

// row is a signed envelope; tampered rows carry a payload their hash does not cover, so the SDK
// excludes them.
func (f *signedWhitelist) row(t *testing.T, id string, tampered bool) map[string]any {
	t.Helper()
	payload := fmt.Sprintf(`{"currency":"ALGO","network":"mainnet","addressType":"individual",`+
		`"address":"ADDR%s","memo":"","label":"row %s","customerId":""}`, id, id)
	hash := crypto.CalculateHexHash(payload)
	hashes, err := json.Marshal([]string{hash})
	if err != nil {
		t.Fatal(err)
	}
	signature, err := crypto.SignData(f.team, hashes)
	if err != nil {
		t.Fatal(err)
	}
	if tampered {
		payload = strings.Replace(payload, "ADDR", "EVIL", 1)
	}
	return map[string]any{
		"id":              id,
		"blockchain":      "ALGO",
		"network":         "mainnet",
		"metadata":        map[string]any{"hash": hash, "payloadAsString": payload},
		"rulesContainer":  f.container,
		"rulesSignatures": f.signatures,
		"signedAddress": map[string]any{"signatures": []any{map[string]any{
			"signature": map[string]any{"userId": "team1@bank.com", "signature": signature},
			"hashes":    []string{hash},
		}}},
	}
}

func (f *signedWhitelist) service(client *openapi.APIClient) *WhitelistedAddressService {
	return NewWhitelistedAddressServiceWithVerification(client, &WhitelistedAddressServiceConfig{
		SuperAdminKeys:     []*ecdsa.PublicKey{&f.superAdmin.PublicKey},
		MinValidSignatures: 1,
	})
}

// The fixture must produce rows the SDK accepts, and reject a tampered one: otherwise the walk
// below would pass for the wrong reason.
func TestSignedWhitelistFixtureVerifies(t *testing.T) {
	f := newSignedWhitelist(t)
	body, _ := json.Marshal(map[string]any{
		"result": []any{f.row(t, "1", false), f.row(t, "2", true)}, "totalItems": "2",
	})
	client := walkClient(t, func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write(body) })

	result, err := f.service(client).ListWhitelistedAddresses(context.Background(), nil)
	if err != nil {
		t.Fatalf("a signed row must verify: %v", err)
	}
	if len(result.Addresses) != 1 || result.Addresses[0].ID != "1" || result.Addresses[0].Address != "ADDR1" {
		t.Fatalf("addresses = %+v, want the signed row only", result.Addresses)
	}
	if len(result.ExcludedUnverified) != 1 || result.ExcludedUnverified[0].ID != "2" {
		t.Fatalf("excluded = %+v, want the tampered row", result.ExcludedUnverified)
	}
}

// Rows the SDK excludes reduce TotalItems but never shift the walk: the next page starts after
// every row the SERVER returned. Advancing by the kept rows instead re-reads excluded rows on the
// next page and repeats the rows after them.
func TestWhitelistedAddressesWalkIsNotShiftedByExclusions(t *testing.T) {
	const n, size = 10, 4
	tampered := map[int]bool{2: true, 5: true, 6: true}
	f := newSignedWhitelist(t)
	rows := make([]map[string]any, n)
	for i := range rows {
		rows[i] = f.row(t, fmt.Sprint(i), tampered[i])
	}

	requests := 0
	client := walkClient(t, func(w http.ResponseWriter, r *http.Request) {
		requests++
		if got := queryInt(t, r, "limit"); got != size {
			t.Errorf("request %d sent limit %d", requests, got)
		}
		offset := queryInt(t, r, "offset")
		page := rows[min(offset, n):min(offset+size, n)]
		body, err := json.Marshal(map[string]any{"result": page, "totalItems": fmt.Sprint(n)})
		if err != nil {
			t.Fatal(err)
		}
		_, _ = w.Write(body)
	})
	svc := f.service(client)

	var kept, excluded []string
	opts := &model.ListWhitelistedAddressesOptions{Limit: size}
	for {
		result, err := svc.ListWhitelistedAddresses(context.Background(), opts)
		if err != nil {
			t.Fatal(err)
		}
		for _, a := range result.Addresses {
			kept = append(kept, a.ID)
		}
		for _, e := range result.ExcludedUnverified {
			excluded = append(excluded, e.ID)
		}
		if want := n - len(result.ExcludedUnverified); result.Pagination.TotalItems != int64(want) {
			t.Errorf("TotalItems = %d, want the server total reduced by this page's exclusions (%d)",
				result.Pagination.TotalItems, want)
		}
		if !result.Pagination.HasMore {
			break
		}
		if requests > maxWalkRequests {
			t.Fatal("the walk does not end")
		}
		opts.Offset = result.Pagination.NextOffset
	}

	if strings.Join(kept, ",") != "0,1,3,4,7,8,9" || strings.Join(excluded, ",") != "2,5,6" {
		t.Errorf("kept %v excluded %v: every row must be seen exactly once", kept, excluded)
	}
	if requests != 3 {
		t.Errorf("sent %d requests, want 3", requests)
	}
}
