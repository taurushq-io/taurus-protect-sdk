package service

import (
	"fmt"
	"math"
	"strconv"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// The pagination contract shared by every list method in all four SDKs (repo-root CLAUDE.md,
// "Pagination (cross-SDK)"). scripts/resources/pagination-vectors.json pins the helpers below and
// list-request-vectors.json pins what each list method sends.
//
//	options ─▶ resolve*Window (page size 0→20, >100 or <0 → error) ─▶ request
//	reply ────▶ offsetPagination(rule) / cursorPage ─▶ model value
//
// validatord omits zero values, so an absent totalItems/offset/hasNext/next means 0/false/none and
// `{}` is a valid empty page.

// maxCount is the largest count a reply may carry. 2^53-1 keeps all four SDKs, TypeScript's
// number included, reading the same value.
const maxCount = 1<<53 - 1

// pageRequestNext is the pageRequest a continuation sends. validatord matches it exactly.
const pageRequestNext = "NEXT"

// resolvePageSize applies the shared page-size rule to a list option.
func resolvePageSize(option string, size int64) (int64, error) {
	return resolveSize(option, size, model.MaxPageSize)
}

// resolveSize maps 0 to model.DefaultPageSize and rejects a negative size or one above max.
// max 0 means the endpoint has no SDK maximum.
func resolveSize(option string, size, max int64) (int64, error) {
	if size < 0 || (max > 0 && size > max) {
		if max > 0 {
			return 0, fmt.Errorf("invalid %s %d: must be between 1 and %d (0 selects the default %d)",
				option, size, max, model.DefaultPageSize)
		}
		return 0, fmt.Errorf("invalid %s %d: must not be negative (0 selects the default %d)",
			option, size, model.DefaultPageSize)
	}
	if size == 0 {
		return model.DefaultPageSize, nil
	}
	return size, nil
}

// offsetWindow is the limit/offset pair an offset list sends.
type offsetWindow struct {
	limit  int64
	offset int64
}

// resolveOffsetWindow validates Limit/Offset before any request is sent.
func resolveOffsetWindow(limit, offset int64) (offsetWindow, error) {
	size, err := resolvePageSize("Limit", limit)
	if err != nil {
		return offsetWindow{}, err
	}
	if offset < 0 {
		return offsetWindow{}, fmt.Errorf("invalid Offset %d: must not be negative", offset)
	}
	return offsetWindow{limit: size, offset: offset}, nil
}

// offsetQueryRequest is satisfied by every generated request builder taking limit/offset.
type offsetQueryRequest[T any] interface {
	Limit(string) T
	Offset(string) T
}

// applyOffsetWindow always sends the limit, and the offset when it is not 0.
func applyOffsetWindow[T offsetQueryRequest[T]](req T, w offsetWindow) T {
	req = req.Limit(strconv.FormatInt(w.limit, 10))
	if w.offset > 0 {
		req = req.Offset(strconv.FormatInt(w.offset, 10))
	}
	return req
}

// offsetRule is how an endpoint's next page offset is derived; each offset endpoint has one.
type offsetRule int

const (
	// ruleReplyOffset: the reply's offset IS the next page's offset (offset + rows when absent).
	// WalletService_GetWalletsV2, WalletService_GetAddresses.
	ruleReplyOffset offsetRule = iota + 1
	// rulePlusRows: offset + rows. Transactions, fee payers, actions.
	rulePlusRows
	// rulePlusMinRowsLimit: offset + min(rows, limit); a synthetic daemon user or technical
	// group can be appended beyond the limit. Users, groups.
	rulePlusMinRowsLimit
	// rulePlusServerRows: offset + the rows the SERVER returned, before SDK exclusions; rows the
	// server dropped for a bad signature are gone for good. Whitelisted addresses.
	rulePlusServerRows
	// rulePlusLimit: offset + limit; rows the server skipped keep their SQL slot, so a short page
	// is not the end. Whitelisted contracts (assets).
	rulePlusLimit
)

// offsetReply is what an offset reply says about paging.
type offsetReply struct {
	TotalItems *string
	Offset     *string
}

// offsetPagination builds the page value every offset list returns.
//
// servedRows is the number of rows the server returned; excluded the rows the SDK withheld from
// them. TotalItems is reduced by the exclusions, NextOffset and HasMore never are: pagination
// position is a server-side fact. HasMore also requires progress, so a trailing empty page under
// an upper-bound total ends a walk.
func offsetPagination(rule offsetRule, w offsetWindow, servedRows, excluded int, reply offsetReply) (*model.Pagination, error) {
	total, err := parseCount("totalItems", reply.TotalItems)
	if err != nil {
		return nil, err
	}

	rows := int64(servedRows)
	var next int64
	switch rule {
	case ruleReplyOffset:
		if reply.Offset == nil {
			next = addSaturated(w.offset, rows)
		} else if next, err = parseCount("offset", reply.Offset); err != nil {
			return nil, err
		}
	case rulePlusRows, rulePlusServerRows:
		next = addSaturated(w.offset, rows)
	case rulePlusMinRowsLimit:
		next = addSaturated(w.offset, min(rows, w.limit))
	case rulePlusLimit:
		next = addSaturated(w.offset, w.limit)
	default:
		return nil, fmt.Errorf("unknown offset pagination rule %d", rule)
	}

	reported := total - int64(excluded)
	if reported < 0 {
		reported = 0
	}
	return &model.Pagination{
		Limit:      w.limit,
		Offset:     w.offset,
		TotalItems: reported,
		NextOffset: next,
		HasMore:    next > w.offset && next < total,
	}, nil
}

// cursorWindow is what a cursor list sends: always a page size, plus either a continuation
// (the caller's Cursor as currentPage with pageRequest NEXT) or the low-level
// CurrentPage/PageRequest options as given.
type cursorWindow struct {
	pageSize    int64
	currentPage string
	pageRequest string
}

// resolveCursorWindow validates the paging options of a cursor list. Cursor always means "the
// page after this one", so combining it with a low-level option is rejected rather than letting
// one of them silently lose.
func resolveCursorWindow(pageSize int64, cursor, currentPage, pageRequest string) (cursorWindow, error) {
	size, err := resolvePageSize("PageSize", pageSize)
	if err != nil {
		return cursorWindow{}, err
	}
	if cursor == "" {
		return cursorWindow{pageSize: size, currentPage: currentPage, pageRequest: pageRequest}, nil
	}
	if currentPage != "" {
		return cursorWindow{}, fmt.Errorf("invalid options: Cursor and CurrentPage cannot both be set; Cursor continues from a previous page's NextCursor")
	}
	if pageRequest != "" {
		return cursorWindow{}, fmt.Errorf("invalid options: Cursor and PageRequest cannot both be set; Cursor always requests the %s page", pageRequestNext)
	}
	return cursorWindow{pageSize: size, currentPage: cursor, pageRequest: pageRequestNext}, nil
}

// cursorQueryRequest is satisfied by every generated request builder taking cursor.* parameters.
type cursorQueryRequest[T any] interface {
	CursorCurrentPage(string) T
	CursorPageRequest(string) T
	CursorPageSize(string) T
}

// applyCursorQuery sends the window as cursor.* query parameters.
func applyCursorQuery[T cursorQueryRequest[T]](req T, w cursorWindow) T {
	req = req.CursorPageSize(strconv.FormatInt(w.pageSize, 10))
	if w.currentPage != "" {
		req = req.CursorCurrentPage(w.currentPage)
	}
	if w.pageRequest != "" {
		req = req.CursorPageRequest(w.pageRequest)
	}
	return req
}

// requestCursorQueryRequest is satisfied by generated builders taking requestCursor.* parameters.
type requestCursorQueryRequest[T any] interface {
	RequestCursorCurrentPage(string) T
	RequestCursorPageRequest(string) T
	RequestCursorPageSize(string) T
}

// applyRequestCursorQuery sends the window as requestCursor.* query parameters.
func applyRequestCursorQuery[T requestCursorQueryRequest[T]](req T, w cursorWindow) T {
	req = req.RequestCursorPageSize(strconv.FormatInt(w.pageSize, 10))
	if w.currentPage != "" {
		req = req.RequestCursorCurrentPage(w.currentPage)
	}
	if w.pageRequest != "" {
		req = req.RequestCursorPageRequest(w.pageRequest)
	}
	return req
}

// body renders the window for endpoints that take the cursor in the request body.
func (w cursorWindow) body() *openapi.TgvalidatordRequestCursor {
	size := strconv.FormatInt(w.pageSize, 10)
	c := &openapi.TgvalidatordRequestCursor{PageSize: &size}
	if w.currentPage != "" {
		currentPage := w.currentPage
		c.CurrentPage = &currentPage
	}
	if w.pageRequest != "" {
		pageRequest := w.pageRequest
		c.PageRequest = &pageRequest
	}
	return c
}

// cursorReply is what a cursor-paged reply says about paging.
type cursorReply struct {
	// Cursor is the reply cursor of the cursor.* and requestCursor lists.
	Cursor *openapi.TgvalidatordResponseCursor
	// TokenOnly marks the lists that page by a bare token instead (wallet tokens, rules
	// history); Token is the reply's next-page token.
	TokenOnly bool
	Token     *string
	// HasTotal marks the lists whose reply carries a count; an absent Total is then 0.
	HasTotal bool
	Total    *string
	// Excluded is the number of rows the SDK withheld; it reduces the total (never below 0).
	Excluded int
}

// cursorPage builds the page value every cursor list returns.
//
// A cursor list continues with the reply's currentPage only when hasNext is set: sending NEXT
// past the end is a 400 on the v2 lists and silently returns page 1 on the keyset lists, so a
// cursor is never derived any other way. A token-only list continues with its token, present
// exactly when another page exists.
func cursorPage(pageSize int64, reply cursorReply) (model.CursorPage, error) {
	page := model.CursorPage{PageSize: pageSize}
	if reply.HasTotal {
		total, err := parseCount("total", reply.Total)
		if err != nil {
			return model.CursorPage{}, err
		}
		total = max(total-int64(reply.Excluded), 0)
		page.TotalItems = &total
	}

	if reply.TokenOnly {
		if reply.Token != nil && *reply.Token != "" {
			page.HasMore = true
			page.NextCursor = *reply.Token
		}
		return page, nil
	}

	c := reply.Cursor
	if c == nil || c.HasNext == nil || !*c.HasNext {
		return page, nil
	}
	if c.CurrentPage == nil || *c.CurrentPage == "" {
		return model.CursorPage{}, fmt.Errorf("%w: the reply cursor has hasNext but no currentPage", model.ErrMalformedPagination)
	}
	page.HasMore = true
	page.NextCursor = *c.CurrentPage
	return page, nil
}

// parseCount reads a count from a reply. Absent means 0; anything but a canonical decimal in
// [0, 2^53-1] is an error, because a count silently read as 0 ends a walk early.
func parseCount(field string, value *string) (int64, error) {
	if value == nil {
		return 0, nil
	}
	s := *value
	if !isCanonicalDecimal(s) {
		return 0, fmt.Errorf("%w: reply %s %q is not a decimal count", model.ErrMalformedPagination, field, s)
	}
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil || n > maxCount {
		return 0, fmt.Errorf("%w: reply %s %q is above %d", model.ErrMalformedPagination, field, s, uint64(maxCount))
	}
	return int64(n), nil
}

// isCanonicalDecimal accepts digits only, without a sign or a leading zero.
func isCanonicalDecimal(s string) bool {
	if s == "" || (len(s) > 1 && s[0] == '0') {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// addSaturated adds a non-negative b to a, stopping at math.MaxInt64.
func addSaturated(a, b int64) int64 {
	if b > 0 && a > math.MaxInt64-b {
		return math.MaxInt64
	}
	return a + b
}

// chunkIDs splits ids into consecutive batches of at most size, so an id-filtered internal read
// stays within both the page-size maximum and the endpoint's own id cap.
func chunkIDs(ids []string, size int) [][]string {
	var chunks [][]string
	for start := 0; start < len(ids); start += size {
		end := min(start+size, len(ids))
		chunks = append(chunks, ids[start:end])
	}
	return chunks
}
