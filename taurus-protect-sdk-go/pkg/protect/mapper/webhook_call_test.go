package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestWebhookCallFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordWebhookCall
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns webhook call with zero values",
			dto:  &openapi.TgvalidatordWebhookCall{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordWebhookCall {
				id := "call-123"
				eventId := "event-456"
				webhookId := "webhook-789"
				payload := `{"type":"transaction","data":{}}`
				status := "SUCCESS"
				statusMessage := "Delivered successfully"
				attempts := "3"
				createdAt := time.Now().Add(-time.Hour)
				updatedAt := time.Now()
				return &openapi.TgvalidatordWebhookCall{
					Id:            &id,
					EventId:       &eventId,
					WebhookId:     &webhookId,
					Payload:       &payload,
					Status:        &status,
					StatusMessage: &statusMessage,
					Attempts:      &attempts,
					CreatedAt:     &createdAt,
					UpdatedAt:     &updatedAt,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WebhookCallFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("WebhookCallFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("WebhookCallFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.EventId != nil && got.EventID != *tt.dto.EventId {
				t.Errorf("EventID = %v, want %v", got.EventID, *tt.dto.EventId)
			}
			if tt.dto.WebhookId != nil && got.WebhookID != *tt.dto.WebhookId {
				t.Errorf("WebhookID = %v, want %v", got.WebhookID, *tt.dto.WebhookId)
			}
			if tt.dto.Payload != nil && got.Payload != *tt.dto.Payload {
				t.Errorf("Payload = %v, want %v", got.Payload, *tt.dto.Payload)
			}
			if tt.dto.Status != nil && got.Status != *tt.dto.Status {
				t.Errorf("Status = %v, want %v", got.Status, *tt.dto.Status)
			}
			if tt.dto.StatusMessage != nil && got.StatusMessage != *tt.dto.StatusMessage {
				t.Errorf("StatusMessage = %v, want %v", got.StatusMessage, *tt.dto.StatusMessage)
			}
			if tt.dto.CreatedAt != nil && !got.CreatedAt.Equal(*tt.dto.CreatedAt) {
				t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, *tt.dto.CreatedAt)
			}
			if tt.dto.UpdatedAt != nil && !got.UpdatedAt.Equal(*tt.dto.UpdatedAt) {
				t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, *tt.dto.UpdatedAt)
			}
		})
	}
}

func TestWebhookCallsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordWebhookCall
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordWebhookCall{},
			want: 0,
		},
		{
			name: "converts multiple webhook calls",
			dtos: func() []openapi.TgvalidatordWebhookCall {
				id1 := "call-1"
				id2 := "call-2"
				return []openapi.TgvalidatordWebhookCall{
					{Id: &id1},
					{Id: &id2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WebhookCallsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("WebhookCallsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("WebhookCallsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestWebhookCallFromDTO_AttemptsField(t *testing.T) {
	tests := []struct {
		name         string
		attempts     *string
		wantAttempts int64
	}{
		{
			name:         "nil attempts defaults to 0",
			attempts:     nil,
			wantAttempts: 0,
		},
		{
			name:         "valid attempts string",
			attempts:     stringPtr("5"),
			wantAttempts: 5,
		},
		{
			name:         "zero attempts",
			attempts:     stringPtr("0"),
			wantAttempts: 0,
		},
		{
			name:         "invalid attempts string defaults to 0",
			attempts:     stringPtr("invalid"),
			wantAttempts: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &openapi.TgvalidatordWebhookCall{
				Attempts: tt.attempts,
			}
			got := WebhookCallFromDTO(dto)
			if got.Attempts != tt.wantAttempts {
				t.Errorf("Attempts = %v, want %v", got.Attempts, tt.wantAttempts)
			}
		})
	}
}

func TestWebhookCallFromDTO_NilTimestamps(t *testing.T) {
	id := "call-123"
	dto := &openapi.TgvalidatordWebhookCall{
		Id:        &id,
		CreatedAt: nil,
		UpdatedAt: nil,
	}

	got := WebhookCallFromDTO(dto)
	if got == nil {
		t.Fatal("WebhookCallFromDTO() returned nil for non-nil input")
	}
	// When timestamps are nil, they should be the zero time value
	if !got.CreatedAt.IsZero() {
		t.Errorf("CreatedAt should be zero time when nil, got %v", got.CreatedAt)
	}
	if !got.UpdatedAt.IsZero() {
		t.Errorf("UpdatedAt should be zero time when nil, got %v", got.UpdatedAt)
	}
}

func TestWebhookCallFromDTO_StatusValues(t *testing.T) {
	validStatuses := []string{"SUCCESS", "FAILED", "PENDING", "RETRY"}

	for _, status := range validStatuses {
		t.Run(status, func(t *testing.T) {
			dto := &openapi.TgvalidatordWebhookCall{
				Status: stringPtr(status),
			}
			got := WebhookCallFromDTO(dto)
			if got.Status != status {
				t.Errorf("Status = %v, want %v", got.Status, status)
			}
		})
	}
}
