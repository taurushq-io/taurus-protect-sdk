package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

// ErrorMapper maps OpenAPI errors to domain errors.
type ErrorMapper struct{}

// NewErrorMapper creates a new ErrorMapper.
func NewErrorMapper() *ErrorMapper {
	return &ErrorMapper{}
}

// MapError converts an OpenAPI error to a domain error.
func (m *ErrorMapper) MapError(err error, resp *http.Response) error {
	if err == nil {
		return nil
	}

	// Try to extract OpenAPI error details
	if openAPIErr, ok := err.(*openapi.GenericOpenAPIError); ok {
		return mapOpenAPIError(openAPIErr, resp)
	}

	return err
}

// errorBody is the flat error payload Taurus-PROTECT returns on any status >= 400.
// "error" is a string (the gRPC status mnemonic), not a nested object; error_code
// is empty on some paths.
type errorBody struct {
	Err       string `json:"error"`
	Message   string `json:"message"`
	Code      int    `json:"code"`
	ErrorCode string `json:"error_code"`
}

// parseErrorBody extracts the server's error detail. GenericOpenAPIError.Error()
// carries only the status line, so without this the server's message is lost.
func parseErrorBody(raw []byte) (errorBody, bool) {
	if len(raw) == 0 {
		return errorBody{}, false
	}

	var body errorBody
	if jsonErr := json.Unmarshal(raw, &body); jsonErr != nil {
		return errorBody{}, false
	}
	if body.Message == "" && body.ErrorCode == "" {
		return errorBody{}, false
	}

	return body, true
}

func mapOpenAPIError(err *openapi.GenericOpenAPIError, resp *http.Response) error {
	if resp == nil {
		return err
	}

	return buildAPIError(resp.StatusCode, err.Error(), err.Body())
}

// buildAPIError is the pure mapping core. Split from mapOpenAPIError because
// GenericOpenAPIError has unexported fields in another package, so it cannot be
// constructed in a test.
func buildAPIError(code int, statusLine string, rawBody []byte) error {
	message := statusLine
	errorCode := ""
	if body, ok := parseErrorBody(rawBody); ok {
		if body.Message != "" {
			message = body.Message
		}
		errorCode = body.ErrorCode
	}

	// Create typed error based on status code
	switch {
	case code == 400:
		return &APIError{Code: code, Message: message, Description: "Bad Request", ErrorCode: errorCode}
	case code == 401:
		return &APIError{Code: code, Message: message, Description: "Unauthorized", ErrorCode: errorCode}
	case code == 403:
		return &AuthorizationError{
			APIError:      &APIError{Code: code, Message: message, Description: "Forbidden", ErrorCode: errorCode},
			RequiredRoles: ParseRequiredRoles(message),
		}
	case code == 404:
		return &APIError{Code: code, Message: message, Description: "Not Found", ErrorCode: errorCode}
	case code == 429:
		return &APIError{Code: code, Message: message, Description: "Rate Limited", ErrorCode: errorCode}
	case code >= 500:
		return &APIError{Code: code, Message: message, Description: "Server Error", ErrorCode: errorCode}
	default:
		return &APIError{Code: code, Message: message, ErrorCode: errorCode}
	}
}

// APIError represents an API error from the service layer. Every service funnels
// through the mapper into this type, so it is what callers actually receive; the
// protect package re-exports it as protect.APIError along with the Err* sentinels.
type APIError struct {
	Code        int
	Message     string
	Description string
	ErrorCode   string
	Err         error
	// RetryAfter is the suggested retry delay for rate-limit errors.
	RetryAfter time.Duration
}

func (e *APIError) Error() string {
	switch {
	case e.Description != "" && e.Message != "":
		return fmt.Sprintf("%s: %s (code=%d)", e.Description, e.Message, e.Code)
	case e.Description != "":
		return fmt.Sprintf("%s (code=%d)", e.Description, e.Code)
	case e.Message != "":
		return fmt.Sprintf("%s (code=%d)", e.Message, e.Code)
	default:
		// Never return just " (code=N)": an error string with no context is
		// useless in a log.
		return fmt.Sprintf("API error (code=%d)", e.Code)
	}
}

func (e *APIError) Unwrap() error {
	return e.Err
}

// Is matches by HTTP status code so the protect.Err* sentinels work with errors.Is.
// The 500 sentinel matches any 5xx; every other code matches exactly.
func (e *APIError) Is(target error) bool {
	t, ok := target.(*APIError)
	if !ok {
		return false
	}
	if t.Code == 500 {
		return e.Code >= 500
	}
	return e.Code == t.Code
}

// IsRetryable returns true if the request might succeed on retry: 429 (rate
// limited) and 5xx (server errors). Mirrors isRetryable() in the other three SDKs.
func (e *APIError) IsRetryable() bool {
	return e.Code == 429 || e.Code >= 500
}

// IsClientError returns true for a 4xx status.
func (e *APIError) IsClientError() bool {
	return e.Code >= 400 && e.Code < 500
}

// IsServerError returns true for a 5xx status.
func (e *APIError) IsServerError() bool {
	return e.Code >= 500
}

// SuggestedRetryDelay returns how long to wait before retrying: the server's
// Retry-After for rate limits when present, a default backoff otherwise, and 0
// for errors that are not retryable.
func (e *APIError) SuggestedRetryDelay() time.Duration {
	if e.Code == 429 {
		if e.RetryAfter > 0 {
			return e.RetryAfter
		}
		return time.Second
	}
	if e.Code >= 500 {
		return 5 * time.Second
	}
	return 0
}

// requiredRolesRe matches how Taurus-PROTECT reports a failed role check, for both
// its all-of and any-of checks.
//
// Unanchored on purpose: the server wraps the gRPC status, so this arrives as
// "pre-filter failed: … desc = one of the '…' role is required". Anchoring it
// matches nothing.
var requiredRolesRe = regexp.MustCompile(`one of the '([^']*)' role is required`)

// ParseRequiredRoles extracts the roles named in a 403 message, or an empty slice when
// the message is not a role check. Role names are lowercase alphanumeric, so " - " is an
// unambiguous separator.
//
// Exported to match Java AuthorizationException.parseRequiredRoles, Python
// parse_required_roles and TS parseRequiredRoles: a caller that has the server message
// from somewhere other than a mapped error (a log line, a wrapped gRPC status) can still
// recover the roles. Returns empty rather than nil so callers need not nil-check, which
// is what the other three do.
func ParseRequiredRoles(message string) []string {
	roles := make([]string, 0, 1)

	match := requiredRolesRe.FindStringSubmatch(message)
	if match == nil {
		return roles
	}

	for _, role := range strings.Split(match[1], " - ") {
		if trimmed := strings.TrimSpace(role); trimmed != "" {
			roles = append(roles, trimmed)
		}
	}

	return roles
}

// AuthorizationError is returned for HTTP 403. RequiredRoles lets a caller say
// which role to ask for instead of just "forbidden".
type AuthorizationError struct {
	*APIError

	// Empty when the denial was not role-based. One entry is a required role;
	// several mean any one of them suffices.
	RequiredRoles []string
}

func (e *AuthorizationError) Error() string {
	if len(e.RequiredRoles) == 0 {
		return e.APIError.Error()
	}

	return fmt.Sprintf("%s; requires one of: %s", e.APIError.Error(), strings.Join(e.RequiredRoles, ", "))
}

// Unwrap returns the embedded *APIError. Required: without it the promoted
// (*APIError).Unwrap returns APIError.Err (nil), so errors.As(err, &APIError{})
// would report false for every 403.
func (e *AuthorizationError) Unwrap() error {
	return e.APIError
}
