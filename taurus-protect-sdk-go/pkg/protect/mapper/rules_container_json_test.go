package mapper

import (
	"encoding/base64"
	"fmt"
	"testing"

	pb "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/proto"
	"google.golang.org/protobuf/proto"
)

func TestRulesContainerJSONRoundTripPreservesProto(t *testing.T) {
	container := createTestRulesContainer()
	encoded := encodeRulesContainerToBase64(t, container)

	jsonData, err := RulesContainerJSONFromBase64(encoded)
	if err != nil {
		t.Fatalf("RulesContainerJSONFromBase64() error = %v", err)
	}

	encodedAgain, err := RulesContainerBase64FromJSON(jsonData)
	if err != nil {
		t.Fatalf("RulesContainerBase64FromJSON() error = %v", err)
	}

	got := decodeRulesContainerProto(t, encodedAgain)
	if !proto.Equal(container, got) {
		t.Fatalf("round-tripped container differs\nwant: %v\ngot: %v", container, got)
	}
}

// The JSON bridge is the path the MCP governance tools drive: decode a container to JSON,
// let a caller edit it, re-encode, submit — and approvers sign the re-encoded bytes. An
// edited container arrives with its map keys in arbitrary order, so encoding MUST be
// order-independent: every runtime emits map<string, bytes> in unspecified order unless
// told otherwise. The input below carries the five `properties` keys in reverse-sorted
// order on purpose; feeding back the canonical (already sorted) JSON would pass even with
// a non-deterministic encoder and prove nothing. Expected bytes are the same vector
// asserted by TestRulesContainerToBase64_IsDeterministicAndCrossSDKStable for the typed
// encoder, and the same reversed input is asserted in all four SDKs.
func TestRulesContainerJSONBridge_IsOrderIndependentAndCrossSDKStable(t *testing.T) {
	const expected = "CmgKAnUxEgNQRU0aAQMiEAoGa0FscGhhEgZrQWxwaGEiEAoGa0JyYXZvEgZrQnJhdm8iFAoIa0NoYXJsaWUSCGtDaGFybGllIhAKBmtEZWx0YRIGa0RlbHRhIg4KBWtFY2hvEgVrRWNob0oQCgZrQWxwaGESBmtBbHBoYUoQCgZrQnJhdm8SBmtCcmF2b0oUCghrQ2hhcmxpZRIIa0NoYXJsaWVKEAoGa0RlbHRhEgZrRGVsdGFKDgoFa0VjaG8SBWtFY2hv"
	const reversedKeyJSON = `{"users":[{"id":"u1","publicKey":"PEM","roles":["SUPERADMIN"],"properties":{"kEcho":"a0VjaG8=","kDelta":"a0RlbHRh","kCharlie":"a0NoYXJsaWU=","kBravo":"a0JyYXZv","kAlpha":"a0FscGhh"}}],"properties":{"kEcho":"a0VjaG8=","kDelta":"a0RlbHRh","kCharlie":"a0NoYXJsaWU=","kBravo":"a0JyYXZv","kAlpha":"a0FscGhh"}}`

	for i := 0; i < 20; i++ {
		out, err := RulesContainerBase64FromJSON([]byte(reversedKeyJSON))
		if err != nil {
			t.Fatalf("run %d: base64 from json: %v", i, err)
		}
		if out != expected {
			t.Fatalf("run %d not stable/cross-SDK-equal:\n got %s\nwant %s", i, out, expected)
		}
	}
}

func TestRuleMessageBase64FromJSONEncodesRuleCells(t *testing.T) {
	walletPayload, err := RuleMessageBase64FromJSON("RuleSourceInternalWallet", []byte(`{"path":"m/44'/60'/0'/0/0"}`))
	if err != nil {
		t.Fatalf("RuleMessageBase64FromJSON(wallet) error = %v", err)
	}
	sourceJSON := []byte(fmt.Sprintf(`{"type":"RuleSourceInternalWallet","payload":%q}`, walletPayload))
	sourceCell, err := RuleMessageBase64FromJSON("RuleSource", sourceJSON)
	if err != nil {
		t.Fatalf("RuleMessageBase64FromJSON(source) error = %v", err)
	}
	sourceBytes, err := base64.StdEncoding.DecodeString(sourceCell)
	if err != nil {
		t.Fatalf("decode source cell: %v", err)
	}
	var source pb.RuleSource
	if err := proto.Unmarshal(sourceBytes, &source); err != nil {
		t.Fatalf("unmarshal source cell: %v", err)
	}
	if source.GetType() != pb.RuleSource_RuleSourceInternalWallet {
		t.Fatalf("source type = %v", source.GetType())
	}
	_, err = RuleMessageBase64FromJSON("NoSuchRuleMessage", []byte(`{}`))
	if err == nil {
		t.Fatal("RuleMessageBase64FromJSON() expected error")
	}
}

func TestRulesContainerJSONFromBase64Errors(t *testing.T) {
	if _, err := RulesContainerJSONFromBase64("!!!not base64!!!"); err == nil {
		t.Fatal("expected error for invalid base64")
	}
	// Valid base64 but invalid protobuf wire bytes (truncated varint tag).
	if _, err := RulesContainerJSONFromBase64(base64.StdEncoding.EncodeToString([]byte{0x08})); err == nil {
		t.Fatal("expected error for invalid protobuf")
	}
}

func TestRulesContainerBase64FromJSONRejectsInvalidJSON(t *testing.T) {
	if _, err := RulesContainerBase64FromJSON([]byte("not json")); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestRuleMessageJSONFromBase64EmptyReturnsNil(t *testing.T) {
	got, err := RuleMessageJSONFromBase64("RuleSource", "")
	if err != nil {
		t.Fatalf("empty input error = %v", err)
	}
	if got != nil {
		t.Fatalf("empty input = %v, want nil", got)
	}
	got, err = RuleMessageJSONFromBase64("RuleSource", "   ")
	if err != nil || got != nil {
		t.Fatalf("whitespace input = (%v, %v), want (nil, nil)", got, err)
	}
}

func TestRuleMessageJSONFromBase64Errors(t *testing.T) {
	if _, err := RuleMessageJSONFromBase64("RuleSource", "!!!not base64!!!"); err == nil {
		t.Fatal("expected error for invalid base64")
	}
	if _, err := RuleMessageJSONFromBase64("NoSuchRuleMessage", "AAAA"); err == nil {
		t.Fatal("expected error for unknown message type")
	}
}

// TestRuleMessageJSONRoundTripStable verifies JSON->base64->JSON->base64 is stable
// for a concrete rule cell type.
func TestRuleMessageJSONRoundTripStable(t *testing.T) {
	encoded, err := RuleMessageBase64FromJSON("RuleFiatAmountRange", []byte(`{"minAmount":"1000"}`))
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	decoded, err := RuleMessageJSONFromBase64("RuleFiatAmountRange", encoded)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if decoded == nil {
		t.Fatal("decoded JSON is nil")
	}
	encodedAgain, err := RuleMessageBase64FromJSON("RuleFiatAmountRange", decoded)
	if err != nil {
		t.Fatalf("re-encode: %v", err)
	}
	if encodedAgain != encoded {
		t.Fatalf("round-trip not stable: %q != %q", encodedAgain, encoded)
	}
}

func TestNewGovernanceRuleMessageRejectsEmpty(t *testing.T) {
	if _, err := newGovernanceRuleMessage("   "); err == nil {
		t.Fatal("expected error for empty message type")
	}
}

func decodeRulesContainerProto(t *testing.T, encoded string) *pb.RulesContainer {
	t.Helper()
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("decode base64: %v", err)
	}
	var container pb.RulesContainer
	if err := proto.Unmarshal(data, &container); err != nil {
		t.Fatalf("unmarshal rules container: %v", err)
	}
	return &container
}
