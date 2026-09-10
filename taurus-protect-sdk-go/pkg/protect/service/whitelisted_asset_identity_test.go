package service

import (
	"math"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestAssetPagination(t *testing.T) {
	total := "250"

	t.Run("more rows remain", func(t *testing.T) {
		p := assetPagination(&total, 100, 0)
		if p == nil || !p.HasMore || p.TotalItems != 250 {
			t.Errorf("got %+v, want HasMore with TotalItems 250", p)
		}
	})

	t.Run("last page", func(t *testing.T) {
		p := assetPagination(&total, 100, 200)
		if p == nil || p.HasMore {
			t.Errorf("got %+v, want HasMore false", p)
		}
	})

	t.Run("absent total yields no pagination", func(t *testing.T) {
		if p := assetPagination(nil, 100, 0); p != nil {
			t.Errorf("got %+v, want nil", p)
		}
	})

	t.Run("offset+limit cannot wrap into a false HasMore", func(t *testing.T) {
		// The unsafe form (offset+limit < total) overflows to a negative here and
		// reports another page that does not exist.
		p := assetPagination(&total, 1, math.MaxInt64)
		if p == nil || p.HasMore {
			t.Errorf("got %+v, want HasMore false", p)
		}
	})

	t.Run("unparseable total leaves the window without a count", func(t *testing.T) {
		bad := "not-a-number"
		p := assetPagination(&bad, 100, 0)
		if p == nil || p.TotalItems != 0 || p.HasMore {
			t.Errorf("got %+v, want a zero count and no HasMore", p)
		}
	})
}

// The DTO mapper cannot supply the identity fields — they are payload-only — so a verified
// asset used to come back with an empty ContractAddress on Get and List.
func TestPopulateVerifiedIdentity(t *testing.T) {
	payload := `{"blockchain":"ETH","network":"mainnet","contractAddress":"0xA0b8","` +
		`name":"USD Coin","symbol":"USDC","decimals":6,"tokenId":"42"}`

	asset := &model.WhitelistedAsset{
		ID:         "1",
		Blockchain: "ETH",
		Network:    "mainnet",
		Metadata:   &model.WhitelistedAssetMetadata{PayloadAsString: payload},
	}

	if err := populateVerifiedIdentity(asset); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if asset.ContractAddress != "0xA0b8" {
		t.Errorf("ContractAddress = %q, want 0xA0b8", asset.ContractAddress)
	}
	if asset.Name != "USD Coin" {
		t.Errorf("Name = %q, want USD Coin", asset.Name)
	}
	if asset.Symbol != "USDC" {
		t.Errorf("Symbol = %q, want USDC", asset.Symbol)
	}
	if asset.Decimals != 6 {
		t.Errorf("Decimals = %d, want 6", asset.Decimals)
	}
	if asset.TokenID != "42" {
		t.Errorf("TokenID = %q, want 42", asset.TokenID)
	}
}

func TestPopulateVerifiedIdentity_RejectsMissingMetadata(t *testing.T) {
	if err := populateVerifiedIdentity(&model.WhitelistedAsset{ID: "1"}); err == nil {
		t.Error("expected an error when metadata is absent")
	}
}

func TestPopulateVerifiedIdentity_RejectsUnparseablePayload(t *testing.T) {
	asset := &model.WhitelistedAsset{
		ID:       "1",
		Metadata: &model.WhitelistedAssetMetadata{PayloadAsString: `{"contractAddress":`},
	}
	if err := populateVerifiedIdentity(asset); err == nil {
		t.Error("expected an error for a malformed payload")
	}
	if asset.ContractAddress != "" {
		t.Error("no field may be assigned from a payload that failed to parse")
	}
}
