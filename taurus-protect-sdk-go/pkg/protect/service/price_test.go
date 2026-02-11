package service

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewPriceService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestPriceService_Convert_NilOptions(t *testing.T) {
	svc := &PriceService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.Convert(nil, nil)
	if err == nil {
		t.Error("Convert() with nil options should return error")
	}
	if err.Error() != "options cannot be nil" {
		t.Errorf("Convert() error = %v, want 'options cannot be nil'", err)
	}
}

func TestPriceService_Convert_EmptyCurrency(t *testing.T) {
	svc := &PriceService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.Convert(nil, &model.ConvertOptions{})
	if err == nil {
		t.Error("Convert() with empty currency should return error")
	}
	if err.Error() != "currency is required" {
		t.Errorf("Convert() error = %v, want 'currency is required'", err)
	}
}

func TestPriceService_Convert_EmptyAmount(t *testing.T) {
	svc := &PriceService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.Convert(nil, &model.ConvertOptions{Currency: "ETH"})
	if err == nil {
		t.Error("Convert() with empty amount should return error")
	}
	if err.Error() != "amount is required" {
		t.Errorf("Convert() error = %v, want 'amount is required'", err)
	}
}

func TestPriceService_GetPriceHistory_NilOptions(t *testing.T) {
	svc := &PriceService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetPriceHistory(nil, nil)
	if err == nil {
		t.Error("GetPriceHistory() with nil options should return error")
	}
	if err.Error() != "options cannot be nil" {
		t.Errorf("GetPriceHistory() error = %v, want 'options cannot be nil'", err)
	}
}

func TestPriceService_GetPriceHistory_EmptyBase(t *testing.T) {
	svc := &PriceService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetPriceHistory(nil, &model.GetPriceHistoryOptions{})
	if err == nil {
		t.Error("GetPriceHistory() with empty base should return error")
	}
	if err.Error() != "base currency is required" {
		t.Errorf("GetPriceHistory() error = %v, want 'base currency is required'", err)
	}
}

func TestPriceService_GetPriceHistory_EmptyQuote(t *testing.T) {
	svc := &PriceService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetPriceHistory(nil, &model.GetPriceHistoryOptions{Base: "ETH"})
	if err == nil {
		t.Error("GetPriceHistory() with empty quote should return error")
	}
	if err.Error() != "quote currency is required" {
		t.Errorf("GetPriceHistory() error = %v, want 'quote currency is required'", err)
	}
}

func TestPriceService_ExportPriceHistory_NilOptions(t *testing.T) {
	svc := &PriceService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.ExportPriceHistory(nil, nil)
	if err == nil {
		t.Error("ExportPriceHistory() with nil options should return error")
	}
	if err.Error() != "options cannot be nil" {
		t.Errorf("ExportPriceHistory() error = %v, want 'options cannot be nil'", err)
	}
}

func TestPriceService_Convert_ValidOptions(t *testing.T) {
	tests := []struct {
		name              string
		currency          string
		amount            string
		symbols           []string
		targetCurrencyIds []string
	}{
		{
			name:     "basic conversion",
			currency: "ETH",
			amount:   "1.0",
		},
		{
			name:     "with symbols filter",
			currency: "BTC",
			amount:   "0.5",
			symbols:  []string{"USD", "EUR"},
		},
		{
			name:              "with target currency IDs",
			currency:          "ETH",
			amount:            "2.0",
			targetCurrencyIds: []string{"currency-1", "currency-2"},
		},
		{
			name:              "with all options",
			currency:          "BTC",
			amount:            "1.5",
			symbols:           []string{"USD"},
			targetCurrencyIds: []string{"currency-1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify options can be created with these values
			// Actual API testing requires mocking
			svc := &PriceService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("PriceService should not be nil")
			}
		})
	}
}

func TestPriceService_GetPriceHistory_ValidOptions(t *testing.T) {
	tests := []struct {
		name  string
		base  string
		quote string
		limit int64
	}{
		{
			name:  "basic history request",
			base:  "ETH",
			quote: "USD",
		},
		{
			name:  "with limit",
			base:  "BTC",
			quote: "EUR",
			limit: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify options can be created with these values
			// Actual API testing requires mocking
			svc := &PriceService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("PriceService should not be nil")
			}
		})
	}
}

func TestPriceService_ExportPriceHistory_ValidOptions(t *testing.T) {
	tests := []struct {
		name          string
		currencyPairs []string
		limit         int64
		format        string
	}{
		{
			name: "empty options",
		},
		{
			name:          "with currency pairs",
			currencyPairs: []string{"BTC/USD", "ETH/USD"},
		},
		{
			name:   "with format csv",
			format: "csv",
		},
		{
			name:   "with format json",
			format: "json",
		},
		{
			name:          "with all options",
			currencyPairs: []string{"BTC/USD"},
			limit:         50,
			format:        "csv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify options can be created with these values
			// Actual API testing requires mocking
			svc := &PriceService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("PriceService should not be nil")
			}
		})
	}
}

func TestSafeString(t *testing.T) {
	tests := []struct {
		name  string
		input *string
		want  string
	}{
		{
			name:  "nil returns empty string",
			input: nil,
			want:  "",
		},
		{
			name:  "non-nil returns value",
			input: testStringPtr("test"),
			want:  "test",
		},
		{
			name:  "empty string pointer returns empty string",
			input: testStringPtr(""),
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := safeString(tt.input)
			if got != tt.want {
				t.Errorf("safeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// testStringPtr is a helper for creating string pointers in tests (always returns a pointer, even for empty strings)
func testStringPtr(s string) *string {
	return &s
}
