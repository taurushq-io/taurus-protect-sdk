package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// ChangeFromDTO converts an OpenAPI Change to a domain Change.
func ChangeFromDTO(dto *openapi.TgvalidatordChange) *model.Change {
	if dto == nil {
		return nil
	}

	change := &model.Change{
		ID:                safeString(dto.Id),
		TenantID:          safeString(dto.TenantId),
		CreatorID:         safeString(dto.CreatorId),
		CreatorExternalID: safeString(dto.CreatorExternalId),
		Action:            safeString(dto.Action),
		EntityID:          safeString(dto.EntityId),
		EntityUUID:        safeString(dto.EntityUUID),
		Entity:            safeString(dto.Entity),
		Comment:           safeString(dto.Comment),
	}

	// Copy changes map
	if dto.Changes != nil {
		change.Changes = make(map[string]string)
		for k, v := range *dto.Changes {
			change.Changes[k] = v
		}
	}

	// Convert creation date
	if dto.CreationDate != nil {
		change.CreationDate = *dto.CreationDate
	}

	return change
}

// ChangesFromDTO converts a slice of OpenAPI Changes to domain Changes.
func ChangesFromDTO(dtos []openapi.TgvalidatordChange) []*model.Change {
	if dtos == nil {
		return nil
	}
	changes := make([]*model.Change, len(dtos))
	for i := range dtos {
		changes[i] = ChangeFromDTO(&dtos[i])
	}
	return changes
}

// CursorPaginationFromDTO converts an OpenAPI ResponseCursor to a domain CursorPagination.
func CursorPaginationFromDTO(dto *openapi.TgvalidatordResponseCursor) *model.CursorPagination {
	if dto == nil {
		return nil
	}

	return &model.CursorPagination{
		CurrentPage: safeString(dto.CurrentPage),
		HasPrevious: safeBool(dto.HasPrevious),
		HasNext:     safeBool(dto.HasNext),
	}
}

// CreateChangeResultFromDTO converts an OpenAPI CreateChangeResult to a domain CreateChangeResult.
func CreateChangeResultFromDTO(dto *openapi.TgvalidatordCreateChangeResult) *model.CreateChangeResult {
	if dto == nil {
		return nil
	}

	return &model.CreateChangeResult{
		ID: safeString(dto.Id),
	}
}
