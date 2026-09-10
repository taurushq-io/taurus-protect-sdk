package mapper

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	pb "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/proto"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// This file bridges the base64 protobuf wire format to canonical protobuf JSON
// for callers outside the SDK (the internal/proto package is not importable).
// The helpers operate purely on the generated proto messages via protojson, so
// they preserve the exact canonical field names and byte fields — independent of
// the struct-based DecodedRulesContainer convenience model in rules_container.go.

// RulesContainerJSONFromBase64 decodes a base64-encoded protobuf RulesContainer into proto JSON.
//
// The returned JSON keeps the canonical protobuf field names and base64-encoded byte
// fields expected by protojson. It is suitable for round-tripping through
// RulesContainerBase64FromJSON without going through the lossy DecodedRulesContainer
// convenience model.
func RulesContainerJSONFromBase64(base64Data string) (json.RawMessage, error) {
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}
	var pbContainer pb.RulesContainer
	if err := proto.Unmarshal(data, &pbContainer); err != nil {
		return nil, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}
	jsonData, err := protojson.Marshal(&pbContainer)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal rules container JSON: %w", err)
	}
	return json.RawMessage(jsonData), nil
}

// RulesContainerBase64FromJSON encodes proto JSON into a base64 protobuf RulesContainer.
func RulesContainerBase64FromJSON(jsonData []byte) (string, error) {
	var pbContainer pb.RulesContainer
	if err := protojson.Unmarshal(jsonData, &pbContainer); err != nil {
		return "", fmt.Errorf("failed to unmarshal rules container JSON: %w", err)
	}
	data, err := deterministicMarshal.Marshal(&pbContainer)
	if err != nil {
		return "", fmt.Errorf("failed to marshal rules container protobuf: %w", err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// RuleMessageBase64FromJSON encodes a governance rule protobuf message into base64.
//
// messageType must be one of the concrete rule message names from request_reply.proto
// such as RuleSource, RuleDestination, RuleFiatAmount, RulesContainer_Line, or
// RulesContainer_TransactionRules. This helper exists so callers outside the SDK
// do not need to import the generated internal proto package.
func RuleMessageBase64FromJSON(messageType string, jsonData []byte) (string, error) {
	message, err := newGovernanceRuleMessage(messageType)
	if err != nil {
		return "", err
	}
	if err := protojson.Unmarshal(jsonData, message); err != nil {
		return "", fmt.Errorf("failed to unmarshal %s JSON: %w", strings.TrimSpace(messageType), err)
	}
	data, err := deterministicMarshal.Marshal(message)
	if err != nil {
		return "", fmt.Errorf("failed to marshal %s protobuf: %w", strings.TrimSpace(messageType), err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// RuleMessageJSONFromBase64 decodes a base64-encoded protobuf rule message
// back into canonical protobuf JSON. This is the reverse of RuleMessageBase64FromJSON.
// Returns nil, nil for empty input (empty cells mean "match any").
func RuleMessageJSONFromBase64(messageType string, base64Data string) (json.RawMessage, error) {
	if strings.TrimSpace(base64Data) == "" {
		return nil, nil
	}
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}
	if len(data) == 0 {
		return nil, nil
	}
	message, err := newGovernanceRuleMessage(messageType)
	if err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(data, message); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s protobuf: %w", strings.TrimSpace(messageType), err)
	}
	jsonData, err := protojson.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %s JSON: %w", strings.TrimSpace(messageType), err)
	}
	return jsonData, nil
}

func newGovernanceRuleMessage(messageType string) (proto.Message, error) {
	name := strings.TrimSpace(messageType)
	if name == "" {
		return nil, fmt.Errorf("messageType cannot be empty")
	}
	fullName := protoreflect.FullName(strings.ReplaceAll(name, "_", "."))
	if pkg := pb.File_request_reply_proto.Package(); pkg != "" {
		fullName = protoreflect.FullName(string(pkg) + "." + string(fullName))
	}
	message, err := protoregistry.GlobalTypes.FindMessageByName(fullName)
	if err != nil {
		return nil, fmt.Errorf("unsupported governance rule message type %q", messageType)
	}
	return message.New().Interface(), nil
}
