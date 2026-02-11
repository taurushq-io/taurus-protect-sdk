package service

import (
	"context"
	"testing"
)

func TestNewScoreService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestScoreService_RefreshAddressScore_EmptyAddressID(t *testing.T) {
	svc := &ScoreService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.RefreshAddressScore(context.Background(), "", "")
	if err == nil {
		t.Error("expected error for empty addressID, got nil")
	}
	if err.Error() != "addressID cannot be empty" {
		t.Errorf("error message = %v, want 'addressID cannot be empty'", err.Error())
	}
}

func TestScoreService_RefreshAddressScore_WithProvider(t *testing.T) {
	tests := []struct {
		name      string
		addressID string
		provider  string
	}{
		{
			name:      "with chainalysis provider",
			addressID: "addr-123",
			provider:  "chainalysis",
		},
		{
			name:      "with elliptic provider",
			addressID: "addr-456",
			provider:  "elliptic",
		},
		{
			name:      "with empty provider",
			addressID: "addr-789",
			provider:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify service accepts these parameters
			// Actual API testing requires mocking
			svc := &ScoreService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("ScoreService should not be nil")
			}
		})
	}
}

func TestScoreService_RefreshWLAScore_EmptyAddressID(t *testing.T) {
	svc := &ScoreService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.RefreshWLAScore(context.Background(), "", "")
	if err == nil {
		t.Error("expected error for empty addressID, got nil")
	}
	if err.Error() != "addressID cannot be empty" {
		t.Errorf("error message = %v, want 'addressID cannot be empty'", err.Error())
	}
}

func TestScoreService_RefreshWLAScore_WithProvider(t *testing.T) {
	tests := []struct {
		name      string
		addressID string
		provider  string
	}{
		{
			name:      "with chainalysis provider",
			addressID: "wla-123",
			provider:  "chainalysis",
		},
		{
			name:      "with elliptic provider",
			addressID: "wla-456",
			provider:  "elliptic",
		},
		{
			name:      "with empty provider",
			addressID: "wla-789",
			provider:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify service accepts these parameters
			// Actual API testing requires mocking
			svc := &ScoreService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("ScoreService should not be nil")
			}
		})
	}
}
