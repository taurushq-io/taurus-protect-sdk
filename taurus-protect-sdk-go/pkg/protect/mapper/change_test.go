package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestChangeFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordChange
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns change with zero values",
			dto:  &openapi.TgvalidatordChange{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordChange {
				id := "change-123"
				tenantID := "tenant-456"
				creatorID := "user-789"
				creatorExternalID := "external@example.com"
				action := "create"
				entityID := "entity-001"
				entityUUID := "550e8400-e29b-41d4-a716-446655440000"
				entity := "user"
				comment := "Creating a new user"
				now := time.Now()
				changes := map[string]string{
					"firstname": "John",
					"lastname":  "Doe",
				}
				return &openapi.TgvalidatordChange{
					Id:                &id,
					TenantId:          &tenantID,
					CreatorId:         &creatorID,
					CreatorExternalId: &creatorExternalID,
					Action:            &action,
					EntityId:          &entityID,
					EntityUUID:        &entityUUID,
					Entity:            &entity,
					Comment:           &comment,
					CreationDate:      &now,
					Changes:           &changes,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ChangeFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ChangeFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ChangeFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.TenantId != nil && got.TenantID != *tt.dto.TenantId {
				t.Errorf("TenantID = %v, want %v", got.TenantID, *tt.dto.TenantId)
			}
			if tt.dto.CreatorId != nil && got.CreatorID != *tt.dto.CreatorId {
				t.Errorf("CreatorID = %v, want %v", got.CreatorID, *tt.dto.CreatorId)
			}
			if tt.dto.CreatorExternalId != nil && got.CreatorExternalID != *tt.dto.CreatorExternalId {
				t.Errorf("CreatorExternalID = %v, want %v", got.CreatorExternalID, *tt.dto.CreatorExternalId)
			}
			if tt.dto.Action != nil && got.Action != *tt.dto.Action {
				t.Errorf("Action = %v, want %v", got.Action, *tt.dto.Action)
			}
			if tt.dto.EntityId != nil && got.EntityID != *tt.dto.EntityId {
				t.Errorf("EntityID = %v, want %v", got.EntityID, *tt.dto.EntityId)
			}
			if tt.dto.EntityUUID != nil && got.EntityUUID != *tt.dto.EntityUUID {
				t.Errorf("EntityUUID = %v, want %v", got.EntityUUID, *tt.dto.EntityUUID)
			}
			if tt.dto.Entity != nil && got.Entity != *tt.dto.Entity {
				t.Errorf("Entity = %v, want %v", got.Entity, *tt.dto.Entity)
			}
			if tt.dto.Comment != nil && got.Comment != *tt.dto.Comment {
				t.Errorf("Comment = %v, want %v", got.Comment, *tt.dto.Comment)
			}
			if tt.dto.CreationDate != nil && !got.CreationDate.Equal(*tt.dto.CreationDate) {
				t.Errorf("CreationDate = %v, want %v", got.CreationDate, *tt.dto.CreationDate)
			}
			if tt.dto.Changes != nil {
				if got.Changes == nil {
					t.Error("Changes should not be nil when DTO has changes")
				} else {
					for k, v := range *tt.dto.Changes {
						if got.Changes[k] != v {
							t.Errorf("Changes[%s] = %v, want %v", k, got.Changes[k], v)
						}
					}
				}
			}
		})
	}
}

func TestChangesFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordChange
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordChange{},
			want: 0,
		},
		{
			name: "converts multiple changes",
			dtos: func() []openapi.TgvalidatordChange {
				id1 := "change-1"
				id2 := "change-2"
				id3 := "change-3"
				return []openapi.TgvalidatordChange{
					{Id: &id1},
					{Id: &id2},
					{Id: &id3},
				}
			}(),
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ChangesFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("ChangesFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("ChangesFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestCursorPaginationFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordResponseCursor
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns cursor with zero values",
			dto:  &openapi.TgvalidatordResponseCursor{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordResponseCursor {
				currentPage := "eyJpZCI6MTIzfQ=="
				hasPrevious := true
				hasNext := true
				return &openapi.TgvalidatordResponseCursor{
					CurrentPage: &currentPage,
					HasPrevious: &hasPrevious,
					HasNext:     &hasNext,
				}
			}(),
		},
		{
			name: "no next page",
			dto: func() *openapi.TgvalidatordResponseCursor {
				currentPage := "eyJpZCI6OTk5fQ=="
				hasPrevious := true
				hasNext := false
				return &openapi.TgvalidatordResponseCursor{
					CurrentPage: &currentPage,
					HasPrevious: &hasPrevious,
					HasNext:     &hasNext,
				}
			}(),
		},
		{
			name: "first page",
			dto: func() *openapi.TgvalidatordResponseCursor {
				currentPage := "eyJpZCI6MX0="
				hasPrevious := false
				hasNext := true
				return &openapi.TgvalidatordResponseCursor{
					CurrentPage: &currentPage,
					HasPrevious: &hasPrevious,
					HasNext:     &hasNext,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CursorPaginationFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("CursorPaginationFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("CursorPaginationFromDTO() returned nil for non-nil input")
			}
			if tt.dto.CurrentPage != nil && got.CurrentPage != *tt.dto.CurrentPage {
				t.Errorf("CurrentPage = %v, want %v", got.CurrentPage, *tt.dto.CurrentPage)
			}
			if tt.dto.HasPrevious != nil && got.HasPrevious != *tt.dto.HasPrevious {
				t.Errorf("HasPrevious = %v, want %v", got.HasPrevious, *tt.dto.HasPrevious)
			}
			if tt.dto.HasNext != nil && got.HasNext != *tt.dto.HasNext {
				t.Errorf("HasNext = %v, want %v", got.HasNext, *tt.dto.HasNext)
			}
		})
	}
}

func TestCreateChangeResultFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordCreateChangeResult
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns result with empty ID",
			dto:  &openapi.TgvalidatordCreateChangeResult{},
		},
		{
			name: "complete DTO maps ID",
			dto: func() *openapi.TgvalidatordCreateChangeResult {
				id := "change-new-123"
				return &openapi.TgvalidatordCreateChangeResult{
					Id: &id,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateChangeResultFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("CreateChangeResultFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("CreateChangeResultFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
		})
	}
}

func TestChangeFromDTO_ChangesMapCopy(t *testing.T) {
	// Verify that the Changes map is copied, not shared
	originalChanges := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	dto := &openapi.TgvalidatordChange{
		Changes: &originalChanges,
	}

	got := ChangeFromDTO(dto)

	// Modify the result
	got.Changes["key1"] = "modified"

	// Original should be unchanged
	if originalChanges["key1"] != "value1" {
		t.Error("ChangeFromDTO() should copy the Changes map, not share it")
	}
}

func TestChangeFromDTO_NilFields(t *testing.T) {
	// Test that nil pointer fields are handled gracefully
	dto := &openapi.TgvalidatordChange{
		// All pointer fields are nil
	}

	got := ChangeFromDTO(dto)

	if got == nil {
		t.Fatal("ChangeFromDTO() returned nil for non-nil input")
	}

	// All string fields should be empty
	if got.ID != "" {
		t.Errorf("ID should be empty string, got %v", got.ID)
	}
	if got.TenantID != "" {
		t.Errorf("TenantID should be empty string, got %v", got.TenantID)
	}
	if got.CreatorID != "" {
		t.Errorf("CreatorID should be empty string, got %v", got.CreatorID)
	}
	if got.Action != "" {
		t.Errorf("Action should be empty string, got %v", got.Action)
	}
	if got.Entity != "" {
		t.Errorf("Entity should be empty string, got %v", got.Entity)
	}
	if got.Changes != nil {
		t.Errorf("Changes should be nil, got %v", got.Changes)
	}
}
