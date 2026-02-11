package model

import (
	"testing"
)

func TestRequestMetadata_ParsePayloadEntries(t *testing.T) {
	tests := []struct {
		name           string
		payloadAsStr   string
		expectedKeys   []string
		expectError    bool
		expectNil      bool
	}{
		{
			name: "valid array payload",
			payloadAsStr: `[
				{"key": "request_id", "value": "123"},
				{"key": "currency", "value": "BTC"},
				{"key": "source", "value": {"payload": {"address": "addr1"}}},
				{"key": "destination", "value": {"payload": {"address": "addr2"}}}
			]`,
			expectedKeys: []string{"request_id", "currency", "source", "destination"},
			expectError:  false,
			expectNil:    false,
		},
		{
			name:         "empty string",
			payloadAsStr: "",
			expectedKeys: nil,
			expectError:  false,
			expectNil:    true,
		},
		{
			name:         "invalid JSON",
			payloadAsStr: "not json",
			expectError:  true,
		},
		{
			name:         "empty array",
			payloadAsStr: "[]",
			expectedKeys: []string{},
			expectError:  false,
			expectNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &RequestMetadata{PayloadAsString: tt.payloadAsStr}
			entries, err := m.ParsePayloadEntries()

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.expectNil {
				if entries != nil {
					t.Errorf("expected nil entries, got %v", entries)
				}
				return
			}

			if len(entries) != len(tt.expectedKeys) {
				t.Errorf("expected %d entries, got %d", len(tt.expectedKeys), len(entries))
			}

			for i, key := range tt.expectedKeys {
				if entries[i].Key != key {
					t.Errorf("entry[%d].Key = %q, want %q", i, entries[i].Key, key)
				}
			}
		})
	}
}

func TestRequestMetadata_GetPayloadValue(t *testing.T) {
	payload := `[
		{"key": "currency", "value": "BTC"},
		{"key": "request_id", "value": 12345},
		{"key": "source", "value": {"payload": {"address": "src_addr"}}}
	]`

	m := &RequestMetadata{PayloadAsString: payload}

	t.Run("existing string value", func(t *testing.T) {
		val := m.GetPayloadValue("currency")
		if val == nil {
			t.Fatal("expected non-nil value")
		}
		if s, ok := val.raw.(string); !ok || s != "BTC" {
			t.Errorf("currency = %v, want 'BTC'", val.raw)
		}
	})

	t.Run("existing nested value", func(t *testing.T) {
		val := m.GetPayloadValue("source")
		if val == nil {
			t.Fatal("expected non-nil value")
		}
		addr := val.GetString("payload", "address")
		if addr != "src_addr" {
			t.Errorf("source address = %q, want 'src_addr'", addr)
		}
	})

	t.Run("non-existent key", func(t *testing.T) {
		val := m.GetPayloadValue("nonexistent")
		if val != nil {
			t.Errorf("expected nil for nonexistent key, got %v", val)
		}
	})
}

func TestRequestMetadata_GetSourceAddress(t *testing.T) {
	tests := []struct {
		name     string
		payload  string
		expected string
	}{
		{
			name: "has source address",
			payload: `[
				{"key": "source", "value": {"payload": {"address": "source_addr_123"}}}
			]`,
			expected: "source_addr_123",
		},
		{
			name:     "no source entry",
			payload:  `[{"key": "currency", "value": "BTC"}]`,
			expected: "",
		},
		{
			name:     "empty payload",
			payload:  "",
			expected: "",
		},
		{
			name: "source without nested address",
			payload: `[
				{"key": "source", "value": {"other": "data"}}
			]`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &RequestMetadata{PayloadAsString: tt.payload}
			got := m.GetSourceAddress()
			if got != tt.expected {
				t.Errorf("GetSourceAddress() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestRequestMetadata_GetDestinationAddress(t *testing.T) {
	payload := `[
		{"key": "destination", "value": {"payload": {"address": "dest_addr_456"}}}
	]`

	m := &RequestMetadata{PayloadAsString: payload}
	got := m.GetDestinationAddress()
	if got != "dest_addr_456" {
		t.Errorf("GetDestinationAddress() = %q, want 'dest_addr_456'", got)
	}
}

func TestRequestMetadata_GetMetadataCurrency(t *testing.T) {
	payload := `[{"key": "currency", "value": "ETH"}]`

	m := &RequestMetadata{PayloadAsString: payload}
	got := m.GetMetadataCurrency()
	if got != "ETH" {
		t.Errorf("GetMetadataCurrency() = %q, want 'ETH'", got)
	}
}

func TestRequestMetadata_GetMetadataRequestID(t *testing.T) {
	payload := `[{"key": "request_id", "value": 12345}]`

	m := &RequestMetadata{PayloadAsString: payload}
	got := m.GetMetadataRequestID()
	if got != 12345 {
		t.Errorf("GetMetadataRequestID() = %d, want 12345", got)
	}
}

func TestRequestMetadata_GetAmount(t *testing.T) {
	t.Run("has amount with numeric values", func(t *testing.T) {
		payload := `[{
			"key": "amount",
			"value": {
				"valueFrom": 1000000,
				"valueTo": 10.5,
				"rate": 0.00001,
				"decimals": 8,
				"currencyFrom": "BTC",
				"currencyTo": "USD"
			}
		}]`

		m := &RequestMetadata{PayloadAsString: payload}
		amount := m.GetAmount()

		if amount == nil {
			t.Fatal("expected non-nil amount")
		}
		if amount.ValueFrom != "1000000" {
			t.Errorf("ValueFrom = %q, want \"1000000\"", amount.ValueFrom)
		}
		if amount.ValueTo != "10.5" {
			t.Errorf("ValueTo = %q, want \"10.5\"", amount.ValueTo)
		}
		if amount.Rate != "0.00001" {
			t.Errorf("Rate = %q, want \"0.00001\"", amount.Rate)
		}
		if amount.Decimals != 8 {
			t.Errorf("Decimals = %d, want 8", amount.Decimals)
		}
		if amount.CurrencyFrom != "BTC" {
			t.Errorf("CurrencyFrom = %q, want \"BTC\"", amount.CurrencyFrom)
		}
		if amount.CurrencyTo != "USD" {
			t.Errorf("CurrencyTo = %q, want \"USD\"", amount.CurrencyTo)
		}
	})

	t.Run("has amount with string values", func(t *testing.T) {
		payload := `[{
			"key": "amount",
			"value": {
				"valueFrom": "1000000",
				"valueTo": "10.5",
				"rate": "0.00001",
				"decimals": 8,
				"currencyFrom": "BTC",
				"currencyTo": "USD"
			}
		}]`

		m := &RequestMetadata{PayloadAsString: payload}
		amount := m.GetAmount()

		if amount == nil {
			t.Fatal("expected non-nil amount")
		}
		if amount.ValueFrom != "1000000" {
			t.Errorf("ValueFrom = %q, want \"1000000\"", amount.ValueFrom)
		}
		if amount.ValueTo != "10.5" {
			t.Errorf("ValueTo = %q, want \"10.5\"", amount.ValueTo)
		}
		if amount.Rate != "0.00001" {
			t.Errorf("Rate = %q, want \"0.00001\"", amount.Rate)
		}
		if amount.Decimals != 8 {
			t.Errorf("Decimals = %d, want 8", amount.Decimals)
		}
		if amount.CurrencyFrom != "BTC" {
			t.Errorf("CurrencyFrom = %q, want \"BTC\"", amount.CurrencyFrom)
		}
		if amount.CurrencyTo != "USD" {
			t.Errorf("CurrencyTo = %q, want \"USD\"", amount.CurrencyTo)
		}
	})

	t.Run("no amount", func(t *testing.T) {
		payload := `[{"key": "currency", "value": "BTC"}]`

		m := &RequestMetadata{PayloadAsString: payload}
		amount := m.GetAmount()

		if amount != nil {
			t.Errorf("expected nil amount, got %+v", amount)
		}
	})
}

func TestPayloadEntryValue_GetString(t *testing.T) {
	val := NewPayloadEntryValue(map[string]interface{}{
		"level1": map[string]interface{}{
			"level2": "deep_value",
		},
	})

	t.Run("nested path", func(t *testing.T) {
		got := val.GetString("level1", "level2")
		if got != "deep_value" {
			t.Errorf("GetString() = %q, want 'deep_value'", got)
		}
	})

	t.Run("invalid path", func(t *testing.T) {
		got := val.GetString("invalid", "path")
		if got != "" {
			t.Errorf("GetString() = %q, want empty string", got)
		}
	})

	t.Run("nil value", func(t *testing.T) {
		var nilVal *PayloadEntryValue
		got := nilVal.GetString("any")
		if got != "" {
			t.Errorf("GetString() on nil = %q, want empty string", got)
		}
	})
}

func TestPayloadEntryValue_GetInt64(t *testing.T) {
	val := NewPayloadEntryValue(map[string]interface{}{
		"count": float64(42), // JSON numbers are float64
	})

	t.Run("numeric value", func(t *testing.T) {
		got := val.GetInt64("count")
		if got != 42 {
			t.Errorf("GetInt64() = %d, want 42", got)
		}
	})

	t.Run("missing key", func(t *testing.T) {
		got := val.GetInt64("missing")
		if got != 0 {
			t.Errorf("GetInt64() = %d, want 0", got)
		}
	})
}

func TestPayloadEntryValue_GetFloat64(t *testing.T) {
	val := NewPayloadEntryValue(map[string]interface{}{
		"rate": float64(0.123),
	})

	got := val.GetFloat64("rate")
	if got != 0.123 {
		t.Errorf("GetFloat64() = %f, want 0.123", got)
	}
}

func TestNilRequestMetadata(t *testing.T) {
	var m *RequestMetadata

	t.Run("ParsePayloadEntries", func(t *testing.T) {
		entries, err := m.ParsePayloadEntries()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if entries != nil {
			t.Errorf("expected nil entries, got %v", entries)
		}
	})

	t.Run("GetPayloadValue", func(t *testing.T) {
		val := m.GetPayloadValue("any")
		if val != nil {
			t.Errorf("expected nil, got %v", val)
		}
	})

	t.Run("GetSourceAddress", func(t *testing.T) {
		got := m.GetSourceAddress()
		if got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})

	t.Run("GetDestinationAddress", func(t *testing.T) {
		got := m.GetDestinationAddress()
		if got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})

	t.Run("GetMetadataCurrency", func(t *testing.T) {
		got := m.GetMetadataCurrency()
		if got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})

	t.Run("GetMetadataRequestID", func(t *testing.T) {
		got := m.GetMetadataRequestID()
		if got != 0 {
			t.Errorf("expected 0, got %d", got)
		}
	})

	t.Run("GetAmount", func(t *testing.T) {
		amount := m.GetAmount()
		if amount != nil {
			t.Errorf("expected nil, got %+v", amount)
		}
	})
}
