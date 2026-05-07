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
