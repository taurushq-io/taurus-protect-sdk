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

// Is matches any IntegrityError so the protect.ErrIntegrity sentinel works with
// errors.Is regardless of the message.
func (e *IntegrityError) Is(target error) bool {
	_, ok := target.(*IntegrityError)
	return ok
}

// ContainerIntegrityError means the governance rules container itself carries
// something this SDK version cannot interpret, so no verification decision taken
// against it can be trusted.
//
// Deliberately distinct from a per-row WhitelistError. A caller listing whitelisted
// addresses excludes a row that fails its own integrity check and keeps the rest —
// but an uninterpretable container invalidates EVERY row judged against it, so the
// whole call must fail instead. Collapsing the two lets a schema-newer container
// empty a whitelist while reporting success.
type ContainerIntegrityError struct {
	Message string
	Err     error
}

func (e *ContainerIntegrityError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("rules container integrity error: %s", e.Message)
	}
	return "rules container integrity error"
}

func (e *ContainerIntegrityError) Unwrap() error {
	return e.Err
}

// Is matches its own type and also *IntegrityError, so a container-level failure
// still satisfies the protect.ErrIntegrity sentinel while remaining separately
// detectable with errors.As.
func (e *ContainerIntegrityError) Is(target error) bool {
	if _, ok := target.(*ContainerIntegrityError); ok {
		return true
	}
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

// Is matches any WhitelistError so the protect.ErrWhitelist sentinel works with
// errors.Is regardless of the message.
func (e *WhitelistError) Is(target error) bool {
	_, ok := target.(*WhitelistError)
	return ok
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
