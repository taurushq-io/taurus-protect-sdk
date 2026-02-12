// Package mapper provides functions to convert between OpenAPI DTOs and domain models.
package mapper

import "time"

// Helper functions for safe pointer dereferencing

// safeString returns an empty string if the pointer is nil, otherwise returns the dereferenced value.
func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// safeBool returns false if the pointer is nil, otherwise returns the dereferenced value.
func safeBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// safeInt64 returns 0 if the pointer is nil, otherwise returns the dereferenced value.
func safeInt64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

// safeFloat32 safely dereferences a float32 pointer, returning 0 if nil.
func safeFloat32(f *float32) float32 {
	if f == nil {
		return 0
	}
	return *f
}

// safeTime returns a zero time if the pointer is nil, otherwise returns the dereferenced value.
func safeTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

// Helper functions for creating pointers

// stringPtr returns a pointer to the string, or nil if the string is empty.
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// boolPtr returns a pointer to the bool.
func boolPtr(b bool) *bool {
	return &b
}

// int64Ptr returns a pointer to the int64.
func int64Ptr(i int64) *int64 {
	return &i
}
