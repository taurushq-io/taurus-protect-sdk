package mapper

import (
	"strconv"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// VisibilityGroupFromDTO converts an OpenAPI InternalVisibilityGroup to a domain VisibilityGroup.
func VisibilityGroupFromDTO(dto *openapi.TgvalidatordInternalVisibilityGroup) *model.VisibilityGroup {
	if dto == nil {
		return nil
	}

	vg := &model.VisibilityGroup{
		ID:          safeString(dto.Id),
		TenantID:    safeString(dto.TenantId),
		Name:        safeString(dto.Name),
		Description: safeString(dto.Description),
	}

	// Parse user count
	if dto.UserCount != nil {
		if count, err := strconv.ParseInt(*dto.UserCount, 10, 64); err == nil {
			vg.UserCount = count
		}
	}

	// Convert timestamps
	if dto.CreationDate != nil {
		vg.CreatedAt = *dto.CreationDate
	}
	if dto.UpdateDate != nil {
		vg.UpdatedAt = *dto.UpdateDate
	}

	// Convert users
	if dto.Users != nil {
		vg.Users = make([]*model.VisibilityGroupUser, len(dto.Users))
		for i := range dto.Users {
			vg.Users[i] = VisibilityGroupUserFromDTO(&dto.Users[i])
		}
	}

	return vg
}

// VisibilityGroupsFromDTO converts a slice of OpenAPI InternalVisibilityGroup to domain VisibilityGroups.
func VisibilityGroupsFromDTO(dtos []openapi.TgvalidatordInternalVisibilityGroup) []*model.VisibilityGroup {
	if dtos == nil {
		return nil
	}
	groups := make([]*model.VisibilityGroup, len(dtos))
	for i := range dtos {
		groups[i] = VisibilityGroupFromDTO(&dtos[i])
	}
	return groups
}

// VisibilityGroupUserFromDTO converts an OpenAPI InternalVisibilityGroupUser to a domain VisibilityGroupUser.
func VisibilityGroupUserFromDTO(dto *openapi.TgvalidatordInternalVisibilityGroupUser) *model.VisibilityGroupUser {
	if dto == nil {
		return nil
	}

	return &model.VisibilityGroupUser{
		ID:             safeString(dto.Id),
		ExternalUserID: safeString(dto.ExternalUserId),
	}
}

// Note: UserFromDTO and UsersFromDTO are defined in user.go
