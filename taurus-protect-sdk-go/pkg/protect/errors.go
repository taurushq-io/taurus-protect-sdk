// Package protect provides the Taurus-PROTECT SDK client.
package protect

import (
	"errors"
	"fmt"
	"time"
)

// APIError is the base error type for all Taurus-PROTECT API errors.
// It captures HTTP status codes, error messages, and provides helper methods
// for common error handling patterns.
type APIError struct {
	// Message is the human-readable error message.
	Message string
	// Code is the HTTP status code.
	Code int
	// Description is the error description from the API.
	Description string
	// ErrorCode is the application-specific error code.
	ErrorCode string
	// Err is the underlying error that caused this error.
	Err error
	// RetryAfter is the suggested retry delay for rate limit errors.
	RetryAfter time.Duration
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s (code=%d)", e.Message, e.Code)
	}
	if e.Description != "" {
		return fmt.Sprintf("%s (code=%d)", e.Description, e.Code)
	}
	return fmt.Sprintf("API error (code=%d)", e.Code)
}

// Unwrap returns the underlying error.
func (e *APIError) Unwrap() error {
	return e.Err
}

// IsRetryable returns true if the request might succeed on retry.
// This is true for 429 (rate limited) and 5xx (server errors).
func (e *APIError) IsRetryable() bool {
	return e.Code == 429 || e.Code >= 500
}

// IsClientError returns true if this is a client error (4xx).
func (e *APIError) IsClientError() bool {
	return e.Code >= 400 && e.Code < 500
}

// IsServerError returns true if this is a server error (5xx).
func (e *APIError) IsServerError() bool {
	return e.Code >= 500
}

// SuggestedRetryDelay returns the suggested delay before retrying.
// For rate limit errors, this returns the Retry-After value if available.
// For server errors, returns a default backoff value.
// For non-retryable errors, returns 0.
func (e *APIError) SuggestedRetryDelay() time.Duration {
	if e.Code == 429 {
		if e.RetryAfter > 0 {
			return e.RetryAfter
		}
		return time.Second // Default 1 second for rate limits
	}
	if e.Code >= 500 {
		return 5 * time.Second // Default 5 seconds for server errors
	}
	return 0
}

// Sentinel errors for type checking with errors.Is().
var (
	// ErrValidation indicates a 400 Bad Request error.
	ErrValidation = &APIError{Code: 400, Message: "validation error"}
	// ErrAuthentication indicates a 401 Unauthorized error.
	ErrAuthentication = &APIError{Code: 401, Message: "authentication error"}
	// ErrAuthorization indicates a 403 Forbidden error.
	ErrAuthorization = &APIError{Code: 403, Message: "authorization error"}
	// ErrNotFound indicates a 404 Not Found error.
	ErrNotFound = &APIError{Code: 404, Message: "not found"}
	// ErrRateLimit indicates a 429 Too Many Requests error.
	ErrRateLimit = &APIError{Code: 429, Message: "rate limit exceeded"}
	// ErrServer indicates a 5xx server error.
	ErrServer = &APIError{Code: 500, Message: "server error"}
)

// Is implements errors.Is for APIError.
func (e *APIError) Is(target error) bool {
	t, ok := target.(*APIError)
	if !ok {
		return false
	}
	// Match by HTTP status code category
	switch t.Code {
	case 400:
		return e.Code == 400
	case 401:
		return e.Code == 401
	case 403:
		return e.Code == 403
	case 404:
		return e.Code == 404
	case 429:
		return e.Code == 429
	case 500:
		return e.Code >= 500
	default:
		return e.Code == t.Code
	}
}

// ValidationError creates a new validation error (400).
func ValidationError(message string, err error) *APIError {
	return &APIError{
		Code:    400,
		Message: message,
		Err:     err,
	}
}

// AuthenticationError creates a new authentication error (401).
func AuthenticationError(message string, err error) *APIError {
	return &APIError{
		Code:    401,
		Message: message,
		Err:     err,
	}
}

// AuthorizationError creates a new authorization error (403).
func AuthorizationError(message string, err error) *APIError {
	return &APIError{
		Code:    403,
		Message: message,
		Err:     err,
	}
}

// NotFoundError creates a new not found error (404).
func NotFoundError(message string, err error) *APIError {
	return &APIError{
		Code:    404,
		Message: message,
		Err:     err,
	}
}

// RateLimitError creates a new rate limit error (429).
func RateLimitError(message string, retryAfter time.Duration, err error) *APIError {
	return &APIError{
		Code:       429,
		Message:    message,
		RetryAfter: retryAfter,
		Err:        err,
	}
}

// ServerError creates a new server error (5xx).
func ServerError(code int, message string, err error) *APIError {
	if code < 500 {
		code = 500
	}
	return &APIError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// IntegrityError indicates a cryptographic verification failure.
// This is a security-critical error that should never be retried.
type IntegrityError struct {
	Message string
	Err     error
}

func (e *IntegrityError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("integrity error: %s", e.Message)
	}
	return "integrity error"
}

func (e *IntegrityError) Unwrap() error {
	return e.Err
}

// ErrIntegrity is the sentinel error for integrity failures.
var ErrIntegrity = &IntegrityError{Message: "integrity verification failed"}

// Is implements errors.Is for IntegrityError.
func (e *IntegrityError) Is(target error) bool {
	_, ok := target.(*IntegrityError)
	return ok
}

// WhitelistError indicates a whitelist verification failure.
type WhitelistError struct {
	Message string
	Err     error
}

func (e *WhitelistError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("whitelist error: %s", e.Message)
	}
	return "whitelist error"
}

func (e *WhitelistError) Unwrap() error {
	return e.Err
}

// ErrWhitelist is the sentinel error for whitelist failures.
var ErrWhitelist = &WhitelistError{Message: "whitelist verification failed"}

// Is implements errors.Is for WhitelistError.
func (e *WhitelistError) Is(target error) bool {
	_, ok := target.(*WhitelistError)
	return ok
}

// RequestMetadataError indicates missing or invalid request metadata.
type RequestMetadataError struct {
	Message string
	Err     error
}

func (e *RequestMetadataError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("request metadata error: %s", e.Message)
	}
	return "request metadata error"
}

func (e *RequestMetadataError) Unwrap() error {
	return e.Err
}

// ErrRequestMetadata is the sentinel error for request metadata failures.
var ErrRequestMetadata = &RequestMetadataError{Message: "invalid request metadata"}

// Is implements errors.Is for RequestMetadataError.
func (e *RequestMetadataError) Is(target error) bool {
	_, ok := target.(*RequestMetadataError)
	return ok
}

// IsAPIError checks if the error is an APIError and returns it.
func IsAPIError(err error) (*APIError, bool) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr, true
	}
	return nil, false
}

// IsIntegrityError checks if the error is an IntegrityError.
func IsIntegrityError(err error) bool {
	return errors.Is(err, ErrIntegrity)
}

// IsWhitelistError checks if the error is a WhitelistError.
func IsWhitelistError(err error) bool {
	return errors.Is(err, ErrWhitelist)
}

// ConfigurationError indicates invalid SDK configuration.
// This error is raised when the SDK is configured with invalid or
// incompatible settings, such as missing credentials or invalid key formats.
type ConfigurationError struct {
	Message string
	Err     error
}

func (e *ConfigurationError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("configuration error: %s", e.Message)
	}
	return "configuration error"
}

func (e *ConfigurationError) Unwrap() error {
	return e.Err
}

// ErrConfiguration is the sentinel error for configuration failures.
var ErrConfiguration = &ConfigurationError{Message: "invalid configuration"}

// Is implements errors.Is for ConfigurationError.
func (e *ConfigurationError) Is(target error) bool {
	_, ok := target.(*ConfigurationError)
	return ok
}

// IsConfigurationError checks if the error is a ConfigurationError.
func IsConfigurationError(err error) bool {
	return errors.Is(err, ErrConfiguration)
}
