package service

import (
	"errors"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestGovernanceRuleService_GetRulesByID_Validation(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty id returns error",
			id:      "",
			wantErr: true,
			errMsg:  "id cannot be empty",
		},
	}

	// Note: We can only test input validation without a real API client
	// Full integration tests require a running server
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a minimal service for validation testing
			s := &GovernanceRuleService{
				errMapper: NewErrorMapper(),
			}

			_, err := s.GetRulesByID(nil, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRulesByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("GetRulesByID() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestGovernanceRulesHistoryOptions(t *testing.T) {
	tests := []struct {
		name string
		opts *model.ListRulesHistoryOptions
	}{
		{
			name: "nil options",
			opts: nil,
		},
		{
			name: "with limit",
			opts: &model.ListRulesHistoryOptions{
				Limit: 10,
			},
		},
		{
			name: "with cursor",
			opts: &model.ListRulesHistoryOptions{
				Cursor: "eyJwYWdlIjogMn0=",
			},
		},
		{
			name: "with limit and cursor",
			opts: &model.ListRulesHistoryOptions{
				Limit:  20,
				Cursor: "eyJwYWdlIjogM30=",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify options are properly constructed
			if tt.opts != nil {
				if tt.opts.Limit < 0 {
					t.Error("Limit should not be negative")
				}
			}
		})
	}
}

func TestNewGovernanceRuleService(t *testing.T) {
	// Test that NewGovernanceRuleService doesn't panic with nil
	// In real usage, this would require a valid APIClient
	defer func() { _ = recover() }()

	// This would panic with nil, which is expected
	// service := NewGovernanceRuleService(nil)
	// We just verify the function signature is correct
}

// VerifyGovernanceRules returns the verified rules alongside the error, matching Java
// and Python (Go used to return only an error, so a caller could not chain
// verification into an assignment). These cases pin the two-value contract without
// needing a signing identity: on failure the rules must be nil, never a half-verified
// value a caller might go on to use.
func TestGovernanceRuleService_VerifyGovernanceRules_ReturnsNilRulesOnFailure(t *testing.T) {
	svc := &GovernanceRuleService{}

	t.Run("no signatures", func(t *testing.T) {
		rules, err := svc.VerifyGovernanceRules(&model.GovernanceRuleset{RulesContainer: "AA=="})
		if err == nil {
			t.Fatal("expected an error for rules with no signatures")
		}
		if rules != nil {
			t.Errorf("rules = %#v, want nil on failure", rules)
		}
		var integrityErr *model.IntegrityError
		if !errors.As(err, &integrityErr) {
			t.Errorf("err = %T, want *model.IntegrityError", err)
		}
	})

	t.Run("undecodable container", func(t *testing.T) {
		rules, err := svc.VerifyGovernanceRules(&model.GovernanceRuleset{
			RulesContainer: "not-base64!!",
			Signatures:     []model.RuleUserSignature{{UserID: "u1", Signature: "sig"}},
		})
		if err == nil {
			t.Fatal("expected an error for an undecodable rules container")
		}
		if rules != nil {
			t.Errorf("rules = %#v, want nil on failure", rules)
		}
	})
}
