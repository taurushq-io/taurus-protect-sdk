package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestRequestFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordRequest
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns request with zero values",
			dto:  &openapi.TgvalidatordRequest{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordRequest {
				id := "req-123"
				tenantId := "tenant-456"
				currency := "ETH"
				status := "CONFIRMED"
				requestType := "transfer"
				rule := "standard"
				requestBundleId := "bundle-789"
				externalRequestId := "ext-req-001"
				needsApprovalFrom := []string{"group-1", "group-2"}
				now := time.Now()
				return &openapi.TgvalidatordRequest{
					Id:                &id,
					TenantId:          &tenantId,
					Currency:          &currency,
					Status:            &status,
					Type:              &requestType,
					Rule:              &rule,
					RequestBundleId:   &requestBundleId,
					ExternalRequestId: &externalRequestId,
					NeedsApprovalFrom: needsApprovalFrom,
					CreationDate:      &now,
					UpdateDate:        &now,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RequestFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("RequestFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("RequestFromDTO() returned nil for non-nil input")
			}
			// Verify key fields
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Status != nil && got.Status != *tt.dto.Status {
				t.Errorf("Status = %v, want %v", got.Status, *tt.dto.Status)
			}
		})
	}
}

func TestRequestFromDTO_WithMetadata(t *testing.T) {
	hash := "abc123"
	payloadString := `{"key":"value"}`
	payload := map[string]interface{}{"key": "value"}
	dto := &openapi.TgvalidatordRequest{
		Metadata: &openapi.TgvalidatordMetadata{
			Hash:            &hash,
			PayloadAsString: &payloadString,
			Payload:         payload,
		},
	}

	got := RequestFromDTO(dto)
	if got.Metadata == nil {
		t.Fatal("Metadata should not be nil")
	}
	if got.Metadata.Hash != hash {
		t.Errorf("Metadata.Hash = %v, want %v", got.Metadata.Hash, hash)
	}
	if got.Metadata.PayloadAsString != payloadString {
		t.Errorf("Metadata.PayloadAsString = %v, want %v", got.Metadata.PayloadAsString, payloadString)
	}
	// SECURITY: Payload field intentionally removed - use PayloadAsString for data extraction
}

func TestRequestFromDTO_WithAttributes(t *testing.T) {
	attrId := "attr-1"
	attrKey := "priority"
	attrValue := "high"
	dto := &openapi.TgvalidatordRequest{
		Attributes: []openapi.TgvalidatordRequestAttribute{
			{Id: &attrId, Key: &attrKey, Value: &attrValue},
		},
	}

	got := RequestFromDTO(dto)
	if len(got.Attributes) != 1 {
		t.Fatalf("Attributes length = %v, want 1", len(got.Attributes))
	}
	if got.Attributes[0].ID != attrId {
		t.Errorf("Attribute ID = %v, want %v", got.Attributes[0].ID, attrId)
	}
}

func TestRequestFromDTO_NeedsApprovalFromCopy(t *testing.T) {
	// Verify that the slice is copied, not referenced
	groups := []string{"group-1", "group-2"}
	dto := &openapi.TgvalidatordRequest{
		NeedsApprovalFrom: groups,
	}

	got := RequestFromDTO(dto)

	// Modify the original slice
	groups[0] = "modified"

	// The converted request should still have the original value
	if got.NeedsApprovalFrom[0] != "group-1" {
		t.Errorf("NeedsApprovalFrom was not properly copied, got %v", got.NeedsApprovalFrom[0])
	}
}

func TestRequestsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordRequest
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1,
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordRequest{},
			want: 0,
		},
		{
			name: "converts multiple requests",
			dtos: func() []openapi.TgvalidatordRequest {
				id1 := "req-1"
				id2 := "req-2"
				return []openapi.TgvalidatordRequest{
					{Id: &id1},
					{Id: &id2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RequestsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("RequestsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("RequestsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestMetadataFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordMetadata
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns metadata with zero values",
			dto:  &openapi.TgvalidatordMetadata{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordMetadata {
				hash := "hash123"
				payloadAsString := `{"data":"test"}`
				payload := map[string]interface{}{"data": "test"}
				return &openapi.TgvalidatordMetadata{
					Hash:            &hash,
					PayloadAsString: &payloadAsString,
					Payload:         payload,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MetadataFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("MetadataFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("MetadataFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Hash != nil && got.Hash != *tt.dto.Hash {
				t.Errorf("Hash = %v, want %v", got.Hash, *tt.dto.Hash)
			}
		})
	}
}

func TestRequestAttributeFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordRequestAttribute
	}{
		{
			name: "nil input returns zero value",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordRequestAttribute {
				id := "attr-123"
				key := "note"
				value := "important transaction"
				return &openapi.TgvalidatordRequestAttribute{
					Id:    &id,
					Key:   &key,
					Value: &value,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RequestAttributeFromDTO(tt.dto)
			if tt.dto == nil {
				if got.ID != "" || got.Key != "" || got.Value != "" {
					t.Errorf("RequestAttributeFromDTO(nil) should return zero value")
				}
				return
			}
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
		})
	}
}

func TestRequestFromDTO_TimestampConversion(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	dto := &openapi.TgvalidatordRequest{
		CreationDate: &now,
		UpdateDate:   &now,
	}

	got := RequestFromDTO(dto)
	if !got.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, now)
	}
	if !got.UpdatedAt.Equal(now) {
		t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, now)
	}
}
