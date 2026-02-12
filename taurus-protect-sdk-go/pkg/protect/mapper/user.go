package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// UserFromDTO converts an OpenAPI InternalUser to a domain User.
func UserFromDTO(dto *openapi.TgvalidatordInternalUser) *model.User {
	if dto == nil {
		return nil
	}

	user := &model.User{
		ID:                       safeString(dto.Id),
		TenantID:                 safeString(dto.TenantId),
		ExternalUserID:           safeString(dto.ExternalUserId),
		Username:                 safeString(dto.Username),
		Email:                    safeString(dto.Email),
		FirstName:                safeString(dto.FirstName),
		LastName:                 safeString(dto.LastName),
		Status:                   safeString(dto.Status),
		PublicKey:                safeString(dto.PublicKey),
		PasswordChanged:          safeBool(dto.PasswordChanged),
		TotpEnabled:              safeBool(dto.TotpEnabled),
		EnforcedInRules:          safeBool(dto.EnforcedInRules),
		PublicKeyEnforcedInRules: safeBool(dto.PublicKeyEnforcedInRules),
	}

	// Copy roles
	if dto.Roles != nil {
		user.Roles = make([]string, len(dto.Roles))
		copy(user.Roles, dto.Roles)
	}

	// Convert groups
	if dto.Groups != nil {
		user.Groups = make([]model.UserGroup, len(dto.Groups))
		for i, g := range dto.Groups {
			user.Groups[i] = UserGroupFromDTO(&g)
		}
	}

	// Convert attributes
	if dto.Attributes != nil {
		user.Attributes = make([]model.UserAttribute, len(dto.Attributes))
		for i, attr := range dto.Attributes {
			user.Attributes[i] = UserAttributeFromDTO(&attr)
		}
	}

	// Convert timestamps
	if dto.CreationDate != nil {
		user.CreatedAt = *dto.CreationDate
	}
	if dto.UpdateDate != nil {
		user.UpdatedAt = *dto.UpdateDate
	}
	if dto.LastLogin != nil {
		user.LastLogin = *dto.LastLogin
	}

	return user
}

// UsersFromDTO converts a slice of OpenAPI InternalUser to domain Users.
func UsersFromDTO(dtos []openapi.TgvalidatordInternalUser) []*model.User {
	if dtos == nil {
		return nil
	}
	users := make([]*model.User, len(dtos))
	for i := range dtos {
		users[i] = UserFromDTO(&dtos[i])
	}
	return users
}

// UserGroupFromDTO converts an OpenAPI InternalUserGroup to a domain UserGroup.
func UserGroupFromDTO(dto *openapi.TgvalidatordInternalUserGroup) model.UserGroup {
	if dto == nil {
		return model.UserGroup{}
	}
	return model.UserGroup{
		ID:              safeString(dto.Id),
		ExternalGroupID: safeString(dto.ExternalGroupId),
		EnforcedInRules: safeBool(dto.EnforcedInRules),
	}
}

// UserGroupsFromDTO converts a slice of OpenAPI InternalUserGroup to domain UserGroups.
func UserGroupsFromDTO(dtos []openapi.TgvalidatordInternalUserGroup) []model.UserGroup {
	if dtos == nil {
		return nil
	}
	groups := make([]model.UserGroup, len(dtos))
	for i := range dtos {
		groups[i] = UserGroupFromDTO(&dtos[i])
	}
	return groups
}

// UserAttributeFromDTO converts an OpenAPI InternalUserAttribute to a domain UserAttribute.
func UserAttributeFromDTO(dto *openapi.InternalUserAttribute) model.UserAttribute {
	if dto == nil {
		return model.UserAttribute{}
	}
	attr := model.UserAttribute{
		ID:          safeString(dto.Id),
		Key:         safeString(dto.Key),
		Value:       safeString(dto.Value),
		ContentType: safeString(dto.ContentType),
		Owner:       safeString(dto.Owner),
		Type:        safeString(dto.Type),
		SubType:     safeString(dto.SubType),
		IsFile:      safeBool(dto.IsFile),
	}

	if dto.CreationDate != nil {
		attr.CreatedAt = *dto.CreationDate
	}
	if dto.UpdateDate != nil {
		attr.UpdatedAt = *dto.UpdateDate
	}

	return attr
}

// UserAttributesFromDTO converts a slice of OpenAPI InternalUserAttribute to domain UserAttributes.
func UserAttributesFromDTO(dtos []openapi.InternalUserAttribute) []model.UserAttribute {
	if dtos == nil {
		return nil
	}
	attrs := make([]model.UserAttribute, len(dtos))
	for i := range dtos {
		attrs[i] = UserAttributeFromDTO(&dtos[i])
	}
	return attrs
}
