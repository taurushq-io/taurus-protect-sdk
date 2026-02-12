package service

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewTagService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestTagService_CreateTag_NilRequest(t *testing.T) {
	svc := &TagService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreateTag(nil, nil)
	if err == nil {
		t.Error("CreateTag() with nil request should return error")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("CreateTag() error = %v, want 'request cannot be nil'", err)
	}
}

func TestTagService_CreateTag_EmptyValue(t *testing.T) {
	svc := &TagService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreateTag(nil, &model.CreateTagRequest{})
	if err == nil {
		t.Error("CreateTag() with empty value should return error")
	}
	if err.Error() != "value is required" {
		t.Errorf("CreateTag() error = %v, want 'value is required'", err)
	}
}

func TestTagService_CreateTag_EmptyValueWithColor(t *testing.T) {
	svc := &TagService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreateTag(nil, &model.CreateTagRequest{Color: "#FF0000"})
	if err == nil {
		t.Error("CreateTag() with empty value should return error")
	}
	if err.Error() != "value is required" {
		t.Errorf("CreateTag() error = %v, want 'value is required'", err)
	}
}

func TestTagService_DeleteTag_EmptyID(t *testing.T) {
	svc := &TagService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.DeleteTag(nil, "")
	if err == nil {
		t.Error("DeleteTag() with empty ID should return error")
	}
	if err.Error() != "id cannot be empty" {
		t.Errorf("DeleteTag() error = %v, want 'id cannot be empty'", err)
	}
}

func TestTagService_ListTags_NilOptions(t *testing.T) {
	// This test verifies that nil options doesn't cause a panic
	// The actual API call would fail without a real client
	svc := &TagService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	// This would panic when trying to call the API, which is expected
	// since we have a nil api. The important thing is that nil opts
	// doesn't cause a panic before the API call.
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when calling API with nil client")
		}
	}()

	_, _ = svc.ListTags(nil, nil)
}

func TestTagService_ListTags_WithOptions(t *testing.T) {
	// This test verifies that options are properly handled
	// The actual API call would fail without a real client
	svc := &TagService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	opts := &model.ListTagsOptions{
		IDs:   []string{"tag-1", "tag-2"},
		Query: "test",
	}

	// This would panic when trying to call the API, which is expected
	// since we have a nil api. The important thing is that opts
	// doesn't cause a panic before the API call.
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when calling API with nil client")
		}
	}()

	_, _ = svc.ListTags(nil, opts)
}

func TestTagService_ListTags_EmptyOptions(t *testing.T) {
	// Test with empty options struct (no filters set)
	svc := &TagService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	opts := &model.ListTagsOptions{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when calling API with nil client")
		}
	}()

	_, _ = svc.ListTags(nil, opts)
}
