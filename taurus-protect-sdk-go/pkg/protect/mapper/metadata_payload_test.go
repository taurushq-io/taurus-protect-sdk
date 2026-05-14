package mapper

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestTgvalidatordMetadataPayload_AcceptsAnyValidJSONRoot(t *testing.T) {
	tests := []struct {
		name    string
		payload string
	}{
		{name: "object", payload: `{"k":"v","n":1}`},
		{name: "array", payload: `[1,"two",{"k":"v"}]`},
		{name: "string", payload: `"plain-string"`},
		{name: "number", payload: `123.45`},
		{name: "boolean", payload: `true`},
		{name: "null", payload: `null`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := openapi.NewTgvalidatordMetadata()
			meta.SetPayload(json.RawMessage(tt.payload))

			if !json.Valid(meta.GetPayload()) {
				t.Fatalf("payload is not valid JSON: %q", meta.GetPayload())
			}

			encoded, err := json.Marshal(meta)
			if err != nil {
				t.Fatalf("marshal metadata: %v", err)
			}

			var decoded openapi.TgvalidatordMetadata
			if err := json.Unmarshal(encoded, &decoded); err != nil {
				t.Fatalf("unmarshal metadata: %v", err)
			}

			if !bytes.Equal(decoded.GetPayload(), []byte(tt.payload)) {
				t.Fatalf("decoded payload mismatch: got=%s want=%s", decoded.GetPayload(), tt.payload)
			}
		})
	}
}
