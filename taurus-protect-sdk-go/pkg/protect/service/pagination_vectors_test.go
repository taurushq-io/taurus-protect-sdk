package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// Cross-SDK vectors for the pagination helpers (repo-root CLAUDE.md, "Pagination (cross-SDK)"):
// reply → page value for every offset rule and the cursor/token families, page-size and offset
// resolution, and the rule of every paged operation — checked by running the Go method itself.
const paginationVectorsRelPath = "../../../../scripts/resources/pagination-vectors.json"

// paginationVectorCounts are the file's declared section sizes, asserted so a vector added
// without being consumed here fails.
var paginationVectorCounts = map[string]int{
	"operations":   53,
	"offset":       26,
	"cursor":       12,
	"page_size":    14,
	"offset_input": 4,
}

type paginationVectorFile struct {
	Constants struct {
		DefaultPageSize int64 `json:"default_page_size"`
		MaxPageSize     int64 `json:"max_page_size"`
		PriceHistoryMax int64 `json:"price_history_max"`
		MaxCount        int64 `json:"max_count"`
	} `json:"constants"`
	Counts     map[string]int `json:"counts"`
	Operations map[string]struct {
		Rule  string `json:"rule"`
		Size  string `json:"size"`
		Total bool   `json:"total"`
	} `json:"operations"`
	Offset      []offsetVector      `json:"offset"`
	Cursor      []cursorVector      `json:"cursor"`
	PageSize    []pageSizeVector    `json:"page_size"`
	OffsetInput []offsetInputVector `json:"offset_input"`
}

type offsetVector struct {
	Description string `json:"description"`
	Rule        string `json:"rule"`
	Request     struct {
		Limit  int64 `json:"limit"`
		Offset int64 `json:"offset"`
	} `json:"request"`
	ServedRows  int `json:"served_rows"`
	SDKExcluded int `json:"sdk_excluded"`
	Reply       struct {
		TotalItems *string `json:"totalItems"`
		Offset     *string `json:"offset"`
	} `json:"reply"`
	Expect *struct {
		Limit      int64 `json:"limit"`
		Offset     int64 `json:"offset"`
		TotalItems int64 `json:"total_items"`
		NextOffset int64 `json:"next_offset"`
		HasMore    bool  `json:"has_more"`
	} `json:"expect"`
	ExpectError bool `json:"expect_error"`
}

type cursorVector struct {
	Description string                              `json:"description"`
	Family      string                              `json:"family"`
	PageSize    int64                               `json:"page_size"`
	HasTotal    bool                                `json:"has_total"`
	ReplyCursor *openapi.TgvalidatordResponseCursor `json:"reply_cursor"`
	ReplyToken  *string                             `json:"reply_token"`
	ReplyTotal  *string                             `json:"reply_total"`
	Expect      *struct {
		PageSize   int64  `json:"page_size"`
		NextCursor string `json:"next_cursor"`
		HasMore    bool   `json:"has_more"`
		TotalItems *int64 `json:"total_items"`
	} `json:"expect"`
	ExpectError bool `json:"expect_error"`
}

type pageSizeVector struct {
	Description string `json:"description"`
	Kind        string `json:"kind"`
	Input       *int64 `json:"input"`
	Expect      *int64 `json:"expect"`
	ExpectError bool   `json:"expect_error"`
}

type offsetInputVector struct {
	Description string `json:"description"`
	Input       *int64 `json:"input"`
	Expect      *int64 `json:"expect"`
	ExpectError bool   `json:"expect_error"`
}

var offsetRulesByName = map[string]offsetRule{
	"reply_offset":        ruleReplyOffset,
	"plus_rows":           rulePlusRows,
	"plus_min_rows_limit": rulePlusMinRowsLimit,
	"plus_server_rows":    rulePlusServerRows,
	"plus_limit":          rulePlusLimit,
}

func loadPaginationVectors(t *testing.T) paginationVectorFile {
	t.Helper()
	path := filepath.Clean(paginationVectorsRelPath)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read shared pagination vectors %s: %v", path, err)
	}
	var file paginationVectorFile
	if err := json.Unmarshal(raw, &file); err != nil {
		t.Fatalf("cannot parse %s: %v", path, err)
	}

	sizes := map[string]int{
		"operations":   len(file.Operations),
		"offset":       len(file.Offset),
		"cursor":       len(file.Cursor),
		"page_size":    len(file.PageSize),
		"offset_input": len(file.OffsetInput),
	}
	for section, declared := range file.Counts {
		want, known := paginationVectorCounts[section]
		if !known {
			t.Fatalf("the file declares a section %q this loader does not consume", section)
		}
		if declared != want || sizes[section] != declared {
			t.Fatalf("%s: file declares %d, holds %d, this loader expects %d", section, declared, sizes[section], want)
		}
	}
	if len(file.Counts) != len(paginationVectorCounts) {
		t.Fatalf("the file declares %d sections, this loader consumes %d", len(file.Counts), len(paginationVectorCounts))
	}
	return file
}

func TestPaginationVectorConstants(t *testing.T) {
	c := loadPaginationVectors(t).Constants
	if c.DefaultPageSize != model.DefaultPageSize || c.MaxPageSize != model.MaxPageSize ||
		c.PriceHistoryMax != model.MaxPriceHistoryLimit || c.MaxCount != maxCount {
		t.Errorf("constants %+v disagree with model.DefaultPageSize=%d, MaxPageSize=%d, MaxPriceHistoryLimit=%d, maxCount=%d",
			c, model.DefaultPageSize, model.MaxPageSize, model.MaxPriceHistoryLimit, int64(maxCount))
	}
}

func TestPaginationOffsetVectors(t *testing.T) {
	for i, v := range loadPaginationVectors(t).Offset {
		v := v
		t.Run(fmt.Sprintf("%02d_%s", i, v.Description), func(t *testing.T) {
			rule, ok := offsetRulesByName[v.Rule]
			if !ok {
				t.Fatalf("unknown offset rule %q", v.Rule)
			}
			window, err := resolveOffsetWindow(v.Request.Limit, v.Request.Offset)
			if err != nil {
				t.Fatalf("request window: %v", err)
			}
			got, err := offsetPagination(rule, window, v.ServedRows, v.SDKExcluded,
				offsetReply{TotalItems: v.Reply.TotalItems, Offset: v.Reply.Offset})
			if v.ExpectError {
				if err == nil {
					t.Fatalf("expected an error, got %+v", got)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			want := model.Pagination{Limit: v.Expect.Limit, Offset: v.Expect.Offset, TotalItems: v.Expect.TotalItems,
				NextOffset: v.Expect.NextOffset, HasMore: v.Expect.HasMore}
			if *got != want {
				t.Errorf("got %+v, want %+v", *got, want)
			}
		})
	}
}

func TestPaginationCursorVectors(t *testing.T) {
	for i, v := range loadPaginationVectors(t).Cursor {
		v := v
		t.Run(fmt.Sprintf("%02d_%s", i, v.Description), func(t *testing.T) {
			if v.Family != "cursor" && v.Family != "token" {
				t.Fatalf("unknown cursor family %q", v.Family)
			}
			got, err := cursorPage(v.PageSize, cursorReply{
				Cursor:    v.ReplyCursor,
				TokenOnly: v.Family == "token",
				Token:     v.ReplyToken,
				HasTotal:  v.HasTotal,
				Total:     v.ReplyTotal,
			})
			if v.ExpectError {
				if err == nil {
					t.Fatalf("expected an error, got %+v", got)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got.PageSize != v.Expect.PageSize || got.NextCursor != v.Expect.NextCursor || got.HasMore != v.Expect.HasMore {
				t.Errorf("got %+v, want %+v", got, *v.Expect)
			}
			switch {
			case v.Expect.TotalItems == nil && got.TotalItems != nil:
				t.Errorf("TotalItems = %d, want none", *got.TotalItems)
			case v.Expect.TotalItems != nil && (got.TotalItems == nil || *got.TotalItems != *v.Expect.TotalItems):
				t.Errorf("TotalItems = %v, want %d", got.TotalItems, *v.Expect.TotalItems)
			}
		})
	}
}

func TestPaginationPageSizeVectors(t *testing.T) {
	maxByKind := map[string]int64{"page": model.MaxPageSize, "price_history": model.MaxPriceHistoryLimit, "export": 0}
	for i, v := range loadPaginationVectors(t).PageSize {
		v := v
		t.Run(fmt.Sprintf("%02d_%s", i, v.Description), func(t *testing.T) {
			max, ok := maxByKind[v.Kind]
			if !ok {
				t.Fatalf("unknown page-size kind %q", v.Kind)
			}
			var input int64
			if v.Input != nil {
				input = *v.Input
			}
			got, err := resolveSize("PageSize", input, max)
			if v.ExpectError {
				if err == nil {
					t.Fatalf("expected an error, got %d", got)
				}
				return
			}
			if err != nil || got != *v.Expect {
				t.Errorf("got %d, %v; want %d", got, err, *v.Expect)
			}
		})
	}
}

func TestPaginationOffsetInputVectors(t *testing.T) {
	for i, v := range loadPaginationVectors(t).OffsetInput {
		v := v
		t.Run(fmt.Sprintf("%02d_%s", i, v.Description), func(t *testing.T) {
			var input int64
			if v.Input != nil {
				input = *v.Input
			}
			window, err := resolveOffsetWindow(0, input)
			if v.ExpectError {
				if err == nil {
					t.Fatalf("expected an error, got %+v", window)
				}
				return
			}
			if err != nil || window.offset != *v.Expect {
				t.Errorf("got %+v, %v; want offset %d", window, err, *v.Expect)
			}
		})
	}
}

// ruleProbe feeds each wrapped list method replies that tell the rules apart and returns the
// rule its page value follows. Offset probes are sent at limit 1, offset 40:
//
//	reply A  {"totalItems":"1000","offset":"77"}   ─▶ 77 reply_offset · 41 plus_limit · 40 otherwise
//	reply B  2 rows + {"totalItems":"1000"}        ─▶ 42 plus_rows · 41 plus_min_rows_limit
//
// Reply B is skipped for the verified whitelists, whose rows cannot be built unsigned; the
// whitelisted-address walk test pins plus_server_rows against plus_rows with exclusions.
func ruleProbe(t *testing.T, s *vectorServices, op string, rowsKey string) (rule string, total bool) {
	t.Helper()
	adapter := listAdapters[op]
	call := func(reply string, o canonicalOptions) observation {
		t.Helper()
		s.server.reset(reply)
		obs, err := adapter.call(context.Background(), s, map[string]string{"assetID": "a1", "id": "7", "base": "BTC", "quote": "USD"}, o)
		if err != nil {
			t.Fatalf("%s: %v (reply %s)", op, err, reply)
		}
		return obs
	}

	switch {
	case containsString(adapter.options, "offset"):
		a := call(`{"totalItems":"1000","offset":"77"}`, canonicalOptions{Limit: 1, Offset: 40})
		if a.offset == nil {
			t.Fatalf("%s: no offset pagination", op)
		}
		total := a.offset.TotalItems == 1000
		switch a.offset.NextOffset {
		case 77:
			return "reply_offset", total
		case 41:
			return "plus_limit", total
		case 40:
		default:
			t.Fatalf("%s: next offset %d fits no rule", op, a.offset.NextOffset)
		}
		if rowsKey == "" {
			return "plus_server_rows", total
		}
		b := call(fmt.Sprintf(`{%q:[{},{}],"totalItems":"1000"}`, rowsKey), canonicalOptions{Limit: 1, Offset: 40})
		switch b.offset.NextOffset {
		case 42:
			return "plus_rows", total
		case 41:
			return "plus_min_rows_limit", total
		}
		t.Fatalf("%s: next offset %d with 2 rows at limit 1 fits no rule", op, b.offset.NextOffset)

	case containsString(adapter.options, "page_size"):
		// The cursor families: a v1/v2 cursor object, or the bare token of wallet tokens and
		// rules history. Every total field name is offered; only the operation's own is read.
		cursorReply := `{"cursor":{"currentPage":"Q1VSU09S","hasNext":true},"total":"5","totalItems":"5"}`
		tokenReply := `{"next":"Q1VSU09S","cursor":"Q1VSU09S","total":"5","totalItems":"5"}`
		reply, family := cursorReply, "cursor"
		if op == "RuleService_GetRulesHistory" || op == "WalletService_GetWalletTokens" {
			reply, family = tokenReply, "token"
		}
		// Currency, Provider and Label are required by some of these methods.
		obs := call(reply, canonicalOptions{Currency: "c", Provider: "p", Label: "l"})
		if obs.cursor == nil || !obs.cursor.HasMore || obs.cursor.NextCursor != "Q1VSU09S" ||
			obs.cursor.PageSize != model.DefaultPageSize {
			t.Fatalf("%s: page %+v does not continue from the reply", op, obs.cursor)
		}
		return family, obs.cursor.TotalItems != nil && *obs.cursor.TotalItems == 5

	default:
		obs := call(`{"totalItems":"5"}`, canonicalOptions{})
		return "limit_only", obs.total != nil && *obs.total == 5
	}
	return "", false
}

// TestPaginationVectorOperationRules checks, per paged operation, that the Go method follows the
// rule the vectors assign it and reports a total exactly where the server returns one.
func TestPaginationVectorOperationRules(t *testing.T) {
	file := loadPaginationVectors(t)
	svc := newVectorServices(t)

	// The rows key of each offset reply whose rows need no signature.
	rowsKeys := map[string]string{
		"ActionService_GetActions":           "result",
		"FeePayerService_GetFeePayers":       "result",
		"TransactionService_GetTransactions": "result",
		"UserService_GetUsers":               "result",
		"UserService_GetGroups":              "result",
		"WalletService_GetWalletsV2":         "result",
		"WalletService_GetAddresses":         "",
	}

	ops := make([]string, 0, len(file.Operations))
	for op := range file.Operations {
		ops = append(ops, op)
	}
	sort.Strings(ops)
	for _, op := range ops {
		want := file.Operations[op]
		t.Run(op, func(t *testing.T) {
			if reason, skip := notWrappedOperations[op]; skip {
				t.Skip(reason)
			}
			if _, ok := listAdapters[op]; !ok {
				t.Fatalf("%s has no Go adapter", op)
			}
			rule, total := ruleProbe(t, svc, op, rowsKeys[op])
			if rule != want.Rule {
				t.Errorf("rule = %s, want %s", rule, want.Rule)
			}
			if total != want.Total {
				t.Errorf("reports a total = %v, want %v", total, want.Total)
			}
		})
	}
}

// validatord omits zero values, so `{}` is a valid empty page. Every wrapped list must answer it
// with a complete page value — never nil — so a consumer never renders an empty object.
func TestEveryListAnswersAnEmptyReplyWithAFullPage(t *testing.T) {
	file := loadPaginationVectors(t)
	svc := newVectorServices(t)
	path := map[string]string{"assetID": "a1", "id": "7", "base": "BTC", "quote": "USD"}
	required := canonicalOptions{Currency: "c", Provider: "p", Label: "l"}

	for op, adapter := range listAdapters {
		want := file.Operations[op]
		t.Run(op, func(t *testing.T) {
			svc.server.reset(`{}`)
			obs, err := adapter.call(context.Background(), svc, path, required)
			if err != nil {
				t.Fatalf("an empty page is not an error: %v", err)
			}
			switch want.Rule {
			case "limit_only":
				if want.Total && (obs.total == nil || *obs.total != 0) {
					t.Errorf("total = %v, want 0", obs.total)
				}
			case "cursor", "token":
				if obs.cursor == nil {
					t.Fatal("no page value")
				}
				page := *obs.cursor
				if page.PageSize != model.DefaultPageSize || page.HasMore || page.NextCursor != "" {
					t.Errorf("page = %+v, want the default size and no next page", page)
				}
				if want.Total != (page.TotalItems != nil) || (page.TotalItems != nil && *page.TotalItems != 0) {
					t.Errorf("TotalItems = %v, want a 0 total exactly when the reply carries one (%v)", page.TotalItems, want.Total)
				}
			default:
				if obs.offset == nil {
					t.Fatal("pagination is nil")
				}
				next := int64(0)
				if want.Rule == "plus_limit" {
					next = model.DefaultPageSize
				}
				wantPage := model.Pagination{Limit: model.DefaultPageSize, NextOffset: next}
				if *obs.offset != wantPage {
					t.Errorf("pagination = %+v, want %+v", *obs.offset, wantPage)
				}
			}
		})
	}
}
