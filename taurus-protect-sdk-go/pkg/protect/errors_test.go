package protect

import (
	"errors"
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/service"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *APIError
		want string
	}{
		{
			name: "with message",
			err:  &APIError{Message: "test error", Code: 400},
			want: "test error (code=400)",
		},
		{
			name: "with description",
			err:  &APIError{Description: "test description", Code: 500},
			want: "test description (code=500)",
		},
		{
			name: "without message or description",
			err:  &APIError{Code: 404},
			want: "API error (code=404)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("APIError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_IsRetryable(t *testing.T) {
	tests := []struct {
		name string
		code int
		want bool
	}{
		{"400", 400, false},
		{"401", 401, false},
		{"403", 403, false},
		{"404", 404, false},
		{"429", 429, true},
		{"500", 500, true},
		{"502", 502, true},
		{"503", 503, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &APIError{Code: tt.code}
			if got := err.IsRetryable(); got != tt.want {
				t.Errorf("APIError.IsRetryable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_IsClientError(t *testing.T) {
	tests := []struct {
		code int
		want bool
	}{
		{399, false},
		{400, true},
		{404, true},
		{499, true},
		{500, false},
	}

	for _, tt := range tests {
		err := &APIError{Code: tt.code}
		if got := err.IsClientError(); got != tt.want {
			t.Errorf("APIError{Code: %d}.IsClientError() = %v, want %v", tt.code, got, tt.want)
		}
	}
}

func TestAPIError_IsServerError(t *testing.T) {
	tests := []struct {
		code int
		want bool
	}{
		{499, false},
		{500, true},
		{502, true},
		{503, true},
	}

	for _, tt := range tests {
		err := &APIError{Code: tt.code}
		if got := err.IsServerError(); got != tt.want {
			t.Errorf("APIError{Code: %d}.IsServerError() = %v, want %v", tt.code, got, tt.want)
		}
	}
}

func TestAPIError_SuggestedRetryDelay(t *testing.T) {
	tests := []struct {
		name    string
		err     *APIError
		wantMin time.Duration
	}{
		{
			name:    "rate limit default",
			err:     &APIError{Code: 429},
			wantMin: time.Second,
		},
		{
			name:    "rate limit with retry-after",
			err:     &APIError{Code: 429, RetryAfter: 5 * time.Second},
			wantMin: 5 * time.Second,
		},
		{
			name:    "server error",
			err:     &APIError{Code: 500},
			wantMin: 5 * time.Second,
		},
		{
			name:    "client error",
			err:     &APIError{Code: 400},
			wantMin: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.SuggestedRetryDelay(); got < tt.wantMin {
				t.Errorf("APIError.SuggestedRetryDelay() = %v, want >= %v", got, tt.wantMin)
			}
		})
	}
}

func TestAPIError_Is(t *testing.T) {
	tests := []struct {
		name   string
		err    *APIError
		target error
		want   bool
	}{
		{
			name:   "validation error matches",
			err:    &APIError{Code: 400},
			target: ErrValidation,
			want:   true,
		},
		{
			name:   "authentication error matches",
			err:    &APIError{Code: 401},
			target: ErrAuthentication,
			want:   true,
		},
		{
			name:   "authorization error matches",
			err:    &APIError{Code: 403},
			target: ErrAuthorization,
			want:   true,
		},
		{
			name:   "not found error matches",
			err:    &APIError{Code: 404},
			target: ErrNotFound,
			want:   true,
		},
		{
			name:   "rate limit error matches",
			err:    &APIError{Code: 429},
			target: ErrRateLimit,
			want:   true,
		},
		{
			name:   "server error matches 500",
			err:    &APIError{Code: 500},
			target: ErrServer,
			want:   true,
		},
		{
			name:   "server error matches 502",
			err:    &APIError{Code: 502},
			target: ErrServer,
			want:   true,
		},
		{
			name:   "different codes don't match",
			err:    &APIError{Code: 400},
			target: ErrNotFound,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := errors.Is(tt.err, tt.target); got != tt.want {
				t.Errorf("errors.Is(%v, %v) = %v, want %v", tt.err, tt.target, got, tt.want)
			}
		})
	}
}

func TestAPIError_Unwrap(t *testing.T) {
	inner := errors.New("inner error")
	err := &APIError{
		Code:    500,
		Message: "outer error",
		Err:     inner,
	}

	if !errors.Is(err, inner) {
		t.Error("APIError should unwrap to inner error")
	}
}

func TestIntegrityError(t *testing.T) {
	err := &IntegrityError{Message: "signature mismatch"}

	if !errors.Is(err, ErrIntegrity) {
		t.Error("IntegrityError should match ErrIntegrity")
	}

	expected := "integrity error: signature mismatch"
	if err.Error() != expected {
		t.Errorf("IntegrityError.Error() = %v, want %v", err.Error(), expected)
	}
}

func TestWhitelistError(t *testing.T) {
	err := &WhitelistError{Message: "address not whitelisted"}

	if !errors.Is(err, ErrWhitelist) {
		t.Error("WhitelistError should match ErrWhitelist")
	}

	expected := "whitelist error: address not whitelisted"
	if err.Error() != expected {
		t.Errorf("WhitelistError.Error() = %v, want %v", err.Error(), expected)
	}
}

// This package used to declare its own APIError / IntegrityError / WhitelistError
// structs alongside the ones the SDK actually returns, so every check below silently
// matched nothing: services return service.APIError and the verification helpers return
// model.IntegrityError. The types are aliases now. These assertions pin that, because a
// regression re-introducing a parallel type disables every documented errors.Is /
// errors.As branch without failing anything.
func TestSentinelsMatchTheErrorsServicesActuallyReturn(t *testing.T) {
	t.Run("status sentinels match a mapper-shaped error", func(t *testing.T) {
		notFound := &service.APIError{Code: 404, Message: "wallet not found"}
		if !errors.Is(notFound, ErrNotFound) {
			t.Error("errors.Is(serviceErr, ErrNotFound) = false, want true")
		}
		if errors.Is(notFound, ErrValidation) {
			t.Error("a 404 must not match ErrValidation")
		}
		for _, code := range []int{500, 502, 503} {
			if !errors.Is(&service.APIError{Code: code}, ErrServer) {
				t.Errorf("code %d must match ErrServer", code)
			}
		}
	})

	t.Run("rate limit exposes retry hints", func(t *testing.T) {
		rateLimited := &service.APIError{Code: 429, RetryAfter: 5 * time.Second}
		if !errors.Is(rateLimited, ErrRateLimit) {
			t.Error("a 429 must match ErrRateLimit")
		}
		if !rateLimited.IsRetryable() {
			t.Error("a 429 must be retryable")
		}
		if got := rateLimited.SuggestedRetryDelay(); got != 5*time.Second {
			t.Errorf("SuggestedRetryDelay() = %v, want 5s", got)
		}
	})

	t.Run("the typed 403 is reachable as both itself and the base error", func(t *testing.T) {
		authz := &service.AuthorizationError{
			APIError:      &service.APIError{Code: 403, Message: "denied"},
			RequiredRoles: []string{"admin"},
		}
		got, ok := IsAuthorizationError(authz)
		if !ok {
			t.Fatal("IsAuthorizationError() = false, want true")
		}
		if len(got.RequiredRoles) != 1 || got.RequiredRoles[0] != "admin" {
			t.Errorf("RequiredRoles = %v, want [admin]", got.RequiredRoles)
		}
		if _, ok := IsAPIError(authz); !ok {
			t.Error("IsAPIError() must also match the typed 403")
		}
		if !errors.Is(authz, ErrAuthorization) {
			t.Error("a 403 must match ErrAuthorization")
		}
	})

	t.Run("verification failures match their sentinels", func(t *testing.T) {
		if !IsIntegrityError(&model.IntegrityError{Message: "hash mismatch"}) {
			t.Error("IsIntegrityError() = false for a verifier-produced error")
		}
		whitelist := &model.WhitelistError{Message: "not whitelisted"}
		if !IsWhitelistError(whitelist) {
			t.Error("IsWhitelistError() = false for a verifier-produced error")
		}
		if IsIntegrityError(whitelist) {
			t.Error("a whitelist error must not match ErrIntegrity")
		}
	})
}

func TestIsAPIError(t *testing.T) {
	apiErr := &APIError{Code: 404, Message: "not found"}

	got, ok := IsAPIError(apiErr)
	if !ok {
		t.Error("IsAPIError should return true for APIError")
	}
	if got != apiErr {
		t.Error("IsAPIError should return the same APIError")
	}

	_, ok = IsAPIError(errors.New("regular error"))
	if ok {
		t.Error("IsAPIError should return false for non-APIError")
	}
}
