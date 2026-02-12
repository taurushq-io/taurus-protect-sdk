package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestAuditTrailFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordAuditTrail
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns audit trail with zero values",
			dto:  &openapi.TgvalidatordAuditTrail{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordAuditTrail {
				id := "audit-123"
				entity := "Wallet"
				action := "Create"
				subAction := "API"
				details := `{"walletId":"wallet-456"}`
				creationDate := time.Now()
				userId := "user-789"
				externalUserId := "ext-user-123"
				email := "test@example.com"
				deleted := false
				return &openapi.TgvalidatordAuditTrail{
					Id:           &id,
					Entity:       &entity,
					Action:       &action,
					SubAction:    &subAction,
					Details:      &details,
					CreationDate: &creationDate,
					User: &openapi.TgvalidatordUserInfo{
						Id:             &userId,
						ExternalUserId: &externalUserId,
						Email:          &email,
						Deleted:        &deleted,
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AuditTrailFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("AuditTrailFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("AuditTrailFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Entity != nil && got.Entity != *tt.dto.Entity {
				t.Errorf("Entity = %v, want %v", got.Entity, *tt.dto.Entity)
			}
			if tt.dto.Action != nil && got.Action != *tt.dto.Action {
				t.Errorf("Action = %v, want %v", got.Action, *tt.dto.Action)
			}
			if tt.dto.SubAction != nil && got.SubAction != *tt.dto.SubAction {
				t.Errorf("SubAction = %v, want %v", got.SubAction, *tt.dto.SubAction)
			}
			if tt.dto.Details != nil && got.Details != *tt.dto.Details {
				t.Errorf("Details = %v, want %v", got.Details, *tt.dto.Details)
			}
			if tt.dto.CreationDate != nil && !got.CreationDate.Equal(*tt.dto.CreationDate) {
				t.Errorf("CreationDate = %v, want %v", got.CreationDate, *tt.dto.CreationDate)
			}
			// Verify user is mapped if present
			if tt.dto.User != nil {
				if got.User == nil {
					t.Error("User should not be nil when DTO has user")
				}
			}
		})
	}
}

func TestAuditTrailsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordAuditTrail
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordAuditTrail{},
			want: 0,
		},
		{
			name: "converts multiple audit trails",
			dtos: func() []openapi.TgvalidatordAuditTrail {
				entity1 := "Wallet"
				entity2 := "Request"
				return []openapi.TgvalidatordAuditTrail{
					{Entity: &entity1},
					{Entity: &entity2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AuditTrailsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("AuditTrailsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("AuditTrailsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestUserInfoFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordUserInfo
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns user info with zero values",
			dto:  &openapi.TgvalidatordUserInfo{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordUserInfo {
				id := "user-123"
				externalUserId := "ext-user-456"
				email := "user@example.com"
				deleted := true
				return &openapi.TgvalidatordUserInfo{
					Id:             &id,
					ExternalUserId: &externalUserId,
					Email:          &email,
					Deleted:        &deleted,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UserInfoFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("UserInfoFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("UserInfoFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.ExternalUserId != nil && got.ExternalUserID != *tt.dto.ExternalUserId {
				t.Errorf("ExternalUserID = %v, want %v", got.ExternalUserID, *tt.dto.ExternalUserId)
			}
			if tt.dto.Email != nil && got.Email != *tt.dto.Email {
				t.Errorf("Email = %v, want %v", got.Email, *tt.dto.Email)
			}
			if tt.dto.Deleted != nil && got.Deleted != *tt.dto.Deleted {
				t.Errorf("Deleted = %v, want %v", got.Deleted, *tt.dto.Deleted)
			}
		})
	}
}

func TestUserInfoFromDTO_DeletedField(t *testing.T) {
	tests := []struct {
		name        string
		deleted     *bool
		wantDeleted bool
	}{
		{
			name:        "nil deleted defaults to false",
			deleted:     nil,
			wantDeleted: false,
		},
		{
			name:        "true deleted",
			deleted:     boolPtr(true),
			wantDeleted: true,
		},
		{
			name:        "false deleted",
			deleted:     boolPtr(false),
			wantDeleted: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &openapi.TgvalidatordUserInfo{
				Deleted: tt.deleted,
			}
			got := UserInfoFromDTO(dto)
			if got.Deleted != tt.wantDeleted {
				t.Errorf("Deleted = %v, want %v", got.Deleted, tt.wantDeleted)
			}
		})
	}
}

func TestAuditTrailFromDTO_NilUser(t *testing.T) {
	entity := "Wallet"
	dto := &openapi.TgvalidatordAuditTrail{
		Entity: &entity,
		User:   nil,
	}

	got := AuditTrailFromDTO(dto)
	if got == nil {
		t.Fatal("AuditTrailFromDTO() returned nil for non-nil input")
	}
	if got.User != nil {
		t.Errorf("User should be nil when DTO user is nil, got %v", got.User)
	}
	if got.Entity != "Wallet" {
		t.Errorf("Entity = %v, want Wallet", got.Entity)
	}
}

func TestAuditTrailFromDTO_NilCreationDate(t *testing.T) {
	entity := "Wallet"
	dto := &openapi.TgvalidatordAuditTrail{
		Entity:       &entity,
		CreationDate: nil,
	}

	got := AuditTrailFromDTO(dto)
	if got == nil {
		t.Fatal("AuditTrailFromDTO() returned nil for non-nil input")
	}
	// When creation date is nil, it should be the zero time value
	if !got.CreationDate.IsZero() {
		t.Errorf("CreationDate should be zero time when nil, got %v", got.CreationDate)
	}
}
