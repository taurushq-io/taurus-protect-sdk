package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// ActionFromDTO converts an OpenAPI ActionEnvelope to a domain Action.
func ActionFromDTO(dto *openapi.TgvalidatordActionEnvelope) *model.Action {
	if dto == nil {
		return nil
	}

	action := &model.Action{
		ID:          safeString(dto.Id),
		TenantID:    safeString(dto.TenantId),
		Label:       safeString(dto.Label),
		Status:      safeString(dto.Status),
		AutoApprove: safeBool(dto.AutoApprove),
	}

	// Convert dates
	if dto.CreationDate != nil {
		action.CreationDate = *dto.CreationDate
	}
	if dto.UpdateDate != nil {
		action.UpdateDate = *dto.UpdateDate
	}
	if dto.Lastcheckeddate != nil {
		action.LastCheckedDate = *dto.Lastcheckeddate
	}

	// Convert action details
	if dto.Action != nil {
		action.ActionDetails = ActionDetailsFromDTO(dto.Action)
	}

	// Convert attributes
	action.Attributes = ActionAttributesFromDTO(dto.Attributes)

	// Convert trails
	action.Trails = ActionTrailsFromDTO(dto.Trails)

	return action
}

// ActionsFromDTO converts a slice of OpenAPI ActionEnvelope to domain Actions.
func ActionsFromDTO(dtos []openapi.TgvalidatordActionEnvelope) []*model.Action {
	if dtos == nil {
		return nil
	}
	actions := make([]*model.Action, len(dtos))
	for i := range dtos {
		actions[i] = ActionFromDTO(&dtos[i])
	}
	return actions
}

// ActionDetailsFromDTO converts an OpenAPI Action to domain ActionDetails.
func ActionDetailsFromDTO(dto *openapi.TgvalidatordAction) *model.ActionDetails {
	if dto == nil {
		return nil
	}

	details := &model.ActionDetails{}

	if dto.Trigger != nil {
		details.Trigger = ActionTriggerFromDTO(dto.Trigger)
	}

	details.Tasks = ActionTasksFromDTO(dto.Tasks)

	return details
}

// ActionTriggerFromDTO converts an OpenAPI ActionTrigger to domain ActionTrigger.
func ActionTriggerFromDTO(dto *openapi.ActionTrigger) *model.ActionTrigger {
	if dto == nil {
		return nil
	}

	trigger := &model.ActionTrigger{
		Kind: safeString(dto.Kind),
	}

	if dto.Balance != nil {
		trigger.Balance = TriggerBalanceFromDTO(dto.Balance)
	}

	return trigger
}

// TriggerBalanceFromDTO converts an OpenAPI TriggerBalance to domain TriggerBalance.
func TriggerBalanceFromDTO(dto *openapi.TriggerBalance) *model.TriggerBalance {
	if dto == nil {
		return nil
	}

	balance := &model.TriggerBalance{}

	if dto.Target != nil {
		balance.Target = ActionTargetFromDTO(dto.Target)
	}
	if dto.Comparator != nil {
		balance.Comparator = ActionComparatorFromDTO(dto.Comparator)
	}
	if dto.Amount != nil {
		balance.Amount = ActionAmountFromDTO(dto.Amount)
	}

	return balance
}

// ActionTargetFromDTO converts an OpenAPI ActionTarget to domain ActionTarget.
func ActionTargetFromDTO(dto *openapi.ActionTarget) *model.ActionTarget {
	if dto == nil {
		return nil
	}

	target := &model.ActionTarget{
		Kind: safeString(dto.Kind),
	}

	if dto.Address != nil {
		target.Address = TargetAddressFromDTO(dto.Address)
	}
	if dto.Wallet != nil {
		target.Wallet = TargetWalletFromDTO(dto.Wallet)
	}

	return target
}

// TargetAddressFromDTO converts an OpenAPI TargetAddress to domain TargetAddress.
func TargetAddressFromDTO(dto *openapi.TargetAddress) *model.TargetAddress {
	if dto == nil {
		return nil
	}

	return &model.TargetAddress{
		Kind:      safeString(dto.Kind),
		AddressID: safeString(dto.AddressID),
	}
}

// TargetWalletFromDTO converts an OpenAPI TargetWallet to domain TargetWallet.
func TargetWalletFromDTO(dto *openapi.TargetWallet) *model.TargetWallet {
	if dto == nil {
		return nil
	}

	return &model.TargetWallet{
		Kind:     safeString(dto.Kind),
		WalletID: safeString(dto.WalletID),
	}
}

// ActionComparatorFromDTO converts an OpenAPI ActionComparator to domain ActionComparator.
func ActionComparatorFromDTO(dto *openapi.ActionComparator) *model.ActionComparator {
	if dto == nil {
		return nil
	}

	return &model.ActionComparator{
		Kind: safeString(dto.Kind),
	}
}

// ActionAmountFromDTO converts an OpenAPI ActionAmount to domain ActionAmount.
func ActionAmountFromDTO(dto *openapi.TgvalidatordActionAmount) *model.ActionAmount {
	if dto == nil {
		return nil
	}

	return &model.ActionAmount{
		Kind:         safeString(dto.Kind),
		CryptoAmount: safeString(dto.CryptoAmount),
	}
}

// ActionTasksFromDTO converts a slice of OpenAPI ActionTask to domain ActionTasks.
func ActionTasksFromDTO(dtos []openapi.ActionTask) []*model.ActionTask {
	if dtos == nil {
		return nil
	}
	tasks := make([]*model.ActionTask, len(dtos))
	for i := range dtos {
		tasks[i] = ActionTaskFromDTO(&dtos[i])
	}
	return tasks
}

// ActionTaskFromDTO converts an OpenAPI ActionTask to domain ActionTask.
func ActionTaskFromDTO(dto *openapi.ActionTask) *model.ActionTask {
	if dto == nil {
		return nil
	}

	task := &model.ActionTask{
		Kind: safeString(dto.Kind),
	}

	if dto.Transfer != nil {
		task.Transfer = TaskTransferFromDTO(dto.Transfer)
	}
	if dto.Notification != nil {
		task.Notification = TaskNotificationFromDTO(dto.Notification)
	}

	return task
}

// TaskTransferFromDTO converts an OpenAPI TaskTransfer to domain TaskTransfer.
func TaskTransferFromDTO(dto *openapi.TaskTransfer) *model.TaskTransfer {
	if dto == nil {
		return nil
	}

	transfer := &model.TaskTransfer{
		TopUp:       safeBool(dto.TopUp),
		UseAllFunds: safeBool(dto.UseAllFunds),
	}

	if dto.From != nil {
		transfer.From = ActionSourceFromDTO(dto.From)
	}
	if dto.To != nil {
		transfer.To = ActionDestinationFromDTO(dto.To)
	}
	if dto.Amount != nil {
		transfer.Amount = ActionAmountFromDTO(dto.Amount)
	}

	return transfer
}

// ActionSourceFromDTO converts an OpenAPI ActionSource to domain ActionSource.
func ActionSourceFromDTO(dto *openapi.TgvalidatordActionSource) *model.ActionSource {
	if dto == nil {
		return nil
	}

	return &model.ActionSource{
		Kind:      safeString(dto.Kind),
		AddressID: safeString(dto.AddressID),
		WalletID:  safeString(dto.WalletID),
	}
}

// ActionDestinationFromDTO converts an OpenAPI ActionDestination to domain ActionDestination.
func ActionDestinationFromDTO(dto *openapi.TgvalidatordActionDestination) *model.ActionDestination {
	if dto == nil {
		return nil
	}

	return &model.ActionDestination{
		Kind:                 safeString(dto.Kind),
		AddressID:            safeString(dto.AddressID),
		WhitelistedAddressID: safeString(dto.WhitelistedAddressID),
		WalletID:             safeString(dto.WalletID),
	}
}

// TaskNotificationFromDTO converts an OpenAPI TaskNotification to domain TaskNotification.
func TaskNotificationFromDTO(dto *openapi.TaskNotification) *model.TaskNotification {
	if dto == nil {
		return nil
	}

	return &model.TaskNotification{
		EmailAddresses:      dto.EmailAddresses,
		NotificationMessage: safeString(dto.NotificationMessage),
		NumberOfReminders:   safeString(dto.NumberOfReminders),
	}
}

// ActionAttributesFromDTO converts a slice of OpenAPI ActionAttribute to domain ActionAttributes.
func ActionAttributesFromDTO(dtos []openapi.TgvalidatordActionAttribute) []*model.ActionAttribute {
	if dtos == nil {
		return nil
	}
	attrs := make([]*model.ActionAttribute, len(dtos))
	for i := range dtos {
		attrs[i] = ActionAttributeFromDTO(&dtos[i])
	}
	return attrs
}

// ActionAttributeFromDTO converts an OpenAPI ActionAttribute to domain ActionAttribute.
func ActionAttributeFromDTO(dto *openapi.TgvalidatordActionAttribute) *model.ActionAttribute {
	if dto == nil {
		return nil
	}

	return &model.ActionAttribute{
		ID:          safeString(dto.Id),
		TenantID:    safeString(dto.TenantId),
		Key:         safeString(dto.Key),
		Value:       safeString(dto.Value),
		ContentType: safeString(dto.ContentType),
	}
}

// ActionTrailsFromDTO converts a slice of OpenAPI ActionEnvelopeTrail to domain ActionTrails.
func ActionTrailsFromDTO(dtos []openapi.TgvalidatordActionEnvelopeTrail) []*model.ActionTrail {
	if dtos == nil {
		return nil
	}
	trails := make([]*model.ActionTrail, len(dtos))
	for i := range dtos {
		trails[i] = ActionTrailFromDTO(&dtos[i])
	}
	return trails
}

// ActionTrailFromDTO converts an OpenAPI ActionEnvelopeTrail to domain ActionTrail.
func ActionTrailFromDTO(dto *openapi.TgvalidatordActionEnvelopeTrail) *model.ActionTrail {
	if dto == nil {
		return nil
	}

	trail := &model.ActionTrail{
		ID:           safeString(dto.Id),
		Action:       safeString(dto.Action),
		Comment:      safeString(dto.Comment),
		ActionStatus: safeString(dto.ActionStatus),
	}

	if dto.Date != nil {
		trail.Date = *dto.Date
	}

	return trail
}
