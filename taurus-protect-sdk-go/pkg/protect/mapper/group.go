package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// GroupFromDTO converts an OpenAPI InternalGroup to a domain Group.
func GroupFromDTO(dto *openapi.TgvalidatordInternalGroup) *model.Group {
	if dto == nil {
		return nil
	}

	group := &model.Group{
		ID:              safeString(dto.Id),
		TenantID:        safeString(dto.TenantId),
		ExternalGroupID: safeString(dto.ExternalGroupId),
		Name:            safeString(dto.Name),
		Email:           safeString(dto.Email),
		Description:     safeString(dto.Description),
		EnforcedInRules: safeBool(dto.EnforcedInRules),
	}

	// Convert timestamps
	if dto.CreationDate != nil {
		group.CreatedAt = *dto.CreationDate
	}
	if dto.UpdateDate != nil {
		group.UpdatedAt = *dto.UpdateDate
	}

	// Convert users
	if dto.Users != nil {
		group.Users = make([]model.GroupUser, len(dto.Users))
		for i, user := range dto.Users {
			group.Users[i] = GroupUserFromDTO(&user)
		}
	}

	return group
}

// GroupsFromDTO converts a slice of OpenAPI InternalGroup to domain Groups.
func GroupsFromDTO(dtos []openapi.TgvalidatordInternalGroup) []*model.Group {
	if dtos == nil {
		return nil
	}
	groups := make([]*model.Group, len(dtos))
	for i := range dtos {
		groups[i] = GroupFromDTO(&dtos[i])
	}
	return groups
}

// GroupUserFromDTO converts an OpenAPI InternalGroupUser to a domain GroupUser.
func GroupUserFromDTO(dto *openapi.TgvalidatordInternalGroupUser) model.GroupUser {
	if dto == nil {
		return model.GroupUser{}
	}
	return model.GroupUser{
		ID:              safeString(dto.Id),
		ExternalUserID:  safeString(dto.ExternalUserId),
		EnforcedInRules: safeBool(dto.EnforcedInRules),
	}
}
