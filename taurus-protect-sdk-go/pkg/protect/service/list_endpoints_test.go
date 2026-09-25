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
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/cache"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/helper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// Endpoint and filter behaviour of the list methods, through the real generated client.

// captureServer records every request and answers each with reply(request).
type captureServer struct {
	mu       sync.Mutex
	requests []recordedRequest
}

func (c *captureServer) client(t *testing.T, reply func(r recordedRequest) string) *openapi.APIClient {
	t.Helper()
	return walkClient(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		rec := recordedRequest{method: r.Method, path: r.URL.Path, query: r.URL.RawQuery, body: body}
		c.mu.Lock()
		c.requests = append(c.requests, rec)
		c.mu.Unlock()
		_, _ = io.WriteString(w, reply(rec))
	})
}

func (c *captureServer) all() []recordedRequest {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]recordedRequest(nil), c.requests...)
}

func emptyReply(recordedRequest) string { return `{}` }

func decodedQuery(t *testing.T, r recordedRequest) url.Values {
	t.Helper()
	v, err := url.ParseQuery(r.query)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func decodedBody(t *testing.T, r recordedRequest) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal(r.body, &body); err != nil {
		t.Fatalf("body %q: %v", r.body, err)
	}
	return body
}

// The approval queue filters by ids, types and currency only; the filters it cannot apply are
// rejected by name before anything is sent, instead of silently returning an unfiltered queue.
func TestListRequestsForApprovalRejectsFiltersItCannotApply(t *testing.T) {
	now := time.Now()
	for name, opts := range map[string]*model.ListRequestsOptions{
		"Statuses":           {Statuses: []string{"CREATED"}},
		"FromDate":           {FromDate: &now},
		"ToDate":             {ToDate: &now},
		"ExternalRequestIDs": {ExternalRequestIDs: []string{"ext-1"}},
	} {
		t.Run(name, func(t *testing.T) {
			var srv captureServer
			svc := NewRequestService(srv.client(t, emptyReply))
			_, err := svc.ListRequestsForApproval(context.Background(), opts)
			if err == nil || !strings.Contains(err.Error(), name) {
				t.Errorf("error = %v, want one naming %s", err, name)
			}
			if n := len(srv.all()); n != 0 {
				t.Errorf("sent %d request(s) for a rejected call", n)
			}
		})
	}

	var srv captureServer
	svc := NewRequestService(srv.client(t, emptyReply))
	if _, err := svc.ListRequestsForApproval(context.Background(), &model.ListRequestsOptions{
		IDs: []string{"1", "2"}, Types: []string{"Payment"}, Currency: "ETH",
	}); err != nil {
		t.Fatal(err)
	}
	q := decodedQuery(t, srv.all()[0])
	if strings.Join(q["ids"], ",") != "1,2" || q.Get("types") != "Payment" || q.Get("currencyID") != "ETH" {
		t.Errorf("ids/types/currency did not reach the wire: %v", q)
	}
}

// Cursor continues from a previous page; combining it with a hand-set CurrentPage is ambiguous
// and must be rejected before any request.
func TestCursorAndCurrentPageAreMutuallyExclusive(t *testing.T) {
	var srv captureServer
	svc := NewChangeService(srv.client(t, emptyReply))
	_, err := svc.ListChanges(context.Background(), &model.ListChangesOptions{Cursor: "abc", CurrentPage: "def"})
	if err == nil || !strings.Contains(err.Error(), "CurrentPage") {
		t.Errorf("error = %v", err)
	}
	if n := len(srv.all()); n != 0 {
		t.Errorf("sent %d request(s)", n)
	}
}

// ListPrices is QueryPricesV2 with the same verification the v1 list had: with a PRICEUPDATER in
// the rules container, every row must carry a signature that covers its rate.
func TestListPricesVerifiesEveryRow(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	container := &model.DecodedRulesContainer{Users: []*model.RuleUser{
		{ID: "updater", PublicKey: &key.PublicKey, Roles: []string{"PRICEUPDATER"}},
	}}
	signed := func(rate string) map[string]any {
		data, err := helper.PriceSignedBytes(&model.Price{Blockchain: "ETH", CurrencyFrom: "ETH", CurrencyTo: "CHF", Decimals: "2", Rate: rate})
		if err != nil {
			t.Fatal(err)
		}
		sig, err := crypto.SignData(key, data)
		if err != nil {
			t.Fatal(err)
		}
		return map[string]any{"id": "p-" + rate, "blockchain": "ETH", "currencyFrom": "ETH", "currencyTo": "CHF",
			"decimals": "2", "rate": rate, "isPrimary": true, "status": "OK",
			"signatures": []any{map[string]any{"userId": "updater", "signature": sig}}}
	}
	rules := cache.NewRulesContainerCache(time.Hour, func(context.Context) (*model.DecodedRulesContainer, error) {
		return container, nil
	})

	for name, tc := range map[string]struct {
		row     map[string]any
		wantErr bool
	}{
		"signed row": {signed("300000"), false},
		"tampered rate": {func() map[string]any {
			r := signed("300000")
			r["rate"] = "1"
			return r
		}(), true},
	} {
		t.Run(name, func(t *testing.T) {
			body, _ := json.Marshal(map[string]any{"result": []any{tc.row}, "baseCurrency": "CHF"})
			var srv captureServer
			svc := NewPriceService(srv.client(t, func(recordedRequest) string { return string(body) }), rules)
			result, err := svc.ListPrices(context.Background(), &model.ListPricesOptions{OnlyPrimary: true, SortOrder: "ASC"})
			if tc.wantErr {
				if err == nil {
					t.Fatal("a price whose signature does not cover its rate was returned")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			p := result.Prices[0]
			if result.BaseCurrency != "CHF" || p.ID != "p-300000" || !p.IsPrimary || p.Status != "OK" {
				t.Errorf("result = %+v, price = %+v", result, p)
			}
			req := decodedBody(t, srv.all()[0])
			if req["onlyPrimary"] != true || req["sortOrder"] != "ASC" {
				t.Errorf("onlyPrimary/sortOrder did not reach the body: %v", req)
			}
			if srv.all()[0].path != "/api/rest/v2/prices/query" {
				t.Errorf("path = %s, want the QueryPricesV2 endpoint", srv.all()[0].path)
			}
		})
	}
}

// The transaction export cannot page: Limit is its only size control, the total is reported,
// and a malformed total is an error rather than a zero.
func TestExportTransactionsReportsTheTotal(t *testing.T) {
	var srv captureServer
	svc := NewTransactionService(srv.client(t, func(recordedRequest) string {
		return `{"result":"id,amount\n1,5\n","totalItems":"37"}`
	}))
	result, err := svc.ExportTransactions(context.Background(), &model.ExportTransactionsOptions{Format: "csv", Limit: 5000})
	if err != nil {
		t.Fatal(err)
	}
	if result.Data != "id,amount\n1,5\n" || result.TotalItems != 37 {
		t.Errorf("result = %+v", result)
	}
	q := decodedQuery(t, srv.all()[0])
	if q.Get("format") != "csv" || q.Get("limit") != "5000" || q.Has("offset") {
		t.Errorf("query = %v", q)
	}

	var bad captureServer
	svc = NewTransactionService(bad.client(t, func(recordedRequest) string { return `{"totalItems":"many"}` }))
	if _, err := svc.ExportTransactions(context.Background(), nil); err == nil {
		t.Error("a malformed total was accepted")
	}
}

// GetUsersByEmail sends at most MaxPageSize emails per request and walks each batch to its last
// page, so a batch matching more users than one page holds is still read in full.
func TestGetUsersByEmailBatchesTheEmails(t *testing.T) {
	var srv captureServer
	svc := NewUserService(srv.client(t, func(r recordedRequest) string {
		q, _ := url.ParseQuery(r.query)
		// Every email matches two users, so the first batch spans two pages.
		var matches []map[string]any
		for _, e := range q["emails"] {
			matches = append(matches, map[string]any{"id": e + "#1"}, map[string]any{"id": e + "#2"})
		}
		limit, _ := strconv.Atoi(q.Get("limit"))
		offset, _ := strconv.Atoi(q.Get("offset"))
		page := matches[min(offset, len(matches)):min(offset+limit, len(matches))]
		body, _ := json.Marshal(map[string]any{"result": page, "totalItems": fmt.Sprint(len(matches))})
		return string(body)
	}))

	emails := make([]string, 150)
	for i := range emails {
		emails[i] = fmt.Sprintf("u%03d@bank.com", i)
	}
	users, err := svc.GetUsersByEmail(context.Background(), emails)
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, u := range users {
		seen[u.ID] = true
	}
	if len(users) != 300 || len(seen) != 300 {
		t.Errorf("got %d users (%d distinct), want 300", len(users), len(seen))
	}
	requests := srv.all()
	if len(requests) != 3 {
		t.Fatalf("sent %d requests, want 3 (two pages for 100 emails, one for 50)", len(requests))
	}
	for i, want := range []struct{ emails, offset int }{{100, 0}, {100, 100}, {50, 0}} {
		q := decodedQuery(t, requests[i])
		if len(q["emails"]) != want.emails || q.Get("limit") != "100" || queryOffset(q) != want.offset {
			t.Errorf("request %d: %d emails, limit %s, offset %d; want %d emails at offset %d",
				i, len(q["emails"]), q.Get("limit"), queryOffset(q), want.emails, want.offset)
		}
	}
}

func queryOffset(q url.Values) int {
	n, _ := strconv.Atoi(q.Get("offset"))
	return n
}

// The approval re-read goes through the verifying list in batches of at most MaxPageSize ids,
// and one signature still covers every reviewed row.
func TestApproveWhitelistedAddressesRereadsInBatches(t *testing.T) {
	f := newSignedWhitelist(t)
	const n = 150
	rows := map[string]map[string]any{}
	hashes := map[string]string{}
	for i := 0; i < n; i++ {
		id := fmt.Sprint(i + 1)
		rows[id] = f.row(t, id, false)
		hashes[id] = rows[id]["metadata"].(map[string]any)["hash"].(string)
	}

	var srv captureServer
	var approved []string
	client := srv.client(t, func(r recordedRequest) string {
		if strings.HasSuffix(r.path, "/approve") {
			var body struct {
				IDs []string `json:"ids"`
			}
			_ = json.Unmarshal(r.body, &body)
			approved = body.IDs
			return `{}`
		}
		q, _ := url.ParseQuery(r.query)
		page := []any{}
		for _, id := range q["ids"] {
			page = append(page, rows[id])
		}
		body, _ := json.Marshal(map[string]any{"result": page, "totalItems": fmt.Sprint(len(page))})
		return string(body)
	})

	if err := f.service(client).ApproveWhitelistedAddresses(context.Background(),
		reviewedSelection(t, hashes), addressApprovalKey(t), "ok"); err != nil {
		t.Fatal(err)
	}

	var reads []recordedRequest
	for _, r := range srv.all() {
		if !strings.HasSuffix(r.path, "/approve") {
			reads = append(reads, r)
		}
	}
	if len(reads) != 2 {
		t.Fatalf("sent %d re-reads, want 2", len(reads))
	}
	for i, want := range []int{100, 50} {
		q := decodedQuery(t, reads[i])
		if len(q["ids"]) != want || q.Get("limit") != fmt.Sprint(want) {
			t.Errorf("re-read %d: %d ids, limit %s; want %d", i, len(q["ids"]), q.Get("limit"), want)
		}
		// The rows being approved are pending, so the re-read must include them.
		if q.Get("includeForApproval") != "true" {
			t.Errorf("re-read %d omits includeForApproval: %v", i, q)
		}
	}
	if len(approved) != n || !sort.SliceIsSorted(approved, func(i, j int) bool {
		return len(approved[i]) < len(approved[j]) || (len(approved[i]) == len(approved[j]) && approved[i] < approved[j])
	}) {
		t.Errorf("approved %d ids (sorted numerically: %v), want all %d", len(approved), approved[:min(5, len(approved))], n)
	}
}

// The asset approval re-read batches the same way; a batch that omits rows aborts before signing.
func TestApproveWhitelistedAssetsRereadsInBatches(t *testing.T) {
	var srv captureServer
	client := srv.client(t, emptyReply)
	svc := NewWhitelistedAssetServiceWithVerification(client, &WhitelistedAssetServiceConfig{
		SuperAdminKeys: []*ecdsa.PublicKey{&addressApprovalKey(t).PublicKey}, MinValidSignatures: 1})

	hashes := map[string]string{}
	for i := 1; i <= 150; i++ {
		hashes[fmt.Sprint(i)] = fmt.Sprintf("%064x", i)
	}
	err := svc.ApproveWhitelistedAssets(context.Background(), reviewedAssetSelection(t, hashes), addressApprovalKey(t), "ok")
	if err == nil {
		t.Fatal("rows the re-read did not return were approved")
	}
	requests := srv.all()
	if len(requests) != 2 {
		t.Fatalf("sent %d requests, want 2 re-reads and no approval", len(requests))
	}
	for i, want := range []int{100, 50} {
		q := decodedQuery(t, requests[i])
		if len(q["whitelistedContractAddressIds"]) != want || q.Get("limit") != fmt.Sprint(want) {
			t.Errorf("re-read %d: %d ids, limit %s; want %d", i, len(q["whitelistedContractAddressIds"]), q.Get("limit"), want)
		}
	}
}

// The v2 asset reads, Earn rewards and fiat entities: every filter reaches the wire, and a row
// maps into the domain model.
func TestNewCursorListsWireFiltersAndMapRows(t *testing.T) {
	ctx := context.Background()

	t.Run("QueryAssets", func(t *testing.T) {
		var srv captureServer
		svc := testAssetService(t, srv.client(t, func(recordedRequest) string {
			return `{"result":[{"id":"a1","label":"Gold","status":"ASSET_VIEW_STATUS_V2_ACTIVE","symbol":"GLD",` +
				`"decimals":"6","attributes":[{"key":"k","value":"v"}],"blockchainAsset":{"cantonNativeTokenAsset":` +
				`{"instrumentID":"i1","configuration":{"cid":"c1","paused":true}}}}]}`
		}), emptyRulesCache())
		result, err := svc.QueryAssets(ctx, &model.QueryAssetsOptions{Blockchain: "CANTON", Network: "mainnet",
			Symbol: "GLD", ContractAddress: "0x1", Label: "Gold", CurrencyName: "gold"})
		if err != nil {
			t.Fatal(err)
		}
		a := result.Assets[0]
		if a.ID != "a1" || a.Status != "ASSET_VIEW_STATUS_V2_ACTIVE" || a.Attributes[0].Key != "k" ||
			a.CantonInstrument == nil || a.CantonInstrument.ContractID != "c1" || !a.CantonInstrument.Paused {
			t.Errorf("asset = %+v", a)
		}
		body := decodedBody(t, srv.all()[0])
		for field, want := range map[string]string{"blockchain": "CANTON", "network": "mainnet", "symbol": "GLD",
			"contractAddress": "0x1", "label": "Gold", "currencyName": "gold"} {
			if body[field] != want {
				t.Errorf("%s = %v, want %s", field, body[field], want)
			}
		}
	})

	t.Run("QueryAssetAddresses", func(t *testing.T) {
		var srv captureServer
		svc := testAssetService(t, srv.client(t, func(recordedRequest) string {
			return `{"result":[{"address":"0xabc","kycStatus":"KYC_STATUS_V2_APPROVED","balance":"12",` +
				`"addressType":"ADDRESS_TYPE_V2_EXTERNAL"}]}`
		}), emptyRulesCache())
		result, err := svc.QueryAssetAddresses(ctx, "a1", &model.QueryAssetAddressesOptions{
			AddressType: "ADDRESS_TYPE_V2_EXTERNAL", KYCStatus: "KYC_STATUS_V2_APPROVED"})
		if err != nil {
			t.Fatal(err)
		}
		if a := result.Addresses[0]; a.Address != "0xabc" || a.KYCStatus != "KYC_STATUS_V2_APPROVED" || a.Balance != "12" {
			t.Errorf("address = %+v", a)
		}
		body := decodedBody(t, srv.all()[0])
		if body["addressType"] != "ADDRESS_TYPE_V2_EXTERNAL" || body["kycStatus"] != "KYC_STATUS_V2_APPROVED" {
			t.Errorf("body = %v", body)
		}
		if srv.all()[0].path != "/api/rest/v2/assets/a1/addresses/query" {
			t.Errorf("path = %s", srv.all()[0].path)
		}
		if _, err := svc.QueryAssetAddresses(ctx, "", nil); err == nil {
			t.Error("an empty assetID was accepted")
		}
	})

	t.Run("ListAssetOperations", func(t *testing.T) {
		var srv captureServer
		svc := testAssetService(t, srv.client(t, func(recordedRequest) string {
			return `{"result":[{"id":"o1","type":"ASSET_OPERATION_TYPE_V2_MINT","status":"ASSET_OPERATION_STATUS_V2_PENDING",` +
				`"mint":{"amount":"5","destination":{"addressID":"9"}}}]}`
		}), emptyRulesCache())
		result, err := svc.ListAssetOperations(ctx, "a1", &model.ListAssetOperationsOptions{
			Type: "ASSET_OPERATION_TYPE_V2_MINT", Status: "ASSET_OPERATION_STATUS_V2_PENDING"})
		if err != nil {
			t.Fatal(err)
		}
		if op := result.Operations[0]; op.Amount != "5" || op.Target == nil || op.Target.AddressID != "9" {
			t.Errorf("operation = %+v", op)
		}
		q := decodedQuery(t, srv.all()[0])
		if q.Get("type") != "ASSET_OPERATION_TYPE_V2_MINT" || q.Get("status") != "ASSET_OPERATION_STATUS_V2_PENDING" {
			t.Errorf("query = %v", q)
		}
	})

	t.Run("ListRewards", func(t *testing.T) {
		var srv captureServer
		svc := NewEarnService(srv.client(t, func(recordedRequest) string {
			return `{"rewards":[{"id":"r1","recipientAddressId":"4","rewardType":"RewardTypeMerklToken",` +
				`"merklTokenReward":{"amount":"10","pending":"2","token":{"symbol":"MRK"}}}]}`
		}))
		result, err := svc.ListRewards(ctx, &model.ListEarnRewardsOptions{RecipientAddressID: "4"})
		if err != nil {
			t.Fatal(err)
		}
		r := result.Rewards[0]
		if r.ID != "r1" || r.RewardType != "RewardTypeMerklToken" || r.MerklToken == nil || r.MerklToken.Token.Symbol != "MRK" {
			t.Errorf("reward = %+v", r)
		}
		if q := decodedQuery(t, srv.all()[0]); q.Get("recipientAddressId") != "4" {
			t.Errorf("query = %v", q)
		}
	})

	t.Run("ListFiatProviderEntities", func(t *testing.T) {
		var srv captureServer
		svc := NewFiatService(srv.client(t, func(recordedRequest) string {
			return `{"result":[{"id":"e1","provider":"circle","name":"ACME","creationDate":"2024-01-02T00:00:00Z"}]}`
		}))
		result, err := svc.ListFiatProviderEntities(ctx, &model.ListFiatProviderEntitiesOptions{
			Provider: "circle", Label: "main", SortOrder: "ASC"})
		if err != nil {
			t.Fatal(err)
		}
		if e := result.Entities[0]; e.ID != "e1" || e.Name != "ACME" || e.CreationDate.Year() != 2024 {
			t.Errorf("entity = %+v", e)
		}
		q := decodedQuery(t, srv.all()[0])
		if q.Get("provider") != "circle" || q.Get("label") != "main" || q.Get("sortOrder") != "ASC" {
			t.Errorf("query = %v", q)
		}
	})
}

// testAssetService builds an AssetService whose verified readers share its client.
func testAssetService(t *testing.T, client *openapi.APIClient, rules *cache.RulesContainerCache) *AssetService {
	t.Helper()
	whitelisted := NewWhitelistedAddressServiceWithVerification(client, &WhitelistedAddressServiceConfig{
		SuperAdminKeys: []*ecdsa.PublicKey{&addressApprovalKey(t).PublicKey}, MinValidSignatures: 1})
	return NewAssetService(client, rules, NewAddressService(client, rules), whitelisted)
}

func emptyRulesCache() *cache.RulesContainerCache {
	return cache.NewRulesContainerCache(time.Hour, func(context.Context) (*model.DecodedRulesContainer, error) {
		return &model.DecodedRulesContainer{}, nil
	})
}
