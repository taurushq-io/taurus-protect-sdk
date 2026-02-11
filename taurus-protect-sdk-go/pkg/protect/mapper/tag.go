package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// TagFromDTO converts an OpenAPI Tag to a domain Tag.
func TagFromDTO(dto *openapi.TgvalidatordTag) *model.Tag {
	if dto == nil {
		return nil
	}

	tag := &model.Tag{
		ID:    safeString(dto.Id),
		Value: safeString(dto.Value),
		Color: safeString(dto.Color),
	}

	// Convert creation date
	if dto.CreationDate != nil {
		tag.CreationDate = *dto.CreationDate
	}

	return tag
}

// TagsFromDTO converts a slice of OpenAPI Tags to domain Tags.
func TagsFromDTO(dtos []openapi.TgvalidatordTag) []*model.Tag {
	if dtos == nil {
		return nil
	}
	tags := make([]*model.Tag, len(dtos))
	for i := range dtos {
		tags[i] = TagFromDTO(&dtos[i])
	}
	return tags
}
