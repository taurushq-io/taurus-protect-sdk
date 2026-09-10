package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// recordingLogger captures what the SDK reports, so a test can assert an exclusion
// was surfaced rather than silently dropped.
type recordingLogger struct {
	warns  []string
	fields [][]Field
	ctxs   []context.Context
}

func (l *recordingLogger) Debug(context.Context, string, ...Field) {}
func (l *recordingLogger) Info(context.Context, string, ...Field)  {}
func (l *recordingLogger) Error(context.Context, string, ...Field) {}
func (l *recordingLogger) Warn(ctx context.Context, msg string, f ...Field) {
	l.warns = append(l.warns, msg)
	l.fields = append(l.fields, f)
	l.ctxs = append(l.ctxs, ctx)
}

func requestsListBody(t *testing.T) string {
	t.Helper()
	good := `[{"key":"currency","value":"BTC"}]`
	tampered := `[{"key":"currency","value":"ETH"}]`

	rows := []map[string]any{
		// verifies
		{"id": "1", "metadata": map[string]any{"hash": crypto.CalculateHexHash(good), "payloadAsString": good}},
		// payload altered, hash left alone — the documented attack
		{"id": "2", "metadata": map[string]any{"hash": crypto.CalculateHexHash(good), "payloadAsString": tampered}},
		// nothing to verify yet: an early-status request
		{"id": "3", "metadata": map[string]any{}},
	}
	b, err := json.Marshal(map[string]any{"result": rows})
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func listService(t *testing.T, body string, logger Logger) *RequestService {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, body)
	}))
	t.Cleanup(srv.Close)

	cfg := openapi.NewConfiguration()
	cfg.Servers = []openapi.ServerConfiguration{{URL: srv.URL}}
	cfg.HTTPClient = srv.Client()
	return NewRequestService(openapi.NewAPIClient(cfg), WithServiceLogger(logger))
}

// The defect this fixes: ListRequests returned every row without verifying any of
// them, so a tampered payload reached the caller — and, through mcpd, an agent's
// context — with no error and no flag.
func TestListRequestsExcludesUnverifiedRows(t *testing.T) {
	for _, tc := range []struct {
		name string
		call func(*RequestService) (*model.RequestResult, error)
	}{
		{"ListRequests", func(s *RequestService) (*model.RequestResult, error) {
			return s.ListRequests(context.Background(), nil)
		}},
		{"ListRequestsForApproval", func(s *RequestService) (*model.RequestResult, error) {
			return s.ListRequestsForApproval(context.Background(), nil)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			logger := &recordingLogger{}
			result, err := tc.call(listService(t, requestsListBody(t), logger))
			if err != nil {
				t.Fatalf("one bad row must not fail the call: %v", err)
			}

			var ids []string
			for _, r := range result.Requests {
				ids = append(ids, r.ID)
			}
			if len(result.Requests) != 2 {
				t.Fatalf("expected the tampered row dropped, got ids %v", ids)
			}
			if ids[0] != "1" || ids[1] != "3" {
				t.Errorf("kept ids = %v, want [1 3]", ids)
			}

			// The verified row carries usable structure; the metadata-less one does not.
			if !result.Requests[0].Metadata.HashVerified {
				t.Error("row 1 should be marked verified")
			}
			got, _ := result.Requests[0].Metadata.GetMetadataCurrency()
			if got != "BTC" {
				t.Errorf("row 1 currency = %q, want BTC", got)
			}
			if result.Requests[1].Metadata.HashVerified {
				t.Error("row 3 has nothing to verify and must not be marked verified")
			}

			// A shortened list must never read as a complete one.
			if len(result.ExcludedUnverified) != 1 || result.ExcludedUnverified[0] != "2" {
				t.Errorf("ExcludedUnverified = %v, want [2]", result.ExcludedUnverified)
			}
			if len(logger.warns) != 1 {
				t.Fatalf("expected exactly one exclusion logged, got %d", len(logger.warns))
			}
			if logger.ctxs[0] == nil {
				t.Error("ctx must reach the logger, or a consumer loses correlation")
			}
			for _, f := range logger.fields[0] {
				if s, ok := f.Value.(string); ok && (s == `[{"key":"currency","value":"ETH"}]`) {
					t.Errorf("log field %q carries the payload", f.Key)
				}
			}
		})
	}
}

// Without a logger the exclusion is silent, but it must still happen — and must not
// panic, which a nil interface would.
func TestListRequestsWithoutLoggerStillExcludes(t *testing.T) {
	svc := listService(t, requestsListBody(t), nil)
	result, err := svc.ListRequests(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Requests) != 2 || len(result.ExcludedUnverified) != 1 {
		t.Errorf("kept %d rows, excluded %v", len(result.Requests), result.ExcludedUnverified)
	}
}
