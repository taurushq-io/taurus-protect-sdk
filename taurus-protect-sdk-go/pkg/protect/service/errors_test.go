package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

// authzVectorsRelPath points at the monorepo-shared vectors file consumed by all
// four SDK test suites (same shared-resources pattern as the governance cell
// vectors), so cross-SDK role-parsing parity is CI-enforced.
const authzVectorsRelPath = "../../../../scripts/resources/authorization-error-vectors.json"

type authzVector struct {
	Description   string   `json:"description"`
	Message       string   `json:"message"`
	ExpectedRoles []string `json:"expected_roles"`
}

// TestAuthorizationErrorVectors validates role extraction against the shared
// cross-SDK vectors. Every SDK must derive the same roles from the same server
// message, so a wording change is caught in one place for all four.
func TestAuthorizationErrorVectors(t *testing.T) {
	path := filepath.Clean(authzVectorsRelPath)

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read shared vectors file %s: %v", path, err)
	}

	var vectors []authzVector
	if err := json.Unmarshal(raw, &vectors); err != nil {
		t.Fatalf("cannot parse shared vectors file: %v", err)
	}
	if len(vectors) == 0 {
		t.Fatal("shared vectors file is empty")
	}

	for _, v := range vectors {
		t.Run(v.Description, func(t *testing.T) {
			got := ParseRequiredRoles(v.Message)

			// An empty expectation means "not a role check" — nil and empty are
			// both acceptable there, so compare lengths before contents.
			if len(got) != len(v.ExpectedRoles) {
				t.Fatalf("ParseRequiredRoles(%q) = %v, want %v", v.Message, got, v.ExpectedRoles)
			}
			if len(got) > 0 && !reflect.DeepEqual(got, v.ExpectedRoles) {
				t.Fatalf("ParseRequiredRoles(%q) = %v, want %v", v.Message, got, v.ExpectedRoles)
			}
		})
	}
}

const rolesBody = `{"error":"PermissionDenied","message":"one of the 'admin - adminreadonly' role is required","code":403,"error_code":""}`

// Guards the explicit AuthorizationError.Unwrap: without it errors.As walks past the
// embedded value, so callers matching on APIError alone break on every 403.
func TestAuthorizationErrorAsBothTargets(t *testing.T) {
	err := buildAPIError(http.StatusForbidden, "403 Forbidden", []byte(rolesBody))

	var authzErr *AuthorizationError
	if !errors.As(err, &authzErr) {
		t.Fatalf("errors.As(*AuthorizationError) = false for %T", err)
	}
	if want := []string{"admin", "adminreadonly"}; !reflect.DeepEqual(authzErr.RequiredRoles, want) {
		t.Fatalf("RequiredRoles = %v, want %v", authzErr.RequiredRoles, want)
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatal("errors.As(*APIError) = false — the explicit AuthorizationError.Unwrap is missing")
	}
	if apiErr.Code != http.StatusForbidden {
		t.Fatalf("APIError.Code = %d, want 403", apiErr.Code)
	}

	// The cause chain must still terminate rather than loop.
	if cause := errors.Unwrap(apiErr); cause != nil {
		t.Fatalf("errors.Unwrap(*APIError) = %v, want nil", cause)
	}
}

func TestBuildAPIErrorUsesResponseBody(t *testing.T) {
	tests := []struct {
		name          string
		status        int
		statusLine    string
		body          string
		wantMessage   string
		wantErrorCode string
		wantDesc      string
		wantRoles     []string
	}{
		{
			name:        "403 with roles",
			status:      http.StatusForbidden,
			statusLine:  "403 Forbidden",
			body:        rolesBody,
			wantMessage: "one of the 'admin - adminreadonly' role is required",
			wantDesc:    "Forbidden",
			wantRoles:   []string{"admin", "adminreadonly"},
		},
		{
			name:        "403 without roles stays typed with no roles",
			status:      http.StatusForbidden,
			statusLine:  "403 Forbidden",
			body:        `{"error":"PermissionDenied","message":"This endpoint has been disabled","code":403}`,
			wantMessage: "This endpoint has been disabled",
			wantDesc:    "Forbidden",
			// Empty, not nil: ParseRequiredRoles always returns a slice so callers
			// need not nil-check, matching Java/Python/TS.
			wantRoles: []string{},
		},
		{
			name:          "400 carries error_code",
			status:        http.StatusBadRequest,
			statusLine:    "400 Bad Request",
			body:          `{"error":"InvalidArgument","message":"Invalid parameters - bad","code":400,"error_code":"E900-INVALID_PARAMETERS"}`,
			wantMessage:   "Invalid parameters - bad",
			wantErrorCode: "E900-INVALID_PARAMETERS",
			wantDesc:      "Bad Request",
		},
		{
			name:        "401 unauthorized",
			status:      http.StatusUnauthorized,
			statusLine:  "401 Unauthorized",
			body:        `{"error":"Unauthenticated","message":"token expired","code":401}`,
			wantMessage: "token expired",
			wantDesc:    "Unauthorized",
		},
		{
			name:        "429 rate limited",
			status:      http.StatusTooManyRequests,
			statusLine:  "429 Too Many Requests",
			body:        `{"error":"ResourceExhausted","message":"slow down","code":429}`,
			wantMessage: "slow down",
			wantDesc:    "Rate Limited",
		},
		{
			name:        "unparseable body falls back to the status line",
			status:      http.StatusInternalServerError,
			statusLine:  "500 Internal Server Error",
			body:        "<html>gateway exploded</html>",
			wantMessage: "500 Internal Server Error",
			wantDesc:    "Server Error",
		},
		{
			name:        "empty body falls back to the status line",
			status:      http.StatusNotFound,
			statusLine:  "404 Not Found",
			body:        "",
			wantMessage: "404 Not Found",
			wantDesc:    "Not Found",
		},
		{
			name:        "well-formed json without message or error_code is ignored",
			status:      http.StatusNotFound,
			statusLine:  "404 Not Found",
			body:        `{"something":"else"}`,
			wantMessage: "404 Not Found",
			wantDesc:    "Not Found",
		},
		{
			name:        "unmapped status keeps an empty description",
			status:      http.StatusTeapot,
			statusLine:  "418 I'm a teapot",
			body:        `{"message":"short and stout"}`,
			wantMessage: "short and stout",
			wantDesc:    "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := buildAPIError(tc.status, tc.statusLine, []byte(tc.body))

			var apiErr *APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("errors.As(*APIError) = false for %T", err)
			}
			if apiErr.Message != tc.wantMessage {
				t.Fatalf("Message = %q, want %q", apiErr.Message, tc.wantMessage)
			}
			if apiErr.ErrorCode != tc.wantErrorCode {
				t.Fatalf("ErrorCode = %q, want %q", apiErr.ErrorCode, tc.wantErrorCode)
			}
			if apiErr.Description != tc.wantDesc {
				t.Fatalf("Description = %q, want %q", apiErr.Description, tc.wantDesc)
			}
			if apiErr.Code != tc.status {
				t.Fatalf("Code = %d, want %d", apiErr.Code, tc.status)
			}

			var authzErr *AuthorizationError
			isAuthz := errors.As(err, &authzErr)
			if wantAuthz := tc.status == http.StatusForbidden; isAuthz != wantAuthz {
				t.Fatalf("errors.As(*AuthorizationError) = %v, want %v", isAuthz, wantAuthz)
			}
			if isAuthz && !reflect.DeepEqual(authzErr.RequiredRoles, tc.wantRoles) {
				t.Fatalf("RequiredRoles = %v, want %v", authzErr.RequiredRoles, tc.wantRoles)
			}
		})
	}
}

// Pins the status -> Description matrix through the adapter, including the
// code >= 500 branch for non-500 codes.
func TestMapOpenAPIErrorStatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantDescr  string
	}{
		{"bad request", 400, "Bad Request"},
		{"unauthorized", 401, "Unauthorized"},
		{"forbidden", 403, "Forbidden"},
		{"not found", 404, "Not Found"},
		{"rate limited", 429, "Rate Limited"},
		{"server error 500", 500, "Server Error"},
		{"server error 502", 502, "Server Error"},
		{"server error 503", 503, "Server Error"},
		{"unknown status", 418, ""}, // I'm a teapot - no special description
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mapOpenAPIError(&openapi.GenericOpenAPIError{}, &http.Response{StatusCode: tt.statusCode})

			var apiErr *APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("errors.As(*APIError) = false for %T", err)
			}
			if apiErr.Code != tt.statusCode {
				t.Errorf("APIError.Code = %v, want %v", apiErr.Code, tt.statusCode)
			}
			if apiErr.Description != tt.wantDescr {
				t.Errorf("APIError.Description = %v, want %v", apiErr.Description, tt.wantDescr)
			}
		})
	}
}

func TestMapOpenAPIErrorNilResponse(t *testing.T) {
	mockErr := &openapi.GenericOpenAPIError{}

	// Should return the original error when response is nil
	if err := mapOpenAPIError(mockErr, nil); err != error(mockErr) {
		t.Errorf("mapOpenAPIError(err, nil) should return original error, got %v", err)
	}
}

func TestAuthorizationErrorMessage(t *testing.T) {
	withRoles := &AuthorizationError{
		APIError:      &APIError{Code: 403, Message: "denied", Description: "Forbidden"},
		RequiredRoles: []string{"admin", "adminreadonly"},
	}
	if got := withRoles.Error(); !strings.Contains(got, "requires one of: admin, adminreadonly") {
		t.Fatalf("Error() = %q, want it to name the roles", got)
	}

	withoutRoles := &AuthorizationError{
		APIError: &APIError{Code: 403, Message: "denied", Description: "Forbidden"},
	}
	if got, want := withoutRoles.Error(), "Forbidden: denied (code=403)"; got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}

	noDescription := &APIError{Code: 418, Message: "short and stout"}
	if got, want := noDescription.Error(), "short and stout (code=418)"; got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestMapErrorPassesThroughNonOpenAPIErrors(t *testing.T) {
	// A transport failure is not a GenericOpenAPIError; it must come back
	// untouched so callers can inspect it (and so nothing invents a status).
	cause := errors.New("dial tcp: connection refused")
	if got := NewErrorMapper().MapError(cause, nil); !errors.Is(got, cause) {
		t.Fatalf("MapError(%v) = %v, want the original error", cause, got)
	}

	if got := NewErrorMapper().MapError(nil, nil); got != nil {
		t.Fatalf("MapError(nil) = %v, want nil", got)
	}
}

// Proves the wiring the pure tests cannot: a real 403 from the generated client
// carries its response body through to the typed error.
func TestMapErrorThroughRealClient(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(rolesBody))
	}))
	defer srv.Close()

	cfg := openapi.NewConfiguration()
	cfg.Servers = []openapi.ServerConfiguration{{URL: srv.URL}}
	cfg.HTTPClient = srv.Client()

	_, err := NewBusinessRuleService(openapi.NewAPIClient(cfg)).ListBusinessRules(context.Background(), nil)
	if err == nil {
		t.Fatal("ListBusinessRules against a 403 server returned no error")
	}

	var authzErr *AuthorizationError
	if !errors.As(err, &authzErr) {
		t.Fatalf("errors.As(*AuthorizationError) = false for %T (%v)", err, err)
	}
	if want := []string{"admin", "adminreadonly"}; !reflect.DeepEqual(authzErr.RequiredRoles, want) {
		t.Fatalf("RequiredRoles = %v, want %v — the response body did not reach the mapper", authzErr.RequiredRoles, want)
	}
}
