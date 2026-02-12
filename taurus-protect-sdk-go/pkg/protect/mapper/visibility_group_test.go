package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestVisibilityGroupFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordInternalVisibilityGroup
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns visibility group with zero values",
			dto:  &openapi.TgvalidatordInternalVisibilityGroup{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordInternalVisibilityGroup {
				id := "vg-123"
				tenantId := "tenant-456"
				name := "Engineering Team"
				description := "Engineering department visibility group"
				userCount := "5"
				creationDate := time.Now().Add(-24 * time.Hour)
				updateDate := time.Now()
				userId := "user-789"
				externalUserId := "ext-user-123"
				return &openapi.TgvalidatordInternalVisibilityGroup{
					Id:           &id,
					TenantId:     &tenantId,
					Name:         &name,
					Description:  &description,
					UserCount:    &userCount,
					CreationDate: &creationDate,
					UpdateDate:   &updateDate,
					Users: []openapi.TgvalidatordInternalVisibilityGroupUser{
						{
							Id:             &userId,
							ExternalUserId: &externalUserId,
						},
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := VisibilityGroupFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("VisibilityGroupFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("VisibilityGroupFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.TenantId != nil && got.TenantID != *tt.dto.TenantId {
				t.Errorf("TenantID = %v, want %v", got.TenantID, *tt.dto.TenantId)
			}
			if tt.dto.Name != nil && got.Name != *tt.dto.Name {
				t.Errorf("Name = %v, want %v", got.Name, *tt.dto.Name)
			}
			if tt.dto.Description != nil && got.Description != *tt.dto.Description {
				t.Errorf("Description = %v, want %v", got.Description, *tt.dto.Description)
			}
			if tt.dto.CreationDate != nil && !got.CreatedAt.Equal(*tt.dto.CreationDate) {
				t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, *tt.dto.CreationDate)
			}
			if tt.dto.UpdateDate != nil && !got.UpdatedAt.Equal(*tt.dto.UpdateDate) {
				t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, *tt.dto.UpdateDate)
			}
			// Verify users are mapped if present
			if tt.dto.Users != nil {
				if got.Users == nil || len(got.Users) != len(tt.dto.Users) {
					t.Errorf("Users length = %d, want %d", len(got.Users), len(tt.dto.Users))
				}
			}
		})
	}
}

func TestVisibilityGroupFromDTO_UserCount(t *testing.T) {
	tests := []struct {
		name      string
		userCount *string
		want      int64
	}{
		{
			name:      "nil user count defaults to 0",
			userCount: nil,
			want:      0,
		},
		{
			name:      "valid user count is parsed",
			userCount: stringPtr("10"),
			want:      10,
		},
		{
			name:      "invalid user count defaults to 0",
			userCount: stringPtr("invalid"),
			want:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &openapi.TgvalidatordInternalVisibilityGroup{
				UserCount: tt.userCount,
			}
			got := VisibilityGroupFromDTO(dto)
			if got.UserCount != tt.want {
				t.Errorf("UserCount = %v, want %v", got.UserCount, tt.want)
			}
		})
	}
}

func TestVisibilityGroupsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordInternalVisibilityGroup
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordInternalVisibilityGroup{},
			want: 0,
		},
		{
			name: "converts multiple visibility groups",
			dtos: func() []openapi.TgvalidatordInternalVisibilityGroup {
				name1 := "Group A"
				name2 := "Group B"
				return []openapi.TgvalidatordInternalVisibilityGroup{
					{Name: &name1},
					{Name: &name2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := VisibilityGroupsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("VisibilityGroupsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("VisibilityGroupsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestVisibilityGroupUserFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordInternalVisibilityGroupUser
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns user with zero values",
			dto:  &openapi.TgvalidatordInternalVisibilityGroupUser{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordInternalVisibilityGroupUser {
				id := "user-123"
				externalUserId := "ext-user-456"
				return &openapi.TgvalidatordInternalVisibilityGroupUser{
					Id:             &id,
					ExternalUserId: &externalUserId,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := VisibilityGroupUserFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("VisibilityGroupUserFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("VisibilityGroupUserFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.ExternalUserId != nil && got.ExternalUserID != *tt.dto.ExternalUserId {
				t.Errorf("ExternalUserID = %v, want %v", got.ExternalUserID, *tt.dto.ExternalUserId)
			}
		})
	}
}

func TestVisibilityGroupFromDTO_NilTimestamps(t *testing.T) {
	dto := &openapi.TgvalidatordInternalVisibilityGroup{
		CreationDate: nil,
		UpdateDate:   nil,
	}

	got := VisibilityGroupFromDTO(dto)
	if !got.CreatedAt.IsZero() {
		t.Errorf("CreatedAt should be zero time when nil, got %v", got.CreatedAt)
	}
	if !got.UpdatedAt.IsZero() {
		t.Errorf("UpdatedAt should be zero time when nil, got %v", got.UpdatedAt)
	}
}

func TestVisibilityGroupFromDTO_NilUsers(t *testing.T) {
	dto := &openapi.TgvalidatordInternalVisibilityGroup{
		Users: nil,
	}

	got := VisibilityGroupFromDTO(dto)
	if got.Users != nil {
		t.Errorf("Users should be nil when DTO users is nil, got %v", got.Users)
	}
}

func TestVisibilityGroupFromDTO_EmptyUsers(t *testing.T) {
	dto := &openapi.TgvalidatordInternalVisibilityGroup{
		Users: []openapi.TgvalidatordInternalVisibilityGroupUser{},
	}

	got := VisibilityGroupFromDTO(dto)
	if got.Users == nil || len(got.Users) != 0 {
		t.Errorf("Users should be empty slice, got %v", got.Users)
	}
}

func TestVisibilityGroupFromDTO_MultipleUsers(t *testing.T) {
	user1Id := "user-1"
	user1ExtId := "ext-1"
	user2Id := "user-2"
	user2ExtId := "ext-2"

	dto := &openapi.TgvalidatordInternalVisibilityGroup{
		Users: []openapi.TgvalidatordInternalVisibilityGroupUser{
			{Id: &user1Id, ExternalUserId: &user1ExtId},
			{Id: &user2Id, ExternalUserId: &user2ExtId},
		},
	}

	got := VisibilityGroupFromDTO(dto)
	if len(got.Users) != 2 {
		t.Fatalf("Users length = %d, want 2", len(got.Users))
	}
	if got.Users[0].ID != user1Id {
		t.Errorf("Users[0].ID = %v, want %v", got.Users[0].ID, user1Id)
	}
	if got.Users[0].ExternalUserID != user1ExtId {
		t.Errorf("Users[0].ExternalUserID = %v, want %v", got.Users[0].ExternalUserID, user1ExtId)
	}
	if got.Users[1].ID != user2Id {
		t.Errorf("Users[1].ID = %v, want %v", got.Users[1].ID, user2Id)
	}
	if got.Users[1].ExternalUserID != user2ExtId {
		t.Errorf("Users[1].ExternalUserID = %v, want %v", got.Users[1].ExternalUserID, user2ExtId)
	}
}
