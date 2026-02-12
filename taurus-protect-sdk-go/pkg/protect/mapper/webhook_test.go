package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestWebhookFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordWebhook
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns webhook with zero values",
			dto:  &openapi.TgvalidatordWebhook{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordWebhook {
				id := "webhook-123"
				webhookType := "request"
				url := "https://example.com/webhook"
				status := "enabled"
				createdAt := time.Now().Add(-24 * time.Hour)
				updatedAt := time.Now()
				timeoutUntil := time.Now().Add(time.Hour)
				return &openapi.TgvalidatordWebhook{
					Id:           &id,
					Type:         &webhookType,
					Url:          &url,
					Status:       &status,
					CreatedAt:    &createdAt,
					UpdatedAt:    &updatedAt,
					TimeoutUntil: &timeoutUntil,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WebhookFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("WebhookFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("WebhookFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Type != nil && got.Type != *tt.dto.Type {
				t.Errorf("Type = %v, want %v", got.Type, *tt.dto.Type)
			}
			if tt.dto.Url != nil && got.URL != *tt.dto.Url {
				t.Errorf("URL = %v, want %v", got.URL, *tt.dto.Url)
			}
			if tt.dto.Status != nil && got.Status != *tt.dto.Status {
				t.Errorf("Status = %v, want %v", got.Status, *tt.dto.Status)
			}
			if tt.dto.CreatedAt != nil && !got.CreatedAt.Equal(*tt.dto.CreatedAt) {
				t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, *tt.dto.CreatedAt)
			}
			if tt.dto.UpdatedAt != nil && !got.UpdatedAt.Equal(*tt.dto.UpdatedAt) {
				t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, *tt.dto.UpdatedAt)
			}
			if tt.dto.TimeoutUntil != nil && !got.TimeoutUntil.Equal(*tt.dto.TimeoutUntil) {
				t.Errorf("TimeoutUntil = %v, want %v", got.TimeoutUntil, *tt.dto.TimeoutUntil)
			}
		})
	}
}

func TestWebhooksFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordWebhook
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordWebhook{},
			want: 0,
		},
		{
			name: "converts multiple webhooks",
			dtos: func() []openapi.TgvalidatordWebhook {
				id1 := "webhook-1"
				id2 := "webhook-2"
				return []openapi.TgvalidatordWebhook{
					{Id: &id1},
					{Id: &id2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WebhooksFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("WebhooksFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("WebhooksFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestWebhookFromDTO_NilTimestamps(t *testing.T) {
	id := "webhook-123"
	dto := &openapi.TgvalidatordWebhook{
		Id:           &id,
		CreatedAt:    nil,
		UpdatedAt:    nil,
		TimeoutUntil: nil,
	}

	got := WebhookFromDTO(dto)
	if got == nil {
		t.Fatal("WebhookFromDTO() returned nil for non-nil input")
	}
	if !got.CreatedAt.IsZero() {
		t.Errorf("CreatedAt should be zero time when nil, got %v", got.CreatedAt)
	}
	if !got.UpdatedAt.IsZero() {
		t.Errorf("UpdatedAt should be zero time when nil, got %v", got.UpdatedAt)
	}
	if !got.TimeoutUntil.IsZero() {
		t.Errorf("TimeoutUntil should be zero time when nil, got %v", got.TimeoutUntil)
	}
}

func TestWebhookFromDTO_PartialFields(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordWebhook
	}{
		{
			name: "only id set",
			dto: func() *openapi.TgvalidatordWebhook {
				id := "webhook-only-id"
				return &openapi.TgvalidatordWebhook{Id: &id}
			}(),
		},
		{
			name: "only type and url set",
			dto: func() *openapi.TgvalidatordWebhook {
				webhookType := "transaction"
				url := "https://example.com/tx-webhook"
				return &openapi.TgvalidatordWebhook{
					Type: &webhookType,
					Url:  &url,
				}
			}(),
		},
		{
			name: "only status set",
			dto: func() *openapi.TgvalidatordWebhook {
				status := "disabled"
				return &openapi.TgvalidatordWebhook{Status: &status}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WebhookFromDTO(tt.dto)
			if got == nil {
				t.Fatal("WebhookFromDTO() returned nil for non-nil input")
			}
			// Verify that non-nil fields are set correctly
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Type != nil && got.Type != *tt.dto.Type {
				t.Errorf("Type = %v, want %v", got.Type, *tt.dto.Type)
			}
			if tt.dto.Url != nil && got.URL != *tt.dto.Url {
				t.Errorf("URL = %v, want %v", got.URL, *tt.dto.Url)
			}
			if tt.dto.Status != nil && got.Status != *tt.dto.Status {
				t.Errorf("Status = %v, want %v", got.Status, *tt.dto.Status)
			}
			// Verify nil fields result in empty strings
			if tt.dto.Id == nil && got.ID != "" {
				t.Errorf("ID should be empty when nil, got %v", got.ID)
			}
			if tt.dto.Type == nil && got.Type != "" {
				t.Errorf("Type should be empty when nil, got %v", got.Type)
			}
			if tt.dto.Url == nil && got.URL != "" {
				t.Errorf("URL should be empty when nil, got %v", got.URL)
			}
			if tt.dto.Status == nil && got.Status != "" {
				t.Errorf("Status should be empty when nil, got %v", got.Status)
			}
		})
	}
}

func TestWebhooksFromDTO_PreservesOrder(t *testing.T) {
	id1 := "first"
	id2 := "second"
	id3 := "third"
	dtos := []openapi.TgvalidatordWebhook{
		{Id: &id1},
		{Id: &id2},
		{Id: &id3},
	}

	got := WebhooksFromDTO(dtos)
	if len(got) != 3 {
		t.Fatalf("WebhooksFromDTO() length = %v, want 3", len(got))
	}

	expectedOrder := []string{"first", "second", "third"}
	for i, expected := range expectedOrder {
		if got[i].ID != expected {
			t.Errorf("WebhooksFromDTO()[%d].ID = %v, want %v", i, got[i].ID, expected)
		}
	}
}
