package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model/taurusnetwork"
)

// PledgeFromDTO converts an OpenAPI TgvalidatordTnPledge to a domain Pledge.
func PledgeFromDTO(dto *openapi.TgvalidatordTnPledge) *taurusnetwork.Pledge {
	if dto == nil {
		return nil
	}

	pledge := &taurusnetwork.Pledge{
		ID:                  safeString(dto.Id),
		OwnerParticipantID:  safeString(dto.OwnerParticipantID),
		TargetParticipantID: safeString(dto.TargetParticipantID),
		SharedAddressID:     safeString(dto.SharedAddressID),
		CurrencyID:          safeString(dto.CurrencyID),
		Blockchain:          safeString(dto.Blockchain),
		Network:             safeString(dto.Network),
		Amount:              safeString(dto.Amount),
		Status:              safeString(dto.Status),
		PledgeType:          safeString(dto.PledgeType),
		Direction:           safeString(dto.Direction),
		ExternalReferenceID: safeString(dto.ExternalReferenceId),
		ReconciliationNote:  safeString(dto.ReconciliationNote),
		WladdressID:         safeString(dto.WladdressID),
		DurationSetup:       PledgeDurationSetupFromDTO(dto.DurationSetup),
		Attributes:          PledgeAttributesFromDTO(dto.Attributes),
		Trails:              PledgeTrailsFromDTO(dto.Trails),
		OriginCreationDate:  safeTime(dto.OriginCreationDate),
		UnpledgeDate:        safeTime(dto.UnpledgeDate),
		CreatedAt:           safeTime(dto.CreatedAt),
		UpdatedAt:           safeTime(dto.UpdatedAt),
	}

	return pledge
}

// PledgesFromDTO converts a slice of OpenAPI TgvalidatordTnPledge to domain Pledges.
func PledgesFromDTO(dtos []openapi.TgvalidatordTnPledge) []*taurusnetwork.Pledge {
	if dtos == nil {
		return nil
	}

	pledges := make([]*taurusnetwork.Pledge, len(dtos))
	for i := range dtos {
		pledges[i] = PledgeFromDTO(&dtos[i])
	}
	return pledges
}

// PledgeDurationSetupFromDTO converts an OpenAPI TnPledgePledgeDurationSetup to a domain PledgeDurationSetup.
func PledgeDurationSetupFromDTO(dto *openapi.TnPledgePledgeDurationSetup) *taurusnetwork.PledgeDurationSetup {
	if dto == nil {
		return nil
	}

	return &taurusnetwork.PledgeDurationSetup{
		MinimumDuration:          safeString(dto.MinimumDuration),
		EndOfMinimumDurationDate: safeTime(dto.EndOfMinimumDurationDate),
		NoticePeriodDuration:     safeString(dto.NoticePeriodDuration),
		EndOfNoticePeriodDate:    safeTime(dto.EndOfNoticePeriodDate),
	}
}

// PledgeAttributeFromDTO converts an OpenAPI TnPledgePledgeAttribute to a domain PledgeAttribute.
func PledgeAttributeFromDTO(dto *openapi.TnPledgePledgeAttribute) *taurusnetwork.PledgeAttribute {
	if dto == nil {
		return nil
	}

	return &taurusnetwork.PledgeAttribute{
		ID:                    safeString(dto.Id),
		Key:                   safeString(dto.Key),
		Value:                 safeString(dto.Value),
		Owner:                 safeString(dto.Owner),
		Type:                  safeString(dto.Type),
		Subtype:               safeString(dto.Subtype),
		ContentType:           safeString(dto.ContentType),
		IsTaurusNetworkShared: safeBool(dto.IsTaurusNetworkShared),
	}
}

// PledgeAttributesFromDTO converts a slice of OpenAPI TnPledgePledgeAttribute to domain PledgeAttributes.
func PledgeAttributesFromDTO(dtos []openapi.TnPledgePledgeAttribute) []taurusnetwork.PledgeAttribute {
	if dtos == nil {
		return nil
	}

	attributes := make([]taurusnetwork.PledgeAttribute, len(dtos))
	for i := range dtos {
		attr := PledgeAttributeFromDTO(&dtos[i])
		if attr != nil {
			attributes[i] = *attr
		}
	}
	return attributes
}

// PledgeTrailFromDTO converts an OpenAPI TgvalidatordTnPledgeTrail to a domain PledgeTrail.
func PledgeTrailFromDTO(dto *openapi.TgvalidatordTnPledgeTrail) *taurusnetwork.PledgeTrail {
	if dto == nil {
		return nil
	}

	return &taurusnetwork.PledgeTrail{
		ID:               safeString(dto.Id),
		PledgeID:         safeString(dto.PledgeID),
		AddressCommandID: safeString(dto.AddressCommandID),
		ParticipantID:    safeString(dto.ParticipantID),
		PledgeAmount:     safeString(dto.PledgeAmount),
		Action:           safeString(dto.Action),
		Comment:          safeString(dto.Comment),
		CreatedAt:        safeTime(dto.CreatedAt),
	}
}

// PledgeTrailsFromDTO converts a slice of OpenAPI TgvalidatordTnPledgeTrail to domain PledgeTrails.
func PledgeTrailsFromDTO(dtos []openapi.TgvalidatordTnPledgeTrail) []taurusnetwork.PledgeTrail {
	if dtos == nil {
		return nil
	}

	trails := make([]taurusnetwork.PledgeTrail, len(dtos))
	for i := range dtos {
		trail := PledgeTrailFromDTO(&dtos[i])
		if trail != nil {
			trails[i] = *trail
		}
	}
	return trails
}

// PledgeWithdrawalFromDTO converts an OpenAPI TgvalidatordTnPledgeWithdrawal to a domain PledgeWithdrawal.
func PledgeWithdrawalFromDTO(dto *openapi.TgvalidatordTnPledgeWithdrawal) *taurusnetwork.PledgeWithdrawal {
	if dto == nil {
		return nil
	}

	return &taurusnetwork.PledgeWithdrawal{
		ID:                         safeString(dto.Id),
		PledgeID:                   safeString(dto.PledgeID),
		DestinationSharedAddressID: safeString(dto.DestinationSharedAddressID),
		Amount:                     safeString(dto.Amount),
		Status:                     safeString(dto.Status),
		TxHash:                     safeString(dto.TxHash),
		TxID:                       safeString(dto.TxID),
		RequestID:                  safeString(dto.RequestID),
		TxBlockNumber:              safeString(dto.TxBlockNumber),
		InitiatorParticipantID:     safeString(dto.InitiatorParticipantID),
		ExternalReferenceID:        safeString(dto.ExternalReferenceID),
		Trails:                     PledgeWithdrawalTrailsFromDTO(dto.Trails),
		CreatedAt:                  safeTime(dto.CreatedAt),
	}
}

// PledgeWithdrawalsFromDTO converts a slice of OpenAPI TgvalidatordTnPledgeWithdrawal to domain PledgeWithdrawals.
func PledgeWithdrawalsFromDTO(dtos []openapi.TgvalidatordTnPledgeWithdrawal) []taurusnetwork.PledgeWithdrawal {
	if dtos == nil {
		return nil
	}

	withdrawals := make([]taurusnetwork.PledgeWithdrawal, len(dtos))
	for i := range dtos {
		withdrawal := PledgeWithdrawalFromDTO(&dtos[i])
		if withdrawal != nil {
			withdrawals[i] = *withdrawal
		}
	}
	return withdrawals
}

// PledgeWithdrawalTrailFromDTO converts an OpenAPI TgvalidatordTnPledgeWithdrawalTrail to a domain PledgeWithdrawalTrail.
func PledgeWithdrawalTrailFromDTO(dto *openapi.TgvalidatordTnPledgeWithdrawalTrail) *taurusnetwork.PledgeWithdrawalTrail {
	if dto == nil {
		return nil
	}

	return &taurusnetwork.PledgeWithdrawalTrail{
		ID:                 safeString(dto.Id),
		PledgeWithdrawalID: safeString(dto.PledgeWithdrawalID),
		AddressCommandID:   safeString(dto.AddressCommandID),
		ParticipantID:      safeString(dto.ParticipantID),
		Action:             safeString(dto.Action),
		Comment:            safeString(dto.Comment),
		CreatedAt:          safeTime(dto.CreatedAt),
	}
}

// PledgeWithdrawalTrailsFromDTO converts a slice of OpenAPI TgvalidatordTnPledgeWithdrawalTrail to domain PledgeWithdrawalTrails.
func PledgeWithdrawalTrailsFromDTO(dtos []openapi.TgvalidatordTnPledgeWithdrawalTrail) []taurusnetwork.PledgeWithdrawalTrail {
	if dtos == nil {
		return nil
	}

	trails := make([]taurusnetwork.PledgeWithdrawalTrail, len(dtos))
	for i := range dtos {
		trail := PledgeWithdrawalTrailFromDTO(&dtos[i])
		if trail != nil {
			trails[i] = *trail
		}
	}
	return trails
}

// PledgeActionFromDTO converts an OpenAPI TgvalidatordTnPledgeAction to a domain PledgeAction.
func PledgeActionFromDTO(dto *openapi.TgvalidatordTnPledgeAction) *taurusnetwork.PledgeAction {
	if dto == nil {
		return nil
	}

	return &taurusnetwork.PledgeAction{
		ID:                 safeString(dto.Id),
		PledgeID:           safeString(dto.PledgeID),
		ActionType:         safeString(dto.ActionType),
		Status:             safeString(dto.Status),
		Metadata:           MetadataFromDTO(dto.Metadata),
		Rule:               safeString(dto.Rule),
		Approvers:          ApproversFromDTO(dto.Approvers),
		NeedsApprovalFrom:  dto.NeedsApprovalFrom,
		Envelope:           safeString(dto.Envelope),
		PledgeWithdrawalID: safeString(dto.PledgeWithdrawalID),
		Trails:             PledgeActionTrailsFromDTO(dto.Trails),
		CreatedAt:          safeTime(dto.CreatedAt),
		LastApprovalDate:   safeTime(dto.LastApprovalDate),
	}
}

// PledgeActionsFromDTO converts a slice of OpenAPI TgvalidatordTnPledgeAction to domain PledgeActions.
func PledgeActionsFromDTO(dtos []openapi.TgvalidatordTnPledgeAction) []taurusnetwork.PledgeAction {
	if dtos == nil {
		return nil
	}

	actions := make([]taurusnetwork.PledgeAction, len(dtos))
	for i := range dtos {
		action := PledgeActionFromDTO(&dtos[i])
		if action != nil {
			actions[i] = *action
		}
	}
	return actions
}

// Note: MetadataFromDTO is defined in request.go
// Note: ApproversFromDTO is defined in whitelisted_asset.go

// PledgeActionTrailFromDTO converts an OpenAPI TgvalidatordTnPledgeActionTrail to a domain PledgeActionTrail.
func PledgeActionTrailFromDTO(dto *openapi.TgvalidatordTnPledgeActionTrail) *taurusnetwork.PledgeActionTrail {
	if dto == nil {
		return nil
	}

	return &taurusnetwork.PledgeActionTrail{
		ID:             safeString(dto.Id),
		PledgeActionID: safeString(dto.PledgeActionID),
		UserID:         safeString(dto.UserID),
		ExternalUserID: safeString(dto.ExternalUserID),
		Action:         safeString(dto.Action),
		Comment:        safeString(dto.Comment),
		CreatedAt:      safeTime(dto.CreatedAt),
	}
}

// PledgeActionTrailsFromDTO converts a slice of OpenAPI TgvalidatordTnPledgeActionTrail to domain PledgeActionTrails.
func PledgeActionTrailsFromDTO(dtos []openapi.TgvalidatordTnPledgeActionTrail) []taurusnetwork.PledgeActionTrail {
	if dtos == nil {
		return nil
	}

	trails := make([]taurusnetwork.PledgeActionTrail, len(dtos))
	for i := range dtos {
		trail := PledgeActionTrailFromDTO(&dtos[i])
		if trail != nil {
			trails[i] = *trail
		}
	}
	return trails
}

// Note: CursorPaginationFromDTO is defined in change.go
