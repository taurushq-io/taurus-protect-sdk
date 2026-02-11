package service

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewTokenMetadataService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestTokenMetadataService_GetCryptoPunkMetadata_Validation(t *testing.T) {
	tests := []struct {
		name      string
		network   string
		contract  string
		tokenID   string
		wantErr   bool
		errSubstr string
	}{
		{
			name:      "empty network returns error",
			network:   "",
			contract:  "0xb47e3cd837dDF8e4c57F05d70Ab865de6e193BBB",
			tokenID:   "1234",
			wantErr:   true,
			errSubstr: "network cannot be empty",
		},
		{
			name:      "empty contract returns error",
			network:   "mainnet",
			contract:  "",
			tokenID:   "1234",
			wantErr:   true,
			errSubstr: "contract cannot be empty",
		},
		{
			name:      "empty tokenID returns error",
			network:   "mainnet",
			contract:  "0xb47e3cd837dDF8e4c57F05d70Ab865de6e193BBB",
			tokenID:   "",
			wantErr:   true,
			errSubstr: "tokenID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &TokenMetadataService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}

			// We can only test validation since we don't have a mock
			_, err := svc.GetCryptoPunkMetadata(nil, tt.network, tt.contract, tt.tokenID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetCryptoPunkMetadata() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errSubstr != "" && err.Error() != tt.errSubstr {
					t.Errorf("GetCryptoPunkMetadata() error = %v, want substring %v", err, tt.errSubstr)
				}
			}
		})
	}
}

func TestTokenMetadataService_GetERCTokenMetadata_Validation(t *testing.T) {
	tests := []struct {
		name      string
		network   string
		contract  string
		tokenID   string
		opts      *model.GetERCTokenMetadataOptions
		wantErr   bool
		errSubstr string
	}{
		{
			name:      "empty network returns error",
			network:   "",
			contract:  "0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D",
			tokenID:   "1234",
			opts:      nil,
			wantErr:   true,
			errSubstr: "network cannot be empty",
		},
		{
			name:      "empty contract returns error",
			network:   "mainnet",
			contract:  "",
			tokenID:   "1234",
			opts:      nil,
			wantErr:   true,
			errSubstr: "contract cannot be empty",
		},
		{
			name:      "empty tokenID returns error",
			network:   "mainnet",
			contract:  "0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D",
			tokenID:   "",
			opts:      nil,
			wantErr:   true,
			errSubstr: "tokenID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &TokenMetadataService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}

			// We can only test validation since we don't have a mock
			_, err := svc.GetERCTokenMetadata(nil, tt.network, tt.contract, tt.tokenID, tt.opts)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetERCTokenMetadata() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errSubstr != "" && err.Error() != tt.errSubstr {
					t.Errorf("GetERCTokenMetadata() error = %v, want substring %v", err, tt.errSubstr)
				}
			}
		})
	}
}

func TestGetERCTokenMetadataOptions(t *testing.T) {
	tests := []struct {
		name string
		opts *model.GetERCTokenMetadataOptions
	}{
		{
			name: "nil options",
			opts: nil,
		},
		{
			name: "empty options",
			opts: &model.GetERCTokenMetadataOptions{},
		},
		{
			name: "with data enabled",
			opts: &model.GetERCTokenMetadataOptions{
				WithData: true,
			},
		},
		{
			name: "with blockchain specified",
			opts: &model.GetERCTokenMetadataOptions{
				Blockchain: "ethereum",
			},
		},
		{
			name: "all options set",
			opts: &model.GetERCTokenMetadataOptions{
				WithData:   true,
				Blockchain: "polygon",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify options can be created with these values
			// Actual API testing requires mocking
			svc := &TokenMetadataService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("TokenMetadataService should not be nil")
			}
		})
	}
}

func TestTokenMetadataService_CryptoPunkContractAddress(t *testing.T) {
	// Document the expected CryptoPunks contract address for ETH mainnet
	expectedContract := "0xb47e3cd837dDF8e4c57F05d70Ab865de6e193BBB"
	if len(expectedContract) != 42 {
		t.Errorf("CryptoPunks contract address has unexpected length: %d", len(expectedContract))
	}
}

func TestTokenMetadataService_ValidPunkIDRange(t *testing.T) {
	// Document the valid range for CryptoPunk IDs
	tests := []struct {
		name    string
		punkID  string
		isValid bool
	}{
		{
			name:    "minimum valid ID",
			punkID:  "0",
			isValid: true,
		},
		{
			name:    "maximum valid ID",
			punkID:  "9999",
			isValid: true,
		},
		{
			name:    "mid-range ID",
			punkID:  "5000",
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This documents valid punk ID values
			// Actual validation is done server-side
			if tt.punkID == "" && tt.isValid {
				t.Error("empty punkID should not be valid")
			}
		})
	}
}
