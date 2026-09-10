package mapper

import (
	"bytes"
	"testing"

	pb "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/proto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// buildRichPbContainer covers every modeled field of the rules container.
func buildRichPbContainer(t *testing.T) *pb.RulesContainer {
	t.Helper()

	fiatRange := mustMarshal(t, &pb.RuleFiatAmountRange{MinAmount: "1000", MaxAmount: "50000"})
	fiatCell := mustMarshal(t, &pb.RuleFiatAmount{Type: pb.RuleFiatAmount_RuleFiatAmountRange, Payload: fiatRange})

	destWallet := mustMarshal(t, &pb.RuleDestinationInternalWallet{Path: "m/44'/60'/1'"})
	destCell := mustMarshal(t, &pb.RuleDestination{Type: pb.RuleDestination_RuleDestinationInternalWallet, Payload: destWallet})

	stringCell := mustMarshal(t, &pb.RuleStringEqual{Type: pb.RuleStringEqual_RuleStringEqualValue, Payload: []byte("contract-42")})

	sourceWallet := mustMarshal(t, &pb.RuleSourceInternalWallet{Path: "m/44'/60'/0'"})
	sourceCell := mustMarshal(t, &pb.RuleSource{Type: pb.RuleSource_RuleSourceInternalWallet, Payload: sourceWallet})
	sourceAnyCell := mustMarshal(t, &pb.RuleSource{Type: pb.RuleSource_RuleSourceAny})

	thresholds := []*pb.SequentialThresholds{{
		Thresholds: []*pb.GroupThreshold{{GroupId: "approvers", MinimumSignatures: 2}},
	}}

	return &pb.RulesContainer{
		Users: []*pb.User{
			{
				Id:         "user-1",
				PublicKey:  "-----BEGIN PUBLIC KEY-----\nnot-a-real-key\n-----END PUBLIC KEY-----",
				Roles:      []pb.Role{pb.Role(1)},
				Properties: map[string][]byte{"team": []byte("ops")},
			},
			{Id: "user-2"},
		},
		Groups: []*pb.Group{
			{Id: "approvers", UserIds: []string{"user-1", "user-2"}, Properties: map[string][]byte{"k": []byte("v")}},
		},
		MinimumDistinctUserSignatures:  2,
		MinimumDistinctGroupSignatures: 1,
		TransactionRules: []*pb.RulesContainer_TransactionRules{
			{
				Key: "ETH/ERC20_transfer",
				Columns: []*pb.RulesContainer_Column{
					{Type: pb.RulesContainer_RuleFiatAmount, Name: "amount", MetadataKey: "amount"},
					{Type: pb.RulesContainer_RuleDestination, Name: "to", MetadataKey: "destination"},
					{Type: pb.RulesContainer_RuleStringEqual, Name: "contract", MetadataKey: "contract_id"},
				},
				Lines: []*pb.RulesContainer_Line{
					{
						Cells:              [][]byte{fiatCell, destCell, stringCell},
						ParallelThresholds: thresholds,
						Priority:           1,
						Properties:         map[string][]byte{"note": []byte("high value")},
					},
					{
						// Empty cells mean "match any".
						Cells:              [][]byte{nil, nil, nil},
						ParallelThresholds: thresholds,
					},
				},
				Details: &pb.RulesContainer_TransactionRules_TransactionRuleDetails{
					Domain:     pb.RulesContainer_TransactionRules_TransactionRuleDetails_RuleDomainTransfer,
					SubDomain:  pb.RulesContainer_TransactionRules_TransactionRuleDetails_RuleSubDomainERC20,
					Blockchain: "ETH",
					Network:    "mainnet",
					EvmCallContract: &pb.RulesContainer_TransactionRules_TransactionRuleDetails_EvmCallContract{
						ContractType:    "ERC20",
						MethodSignature: "transfer(address,uint256)",
					},
					CashSettlement: &pb.RulesContainer_TransactionRules_TransactionRuleDetails_CashSettlement{
						Provider:    "prov",
						RequestType: "settle",
					},
					CosmosDetails: &pb.RulesContainer_TransactionRules_TransactionRuleDetails_CosmosDetails{
						MethodSignatures: []string{"/cosmos.bank.v1beta1.MsgSend"},
					},
				},
			},
		},
		AddressWhitelistingRules: []*pb.RulesContainer_AddressWhitelistingRules{
			{
				Currency:           "ETH",
				Network:            "mainnet",
				ParallelThresholds: thresholds,
				Properties:         map[string][]byte{"p": []byte("q")},
				Lines: []*pb.RulesContainer_AddressWhitelistingRules_Line{
					{
						Cells:              [][]byte{sourceCell, sourceAnyCell},
						ParallelThresholds: thresholds,
						Properties:         map[string][]byte{"lp": []byte("lv")},
					},
				},
			},
		},
		ContractAddressWhitelistingRules: []*pb.RulesContainer_ContractAddressWhitelistingRules{
			{
				Blockchain:         pb.Blockchain_ETH,
				Network:            "mainnet",
				ParallelThresholds: thresholds,
				Properties:         map[string][]byte{"c": []byte("d")},
			},
		},
		EnforcedRulesHash:           "server-computed-hash",
		Properties:                  map[string][]byte{"tenant": []byte("42")},
		Timestamp:                   1750000000,
		MinimumCommitmentSignatures: 1,
		EngineIdentities:            []string{"hsm-1", "hsm-2"},
		HsmSlotId:                   7,
	}
}

// TestRulesContainer_RoundTripLossless is the core guarantee: decode → encode
// → decode yields the same model (modulo the server-controlled fields the
// encoder strips), and re-encoding is byte-stable.
func TestRulesContainer_RoundTripLossless(t *testing.T) {
	data, err := deterministicMarshal.Marshal(buildRichPbContainer(t))
	mustNoErr(t, err)

	first, err := RulesContainerFromBytes(data)
	mustNoErr(t, err)
	if first.HasUnknownFields() {
		t.Fatal("fully-known container should not report unknown fields")
	}

	encoded, err := RulesContainerToBytes(first)
	mustNoErr(t, err)

	second, err := RulesContainerFromBytes(encoded)
	mustNoErr(t, err)

	// The encoder strips server-controlled fields.
	if second.EnforcedRulesHash != "" || second.Timestamp != 0 {
		t.Fatalf("server-controlled fields not stripped: hash=%q ts=%d", second.EnforcedRulesHash, second.Timestamp)
	}
	first.EnforcedRulesHash = ""
	first.Timestamp = 0

	wantEqual(t, first, second)

	// Typed cells decoded as expected (spot checks).
	line := second.TransactionRules[0].Lines[0]
	wantEqual(t, model.RuleCell(model.FiatAmountRange{MinAmount: "1000", MaxAmount: "50000"}), line.Cells[0])
	wantEqual(t, model.RuleCell(model.DestinationInternalWallet{Path: "m/44'/60'/1'"}), line.Cells[1])
	wantEqual(t, model.RuleCell(model.StringEqualValue{Value: "contract-42"}), line.Cells[2])
	// Empty cells decode to the columns' typed *Any values.
	wantEqual(t,
		[]model.RuleCell{model.FiatAmountAny{}, model.DestinationAny{}, model.StringEqualAny{}},
		second.TransactionRules[0].Lines[1].Cells)
	wantEqual(t, uint32(1), line.Priority)
	wantEqual(t, "amount", second.TransactionRules[0].Columns[0].Name)
	wantEqual(t, "destination", second.TransactionRules[0].Columns[1].MetadataKey)
	wantEqual(t, "ETH", second.TransactionRules[0].Details.Blockchain)
	wantEqual(t, "transfer(address,uint256)", second.TransactionRules[0].Details.EvmCallContract.MethodSignature)

	// Whitelisting sources decoded with payload variants.
	source := second.AddressWhitelistingRules[0].Lines[0].Cells[0]
	wantEqual(t, model.RuleSourceTypeInternalWallet, source.Type)
	wantEqual(t, "m/44'/60'/0'", source.InternalWallet.Path)

	// Encoding the same model twice is byte-identical (deterministic marshal).
	reEncoded, err := RulesContainerToBytes(second)
	mustNoErr(t, err)
	if !bytes.Equal(encoded, reEncoded) {
		t.Fatal("re-encoding the same container is not byte-stable")
	}
}

// TestRulesContainer_EncodeStripsServerControlledFields pins the proposal
// submission requirement: enforcedRulesHash and timestamp must not survive.
func TestRulesContainer_EncodeStripsServerControlledFields(t *testing.T) {
	decoded, err := RulesContainerFromBytes(mustMarshal(t, buildRichPbContainer(t)))
	mustNoErr(t, err)
	wantEqual(t, "server-computed-hash", decoded.EnforcedRulesHash)
	if decoded.Timestamp == 0 {
		t.Fatal("fixture should carry a timestamp")
	}

	encoded, err := RulesContainerToBytes(decoded)
	mustNoErr(t, err)

	var back pb.RulesContainer
	mustNoErr(t, proto.Unmarshal(encoded, &back))
	if back.GetEnforcedRulesHash() != "" || back.GetTimestamp() != 0 {
		t.Fatalf("server-controlled fields survived encode: hash=%q ts=%d", back.GetEnforcedRulesHash(), back.GetTimestamp())
	}
}

// TestRulesContainer_UnknownFieldsSurviveRoundTrip is the schema-skew
// guarantee: fields from a newer request_reply.proto are never silently
// dropped — they are preserved per node, surviving caller edits and
// reordering.
func TestRulesContainer_UnknownFieldsSurviveRoundTrip(t *testing.T) {
	topLevelUnknown := protoreflect.RawFields{0xC0, 0x0C, 0x2A}   // field 200, varint 42
	nestedRuleUnknown := protoreflect.RawFields{0xE0, 0x12, 0x07} // field 300, varint 7
	nestedUserUnknown := protoreflect.RawFields{0xA0, 0x19, 0x01} // field 404, varint 1

	pbContainer := buildRichPbContainer(t)
	pbContainer.ProtoReflect().SetUnknown(topLevelUnknown)
	pbContainer.TransactionRules[0].ProtoReflect().SetUnknown(nestedRuleUnknown)
	pbContainer.Users[0].ProtoReflect().SetUnknown(nestedUserUnknown)

	decoded, err := RulesContainerFromBytes(mustMarshal(t, pbContainer))
	mustNoErr(t, err)
	if !decoded.HasUnknownFields() {
		t.Fatal("expected HasUnknownFields to report true")
	}

	// Edit the container the way a caller would: prepend a brand-new rule so
	// the original rule moves — unknown fields must travel with their node.
	decoded.TransactionRules = append([]*model.TransactionRules{{Key: "NEW/rule"}}, decoded.TransactionRules...)

	encoded, err := RulesContainerToBytes(decoded)
	mustNoErr(t, err)

	var back pb.RulesContainer
	mustNoErr(t, proto.Unmarshal(encoded, &back))
	wantEqual(t, []byte(topLevelUnknown), []byte(back.ProtoReflect().GetUnknown()))

	if len(back.GetTransactionRules()) != 2 {
		t.Fatalf("expected 2 transaction rules, got %d", len(back.GetTransactionRules()))
	}
	wantEqual(t, "NEW/rule", back.GetTransactionRules()[0].GetKey())
	if len(back.GetTransactionRules()[0].ProtoReflect().GetUnknown()) != 0 {
		t.Fatal("new rule must not inherit unknown fields")
	}
	wantEqual(t, "ETH/ERC20_transfer", back.GetTransactionRules()[1].GetKey())
	wantEqual(t, []byte(nestedRuleUnknown), []byte(back.GetTransactionRules()[1].ProtoReflect().GetUnknown()))

	wantEqual(t, []byte(nestedUserUnknown), []byte(back.GetUsers()[0].ProtoReflect().GetUnknown()))
}

// TestRulesContainer_UnknownEnumValuePassesThroughNumerically covers enum
// values newer than this SDK: they decode to their numeric string and
// re-encode to the same number.
func TestRulesContainer_UnknownEnumValuePassesThroughNumerically(t *testing.T) {
	pbContainer := buildRichPbContainer(t)
	pbContainer.TransactionRules[0].Columns = append(pbContainer.TransactionRules[0].Columns,
		&pb.RulesContainer_Column{Type: pb.RulesContainer_ColumnType(999), Name: "future"})

	decoded, err := RulesContainerFromBytes(mustMarshal(t, pbContainer))
	mustNoErr(t, err)

	cols := decoded.TransactionRules[0].Columns
	wantEqual(t, "999", cols[len(cols)-1].Type)

	encoded, err := RulesContainerToBytes(decoded)
	mustNoErr(t, err)

	var back pb.RulesContainer
	mustNoErr(t, proto.Unmarshal(encoded, &back))
	backCols := back.GetTransactionRules()[0].GetColumns()
	wantEqual(t, pb.RulesContainer_ColumnType(999), backCols[len(backCols)-1].GetType())
}

// TestRulesContainer_UnknownSourceTypeSurvivesVerbatim covers whitelisting
// source cells of a type newer than this SDK.
func TestRulesContainer_UnknownSourceTypeSurvivesVerbatim(t *testing.T) {
	futureSource := mustMarshal(t, &pb.RuleSource{
		Type:    pb.RuleSource_RuleSourceType(99),
		Payload: []byte("opaque-future-payload"),
	})

	pbContainer := buildRichPbContainer(t)
	pbContainer.AddressWhitelistingRules[0].Lines[0].Cells = [][]byte{futureSource}

	decoded, err := RulesContainerFromBytes(mustMarshal(t, pbContainer))
	mustNoErr(t, err)
	if !decoded.HasUnknownFields() {
		t.Fatal("raw source should surface via HasUnknownFields")
	}

	source := decoded.AddressWhitelistingRules[0].Lines[0].Cells[0]
	wantEqual(t, futureSource, source.Raw)

	encoded, err := RulesContainerToBytes(decoded)
	mustNoErr(t, err)
	var back pb.RulesContainer
	mustNoErr(t, proto.Unmarshal(encoded, &back))
	wantEqual(t, futureSource, back.GetAddressWhitelistingRules()[0].GetLines()[0].GetCells()[0])
}

// TestRulesContainer_AuthoredUnknownEnumErrors pins the fail-loud rule for
// caller-authored values (as opposed to wire-decoded skew, which passes
// through numerically).
func TestRulesContainer_AuthoredUnknownEnumErrors(t *testing.T) {
	container := &model.DecodedRulesContainer{
		TransactionRules: []*model.TransactionRules{{
			Key:     "k",
			Columns: []*model.RuleColumn{{Type: "RuleTypoAmount"}},
		}},
	}
	_, err := RulesContainerToBytes(container)
	wantErrContaining(t, err, `unknown column type "RuleTypoAmount"`)
}

func TestRulesContainer_EncodeNilContainerErrors(t *testing.T) {
	if _, err := RulesContainerToBytes(nil); err == nil {
		t.Fatal("expected error for nil container")
	}
}

func mustMarshal(t *testing.T, m proto.Message) []byte {
	t.Helper()
	data, err := proto.Marshal(m)
	mustNoErr(t, err)
	return data
}
