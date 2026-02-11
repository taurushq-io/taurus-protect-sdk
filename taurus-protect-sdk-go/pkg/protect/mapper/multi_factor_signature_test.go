package mapper

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestMultiFactorSignatureInfoFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordGetMultiFactorSignatureEntitiesInfoReply
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns info with zero values",
			dto:  &openapi.TgvalidatordGetMultiFactorSignatureEntitiesInfoReply{},
		},
		{
			name: "complete DTO maps all fields",
			dto: &openapi.TgvalidatordGetMultiFactorSignatureEntitiesInfoReply{
				Id:            "mfs-123",
				PayloadToSign: []string{"payload-1", "payload-2"},
				EntityType:    "WHITELISTED_ADDRESS",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MultiFactorSignatureInfoFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("MultiFactorSignatureInfoFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("MultiFactorSignatureInfoFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Id != "" && got.ID != tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, tt.dto.Id)
			}
			if len(tt.dto.PayloadToSign) > 0 && len(got.PayloadToSign) != len(tt.dto.PayloadToSign) {
				t.Errorf("PayloadToSign length = %v, want %v", len(got.PayloadToSign), len(tt.dto.PayloadToSign))
			}
			if tt.dto.EntityType != "" && string(got.EntityType) != string(tt.dto.EntityType) {
				t.Errorf("EntityType = %v, want %v", got.EntityType, tt.dto.EntityType)
			}
		})
	}
}

func TestMultiFactorSignatureResultFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordCreateMultiFactorSignaturesReply
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns result with empty ID",
			dto:  &openapi.TgvalidatordCreateMultiFactorSignaturesReply{},
		},
		{
			name: "complete DTO maps ID",
			dto: &openapi.TgvalidatordCreateMultiFactorSignaturesReply{
				Id: strPtr("result-456"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MultiFactorSignatureResultFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("MultiFactorSignatureResultFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("MultiFactorSignatureResultFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
		})
	}
}

func TestMultiFactorSignatureApprovalResultFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordApproveMultiFactorSignatureReply
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns result with nil signature count",
			dto:  &openapi.TgvalidatordApproveMultiFactorSignatureReply{},
		},
		{
			name: "complete DTO maps signature count",
			dto: &openapi.TgvalidatordApproveMultiFactorSignatureReply{
				Signatures: "3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MultiFactorSignatureApprovalResultFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("MultiFactorSignatureApprovalResultFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("MultiFactorSignatureApprovalResultFromDTO() returned nil for non-nil input")
			}
		})
	}
}

func TestEntityTypeToDTO(t *testing.T) {
	entityType := model.MultiFactorSignatureEntityType("WHITELISTED_ADDRESS")
	result := EntityTypeToDTO(entityType)
	if string(result) != "WHITELISTED_ADDRESS" {
		t.Errorf("EntityTypeToDTO() = %v, want WHITELISTED_ADDRESS", result)
	}
}
