package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// AuditTrailFromDTO converts an OpenAPI AuditTrail to a domain AuditTrail.
func AuditTrailFromDTO(dto *openapi.TgvalidatordAuditTrail) *model.AuditTrail {
	if dto == nil {
		return nil
	}

	trail := &model.AuditTrail{
		ID:        safeString(dto.Id),
		Entity:    safeString(dto.Entity),
		Action:    safeString(dto.Action),
		SubAction: safeString(dto.SubAction),
		Details:   safeString(dto.Details),
	}

	// Convert user info
	if dto.User != nil {
		trail.User = UserInfoFromDTO(dto.User)
	}

	// Convert creation date
	if dto.CreationDate != nil {
		trail.CreationDate = *dto.CreationDate
	}

	return trail
}

// AuditTrailsFromDTO converts a slice of OpenAPI AuditTrail to domain AuditTrails.
func AuditTrailsFromDTO(dtos []openapi.TgvalidatordAuditTrail) []*model.AuditTrail {
	if dtos == nil {
		return nil
	}
	trails := make([]*model.AuditTrail, len(dtos))
	for i := range dtos {
		trails[i] = AuditTrailFromDTO(&dtos[i])
	}
	return trails
}

// UserInfoFromDTO converts an OpenAPI UserInfo to a domain UserInfo.
func UserInfoFromDTO(dto *openapi.TgvalidatordUserInfo) *model.UserInfo {
	if dto == nil {
		return nil
	}

	return &model.UserInfo{
		ID:             safeString(dto.Id),
		ExternalUserID: safeString(dto.ExternalUserId),
		Email:          safeString(dto.Email),
		Deleted:        safeBool(dto.Deleted),
	}
}
