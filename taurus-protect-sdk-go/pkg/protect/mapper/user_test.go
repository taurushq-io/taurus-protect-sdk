package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestUserFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordInternalUser
		want func(t *testing.T, got interface{})
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
			want: func(t *testing.T, got interface{}) {
				if got != nil {
					t.Errorf("expected nil, got %v", got)
				}
			},
		},
		{
			name: "empty DTO returns user with zero values",
			dto:  &openapi.TgvalidatordInternalUser{},
			want: func(t *testing.T, got interface{}) {
				if got == nil {
					t.Error("expected non-nil user")
				}
			},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordInternalUser {
				id := "user-123"
				tenantId := "tenant-456"
				externalUserId := "ext-user-789"
				username := "johndoe"
				email := "john@example.com"
				firstName := "John"
				lastName := "Doe"
				status := "ACTIVE"
				publicKey := "pk-abc123"
				passwordChanged := true
				totpEnabled := true
				enforcedInRules := true
				publicKeyEnforcedInRules := false
				now := time.Now()
				return &openapi.TgvalidatordInternalUser{
					Id:                       &id,
					TenantId:                 &tenantId,
					ExternalUserId:           &externalUserId,
					Username:                 &username,
					Email:                    &email,
					FirstName:                &firstName,
					LastName:                 &lastName,
					Status:                   &status,
					PublicKey:                &publicKey,
					PasswordChanged:          &passwordChanged,
					TotpEnabled:              &totpEnabled,
					EnforcedInRules:          &enforcedInRules,
					PublicKeyEnforcedInRules: &publicKeyEnforcedInRules,
					Roles:                    []string{"admin", "operator"},
					CreationDate:             &now,
					UpdateDate:               &now,
					LastLogin:                &now,
				}
			}(),
			want: func(t *testing.T, got interface{}) {
				// Type assertion handled in test execution
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UserFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("UserFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("UserFromDTO() returned nil for non-nil input")
			}
			// Verify specific fields for complete DTO
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.TenantId != nil && got.TenantID != *tt.dto.TenantId {
				t.Errorf("TenantID = %v, want %v", got.TenantID, *tt.dto.TenantId)
			}
			if tt.dto.Email != nil && got.Email != *tt.dto.Email {
				t.Errorf("Email = %v, want %v", got.Email, *tt.dto.Email)
			}
			if tt.dto.Username != nil && got.Username != *tt.dto.Username {
				t.Errorf("Username = %v, want %v", got.Username, *tt.dto.Username)
			}
			if tt.dto.Roles != nil && len(got.Roles) != len(tt.dto.Roles) {
				t.Errorf("Roles length = %v, want %v", len(got.Roles), len(tt.dto.Roles))
			}
		})
	}
}

func TestUserFromDTO_WithGroups(t *testing.T) {
	groupID := "group-123"
	externalGroupID := "ext-group-456"
	enforcedInRules := true

	dto := &openapi.TgvalidatordInternalUser{
		Groups: []openapi.TgvalidatordInternalUserGroup{
			{
				Id:              &groupID,
				ExternalGroupId: &externalGroupID,
				EnforcedInRules: &enforcedInRules,
			},
		},
	}

	got := UserFromDTO(dto)
	if got == nil {
		t.Fatal("UserFromDTO() returned nil")
	}
	if len(got.Groups) != 1 {
		t.Fatalf("Groups length = %v, want 1", len(got.Groups))
	}
	if got.Groups[0].ID != groupID {
		t.Errorf("Group.ID = %v, want %v", got.Groups[0].ID, groupID)
	}
	if got.Groups[0].ExternalGroupID != externalGroupID {
		t.Errorf("Group.ExternalGroupID = %v, want %v", got.Groups[0].ExternalGroupID, externalGroupID)
	}
	if got.Groups[0].EnforcedInRules != enforcedInRules {
		t.Errorf("Group.EnforcedInRules = %v, want %v", got.Groups[0].EnforcedInRules, enforcedInRules)
	}
}

func TestUserFromDTO_WithAttributes(t *testing.T) {
	attrID := "attr-123"
	key := "department"
	value := "engineering"
	contentType := "text/plain"
	isFile := false
	now := time.Now()

	dto := &openapi.TgvalidatordInternalUser{
		Attributes: []openapi.InternalUserAttribute{
			{
				Id:           &attrID,
				Key:          &key,
				Value:        &value,
				ContentType:  &contentType,
				IsFile:       &isFile,
				CreationDate: &now,
				UpdateDate:   &now,
			},
		},
	}

	got := UserFromDTO(dto)
	if got == nil {
		t.Fatal("UserFromDTO() returned nil")
	}
	if len(got.Attributes) != 1 {
		t.Fatalf("Attributes length = %v, want 1", len(got.Attributes))
	}
	if got.Attributes[0].ID != attrID {
		t.Errorf("Attribute.ID = %v, want %v", got.Attributes[0].ID, attrID)
	}
	if got.Attributes[0].Key != key {
		t.Errorf("Attribute.Key = %v, want %v", got.Attributes[0].Key, key)
	}
	if got.Attributes[0].Value != value {
		t.Errorf("Attribute.Value = %v, want %v", got.Attributes[0].Value, value)
	}
}

func TestUsersFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordInternalUser
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordInternalUser{},
			want: 0,
		},
		{
			name: "converts multiple users",
			dtos: func() []openapi.TgvalidatordInternalUser {
				id1 := "user-1"
				id2 := "user-2"
				return []openapi.TgvalidatordInternalUser{
					{Id: &id1},
					{Id: &id2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UsersFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("UsersFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("UsersFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestUserGroupFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordInternalUserGroup
	}{
		{
			name: "nil input returns zero value",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordInternalUserGroup {
				id := "group-123"
				externalGroupId := "ext-group-456"
				enforcedInRules := true
				return &openapi.TgvalidatordInternalUserGroup{
					Id:              &id,
					ExternalGroupId: &externalGroupId,
					EnforcedInRules: &enforcedInRules,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UserGroupFromDTO(tt.dto)
			if tt.dto == nil {
				if got.ID != "" || got.ExternalGroupID != "" || got.EnforcedInRules != false {
					t.Errorf("UserGroupFromDTO(nil) should return zero value")
				}
				return
			}
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.ExternalGroupId != nil && got.ExternalGroupID != *tt.dto.ExternalGroupId {
				t.Errorf("ExternalGroupID = %v, want %v", got.ExternalGroupID, *tt.dto.ExternalGroupId)
			}
			if tt.dto.EnforcedInRules != nil && got.EnforcedInRules != *tt.dto.EnforcedInRules {
				t.Errorf("EnforcedInRules = %v, want %v", got.EnforcedInRules, *tt.dto.EnforcedInRules)
			}
		})
	}
}

func TestUserGroupsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordInternalUserGroup
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1,
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordInternalUserGroup{},
			want: 0,
		},
		{
			name: "converts multiple groups",
			dtos: func() []openapi.TgvalidatordInternalUserGroup {
				id1 := "group-1"
				id2 := "group-2"
				return []openapi.TgvalidatordInternalUserGroup{
					{Id: &id1},
					{Id: &id2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UserGroupsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("UserGroupsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("UserGroupsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestUserAttributeFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.InternalUserAttribute
	}{
		{
			name: "nil input returns zero value",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.InternalUserAttribute {
				id := "attr-123"
				key := "department"
				value := "engineering"
				contentType := "text/plain"
				owner := "user-456"
				attrType := "custom"
				subType := "user-defined"
				isFile := false
				now := time.Now()
				return &openapi.InternalUserAttribute{
					Id:           &id,
					Key:          &key,
					Value:        &value,
					ContentType:  &contentType,
					Owner:        &owner,
					Type:         &attrType,
					SubType:      &subType,
					IsFile:       &isFile,
					CreationDate: &now,
					UpdateDate:   &now,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UserAttributeFromDTO(tt.dto)
			if tt.dto == nil {
				if got.ID != "" || got.Key != "" || got.Value != "" {
					t.Errorf("UserAttributeFromDTO(nil) should return zero value")
				}
				return
			}
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Key != nil && got.Key != *tt.dto.Key {
				t.Errorf("Key = %v, want %v", got.Key, *tt.dto.Key)
			}
			if tt.dto.Value != nil && got.Value != *tt.dto.Value {
				t.Errorf("Value = %v, want %v", got.Value, *tt.dto.Value)
			}
			if tt.dto.ContentType != nil && got.ContentType != *tt.dto.ContentType {
				t.Errorf("ContentType = %v, want %v", got.ContentType, *tt.dto.ContentType)
			}
			if tt.dto.IsFile != nil && got.IsFile != *tt.dto.IsFile {
				t.Errorf("IsFile = %v, want %v", got.IsFile, *tt.dto.IsFile)
			}
		})
	}
}

func TestUserAttributesFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.InternalUserAttribute
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1,
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.InternalUserAttribute{},
			want: 0,
		},
		{
			name: "converts multiple attributes",
			dtos: func() []openapi.InternalUserAttribute {
				id1 := "attr-1"
				id2 := "attr-2"
				return []openapi.InternalUserAttribute{
					{Id: &id1},
					{Id: &id2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UserAttributesFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("UserAttributesFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("UserAttributesFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}
