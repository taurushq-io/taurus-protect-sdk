package mapper

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestCryptoPunkMetadataFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordCryptoPunkMetadata
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns metadata with zero values",
			dto:  &openapi.TgvalidatordCryptoPunkMetadata{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordCryptoPunkMetadata {
				punkId := "1234"
				attributes := "Male, Mohawk, Smile"
				image := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
				return &openapi.TgvalidatordCryptoPunkMetadata{
					PunkId:         &punkId,
					PunkAttributes: &attributes,
					Image:          &image,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CryptoPunkMetadataFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("CryptoPunkMetadataFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("CryptoPunkMetadataFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.PunkId != nil && got.PunkID != *tt.dto.PunkId {
				t.Errorf("PunkID = %v, want %v", got.PunkID, *tt.dto.PunkId)
			}
			if tt.dto.PunkAttributes != nil && got.Attributes != *tt.dto.PunkAttributes {
				t.Errorf("Attributes = %v, want %v", got.Attributes, *tt.dto.PunkAttributes)
			}
			if tt.dto.Image != nil && got.Image != *tt.dto.Image {
				t.Errorf("Image = %v, want %v", got.Image, *tt.dto.Image)
			}
		})
	}
}

func TestCryptoPunkMetadataFromDTO_PartialFields(t *testing.T) {
	punkId := "42"
	dto := &openapi.TgvalidatordCryptoPunkMetadata{
		PunkId: &punkId,
		// Other fields are nil
	}

	got := CryptoPunkMetadataFromDTO(dto)
	if got == nil {
		t.Fatal("CryptoPunkMetadataFromDTO() returned nil for non-nil input")
	}
	if got.PunkID != "42" {
		t.Errorf("PunkID = %v, want 42", got.PunkID)
	}
	if got.Attributes != "" {
		t.Errorf("Attributes should be empty string when nil, got %v", got.Attributes)
	}
	if got.Image != "" {
		t.Errorf("Image should be empty string when nil, got %v", got.Image)
	}
}

func TestERCTokenMetadataFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordERCTokenMetadata
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns metadata with zero values",
			dto:  &openapi.TgvalidatordERCTokenMetadata{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordERCTokenMetadata {
				name := "Bored Ape #1234"
				description := "A unique Bored Ape from the BAYC collection"
				decimals := "0"
				dataType := "image/png"
				base64Data := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
				uri := "ipfs://QmeSjSinHpPnmXmspMjwiXyN6zS4E9zccariGR3jxcaWtq/1234"
				return &openapi.TgvalidatordERCTokenMetadata{
					Name:        &name,
					Description: &description,
					Decimals:    &decimals,
					DataType:    &dataType,
					Base64Data:  &base64Data,
					Uri:         &uri,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ERCTokenMetadataFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ERCTokenMetadataFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ERCTokenMetadataFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Name != nil && got.Name != *tt.dto.Name {
				t.Errorf("Name = %v, want %v", got.Name, *tt.dto.Name)
			}
			if tt.dto.Description != nil && got.Description != *tt.dto.Description {
				t.Errorf("Description = %v, want %v", got.Description, *tt.dto.Description)
			}
			if tt.dto.Decimals != nil && got.Decimals != *tt.dto.Decimals {
				t.Errorf("Decimals = %v, want %v", got.Decimals, *tt.dto.Decimals)
			}
			if tt.dto.DataType != nil && got.DataType != *tt.dto.DataType {
				t.Errorf("DataType = %v, want %v", got.DataType, *tt.dto.DataType)
			}
			if tt.dto.Base64Data != nil && got.Base64Data != *tt.dto.Base64Data {
				t.Errorf("Base64Data = %v, want %v", got.Base64Data, *tt.dto.Base64Data)
			}
			if tt.dto.Uri != nil && got.URI != *tt.dto.Uri {
				t.Errorf("URI = %v, want %v", got.URI, *tt.dto.Uri)
			}
		})
	}
}

func TestERCTokenMetadataFromDTO_PartialFields(t *testing.T) {
	name := "My NFT"
	uri := "https://example.com/token/1"
	dto := &openapi.TgvalidatordERCTokenMetadata{
		Name: &name,
		Uri:  &uri,
		// Other fields are nil
	}

	got := ERCTokenMetadataFromDTO(dto)
	if got == nil {
		t.Fatal("ERCTokenMetadataFromDTO() returned nil for non-nil input")
	}
	if got.Name != "My NFT" {
		t.Errorf("Name = %v, want My NFT", got.Name)
	}
	if got.URI != "https://example.com/token/1" {
		t.Errorf("URI = %v, want https://example.com/token/1", got.URI)
	}
	if got.Description != "" {
		t.Errorf("Description should be empty string when nil, got %v", got.Description)
	}
	if got.Decimals != "" {
		t.Errorf("Decimals should be empty string when nil, got %v", got.Decimals)
	}
	if got.DataType != "" {
		t.Errorf("DataType should be empty string when nil, got %v", got.DataType)
	}
	if got.Base64Data != "" {
		t.Errorf("Base64Data should be empty string when nil, got %v", got.Base64Data)
	}
}
