package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// newFilterTestClient wires an OpenAPI client at an httptest server that records
// the query string of the request it receives and returns an empty result body.
func newFilterTestClient(t *testing.T, capture *string) (*openapi.APIClient, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*capture = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	cfg := openapi.NewConfiguration()
	cfg.Servers = []openapi.ServerConfiguration{{URL: srv.URL}}
	cfg.HTTPClient = srv.Client()
	return openapi.NewAPIClient(cfg), srv.Close
}

func TestListTransactionsWiresDateAndNetworkFilters(t *testing.T) {
	var rawQuery string
	apiClient, closeFn := newFilterTestClient(t, &rawQuery)
	defer closeFn()

	from := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 3, 4, 0, 0, 0, 0, time.UTC)
	_, _, err := NewTransactionService(apiClient).ListTransactions(context.Background(), &model.ListTransactionsOptions{
		Network:  "sepolia",
		FromDate: &from,
		ToDate:   &to,
	})
	if err != nil {
		t.Fatalf("ListTransactions: %v", err)
	}
	for _, want := range []string{"sepolia", "2024-01-02", "2024-03-04"} {
		if !strings.Contains(rawQuery, want) {
			t.Errorf("query %q missing %q", rawQuery, want)
		}
	}
}

func TestListRequestsWiresDateFilters(t *testing.T) {
	var rawQuery string
	apiClient, closeFn := newFilterTestClient(t, &rawQuery)
	defer closeFn()

	from := time.Date(2024, 5, 6, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 7, 8, 0, 0, 0, 0, time.UTC)
	_, err := NewRequestService(apiClient).ListRequests(context.Background(), &model.ListRequestsOptions{
		Statuses: []string{"PENDING"},
		FromDate: &from,
		ToDate:   &to,
	})
	if err != nil {
		t.Fatalf("ListRequests: %v", err)
	}
	for _, want := range []string{"2024-05-06", "2024-07-08"} {
		if !strings.Contains(rawQuery, want) {
			t.Errorf("query %q missing %q", rawQuery, want)
		}
	}
}
