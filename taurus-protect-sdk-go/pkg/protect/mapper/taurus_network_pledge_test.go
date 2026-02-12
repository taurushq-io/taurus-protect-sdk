package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestPledgeFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTnPledge
		want bool // true if result should be non-nil
	}{
		{
			name: "nil input",
			dto:  nil,
			want: false,
		},
		{
			name: "empty pledge",
			dto:  &openapi.TgvalidatordTnPledge{},
			want: true,
		},
		{
			name: "pledge with all fields",
			dto: func() *openapi.TgvalidatordTnPledge {
				id := "pledge-123"
				ownerID := "owner-456"
				targetID := "target-789"
				sharedAddrID := "addr-111"
				currencyID := "ETH"
				amount := "1000000000000000000"
				status := "ACCEPTED_BY_TARGET"
				pledgeType := "COLLATERAL"
				direction := "OUTGOING"
				extRef := "ext-ref-001"
				note := "Test pledge"
				now := time.Now()
				return &openapi.TgvalidatordTnPledge{
					Id:                  &id,
					OwnerParticipantID:  &ownerID,
					TargetParticipantID: &targetID,
					SharedAddressID:     &sharedAddrID,
					CurrencyID:          &currencyID,
					Amount:              &amount,
					Status:              &status,
					PledgeType:          &pledgeType,
					Direction:           &direction,
					ExternalReferenceId: &extRef,
					ReconciliationNote:  &note,
					CreatedAt:           &now,
				}
			}(),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PledgeFromDTO(tt.dto)
			if (result != nil) != tt.want {
				t.Errorf("PledgeFromDTO() = %v, want non-nil: %v", result, tt.want)
			}
			if result != nil && tt.dto != nil {
				if tt.dto.Id != nil && result.ID != *tt.dto.Id {
					t.Errorf("PledgeFromDTO().ID = %v, want %v", result.ID, *tt.dto.Id)
				}
				if tt.dto.OwnerParticipantID != nil && result.OwnerParticipantID != *tt.dto.OwnerParticipantID {
					t.Errorf("PledgeFromDTO().OwnerParticipantID = %v, want %v", result.OwnerParticipantID, *tt.dto.OwnerParticipantID)
				}
				if tt.dto.Amount != nil && result.Amount != *tt.dto.Amount {
					t.Errorf("PledgeFromDTO().Amount = %v, want %v", result.Amount, *tt.dto.Amount)
				}
			}
		})
	}
}

func TestPledgesFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordTnPledge
		want int
	}{
		{
			name: "nil input",
			dtos: nil,
			want: 0,
		},
		{
			name: "empty slice",
			dtos: []openapi.TgvalidatordTnPledge{},
			want: 0,
		},
		{
			name: "multiple pledges",
			dtos: []openapi.TgvalidatordTnPledge{
				{},
				{},
				{},
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PledgesFromDTO(tt.dtos)
			if tt.dtos == nil {
				if result != nil {
					t.Errorf("PledgesFromDTO() = %v, want nil", result)
				}
			} else if len(result) != tt.want {
				t.Errorf("PledgesFromDTO() len = %v, want %v", len(result), tt.want)
			}
		})
	}
}

func TestPledgeDurationSetupFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TnPledgePledgeDurationSetup
		want bool
	}{
		{
			name: "nil input",
			dto:  nil,
			want: false,
		},
		{
			name: "empty setup",
			dto:  &openapi.TnPledgePledgeDurationSetup{},
			want: true,
		},
		{
			name: "with all fields",
			dto: func() *openapi.TnPledgePledgeDurationSetup {
				minDur := "2592000s"
				noticeDur := "172800s"
				endMin := time.Now().Add(30 * 24 * time.Hour)
				endNotice := time.Now().Add(32 * 24 * time.Hour)
				return &openapi.TnPledgePledgeDurationSetup{
					MinimumDuration:          &minDur,
					NoticePeriodDuration:     &noticeDur,
					EndOfMinimumDurationDate: &endMin,
					EndOfNoticePeriodDate:    &endNotice,
				}
			}(),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PledgeDurationSetupFromDTO(tt.dto)
			if (result != nil) != tt.want {
				t.Errorf("PledgeDurationSetupFromDTO() = %v, want non-nil: %v", result, tt.want)
			}
			if result != nil && tt.dto != nil {
				if tt.dto.MinimumDuration != nil && result.MinimumDuration != *tt.dto.MinimumDuration {
					t.Errorf("PledgeDurationSetupFromDTO().MinimumDuration = %v, want %v", result.MinimumDuration, *tt.dto.MinimumDuration)
				}
			}
		})
	}
}

func TestPledgeAttributeFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TnPledgePledgeAttribute
		want bool
	}{
		{
			name: "nil input",
			dto:  nil,
			want: false,
		},
		{
			name: "empty attribute",
			dto:  &openapi.TnPledgePledgeAttribute{},
			want: true,
		},
		{
			name: "with all fields",
			dto: func() *openapi.TnPledgePledgeAttribute {
				id := "attr-123"
				key := "customKey"
				value := "customValue"
				owner := "owner-id"
				attrType := "STRING"
				shared := true
				return &openapi.TnPledgePledgeAttribute{
					Id:                    &id,
					Key:                   &key,
					Value:                 &value,
					Owner:                 &owner,
					Type:                  &attrType,
					IsTaurusNetworkShared: &shared,
				}
			}(),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PledgeAttributeFromDTO(tt.dto)
			if (result != nil) != tt.want {
				t.Errorf("PledgeAttributeFromDTO() = %v, want non-nil: %v", result, tt.want)
			}
			if result != nil && tt.dto != nil {
				if tt.dto.Key != nil && result.Key != *tt.dto.Key {
					t.Errorf("PledgeAttributeFromDTO().Key = %v, want %v", result.Key, *tt.dto.Key)
				}
				if tt.dto.IsTaurusNetworkShared != nil && result.IsTaurusNetworkShared != *tt.dto.IsTaurusNetworkShared {
					t.Errorf("PledgeAttributeFromDTO().IsTaurusNetworkShared = %v, want %v", result.IsTaurusNetworkShared, *tt.dto.IsTaurusNetworkShared)
				}
			}
		})
	}
}

func TestPledgeTrailFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTnPledgeTrail
		want bool
	}{
		{
			name: "nil input",
			dto:  nil,
			want: false,
		},
		{
			name: "empty trail",
			dto:  &openapi.TgvalidatordTnPledgeTrail{},
			want: true,
		},
		{
			name: "with all fields",
			dto: func() *openapi.TgvalidatordTnPledgeTrail {
				id := "trail-123"
				pledgeID := "pledge-456"
				action := "CREATED"
				comment := "Test trail"
				now := time.Now()
				return &openapi.TgvalidatordTnPledgeTrail{
					Id:        &id,
					PledgeID:  &pledgeID,
					Action:    &action,
					Comment:   &comment,
					CreatedAt: &now,
				}
			}(),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PledgeTrailFromDTO(tt.dto)
			if (result != nil) != tt.want {
				t.Errorf("PledgeTrailFromDTO() = %v, want non-nil: %v", result, tt.want)
			}
		})
	}
}

func TestPledgeWithdrawalFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTnPledgeWithdrawal
		want bool
	}{
		{
			name: "nil input",
			dto:  nil,
			want: false,
		},
		{
			name: "empty withdrawal",
			dto:  &openapi.TgvalidatordTnPledgeWithdrawal{},
			want: true,
		},
		{
			name: "with all fields",
			dto: func() *openapi.TgvalidatordTnPledgeWithdrawal {
				id := "withdrawal-123"
				pledgeID := "pledge-456"
				amount := "500000000000000000"
				status := "PENDING"
				now := time.Now()
				return &openapi.TgvalidatordTnPledgeWithdrawal{
					Id:        &id,
					PledgeID:  &pledgeID,
					Amount:    &amount,
					Status:    &status,
					CreatedAt: &now,
				}
			}(),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PledgeWithdrawalFromDTO(tt.dto)
			if (result != nil) != tt.want {
				t.Errorf("PledgeWithdrawalFromDTO() = %v, want non-nil: %v", result, tt.want)
			}
			if result != nil && tt.dto != nil {
				if tt.dto.Amount != nil && result.Amount != *tt.dto.Amount {
					t.Errorf("PledgeWithdrawalFromDTO().Amount = %v, want %v", result.Amount, *tt.dto.Amount)
				}
			}
		})
	}
}

func TestPledgeActionFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTnPledgeAction
		want bool
	}{
		{
			name: "nil input",
			dto:  nil,
			want: false,
		},
		{
			name: "empty action",
			dto:  &openapi.TgvalidatordTnPledgeAction{},
			want: true,
		},
		{
			name: "with basic fields",
			dto: func() *openapi.TgvalidatordTnPledgeAction {
				id := "action-123"
				pledgeID := "pledge-456"
				actionType := "CREATE"
				status := "PENDING_APPROVAL"
				now := time.Now()
				return &openapi.TgvalidatordTnPledgeAction{
					Id:         &id,
					PledgeID:   &pledgeID,
					ActionType: &actionType,
					Status:     &status,
					CreatedAt:  &now,
				}
			}(),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PledgeActionFromDTO(tt.dto)
			if (result != nil) != tt.want {
				t.Errorf("PledgeActionFromDTO() = %v, want non-nil: %v", result, tt.want)
			}
			if result != nil && tt.dto != nil {
				if tt.dto.Id != nil && result.ID != *tt.dto.Id {
					t.Errorf("PledgeActionFromDTO().ID = %v, want %v", result.ID, *tt.dto.Id)
				}
				if tt.dto.ActionType != nil && result.ActionType != *tt.dto.ActionType {
					t.Errorf("PledgeActionFromDTO().ActionType = %v, want %v", result.ActionType, *tt.dto.ActionType)
				}
			}
		})
	}
}

func TestPledgeActionTrailFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTnPledgeActionTrail
		want bool
	}{
		{
			name: "nil input",
			dto:  nil,
			want: false,
		},
		{
			name: "empty trail",
			dto:  &openapi.TgvalidatordTnPledgeActionTrail{},
			want: true,
		},
		{
			name: "with all fields",
			dto: func() *openapi.TgvalidatordTnPledgeActionTrail {
				id := "trail-123"
				actionID := "action-456"
				userID := "user-789"
				extUserID := "ext-user-001"
				action := "APPROVED"
				comment := "Approved by admin"
				now := time.Now()
				return &openapi.TgvalidatordTnPledgeActionTrail{
					Id:             &id,
					PledgeActionID: &actionID,
					UserID:         &userID,
					ExternalUserID: &extUserID,
					Action:         &action,
					Comment:        &comment,
					CreatedAt:      &now,
				}
			}(),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PledgeActionTrailFromDTO(tt.dto)
			if (result != nil) != tt.want {
				t.Errorf("PledgeActionTrailFromDTO() = %v, want non-nil: %v", result, tt.want)
			}
			if result != nil && tt.dto != nil {
				if tt.dto.UserID != nil && result.UserID != *tt.dto.UserID {
					t.Errorf("PledgeActionTrailFromDTO().UserID = %v, want %v", result.UserID, *tt.dto.UserID)
				}
			}
		})
	}
}

func TestPledgeWithdrawalTrailFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTnPledgeWithdrawalTrail
		want bool
	}{
		{
			name: "nil input",
			dto:  nil,
			want: false,
		},
		{
			name: "empty trail",
			dto:  &openapi.TgvalidatordTnPledgeWithdrawalTrail{},
			want: true,
		},
		{
			name: "with all fields",
			dto: func() *openapi.TgvalidatordTnPledgeWithdrawalTrail {
				id := "trail-123"
				withdrawalID := "withdrawal-456"
				participantID := "participant-789"
				action := "INITIATED"
				comment := "Withdrawal initiated"
				now := time.Now()
				return &openapi.TgvalidatordTnPledgeWithdrawalTrail{
					Id:                 &id,
					PledgeWithdrawalID: &withdrawalID,
					ParticipantID:      &participantID,
					Action:             &action,
					Comment:            &comment,
					CreatedAt:          &now,
				}
			}(),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PledgeWithdrawalTrailFromDTO(tt.dto)
			if (result != nil) != tt.want {
				t.Errorf("PledgeWithdrawalTrailFromDTO() = %v, want non-nil: %v", result, tt.want)
			}
			if result != nil && tt.dto != nil {
				if tt.dto.ParticipantID != nil && result.ParticipantID != *tt.dto.ParticipantID {
					t.Errorf("PledgeWithdrawalTrailFromDTO().ParticipantID = %v, want %v", result.ParticipantID, *tt.dto.ParticipantID)
				}
			}
		})
	}
}
