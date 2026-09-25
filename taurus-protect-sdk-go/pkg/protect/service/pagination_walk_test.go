package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/cache"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// Walks through the real generated client: passing NextOffset / NextCursor back until HasMore is
// false visits every row exactly once, in server order, and the walk stops. The fakes page like
// validatord — zero values omitted, so the last page's reply leaves out whatever is zero.

// walkClient points a generated client at handler.
func walkClient(t *testing.T, handler http.HandlerFunc) *openapi.APIClient {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler(w, r)
	}))
	t.Cleanup(srv.Close)
	cfg := openapi.NewConfiguration()
	cfg.Servers = []openapi.ServerConfiguration{{URL: srv.URL}}
	cfg.HTTPClient = srv.Client()
	return openapi.NewAPIClient(cfg)
}

// replyJSON renders a reply the way validatord does: zero-valued fields are omitted.
func replyJSON(t *testing.T, fields map[string]any) string {
	t.Helper()
	out := map[string]any{}
	for k, v := range fields {
		switch x := v.(type) {
		case int:
			if x != 0 {
				out[k] = strconv.Itoa(x)
			}
		case string:
			if x != "" {
				out[k] = x
			}
		case bool:
			if x {
				out[k] = true
			}
		case []map[string]any:
			if len(x) > 0 {
				out[k] = x
			}
		case map[string]any:
			if len(x) > 0 {
				out[k] = x
			}
		default:
			t.Fatalf("replyJSON: unsupported %T", v)
		}
	}
	b, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func idRows(from, to int) []map[string]any {
	rows := []map[string]any{}
	for i := from; i < to; i++ {
		rows = append(rows, map[string]any{"id": strconv.Itoa(i)})
	}
	return rows
}

// queryInt reads a required integer query parameter.
func queryInt(t *testing.T, r *http.Request, name string) int {
	t.Helper()
	v := r.URL.Query().Get(name)
	if v == "" {
		return 0
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		t.Errorf("%s=%q is not an integer", name, v)
	}
	return n
}

// pageToken is an opaque continuation that holds + / and = once standard-base64 encoded, so a
// cursor that is not URL-encoded on the wire arrives altered and fails the walk.
func pageToken(t *testing.T, next int) string {
	t.Helper()
	token := base64.StdEncoding.EncodeToString([]byte{0xfb, 0xff, byte(next), 0xfe})
	if !strings.ContainsAny(token, "+") || !strings.ContainsAny(token, "/") || !strings.HasSuffix(token, "=") {
		t.Fatalf("token %q does not exercise + / =", token)
	}
	return token
}

func tokenIndex(t *testing.T, token string) int {
	t.Helper()
	raw, err := base64.StdEncoding.DecodeString(token)
	if err != nil || len(raw) != 4 || raw[0] != 0xfb || raw[1] != 0xff || raw[3] != 0xfe {
		t.Errorf("continuation %q is not a token this server issued (%v)", token, err)
		return -1
	}
	return int(raw[2])
}

func assertVisitedInOrder(t *testing.T, got []string, n int) {
	t.Helper()
	if len(got) != n {
		t.Fatalf("visited %d rows, want %d: %v", len(got), n, got)
	}
	for i, id := range got {
		if id != strconv.Itoa(i) {
			t.Fatalf("row %d is %q: rows were skipped, repeated or reordered: %v", i, id, got)
		}
	}
}

// maxWalkRequests guards a walk that would never end.
const maxWalkRequests = 200

func TestWalletsWalkFollowsTheReplyOffset(t *testing.T) {
	for _, n := range []int{0, 1, 19, 20, 21, 40, 45} {
		for _, size := range []int{1, 7, 20} {
			t.Run(fmt.Sprintf("rows=%d/size=%d", n, size), func(t *testing.T) {
				requests := 0
				client := walkClient(t, func(w http.ResponseWriter, r *http.Request) {
					requests++
					if got := queryInt(t, r, "limit"); got != size {
						t.Errorf("request %d sent limit %d, want %d", requests, got, size)
					}
					offset := queryInt(t, r, "offset")
					end := min(offset+size, n)
					start := min(offset, n)
					// validatord's reply offset is where the NEXT page starts.
					_, _ = io.WriteString(w, replyJSON(t, map[string]any{
						"result": idRows(start, end), "totalItems": n, "offset": end,
					}))
				})
				svc := NewWalletService(client)

				var ids []string
				opts := &model.ListWalletsOptions{Limit: int64(size)}
				for {
					wallets, page, err := svc.ListWallets(context.Background(), opts)
					if err != nil {
						t.Fatal(err)
					}
					for _, w := range wallets {
						ids = append(ids, w.ID)
					}
					if !page.HasMore {
						break
					}
					if requests > maxWalkRequests {
						t.Fatal("the walk does not end")
					}
					opts.Offset = page.NextOffset
				}
				assertVisitedInOrder(t, ids, n)
				wantRequests := max(1, (n+size-1)/size)
				if requests != wantRequests {
					t.Errorf("sent %d requests, want %d: a page was fetched after the last one", requests, wantRequests)
				}
			})
		}
	}
}

// validatord can append a synthetic daemon user beyond the limit; the next page must start at
// offset + limit, or the row after the page is skipped.
func TestUsersWalkIgnoresTheSyntheticRowBeyondTheLimit(t *testing.T) {
	const n, size = 45, 20
	client := walkClient(t, func(w http.ResponseWriter, r *http.Request) {
		offset := queryInt(t, r, "offset")
		rows := idRows(min(offset, n), min(offset+size, n))
		if offset == 0 {
			rows = append(rows, map[string]any{"id": "daemon"})
		}
		_, _ = io.WriteString(w, replyJSON(t, map[string]any{"result": rows, "totalItems": n}))
	})
	svc := NewUserService(client)

	var ids []string
	opts := &model.ListUsersOptions{Limit: size}
	for requests := 0; ; requests++ {
		result, err := svc.ListUsers(context.Background(), opts)
		if err != nil {
			t.Fatal(err)
		}
		for _, u := range result.Users {
			if u.ID != "daemon" {
				ids = append(ids, u.ID)
			}
		}
		if !result.Pagination.HasMore {
			break
		}
		if requests > maxWalkRequests {
			t.Fatal("the walk does not end")
		}
		opts.Offset = result.Pagination.NextOffset
	}
	assertVisitedInOrder(t, ids, n)
}

// On the post-filter path the transactions total is an upper bound: a trailing page comes back
// empty, and a page that made no progress ends the walk.
func TestTransactionsWalkEndsOnATrailingEmptyPage(t *testing.T) {
	const existing, reportedTotal, size = 40, 45, 20
	requests := 0
	client := walkClient(t, func(w http.ResponseWriter, r *http.Request) {
		requests++
		offset := queryInt(t, r, "offset")
		_, _ = io.WriteString(w, replyJSON(t, map[string]any{
			"result": idRows(min(offset, existing), min(offset+size, existing)), "totalItems": reportedTotal,
		}))
	})
	svc := NewTransactionService(client)

	var ids []string
	opts := &model.ListTransactionsOptions{Limit: size}
	for {
		txs, page, err := svc.ListTransactions(context.Background(), opts)
		if err != nil {
			t.Fatal(err)
		}
		for _, tx := range txs {
			ids = append(ids, tx.ID)
		}
		if !page.HasMore {
			break
		}
		if requests > maxWalkRequests {
			t.Fatal("the walk does not end")
		}
		opts.Offset = page.NextOffset
	}
	assertVisitedInOrder(t, ids, existing)
	if requests != 3 {
		t.Errorf("sent %d requests, want 3 (two full pages and the empty one)", requests)
	}
}

// A v2 cursor list: the continuation is sent byte-exact with pageRequest=NEXT and the page size
// on every request, and the last page (currentPage present, hasNext omitted) ends the walk.
func TestRequestsWalkSendsTheCursorBack(t *testing.T) {
	for _, n := range []int{0, 1, 5, 6, 12} {
		t.Run(fmt.Sprintf("rows=%d", n), func(t *testing.T) {
			const size = 6
			requests := 0
			client := walkClient(t, func(w http.ResponseWriter, r *http.Request) {
				requests++
				q := r.URL.Query()
				if q.Get("cursor.pageSize") != strconv.Itoa(size) {
					t.Errorf("request %d: cursor.pageSize = %q", requests, q.Get("cursor.pageSize"))
				}
				start := 0
				if token := q.Get("cursor.currentPage"); token != "" {
					if q.Get("cursor.pageRequest") != "NEXT" {
						t.Errorf("request %d: pageRequest = %q, want NEXT", requests, q.Get("cursor.pageRequest"))
					}
					// The raw query must carry the escaped bytes, not a literal + / =.
					if !strings.Contains(r.URL.RawQuery, "%2B") || !strings.Contains(r.URL.RawQuery, "%2F") ||
						!strings.Contains(r.URL.RawQuery, "%3D") {
						t.Errorf("request %d: cursor not URL-encoded: %s", requests, r.URL.RawQuery)
					}
					start = tokenIndex(t, token)
				} else if q.Get("cursor.pageRequest") != "" {
					t.Errorf("request %d: a first page sends no pageRequest", requests)
				}
				end := min(start+size, n)
				cursor := map[string]any{"currentPage": pageToken(t, start)}
				if end < n {
					cursor = map[string]any{"currentPage": pageToken(t, end), "hasNext": true, "hasPrevious": start > 0}
				}
				_, _ = io.WriteString(w, replyJSON(t, map[string]any{"result": idRows(start, end), "cursor": cursor}))
			})
			svc := NewRequestService(client)

			var ids []string
			opts := &model.ListRequestsOptions{PageSize: size}
			for {
				result, err := svc.ListRequests(context.Background(), opts)
				if err != nil {
					t.Fatal(err)
				}
				for _, r := range result.Requests {
					ids = append(ids, r.ID)
				}
				if !result.Page.HasMore {
					if result.Page.NextCursor != "" {
						t.Errorf("the last page carries a cursor %q", result.Page.NextCursor)
					}
					break
				}
				if requests > maxWalkRequests {
					t.Fatal("the walk does not end")
				}
				opts.Cursor = result.Page.NextCursor
			}
			assertVisitedInOrder(t, ids, n)
			if want := max(1, (n+size-1)/size); requests != want {
				t.Errorf("sent %d requests, want %d", requests, want)
			}
		})
	}
}

// The prices list takes its cursor in the request body.
func TestPricesWalkSendsTheCursorInTheBody(t *testing.T) {
	const n, size = 5, 2
	requests := 0
	client := walkClient(t, func(w http.ResponseWriter, r *http.Request) {
		requests++
		var body struct {
			Cursor struct {
				CurrentPage string `json:"currentPage"`
				PageRequest string `json:"pageRequest"`
				PageSize    string `json:"pageSize"`
			} `json:"cursor"`
			From *struct {
				CurrencyFromID string `json:"currencyFromId"`
			} `json:"from"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("request %d: body: %v", requests, err)
		}
		if body.Cursor.PageSize != strconv.Itoa(size) || body.From == nil || body.From.CurrencyFromID != "c1" {
			t.Errorf("request %d: page size or filter lost on the continuation: %+v", requests, body)
		}
		start := 0
		if body.Cursor.CurrentPage != "" {
			if body.Cursor.PageRequest != "NEXT" {
				t.Errorf("request %d: pageRequest = %q", requests, body.Cursor.PageRequest)
			}
			start = tokenIndex(t, body.Cursor.CurrentPage)
		}
		end := min(start+size, n)
		rows := []map[string]any{}
		for i := start; i < end; i++ {
			rows = append(rows, map[string]any{"id": strconv.Itoa(i), "rate": "1"})
		}
		cursor := map[string]any{}
		if end < n {
			cursor = map[string]any{"currentPage": pageToken(t, end), "hasNext": true}
		}
		_, _ = io.WriteString(w, replyJSON(t, map[string]any{"result": rows, "cursor": cursor, "baseCurrency": "CHF"}))
	})
	// No PRICEUPDATER in the container: this tenant does not sign prices.
	rules := cache.NewRulesContainerCache(time.Hour, func(context.Context) (*model.DecodedRulesContainer, error) {
		return &model.DecodedRulesContainer{}, nil
	})
	svc := NewPriceService(client, rules)

	var ids []string
	opts := &model.ListPricesOptions{PageSize: size, FromCurrencyID: "c1"}
	for {
		result, err := svc.ListPrices(context.Background(), opts)
		if err != nil {
			t.Fatal(err)
		}
		for _, p := range result.Prices {
			ids = append(ids, p.ID)
		}
		if !result.Page.HasMore {
			break
		}
		if requests > maxWalkRequests {
			t.Fatal("the walk does not end")
		}
		opts.Cursor = result.Page.NextCursor
	}
	assertVisitedInOrder(t, ids, n)
}

// Wallet tokens page by a bare token: the reply's next is sent back as cursor, and its absence
// ends the list. The total is the server's.
func TestWalletTokensWalkFollowsTheToken(t *testing.T) {
	const n, size = 7, 3
	requests := 0
	client := walkClient(t, func(w http.ResponseWriter, r *http.Request) {
		requests++
		if !strings.HasSuffix(r.URL.Path, "/wallets/9/tokens") {
			t.Errorf("path %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("limit") != strconv.Itoa(size) {
			t.Errorf("request %d: limit = %q", requests, q.Get("limit"))
		}
		start := 0
		if token := q.Get("cursor"); token != "" {
			start = tokenIndex(t, token)
		}
		end := min(start+size, n)
		rows := []map[string]any{}
		for i := start; i < end; i++ {
			rows = append(rows, map[string]any{"asset": map[string]any{"currency": strconv.Itoa(i)}})
		}
		next := ""
		if end < n {
			next = pageToken(t, end)
		}
		_, _ = io.WriteString(w, replyJSON(t, map[string]any{"balances": rows, "next": next, "total": n}))
	})
	svc := NewWalletService(client)

	var ids []string
	opts := &model.GetWalletTokensOptions{PageSize: size}
	for {
		result, err := svc.GetWalletTokens(context.Background(), "9", opts)
		if err != nil {
			t.Fatal(err)
		}
		for _, tok := range result.Tokens {
			ids = append(ids, tok.Asset.Currency)
		}
		if result.Page.TotalItems == nil || *result.Page.TotalItems != n {
			t.Errorf("TotalItems = %v, want %d", result.Page.TotalItems, n)
		}
		if !result.Page.HasMore {
			break
		}
		if requests > maxWalkRequests {
			t.Fatal("the walk does not end")
		}
		opts.Cursor = result.Page.NextCursor
	}
	assertVisitedInOrder(t, ids, n)
	if requests != 3 {
		t.Errorf("sent %d requests, want 3", requests)
	}
}

// Balances take requestCursor.* only — the legacy limit and bytes cursor are never sent.
func TestBalancesWalkUsesRequestCursorOnly(t *testing.T) {
	const n, size = 4, 3
	requests := 0
	client := walkClient(t, func(w http.ResponseWriter, r *http.Request) {
		requests++
		q := r.URL.Query()
		if q.Has("limit") || q.Has("cursor") {
			t.Errorf("request %d sent the legacy paging parameters: %s", requests, r.URL.RawQuery)
		}
		if q.Get("requestCursor.pageSize") != strconv.Itoa(size) {
			t.Errorf("request %d: requestCursor.pageSize = %q", requests, q.Get("requestCursor.pageSize"))
		}
		start := 0
		if token := q.Get("requestCursor.currentPage"); token != "" {
			if q.Get("requestCursor.pageRequest") != "NEXT" {
				t.Errorf("request %d: pageRequest = %q", requests, q.Get("requestCursor.pageRequest"))
			}
			start = tokenIndex(t, token)
		}
		end := min(start+size, n)
		rows := []map[string]any{}
		for i := start; i < end; i++ {
			rows = append(rows, map[string]any{"asset": map[string]any{"currency": strconv.Itoa(i)}})
		}
		cursor := map[string]any{}
		if end < n {
			cursor = map[string]any{"currentPage": pageToken(t, end), "hasNext": true}
		}
		_, _ = io.WriteString(w, replyJSON(t, map[string]any{"balances": rows, "cursor": cursor, "total": n}))
	})
	svc := NewBalanceService(client)

	var ids []string
	opts := &model.GetBalancesOptions{PageSize: size}
	for {
		result, err := svc.GetBalances(context.Background(), opts)
		if err != nil {
			t.Fatal(err)
		}
		for _, b := range result.Balances {
			ids = append(ids, b.Asset.Currency)
		}
		if result.Page.TotalItems == nil || *result.Page.TotalItems != n {
			t.Errorf("TotalItems = %v, want %d", result.Page.TotalItems, n)
		}
		if !result.Page.HasMore {
			break
		}
		if requests > maxWalkRequests {
			t.Fatal("the walk does not end")
		}
		opts.Cursor = result.Page.NextCursor
	}
	assertVisitedInOrder(t, ids, n)
}
