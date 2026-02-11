package model

import "fmt"

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

// APIError represents an error response from the Taurus-PROTECT API.
// It supports errors.As for type-based error handling, and errors.Is via
// Unwrap for cause chain matching.
type APIError struct {
	// StatusCode is the HTTP status code.
	StatusCode int
	// ErrorCode is the API error code.
	ErrorCode string
	// Message is the human-readable error message.
	Message string
	// Err is the underlying error, if any.
	Err error
}

func (e *APIError) Error() string {
	if e.ErrorCode != "" {
		return fmt.Sprintf("API error %d (%s): %s", e.StatusCode, e.ErrorCode, e.Message)
	}
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

func (e *APIError) Unwrap() error {
	return e.Err
}

// IsRetryable returns true if the error is retryable (429 or 5xx).
func (e *APIError) IsRetryable() bool {
	return e.StatusCode == 429 || (e.StatusCode >= 500 && e.StatusCode < 600)
}

// IsClientError returns true for 4xx errors.
func (e *APIError) IsClientError() bool {
	return e.StatusCode >= 400 && e.StatusCode < 500
}

// IsServerError returns true for 5xx errors.
func (e *APIError) IsServerError() bool {
	return e.StatusCode >= 500 && e.StatusCode < 600
}

// ValidationError represents a 400 Bad Request error.
type ValidationError struct {
	*APIError
}

// AuthenticationError represents a 401 Unauthorized error.
type AuthenticationError struct {
	*APIError
}

// AuthorizationError represents a 403 Forbidden error.
type AuthorizationError struct {
	*APIError
}

// NotFoundError represents a 404 Not Found error.
type NotFoundError struct {
	*APIError
}

// RateLimitError represents a 429 Too Many Requests error.
type RateLimitError struct {
	*APIError
}

// ServerError represents a 5xx server error.
type ServerError struct {
	*APIError
}

// ConfigurationError represents a client configuration error.
type ConfigurationError struct {
	Message string
	Err     error
}

func (e *ConfigurationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("configuration error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("configuration error: %s", e.Message)
}

func (e *ConfigurationError) Unwrap() error {
	return e.Err
}

// RequestMetadataError represents an error related to request metadata.
type RequestMetadataError struct {
	Message string
	Err     error
}

func (e *RequestMetadataError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("request metadata error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("request metadata error: %s", e.Message)
}

func (e *RequestMetadataError) Unwrap() error {
	return e.Err
}

// NewAPIError creates the appropriate typed error based on the HTTP status code.
func NewAPIError(statusCode int, errorCode string, message string, err error) error {
	base := &APIError{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
		Err:        err,
	}
	switch {
	case statusCode == 400:
		return &ValidationError{APIError: base}
	case statusCode == 401:
		return &AuthenticationError{APIError: base}
	case statusCode == 403:
		return &AuthorizationError{APIError: base}
	case statusCode == 404:
		return &NotFoundError{APIError: base}
	case statusCode == 429:
		return &RateLimitError{APIError: base}
	case statusCode >= 500:
		return &ServerError{APIError: base}
	default:
		return base
	}
}
