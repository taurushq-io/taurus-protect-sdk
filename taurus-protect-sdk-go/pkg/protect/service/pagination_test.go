package service

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func strPtr(s string) *string { return &s }
func boolPtr(b bool) *bool    { return &b }

func TestResolvePageSize(t *testing.T) {
	for _, tc := range []struct {
		in, want int64
	}{{0, model.DefaultPageSize}, {1, 1}, {20, 20}, {model.MaxPageSize, model.MaxPageSize}} {
		got, err := resolvePageSize("PageSize", tc.in)
		if err != nil || got != tc.want {
			t.Errorf("resolvePageSize(%d) = %d, %v; want %d", tc.in, got, err, tc.want)
		}
	}
	for _, in := range []int64{model.MaxPageSize + 1, -1, math.MaxInt64, math.MinInt64} {
		_, err := resolvePageSize("PageSize", in)
		if err == nil {
			t.Errorf("resolvePageSize(%d) accepted a size outside [0, %d]", in, model.MaxPageSize)
			continue
		}
		// The error must name the option and the maximum, so a caller can fix the call.
		if !strings.Contains(err.Error(), "PageSize") || !strings.Contains(err.Error(), "100") {
			t.Errorf("resolvePageSize(%d) error %q does not name the option and the maximum", in, err)
		}
	}
}

func TestResolveSizeWithoutMaximum(t *testing.T) {
	// The prices-history export has no SDK maximum; 0 still selects the default.
	for _, tc := range []struct{ in, want int64 }{{0, model.DefaultPageSize}, {5000, 5000}, {math.MaxInt64, math.MaxInt64}} {
		got, err := resolveSize("Limit", tc.in, 0)
		if err != nil || got != tc.want {
			t.Errorf("resolveSize(%d, no max) = %d, %v; want %d", tc.in, got, err, tc.want)
		}
	}
	if _, err := resolveSize("Limit", -1, 0); err == nil || !strings.Contains(err.Error(), "Limit") {
		t.Errorf("a negative limit must be rejected by name, got %v", err)
	}
	if _, err := resolveSize("Limit", model.MaxPriceHistoryLimit+1, model.MaxPriceHistoryLimit); err == nil {
		t.Error("price history above 365 must be rejected")
	}
}

func TestResolveOffsetWindow(t *testing.T) {
	w, err := resolveOffsetWindow(0, 0)
	if err != nil || w != (offsetWindow{limit: model.DefaultPageSize}) {
		t.Errorf("defaults: got %+v, %v", w, err)
	}
	if _, err := resolveOffsetWindow(0, -1); err == nil || !strings.Contains(err.Error(), "Offset") {
		t.Errorf("a negative offset must be rejected by name, got %v", err)
	}
	if _, err := resolveOffsetWindow(101, 0); err == nil || !strings.Contains(err.Error(), "Limit") {
		t.Errorf("limit 101 must be rejected by name, got %v", err)
	}
}

func TestResolveCursorWindow(t *testing.T) {
	cursor := "eyJpZCI6NDF9+/8="

	w, err := resolveCursorWindow(0, cursor, "", "")
	if err != nil || w != (cursorWindow{pageSize: model.DefaultPageSize, currentPage: cursor, pageRequest: "NEXT"}) {
		t.Errorf("a continuation must send the cursor with NEXT and the default size, got %+v, %v", w, err)
	}

	// The low-level options pass through untouched.
	w, err = resolveCursorWindow(5, "", "abc", "PREVIOUS")
	if err != nil || w != (cursorWindow{pageSize: 5, currentPage: "abc", pageRequest: "PREVIOUS"}) {
		t.Errorf("low-level paging: got %+v, %v", w, err)
	}

	if _, err := resolveCursorWindow(0, cursor, "abc", ""); err == nil || !strings.Contains(err.Error(), "CurrentPage") {
		t.Errorf("Cursor with CurrentPage must be rejected by name, got %v", err)
	}
	if _, err := resolveCursorWindow(0, cursor, "", "FIRST"); err == nil || !strings.Contains(err.Error(), "PageRequest") {
		t.Errorf("Cursor with PageRequest must be rejected by name, got %v", err)
	}
	if _, err := resolveCursorWindow(-1, "", "", ""); err == nil {
		t.Error("a negative page size must be rejected")
	}
}

func TestCursorWindowBody(t *testing.T) {
	body := cursorWindow{pageSize: 7, currentPage: "ab+/cQ==", pageRequest: "NEXT"}.body()
	if body.PageSize == nil || *body.PageSize != "7" || body.CurrentPage == nil || *body.CurrentPage != "ab+/cQ==" ||
		body.PageRequest == nil || *body.PageRequest != "NEXT" {
		t.Errorf("body = %+v", body)
	}
	first := cursorWindow{pageSize: 20}.body()
	if first.CurrentPage != nil || first.PageRequest != nil || first.PageSize == nil || *first.PageSize != "20" {
		t.Errorf("a first page sends the page size alone, got %+v", first)
	}
}

func TestOffsetPaginationRules(t *testing.T) {
	w := offsetWindow{limit: 20, offset: 40}
	for _, tc := range []struct {
		name     string
		rule     offsetRule
		served   int
		excluded int
		reply    offsetReply
		want     model.Pagination
	}{
		{"reply offset is the next offset", ruleReplyOffset, 20, 0, offsetReply{strPtr("100"), strPtr("60")},
			model.Pagination{Limit: 20, Offset: 40, TotalItems: 100, NextOffset: 60, HasMore: true}},
		{"reply offset absent falls back to offset + rows", ruleReplyOffset, 7, 0, offsetReply{strPtr("100"), nil},
			model.Pagination{Limit: 20, Offset: 40, TotalItems: 100, NextOffset: 47, HasMore: true}},
		{"reply offset at the total ends the walk", ruleReplyOffset, 20, 0, offsetReply{strPtr("60"), strPtr("60")},
			model.Pagination{Limit: 20, Offset: 40, TotalItems: 60, NextOffset: 60, HasMore: false}},
		{"plus rows", rulePlusRows, 20, 0, offsetReply{TotalItems: strPtr("100")},
			model.Pagination{Limit: 20, Offset: 40, TotalItems: 100, NextOffset: 60, HasMore: true}},
		{"plus rows: no progress ends the walk", rulePlusRows, 0, 0, offsetReply{TotalItems: strPtr("100")},
			model.Pagination{Limit: 20, Offset: 40, TotalItems: 100, NextOffset: 40, HasMore: false}},
		{"plus min rows limit: a synthetic row beyond the limit", rulePlusMinRowsLimit, 21, 0, offsetReply{TotalItems: strPtr("100")},
			model.Pagination{Limit: 20, Offset: 40, TotalItems: 100, NextOffset: 60, HasMore: true}},
		{"plus server rows: exclusions reduce the total only", rulePlusServerRows, 20, 3, offsetReply{TotalItems: strPtr("100")},
			model.Pagination{Limit: 20, Offset: 40, TotalItems: 97, NextOffset: 60, HasMore: true}},
		{"plus limit: a short page is not the end", rulePlusLimit, 3, 0, offsetReply{TotalItems: strPtr("100")},
			model.Pagination{Limit: 20, Offset: 40, TotalItems: 100, NextOffset: 60, HasMore: true}},
		{"exclusions clamp at zero", rulePlusServerRows, 5, 9, offsetReply{TotalItems: strPtr("5")},
			model.Pagination{Limit: 20, Offset: 40, TotalItems: 0, NextOffset: 45, HasMore: false}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := offsetPagination(tc.rule, w, tc.served, tc.excluded, tc.reply)
			if err != nil {
				t.Fatal(err)
			}
			if *got != tc.want {
				t.Errorf("got %+v, want %+v", *got, tc.want)
			}
		})
	}
}

// `{}` is a valid empty page for every rule: never nil, zero total, no next page.
func TestOffsetPaginationEmptyReplyForEveryRule(t *testing.T) {
	w := offsetWindow{limit: 20, offset: 40}
	for rule, next := range map[offsetRule]int64{
		ruleReplyOffset: 40, rulePlusRows: 40, rulePlusMinRowsLimit: 40, rulePlusServerRows: 40, rulePlusLimit: 60,
	} {
		got, err := offsetPagination(rule, w, 0, 0, offsetReply{})
		if err != nil {
			t.Fatalf("rule %d: %v", rule, err)
		}
		want := model.Pagination{Limit: 20, Offset: 40, NextOffset: next}
		if got == nil || *got != want {
			t.Errorf("rule %d: got %+v, want %+v", rule, got, want)
		}
	}
}

// The next offset saturates instead of wrapping, so no rule reports a page past MaxInt64.
func TestOffsetPaginationNeverOverflows(t *testing.T) {
	w := offsetWindow{limit: 100, offset: math.MaxInt64 - 10}
	total := strPtr(strconv.FormatInt(maxCount, 10))
	for _, rule := range []offsetRule{ruleReplyOffset, rulePlusRows, rulePlusMinRowsLimit, rulePlusServerRows, rulePlusLimit} {
		got, err := offsetPagination(rule, w, 50, 0, offsetReply{TotalItems: total})
		if err != nil {
			t.Fatalf("rule %d: %v", rule, err)
		}
		if got.NextOffset < w.offset || got.HasMore {
			t.Errorf("rule %d: got %+v, want a saturated next offset and no next page", rule, got)
		}
	}
}

func TestOffsetPaginationRejectsMalformedCounts(t *testing.T) {
	w := offsetWindow{limit: 20}
	for _, bad := range []string{"", "abc", "-1", "+1", " 1", "1 ", "020", "00", "1e3", "0x10", "1.5",
		"9007199254740992", "18446744073709551616", "99999999999999999999999"} {
		if _, err := offsetPagination(rulePlusRows, w, 0, 0, offsetReply{TotalItems: strPtr(bad)}); err == nil {
			t.Errorf("totalItems %q was accepted", bad)
		}
		if _, err := offsetPagination(ruleReplyOffset, w, 0, 0, offsetReply{TotalItems: strPtr("1"), Offset: strPtr(bad)}); err == nil {
			t.Errorf("reply offset %q was accepted", bad)
		}
	}
	if _, err := offsetPagination(offsetRule(0), w, 0, 0, offsetReply{}); err == nil {
		t.Error("an unknown rule must be an error, not a silent default")
	}
}

// A caller can tell a malformed reply from a transport or API error.
func TestMalformedPaginationIsMatchable(t *testing.T) {
	_, countErr := offsetPagination(rulePlusRows, offsetWindow{limit: 20}, 0, 0, offsetReply{TotalItems: strPtr("x")})
	_, bigErr := parseCount("total", strPtr("9007199254740992"))
	_, cursorErr := cursorPage(20, cursorReply{Cursor: &openapi.TgvalidatordResponseCursor{HasNext: boolPtr(true)}})
	for name, err := range map[string]error{"count": countErr, "count above 2^53-1": bigErr, "cursor": cursorErr} {
		if !errors.Is(err, model.ErrMalformedPagination) {
			t.Errorf("%s: %v does not match model.ErrMalformedPagination", name, err)
		}
	}
}

func TestParseCount(t *testing.T) {
	for in, want := range map[string]int64{"0": 0, "7": 7, "9007199254740991": maxCount} {
		got, err := parseCount("total", strPtr(in))
		if err != nil || got != want {
			t.Errorf("parseCount(%q) = %d, %v; want %d", in, got, err, want)
		}
	}
	if got, err := parseCount("total", nil); err != nil || got != 0 {
		t.Errorf("an absent count is 0, got %d, %v", got, err)
	}
}

func TestCursorPage(t *testing.T) {
	for _, tc := range []struct {
		name  string
		reply *openapi.TgvalidatordResponseCursor
		want  model.CursorPage
	}{
		{"no cursor", nil, model.CursorPage{PageSize: 20}},
		{"empty cursor", &openapi.TgvalidatordResponseCursor{}, model.CursorPage{PageSize: 20}},
		{"last page keeps no cursor", &openapi.TgvalidatordResponseCursor{CurrentPage: strPtr("abc")},
			model.CursorPage{PageSize: 20}},
		{"hasNext false keeps no cursor", &openapi.TgvalidatordResponseCursor{CurrentPage: strPtr("abc"), HasNext: boolPtr(false)},
			model.CursorPage{PageSize: 20}},
		{"next page", &openapi.TgvalidatordResponseCursor{CurrentPage: strPtr("ab+/cQ=="), HasNext: boolPtr(true), HasPrevious: boolPtr(true)},
			model.CursorPage{PageSize: 20, NextCursor: "ab+/cQ==", HasMore: true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := cursorPage(20, cursorReply{Cursor: tc.reply})
			if err != nil || got != tc.want {
				t.Errorf("got %+v, %v; want %+v", got, err, tc.want)
			}
		})
	}
	for _, bad := range []*openapi.TgvalidatordResponseCursor{
		{HasNext: boolPtr(true)},
		{HasNext: boolPtr(true), CurrentPage: strPtr("")},
	} {
		if _, err := cursorPage(20, cursorReply{Cursor: bad}); err == nil {
			t.Errorf("hasNext without a cursor was accepted: %+v", bad)
		}
	}
}

func TestCursorPageWithTotal(t *testing.T) {
	got, err := cursorPage(20, cursorReply{HasTotal: true})
	if err != nil || got.TotalItems == nil || *got.TotalItems != 0 {
		t.Errorf("an absent total on a list that has one is 0, got %+v, %v", got, err)
	}
	if got, _ := cursorPage(20, cursorReply{Total: strPtr("5")}); got.TotalItems != nil {
		t.Errorf("a list without a total must not report one, got %d", *got.TotalItems)
	}
	if _, err := cursorPage(20, cursorReply{HasTotal: true, Total: strPtr("x")}); err == nil {
		t.Error("a non-numeric total was accepted")
	}
}

func TestCursorPageTokenOnly(t *testing.T) {
	got, err := cursorPage(20, cursorReply{TokenOnly: true, Token: strPtr("ab+/cQ=="), HasTotal: true, Total: strPtr("10"), Excluded: 3})
	if err != nil || !got.HasMore || got.NextCursor != "ab+/cQ==" || got.TotalItems == nil || *got.TotalItems != 7 {
		t.Errorf("got %+v, %v", got, err)
	}
	got, err = cursorPage(20, cursorReply{TokenOnly: true, Token: strPtr(""), HasTotal: true, Total: strPtr("2"), Excluded: 5})
	if err != nil || got.HasMore || got.NextCursor != "" || *got.TotalItems != 0 {
		t.Errorf("an empty token ends the list and exclusions clamp at 0, got %+v, %v", got, err)
	}
	// A token-only reply has no cursor object; one that appears is not a continuation.
	got, err = cursorPage(20, cursorReply{TokenOnly: true, Cursor: &openapi.TgvalidatordResponseCursor{
		CurrentPage: strPtr("x"), HasNext: boolPtr(true)}})
	if err != nil || got.HasMore {
		t.Errorf("a token-only list continued from a cursor object: %+v, %v", got, err)
	}
	if _, err := cursorPage(20, cursorReply{TokenOnly: true, HasTotal: true, Total: strPtr("-3")}); err == nil {
		t.Error("a negative total was accepted")
	}
}

func TestChunkIDs(t *testing.T) {
	ids := func(n int) []string {
		out := make([]string, n)
		for i := range out {
			out[i] = strconv.Itoa(i)
		}
		return out
	}
	for n, want := range map[int][]int{0: nil, 1: {1}, 100: {100}, 101: {100, 1}, 250: {100, 100, 50}} {
		chunks := chunkIDs(ids(n), model.MaxPageSize)
		var sizes []int
		seen := 0
		for _, c := range chunks {
			sizes = append(sizes, len(c))
			for _, id := range c {
				if id != strconv.Itoa(seen) {
					t.Fatalf("n=%d: ids out of order or lost at %d", n, seen)
				}
				seen++
			}
		}
		if len(sizes) != len(want) || seen != n {
			t.Fatalf("n=%d: chunk sizes %v, want %v", n, sizes, want)
		}
		for i := range want {
			if sizes[i] != want[i] {
				t.Errorf("n=%d: chunk sizes %v, want %v", n, sizes, want)
			}
		}
	}
}
