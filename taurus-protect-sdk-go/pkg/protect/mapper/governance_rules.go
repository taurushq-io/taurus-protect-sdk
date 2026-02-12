package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// GovernanceRulesetFromDTO converts an OpenAPI TgvalidatordRules to a domain GovernanceRuleset.
func GovernanceRulesetFromDTO(dto *openapi.TgvalidatordRules) *model.GovernanceRuleset {
	if dto == nil {
		return nil
	}

	rules := &model.GovernanceRuleset{
		RulesContainer: safeString(dto.RulesContainer),
		Locked:         safeBool(dto.Locked),
	}

	if dto.CreationDate != nil {
		rules.CreatedAt = *dto.CreationDate
	}
	if dto.UpdateDate != nil {
		rules.UpdatedAt = *dto.UpdateDate
	}

	if dto.RulesSignatures != nil {
		rules.Signatures = make([]model.RuleUserSignature, len(dto.RulesSignatures))
		for i, sig := range dto.RulesSignatures {
			rules.Signatures[i] = RuleUserSignatureFromDTO(&sig)
		}
	}

	if dto.Trails != nil {
		rules.Trails = make([]model.RulesTrail, len(dto.Trails))
		for i, trail := range dto.Trails {
			rules.Trails[i] = RulesTrailFromDTO(&trail)
		}
	}

	return rules
}

// GovernanceRulesetsFromDTO converts a slice of OpenAPI TgvalidatordRules to domain GovernanceRuleset.
func GovernanceRulesetsFromDTO(dtos []openapi.TgvalidatordRules) []*model.GovernanceRuleset {
	if dtos == nil {
		return nil
	}

	result := make([]*model.GovernanceRuleset, len(dtos))
	for i := range dtos {
		result[i] = GovernanceRulesetFromDTO(&dtos[i])
	}
	return result
}

// RuleUserSignatureFromDTO converts an OpenAPI TgvalidatordRuleUserSignature to a domain RuleUserSignature.
func RuleUserSignatureFromDTO(dto *openapi.TgvalidatordRuleUserSignature) model.RuleUserSignature {
	if dto == nil {
		return model.RuleUserSignature{}
	}

	return model.RuleUserSignature{
		UserID:    safeString(dto.UserId),
		Signature: safeString(dto.Signature),
	}
}

// RulesTrailFromDTO converts an OpenAPI TgvalidatordRulesTrail to a domain RulesTrail.
func RulesTrailFromDTO(dto *openapi.TgvalidatordRulesTrail) model.RulesTrail {
	if dto == nil {
		return model.RulesTrail{}
	}

	trail := model.RulesTrail{
		ID:             safeString(dto.Id),
		UserID:         safeString(dto.UserId),
		ExternalUserID: safeString(dto.ExternalUserId),
		Action:         safeString(dto.Action),
		Comment:        safeString(dto.Comment),
	}

	if dto.Date != nil {
		trail.Date = *dto.Date
	}

	return trail
}

// SuperAdminPublicKeyFromDTO converts an OpenAPI GetPublicKeysReplyPublicKey to a domain SuperAdminPublicKey.
func SuperAdminPublicKeyFromDTO(dto *openapi.GetPublicKeysReplyPublicKey) *model.SuperAdminPublicKey {
	if dto == nil {
		return nil
	}

	return &model.SuperAdminPublicKey{
		UserID:    safeString(dto.UserID),
		PublicKey: safeString(dto.PublicKey),
	}
}

// SuperAdminPublicKeysFromDTO converts a slice of OpenAPI GetPublicKeysReplyPublicKey to domain SuperAdminPublicKey.
func SuperAdminPublicKeysFromDTO(dtos []openapi.GetPublicKeysReplyPublicKey) []*model.SuperAdminPublicKey {
	if dtos == nil {
		return nil
	}

	result := make([]*model.SuperAdminPublicKey, len(dtos))
	for i := range dtos {
		result[i] = SuperAdminPublicKeyFromDTO(&dtos[i])
	}
	return result
}
