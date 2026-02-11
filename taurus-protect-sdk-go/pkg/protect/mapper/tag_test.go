package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestTagFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTag
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns tag with zero values",
			dto:  &openapi.TgvalidatordTag{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordTag {
				id := "tag-123"
				value := "Production"
				color := "#FF0000"
				now := time.Now()
				return &openapi.TgvalidatordTag{
					Id:           &id,
					Value:        &value,
					Color:        &color,
					CreationDate: &now,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TagFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("TagFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("TagFromDTO() returned nil for non-nil input")
			}
			// Verify specific fields for complete DTO
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Value != nil && got.Value != *tt.dto.Value {
				t.Errorf("Value = %v, want %v", got.Value, *tt.dto.Value)
			}
			if tt.dto.Color != nil && got.Color != *tt.dto.Color {
				t.Errorf("Color = %v, want %v", got.Color, *tt.dto.Color)
			}
			if tt.dto.CreationDate != nil && !got.CreationDate.Equal(*tt.dto.CreationDate) {
				t.Errorf("CreationDate = %v, want %v", got.CreationDate, *tt.dto.CreationDate)
			}
		})
	}
}

func TestTagFromDTO_PartialFields(t *testing.T) {
	// Test with only some fields set
	id := "tag-456"
	dto := &openapi.TgvalidatordTag{
		Id: &id,
	}

	got := TagFromDTO(dto)
	if got == nil {
		t.Fatal("TagFromDTO() returned nil for non-nil input")
	}
	if got.ID != id {
		t.Errorf("ID = %v, want %v", got.ID, id)
	}
	if got.Value != "" {
		t.Errorf("Value = %v, want empty string", got.Value)
	}
	if got.Color != "" {
		t.Errorf("Color = %v, want empty string", got.Color)
	}
	if !got.CreationDate.IsZero() {
		t.Errorf("CreationDate = %v, want zero time", got.CreationDate)
	}
}

func TestTagsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordTag
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordTag{},
			want: 0,
		},
		{
			name: "converts multiple tags",
			dtos: func() []openapi.TgvalidatordTag {
				id1 := "tag-1"
				id2 := "tag-2"
				id3 := "tag-3"
				return []openapi.TgvalidatordTag{
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
			got := TagsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("TagsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("TagsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestTagsFromDTO_PreservesOrder(t *testing.T) {
	id1 := "tag-first"
	id2 := "tag-second"
	id3 := "tag-third"
	dtos := []openapi.TgvalidatordTag{
		{Id: &id1},
		{Id: &id2},
		{Id: &id3},
	}

	got := TagsFromDTO(dtos)
	if len(got) != 3 {
		t.Fatalf("TagsFromDTO() length = %v, want 3", len(got))
	}
	if got[0].ID != id1 {
		t.Errorf("First tag ID = %v, want %v", got[0].ID, id1)
	}
	if got[1].ID != id2 {
		t.Errorf("Second tag ID = %v, want %v", got[1].ID, id2)
	}
	if got[2].ID != id3 {
		t.Errorf("Third tag ID = %v, want %v", got[2].ID, id3)
	}
}
