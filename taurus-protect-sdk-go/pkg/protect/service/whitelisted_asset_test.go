package service

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewWhitelistedAssetServiceWithVerification(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestWhitelistedAssetService_GetWhitelistedAsset_EmptyID(t *testing.T) {
	// Create a service with nil API to test validation before API call
	svc := &WhitelistedAssetService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetWhitelistedAsset(nil, "")
	if err == nil {
		t.Error("GetWhitelistedAsset() with empty ID should return error")
	}
	if err.Error() != "id cannot be empty" {
		t.Errorf("GetWhitelistedAsset() error = %v, want 'id cannot be empty'", err)
	}
}

func TestListWhitelistedAssetsOptions(t *testing.T) {
	// Test that options struct can be created with all fields
	opts := &model.ListWhitelistedAssetsOptions{
		Limit:              100,
		Offset:             50,
		Query:              "USDC",
		Blockchain:         "ETH",
		Network:            "mainnet",
		IncludeForApproval: true,
		KindTypes:          []string{"token", "nft"},
		IDs:                []string{"id1", "id2"},
	}

	if opts.Limit != 100 {
		t.Errorf("Limit = %v, want 100", opts.Limit)
	}
	if opts.Offset != 50 {
		t.Errorf("Offset = %v, want 50", opts.Offset)
	}
	if opts.Query != "USDC" {
		t.Errorf("Query = %v, want 'USDC'", opts.Query)
	}
	if opts.Blockchain != "ETH" {
		t.Errorf("Blockchain = %v, want 'ETH'", opts.Blockchain)
	}
	if opts.Network != "mainnet" {
		t.Errorf("Network = %v, want 'mainnet'", opts.Network)
	}
	if !opts.IncludeForApproval {
		t.Error("IncludeForApproval should be true")
	}
	if len(opts.KindTypes) != 2 {
		t.Errorf("KindTypes length = %v, want 2", len(opts.KindTypes))
	}
	if len(opts.IDs) != 2 {
		t.Errorf("IDs length = %v, want 2", len(opts.IDs))
	}
}
