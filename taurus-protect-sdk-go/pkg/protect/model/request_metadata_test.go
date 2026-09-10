package model

import (
	"errors"
	"strings"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
)

func TestRequestMetadata_ParsePayloadEntries(t *testing.T) {
	tests := []struct {
		name         string
		payloadAsStr string
		expectedKeys []string
		expectError  bool
		expectNil    bool
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
			m := &RequestMetadata{PayloadAsString: tt.payloadAsStr} // fallback path, deliberately unverified
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

	m := verifiedMetadata(t, payload)

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
			name: "source without nested address",
			payload: `[
				{"key": "source", "value": {"other": "data"}}
			]`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := verifiedMetadata(t, tt.payload)
			got, _ := m.GetSourceAddress()
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

	m := verifiedMetadata(t, payload)
	got, _ := m.GetDestinationAddress()
	if got != "dest_addr_456" {
		t.Errorf("GetDestinationAddress() = %q, want 'dest_addr_456'", got)
	}
}

func TestRequestMetadata_GetMetadataCurrency(t *testing.T) {
	payload := `[{"key": "currency", "value": "ETH"}]`

	m := verifiedMetadata(t, payload)
	got, _ := m.GetMetadataCurrency()
	if got != "ETH" {
		t.Errorf("GetMetadataCurrency() = %q, want 'ETH'", got)
	}
}

func TestRequestMetadata_GetMetadataRequestID(t *testing.T) {
	payload := `[{"key": "request_id", "value": 12345}]`

	m := verifiedMetadata(t, payload)
	got, _ := m.GetMetadataRequestID()
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

		m := verifiedMetadata(t, payload)
		amount, _ := m.GetAmount()

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

		m := verifiedMetadata(t, payload)
		amount, _ := m.GetAmount()

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

		m := verifiedMetadata(t, payload)
		amount, _ := m.GetAmount()

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
		got, _ := m.GetSourceAddress()
		if got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})

	t.Run("GetDestinationAddress", func(t *testing.T) {
		got, _ := m.GetDestinationAddress()
		if got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})

	t.Run("GetMetadataCurrency", func(t *testing.T) {
		got, _ := m.GetMetadataCurrency()
		if got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})

	t.Run("GetMetadataRequestID", func(t *testing.T) {
		got, _ := m.GetMetadataRequestID()
		if got != 0 {
			t.Errorf("expected 0, got %d", got)
		}
	})

	t.Run("GetAmount", func(t *testing.T) {
		amount, _ := m.GetAmount()
		if amount != nil {
			t.Errorf("expected nil, got %+v", amount)
		}
	})
}

// verifiedMetadata builds metadata that has been through the real verification
// path, which is the only way accessors return anything: they read the unexported
// `entries`, and only VerifyAndMaterialise populates it.
// A hash with nothing to check it against is rejected rather than passing with an
// empty result: there is no payload to verify, so accepting would mean trusting the
// response's own claim.
func TestRequestMetadata_HashWithNoPayloadIsRejected(t *testing.T) {
	m := &RequestMetadata{PayloadAsString: "", Hash: crypto.CalculateHexHash("")}

	err := m.VerifyAndMaterialise()
	if err == nil {
		t.Fatal("VerifyAndMaterialise() should reject a hash with no payload")
	}
	if !strings.Contains(err.Error(), "payload is missing") {
		t.Errorf("error should name the missing payload, got %q", err.Error())
	}
	if m.HashVerified {
		t.Error("HashVerified must stay false")
	}
}

func verifiedMetadata(t *testing.T, payload string) *RequestMetadata {
	t.Helper()
	m := &RequestMetadata{
		PayloadAsString: payload,
		Hash:            crypto.CalculateHexHash(payload),
	}
	if err := m.VerifyAndMaterialise(); err != nil {
		t.Fatalf("VerifyAndMaterialise: %v", err)
	}
	return m
}

// Accessors must yield NOTHING from metadata nobody verified. This is the property
// that makes unverified structure unrepresentable rather than merely discouraged,
// and it is what the list paths were silently violating.
func TestRequestMetadata_AccessorsReturnZeroWithoutVerification(t *testing.T) {
	payload := `[
		{"key": "source", "value": {"payload": {"address": "src_addr"}}},
		{"key": "destination", "value": {"payload": {"address": "dst_addr"}}},
		{"key": "currency", "value": "ETH"},
		{"key": "request_id", "value": 12345}
	]`
	// PayloadAsString is set and parseable, but VerifyAndMaterialise never ran.
	m := &RequestMetadata{PayloadAsString: payload}

	// Each accessor must SAY it could not verify, rather than returning an empty
	// value that reads identically to "the payload has no source address".
	if _, err := m.GetSourceAddress(); !errors.Is(err, ErrMetadataUnverified) {
		t.Errorf("GetSourceAddress() err = %v, want ErrMetadataUnverified", err)
	}
	if _, err := m.GetDestinationAddress(); !errors.Is(err, ErrMetadataUnverified) {
		t.Errorf("GetDestinationAddress() err = %v, want ErrMetadataUnverified", err)
	}
	if _, err := m.GetMetadataCurrency(); !errors.Is(err, ErrMetadataUnverified) {
		t.Errorf("GetMetadataCurrency() err = %v, want ErrMetadataUnverified", err)
	}
	if _, err := m.GetMetadataRequestID(); !errors.Is(err, ErrMetadataUnverified) {
		t.Errorf("GetMetadataRequestID() err = %v, want ErrMetadataUnverified", err)
	}
	if _, err := m.GetAmount(); !errors.Is(err, ErrMetadataUnverified) {
		t.Errorf("GetAmount() err = %v, want ErrMetadataUnverified", err)
	}
	if got := m.GetPayloadValue("currency"); got != nil {
		t.Errorf("GetPayloadValue() = %v on unverified metadata, want nil", got)
	}
	if m.HashVerified {
		t.Error("HashVerified is true without verification")
	}

	// Metadata carrying no payload at all is not an error: a request in an early
	// status has nothing to verify and nothing to read.
	empty := &RequestMetadata{}
	if got, err := empty.GetSourceAddress(); err != nil || got != "" {
		t.Errorf("GetSourceAddress() on empty metadata = %q, %v; want \"\", nil", got, err)
	}
}

func TestRequestMetadata_VerifyAndMaterialise(t *testing.T) {
	payload := `[{"key": "currency", "value": "BTC"}]`

	t.Run("hash matches: entries populated and marked verified", func(t *testing.T) {
		m := &RequestMetadata{PayloadAsString: payload, Hash: crypto.CalculateHexHash(payload)}
		if err := m.VerifyAndMaterialise(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !m.HashVerified {
			t.Error("HashVerified not set")
		}
		if cur, err := m.GetMetadataCurrency(); err != nil || cur != "BTC" {
			t.Errorf("accessor returned %q, %v after verification; want BTC, nil", cur, err)
		}
	})

	t.Run("payload altered without updating hash: rejected", func(t *testing.T) {
		m := &RequestMetadata{
			PayloadAsString: `[{"key": "currency", "value": "ETH"}]`, // tampered
			Hash:            crypto.CalculateHexHash(payload),        // hash of the original
		}
		if err := m.VerifyAndMaterialise(); err == nil {
			t.Fatal("expected an IntegrityError for a tampered payload")
		}
		if m.HashVerified {
			t.Error("HashVerified set despite a failed verification")
		}
		if _, err := m.GetMetadataCurrency(); !errors.Is(err, ErrMetadataUnverified) {
			t.Errorf("accessor err = %v after a failed verification, want ErrMetadataUnverified", err)
		}
	})

	t.Run("verified but unparseable: error, not a swallowed nil", func(t *testing.T) {
		bad := `{"not": "an array"}`
		m := &RequestMetadata{PayloadAsString: bad, Hash: crypto.CalculateHexHash(bad)}
		if err := m.VerifyAndMaterialise(); err == nil {
			t.Fatal("expected an error for a verified but malformed payload")
		}
	})

	t.Run("no metadata yet: not an error", func(t *testing.T) {
		m := &RequestMetadata{}
		if err := m.VerifyAndMaterialise(); err != nil {
			t.Fatalf("early-status metadata must not be an error, got %v", err)
		}
		if m.HashVerified {
			t.Error("HashVerified set for metadata with nothing to verify")
		}
	})
}

// ParsePayloadEntries has two arms and they must stay distinguishable: the
// materialised one carries the guarantee, the fallback does not.
func TestRequestMetadata_ParsePayloadEntriesArms(t *testing.T) {
	payload := `[{"key": "currency", "value": "BTC"}]`

	t.Run("materialised entries are returned", func(t *testing.T) {
		m := verifiedMetadata(t, payload)
		entries, err := m.ParsePayloadEntries()
		if err != nil || len(entries) != 1 || entries[0].Key != "currency" {
			t.Fatalf("entries = %v, err = %v", entries, err)
		}
	})

	t.Run("unverified metadata falls back to parsing", func(t *testing.T) {
		m := &RequestMetadata{PayloadAsString: payload}
		entries, err := m.ParsePayloadEntries()
		if err != nil || len(entries) != 1 {
			t.Fatalf("fallback should still parse: entries = %v, err = %v", entries, err)
		}
		// ...but it grants no access through the accessors.
		if _, err := m.GetMetadataCurrency(); !errors.Is(err, ErrMetadataUnverified) {
			t.Errorf("fallback parsing must not feed the accessors, err = %v", err)
		}
	})
}
