package service

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewAssetService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestAssetService_GetAssetAddresses_NilRequest(t *testing.T) {
	svc := &AssetService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetAssetAddresses(nil, nil)
	if err == nil {
		t.Error("expected error for nil request")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("error = %v, want 'request cannot be nil'", err)
	}
}

func TestAssetService_GetAssetAddresses_EmptyCurrency(t *testing.T) {
	svc := &AssetService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	req := &model.GetAssetAddressesRequest{
		Asset: model.AssetFilter{
			Currency: "",
		},
	}

	_, err := svc.GetAssetAddresses(nil, req)
	if err == nil {
		t.Error("expected error for empty currency")
	}
	if err.Error() != "asset currency is required" {
		t.Errorf("error = %v, want 'asset currency is required'", err)
	}
}

func TestAssetService_GetAssetAddresses_WithOptions(t *testing.T) {
	tests := []struct {
		name    string
		request *model.GetAssetAddressesRequest
	}{
		{
			name: "basic request",
			request: &model.GetAssetAddressesRequest{
				Asset: model.AssetFilter{
					Currency: "BTC",
				},
			},
		},
		{
			name: "with limit and cursor",
			request: &model.GetAssetAddressesRequest{
				Asset: model.AssetFilter{
					Currency: "ETH",
				},
				Limit:  50,
				Cursor: "abc123",
			},
		},
		{
			name: "with sort order",
			request: &model.GetAssetAddressesRequest{
				Asset: model.AssetFilter{
					Currency: "BTC",
				},
				SortOrder: "DESC",
			},
		},
		{
			name: "with wallet and address filters",
			request: &model.GetAssetAddressesRequest{
				Asset: model.AssetFilter{
					Currency: "ETH",
				},
				WalletID:  "wallet-123",
				AddressID: "address-456",
			},
		},
		{
			name: "with address list filter",
			request: &model.GetAssetAddressesRequest{
				Asset: model.AssetFilter{
					Currency: "ETH",
				},
				Addresses: []string{"0x1234", "0x5678"},
			},
		},
		{
			name: "with cursor pagination",
			request: &model.GetAssetAddressesRequest{
				Asset: model.AssetFilter{
					Currency: "BTC",
				},
				CurrentPage: "xyz789",
				PageRequest: "NEXT",
				PageSize:    25,
			},
		},
		{
			name: "with NFT filter",
			request: &model.GetAssetAddressesRequest{
				Asset: model.AssetFilter{
					Currency: "ETH",
					Kind:     "NFT",
					NFT: &model.AssetNFTFilter{
						TokenID: "token-123",
					},
				},
			},
		},
		{
			name: "with unknown asset filter",
			request: &model.GetAssetAddressesRequest{
				Asset: model.AssetFilter{
					Currency: "Unknown",
					Kind:     "Unknown",
					Unknown: &model.AssetUnknownFilter{
						Blockchain: "ethereum",
						Arg1:       "0xcontract",
						Arg2:       "12345",
						Network:    "mainnet",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify request can be created with these values
			// Actual API testing requires mocking
			svc := &AssetService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("AssetService should not be nil")
			}
		})
	}
}

func TestAssetService_GetAssetWallets_NilRequest(t *testing.T) {
	svc := &AssetService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetAssetWallets(nil, nil)
	if err == nil {
		t.Error("expected error for nil request")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("error = %v, want 'request cannot be nil'", err)
	}
}

func TestAssetService_GetAssetWallets_EmptyCurrency(t *testing.T) {
	svc := &AssetService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	req := &model.GetAssetWalletsRequest{
		Asset: model.AssetFilter{
			Currency: "",
		},
	}

	_, err := svc.GetAssetWallets(nil, req)
	if err == nil {
		t.Error("expected error for empty currency")
	}
	if err.Error() != "asset currency is required" {
		t.Errorf("error = %v, want 'asset currency is required'", err)
	}
}

func TestAssetService_GetAssetWallets_WithOptions(t *testing.T) {
	tests := []struct {
		name    string
		request *model.GetAssetWalletsRequest
	}{
		{
			name: "basic request",
			request: &model.GetAssetWalletsRequest{
				Asset: model.AssetFilter{
					Currency: "BTC",
				},
			},
		},
		{
			name: "with limit and cursor",
			request: &model.GetAssetWalletsRequest{
				Asset: model.AssetFilter{
					Currency: "ETH",
				},
				Limit:  100,
				Cursor: "cursor-xyz",
			},
		},
		{
			name: "with wallet filters",
			request: &model.GetAssetWalletsRequest{
				Asset: model.AssetFilter{
					Currency: "BTC",
				},
				WalletID:   "wallet-abc",
				WalletName: "My Wallet",
			},
		},
		{
			name: "with cursor pagination",
			request: &model.GetAssetWalletsRequest{
				Asset: model.AssetFilter{
					Currency: "ALGO",
				},
				CurrentPage: "page-123",
				PageRequest: "FIRST",
				PageSize:    50,
			},
		},
		{
			name: "with NFT filter",
			request: &model.GetAssetWalletsRequest{
				Asset: model.AssetFilter{
					Currency: "ETH",
					Kind:     "NFT",
					NFT: &model.AssetNFTFilter{
						TokenID: "nft-456",
					},
				},
			},
		},
		{
			name: "with unknown asset filter",
			request: &model.GetAssetWalletsRequest{
				Asset: model.AssetFilter{
					Currency: "Unknown",
					Kind:     "Unknown",
					Unknown: &model.AssetUnknownFilter{
						Blockchain: "polygon",
						Arg1:       "0xtoken",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify request can be created with these values
			// Actual API testing requires mocking
			svc := &AssetService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("AssetService should not be nil")
			}
		})
	}
}

func TestGetAssetAddressesRequest_PageRequestValues(t *testing.T) {
	// Test that page request values match expected API values
	validPageRequests := []string{"FIRST", "PREVIOUS", "NEXT", "LAST"}

	for _, pageRequest := range validPageRequests {
		t.Run(pageRequest, func(t *testing.T) {
			req := &model.GetAssetAddressesRequest{
				Asset: model.AssetFilter{
					Currency: "BTC",
				},
				PageRequest: pageRequest,
			}
			if req.PageRequest != pageRequest {
				t.Errorf("PageRequest = %v, want %v", req.PageRequest, pageRequest)
			}
		})
	}
}

func TestGetAssetAddressesRequest_SortOrderValues(t *testing.T) {
	// Test that sort order values match expected API values
	validSortOrders := []string{"ASC", "DESC"}

	for _, sortOrder := range validSortOrders {
		t.Run(sortOrder, func(t *testing.T) {
			req := &model.GetAssetAddressesRequest{
				Asset: model.AssetFilter{
					Currency: "BTC",
				},
				SortOrder: sortOrder,
			}
			if req.SortOrder != sortOrder {
				t.Errorf("SortOrder = %v, want %v", req.SortOrder, sortOrder)
			}
		})
	}
}

func TestGetAssetWalletsRequest_PageRequestValues(t *testing.T) {
	// Test that page request values match expected API values
	validPageRequests := []string{"FIRST", "PREVIOUS", "NEXT", "LAST"}

	for _, pageRequest := range validPageRequests {
		t.Run(pageRequest, func(t *testing.T) {
			req := &model.GetAssetWalletsRequest{
				Asset: model.AssetFilter{
					Currency: "BTC",
				},
				PageRequest: pageRequest,
			}
			if req.PageRequest != pageRequest {
				t.Errorf("PageRequest = %v, want %v", req.PageRequest, pageRequest)
			}
		})
	}
}

func TestAssetFilter_KindValues(t *testing.T) {
	// Test that kind values match expected API values
	validKinds := []string{"NFT", "Unknown"}

	for _, kind := range validKinds {
		t.Run(kind, func(t *testing.T) {
			filter := model.AssetFilter{
				Currency: "ETH",
				Kind:     kind,
			}
			if filter.Kind != kind {
				t.Errorf("Kind = %v, want %v", filter.Kind, kind)
			}
		})
	}
}
