package mapper

import (
	"encoding/json"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

// The wire sends metadata.payload as a JSON ARRAY of entries. The proto declares it
// `google.protobuf.Value` (any JSON value), the swagger correctly emits an untyped
// `{}`, but openapi-generator's Go target maps that to map[string]interface{} — so
// decoding failed with "cannot unmarshal array into ... map[string]interface {}"
// and EVERY requests read broke. The field is now omitted entirely, which also
// keeps the tamperable object out of the process.
//
// Removing the field again would reintroduce both problems, so this pins the shape
// against the payload the live backend actually returns.
func TestMetadataDecodesArrayPayload(t *testing.T) {
	// Trimmed from a live tg-validatord response.
	body := []byte(`{
		"hash": "7d8428814511d4f8f0bddc421443214d003d19e2d5fa08f9be9dc9d2c1400d60",
		"payload": [
			{"column": "", "key": "request_id", "type": "String", "value": "227"},
			{"column": "", "key": "currency", "type": "String", "value": "CC"}
		],
		"payloadAsString": "[{\"key\":\"request_id\",\"value\":\"227\"}]"
	}`)

	var md openapi.TgvalidatordMetadata
	if err := json.Unmarshal(body, &md); err != nil {
		t.Fatalf("array payload must decode, got: %v", err)
	}
	if md.GetHash() == "" || md.GetPayloadAsString() == "" {
		t.Error("the verified fields must survive decoding")
	}
}

// The raw payload object must not be reachable: mapper/request.go maps only
// PayloadAsString, on the grounds that an interceptor could alter the object while
// leaving the hashed string intact. Omitting the field makes that structural.
func TestMetadataHasNoPayloadField(t *testing.T) {
	round, err := json.Marshal(openapi.TgvalidatordMetadata{})
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal(round, &out); err != nil {
		t.Fatal(err)
	}
	if _, ok := out["payload"]; ok {
		t.Error("serialized metadata still carries a `payload` key")
	}
}
