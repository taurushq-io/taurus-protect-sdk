package openapi

import (
	"encoding/json"
	"testing"
)

func TestTgvalidatordMetadata_UnmarshalJSON_ArrayPayload(t *testing.T) {
	// JSON with payload as array (actual API response format)
	jsonData := `{
		"hash": "abc123",
		"payload": [
			{"key": "source", "value": {"payload": {"address": "addr1"}}},
			{"key": "destination", "value": {"payload": {"address": "addr2"}}}
		],
		"payloadAsString": "[{\"key\":\"source\"}]"
	}`

	var metadata TgvalidatordMetadata
	err := json.Unmarshal([]byte(jsonData), &metadata)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	// Hash should be preserved
	if metadata.GetHash() != "abc123" {
		t.Errorf("Hash = %q, want 'abc123'", metadata.GetHash())
	}

	// PayloadAsString should be preserved
	if metadata.GetPayloadAsString() != "[{\"key\":\"source\"}]" {
		t.Errorf("PayloadAsString = %q, want '[{\"key\":\"source\"}]'", metadata.GetPayloadAsString())
	}

	// Payload should be nil since it was an array
	if metadata.Payload != nil {
		t.Errorf("Payload should be nil for array format, got %v", metadata.Payload)
	}
}

func TestTgvalidatordMetadata_UnmarshalJSON_MapPayload(t *testing.T) {
	// JSON with payload as map - Payload should still be nil for security reasons
	// (all data extraction must use PayloadAsString, the cryptographically verified source)
	jsonData := `{
		"hash": "def456",
		"payload": {"source": "addr1", "destination": "addr2"},
		"payloadAsString": "{\"source\":\"addr1\"}"
	}`

	var metadata TgvalidatordMetadata
	err := json.Unmarshal([]byte(jsonData), &metadata)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	// Hash should be preserved
	if metadata.GetHash() != "def456" {
		t.Errorf("Hash = %q, want 'def456'", metadata.GetHash())
	}

	// PayloadAsString should be preserved
	if metadata.GetPayloadAsString() != "{\"source\":\"addr1\"}" {
		t.Errorf("PayloadAsString = %q, want '{\"source\":\"addr1\"}'", metadata.GetPayloadAsString())
	}

	// SECURITY: Payload MUST be nil even for map format.
	// Extracting data from the raw payload field bypasses integrity verification.
	// All data extraction must use PayloadAsString (the cryptographically verified source).
	if metadata.Payload != nil {
		t.Errorf("Payload should be nil for security (was %v); extract from PayloadAsString instead", metadata.Payload)
	}
}

func TestTgvalidatordMetadata_UnmarshalJSON_NullPayload(t *testing.T) {
	// JSON with null payload
	jsonData := `{
		"hash": "ghi789",
		"payload": null,
		"payloadAsString": "test"
	}`

	var metadata TgvalidatordMetadata
	err := json.Unmarshal([]byte(jsonData), &metadata)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	if metadata.GetHash() != "ghi789" {
		t.Errorf("Hash = %q, want 'ghi789'", metadata.GetHash())
	}

	if metadata.Payload != nil {
		t.Errorf("Payload should be nil, got %v", metadata.Payload)
	}
}

func TestTgvalidatordMetadata_UnmarshalJSON_MissingFields(t *testing.T) {
	// JSON with missing fields
	jsonData := `{
		"hash": "jkl012"
	}`

	var metadata TgvalidatordMetadata
	err := json.Unmarshal([]byte(jsonData), &metadata)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	if metadata.GetHash() != "jkl012" {
		t.Errorf("Hash = %q, want 'jkl012'", metadata.GetHash())
	}

	if metadata.Payload != nil {
		t.Errorf("Payload should be nil, got %v", metadata.Payload)
	}

	if metadata.GetPayloadAsString() != "" {
		t.Errorf("PayloadAsString should be empty, got %q", metadata.GetPayloadAsString())
	}
}

func TestTgvalidatordMetadata_UnmarshalJSON_EmptyObject(t *testing.T) {
	jsonData := `{}`

	var metadata TgvalidatordMetadata
	err := json.Unmarshal([]byte(jsonData), &metadata)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	if metadata.Hash != nil {
		t.Errorf("Hash should be nil, got %v", metadata.Hash)
	}
}

func TestTgvalidatordMetadata_UnmarshalJSON_InvalidJSON(t *testing.T) {
	jsonData := `not valid json`

	var metadata TgvalidatordMetadata
	err := json.Unmarshal([]byte(jsonData), &metadata)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
