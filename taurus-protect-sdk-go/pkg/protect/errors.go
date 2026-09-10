// Package protect provides the Taurus-PROTECT SDK client.
package protect

import (
	"errors"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/service"
)

// The error types callers actually receive, re-exported here so the package a caller
// already imports is enough to match them.
//
// These are aliases, not copies: every service funnels through the error mapper into
// service.APIError, and the verification helpers return model.IntegrityError /
// model.WhitelistError. This package previously declared its own parallel structs with
// the same names, so errors.Is/As against them could never match a real error and every
// documented error-handling branch was dead code.
type (
	// APIError is an error response from the Taurus-PROTECT API.
	APIError = service.APIError
	// AuthorizationError is a 403 that also names the roles which would satisfy the
	// failed check (empty when the denial was not role-based).
	AuthorizationError = service.AuthorizationError
	// IntegrityError is a cryptographic verification failure. Never retry it.
	IntegrityError = model.IntegrityError
	// WhitelistError is a whitelist verification failure.
	WhitelistError = model.WhitelistError
	// ConfigurationError is invalid SDK configuration.
	ConfigurationError = model.ConfigurationError
	// RequestMetadataError is missing or invalid request metadata.
	RequestMetadataError = model.RequestMetadataError
)

// Sentinel errors for use with errors.Is. Matching is by HTTP status code
// (see (*service.APIError).Is), so the message on these values is irrelevant.
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
	// ErrServer matches any 5xx server error.
	ErrServer = &APIError{Code: 500, Message: "server error"}

	// ErrIntegrity is the sentinel for integrity failures.
	ErrIntegrity = &IntegrityError{Message: "integrity verification failed"}
	// ErrWhitelist is the sentinel for whitelist failures.
	ErrWhitelist = &WhitelistError{Message: "whitelist verification failed"}
)

// IsAPIError reports whether err is an API error, returning it for access to the
// status code, retry hints and error code.
func IsAPIError(err error) (*APIError, bool) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr, true
	}
	return nil, false
}

// IsAuthorizationError reports whether err is a 403, returning it so the caller can
// read RequiredRoles.
func IsAuthorizationError(err error) (*AuthorizationError, bool) {
	var authErr *AuthorizationError
	if errors.As(err, &authErr) {
		return authErr, true
	}
	return nil, false
}

// IsIntegrityError reports whether err is a cryptographic verification failure.
func IsIntegrityError(err error) bool {
	return errors.Is(err, ErrIntegrity)
}

// IsWhitelistError reports whether err is a whitelist verification failure.
func IsWhitelistError(err error) bool {
	return errors.Is(err, ErrWhitelist)
}
