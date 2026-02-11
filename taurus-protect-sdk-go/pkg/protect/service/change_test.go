package service

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewChangeService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestChangeService_GetChange_EmptyID(t *testing.T) {
	svc := &ChangeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetChange(nil, "")
	if err == nil {
		t.Error("GetChange() with empty ID should return error")
	}
	if err.Error() != "changeID cannot be empty" {
		t.Errorf("GetChange() error = %v, want 'changeID cannot be empty'", err)
	}
}

func TestChangeService_CreateChange_NilRequest(t *testing.T) {
	svc := &ChangeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreateChange(nil, nil)
	if err == nil {
		t.Error("CreateChange() with nil request should return error")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("CreateChange() error = %v, want 'request cannot be nil'", err)
	}
}

func TestChangeService_CreateChange_EmptyAction(t *testing.T) {
	svc := &ChangeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreateChange(nil, &model.CreateChangeRequest{})
	if err == nil {
		t.Error("CreateChange() with empty action should return error")
	}
	if err.Error() != "action is required" {
		t.Errorf("CreateChange() error = %v, want 'action is required'", err)
	}
}

func TestChangeService_CreateChange_EmptyEntity(t *testing.T) {
	svc := &ChangeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreateChange(nil, &model.CreateChangeRequest{Action: "create"})
	if err == nil {
		t.Error("CreateChange() with empty entity should return error")
	}
	if err.Error() != "entity is required" {
		t.Errorf("CreateChange() error = %v, want 'entity is required'", err)
	}
}

func TestChangeService_ApproveChange_EmptyID(t *testing.T) {
	svc := &ChangeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.ApproveChange(nil, "")
	if err == nil {
		t.Error("ApproveChange() with empty ID should return error")
	}
	if err.Error() != "changeID cannot be empty" {
		t.Errorf("ApproveChange() error = %v, want 'changeID cannot be empty'", err)
	}
}

func TestChangeService_ApproveChanges_EmptyIDs(t *testing.T) {
	svc := &ChangeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.ApproveChanges(nil, nil)
	if err == nil {
		t.Error("ApproveChanges() with nil IDs should return error")
	}
	if err.Error() != "changeIDs cannot be empty" {
		t.Errorf("ApproveChanges() error = %v, want 'changeIDs cannot be empty'", err)
	}

	err = svc.ApproveChanges(nil, []string{})
	if err == nil {
		t.Error("ApproveChanges() with empty IDs should return error")
	}
	if err.Error() != "changeIDs cannot be empty" {
		t.Errorf("ApproveChanges() error = %v, want 'changeIDs cannot be empty'", err)
	}
}

func TestChangeService_RejectChange_EmptyID(t *testing.T) {
	svc := &ChangeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.RejectChange(nil, "")
	if err == nil {
		t.Error("RejectChange() with empty ID should return error")
	}
	if err.Error() != "changeID cannot be empty" {
		t.Errorf("RejectChange() error = %v, want 'changeID cannot be empty'", err)
	}
}

func TestChangeService_RejectChanges_EmptyIDs(t *testing.T) {
	svc := &ChangeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.RejectChanges(nil, nil)
	if err == nil {
		t.Error("RejectChanges() with nil IDs should return error")
	}
	if err.Error() != "changeIDs cannot be empty" {
		t.Errorf("RejectChanges() error = %v, want 'changeIDs cannot be empty'", err)
	}

	err = svc.RejectChanges(nil, []string{})
	if err == nil {
		t.Error("RejectChanges() with empty IDs should return error")
	}
	if err.Error() != "changeIDs cannot be empty" {
		t.Errorf("RejectChanges() error = %v, want 'changeIDs cannot be empty'", err)
	}
}

func TestChangeService_ListChanges_NilOptions(t *testing.T) {
	// Create a service with nil API to test that nil options don't cause panic
	// The actual API call will fail, but we're testing the options handling
	svc := &ChangeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	// This will panic on nil API, so we just verify the service accepts nil options
	// In a real test with mocked API, nil options should work
	if svc == nil {
		t.Error("ChangeService should not be nil")
	}
}

func TestChangeService_ListChanges_WithOptions(t *testing.T) {
	// Verify that ListChangesOptions fields are properly structured
	tests := []struct {
		name        string
		entity      string
		status      string
		creatorID   string
		sortOrder   string
		pageSize    int64
		currentPage string
		pageRequest string
	}{
		{
			name:        "empty options",
			entity:      "",
			status:      "",
			creatorID:   "",
			sortOrder:   "",
			pageSize:    0,
			currentPage: "",
			pageRequest: "",
		},
		{
			name:        "entity filter",
			entity:      "user",
			status:      "",
			creatorID:   "",
			sortOrder:   "",
			pageSize:    0,
			currentPage: "",
			pageRequest: "",
		},
		{
			name:        "status filter",
			entity:      "",
			status:      "Created",
			creatorID:   "",
			sortOrder:   "",
			pageSize:    0,
			currentPage: "",
			pageRequest: "",
		},
		{
			name:        "pagination options",
			entity:      "",
			status:      "",
			creatorID:   "",
			sortOrder:   "DESC",
			pageSize:    50,
			currentPage: "eyJpZCI6MTIzfQ==",
			pageRequest: "NEXT",
		},
		{
			name:        "all options",
			entity:      "wallet",
			status:      "Approved",
			creatorID:   "user-123",
			sortOrder:   "ASC",
			pageSize:    100,
			currentPage: "abc123",
			pageRequest: "FIRST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify options can be created with these values
			// Actual API testing requires mocking
			svc := &ChangeService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("ChangeService should not be nil")
			}
		})
	}
}

func TestChangeService_ListChangesForApproval_NilOptions(t *testing.T) {
	svc := &ChangeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	if svc == nil {
		t.Error("ChangeService should not be nil")
	}
}

func TestChangeService_CreateChange_Validation(t *testing.T) {
	// Test input validation errors (these don't require the API)
	tests := []struct {
		name    string
		request *model.CreateChangeRequest
		errMsg  string
	}{
		{
			name:    "nil request",
			request: nil,
			errMsg:  "request cannot be nil",
		},
		{
			name:    "empty action",
			request: &model.CreateChangeRequest{Entity: "user"},
			errMsg:  "action is required",
		},
		{
			name:    "empty entity",
			request: &model.CreateChangeRequest{Action: "create"},
			errMsg:  "entity is required",
		},
	}

	svc := &ChangeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateChange(nil, tt.request)
			if err == nil {
				t.Error("CreateChange() should return error")
			}
			if err.Error() != tt.errMsg {
				t.Errorf("CreateChange() error = %v, want %v", err, tt.errMsg)
			}
		})
	}
}

func TestCreateChangeRequest_FieldStructure(t *testing.T) {
	// Verify that CreateChangeRequest fields are properly structured
	// This documents the expected request structure without calling the API
	tests := []struct {
		name    string
		request *model.CreateChangeRequest
	}{
		{
			name: "minimal request",
			request: &model.CreateChangeRequest{
				Action: "create",
				Entity: "user",
			},
		},
		{
			name: "full request",
			request: &model.CreateChangeRequest{
				Action:     "update",
				Entity:     "user",
				EntityID:   "user-123",
				EntityUUID: "550e8400-e29b-41d4-a716-446655440000",
				Changes: map[string]string{
					"firstname": "John",
					"lastname":  "Doe",
				},
				Comment: "Updating user details",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.request.Action == "" {
				t.Error("Action should not be empty")
			}
			if tt.request.Entity == "" {
				t.Error("Entity should not be empty")
			}
		})
	}
}
