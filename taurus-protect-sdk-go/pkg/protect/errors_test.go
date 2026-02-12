package protect

import (
	"errors"
	"testing"
	"time"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     *APIError
		want    string
	}{
		{
			name:    "with message",
			err:     &APIError{Message: "test error", Code: 400},
			want:    "test error (code=400)",
		},
		{
			name:    "with description",
			err:     &APIError{Description: "test description", Code: 500},
			want:    "test description (code=500)",
		},
		{
			name:    "without message or description",
			err:     &APIError{Code: 404},
			want:    "API error (code=404)",
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
		name       string
		err        *APIError
		wantMin    time.Duration
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

func TestErrorFactoryFunctions(t *testing.T) {
	t.Run("ValidationError", func(t *testing.T) {
		err := ValidationError("invalid input", nil)
		if err.Code != 400 {
			t.Errorf("ValidationError code = %d, want 400", err.Code)
		}
	})

	t.Run("AuthenticationError", func(t *testing.T) {
		err := AuthenticationError("invalid token", nil)
		if err.Code != 401 {
			t.Errorf("AuthenticationError code = %d, want 401", err.Code)
		}
	})

	t.Run("AuthorizationError", func(t *testing.T) {
		err := AuthorizationError("access denied", nil)
		if err.Code != 403 {
			t.Errorf("AuthorizationError code = %d, want 403", err.Code)
		}
	})

	t.Run("NotFoundError", func(t *testing.T) {
		err := NotFoundError("wallet not found", nil)
		if err.Code != 404 {
			t.Errorf("NotFoundError code = %d, want 404", err.Code)
		}
	})

	t.Run("RateLimitError", func(t *testing.T) {
		err := RateLimitError("too many requests", 5*time.Second, nil)
		if err.Code != 429 {
			t.Errorf("RateLimitError code = %d, want 429", err.Code)
		}
		if err.RetryAfter != 5*time.Second {
			t.Errorf("RateLimitError RetryAfter = %v, want 5s", err.RetryAfter)
		}
	})

	t.Run("ServerError", func(t *testing.T) {
		err := ServerError(502, "bad gateway", nil)
		if err.Code != 502 {
			t.Errorf("ServerError code = %d, want 502", err.Code)
		}
	})

	t.Run("ServerError clamps code", func(t *testing.T) {
		err := ServerError(400, "should be 500", nil)
		if err.Code != 500 {
			t.Errorf("ServerError should clamp code to 500, got %d", err.Code)
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
