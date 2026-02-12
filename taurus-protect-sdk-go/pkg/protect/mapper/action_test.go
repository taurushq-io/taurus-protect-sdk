package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestActionFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordActionEnvelope
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns action with zero values",
			dto:  &openapi.TgvalidatordActionEnvelope{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordActionEnvelope {
				id := "action-123"
				tenantId := "tenant-456"
				label := "Test Action"
				status := "ACTIVE"
				autoApprove := true
				creationDate := time.Now()
				updateDate := time.Now().Add(time.Hour)
				lastCheckedDate := time.Now().Add(2 * time.Hour)

				return &openapi.TgvalidatordActionEnvelope{
					Id:              &id,
					TenantId:        &tenantId,
					Label:           &label,
					Status:          &status,
					AutoApprove:     &autoApprove,
					CreationDate:    &creationDate,
					UpdateDate:      &updateDate,
					Lastcheckeddate: &lastCheckedDate,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ActionFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ActionFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ActionFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.TenantId != nil && got.TenantID != *tt.dto.TenantId {
				t.Errorf("TenantID = %v, want %v", got.TenantID, *tt.dto.TenantId)
			}
			if tt.dto.Label != nil && got.Label != *tt.dto.Label {
				t.Errorf("Label = %v, want %v", got.Label, *tt.dto.Label)
			}
			if tt.dto.Status != nil && got.Status != *tt.dto.Status {
				t.Errorf("Status = %v, want %v", got.Status, *tt.dto.Status)
			}
			if tt.dto.AutoApprove != nil && got.AutoApprove != *tt.dto.AutoApprove {
				t.Errorf("AutoApprove = %v, want %v", got.AutoApprove, *tt.dto.AutoApprove)
			}
			if tt.dto.CreationDate != nil && !got.CreationDate.Equal(*tt.dto.CreationDate) {
				t.Errorf("CreationDate = %v, want %v", got.CreationDate, *tt.dto.CreationDate)
			}
			if tt.dto.UpdateDate != nil && !got.UpdateDate.Equal(*tt.dto.UpdateDate) {
				t.Errorf("UpdateDate = %v, want %v", got.UpdateDate, *tt.dto.UpdateDate)
			}
			if tt.dto.Lastcheckeddate != nil && !got.LastCheckedDate.Equal(*tt.dto.Lastcheckeddate) {
				t.Errorf("LastCheckedDate = %v, want %v", got.LastCheckedDate, *tt.dto.Lastcheckeddate)
			}
		})
	}
}

func TestActionsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordActionEnvelope
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordActionEnvelope{},
			want: 0,
		},
		{
			name: "converts multiple actions",
			dtos: func() []openapi.TgvalidatordActionEnvelope {
				label1 := "Action 1"
				label2 := "Action 2"
				return []openapi.TgvalidatordActionEnvelope{
					{Label: &label1},
					{Label: &label2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ActionsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("ActionsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("ActionsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestActionDetailsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordAction
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns action details with nil fields",
			dto:  &openapi.TgvalidatordAction{},
		},
		{
			name: "DTO with trigger",
			dto: func() *openapi.TgvalidatordAction {
				kind := "BALANCE"
				return &openapi.TgvalidatordAction{
					Trigger: &openapi.ActionTrigger{
						Kind: &kind,
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ActionDetailsFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ActionDetailsFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ActionDetailsFromDTO() returned nil for non-nil input")
			}
		})
	}
}

func TestActionTriggerFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.ActionTrigger
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns trigger with zero values",
			dto:  &openapi.ActionTrigger{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.ActionTrigger {
				kind := "BALANCE"
				return &openapi.ActionTrigger{
					Kind:    &kind,
					Balance: &openapi.TriggerBalance{},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ActionTriggerFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ActionTriggerFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ActionTriggerFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Kind != nil && got.Kind != *tt.dto.Kind {
				t.Errorf("Kind = %v, want %v", got.Kind, *tt.dto.Kind)
			}
		})
	}
}

func TestTriggerBalanceFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TriggerBalance
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns balance with nil fields",
			dto:  &openapi.TriggerBalance{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TriggerBalance {
				targetKind := "ADDRESS"
				comparatorKind := "LESS_THAN"
				amountKind := "CRYPTO"
				cryptoAmount := "1000000"
				return &openapi.TriggerBalance{
					Target:     &openapi.ActionTarget{Kind: &targetKind},
					Comparator: &openapi.ActionComparator{Kind: &comparatorKind},
					Amount:     &openapi.TgvalidatordActionAmount{Kind: &amountKind, CryptoAmount: &cryptoAmount},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TriggerBalanceFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("TriggerBalanceFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("TriggerBalanceFromDTO() returned nil for non-nil input")
			}
		})
	}
}

func TestActionTargetFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.ActionTarget
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns target with zero values",
			dto:  &openapi.ActionTarget{},
		},
		{
			name: "complete DTO with address",
			dto: func() *openapi.ActionTarget {
				kind := "ADDRESS"
				addressKind := "SPECIFIC"
				addressID := "addr-123"
				return &openapi.ActionTarget{
					Kind: &kind,
					Address: &openapi.TargetAddress{
						Kind:      &addressKind,
						AddressID: &addressID,
					},
				}
			}(),
		},
		{
			name: "complete DTO with wallet",
			dto: func() *openapi.ActionTarget {
				kind := "WALLET"
				walletKind := "SPECIFIC"
				walletID := "wallet-123"
				return &openapi.ActionTarget{
					Kind: &kind,
					Wallet: &openapi.TargetWallet{
						Kind:     &walletKind,
						WalletID: &walletID,
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ActionTargetFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ActionTargetFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ActionTargetFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Kind != nil && got.Kind != *tt.dto.Kind {
				t.Errorf("Kind = %v, want %v", got.Kind, *tt.dto.Kind)
			}
		})
	}
}

func TestTargetAddressFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TargetAddress
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns address with zero values",
			dto:  &openapi.TargetAddress{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TargetAddress {
				kind := "SPECIFIC"
				addressID := "addr-123"
				return &openapi.TargetAddress{
					Kind:      &kind,
					AddressID: &addressID,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TargetAddressFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("TargetAddressFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("TargetAddressFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Kind != nil && got.Kind != *tt.dto.Kind {
				t.Errorf("Kind = %v, want %v", got.Kind, *tt.dto.Kind)
			}
			if tt.dto.AddressID != nil && got.AddressID != *tt.dto.AddressID {
				t.Errorf("AddressID = %v, want %v", got.AddressID, *tt.dto.AddressID)
			}
		})
	}
}

func TestTargetWalletFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TargetWallet
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns wallet with zero values",
			dto:  &openapi.TargetWallet{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TargetWallet {
				kind := "SPECIFIC"
				walletID := "wallet-123"
				return &openapi.TargetWallet{
					Kind:     &kind,
					WalletID: &walletID,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TargetWalletFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("TargetWalletFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("TargetWalletFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Kind != nil && got.Kind != *tt.dto.Kind {
				t.Errorf("Kind = %v, want %v", got.Kind, *tt.dto.Kind)
			}
			if tt.dto.WalletID != nil && got.WalletID != *tt.dto.WalletID {
				t.Errorf("WalletID = %v, want %v", got.WalletID, *tt.dto.WalletID)
			}
		})
	}
}

func TestActionComparatorFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.ActionComparator
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns comparator with zero values",
			dto:  &openapi.ActionComparator{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.ActionComparator {
				kind := "LESS_THAN"
				return &openapi.ActionComparator{
					Kind: &kind,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ActionComparatorFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ActionComparatorFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ActionComparatorFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Kind != nil && got.Kind != *tt.dto.Kind {
				t.Errorf("Kind = %v, want %v", got.Kind, *tt.dto.Kind)
			}
		})
	}
}

func TestActionAmountFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordActionAmount
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns amount with zero values",
			dto:  &openapi.TgvalidatordActionAmount{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordActionAmount {
				kind := "CRYPTO"
				cryptoAmount := "1000000"
				return &openapi.TgvalidatordActionAmount{
					Kind:         &kind,
					CryptoAmount: &cryptoAmount,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ActionAmountFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ActionAmountFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ActionAmountFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Kind != nil && got.Kind != *tt.dto.Kind {
				t.Errorf("Kind = %v, want %v", got.Kind, *tt.dto.Kind)
			}
			if tt.dto.CryptoAmount != nil && got.CryptoAmount != *tt.dto.CryptoAmount {
				t.Errorf("CryptoAmount = %v, want %v", got.CryptoAmount, *tt.dto.CryptoAmount)
			}
		})
	}
}

func TestActionTasksFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.ActionTask
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1,
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.ActionTask{},
			want: 0,
		},
		{
			name: "converts multiple tasks",
			dtos: func() []openapi.ActionTask {
				kind1 := "TRANSFER"
				kind2 := "NOTIFICATION"
				return []openapi.ActionTask{
					{Kind: &kind1},
					{Kind: &kind2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ActionTasksFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("ActionTasksFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("ActionTasksFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestActionTaskFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.ActionTask
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns task with zero values",
			dto:  &openapi.ActionTask{},
		},
		{
			name: "complete DTO with transfer",
			dto: func() *openapi.ActionTask {
				kind := "TRANSFER"
				topUp := true
				return &openapi.ActionTask{
					Kind: &kind,
					Transfer: &openapi.TaskTransfer{
						TopUp: &topUp,
					},
				}
			}(),
		},
		{
			name: "complete DTO with notification",
			dto: func() *openapi.ActionTask {
				kind := "NOTIFICATION"
				message := "Test notification"
				return &openapi.ActionTask{
					Kind: &kind,
					Notification: &openapi.TaskNotification{
						NotificationMessage: &message,
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ActionTaskFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ActionTaskFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ActionTaskFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Kind != nil && got.Kind != *tt.dto.Kind {
				t.Errorf("Kind = %v, want %v", got.Kind, *tt.dto.Kind)
			}
		})
	}
}

func TestTaskTransferFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TaskTransfer
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns transfer with zero values",
			dto:  &openapi.TaskTransfer{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TaskTransfer {
				topUp := true
				useAllFunds := false
				sourceKind := "ADDRESS"
				destKind := "WALLET"
				return &openapi.TaskTransfer{
					TopUp:       &topUp,
					UseAllFunds: &useAllFunds,
					From:        &openapi.TgvalidatordActionSource{Kind: &sourceKind},
					To:          &openapi.TgvalidatordActionDestination{Kind: &destKind},
					Amount:      &openapi.TgvalidatordActionAmount{},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TaskTransferFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("TaskTransferFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("TaskTransferFromDTO() returned nil for non-nil input")
			}
			if tt.dto.TopUp != nil && got.TopUp != *tt.dto.TopUp {
				t.Errorf("TopUp = %v, want %v", got.TopUp, *tt.dto.TopUp)
			}
			if tt.dto.UseAllFunds != nil && got.UseAllFunds != *tt.dto.UseAllFunds {
				t.Errorf("UseAllFunds = %v, want %v", got.UseAllFunds, *tt.dto.UseAllFunds)
			}
		})
	}
}

func TestActionSourceFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordActionSource
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns source with zero values",
			dto:  &openapi.TgvalidatordActionSource{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordActionSource {
				kind := "ADDRESS"
				addressID := "addr-123"
				walletID := "wallet-456"
				return &openapi.TgvalidatordActionSource{
					Kind:      &kind,
					AddressID: &addressID,
					WalletID:  &walletID,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ActionSourceFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ActionSourceFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ActionSourceFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Kind != nil && got.Kind != *tt.dto.Kind {
				t.Errorf("Kind = %v, want %v", got.Kind, *tt.dto.Kind)
			}
			if tt.dto.AddressID != nil && got.AddressID != *tt.dto.AddressID {
				t.Errorf("AddressID = %v, want %v", got.AddressID, *tt.dto.AddressID)
			}
			if tt.dto.WalletID != nil && got.WalletID != *tt.dto.WalletID {
				t.Errorf("WalletID = %v, want %v", got.WalletID, *tt.dto.WalletID)
			}
		})
	}
}

func TestActionDestinationFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordActionDestination
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns destination with zero values",
			dto:  &openapi.TgvalidatordActionDestination{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordActionDestination {
				kind := "WHITELISTED_ADDRESS"
				addressID := "addr-123"
				whitelistedID := "whitelist-456"
				walletID := "wallet-789"
				return &openapi.TgvalidatordActionDestination{
					Kind:                 &kind,
					AddressID:            &addressID,
					WhitelistedAddressID: &whitelistedID,
					WalletID:             &walletID,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ActionDestinationFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ActionDestinationFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ActionDestinationFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Kind != nil && got.Kind != *tt.dto.Kind {
				t.Errorf("Kind = %v, want %v", got.Kind, *tt.dto.Kind)
			}
			if tt.dto.AddressID != nil && got.AddressID != *tt.dto.AddressID {
				t.Errorf("AddressID = %v, want %v", got.AddressID, *tt.dto.AddressID)
			}
			if tt.dto.WhitelistedAddressID != nil && got.WhitelistedAddressID != *tt.dto.WhitelistedAddressID {
				t.Errorf("WhitelistedAddressID = %v, want %v", got.WhitelistedAddressID, *tt.dto.WhitelistedAddressID)
			}
			if tt.dto.WalletID != nil && got.WalletID != *tt.dto.WalletID {
				t.Errorf("WalletID = %v, want %v", got.WalletID, *tt.dto.WalletID)
			}
		})
	}
}

func TestTaskNotificationFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TaskNotification
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns notification with zero values",
			dto:  &openapi.TaskNotification{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TaskNotification {
				message := "Test notification message"
				reminders := "3"
				return &openapi.TaskNotification{
					EmailAddresses:      []string{"test@example.com", "user@example.com"},
					NotificationMessage: &message,
					NumberOfReminders:   &reminders,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TaskNotificationFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("TaskNotificationFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("TaskNotificationFromDTO() returned nil for non-nil input")
			}
			if tt.dto.NotificationMessage != nil && got.NotificationMessage != *tt.dto.NotificationMessage {
				t.Errorf("NotificationMessage = %v, want %v", got.NotificationMessage, *tt.dto.NotificationMessage)
			}
			if tt.dto.NumberOfReminders != nil && got.NumberOfReminders != *tt.dto.NumberOfReminders {
				t.Errorf("NumberOfReminders = %v, want %v", got.NumberOfReminders, *tt.dto.NumberOfReminders)
			}
			if len(tt.dto.EmailAddresses) != len(got.EmailAddresses) {
				t.Errorf("EmailAddresses length = %v, want %v", len(got.EmailAddresses), len(tt.dto.EmailAddresses))
			}
		})
	}
}

func TestActionAttributesFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordActionAttribute
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1,
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordActionAttribute{},
			want: 0,
		},
		{
			name: "converts multiple attributes",
			dtos: func() []openapi.TgvalidatordActionAttribute {
				key1 := "key1"
				key2 := "key2"
				return []openapi.TgvalidatordActionAttribute{
					{Key: &key1},
					{Key: &key2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ActionAttributesFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("ActionAttributesFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("ActionAttributesFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestActionAttributeFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordActionAttribute
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns attribute with zero values",
			dto:  &openapi.TgvalidatordActionAttribute{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordActionAttribute {
				id := "attr-123"
				tenantId := "tenant-456"
				key := "custom_key"
				value := "custom_value"
				contentType := "text/plain"
				return &openapi.TgvalidatordActionAttribute{
					Id:          &id,
					TenantId:    &tenantId,
					Key:         &key,
					Value:       &value,
					ContentType: &contentType,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ActionAttributeFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ActionAttributeFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ActionAttributeFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.TenantId != nil && got.TenantID != *tt.dto.TenantId {
				t.Errorf("TenantID = %v, want %v", got.TenantID, *tt.dto.TenantId)
			}
			if tt.dto.Key != nil && got.Key != *tt.dto.Key {
				t.Errorf("Key = %v, want %v", got.Key, *tt.dto.Key)
			}
			if tt.dto.Value != nil && got.Value != *tt.dto.Value {
				t.Errorf("Value = %v, want %v", got.Value, *tt.dto.Value)
			}
			if tt.dto.ContentType != nil && got.ContentType != *tt.dto.ContentType {
				t.Errorf("ContentType = %v, want %v", got.ContentType, *tt.dto.ContentType)
			}
		})
	}
}

func TestActionTrailsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordActionEnvelopeTrail
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1,
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordActionEnvelopeTrail{},
			want: 0,
		},
		{
			name: "converts multiple trails",
			dtos: func() []openapi.TgvalidatordActionEnvelopeTrail {
				action1 := "CREATE"
				action2 := "APPROVE"
				return []openapi.TgvalidatordActionEnvelopeTrail{
					{Action: &action1},
					{Action: &action2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ActionTrailsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("ActionTrailsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("ActionTrailsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestActionTrailFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordActionEnvelopeTrail
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns trail with zero values",
			dto:  &openapi.TgvalidatordActionEnvelopeTrail{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordActionEnvelopeTrail {
				id := "trail-123"
				action := "APPROVE"
				comment := "Approved by admin"
				date := time.Now()
				actionStatus := "APPROVED"
				return &openapi.TgvalidatordActionEnvelopeTrail{
					Id:           &id,
					Action:       &action,
					Comment:      &comment,
					Date:         &date,
					ActionStatus: &actionStatus,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ActionTrailFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ActionTrailFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ActionTrailFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Action != nil && got.Action != *tt.dto.Action {
				t.Errorf("Action = %v, want %v", got.Action, *tt.dto.Action)
			}
			if tt.dto.Comment != nil && got.Comment != *tt.dto.Comment {
				t.Errorf("Comment = %v, want %v", got.Comment, *tt.dto.Comment)
			}
			if tt.dto.Date != nil && !got.Date.Equal(*tt.dto.Date) {
				t.Errorf("Date = %v, want %v", got.Date, *tt.dto.Date)
			}
			if tt.dto.ActionStatus != nil && got.ActionStatus != *tt.dto.ActionStatus {
				t.Errorf("ActionStatus = %v, want %v", got.ActionStatus, *tt.dto.ActionStatus)
			}
		})
	}
}

func TestActionTrailFromDTO_NilDate(t *testing.T) {
	action := "CREATE"
	dto := &openapi.TgvalidatordActionEnvelopeTrail{
		Action: &action,
		Date:   nil,
	}

	got := ActionTrailFromDTO(dto)
	if got == nil {
		t.Fatal("ActionTrailFromDTO() returned nil for non-nil input")
	}
	// When date is nil, it should be the zero time value
	if !got.Date.IsZero() {
		t.Errorf("Date should be zero time when nil, got %v", got.Date)
	}
	if got.Action != "CREATE" {
		t.Errorf("Action = %v, want CREATE", got.Action)
	}
}

func TestActionFromDTO_WithCompleteAction(t *testing.T) {
	// Test a complete action with all nested structures
	id := "action-123"
	triggerKind := "BALANCE"
	taskKind := "TRANSFER"
	topUp := true

	dto := &openapi.TgvalidatordActionEnvelope{
		Id: &id,
		Action: &openapi.TgvalidatordAction{
			Trigger: &openapi.ActionTrigger{
				Kind: &triggerKind,
			},
			Tasks: []openapi.ActionTask{
				{
					Kind: &taskKind,
					Transfer: &openapi.TaskTransfer{
						TopUp: &topUp,
					},
				},
			},
		},
	}

	got := ActionFromDTO(dto)
	if got == nil {
		t.Fatal("ActionFromDTO() returned nil for non-nil input")
	}
	if got.ID != "action-123" {
		t.Errorf("ID = %v, want action-123", got.ID)
	}
	if got.ActionDetails == nil {
		t.Fatal("ActionDetails should not be nil")
	}
	if got.ActionDetails.Trigger == nil {
		t.Fatal("Trigger should not be nil")
	}
	if got.ActionDetails.Trigger.Kind != "BALANCE" {
		t.Errorf("Trigger.Kind = %v, want BALANCE", got.ActionDetails.Trigger.Kind)
	}
	if len(got.ActionDetails.Tasks) != 1 {
		t.Fatalf("Tasks length = %v, want 1", len(got.ActionDetails.Tasks))
	}
	if got.ActionDetails.Tasks[0].Kind != "TRANSFER" {
		t.Errorf("Tasks[0].Kind = %v, want TRANSFER", got.ActionDetails.Tasks[0].Kind)
	}
	if got.ActionDetails.Tasks[0].Transfer == nil {
		t.Fatal("Transfer should not be nil")
	}
	if !got.ActionDetails.Tasks[0].Transfer.TopUp {
		t.Errorf("Transfer.TopUp = %v, want true", got.ActionDetails.Tasks[0].Transfer.TopUp)
	}
}
