package service

import (
	"math"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// The contracts list's next page starts at offset + limit: validatord skips rows it cannot
// decode without giving up their SQL slot, so a short page is not the end of the list.
func TestWhitelistedAssetPagination(t *testing.T) {
	total := "250"
	page := func(t *testing.T, total *string, served int, limit, offset int64) *model.Pagination {
		t.Helper()
		p, err := offsetPagination(rulePlusLimit, offsetWindow{limit: limit, offset: offset}, served, 0,
			offsetReply{TotalItems: total})
		if err != nil {
			t.Fatal(err)
		}
		return p
	}

	t.Run("more rows remain", func(t *testing.T) {
		p := page(t, &total, 100, 100, 0)
		if !p.HasMore || p.TotalItems != 250 || p.NextOffset != 100 {
			t.Errorf("got %+v, want HasMore, TotalItems 250 and NextOffset 100", p)
		}
	})

	t.Run("a short page is not the end", func(t *testing.T) {
		p := page(t, &total, 97, 100, 0)
		if !p.HasMore || p.NextOffset != 100 {
			t.Errorf("got %+v, want HasMore with NextOffset 100: skipped rows keep their slot", p)
		}
	})

	t.Run("last page", func(t *testing.T) {
		p := page(t, &total, 50, 100, 200)
		if p.HasMore {
			t.Errorf("got %+v, want HasMore false", p)
		}
	})

	t.Run("an empty reply is a full empty page", func(t *testing.T) {
		p := page(t, nil, 0, 100, 0)
		want := model.Pagination{Limit: 100, Offset: 0, TotalItems: 0, NextOffset: 100, HasMore: false}
		if *p != want {
			t.Errorf("got %+v, want %+v", p, want)
		}
	})

	t.Run("offset+limit cannot wrap into a false HasMore", func(t *testing.T) {
		// The unsafe form (offset+limit) overflows to a negative here.
		p := page(t, &total, 0, 1, math.MaxInt64)
		if p.HasMore || p.NextOffset != math.MaxInt64 {
			t.Errorf("got %+v, want HasMore false with a saturated NextOffset", p)
		}
	})

	t.Run("an unparseable total is an error, not a zero count", func(t *testing.T) {
		bad := "not-a-number"
		if _, err := offsetPagination(rulePlusLimit, offsetWindow{limit: 100}, 0, 0,
			offsetReply{TotalItems: &bad}); err == nil {
			t.Error("a count silently read as 0 ends a walk early; it must be an error")
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

	if err := populateVerifiedIdentity(asset, payload); err != nil {
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
	if err := populateVerifiedIdentity(&model.WhitelistedAsset{ID: "1"}, "{}"); err == nil {
		t.Error("expected an error when metadata is absent")
	}
}

func TestPopulateVerifiedIdentity_RejectsUnparseablePayload(t *testing.T) {
	asset := &model.WhitelistedAsset{
		ID:       "1",
		Metadata: &model.WhitelistedAssetMetadata{PayloadAsString: `{"contractAddress":`},
	}
	if err := populateVerifiedIdentity(asset, `{"contractAddress":`); err == nil {
		t.Error("expected an error for a malformed payload")
	}
	if asset.ContractAddress != "" {
		t.Error("no field may be assigned from a payload that failed to parse")
	}
}
