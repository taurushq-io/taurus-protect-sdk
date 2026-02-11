package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestGroupFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordInternalGroup
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns group with zero values",
			dto:  &openapi.TgvalidatordInternalGroup{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordInternalGroup {
				id := "group-123"
				tenantId := "tenant-456"
				externalGroupId := "ext-group-789"
				name := "Test Group"
				email := "group@example.com"
				description := "A test group"
				enforcedInRules := true
				creationDate := time.Now()
				updateDate := time.Now().Add(time.Hour)
				userId := "user-123"
				externalUserId := "ext-user-456"
				userEnforced := false
				return &openapi.TgvalidatordInternalGroup{
					Id:              &id,
					TenantId:        &tenantId,
					ExternalGroupId: &externalGroupId,
					Name:            &name,
					Email:           &email,
					Description:     &description,
					EnforcedInRules: &enforcedInRules,
					CreationDate:    &creationDate,
					UpdateDate:      &updateDate,
					Users: []openapi.TgvalidatordInternalGroupUser{
						{
							Id:              &userId,
							ExternalUserId:  &externalUserId,
							EnforcedInRules: &userEnforced,
						},
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GroupFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("GroupFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("GroupFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.TenantId != nil && got.TenantID != *tt.dto.TenantId {
				t.Errorf("TenantID = %v, want %v", got.TenantID, *tt.dto.TenantId)
			}
			if tt.dto.ExternalGroupId != nil && got.ExternalGroupID != *tt.dto.ExternalGroupId {
				t.Errorf("ExternalGroupID = %v, want %v", got.ExternalGroupID, *tt.dto.ExternalGroupId)
			}
			if tt.dto.Name != nil && got.Name != *tt.dto.Name {
				t.Errorf("Name = %v, want %v", got.Name, *tt.dto.Name)
			}
			if tt.dto.Email != nil && got.Email != *tt.dto.Email {
				t.Errorf("Email = %v, want %v", got.Email, *tt.dto.Email)
			}
			if tt.dto.Description != nil && got.Description != *tt.dto.Description {
				t.Errorf("Description = %v, want %v", got.Description, *tt.dto.Description)
			}
			if tt.dto.EnforcedInRules != nil && got.EnforcedInRules != *tt.dto.EnforcedInRules {
				t.Errorf("EnforcedInRules = %v, want %v", got.EnforcedInRules, *tt.dto.EnforcedInRules)
			}
			if tt.dto.CreationDate != nil && !got.CreatedAt.Equal(*tt.dto.CreationDate) {
				t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, *tt.dto.CreationDate)
			}
			if tt.dto.UpdateDate != nil && !got.UpdatedAt.Equal(*tt.dto.UpdateDate) {
				t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, *tt.dto.UpdateDate)
			}
			// Verify users are mapped if present
			if tt.dto.Users != nil && len(got.Users) != len(tt.dto.Users) {
				t.Errorf("Users length = %v, want %v", len(got.Users), len(tt.dto.Users))
			}
		})
	}
}

func TestGroupsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordInternalGroup
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordInternalGroup{},
			want: 0,
		},
		{
			name: "converts multiple groups",
			dtos: func() []openapi.TgvalidatordInternalGroup {
				name1 := "Group 1"
				name2 := "Group 2"
				return []openapi.TgvalidatordInternalGroup{
					{Name: &name1},
					{Name: &name2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GroupsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("GroupsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("GroupsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestGroupUserFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordInternalGroupUser
	}{
		{
			name: "nil input returns empty group user",
			dto:  nil,
		},
		{
			name: "empty DTO returns group user with zero values",
			dto:  &openapi.TgvalidatordInternalGroupUser{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordInternalGroupUser {
				id := "user-123"
				externalUserId := "ext-user-456"
				enforcedInRules := true
				return &openapi.TgvalidatordInternalGroupUser{
					Id:              &id,
					ExternalUserId:  &externalUserId,
					EnforcedInRules: &enforcedInRules,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GroupUserFromDTO(tt.dto)
			if tt.dto == nil {
				// nil input returns empty struct
				if got.ID != "" || got.ExternalUserID != "" || got.EnforcedInRules != false {
					t.Errorf("GroupUserFromDTO(nil) should return empty struct, got %+v", got)
				}
				return
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.ExternalUserId != nil && got.ExternalUserID != *tt.dto.ExternalUserId {
				t.Errorf("ExternalUserID = %v, want %v", got.ExternalUserID, *tt.dto.ExternalUserId)
			}
			if tt.dto.EnforcedInRules != nil && got.EnforcedInRules != *tt.dto.EnforcedInRules {
				t.Errorf("EnforcedInRules = %v, want %v", got.EnforcedInRules, *tt.dto.EnforcedInRules)
			}
		})
	}
}

func TestGroupFromDTO_NilDates(t *testing.T) {
	name := "Test Group"
	dto := &openapi.TgvalidatordInternalGroup{
		Name:         &name,
		CreationDate: nil,
		UpdateDate:   nil,
	}

	got := GroupFromDTO(dto)
	if got == nil {
		t.Fatal("GroupFromDTO() returned nil for non-nil input")
	}
	// When dates are nil, they should be the zero time value
	if !got.CreatedAt.IsZero() {
		t.Errorf("CreatedAt should be zero time when nil, got %v", got.CreatedAt)
	}
	if !got.UpdatedAt.IsZero() {
		t.Errorf("UpdatedAt should be zero time when nil, got %v", got.UpdatedAt)
	}
}

func TestGroupFromDTO_NilUsers(t *testing.T) {
	name := "Test Group"
	dto := &openapi.TgvalidatordInternalGroup{
		Name:  &name,
		Users: nil,
	}

	got := GroupFromDTO(dto)
	if got == nil {
		t.Fatal("GroupFromDTO() returned nil for non-nil input")
	}
	if got.Users != nil {
		t.Errorf("Users should be nil when DTO users is nil, got %v", got.Users)
	}
}

func TestGroupFromDTO_EmptyUsers(t *testing.T) {
	name := "Test Group"
	dto := &openapi.TgvalidatordInternalGroup{
		Name:  &name,
		Users: []openapi.TgvalidatordInternalGroupUser{},
	}

	got := GroupFromDTO(dto)
	if got == nil {
		t.Fatal("GroupFromDTO() returned nil for non-nil input")
	}
	if got.Users == nil {
		t.Error("Users should not be nil when DTO users is empty slice")
	}
	if len(got.Users) != 0 {
		t.Errorf("Users length = %v, want 0", len(got.Users))
	}
}

func TestGroupUserFromDTO_EnforcedInRulesField(t *testing.T) {
	tests := []struct {
		name            string
		enforcedInRules *bool
		want            bool
	}{
		{
			name:            "nil enforcedInRules defaults to false",
			enforcedInRules: nil,
			want:            false,
		},
		{
			name:            "true enforcedInRules",
			enforcedInRules: boolPtr(true),
			want:            true,
		},
		{
			name:            "false enforcedInRules",
			enforcedInRules: boolPtr(false),
			want:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &openapi.TgvalidatordInternalGroupUser{
				EnforcedInRules: tt.enforcedInRules,
			}
			got := GroupUserFromDTO(dto)
			if got.EnforcedInRules != tt.want {
				t.Errorf("EnforcedInRules = %v, want %v", got.EnforcedInRules, tt.want)
			}
		})
	}
}
