package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/cache"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// The v2 asset-holders rows carry no signature, so QueryAssetAddresses returns an address as a
// Taurus-PROTECT address only after its verified counterpart confirms it.

const (
	holderInternal    = "ADDRESS_TYPE_V2_INTERNAL"
	holderWhitelisted = "ADDRESS_TYPE_V2_WHITELISTED"
	holderExternal    = "ADDRESS_TYPE_V2_EXTERNAL"
	holdersPath       = "/api/rest/v2/assets/a1/addresses/query"
	managedPath       = "/api/rest/v1/addresses"
	whitelistPath     = "/api/rest/v1/whitelists/addresses"
	holdersNextCursor = "bmV4dA=="
)

// holdersFixture fakes the holders page and the two verified readers behind it.
type holdersFixture struct {
	t         *testing.T
	hsm       *ecdsa.PrivateKey
	whitelist *signedWhitelist
	// holders are the rows of the holders page; the page always claims a next page.
	holders []map[string]any
	// managed and whitelisted are what the verified readers return, by id.
	managed     map[string]map[string]any
	whitelisted map[string]map[string]any
	// managedStatus, when set, is the status the managed-address read answers with.
	managedStatus int
	// noHSMKey serves a rules container without an HSMSLOT user.
	noHSMKey bool

	mu       sync.Mutex
	requests []recordedRequest
}

func newHoldersFixture(t *testing.T) *holdersFixture {
	t.Helper()
	hsm, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return &holdersFixture{
		t: t, hsm: hsm, whitelist: newSignedWhitelist(t),
		managed: map[string]map[string]any{}, whitelisted: map[string]map[string]any{},
	}
}

func holderRow(address, addressType, addressID, whitelistedID string) map[string]any {
	row := map[string]any{"address": address, "balance": "5", "kycStatus": "KYC_STATUS_V2_APPROVED"}
	if addressType != "" {
		row["addressType"] = addressType
	}
	if addressID != "" {
		row["addressID"] = addressID
	}
	if whitelistedID != "" {
		row["whitelistedAddressID"] = whitelistedID
	}
	return row
}

// manage makes the managed-address read return id with an HSM signature over signedAddress.
func (f *holdersFixture) manage(id, address, signedAddress string) {
	signature, err := crypto.SignData(f.hsm, []byte(signedAddress))
	if err != nil {
		f.t.Fatal(err)
	}
	f.managed[id] = map[string]any{"id": id, "address": address, "signature": signature, "status": "confirmed"}
}

// whitelist makes the whitelist read return a signed envelope for id, whose address is ADDR<id>.
func (f *holdersFixture) whitelistRow(id string, tampered bool) {
	f.whitelisted[id] = f.whitelist.row(f.t, id, tampered)
}

func (f *holdersFixture) serve(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	f.mu.Lock()
	f.requests = append(f.requests, recordedRequest{method: r.Method, path: r.URL.Path, query: r.URL.RawQuery, body: body})
	f.mu.Unlock()

	q := r.URL.Query()
	var reply any
	switch r.URL.Path {
	case holdersPath:
		reply = map[string]any{"result": f.holders,
			"cursor": map[string]any{"currentPage": holdersNextCursor, "hasNext": true}}
	case managedPath:
		if f.managedStatus != 0 {
			w.WriteHeader(f.managedStatus)
			_, _ = io.WriteString(w, `{"code":13,"message":"unavailable"}`)
			return
		}
		reply = f.page(f.managed, q["addressIds"])
	case whitelistPath:
		reply = f.page(f.whitelisted, q["ids"])
	default:
		f.t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	out, err := json.Marshal(reply)
	if err != nil {
		f.t.Fatal(err)
	}
	_, _ = w.Write(out)
}

func (f *holdersFixture) page(rows map[string]map[string]any, ids []string) map[string]any {
	result := []any{}
	for _, id := range ids {
		if row, ok := rows[id]; ok {
			result = append(result, row)
		}
	}
	return map[string]any{"result": result, "totalItems": fmt.Sprint(len(result))}
}

func (f *holdersFixture) service() *AssetService {
	client := walkClient(f.t, f.serve)
	container := &model.DecodedRulesContainer{}
	if !f.noHSMKey {
		container.Users = []*model.RuleUser{{ID: "hsm", PublicKey: &f.hsm.PublicKey, Roles: []string{"HSMSLOT"}}}
	}
	rules := cache.NewRulesContainerCache(time.Hour, func(context.Context) (*model.DecodedRulesContainer, error) {
		return container, nil
	})
	return NewAssetService(client, rules, NewAddressService(client, rules), f.whitelist.service(client))
}

func (f *holdersFixture) requestsTo(path string) []url.Values {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []url.Values
	for _, r := range f.requests {
		if r.path == path {
			q, _ := url.ParseQuery(r.query)
			out = append(out, q)
		}
	}
	return out
}

func (f *holdersFixture) query(t *testing.T) (*model.QueryAssetAddressesResult, error) {
	t.Helper()
	return f.service().QueryAssetAddresses(context.Background(), "a1", nil)
}

// resultIDs lists the kept rows as address:verified and the exclusions as id.
func resultIDs(result *model.QueryAssetAddressesResult) (kept, excluded []string) {
	for _, a := range result.Addresses {
		kept = append(kept, fmt.Sprintf("%s:%v", a.Address, a.Verified))
	}
	for _, e := range result.ExcludedUnverified {
		excluded = append(excluded, e.ID)
	}
	return kept, excluded
}

func TestQueryAssetAddressesConfirmsInternalRows(t *testing.T) {
	for _, tc := range []struct {
		name         string
		row          map[string]any
		setup        func(f *holdersFixture)
		wantExcluded string
		wantReason   string
		wantReadIDs  string
	}{
		{name: "confirmed row is kept", row: holderRow("party::8", holderInternal, "8", ""),
			setup:       func(f *holdersFixture) { f.manage("8", "party::8", "party::8") },
			wantReadIDs: "7,8"},
		{name: "address differs from the verified one", row: holderRow("party::8", holderInternal, "8", ""),
			setup:        func(f *holdersFixture) { f.manage("8", "party::other", "party::other") },
			wantExcluded: "8", wantReason: "differs", wantReadIDs: "7,8"},
		{name: "id not returned by the verified read", row: holderRow("party::9", holderInternal, "9", ""),
			wantExcluded: "9", wantReason: "did not return", wantReadIDs: "7,9"},
		{name: "row without an addressID", row: holderRow("party::10", holderInternal, "", ""),
			wantExcluded: "party::10", wantReason: "no addressID", wantReadIDs: "7"},
		{name: "row with a zero addressID", row: holderRow("party::11", holderInternal, "0", ""),
			wantExcluded: "party::11", wantReason: "no addressID", wantReadIDs: "7"},
		{name: "HSM signature does not verify", row: holderRow("party::12", holderInternal, "12", ""),
			setup:        func(f *holdersFixture) { f.manage("12", "party::12", "party::forged") },
			wantExcluded: "12", wantReason: "did not verify", wantReadIDs: "7,12"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newHoldersFixture(t)
			f.manage("7", "party::7", "party::7")
			f.holders = []map[string]any{holderRow("party::7", holderInternal, "7", ""), tc.row}
			if tc.setup != nil {
				tc.setup(f)
			}

			result, err := f.query(t)
			if err != nil {
				t.Fatal(err)
			}
			kept, excluded := resultIDs(result)
			wantKept := "party::7:true"
			if tc.wantExcluded == "" {
				wantKept = "party::7:true,party::8:true"
			}
			if strings.Join(kept, ",") != wantKept {
				t.Errorf("kept %v, want %s", kept, wantKept)
			}
			if strings.Join(excluded, ",") != tc.wantExcluded {
				t.Errorf("excluded %v, want %q", excluded, tc.wantExcluded)
			}
			if tc.wantReason != "" && !strings.Contains(result.ExcludedUnverified[0].Reason, tc.wantReason) {
				t.Errorf("reason %q does not say %q", result.ExcludedUnverified[0].Reason, tc.wantReason)
			}
			// Exclusions never move the cursor.
			if !result.Page.HasMore || result.Page.NextCursor != holdersNextCursor {
				t.Errorf("page = %+v, want the holders cursor untouched", result.Page)
			}
			reads := f.requestsTo(managedPath)
			if len(reads) != 1 || strings.Join(reads[0]["addressIds"], ",") != tc.wantReadIDs ||
				reads[0].Get("limit") != fmt.Sprint(len(reads[0]["addressIds"])) {
				t.Errorf("managed-address reads = %v, want one for ids %s", reads, tc.wantReadIDs)
			}
			if n := len(f.requestsTo(whitelistPath)); n != 0 {
				t.Errorf("sent %d whitelist reads for a page without WHITELISTED rows", n)
			}
		})
	}
}

func TestQueryAssetAddressesConfirmsWhitelistedRows(t *testing.T) {
	for _, tc := range []struct {
		name         string
		row          map[string]any
		setup        func(f *holdersFixture)
		wantExcluded string
		wantReason   string
		wantReadIDs  string
	}{
		{name: "confirmed row is kept", row: holderRow("ADDR2", holderWhitelisted, "", "2"),
			setup: func(f *holdersFixture) { f.whitelistRow("2", false) }, wantReadIDs: "1,2"},
		{name: "address differs from the verified one", row: holderRow("ADDR2x", holderWhitelisted, "", "2"),
			setup:        func(f *holdersFixture) { f.whitelistRow("2", false) },
			wantExcluded: "2", wantReason: "differs", wantReadIDs: "1,2"},
		{name: "id not returned by the verified read", row: holderRow("ADDR3", holderWhitelisted, "", "3"),
			wantExcluded: "3", wantReason: "did not return", wantReadIDs: "1,3"},
		{name: "row without a whitelistedAddressID", row: holderRow("ADDR4", holderWhitelisted, "", ""),
			wantExcluded: "ADDR4", wantReason: "no whitelistedAddressID", wantReadIDs: "1"},
		{name: "envelope does not verify", row: holderRow("ADDR5", holderWhitelisted, "", "5"),
			setup:        func(f *holdersFixture) { f.whitelistRow("5", true) },
			wantExcluded: "5", wantReason: "did not verify", wantReadIDs: "1,5"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newHoldersFixture(t)
			f.whitelistRow("1", false)
			f.holders = []map[string]any{holderRow("ADDR1", holderWhitelisted, "", "1"), tc.row}
			if tc.setup != nil {
				tc.setup(f)
			}

			result, err := f.query(t)
			if err != nil {
				t.Fatal(err)
			}
			kept, excluded := resultIDs(result)
			wantKept := "ADDR1:true"
			if tc.wantExcluded == "" {
				wantKept = "ADDR1:true,ADDR2:true"
			}
			if strings.Join(kept, ",") != wantKept {
				t.Errorf("kept %v, want %s", kept, wantKept)
			}
			if strings.Join(excluded, ",") != tc.wantExcluded {
				t.Errorf("excluded %v, want %q", excluded, tc.wantExcluded)
			}
			if tc.wantReason != "" && !strings.Contains(result.ExcludedUnverified[0].Reason, tc.wantReason) {
				t.Errorf("reason %q does not say %q", result.ExcludedUnverified[0].Reason, tc.wantReason)
			}
			reads := f.requestsTo(whitelistPath)
			if len(reads) != 1 || strings.Join(reads[0]["ids"], ",") != tc.wantReadIDs {
				t.Errorf("whitelist reads = %v, want one for ids %s", reads, tc.wantReadIDs)
			}
			// A whitelisted address still awaiting approval is not one a holder can be confirmed by.
			if len(reads) == 1 && reads[0].Has("includeForApproval") {
				t.Errorf("the holders re-read asked for rows awaiting approval: %v", reads[0])
			}
			if n := len(f.requestsTo(managedPath)); n != 0 {
				t.Errorf("sent %d managed-address reads for a page without INTERNAL rows", n)
			}
		})
	}
}

// EXTERNAL and untyped rows are on-chain data nothing signs: returned with Verified false, and
// no verified reader is asked — not even for an EXTERNAL row that carries a platform id.
func TestQueryAssetAddressesReturnsOtherRowsUnverified(t *testing.T) {
	f := newHoldersFixture(t)
	f.manage("7", "party::7", "party::7")
	f.holders = []map[string]any{
		holderRow("party::ext", holderExternal, "", ""),
		holderRow("party::untyped", "", "", ""),
		holderRow("party::7", holderExternal, "7", ""),
	}

	result, err := f.query(t)
	if err != nil {
		t.Fatal(err)
	}
	kept, excluded := resultIDs(result)
	if strings.Join(kept, ",") != "party::ext:false,party::untyped:false,party::7:false" || len(excluded) != 0 {
		t.Errorf("kept %v excluded %v, want all three rows unverified", kept, excluded)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.requests) != 1 {
		t.Errorf("sent %d requests, want only the holders page", len(f.requests))
	}
}

// The managed-address read takes at most 50 ids, the whitelist read at most 100.
func TestQueryAssetAddressesBatchesTheVerifiedReads(t *testing.T) {
	f := newHoldersFixture(t)
	for i := 1; i <= 51; i++ {
		id, address := fmt.Sprint(i), fmt.Sprintf("party::%d", i)
		f.manage(id, address, address)
		f.holders = append(f.holders, holderRow(address, holderInternal, id, ""))
	}
	for i := 1; i <= 101; i++ {
		id := fmt.Sprint(1000 + i)
		f.whitelistRow(id, false)
		f.holders = append(f.holders, holderRow("ADDR"+id, holderWhitelisted, "", id))
	}

	result, err := f.query(t)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Addresses) != 152 || len(result.ExcludedUnverified) != 0 {
		t.Fatalf("kept %d excluded %d, want 152 and 0", len(result.Addresses), len(result.ExcludedUnverified))
	}
	for _, a := range result.Addresses {
		if !a.Verified {
			t.Fatalf("%s is not verified", a.Address)
		}
	}
	for path, want := range map[string][]int{managedPath: {50, 1}, whitelistPath: {100, 1}} {
		idsParam := "addressIds"
		if path == whitelistPath {
			idsParam = "ids"
		}
		reads := f.requestsTo(path)
		if len(reads) != len(want) {
			t.Fatalf("%s: %d reads, want %d", path, len(reads), len(want))
		}
		for i, n := range want {
			if len(reads[i][idsParam]) != n || reads[i].Get("limit") != fmt.Sprint(n) {
				t.Errorf("%s read %d: %d ids, limit %s; want %d", path, i, len(reads[i][idsParam]), reads[i].Get("limit"), n)
			}
		}
	}
}

// Rows came back but none can be confirmed: an error, never an empty page.
func TestQueryAssetAddressesFailsWhenNoRowSurvives(t *testing.T) {
	f := newHoldersFixture(t)
	f.holders = []map[string]any{
		holderRow("party::1", holderInternal, "", ""),
		holderRow("ADDR2", holderWhitelisted, "", "2"),
	}
	result, err := f.query(t)
	if err == nil {
		t.Fatalf("every row was excluded, yet the call returned %+v", result)
	}
	var integrity *model.IntegrityError
	if !asIntegrityError(err, &integrity) || !strings.Contains(err.Error(), "all 2 asset address(es)") {
		t.Errorf("error = %v, want an IntegrityError naming the rows", err)
	}
}

// A verified reader that cannot answer aborts the call; it never turns into exclusions that let
// the rest of the page through.
func TestQueryAssetAddressesAbortsWhenAVerifiedReaderFails(t *testing.T) {
	for _, tc := range []struct {
		name  string
		setup func(f *holdersFixture)
	}{
		{"managed-address read fails", func(f *holdersFixture) {
			f.managedStatus = http.StatusInternalServerError
			f.holders = append(f.holders, holderRow("party::7", holderInternal, "7", ""))
		}},
		{"rules container has no HSMSLOT key", func(f *holdersFixture) {
			f.noHSMKey = true
			f.manage("7", "party::7", "party::7")
			f.holders = append(f.holders, holderRow("party::7", holderInternal, "7", ""))
		}},
		{"whitelist read fails for the whole batch", func(f *holdersFixture) {
			f.whitelistRow("5", true)
			f.holders = append(f.holders, holderRow("ADDR5", holderWhitelisted, "", "5"))
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newHoldersFixture(t)
			f.holders = []map[string]any{holderRow("party::ext", holderExternal, "", "")}
			tc.setup(f)
			if result, err := f.query(t); err == nil {
				t.Fatalf("the call succeeded with %+v", result)
			}
		})
	}
}
